package testing

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/homecooking/backend/internal/db/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*sql.DB, *sqlc.Queries, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, err
	}

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, nil, err
	}

	// Read migration files (SQLite-compatible for testing)
	migrations := []string{
		"001_init_sqlite.up.sql",
		"002_add_recipe_variations_sqlite.up.sql",
	}

	for _, migration := range migrations {
		migrationPath := filepath.Join("..", "db", "migrations", migration)
		migrationSQL, err := os.ReadFile(migrationPath)
		if err != nil {
			db.Close()
			return nil, nil, err
		}

		// Execute migration
		_, err = db.Exec(string(migrationSQL))
		if err != nil {
			db.Close()
			return nil, nil, err
		}
	}

	q := sqlc.New(db)
	return db, q, nil
}

// TeardownTestDB closes the test database connection
func TeardownTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
