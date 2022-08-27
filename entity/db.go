package entity

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"myPhotos/config"
	"path/filepath"
)

var DB *gorm.DB

func init() {
	db, _ := gorm.Open(sqlite.Open(filepath.Join(config.DataPath, "data.db")), &gorm.Config{
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
