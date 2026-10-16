package routes

import (
	notifCtrl "back-end/app/controllers/notification"
	"back-end/app/middleware"

	"github.com/gofiber/fiber/v2"
)

// Tab Notifikasi muncul di dashboard sekolah, SPPG, dan umum/siswa, jadi
// rutenya tidak dibatasi role — cukup token valid. Setiap query sudah
// disaring per user_id di controller.
func setupNotificationRoutes(api fiber.Router) {
	route := api.Group("/notifications", middleware.Authenticated)

	route.Get("/", notifCtrl.GetNotifications)
	route.Get("/unread-count", notifCtrl.GetUnreadCount)
	route.Patch("/read-all", notifCtrl.MarkAllAsRead)
	route.Patch("/:id/read", notifCtrl.MarkAsRead)
}
