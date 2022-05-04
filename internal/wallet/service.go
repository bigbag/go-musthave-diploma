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
