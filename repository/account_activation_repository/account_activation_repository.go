package accountactivationrepository

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type AccountActivationRepository interface {
	Create(ctx *gin.Context, tx *sql.Tx, account entity.AccountActivation) errs.Error
	FindOne(ctx *gin.Context, tx *sql.Tx, token string) (*entity.AccountActivation, errs.Error)
	Delete(ctx *gin.Context, tx *sql.Tx, token string) errs.Error
}
