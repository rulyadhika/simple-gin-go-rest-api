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
		StatusCode: http.StatusCreated,
		Status:     http.StatusText(http.StatusCreated),
		Message:    "successfully registered a new user. please check your email to activate the account",
		Data:       result,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (a *AuthHandlerImpl) Login(ctx *gin.Context) {
	userDto := &dto.LoginUserRequest{}

	if err := ctx.ShouldBindJSON(userDto); err != nil {
		log.Printf("[Login - Handler], err: %s", err.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("invalid json request body")

		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := a.AuthService.Login(ctx, userDto)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "login successful",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *AuthHandlerImpl) RefreshToken(ctx *gin.Context) {
	refreshToken, errCookie := ctx.Cookie("refresh-token")

	if errCookie != nil {
		unauthorizedError := errs.NewUnauthorizedError("refresh-token cookie not present")
		ctx.AbortWithStatusJSON(unauthorizedError.StatusCode(), unauthorizedError)
		return
	}

	userDto := &dto.RefreshTokenRequest{
		Token: refreshToken,
	}

	result, err := a.AuthService.RefreshToken(ctx, userDto)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully obtain a new access token",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}
