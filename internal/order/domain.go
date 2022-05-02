package order

import (
	"errors"
	"time"
)

const (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)

var (
	ErrOrderAlreadyExist            = errors.New("order already exist")
	ErrOrderAlreadyCreatedOtherUser = errors.New("order already created by other user")
	ErrOrderNotFound                = errors.New("order not found")
	ErrOrdersNotFound               = errors.New("orders not found")
)

type Order struct {
	ID         string
	UserID     string
	Amount     float64
	Status     string
	UploadedAt time.Time
	IsChecked  bool
}

type ResponseOrder struct {
	ID         string  `json:"number"`
	Amount     float64 `json:"accrual,omitempty"`
	Status     string  `json:"status"`
	UploadedAt string  `json:"uploaded_at"`
}

func NewResponseOrder(order *Order) *ResponseOrder {
	return &ResponseOrder{
		ID:         order.ID,
		Amount:     float64(order.Amount),
		Status:     order.Status,
		UploadedAt: order.UploadedAt.Format(time.RFC3339),
	}
}
