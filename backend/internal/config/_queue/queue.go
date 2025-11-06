package _queue

import "github.com/hibiken/asynq"

func Init() {
	redisAddress := "redis:6379"

	// Initialize Queue Client and Server
	InitQueueClient(asynq.RedisClientOpt{Addr: redisAddress})

	InitQueueServer(asynq.RedisClientOpt{Addr: redisAddress},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 2,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		})
}
