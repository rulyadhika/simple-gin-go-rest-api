package routes

import (
	"github.com/gin-gonic/gin"
	tickethandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/ticket_handler"
	authmiddleware "github.com/rulyadhika/simple-gin-go-rest-api/infra/middleware/auth_middleware"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

func NewTicketRoutes(r *gin.Engine, handler tickethandler.TicketHandler, authMiddleware authmiddleware.AuthMiddleware) {
	ticketGroup := r.Group("/tickets")
	{
		ticketGroup.Use(authMiddleware.Authentication())
		ticketGroup.POST("/", handler.Create)
		ticketGroup.GET("/", handler.FindAll)

		ticketGroup.GET("/:ticketId", authMiddleware.AuthorizationTicket(), handler.FindOneByTicketId)
		ticketGroup.PUT("/:ticketId/assignments/:userId", authMiddleware.RoleAuthorization(entity.Role_SUPPORT_SUPERVISOR), handler.AssignTicketToUser)
		ticketGroup.PUT("/:ticketId/status/", authMiddleware.RoleAuthorization(entity.Role_SUPPORT_SUPERVISOR, entity.Role_CLIENT), handler.UpdateTicketStatus)
	}
}
