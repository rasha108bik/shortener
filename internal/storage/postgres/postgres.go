package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// const (
// 	defaultMaxPoolSize  = 1
// 	defaultConnAttempts = 2
// 	defaultConnTimeout  = time.Second
// )

// type Postgres struct {
// 	maxPoolSize  int
// 	connAttempts int
// 	connTimeout  time.Duration

// 	Pool *pgxpool.Pool
// }

// func New(ctx context.Context, url string, opts ...Options) (*Postgres, error) {
// 	pg := &Postgres{
// 		maxPoolSize:  defaultConnAttempts,
// 		connAttempts: defaultConnAttempts,
// 		connTimeout:  defaultConnTimeout,
// 	}

// 	for _, opt := range opts {
// 		opt(pg)
// 	}

// 	poolConfig, err := pgxpool.ParseConfig(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
// 	}

// 	poolConfig.MaxConns = int32(pg.maxPoolSize)

// 	for pg.connAttempts > 0 {
// 		pg.Pool, err = pgxpool.ConnectConfig(ctx, poolConfig)
// 		if err == nil {
// 			break
// 		}

// 		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

// 		time.Sleep(pg.connTimeout)
// 		pg.connAttempts--
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
// 	}

// 	return pg, nil
// }

// func (p *Postgres) Close() {
// 	if p.Pool != nil {
// 		p.Pool.Close()
// 	}
// }

//go:generate mockgen -source=Postgres -package=$GOPACKAGE -destination=postgres_interface_mock.go

type Postgres interface {
	StoreURL(originalURL string, shortURL string) error
	GetOriginalURLByShortURL(shortURL string) (string, error)
	GetAllURLs() (map[string]string, error)
	Ping(ctx context.Context) error
}

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

func (p *postgres) Ping(ctx context.Context) error {
	err := p.Postgres.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
