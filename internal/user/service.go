package user

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	l  logrus.FieldLogger
	ur *UserRepository
}

func NewUserService(l logrus.FieldLogger, ur *UserRepository) *UserService {
	return &UserService{l: l, ur: ur}
}

func (us *UserService) Get(requestUser *RequestUser) (*User, error) {
	user, err := us.ur.Get(requestUser.Login)
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

func (us *UserService) Save(requestUser *RequestUser) (*User, error) {
	err := us.ur.Save(requestUser)
	if err != nil {
		return nil, err
	}

	return us.Get(requestUser)
}
