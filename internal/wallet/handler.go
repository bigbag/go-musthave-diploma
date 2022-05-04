package wallet

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

	r.Get("", handler.getWallet)
	r.Post("withdraw", handler.createWithdrawal)
	r.Get("withdrawals", handler.getWithdrawals)
}

func (h *handler) getWallet(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)
	return c.Status(fiber.StatusOK).SendString(userID)
}

func (h *handler) createWithdrawal(c *fiber.Ctx) error {
	rw := new(RequestWithdrawal)
	if err := c.BodyParser(rw); err != nil {
		return utils.SendJSONError(
			c, fiber.StatusBadRequest, "Please specify a valid withdrawal parameters",
		)
	}

	orderID, err := strconv.Atoi(rw.ID)
	if err != nil {
		return utils.SendJSONError(
			c, fiber.StatusUnprocessableEntity, "order id is invalid",
		)
	}

	if !luhn.Valid(orderID) {
		return utils.SendJSONError(
			c, fiber.StatusUnprocessableEntity, "order id is invalid",
		)
	}

	rw.UserID = c.Locals(h.cfg.Auth.ContextKey).(string)

	switch err := h.s.CreateWithdrawal(rw); err {
	case ErrWithdrawalAlreadyExist:
		return utils.SendJSONError(c, fiber.StatusConflict, err.Error())
	case ErrWithdrawalsNotEnoughMoney:
		return utils.SendJSONError(c, fiber.StatusPaymentRequired, err.Error())
	case nil:
		return c.Status(fiber.StatusOK).SendString("")
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}

func (h *handler) getWithdrawals(c *fiber.Ctx) error {
	userID := c.Locals(h.cfg.Auth.ContextKey).(string)

	switch ws, err := h.s.FetchUserWithdrawals(userID); err {
	case ErrWithdrawalsNotFound:
		return utils.SendJSONError(c, fiber.StatusNoContent, err.Error())
	case nil:
		return c.Status(fiber.StatusOK).JSON(ws)
	default:
		return utils.SendJSONError(c, fiber.StatusInternalServerError, err.Error())
	}
}
