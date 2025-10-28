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
	App      AppConfig
}

type ServerConfig struct {
	Host string
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AppConfig struct {
	LogLevel string
}

var (
	config *Config
	once   sync.Once
)

func Load() (*Config, error) {
	var err error

	once.Do(func() {
		_ = godotenv.Load()

		config = &Config{
			Server: ServerConfig{
				Host: getEnv("SERVER_HOST", "localhost"),
				Port: getEnv("SERVER_PORT", "8080"),
				Mode: getEnv("GIN_MODE", "debug"),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnv("DB_PORT", "5432"),
				User:     getEnv("DB_USER", "postgres"),
				Password: getEnv("DB_PASSWORD", "postgres"),
				Name:     getEnv("DB_NAME", "ontrack_db"),
				SSLMode:  getEnv("DB_SSLMODE", "disable"),
			},
			App: AppConfig{
				getEnv("LOG_LEVEL", "info"),
			},
		}

		if config.Database.User == "" || config.Database.Name == "" {
			err = fmt.Errorf("database user and name are required")
		}
	})

	return config, err
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}
