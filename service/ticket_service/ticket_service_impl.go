package ticketservice

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
)

type ticketServiceImpl struct {
	db        *sql.DB
	tr        ticketrepository.TicketRepository
	validator *validator.Validate
}

func NewTicketServiceImpl(tr ticketrepository.TicketRepository, db *sql.DB, validator *validator.Validate) TicketService {
	return &ticketServiceImpl{
		db,
		tr,
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
