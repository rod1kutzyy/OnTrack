package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Host string
	Port string
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
				Host: getEnv("SERVER_HOST", "localhost"),
				Port: getEnv("SERVER_PORT", "8080"),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", ""),
				Name:     getEnv("DB_NAME", "ontrack_db"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},
			Logger: LoggerConfig{
				Level: getEnv("LOG_LEVEL", "info"),
			},
		}

		if config.Database.User == "" || config.Database.Name == "" {
			err = fmt.Errorf("database user and name are required")
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
