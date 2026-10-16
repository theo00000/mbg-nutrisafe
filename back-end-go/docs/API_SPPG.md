# API Documentation — SPPG (NutriSafe)

Dokumentasi lengkap alur SPPG, dari pendaftaran sampai pelaporan. Semua contoh siap di-copy ke Postman.

---

## 0. Informasi Umum

| Item | Nilai |
|---|---|
| Base URL (local) | `http://localhost:3000` |
| Base URL (produksi) | `https://api.nutrisafe.id` |
| Format response | JSON |
| Auth | JWT Bearer token, berlaku **24 jam** |
| Body limit | 10 MB per request (nginx produksi: 25 MB) |

### Akses Produksi

Di VPS, trafik publik masuk lewat **nginx host** (`deploy/nginx/nutrisafe.conf`), bukan langsung ke container:

| URL | Hasil |
|---|---|
| `https://api.nutrisafe.id` | ✅ Backend Go — **ini yang dipakai** |
| `http://api.nutrisafe.id` | `301` redirect ke HTTPS |
| `http://139.190.98.31:3000` | ❌ Unreachable — port 3000 di-bind ke `127.0.0.1` saja di compose |
| `http://139.190.98.31` | ❌ `444` — nginx `default_server` tolak akses via IP mentah |

Domain lain di VPS yang sama: `nutrisafe.id` (FE user, port 5173), `admin.nutrisafe.id` (FE admin, port 5174).

Timeout proxy nginx = **300s** untuk `api.nutrisafe.id`, jadi `POST /sppg/scan` yang lama tidak terputus.

### Bentuk Response

Sukses:
```json
{
  "status": "success",
  "message": "...",
  "data": { }
}
```

Error:
```json
{
  "status": "error",
  "message": "Pesan error dalam Bahasa Indonesia"
}
```

### Header Standar

| Endpoint jenis | Header |
|---|---|
| Public, body JSON | `Content-Type: application/json` |
| Public, body file | `Content-Type: multipart/form-data` *(Postman set otomatis — jangan diisi manual)* |
| Protected, body JSON | `Content-Type: application/json`<br>`Authorization: Bearer <token>` |
| Protected, body file | `Authorization: Bearer <token>` *(Content-Type otomatis)* |

> **Catatan Postman:** saat pakai `form-data`, **hapus** header `Content-Type` manual. Kalau diisi manual, boundary multipart hilang dan server balas `Format request tidak valid`.

### Variabel Postman yang Disarankan

Bikin Environment dengan 3 variabel, dipakai di semua request:

| Variable | Initial value |
|---|---|
| `base_url` | `http://localhost:3000` *(produksi: `https://api.nutrisafe.id`)* |
| `token` | *(kosong, diisi setelah login)* |
| `admin_token` | *(kosong, diisi setelah login admin)* |

Script di tab **Tests** request Login, biar token auto-tersimpan:
```javascript
const res = pm.response.json();
if (res.token) pm.environment.set("token", res.token);
```

---

## 1. Peta Alur SPPG

```
[1] POST /register/sppg              → SPPG kirim formulir (status: pending)
         ↓
[2] GET  /admin/registrations        → Admin lihat daftar pendaftar
[3] PATCH /admin/registrations/:id/status {"status":"disetujui"}
         ↓ sistem auto-buat akun user role=sppg
         ↓ kredensial dikirim ke email PIC
[4] POST /login                      → SPPG dapat JWT token
         ↓
[5] GET  /sppg/profile               → cek profil
[6] PUT  /sppg/profile               → lengkapi profil
[7] GET  /sppg/status                → cek status kemitraan
         ↓
[8] GET  /sppg/available-schools     → lihat sekolah belum punya SPPG
[9] POST /sppg/assign-schools        → klaim sekolah
[10] GET /sppg/schools               → daftar sekolah yang dilayani
[11] GET /sppg/allergy-dashboard     → rekap alergi siswa
         ↓
[12] POST /sppg/scan                 → scan foto menu (ML + cek alergen)
[13] POST /sppg/menu-reports         → simpan hasil scan jadi laporan gizi
         ↓
[14] POST /sppg/delivery-reports     → lapor pengiriman ke sekolah
[15] GET  /sppg/delivery-reports     → riwayat pengiriman
[16] GET  /sppg/delivery-reports/:id → detail pengiriman
         ↓
[17] POST /sppg/food-problems        → lapor masalah makanan (internal)
[18] GET  /sppg/food-reports         → lihat komplain dari sekolah
         ↓
[19] POST /change-password            → ganti password
[20] POST /logout                     → keluar
```

---

# BAGIAN A — Registrasi & Aktivasi Akun

## [1] Daftar Mitra SPPG

