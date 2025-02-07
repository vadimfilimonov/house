package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"

	"github.com/vadimfilimonov/house/internal/models"
)

var (
	ErrUserNotFound = errors.New("user is not found")
)

type Database struct {
	db *sql.DB
}

func New(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		if err := db.Close(); err != nil {
			log.Println(err)
		}

		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) Add(ctx context.Context, email, hashedPassword, userType string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id := uuid.New().String()

	query := `INSERT INTO users (user_id, email, password, user_type) VALUES ($1, $2, $3, $4)`
	_, err := d.db.ExecContext(ctx, query, id, email, hashedPassword, userType)
	if err != nil {
		return nil, fmt.Errorf("cannot add user to database: %w", err)
	}

	return &id, nil
}

func (d *Database) Get(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var email string
	var hashedPassword string
	var userType string

	query := "SELECT email, password, user_type FROM users WHERE user_id = $1 LIMIT 1"

	sqlRow := d.db.QueryRowContext(ctx, query, id)
	if sqlRow == nil {
		return nil, fmt.Errorf("sql row is nil")
	}

	if err := sqlRow.Scan(&email, &hashedPassword, &userType); err != nil {
		return nil, ErrUserNotFound
	}

	return &models.User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
		UserType: userType,
	}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
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
