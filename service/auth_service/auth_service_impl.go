package authservice

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/helper"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/jwt"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	accountactivationrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/account_activation_repository"
	rolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/role_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	userrolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_role_repository"
)

type AuthServiceImpl struct {
	UserRepository              userrepository.UserRepository
	DB                          *sql.DB
	Validate                    *validator.Validate
	UserRoleRepository          userrolerepository.UserRoleRepository
	RoleRepository              rolerepository.RoleRepository
	AccountActivationRepository accountactivationrepository.AccountActivationRepository
}

func NewAuthServiceImpl(userRepository userrepository.UserRepository, userRoleRepository userrolerepository.UserRoleRepository, roleRepository rolerepository.RoleRepository, accountActivationRepository accountactivationrepository.AccountActivationRepository, db *sql.DB, validator *validator.Validate) AuthService {
	return &AuthServiceImpl{
		UserRepository:              userRepository,
		DB:                          db,
		Validate:                    validator,
		UserRoleRepository:          userRoleRepository,
		RoleRepository:              roleRepository,
		AccountActivationRepository: accountActivationRepository,
	}
}

func (a *AuthServiceImpl) Register(ctx *gin.Context, userDto *dto.RegisterUserRequest) (*dto.RegisterUserResponse, errs.Error) {
	if errValidate := a.Validate.Struct(userDto); errValidate != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(errValidate))
	}

	user := entity.User{
		Username: userDto.Username,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	if err := user.HashPassword(); err != nil {
		log.Printf("[Register - Service] err: %v", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	// check if user already exists
	if findByEmail, err := a.UserRepository.FindByEmail(ctx, a.DB, user.Email); err != nil {
		if err.Status() == http.StatusText(http.StatusInternalServerError) {
			return nil, err
		}
	} else if findByEmail != nil {
		return nil, errs.NewConflictError("email has already been taken")
	}

	if findByUsername, err := a.UserRepository.FindByUsername(ctx, a.DB, user.Username); err != nil {
		if err.Status() == http.StatusText(http.StatusInternalServerError) {
			return nil, err
		}
	} else if findByUsername != nil {
		return nil, errs.NewConflictError("username has already been taken")
	}
	// end of check if user already exists

	// fetch roles data
	roles, err := a.RoleRepository.FindRolesByName(ctx, a.DB, []entity.UserType{entity.Role_CLIENT})
	if err != nil {
		return nil, err
	}
	// end of fetch roles data

	// create new user with client roles
	tx, txErr := a.DB.Begin()
	if txErr != nil {
		log.Printf("[Register - Service] err: %v", txErr.Error())
		tx.Rollback()
		return nil, errs.NewInternalServerError("something went wrong")
	}

	// create user
	result, err := a.UserRepository.Create(ctx, tx, user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// assign user roles
	userRoles := []entity.UserRole{}
	for _, role := range *roles {
		userRoles = append(userRoles, entity.UserRole{UserId: result.Id, RoleId: role.Id})
	}

	if err = a.UserRoleRepository.AssignRolesToUser(ctx, tx, userRoles); err != nil {
		tx.Rollback()
		return nil, err
	}
	// end of assign user roles

	// user account activation
	accountActivation := entity.AccountActivation{
		UserId:         result.Id,
		Token:          helper.GenerateRandomHashString(),
		RequestTime:    time.Now(),
		ExpirationTime: time.Now().Add(config.GetAppConfig().ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION),
	}

	if err := a.AccountActivationRepository.Create(ctx, tx, accountActivation); err != nil {
		tx.Rollback()
		return nil, err
	}

	// send activation link via email
	go func() {
		helper.SendTokenEmail(dto.SendTokenEmailRequest{ToEmailAddress: result.Email, Subject: "Account Activation", Username: result.Username, Token: accountActivation.Token})
	}()

	// end of user account activation

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[Register - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}
	// end of create user

	return &dto.RegisterUserResponse{Username: result.Username, Email: result.Email}, nil
}

func (a *AuthServiceImpl) Login(ctx *gin.Context, userDto *dto.LoginUserRequest) (*dto.LoginUserResponse, errs.Error) {
	if errValidate := a.Validate.Struct(userDto); errValidate != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(errValidate))
	}

	var userFound bool
	var user userrepository.UserRoles

	// check if user exists with findByEmail
	if findByEmail, err := a.UserRepository.FindByEmail(ctx, a.DB, userDto.UsernameOrEmail); err != nil {
		if err.Status() == http.StatusText(http.StatusInternalServerError) {
			return nil, err
		}
	} else if findByEmail != nil {
		userFound = true
		user = *findByEmail
	}

	// if user still not found with findByEmail then check if user exists with findByUsername
	if !userFound {
		if findByUsername, err := a.UserRepository.FindByUsername(ctx, a.DB, userDto.UsernameOrEmail); err != nil {
			if err.Status() == http.StatusText(http.StatusInternalServerError) {
				return nil, err
			}
		} else if findByUsername != nil {
			userFound = true
			user = *findByUsername
		}
	}

	// if user still not found then return error
	if !userFound {
		return nil, errs.NewBadRequestError("invalid login credential")
	}

	// check if password is valid
	if passwordIsValid := user.ValidatePassword(userDto.Password); !passwordIsValid {
		return nil, errs.NewBadRequestError("invalid login credential")
	}

	// check if account is activated or not
	if !user.ActivatedAt.Valid {
		return nil, errs.NewForbiddenError("Your account has not been activated. Please check your email for the activation link.")
	}

	appConfig := config.GetAppConfig()

	accessToken, errGenerateToken := jwt.NewJWTToken(&user).GenerateToken(appConfig.ACCESS_TOKEN_SECRET, time.Now().Add(appConfig.ACCESS_TOKEN_EXPIRATION_DURATION))
	if errGenerateToken != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	refreshToken, errGenerateToken := jwt.NewJWTToken(&user).GenerateToken(appConfig.REFRESH_TOKEN_SECRET, time.Now().Add(appConfig.REFRESH_TOKEN_EXPIRATION_DURATION))
	if errGenerateToken != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	// set refresh token to cookies
	ctx.SetCookie("refresh-token", refreshToken.(string), 24*60*60, "/", "", false, true) //max age: a day

	return &dto.LoginUserResponse{Token: accessToken.(string)}, nil
}

func (a *AuthServiceImpl) RefreshToken(ctx *gin.Context, userDto *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, errs.Error) {
	userData := jwt.NewJWTTokenParser()

	if err := userData.ParseToken(userDto.Token, config.GetAppConfig().REFRESH_TOKEN_SECRET); err != nil {
		return nil, errs.NewForbiddenError("token is invalid or expired")
	}

	user, err := a.UserRepository.FindById(ctx, a.DB, userData.Id)

	if err != nil {
		log.Printf("[RefreshToken - Service] err: user with id: %v not found\n", userData.Id)
		return nil, errs.NewInternalServerError("something went wrong")
	}

	accessToken, errGenerateToken := jwt.NewJWTToken(user).GenerateToken(config.GetAppConfig().ACCESS_TOKEN_SECRET, time.Now().Add(config.GetAppConfig().REFRESH_TOKEN_EXPIRATION_DURATION))
	if errGenerateToken != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &dto.RefreshTokenResponse{Token: accessToken.(string)}, nil
}
