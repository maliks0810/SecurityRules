package main

import (
	"SecurityRules/config"
	"SecurityRules/database"
	"SecurityRules/middleware"
	"SecurityRules/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	defer database.Close()
	database.Migrate()

	app := fiber.New(fiber.Config{
		AppName: "SecurityRules API",
	})

	middleware.Setup(app)
	routes.Setup(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Printf("SecurityRules API starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
