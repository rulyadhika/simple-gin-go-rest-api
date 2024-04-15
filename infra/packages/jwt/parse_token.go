package jwt

import (
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

func NewJWTTokenParser() *jwtToken {
	return &jwtToken{}
}

func (j *jwtToken) ParseToken(tokenString string, secretKey string) error {
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

func (j *jwtToken) bindTokenToStruct(token *jwt.Token) error {
	var claims jwt.MapClaims

	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
		claims = mapClaims
	} else {
		return errors.New("invalid token")
	}

	roles := []string{}
	for _, data := range claims["roles"].([]any) {
		roles = append(roles, data.(string))
	}

	j.Id = uint32(claims["id"].(float64))
	j.Email = claims["email"].(string)
	j.Username = claims["username"].(string)
	j.Roles = roles
	j.Iat, _ = claims.GetIssuedAt()
	j.Exp, _ = claims.GetExpirationTime()

	return nil
}
