package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetInformation(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("Welcome to Go microservices using Fiber")
}

func GetMukesh(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("Mukesh work faster")
}
