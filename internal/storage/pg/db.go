package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

type Storage struct {
	db *sql.DB
}

var (
	storage  = &Storage{}
	syncErr  error
	initOnce = sync.Once{}
)

func New(connectionString string) (*Storage, error) {
	initOnce.Do(func() {
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			syncErr = err
			return
		}

		if err := runMigrations(db); err != nil {
			if err := db.Close(); err != nil {
				log.Println(err)
			}

			syncErr = err
			return
		}

		if err := db.Ping(); err != nil {
			syncErr = fmt.Errorf("cannot ping db: %w", err)
			return
		}
	})

	return storage, syncErr
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		log.Println(err)
	}

	return nil
}
