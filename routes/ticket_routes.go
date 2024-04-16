package routes

import (
	"github.com/gin-gonic/gin"
	tickethandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/ticket_handler"
	authmiddleware "github.com/rulyadhika/simple-gin-go-rest-api/infra/middleware/auth_middleware"
)

func NewTicketRoutes(r *gin.Engine, handler tickethandler.TicketHandler, authMiddleware authmiddleware.AuthMiddleware) {
	ticketGroup := r.Group("/tickets")
	{
		ticketGroup.Use(authMiddleware.Authentication())
		ticketGroup.POST("/", handler.Create)
	}
}
