package services

import (
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

func (t *Token) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "username"

	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *Token) VerifyToken(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		if request.Header["Authorization"] != nil {
			providedToken := request.Header["Authorization"][0]
			if strings.Contains(providedToken, "Bearer") {
				providedToken = strings.Split(providedToken, " ")[1]
			}

			jwtToken, err := jwt.Parse(providedToken, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodECDSA)
				if !ok {
					if err := api.RespondError(ctx, writer, errors.New("Unauthorized", http.StatusUnauthorized)); err != nil {
						return nil, err

					}
				}
				return "", nil

			})

			if err != nil {
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
				endpointHandler(writer, request)
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
