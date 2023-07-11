package controllers

import (
	"strings"

	"github.com/amteja/ryuk/ecodes"
	"github.com/amteja/ryuk/intializers"
	"github.com/amteja/ryuk/models"
	"github.com/gin-gonic/gin"
)

func CreateProfile(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var profile models.Profile

	if profile.ID != 0 {
		c.JSON(400, models.ServerResponse{
			Message: "Profile already exists",
			Code:    ecodes.PROFILE_ALREADY_EXISTS,
			Data:    nil,
		})
		return
	}

	var body struct {
		DisplayName string `json:"display_name" binding:"required"`
		UserName    string `json:"username" binding:"required"`
		Bio         string `json:"bio,omitempty"`
		DisplayURL  string `json:"display_url,omitempty"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(400, models.ServerResponse{
			Message: "Bad Body Format",
			Code:    ecodes.DEFAULT_BAD_REQUEST,
			Data:    nil,
		})
		return
	}

	// create profile
	profile = models.Profile{
		DisplayName: body.DisplayName,
		UserName:    body.UserName,
		Bio:         body.Bio,
		DisplayURL:  body.DisplayURL,
		UserID:      userId,
	}

	result := intializers.DB.Create(&profile)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "unique") {
			c.JSON(400, models.ServerResponse{
				Message: "Profile already exists",
				Code:    ecodes.PROFILE_ALREADY_EXISTS,
				Data:    nil,
			})
			return
		}
		c.JSON(500, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	c.JSON(200, models.ServerResponse{
		Message: "Profile created successfully",
		Code:    ecodes.NO_ERROR,
		Data:    profile,
	})
}

func GetProfile(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var profile models.Profile

	result := intializers.DB.Where("user_id = ?", userId).First(&profile)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "record not found") {
			c.JSON(400, models.ServerResponse{
				Message: "Profile not found",
				Code:    ecodes.PROFILE_NOT_FOUND,
				Data:    nil,
			})
			return
		}

		c.JSON(500, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	c.JSON(200, models.ServerResponse{
		Message: "Profile fetched successfully",
		Code:    ecodes.NO_ERROR,
		Data:    profile,
	})

}

func UpdateProfile(c *gin.Context) {
	userId := c.MustGet("userId").(uint)

	var profile models.Profile

	var body struct {
		DisplayName string `json:"display_name,omitempty"`
		UserName    string `json:"username,omitempty"`
		Bio         string `json:"bio,omitempty"`
		DisplayURL  string `json:"display_url,omitempty"`
	}

	if c.BindJSON(&body) != nil {
		c.JSON(400, models.ServerResponse{
			Message: "Bad Body Format",
			Code:    ecodes.DEFAULT_BAD_REQUEST,
			Data:    nil,
		})
		return
	}

	if body.DisplayName == "" && body.UserName == "" && body.Bio == "" && body.DisplayURL == "" {
		c.JSON(400, models.ServerResponse{
			Message: "Bad Body Format",
			Code:    ecodes.DEFAULT_BAD_REQUEST,
			Data:    nil,
		})
		return
	}

	result := intializers.DB.Where("user_id = ?", userId).First(&profile)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "record not found") {
			c.JSON(400, models.ServerResponse{
				Message: "Profile not found",
				Code:    ecodes.PROFILE_NOT_FOUND,
				Data:    nil,
			})
			return
		}

		c.JSON(500, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	if body.DisplayName != "" {
		profile.DisplayName = body.DisplayName
	}

	if body.UserName != "" {
		profile.UserName = body.UserName
	}

	if body.Bio != "" {
		profile.Bio = body.Bio
	}

	if body.DisplayURL != "" {
		profile.DisplayURL = body.DisplayURL
	}

	result = intializers.DB.Save(&profile)

	if result.Error != nil {
		c.JSON(500, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	c.JSON(200, models.ServerResponse{
		Message: "Profile updated successfully",
		Code:    ecodes.NO_ERROR,
		Data:    profile,
	})

}

func DeleteProfile(c *gin.Context) {

	userId := c.MustGet("userId").(uint)

	var profile models.Profile

	result := intializers.DB.Where("user_id = ?", userId).First(&profile)

	if result.Error != nil {
		c.JSON(500, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	result = intializers.DB.Delete(&profile)

	if result.Error != nil {
		c.JSON(500, models.ServerResponse{
			Message: "Something went wrong",
			Code:    ecodes.INTERNAL_SERVER_ERROR,
			Data:    nil,
		})
		return
	}

	c.JSON(200, models.ServerResponse{
		Message: "Profile deleted successfully",
		Code:    ecodes.NO_ERROR,
		Data:    nil,
	})
}
