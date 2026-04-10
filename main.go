package main

import (
	"SecurityRules/config"
	"SecurityRules/database"
	"SecurityRules/internal/middleware"
	"SecurityRules/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func connectPostgres() {
	cfg := config.Load()
	database.Connect(cfg)
	database.Migrate()
	log.Println("PostgreSQL connection established")
}

func connectSnowflake() {
	sfCfg := config.LoadSnowflake()
	database.ConnectSnowflake(sfCfg)
	log.Println("Snowflake connection established")
}

func main() {
	cfg := config.Load()

	//connectPostgres()
	//defer database.Close()

	connectSnowflake()
	defer database.CloseSnowflake()

	app := fiber.New(fiber.Config{
		AppName: "SecurityRules API",
	})

	middleware.Setup(app)
	routes.PublicRoutes(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	log.Printf("SecurityRules API starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
