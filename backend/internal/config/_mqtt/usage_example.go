package _mqtt

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// ExampleUsageInMultiplePlaces demonstrates how to use the MQTT client
// in different parts of your application after initialization
func ExampleUsageInMultiplePlaces() {
	// Note: mqtt.InitGlobal(cfg) must be called first in main.go or config.Init()

	// Example 1: Use in a handler
	handleDeviceMessage("device-123", "temperature: 25.5")

	// Example 2: Use in a service
	sendAlert("High temperature detected!")

	// Example 3: Subscribe in a subscriber
	startDeviceMonitoring()
}

// handleDeviceMessage simulates a handler that publishes device data
func handleDeviceMessage(deviceID string, data string) {
	// Get the global MQTT client (no need to pass it around!)
	client := GetClient()

	topic := fmt.Sprintf("device/%s/data", deviceID)
	err := client.Publish(topic, 0, false, data)
	if err != nil {
		log.Printf("Error publishing device data: %v", err)
		return
	}

	fmt.Printf("✅ Published device data to %s\n", topic)
}

// sendAlert simulates a service that publishes alerts
func sendAlert(message string) {
	// Get the global MQTT client
	client := GetClient()

	err := client.Publish("alerts/critical", 1, false, message)
	if err != nil {
		log.Printf("Error publishing alert: %v", err)
		return
	}

	fmt.Printf("🚨 Alert sent: %s\n", message)
}

// startDeviceMonitoring simulates a subscriber that listens to device topics
func startDeviceMonitoring() {
	// Get the global MQTT client
	client := GetClient()

	// Subscribe to all device topics
	err := client.Subscribe("device/#", func(c mqtt.Client, msg mqtt.Message) {
		fmt.Printf("📨 [Monitor] Received on %s: %s\n", msg.Topic(), string(msg.Payload()))
	})

	if err != nil {
		log.Printf("Error subscribing to device topics: %v", err)
		return
	}

	fmt.Println("📡 Device monitoring started")
}

// Example: Device Handler Struct
type DeviceHandler struct {
	// No need to store MQTT client as a field!
}

func NewDeviceHandler() *DeviceHandler {
	return &DeviceHandler{}
}

func (h *DeviceHandler) PublishStatus(deviceID string, status string) error {
	client := GetClient()
	topic := fmt.Sprintf("device/%s/status", deviceID)
	return client.Publish(topic, 1, true, status) // Retained message
}

func (h *DeviceHandler) PublishCommand(deviceID string, command string) error {
	client := GetClient()
	topic := fmt.Sprintf("device/%s/command", deviceID)
	return client.Publish(topic, 1, false, command)
}

// Example: Alert Service Struct
type AlertService struct{}

func NewAlertService() *AlertService {
	svc := &AlertService{}
	svc.startListening()
	return svc
}

func (s *AlertService) startListening() {
	client := GetClient()

	// Subscribe to alerts from all devices
	client.Subscribe("device/+/alerts", func(c mqtt.Client, msg mqtt.Message) {
		// deviceID := extractDeviceID(msg.Topic())
		// s.handleAlert(deviceID, string(msg.Payload()))
	})
}

func (s *AlertService) handleAlert(deviceID string, alert string) {
	fmt.Printf("⚠️  [AlertService] Device %s: %s\n", deviceID, alert)
	// Process alert logic here
}

func (s *AlertService) BroadcastAlert(message string) error {
	client := GetClient()
	return client.Publish("alerts/broadcast", 2, false, message)
}

// Example: Metrics Collector
type MetricsCollector struct{}

func NewMetricsCollector() *MetricsCollector {
	collector := &MetricsCollector{}
	collector.subscribeToMetrics()
	return collector
}

func (m *MetricsCollector) subscribeToMetrics() {
	client := GetClient()

	// Subscribe to metrics from all devices
	client.Subscribe("device/+/metrics", func(c mqtt.Client, msg mqtt.Message) {
		// deviceID := extractDeviceID(msg.Topic())
		// m.processMetrics(deviceID, string(msg.Payload()))
	})

	fmt.Println("📊 Metrics collector started")
}

func (m *MetricsCollector) processMetrics(deviceID string, metrics string) {
	fmt.Printf("📊 [Metrics] Device %s: %s\n", deviceID, metrics)
	// Store metrics in database, etc.
}

func (m *MetricsCollector) PublishSystemMetrics(metrics string) error {
	client := GetClient()
	return client.Publish("system/metrics", 0, false, metrics)
}

// ExampleRealWorldUsage shows a more complete example
func ExampleRealWorldUsage() {
	// Ensure MQTT is initialized (done in config.Init())
	if !IsInitialized() {
		log.Fatal("MQTT not initialized! Call mqtt.InitGlobal() first")
	}

	// Create services
	deviceHandler := NewDeviceHandler()
	alertService := NewAlertService()
	metricsCollector := NewMetricsCollector()

	// Use the services
	deviceHandler.PublishStatus("device-001", "online")
	deviceHandler.PublishCommand("device-001", "START_PROCESS")

	alertService.BroadcastAlert("System maintenance in 5 minutes")

	metricsCollector.PublishSystemMetrics(`{"cpu": 45, "memory": 60}`)

	fmt.Println("\n✅ All services are using the same MQTT client instance!")
	fmt.Printf("📡 MQTT Config: %+v\n", GetConfig())
}
