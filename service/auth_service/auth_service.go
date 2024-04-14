package authservice

import (
	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
)

type AuthService interface {
	Register(ctx *gin.Context, userDto *dto.RegisterUserRequest) (*dto.RegisterUserResponse, errs.Error)
	Login(ctx *gin.Context, userDto *dto.LoginUserRequest) (*dto.LoginUserResponse, errs.Error)
}
