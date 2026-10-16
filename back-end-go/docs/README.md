# NutriSafe API — Dokumentasi

Dokumentasi REST API backend Go (Fiber v2). Semua contoh siap di-copy ke Postman.

## Daftar Dokumen

| Dokumen | Isi |
|---|---|
| [API_AUTH.md](API_AUTH.md) | Login, register (semua role), lupa password, profil generik, ganti password, logout |
| [API_SPPG.md](API_SPPG.md) | Alur mitra SPPG: pendaftaran → kelola sekolah → scan menu → laporan pengiriman |
| [API_ADMIN.md](API_ADMIN.md) | Dashboard, verifikasi pendaftaran SPPG, kelola sekolah, laporan, distribusi |
| [API_SCHOOL.md](API_SCHOOL.md) | Profil sekolah, siswa, kelas, guru, data alergi, terima kiriman, lapor masalah |
| [API_UMUM.md](API_UMUM.md) | Akun umum & siswa: profil, daftar sekolah, lapor masalah publik |

---

## 1. Informasi Umum

| Item | Nilai |
|---|---|
| Base URL (local) | `http://localhost:3000` |
| Base URL (produksi) | `https://api.nutrisafe.id` |
| Framework | Go + Fiber v2 |
| Database | PostgreSQL (Supabase) |
| Format response | JSON |
| Auth | JWT Bearer token, berlaku **24 jam** |
| Body limit | 10 MB per request (nginx produksi: 25 MB) |

### Akses Produksi

Trafik publik masuk lewat **nginx host** (`deploy/nginx/nutrisafe.conf`), bukan langsung ke container:

| URL | Hasil |
|---|---|
| `https://api.nutrisafe.id` | ✅ Backend Go — **ini yang dipakai** |
| `http://api.nutrisafe.id` | `301` redirect ke HTTPS |
| `http://<ip-vps>:3000` | ❌ Unreachable — port di-bind ke `127.0.0.1` saja |
| `http://<ip-vps>` | ❌ `444` — nginx `default_server` tolak akses via IP mentah |

Domain lain di VPS yang sama: `nutrisafe.id` (FE user, port 5173), `admin.nutrisafe.id` (FE admin, port 5174).

Timeout proxy nginx = **300s**, jadi request lama seperti `POST /sppg/scan` tidak terputus.

---

## 2. Bentuk Response

Sukses:
```json
{
  "status": "success",
  "message": "...",
  "data": { }
}
```

Error — **selalu** bentuk yang sama, apapun kodenya:
```json
{
  "status": "error",
  "message": "Pesan error dalam Bahasa Indonesia"
}
```

> Tidak ada kode error mesin (`error_code`). Client harus membedakan error lewat **HTTP status** + string `message`. String ini bisa berubah, jadi jangan dipakai sebagai logika percabangan.

Pengecualian: `GET /dashboard/stats` menambahkan `error_detail` berisi pesan error database mentah saat gagal.

---

## 3. Autentikasi

### Dapatkan Token

`POST /login` → response memuat `token`. Lihat [API_AUTH.md](API_AUTH.md).

### Pakai Token

```
Authorization: Bearer <token>
```

Isi JWT: `user_id`, `email`, `role`, `exp` (+24 jam). Signed HS256 dengan env `JWT_SECRET`.

### Peta Role → Middleware

| Prefix | Middleware | Role yang diterima |
|---|---|---|
| `/admin/*` | `AdminOnly` | `admin` |
| `/sppg/*` | `SppgOnly` | `sppg` |
| `/school/*` | `SchoolOnly` | `school` |
| `/umum/*` | `UmumOnly` | `umum` **atau** `siswa` |
| sisanya | — | publik / cek token sendiri di handler |

### Error Auth Global

Sama di semua middleware (beda kapitalisasi pesan saja):

