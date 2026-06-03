package routes

import (
	sppgCtrl "back-end/app/controllers/sppg"
	"back-end/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupsppgRoutes(api fiber.Router) {
	route := api.Group("/sppg", middleware.SppgOnly)

	route.Get("/profile", sppgCtrl.GetsppgProfile)
	route.Put("/profile", sppgCtrl.UpdatesppgProfile)

	route.Get("/status", sppgCtrl.GetPartnershipStatus)
	route.Get("/schools", sppgCtrl.GetsppgSchools)
	route.Get("/available-schools", sppgCtrl.GetAvailableSchools)
	route.Post("/assign-schools", sppgCtrl.AssignSchools)
	route.Get("/allergy-dashboard", sppgCtrl.GetAllergyDashboard)

	route.Post("/scan", sppgCtrl.ScanMenu)
	route.Post("/menu-reports", sppgCtrl.SubmitMenuReport)

	route.Get("/delivery-reports", sppgCtrl.GetDeliveryReports)
	route.Get("/delivery-reports/:id", sppgCtrl.GetDeliveryReportDetail)
	route.Post("/delivery-reports", sppgCtrl.CreateDeliveryReport)

	route.Post("/food-problems", sppgCtrl.CreateFoodProblem)
	route.Get("/food-reports", sppgCtrl.GetSchoolFoodReports)
}
