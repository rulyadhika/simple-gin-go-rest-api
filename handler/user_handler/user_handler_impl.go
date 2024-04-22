package userhandler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	userservice "github.com/rulyadhika/simple-gin-go-rest-api/service/user_service"
)

type UserHandlerImpl struct {
	us userservice.UserService
}

func NewUserHandlerImpl(us userservice.UserService) UserHandler {
	return &UserHandlerImpl{
		us,
	}
}

func (u *UserHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := u.us.FindAll(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully get all users",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UserHandlerImpl) FindOneByUsername(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.Param("username"))

	if username == "" {
		unprocessableEntityError := errs.NewUnprocessableEntityError("param username must be a valid string")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := u.us.FindOneByUsername(ctx, username)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully get a user",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}
