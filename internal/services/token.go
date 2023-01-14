package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mathif92/auth-service/internal/api"
	"github.com/mathif92/auth-service/internal/errors"
)

type Token struct {
	secretKey string
}

func NewToken(secretKey string) *Token {
	return &Token{secretKey: secretKey}
}

func (t *Token) GenerateToken(username string, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)

	// The "iat" claims field needs to be slightly less than now in order to make the verification work
	claims["iat"] = time.Now().Add(-(time.Second * 2)).Unix()

	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["username"] = username
	claims["email"] = email

	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Token) VerifyToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if request.Header["Authorization"] != nil {
			providedToken := request.Header["Authorization"][0]
			if strings.Contains(providedToken, "Bearer") {
				providedToken = strings.Split(providedToken, " ")[1]
			}

			jwtToken, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(t.secretKey), nil
			})

			if err != nil {
				fmt.Printf("Error parsing JWT: %e\n", err)
				if err := api.RespondError(ctx, writer, errors.New("Unauthorized: error parsing the JWT", http.StatusUnauthorized)); err != nil {
					return
				}
			}

			// if there's a token
			if jwtToken.Valid {
				if err := t.VerifyPermissions(jwtToken); err != nil {
					if err := api.RespondError(ctx, writer, errors.New("Lacking permissions", http.StatusForbidden)); err != nil {
						return
					}
				}
				next.ServeHTTP(writer, request)
			} else {
				if err := api.RespondError(ctx, writer, errors.New("Invalid token", http.StatusUnauthorized)); err != nil {
					return
				}
			}
		} else {
			if err := api.RespondError(ctx, writer, errors.New("Missing token", http.StatusUnauthorized)); err != nil {
				return
			}
		}
	})
}

func (t *Token) VerifyPermissions(token *jwt.Token) error {
	// Not implemented yet. In the future, this needs to retrieve the username/email & verify the permissions
	return nil
}
