package ticketservice

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type TicketService interface {
	Create(ctx *gin.Context, ticketDto dto.NewTicketRequest) (*dto.NewTicketResponse, errs.Error)
	FindAll(ctx *gin.Context, userId uuid.UUID, userRoles []entity.UserType) (*[]dto.TicketResponse, errs.Error)
	FindOneByTicketId(ctx *gin.Context, ticketId string) (*dto.TicketResponse, errs.Error)
	AssignTicketToUser(ctx *gin.Context, ticketDto dto.AssignTicketToUserRequest) (*dto.TicketResponse, errs.Error)
	UpdateTicketStatus(ctx *gin.Context, ticketDto dto.UpdateTicketStatusRequest, userRoles []entity.UserType) (*dto.TicketResponse, errs.Error)
}
