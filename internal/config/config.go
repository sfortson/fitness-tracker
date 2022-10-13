package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv     string `mapstructure:"APP_ENV"`
	DBUser     string `mapstructure:"DB_USER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBURI      string `mapstructure:"DB_URI"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("app")

	// viper.SetEnvPrefix("ft")
	// viper.SetConfigType("env")
	viper.BindEnv("app_env")
	viper.BindEnv("db_user")
	viper.BindEnv("db_host")
	viper.BindEnv("db_name")
	viper.BindEnv("db_port")
	viper.BindEnv("db_password")
	viper.BindEnv("db_uri")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}

	err = viper.Unmarshal(&config)
	return
}
