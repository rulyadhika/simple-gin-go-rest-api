package userservice

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/helper"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	rolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/role_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	userrolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_role_repository"
)

type UserServiceImpl struct {
	ur       userrepository.UserRepository
	rr       rolerepository.RoleRepository
	urr      userrolerepository.UserRoleRepository
	db       *sql.DB
	validate *validator.Validate
}

func NewUserServiceImpl(ur userrepository.UserRepository, rr rolerepository.RoleRepository,
	urr userrolerepository.UserRoleRepository, db *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		ur,
		rr,
		urr,
		db,
		validate,
	}
}

func (u *UserServiceImpl) FindAll(ctx *gin.Context) (*[]dto.UserResponse, errs.Error) {
	result, err := u.ur.FindAll(ctx, u.db)

	if err != nil {
		return nil, err
	}

	return helper.ToDtoUsersResponse(result), nil
}

func (u *UserServiceImpl) FindOneByUsername(ctx *gin.Context, username string) (*dto.UserResponse, errs.Error) {
	result, err := u.ur.FindByUsername(ctx, u.db, username)

	if err != nil {
		return nil, err
	}

	return helper.ToDtoUserResponse(result), nil
}

func (u *UserServiceImpl) Create(ctx *gin.Context, userDto *dto.CreateNewUserRequest) (*dto.UserResponse, errs.Error) {
	if validationErr := u.validate.Struct(userDto); validationErr != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(validationErr))
	}

	user := entity.User{
		Username: userDto.Username,
		Email:    userDto.Username,
		Password: userDto.Password,
	}

	if err := user.HashPassword(); err != nil {
		log.Printf("[CreateNewUser - Repo] err: %v", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	// check if user already exists
	if findByEmail, err := u.ur.FindByEmail(ctx, u.db, user.Email); err != nil {
		if err.Status() == http.StatusText(http.StatusInternalServerError) {
			return nil, err
		}
	} else if findByEmail != nil {
		return nil, errs.NewConflictError("email has already been taken")
	}

	if findByUsername, err := u.ur.FindByUsername(ctx, u.db, user.Username); err != nil {
		if err.Status() == http.StatusText(http.StatusInternalServerError) {
			return nil, err
		}
	} else if findByUsername != nil {
		return nil, errs.NewConflictError("username has already been taken")
	}
	// end of check if user already exists

	// fetch roles data
	roles, err := u.rr.FindRolesByName(ctx, u.db, userDto.Roles)
	if err != nil {
		return nil, err
	}
	// end of fetch roles data

	// create user
	tx, errTx := u.db.Begin()
	if errTx != nil {
		log.Printf("[CreateNewUser - Repo], err: %s\n", errTx.Error())
		tx.Rollback()
		return nil, errs.NewInternalServerError("something went wrong")
	}

	result, err := u.ur.Create(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	userRoles := []entity.UserRole{}
	for _, role := range *roles {
		userRoles = append(userRoles, entity.UserRole{UserId: result.Id, RoleId: role.Id})
	}

	if err = u.urr.AssignRolesToUser(ctx, tx, userRoles); err != nil {
		tx.Rollback()
		return nil, err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[CreateNewUser - Repo] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}
	// end of create user

	newUser := userrepository.UserRoles{
		User:  *result,
		Roles: *roles,
	}

	return helper.ToDtoUserResponse(&newUser), nil
}
