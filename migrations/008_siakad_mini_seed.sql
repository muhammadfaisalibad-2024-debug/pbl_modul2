-- Migration 008: SIAKAD Mini Seed Data
BEGIN;

-- 1. Roles
INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator SIAKAD Mini'),
    ('mahasiswa', 'Mahasiswa SIAKAD Mini')
ON CONFLICT (name) DO NOTHING;

-- 2. Admin User (Password: Admin123!)
DO $$
DECLARE
    v_admin_id INTEGER;
BEGIN
    SELECT id INTO v_admin_id FROM users WHERE email = 'admin@siakad.ac.id';
    IF v_admin_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('admin_siakad', 'admin@siakad.ac.id', '$2a$10$PNULYYo5o63Vxs2nohq8AezksasbhytKGCdAOarOl6ORNcyMDKk6O', 'admin', true);
    ELSE
        UPDATE users SET password = '$2a$10$PNULYYo5o63Vxs2nohq8AezksasbhytKGCdAOarOl6ORNcyMDKk6O', role = 'admin', is_active = true WHERE id = v_admin_id;
    END IF;
END $$;

-- 3. 20 Mahasiswa Users & Students (Password = NIM di-hash)
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'rina.putri@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000001', 'rina.putri@siakad.ac.id', '$2a$10$H9bZc7juK8n1yWQC45gCj.vHr9RirQUjjDimN9Y0MNKpM/tSQ8PhG', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$H9bZc7juK8n1yWQC45gCj.vHr9RirQUjjDimN9Y0MNKpM/tSQ8PhG', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000001';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000001', 'Rina Putri', 'Rina Putri', 3.45, 'Sistem Informasi', 2022, 3.45);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Rina Putri',
            name = 'Rina Putri',
            grade = 3.45,
            prodi = 'Sistem Informasi',
            angkatan = 2022,
            ipk_terakhir = 3.45,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'budi.santoso@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000002', 'budi.santoso@siakad.ac.id', '$2a$10$RetMpB54vbUkuRaNp/1QNep6ZgFC9J9odFfK2YMMDO29zrcWHuDHi', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$RetMpB54vbUkuRaNp/1QNep6ZgFC9J9odFfK2YMMDO29zrcWHuDHi', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000002';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000002', 'Budi Santoso', 'Budi Santoso', 3.80, 'Informatika', 2022, 3.80);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Budi Santoso',
            name = 'Budi Santoso',
            grade = 3.80,
            prodi = 'Informatika',
            angkatan = 2022,
            ipk_terakhir = 3.80,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'citra.lestari@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000003', 'citra.lestari@siakad.ac.id', '$2a$10$IG0GSyyu4o36R038AD8nAealJAl5nfsNd4wviCJdM0z0zmuJCMHHC', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$IG0GSyyu4o36R038AD8nAealJAl5nfsNd4wviCJdM0z0zmuJCMHHC', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000003';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000003', 'Citra Lestari', 'Citra Lestari', 2.85, 'Sistem Informasi', 2023, 2.85);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Citra Lestari',
            name = 'Citra Lestari',
            grade = 2.85,
            prodi = 'Sistem Informasi',
            angkatan = 2023,
            ipk_terakhir = 2.85,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'dimas.pratama@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000004', 'dimas.pratama@siakad.ac.id', '$2a$10$ULng.HLeK0uu9f3HVA4W6OQAAxalaHD3JDG/0A0XdeexalZ60kJIO', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$ULng.HLeK0uu9f3HVA4W6OQAAxalaHD3JDG/0A0XdeexalZ60kJIO', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000004';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000004', 'Dimas Pratama', 'Dimas Pratama', 2.40, 'Teknologi Informasi', 2023, 2.40);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Dimas Pratama',
            name = 'Dimas Pratama',
            grade = 2.40,
            prodi = 'Teknologi Informasi',
            angkatan = 2023,
            ipk_terakhir = 2.40,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'eka.saputri@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000005', 'eka.saputri@siakad.ac.id', '$2a$10$/pDh4ybeSMAs4T1g9W99LejCUl6Zrb0cJOpDr6b1aARdig/6cMpjS', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$/pDh4ybeSMAs4T1g9W99LejCUl6Zrb0cJOpDr6b1aARdig/6cMpjS', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000005';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000005', 'Eka Saputri', 'Eka Saputri', 3.10, 'Informatika', 2023, 3.10);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Eka Saputri',
            name = 'Eka Saputri',
            grade = 3.10,
            prodi = 'Informatika',
            angkatan = 2023,
            ipk_terakhir = 3.10,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'fajar.hidayat@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000006', 'fajar.hidayat@siakad.ac.id', '$2a$10$..rrMSXjv25gup4eUOv3feeNx4j/xBT.Rq38T4ieScq9ROVipOgOW', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$..rrMSXjv25gup4eUOv3feeNx4j/xBT.Rq38T4ieScq9ROVipOgOW', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000006';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000006', 'Fajar Hidayat', 'Fajar Hidayat', 3.65, 'Sistem Informasi', 2024, 3.65);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Fajar Hidayat',
            name = 'Fajar Hidayat',
            grade = 3.65,
            prodi = 'Sistem Informasi',
            angkatan = 2024,
            ipk_terakhir = 3.65,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'gita.gutawa@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000007', 'gita.gutawa@siakad.ac.id', '$2a$10$VuVldKE5MLtip1k2epwmcuBc7eEhVUBxPc4VnnJyz7j9xl.jJivqq', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$VuVldKE5MLtip1k2epwmcuBc7eEhVUBxPc4VnnJyz7j9xl.jJivqq', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000007';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000007', 'Gita Gutawa', 'Gita Gutawa', 2.90, 'Teknologi Informasi', 2024, 2.90);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Gita Gutawa',
            name = 'Gita Gutawa',
            grade = 2.90,
            prodi = 'Teknologi Informasi',
            angkatan = 2024,
            ipk_terakhir = 2.90,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'hadi.wijaya@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000008', 'hadi.wijaya@siakad.ac.id', '$2a$10$Hv4nVDpAXJjdTRcrrhrY1OpNby3neE9tisEQdEbYqYP6yv.rICZV.', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$Hv4nVDpAXJjdTRcrrhrY1OpNby3neE9tisEQdEbYqYP6yv.rICZV.', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000008';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000008', 'Hadi Wijaya', 'Hadi Wijaya', 2.30, 'Informatika', 2024, 2.30);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Hadi Wijaya',
            name = 'Hadi Wijaya',
            grade = 2.30,
            prodi = 'Informatika',
            angkatan = 2024,
            ipk_terakhir = 2.30,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'indah.permata@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000009', 'indah.permata@siakad.ac.id', '$2a$10$cudIRJbixkpd2MnkgaR7G.Br4uGQ32WqFoCLOwybmwDo0IhHA41LS', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$cudIRJbixkpd2MnkgaR7G.Br4uGQ32WqFoCLOwybmwDo0IhHA41LS', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000009';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000009', 'Indah Permata', 'Indah Permata', 3.75, 'Sistem Informasi', 2024, 3.75);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Indah Permata',
            name = 'Indah Permata',
            grade = 3.75,
            prodi = 'Sistem Informasi',
            angkatan = 2024,
            ipk_terakhir = 3.75,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'joko.susilo@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000010', 'joko.susilo@siakad.ac.id', '$2a$10$bqyFnp1D6/GZm4bDnpJjleGKfPnzgDbrRiZr9SaCF9vnFaxOvPrQC', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$bqyFnp1D6/GZm4bDnpJjleGKfPnzgDbrRiZr9SaCF9vnFaxOvPrQC', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000010';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000010', 'Joko Susilo', 'Joko Susilo', 3.20, 'Informatika', 2024, 3.20);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Joko Susilo',
            name = 'Joko Susilo',
            grade = 3.20,
            prodi = 'Informatika',
            angkatan = 2024,
            ipk_terakhir = 3.20,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'kartika.sari@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000011', 'kartika.sari@siakad.ac.id', '$2a$10$HYHqM3FnNYzPaO2383ft1OhK3HwNVmoapybZDQsWTervkP9/o464i', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$HYHqM3FnNYzPaO2383ft1OhK3HwNVmoapybZDQsWTervkP9/o464i', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000011';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000011', 'Kartika Sari', 'Kartika Sari', 3.05, 'Teknologi Informasi', 2024, 3.05);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Kartika Sari',
            name = 'Kartika Sari',
            grade = 3.05,
            prodi = 'Teknologi Informasi',
            angkatan = 2024,
            ipk_terakhir = 3.05,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'lukman.hakim@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000012', 'lukman.hakim@siakad.ac.id', '$2a$10$1gCQ37ZQGIPf4Em7n/A08eBW.9bor.9FV3VjKNjNO6QQdfcVWNAPe', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$1gCQ37ZQGIPf4Em7n/A08eBW.9bor.9FV3VjKNjNO6QQdfcVWNAPe', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000012';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000012', 'Lukman Hakim', 'Lukman Hakim', 2.70, 'Sistem Informasi', 2025, 2.70);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Lukman Hakim',
            name = 'Lukman Hakim',
            grade = 2.70,
            prodi = 'Sistem Informasi',
            angkatan = 2025,
            ipk_terakhir = 2.70,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'maya.anggraini@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000013', 'maya.anggraini@siakad.ac.id', '$2a$10$Ntld31U.mRXD5UPmNfhQnOlkbUUYAeAtEIhpDlU.U4a1g0aLYH.RW', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$Ntld31U.mRXD5UPmNfhQnOlkbUUYAeAtEIhpDlU.U4a1g0aLYH.RW', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000013';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000013', 'Maya Anggraini', 'Maya Anggraini', 3.50, 'Informatika', 2025, 3.50);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Maya Anggraini',
            name = 'Maya Anggraini',
            grade = 3.50,
            prodi = 'Informatika',
            angkatan = 2025,
            ipk_terakhir = 3.50,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'nanda.pratama@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000014', 'nanda.pratama@siakad.ac.id', '$2a$10$zULZsAUoG6ial7TukeN7yuUk/uNfbd5rtUPqrpXwpkwvXiDGyotku', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$zULZsAUoG6ial7TukeN7yuUk/uNfbd5rtUPqrpXwpkwvXiDGyotku', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000014';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000014', 'Nanda Pratama', 'Nanda Pratama', 2.15, 'Teknologi Informasi', 2025, 2.15);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Nanda Pratama',
            name = 'Nanda Pratama',
            grade = 2.15,
            prodi = 'Teknologi Informasi',
            angkatan = 2025,
            ipk_terakhir = 2.15,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'oki.setiawan@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000015', 'oki.setiawan@siakad.ac.id', '$2a$10$rfrMQcb7tUYqlQ7AMGLB..WTPfClmKPh2OGq18F1P5/KeBV21LxKy', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$rfrMQcb7tUYqlQ7AMGLB..WTPfClmKPh2OGq18F1P5/KeBV21LxKy', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000015';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000015', 'Oki Setiawan', 'Oki Setiawan', 3.90, 'Sistem Informasi', 2025, 3.90);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Oki Setiawan',
            name = 'Oki Setiawan',
            grade = 3.90,
            prodi = 'Sistem Informasi',
            angkatan = 2025,
            ipk_terakhir = 3.90,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'putri.rahayu@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000016', 'putri.rahayu@siakad.ac.id', '$2a$10$MQzx2nQWHjvQ6AO55kfTDe6HP5Nr3BY68W.QaO7xPX/xKOACXi.fy', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$MQzx2nQWHjvQ6AO55kfTDe6HP5Nr3BY68W.QaO7xPX/xKOACXi.fy', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000016';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000016', 'Putri Rahayu', 'Putri Rahayu', 3.35, 'Informatika', 2025, 3.35);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Putri Rahayu',
            name = 'Putri Rahayu',
            grade = 3.35,
            prodi = 'Informatika',
            angkatan = 2025,
            ipk_terakhir = 3.35,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'qori.sandioriva@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000017', 'qori.sandioriva@siakad.ac.id', '$2a$10$MozZL0.HyVcYr9UQio0GTugfKetz1cWm4fA78qv20istUguz.kn7y', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$MozZL0.HyVcYr9UQio0GTugfKetz1cWm4fA78qv20istUguz.kn7y', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000017';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000017', 'Qori Sandioriva', 'Qori Sandioriva', 2.95, 'Teknologi Informasi', 2025, 2.95);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Qori Sandioriva',
            name = 'Qori Sandioriva',
            grade = 2.95,
            prodi = 'Teknologi Informasi',
            angkatan = 2025,
            ipk_terakhir = 2.95,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'rizky.ramadhan@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000018', 'rizky.ramadhan@siakad.ac.id', '$2a$10$lza01yDdlT0ATVI.Mx4I6./mADIPM9RkBTlZZqp7/oBgT1vhVe2hm', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$lza01yDdlT0ATVI.Mx4I6./mADIPM9RkBTlZZqp7/oBgT1vhVe2hm', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000018';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000018', 'Rizky Ramadhan', 'Rizky Ramadhan', 3.60, 'Sistem Informasi', 2025, 3.60);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Rizky Ramadhan',
            name = 'Rizky Ramadhan',
            grade = 3.60,
            prodi = 'Sistem Informasi',
            angkatan = 2025,
            ipk_terakhir = 3.60,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'siti.nurhaliza@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000019', 'siti.nurhaliza@siakad.ac.id', '$2a$10$S6LzGKTTdayXFRNHa0b7x.6j5nBe02GspwU2aQEwPacagrDfbye0.', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$S6LzGKTTdayXFRNHa0b7x.6j5nBe02GspwU2aQEwPacagrDfbye0.', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000019';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000019', 'Siti Nurhaliza', 'Siti Nurhaliza', 3.15, 'Informatika', 2025, 3.15);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Siti Nurhaliza',
            name = 'Siti Nurhaliza',
            grade = 3.15,
            prodi = 'Informatika',
            angkatan = 2025,
            ipk_terakhir = 3.15,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;
