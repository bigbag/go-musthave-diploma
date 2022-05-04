package order

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	l logrus.FieldLogger
	r *Repository
	p *TaskPool
}

func NewService(l logrus.FieldLogger, r *Repository, p *TaskPool) *Service {
	return &Service{l: l, r: r, p: p}
}

func (s *Service) CreateOrder(userID string, orderID string) error {
	order, err := s.r.Get(orderID)
	if err != nil {
		return err
	}

	if order.UserID != "" && order.UserID != userID {
		return ErrOrderAlreadyCreatedOtherUser
	}
	if order.ID != "" {
		return ErrOrderAlreadyExist
	}
	err = s.r.CreateNew(userID, orderID)
	if err != nil {
		return err
	}

	return s.p.Push(NewTask(orderID))
}

func (s *Service) FetchUserOrders(userID string) ([]*ResponseOrder, error) {
	orders, err := s.r.GetAllByUserID(userID)
	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
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
