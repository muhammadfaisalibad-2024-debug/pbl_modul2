-- Migration 008: SIAKAD Mini Seed Data
BEGIN;

-- 1. Roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator SIAKAD Mini'),
    ('mahasiswa', 'Mahasiswa SIAKAD Mini')
ON CONFLICT (name) DO NOTHING;

-- 2. Admin User (Password: Admin123!)
INSERT INTO users (username, email, password, role, is_active)
VALUES ('admin_siakad', 'admin@siakad.ac.id', '$2a$10$d3.r0ATVtFu6jprfakixKuO/0vU6YgrV9GyTuFPCQb7KKhyoXms0u', 'admin', true)
ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'admin';

-- 3. 20 Mahasiswa Users & Students (Password = NIM di-hash)
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000001', 'rina.putri@siakad.ac.id', '$2a$10$a2paBOUDGKhWI3rcDYktkOV5To1DFzdsWTME6B6AM8zHO1/QwQ7ju', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'rina.putri@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000001', 'Rina Putri', 'Sistem Informasi', 2022, 3.45)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000002', 'budi.santoso@siakad.ac.id', '$2a$10$HLIGXeMfdwFc6vUWFKWlXuirXvQD0QFe1/6oFSihnbcpT/hf2z/rC', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'budi.santoso@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000002', 'Budi Santoso', 'Informatika', 2022, 3.80)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000003', 'citra.lestari@siakad.ac.id', '$2a$10$r5enxyTIZczsHWrS1b2WtereYCq69.PQqIICkcr3mZwNQZvkA6l8K', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'citra.lestari@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000003', 'Citra Lestari', 'Sistem Informasi', 2023, 2.85)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000004', 'dimas.pratama@siakad.ac.id', '$2a$10$TVWY/WXiI767E6HZH957SOCagmR/IG39WhmCKTFUeW1kC9dZLrZwm', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'dimas.pratama@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000004', 'Dimas Pratama', 'Teknologi Informasi', 2023, 2.40)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000005', 'eka.saputri@siakad.ac.id', '$2a$10$g6yFIZSMqry6sZ5Eny/qb..GtaZsx1/D09tDko42IycltSJmzMfIW', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'eka.saputri@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000005', 'Eka Saputri', 'Informatika', 2023, 3.10)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000006', 'fajar.hidayat@siakad.ac.id', '$2a$10$vePDpX5F3Hv2BspQ4VeRreR36fuvrLSKUP9tNwQbYk6FCxhDMT5Ji', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'fajar.hidayat@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000006', 'Fajar Hidayat', 'Sistem Informasi', 2024, 3.65)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000007', 'gita.gutawa@siakad.ac.id', '$2a$10$CcLo1S532qliRAgZHWjJsuOFnTo2M3UHWUNk1zlkW83NLMx8BebNW', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'gita.gutawa@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000007', 'Gita Gutawa', 'Teknologi Informasi', 2024, 2.90)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000008', 'hadi.wijaya@siakad.ac.id', '$2a$10$lTsZc66jRrqra4tcwI34NuvM92GD5EHFZfCr/MSawCVw56fd2za2G', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'hadi.wijaya@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000008', 'Hadi Wijaya', 'Informatika', 2024, 2.30)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000009', 'indah.permata@siakad.ac.id', '$2a$10$MSNygqbSE1mQPC2/sJZCG.NuLY9zWYYF9iIxjfKdkPafXbMD/FtX2', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'indah.permata@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000009', 'Indah Permata', 'Sistem Informasi', 2024, 3.75)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000010', 'joko.susilo@siakad.ac.id', '$2a$10$giTX3QsYx.nFMz2C6dA7XeiHJkNmz3RoHT/xukTkYIgOue7T.piQC', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'joko.susilo@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000010', 'Joko Susilo', 'Informatika', 2024, 3.20)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000011', 'kartika.sari@siakad.ac.id', '$2a$10$D5JR09Afg6E7y5DYEnBYZe7YFqftYafL0orYgWvNoaI0s1Ul7uNdm', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'kartika.sari@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000011', 'Kartika Sari', 'Teknologi Informasi', 2024, 3.05)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000012', 'lukman.hakim@siakad.ac.id', '$2a$10$BAM2IwL3GhDvGVLMVJuiWufqiuVREG5RdyOyGhhzf6LDqyIHdT5/6', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'lukman.hakim@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000012', 'Lukman Hakim', 'Sistem Informasi', 2025, 2.70)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000013', 'maya.anggraini@siakad.ac.id', '$2a$10$kqQDTBnG0QO4uPs65FmmeOKLvCwgloUIh/U.9YjBlnZ26Kq48oQL2', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'maya.anggraini@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000013', 'Maya Anggraini', 'Informatika', 2025, 3.50)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000014', 'nanda.pratama@siakad.ac.id', '$2a$10$04I34HJnPv3w3vehB9vJs.ppKpQFHoJdez5/YtVF2iXZFW3c27/OG', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'nanda.pratama@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000014', 'Nanda Pratama', 'Teknologi Informasi', 2025, 2.15)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000015', 'oki.setiawan@siakad.ac.id', '$2a$10$.s5QXAD1q.sxDraH9mu.ZuR.fP.8Ra7mdwv7rB7QYxwM0VPAcp3oi', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'oki.setiawan@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000015', 'Oki Setiawan', 'Sistem Informasi', 2025, 3.90)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000016', 'putri.rahayu@siakad.ac.id', '$2a$10$6M2pIt45hrrVXnWSGT9U/exW1N8aA1g4w6/QHqaIeWcF6AbvXXx6W', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'putri.rahayu@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000016', 'Putri Rahayu', 'Informatika', 2025, 3.35)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000017', 'qori.sandioriva@siakad.ac.id', '$2a$10$gKE.WrEwWPNyymvY5lITmO.0U9EIsb2YxOndJ3HB2nVJvlkS.Nwra', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'qori.sandioriva@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000017', 'Qori Sandioriva', 'Teknologi Informasi', 2025, 2.95)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000018', 'rizky.ramadhan@siakad.ac.id', '$2a$10$rFC3eOM9vhAE1qsIa36Y.uFJEd9aUEmdU6A5LASpeNV5Vax9bPPiq', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'rizky.ramadhan@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000018', 'Rizky Ramadhan', 'Sistem Informasi', 2025, 3.60)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000019', 'siti.nurhaliza@siakad.ac.id', '$2a$10$foNd1RFkVIJlIYaZVinFneU6PgcgOUs3EE4n9MD0pUoKPj0d827yC', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'siti.nurhaliza@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000019', 'Siti Nurhaliza', 'Informatika', 2025, 3.15)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
BEGIN
    INSERT INTO users (username, email, password, role, is_active)
    VALUES ('187221000020', 'taufik.hidayat@siakad.ac.id', '$2a$10$Z/3Lw4QYdIeWNeNvOdNjVee2otNNrvsbqNGfiJimr97rr7XzbhV4.', 'mahasiswa', true)
    ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password, role = 'mahasiswa'
    RETURNING id INTO v_user_id;

    IF v_user_id IS NULL THEN
        SELECT id INTO v_user_id FROM users WHERE email = 'taufik.hidayat@siakad.ac.id';
    END IF;

    INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
    VALUES (v_user_id, '187221000020', 'Taufik Hidayat', 'Teknologi Informasi', 2025, 2.45)
    ON CONFLICT (nim) DO UPDATE SET
        user_id = EXCLUDED.user_id,
        nama = EXCLUDED.nama,
        prodi = EXCLUDED.prodi,
        angkatan = EXCLUDED.angkatan,
        ipk_terakhir = EXCLUDED.ipk_terakhir,
        deleted_at = NULL;
END $$;

-- 4. 10 Mata Kuliah
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK001', 'Pemrograman Web', 3, 3, 30)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK002', 'Basis Data Lanjut', 3, 3, 25)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK003', 'Algoritma dan Pemrograman', 4, 1, 35)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK004', 'Jaringan Komputer', 3, 3, 20)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK005', 'Rekayasa Perangkat Lunak', 3, 5, 30)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK006', 'Kecerdasan Buatan', 3, 5, 25)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK007', 'Keamanan Informasi', 3, 5, 20)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK008', 'Cloud Computing', 3, 7, 20)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK009', 'Proyek Sistem Informasi', 4, 7, 15)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
VALUES ('MK010', 'Etika Profesi IT', 2, 1, 40)
ON CONFLICT (kode_mk) DO UPDATE SET
    nama_mk = EXCLUDED.nama_mk,
    sks = EXCLUDED.sks,
    semester = EXCLUDED.semester,
    kuota = EXCLUDED.kuota;

COMMIT;
