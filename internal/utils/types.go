package utils

import "github.com/gofiber/fiber/v2"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendJSONError(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(Error{code, msg})
}
