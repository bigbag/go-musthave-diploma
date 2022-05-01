package order

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/bigbag/go-musthave-diploma/internal/config"
)

type orderHandler struct {
	log logrus.FieldLogger
	cfg *config.Config
}

func NewOrderHandler(
	r fiber.Router,
	l logrus.FieldLogger,
	cfg *config.Config,

) {
	handler := &orderHandler{log: l, cfg: cfg}

	r.Post("orders", handler.saveOrder)
	r.Get("orders", handler.getOrders)
}

func (h *orderHandler) saveOrder(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)
	return c.Status(fiber.StatusOK).SendString(userID)
}

func (h *orderHandler) getOrders(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("")
}
