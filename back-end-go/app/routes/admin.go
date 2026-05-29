package routes

import (
	adminCtrl "back-end/app/controllers/admin"
	"back-end/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupAdminRoutes(api fiber.Router) {
	route := api.Group("/admin", middleware.AdminOnly)

	route.Get("/dashboard", adminCtrl.GetAdminDashboard)

	route.Get("/registrations", adminCtrl.GetRegistrations)
	route.Get("/registrations/:id", adminCtrl.GetRegistrationDetail)
	route.Patch("/registrations/:id/status", adminCtrl.UpdateRegistrationStatus)

	route.Get("/schools", adminCtrl.GetSchools)
	route.Get("/schools/:id", adminCtrl.GetSchoolDetail)

	route.Get("/reports", adminCtrl.GetReports)
	route.Get("/reports/:id", adminCtrl.GetReportDetail)
	route.Patch("/reports/:id/status", adminCtrl.UpdateReportStatus)

	route.Get("/distributions", adminCtrl.GetDistributions)
	route.Patch("/distributions/:id/status", adminCtrl.UpdateDistributionStatus)
}
