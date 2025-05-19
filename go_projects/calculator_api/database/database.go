package database

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = slog.New(slog.NewTextHandler(os.Stdout, nil))

// func start() {
// 	cmd := exec.Command("docker-compose", "up", "-d")
// 	output := cmd.Stdout
// 	fmt.Println(output)
// }

func StartDbConnection() {
	//start()

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=calculatorDB port=5432 sslmode=disable", user, password)

	gorm.Open(postgres.Open(dsn))
}
