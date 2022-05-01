package user

import (
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/config"
	"github.com/bigbag/go-musthave-diploma/internal/utils"
)

type UserHandler struct {
	log logrus.FieldLogger
	cfg *config.Config
	us  *UserService
}

func NewUserHandler(
	r fiber.Router,
	l logrus.FieldLogger,
	cfg *config.Config,
	us *UserService,

) {
	handler := &UserHandler{log: l, cfg: cfg, us: us}

	r.Post("register", handler.createUser)
	r.Post("login", handler.authUser)
}

func (h *UserHandler) saveAuthCookie(c *fiber.Ctx, userID int) error {
	expiresTime := time.Now().Add(time.Hour * h.cfg.Auth.ExpiresTime)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(userID),
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

func (h *UserHandler) authUser(c *fiber.Ctx) error {
	requestUser := new(RequestUser)
	if err := c.BodyParser(requestUser); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}
	if requestUser.Login == "" || requestUser.Password == "" {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}

	switch user, err := h.us.Get(requestUser); err {
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

func (h *UserHandler) createUser(c *fiber.Ctx) error {
	requestUser := new(RequestUser)
	if err := c.BodyParser(requestUser); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}
	if requestUser.Login == "" || requestUser.Password == "" {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid user parameters",
		)
	}

	switch user, err := h.us.Save(requestUser); err {
	case ErrLoginAlreadyExist:
		return utils.SendJSONError(c, fiber.StatusConflict, err.Error())
	case nil:
		h.saveAuthCookie(c, user.ID)
		return c.Status(fiber.StatusOK).SendString("")
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}
