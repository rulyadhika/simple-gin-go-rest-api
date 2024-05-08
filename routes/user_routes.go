package routes

import (
	"github.com/gin-gonic/gin"
	userhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/user_handler"
	authmiddleware "github.com/rulyadhika/simple-gin-go-rest-api/infra/middleware/auth_middleware"
	"github.com/rulyadhika/simple-gin-go-rest-api/model/entity"
)

func NewUserRoutes(r *gin.Engine, handler userhandler.UserHandler, authMiddleware authmiddleware.AuthMiddleware) {
	userRoute := r.Group("/users")
	{
		userRoute.Use(authMiddleware.Authentication())
		userRoute.GET("/", authMiddleware.RoleAuthorization(entity.Role_ADMINISTRATOR), handler.FindAll)
		userRoute.GET("/:username", authMiddleware.RoleAuthorization(entity.Role_ADMINISTRATOR), handler.FindOneByUsername)
		userRoute.POST("/", authMiddleware.RoleAuthorization(entity.Role_ADMINISTRATOR), handler.Create)
		userRoute.PATCH("/:userId/roles/:roleName", authMiddleware.RoleAuthorization(entity.Role_ADMINISTRATOR), handler.AssignReassignRoleToUser)
		userRoute.DELETE("/:userId", authMiddleware.RoleAuthorization(entity.Role_ADMINISTRATOR), handler.Delete)
	}
}
