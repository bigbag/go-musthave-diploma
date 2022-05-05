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

func (r *Repository) GetWallet(userID string) (*Wallet, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	wallet := &Wallet{}

	query := `SELECT id, user_id, balance, withdrawal 
						FROM wallets WHERE user_id=$1;`
	row := r.conn.QueryRowContext(ctx, query, userID)
	err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Withdrawal)

	return wallet, err
}

func (r *Repository) CreateWithdrawal(rw *RequestWithdrawal) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	tx, err := r.conn.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	walletsQuery := `UPDATE wallets
						SET withdrawal = withdrawal + $1, balance = balance - $1
						WHERE user_id = $2;`
	if _, err := tx.ExecContext(
		ctx, walletsQuery, rw.Amount, rw.UserID,
	); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.CheckViolation {
				return ErrNotEnoughMoney
			}
		}
		return err
	}

	withdrawalsQuery := `INSERT INTO
						withdrawals(id, user_id, amount, processed_at)
						VALUES($1, $2, $3, $4);`
	if _, err := tx.ExecContext(
		ctx, withdrawalsQuery, rw.ID, rw.UserID, rw.Amount, time.Now(),
	); err != nil {
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
