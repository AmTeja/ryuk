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
	r.POST("/loginToken", controllers.LoginWithToken)

	r.GET("/profile", middleware.RequireAuth, controllers.GetProfile)
	r.POST("/profile", middleware.RequireAuth, controllers.CreateProfile)
	r.PUT("/profile", middleware.RequireAuth, controllers.UpdateProfile)
	r.DELETE("/profile", middleware.RequireAuth, controllers.DeleteProfile)

	r.Run()
}
