package middleware

import (
	"context"
	"learn-conventional-commits/initializers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SessionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the session token from Cookie.
        sessionToken, err := c.Cookie("Authorization") // Use the correct cookie name
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Session token not found or expired",
            })
            c.Abort()
            return
        }

        // Initialize the Redis client and retrieve session data based on the session token.
        client, err := initializers.InitializeRedisClient()
        if err != nil {
            log.Println("Failed to connect to Redis:", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to connect to Redis",
            })
            c.Abort()
            return
        }

        // Retrieve the session data using the session token.
        email, err := client.HGet(context.Background(), sessionToken, "email").Result()
        if err != nil {
            log.Println("Error retrieving email from Redis:", err)
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired session",
            })
            c.Abort()
            return
        }

        id, err := client.HGet(context.Background(), sessionToken, "id").Result()
        if err != nil {
            log.Println("Error retrieving id from Redis:", err)
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired session",
            })
            c.Abort()
            return
        }

        token, err := client.HGet(context.Background(), sessionToken, "token").Result()
        if err != nil {
            log.Println("Error retrieving token from Redis:", err)
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired session",
            })
            c.Abort()
            return
        }

        sessionData := map[string]interface{}{
            "email": email,
            "id":    id,
            "token": token,
            // Add other session-related information as needed.
        }

        c.Set("userData", sessionData)

        c.Next()
    }
}