package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AutoMigrateAndSeed(db *pgxpool.Pool) error {
	ctx := context.Background()

	// Try reading schema migration file
	var schemaSql string
	for _, p := range []string{
		"migrations/007_siakad_mini_schema.sql",
		"../migrations/007_siakad_mini_schema.sql",
		"../../migrations/007_siakad_mini_schema.sql",
		"../../../migrations/007_siakad_mini_schema.sql",
	} {
		data, err := os.ReadFile(p)
		if err == nil {
			schemaSql = string(data)
			break
		}
	}

	if schemaSql != "" {
		if _, err := db.Exec(ctx, schemaSql); err != nil {
			return fmt.Errorf("gagal menjalankan migrasi skema: %w", err)
		}
	} else {
		// Fallback embedded schema SQL
		fallbackSchema := `
			INSERT INTO roles (name, description) VALUES
				('admin', 'Administrator SIAKAD Mini'),
				('mahasiswa', 'Mahasiswa SIAKAD Mini')
			ON CONFLICT (name) DO NOTHING;

			ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'mahasiswa';
			ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_fkey;
			ALTER TABLE users ADD CONSTRAINT users_role_fkey FOREIGN KEY (role) REFERENCES roles(name) ON UPDATE CASCADE;

			ALTER TABLE students ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
			ALTER TABLE students ADD COLUMN IF NOT EXISTS nama VARCHAR(120);
			UPDATE students SET nama = COALESCE(name, '') WHERE nama IS NULL;
			ALTER TABLE students ALTER COLUMN name DROP NOT NULL;
			ALTER TABLE students ALTER COLUMN grade DROP NOT NULL;
			ALTER TABLE students ADD COLUMN IF NOT EXISTS prodi VARCHAR(100) NOT NULL DEFAULT 'Sistem Informasi';
			ALTER TABLE students ADD COLUMN IF NOT EXISTS angkatan INTEGER NOT NULL DEFAULT 2024;
			ALTER TABLE students ADD COLUMN IF NOT EXISTS ipk_terakhir NUMERIC(3,2) DEFAULT 0.00;
			ALTER TABLE students ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

			CREATE TABLE IF NOT EXISTS courses (
				id SERIAL PRIMARY KEY,
				kode_mk VARCHAR(20) UNIQUE NOT NULL,
				nama_mk VARCHAR(120) NOT NULL,
				sks INTEGER NOT NULL CHECK (sks > 0),
				semester INTEGER NOT NULL CHECK (semester > 0),
				kuota INTEGER NOT NULL CHECK (kuota >= 0),
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS enrollments (
				id SERIAL PRIMARY KEY,
				student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
				course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
				tahun_akademik VARCHAR(30) NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				CONSTRAINT uq_student_course_year UNIQUE (student_id, course_id, tahun_akademik)
			);
		`
		if _, err := db.Exec(ctx, fallbackSchema); err != nil {
			return fmt.Errorf("gagal menjalankan fallback migrasi skema: %w", err)
		}
	}

	// Try reading seed file
	var seedSql string
	for _, p := range []string{
		"migrations/008_siakad_mini_seed.sql",
		"../migrations/008_siakad_mini_seed.sql",
		"../../migrations/008_siakad_mini_seed.sql",
		"../../../migrations/008_siakad_mini_seed.sql",
	} {
		data, err := os.ReadFile(p)
		if err == nil {
			seedSql = string(data)
			break
		}
	}

	if seedSql != "" {
		if _, err := db.Exec(ctx, seedSql); err != nil {
			return fmt.Errorf("gagal menjalankan seeder: %w", err)
		}
	}

	return nil
}
