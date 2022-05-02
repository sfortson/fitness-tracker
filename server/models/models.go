package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BodyFatMeasurements []BodyFat
	FirstName           string
	LastName            string
	Email               string
}

type BodyFat struct {
	Neck   float32
	Waist  float32
	UserID uint
}

func Model() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &BodyFat{})
}