Pendaftaran SPPG **tidak** langsung bikin akun. Yang tersimpan adalah formulir pendaftaran status `pending`, menunggu verifikasi admin.

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/register/sppg` |
| Auth | Tidak perlu |
| Body type | **form-data** (karena ada upload file) |

### Header
```
(kosongkan — Postman isi otomatis untuk form-data)
```

### Body (form-data)

| Key | Type | Wajib | Aturan | Contoh |
|---|---|---|---|---|
| `registrant_name` | Text | ✅ | Nama PIC/pendaftar | `Budi Santoso` |
| `institution_name` | Text | ❌ | Nama lembaga/yayasan | `Yayasan Gizi Nusantara` |
| `nik_npwp` | Text | ✅ | **Tepat 16 digit** | `3201234567890001` |
| `email` | Text | ✅ | Email PIC, belum pernah dipakai daftar | `budi@gmail.com` |
| `phone` | Text | ✅ | Nomor HP | `081234567890` |
| `sppg_name` | Text | ✅ | Nama dapur SPPG | `SPPG Dapur Sehat Bandung` |
| `sppg_address` | Text | ✅ | Alamat lengkap dapur | `Jl. Merdeka No. 10, Bandung` |
| `production_capacity` | Text | ✅ | Angka > 0 (porsi/hari) | `1500` |
| `proposal` | **File** | ✅ | **Harus .pdf** | `proposal.pdf` |
| `kitchen_photo` | **File** | ✅ | Maks **5 MB** | `dapur.jpg` |

### Response `201 Created`
```json
{
  "status": "success",
  "message": "Formulir pendaftaran berhasil dikirim. Silahkan tunggu verifikasi selanjutnya melalui e-mail.",
  "data": {
    "id": 1,
    "email": "budi@gmail.com"
  }
}
```

### Kemungkinan Error

| Code | Message | Penyebab |
|---|---|---|
| 400 | `Semua kolom wajib (kecuali institution_name) harus diisi.` | Ada field kosong |
| 400 | `Kapasitas produksi harus berupa angka positif.` | `production_capacity` bukan angka atau ≤ 0 |
| 400 | `NIK harus terdiri dari 16 digit.` | Panjang `nik_npwp` ≠ 16 |
| 400 | `File proposal wajib diupload.` | Field `proposal` tidak ada |
| 400 | `File proposal harus berformat .pdf.` | Ekstensi bukan `.pdf` |
| 400 | `File foto dapur wajib diupload.` | Field `kitchen_photo` tidak ada |
| 400 | `Ukuran foto dapur tidak boleh melebihi 5 MB.` | File terlalu besar |
| 409 | `Email ini sudah pernah mendaftar sebagai mitra sppg.` | Email duplikat |

> **Simpan `data.id`** — dipakai admin di step [3].

---

## [2] Admin: Lihat Daftar Pendaftaran

Perlu login admin dulu (`POST /login` dengan akun role `admin` → simpan ke `{{admin_token}}`).

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/admin/registrations` |
| Query opsional | `?status=pending` |
| Auth | `Authorization: Bearer {{admin_token}}` |

Nilai `status` yang valid: `pending`, `belum_ditinjau`, `verifikasi`, `disetujui`.

### Response `200 OK`
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "registrant_name": "Budi Santoso",
      "institution_name": "Yayasan Gizi Nusantara",
      "email": "budi@gmail.com",
      "phone": "081234567890",
      "sppg_name": "SPPG Dapur Sehat Bandung",
      "sppg_address": "Jl. Merdeka No. 10, Bandung",
      "production_capacity": 1500,
      "status": "pending",
      "created_at": "2026-10-03T08:00:00Z",
      "updated_at": "2026-10-03T08:00:00Z"
    }
  ]
}
```

> `nik_npwp`, `proposal_path`, dan `kitchen_photo_path` **tidak** ikut di response (di-hide via `json:"-"`).

Detail satu pendaftaran: `GET {{base_url}}/admin/registrations/1`

---

## [3] Admin: Setujui Pendaftaran (Akun SPPG Dibuat di Sini)

| | |
|---|---|
| Method | `PATCH` |
| URL | `{{base_url}}/admin/registrations/1/status` |
| Auth | `Authorization: Bearer {{admin_token}}` |
| Body type | **raw → JSON** |

### Header
```
Content-Type: application/json
Authorization: Bearer {{admin_token}}
```

### Body
```json
{
  "status": "disetujui"
}
```

Status yang valid: `belum_ditinjau`, `verifikasi`, `disetujui`. Hanya `disetujui` yang memicu pembuatan akun.

### Response `200 OK`
```json
{
  "status": "success",
  "message": "Pendaftaran disetujui. Akun sppg berhasil dibuat. Kredensial login telah dikirim ke email PIC.",
  "data": {
    "registration": { "id": 1, "status": "disetujui" },
    "user_id": 12,
    "email": "sppg-dapur-sehat-bandung@sppg.id"
  }
}
```

### Yang Terjadi di Balik Layar

1. Login email di-generate dari slug `sppg_name` + domain `@sppg.id`.
   Contoh: `SPPG Dapur Sehat Bandung` → `sppg-dapur-sehat-bandung@sppg.id`.
   Kalau tabrakan, ditambah counter: `...sppg.id` → `....2@sppg.id`.
2. Password temporary = gabungan kata `sppg_name` (kapital di huruf awal, tanpa karakter non-alfanumerik) + 4 digit random.
   Contoh: `SPPGDapurSehatBandung4821`.
3. Record `users` dibuat dengan `role_name = "sppg"`.
4. Record `sppg_profiles` dibuat: `institution_name`, `capacity`, `address`, `contact_email` (email PIC), `status = "active"`.
5. Email kredensial dikirim **asinkron** ke email PIC.

> **Untuk testing manual:** email login ada di `data.email` response ini. Password temporary **tidak** ada di response — harus dicek di inbox email PIC, atau reset lewat `POST /forgot-password/request`.

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Request tidak valid, pilih: ['belum_ditinjau', 'verifikasi', 'disetujui'].` |
| 400 | `Status 'xxx' tidak valid. Pilihan yang tersedia: [...]` |
| 404 | `Data pendaftaran dengan ID 1 tidak ditemukan.` |
| 500 | `Status disetujui, tapi gagal membuat akun sppg: ...` |

