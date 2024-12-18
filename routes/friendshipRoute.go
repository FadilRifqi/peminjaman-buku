package routes

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func FriendshipRoutes(router *gin.Engine) {
	friendshipRoutes := router.Group("/friendship")
	friendshipRoutes.Use(middlewares.RequireAuth)
	{
		friendshipRoutes.POST("/send/:id", controllers.SendFriendRequest)
		friendshipRoutes.GET("", controllers.GetMyFriends)
		friendshipRoutes.POST("/accept/:id", controllers.AcceptFriendRequest)
		friendshipRoutes.DELETE("/delete/:id", controllers.DeleteFriend)
	}
}
