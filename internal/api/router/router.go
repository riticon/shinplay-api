package router

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
	"go.uber.org/zap"

	"github.com/shinplay/internal/config"
)

func CreateNewFiberApp(config *config.Config) *fiber.App {
	fiberOptions := fiber.Config{
		ReadTimeout: 30 * time.Second,
	}

	app := fiber.New(fiberOptions)

	app.Use(recover.New())
	app.Use(requestid.New())

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: config.Logger,
	}))

	app.Use(cors.New(cors.Config{ // CORS configuration
		AllowOrigins:     config.Env.CORS,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		ExposeHeaders:    "Content-Length, Content-Type",
		AllowCredentials: false,
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

	return app
}

func StartServer(app *fiber.App, config *config.Config) {
	port := config.Env.ServerPort
	host := config.Env.ServerHost

	serverAddr := host + ":" + port
	config.Logger.Info("Starting server", zap.String("address", serverAddr), zap.String("environment", config.Env.Environment))

	// Start server
	if err := app.Listen(serverAddr); err != nil {
		config.Logger.Fatal("Error starting server", zap.Error(err))
	}
}

func Routes(app *fiber.App) {
	// Define your routes here
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Add more routes as needed
}
