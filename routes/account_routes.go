package routes

import (
	"github.com/gin-gonic/gin"
	accounthandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/account_handler"
)

func NewAccountRoutes(r *gin.Engine, handler accounthandler.AccountHandler) {
	accountRoute := r.Group("/accounts")
	{
		accountRoute.POST("/activation/:token", handler.Activation)
		accountRoute.POST("/resend-activation-token/", handler.ResendActivationToken)

		accountRoute.POST("/forgot-password/", handler.ForgotPassword)
	}
}
