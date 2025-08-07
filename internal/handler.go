package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/internal/auth"
	"github.com/shinplay/internal/user"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	AuthHandler *auth.AuthHandler
	UserHandler *user.UserHandler
}

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "API is running",
	})
}
