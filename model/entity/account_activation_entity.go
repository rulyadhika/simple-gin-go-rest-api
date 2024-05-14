package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccountActivation struct {
	UserId         uuid.UUID
	Token          string
	RequestTime    time.Time
	ExpirationTime time.Time
}
