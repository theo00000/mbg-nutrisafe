# Backend (Go)
---

## 🛠️ Langkah-Langkah Setelah `git pull`

Ikuti langkah berikut secara berurutan setelah kamu menarik (*pull*) branch ini ke lokalmu:

### 1. Masuk ke Direktori Backend
Buka terminal dan arahkan ke folder backend:
```bash
cd back-end
```

### 2. Instalasi Dependensi
Jalankan perintah ini untuk mengunduh semua library yang dibutuhkan:
```bash
go mod tidy
```

### 3. Setup Environment Variables (`.env`)
File `.env` tidak ikut terunggah ke Git demi keamanan. Kamu harus membuat file baru bernama **`.env`** di root folder `back-end` dan isi dengan konfigurasi berikut (sesuaikan dengan password PostgreSQL di laptopmu masing-masing):

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=[isi_password_postgresql_kamu]
DB_NAME=nutrisafe
DB_PORT=5432
```

### 4. Siapkan Database
Buka aplikasi database manager kamu (HeidiSQL/pgAdmin/DBeaver), lalu buat database baru bernama:
**`nutrisafe`**

> **Catatan:** Kamu tidak perlu membuat tabel secara manual. Sistem ini sudah menggunakan **GORM AutoMigrate**, jadi tabel akan otomatis dibuat saat server dijalankan pertama kali.

### 5. Jalankan Server
Setelah database dibuat dan `.env` siap, jalankan server dengan perintah:
```bash
go run main.go
```

Jika muncul pesan **"Koneksi ke PostgreSQL berhasil!"**, artinya backend sudah siap digunakan.

---

## 📂 Struktur Proyek
Proyek ini mengikuti arsitektur yang terorganisir untuk memisahkan tanggung jawab:

```text
back-end/
├── app/
│   ├── controllers/   
│   ├── middleware/   
│   ├── models/       
│   ├── routes/     
│   └── services/     
├── config/         
├── utils/            
├── .env            
└── main.go         
```

---

## 🚦 Endpoint Uji Coba
Untuk memastikan API berjalan, kamu bisa mengakses:
* **URL:** `http://localhost:3000/api/ping`
* **Method:** `GET`
* **Expected Response:** `{"status": "success", "message": "Sistem API NutriSafe berjalan lancar!"}`

---

## 👥 Kontribusi
Jika kamu melakukan perubahan pada struktur database, jangan lupa untuk menambahkan model baru tersebut di file `config/database.go` pada bagian `AutoMigrate`.

```
Sebagai tambahan pengingat, pastikan kamu juga membuat file **`.gitignore`** (jika belum ada) di dalam folder `back-end` dan isi dengan kode berikut agar rahasia *database* kalian tetap aman:

```text
# Konfigurasi Environment Lokal
.env

# Build file Go (opsional, agar file exe tidak ikut ter-push)
*.exe
*.exe~
*.dll
*.so
*.dylib
```
