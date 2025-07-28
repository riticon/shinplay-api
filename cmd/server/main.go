package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/internal/auth"
	"github.com/shinplay/internal/config"
	"github.com/shinplay/internal/db"
)

func main() {
	app := fiber.New(
		fiber.Config{
			AppName: "Shinplay API",
		},
	)

	// health check route
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	api := app.Group("/api")
	auth.Routes("/auth", api)

	// initialize postgres database and auto migrate schema
	db.InitializeDatabase()

	app.Listen(":" + config.GetConfig().Server.Port)
}
