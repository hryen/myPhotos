package entity

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"myPhotos/logger"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func init() {
	path, err := os.Executable()
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}

	db, _ := gorm.Open(sqlite.Open(filepath.Join(filepath.Dir(path), "data.db")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// https://gorm.io/docs/performance.html
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	_ = db.AutoMigrate(Media{})

	DB = db
}
