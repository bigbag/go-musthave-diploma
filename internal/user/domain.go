package user

import (
	"encoding/hex"
	"errors"
)

var (
	ErrLoginAlreadyExist   = errors.New("user already exist")
	ErrUserNotFound        = errors.New("user not exist")
	ErrNotValidCredentials = errors.New("not valid users credentials")
)

type RequestUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *RequestUser) HexPassword() string {
	return hex.EncodeToString([]byte(u.Password))
}

type User struct {
	ID       int
	Login    string
	Password string
}
