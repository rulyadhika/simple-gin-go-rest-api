package main

import (
	"github.com/gin-gonic/gin"
	authhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/auth_handler"
	tickethandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/ticket_handler"
	userhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/user_handler"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/db"
	authmiddleware "github.com/rulyadhika/simple-gin-go-rest-api/infra/middleware/auth_middleware"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation"
	rolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/role_repository"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	userrolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_role_repository"
	"github.com/rulyadhika/simple-gin-go-rest-api/routes"
	authservice "github.com/rulyadhika/simple-gin-go-rest-api/service/auth_service"
	ticketservice "github.com/rulyadhika/simple-gin-go-rest-api/service/ticket_service"
	userservice "github.com/rulyadhika/simple-gin-go-rest-api/service/user_service"
)

func main() {
	appConfig := config.GetAppConfig()
	db := db.InitDB()
	validator := validation.NewValidator()
	app := gin.Default()

	userRepository := userrepository.NewUserRepositoryImpl()
	userRoleRepository := userrolerepository.NewUserRoleRepositoryImpl()
	roleRepository := rolerepository.NewRoleRepositoryImpl()
	ticketRepository := ticketrepository.NewTicketRepositoryImpl()

	authService := authservice.NewAuthServiceImpl(userRepository, userRoleRepository, roleRepository, db, validator)
	ticketService := ticketservice.NewTicketServiceImpl(ticketRepository, userRepository, db, validator)
	userService := userservice.NewUserServiceImpl(userRepository, db, validator)

	authHandler := authhandler.NewAuthHandlerImpl(authService)
	ticketHandler := tickethandler.NewTicketHandlerImpl(ticketService)
	userHander := userhandler.NewUserHandlerImpl(userService)

	// middlewares
	authMiddleware := authmiddleware.NewAuthMiddlewareImpl(ticketRepository, db)
	// end of middlewares

	// routes
	routes.NewAuthRoutes(app, authHandler)
	routes.NewTicketRoutes(app, ticketHandler, authMiddleware)
	routes.NewUserRoutes(app, userHander, authMiddleware)
	// end of routes

	app.Run(":" + appConfig.APP_PORT)
}
