package router

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/shinplay/internal/config"
)

func CreateNewFiberApp(config *config.Config) *fiber.App {
	fiberOptions := fiber.Config{
		ReadTimeout: 30 * time.Second,
	}

	app := fiber.New(fiberOptions)

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${ip} ${status} - ${method} ${path} ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// app.Use(cors.New(cors.Config{ // CORS configuration
	// 	AllowOrigins:     config.Env.CORS,
	// 	AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
	// 	AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Request-ID",
	// 	ExposeHeaders:    "Content-Length, Content-Type",
	// 	AllowCredentials: true,
	// 	MaxAge:           86400, // 24 hours
	// }))

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
	log.Printf("Starting server on %s in %s mode", serverAddr, config.Env.Environment)

	// Start server
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
