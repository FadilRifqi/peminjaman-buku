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
	routes.BukuRoutes(r)
	routes.PeminjamanRoutes(r)
	routes.FriendshipRoutes(r)
	routes.ChatRoutes(r)

	r.Run()
}
