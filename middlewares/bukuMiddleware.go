package middlewares

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsOwner(c *gin.Context) {
	// Retrieve the user from the context
	user, _ := c.Get("user")

	u, ok := user.(models.User)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// Retrieve the book id from the context
	bookId := c.Param("id")

	// Check if the user is the owner of the book
	var book models.Buku
	result := database.DB.First(&book, bookId)
	if result.Error != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.UserID != u.ID {
		if u.Role.Name == "admin" {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	c.Next()
}
