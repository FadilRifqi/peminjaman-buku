package middlewares

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Admin(c *gin.Context) {
    // Retrieve the user from the context
    user, _ := c.Get("user")

	u,ok := user.(models.User)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

    // Check if the user's role is "admin"
    if u.Role.Name != "admin" {
        c.AbortWithStatus(http.StatusForbidden)
        return
    }

    c.Next()
}
