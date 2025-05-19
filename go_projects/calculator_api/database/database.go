package database

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))
var DB *gorm.DB

func StartDbConnection() {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=calculatorDB port=5432 sslmode=disable", user, password)

	log.Info("connecting to database...")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err.Error())
	}

}
