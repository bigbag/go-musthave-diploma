package user

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

type Service struct {
	l logrus.FieldLogger
	r *Repository
}

func NewService(l logrus.FieldLogger, r *Repository) *Service {
	return &Service{l: l, r: r}
}

func (s *Service) Get(requestUser *RequestUser) (*User, error) {
	user, err := s.r.Get(requestUser.ID)
	if err != nil && err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	if user.Password != requestUser.HexPassword() {
		return nil, ErrNotValidCredentials
	}

	return user, nil
}

func (s *Service) Save(requestUser *RequestUser) (*User, error) {
	err := s.r.Save(requestUser)
	if err != nil {
		return nil, err
	}

	return s.Get(requestUser)
}
