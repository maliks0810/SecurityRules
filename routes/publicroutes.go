package routes

import (
	"SecurityRules/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	handler := handlers.NewSecurityRuleHandler()

	api := app.Group("/api/v1")

	rules := api.Group("/rules")
	rules.Get("/", handler.GetAll)
	rules.Get("/getSecurityExceptions", handler.GetSecurityException)
	rules.Get("/:id", handler.GetByID)
	rules.Post("/", handler.Create)
	rules.Post("/insertSecurityException", handler.InsertSecurityException)
	rules.Put("/:id", handler.Update)
	rules.Delete("/:id", handler.Delete)
}
