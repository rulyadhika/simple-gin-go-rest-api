package main

import (
	"github.com/gin-gonic/gin"
	authhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/auth_handler"
	tickethandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/ticket_handler"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/db"
	authmiddleware "github.com/rulyadhika/simple-gin-go-rest-api/infra/middleware/auth_middleware"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/packages/validation"
	ticketrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/ticket_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	"github.com/rulyadhika/simple-gin-go-rest-api/routes"
	authservice "github.com/rulyadhika/simple-gin-go-rest-api/service/auth_service"
	ticketservice "github.com/rulyadhika/simple-gin-go-rest-api/service/ticket_service"
)

func main() {
	appConfig := config.GetAppConfig()
	db := db.InitDB()
	validator := validation.NewValidator()
	authMiddleware := authmiddleware.NewAuthMiddlewareImpl()
	app := gin.Default()

	userRepository := userrepository.NewUserRepositoryImpl()
	authService := authservice.NewAuthServiceImpl(userRepository, db, validator)
	authHandler := authhandler.NewAuthHandlerImpl(authService)

	ticketRepository := ticketrepository.NewTicketRepositoryImpl()
	ticketService := ticketservice.NewTicketServiceImpl(ticketRepository, db, validator)
	ticketHandler := tickethandler.NewTicketHandlerImpl(ticketService)

	// routes
	routes.NewAuthRoutes(app, authHandler)
	routes.NewTicketRoutes(app, ticketHandler, authMiddleware)
	// end of routes

	app.Run(":" + appConfig.APP_PORT)
}
