package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"api-students/app/model"
	"api-students/config"
	"api-students/database"
)

func main() {
	outboxID := flag.Int64("outbox", 0, "ID outbox event yang akan dikirim ulang")
	copies := flag.Int("copies", 1, "Jumlah salinan task yang akan di-enqueue ke Redis")
	flag.Parse()

	if *outboxID <= 0 {
		log.Fatal("flag -outbox harus berupa angka positif")
	}
	if *copies <= 0 {
		log.Fatal("flag -copies harus minimal 1")
	}

	config.LoadEnv()
	cfg := config.Load()

	pool, err := database.NewDBPool(cfg)
	if err != nil {
		log.Fatalf("koneksi database gagal: %v", err)
	}
	defer pool.Close()

	var event model.OutboxEvent
	err = pool.QueryRow(context.Background(), `
		SELECT id, event_type, payload, created_at
		FROM outbox
		WHERE id = $1
	`, *outboxID).Scan(&event.ID, &event.EventType, &event.Payload, &event.CreatedAt)
	if err != nil {
		log.Fatalf("mengambil outbox %d gagal: %v", *outboxID, err)
	}

	job := model.NotificationJob{
		OutboxID:  event.ID,
		EventType: event.EventType,
		Payload:   event.Payload,
	}
	jobBytes, err := json.Marshal(job)
	if err != nil {
		log.Fatalf("marshal job gagal: %v", err)
	}

	client := asynq.NewClient(config.RedisOpt())
	defer client.Close()

	fmt.Printf("Mengirim %d salinan task untuk outbox ID %d ke Redis...\n", *copies, *outboxID)
	for i := 1; i <= *copies; i++ {
		// TaskID unik per salinan agar tidak ditolak ErrTaskIDConflict oleh Redis broker
		task := asynq.NewTask(
			model.TaskSendNotification,
			jobBytes,
			asynq.Queue(model.QueueNotifications),
			asynq.TaskID(fmt.Sprintf("redeliver-%d-%d-%d", *outboxID, i, time.Now().UnixNano())),
			asynq.MaxRetry(5),
			asynq.Timeout(30*time.Second),
			asynq.Retention(24*time.Hour),
		)

		info, err := client.Enqueue(task)
		if err != nil {
			log.Fatalf("gagal enqueue salinan %d: %v", i, err)
		}
		fmt.Printf("Salinan %d terkirim: ID=%s Queue=%s\n", i, info.ID, info.Queue)
	}

	fmt.Printf("Selesai: %d salinan berhasil dimasukkan ke queue '%s'.\n", *copies, model.QueueNotifications)
}
