package helper

import (
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

func ToDtoTicketResponse(data *ticketrepository.TicketUser) *dto.TicketResponse {
	ticketResponse := dto.TicketResponse{
		Id:          data.Id,
		TicketId:    data.TicketId,
		Title:       data.Title,
		Description: data.Description,
		Priority:    data.Priority,
		Status:      data.Status,
		CreatedBy: dto.TicketResponseUserData{
			Username: data.CreatedBy.Username.String,
			Email:    data.CreatedBy.Email.String,
		},
		AssignTo: dto.TicketResponseUserData{
			Username: data.AssignTo.Username.String,
			Email:    data.AssignTo.Email.String,
		},
		AssignBy: dto.TicketResponseUserData{
			Username: data.AssignBy.Username.String,
			Email:    data.AssignBy.Email.String,
		},
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return &ticketResponse
}

func ToDtoUsersResponse(data *[]userrepository.UserRoles) *[]dto.UserResponse {
	userResponse := []dto.UserResponse{}

	for _, user := range *data {
		roles := []dto.UserRolesResponse{}

		for _, role := range user.Roles {
			roles = append(roles, dto.UserRolesResponse{
				Id:   role.Id,
				Name: string(role.RoleName),
			})
		}

		ur := dto.UserResponse{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}

		userResponse = append(userResponse, ur)
	}

	return &userResponse
}

func ToDtoUserResponse(data *userrepository.UserRoles) *dto.UserResponse {
	roles := []dto.UserRolesResponse{}

	for _, role := range data.Roles {
		roles = append(roles, dto.UserRolesResponse{
			Id:   role.Id,
			Name: string(role.RoleName),
		})
	}

	userResponse := dto.UserResponse{
		Id:        data.Id,
		Username:  data.Username,
		Email:     data.Email,
		Roles:     roles,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return &userResponse
}
