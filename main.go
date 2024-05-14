package main

import (
	"github.com/gin-gonic/gin"
	accounthandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/account_handler"
	authhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/auth_handler"
	tickethandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/ticket_handler"
	userhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/user_handler"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/db"
	authmiddleware "github.com/rulyadhika/simple-gin-go-rest-api/infra/middleware/auth_middleware"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation"
	accountactivationrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/account_activation_repository"
	rolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/role_repository"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	userrolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_role_repository"
	"github.com/rulyadhika/simple-gin-go-rest-api/routes"
	accountservice "github.com/rulyadhika/simple-gin-go-rest-api/service/account_service"
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
	accountActivationRepository := accountactivationrepository.NewAccountActivationRepositoryImpl()
	userRoleRepository := userrolerepository.NewUserRoleRepositoryImpl()
	roleRepository := rolerepository.NewRoleRepositoryImpl()
	ticketRepository := ticketrepository.NewTicketRepositoryImpl()

	authService := authservice.NewAuthServiceImpl(userRepository, userRoleRepository, roleRepository, accountActivationRepository, db, validator)
	ticketService := ticketservice.NewTicketServiceImpl(ticketRepository, userRepository, db, validator)
	userService := userservice.NewUserServiceImpl(userRepository, roleRepository, userRoleRepository, db, validator)
	accountService := accountservice.NewAccountServiceImpl(accountActivationRepository, userRepository, db, validator)

	authHandler := authhandler.NewAuthHandlerImpl(authService)
	ticketHandler := tickethandler.NewTicketHandlerImpl(ticketService)
	userHander := userhandler.NewUserHandlerImpl(userService)
	accountHandler := accounthandler.NewAccountHandlerImpl(accountService)

	// middlewares
	authMiddleware := authmiddleware.NewAuthMiddlewareImpl(ticketRepository, db)
	// end of middlewares

	// routes
	routes.NewAuthRoutes(app, authHandler)
	routes.NewTicketRoutes(app, ticketHandler, authMiddleware)
	routes.NewUserRoutes(app, userHander, authMiddleware)
	routes.NewAccountRoutes(app, accountHandler)
	// end of routes

	app.Run(":" + appConfig.APP_PORT)
}
