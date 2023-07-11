package utils

import (
	"github.com/amteja/ryuk/ecodes"
	"github.com/amteja/ryuk/intializers"
	"github.com/amteja/ryuk/models"
	"github.com/gin-gonic/gin"
)

func CheckProfileCompletetion(c *gin.Context) int {
	user := c.MustGet("user").(models.User)

	var profile models.Profile
	intializers.DB.Where("user_id = ?", user.ID).First(&profile)

	if profile.ID == 0 {
		return ecodes.PROFILE_NOT_FOUND
	}

	if profile.DisplayName == "" || profile.UserName == "" || profile.Bio == "" || profile.DisplayURL == "" {
		return ecodes.PROFILE_INCOMPLETE
	}

	return ecodes.NO_ERROR
}
