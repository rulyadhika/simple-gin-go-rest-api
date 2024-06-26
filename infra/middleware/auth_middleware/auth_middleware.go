package authmiddleware

import (
	"database/sql"
	"log"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/jwt"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
)

type AuthMiddleware interface {
	Authentication() gin.HandlerFunc
	AuthorizationTicket() gin.HandlerFunc
	RoleAuthorization(rolesAllowed ...entity.UserType) gin.HandlerFunc
}

type AuthMiddlewareImpl struct {
	tr ticketrepository.TicketRepository
	db *sql.DB
}

func NewAuthMiddlewareImpl(tr ticketrepository.TicketRepository, db *sql.DB) AuthMiddleware {
	return &AuthMiddlewareImpl{tr, db}
}

func (a *AuthMiddlewareImpl) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")
		hasPrefixBearer := strings.HasPrefix(authorizationHeader, "Bearer")

		unauthorizedError := errs.NewUnauthorizedError("invalid token")

		if !hasPrefixBearer {
			ctx.AbortWithStatusJSON(unauthorizedError.StatusCode(), unauthorizedError)
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) != 2 {
			ctx.AbortWithStatusJSON(unauthorizedError.StatusCode(), unauthorizedError)
			return
		}

		token := bearerToken[1]

		userData := jwt.NewJWTTokenParser()
		if err := userData.ParseToken(token, config.GetAppConfig().ACCESS_TOKEN_SECRET); err != nil {
			ctx.AbortWithStatusJSON(unauthorizedError.StatusCode(), unauthorizedError)
			return
		}

		ctx.Set("userData", userData)
		ctx.Next()
	}
}

func (a *AuthMiddlewareImpl) AuthorizationTicket() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, ok := ctx.MustGet("userData").(*jwt.JWTPayload)

		if !ok {
			log.Printf("[AuthorizationTicket - Middleware] err: %s\n", "failed type casting to '*jwt.JWTPayload'")
			internalServerErr := errs.NewInternalServerError("something went wrong")
			ctx.AbortWithStatusJSON(internalServerErr.StatusCode(), internalServerErr)
			return
		}

		ticketId, errParseUUID := uuid.Parse(ctx.Param("ticketId"))
		if errParseUUID != nil {
			log.Printf("[AuthorizationTicket - middleware], err: %s\n", errParseUUID.Error())
			unprocessableEntityError := errs.NewUnprocessableEntityError("param ticketId must be a valid id")
			ctx.AbortWithStatusJSON(unprocessableEntityError.StatusCode(), unprocessableEntityError)
			return
		}

		result, err := a.tr.FindOneByTicketId(ctx, a.db, ticketId)

		if err != nil {
			ctx.AbortWithStatusJSON(err.StatusCode(), err)
			return
		}

		// if user's roles is only client
		if slices.Contains(userData.Roles, "client") && len(userData.Roles) == 1 {
			// then check whether the data belongs to that user
			if result.CreatedBy.Email.String != userData.Email {
				forbiddenErr := errs.NewForbiddenError("you are not authorized to access/modify this data")
				ctx.AbortWithStatusJSON(forbiddenErr.StatusCode(), forbiddenErr)
				return
			}
		}

		ctx.Next()
	}
}

func (a *AuthMiddlewareImpl) RoleAuthorization(rolesAllowed ...entity.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData, ok := ctx.MustGet("userData").(*jwt.JWTPayload)

		_ = userData

		if !ok {
			log.Printf("[RoleAuthorization - Middleware] err: %s\n", "failed type casting to '*jwt.JWTPayload'")
			internalServerErr := errs.NewInternalServerError("something went wrong")
			ctx.AbortWithStatusJSON(internalServerErr.StatusCode(), internalServerErr)
			return
		}

		userIsEligible := []bool{}

		for _, r := range userData.Roles {
			userIsEligible = append(userIsEligible, slices.Contains(rolesAllowed, r))
		}

		if !slices.Contains(userIsEligible, true) {
			forbiddenErr := errs.NewForbiddenError("you are not authorized to access/modify this data")
			ctx.AbortWithStatusJSON(forbiddenErr.StatusCode(), forbiddenErr)
			return
		}

		ctx.Next()
	}
}
