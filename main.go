package main

import (
	"learn-conventional-commits/controllers"
	"learn-conventional-commits/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDb()
}

func main() {
	r := gin.Default()
	r.POST("/register", controllers.Register)
	r.Run()
}