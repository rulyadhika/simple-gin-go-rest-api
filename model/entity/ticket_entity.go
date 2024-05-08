package entity

import (
	"time"

	"github.com/google/uuid"
)

type TicketPriority string
type TicketStatus string

var (
	TicketStatus_OPEN        TicketStatus = "open"
	TicketStatus_IN_PROGRESS TicketStatus = "in progress"
	TicketStatus_RESOLVED    TicketStatus = "resolved"
	TicketStatus_CLOSED      TicketStatus = "closed"
)

var (
	TicketPriority_LOW      TicketPriority = "low"
	TicketPriority_MED      TicketPriority = "med"
	TicketPriority_HIGH     TicketPriority = "high"
	TicketPriority_CRITICAL TicketPriority = "critical"
)

type Ticket struct {
	Id          uuid.UUID
	Title       string
	Description string
	Priority    TicketPriority
	Status      TicketStatus
	CreatedBy   uuid.UUID
	AssignTo    uuid.UUID
	AssignBy    uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
