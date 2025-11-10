package main

import (
	"go/hioto-logger/config"
	"go/hioto-logger/utils"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
	db, errDb := config.DBConnection()

	if errDb != nil {
		log.Fatal(errDb)
	}

	utils.AutoMigrateDb(db)
}