---

## [4] Login SPPG

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/login` |
| Auth | Tidak perlu |
| Body type | **raw → JSON** |

### Header
```
Content-Type: application/json
```

### Body
```json
{
  "email": "sppg-dapur-sehat-bandung@sppg.id",
  "password": "SPPGDapurSehatBandung4821"
}
```

### Response `200 OK`
```json
{
  "status": "success",
  "message": "Login berhasil",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "role": "sppg"
}
```

Isi JWT: `user_id`, `email`, `role`, `exp` (+24 jam).

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Email dan kata sandi tidak boleh kosong!` |
| 401 | `Email atau kata sandi salah.` |
| 500 | `Terjadi kesalahan internal server (Konfigurasi keamanan tidak lengkap).` — `JWT_SECRET` belum di-set |

> **Mulai dari sini, semua endpoint `/sppg/*` wajib header `Authorization: Bearer {{token}}`.**
> Token non-SPPG ditolak `403` dengan message `Hanya mitra sppg yang dapat mengakses fitur ini`.

### Lupa Password (opsional)

**a. Minta kode**
```
POST {{base_url}}/forgot-password/request
Content-Type: application/json

{ "email": "budi@gmail.com" }
```
Boleh pakai email PIC (`contact_email`) **atau** login email (`...@sppg.id`) — sistem resolve keduanya. Response selalu `200 success` (anti email enumeration). Kode 6 digit dikirim ke inbox, berlaku **15 menit**, maks **5 kali salah**.

**b. Verifikasi + password baru**
```
POST {{base_url}}/forgot-password/verify
Content-Type: application/json

{
  "email": "budi@gmail.com",
  "code": "483201",
  "new_password": "PasswordBaru123"
}
```
`new_password` minimal **8 karakter**.

---

# BAGIAN B — Profil & Status Kemitraan

## [5] Lihat Profil SPPG

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/profile` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "data": {
    "id": 12,
    "name": "SPPG Dapur Sehat Bandung",
    "email": "sppg-dapur-sehat-bandung@sppg.id",
    "phone": "081234567890",
    "role": "sppg",
    "lokasi": "Jl. Merdeka No. 10, Bandung",
    "institution_name": "Yayasan Gizi Nusantara",
    "capacity": 1500,
    "sppg_status": "active"
  }
}
```

---

## [6] Update Profil SPPG

Semua field **opsional** — kirim hanya yang mau diubah. Field kosong/tidak dikirim diabaikan.

| | |
|---|---|
| Method | `PUT` |
| URL | `{{base_url}}/sppg/profile` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **raw → JSON** |

### Header
```
Content-Type: application/json
Authorization: Bearer {{token}}
```

### Body
```json
{
  "name": "SPPG Dapur Sehat Bandung Raya",
  "email": "dapursehat@sppg.id",
  "phone": "081298765432",
  "lokasi": "Jl. Merdeka No. 12, Bandung",
  "institution_name": "Yayasan Gizi Nusantara"
}
```

| Field | Target tabel | Catatan |
|---|---|---|
| `name` | `users.name` | — |
| `email` | `users.email` | Harus unik |
| `phone` | `users.phone` | Harus unik |
| `lokasi` | `sppg_profiles.address` | — |
| `institution_name` | `sppg_profiles.institution_name` | — |

