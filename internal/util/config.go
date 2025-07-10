// Package util
package util

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	SMTPHost             string        `mapstructure:"SMTP_HOST"`
	SMTPPort             int           `mapstructure:"SMTP_PORT"`
	SMTPUsername         string        `mapstructure:"SMTP_USERNAME"`
	SMTPPassword         string        `mapstructure:"SMTP_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	absPath, err := filepath.Abs(path)
	if err != nil {
		return config, fmt.Errorf("failed to resolve absolute path: %v", err)
	}

	env := viper.GetString("ENVIRONMENT")

	if env == "" || env == "development" {
		viper.AddConfigPath(absPath) // Path to look for the file
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		if err = viper.ReadInConfig(); err != nil {
			fmt.Println("No .env file found, relying on environment variables")
		} else {
			fmt.Println("Loaded .env file for local development")
		}
	}

	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("HTTP_SERVER_ADDRESS")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("TESTING_DB_SOURCE")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("SMTP_USERNAME")
	viper.BindEnv("SMTP_PASSWORD")

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Error unmarshalling config")
	}
	return
}
