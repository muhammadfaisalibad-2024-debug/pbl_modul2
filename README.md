# API Students

REST API sederhana untuk mengelola data mahasiswa menggunakan **Go** dan **Fiber v2**.

## Teknologi

- Go 1.26.5
- Fiber v2
- REST API
- Git & GitHub

## Menjalankan Project

Clone repository lalu masuk ke folder project:

```bash
git clone <URL_REPOSITORY>
cd api-students
```

Download dependency:

```bash
go mod tidy
```

Jalankan server:

```bash
go run .
```

Server berjalan di:

```text
http://localhost:3000
```

## Struktur Data Student

Setiap student memiliki data:

```json
{
  "id": 1,
  "nim": "0123456789",
  "name": "Faisal",
  "grade": 90,
  "is_active": true
}
```

# API Endpoint

Base URL:

```text
http://localhost:3000/api/v1
```

## 1. Get All Students

**GET**

```text
/api/v1/students
```

Contoh:

```text
GET /api/v1/students
```

### Pagination

```text
GET /api/v1/students?page=1&limit=2
```

### Search

```text
GET /api/v1/students?search=faisal
```

### Sorting

```text
GET /api/v1/students?sort=grade&order=desc
```

### Filter Status

```text
GET /api/v1/students?is_active=true
```

### Kombinasi Query

```text
GET /api/v1/students?page=1&limit=2&search=faisal&sort=grade&order=desc&is_active=true
```

### Response

```json
{
  "data": [
    {
      "id": 1,
      "nim": "0123456789",
      "name": "Faisal",
      "grade": 90,
      "is_active": true
    }
  ],
  "message": "Students retrieved successfully",
  "meta": {
    "limit": 2,
    "page": 1,
    "total": 1,
    "total_pages": 1
  }
}
```

## 2. Get Student by ID

**GET**

```text
/api/v1/students/:id
```

Contoh:

```text
GET /api/v1/students/1
```

Kemungkinan status:

- `200 OK` — data ditemukan
- `400 Bad Request` — ID tidak valid
- `404 Not Found` — student tidak ditemukan

## 3. Create Student

**POST**

```text
/api/v1/students
```

### Request Body

```json
{
  "nim": "0123456799",
  "name": "Citra",
  "grade": 92,
  "is_active": true
}
```

### Response

Status:

```text
201 Created
```

Response:

```json
{
  "success": true,
  "message": "Student created successfully",
  "data": {
    "id": 4,
    "nim": "0123456799",
    "name": "Citra",
    "grade": 92,
    "is_active": true
  }
}
```

Response juga menggunakan header:

```text
Location: /api/v1/students/4
```

Kemungkinan status:

- `201 Created`
- `409 Conflict`
- `415 Unsupported Media Type`
- `422 Unprocessable Entity`

## 4. Update Student

**PUT**

```text
/api/v1/students/:id
```

Contoh:

```text
PUT /api/v1/students/1
```

Request Body:

```json
{
  "nim": "0123456789",
  "name": "Faisal Updated",
  "grade": 95,
  "is_active": true
}
```

Kemungkinan status:

- `200 OK`
- `400 Bad Request`
- `404 Not Found`
- `415 Unsupported Media Type`
- `422 Unprocessable Entity`

## 5. Patch Student

**PATCH**

```text
/api/v1/students/:id
```

Contoh:

```text
PATCH /api/v1/students/1
```

Request Body:

```json
{
  "grade": 95
}
```

PATCH digunakan untuk mengubah sebagian data student.

Kemungkinan status:

- `200 OK`
- `400 Bad Request`
- `404 Not Found`
- `415 Unsupported Media Type`
- `422 Unprocessable Entity`

## 6. Delete Student

**DELETE**

```text
/api/v1/students/:id
```

Contoh:

```text
DELETE /api/v1/students/4
```

Jika berhasil:

```text
204 No Content
```

Response tidak memiliki body.

Kemungkinan status:

- `204 No Content`
- `400 Bad Request`
- `404 Not Found`

# HTTP Status yang Digunakan

| Status | Keterangan |
|---|---|
| 200 | Request berhasil |
| 201 | Data berhasil dibuat |
| 204 | Data berhasil dihapus tanpa response body |
| 400 | Request tidak valid |
| 404 | Data tidak ditemukan |
| 409 | Terjadi konflik, misalnya NIM sudah digunakan |
| 415 | Content-Type tidak sesuai |
| 422 | Data gagal validasi |

# Query Parameter

Endpoint:

```text
GET /api/v1/students
```

Mendukung parameter:

| Parameter | Fungsi | Contoh |
|---|---|---|
| `page` | Nomor halaman | `page=1` |
| `limit` | Jumlah data per halaman | `limit=2` |
| `search` | Mencari berdasarkan nama/NIM | `search=faisal` |
| `sort` | Field untuk sorting | `sort=grade` |
| `order` | Urutan data | `order=desc` |
| `is_active` | Filter status aktif | `is_active=true` |

# Testing

API diuji menggunakan **Postman**.

Pengujian meliputi:

- GET students
- GET student berdasarkan ID
- POST student
- PUT student
- PATCH student
- DELETE student
- Pagination
- Search
- Sorting
- Filtering
- HTTP status code
- HTTP response header

# Author

**Nama:** Faisal

**NIM:** 434241125