| Code | Message | Penyebab |
|---|---|---|
| 401 | `Token tidak ditemukan` | Header `Authorization` tidak ada |
| 401 | `Format token tidak valid` | Bukan `Bearer <token>` |
| 401 | `Token tidak valid atau sudah kadaluarsa` | Signature salah atau > 24 jam |
| 401 | `Gagal membaca token` | Claims tidak bisa di-parse |
| 403 | `Hanya <role> yang dapat mengakses fitur ini` | Role di token tidak cocok |
| 500 | `Konfigurasi keamanan server tidak lengkap` | Env `JWT_SECRET` kosong |

> `/sppg/*` dan `/school/*` memakai pesan huruf kecil (`token tidak ditemukan`), `/admin/*` dan `/umum/*` huruf besar (`Token tidak ditemukan`). Jangan cocokkan string persis.

### Logout Tidak Mencabut Token

`POST /logout` hanya **memvalidasi** token, tidak ada blacklist server-side. Token yang sama tetap bisa dipakai sampai `exp`. Hapus token di sisi client. Hal yang sama berlaku setelah ganti password.

---

## 4. Enam Role

| Role | Cara akun dibuat | Format login email |
|---|---|---|
| `admin` | `POST /register` dengan `role_name: "admin"` ⚠️ (lihat §7) | email asli |
| `umum` | `POST /register` dengan `role_name: "umum"` | email asli |
| `school` | `POST /register/school` — akun dibuat langsung | `<slug-nama-sekolah>@sch.id` |
| `sppg` | `POST /register/sppg` → **harus disetujui admin** | `<slug-nama-sppg>@sppg.id` |
| `siswa` | Otomatis oleh cron, dari data siswa yang diinput sekolah | `<slug-nama>.<4-digit-nisn>@siswa.sch.id` |
| — | `role_name` lain ditolak | — |

Login email untuk `school`, `sppg`, dan `siswa` **di-generate sistem**, bukan diisi user. Password temporary = gabungan kata nama (kapital di awal tiap kata, non-alfanumerik dibuang) + 4 digit acak — contoh `SPPGDapurSehat4821` (`utils/helpers.go:43`).

Kalau slug tabrakan, sistem menambah counter: `nama@sppg.id` → `nama.2@sppg.id` → `nama.3@sppg.id`.

---

## 5. Konvensi Postman

### Environment

| Variable | Initial value |
|---|---|
| `base_url` | `http://localhost:3000` |
| `token` | *(kosong — diisi otomatis)* |
| `admin_token` | *(kosong)* |

Script di tab **Tests** pada request Login:
```javascript
const res = pm.response.json();
if (res.token) pm.environment.set("token", res.token);
```

### Header per Jenis Body

| Jenis | Header |
|---|---|
| JSON | `Content-Type: application/json` |
| File upload | *(kosongkan — Postman isi otomatis)* |
| Protected | tambah `Authorization: Bearer {{token}}` |

> **Saat pakai `form-data`, hapus header `Content-Type` manual.** Kalau diisi manual, boundary multipart hilang dan server balas `Format request tidak valid`.

### Endpoint yang Pakai form-data (bukan JSON)

Semua endpoint yang menerima file. Perhatikan — sisanya JSON:

| Endpoint | Field file |
|---|---|
| `POST /register/sppg` | `proposal` (.pdf), `kitchen_photo` (≤5 MB) |
| `POST /sppg/scan` | `image` |
| `POST /sppg/delivery-reports` | `photo1`, `photo2` |
| `POST /sppg/food-problems` | `fotoMasalah` |
| `POST /school/reports` | `photo` |
| `POST /school/delivery-receipt` | `photo` |
| `POST /umum/reports` | `photo` |

---

## 6. Konvensi Penamaan — Tidak Konsisten

Daftar jebakan nyata di kode. Tidak diseragamkan supaya front-end yang sudah jalan tidak pecah.

