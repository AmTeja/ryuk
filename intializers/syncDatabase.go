package intializers

import "github.com/amteja/ryuk/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Profile{})
}
