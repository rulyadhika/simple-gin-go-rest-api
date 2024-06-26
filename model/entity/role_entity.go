package entity

import "github.com/google/uuid"

type UserType string

func (u *UserType) ToString() string {
	return string(*u)
}

const (
	Role_ADMINISTRATOR      UserType = "administrator"
	Role_SUPPORT_SUPERVISOR UserType = "support-supervisor"
	Role_SUPPORT_AGENT      UserType = "support-agent"
	Role_CLIENT             UserType = "client"
)

type Role struct {
	Id       uuid.UUID
	RoleName UserType
}
