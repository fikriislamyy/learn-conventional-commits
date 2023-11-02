package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"learn-conventional-commits/initializers"
	"learn-conventional-commits/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func Login(c *gin.Context) {
    // receive email and password from request, check if email exists, store session and user info in Redis, send session token in response
    var body struct {
        Email    string
        Password string
    }
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to read body",
        })
        return
    }

    var user models.User
    initializers.DB.First(&user, "email = ?", body.Email)

    if user.ID == uuid.Nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid email or password",
        })
        return
    }

    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid email or password",
        })
        return
    }

    userEmail := user.Email
    userID := user.ID

    token, err := generateSessionToken(32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to generate session token",
        })
        return
    }

    client, err := initializers.InitializeRedisClient()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to connect to Redis",
        })
        return
    }

    err = client.HMSet(context.Background(), token, map[string]interface{}{
        "email": userEmail,
        "id":    userID,
        "token": token, // Include the token in the Redis hash.
        // Add other session-related information as needed.
    }).Err()

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to store session data in Redis",
        })
        return
    }

    err = client.Expire(context.Background(), token, time.Hour*2).Err()

    c.SetCookie("Authorization", token, int(time.Hour*2), "/", "", false, true)  // Set the cookie

    c.JSON(http.StatusOK, gin.H{
        "session_token": token,
    })
}

func Logout(c *gin.Context) {
    // Get the session token from Cookie.
    sessionToken, err := c.Cookie("Authorization")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Session token not found or expired",
        })
        return
    }

    // Initialize the Redis client and remove the session data using the session token.
    client, err := initializers.InitializeRedisClient()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to connect to Redis",
        })
        return
    }

    // Delete the session data and cookie.
    err = client.Del(context.Background(), sessionToken).Err()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete session data in Redis",
        })
        return
    }

    // Clear the cookie on the client side.
    c.SetCookie("Authorization", "", -1, "/", "", false, true)

    c.JSON(http.StatusOK, gin.H{
        "message": "Logout successful",
    })
}

func GetUserData(c *gin.Context) {
    // Access this endpoint through session middleware.
    userData, exists := c.Get("userData")

    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "User data not found",
        })
        return
    }

    // Ensure userData is of type map[string]interface{}.
    userDataMap, ok := userData.(map[string]interface{})
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Invalid user data format",
        })
        return
    }

    // You now have access to the user data.
    c.JSON(http.StatusOK, userDataMap)
}

func generateSessionToken(length int) (string, error) {
	// Generate a random byte slice of the specified length.
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes as a base64 string.
	token := base64.StdEncoding.EncodeToString(randomBytes)

	// Optionally, remove any characters that might be problematic in URLs.
	token = strings.ReplaceAll(token, "/", "-")
	token = strings.ReplaceAll(token, "+", "_")

	return token, nil
}