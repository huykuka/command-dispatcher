package _queue

import (
	"sync"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

var (
	queueServer     *asynq.Server
	queueServerOnce sync.Once
)

func GetQueueServer() *asynq.Server {
	if queueServer == nil {
		log.Fatalf("queue server not initialized")
	}
	return queueServer
}

func CloseQueueServer() {
	if queueServer == nil {
		log.Fatalf("queue server not initialized")
	}
	queueServer.Shutdown()
}

func InitQueueServer(redisOption asynq.RedisConnOpt, opts asynq.Config) {
	// Initialize the queue server
	queueServerOnce.Do(func() {
		srv := asynq.NewServer(
			redisOption,
			opts,
		)
		queueServer = srv
	})
}