### Response `200 OK`
```json
{
  "status": "success",
  "message": "Profil berhasil diperbarui.",
  "data": {
    "id": 12,
    "name": "SPPG Dapur Sehat Bandung Raya",
    "email": "dapursehat@sppg.id",
    "phone": "081298765432",
    "lokasi": "Jl. Merdeka No. 12, Bandung",
    "institution_name": "Yayasan Gizi Nusantara"
  }
}
```

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Format request tidak valid.` |
| 404 | `User tidak ditemukan.` |
| 409 | `Email sudah digunakan oleh pengguna lain.` |
| 409 | `Nomor telepon sudah digunakan oleh pengguna lain.` |

> **Penting:** kalau `email` diganti, token lama masih valid (klaim `email` di JWT tidak dipakai untuk lookup), tapi login berikutnya harus pakai email baru.

---

## [7] Cek Status Kemitraan

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/status` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "data": {
    "is_partner": true,
    "sppg_count": 1,
    "school_count": 3,
    "sppg_list": [
      {
        "id": 5,
        "name": "SPPG Dapur Sehat Bandung",
        "status": "active",
        "address": "Jl. Merdeka No. 10, Bandung",
        "capacity": 1500
      }
    ]
  }
}
```

| Field | Arti |
|---|---|
| `is_partner` | `true` kalau `sppg_profiles` sudah ada |
| `sppg_count` | Jumlah profil SPPG milik user ini |
| `school_count` | Jumlah sekolah yang sudah di-assign ke SPPG ini |
| `sppg_list[].id` | ID `sppg_profiles` (**bukan** `user_id`) |

---

# BAGIAN C — Manajemen Sekolah

## [8] Lihat Sekolah yang Tersedia

Daftar sekolah yang **belum punya SPPG** (`sppg_user_id IS NULL`). Tidak difilter per wilayah — semua SPPG lihat pool yang sama.

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/available-schools` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "data": [
    { "id": 7, "name": "SDN Merdeka 1", "grade": "SD/MI", "student_count": 240 },
    { "id": 8, "name": "SMPN 3 Bandung", "grade": "SMP/MTs/MTsN", "student_count": 410 }
  ]
}
```

> `data[].id` adalah **`users.id` sekolah** — ini yang dipakai di `school_ids` step [9] dan `school_id` step [14].

Kalau tidak ada sekolah tersedia, `data` = `[]`.

---

## [9] Klaim / Assign Sekolah

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/sppg/assign-schools` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **raw → JSON** |

### Header
```
Content-Type: application/json
Authorization: Bearer {{token}}
```

### Body
```json
{
  "school_ids": [7, 8]
}
```

| Field | Type | Wajib | Catatan |
|---|---|---|---|
| `school_ids` | array of integer | ✅ | Tidak boleh array kosong. Isi = `users.id` sekolah |

### Response `200 OK`
```json
{
  "status": "success",
  "message": "Sekolah berhasil ditambahkan."
}
```

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `school_ids tidak valid` — body gagal parse atau array kosong |
| 500 | `Gagal menetapkan sekolah.` |

> **Hati-hati:** update hanya kena sekolah yang `sppg_user_id IS NULL`. ID yang sudah diklaim SPPG lain **di-skip tanpa error** — response tetap `success`. Verifikasi hasil nyata lewat step [10].

---

## [10] Daftar Sekolah yang Dilayani

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/schools` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "data": {
    "total_school": 2,
    "total_students": 650,
    "schools": [
      {
        "id": 7,
        "name": "SDN Merdeka 1",
        "address": "Jl. Sekolah No. 1, Bandung",
        "Grade": "SD/MI",
        "teacher_count": 18,
        "student_count": 240,
        "allergy_types": ["Telur", "Kacang", "Susu"]
      },
      {
        "id": 8,
        "name": "SMPN 3 Bandung",
        "address": "Jl. Pelajar No. 3, Bandung",
        "Grade": "SMP/MTs/MTsN",
        "teacher_count": 32,
        "student_count": 410,
        "allergy_types": ["Seafood"]
      }
    ]
  }
}
```

> ⚠️ Key `Grade` di-serialize dengan **huruf kapital G** (beda dari endpoint lain yang pakai `grade`). Sesuaikan parsing di client.

---

## [11] Dashboard Alergi

Rekap alergi gabungan dari semua sekolah yang dilayani.

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/allergy-dashboard` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "data": {
    "total_allergy_types": 4,
    "total_allergy_students": 37,
    "allergy_breakdown": [
      { "allergy_type": "Telur", "count": 15 },
      { "allergy_type": "Kacang", "count": 12 },
      { "allergy_type": "Susu", "count": 6 },
      { "allergy_type": "Seafood", "count": 4 }
    ],
    "schools": [
      { "id": 7, "name": "SDN Merdeka 1", "total_students": 240, "allergy_count": 33 },
      { "id": 8, "name": "SMPN 3 Bandung", "total_students": 410, "allergy_count": 4 }
    ]
  }
}
```

Belum ada sekolah ter-assign → semua nilai `0` dan array kosong.

---

# BAGIAN D — Scan Menu & Laporan Gizi

## [12] Scan Foto Menu

