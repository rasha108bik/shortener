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

type psg struct {
	Psg *sqlx.DB
}

func NewPostgres(connString string) (*psg, error) {
	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = migrateUP(db)
	if err != nil {
		return nil, err
	}

	return &psg{
		Psg: db,
	}, nil
}

func (p *psg) Close() error {
	return p.Psg.Close()
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
	}

	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(fmt.Errorf("migrate failed: %v", err))
	}

	return nil
}

type ShortLink struct {
	ID          int    `db:"id"`
	ShortURL    string `db:"short_url"`
	OriginalURL string `db:"original_url"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func (p *psg) StoreURL(originalURL string, shortURL string) error {
	_, err := p.Psg.NamedExec(`INSERT INTO short_links (short_url, original_url)
	VALUES (:short_url, :original_url)`, &ShortLink{ShortURL: shortURL, OriginalURL: originalURL})
	if err != nil {
		return err
	}

	return nil
}

func (p *psg) GetOriginalURLByShortURL(shortURL string) (string, error) {
	var shortLink ShortLink
	err := p.Psg.Get(&shortLink, "SELECT * FROM short_links WHERE short_url=$1", shortURL)
	if err == sql.ErrNoRows {
		return "", sql.ErrNoRows
	}

	return shortLink.OriginalURL, nil
}

func (p *psg) GetAllURLs() (map[string]string, error) {
	var shortLink []ShortLink
	err := p.Psg.Select(&shortLink, "SELECT * FROM short_links")
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(shortLink))
	for _, v := range shortLink {
		res[v.ShortURL] = v.OriginalURL
	}

	return res, nil
}

func (p *psg) GetShortURLByOriginalURL(originalURL string) (string, error) {
	var shortLink ShortLink
	err := p.Psg.Get(&shortLink, "SELECT * FROM short_links WHERE original_url=$1", originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	return shortLink.ShortURL, appErr.ErrOriginalURLExist
}

func (p *psg) Ping(ctx context.Context) error {
	err := p.Psg.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
