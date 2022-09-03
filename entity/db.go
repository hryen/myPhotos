package entity

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"myPhotos/config"
	"myPhotos/logger"
	"os"
	"time"
)

var DB *gorm.DB

func InitializeDatabase() {
	logger.InfoLogger.Println("db init...")
	db, err := gorm.Open(mysql.Open(config.DataSourceName), &gorm.Config{
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

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)

	DB = db
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
