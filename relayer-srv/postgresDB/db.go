package postgresdb

import (
	"github.com/sirupsen/logrus"
	storage "github.com/volmexfinance/relayers/relayer-srv/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// StorageConfig contains configurations for storage, postgreSQL
type PostgresStorageConfig struct {
	URL        string // DataBase URL for connection
	DBDriver   string // DataBase driver
	DBHOST     string // DataBase host
	DBPORT     int64
	DBSSL      string // DataBase sslmode
	DBName     string // DataBase name
	DBUser     string // DataBase's user
	DBPassword string // User's password
}

type PostgresDataBase struct {
	DB     *gorm.DB
	Logger *logrus.Logger
}

// TODO: Tx sent confirmed to be moved to postgres
func InitialMigration(dbURL string) (*PostgresDataBase, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&storage.Order{}, &storage.Assets{})

	return &PostgresDataBase{DB: db}, err
}
