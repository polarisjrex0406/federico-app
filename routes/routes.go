package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/polarisjrex0406/federico-app/handlers"
	"github.com/polarisjrex0406/federico-app/middleware"
)

func SetupRouter(
	r *gin.Engine,
	userHandler handlers.UserHandler,
) *gin.Engine {
	r.Use(middleware.LoggingRequests())

	userGroup := r.Group("/user")
	{
		userGroup.POST("/:userId/transaction", middleware.HeaderValidation(), middleware.UserDoTransactionValidation(), userHandler.DoTransaction)
		userGroup.GET("/:userId/balance", userHandler.GetBalance)
	}

	return r
}
