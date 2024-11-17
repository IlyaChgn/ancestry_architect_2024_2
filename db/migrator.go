package migrator

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Структура для применения миграций
type Migrator struct {
	srcDriver source.Driver
}

func mustGetNewMigrator(sqlFiles embed.FS, dirName string) (*Migrator, error) {
	driver, err := iofs.New(sqlFiles, dirName)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		srcDriver: driver,
	}, nil
}

func (m *Migrator) applyMigrations(db *sql.DB, version uint) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		return err
	}

	defer migrator.Close()

	if err = migrator.Migrate(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
