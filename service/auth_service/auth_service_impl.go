package authservice

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

type AuthServiceImpl struct {
	UserRepository userrepository.UserRepository
	DB             *sql.DB
	Validate       *validator.Validate
}

func NewAuthServiceImpl(userRepository userrepository.UserRepository, db *sql.DB, validator *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		DB:             db,
		Validate:       validator,
	}
}

func (a *AuthServiceImpl) Register(ctx *gin.Context, userDto *dto.RegisterUserRequest) (*dto.RegisterUserResponse, errs.Error) {
	if errValidate := a.Validate.Struct(userDto); errValidate != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(errValidate))
	}

	tx, txErr := a.DB.Begin()

	if txErr != nil {
		log.Printf("[Register - Service] err: %v", txErr.Error())

		if errRollback := tx.Rollback(); errRollback != nil {
			log.Printf("[Register - Service] err: %v", errRollback.Error())
		}
		return nil, errs.NewInternalServerError("something went wrong")
	}

	user := entity.User{
		Username: userDto.Username,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	findByEmail, err := a.UserRepository.FindByEmail(ctx, a.DB, user.Email)

	if err != nil {
		return nil, err
	}

	if findByEmail != nil {
		return nil, errs.NewConflictError("email has already been taken")
	}

	findByUsername, err := a.UserRepository.FindByUsername(ctx, a.DB, user.Username)
	if err != nil {
		return nil, err
	}

	if findByUsername != nil {
		return nil, errs.NewConflictError("username has already been taken")
	}

	if err := user.HashPassword(); err != nil {
		log.Printf("[Register - Service] err: %v", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	result, err := a.UserRepository.Create(ctx, tx, user)
	if err != nil {
		return nil, err
	}

	usersRole := []entity.UserRole{}
	usersRole = append(usersRole, entity.UserRole{UserId: result.Id, RoleId: uint32(entity.Role_CLIENT)})

	err = a.UserRepository.AssignRolesToUser(ctx, tx, usersRole)
	if err != nil {
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[Register - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &dto.RegisterUserResponse{Username: result.Username, Email: result.Email}, nil
}
