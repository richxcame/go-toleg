package routes

import (
	"gotoleg/web/handlers"
	"gotoleg/web/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	routes := gin.Default()

	// Cors defualt config
	routes.Use(cors.Default())

	api := routes.Group("/api")

	{
		api.GET("/transactions", middlewares.Auth(), handlers.GetTransactions)
		api.POST("/transactions", middlewares.Auth(), handlers.SendTransactions)
		api.POST("/transactions/:uuid", middlewares.Auth(), handlers.SendTransaction)
		api.POST("/declined-transactions", middlewares.Auth(), handlers.ResendDeclinedTrxns)
		api.POST("/declined-transactions/:uuid", middlewares.Auth(), handlers.ResendDeclinedTrxn)

		api.POST("/auth/login", handlers.Login)
		api.POST("/auth/token", handlers.Token)
	}

	return routes

}
