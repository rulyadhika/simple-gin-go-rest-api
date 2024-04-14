package routes

import (
	"github.com/gin-gonic/gin"
	authhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/auth_handler"
)

func NewAuthRoutes(r *gin.Engine, handler authhandler.AuthHandler) {
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/register", handler.Register)
		authRoute.POST("/login", handler.Login)
	}
}
