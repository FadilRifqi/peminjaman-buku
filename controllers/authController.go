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

func Register(c *gin.Context){
	// Get Email and Password
	var body struct {
		Username 	string	`json:"username"`
		Email 		string	`json:"email"`
		Password 	string	`json:"password"`
	}

	// Validate Email and Password
	if c.Bind(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }

    if body.Email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email cannot be empty"})
        return
    }

	if body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	if body.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username cannot be empty"})
		return
	}

	// Hash the Password
	hash,err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Save the User in the Database
	user := models.User{
		Username: body.Username,
		Email: body.Email,
		Password: string(hash),
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	// Respond
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

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
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Credentials"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
		return
	}

	//Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"sub": user.ID,
	"exp": time.Now().Add(time.Minute * 60).Unix(),
	"iat": time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	//Send it back to the user
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Set Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24, "/", "localhost", false, true)
	c.SetCookie("RefreshToken", refreshTokenString, 3600 * 24 * 7, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Validate(c *gin.Context){
	user, _ := c.Get("user")

    // Create new structs to exclude sensitive fields
    type RoleResponse struct {
        ID   uint   `json:"id"`
        Name string `json:"name"`
    }

    type UserResponse struct {
        ID       uint         `json:"id"`
        Username string       `json:"username"`
        Email    string       `json:"email"`
        RoleID   uint         `json:"role_id"`
        Role     RoleResponse `json:"role"`
    }

    // Cast the user to the models.User type
    userModel, ok := user.(models.User)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User type assertion failed"})
        return
    }

    // Create a response without the password
    userResponse := UserResponse{
        ID:       userModel.ID,
        Username: userModel.Username,
        Email:    userModel.Email,
        RoleID:   userModel.RoleID,
        Role: RoleResponse{
            ID:   userModel.Role.ID,
            Name: userModel.Role.Name,
        },
    }

    // Return the user without the password
    c.JSON(http.StatusOK, gin.H{"user": userResponse})
}


func Logout(c *gin.Context){
	// Delete the Authorization cookie and RefreshToken cookie
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.SetCookie("RefreshToken", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
