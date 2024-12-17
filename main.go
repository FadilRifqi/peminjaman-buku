package main

import (
	"main/database"
	"main/initializers"
	"main/routes"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnvVariables()
	database.Connect()
	database.SyncDB()
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},              // Change this to the origin you want
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        AllowCredentials: true,                       // Required for withCredentials
        MaxAge:           12 * time.Hour,             // Optional, max age for pre-flight requests
    }))

	routes.UserRoutes(r)
	routes.BukuRoutes(r)
	routes.PeminjamanRoutes(r)
	routes.FriendshipRoutes(r)
	routes.ChatRoutes(r)
	routes.RoomRoutes(r)
	routes.WebSocketRoutes(r)

	r.Run()
}
