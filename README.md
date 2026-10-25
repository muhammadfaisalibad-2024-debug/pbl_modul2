# SIAKAD Mini - Polyglot Persistence RESTful API

RESTful API back end untuk **SIAKAD Mini**, sistem informasi akademik modern yang menerapkan arsitektur **Polyglot Persistence** (relasional PostgreSQL + dokumen NoSQL MongoDB) menggunakan **Go (Golang)** dan framework **Fiber v2**.

---

## 📌 Ketentuan Teknis & Arsitektur

- **Bahasa**: Go 1.26.5
- **Framework Web**: Fiber v2 (`github.com/gofiber/fiber/v2`)
- **Database Relasional**: PostgreSQL dengan driver `pgx/v5` (`github.com/jackc/pgx/v5`) & connection pooling
- **Database NoSQL**: MongoDB dengan official driver v2 (`go.mongodb.org/mongo-driver/v2`)
- **Arsitektur Persistence**:
  - **PostgreSQL**: Data relasional inti (Users, Roles, Permissions, Students, Courses, Enrollments, Supervisions).
  - **MongoDB**: Data fleksibel berbasis dokumen (Portofolio Mahasiswa, Proyek, Sertifikat, dan Ulasan Kelas / Reviews).
- **Clean Architecture**: Package `app/model` murni domain logic tanpa dependensi BSON/driver MongoDB.
- **Autentikasi & Otorisasi**: JWT (Bearer Token) + Role-Based Access Control (RBAC) dinamis dari tabel `role_permissions`.
- **Keamanan**: Hash bcrypt, Rate Limiting login (HTTP 429), penolakan NoSQL injection pada body parser (HTTP 400).
- **Konkurensi & Integritas**:
  - PostgreSQL: Row-level locking (`SELECT ... FOR UPDATE`) pada transaksi KRS atomik dan alokasi pembimbing tugas akhir.
  - MongoDB: Atomik conditional update (`projects.19: {$exists: false}`) untuk pencegahan race condition batas 20 proyek.

---

## 🗄️ Model Data & Penyimpanan

### 1. PostgreSQL (Relasional)
1. **`users`**: Akun pengguna (`id`, `username`, `email`, `password`, `role`, `is_active`, `created_at`).
2. **`roles` & `permissions`**: Pengaturan RBAC (`role_permissions`).
3. **`students`**: Biodata mahasiswa (`id`, `user_id`, `nim`, `nama`, `prodi`, `angkatan`, `ipk_terakhir`, `deleted_at`).
4. **`courses`**: Mata kuliah (`id`, `kode_mk`, `nama_mk`, `sks`, `semester`, `kuota`).
5. **`enrollments`**: Transaksi KRS mahasiswa (`id`, `student_id`, `course_id`, `tahun_akademik`, `created_at`).
6. **`supervisions`**: Pembimbing Tugas Akhir (`id`, `student_id`, `supervisor_id`, `status`, `assigned_at`).

### 2. MongoDB (NoSQL Dokumen)
1. **`student_portfolios`**:
   - `student_id` (Integer, Unique Index: `student_unique`)
   - `headline`, `bio`, `linkedin_url`, `github_url`, `skills`
   - `projects` (Embedded Array): `_id` (ObjectID), `title`, `technologies` (Multikey Index: `project_tech`), `url`, `year`
   - `certificates` (Embedded Array): `_id` (ObjectID), `name`, `issuer`, `date`
2. **`course_reviews`**:
   - `course_id` (Integer), `student_id` (Integer), `rating` (1-5), `comment`, `difficulty`, `helpful_count`
   - `replies` (Embedded Array): tanggapan dari dosen / admin.

---

## ⚖️ Aturan Bisnis (Business Rules)

1. **Batas SKS berdasarkan IPK Terakhir**:
   - $\text{IPK} \ge 3.00 \rightarrow$ Maksimal **24 SKS**
   - $2.50 \le \text{IPK} < 3.00 \rightarrow$ Maksimal **21 SKS**
   - $\text{IPK} < 2.50 \rightarrow$ Maksimal **18 SKS**
2. **KRS Atomik & Anti-Overbooking**:
   - Pengecekan kuota dan batas SKS dieksekusi dalam satu transaksi atomik berisolasi tinggi.
   - Pembatalan KRS mengembalikan kuota secara otomatis.
3. **Batas Proyek Portofolio (Modul 10)**:
   - Satu mahasiswa maksimal memiliki **20 proyek**. Penambahan proyek ke-21 ditolak atomik (`409 Conflict`).
4. **Independensi Layanan (Fault Isolation)**:
   - Jika MongoDB mati, endpoint akademik inti (Login, Siswa, Kelas, KRS) yang bergantung pada PostgreSQL tetap berjalan normal (**200 OK**). Endpoint `/health` secara transparan melaporkan status **503 Service Unavailable**.

---

## 🚀 Panduan Menjalankan Aplikasi

### 1. Konfigurasi Lingkungan (`.env`)
Salin file `.env.example` menjadi `.env` lalu sesuaikan konfigurasinya:
```env
APP_NAME=SIAKAD Mini
APP_PORT=3000

# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_db_password
DB_NAME=praktikum_backend
DB_SSLMODE=disable
DB_MAX_CONNS=10

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=praktikum_backend

# JWT (Gunakan string acak aman minimal 32 karakter)
JWT_SECRET=your_super_secret_jwt_key_change_in_production
JWT_ISSUER=siakad-mini
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=7
```

### 2. Migrasi Database PostgreSQL & Seeding
```powershell
# Jalankan migrasi dan seeder awal
go run cmd/seeder/main.go
```

