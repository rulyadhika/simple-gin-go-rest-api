package entity

import "time"

type UserRole struct {
	Id        uint32
	UserId    uint32
	RoleId    uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}
