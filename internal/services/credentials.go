package services

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	errs "github.com/mathif92/auth-service/internal/errors"
	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	_insertCredentials = `INSERT INTO credentials (username, password, email) 
								VALUES (?, ?, ?)`

	_selectToken = `SELECT id FROM token WHERE credentials_id = ?`

	_selectCredentials = `SELECT c.* FROM credentials as c WHERE c.id = ?`
)

type Credentials struct {
	db           *sqlx.DB
	tokenService *Token
}

func NewCredentials(db *sqlx.DB, tokenService *Token) *Credentials {
	return &Credentials{db: db, tokenService: tokenService}
}

// SaveCredentials saves the credentials provided in the credentials parameter into the database, returns an error in case there's one
func (c *Credentials) SaveCredentials(ctx context.Context, credentials CredentialsModel) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "creating db transacion")
	}
	credentials.Password = hashAndSalt([]byte(credentials.Password))

	_, err = tx.ExecContext(ctx, _insertCredentials, credentials.Username, credentials.Password, credentials.Email)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "inserting credentials into db")
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "commiting db transaction")
	}

	return nil
}

// Authenticate returns a valid access_token for the provided credentials or an error instead
func (c *Credentials) Authenticate(ctx context.Context, credentials CredentialsModel) (string, error) {
	query := `SELECT * FROM credentials WHERE email = ?`
	args := []interface{}{credentials.Email}
	if credentials.Username != "" {
		query = `SELECT * FROM credentials WHERE username = ?`
		args[0] = credentials.Username
	}

	var storedCredentials CredentialsModel
	if err := c.db.GetContext(ctx, &storedCredentials, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.New("invalid credentials", http.StatusUnauthorized)
		}
		return "", errors.Wrap(err, "retrieving credentials from db")
	}

	if !comparePasswords(storedCredentials.Password, []byte(credentials.Password)) {
		return "", errs.New("invalid credentials", http.StatusUnauthorized)
	}

	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", errors.Wrap(err, "creating db transaction")
	}

	token, err := c.tokenService.GenerateToken(credentials.Username, credentials.Email)
	if err != nil {
		return "", errors.Wrap(err, "generating token")
	}

	var tokenID int64
	if err := c.db.GetContext(ctx, &tokenID, _selectToken, storedCredentials.ID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", errors.Wrap(err, "retrieving token from db")
		}
		// Token does not exist, we need to create it
		const (
			_insertToken = `INSERT INTO token (credentials_id, token, updated_at) 
									VALUES (?, ?, current_timestamp())`
		)

		if _, err := tx.ExecContext(ctx, _insertToken, storedCredentials.ID, token); err != nil {
			tx.Rollback()
			return "", errors.Wrap(err, "creating new token")
		}
	} else {
		// There's a token row in the db, we need to update the token value
		const (
			_updateToken = `UPDATE token SET token = ?, updated_at = current_timestamp() WHERE id = ?`
		)
		if _, err := tx.ExecContext(ctx, _updateToken, token, tokenID); err != nil {
			tx.Rollback()
			return "", errors.Wrap(err, "updating token")
		}
	}

	if err := tx.Commit(); err != nil {
		return "", errors.Wrap(err, "commiting db transaction")
	}

	return token, nil
}

func (c *Credentials) GetCredentials(ctx context.Context, credentialsID int64) (CredentialsModel, error) {
	var credentials CredentialsModel
	if err := c.db.GetContext(ctx, &credentials, _selectCredentials, credentialsID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CredentialsModel{}, errs.New("credentials not found", http.StatusNotFound)
		}
		return CredentialsModel{}, errors.Wrap(err, "retrieving credentials from DB")
	}

	return credentials, nil
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
