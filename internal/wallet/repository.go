package wallet

import (
	"context"
	"time"

	"database/sql"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	ctx         context.Context
	conn        *sql.DB
	connTimeout time.Duration
}

func NewRepository(
	ctx context.Context,
	l logrus.FieldLogger,
	conn *sql.DB,
	connTimeout time.Duration,
) *Repository {
	return &Repository{
		ctx:         ctx,
		connTimeout: connTimeout,
		conn:        conn,
	}
}
