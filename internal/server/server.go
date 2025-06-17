package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"minion/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// StartServer запускает HTTP сервер
func StartServer(port string) error {
	// Создаем Fiber приложение
	app := fiber.New(fiber.Config{
		AppName: "🍌 Minion API v2.1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Printf("❌ Ошибка API: %v", err)

			return c.Status(code).JSON(handlers.APIResponse{
				Success: false,
				Message: "Внутренняя ошибка сервера",
				Error:   err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "🍌 ${time} | ${status} | ${latency} | ${ip} | ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// API Routes
	api := app.Group("/api")

	// Health check
	api.Get("/health", handlers.HealthCheck)

	// Configuration
	api.Get("/config", handlers.GetConfig)

	// Main operations
	api.Post("/extend-keys", handlers.ExtendKeys)
	api.Post("/refresh-menus", handlers.RefreshMenus)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🍌 BELLO! Minion API работает",
			"version": "2.1.0",
			"endpoints": []string{
				"GET  /api/health",
				"GET  /api/config",
				"POST /api/extend-keys",
				"POST /api/refresh-menus",
			},
		})
	})

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("🛑 Получен сигнал остановки, завершаем сервер...")
		if err := app.Shutdown(); err != nil {
			log.Printf("❌ Ошибка остановки сервера: %v", err)
		}
	}()

	// Запускаем сервер
	log.Printf("🚀 Запуск Minion API сервера на порту :%s\n", port)
	log.Println("📍 Доступные эндпоинты:")
	log.Println("   GET  /api/health")
	log.Println("   POST /api/extend-keys")
	log.Println("   POST /api/refresh-menus")

	return app.Listen(":" + port)
}
