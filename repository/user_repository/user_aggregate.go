package userrepository

import "github.com/rulyadhika/simple-gin-go-rest-api/model/entity"

type UserRole struct {
	entity.User
	entity.Role
}

type UserRoles struct {
	entity.User
	roles []entity.Role
}

func (u *UserRoles) HandleMappingUserRoles(userRole []UserRole) {
	roles := []entity.Role{}

	for index, data := range userRole {
		if index == 0 {
			u.Id = data.User.Id
			u.Username = data.Username
			u.Email = data.Email
			u.Password = data.Password
			u.CreatedAt = data.CreatedAt
			u.UpdatedAt = data.UpdatedAt
		}

		roles = append(roles, entity.Role{
			Id:       data.Role.Id,
			RoleName: data.RoleName,
		})
	}

	u.roles = roles
}
