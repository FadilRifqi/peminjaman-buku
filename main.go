package main

import (
	"main/database"
	"main/initializers"

	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVariables()
	database.Connect()
	database.SyncDB()
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
