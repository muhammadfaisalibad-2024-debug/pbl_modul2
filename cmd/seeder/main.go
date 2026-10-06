package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"api-students/config"
	"api-students/database"

	"golang.org/x/crypto/bcrypt"
)

type studentSeed struct {
	NIM         string
	Nama        string
	Email       string
	Prodi       string
	Angkatan    int
	IPKTerakhir float64
}

type courseSeed struct {
	KodeMK   string
	NamaMK   string
	SKS      int
	Semester int
	Kuota    int
}

func main() {
	cfg := config.Load()
	log.Printf("Connecting to DB: %s@%s:%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := database.NewDBPool(cfg)
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	var currentDB, currentSchema, searchPath string
	_ = db.QueryRow(ctx, "SELECT current_database(), current_schema(), current_setting('search_path')").Scan(&currentDB, &currentSchema, &searchPath)
	log.Printf("DB: %s, Schema: %s, SearchPath: %s", currentDB, currentSchema, searchPath)

	log.Println("Memulai proses seeding SIAKAD Mini...")

	// 1. Seed Roles
	_, err = db.Exec(ctx, `
		INSERT INTO roles (name, description) VALUES
			('admin', 'Administrator SIAKAD Mini'),
			('mahasiswa', 'Mahasiswa SIAKAD Mini')
		ON CONFLICT (name) DO NOTHING;
	`)
	if err != nil {
		log.Fatalf("Gagal seed roles: %v", err)
	}

	// 2. Seed Admin User
	adminPassHash, err := bcrypt.GenerateFromPassword([]byte("Admin123!"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password admin: %v", err)
	}

	var adminID int
	err = db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, "admin@siakad.ac.id").Scan(&adminID)
	if err != nil {
		// Insert admin baru
		err = db.QueryRow(ctx, `
			INSERT INTO users (username, email, password, role, is_active)
			VALUES ($1, $2, $3, 'admin', true)
			RETURNING id;
		`, "admin_siakad", "admin@siakad.ac.id", string(adminPassHash)).Scan(&adminID)
		if err != nil {
			log.Fatalf("Gagal insert admin: %v", err)
		}
	} else {
		// Update admin yang ada
		_, err = db.Exec(ctx, `
			UPDATE users SET password = $1, role = 'admin', is_active = true WHERE id = $2;
		`, string(adminPassHash), adminID)
		if err != nil {
			log.Fatalf("Gagal update admin: %v", err)
		}
	}
	log.Printf("Admin terdaftar (ID: %d, Email: admin@siakad.ac.id, Password: Admin123!)", adminID)

	// 3. 20 Mahasiswa Seed Data
	students := []studentSeed{
		{"187221000001", "Rina Putri", "rina.putri@siakad.ac.id", "Sistem Informasi", 2022, 3.45},
		{"187221000002", "Budi Santoso", "budi.santoso@siakad.ac.id", "Informatika", 2022, 3.80},
		{"187221000003", "Citra Lestari", "citra.lestari@siakad.ac.id", "Sistem Informasi", 2023, 2.85},
		{"187221000004", "Dimas Pratama", "dimas.pratama@siakad.ac.id", "Teknologi Informasi", 2023, 2.40},
		{"187221000005", "Eka Saputri", "eka.saputri@siakad.ac.id", "Informatika", 2023, 3.10},
		{"187221000006", "Fajar Hidayat", "fajar.hidayat@siakad.ac.id", "Sistem Informasi", 2024, 3.65},
		{"187221000007", "Gita Gutawa", "gita.gutawa@siakad.ac.id", "Teknologi Informasi", 2024, 2.90},
		{"187221000008", "Hadi Wijaya", "hadi.wijaya@siakad.ac.id", "Informatika", 2024, 2.30},
		{"187221000009", "Indah Permata", "indah.permata@siakad.ac.id", "Sistem Informasi", 2024, 3.75},
		{"187221000010", "Joko Susilo", "joko.susilo@siakad.ac.id", "Informatika", 2024, 3.20},
		{"187221000011", "Kartika Sari", "kartika.sari@siakad.ac.id", "Teknologi Informasi", 2024, 3.05},
		{"187221000012", "Lukman Hakim", "lukman.hakim@siakad.ac.id", "Sistem Informasi", 2025, 2.70},
		{"187221000013", "Maya Anggraini", "maya.anggraini@siakad.ac.id", "Informatika", 2025, 3.50},
		{"187221000014", "Nanda Pratama", "nanda.pratama@siakad.ac.id", "Teknologi Informasi", 2025, 2.15},
		{"187221000015", "Oki Setiawan", "oki.setiawan@siakad.ac.id", "Sistem Informasi", 2025, 3.90},
		{"187221000016", "Putri Rahayu", "putri.rahayu@siakad.ac.id", "Informatika", 2025, 3.35},
		{"187221000017", "Qori Sandioriva", "qori.sandioriva@siakad.ac.id", "Teknologi Informasi", 2025, 2.95},
		{"187221000018", "Rizky Ramadhan", "rizky.ramadhan@siakad.ac.id", "Sistem Informasi", 2025, 3.60},
		{"187221000019", "Siti Nurhaliza", "siti.nurhaliza@siakad.ac.id", "Informatika", 2025, 3.15},
		{"187221000020", "Taufik Hidayat", "taufik.hidayat@siakad.ac.id", "Teknologi Informasi", 2025, 2.45},
	}

	for _, s := range students {
		// Default password = NIM (di-hash)
		nimHash, err := bcrypt.GenerateFromPassword([]byte(s.NIM), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Gagal hash NIM %s: %v", s.NIM, err)
		}

		var uID int
		err = db.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, s.Email).Scan(&uID)
		if err != nil {
			// Also check by username (NIM)
			err = db.QueryRow(ctx, `SELECT id FROM users WHERE username = $1`, s.NIM).Scan(&uID)
		}

		if err != nil {
			err = db.QueryRow(ctx, `
				INSERT INTO users (username, email, password, role, is_active)
				VALUES ($1, $2, $3, 'mahasiswa', true)
				RETURNING id;
			`, s.NIM, s.Email, string(nimHash)).Scan(&uID)
			if err != nil {
				log.Fatalf("Gagal insert user mahasiswa %s: %v", s.Email, err)
			}
		} else {
			_, err = db.Exec(ctx, `
				UPDATE users SET username = $1, email = $2, password = $3, role = 'mahasiswa', is_active = true
				WHERE id = $4;
			`, s.NIM, s.Email, string(nimHash), uID)
			if err != nil {
				log.Fatalf("Gagal update user mahasiswa %s: %v", s.Email, err)
			}
		}

		// Insert or update student record
		var sID int
		err = db.QueryRow(ctx, `SELECT id FROM students WHERE nim = $1`, s.NIM).Scan(&sID)
		if err != nil {
			_, err = db.Exec(ctx, `
				INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
				VALUES ($1, $2, $3, $4, $5, $6);
			`, uID, s.NIM, s.Nama, s.Prodi, s.Angkatan, s.IPKTerakhir)
			if err != nil {
				log.Fatalf("Gagal insert student %s: %v", s.NIM, err)
			}
		} else {
			_, err = db.Exec(ctx, `
				UPDATE students SET
					user_id = $1,
					nama = $2,
					prodi = $3,
					angkatan = $4,
					ipk_terakhir = $5,
					deleted_at = NULL
				WHERE id = $6;
			`, uID, s.Nama, s.Prodi, s.Angkatan, s.IPKTerakhir, sID)
			if err != nil {
				log.Fatalf("Gagal update student %s: %v", s.NIM, err)
			}
		}
	}
	log.Printf("Berhasil seed %d mahasiswa!", len(students))

	// 4. 10 Mata Kuliah Seed Data
	courses := []courseSeed{
		{"MK001", "Pemrograman Web", 3, 3, 30},
		{"MK002", "Basis Data Lanjut", 3, 3, 25},
		{"MK003", "Algoritma dan Pemrograman", 4, 1, 35},
		{"MK004", "Jaringan Komputer", 3, 3, 20},
		{"MK005", "Rekayasa Perangkat Lunak", 3, 5, 30},
		{"MK006", "Kecerdasan Buatan", 3, 5, 25},
		{"MK007", "Keamanan Informasi", 3, 5, 20},
		{"MK008", "Cloud Computing", 3, 7, 20},
		{"MK009", "Proyek Sistem Informasi", 4, 7, 15},
		{"MK010", "Etika Profesi IT", 2, 1, 40},
	}

	for _, c := range courses {
		var cID int
		err = db.QueryRow(ctx, `SELECT id FROM courses WHERE kode_mk = $1`, c.KodeMK).Scan(&cID)
		if err != nil {
			_, err = db.Exec(ctx, `
				INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
				VALUES ($1, $2, $3, $4, $5);
			`, c.KodeMK, c.NamaMK, c.SKS, c.Semester, c.Kuota)
			if err != nil {
				log.Fatalf("Gagal insert course %s: %v", c.KodeMK, err)
			}
		} else {
			_, err = db.Exec(ctx, `
				UPDATE courses SET
					nama_mk = $1,
					sks = $2,
					semester = $3,
					kuota = $4
				WHERE id = $5;
			`, c.NamaMK, c.SKS, c.Semester, c.Kuota, cID)
			if err != nil {
				log.Fatalf("Gagal update course %s: %v", c.KodeMK, err)
			}
		}
	}
	log.Printf("Berhasil seed %d mata kuliah!", len(courses))

	fmt.Println("=== Seeding Selesai dengan Sukses ===")
	os.Exit(0)
}
