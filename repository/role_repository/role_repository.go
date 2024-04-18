package rolerepository

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type RoleRepository interface {
	FindRolesByName(ctx *gin.Context, db *sql.DB, rolesList []string) (*[]entity.Role, errs.Error)
}
