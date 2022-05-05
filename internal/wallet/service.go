package wallet

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	l logrus.FieldLogger
	r *Repository
}

func NewService(l logrus.FieldLogger, r *Repository) *Service {
	return &Service{l: l, r: r}
}

func (s *Service) FetchWallet(userID string) (*ResponseWallet, error) {
	w, err := s.r.GetWallet(userID)
	if err != nil {
		return nil, err
	}
	return &ResponseWallet{Balance: w.Balance, Withdrawal: w.Withdrawal}, nil
}

func (s *Service) CreateWithdrawal(rw *RequestWithdrawal) error {
	return s.r.CreateWithdrawal(rw)
}

func (s *Service) FetchUserWithdrawals(userID string) ([]*ResponseWithdrawal, error) {
	ws, err := s.r.GetWithdrawalsByUserID(userID)
	if len(ws) == 0 {
		return nil, ErrWithdrawalsNotFound
	}

	result := make([]*ResponseWithdrawal, 0, 100)
	for _, w := range ws {
		result = append(result, NewResponseWithdrawal(w))
	}

	return result, err
}
