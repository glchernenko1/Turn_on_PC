package DTO

import "github.com/dgrijalva/jwt-go"

type JWTUser struct {
	UserId uint `json:"userId"`
	jwt.StandardClaims
	Scope string `json:"scope"`
}
