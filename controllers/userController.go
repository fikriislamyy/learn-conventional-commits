package controllers

import (
	"learn-conventional-commits/initializers"
	"learn-conventional-commits/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Firstname string
		Lastname string
		Email string
		Password string
		Phone string
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var existingEmailUser models.User
	emailErr := initializers.DB.Where("email = ?", body.Email).First(&existingEmailUser).Error
	if emailErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already in use",
		})
		return
	}

	var existingPhoneUser models.User
	phoneErr := initializers.DB.Where("phone = ?", body.Phone).First(&existingPhoneUser).Error
	if phoneErr == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone already in use",
		})
		return
	}


	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Firstname: body.Firstname,
		Lastname: body.Lastname,
		Email: body.Email,
		Password: string(hash),
		Phone: body.Phone,
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}