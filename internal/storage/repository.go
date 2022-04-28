package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type StorageRepository struct {
	ctx         context.Context
	conn        *sql.DB
	connTimeout time.Duration
}

func NewStorageRepository(
	ctx context.Context,
	databaseDSN string,
	connTimeout time.Duration,
) (*StorageRepository, error) {
	r := &StorageRepository{
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

func (r *StorageRepository) GetConnect() *sql.DB {
	return r.conn
}

func (r *StorageRepository) Status() error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	return r.conn.PingContext(ctx)
}

func (r *StorageRepository) Close() error {
	return r.conn.Close()
}