Upload foto makanan → backend forward ke ML service (`ML_SERVICE_URL`, default `http://localhost:8000/predict`) → hasil nutrisi dicocokkan dengan data alergi siswa di sekolah yang dilayani.

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/sppg/scan` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **form-data** |
| Timeout | ML service dibatasi **90 detik** — set timeout Postman ≥ 100s |

### Header
```
Authorization: Bearer {{token}}
```
*(Content-Type otomatis)*

### Body (form-data)

| Key | Type | Wajib | Aturan | Contoh |
|---|---|---|---|---|
| `image` | **File** | ✅ | Foto menu (jpg/png) | `menu-ayam-goreng.jpg` |
| `total_portions` | Text | ✅ | Angka > 0 | `1500` |

### Response `200 OK`
```json
{
  "status": "success",
  "data": {
    "menu_detail": {
      "kalori": 750000,
      "protein": 30000,
      "karbohidrat": 105000,
      "lemak": 22500
    },
    "kalori_per_porsi": "500.0",
    "total_porsi": 1500,
    "status_gizi": "Baik",
    "bahan_terdeteksi": [
      { "name": "ayam", "is_allergen": false },
      { "name": "telur", "is_allergen": true },
      { "name": "nasi", "is_allergen": false }
    ],
    "peringatan_alergi": [
      {
        "allergen_type": "Telur",
        "school_id": 7,
        "school_name": "SDN Merdeka 1",
        "student_count": 15
      }
    ],
    "gemini_failed": false
  }
}
```

### Penjelasan Field Response

| Field | Arti |
|---|---|
| `menu_detail` | Nilai gizi mentah dari ML service |
| `kalori_per_porsi` | **String** hasil `kalori / total_porsi`, 1 desimal |
| `bahan_terdeteksi[].is_allergen` | `true` kalau bahan match (`ILIKE %bahan%`) dengan `allergy_type` siswa di sekolah yang dilayani |
| `peringatan_alergi` | Satu entri per kombinasi (alergen × sekolah) yang punya siswa terdampak |
| `gemini_failed` | `true` kalau ML service gagal/tidak terjangkau → data di atas adalah **fallback dummy** |

### Fallback Saat ML Down

Kalau ML service error atau timeout, request **tetap balas `200`** dengan nilai dummy:
```json
{
  "menu_detail": { "kalori": 500, "protein": 20, "karbohidrat": 70, "lemak": 15 },
  "status_gizi": "Baik",
  "bahan_terdeteksi": [],
  "peringatan_alergi": [],
  "gemini_failed": true
}
```
> **Selalu cek `gemini_failed`** sebelum percaya angkanya. `true` = data bukan hasil scan asli.

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `total_portions harus diisi dan bernilai positif.` |
| 400 | `Field 'image' wajib diisi dengan file gambar.` |
| 500 | `Gagal membaca file gambar.` / `Gagal memproses file gambar.` |

---

## [13] Simpan Laporan Menu

Simpan hasil scan step [12] jadi record permanen. Field diisi dari response scan.

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/sppg/menu-reports` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **raw → JSON** |

### Header
```
Content-Type: application/json
Authorization: Bearer {{token}}
```

### Body
```json
{
  "menu_name": "Ayam Goreng + Nasi + Sayur Bening",
  "kalori": 750000,
  "protein": 30000,
  "karbohidrat": 105000,
  "lemak": 22500,
  "total_porsi": 1500,
  "kalori_per_porsi": 500,
  "status_gizi": "Baik"
}
```

| Field | Type | Wajib | Catatan |
|---|---|---|---|
| `menu_name` | string | ❌ | Maks 200 karakter. Diinput manual oleh SPPG |
| `kalori` | number | ❌ | Dari `menu_detail.kalori` |
| `protein` | number | ❌ | Dari `menu_detail.protein` |
| `karbohidrat` | number | ❌ | Dari `menu_detail.karbohidrat` |
| `lemak` | number | ❌ | Dari `menu_detail.lemak` |
| `total_porsi` | integer | ✅ | **Harus > 0** — satu-satunya field yang divalidasi |
| `kalori_per_porsi` | number | ❌ | **Kirim sebagai angka**, bukan string. Response scan memberi string `"500.0"` → convert dulu |
| `status_gizi` | string | ❌ | Maks 50 karakter |

### Response `201 Created`
```json
{
  "status": "success",
  "message": "Laporan menu berhasil dikirim.",
  "data": {
    "id": 4,
    "sppg_id": 12,
    "menu_name": "Ayam Goreng + Nasi + Sayur Bening",
    "kalori": 750000,
    "protein": 30000,
    "karbohidrat": 105000,
    "lemak": 22500,
    "total_porsi": 1500,
    "kalori_per_porsi": 500,
    "status_gizi": "Baik",
    "report_date": "2026-10-03T00:00:00Z",
    "created_at": "2026-10-03T10:15:00Z",
    "updated_at": "2026-10-03T10:15:00Z"
  }
}
```

`report_date` di-set otomatis ke waktu server.

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Format body tidak valid.` |
| 400 | `total_porsi harus bernilai positif.` |
| 500 | `Gagal menyimpan laporan menu.` |

---

# BAGIAN E — Laporan Pengiriman

## [14] Buat Laporan Pengiriman

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/sppg/delivery-reports` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **form-data** (ada upload foto) |

### Header
```
Authorization: Bearer {{token}}
```

### Body (form-data)

| Key | Type | Wajib | Aturan | Contoh |
|---|---|---|---|---|
| `school_id` | Text | ✅ | `users.id` sekolah, > 0 | `7` |
| `delivery_date` | Text | ✅ | Format **`YYYY-MM-DD`** | `2026-10-03` |
| `total_portions` | Text | ✅ | Angka > 0 | `240` |
| `distributed_portions` | Text | ❌ | Default = `total_portions` | `235` |
| `menu_items` | Text | ❌ | **JSON string** (lihat bawah) | `[{...}]` |
| `photo1` | File | ❌ | Foto bukti 1 | `kirim-1.jpg` |
| `photo2` | File | ❌ | Foto bukti 2 | `kirim-2.jpg` |

