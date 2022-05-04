package wallet

import (
	"context"
	"time"

	"database/sql"
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

func (r *Repository) CreateWithdrawal(rw *RequestWithdrawal) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	query := `INSERT INTO 
			withdrawals(id, user_id, amount, processed_at) 
			VALUES($1, $2, $3, $4);`

	if _, err := r.conn.ExecContext(
		ctx,
		query,
		rw.ID,
		rw.UserID,
		rw.Amount,
		time.Now(),
	); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return ErrWithdrawalAlreadyExist
			}
			return err
		}
	}

	return nil
}

func (r *Repository) GetWithdrawalsByUserID(userID string) ([]*Withdrawal, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	sqlStatement := `SELECT id, amount, processed_at 
					FROM withdrawals where user_id=$1 
					ORDER BY processed_at DESC;`
	rows, err := r.conn.QueryContext(ctx, sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var w *Withdrawal
	ws := make([]*Withdrawal, 0, 100)
	for rows.Next() {
		w = &Withdrawal{}
		err = rows.Scan(&w.ID, &w.Amount, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ws, nil
}
