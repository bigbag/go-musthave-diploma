package order

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	l logrus.FieldLogger
	r *Repository
	w *Worker
}

func NewService(l logrus.FieldLogger, r *Repository, w *Worker) *Service {
	return &Service{l: l, r: r, w: w}
}

func (s *Service) CreateOrder(userID string, orderID string) error {
	order, err := s.r.Get(orderID)
	if err != nil {
		return err
	}

	if order.UserID != "" && order.UserID != userID {
		return ErrAlreadyCreatedOtherUser
	}
	if order.ID != "" {
		return ErrAlreadyExist
	}

	err = s.r.CreateNew(userID, orderID)
	if err != nil {
		return err
	}

	return s.w.Add(&Task{OrderID: orderID, UserID: userID})
}

func (s *Service) FetchUserOrders(userID string) ([]*ResponseOrder, error) {
	orders, err := s.r.GetAllByUserID(userID)
	if len(orders) == 0 {
		return nil, ErrNotFound
	}

	result := make([]*ResponseOrder, 0, 100)
	for _, order := range orders {
		result = append(
			result,
			NewResponseOrder(order),
		)
	}

	return result, err
}
