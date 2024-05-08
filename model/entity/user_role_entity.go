package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	RoleId    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}
