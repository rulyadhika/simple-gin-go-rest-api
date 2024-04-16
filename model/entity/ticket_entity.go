package entity

import "time"

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
	Id          uint32
	TicketId    string
	Title       string
	Description string
	Priority    TicketPriority
	Status      TicketStatus
	CreatedBy   uint32
	AssignTo    uint32
	AssignBy    uint32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
