package middlewares

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PeminjamanMiddleware(c *gin.Context) {
	user, _ := c.Get("user")

	u, ok := user.(models.User)

	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	peminjamanId := c.Param("id")
	var peminjaman models.Peminjaman
	result := database.DB.First(&peminjaman, peminjamanId)

	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Peminjaman not found"})
		return
	}

	if peminjaman.UserID != u.ID {
		if u.Role.Name == "admin" {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}
	c.Next()
}
