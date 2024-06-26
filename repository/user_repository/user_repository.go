package userrepository

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserRepository interface {
	Create(ctx *gin.Context, tx *sql.Tx, user entity.User) (*entity.User, errs.Error)
	FindAll(ctx *gin.Context, db *sql.DB) (*[]UserRoles, errs.Error)
	FindById(ctx *gin.Context, db *sql.DB, id uuid.UUID) (*UserRoles, errs.Error)
	FindByEmail(ctx *gin.Context, db *sql.DB, email string) (*UserRoles, errs.Error)
	FindByUsername(ctx *gin.Context, db *sql.DB, username string) (*UserRoles, errs.Error)
	Delete(ctx *gin.Context, db *sql.DB, id uuid.UUID) errs.Error
	UpdateUserActivation(ctx *gin.Context, tx *sql.Tx, user entity.User) errs.Error
}
