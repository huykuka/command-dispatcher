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

type CommandWorker struct {
	jobName string
}

// Define a struct that includes the original DTO and the TaskID
type commandPayload struct {
	models.CommandCreateDTO
	TaskID string `json:"taskId"`
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

	log.Infof("Executing command for device %s, task %s", p.DeviceID, taskId)

	publishCommand(p.DeviceID, taskId, t.Payload())

	if err := waitForAcknowledgement(ctx, p.DeviceID, taskId); err != nil {
		return err
	}

	if err := waitForCompletion(ctx, p.DeviceID, taskId); err != nil {
		return err
	}

	log.Infof("Finished processing command for device %s, task %s", p.DeviceID, taskId)
	return nil
}

// publishCommand publishes the command payload to the specified dispatch topic.
func publishCommand(deviceID, taskId string, originalPayload []byte) {
	dispatchTopic := fmt.Sprintf("device/%s/dispatch", deviceID)

	var originalPayloadDTO models.CommandCreateDTO
	if err := json.Unmarshal(originalPayload, &originalPayloadDTO); err != nil {
		log.Errorf("Failed to unmarshal original command payload for device %s, task %s: %v", deviceID, taskId, err)
		return // Or handle error appropriately
	}

	payloadWithTaskID := commandPayload{
		CommandCreateDTO: originalPayloadDTO,
		TaskID:           taskId,
	}

	finalPayload, err := json.Marshal(payloadWithTaskID)
	if err != nil {
		log.Errorf("Failed to marshal final command payload for device %s, task %s: %v", deviceID, taskId, err)
		return // Or handle error appropriately
	}

	_mqtt.GetClient().Publish(dispatchTopic, 2, false, finalPayload)
}

// waitForAcknowledgement waits for an acknowledgment from the device or times out.
func waitForAcknowledgement(ctx context.Context, deviceID, taskId string) error {
	acknowledgeTopic := fmt.Sprintf("device/%s/acknowledge/%s", deviceID, taskId)
	ackCh := make(chan struct{})

	_mqtt.GetClient().Subscribe(acknowledgeTopic, func(client mqtt.Client, msg mqtt.Message) {
		ackCh <- struct{}{}
	}, 2)
	defer _mqtt.GetClient().Unsubscribe(acknowledgeTopic) // Ensure unsubscribe happens

	select {
	case <-ackCh:
		log.Infof("Command aknowledged for device %s, task %s", deviceID, taskId)
		return nil
	case <-time.After(20 * time.Second):
		msg := fmt.Sprintf("Command aknowledgment timed out by device %s, task %s", deviceID, taskId)
		log.Error(msg)
		return errors.New(msg)
	}
}

// waitForCompletion waits for command completion from the device or times out.
func waitForCompletion(ctx context.Context, deviceID, taskId string) error {
	completeTopic := fmt.Sprintf("device/%s/complete/%s", deviceID, taskId)
	completeCh := make(chan struct{})

	_mqtt.GetClient().Subscribe(completeTopic, func(client mqtt.Client, msg mqtt.Message) {
		completeCh <- struct{}{}
	}, 2)
	defer _mqtt.GetClient().Unsubscribe(completeTopic) // Ensure unsubscribe happens

	select {
	case <-completeCh:
		log.Infof("Command completed for device %s, task %s", deviceID, taskId)
		return nil
	case <-time.After(5 * time.Second): // TODO: make timeout configurable
		msg := fmt.Sprintf("Command completion timed out by device %s, task %s", deviceID, taskId)
		log.Error(msg)
		return errors.New(msg)
	}
}
