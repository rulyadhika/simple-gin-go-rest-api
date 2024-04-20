package jwt

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type JWTPayload struct {
	Id       uint32            `json:"id"`
	Username string            `json:"username"`
	Email    string            `json:"email"`
	Roles    []entity.UserType `json:"roles"`
}

func NewJWTTokenParser() *JWTPayload {
	return &JWTPayload{}
}

func (j *JWTPayload) ParseToken(tokenString string, secretKey string) error {
	token, errParse := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(secretKey), nil
	})

	if errParse != nil {
		log.Printf("[ParseToken - JWT] err: %s\n", errParse.Error())
		return errParse
	}

	if err := j.bindTokenToStruct(token); err != nil {
		log.Printf("[ParseToken - JWT] err: %s\n", err.Error())
		return err
	}

	return nil
}

func (j *JWTPayload) bindTokenToStruct(token *jwt.Token) error {
	var claims jwt.MapClaims

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		claims = mapClaims
	} else {
		return errors.New("invalid token")
	}

	roles := []entity.UserType{}
	for _, data := range claims["roles"].([]any) {
		roles = append(roles, entity.UserType(data.(string)))
	}

	j.Id = uint32(claims["id"].(float64))
	j.Email = claims["email"].(string)
	j.Username = claims["username"].(string)
	j.Roles = roles

	return nil
}
