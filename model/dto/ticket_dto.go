package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type NewTicketRequest struct {
	Title       string                `json:"title" validate:"required"`
	Description string                `json:"description" validate:"required"`
	Priority    entity.TicketPriority `json:"priority" validate:"required,ticket_priority_custom_validation"`
	Status      entity.TicketStatus   `validate:"required,ticket_status_custom_validation"`
	CreatedBy   uuid.UUID             `validate:"required"`
}

type NewTicketResponse struct {
	Id          uuid.UUID             `json:"id"`
	TicketId    string                `json:"ticket_id"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Priority    entity.TicketPriority `json:"priority"`
	Status      entity.TicketStatus   `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type TicketResponse struct {
	Id          uuid.UUID              `json:"id"`
	TicketId    string                 `json:"ticket_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Priority    entity.TicketPriority  `json:"priority"`
	Status      entity.TicketStatus    `json:"status"`
	CreatedBy   TicketResponseUserData `json:"created_by"`
	AssignTo    TicketResponseUserData `json:"assign_to"`
	AssignBy    TicketResponseUserData `json:"assign_by"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type TicketResponseUserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AssignTicketToUserRequest struct {
	TicketId   string
	AssignToId uuid.UUID
	AssignById uuid.UUID
}

type UpdateTicketStatusRequest struct {
	TicketId string
	Status   entity.TicketStatus `json:"status" validate:"required,ticket_status_custom_validation"`
}
