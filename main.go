package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	authhandler "github.com/rulyadhika/simple-gin-go-rest-api/handler/auth_handler"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/config"
	"github.com/rulyadhika/simple-gin-go-rest-api/infra/db"
	rolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/role_repository"
	userrepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_repository"
	userrolerepository "github.com/rulyadhika/simple-gin-go-rest-api/repository/user_role_repository"
	"github.com/rulyadhika/simple-gin-go-rest-api/routes"
	authservice "github.com/rulyadhika/simple-gin-go-rest-api/service/auth_service"
)

func main() {
	appConfig := config.GetAppConfig()
	db := db.InitDB()
	validator := validator.New()
	app := gin.Default()

	userRepository := userrepository.NewUserRepositoryImpl()
	userRoleRepository := userrolerepository.NewUserRoleRepositoryImpl()
	roleRepository := rolerepository.NewRoleRepositoryImpl()
	authService := authservice.NewAuthServiceImpl(userRepository, userRoleRepository, roleRepository, db, validator)
	authHandler := authhandler.NewAuthHandlerImpl(authService)

	// routes
	routes.NewAuthRoutes(app, authHandler)
	// end of routes

	app.Run(":" + appConfig.APP_PORT)
}
