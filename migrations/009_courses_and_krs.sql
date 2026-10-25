-- Migration 009: courses and KRS (Modul 9)
ALTER TABLE courses ADD COLUMN IF NOT EXISTS code VARCHAR(20) UNIQUE;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS name VARCHAR(120);
ALTER TABLE courses ADD COLUMN IF NOT EXISTS credits SMALLINT CHECK (credits BETWEEN 1 AND 6);
ALTER TABLE courses ADD COLUMN IF NOT EXISTS quota INTEGER CHECK (quota > 0);
ALTER TABLE courses ADD COLUMN IF NOT EXISTS enrolled_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE courses ADD COLUMN IF NOT EXISTS version INTEGER NOT NULL DEFAULT 1;

UPDATE courses SET
    code = COALESCE(code, kode_mk),
    name = COALESCE(name, nama_mk),
    credits = COALESCE(credits, sks::smallint),
    quota = COALESCE(quota, kuota)
WHERE code IS NULL;

ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_enrolled_within_quota;
ALTER TABLE courses ADD CONSTRAINT courses_enrolled_within_quota
    CHECK (enrolled_count >= 0 AND enrolled_count <= quota);

-- Adjust enrollments for Modul 9 (supports user_id)
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE enrollments ALTER COLUMN student_id DROP NOT NULL;
ALTER TABLE enrollments ALTER COLUMN tahun_akademik DROP NOT NULL;
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_course_user_unique;
ALTER TABLE enrollments ADD CONSTRAINT enrollments_course_user_unique UNIQUE (course_id, user_id);
CREATE INDEX IF NOT EXISTS enrollments_user_idx ON enrollments (user_id);

-- Insert courses IF101 to IF399
INSERT INTO courses (code, name, credits, quota, semester, kode_mk, nama_mk, sks, kuota) VALUES
 ('IF101', 'Pemrograman Backend Lanjut', 3, 40, 1, 'IF101', 'Pemrograman Backend Lanjut', 3, 40),
 ('IF102', 'Basis Data Lanjut', 3, 40, 1, 'IF102', 'Basis Data Lanjut', 3, 40),
 ('IF301', 'Kecerdasan Buatan', 3, 40, 3, 'IF301', 'Kecerdasan Buatan', 3, 40),
 ('IF302', 'Keamanan Informasi', 3, 40, 3, 'IF302', 'Keamanan Informasi', 3, 40),
 ('IF303', 'Komputasi Awan', 3, 40, 3, 'IF303', 'Komputasi Awan', 3, 40),
 ('IF304', 'Data Mining', 3, 40, 3, 'IF304', 'Data Mining', 3, 40),
 ('IF305', 'Pemrograman Mobile', 3, 40, 3, 'IF305', 'Pemrograman Mobile', 3, 40),
 ('IF306', 'Interaksi Manusia Komputer', 3, 40, 3, 'IF306', 'Interaksi Manusia Komputer', 3, 40),
 ('IF307', 'Sistem Terdistribusi', 3, 40, 3, 'IF307', 'Sistem Terdistribusi', 3, 40),
 ('IF308', 'Visi Komputer', 3, 40, 3, 'IF308', 'Visi Komputer', 3, 40),
 ('IF309', 'Pemrosesan Bahasa Alami', 3, 40, 3, 'IF309', 'Pemrosesan Bahasa Alami', 3, 40),
 ('IF310', 'Manajemen Proyek TI', 3, 40, 3, 'IF310', 'Manajemen Proyek TI', 3, 40),
 ('IF399', 'Kelas Rebutan', 3, 5, 3, 'IF399', 'Kelas Rebutan', 3, 5)
ON CONFLICT (code) DO NOTHING;

-- Permissions
INSERT INTO permissions (name, description) VALUES
 ('course:manage', 'Mengubah kuota dan data kelas')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
 ('admin', 'course:manage')
ON CONFLICT DO NOTHING;
