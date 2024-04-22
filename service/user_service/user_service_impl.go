package userservice

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/helper"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

type UserServiceImpl struct {
	ur        userrepository.UserRepository
	db        *sql.DB
	validator *validator.Validate
}

func NewUserServiceImpl(ur userrepository.UserRepository, db *sql.DB, validator *validator.Validate) UserService {
	return &UserServiceImpl{
		ur,
		db,
		validator,
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
