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
	"github.com/shinplay/internal/config"
	"go.uber.org/zap"
)

type Server struct {
	appConfig *config.Config
	app       *fiber.App
}

func CreateNewFiberApp() Server {
	config := config.GetConfig()

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
		AllowOrigins:     config.Server.CORS,
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

	return Server{
		appConfig: config,
		app:       app,
	}
}

func (s *Server) StartServer() {
	port := s.appConfig.Server.Port
	host := s.appConfig.Server.Host

	serverAddr := host + ":" + port
	s.appConfig.Logger.Info(
		"Starting server",
		zap.String("address", serverAddr),
		zap.String("environment", s.appConfig.Environment),
	)

	// Start server
	if err := s.app.Listen(serverAddr); err != nil {
		s.appConfig.Logger.Fatal("Error starting server", zap.Error(err))
	}
}

func (s *Server) Routes() {
	// Define your routes here
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Add more routes as needed
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}
