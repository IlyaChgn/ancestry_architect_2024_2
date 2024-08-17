package migrator

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"

	pool "github.com/IlyaChgn/ancestry_architect_2024_2/internal/pkg/server/repository"
)

const migrationsDir = "migrations"

//go:embed migrations/*.sql
var migrationsFS embed.FS

func ApplyMigrations(version uint) {
	migrator, err := mustGetNewMigrator(migrationsFS, migrationsDir)
	if err != nil {
		log.Fatal("Something went wrong while initializing migrator: ", err)
		os.Exit(1)
	}

	postgresURL := pool.NewConnectionString(os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("DATABASE_NAME"))
	connectionStr := fmt.Sprintf("%s?sslmode=disable", postgresURL)

	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal("Unable to connect to database for migrations")
		os.Exit(1)
	}

	defer conn.Close()

	err = migrator.applyMigrations(conn, version)
	if err != nil {
		log.Fatal("something went wrong while applying migrations: ", err)
		os.Exit(1)
	}

	log.Printf("Migration of version %d has been applied", version)
}
