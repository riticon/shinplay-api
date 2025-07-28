package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/internal/auth"
	"github.com/shinplay/internal/config"
	"github.com/shinplay/internal/db"
)

func main() {
	// ctx := context.Background()
	app := fiber.New(
		fiber.Config{
			AppName: "Shinplay API",
		},
	)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	api := app.Group("/api")

	auth.Routes("/auth", api)

	println("Starting server on port:", config.GetConfig().Server.Port)

	db, err := db.NewPostgresDB()
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	app.Listen(":" + config.GetConfig().Server.Port)
}
