package worker

import (
	"command-dispatcher/internal/config/_queue"
	"command-dispatcher/internal/models"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

// TypeCommandExecutionJob Task type identifiers for the command domain.
const (
	TypeCommandExecutionJob = "command:execute"
)

// Init starts the asynq server and registers all domain worker handlers.
func Init() {
	srv := _queue.GetQueueServer()
	mux := asynq.NewServeMux()

	commandWorker := NewCommandWorker(TypeCommandExecutionJob)
	mux.HandleFunc(TypeCommandExecutionJob, commandWorker.Process)

	log.Info("Worker server starting...")
	if err := srv.Run(mux); err != nil {
		log.Fatalf("Could not run worker server: %v", err)
	}
}

// EnqueueTask enqueues a pre-built task.
func EnqueueTask(task *asynq.Task) error {
	_, err := _queue.GetQueueClient().Enqueue(task)
	if err != nil {
		log.Errorf("Could not enqueue task: %v", err)
		return err
	}
	return nil
}

// EnqueueCommandExecutionTask is a convenience wrapper: generate then enqueue.
func EnqueueCommandExecutionTask(dto models.CommandCreateDTO) error {
	cw := NewCommandWorker(TypeCommandExecutionJob)
	t, err := cw.Generate(dto)
	if err != nil {
		return err
	}
	if err := EnqueueTask(t); err != nil {
		return err
	}
	return nil
}
