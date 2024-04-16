package tickethandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/jwt"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	ticketservice "github.com/rulyadhika/simple-gin-go-rest-api/service/ticket_service"
)

type ticketHandlerImpl struct {
	ts ticketservice.TicketService
}

func NewTicketHandlerImpl(ts ticketservice.TicketService) TicketHandler {
	return &ticketHandlerImpl{
		ts,
	}
}

func (t *ticketHandlerImpl) Create(ctx *gin.Context) {
	ticketDto := &dto.NewTicketRequest{}

	if err := ctx.ShouldBindJSON(ticketDto); err != nil {
		unprocessableEntityError := errs.NewUnprocessableEntityError("invalid json request body")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	userData := ctx.MustGet("userData").(*jwt.JWTPayload)
	ticketDto.CreatedBy = userData.Id

	result, err := t.ts.Create(ctx, *ticketDto)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusCreated,
		Status:     http.StatusText(http.StatusCreated),
		Message:    "successfully created new ticket",
		Data:       result,
	}

	ctx.JSON(http.StatusCreated, response)
}
