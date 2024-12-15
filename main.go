package main

import (
	"main/database"
	"main/initializers"
	"main/routes"

	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVariables()
	database.Connect()
	database.SyncDB()
}

func main() {
	r := gin.Default()

	routes.UserRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
