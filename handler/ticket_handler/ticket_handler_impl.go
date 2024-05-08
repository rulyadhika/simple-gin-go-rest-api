package tickethandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	userData, ok := ctx.MustGet("userData").(*jwt.JWTPayload)

	if !ok {
		log.Printf("[CreateTicket - Handler] err: %s\n", "failed type casting to '*jwt.JWTPayload'")
		internalServerErr := errs.NewInternalServerError("something went wrong")
		ctx.AbortWithStatusJSON(internalServerErr.StatusCode(), internalServerErr)
		return
	}

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

func (t *ticketHandlerImpl) FindAll(ctx *gin.Context) {
	userData, ok := ctx.MustGet("userData").(*jwt.JWTPayload)

	if !ok {
		log.Printf("[FindAllTickets - Handler] err: %s\n", "failed type casting to '*jwt.JWTPayload'")
		internalServerErr := errs.NewInternalServerError("something went wrong")
		ctx.AbortWithStatusJSON(internalServerErr.StatusCode(), internalServerErr)
		return
	}

	result, err := t.ts.FindAll(ctx, userData.Id, userData.Roles)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully get all tickets",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (t *ticketHandlerImpl) FindOneByTicketId(ctx *gin.Context) {
	ticketId, errParseUUID := uuid.Parse(ctx.Param("ticketId"))
	if errParseUUID != nil {
		log.Printf("[FindOneByTicketId - Handler], err: %s\n", errParseUUID.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("param ticketId must be a valid id")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := t.ts.FindOneByTicketId(ctx, ticketId)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully get a ticket",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (t *ticketHandlerImpl) AssignTicketToUser(ctx *gin.Context) {
	user, ok := ctx.MustGet("userData").(*jwt.JWTPayload)

	if !ok {
		log.Printf("[AssignTicketToUser - Handler] err: %s\n", "failed type casting to '*jwt.JWTPayload'")
		internalServerErr := errs.NewInternalServerError("something went wrong")
		ctx.AbortWithStatusJSON(internalServerErr.StatusCode(), internalServerErr)
		return
	}

	ticketId, errParseUUID := uuid.Parse(ctx.Param("ticketId"))
	if errParseUUID != nil {
		log.Printf("[AssignTicketToUser - Handler], err: %s\n", errParseUUID.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("param ticketId must be a valid id")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	userId, errConvert := uuid.Parse(ctx.Param("userId"))
	if errConvert != nil {
		log.Printf("[AssignTicketToUser - Handler], err: %s\n", errConvert.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("param userId must be a valid id")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	ticketDto := &dto.AssignTicketToUserRequest{
		TicketId:   ticketId,
		AssignToId: userId,
		AssignById: user.Id,
	}

	result, err := t.ts.AssignTicketToUser(ctx, *ticketDto)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := &dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully asign ticket to support agent",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)

}

func (t *ticketHandlerImpl) UpdateTicketStatus(ctx *gin.Context) {
	ticketDto := &dto.UpdateTicketStatusRequest{}

	if err := ctx.ShouldBindJSON(ticketDto); err != nil {
		unprocessableEntityError := errs.NewUnprocessableEntityError("invalid json request body")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	user, ok := ctx.MustGet("userData").(*jwt.JWTPayload)

	if !ok {
		log.Printf("[UpdateTicketStatus - Handler] err: %s\n", "failed type casting to '*jwt.JWTPayload'")
		internalServerErr := errs.NewInternalServerError("something went wrong")
		ctx.AbortWithStatusJSON(internalServerErr.StatusCode(), internalServerErr)
		return
	}

	ticketId, errParseUUID := uuid.Parse(ctx.Param("ticketId"))
	if errParseUUID != nil {
		log.Printf("[UpdateTicketStatus - Handler], err: %s\n", errParseUUID.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("param ticketId must be a valid id")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}
	ticketDto.TicketId = ticketId

	result, err := t.ts.UpdateTicketStatus(ctx, *ticketDto, user.Roles)
	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := &dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully update ticket status",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}
