package ticketrepository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type ticketRepositoryImpl struct{}

func NewTicketRepositoryImpl() TicketRepository {
	return &ticketRepositoryImpl{}
}

func (t *ticketRepositoryImpl) Create(ctx *gin.Context, db *sql.DB, ticket entity.Ticket) (*entity.Ticket, errs.Error) {
	sqlQuery := createNewTicketQuery

	err := db.QueryRowContext(ctx, sqlQuery, ticket.TicketId, ticket.Title, ticket.Description, ticket.Priority, ticket.Status, ticket.CreatedBy).Scan(&ticket.Id, &ticket.CreatedAt, &ticket.UpdatedAt)

	if err != nil {
		log.Printf("[CreateTicket - Repo] err: %s\n", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &ticket, nil
}

func (t *ticketRepositoryImpl) FindAll(ctx *gin.Context, db *sql.DB) (*[]TicketUser, errs.Error) {
	sqlQuery := findAllTicketQuery

	ticketsUser := []TicketUser{}

	rows, err := db.QueryContext(ctx, sqlQuery)

	if err != nil {
		log.Printf("[FindAllTickets - Repo] err: %s\n", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	for rows.Next() {
		ticketUser := TicketUser{}

		err := rows.Scan(&ticketUser.Id, &ticketUser.TicketId, &ticketUser.Title, &ticketUser.Description, &ticketUser.Priority, &ticketUser.Status,
			&ticketUser.CreatedAt, &ticketUser.UpdatedAt, &ticketUser.CreatedBy.Username, &ticketUser.CreatedBy.Email,
			&ticketUser.AssignTo.Username, &ticketUser.AssignTo.Email, &ticketUser.AssignBy.Username, &ticketUser.AssignBy.Email)

		if err != nil {
			log.Printf("[FindAllTickets - Repo] err: %s\n", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		ticketsUser = append(ticketsUser, ticketUser)
	}

	// if the result is empty
	if len(ticketsUser) == 0 {
		return nil, errs.NewNotFoundError("not tickets found")
	}

	return &ticketsUser, nil
}

func (t *ticketRepositoryImpl) FindAllByUserId(ctx *gin.Context, db *sql.DB, userId uuid.UUID) (*[]TicketUser, errs.Error) {
	sqlQuery := findAllTicketByUserIdQuery

	ticketsUser := []TicketUser{}

	rows, err := db.QueryContext(ctx, sqlQuery, userId)

	if err != nil {
		log.Printf("[FindAllTicketsByUserId - Repo] err: %s\n", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	for rows.Next() {
		ticketUser := TicketUser{}

		err := rows.Scan(&ticketUser.Id, &ticketUser.TicketId, &ticketUser.Title, &ticketUser.Description, &ticketUser.Priority, &ticketUser.Status,
			&ticketUser.CreatedAt, &ticketUser.UpdatedAt, &ticketUser.CreatedBy.Username, &ticketUser.CreatedBy.Email,
			&ticketUser.AssignTo.Username, &ticketUser.AssignTo.Email, &ticketUser.AssignBy.Username, &ticketUser.AssignBy.Email)

		if err != nil {
			log.Printf("[FindAllTicketsByUserId - Repo] err: %s\n", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		ticketsUser = append(ticketsUser, ticketUser)
	}

	// if the result is empty
	if len(ticketsUser) == 0 {
		return nil, errs.NewNotFoundError("not tickets found")
	}

	return &ticketsUser, nil
}

func (t *ticketRepositoryImpl) FindOneByTicketId(ctx *gin.Context, db *sql.DB, ticketId string) (*TicketUser, errs.Error) {
	sqlQuery := findOneTicketByTicketIdQuery

	ticketUser := TicketUser{}

	err := db.QueryRowContext(ctx, sqlQuery, ticketId).Scan(&ticketUser.Id, &ticketUser.TicketId, &ticketUser.Title, &ticketUser.Description, &ticketUser.Priority, &ticketUser.Status,
		&ticketUser.CreatedAt, &ticketUser.UpdatedAt, &ticketUser.CreatedBy.Username, &ticketUser.CreatedBy.Email,
		&ticketUser.AssignTo.Username, &ticketUser.AssignTo.Email, &ticketUser.AssignBy.Username, &ticketUser.AssignBy.Email)

	if err != nil {
		log.Printf("[FindOneByTicketId - Repo] err: %s\n", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("ticket not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}
	return &ticketUser, nil
}

func (t *ticketRepositoryImpl) AssignTicketToUser(ctx *gin.Context, db *sql.DB, ticket entity.Ticket) (*TicketUser, errs.Error) {
	sqlQuery := assignTicketToUserQuery

	if err := db.QueryRowContext(ctx, sqlQuery, ticket.AssignTo, ticket.AssignBy, ticket.Status, time.Now(), ticket.TicketId).Scan(&ticket.Id); err != nil {
		log.Printf("[AssignTicketToUser - Repo], err: %s\n", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("ticket not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}

	queryGetData := findOneTicketByTicketIdQuery

	ticketUser := TicketUser{}

	err := db.QueryRowContext(ctx, queryGetData, ticket.TicketId).Scan(&ticketUser.Id, &ticketUser.TicketId, &ticketUser.Title, &ticketUser.Description, &ticketUser.Priority, &ticketUser.Status,
		&ticketUser.CreatedAt, &ticketUser.UpdatedAt, &ticketUser.CreatedBy.Username, &ticketUser.CreatedBy.Email,
		&ticketUser.AssignTo.Username, &ticketUser.AssignTo.Email, &ticketUser.AssignBy.Username, &ticketUser.AssignBy.Email)

	if err != nil {
		log.Printf("[AssignTicketToUser - Repo], err: %s\n", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &ticketUser, nil
}

func (t *ticketRepositoryImpl) UpdateTicketStatus(ctx *gin.Context, db *sql.DB, ticket entity.Ticket) (*TicketUser, errs.Error) {
	sqlQuery := updateTicketStatusQuery

	if err := db.QueryRowContext(ctx, sqlQuery, ticket.Status, time.Now(), ticket.TicketId).Scan(&ticket.Id); err != nil {
		log.Printf("[UpdateTicketStatus - Repo], err: %s\n", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("ticket not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}

	queryGetData := findOneTicketByTicketIdQuery

	ticketUser := TicketUser{}

	err := db.QueryRowContext(ctx, queryGetData, ticket.TicketId).Scan(&ticketUser.Id, &ticketUser.TicketId, &ticketUser.Title, &ticketUser.Description, &ticketUser.Priority, &ticketUser.Status,
		&ticketUser.CreatedAt, &ticketUser.UpdatedAt, &ticketUser.CreatedBy.Username, &ticketUser.CreatedBy.Email,
		&ticketUser.AssignTo.Username, &ticketUser.AssignTo.Email, &ticketUser.AssignBy.Username, &ticketUser.AssignBy.Email)

	if err != nil {
		log.Printf("[UpdateTicketStatus - Repo], err: %s\n", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &ticketUser, nil
}
