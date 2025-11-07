package device

import (
	"command-dispatcher/internal/config/_mqtt"

	log "github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Register() {
	log.Println("device subscriber registered")
	mqttClient := _mqtt.GetClient()

	// Use MQTT single-level wildcard `+`
	mqttClient.Subscribe("devices/+/status", func(c mqtt.Client, m mqtt.Message) {
		// TODO: dispatch to service, update DB, etc.
	})
}