DO $$
DECLARE
    v_user_id INTEGER;
    v_student_id INTEGER;
BEGIN
    SELECT id INTO v_user_id FROM users WHERE email = 'taufik.hidayat@siakad.ac.id';
    IF v_user_id IS NULL THEN
        INSERT INTO users (username, email, password, role, is_active)
        VALUES ('187221000020', 'taufik.hidayat@siakad.ac.id', '$2a$10$z0DIjrCBSM1X5d5EgbPXOu4/L6s26gOHXCRxkH3iVYPRayLNFn10a', 'mahasiswa', true)
        RETURNING id INTO v_user_id;
    ELSE
        UPDATE users SET password = '$2a$10$z0DIjrCBSM1X5d5EgbPXOu4/L6s26gOHXCRxkH3iVYPRayLNFn10a', role = 'mahasiswa', is_active = true WHERE id = v_user_id;
    END IF;

    SELECT id INTO v_student_id FROM students WHERE nim = '187221000020';
    IF v_student_id IS NULL THEN
        INSERT INTO students (user_id, nim, nama, name, grade, prodi, angkatan, ipk_terakhir)
        VALUES (v_user_id, '187221000020', 'Taufik Hidayat', 'Taufik Hidayat', 2.45, 'Teknologi Informasi', 2025, 2.45);
    ELSE
        UPDATE students SET
            user_id = v_user_id,
            nama = 'Taufik Hidayat',
            name = 'Taufik Hidayat',
            grade = 2.45,
            prodi = 'Teknologi Informasi',
            angkatan = 2025,
            ipk_terakhir = 2.45,
            deleted_at = NULL
        WHERE id = v_student_id;
    END IF;
