package ticketservice

import (
	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
)

type TicketService interface {
	Create(ctx *gin.Context, ticketDto dto.NewTicketRequest) (*dto.NewTicketResponse, errs.Error)
	FindAll(ctx *gin.Context, userId uint32, userRoles []string) (*[]dto.TicketResponse, errs.Error)
	FindOneByTicketId(ctx *gin.Context, ticketId string) (*dto.TicketResponse, errs.Error)
}
