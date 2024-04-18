package authmiddleware

import (
	"log"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/errs"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/jwt"
)

type AuthMiddleware interface {
	Authentication() gin.HandlerFunc
	RoleAuthorization(rolesAllowed []string) gin.HandlerFunc
}

type AuthMiddlewareImpl struct{}

func NewAuthMiddlewareImpl() AuthMiddleware {
	return &AuthMiddlewareImpl{}
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

func (a *AuthMiddlewareImpl) RoleAuthorization(rolesAllowed []string) gin.HandlerFunc {
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
