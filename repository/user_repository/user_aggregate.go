package userrepository

import (
	"slices"

	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserRole struct {
	entity.User
	entity.Role
}

type UserRoles struct {
	entity.User
	Roles []entity.Role
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

	u.Roles = roles
}

func (u *UserRoles) HandleMappingUsersRoles(userRoles []UserRole) *[]UserRoles {
	allUsers := []UserRoles{}

	for _, user := range userRoles {
		index := slices.IndexFunc(allUsers, func(ur UserRoles) bool {
			return ur.Id == user.User.Id
		})

		if index == -1 {
			user := UserRoles{
				User: entity.User{
					Id:        user.User.Id,
					Username:  user.Username,
					Email:     user.Email,
					Password:  user.Password,
					CreatedAt: user.CreatedAt,
					UpdatedAt: user.UpdatedAt,
				}, Roles: []entity.Role{{
					Id:       user.Role.Id,
					RoleName: user.Role.RoleName,
				}},
			}
			allUsers = append(allUsers, user)
		} else {
			allUsers[index].Roles = append(allUsers[index].Roles, user.Role)
		}
	}

	return &allUsers
}
