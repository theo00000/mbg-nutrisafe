package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendSppgCredentials(toEmail, sppgName, loginEmail, password string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" {
		return fmt.Errorf("konfigurasi SMTP belum lengkap")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Akun Login NutriSafe SPPG - " + sppgName
	body := fmt.Sprintf(
		"Halo!\r\n\r\n"+
			"Pendaftaran SPPG \"%s\" di NutriSafe telah disetujui.\r\n"+
			"Berikut adalah informasi akun login Anda:\r\n\r\n"+
			"Email Login : %s\r\n"+
			"Password    : %s\r\n\r\n"+
			"Silakan login di aplikasi NutriSafe menggunakan credentials di atas.\r\n"+
			"Demi keamanan, segera ganti password setelah login pertama kali.\r\n\r\n"+
			"Salam,\r\nTim NutriSafe",
		sppgName, loginEmail, password,
	)

	msg := fmt.Sprintf(
		"From: NutriSafe <%s>\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, body,
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
}

type StudentAccountInfo struct {
	Name       string
	LoginEmail string
	Password   string
}

func SendStudentAccountList(toEmail, schoolName string, accounts []StudentAccountInfo) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" {
		return fmt.Errorf("konfigurasi SMTP belum lengkap")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	var list strings.Builder
	for i, acc := range accounts {
		fmt.Fprintf(&list, "%d. %s\r\n   Email Login : %s\r\n   Password    : %s\r\n\r\n",
			i+1, acc.Name, acc.LoginEmail, acc.Password)
	}

	subject := fmt.Sprintf("Daftar Akun Siswa NutriSafe - %s", schoolName)
	body := fmt.Sprintf(
		"Halo %s!\r\n\r\n"+
			"Masa pengisian data siswa selama 7 hari telah selesai.\r\n"+
			"Berikut adalah daftar akun login untuk masing-masing siswa:\r\n\r\n"+
			"%s"+
			"Harap sampaikan informasi akun ini kepada masing-masing siswa.\r\n\r\n"+
			"Salam,\r\nTim NutriSafe",
		schoolName, list.String(),
	)

	msg := fmt.Sprintf(
		"From: NutriSafe <%s>\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, body,
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
}

func SendPasswordResetCode(toEmail, code string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" {
		return fmt.Errorf("konfigurasi SMTP belum lengkap")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Kode Reset Password NutriSafe"
	body := fmt.Sprintf(
		"Halo!\r\n\r\n"+
			"Kami menerima permintaan untuk mereset password akun NutriSafe Anda.\r\n"+
			"Berikut adalah kode verifikasi Anda:\r\n\r\n"+
			"  %s\r\n\r\n"+
			"Kode ini berlaku selama 15 menit.\r\n"+
			"Jika Anda tidak meminta reset password, abaikan email ini.\r\n\r\n"+
			"Salam,\r\nTim NutriSafe",
		code,
	)

	msg := fmt.Sprintf(
		"From: NutriSafe <%s>\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, body,
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
}

func SendSchoolCredentials(toEmail, schoolName, loginEmail, password string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" {
		return fmt.Errorf("konfigurasi SMTP belum lengkap")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Akun Login NutriSafe - " + schoolName
	body := fmt.Sprintf(
		"Halo!\r\n\r\n"+
			"Pendaftaran sekolah \"%s\" di NutriSafe telah berhasil.\r\n"+
			"Berikut adalah informasi akun login Anda:\r\n\r\n"+
			"Email Login : %s\r\n"+
			"Password    : %s\r\n\r\n"+
			"Silakan login di aplikasi NutriSafe menggunakan credentials di atas.\r\n"+
			"Demi keamanan, segera ganti password setelah login pertama kali.\r\n\r\n"+
			"Salam,\r\nTim NutriSafe",
		schoolName, loginEmail, password,
	)

	msg := fmt.Sprintf(
		"From: NutriSafe <%s>\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, body,
	)

	return smtp.SendMail(host+":"+port, auth, from, []string{toEmail}, []byte(msg))
}
