package user

import (
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/config"
	"github.com/bigbag/go-musthave-diploma/internal/utils"
)

type handler struct {
	log logrus.FieldLogger
	cfg *config.Config
	s   *Service
}

func NewHandler(
	r fiber.Router,
	l logrus.FieldLogger,
	cfg *config.Config,
	s *Service,

) {
	handler := &handler{log: l, cfg: cfg, s: s}

	r.Post("register", handler.createUser)
	r.Post("login", handler.authUser)
}

func (h *handler) saveAuthCookie(c *fiber.Ctx, userID string) error {
	expiresTime := time.Now().Add(time.Hour * h.cfg.Auth.ExpiresTime)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    userID,
		ExpiresAt: expiresTime.Unix(),
	})

	token, err := claims.SignedString([]byte(h.cfg.Auth.SecretKey))
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     h.cfg.Auth.CookieName,
		Value:    token,
		Expires:  expiresTime,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return nil

}

func (h *handler) authUser(c *fiber.Ctx) error {
	requestUser := new(RequestUser)
	if err := c.BodyParser(requestUser); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}
	if requestUser.ID == "" || requestUser.Password == "" {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}

	switch user, err := h.s.Get(requestUser); err {
	case ErrUserNotFound:
		return utils.SendJSONError(c, fiber.StatusNotFound, err.Error())
	case ErrNotValidCredentials:
		return utils.SendJSONError(c, fiber.StatusUnauthorized, err.Error())
	case nil:
		h.saveAuthCookie(c, user.ID)
		return c.Status(fiber.StatusOK).SendString("")
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}

func (h *handler) createUser(c *fiber.Ctx) error {
	requestUser := new(RequestUser)
	if err := c.BodyParser(requestUser); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}
	if requestUser.ID == "" || requestUser.Password == "" {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}

	switch user, err := h.s.Save(requestUser); err {
	case ErrAlreadyExist:
		return utils.SendJSONError(c, fiber.StatusConflict, err.Error())
	case nil:
		h.saveAuthCookie(c, user.ID)
		return c.Status(fiber.StatusOK).SendString("")
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}
