package entity

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"myPhotos/config"
	"myPhotos/logger"
	"os"
	"path/filepath"
	"time"
)

var DB *gorm.DB

func InitializeDatabase() {
	logger.InfoLogger.Println("db init...")
	db, err := gorm.Open(sqlite.Open(filepath.Join(config.DataPath, "data.db")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// https://gorm.io/docs/performance.html
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logger.ErrorLogger.Println("gorm initialize error", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(Media{})
	if err != nil {
		logger.ErrorLogger.Println("gorm migration error", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.ErrorLogger.Println("gorm get db error", err)
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(6)
	sqlDB.SetMaxOpenConns(60)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
