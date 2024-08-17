package main

import (
	"flag"
	"log"

	migrator "github.com/IlyaChgn/ancestry_architect_2024_2/db"
	app "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server"
	"github.com/joho/godotenv"
)

var migrationFlag = flag.Uint("migrate", 0, "Применяет миграцию указанной версии. Если указан 0, то флаг игонорируется")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading env file", err)
	}

	flag.Parse()
	if *migrationFlag != 0 {
		migrator.ApplyMigrations(*migrationFlag)
	}

	srv := new(app.Server)

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
