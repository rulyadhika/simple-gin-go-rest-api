package rolerepository

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type RoleRepositoryImpl struct{}

func NewRoleRepositoryImpl() RoleRepository {
	return &RoleRepositoryImpl{}
}

func (r *RoleRepositoryImpl) FindRolesByName(ctx *gin.Context, db *sql.DB, rolesList []entity.UserType) (*[]entity.Role, errs.Error) {
	sqlQuery := "SELECT id, role_name FROM roles WHERE role_name = ANY($1)"

	rows, err := db.QueryContext(ctx, sqlQuery, pq.Array(rolesList))

	if err != nil {
		log.Printf("[FindRolesByName - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	roles := []entity.Role{}

	for rows.Next() {
		role := entity.Role{}
		err := rows.Scan(&role.Id, &role.RoleName)

		if err != nil {
			log.Printf("[FindByEmail - Repo] err: %s", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		roles = append(roles, role)
	}

	// if the result is empty
	if len(roles) == 0 {
		return nil, errs.NewNotFoundError("user not found")
	}

	return &roles, nil
}
