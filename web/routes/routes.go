package routes

import (
	"gotoleg/web/handlers"
	"gotoleg/web/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	router := gin.Default()
	router.Use(cors.Default())

	api := router.Group("/api")

	{
		setupTransactionRoutes(api)
		setupAuthRoutes(api)
	}

	return router
}

func setupTransactionRoutes(router *gin.RouterGroup) {

	router.POST("/trxns", handlers.AddTransaction)
	router.GET("/transactions", middlewares.Auth(), handlers.GetTransactions)
	router.POST("/transactions", middlewares.Auth(), handlers.SendTransactions)
	router.POST("/transactions/:uuid", middlewares.Auth(), handlers.SendTransaction)
	router.POST("/declined-transactions", middlewares.Auth(), handlers.ResendDeclinedTrxns)
	router.POST("/declined-transactions/:uuid", middlewares.Auth(), handlers.ResendDeclinedTrxn)
	router.GET("/check-transactions/:uuid", middlewares.Auth(), handlers.CheckTrxnStatus)
	router.POST("/force-add-transactions/:uuid", middlewares.Auth(), handlers.ForceAddDeclinedTransaction)
	router.POST("/force-add-transactions", middlewares.Auth(), handlers.ForceAddDeclinedTransactions)

}

func setupAuthRoutes(router *gin.RouterGroup) {

	router.POST("/auth/login", handlers.Login)
	router.POST("/auth/token", handlers.Token)

}
