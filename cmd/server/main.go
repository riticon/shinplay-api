package main

import (
	"time"

	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/shinplay/internal/api/handlers"
	"github.com/shinplay/internal/auth"
	"github.com/shinplay/internal/config"
	"github.com/shinplay/internal/db"
	"github.com/shinplay/internal/user"
	"go.uber.org/dig"
)

func main() {
	cnf := config.GetConfig()

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

	app.Use(recover.New())
	app.Use(requestid.New())

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: cnf.Logger,
	}))

	app.Use(cors.New(cors.Config{ // CORS configuration
		AllowOrigins:     cnf.Server.CORS, // Explicitly allow development origin
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		ExposeHeaders:    "Content-Length, Content-Type",
		AllowCredentials: true,  // Allow credentials for development
		MaxAge:           86400, // 24 hours
	}))

	app.Use(helmet.New())   // Security headers
	app.Use(compress.New()) // Response compression
	app.Use(limiter.New(limiter.Config{
		Max:        100,             // Max number of requests
		Expiration: 1 * time.Minute, // Per minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Use IP as identifier
		},
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "API is running",
		})
	})

	// all the routes goes here
	err := container.Invoke(func(r handlers.Routes) {
		api := app.Group("/api")
		api.Post("/auth/whatsapp/send-otp", r.AuthHandler.SendWhatsAppOTP)
	})

	if err != nil {
		panic(err)
	}

	cnf.Logger.Info("Starting Shinplay API...")

	app.Listen(":" + cnf.Server.Port)

}
