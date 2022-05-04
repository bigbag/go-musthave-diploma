package userid

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type storage struct {
	ctx        *fiber.Ctx
	secret     string
	cookieName string
	userID     string
}

func NewStorage(ctx *fiber.Ctx, cookieName string, secret string) *storage {
	return &storage{
		ctx:        ctx,
		secret:     secret,
		cookieName: cookieName,
	}
}

func (s *storage) Get() (string, error) {
	token := s.ctx.Cookies(s.cookieName)
	if token == "" {
		return "", errors.New("missing or malformed JWT")
	}

	return s.getUserID(token)
}

func (s *storage) getUserID(tokenStr string) (string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.secret), nil
		},
	)

	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("user unauthorized")
	}

	return claims.Issuer, nil
}
