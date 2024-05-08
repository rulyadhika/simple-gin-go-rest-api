package userrepository

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type UserRepositoryImpl struct{}

func NewUserRepositoryImpl() UserRepository {
	return &UserRepositoryImpl{}
}

func (u *UserRepositoryImpl) Create(ctx *gin.Context, tx *sql.Tx, user entity.User) (*entity.User, errs.Error) {
	sqlQuery := createNewUserQuery

	err := tx.QueryRowContext(ctx, sqlQuery, user.Username, user.Email, user.Password).Scan(&user.Id, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("[CreateUser - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return &user, nil
}

func (u *UserRepositoryImpl) FindAll(ctx *gin.Context, db *sql.DB) (*[]UserRoles, errs.Error) {
	sqlQuery := findAllUserQuery

	rows, err := db.QueryContext(ctx, sqlQuery)

	if err != nil {
		log.Printf("[FindAll - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	usersRoles := []UserRole{}

	for rows.Next() {
		userRole := UserRole{}
		err := rows.Scan(&userRole.User.Id, &userRole.Username, &userRole.Email, &userRole.Password, &userRole.CreatedAt, &userRole.UpdatedAt, &userRole.Role.Id, &userRole.RoleName)

		if err != nil {
			log.Printf("[FindAll - Repo] err: %s", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		usersRoles = append(usersRoles, userRole)
	}

	// if the result is empty
	if len(usersRoles) == 0 {
		return nil, errs.NewNotFoundError("no user data found")
	}

	allUsers := UserRoles{}

	return allUsers.HandleMappingUsersRoles(usersRoles), nil
}

func (u *UserRepositoryImpl) FindByEmail(ctx *gin.Context, db *sql.DB, email string) (*UserRoles, errs.Error) {
	sqlQuery := findOneUserByEmailQuery

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
		return nil, errs.NewNotFoundError("user not found")
	}

	userRoles := UserRoles{}
	userRoles.HandleMappingUserRoles(usersRoles)

	return &userRoles, nil
}

func (u *UserRepositoryImpl) FindByUsername(ctx *gin.Context, db *sql.DB, username string) (*UserRoles, errs.Error) {
	sqlQuery := findOneUserByUsernameQuery

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
		return nil, errs.NewNotFoundError("user not found")
	}

	userRoles := UserRoles{}
	userRoles.HandleMappingUserRoles(usersRoles)

	return &userRoles, nil
}

func (u *UserRepositoryImpl) FindById(ctx *gin.Context, db *sql.DB, id uint32) (*UserRoles, errs.Error) {
	sqlQuery := findOneUserByIdQuery

	rows, err := db.QueryContext(ctx, sqlQuery, id)

	if err != nil {
		log.Printf("[FindById - Repo] err: %s", err.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	defer rows.Close()

	usersRoles := []UserRole{}

	for rows.Next() {
		userRole := UserRole{}
		err := rows.Scan(&userRole.User.Id, &userRole.Username, &userRole.Email, &userRole.Password, &userRole.CreatedAt, &userRole.UpdatedAt, &userRole.Role.Id, &userRole.RoleName)

		if err != nil {
			log.Printf("[FindById - Repo] err: %s", err.Error())
			return nil, errs.NewInternalServerError("something went wrong")
		}

		usersRoles = append(usersRoles, userRole)
	}

	// if the result is empty
	if len(usersRoles) == 0 {
		return nil, errs.NewNotFoundError("user not found")
	}

	userRoles := UserRoles{}
	userRoles.HandleMappingUserRoles(usersRoles)

	return &userRoles, nil
}

func (u *UserRepositoryImpl) Delete(ctx *gin.Context, db *sql.DB, id uint32) errs.Error {
	err := db.QueryRowContext(ctx, deleteUserQuery, id).Scan(&id)

	if err != nil {
		log.Printf("[DeleteUser - Repo] err: %s", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return errs.NewNotFoundError("user not found")
		} else if strings.Contains(err.Error(), "foreign key constraint") {
			return errs.NewConflictError("cannot delete this user because of data constraints. alternatively, you can deactivate the user instead of deleting it")
		}

		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
