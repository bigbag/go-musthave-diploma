package order

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

func (r *Repository) Get(orderID string) (*Order, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	order := &Order{}

	sqlStatement := `SELECT id, user_id, amount, status, uploaded_at
						FROM orders WHERE id=$1;`
	row := r.conn.QueryRowContext(ctx, sqlStatement, orderID)
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.Amount,
		&order.Status,
		&order.UploadedAt,
	)

	if err != nil && err == sql.ErrNoRows {
		return &Order{}, nil
	}

	return order, err
}

func (r *Repository) GetAllByUserID(userID string) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	query := `SELECT id, amount, status, uploaded_at 
						FROM orders where user_id=$1 
						ORDER BY uploaded_at DESC;`
	rows, err := r.conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order *Order
	orders := make([]*Order, 0, 100)
	for rows.Next() {
		order = &Order{}
		err = rows.Scan(&order.ID, &order.Amount, &order.Status, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetTaskForChecking() ([]*Task, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	query := `SELECT id, user_id FROM orders where is_final=false;`
	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil && err == sql.ErrNoRows {
		return make([]*Task, 0), nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var task *Task
	tasks := make([]*Task, 0, 100)
	for rows.Next() {
		task = &Task{}
		err = rows.Scan(&task.OrderID, &task.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repository) CreateNew(userID string, orderID string) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	query := `INSERT INTO 
				orders(id, user_id, amount, uploaded_at, status) 
				VALUES($1, $2, $3, $4, $5);`
	if _, err := r.conn.ExecContext(
		ctx,
		query,
		orderID,
		userID,
		0,
		time.Now(),
		NEW,
	); err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == pgerrcode.UniqueViolation {
				return ErrAlreadyExist
			}
			return err
		}
	}

	return nil
}

func (r *Repository) UpdateOrder(order *Order) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	tx, err := r.conn.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	orderQuery := `UPDATE orders 
				SET status = $1, amount = $2, is_final = $3
				WHERE id = $4;`
	if _, err := tx.ExecContext(
		ctx, orderQuery, order.Status, order.Amount, order.IsFinal, order.ID,
	); err != nil {
		return err
	}

	walletQuery := `UPDATE wallets 
				SET balance = balance + $1
				WHERE user_id = $2;`
	if _, err := tx.ExecContext(
		ctx, walletQuery, order.Amount, order.UserID,
	); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
