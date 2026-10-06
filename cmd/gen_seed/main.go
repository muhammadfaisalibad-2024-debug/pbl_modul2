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
	fmt.Fprintf(f, `DO $$
DECLARE
    v_admin_id INTEGER;
BEGIN
    SELECT id INTO v_admin_id FROM users WHERE email = 'admin@siakad.ac.id';
    IF v_admin_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('admin_siakad', 'admin@siakad.ac.id', '%s', 'admin', true);
    ELSE
        UPDATE users SET password = '%s', role = 'admin', is_active = true WHERE id = v_admin_id;
    END IF;
END $$;
`, string(adminHash), string(adminHash))
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
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = '%s';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('%s', '%s', '%s', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '%s', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '%s';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '%s', '%s', '%s', %.2f, '%s', %d, %.2f);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = '%s',
            name = '%s',
            grade = %.2f,
            prodi = '%s',
            angkatan = %d,
            ipk_terakhir = %.2f,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
`, s.Email, s.NIM, s.Email, string(hash), string(hash), s.NIM, s.NIM, s.Nama, s.Nama, s.IPKTerakhir, s.Prodi, s.Angkatan, s.IPKTerakhir, s.Nama, s.Nama, s.IPKTerakhir, s.Prodi, s.Angkatan, s.IPKTerakhir)
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
		fmt.Fprintf(f, `DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = '%s';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('%s', '%s', %d, %d, %d);
    ELSE
        UPDATE courses SET
            nama_mk = '%s',
            sks = %d,
            semester = %d,
            kuota = %d
        WHERE id = v_course_id;
    END IF;
END $$;
`, c.KodeMK, c.KodeMK, c.NamaMK, c.SKS, c.Semester, c.Kuota, c.NamaMK, c.SKS, c.Semester, c.Kuota)
	}

	fmt.Fprintln(f)
	fmt.Fprintln(f, "COMMIT;")
	fmt.Println("Generated migrations/008_siakad_mini_seed.sql successfully")
}
