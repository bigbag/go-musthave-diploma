package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	ctx         context.Context
	l           logrus.FieldLogger
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
		l:           l,
		connTimeout: connTimeout,
		conn:        conn,
	}
}

func (r *Repository) Get(userID string) (*User, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	user := &User{}

	query := `SELECT id, password FROM users WHERE id=$1;`
	row := r.conn.QueryRowContext(ctx, query, userID)
	err := row.Scan(&user.ID, &user.Password)

	return user, err
}

func (r *Repository) Save(user *RequestUser) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	tx, err := r.conn.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	usersQuery := `INSERT INTO users(id, password) VALUES($1, $2);`
	if _, err := tx.ExecContext(
		ctx, usersQuery, user.ID, user.HexPassword(),
	); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyExist
			}
		}
		return err
	}

	walletsQuery := `INSERT INTO wallets(user_id) VALUES($1);`
	if _, err := tx.ExecContext(ctx, walletsQuery, user.ID); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyExist
			}
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
