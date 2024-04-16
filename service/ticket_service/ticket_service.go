package ticketservice

import (
	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
)

type TicketService interface {
	Create(ctx *gin.Context, ticketDto dto.NewTicketRequest) (*dto.NewTicketResponse, errs.Error)
}
