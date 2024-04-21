package tickethandler

import "github.com/gin-gonic/gin"

type TicketHandler interface {
	Create(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindOneByTicketId(ctx *gin.Context)
	AssignTicketToUser(ctx *gin.Context)
	UpdateTicketStatus(ctx *gin.Context)
}
