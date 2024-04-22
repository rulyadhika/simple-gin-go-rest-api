package userservice

import (
	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
)

type UserService interface {
	FindAll(ctx *gin.Context) (*[]dto.UserResponse, errs.Error)
	FindOneByUsername(ctx *gin.Context, username string) (*dto.UserResponse, errs.Error)
}
