-- Migration 009: Penyesuaian idempotency key notifikasi untuk mendukung broadcast 1-to-N
-- (Event course.quota_changed menghasilkan notifikasi ke seluruh mahasiswa peserta kelas)

ALTER TABLE notifications DROP CONSTRAINT IF EXISTS notifications_outbox_unique;

ALTER TABLE notifications ADD CONSTRAINT notifications_outbox_user_unique UNIQUE (outbox_id, user_id);

