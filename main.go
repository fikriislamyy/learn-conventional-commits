package main

import (
	"learn-conventional-commits/controllers"
	"learn-conventional-commits/initializers"
	"learn-conventional-commits/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDb()
	initializers.InitializeRedisClient()
}

func main() {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/user", middleware.SessionMiddleware() ,controllers.GetUserData)
	r.POST("/logout", middleware.SessionMiddleware(),controllers.Logout)
	r.Run()
}