package authhandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	authservice "github.com/rulyadhika/simple-gin-go-rest-api/service/auth_service"
)

type AuthHandlerImpl struct {
	AuthService authservice.AuthService
}

func NewAuthHandlerImpl(authService authservice.AuthService) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
	}
}

func (a *AuthHandlerImpl) Register(ctx *gin.Context) {
	userDto := &dto.RegisterUserRequest{}

	if err := ctx.ShouldBindJSON(userDto); err != nil {
		log.Printf("[Register - Handler], err: %s", err.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("invalid json request body")

		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := a.AuthService.Register(ctx, userDto)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully registered a new user",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}
