package accountservice

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	accountactivationrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/account_activation_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

type accountServiceImpl struct {
	validate *validator.Validate
	db       *sql.DB
	aar      accountactivationrepository.AccountActivationRepository
	ur       userrepository.UserRepository
}

func NewAccountServiceImpl(aar accountactivationrepository.AccountActivationRepository, ur userrepository.UserRepository, db *sql.DB, validator *validator.Validate) AccountService {
	return &accountServiceImpl{
		validate: validator,
		db:       db,
		aar:      aar,
		ur:       ur,
	}
}

func (a *accountServiceImpl) Activation(ctx *gin.Context, token string) errs.Error {
	tx, errTx := a.db.Begin()

	if errTx != nil {
		log.Printf("[AccountActivation - Service] err: %s", errTx.Error())

		return errs.NewInternalServerError("something went wrong")
	}

	accountActivation, err := a.aar.FindOne(ctx, tx, token)

	if err != nil {
		tx.Rollback()

		if err.Status() != http.StatusText(http.StatusInternalServerError) {
			return errs.NewConflictError("token is invalid or expired")
		}

		return err
	}

	// check if token is expired or not
	if !time.Now().Before(accountActivation.ExpirationTime) {
		tx.Rollback()
		return errs.NewConflictError("token is invalid or expired")
	}
	// end of check if token is expired or not

	// update user activation timestamps table
	user := entity.User{
		Id:          accountActivation.UserId,
		ActivatedAt: time.Now(),
	}

	if err := a.ur.UpdateUserActivation(ctx, tx, user); err != nil {
		tx.Rollback()
		return err
	}
	// end of update user activation timestamps table

	// delete account activation data
	if err := a.aar.Delete(ctx, tx, token); err != nil {
		tx.Rollback()
		return err
	}
	// end of delete account activation data

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[AccountActivation - Service] err: %v", commitErr.Error())
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
