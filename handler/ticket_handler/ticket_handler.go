package tickethandler

import "github.com/gin-gonic/gin"

type TicketHandler interface {
	Create(ctx *gin.Context)
}
