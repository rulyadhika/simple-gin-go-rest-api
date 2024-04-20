package ticketrepository

import (
	"database/sql"
	"errors"
	"log"
	"time"

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

func (t *ticketRepositoryImpl) FindAll(ctx *gin.Context, db *sql.DB) (*[]TicketUser, errs.Error) {
	sqlQuery := `SELECT tickets.id, ticket_id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id`

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

func (t *ticketRepositoryImpl) FindAllByUserId(ctx *gin.Context, db *sql.DB, userId uint32) (*[]TicketUser, errs.Error) {
	sqlQuery := `SELECT tickets.id, ticket_id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id 
	WHERE tickets.created_by = $1`

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
	sqlQuery := `SELECT tickets.id, ticket_id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id 
	WHERE tickets.ticket_id = $1`

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
	sqlQuery := `UPDATE tickets SET assign_to=$1, assign_by=$2, status=$3, updated_at=$4 WHERE ticket_id=$5 RETURNING id`

	if err := db.QueryRowContext(ctx, sqlQuery, ticket.AssignTo, ticket.AssignBy, ticket.Status, time.Now(), ticket.TicketId).Scan(&ticket.Id); err != nil {
		log.Printf("[AssignTicketToUser - Repo], err: %s\n", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("ticket not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}

	queryGetData := `SELECT tickets.id, ticket_id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id 
	WHERE tickets.ticket_id = $1`

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
	sqlQuery := `UPDATE tickets SET status=$1, updated_at=$2 WHERE ticket_id=$3 RETURNING id`

	if err := db.QueryRowContext(ctx, sqlQuery, ticket.Status, time.Now(), ticket.TicketId).Scan(&ticket.Id); err != nil {
		log.Printf("[UpdateTicketStatus - Repo], err: %s\n", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("ticket not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}

	queryGetData := `SELECT tickets.id, ticket_id, title, description, priority, status, tickets.created_at, tickets.updated_at,
	a.username, a.email, b.username, b.email, c.username, c.email
	FROM tickets JOIN users AS a ON tickets.created_by = a.id
	LEFT JOIN users AS b ON tickets.assign_to = b.id
	LEFT JOIN users AS c ON tickets.assign_by = c.id 
	WHERE tickets.ticket_id = $1`

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
