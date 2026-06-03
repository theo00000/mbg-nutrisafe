package school

import (
	"back-end/app/middleware"
	"back-end/app/models"
	"back-end/config"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetPendingDelivery mengembalikan pengiriman terbaru yang berstatus "pending"
// dari SPPG ke sekolah ini. Dipakai untuk mengisi form Penerimaan Makanan.
func GetPendingDelivery(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	var report models.DeliveryReport
	if err := config.DB.Where("school_id = ? AND status = ?", userID, "pending").
		Order("delivery_date desc").First(&report).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   nil,
		})
	}

	var sppgUser models.User
	config.DB.Select("id, name").First(&sppgUser, report.SppgID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"id":             report.ID,
			"sppg_id":        report.SppgID,
			"sppg_name":      sppgUser.Name,
			"delivery_date":  report.DeliveryDate.Format("2006-01-02"),
			"total_portions": report.TotalPortions,
			"status":         report.Status,
		},
	})
}

// ConfirmDelivery dipanggil sekolah untuk mengkonfirmasi makanan diterima.
// Form-data: received_portions (int, wajib), delivery_id (opsional - kalau
// tidak diisi pakai pengiriman pending terbaru), photo (opsional).
func ConfirmDelivery(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	receivedStr := c.FormValue("received_portions")
	received, err := strconv.Atoi(receivedStr)
	if err != nil || received < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "received_portions harus berupa angka >= 0."})
	}

	var report models.DeliveryReport
	q := config.DB.Where("school_id = ?", userID)
	if idStr := c.FormValue("delivery_id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "delivery_id tidak valid."})
		}
		q = q.Where("id = ?", id)
	} else {
		q = q.Where("status = ?", "pending").Order("delivery_date desc")
	}
	if err := q.First(&report).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": "Pengiriman tidak ditemukan."})
	}

	if photo, err := c.FormFile("photo"); err == nil {
		uploadDir := "./uploads/delivery-receipt"
		os.MkdirAll(uploadDir, os.ModePerm)
		filename := fmt.Sprintf("receipt_%d_%d_%s", userID, time.Now().UnixNano(), filepath.Base(photo.Filename))
		path := filepath.Join(uploadDir, filename)
		if err := c.SaveFile(photo, path); err == nil {
			if report.Photo2Path == "" {
				report.Photo2Path = path
			} else {
				report.Photo1Path = path
			}
		}
	}

	report.DistributedPortions = received
	report.Status = "diterima"
	if err := config.DB.Save(&report).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan konfirmasi."})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Konfirmasi penerimaan berhasil disimpan.",
		"data": fiber.Map{
			"id":                   report.ID,
			"status":               report.Status,
			"distributed_portions": report.DistributedPortions,
		},
	})
}
