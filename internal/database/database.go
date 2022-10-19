package database

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/sfortson/fitness-tracker/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	// ErrUserFound should be returned from Create (see ConfirmUser)
	// when the primaryID of the record is found.
	ErrUserFound = errors.New("user found")
	// ErrUserNotFound should be returned from Get when the record is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrTokenNotFound should be returned from UseToken when the
	// record is not found.
	ErrTokenNotFound = errors.New("token not found")
)

type User struct {
	gorm.Model
	BodyFatMeasurements []BodyFat
	Email               string
	Username            string
	Birthdate           time.Time
	Password            []byte
	Sex                 string
	Height              float64
}

type BodyFat struct {
	gorm.Model
	UserID     uint
	Year       int
	Month      time.Month
	Day        int
	Neck       float64
	Waist      float64
	Weight     float64
	Height     float64
	Percentage float64
	BMI        float64
}

type Session struct {
	gorm.Model
	Username     string
	Expiry       time.Time
	SessionToken string
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func Open(config config.Config) error {
	var err error

	// if config.AppEnv == "prod" || config.AppEnv == "production" {
	// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/New_York", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	// 	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// } else {
	// 	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// }
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/New_York", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}
	return nil
}

func Migrate() {
	// Migrate the schema
	DB.AutoMigrate(&User{}, &BodyFat{}, &Session{})
}

func LookupUserByToken(_ context.Context, tok string) (*User, error) {
	var session Session
	tokenResult := DB.Where("session_token = ?", tok).First(&session)
	if tokenResult.Error != nil {
		return nil, ErrTokenNotFound
	}

	var user User
	userResult := DB.Where("username = ?", session.Username).First(&user)
	if userResult.Error != nil {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

func (user *User) GetAge() int {
	// Get Age
	year2, _, _ := time.Now().Date()
	year1, _, _ := user.Birthdate.Date()
	year := int(math.Abs(float64(int(year2 - year1))))
	return year
}
