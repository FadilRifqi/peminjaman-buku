package controllers

import (
	"main/database"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TODO: Jwt token

func GenerateToken(c *gin.Context) {
	//Get Email and Password
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password are required"})
		return
	}
	//Check if the email and password are correct
	var user models.User
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}

	// refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"sub": user.ID,
	// 	"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	// })

	// refreshToken, err := refresh.SignedString(os.Getenv("SECRET"))

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
	// 	return
	// }



	//Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"sub": user.ID,
	// "refresh": refresh,
	"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	//Send it back to the user
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	// Set Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}

func Validate(c *gin.Context){
	user,_ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"user": user})
}
