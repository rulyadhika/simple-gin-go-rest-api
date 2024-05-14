package routes

import (
	"github.com/gin-gonic/gin"
	accounthandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/account_handler"
)

func NewAccountRoutes(r *gin.Engine, handler accounthandler.AccountHandler) {
	accountRoute := r.Group("/accounts")
	{
		accountRoute.GET("/activation/:token", handler.Activation)
	}
}
