package jwt

import "github.com/golang-jwt/jwt/v5"

type jwtToken struct {
	Id       uint32           `json:"id"`
	Username string           `json:"username"`
	Email    string           `json:"email"`
	Roles    []string         `json:"roles"`
	Exp      *jwt.NumericDate `json:"exp"`
	Iat      *jwt.NumericDate `json:"iat"`
}
