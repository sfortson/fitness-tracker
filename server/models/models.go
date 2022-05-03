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

type Session struct {
	Username     string
	Expiry       time.Time
	SessionToken string
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
