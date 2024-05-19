package userservice

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/helper"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	accountactivationrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/account_activation_repository"
	rolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/role_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	userrolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_role_repository"
)

type UserServiceImpl struct {
	ur       userrepository.UserRepository
	rr       rolerepository.RoleRepository
	urr      userrolerepository.UserRoleRepository
	aar      accountactivationrepository.AccountActivationRepository
	db       *sql.DB
	validate *validator.Validate
}

func NewUserServiceImpl(ur userrepository.UserRepository, rr rolerepository.RoleRepository,
	urr userrolerepository.UserRoleRepository, aar accountactivationrepository.AccountActivationRepository, db *sql.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		ur,
		rr,
		urr,
		aar,
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
		log.Printf("[CreateNewUser - Service] err: %v", err.Error())
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
		log.Printf("[CreateNewUser - Service], err: %s\n", errTx.Error())
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

	// user account activation
	accountActivation := entity.AccountActivation{
		UserId:         result.Id,
		Token:          helper.GenerateRandomHashString(),
		RequestTime:    time.Now(),
		ExpirationTime: time.Now().Add(config.GetAppConfig().ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION),
	}

	if err := u.aar.Create(ctx, tx, accountActivation); err != nil {
		tx.Rollback()
		return nil, err
	}

	// send activation link via email
	go func() {
		helper.SendTokenEmail(dto.SendTokenEmailRequest{ToEmailAddress: userDto.Email, Subject: "Account Activation", Username: userDto.Username, Token: accountActivation.Token})
	}()

	// end of user account activation

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[CreateNewUser - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}
	// end of create user

	newUser := userrepository.UserRoles{
		User:  *result,
		Roles: *roles,
	}

	return helper.ToDtoUserResponse(&newUser), nil
}

func (u *UserServiceImpl) AssignOrRemoveUserRole(ctx *gin.Context, userDto *dto.AssignRoleToUserRequest) (*dto.UserResponse, errs.Error) {
	if errValidation := u.validate.Struct(userDto); errValidation != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(errValidation))
	}

	// fetch roles data
	rolesResult, err := u.rr.FindRolesByName(ctx, u.db, []entity.UserType{userDto.Role})
	if err != nil {
		return nil, err
	}
	// end of fetch roles data

	roles := *rolesResult

	userRole := entity.UserRole{
		UserId: userDto.UserId,
		RoleId: roles[0].Id,
	}

	tx, errTx := u.db.Begin()
	if errTx != nil {
		log.Printf("[AssignOrRemoveUserRole - Service], err: %s\n", errTx.Error())
		tx.Rollback()
		return nil, errs.NewInternalServerError("something went wrong")
	}

	if err := u.urr.RemoveRoleFromUser(ctx, tx, userRole); err != nil {
		// if when revoking role from user return not found then assign role to user
		if err.Status() == http.StatusText(http.StatusNotFound) {
			if err := u.urr.AssignRolesToUser(ctx, tx, []entity.UserRole{userRole}); err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			tx.Rollback()
			return nil, err
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[AssignOrRemoveUserRole - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	user, err := u.ur.FindById(ctx, u.db, userDto.UserId)
	if err != nil {
		return nil, err
	}

	return helper.ToDtoUserResponse(user), nil
}

func (u *UserServiceImpl) Delete(ctx *gin.Context, userId uuid.UUID) errs.Error {
	if err := u.ur.Delete(ctx, u.db, userId); err != nil {
		return err
	}

	return nil
}
