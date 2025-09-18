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
	HTTPPort             int           `mapstructure:"HTTP_PORT"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	LogLevel             string        `mapstructure:"LOG_LEVEL"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	RedisURL             string        `mapstructure:"REDIS_URL"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	SMTPHost             string        `mapstructure:"SMTP_HOST"`
	SMTPPort             int           `mapstructure:"SMTP_PORT"`
	SMTPUsername         string        `mapstructure:"SMTP_USERNAME"`
	SMTPPassword         string        `mapstructure:"SMTP_PASSWORD"`

	// OAuth Configuration
	GoogleOAuthClientID     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOAuthClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectURL  string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
	GitHubOAuthClientID     string `mapstructure:"GITHUB_OAUTH_CLIENT_ID"`
	GitHubOAuthClientSecret string `mapstructure:"GITHUB_OAUTH_CLIENT_SECRET"`
	GitHubOAuthRedirectURL  string `mapstructure:"GITHUB_OAUTH_REDIRECT_URL"`

	//CORS + TLS
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
	EnableTLS      bool     `mapstructure:"ENABLE_TLS"`
	TLSCertFile    string   `mapstructure:"TLS_CERT_FILE"`
	TLSKeyFile     string   `mapstructure:"TLS_KEY_FILE"`
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
	viper.BindEnv("LOG_LEVEL")
	viper.BindEnv("HTTP_SERVER_ADDRESS")
	viper.BindEnv("HTTP_PORT")
	viper.BindEnv("REDIS_URL")
	viper.BindEnv("DB_SOURCE")
	viper.BindEnv("TESTING_DB_SOURCE")
	viper.BindEnv("ACCESS_TOKEN_DURATION")
	viper.BindEnv("REFRESH_TOKEN_DURATION")
	viper.BindEnv("TOKEN_SYMMETRIC_KEY")
	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("SMTP_USERNAME")
	viper.BindEnv("SMTP_PASSWORD")

	//OAuth
	viper.BindEnv("GOOGLE_OAUTH_CLIENT_ID")
	viper.BindEnv("GOOGLE_OAUTH_CLIENT_SECRET")
	viper.BindEnv("GOOGLE_OAUTH_REDIRECT_URL")
	viper.BindEnv("GITHUB_OAUTH_CLIENT_ID")
	viper.BindEnv("GITHUB_OAUTH_CLIENT_SECRET")
	viper.BindEnv("GITHUB_OAUTH_REDIRECT_URL")

	//CORS + TLS
	viper.BindEnv("ALLOWED_ORIGINS")
	viper.BindEnv("ENABLE_TLS")
	viper.BindEnv("TLS_CERT_FILE")
	viper.BindEnv("TLS_KEY_FILE")

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}
	return
}
