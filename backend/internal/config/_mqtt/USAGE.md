# MQTT Client - Usage Guide

A simplified, singleton-based MQTT client for Go with support for wildcard topics.

## Table of Contents
- [Quick Start](#quick-start)
- [Using in Multiple Places](#using-in-multiple-places)
- [API Reference](#api-reference)
- [Wildcard Patterns](#wildcard-patterns)
- [Examples](#examples)

## Quick Start

```go
package main

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    mqttconfig "your-project/internal/config/mqtt"
)

func main() {
    // 1. Configure the client
    cfg := mqttconfig.MQTTConfig{
        Broker:    "tcp://localhost:1883",
        ClientID:  "my-app",
        Username:  "",
        Password:  "",
        CleanSess: true,
        StoreDir:  ":memory:",
    }

    // 2. Get singleton instance
    client := mqttconfig.GetMQTTClient(cfg)
    defer client.Disconnect(250)

    // 3. Subscribe to a topic
    client.Subscribe("device/sensors", 0, func(c mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Received: %s\n", string(msg.Payload()))
    })

    // 4. Publish a message
    client.Publish("device/sensors", 0, false, "Temperature: 25°C")
}
```

---

## Using in Multiple Places

The singleton pattern allows you to use the MQTT client anywhere in your application without passing it around.

### Initialize Once in `main.go` or `config.Init()`

```go
// In your main.go or config package
package main

import (
    mqttconfig "your-project/internal/config/mqtt"
)

func main() {
    // Initialize MQTT client once at startup
    cfg := mqttconfig.MQTTConfig{
        Broker:    "tcp://localhost:1883",
        ClientID:  "command-dispatcher",
        Username:  "",
        Password:  "",
        CleanSess: true,
        StoreDir:  ":memory:",
    }

    // This creates the singleton and makes it available globally
    mqttconfig.InitGlobal(cfg)
    defer mqttconfig.GetClient().Disconnect(250)

    // Now you can use mqtt.GetClient() anywhere in your app
}
```

### Use Anywhere with `GetClient()`

After initialization, you can access the client from any package:

#### In a Handler
```go
// internal/handlers/device_handler.go
package handlers

import (
    mqttconfig "your-project/internal/config/mqtt"
)

type DeviceHandler struct{}

func (h *DeviceHandler) PublishData(deviceID string, data string) error {
    // Get the global client - no need to pass it around!
    client := mqttconfig.GetClient()
    
    topic := fmt.Sprintf("device/%s/data", deviceID)
    return client.Publish(topic, 0, false, data)
}
```

#### In a Service
```go
// internal/services/alert_service.go
package services

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
    mqttconfig "your-project/internal/config/mqtt"
)

type AlertService struct{}

func NewAlertService() *AlertService {
    // Subscribe when service is created
    client := mqttconfig.GetClient()
    
    client.Subscribe("device/+/alerts", 1, func(c mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Alert: %s\n", string(msg.Payload()))
    })
    
    return &AlertService{}
}

func (s *AlertService) SendAlert(message string) error {
    client := mqttconfig.GetClient()
    return client.Publish("alerts/critical", 1, false, message)
}
```

#### In Multiple Subscribers
```go
// internal/subscribers/device_subscriber.go
package subscribers

func StartDeviceMonitoring() {
    client := mqttconfig.GetClient()
    
    client.Subscribe("device/#", 0, func(c mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Device: %s = %s\n", msg.Topic(), string(msg.Payload()))
    })
}

// internal/subscribers/metrics_subscriber.go
package subscribers

func StartMetricsCollection() {
    client := mqttconfig.GetClient()
    
    client.Subscribe("metrics/#", 0, func(c mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Metrics: %s = %s\n", msg.Topic(), string(msg.Payload()))
    })
}
```

### Check Initialization Status

```go
if mqttconfig.IsInitialized() {
    fmt.Println("MQTT client is ready")
    client := mqttconfig.GetClient()
} else {
    log.Fatal("MQTT not initialized! Call InitGlobal() first")
}
```

### Get Current Configuration

```go
config := mqttconfig.GetConfig()
fmt.Printf("Broker: %s, ClientID: %s\n", config.Broker, config.ClientID)
```

## API Reference

### GetMQTTClient(cfg MQTTConfig) *MQTTClient
Returns the singleton MQTT client instance. Creates and connects the client on first call.

```go
cfg := mqttconfig.MQTTConfig{
    Broker:    "tcp://localhost:1883",  // MQTT broker URL
    ClientID:  "unique-client-id",       // Unique identifier
    Username:  "",                        // Optional
    Password:  "",                        // Optional
    CleanSess: true,                      // Clean session flag
    StoreDir:  ":memory:",                // ":memory:" or file path
}

client := mqttconfig.GetMQTTClient(cfg)
```

### Subscribe(topic string, qos byte, handler MessageHandler) error
Subscribes to a topic with a message handler.

**Parameters:**
- `topic` - MQTT topic (supports wildcards `#` and `+`)
- `qos` - Quality of Service (0, 1, or 2)
- `handler` - Function to handle incoming messages

**Example - Fixed Topic:**
```go
client.Subscribe("device/sensors", 0, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Sensor data: %s\n", string(msg.Payload()))
})
```

**Example - Wildcard with Routing:**
```go
client.Subscribe("device/#", 1, func(c mqtt.Client, msg mqtt.Message) {
    topic := msg.Topic()
    payload := string(msg.Payload())
    
    switch topic {
    case "device/alerts":
        fmt.Printf("🚨 Alert: %s\n", payload)
    case "device/metrics":
        fmt.Printf("📊 Metrics: %s\n", payload)
    case "device/status":
        fmt.Printf("ℹ️ Status: %s\n", payload)
    default:
        fmt.Printf("Other: %s = %s\n", topic, payload)
    }
})
```

### Publish(topic string, qos byte, retained bool, payload interface{}) error
Publishes a message to a topic.

**Parameters:**
- `topic` - MQTT topic to publish to
- `qos` - Quality of Service (0, 1, or 2)
- `retained` - Whether the broker should retain the message
- `payload` - Message payload (string, []byte, or any serializable type)

**Examples:**
```go
// Simple message
client.Publish("device/sensors", 0, false, "25.5°C")

// Retained message (delivered to future subscribers)
client.Publish("device/status", 0, true, "online")

// QoS 1 for important messages
client.Publish("device/alerts", 1, false, "High temperature!")

// JSON payload
client.Publish("device/data", 0, false, `{"temp": 25.5, "humidity": 60}`)
```

### Unsubscribe(topics ...string) error
Unsubscribes from one or more topics.

```go
client.Unsubscribe("device/sensors")
client.Unsubscribe("device/#", "sensor/+/temperature")
```

### IsConnected() bool
Returns the current connection status.

```go
if client.IsConnected() {
    fmt.Println("Connected to MQTT broker")
}
```

### Disconnect(quiesce uint)
Gracefully disconnects the client.

```go
client.Disconnect(250) // Wait 250ms for pending messages
```

## Wildcard Patterns

### Multi-level Wildcard (#)
Matches **multiple levels** in the topic hierarchy.

| Pattern | Matches | Does NOT Match |
|---------|---------|----------------|
| `device/#` | `device/sensors`<br>`device/alerts/critical`<br>`device/floor1/room2/sensor` | `devices/sensor`<br>`home/device` |
| `home/+/temperature` | `home/living-room/temperature`<br>`home/bedroom/temperature` | `home/temperature`<br>`home/room1/sensor/temperature` |

**Example:**
```go
// Subscribe to all device topics
client.Subscribe("device/#", 0, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Device topic: %s = %s\n", msg.Topic(), string(msg.Payload()))
})
```

### Single-level Wildcard (+)
Matches **exactly one level** in the topic hierarchy.

| Pattern | Matches | Does NOT Match |
|---------|---------|----------------|
| `sensor/+/temperature` | `sensor/room1/temperature`<br>`sensor/kitchen/temperature` | `sensor/temperature`<br>`sensor/room1/living/temperature` |
| `home/+/+/status` | `home/floor1/room1/status`<br>`home/floor2/room3/status` | `home/floor1/status`<br>`home/floor1/room1/sensor/status` |

**Example:**
```go
// Subscribe to temperature from any room (one level)
client.Subscribe("sensor/+/temperature", 0, func(c mqtt.Client, msg mqtt.Message) {
    room := strings.Split(msg.Topic(), "/")[1]
    fmt.Printf("Room %s temperature: %s\n", room, string(msg.Payload()))
})
```

## Complete Examples

### Example 1: Simple Temperature Monitor

```go
package main

import (
    "fmt"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    mqttconfig "your-project/internal/config/mqtt"
)

func main() {
    cfg := mqttconfig.MQTTConfig{
        Broker:    "tcp://localhost:1883",
        ClientID:  "temp-monitor",
        CleanSess: true,
        StoreDir:  ":memory:",
    }

    client := mqttconfig.GetMQTTClient(cfg)
    defer client.Disconnect(250)

    // Subscribe to temperature sensors
    client.Subscribe("sensor/+/temperature", 0, func(c mqtt.Client, msg mqtt.Message) {
        fmt.Printf("%s: %s\n", msg.Topic(), string(msg.Payload()))
    })

    // Simulate sensor data
    client.Publish("sensor/room1/temperature", 0, false, "22.5°C")
    client.Publish("sensor/room2/temperature", 0, false, "24.0°C")

    time.Sleep(2 * time.Second)
}
```

### Example 2: Home Automation with Routing

```go
package main

import (
    "fmt"
    "time"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    mqttconfig "your-project/internal/config/mqtt"
)

func main() {
    cfg := mqttconfig.MQTTConfig{
        Broker:    "tcp://localhost:1883",
        ClientID:  "home-automation",
        CleanSess: true,
        StoreDir:  ":memory:",
    }

    client := mqttconfig.GetMQTTClient(cfg)
    defer client.Disconnect(250)

    // Subscribe to all home topics
    client.Subscribe("home/#", 1, func(c mqtt.Client, msg mqtt.Message) {
        topic := msg.Topic()
        payload := string(msg.Payload())

        switch topic {
        case "home/lights/living-room":
            fmt.Printf("💡 Living room lights: %s\n", payload)
        case "home/thermostat/target":
            fmt.Printf("🌡️  Thermostat set to: %s\n", payload)
        case "home/security/door":
            fmt.Printf("🚪 Door status: %s\n", payload)
        case "home/alerts":
            fmt.Printf("⚠️  ALERT: %s\n", payload)
        default:
            fmt.Printf("📬 %s: %s\n", topic, payload)
        }
    })

    // Simulate home automation events
    client.Publish("home/lights/living-room", 0, false, "ON")
    client.Publish("home/thermostat/target", 1, false, "22°C")
    client.Publish("home/security/door", 1, false, "LOCKED")
    client.Publish("home/alerts", 1, false, "Motion detected in garage")

    time.Sleep(2 * time.Second)
}
```

### Example 3: Multiple Subscriptions

```go
// Subscribe to different topic patterns separately
client.Subscribe("sensors/temperature/#", 0, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Temperature: %s = %s\n", msg.Topic(), string(msg.Payload()))
})

client.Subscribe("sensors/humidity/#", 0, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Humidity: %s = %s\n", msg.Topic(), string(msg.Payload()))
})

client.Subscribe("alerts/+", 1, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("🚨 ALERT from %s: %s\n", msg.Topic(), string(msg.Payload()))
})
```

## QoS Levels

| QoS | Name | Description | Use Case |
|-----|------|-------------|----------|
| **0** | At most once | Fire and forget, no acknowledgment | Non-critical sensor data, frequent updates |
| **1** | At least once | Acknowledged, may receive duplicates | Important messages that can tolerate duplicates |
| **2** | Exactly once | Assured delivery, no duplicates (slowest) | Critical commands, financial transactions |

**Example:**
```go
// QoS 0 for frequent sensor updates
client.Publish("sensor/temperature", 0, false, "25.5")

// QoS 1 for important notifications
client.Publish("alerts/high-temp", 1, false, "Temperature exceeded!")

// QoS 2 for critical commands
client.Publish("device/shutdown", 2, false, "SHUTDOWN")
```

## Docker Configuration

For Docker environments, use `host.docker.internal` to access the host machine:

```go
cfg := mqttconfig.MQTTConfig{
    Broker:    "tcp://host.docker.internal:1883",
    ClientID:  "docker-app",
    Username:  "admin",
    Password:  "secret",
    CleanSess: false,
    StoreDir:  "./mqtt-store",
}
```

**Start Mosquitto broker:**
```bash
docker run -d -p 1883:1883 -p 9001:9001 eclipse-mosquitto
```

## Best Practices

1. **Unique Client IDs** - Each client must have a unique ID to avoid conflicts
2. **Use Wildcards Wisely** - Be specific to avoid receiving too many messages
3. **Handle Errors** - Check errors from Publish and Subscribe
4. **Graceful Shutdown** - Always call `Disconnect()` before exiting
5. **Appropriate QoS** - Use QoS 0 for non-critical data, QoS 1-2 for important messages
6. **Retained Messages** - Use for status messages that new subscribers should receive
7. **Topic Routing** - Use switch statements in handlers for wildcard subscriptions

## Troubleshooting

### No Messages Received
- Verify broker is running
- Check topic names match exactly (case-sensitive)
- Ensure subscription happens before publishing
- Check QoS compatibility

### Connection Failed
```go
if !client.IsConnected() {
    fmt.Println("Not connected to broker")
}
```

**Solutions:**
- Verify broker URL: `tcp://localhost:1883`
- Check broker is running: `docker ps` or `systemctl status mosquitto`
- Verify firewall allows port 1883
- Check credentials if authentication is enabled

### Message Not Published
- Check return error from `Publish()`
- Verify client is connected
- Check QoS settings
- Ensure topic name is valid (no wildcards in publish)

### Testing with CLI
```bash
# Subscribe to test
mosquitto_sub -h localhost -t "device/#" -v

# Publish test message
mosquitto_pub -h localhost -t "device/sensors" -m "test message"
```

## Advanced Topics

### Retained Messages
Retained messages are delivered to new subscribers immediately:

```go
// Publish retained status
client.Publish("device/status", 0, true, "online")

// Clear retained message
client.Publish("device/status", 0, true, "")
```

### Clean Session
- `CleanSess: true` - Don't persist subscriptions across reconnects
- `CleanSess: false` - Maintain subscriptions and queued messages

### Message Store
- `:memory:` - In-memory store (default)
- `./mqtt-store` - File-based store for persistence

```go
cfg := mqttconfig.MQTTConfig{
    CleanSess: false,
    StoreDir:  "./mqtt-store", // Persists messages across restarts
}
```
