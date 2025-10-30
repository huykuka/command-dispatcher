package config

import (
	"command-dispatcher/internal/config/_mqtt"
	"command-dispatcher/internal/config/environments"
	"command-dispatcher/internal/config/log"
	"crypto/rand"
)

var mqttCfg = _mqtt.MQTTConfig{
	Broker:    "tcp://host.docker.internal:1883",
	ClientID:  rand.Text(),
	Username:  "",
	Password:  "",
	CleanSess: true,
	StoreDir:  ":memory:",
}

func Init() {
	log.Init()
	environments.Init()
	_mqtt.Init(mqttCfg)
}