END $$;

-- 4. 10 Mata Kuliah
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK001';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK001', 'Pemrograman Web', 3, 3, 30);
    ELSE
        UPDATE courses SET
            nama_mk = 'Pemrograman Web',
            sks = 3,
            semester = 3,
            kuota = 30
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK002';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK002', 'Basis Data Lanjut', 3, 3, 25);
    ELSE
        UPDATE courses SET
            nama_mk = 'Basis Data Lanjut',
            sks = 3,
            semester = 3,
            kuota = 25
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK003';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK003', 'Algoritma dan Pemrograman', 4, 1, 35);
    ELSE
        UPDATE courses SET
            nama_mk = 'Algoritma dan Pemrograman',
            sks = 4,
            semester = 1,
            kuota = 35
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK004';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK004', 'Jaringan Komputer', 3, 3, 20);
    ELSE
        UPDATE courses SET
            nama_mk = 'Jaringan Komputer',
            sks = 3,
            semester = 3,
            kuota = 20
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK005';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK005', 'Rekayasa Perangkat Lunak', 3, 5, 30);
    ELSE
        UPDATE courses SET
            nama_mk = 'Rekayasa Perangkat Lunak',
            sks = 3,
            semester = 5,
            kuota = 30
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK006';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK006', 'Kecerdasan Buatan', 3, 5, 25);
    ELSE
        UPDATE courses SET
            nama_mk = 'Kecerdasan Buatan',
            sks = 3,
            semester = 5,
            kuota = 25
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK007';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK007', 'Keamanan Informasi', 3, 5, 20);
    ELSE
        UPDATE courses SET
            nama_mk = 'Keamanan Informasi',
            sks = 3,
            semester = 5,
            kuota = 20
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK008';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK008', 'Cloud Computing', 3, 7, 20);
    ELSE
        UPDATE courses SET
            nama_mk = 'Cloud Computing',
            sks = 3,
            semester = 7,
            kuota = 20
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK009';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK009', 'Proyek Sistem Informasi', 4, 7, 15);
    ELSE
        UPDATE courses SET
            nama_mk = 'Proyek Sistem Informasi',
            sks = 4,
            semester = 7,
            kuota = 15
        WHERE id = v_course_id;
    END IF;
END $$;
DO $$
DECLARE
    v_course_id INTEGER;
BEGIN
    SELECT id INTO v_course_id FROM courses WHERE kode_mk = 'MK010';
    IF v_course_id IS NULL THEN
        INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota)
        VALUES ('MK010', 'Etika Profesi IT', 2, 1, 40);
    ELSE
        UPDATE courses SET
            nama_mk = 'Etika Profesi IT',
            sks = 2,
            semester = 1,
            kuota = 40
        WHERE id = v_course_id;
    END IF;
END $$;

COMMIT;
