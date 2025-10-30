package device

import (
	"command-dispatcher/internal/config/_mqtt"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func Register() {
	fmt.Print("hello world")
	mqttClient := _mqtt.GetClient()

	mqttClient.Subscribe("devices/+/commands", func(client mqtt.Client, msg mqtt.Message) {
		logrus.Infof("Received command on topic: %s, payload: %s", msg.Topic(), string(msg.Payload()))
		// deviceID := extractDeviceID(msg.Topic())
		// handleCommand(deviceID, string(msg.Payload()))
	}, 2)
}
