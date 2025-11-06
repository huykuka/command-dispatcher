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
	log.Println("Executing command task...", string(t.Payload()))
	var p models.CommandCreateDTO

	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Publish the original task payload to MQTT (device-specific topic)
	if !_mqtt.IsInitialized() {
		return fmt.Errorf("mqtt client not initialized")
	}

	mqttClient := _mqtt.GetClient()

	disPatchTopic := "commands/" + p.DeviceID + "/dispatch"
	completeTopic := "commands/" + p.DeviceID + "/complete"
	acknowledgeTopic := "commands/" + p.DeviceID + "/acknowledge"

	ackCh := make(chan struct{})
	completeCh := make(chan struct{})

	mqttClient.Subscribe(acknowledgeTopic, func(client mqtt.Client, msg mqtt.Message) {
		completeCh <- struct{}{}
	}, 2)
	defer mqttClient.Unsubscribe(acknowledgeTopic)

	if err := mqttClient.Publish(disPatchTopic, 2, false, t.Payload()); err != nil {
		return fmt.Errorf("mqtt publish failed: %w", err)
	}

	select {
	case <-ackCh:
		// Acknowledgment received
		log.Infof("Command acknowledged by device %s", p.DeviceID)
	case <-time.After(5 * time.Second):
		log.Infof("Command acknowledgment timed out by device %s", p.DeviceID)
		return nil
	}

	mqttClient.Subscribe(completeTopic, func(client mqtt.Client, msg mqtt.Message) {
		//fmt.Println(string(msg.Payload()))
		ackCh <- struct{}{}
	}, 2)
	defer mqttClient.Unsubscribe(completeTopic)

	//select {
	//case <-completeCh:
	//	// Acknowledgment received
	//	log.Infof("Command completed by device %s", p.DeviceID)
	//case <-time.After(5 * time.Second):
	//	log.Infof("Command timed out by device %s", p.DeviceID)
	//}

	return nil
}
