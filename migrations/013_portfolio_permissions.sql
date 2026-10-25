-- Migration 013: permission portofolio (Modul 10)
INSERT INTO permissions (name, description) VALUES
    ('portfolio:manage', 'Mengelola portofolio mahasiswa mana pun')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role_name, permission_name) VALUES
    ('admin', 'portfolio:manage')
ON CONFLICT DO NOTHING;
