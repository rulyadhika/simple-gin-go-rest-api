package accountpasswordresetrepository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type accountPasswordResetRepositoryImpl struct{}

func NewAccountPasswordResetRepositoryImpl() AccountPasswordResetRepository {
	return &accountPasswordResetRepositoryImpl{}
}

func (a *accountPasswordResetRepositoryImpl) Create(ctx *gin.Context, tx *sql.Tx, account entity.AccountPasswordReset) errs.Error {
	_, err := tx.ExecContext(ctx, createAccountPasswordResetDataQuery, account.UserId, account.Token, account.ExpirationTime, account.NextRequestAvailableAt)

	if err != nil {
		log.Printf("[CreateAccountResetPasswordData - Repo] err: %s", err.Error())

		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}

func (a *accountPasswordResetRepositoryImpl) FindOneByUserId(ctx *gin.Context, tx *sql.Tx, userId uuid.UUID) (*entity.AccountPasswordReset, errs.Error) {
	data := new(entity.AccountPasswordReset)

	err := tx.QueryRowContext(ctx, findOneAccountPasswordResetDataByUserIdQuery, userId).Scan(&data.UserId, &data.Token, &data.RequestTime, &data.ExpirationTime, &data.NextRequestAvailableAt)

	if err != nil {
		log.Printf("[FindOneAccountPasswordResetByUserId - Repo] err: %s", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundError("password reset data not found")
		}

		return nil, errs.NewInternalServerError("something went wrong")
	}

	return data, nil
}

func (a *accountPasswordResetRepositoryImpl) UpdateRequestTime(ctx *gin.Context, tx *sql.Tx, account entity.AccountPasswordReset) errs.Error {
	_, err := tx.ExecContext(ctx, updateRequestTimeAccountPasswordResetDataQuery, account.RequestTime, account.NextRequestAvailableAt, account.Token)

	if err != nil {
		log.Printf("[UpdateRequestTomeAccountPasswordResetData - Repo] err: %s", err.Error())
		return errs.NewInternalServerError("something went wrong")
	}

	return nil
}
