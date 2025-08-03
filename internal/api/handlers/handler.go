package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	AuthHandler *AuthHandler
}

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "API is running",
	})
}
