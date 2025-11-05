package _queue

import (
	"sync"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

var (
	queueClient     *asynq.Client
	queueClientOnce sync.Once
)

func InitQueueClient(cfg asynq.RedisClientOpt) {
	queueClientOnce.Do(func() {
		queueClient = asynq.NewClient(cfg)
		log.Infof("Asynq Redis queue client initialized at: %s", cfg.Addr)
	})
}

func GetQueueClient() *asynq.Client {
	if queueClient == nil {
		log.Fatalf("Queue client has not been initialized. Call InitQueueClient first.")
	}
	return queueClient
}

func CloseQueueClient() {
	if queueClient != nil {
		if err := queueClient.Close(); err != nil {
			log.Errorf("Failed to close queue client: %v", err)
		} else {
			log.Info("Queue client disconnected.")
		}
	}
}
