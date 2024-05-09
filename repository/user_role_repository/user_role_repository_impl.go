package userrolerepository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserRoleRepositoryImpl struct{}

func NewUserRoleRepositoryImpl() UserRoleRepository {
	return &UserRoleRepositoryImpl{}
}

func (u *UserRoleRepositoryImpl) AssignRolesToUser(ctx *gin.Context, tx *sql.Tx, userRole []entity.UserRole) errs.Error {
	sqlQuery := `INSERT INTO users_roles(user_id, role_id) VALUES($1,$2)`
	statement, err := tx.PrepareContext(ctx, sqlQuery)

	if err != nil {
		log.Printf("[AssignRolesToUser - Repo] err: %s", err.Error())
		return errs.NewInternalServerError("something went wrong")
	}

	defer statement.Close()

	for _, data := range userRole {
		_, err = statement.ExecContext(ctx, data.UserId, data.RoleId)

		if err != nil {
			log.Printf("[AssignRolesToUser - Repo] err: %s", err.Error())
			return errs.NewInternalServerError("something went wrong")
		}
	}

	return nil
}

func (u *UserRoleRepositoryImpl) RemoveRoleFromUser(ctx *gin.Context, tx *sql.Tx, userRole entity.UserRole) errs.Error {
	sqlQuery := `DELETE FROM users_roles WHERE user_id=$1 AND role_id=$2 RETURNING id`

	if err := tx.QueryRowContext(ctx, sqlQuery, userRole.UserId, userRole.RoleId).Scan(&userRole.Id); err != nil {
		log.Printf("[RemoveRoleFromUser - Repo], err: %s\n", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewNotFoundError("no user role found")
		}

		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
