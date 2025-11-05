package worker

import (
	"command-dispatcher/internal/config/_mqtt"
	"command-dispatcher/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

// TaskWorker enforces all workers to implement Generate and Process plus expose a JobName.
type TaskWorker interface {
	Generate(models.CommandCreateDTO) (*asynq.Task, error)
	Process(context.Context, *asynq.Task) error
	JobName() string
}

type CommandWorker struct {
	jobName string
}

func NewCommandWorker(jobName string) *CommandWorker {
	return &CommandWorker{jobName: jobName}
}

func (cw *CommandWorker) JobName() string { return cw.jobName }

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

	log.Infof("Processing command for deviceId=%s type=%s", p.DeviceID, p.Type)

	// Publish the original task payload to MQTT (device-specific topic)
	if !_mqtt.IsInitialized() {
		return fmt.Errorf("mqtt client not initialized")
	}

	topic := "commands/" + p.DeviceID + "/dispatch"

	if err := _mqtt.GetClient().Publish(topic, 2, false, t.Payload()); err != nil {
		return fmt.Errorf("mqtt publish failed: %w", err)
	}

	log.Infof("Published command to MQTT topic=%s deviceId=%s type=%s", topic, p.DeviceID, p.Type)
	return nil
}
