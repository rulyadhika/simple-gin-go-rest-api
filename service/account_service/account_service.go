package accountservice

import (
	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
)

type AccountService interface {
	Activation(ctx *gin.Context, token string) errs.Error
}
