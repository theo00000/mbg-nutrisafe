package controllers

import "github.com/gofiber/fiber/v2"

// Ping berfungsi untuk mengecek apakah server berjalan dengan baik
func Ping(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Sistem API MBG berjalan lancar!",
	})
}