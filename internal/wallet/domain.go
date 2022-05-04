package wallet

import (
	"errors"
	"time"
)

var (
	ErrWithdrawalAlreadyExist    = errors.New("withdrawal already exist")
	ErrWithdrawalsNotFound       = errors.New("withdrawals not found")
	ErrWithdrawalsNotEnoughMoney = errors.New("not enough money for withdrawal")
)

type RequestWithdrawal struct {
	ID     string  `json:"order"`
	UserID string  `json:"-"`
	Amount float64 `json:"sum"`
}

type Withdrawal struct {
	ID          string
	UserID      string
	Amount      float64
	ProcessedAt time.Time
}

type ResponseWithdrawal struct {
	ID          string  `json:"order"`
	Amount      float64 `json:"sum,omitempty"`
	ProcessedAt string  `json:"processed_at"`
}

func NewResponseWithdrawal(w *Withdrawal) *ResponseWithdrawal {
	return &ResponseWithdrawal{
		ID:          w.ID,
		Amount:      float64(w.Amount),
		ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
	}
}
