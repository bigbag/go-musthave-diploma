package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	ctx         context.Context
	conn        *sql.DB
	connTimeout time.Duration
}

func NewRepository(
	ctx context.Context,
	databaseDSN string,
	connTimeout time.Duration,
) (*Repository, error) {
	r := &Repository{
		ctx:         ctx,
		connTimeout: connTimeout * time.Second,
	}

	conn, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return r, err
	}

	r.conn = conn
	return r, nil
}

func (r *Repository) GetConnect() *sql.DB {
	return r.conn
}

func (r *Repository) Status() error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	return r.conn.PingContext(ctx)
}

func (r *Repository) Close() error {
	return r.conn.Close()
}
