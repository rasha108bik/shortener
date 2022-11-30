package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

//go:generate bin/mockgen -source=postgres.go -package=$GOPACKAGE -destination=postgres_interface_mock.go

type postgres struct {
	Postgres *sqlx.DB
}

func New(dns string) (*postgres, error) {
	db, err := sqlx.Connect("pgx", dns)
	if err != nil {
		return nil, err
	}

	return &postgres{
		Postgres: db,
	}, nil
}

func (p *postgres) Close() {
	p.Postgres.Close()
}

type ShortLink struct {
	ID          int    `db:"id"`
	ShortURL    string `db:"short_url"`
	OriginalURL string `db:"original_url"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

func (p *postgres) StoreURL(originalURL string, shortURL string) error {
	_, err := p.Postgres.NamedExec(`INSERT INTO short_links (short_url, original_url)
	VALUES (:short_url, :original_url)`, &ShortLink{ShortURL: shortURL, OriginalURL: originalURL})
	if err != nil {
		return err
	}

	return nil
}

func (p *postgres) GetOriginalURLByShortURL(shortURL string) (string, error) {
	var shortLink ShortLink
	err := p.Postgres.Get(&shortLink, "SELECT * FROM short_links WHERE short_url=$1", shortURL)
	if err == sql.ErrNoRows {
		return "", sql.ErrNoRows
	}

	return shortLink.OriginalURL, nil
}

func (p *postgres) GetAllURLs() (map[string]string, error) {
	var shortLink []ShortLink
	err := p.Postgres.Select(&shortLink, "SELECT * FROM short_links")
	if err != nil {
		return nil, err
	}

	res := make(map[string]string, len(shortLink))
	for _, v := range shortLink {
		res[v.ShortURL] = v.OriginalURL
	}

	return res, nil
}

func (p *postgres) GetShortURLByOriginalURL(originalURL string) (string, error) {
	var shortLink ShortLink
	err := p.Postgres.Get(&shortLink, "SELECT * FROM short_links WHERE original_url=$1", originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", sql.ErrNoRows
		}
		return "", err
	}

	return shortLink.ShortURL, nil
}

func (p *postgres) Ping(ctx context.Context) error {
	err := p.Postgres.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
