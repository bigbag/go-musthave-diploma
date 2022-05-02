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

	sqlStatement := `SELECT id, amount, status, uploaded_at 
					FROM orders where user_id=$1 
					ORDER BY uploaded_at DESC;`
	rows, err := r.conn.QueryContext(ctx, sqlStatement, userID)
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

func (r *Repository) GetAllForChecking() ([]string, error) {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()

	sqlStatement := `SELECT id FROM orders where is_final=false;`
	rows, err := r.conn.QueryContext(ctx, sqlStatement)
	if err != nil && err == sql.ErrNoRows {
		return make([]string, 0), nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderID string
	orderIDs := make([]string, 0, 100)
	for rows.Next() {
		err = rows.Scan(&orderID)
		if err != nil {
			return nil, err
		}
		orderIDs = append(orderIDs, orderID)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return orderIDs, nil
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
				return ErrOrderAlreadyExist
			}
			return err
		}
	}

	return nil
}

func (r *Repository) Update(
	orderID string,
	status string,
	amount float64,
	isFinal bool,
) error {
	ctx, cancel := context.WithTimeout(r.ctx, r.connTimeout)
	defer cancel()
	query := `UPDATE orders 
			SET status = $1, amount = $2, is_final = $3
			WHERE id = $4;`
	if _, err := r.conn.ExecContext(
		ctx,
		query,
		status,
		amount,
		isFinal,
		orderID,
	); err != nil {
		if err, ok := err.(*pq.Error); ok {
			return err
		}
	}

	return nil
}
