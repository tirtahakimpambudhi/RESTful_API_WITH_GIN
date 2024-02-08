package db

import (
	"fmt"
	"go_gin/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"time"
)

type PGStore struct {
	Database *config.DB
}

func NewPGStore(database *config.DB) *PGStore {
	return &PGStore{Database: database}
}

func (s *PGStore) Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=America/New_York", s.Database.Host, s.Database.Username, s.Database.Password, s.Database.DatabaseName, s.Database.Port)
	//Best Perfomance Config GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	return db, nil
}
