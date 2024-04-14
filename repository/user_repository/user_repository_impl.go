package userrepository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserRepositoryImpl struct{}

func NewUserRepositoryImpl() UserRepository {
	return &UserRepositoryImpl{}
}

func (a *UserRepositoryImpl) Create(ctx *gin.Context, tx *sql.Tx, user entity.User) (*entity.User, errs.Error) {
	sqlQuery := `INSERT INTO users(username, email, password) VALUES($1,$2,$3) RETURNING id`

	err := tx.QueryRowContext(ctx, sqlQuery, user.Username, user.Email, user.Password).Scan(&user.Id)

	if err != nil {
		log.Printf("[CreateUser - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &user, nil
}

func (a *UserRepositoryImpl) AssignRolesToUser(ctx *gin.Context, tx *sql.Tx, userRole []entity.UserRole) errs.Error {
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

func (a *UserRepositoryImpl) FindByEmail(ctx *gin.Context, db *sql.DB, email string) (*UserRoles, errs.Error) {
	sqlQuery := `SELECT users.id, username, email, password, users.created_at, users.updated_at, roles.id, roles.role_name
	FROM users JOIN users_roles ON users.id=users_roles.user_id
	JOIN roles on users_roles.role_id=roles.id WHERE email=$1`

	rows, err := db.QueryContext(ctx, sqlQuery, email)

	if err != nil {
		log.Printf("[FindByEmail - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	usersRoles := []UserRole{}

	for rows.Next() {
		userRole := UserRole{}
		err := rows.Scan(&userRole.User.Id, &userRole.Username, &userRole.Email, &userRole.Password, &userRole.CreatedAt, &userRole.UpdatedAt, &userRole.Role.Id, &userRole.RoleName)

		if err != nil {
			log.Printf("[FindByEmail - Repo] err: %s", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		usersRoles = append(usersRoles, userRole)
	}

	// if the result is empty
	if len(usersRoles) == 0 {
		return nil, nil
	}

	userRoles := UserRoles{}
	userRoles.HandleMappingUserRoles(usersRoles)

	return &userRoles, nil
}

func (a *UserRepositoryImpl) FindByUsername(ctx *gin.Context, db *sql.DB, username string) (*UserRoles, errs.Error) {
	sqlQuery := `SELECT users.id, username, email, password, users.created_at, users.updated_at, roles.id, roles.role_name
	FROM users JOIN users_roles ON users.id=users_roles.user_id
	JOIN roles on users_roles.role_id=roles.id WHERE username=$1`

	rows, err := db.QueryContext(ctx, sqlQuery, username)

	if err != nil {
		log.Printf("[FindByUsername - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	usersRoles := []UserRole{}

	for rows.Next() {
		userRole := UserRole{}
		err := rows.Scan(&userRole.User.Id, &userRole.Username, &userRole.Email, &userRole.Password, &userRole.CreatedAt, &userRole.UpdatedAt, &userRole.Role.Id, &userRole.RoleName)

		if err != nil {
			log.Printf("[FindByUsername - Repo] err: %s", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		usersRoles = append(usersRoles, userRole)
	}

	// if the result is empty
	if len(usersRoles) == 0 {
		return nil, nil
	}

	userRoles := UserRoles{}
	userRoles.HandleMappingUserRoles(usersRoles)

	fmt.Println(userRoles)

	return &userRoles, nil
}
