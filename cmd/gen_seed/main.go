package main

import (
	"fmt"
	"os"

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
	f, err := os.Create("migrations/008_siakad_mini_seed.sql")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "-- Migration 008: SIAKAD Mini Seed Data")
	fmt.Fprintln(f, "BEGIN;")
	fmt.Fprintln(f)

	// Roles
	fmt.Fprintln(f, "-- 1. Roles")
	fmt.Fprintln(f, "INSERT INTO roles (name, description) VALUES")
	fmt.Fprintln(f, "    ('admin', 'Administrator SIAKAD Mini'),")
	fmt.Fprintln(f, "    ('mahasiswa', 'Mahasiswa SIAKAD Mini')")
	fmt.Fprintln(f, "ON CONFLICT (name) DO NOTHING;")
	fmt.Fprintln(f)

	// Admin
	adminHash, _ := bcrypt.GenerateFromPassword([]byte("Admin123!"), 10)
	fmt.Fprintln(f, "-- 2. Admin User (Password: Admin123!)")
	fmt.Fprintf(f, `INSERT INTO users (username, email, password, role, is_active)
VALUES ('admin_siakad', 'admin@siakad.ac.id', '%s', 'admin', true)
ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'admin';
`, string(adminHash))
	fmt.Fprintln(f)

	// Mahasiswa 20
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

	fmt.Fprintln(f, "-- 3. 20 Mahasiswa Users & Students (Password = NIM di-hash)")
	for _, s := range students {
		hash, _ := bcrypt.GenerateFromPassword([]byte(s.NIM), 10)
		fmt.Fprintf(f, `DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('%s', '%s', '%s', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = '%s';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '%s', '%s', '%s', %d, %.2f)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
`, s.NIM, s.Email, string(hash), s.Email, s.NIM, s.Nama, s.Prodi, s.Angkatan, s.IPKTerakhir)
	}
	fmt.Fprintln(f)

	// 10 Courses
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

	fmt.Fprintln(f, "-- 4. 10 Mata Kuliah")
	for _, c := range courses {
		fmt.Fprintf(f, `INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('%s', '%s', %d, %d, %d)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
`, c.KodeMK, c.NamaMK, c.SKS, c.Semester, c.Kuota)
	}

	fmt.Fprintln(f)
	fmt.Fprintln(f, "COMMIT;")
	fmt.Println("Generated migrations/008_siakad_mini_seed.sql successfully")
}
