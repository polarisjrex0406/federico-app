package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/polarisjrex0406/federico-app/cmd"
	"github.com/polarisjrex0406/federico-app/config"
	"github.com/polarisjrex0406/federico-app/database"
	"github.com/polarisjrex0406/federico-app/handlers"
	"github.com/polarisjrex0406/federico-app/repositories"
	"github.com/polarisjrex0406/federico-app/routes"
	"github.com/polarisjrex0406/federico-app/services"
	"github.com/polarisjrex0406/federico-app/utils"
	"gorm.io/gorm"

	_ "github.com/polarisjrex0406/federico-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger
//
//	@title                       Federico App API
//	@version                     1.0
//	@description                 A comprehensive API for test.
//	@termsOfService              http://example.com/terms
//	@contact.name                API Support Team
//	@contact.url                 http://example.com/support
//	@contact.email               support@example.com
//	@license.name                Apache 2.0
//	@license.url                 http://www.apache.org/licenses/LICENSE-2.0.html
//	@schemes                     http https
//	@securityDefinitions.apiKey  BearerAuth
//	@in                          header
//	@name                        Authorization
//	@description                 JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	setupValidators()

	database.Connect()

	if !cmd.Commands(database.DB) {
		log.Fatalf("Error while running commands: %v", err)
	}

	userHandler := setupDependencyInjections(database.DB)

	r := routes.SetupRouter(
		gin.Default(),
		userHandler,
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

func setupDependencyInjections(db *gorm.DB) handlers.UserHandler {
	// DI to repositories
	balanceRepo := repositories.NewBalanceRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	// DI to services
	userService := services.NewUserService(balanceRepo, transactionRepo)
	// DI to handlers
	userHandler := handlers.NewUserHandler(userService)

	return userHandler
}

func setupValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("amount", utils.AmountValidator)
	}
}
