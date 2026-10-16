package utils

import (
	"os"
	"strings"
)

// PublicURL mengubah path penyimpanan file di disk menjadi URL yang bisa
// dibuka klien (mobile/web).
//
// Controller menyimpan foto dengan bentuk path yang tidak seragam:
//
//	"./uploads/reports/x.jpg"  (umum.CreatePublicReport)
//	"uploads/reports/x.jpg"    (hasil filepath.Join di controller lain)
//	"uploads\\reports\\x.jpg"  (kalau server jalan di Windows)
//
// Semuanya dinormalkan ke "<APP_URL>/uploads/reports/x.jpg".
// Nilai kosong tetap dikembalikan kosong supaya klien bisa membedakan
// "tidak ada foto" dari "foto gagal dimuat".
func PublicURL(storedPath string) string {
	p := strings.TrimSpace(storedPath)
	if p == "" {
		return ""
	}

	// Sudah berupa URL absolut (mis. kalau nanti pindah ke object storage).
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}

	p = strings.ReplaceAll(p, "\\", "/")
	p = strings.TrimPrefix(p, "./")
	p = "/" + strings.TrimPrefix(p, "/")

	return strings.TrimRight(AppURL(), "/") + p
}

// AppURL mengembalikan base URL publik backend. Dipakai PublicURL dan untuk
// menyusun tautan di email/notifikasi.
func AppURL() string {
	if base := strings.TrimSpace(os.Getenv("APP_URL")); base != "" {
		return base
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return "http://localhost:" + port
}
