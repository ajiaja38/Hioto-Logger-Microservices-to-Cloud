package utils

import (
	"go/hioto-logger/pkg/model"

	"gorm.io/gorm"
)

func AutoMigrateDb(db *gorm.DB) {
	db.AutoMigrate(&model.Log{})
	db.AutoMigrate(&model.LogAktuator{})
	db.AutoMigrate(&model.MonitoringHistory{})
}
