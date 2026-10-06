# SIAKAD Mini - RESTful API Back End

RESTful API back end untuk **SIAKAD Mini**, layanan akademik sederhana yang mengelola data mahasiswa, mata kuliah, dan Kartu Rencana Studi (KRS) menggunakan **Go (Golang)**, **Fiber v2**, dan **PostgreSQL**.

---

## 📌 Ketentuan Teknis & Arsitektur

- **Bahasa**: Go 1.26.5
- **Framework Web**: Fiber v2 (`github.com/gofiber/fiber/v2`)
- **Database Driver**: pgx v5 (`github.com/jackc/pgx/v5`) dengan connection pooling
- **Database**: PostgreSQL
- **Autentikasi**: JWT (JSON Web Token) dengan Bearer Token
- **Keamanan Password**: Hash bcrypt
- **Rate Limiting**: Maksimal 5x gagal login per menit (HTTP 429)
- **Concurrency Control**: Row Locking (`SELECT ... FOR UPDATE`) pada pengambilan KRS (Enrollments)
- **Data Seeder**: 1 Admin, 20 Mahasiswa, 10 Mata Kuliah

---

## 🗄️ Model Data (Database Schema)

1. **`users`**:
   - `id` (PK, Serial)
   - `email` (Unique, Varchar)
   - `password` (Hash Bcrypt)
   - `role` (`admin` / `mahasiswa`)
2. **`students`**:
   - `id` (PK, Serial)
   - `user_id` (FK `users.id`, 1-to-1)
   - `nim` (Unique, 12 digit)
   - `nama` (Varchar)
   - `prodi` (Varchar)
   - `angkatan` (Integer, 4 digit)
   - `ipk_terakhir` (Numeric 3.2, rentang 0.00–4.00)
   - `deleted_at` (Timestamptz, soft delete)
3. **`courses`**:
   - `id` (PK, Serial)
   - `kode_mk` (Unique, Varchar)
   - `nama_mk` (Varchar)
   - `sks` (Integer)
   - `semester` (Integer)
   - `kuota` (Integer)
4. **`enrollments`**:
   - `id` (PK, Serial)
   - `student_id` (FK `students.id`)
   - `course_id` (FK `courses.id`)
   - `tahun_akademik` (Varchar, e.g., `2026/2027-Ganjil`)
   - `created_at` (Timestamptz)
   - *Constraint Unik*: `(student_id, course_id, tahun_akademik)`

---

## ⚖️ Business Rules

1. **Batas SKS per Semester berdasarkan IPK Terakhir**:
   - $\text{IPK} \ge 3.00 \rightarrow$ Maksimal **24 SKS**
   - $2.50 \le \text{IPK} < 3.00 \rightarrow$ Maksimal **21 SKS**
   - $\text{IPK} < 2.50 \rightarrow$ Maksimal **18 SKS**
2. Mahasiswa tidak dapat mengambil mata kuliah yang sama dua kali pada tahun akademik yang sama (HTTP `409 Conflict`).
3. Mata kuliah yang kuotanya penuh tidak dapat diambil (HTTP `422 Unprocessable Entity`).
4. Mahasiswa hanya dapat mengakses data dan membatalkan KRS miliknya sendiri (HTTP `403 Forbidden` jika mengakses milik orang lain).
5. Mahasiswa yang telah di-*soft delete* tidak dapat login (HTTP `401 Unauthorized`) dan tidak muncul di daftar mahasiswa.

---

## 🚀 Menjalankan Server

1. **Konfigurasi Environment (`.env`)**:
   ```env
   APP_NAME=SIAKAD Mini
   APP_PORT=3000

   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=praktikum_backend
   DB_SSLMODE=disable

   JWT_SECRET=65b458dc012be6c16665db71ee6ec7f53dc9726237bbbd4e7288bacd789955cf
   JWT_ISSUER=siakad-mini
   JWT_ACCESS_TTL_MINUTES=15
   JWT_REFRESH_TTL_DAYS=7
   ```

2. **Jalankan Migrasi & Seeding Data**:
   ```bash
   go run cmd/gen_seed/main.go
   go run cmd/seeder/main.go
   ```

3. **Jalankan Aplikasi**:
   ```bash
   go run main.go
   ```

4. **Jalankan Unit & Integration Test**:
   ```bash
   go test -v ./...
   ```

---

## 📋 Daftar 10 Endpoint RESTful API

| No | Method | Endpoint | Akses | Deskripsi & Status Sukses |
|---|---|---|---|---|
| 1 | `POST` | `/api/v1/auth/login` | Publik | Login & mendapatkan JWT Access Token. (`200 OK`) |
| 2 | `GET` | `/api/v1/auth/me` | Semua Role | Profil user login (+ data student jika mahasiswa). (`200 OK`) |
| 3 | `GET` | `/api/v1/students` | Admin | Daftar mahasiswa + pagination, filter, search, sort. (`200 OK`) |
| 4 | `POST` | `/api/v1/students` | Admin | Menambah mahasiswa + akun user (password = hash NIM). (`201 Created`) |
| 5 | `GET` | `/api/v1/students/:id` | Admin / Mahasiswa (sendiri) | Detail mahasiswa + daftar MK + `total_sks` + `batas_sks`. (`200 OK`) |
| 6 | `PUT` | `/api/v1/students/:id` | Admin | Update mahasiswa (nama, prodi, angkatan, ipk; NIM tetap). (`200 OK`) |
| 7 | `DELETE` | `/api/v1/students/:id` | Admin | Soft delete mahasiswa. (`204 No Content`) |
| 8 | `GET` | `/api/v1/courses` | Semua Role | Daftar mata kuliah + hitung `terisi` & `sisa_kuota`. (`200 OK`) |
| 9 | `POST` | `/api/v1/enrollments` | Mahasiswa | Ambil MK (KRS) dengan row locking & validasi SKS. (`201 Created`) |
| 10 | `DELETE` | `/api/v1/enrollments/:id` | Mahasiswa (sendiri) | Batalkan MK dari KRS (kuota kembali bertambah). (`204 No Content`) |

---

## 🧑‍💻 Kredensial Akun Default (Seeder)

- **Admin**:
  - Email: `admin@siakad.ac.id`
  - Password: `Admin123!`
- **Mahasiswa (20 akun)**:
  - Email: `rina.putri@siakad.ac.id` (NIM: `187221000001`, Password: `187221000001`)
  - Email: `budi.santoso@siakad.ac.id` (NIM: `187221000002`, Password: `187221000002`)
  - ... hingga `taufik.hidayat@siakad.ac.id` (NIM: `187221000020`, Password: `187221000020`)
