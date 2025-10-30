package _mqtt

import (
	"log"
	"sync"
)

var (
	globalConfig MQTTConfig
	configOnce   sync.Once
	initDone     bool
)

// InitGlobal initializes the global MQTT configuration once and connects the client.
// This should be called once at application startup.
func Init(cfg MQTTConfig) *MQTTClient {
	configOnce.Do(func() {
		globalConfig = cfg
		initDone = true
		log.Printf("MQTT global config initialized: broker=%s, clientID=%s", cfg.Broker, cfg.ClientID)
	})
	return GetMQTTClient(cfg)
}

// GetClient returns the singleton MQTT client using the global configuration.
// InitGlobal must be called before using this method.
func GetClient() *MQTTClient {
	if !initDone {
		log.Fatal("MQTT client not initialized. Call mqtt.InitGlobal() first in your main.go")
	}
	return GetMQTTClient(globalConfig)
}

// IsInitialized returns true if InitGlobal has been called.
func IsInitialized() bool {
	return initDone
}

// GetConfig returns a copy of the global configuration.
// Returns an empty config if not initialized.
func GetConfig() MQTTConfig {
	if !initDone {
		return MQTTConfig{}
	}
	return globalConfig
}
