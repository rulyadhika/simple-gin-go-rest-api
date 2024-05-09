package userhandler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	userservice "github.com/rulyadhika/simple-gin-go-rest-api/service/user_service"
)

type UserHandlerImpl struct {
	us userservice.UserService
}

func NewUserHandlerImpl(us userservice.UserService) UserHandler {
	return &UserHandlerImpl{
		us,
	}
}

func (u *UserHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := u.us.FindAll(ctx)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully get all users",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UserHandlerImpl) FindOneByUsername(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.Param("username"))

	if username == "" {
		unprocessableEntityError := errs.NewUnprocessableEntityError("param username must be a valid string")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := u.us.FindOneByUsername(ctx, username)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully get a user",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UserHandlerImpl) Create(ctx *gin.Context) {
	userDto := &dto.CreateNewUserRequest{}

	if err := ctx.ShouldBindJSON(userDto); err != nil {
		log.Printf("[CreateNewUser - Handler], err: %s\n", err.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("invalid json request body")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	result, err := u.us.Create(ctx, userDto)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusCreated,
		Status:     http.StatusText(http.StatusCreated),
		Message:    "successfully created a new user",
		Data:       result,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (u *UserHandlerImpl) AssignOrRemoveUserRole(ctx *gin.Context) {
	roleName := strings.TrimSpace(ctx.Param("roleName"))

	if roleName == "" {
		unprocessableEntityError := errs.NewUnprocessableEntityError("param roleName must be a valid string")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	userId, errConv := uuid.Parse(ctx.Param("userId"))
	if errConv != nil {
		log.Printf("[AssignOrRemoveUserRole - Handler], err: %s\n", errConv.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("param userId must be a valid id")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	userDto := &dto.AssignRoleToUserRequest{
		UserId: userId,
		Role:   entity.UserType(roleName),
	}

	result, err := u.us.AssignOrRemoveUserRole(ctx, userDto)

	if err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    "successfully assign or remove role to or from user",
		Data:       result,
	}

	ctx.JSON(http.StatusOK, response)
}

func (u *UserHandlerImpl) Delete(ctx *gin.Context) {
	userId, errConv := uuid.Parse(ctx.Param("userId"))

	if errConv != nil {
		log.Printf("[DeleteUser - Handler], err: %s\n", errConv.Error())
		unprocessableEntityError := errs.NewUnprocessableEntityError("param userId must be a valid id")
		ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
		return
	}

	if err := u.us.Delete(ctx, userId); err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode(), err)
		return
	}

	response := dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Message:    fmt.Sprintf("successfully delete user with id: %s", userId),
		Data:       nil,
	}

	ctx.JSON(http.StatusOK, response)
}
