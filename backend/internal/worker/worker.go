package worker

import (
	"command-dispatcher/internal/config/_queue"
	"command-dispatcher/internal/models"
	"context"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

// TypeCommandExecutionJob Task type identifiers for the command domain.
const (
	TypeCommandExecutionJob = "command:execute"
)

type TaskWorker interface {
	Generate(models.CommandCreateDTO) (*asynq.Task, error)
	Process(context.Context, *asynq.Task) error
	JobName() string
}

// commandWorker implements TaskWorker (compile-time assertion)
var commandWorker TaskWorker = NewCommandWorker(TypeCommandExecutionJob)

// Init starts the asynq server and registers all domain worker handlers.
func Init() {
	srv := _queue.GetQueueServer()
	mux := asynq.NewServeMux()

	mux.HandleFunc(commandWorker.JobName(), commandWorker.Process)

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

// EnqueueCommandExecutionTask generates and enqueues a command execution task using the singleton worker.
func EnqueueCommandExecutionTask(dto models.CommandCreateDTO) error {
	t, err := commandWorker.Generate(dto)
	if err != nil {
		return err
	}
	return EnqueueTask(t)
}
