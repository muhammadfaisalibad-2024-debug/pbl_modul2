-- Migration 007: SIAKAD Mini Schema

-- 1. Ensure roles exist
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator dengan hak akses penuh data mahasiswa dan mata kuliah'),
    ('mahasiswa', 'Mahasiswa dengan hak akses profil dan KRS pribadi')
ON CONFLICT (name) DO NOTHING;

-- 2. Adjust users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) NOT NULL DEFAULT 'mahasiswa';
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_role_fkey;
ALTER TABLE users ADD CONSTRAINT users_role_fkey FOREIGN KEY (role) REFERENCES roles(name) ON UPDATE CASCADE;

-- 3. Create or update students table
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    nim VARCHAR(12) UNIQUE NOT NULL,
    nama VARCHAR(120) NOT NULL,
    prodi VARCHAR(100) NOT NULL,
    angkatan INTEGER NOT NULL,
    ipk_terakhir NUMERIC(3,2) DEFAULT 0.00,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ensure all columns exist if students table already existed
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'user_id') THEN
        ALTER TABLE students ADD COLUMN user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'nama') THEN
        ALTER TABLE students ADD COLUMN nama VARCHAR(120);
        UPDATE students SET nama = name WHERE nama IS NULL AND name IS NOT NULL;
    END IF;
    ALTER TABLE students ALTER COLUMN name DROP NOT NULL;
    ALTER TABLE students ALTER COLUMN grade DROP NOT NULL;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'prodi') THEN
        ALTER TABLE students ADD COLUMN prodi VARCHAR(100) NOT NULL DEFAULT 'Sistem Informasi';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'angkatan') THEN
        ALTER TABLE students ADD COLUMN angkatan INTEGER NOT NULL DEFAULT 2024;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'ipk_terakhir') THEN
        ALTER TABLE students ADD COLUMN ipk_terakhir NUMERIC(3,2) DEFAULT 0.00;
        UPDATE students SET ipk_terakhir = CASE WHEN grade <= 4.0 THEN grade ELSE 3.50 END WHERE ipk_terakhir = 0.00 AND grade IS NOT NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'students' AND column_name = 'deleted_at') THEN
        ALTER TABLE students ADD COLUMN deleted_at TIMESTAMPTZ;
    END IF;
END $$;

-- 4. Create courses table
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    kode_mk VARCHAR(20) UNIQUE NOT NULL,
    nama_mk VARCHAR(120) NOT NULL,
    sks INTEGER NOT NULL CHECK (sks > 0),
    semester INTEGER NOT NULL CHECK (semester > 0),
    kuota INTEGER NOT NULL CHECK (kuota >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. Create enrollments table
CREATE TABLE IF NOT EXISTS enrollments (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    tahun_akademik VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_course_year UNIQUE (student_id, course_id, tahun_akademik)
);

-- 6. Indexes for performance
CREATE INDEX IF NOT EXISTS idx_students_user_id ON students(user_id);
CREATE INDEX IF NOT EXISTS idx_students_nim ON students(nim);
CREATE INDEX IF NOT EXISTS idx_students_deleted_at ON students(deleted_at);
CREATE INDEX IF NOT EXISTS idx_courses_kode_mk ON courses(kode_mk);
CREATE INDEX IF NOT EXISTS idx_courses_semester ON courses(semester);
CREATE INDEX IF NOT EXISTS idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_tahun_akademik ON enrollments(tahun_akademik);
