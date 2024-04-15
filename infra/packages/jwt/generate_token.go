package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

func NewJWTToken(userData *userrepository.UserRoles) *jwtToken {
	timeIssuedAt := jwt.NewNumericDate(time.Now())

	userRoles := []string{}

	for _, role := range userData.Roles {
		userRoles = append(userRoles, role.RoleName)
	}

	return &jwtToken{
		Id:       userData.Id,
		Username: userData.Username,
		Email:    userData.Email,
		Roles:    userRoles,
		Iat:      timeIssuedAt,
	}
}

func (j *jwtToken) generateClaims(timeExpiredAt *jwt.NumericDate) jwt.Claims {
	return jwt.MapClaims{
		"id":       j.Id,
		"username": j.Username,
		"email":    j.Email,
		"roles":    j.Roles,
		"exp":      timeExpiredAt,
		"iat":      j.Iat,
	}
}

func (j *jwtToken) signToken(claims jwt.Claims, secretKey string) (any, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	stringToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Printf("[SignToken - JWT] err:%s\n", err.Error())
		return nil, err
	}

	return stringToken, nil
}

func (j *jwtToken) GenerateToken(jwtTokenSecret string, timeExpiredAt time.Time) (any, error) {
	exp := jwt.NewNumericDate(timeExpiredAt)

	payload := j.generateClaims(exp)
	return j.signToken(payload, jwtTokenSecret)
}
