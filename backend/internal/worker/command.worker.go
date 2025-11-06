package worker

import (
	"command-dispatcher/internal/config/_mqtt"
	"command-dispatcher/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

var mqttClient = _mqtt.GetClient()

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
	return asynq.NewTask(cw.jobName, b, asynq.MaxRetry(-1), asynq.Timeout(30*time.Second)), nil
}

// Process executes the queued command.
func (*CommandWorker) Process(ctx context.Context, t *asynq.Task) error {
	var p models.CommandCreateDTO

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Publish the original task payload to MQTT (device-specific topic)
	if !_mqtt.IsInitialized() {
		return fmt.Errorf("mqtt client not initialized")
	}
	taskId := t.ResultWriter().TaskID()

	publishCommand(p.DeviceID, t.Payload())

	if err := waitForAcknowledgement(ctx, p.DeviceID); err != nil {
		return err
	}

	if err := waitForCompletion(ctx, p.DeviceID); err != nil {
		return err
	}

	log.Infof("Command completed by device %s %s", p.DeviceID, taskId)

	return nil
}

// publishCommand publishes the command payload to the specified dispatch topic.
func publishCommand(deviceID string, payload []byte) {
	disPatchTopic := "device/" + deviceID + "/dispatch"
	mqttClient.Publish(disPatchTopic, 2, false, payload)
}

// waitForAcknowledgement waits for an acknowledgment from the device or times out.
func waitForAcknowledgement(ctx context.Context, deviceID string) error {
	acknowledgeTopic := "device/" + deviceID + "/acknowledge"
	ackCh := make(chan struct{})

	mqttClient.Subscribe(acknowledgeTopic, func(client mqtt.Client, msg mqtt.Message) {
		ackCh <- struct{}{}
	}, 2)
	defer mqttClient.Unsubscribe(acknowledgeTopic) // Ensure unsubscribe happens

	select {
	case <-ackCh:
		log.Infof("Command acknowledged by device %s", deviceID)
		return nil
	case <-time.After(20 * time.Second):
		return fmt.Errorf("command acknowledgment timed out by device %s", deviceID)
	}
}

// waitForCompletion waits for command completion from the device or times out.
func waitForCompletion(ctx context.Context, deviceID string) error {
	completeTopic := "device/" + deviceID + "/complete"
	completeCh := make(chan struct{})

	mqttClient.Subscribe(completeTopic, func(client mqtt.Client, msg mqtt.Message) {
		completeCh <- struct{}{}
	}, 2)
	defer mqttClient.Unsubscribe(completeTopic) // Ensure unsubscribe happens

	select {
	case <-completeCh:
		log.Infof("Command completed by device %s", deviceID)
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("command completion timed out by device %s", deviceID)
	}
}
