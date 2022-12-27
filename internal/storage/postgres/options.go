package postgres

import "time"

type Options func(*Postgres)

func MaxPoolSize(size int) Options {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}

func ConnAttempts(attempts int) Options {
	return func(c *Postgres) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Options {
	return func(c *Postgres) {
		c.connTimeout = timeout
	}
}
