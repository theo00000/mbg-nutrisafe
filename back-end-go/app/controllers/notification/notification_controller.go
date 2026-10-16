package notification

import (
	"strconv"
	"time"

	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"back-end/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type notificationRow struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	RefType   string    `json:"ref_type"`
	RefID     uint      `json:"ref_id"`
	IsRead    bool      `json:"is_read"`
	TimeLabel string    `json:"time_label"` // siap tampil: "12:00" atau "1 Februari 2026"
	CreatedAt time.Time `json:"created_at"`
}

// timeLabel meniru gaya daftar notifikasi di aplikasi: notifikasi hari ini
// cukup jamnya, yang lebih lama pakai tanggal penuh.
func timeLabel(t time.Time) string {
	now := time.Now()
	y1, m1, d1 := t.Date()
	y2, m2, d2 := now.Date()
	if y1 == y2 && m1 == m2 && d1 == d2 {
		return t.Format("15:04")
	}
	return utils.TanggalID(t)
}

// GetNotifications mengembalikan notifikasi milik user yang sedang login,
// terbaru lebih dulu. Query: ?unread_only=true&limit=20&offset=0
func GetNotifications(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	limit := defaultLimit
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		limit = v
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	offset := 0
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v > 0 {
		offset = v
	}

	// Chain gorm tidak aman dipakai ulang untuk dua finisher (Count lalu Find):
	// kondisinya bisa menumpuk. Jadi rakit ulang query untuk tiap pemanggilan.
	unreadOnly := c.Query("unread_only") == "true"
	baseQuery := func() *gorm.DB {
		q := config.DB.Model(&models.Notification{}).Where("user_id = ?", userID)
		if unreadOnly {
			q = q.Where("read_at IS NULL")
		}
		return q
	}

	var total int64
	baseQuery().Count(&total)

	var notifications []models.Notification
	if err := baseQuery().Order("created_at desc").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil notifikasi."})
	}

	var unreadCount int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&unreadCount)

	rows := make([]notificationRow, 0, len(notifications))
	for _, n := range notifications {
		rows = append(rows, notificationRow{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			RefType:   n.RefType,
			RefID:     n.RefID,
			IsRead:    n.IsRead(),
			TimeLabel: timeLabel(n.CreatedAt),
			CreatedAt: n.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"data":         rows,
		"unread_count": unreadCount,
		"total":        total,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetUnreadCount dipakai badge pada tab Notifikasi.
func GetUnreadCount(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var unreadCount int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&unreadCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"unread_count": unreadCount,
	})
}

// MarkAsRead menandai satu notifikasi sudah dibaca. Update dibatasi pada
// baris milik user yang login, jadi ID milik orang lain tidak bisa disentuh.
func MarkAsRead(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "ID notifikasi tidak valid."})
	}

	now := time.Now()
	result := config.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ? AND read_at IS NULL", id, userID).
		Update("read_at", now)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memperbarui notifikasi."})
	}

	// RowsAffected 0 bisa berarti sudah dibaca sebelumnya — itu bukan error,
	// tapi ID yang memang tidak ada harus tetap 404.
	if result.RowsAffected == 0 {
		var exists int64
		config.DB.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Count(&exists)
		if exists == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Notifikasi tidak ditemukan."})
		}
	}

	var unreadCount int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND read_at IS NULL", userID).Count(&unreadCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"message":      "Notifikasi ditandai sudah dibaca.",
		"unread_count": unreadCount,
	})
}

// MarkAllAsRead dipakai tombol "tandai semua dibaca".
func MarkAllAsRead(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	result := config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Update("read_at", time.Now())
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal memperbarui notifikasi."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"message":      "Semua notifikasi ditandai sudah dibaca.",
		"marked":       result.RowsAffected,
		"unread_count": 0,
	})
}
