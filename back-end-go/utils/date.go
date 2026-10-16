package utils

import (
	"fmt"
	"time"
)

var namaBulanID = [...]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

var namaHariID = map[time.Weekday]string{
	time.Sunday:    "Minggu",
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
}

// StartOfDay mengembalikan tengah malam pada zona waktu t.
//
// Jangan pakai t.Truncate(24*time.Hour) untuk ini: Truncate menghitung dari
// epoch UTC, jadi di WIB (+07) setiap waktu sebelum pukul 07:00 akan jatuh ke
// tanggal sebelumnya.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// Today adalah StartOfDay untuk waktu sekarang.
func Today() time.Time { return StartOfDay(time.Now()) }

// TanggalID memformat tanggal gaya Indonesia: "1 Februari 2026".
// Dipakai di badan notifikasi yang dibaca pengguna.
func TanggalID(t time.Time) string {
	return fmt.Sprintf("%d %s %d", t.Day(), namaBulanID[int(t.Month())-1], t.Year())
}

// TanggalHariID menambahkan nama hari: "Rabu, 4 Februari 2026".
// Dipakai header layar Penerimaan Makanan & Pelaporan Makanan.
func TanggalHariID(t time.Time) string {
	return fmt.Sprintf("%s, %s", namaHariID[t.Weekday()], TanggalID(t))
}