| Endpoint | Jebakan |
|---|---|
| `POST /sppg/food-problems` | Request pakai **camelCase Indonesia** (`namaSekolah`, `porsiTerdampak`, `fotoMasalah`), tapi response snake_case (`school_name`) |
| `GET /sppg/schools` | Response pakai key **`Grade`** (kapital G) |
| `GET /umum/schools` | Response pakai key **`Grade`** (kapital G) |
| `GET /sppg/available-schools` | Response pakai `grade` (huruf kecil) — beda dari dua di atas |
| `POST /sppg/scan` | `kalori_per_porsi` keluar sebagai **string** `"500.0"`, tapi `POST /sppg/menu-reports` minta **number** |
| Profil semua role | Field alamat bernama `lokasi` di API, kolomnya `address` di DB |
| `POST /school/students` | Gender valid: `Laki-laki` / `Perempuan` (huruf **l** kecil) |
| `POST /school/teachers` | Gender valid: `Laki-Laki` / `Perempuan` (huruf **L** besar) — beda dari siswa |
| `GET /admin/distributions` | Status diterjemahkan: DB `pending` ↔ API `dalam_pengiriman` |
| `POST /school/allergy-data` | Response bungkus sekolah dalam objek dengan key `name_school` (bukan `school_name`) |

### Kosakata Status

| Konteks | Nilai |
|---|---|
| Pendaftaran SPPG | `pending`, `belum_ditinjau`, `verifikasi`, `disetujui` |
| Laporan makanan (`FoodReport`) | `belum_ditinjau`, `sedang_ditinjau`, `selesai` |
| Pengiriman (`DeliveryReport`) | `pending`, `diterima` — tapi API admin pakai `dalam_pengiriman`, `diterima` |
| Masalah makanan SPPG | `baru` (satu-satunya; tidak ada endpoint yang mengubahnya) |
| Prioritas masalah | `Rendah`, `Sedang`, `Tinggi` (kapital di awal) |
| Profil school/sppg | `active`, `inactive` |

### Jenis Masalah (`issue_type`)

Dipakai `POST /school/reports` dan `POST /umum/reports`:
`kekurangan_makanan`, `makanan_basi`, `masalah_alergen`, `terlambat`, `lainnya`

> Komentar di `models/report.go:9` menyebut nilai lain (`food_shortage`, `spoiled_food`, `allergen`, `late_delivery`) — itu **stale**, bukan yang divalidasi kode.

### Tingkat Sekolah (`grade`)

`SD/MI`, `SMP/MTs/MTsN`, `SMA/SMK/MA/MAN` — validasi persis, termasuk garis miring.

---

## 7. ⚠️ Peringatan Keamanan

Temuan dari pembacaan kode. **Belum diperbaiki** — didokumentasikan supaya tidak tertelan.

### 7.1 Siapa pun bisa mendaftar sebagai admin

`POST /register` adalah endpoint **publik tanpa auth** (`routes/auth.go:16`). Validasi rolenya di `controllers/auth/register_controller.go:57`:

```go
if input.RoleName != "school" && input.RoleName != "umum" && input.RoleName != "admin" {
    // ditolak
}
```

`"admin"` lolos validasi. Artinya request ini memberi akses penuh ke seluruh `/admin/*`:

```json
POST /register
{ "name": "x", "email": "x@x.com", "phone": "0811", "password": "x", "role_name": "admin" }
```

Pesan error di endpoint yang sama berbunyi `"Pilihan yang tersedia: school, umum."` — jadi `admin` tampaknya **tidak disengaja** ikut di daftar.

Dampak: pembuatan/persetujuan mitra SPPG, lihat semua data sekolah & siswa, ubah status semua laporan.

Perbaikan minimal: buang `&& input.RoleName != "admin"` dari kondisi itu, lalu buat akun admin lewat seeder atau langsung di database.

### 7.2 Token tidak bisa dicabut

Tidak ada blacklist/revocation. Token tetap sah 24 jam penuh meski user logout atau ganti password. Pencurian token = akses 24 jam.

