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
	accountpasswordresetrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/account_password_reset_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
)

type accountServiceImpl struct {
	validate *validator.Validate
	db       *sql.DB
	aar      accountactivationrepository.AccountActivationRepository
	ur       userrepository.UserRepository
	aprr     accountpasswordresetrepository.AccountPasswordResetRepository
}

func NewAccountServiceImpl(aar accountactivationrepository.AccountActivationRepository, ur userrepository.UserRepository, aprr accountpasswordresetrepository.AccountPasswordResetRepository, db *sql.DB, validate *validator.Validate) AccountService {
	return &accountServiceImpl{
		validate, db, aar, ur, aprr,
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

func (a *accountServiceImpl) ResendActivationToken(ctx *gin.Context, resendTokenDto *dto.ResendActivationTokenRequest) (*dto.ResendActivationTokenRespone, errs.Error) {
	if validationErr := a.validate.Struct(resendTokenDto); validationErr != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(validationErr))
	}

	tx, errTx := a.db.Begin()
	if errTx != nil {
		log.Printf("[ResendActivationToken - Service] err: %s", errTx.Error())
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
	var accountActivation entity.AccountActivation

	// check if token is expired or not
	if time.Now().Before(accountActivationData.ExpirationTime) {
		// if not yet expired
		if !time.Now().After(accountActivationData.NextRequestAvailableAt) {
			// if next request not available
			timeElapsed := time.Until(accountActivationData.NextRequestAvailableAt)
			seconds := int(timeElapsed.Seconds()) % 60

			return nil, errs.NewConflictError(fmt.Sprintf("Cannot resend activation link. Please wait %ds then try again", seconds))
		}

		accountActivation = *accountActivationData

		// if expired then update request time and next available request at data
		updatedAccountActivationData := entity.AccountActivation{RequestTime: time.Now(), NextRequestAvailableAt: time.Now().Add(1 * time.Minute), Token: accountActivationData.Token}

		if err := a.aar.UpdateRequestTime(ctx, tx, updatedAccountActivationData); err != nil {
			tx.Rollback()
			return nil, err
		}

		resendActivationTokenRespone.RequestTime = updatedAccountActivationData.RequestTime
		resendActivationTokenRespone.NextRequestAvailableAt = updatedAccountActivationData.NextRequestAvailableAt
	} else {
		// else then generate new token and send to email
		accountActivationData := entity.AccountActivation{
			UserId:                 user.Id,
			Token:                  helper.GenerateRandomHashString(),
			ExpirationTime:         time.Now().Add(config.GetAppConfig().ACCOUNT_ACTIVATION_TOKEN_EXPIRATION_DURATION),
			NextRequestAvailableAt: time.Now().Add(1 * time.Minute),
		}

		if err := a.aar.Create(ctx, tx, accountActivationData); err != nil {
			tx.Rollback()
			return nil, err
		}

		accountActivation = accountActivationData

		resendActivationTokenRespone.RequestTime = time.Now()
		resendActivationTokenRespone.NextRequestAvailableAt = accountActivationData.NextRequestAvailableAt
	}
	// end of check if token is expired or not

	// commit tx
	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[ResendActivationToken - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	// send activation link via email
	go func() {
		helper.SendActivationTokenEmail(dto.SendActivationTokenEmailRequest{
			ToEmailAddress: user.Email,
			Subject:        "Account Activation",
			Username:       user.Username,
			Token:          accountActivation.Token,
		})
	}()

	return resendActivationTokenRespone, nil
}

func (a *accountServiceImpl) ForgotPassword(ctx *gin.Context, forgotPasswordDto *dto.ForgotPasswordRequest) (*dto.ForgotPasswordRespone, errs.Error) {
	if validationErr := a.validate.Struct(forgotPasswordDto); validationErr != nil {
		return nil, errs.NewBadRequestError(validationformatter.FormatValidationError(validationErr))
	}

	tx, errTx := a.db.Begin()
	if errTx != nil {
		log.Printf("[ForgotPassword - Service] err: %s", errTx.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	// get user data
	user, err := a.ur.FindByEmail(ctx, a.db, forgotPasswordDto.Email)

	if err != nil {
		return nil, err
	}

	// check whether user is activated or not
	if !user.ActivatedAt.Valid {
		return nil, errs.NewConflictError("Unable to send a password reset request. Your account is not activated.")
	}

	// check if user has been requested password reset
	passwordResetData, err := a.aprr.FindOneByUserId(ctx, tx, user.Id)
	generateNewPasswordResetData := false

	if err != nil {
		if err.Status() == http.StatusText(http.StatusInternalServerError) {
			return nil, err
		}

		// if err is returning not found data then need to generate password reset data
		generateNewPasswordResetData = true
	} else {
		// if password reset data is exists then check if the data is expired or not for resend token feature
		if time.Now().Before(passwordResetData.ExpirationTime) {
			// if not yet expired then check if next request is available
			if !time.Now().After(passwordResetData.NextRequestAvailableAt) {
				// if next request not available
				timeElapsed := time.Until(passwordResetData.NextRequestAvailableAt)
				seconds := int(timeElapsed.Seconds()) % 60

				return nil, errs.NewConflictError(fmt.Sprintf("Cannot resend password reset token. Please wait %ds then try again", seconds))
			}

			// if next request available no need to generate new password reset data
			generateNewPasswordResetData = false
		} else {
			// if expired then need to generate new password reset data
			generateNewPasswordResetData = true
		}
	}

	forgotPasswordResponse := new(dto.ForgotPasswordRespone)
	var accountPasswordReset entity.AccountPasswordReset

	if generateNewPasswordResetData {
		accountPasswordReset = entity.AccountPasswordReset{
			UserId:                 user.Id,
			Token:                  helper.GenerateFixedLengthRandomNumber(int(config.GetAppConfig().PASSWORD_RESET_TOKEN_LENGTH)),
			ExpirationTime:         time.Now().Add(config.GetAppConfig().PASSWORD_RESET_TOKEN_EXPIRATION_DURATION),
			NextRequestAvailableAt: time.Now().Add(1 * time.Minute),
		}

		if err := a.aprr.Create(ctx, tx, accountPasswordReset); err != nil {
			return nil, err
		}

		forgotPasswordResponse.RequestTime = time.Now()
		forgotPasswordResponse.NextRequestAvailableAt = accountPasswordReset.NextRequestAvailableAt
	} else {
		accountPasswordReset = *passwordResetData

		// update request time and next available request at data
		updatedPasswordResetData := entity.AccountPasswordReset{RequestTime: time.Now(), NextRequestAvailableAt: time.Now().Add(1 * time.Minute), Token: accountPasswordReset.Token}

		forgotPasswordResponse.RequestTime = updatedPasswordResetData.RequestTime
		forgotPasswordResponse.NextRequestAvailableAt = updatedPasswordResetData.NextRequestAvailableAt

		if err := a.aprr.UpdateRequestTime(ctx, tx, updatedPasswordResetData); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// send token to mail
	go func() {

		helper.SendPasswordResetTokenEmail(dto.SendResetPasswordTokenEmailRequest{
			ToEmailAddress: user.Email,
			Username:       user.Username,
			Subject:        "Password Reset Instructions",
			Token:          accountPasswordReset.Token,
			ExpiredAt:      accountPasswordReset.ExpirationTime.Format(time.RFC1123),
		})
	}()

	if commitErr := tx.Commit(); commitErr != nil {
		log.Printf("[ForgotPassword - Service] err: %v", commitErr.Error())
		return nil, errs.NewInternalServerError("something went wrong")
	}

	return forgotPasswordResponse, nil
}
