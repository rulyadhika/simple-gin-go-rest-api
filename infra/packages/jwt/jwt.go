package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

func generatePayload(data userrepository.UserRoles) jwt.Claims {
	timeIssuedAt := jwt.NewNumericDate(time.Now())
	timeExpiredAt := jwt.NewNumericDate(time.Now().Add(1 * time.Minute))

	userRoles := []string{}

	for _, role := range data.Roles {
		userRoles = append(userRoles, role.RoleName)
	}

	userData := map[string]any{
		"id":       data.Id,
		"username": data.Username,
		"email":    data.Email,
		"roles":    userRoles,
	}

	return jwt.MapClaims{
		"exp":  timeIssuedAt,
		"iat":  timeExpiredAt,
		"data": userData,
	}
}

func signToken(payload jwt.Claims) (any, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	stringToken, err := token.SignedString([]byte(config.GetAppConfig().JWT_SECRET_KEY))

	if err != nil {
		log.Printf("[SignToken - JWT] err:%s\n", err.Error())
		return nil, err
	}

	return stringToken, nil
}

func GenerateToken(data userrepository.UserRoles) (any, error) {
	payload := generatePayload(data)
	return signToken(payload)
}
