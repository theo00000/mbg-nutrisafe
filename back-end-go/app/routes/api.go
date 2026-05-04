package routes

import (
	"back-end/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/ping", controllers.Ping)

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)

	api.Get("/profile/detail", controllers.GetProfile)
	api.Put("/profile/detail", controllers.UpdateProfile)

	api.Get("/dashboard/stats", controllers.GetPublicStats)
}