package migrations

import (
	"github.com/sfortson/fitness-tracker/server/database"
	"github.com/sfortson/fitness-tracker/server/models"
)

func Migrate() {
	// Migrate the schema
	database.DB.AutoMigrate(&models.User{}, &models.BodyFat{})
}
