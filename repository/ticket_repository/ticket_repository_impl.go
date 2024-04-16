package ticketrepository

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type ticketRepositoryImpl struct{}

func NewTicketRepositoryImpl() TicketRepository {
	return &ticketRepositoryImpl{}
}

func (t *ticketRepositoryImpl) Create(ctx *gin.Context, db *sql.DB, ticket entity.Ticket) (*entity.Ticket, errs.Error) {
	sqlQuery := `INSERT INTO tickets(ticket_id, title, description, priority, status, created_by) 
	VALUES($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`

	err := db.QueryRowContext(ctx, sqlQuery, ticket.TicketId, ticket.Title, ticket.Description, ticket.Priority, ticket.Status, ticket.CreatedBy).Scan(&ticket.Id, &ticket.CreatedAt, &ticket.UpdatedAt)

	if err != nil {
		log.Printf("[CreateTicket - Repo] err: %s\n", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &ticket, nil
}
