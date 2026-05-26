package routes

import (
	"back-end/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/ping", controllers.Ping)

	api.Post("/register", controllers.Register)
	api.Post("/register/spgg", controllers.RegisterSPGG)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)

	api.Get("/profile/detail", controllers.GetProfile)
	api.Put("/profile/detail", controllers.UpdateProfile)

	api.Get("/dashboard/stats", controllers.GetPublicStats)

	api.Post("/allergies", controllers.CreateAllergy)
	api.Get("/allergies", controllers.GetAllergies)

	admin := api.Group("/admin")
	admin.Get("/dashboard", controllers.GetAdminDashboard)

	admin.Get("/registrations", controllers.GetRegistrations)
	admin.Get("/registrations/:id", controllers.GetRegistrationDetail)
	admin.Patch("/registrations/:id/status", controllers.UpdateRegistrationStatus)

	admin.Get("/schools", controllers.GetSchools)
	admin.Get("/schools/:id", controllers.GetSchoolDetail)

	admin.Get("/users", controllers.GetUsers)
	admin.Put("/users/:id", controllers.AdminUpdateUser)

	admin.Get("/reports", controllers.GetReports)
	admin.Get("/reports/:id", controllers.GetReportDetail)
	admin.Patch("/reports/:id/status", controllers.UpdateReportStatus)

	admin.Get("/distributions", controllers.GetDistributions)
	admin.Patch("/distributions/:id/status", controllers.UpdateDistributionStatus)
}