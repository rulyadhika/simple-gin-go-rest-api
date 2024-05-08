package ticketservice

import (
	"database/sql"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/helper"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

type ticketServiceImpl struct {
	db        *sql.DB
	tr        ticketrepository.TicketRepository
	ur        userrepository.UserRepository
	validator *validator.Validate
}

func NewTicketServiceImpl(tr ticketrepository.TicketRepository, ur userrepository.UserRepository, db *sql.DB, validator *validator.Validate) TicketService {
	return &ticketServiceImpl{
		db,
		tr,
		ur,
		validator,
	}
}

func (t *ticketServiceImpl) Create(ctx *gin.Context, ticketDto dto.NewTicketRequest) (*dto.NewTicketResponse, errs.Error) {
	ticketDto.Status = entity.TicketStatus_OPEN

	if errValidate := t.validator.Struct(ticketDto); errValidate != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(errValidate))
	}

	ticketId := uuid.NewString()

	ticket := entity.Ticket{
		Title:       ticketDto.Title,
		Description: ticketDto.Description,
		Priority:    ticketDto.Priority,
		Status:      ticketDto.Status,
		CreatedBy:   ticketDto.CreatedBy,
		TicketId:    ticketId,
	}

	result, err := t.tr.Create(ctx, t.db, ticket)

	if err != nil {
		return nil, err
	}

	return &dto.NewTicketResponse{
		Id:          result.Id,
		TicketId:    result.TicketId,
		Title:       result.Title,
		Description: result.Description,
		Priority:    result.Priority,
		Status:      result.Status,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (t *ticketServiceImpl) FindAll(ctx *gin.Context, userId uuid.UUID, userRoles []entity.UserType) (*[]dto.TicketResponse, errs.Error) {
	var result *[]ticketrepository.TicketUser
	var err errs.Error

	// if user's roles is only client
	if slices.Contains(userRoles, "client") && len(userRoles) == 1 {
		// then only select all the tickets created by that user
		result, err = t.tr.FindAllByUserId(ctx, t.db, userId)
	} else {
		// else select all tickets with no exception
		result, err = t.tr.FindAll(ctx, t.db)
	}

	if err != nil {
		return nil, err
	}

	ticketsResponse := []dto.TicketResponse{}

	for _, data := range *result {
		ticketResponse := *helper.ToDtoTicketResponse(&data)

		ticketsResponse = append(ticketsResponse, ticketResponse)
	}

	return &ticketsResponse, nil
}

func (t *ticketServiceImpl) FindOneByTicketId(ctx *gin.Context, ticketId string) (*dto.TicketResponse, errs.Error) {
	result, err := t.tr.FindOneByTicketId(ctx, t.db, ticketId)

	if err != nil {
		return nil, err
	}

	ticketResponse := helper.ToDtoTicketResponse(result)

	return ticketResponse, nil
}

func (t *ticketServiceImpl) AssignTicketToUser(ctx *gin.Context, ticketDto dto.AssignTicketToUserRequest) (*dto.TicketResponse, errs.Error) {
	ticket := entity.Ticket{
		TicketId: ticketDto.TicketId,
		AssignTo: ticketDto.AssignToId,
		AssignBy: ticketDto.AssignById,
		Status:   entity.TicketStatus_IN_PROGRESS,
	}

	user, err := t.ur.FindById(ctx, t.db, ticket.AssignTo)
	if err != nil {
		return nil, err
	}

	userRoles := []entity.UserType{}
	for _, role := range user.Roles {
		userRoles = append(userRoles, role.RoleName)
	}

	if isSupportAgent := slices.Contains(userRoles, entity.Role_SUPPORT_AGENT); !isSupportAgent {
		return nil, errs.NewConflictError("user is not a support agent")
	}

	result, err := t.tr.AssignTicketToUser(ctx, t.db, ticket)

	if err != nil {
		return nil, err
	}

	ticketResponse := helper.ToDtoTicketResponse(result)

	return ticketResponse, nil
}

func (t *ticketServiceImpl) UpdateTicketStatus(ctx *gin.Context, ticketDto dto.UpdateTicketStatusRequest, userRoles []entity.UserType) (*dto.TicketResponse, errs.Error) {
	if validationErr := t.validator.Struct(ticketDto); validationErr != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(validationErr))
	}

	ticket := entity.Ticket{
		TicketId: ticketDto.TicketId,
		Status:   ticketDto.Status,
	}

	// check whether the user is eligible to update ticket status based on their roles and what ticket status choosen
	// entity.Role_SUPPORT_SUPERVISOR is only allowed to change ticket status to entity.TicketStatus_CLOSED
	// entity.Role_CLIENT is only allowed to change ticket status to entity.TicketStatus_RESOLVED
	if eligible := slices.Contains(userRoles, entity.Role_SUPPORT_SUPERVISOR) && ticket.Status == entity.TicketStatus_CLOSED || slices.Contains(userRoles, entity.Role_CLIENT) && ticket.Status == entity.TicketStatus_RESOLVED; !eligible {
		return nil, errs.NewForbiddenError("you're not allowed to change this support ticket status to: " + string(ticket.Status))
	}

	result, err := t.tr.UpdateTicketStatus(ctx, t.db, ticket)
	if err != nil {
		return nil, err
	}

	ticketResponse := helper.ToDtoTicketResponse(result)

	return ticketResponse, nil
}
