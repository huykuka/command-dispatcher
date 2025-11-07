package _queue

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)
type noopLogger struct{}
func (noopLogger) Printf(ctx context.Context, format string, v ...any) {}

func Init() {
	redisAddress := "redis:6379"
	redis.SetLogger(noopLogger{})	// Initialize Queue Client and Server
	InitQueueClient(asynq.RedisClientOpt{Addr: redisAddress})

	InitQueueServer(asynq.RedisClientOpt{Addr: redisAddress},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		})
}
