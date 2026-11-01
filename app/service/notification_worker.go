package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

// NotificationWorker mengerjakan task "notification:send".
//
// Bagi worker, Asynq adalah framework-nya — sama seperti Fiber bagi API.
// Karena itu method di sini menerima *asynq.Task, seperti method service
// lain menerima *fiber.Ctx.
type NotificationWorker struct {
	users         *repository.UserRepository
	notifications repository.NotificationRepository
	tx            *repository.TxManager
	mailer        helper.Mailer
	logger        *slog.Logger
	SimulateBug9  bool
}

func NewNotificationWorker(
	users *repository.UserRepository,
	notifications repository.NotificationRepository,
	tx *repository.TxManager,
	mailer helper.Mailer,
	logger *slog.Logger,
) *NotificationWorker {
	return &NotificationWorker{users: users, notifications: notifications,
		tx: tx, mailer: mailer, logger: logger}
}

// HandleSend mengembalikan error untuk meminta Asynq melakukan retry nanti,
// atau error yang membungkus asynq.SkipRetry untuk menyatakan "jangan di-retry,
// task ini tidak mungkin berhasil".
//
// FIX BUG #3: Parameter menerima *asynq.Task (pointer), bukan asynq.Task (value).
func (w *NotificationWorker) HandleSend(ctx context.Context, t *asynq.Task) error {
	var job model.NotificationJob
	if err := json.Unmarshal(t.Payload(), &job); err != nil {
		return fmt.Errorf("%w: payload task rusak: %v", asynq.SkipRetry, err)
	}

	var event model.EnrollmentEvent
	if err := json.Unmarshal(job.Payload, &event); err != nil {
		return fmt.Errorf("%w: payload event %d rusak: %v", asynq.SkipRetry, job.OutboxID, err)
	}

	// Penanganan khusus Tugas Mandiri D.2: Event perubahan kuota untuk seluruh peserta
	if job.EventType == model.EventCourseQuotaChanged {
		return w.handleQuotaChanged(ctx, job, event)
	}

	user, err := w.users.FindByID(ctx, event.UserID)
	if errors.Is(err, repository.ErrNotFound) {
		// Akunnya sudah dihapus. Melakukan retry tidak akan memunculkannya.
		return fmt.Errorf("%w: user %d tidak ditemukan", asynq.SkipRetry, event.UserID)
	}
	if err != nil {
		return err // database sedang bermasalah: layak di-retry
	}

	msg, err := BuildMessage(job.EventType, event, *user)
	if err != nil {
		return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
	}

	alreadySent := false
	err = w.tx.WithTxOnce(ctx, func(tx repository.DBTX) error {
		// Claim dan pengiriman berada di dalam SATU transaction:
		// - pengiriman gagal -> claim ikut dibatalkan -> attempt
		//   berikutnya masih dapat mengirim;
		// - pengiriman sukses -> claim di-commit -> attempt berikutnya
		//   melihat claim ini dan tidak mengirim lagi.
		//
		// FIX BUG #8: Panggil WithTx(tx).Claim agar claim berada di dalam transaksi yang sama
		claimed, err := w.notifications.WithTx(tx).Claim(ctx, model.Notification{
			OutboxID: job.OutboxID, UserID: user.ID, Channel: "email",
			Subject: msg.Subject, Body: msg.Body,
		})
		if err != nil {
			return err
		}
		if !claimed {
			alreadySent = true
			return nil
		}
		return w.mailer.Send(ctx, msg)
	})

	if errors.Is(err, helper.ErrPermanentDelivery) {
		if w.SimulateBug9 {
			// Bug #9 sebelum perbaikan: %v membungkus SkipRetry sehingga errors.Is gagal mendeteksinya
			return fmt.Errorf("%v: %w", asynq.SkipRetry, err)
		}
		// FIX BUG #9: Bungkus SkipRetry menggunakan %w agar errors.Is mengenalinya
		return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
	}
	if err != nil {
		return err
	}

	retry, _ := asynq.GetRetryCount(ctx)
	w.logger.Info("notification_processed",
		slog.Int64("outbox_id", job.OutboxID),
		slog.String("event", job.EventType),
		slog.Int("retry", retry),
		slog.Bool("already_sent", alreadySent))
	return nil
}

// handleQuotaChanged mengirimkan notifikasi perubahan kuota ke seluruh peserta kelas.
// Setiap pengiriman peserta dijalankan dalam WithTxOnce terpisah dengan claim idempotency (outbox_id, user_id).
// Bila email ke-7 gagal, retry tidak mengirim ulang email ke-1 sampai ke-6.
func (w *NotificationWorker) handleQuotaChanged(ctx context.Context, job model.NotificationJob, event model.EnrollmentEvent) error {
	users, err := w.users.FindUsersByCourseID(ctx, event.CourseID)
	if err != nil {
		return err
	}

	retry, _ := asynq.GetRetryCount(ctx)
	for _, user := range users {
		msg, err := BuildMessage(job.EventType, event, user)
		if err != nil {
			return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
		}

		alreadySent := false
		err = w.tx.WithTxOnce(ctx, func(tx repository.DBTX) error {
			claimed, err := w.notifications.WithTx(tx).Claim(ctx, model.Notification{
				OutboxID: job.OutboxID, UserID: user.ID, Channel: "email",
				Subject: msg.Subject, Body: msg.Body,
			})
			if err != nil {
				return err
			}
			if !claimed {
				alreadySent = true
				return nil
			}
			return w.mailer.Send(ctx, msg)
		})

		if errors.Is(err, helper.ErrPermanentDelivery) {
			if w.SimulateBug9 {
				return fmt.Errorf("%v: %w", asynq.SkipRetry, err)
			}
			return fmt.Errorf("%w: %v", asynq.SkipRetry, err)
		}
		if err != nil {
			return err // Temporary failure: Asynq akan me-retry task ini
		}

		w.logger.Info("notification_processed",
			slog.Int64("outbox_id", job.OutboxID),
			slog.String("event", job.EventType),
			slog.Int("retry", retry),
			slog.Bool("already_sent", alreadySent))
	}

	return nil
}
