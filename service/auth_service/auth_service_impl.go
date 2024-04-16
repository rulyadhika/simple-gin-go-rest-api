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
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/jwt"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
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

	appConfig := config.GetAppConfig()

	accessToken, errGenerateToken := jwt.NewJWTToken(&user).GenerateToken(appConfig.ACCESS_TOKEN_SECRET, time.Now().Add(1*time.Minute))
	if errGenerateToken != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	refreshToken, errGenerateToken := jwt.NewJWTToken(&user).GenerateToken(appConfig.REFRESH_TOKEN_SECRET, time.Now().Add(24*time.Hour))
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

	accessToken, errGenerateToken := jwt.NewJWTToken(user).GenerateToken(config.GetAppConfig().ACCESS_TOKEN_SECRET, time.Now().Add(1*time.Minute))
	if errGenerateToken != nil {
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &dto.RefreshTokenResponse{Token: accessToken.(string)}, nil
}
