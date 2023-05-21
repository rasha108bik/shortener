// module postgres, which implement the Storager interface methods
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	appErr "github.com/rasha108bik/tiny_url/internal/errors"
)

// DB struct with postgres driver *sqlx.DB.
type DB struct {
	db *sqlx.DB
}

// NewPostgres returns a newly initialized handler objects that implements the Storager
// interface.
func NewPostgres(connString string) (*DB, error) {
	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = migrateUP(db)
	if err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

// Close closing db descriptor.
func (db *DB) Close() error {
	return db.db.Close()
}

func migrateUP(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Printf("postgres.WithInstance: %v\n", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"pgx", driver)
	if err != nil {
		log.Printf("migrate.NewWithDatabaseInstance: %v\n", err)
		return err
	}

	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(fmt.Errorf("migrate failed: %v", err))
		return err
	}

	return nil
}

// ShortLink struct which use in postgres operations: StoreURL, GetOriginalURLByShortURL.
type ShortLink struct {
	ID          int    `db:"id"`
	ShortURL    string `db:"short_url"`
	OriginalURL string `db:"original_url"`
	Deleted     bool   `db:"deleted"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

// StoreURL save original URL and short URL in databases.
func (db *DB) StoreURL(ctx context.Context, originalURL string, shortURL string) error {
	_, err := db.db.NamedExecContext(ctx, `INSERT INTO short_links (short_url, original_url)
	VALUES (:short_url, :original_url)`, &ShortLink{ShortURL: shortURL, OriginalURL: originalURL})
	if err != nil {
		return err
	}

	return nil
}

// GetOriginalURLByShortURL get original URL by short URL.
func (db *DB) GetOriginalURLByShortURL(ctx context.Context, shortURL string) (string, error) {
	var shortLink ShortLink
	err := db.db.GetContext(ctx, &shortLink, "SELECT * FROM short_links WHERE short_url=$1", shortURL)
	if err == sql.ErrNoRows {
		return "", sql.ErrNoRows
	}

	if shortLink.Deleted {
		return "", appErr.ErrURLDeleted
	}

	return shortLink.OriginalURL, nil
}

// GetAllURLs get all URLs from databases.
func (db *DB) GetAllURLs(ctx context.Context) (map[string]string, error) {
	var shortLink []ShortLink
	err := db.db.SelectContext(ctx, &shortLink, "SELECT * FROM short_links")
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(shortLink))
	for _, v := range shortLink {
		res[v.ShortURL] = v.OriginalURL
	}

	return res, nil
}

// GetShortURLByOriginalURL get URL by original URL
func (db *DB) GetShortURLByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	var shortLink ShortLink
	err := db.db.GetContext(ctx, &shortLink, "SELECT * FROM short_links WHERE original_url=$1", originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	return shortLink.ShortURL, appErr.ErrOriginalURLExist
}

// Ping check connect with databases.
func (db *DB) Ping(ctx context.Context) error {
	err := db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteURLByShortURL delete long URL by short URL.
func (db *DB) DeleteURLByShortURL(ctx context.Context, shortURLs string) error {
	_, err := db.db.ExecContext(ctx, `UPDATE short_links SET deleted=TRUE WHERE short_url=$1`, shortURLs)
	if err != nil {
		return err
	}

	return nil
}
