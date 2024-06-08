package accountpasswordresetrepository

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

type AccountPasswordResetRepository interface {
	Create(ctx *gin.Context, tx *sql.Tx, account entity.AccountPasswordReset) errs.Error
	FindOneByUserId(ctx *gin.Context, tx *sql.Tx, userId uuid.UUID) (*entity.AccountPasswordReset, errs.Error)
	UpdateRequestTime(ctx *gin.Context, tx *sql.Tx, account entity.AccountPasswordReset) errs.Error
}
