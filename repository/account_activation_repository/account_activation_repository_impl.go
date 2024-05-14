package accountactivationrepository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type AccountActivationRepositoryImpl struct{}

func NewAccountActivationRepositoryImpl() AccountActivationRepository {
	return &AccountActivationRepositoryImpl{}
}

func (a *AccountActivationRepositoryImpl) Create(ctx *gin.Context, tx *sql.Tx, account entity.AccountActivation) errs.Error {
	_, err := tx.ExecContext(ctx, createAccountActivationDataQuery, account.UserId, account.Token, account.RequestTime, account.ExpirationTime)

	if err != nil {
		log.Printf("[CreateAccountActivationData - Repo] err: %s", err.Error())

		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}

func (a *AccountActivationRepositoryImpl) FindOne(ctx *gin.Context, tx *sql.Tx, token string) (*entity.AccountActivation, errs.Error) {
	account := new(entity.AccountActivation)

	err := tx.QueryRowContext(ctx, findOneAccountActivationDataQuery, token).Scan(&account.UserId, &account.Token, &account.RequestTime, &account.ExpirationTime)
	if err != nil {
		log.Printf("[FindOneAccountActivationData - Repo] err: %s", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("account activation data not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}

	return account, nil

}

func (a *AccountActivationRepositoryImpl) Delete(ctx *gin.Context, tx *sql.Tx, token string) errs.Error {
	_, err := tx.ExecContext(ctx, deleteAccountActivationDataQuery, token)

	if err != nil {
		log.Printf("[DeleteAccountActivationData - Repo] err: %s", err.Error())

		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