### Format `menu_items`

Kirim sebagai **satu field Text** berisi JSON array yang di-stringify (bukan nested form-data):

```json
[
  {
    "menu_name": "Nasi Putih",
    "category": "Karbohidrat",
    "portions": 240,
    "is_allergy_substitute": false,
    "allergy_note": ""
  },
  {
    "menu_name": "Ayam Goreng",
    "category": "Protein Hewani",
    "portions": 225,
    "is_allergy_substitute": false,
    "allergy_note": ""
  },
  {
    "menu_name": "Tempe Goreng",
    "category": "Protein Nabati",
    "portions": 15,
    "is_allergy_substitute": true,
    "allergy_note": "pengganti ayam, tanpa telur"
  },
  {
    "menu_name": "Sayur Bening",
    "category": "Sayuran",
    "portions": 240,
    "is_allergy_substitute": false,
    "allergy_note": ""
  },
  {
    "menu_name": "Pisang",
    "category": "Buah",
    "portions": 240,
    "is_allergy_substitute": false,
    "allergy_note": ""
  }
]
```

| Field item | Type | Catatan |
|---|---|---|
| `menu_name` | string | Maks 100 karakter |
| `category` | string | Maks 50. Contoh: `Karbohidrat`, `Protein Hewani`, `Protein Nabati`, `Sayuran`, `Buah` |
| `portions` | integer | Jumlah porsi item ini |
| `is_allergy_substitute` | boolean | `true` kalau ini menu pengganti untuk siswa alergi |
| `allergy_note` | string | Maks 100. Contoh: `tanpa telur` |

### Response `201 Created`
```json
{
  "status": "success",
  "message": "Laporan pengiriman berhasil dikirim.",
  "data": {
    "id": 9,
    "school_id": 7,
    "delivery_date": "2026-10-03",
    "total_portions": 240,
    "status": "pending"
  }
}
```

Status awal selalu `pending`. Berubah jadi `diterima` setelah sekolah konfirmasi penerimaan (endpoint sisi sekolah).

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `school_id tidak valid.` |
| 400 | `Format tanggal tidak valid. Gunakan YYYY-MM-DD.` |
| 400 | `total_portions tidak valid.` |
| 400 | `Format menu_items tidak valid.` — JSON string salah |
| 500 | `Gagal menyimpan laporan.` / `Gagal menyimpan item menu.` |

> Simpan + item menu dijalankan dalam satu transaksi — kalau satu item gagal, seluruh laporan di-rollback.
> Upload foto **tidak** transaksional dan gagal secara silent: kalau save file error, path disimpan kosong tanpa error ke client.

---

## [15] Riwayat Laporan Pengiriman

Hasil dikelompokkan per tanggal pengiriman.

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/delivery-reports` |
| Query opsional | `?month=10&year=2026` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

| Query param | Type | Catatan |
|---|---|---|
| `month` | integer `1`–`12` | Tanpa leading zero juga boleh |
| `year` | integer 4 digit | Contoh `2026` |

Contoh: `{{base_url}}/sppg/delivery-reports?month=10&year=2026`

### Response `200 OK`
```json
{
  "status": "success",
  "data": [
    {
      "date": "2026-10-03",
      "schools": [
        {
          "report_id": 9,
          "school_id": 7,
          "school_name": "SDN Merdeka 1",
          "address": "Jl. Sekolah No. 1, Bandung",
          "total_portions": 240,
          "status": "pending"
        },
        {
          "report_id": 10,
          "school_id": 8,
          "school_name": "SMPN 3 Bandung",
          "address": "Jl. Pelajar No. 3, Bandung",
          "total_portions": 410,
          "status": "diterima"
        }
      ]
    }
  ]
}
```

> ⚠️ Grouping pakai Go map, jadi **urutan tanggal di array `data` tidak deterministik**. Sort di sisi client kalau butuh urut.

---

## [16] Detail Laporan Pengiriman

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/delivery-reports/9` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

`:id` = `report_id` dari step [15]. Hanya bisa akses laporan milik SPPG sendiri.

### Response `200 OK`
```json
{
  "status": "success",
  "data": {
    "id": 9,
    "school": {
      "id": 7,
      "name": "SDN Merdeka 1",
      "address": "Jl. Sekolah No. 1, Bandung"
    },
    "delivery_date": "2026-10-03",
    "status": "pending",
    "total_portions": 240,
    "distributed_portions": 235,
    "distribution_pct": 97,
    "allergy_students": {
      "total": 33,
      "breakdown": [
        { "allergy_type": "Telur", "count": 15 },
        { "allergy_type": "Kacang", "count": 12 },
        { "allergy_type": "Susu", "count": 6 }
      ]
    },
    "menu_items": [
      {
        "id": 21,
        "delivery_report_id": 9,
        "menu_name": "Nasi Putih",
        "category": "Karbohidrat",
        "portions": 240,
        "is_allergy_substitute": false,
        "allergy_note": ""
      }
    ],
    "photos": {
      "photo1": "uploads/delivery/delivery_12_1760000000_kirim-1.jpg",
      "photo2": ""
    }
  }
}
```

