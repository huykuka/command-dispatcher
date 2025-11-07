package config

import (
	"command-dispatcher/internal/config/_mqtt"
	"command-dispatcher/internal/config/_queue"
	"command-dispatcher/internal/config/db"
	"command-dispatcher/internal/config/log"
	"crypto/rand"
	"fmt"
)

var mqttCfg = _mqtt.MQTTConfig{
	Broker:    "tcp://mqtt:1883",
	ClientID:  fmt.Sprintf(rand.Text(), "-backend"),
	Username:  "",
	Password:  "",
	CleanSess: true,
	StoreDir:  ":memory:",
}

func Init() {
	log.Init()
	db.Init()
	//environments.Init()?
	_mqtt.Init(mqttCfg)
	_queue.Init()
}
