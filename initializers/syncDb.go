package initializers

import (
	"learn-conventional-commits/models"
)

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}