package routes

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	user := router.Group("/user")
	{
		user.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Not implemented",
			})
		})
		user.POST("/register", controllers.Register)
		user.PATCH("/update/:id", controllers.UpdateUser)
		user.DELETE("/delete/:id", controllers.DeleteUser)
	}
	auth := router.Group("/auth")
	{
		auth.POST("/login", controllers.GenerateToken)
		auth.POST("/refresh", controllers.RefreshToken)
	}
}
