package routes

import (
	"back-end/app/controllers/umum"
	"back-end/app/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupUmumRoutes(api fiber.Router) {
	route := api.Group("/umum", middleware.UmumOnly)

	route.Get("/profile", umum.GetUmumProfile)
	route.Put("/profile", umum.UpdateUmumProfile)

	route.Get("/schools", umum.GetSchoolsList)
	route.Get("/siswa-school", umum.GetSiswaSchool)
	route.Post("/reports", umum.CreatePublicReport)
}
