package order

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/theplant/luhn"

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

	r.Post("orders", handler.saveOrder)
	r.Get("orders", handler.getOrders)
}

func (h *handler) saveOrder(c *fiber.Ctx) error {
	orderID := string(c.Body())
	id, err := strconv.Atoi(orderID)
	if err != nil {
		return utils.SendJSONError(c, fiber.StatusBadRequest, "order id is invalid")
	}
	if !luhn.Valid(id) {
		return utils.SendJSONError(c, fiber.StatusUnprocessableEntity, "order id is invalid")
	}

	userID := c.Locals(h.cfg.Auth.ContextKey).(string)

	switch err = h.s.CreateOrder(userID, orderID); err {
	case ErrAlreadyExist:
		return c.Status(fiber.StatusOK).SendString("")
	case ErrAlreadyCreatedOtherUser:
		return utils.SendJSONError(c, fiber.StatusConflict, err.Error())
	case nil:
		return c.Status(fiber.StatusAccepted).SendString("")
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}

func (h *handler) getOrders(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)

	switch orders, err := h.s.FetchUserOrders(userID); err {
	case ErrNotFound:
		return utils.SendJSONError(c, fiber.StatusNoContent, err.Error())
	case nil:
		return c.Status(fiber.StatusOK).JSON(orders)
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}
