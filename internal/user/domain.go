package user

import (
	"encoding/hex"
	"errors"
)

var (
	ErrAlreadyExist        = errors.New("user already exist")
	ErrUserNotFound        = errors.New("user not exist")
	ErrNotValidCredentials = errors.New("not valid users credentials")
)

type RequestUser struct {
	ID       string `json:"login"`
	Password string `json:"password"`
}

func (u *RequestUser) HexPassword() string {
	return hex.EncodeToString([]byte(u.Password))
}

type User struct {
	ID       string
	Password string
}
