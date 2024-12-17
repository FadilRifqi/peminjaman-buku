package middlewares

import (
	"errors"
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Get the "Authorization" cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		// Abort with status and message
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Token Provided"})
		return
	}

	// Parse the "Authorization" token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// Check for parsing errors or expired tokens
	if err != nil {
		// Handle token expiration specifically
		if errors.Is(err, jwt.ErrTokenExpired) {
			// Handle token regeneration here
			handleTokenRefresh(c)
			return
		}

		// Any other errors mean the token is invalid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	// Find the user based on the "sub" claim
	var user models.User
    if err := database.DB.Preload("Role").First(&user, "id = ?", claims["sub"]).Error; err != nil || user.ID == 0 {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

	// Attach the user to the context
	c.Set("user", user)

	// Proceed to the next handler
	c.Next()
}

// handleTokenRefresh regenerates an access token using a valid refresh token
func handleTokenRefresh(c *gin.Context) {
	// Get the "RefreshToken" cookie
	refreshTokenString, err := c.Cookie("RefreshToken")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "You Must Be Logged In"})
		return
	}

	// Parse and validate the "RefreshToken"
	refreshToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	// Check if the refresh token is valid
	if err != nil || !refreshToken.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	// Extract claims from the refresh token
	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || float64(time.Now().Unix()) > refreshClaims["exp"].(float64) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
		return
	}

	// Generate a new "Authorization" token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": refreshClaims["sub"], // Use the same user ID from the refresh token
		"exp": time.Now().Add(time.Minute * 60).Unix(),
		"iat": time.Now().Unix(),
	})

	newTokenString, err := newToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Set the new "Authorization" token in the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", newTokenString, 3600 * 24, "", "", false, true)

	//Get User from the database
	var user models.User
    if err := database.DB.Preload("Role").First(&user, "id = ?", refreshClaims["sub"]).Error; err != nil || user.ID == 0 {
        c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Attach the user to the context
    c.Set("user", user)
	c.Next()
}

func IsMe(c *gin.Context) {
	// Check if the user is the same as the one in the URL
	user, _ := c.Get("user")
	u, ok := user.(models.User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	// Convert the URL parameter from string to uint
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    if u.ID != uint(id) {
		if u.Role.Name == "admin" {
			c.Next()
			return
		}
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
        return
    }

    c.Next()
}
