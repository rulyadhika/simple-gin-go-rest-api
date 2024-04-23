package userrolerepository

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserRoleRepository interface {
	AssignRolesToUser(ctx *gin.Context, tx *sql.Tx, userRole []entity.UserRole) errs.Error
	RevokeRoleFromUser(ctx *gin.Context, tx *sql.Tx, userRole entity.UserRole) errs.Error
}
