package database

import (
	"fmt"
	"gin-template/config"
	"gin-template/logging"
	"gin-template/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// NewDatabase returns a new database connection
// The database connection is made with the configuration given in parameter
// The database connection is made with the postgres driver
// The database connection is made with the gorm library
func NewDatabase(config config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.Host,
		config.Username,
		config.Password,
		config.DatabaseName,
		config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,       // Disable color
			},
		),
	})
	if err != nil {
		logging.Error.Fatal(err)
	}

	// Create enum types
	db.Exec("CREATE TYPE role AS ENUM ('superadmin', 'user', 'admin');")

	// Migrate the schema
	err = db.AutoMigrate(
		model.User{},
		model.Account{},
		model.Token{},
	)
	if err != nil {
		logging.Error.Fatal(err)
	}

	return db
}
