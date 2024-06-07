package accountservice

import (
	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
)

type AccountService interface {
	Activation(ctx *gin.Context, token string) errs.Error
	ResendActivationToken(ctx *gin.Context, resendTokenDto dto.ResendActivationTokenRequest) (*dto.ResendActivationTokenRespone, errs.Error)
}
