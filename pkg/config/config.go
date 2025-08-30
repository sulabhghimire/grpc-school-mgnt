package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	DB            MongoDBConfig
	Server        ServerConfig
	TLS           TLSConfig
	JWT           JWTConfig
	PasswordReset PasswordReset
}

type ServerConfig struct {
	Port string
}

type TLSConfig struct {
	CertPath string
	KeyPath  string
}

type MongoDBConfig struct {
	URI string
}

type JWTConfig struct {
	SecretKey      string
	ExpirationTime time.Duration
}

type PasswordReset struct {
	ResetTokenExpiryDuration time.Duration
	FromEmail                string
	EmailHost                string
	EmailPort                int
}

var GlobalConfig *AppConfig

func GetAppConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		DB: MongoDBConfig{
			URI: os.Getenv("DB_URI"),
		},
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		TLS: TLSConfig{
			CertPath: os.Getenv("TLS_CERT_PATH"),
			KeyPath:  os.Getenv("TLS_KEY_PATH"),
		},
		JWT: JWTConfig{
			SecretKey:      os.Getenv("JWT_SECRET_KEY"),
			ExpirationTime: time.Duration(1 * time.Hour), // Default value, can be overridden by env variable
		},
		PasswordReset: PasswordReset{
			ResetTokenExpiryDuration: time.Duration(1 * time.Hour), // Default value, can be overridden by env variable
			FromEmail:                os.Getenv("EMAIL_FROM"),
			EmailHost:                os.Getenv("EMAIL_HOST"),
		},
	}

	// Basic validation
	missing := []string{}

	if cfg.DB.URI == "" {
		missing = append(missing, "URI")
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "3000"
	}
	if cfg.TLS.CertPath == "" {
		missing = append(missing, "TLS_CERT_PATH")
	}
	if cfg.TLS.KeyPath == "" {
		missing = append(missing, "TLS_KEY_PATH")
	}
	if cfg.JWT.SecretKey == "" {
		missing = append(missing, "JWT_SECRET_KEY")
	}
	if cfg.PasswordReset.FromEmail == "" {
		missing = append(missing, "EMAIL_FROM")
	}
	if cfg.PasswordReset.EmailHost == "" {
		missing = append(missing, "EMAIL_HOST")
	}

	if emailPortStr := os.Getenv("EMAIL_PORT"); emailPortStr == "" {
		missing = append(missing, "EMAIL_PORT")
	} else {
		if emailPort, err := strconv.Atoi(emailPortStr); err != nil {
			missing = append(missing, "EMAIL_PORT")
		} else {
			cfg.PasswordReset.EmailPort = emailPort
		}
	}

	if expirationTimeStr := os.Getenv("JWT_EXPIRATION_TIME"); expirationTimeStr != "" {
		expirationTime, err := time.ParseDuration(expirationTimeStr + "m")
		if err == nil {
			cfg.JWT.ExpirationTime = expirationTime
		}
	}

	if resetTokenExpiryDurationStr := os.Getenv("RESET_TOKEN_EXP_DURATION"); resetTokenExpiryDurationStr != "" {
		expirationTime, err := time.ParseDuration(resetTokenExpiryDurationStr + "m")
		if err == nil {
			cfg.PasswordReset.ResetTokenExpiryDuration = expirationTime
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required environment variables: %v", missing)
	}

	GlobalConfig = cfg

	return cfg, nil
}
