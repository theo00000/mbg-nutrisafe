package routes

import (
	allergyCtrl "back-end/app/controllers/allergy"
	schoolCtrl "back-end/app/controllers/school"
	"back-end/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupSchoolRoutes(api fiber.Router) {
	school := api.Group("/school", middleware.SchoolOnly)

	school.Get("/profile", schoolCtrl.GetSchoolProfile)
	school.Put("/profile", schoolCtrl.UpdateSchoolProfile)

	school.Post("/allergy-data", allergyCtrl.CreateAllergy)
	school.Get("/allergy-data", allergyCtrl.GetAllergies)

	school.Get("/sppg-info", schoolCtrl.GetAssignedsppg)
	school.Post("/reports", schoolCtrl.CreateFoodReport)
	school.Get("/reports", schoolCtrl.GetFoodReports)

	school.Get("/allergy-menu", schoolCtrl.GetAllergyMenu)

	school.Get("/classes", schoolCtrl.GetClasses)
	school.Post("/classes", schoolCtrl.AddClass)
	school.Delete("/classes/:id", schoolCtrl.DeleteClass)

	school.Post("/students", schoolCtrl.AddStudent)
	school.Get("/students", schoolCtrl.GetStudents)
	school.Put("/students/:id", schoolCtrl.UpdateStudent)
	school.Delete("/students/:id", schoolCtrl.DeleteStudent)
}
