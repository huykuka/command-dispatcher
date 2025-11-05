// Deprecated: logic split across task_types.go, payload.go, generate.go, process.go, register.go
package worker

import (
	"command-dispatcher/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

type CommandWorker struct {
	jobName string
}

func NewCommandWorker(jobName string) *CommandWorker {
	return &CommandWorker{
		jobName: jobName,
	}
}

// Generate builds an asynq.Task for the worker's jobName using the provided DTO payload.
func (cw *CommandWorker) Generate(dto models.CommandCreateDTO) (*asynq.Task, error) {
	if cw.jobName == "" {
		return nil, errors.New("jobName is empty")
	}
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("marshal command execution payload: %w", err)
	}
	log.Debugf("Generate command execution task type=%s deviceId=%s cmdType=%s", cw.jobName, dto.DeviceID, dto.Type)
	return asynq.NewTask(cw.jobName, b, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute)), nil
}

// Process executes the queued command.
func (*CommandWorker) Process(ctx context.Context, t *asynq.Task) error {
	var p models.CommandCreateDTO
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Infof("Processing command task deviceId=%s type=%s", p.DeviceID, p.Type)
	// TODO: Implement domain-specific execution logic here.
	return nil
}
