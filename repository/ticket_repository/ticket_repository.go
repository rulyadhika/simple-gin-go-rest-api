package ticketrepository

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type TicketRepository interface {
	Create(ctx *gin.Context, db *sql.DB, ticket entity.Ticket) (*entity.Ticket, errs.Error)
	FindAll(ctx *gin.Context, db *sql.DB) (*[]TicketUser, errs.Error)
	FindAllByUserId(ctx *gin.Context, db *sql.DB, userId uint32) (*[]TicketUser, errs.Error)
	FindOneByTicketId(ctx *gin.Context, db *sql.DB, ticketId string) (*TicketUser, errs.Error)
}