### 7.3 Enumerasi email tertutup, tapi tidak di semua jalur

`POST /forgot-password/request` sengaja selalu balas `200` supaya email tidak bisa dienumerasi. Tapi `POST /register` membedakan `409 Email ini sudah terdaftar` vs sukses — jadi enumerasi tetap mungkin lewat sana.

### 7.4 Path file diekspos apa adanya

Response memuat path filesystem server (`uploads/delivery/delivery_12_176...jpg`), bukan URL publik. Tidak ada endpoint yang menyajikan file tersebut — belum ada static file serving di `main.go`. Jadi foto yang diupload **tidak bisa diakses lewat API**.

---

## 8. Catatan Perilaku

### Cron Berjalan Otomatis

`cron.StartStudentAccountCron()` (`main.go:19`) jalan tiap **5 menit**, melakukan dua hal:

1. **Sinkronisasi `student_count`** di semua `school_profiles` dari `COUNT` aktual tabel `students`.
2. **Membuat akun siswa** untuk siswa yang `account_generated = false` dan `updated_at` sudah lebih dari **15 menit** lalu. Kredensial dikirim ke `contact_email` sekolah sebagai satu email berisi daftar.

Efek yang perlu disadari saat testing: siswa yang baru diinput **berubah jadi tidak bisa diedit/dihapus** setelah ±15–20 menit, karena `account_generated` jadi `true` dan `PUT`/`DELETE /school/students/:id` menolak dengan `403`.

### ML Service

`POST /sppg/scan` memanggil `ML_SERVICE_URL` (default `http://localhost:8000`, di Docker `http://ml-service:8000`). Timeout 90 detik. Kalau gagal, **tetap balas `200`** dengan nilai gizi dummy dan `gemini_failed: true`.

### Email Asinkron

Semua pengiriman email jalan di goroutine terpisah (`go func()`). Request tidak menunggu, dan **kegagalan kirim tidak muncul di response** — hanya di log server (`[mailer] gagal kirim ...`).

Produksi: `SMTP_HOST=smtp.gmail.com`, `SMTP_FROM=admin@nutrisafe.com`. Perhatikan domainnya `.com`, bukan `.id` seperti domain situs. Lihat bagian Troubleshooting di [API_SPPG.md](API_SPPG.md#troubleshooting) kalau email kredensial tidak sampai.

---

## 9. Ringkasan Semua Endpoint

### Publik

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/ping` | Health check |
| GET | `/dashboard/stats` | Statistik publik (jumlah sekolah/SPPG/siswa) |
| POST | `/register` | Daftar umum/admin ⚠️ |
| POST | `/register/school` | Daftar sekolah |
| POST | `/register/sppg` | Daftar mitra SPPG (form-data) |
| POST | `/login` | Login semua role |
| POST | `/forgot-password/request` | Minta kode reset |
| POST | `/forgot-password/verify` | Verifikasi kode + set password baru |

### Butuh Token (role apa saja)

| Method | Endpoint |
|---|---|
| GET | `/profile/detail` |
| PUT | `/profile/detail` |
| POST | `/change-password` |
| POST | `/logout` |

### Per Role

| Role | Jumlah endpoint | Dokumen |
|---|---|---|
| `admin` | 11 | [API_ADMIN.md](API_ADMIN.md) |
| `sppg` | 14 | [API_SPPG.md](API_SPPG.md) |
| `school` | 22 | [API_SCHOOL.md](API_SCHOOL.md) |
| `umum` / `siswa` | 5 | [API_UMUM.md](API_UMUM.md) |

Total 52 endpoint role-scoped + 12 publik/generik.

---

## 10. Health Check

```
GET {{base_url}}/ping
```

Tanpa auth, tanpa body.

```json
{
  "status": "success",
  "message": "API Nutrisafe successfully running!"
}
```

Pakai ini dulu untuk memastikan base URL benar sebelum debug hal lain.
