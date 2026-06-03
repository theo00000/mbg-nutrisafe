package routes

import (
	"back-end/app/controllers"
	authCtrl "back-end/app/controllers/auth"
	dashCtrl "back-end/app/controllers/dashboard"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	route := app.Group("")

	route.Get("/ping", controllers.Ping)

	route.Post("/register", authCtrl.Register)
	route.Post("/register/school", authCtrl.RegisterSchool)
	route.Post("/register/sppg", authCtrl.Registersppg)
	route.Post("/login", authCtrl.Login)
	route.Post("/logout", authCtrl.Logout)

	route.Post("/forgot-password/request", authCtrl.RequestPasswordReset)
	route.Post("/forgot-password/verify", authCtrl.VerifyPasswordReset)

	route.Get("/profile/detail", authCtrl.GetProfile)
	route.Put("/profile/detail", authCtrl.UpdateProfile)
	route.Post("/change-password", authCtrl.ChangePassword)

	route.Get("/dashboard/stats", dashCtrl.GetPublicStats)

	setupAdminRoutes(route)
	setupsppgRoutes(route)
	setupSchoolRoutes(route)
	setupUmumRoutes(route)
}
