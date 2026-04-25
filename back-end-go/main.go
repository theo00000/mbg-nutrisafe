package main

import (
	"log"
	"back-end/app/routes"
	"back-end/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDB()

	app := fiber.New()

	app.Use(logger.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}