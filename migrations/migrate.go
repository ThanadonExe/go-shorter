package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func Run(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/sql",
		"postgres", driver)

	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	m.Up()
}
