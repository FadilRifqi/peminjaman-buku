package routes

import (
	"main/controllers"
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func BukuRoutes(router *gin.Engine) {
	buku := router.Group("/buku")
	buku.Use(middlewares.RequireAuth )
	{
		buku.POST("/add", controllers.AddBuku)
		buku.GET("/all", controllers.GetAllBukus)
		buku.GET("/:id", controllers.GetBukuById)
		buku.GET("/me", controllers.GetMyBukus)
		buku.PATCH("/update/:id", middlewares.IsOwner, controllers.UpdateBuku)
		buku.DELETE("/delete/:id", middlewares.IsOwner, controllers.DeleteBuku)
	}
}
