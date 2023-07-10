package main

import (
	"github.com/amteja/ryuk/controllers"
	"github.com/amteja/ryuk/intializers"
	"github.com/amteja/ryuk/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	intializers.LoadEnvVariables()
	intializers.ConnectToDb()
	intializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Valid Token",
		})
	})
	r.Run()
}
