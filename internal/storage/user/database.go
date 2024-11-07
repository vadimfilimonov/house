package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/vadimfilimonov/house/internal/models"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		if err := db.Close(); err != nil {
			log.Println(err)
		}

		return nil, err
	}

	err = runMigrations(db)
	if err != nil {
		if err := db.Close(); err != nil {
			log.Println(err)
		}

		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}

func (d *Database) Add(email, hashedPassword, userType string) (*string, error) {
	if err := d.db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}

	id, err := generateUserID()
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (user_id, email, password, user_type) VALUES ($1, $2, $3, $4)`
	_, err = d.db.Exec(query, *id, email, hashedPassword, userType)
	if err != nil {
		return nil, fmt.Errorf("cannot add user to database: %w", err)
	}

	return id, nil
}

func (d *Database) Get(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var email string
	var hashedPassword string
	var userType string
	query := "SELECT email, password, user_type FROM users WHERE user_id = $1 LIMIT 1"
	err := d.db.QueryRowContext(ctx, query, id).Scan(&email, &hashedPassword, &userType)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return &models.User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
		UserType: userType,
	}, nil
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

	m.Up()

	return nil
}
