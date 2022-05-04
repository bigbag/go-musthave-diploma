package wallet

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/config"
	// "github.com/bigbag/go-musthave-diploma/internal/utils"
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

	r.Get("", handler.getWallet)
	r.Post("withdraw", handler.createWithdrawal)
	r.Get("withdrawals", handler.getWithdrawals)
}

func (h *handler) getWallet(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)
	return c.Status(fiber.StatusOK).SendString(userID)
}

func (h *handler) createWithdrawal(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)
	return c.Status(fiber.StatusOK).SendString(userID)
}

func (h *handler) getWithdrawals(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)
	return c.Status(fiber.StatusOK).SendString(userID)
}
