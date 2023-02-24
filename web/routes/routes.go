package routes

import (
	"gotoleg/web/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	routes := gin.Default()

	// Cors defualt config
	routes.Use(cors.Default())

	api := routes.Group("/api")

	{
		api.GET("/transactions", handlers.GetTransactions)
	}

	return routes

}
