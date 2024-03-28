package main

import (
	"errors"
	"flag"
	"fmt"

	//  Library for migrations
	"github.com/golang-migrate/migrate/v4"
	//  Driver for executing migrations SQLite3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	//  Driver for getting migrations from files
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// function for execution database migration
func main() {
	var storagePath, migrationsPath, migrationsTable string

	// reading parameters for database migration
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migration table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	// preparing database migration
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	// executing up database migration
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}

		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
