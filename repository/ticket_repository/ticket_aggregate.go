package ticketrepository

import (
	"database/sql"

	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type TicketUser struct {
	entity.Ticket
	CreatedBy user
	AssignBy  user
	AssignTo  user
}

type user struct {
	Username sql.NullString
	Email    sql.NullString
}