| Field | Catatan |
|---|---|
| `distribution_pct` | Integer hasil `(distributed × 100) / total` — pembagian bulat, bukan desimal |
| `photos.photo1/2` | **Path di server**, bukan URL publik. String kosong kalau tidak ada foto |

### Kemungkinan Error

| Code | Message |
|---|---|
| 404 | `Laporan tidak ditemukan.` — ID salah **atau** laporan milik SPPG lain |

---

# BAGIAN F — Masalah Makanan

## [17] Lapor Masalah Makanan (Internal SPPG)

Dipakai SPPG untuk melaporkan masalah yang ditemukan sendiri, sebelum/saat distribusi.

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/sppg/food-problems` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **form-data** |

### Header
```
Authorization: Bearer {{token}}
```

### Body (form-data)

> ⚠️ **Nama field di endpoint ini pakai camelCase Bahasa Indonesia**, beda dari endpoint lain. Perhatikan betul.

| Key | Type | Wajib | Aturan | Contoh |
|---|---|---|---|---|
| `namaSekolah` | Text | ✅ | Nama sekolah (free text, bukan ID) | `SDN Merdeka 1` |
| `namaMakanan` | Text | ✅ | Nama makanan bermasalah | `Ayam Goreng` |
| `jenisMasalah` | Text | ✅ | Free text, maks 80 karakter | `Makanan basi` |
| `porsiTerdampak` | Text | ✅ | Angka > 0 | `25` |
| `deskripsiMasalah` | Text | ✅ | Penjelasan detail | `Ditemukan bau tidak sedap pada 25 porsi saat pengemasan.` |
| `prioritas` | Text | ❌ | **Hanya** `Rendah`, `Sedang`, `Tinggi`. Default `Sedang` | `Tinggi` |
| `fotoMasalah` | File | ❌ | Foto bukti | `basi.jpg` |

### Response `201 Created`
```json
{
  "status": "success",
  "message": "Laporan masalah makanan berhasil dikirim.",
  "data": {
    "id": 3,
    "sppg_id": 12,
    "school_name": "SDN Merdeka 1",
    "food_name": "Ayam Goreng",
    "problem_type": "Makanan basi",
    "affected_portions": 25,
    "priority": "Tinggi",
    "description": "Ditemukan bau tidak sedap pada 25 porsi saat pengemasan.",
    "photo_path": "uploads/sppg-problems/sppg_problem_12_1760000000_basi.jpg",
    "status": "baru",
    "created_at": "2026-10-03T11:00:00Z",
    "updated_at": "2026-10-03T11:00:00Z"
  }
}
```

> Response pakai **snake_case** (`school_name`), sedangkan request pakai **camelCase** (`namaSekolah`). Tidak simetris — jangan bingung.

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Lengkapi semua data wajib (nama sekolah, nama makanan, jenis masalah, porsi, deskripsi).` |
| 400 | `porsiTerdampak harus berupa angka > 0.` |
| 400 | `prioritas tidak valid. Pilihan: Rendah, Sedang, Tinggi.` |
| 500 | `Gagal menyimpan laporan.` |

---

## [18] Lihat Komplain dari Sekolah

Laporan makanan bermasalah yang **dikirim sekolah** ke SPPG ini. Sumber data beda dari step [17] — ini tabel `food_reports`, bukan `sppg_food_problems`.

| | |
|---|---|
| Method | `GET` |
| URL | `{{base_url}}/sppg/food-reports` |
| Query opsional | `?month=10&year=2026` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "data": [
    {
      "id": 5,
      "report_date": "2026-10-02",
      "reporter": "SDN Merdeka 1",
      "issue_type": "spoiled_food",
      "description": "Nasi terasa asam saat dibagikan ke siswa.",
      "photo_url": "uploads/reports/report_7_1759999999_nasi.jpg",
      "status": "belum_ditinjau"
    }
  ]
}
```

| Field | Nilai yang mungkin |
|---|---|
| `issue_type` | `food_shortage`, `spoiled_food`, `allergen`, `late_delivery` |
| `status` | `belum_ditinjau`, `sedang_ditinjau`, `selesai` |
| `reporter` | Nama sekolah pengirim laporan |

Diurutkan `report_date` **descending**.

---

# BAGIAN G — Akun & Keluar

## [19] Ganti Password

Endpoint umum, bisa dipakai semua role termasuk SPPG.

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/change-password` |
| Auth | `Authorization: Bearer {{token}}` |
| Body type | **raw → JSON** |

### Header
```
Content-Type: application/json
Authorization: Bearer {{token}}
```

### Body
```json
{
  "old_password": "SPPGDapurSehatBandung4821",
  "new_password": "PasswordBaruAman123"
}
```

| Field | Wajib | Aturan |
|---|---|---|
| `old_password` | ✅ | Harus cocok dengan password saat ini |
| `new_password` | ✅ | Minimal **8 karakter** |

