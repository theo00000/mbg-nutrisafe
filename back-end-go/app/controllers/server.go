package controllers

import "github.com/gofiber/fiber/v2"

func Ping(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "API Nutrisafe successfully running!",
	})
}