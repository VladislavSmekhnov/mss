package initializers

import (
	"music_streaming_service/auth-service/app/models"
)

func SyncDb() {
	// Migrate the schema
	DB.AutoMigrate(&models.User{})
}
