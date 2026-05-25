package main

import (
	"back-end/app/routes"
	"back-end/config"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDB()

	app := fiber.New()
	
	app.Use(cors.New(cors.Config{
	AllowOrigins: "http://localhost:3001,http://127.0.0.1:3001",
	AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
}))

	app.Use(logger.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}