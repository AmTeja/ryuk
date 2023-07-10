package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {

	// Get the user ID from the JWT
	userId := c.MustGet("userId").(uint)

	fmt.Print("IN CREATE PROFILE : ", userId)
}

func GetProfile(c *gin.Context) {}

func UpdateProfile(c *gin.Context) {}

func DeleteProfile(c *gin.Context) {}
