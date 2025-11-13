package config

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func DBConnection() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(DB_PATH.GetValue()+"?_loc=Asia%2FJakarta"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to SQLite database 🗃️")

	return db, nil
}
