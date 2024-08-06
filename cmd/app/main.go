package main

import (
	"log"

	app "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading env file", err)
	}

	srv := new(app.Server)

	if err := srv.Run(); err != nil {
		log.Fatal("Error occurred while starting server:", err.Error())
	}
}
