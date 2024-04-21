package helper

import (
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
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
