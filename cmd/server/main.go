package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/shinplay/internal/api/handlers"
	"github.com/shinplay/internal/auth"
	"github.com/shinplay/internal/config"
	"github.com/shinplay/internal/db"
	"github.com/shinplay/internal/user"
	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	container.Provide(config.GetConfig)
	container.Provide(db.InitializeDatabase)

	container.Provide(user.NewUserRepository)
	container.Provide(user.NewUserService)

	container.Provide(auth.NewOTPRepository)
	container.Provide(auth.NewOTPService)

	container.Provide(auth.NewAuthService)
	container.Provide(handlers.NewAuthHandler)

	app := fiber.New(
		fiber.Config{
			AppName: "Shinplay API",
		},
	)

	err := container.Invoke(func(r handlers.Routes) {

		api := app.Group("/api")

		app.Get("/", func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": "API is running",
			})
		})

		api.Post("/auth/whatsapp/send-otp", r.AuthHandler.SendWhatsAppOTP)

	})

	if err != nil {
		panic(err)
	}

	logger := config.GetConfig().Logger
	logger.Info("Starting Shinplay API...")

	app.Listen(":" + config.GetConfig().Server.Port)

}
