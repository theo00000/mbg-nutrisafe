package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// parseClaims memvalidasi bearer token dan mengembalikan claims-nya.
// Metode penandatanganan dikunci ke HMAC, sama seperti middleware per-role,
// supaya token dengan alg lain tidak pernah sampai ke tahap verifikasi.
func parseClaims(c *fiber.Ctx) (jwt.MapClaims, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("token tidak ditemukan")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("format token tidak valid")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("konfigurasi keamanan server tidak lengkap")
	}
	token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode token tidak didukung")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("gagal membaca claims")
	}
	return claims, nil
}

func GetUserIDFromToken(c *fiber.Ctx) (uint, error) {
	claims, err := parseClaims(c)
	if err != nil {
		return 0, err
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("id tidak valid")
	}
	return uint(userIDFloat), nil
}

// GetRoleFromToken mengembalikan role pada token aktif. Dipakai endpoint yang
// dipakai bersama lintas role tapi perlu menyesuaikan isi respons.
func GetRoleFromToken(c *fiber.Ctx) (string, error) {
	claims, err := parseClaims(c)
	if err != nil {
		return "", err
	}
	role, _ := claims["role"].(string)
	return role, nil
}

// Authenticated hanya memastikan token valid, tanpa membatasi role.
// Dipakai fitur yang muncul di semua dashboard (mis. tab Notifikasi).
func Authenticated(c *fiber.Ctx) error {
	if _, err := parseClaims(c); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	return c.Next()
}
