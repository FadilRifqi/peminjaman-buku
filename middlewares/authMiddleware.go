package middlewares

import (
	"errors"
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Get the "Authorization" cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		fmt.Println("Authorization cookie missing:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
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
			fmt.Println("Authorization token expired. Checking refresh token...")
			// Handle token regeneration here
			handleTokenRefresh(c)
			return
		}

		// Any other errors mean the token is invalid
		fmt.Println("Invalid Authorization token:", err.Error())
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Find the user based on the "sub" claim
	var user models.User
	if err := database.DB.First(&user, "id = ?", claims["sub"]).Error; err != nil || user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
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
		fmt.Println("Refresh token missing:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
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
		fmt.Println("Invalid Refresh token:", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Extract claims from the refresh token
	refreshClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || float64(time.Now().Unix()) > refreshClaims["exp"].(float64) {
		fmt.Println("Refresh token expired")
		c.AbortWithStatus(http.StatusUnauthorized)
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
		fmt.Println("Failed to generate new Authorization token:", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Set the new "Authorization" token in the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", newTokenString, 3600 * 24, "/", "", false, true)

	//Get User from the database
	var user models.User
	if err := database.DB.First(&user, "id = ?", refreshClaims["sub"]).Error; err != nil || user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Attach the user ID to the context and proceed
	c.Set("user", user)
	c.Next()
}
