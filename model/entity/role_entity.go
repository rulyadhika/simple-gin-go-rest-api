package entity

type userRole uint32

const (
	Role_ADMINISTRATOR      userRole = 1
	Role_SUPPORT_SUPERVISOR userRole = 2
	Role_SUPPORT_AGENT      userRole = 3
	Role_CLIENT             userRole = 4
)

type Role struct {
	Id       uint32
	RoleName string
}
