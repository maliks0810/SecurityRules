package routes

import (
	"SecurityRules/handlers"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(app *fiber.App) {
	handler := handlers.NewSecurityRuleHandler()

	api := app.Group("/api/v1")

	rules := api.Group("/rules")
	rules.Get("/getSecurityExceptions", handler.GetSecurityException)
	rules.Post("/insertSecurityException", handler.InsertSecurityException)

}
