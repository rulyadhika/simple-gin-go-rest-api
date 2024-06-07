package accounthandler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	accountservice "github.com/rulyadhika/simple-gin-go-rest-api/service/account_service"
)

type accountHandlerImpl struct {
	as accountservice.AccountService
}

func NewAccountHandlerImpl(as accountservice.AccountService) AccountHandler {
	return &accountHandlerImpl{as: as}
}

func (a *accountHandlerImpl) Activation(ctx *gin.Context) {
	token := strings.TrimSpace(ctx.Param("token"))

	if token == "" {
		unprocessableEntityError := errs.NewUnprocessableEntityError("param token must be a valid string")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	if err := a.as.Activation(ctx, token); err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully activated the account",
		Data:       nil,
	}

	ctx.JSON(http.StatusOK, response)
}

func (a *accountHandlerImpl) ResendActivationToken(ctx *gin.Context) {
	resendTokenDto := &dto.ResendActivationTokenRequest{}

	if err := ctx.ShouldBindJSON(resendTokenDto); err != nil {
		unprocessableEntityError := errs.NewUnprocessableEntityError("invalid json request body")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := a.as.ResendActivationToken(ctx, *resendTokenDto)

	if err != nil {
		// only show error with status code other than 404
		if err.Status() != http.StatusText(http.StatusNotFound) {
			ctx.AbortWithStatusJSON(err.StatusCode(), err)
			return
		}
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "Success! If the email you submitted is valid and registered in our system, you will receive an account activation link shortly.",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}