### 3. Menjalankan Server
```powershell
go run main.go
```
*Aplikasi akan berjalan pada `http://localhost:3000` dan otomatis menginisialisasi indeks unik MongoDB.*

---

## 📋 Daftar Endpoint RESTful API

### 1. Autentikasi & Pemantauan Sistem
| Method | Endpoint | Akses | Keterangan |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` / `/api/v1/health` | Publik | Cek kesehatan koneksi PostgreSQL dan MongoDB |
| `POST` | `/api/v1/auth/login` | Publik | Login & mendapatkan JWT Access Token |
| `GET` | `/api/v1/auth/me` | Terautentikasi | Profil pengguna login (+ biodata mahasiswa) |
| `POST` | `/api/v1/auth/refresh` | Publik | Refresh access token menggunakan refresh token |
| `POST` | `/api/v1/auth/logout` | Terautentikasi | Logout & pencabutan sesi |

### 2. Mahasiswa & Portofolio (NoSQL MongoDB)
| Method | Endpoint | Akses | Keterangan |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/students` | Admin | Daftar mahasiswa + pagination, filter prodi/angkatan, search & sort |
| `POST` | `/api/v1/students` | Admin | Tambah mahasiswa baru + akun user otomatis |
| `GET` | `/api/v1/students/:id` | Admin / Mahasiswa Sendiri | Detail mahasiswa + status SKS |
| `PUT` | `/api/v1/students/:id` | Admin | Perbarui biodata mahasiswa |
| `DELETE` | `/api/v1/students/:id` | Admin | Soft delete mahasiswa (`deleted_at`) |
| `GET` | `/api/v1/students/:id/portfolio` | Terautentikasi | Mengambil portofolio lengkap dari MongoDB |
| `PUT` | `/api/v1/students/:id/portfolio` | Pemilik / Admin | Create / Upsert / Replace data portofolio |
| `POST` | `/api/v1/students/:id/portfolio/projects` | Pemilik / Admin | Tambah proyek baru (atomik, maks 20 proyek) |
| `DELETE` | `/api/v1/students/:id/portfolio/projects/:projectId` | Pemilik / Admin | Hapus proyek berdasarkan ObjectID |
| `GET` | `/api/v1/portfolios/stats/technologies` | Terautentikasi | Agregasi agregat Top 10 teknologi proyek mahasiswa |

### 3. Mata Kuliah, KRS (Enrollments), & Reviews
| Method | Endpoint | Akses | Keterangan |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/courses` | Terautentikasi | Daftar mata kuliah beserta sisa kuota |
| `GET` | `/api/v1/courses/mine` | Mahasiswa | Daftar mata kuliah yang diambil mahasiswa login |
| `POST` | `/api/v1/enrollments` | Mahasiswa | Ambil mata kuliah (transaksi atomik PostgreSQL) |
| `DELETE` | `/api/v1/enrollments/:id` | Mahasiswa | Batalkan KRS (kuota kembali bertambah) |
| `POST` | `/api/v1/enrollments/swap` | Mahasiswa | Tukar mata kuliah secara atomik dalam satu transaksi |
| `POST` | `/api/v1/courses/:id/reviews` | Mahasiswa | Beri ulasan dan rating mata kuliah (MongoDB) |
| `GET` | `/api/v1/courses/:id/reviews` | Terautentikasi | Daftar ulasan pada mata kuliah tertentu |
| `GET` | `/api/v1/courses/:id/reviews/summary` | Terautentikasi | Agregasi rata-rata rating dan distribusi bintang |
| `PATCH` | `/api/v1/reviews/:id` | Pemilik | Perbarui komentar/rating ulasan sendiri |
| `DELETE` | `/api/v1/reviews/:id` | Pemilik / Admin | Hapus ulasan dari MongoDB |
| `POST` | `/api/v1/reviews/:id/helpful` | Terautentikasi | Vote ulasan bermanfaat |
| `POST` | `/api/v1/reviews/:id/replies` | Dosen / Admin | Balas ulasan mahasiswa |

### 4. Pembimbing Tugas Akhir (Supervisions)
| Method | Endpoint | Akses | Keterangan |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/supervisors/:id/assign` | Admin / Permission | Penetapan dosen pembimbing tugas akhir |
| `PATCH` | `/api/v1/supervisors/:id` | Admin / Permission | Update kuota mahasiswa bimbingan |
| `POST` | `/api/v1/supervisions/switch` | Admin / Permission | Pengalihan pembimbing tugas akhir |

---

## 🧪 Pengujian & Koleksi Postman

Repository ini menyediakan koleksi pengujian siap pakai di folder `postman/`:
1. **`postman/SIAKAD_MINI_UTS.postman_collection.json`**: Pengujian komprehensif 10 endpoint relasional (UTS).
2. **`postman/SIAKAD_MINI_MODUL_10.postman_collection.json`**: Pengujian lengkap fitur NoSQL MongoDB (Portofolio, Ulasan, Agregasi, Kasus Error 401, 403, 404, 422, dan NoSQL Injection).

---

## 🧑‍💻 Kredensial Akun Seeder Standar

- **Administrator**:
  - Email: `admin@siakad.ac.id`
  - Password: `Admin123!`
- **Mahasiswa (20 Akun Seeder)**:
  - `rina.putri@siakad.ac.id` (NIM: `187221000001`, Password: `187221000001`) — *Student ID: 32*
  - `budi.santoso@siakad.ac.id` (NIM: `187221000002`, Password: `187221000002`)
  - `fajar.hidayat@siakad.ac.id` (NIM: `187221000006`, Password: `187221000006`) — *Student ID: 37*
  - ... hingga `taufik.hidayat@siakad.ac.id` (NIM: `187221000020`, Password: `187221000020`)
