package _mqtt

import (
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

type MQTTConfig struct {
	Broker    string
	ClientID  string
	Username  string
	Password  string
	CleanSess bool
	StoreDir  string // "" or ":memory:" for in-memory
}

type MQTTClient struct {
	client mqtt.Client
	mu     sync.RWMutex
}

var (
	instance *MQTTClient
	once     sync.Once
)

// GetMQTTClient returns the singleton MQTT client instance.
// It initializes the connection only once and reuses it for all subsequent calls.
func getMQTTClient(cfg MQTTConfig) *MQTTClient {
	once.Do(func() {
		opts := mqtt.NewClientOptions()
		opts.AddBroker(cfg.Broker)
		opts.SetClientID(cfg.ClientID)
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
		opts.SetCleanSession(cfg.CleanSess)
		opts.KeepAlive = 60
		opts.AutoReconnect = true
		opts.SetConnectionNotificationHandler(func(client mqtt.Client, notification mqtt.ConnectionNotification) {
			switch n := notification.(type) {
			case mqtt.ConnectionNotificationConnected:
				log.Info("[NOTIFICATION] connected")
			case mqtt.ConnectionNotificationConnecting:
				log.Infof("[NOTIFICATION] connecting (isReconnect=%t) [%d]", n.IsReconnect, n.Attempt)
			case mqtt.ConnectionNotificationFailed:
				log.Warnf("[NOTIFICATION] connection failed: %v", n.Reason)
			case mqtt.ConnectionNotificationLost:
				log.Errorf("[NOTIFICATION] connection lost: %v", n.Reason)
			case mqtt.ConnectionNotificationBroker:
				log.Infof("[NOTIFICATION] broker connection: %s", n.Broker.String())
			case mqtt.ConnectionNotificationBrokerFailed:
				log.Errorf("[NOTIFICATION] broker connection failed: %v [%s]", n.Reason, n.Broker.String())
			}
		})

		// Set up store if specified
		if cfg.StoreDir != "" && cfg.StoreDir != ":memory:" {
			opts.SetStore(mqtt.NewFileStore(cfg.StoreDir))
		}

		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Fatalf("MQTT connect error: %v", token.Error())
		}

		instance = &MQTTClient{client: client}
		log.Printf("MQTT client initialized: broker=%s, clientID=%s", cfg.Broker, cfg.ClientID)
	})
	return instance
}

// Publish sends a message to a topic.
// qos: Quality of Service (0, 1, or 2)
// retained: Whether the broker should retain this message for future subscribers
func (m *MQTTClient) Publish(topic string, qos byte, retained bool, payload any) error {
	if m.client == nil || !m.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	token := m.client.Publish(topic, qos, retained, payload)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Failed to publish to topic %s: %v", topic, token.Error())
		return token.Error()
	}
	return nil
}

// Subscribe subscribes to a topic with a message handler.
//
// Parameters:
//
//	topic: MQTT topic to subscribe to (supports wildcards # and +)
//	qos: Quality of Service (0, 1, or 2)
//	handler: function to handle incoming messages
func (m *MQTTClient) Subscribe(topic string, handler mqtt.MessageHandler, qos ...byte) error {
	var qosLevel byte = 2 // default QoS
	if len(qos) > 0 {
		qosLevel = qos[0]
	}

	if m.client == nil || !m.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	token := m.client.Subscribe(topic, qosLevel, handler)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Failed to subscribe to topic %s: %v", topic, token.Error())
		return token.Error()
	}

	log.Printf("Subscribed to topic: %s (QoS %d)", topic, qosLevel)
	return nil
}

// Unsubscribe removes subscription from one or more topics.
func (m *MQTTClient) Unsubscribe(topics ...string) error {
	if m.client == nil || !m.client.IsConnected() {
		return fmt.Errorf("MQTT client is not connected")
	}

	token := m.client.Unsubscribe(topics...)
	token.Wait()

	if token.Error() != nil {
		log.Printf("Failed to unsubscribe from topics: %v", token.Error())
		return token.Error()
	}

	log.Printf("Unsubscribed from topics: %v", topics)
	return nil
}

// IsConnected returns the connection status.
func (m *MQTTClient) IsConnected() bool {
	if m.client == nil {
		return false
	}
	return m.client.IsConnected()
}

// Disconnect cleanly disconnects the client.
// quiesce: milliseconds to wait for pending messages to complete
func (m *MQTTClient) Disconnect(quiesce uint) {
	if m.client != nil {
		m.client.Disconnect(quiesce)
		log.Println("MQTT client disconnected")
	}
}
