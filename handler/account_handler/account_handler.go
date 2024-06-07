package accounthandler

import "github.com/gin-gonic/gin"

type AccountHandler interface {
	Activation(ctx *gin.Context)
	ResendActivationToken(ctx *gin.Context)
}
