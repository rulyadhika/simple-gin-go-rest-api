package dto

import (
	"time"

	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type NewTicketRequest struct {
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Priority    entity.TicketPriority `json:"priority" validate:"required,ticket_priority_custom_validation"`
	Status      entity.TicketStatus   `validate:"required,ticket_status_custom_validation"`
	CreatedBy   uint32                `validate:"required"`
}

type NewTicketResponse struct {
	Id          uint32                `json:"id"`
	TicketId    string                `json:"ticket_id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Priority    entity.TicketPriority `json:"priority"`
	Status      entity.TicketStatus   `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}
