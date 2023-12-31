package db

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

// Config sqlite database configuration
type Config struct {
	FileName string
}

// DataBase is a struct for work with postgres DataBase via gorm
type SQLiteDataBase struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

// NewDataBase creates a new DataBase instance
func NewDataBase(logger *logrus.Logger, config Config) (*SQLiteDataBase, error) {
	// TODO: Hook gorm logger with logrus logger
	gormConfig := &gorm.Config{Logger: gorm_logger.Default.LogMode(gorm_logger.Silent)}

	db, err := gorm.Open(sqlite.Open(config.FileName), gormConfig)
	if err != nil {
		return nil, err
	}

	return &SQLiteDataBase{DB: db, Logger: logger}, nil
}

// InitialMigration initialises the SQLite database schema
func (db *SQLiteDataBase) SQLiteInitialMigration() error {
	if err := db.DB.AutoMigrate(&Order{}, &Assets{}, &TransactionSent{}, &TransactionLog{}); err != nil {
		return err
	}

	return nil
}
