package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BodyFatMeasurements []BodyFat
	Email               string
	Username            string
	Birthdate           time.Time
	Password            []byte
	Sex                 string
}

type BodyFat struct {
	Neck   float32
	Waist  float32
	UserID uint
}