### Response `200 OK`
```json
{
  "status": "success",
  "message": "Password berhasil diubah."
}
```

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Password lama dan password baru wajib diisi.` |
| 400 | `Password baru minimal 8 karakter.` |
| 401 | `Password lama tidak sesuai.` |
| 404 | `User tidak ditemukan.` |

> Token lama **tetap valid** setelah ganti password (tidak ada revocation). Berlaku sampai `exp`.

---

## [20] Logout

| | |
|---|---|
| Method | `POST` |
| URL | `{{base_url}}/logout` |
| Auth | `Authorization: Bearer {{token}}` |
| Body | — |

### Response `200 OK`
```json
{
  "status": "success",
  "message": "Logout berhasil."
}
```

### Kemungkinan Error

| Code | Message |
|---|---|
| 400 | `Gagal logout: Tidak ada token yang dikirim.` |
| 401 | `Format token tidak valid. Gunakan format 'Bearer <token>'.` |
| 401 | `Token tidak valid atau sudah kadaluarsa.` |

> **Logout hanya validasi token, tidak ada blacklist server-side.** Token yang sama masih bisa dipakai sampai expired. Hapus token di sisi client.

---

# Lampiran

## A. Ringkasan Semua Endpoint SPPG

| # | Method | Endpoint | Body | Auth |
|---|---|---|---|---|
| 1 | POST | `/register/sppg` | form-data | — |
| 2 | GET | `/admin/registrations` | — | admin |
| 3 | PATCH | `/admin/registrations/:id/status` | JSON | admin |
| 4 | POST | `/login` | JSON | — |
| 5 | GET | `/sppg/profile` | — | sppg |
| 6 | PUT | `/sppg/profile` | JSON | sppg |
| 7 | GET | `/sppg/status` | — | sppg |
| 8 | GET | `/sppg/available-schools` | — | sppg |
| 9 | POST | `/sppg/assign-schools` | JSON | sppg |
| 10 | GET | `/sppg/schools` | — | sppg |
| 11 | GET | `/sppg/allergy-dashboard` | — | sppg |
| 12 | POST | `/sppg/scan` | form-data | sppg |
| 13 | POST | `/sppg/menu-reports` | JSON | sppg |
| 14 | POST | `/sppg/delivery-reports` | form-data | sppg |
| 15 | GET | `/sppg/delivery-reports` | — | sppg |
| 16 | GET | `/sppg/delivery-reports/:id` | — | sppg |
| 17 | POST | `/sppg/food-problems` | form-data | sppg |
| 18 | GET | `/sppg/food-reports` | — | sppg |
| 19 | POST | `/change-password` | JSON | any |
| 20 | POST | `/logout` | — | any |

Endpoint penunjang: `GET /ping` (health check), `GET /dashboard/stats` (statistik publik), `POST /forgot-password/request`, `POST /forgot-password/verify`.

## B. Error Auth Global

Semua endpoint `/sppg/*` lewat middleware `SppgOnly`:

| Code | Message | Penyebab |
|---|---|---|
| 401 | `token tidak ditemukan` | Header `Authorization` tidak ada |
| 401 | `format token tidak valid` | Bukan format `Bearer <token>` |
| 401 | `Token tidak valid atau sudah kadaluarsa` | Signature salah atau sudah > 24 jam |
| 401 | `Gagal membaca token` | Claims tidak bisa di-parse |
| 403 | `Hanya mitra sppg yang dapat mengakses fitur ini` | Role di token bukan `sppg` |
| 500 | `konfigurasi keamanan server tidak lengkap` | Env `JWT_SECRET` kosong |

## C. Konvensi Penamaan Field (Tidak Konsisten — Perhatikan)

| Endpoint | Konvensi request | Catatan |
|---|---|---|
| `/sppg/food-problems` | camelCase ID (`namaSekolah`) | Response snake_case |
| `/sppg/profile` | snake_case, tapi alamat = `lokasi` | DB kolomnya `address` |
| `/sppg/schools` | — | Response pakai `Grade` (kapital G) |
| `/sppg/scan` | snake_case | Response campur ID (`kalori`, `bahan_terdeteksi`) |
| Endpoint lain | snake_case | — |

## D. Urutan Testing Minimal di Postman

Untuk tes alur SPPG end-to-end tanpa data awal:

1. `POST /register/sppg` → catat `data.id`
2. Login admin → simpan `{{admin_token}}`
3. `PATCH /admin/registrations/{id}/status` body `{"status":"disetujui"}` → catat `data.email`
4. Ambil password dari email PIC, **atau** reset via `/forgot-password/request` + `/forgot-password/verify`
5. `POST /login` → token auto-tersimpan ke `{{token}}`
6. `GET /sppg/profile` → pastikan `200`
7. `GET /sppg/available-schools` → catat ID sekolah
8. `POST /sppg/assign-schools` → `{"school_ids":[7]}`
9. `GET /sppg/schools` → verifikasi sekolah masuk
10. `POST /sppg/scan` (ML service harus jalan, atau cek `gemini_failed: true`)
11. `POST /sppg/menu-reports`
12. `POST /sppg/delivery-reports` → catat `data.id`
13. `GET /sppg/delivery-reports/{id}`
14. `POST /sppg/food-problems`
15. `GET /sppg/food-reports`
