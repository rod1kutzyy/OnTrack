package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host        string
	Port        string
	FrontedURLs []string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type LoggerConfig struct {
	Level string
}

var (
	config *Config
	once   sync.Once
)

func Load() (*Config, error) {
	var err error

	once.Do(func() {
		err = godotenv.Load()

		config = &Config{
			Server: ServerConfig{
				Host:        getEnv("SERVER_HOST", "0.0.0.0"),
				Port:        getEnv("SERVER_PORT", "8080"),
				FrontedURLs: strings.Split(getEnv("FRONTEND_URLS", "http://localhost:5173,http://127.0.0.1:5173"), ","),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "db"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "password"),
				Name:     getEnv("DB_NAME", "ontrack"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},
			Logger: LoggerConfig{
				Level: getEnv("LOG_LEVEL", "info"),
			},
		}
	})

	return config, err
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
