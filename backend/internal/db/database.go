package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/homecooking/backend/internal/config"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func New(cfg *config.DatabaseConfig) (*Database, error) {
	var db *sql.DB
	var err error

	switch cfg.Type {
	case "postgres":
		connStr := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
		)
		db, err = sql.Open("postgres", connStr)
	case "sqlite":
		db, err = sql.Open("sqlite3", cfg.Path)
	case "mysql":
		return nil, fmt.Errorf("mysql not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Printf("Connected to %s database", cfg.Type)

	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}

func (db *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}
