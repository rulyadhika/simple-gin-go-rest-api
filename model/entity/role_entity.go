package entity

type userRole string

const (
	Role_ADMINISTRATOR      userRole = "administrator"
	Role_SUPPORT_SUPERVISOR userRole = "support supervisor"
	Role_SUPPORT_AGENT      userRole = "support agent"
	Role_CLIENT             userRole = "client"
)

type Role struct {
	Id       uint32
	RoleName string
}
