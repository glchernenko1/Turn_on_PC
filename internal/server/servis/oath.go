package servis

import (
	"Turn_on_PC/internal/DTO"
	"Turn_on_PC/pkg/hash"
	"os"
	"github.com/dgrijalva/jwt-go"
	"Turn_on_PC/internal/server/apperror"
	"Turn_on_PC/internal/server/DB"
)

func Register(db DB.DB, user *DTO.CreateUser) (uint, error) {
	//todo добавить проверку коректности полей
	unit, _ := DTO.Create(user)
	return db.AddUser(unit)
}

func SingIn(db DB.DB, login string, password string, scope string) (string, error) {

	user, err := db.FiendUserByLogin(login)
	if err != nil {
		return "", apperror.Unauthorized
	}
	if hash.CheckPasswordHash(password, user.PasswordHash) {
		tk := &DTO.JWTUser{UserId: user.ID, Scope: scope}
		token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
		tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
		return tokenString, err
	}
	return "", apperror.Unauthorized
}
