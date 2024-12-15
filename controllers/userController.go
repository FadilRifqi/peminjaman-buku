package controllers

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context){
	// Get Email and Password
	var body struct {
		Username 	string
		Email 		string
		Password 	string
	}

	// Validate Email and Password
	if c.Bind(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }

    if body.Email == "" || body.Password == "" || body.Username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password cannot be empty"})
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

func UpdateUser(c *gin.Context){
    // Find the User in the Database with url parameter
    var user models.User
    result := database.DB.First(&user, c.Param("id"))

    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Get Email and Username
    var body struct {
        Username string
        Email    string
    }

    // Validate Email and Username
    if c.Bind(&body) != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
        return
    }

    // Update the User from body if body fields are not empty, otherwise use existing data
    if body.Username != "" {
        user.Username = body.Username
    }
    if body.Email != "" {
        user.Email = body.Email
    }

    result = database.DB.Save(&user)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
        return
    }

    // Respond
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context){
	// Find the User in the Database with url parameter
	var user models.User
	result := database.DB.First(&user, c.Param("id"))

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete the User from the Database
	result = database.DB.Delete(&user)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
