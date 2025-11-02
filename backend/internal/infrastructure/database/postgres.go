package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rod1kutzyy/OnTrack/internal/config"
	appLogger "github.com/rod1kutzyy/OnTrack/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB *gorm.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	appLogger.Logger.Info("Connecting to PostgreSQL database...")

	var gormLogLevel logger.LogLevel
	switch cfg.Logger.Level {
	case "debug", "trace":
		gormLogLevel = logger.Info
	case "warn":
		gormLogLevel = logger.Warn
	case "error":
		gormLogLevel = logger.Error
	default:
		gormLogLevel = logger.Warn
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	appLogger.Logger.Info("Successfully connected to PostgreSQL database")
	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	appLogger.Logger.Info("Closing database connection...")

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	appLogger.Logger.Info("Database connection closed successfully")
	return nil
}

func (p *PostgresDB) GetDB() *gorm.DB {
	return p.DB
}

func (p *PostgresDB) AutoMigrate(models ...interface{}) error {
	appLogger.Logger.Info("Running database migrations...")

	if err := p.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	appLogger.Logger.Info("Database migrations completed successfully")
	return nil
}
