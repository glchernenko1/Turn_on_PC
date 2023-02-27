package middleware

import (
	"net/http"
	"strings"
	"Turn_on_PC/internal/DTO"
	"os"
	"github.com/dgrijalva/jwt-go"
	"Turn_on_PC/internal/server/apperror"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error
type appHandlerWithToken func(w http.ResponseWriter, r *http.Request, token *DTO.JWTUser) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			switch err.(type) {
			case *apperror.AppError:
				switch err {
				case apperror.ErrNotFound:
					{
						w.WriteHeader(http.StatusNotFound)
						w.Write(apperror.ErrNotFound.Marshal())
					}
				case apperror.Unauthorized:
					{
						w.WriteHeader(http.StatusUnauthorized)
						w.Write(apperror.Unauthorized.Marshal())
					}
				case apperror.NameTaken:
					{
						w.WriteHeader(http.StatusBadRequest)
						w.Write(err.(*apperror.AppError).Marshal())
					}
				default:
					{
						w.WriteHeader(http.StatusBadRequest)
						w.Write(err.(*apperror.AppError).Marshal())
					}
				}
			default:
				{
					w.WriteHeader(http.StatusTeapot)
					w.Write(apperror.SystemError(err).Marshal())
				}
			}
		}
	}
}

func MiddlewareAuth(h appHandlerWithToken, scopes ...string) http.HandlerFunc {
	return Middleware(func(w http.ResponseWriter, r *http.Request) error {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			return apperror.Unauthorized
		}
		splitted := strings.Split(tokenHeader, " ") //Токен обычно поставляется в формате `Bearer {token-body}`, мы проверяем, соответствует ли полученный токен этому требованию
		if len(splitted) != 2 {
			return apperror.Unauthorized
		}
		tokenPart := splitted[1]
		tk := &DTO.JWTUser{}
		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil {
			return apperror.Unauthorized
		}
		if !token.Valid {
			return apperror.Unauthorized
		}
		for _, scope := range scopes {
			if tk.Scope != scope {
				return apperror.Unauthorized
			}
		}
		return h(w, r, tk)
	})
}
