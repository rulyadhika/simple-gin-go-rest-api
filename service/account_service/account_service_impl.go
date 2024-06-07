package accountservice

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/helper"
	validationformatter "github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation/validation_formatter"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/dto"
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
		Id: accountActivation.UserId,
		ActivatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
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

func (a *accountServiceImpl) ResendActivationToken(ctx *gin.Context, resendTokenDto dto.ResendActivationTokenRequest) (*dto.ResendActivationTokenRespone, errs.Error) {
	if validationErr := a.validate.Struct(resendTokenDto); validationErr != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(validationErr))
	}

	tx, errTx := a.db.Begin()
	if errTx != nil {
		log.Printf("[ResendToken - Service] err: %s", errTx.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	user, err := a.ur.FindByEmail(ctx, a.db, resendTokenDto.Email)

	if err != nil {
		return nil, err
	}

	if user.ActivatedAt.Valid {
		return nil, errs.NewConflictError("the account has already been activated.")
	}

	accountActivationData, err := a.aar.FindOneByUserId(ctx, tx, user.Id)
	if err != nil {
		return nil, err
	}

	resendActivationTokenRespone := new(dto.ResendActivationTokenRespone)

	// check if token is expired or not
	if time.Now().Before(accountActivationData.ExpirationTime) {
		// if not yet expired
		if !time.Now().After(accountActivationData.NextRequestAvailableAt) {
			// if next request not available
			timeElapsed := time.Until(accountActivationData.NextRequestAvailableAt)
			seconds := int(timeElapsed.Seconds()) % 60

			return nil, errs.NewConflictError(fmt.Sprintf("Cannot resend activation link. Please wait %ds then try again", seconds))
		}

		// update request time and next available request at data
		updatedAccountActivationData := entity.AccountActivation{RequestTime: time.Now(), NextRequestAvailableAt: time.Now().Add(1 * time.Minute), Token: accountActivationData.Token}

		if err := a.aar.UpdateRequestTime(ctx, tx, updatedAccountActivationData); err != nil {
			tx.Rollback()
			return nil, err
		}

		resendActivationTokenRespone.RequestTime = updatedAccountActivationData.RequestTime
		resendActivationTokenRespone.NextRequestAvailableAt = updatedAccountActivationData.NextRequestAvailableAt

		// resend activation link via email
		go func() {
			helper.SendTokenEmail(dto.SendTokenEmailRequest{
				ToEmailAddress: user.Email,
				Subject:        "Account Activation",
				Username:       user.Username,
				Token:          accountActivationData.Token,
			})
		}()
	} else {
		// else then generate new token and send to email
		accountActivation := entity.AccountActivation{
			UserId:                 user.Id,
			Token:                  helper.GenerateRandomHashString(),
			ExpirationTime:         time.Now().Add(config.GetAppConfig().ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION),
			NextRequestAvailableAt: time.Now().Add(1 * time.Minute),
		}

		if err := a.aar.Create(ctx, tx, accountActivation); err != nil {
			tx.Rollback()
			return nil, err
		}

		resendActivationTokenRespone.RequestTime = time.Now()
		resendActivationTokenRespone.NextRequestAvailableAt = accountActivation.NextRequestAvailableAt

		// send activation link via email
		go func() {
			helper.SendTokenEmail(dto.SendTokenEmailRequest{
				ToEmailAddress: user.Email,
				Subject:        "Account Activation",
				Username:       user.Username,
				Token:          accountActivationData.Token,
			})
		}()
	}
	// end of check if token is expired or not

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[ResendToken - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return resendActivationTokenRespone, nil
}
