package entity

import (
	"time"

	"github.com/google/uuid"
)

type AccountPasswordReset struct {
	UserId                 uuid.UUID
	Token                  string
	RequestTime            time.Time
	ExpirationTime         time.Time
	NextRequestAvailableAt time.Time
}
