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

func (r *Repository) Get(login string) (*User, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	user := &User{}

	sqlStatement := `SELECT id, login, password FROM users WHERE login=$1;`
	row := r.conn.QueryRowContext(ctx, sqlStatement, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)

	return user, err
}

func (r *Repository) Save(user *RequestUser) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	query := `INSERT INTO users(login, password) VALUES($1, $2);`
	if _, err := r.conn.ExecContext(ctx, query, user.Login, user.HexPassword()); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return ErrLoginAlreadyExist
			}
			return err
		}
	}

	return nil
}
