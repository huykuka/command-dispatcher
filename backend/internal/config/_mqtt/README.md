# MQTT Client Package

A simplified, singleton-based MQTT client for Go with support for wildcards and easy topic routing.

## Features

- ✅ **Singleton Pattern** - Single client instance shared across application
- ✅ **Thread-Safe** - Safe for concurrent use
- ✅ **Wildcard Support** - Subscribe to multiple topics with `#` and `+`
- ✅ **Simple API** - Just 3 main methods: Subscribe, Publish, Disconnect
- ✅ **Auto-Reconnect** - Built-in connection management
- ✅ **QoS Support** - Quality of Service levels 0, 1, and 2

## Installation

```bash
go get github.com/eclipse/paho.mqtt.golang
```

## Quick Start

```go
package main

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    mqttconfig "your-project/internal/config/mqtt"
)

func main() {
    // Configure client
    cfg := mqttconfig.MQTTConfig{
        Broker:    "tcp://localhost:1883",
        ClientID:  "my-app",
        Username:  "",
        Password:  "",
        CleanSess: true,
        StoreDir:  ":memory:",
    }

    // Get singleton instance
    client := mqttconfig.GetMQTTClient(cfg)
    defer client.Disconnect(250)

    // Subscribe
    client.Subscribe("device/sensors", 0, func(c mqtt.Client, msg mqtt.Message) {
        fmt.Printf("Received: %s\n", string(msg.Payload()))
    })

    // Publish
    client.Publish("device/sensors", 0, false, "Temperature: 25°C")
}
```

## API

### GetMQTTClient(cfg MQTTConfig)

Returns the singleton MQTT client instance.

```go
client := mqttconfig.GetMQTTClient(cfg)
```

### Subscribe(topic, qos, handler)

Subscribe to a topic (fixed or wildcard) with a message handler.

**Fixed topic:**
```go
client.Subscribe("device/sensors", 0, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Println(string(msg.Payload()))
})
```

**Wildcard with routing:**
```go
client.Subscribe("device/#", 1, func(c mqtt.Client, msg mqtt.Message) {
    topic := msg.Topic()
    payload := string(msg.Payload())
    
    switch topic {
    case "device/alerts":
        handleAlert(payload)
    case "device/metrics":
        handleMetrics(payload)
    default:
        handleOther(topic, payload)
    }
})
```

### Publish(topic, qos, retained, payload)

Publish a message to a topic.

```go
// Simple message
client.Publish("device/sensors", 0, false, "25.5°C")

// Retained message
client.Publish("device/status", 0, true, "online")

// Important message with QoS 1
client.Publish("device/alerts", 1, false, "High CPU!")
```

### Unsubscribe(topics...)

```go
client.Unsubscribe("device/sensors", "device/#")
```

### IsConnected()

```go
if client.IsConnected() {
    fmt.Println("Connected")
}
```

### Disconnect(quiesce)

```go
client.Disconnect(250) // Wait 250ms for pending messages
```

## Wildcards

### Multi-level (#)

Matches **multiple levels** in topic hierarchy.

```go
// Subscribe to all device topics
client.Subscribe("device/#", 0, handler)

// Matches:
// - device/sensors
// - device/alerts/critical
// - device/floor1/room2/temp
```

### Single-level (+)

Matches **exactly one level**.

```go
// Subscribe to temperature from any room
client.Subscribe("sensor/+/temperature", 0, handler)

// Matches:
// - sensor/room1/temperature
// - sensor/kitchen/temperature
//
// Does NOT match:
// - sensor/temperature (missing level)
// - sensor/room1/living/temperature (too many levels)
```

## QoS Levels

| QoS | Name | Use Case |
|-----|------|----------|
| **0** | At most once | Non-critical sensor data |
| **1** | At least once | Important notifications |
| **2** | Exactly once | Critical commands |

## Configuration

```go
type MQTTConfig struct {
    Broker    string  // "tcp://localhost:1883"
    ClientID  string  // Unique identifier
    Username  string  // Optional
    Password  string  // Optional
    CleanSess bool    // true = don't persist, false = persist
    StoreDir  string  // ":memory:" or file path
}
```

## Examples

### Example 1: Temperature Monitor

```go
client.Subscribe("sensor/+/temperature", 0, func(c mqtt.Client, msg mqtt.Message) {
    fmt.Printf("%s: %s\n", msg.Topic(), string(msg.Payload()))
})

client.Publish("sensor/room1/temperature", 0, false, "22.5°C")
client.Publish("sensor/room2/temperature", 0, false, "24.0°C")
```

### Example 2: Home Automation

```go
client.Subscribe("home/#", 1, func(c mqtt.Client, msg mqtt.Message) {
    topic := msg.Topic()
    payload := string(msg.Payload())
    
    switch topic {
    case "home/lights/living-room":
        fmt.Printf("💡 Lights: %s\n", payload)
    case "home/thermostat":
        fmt.Printf("🌡️ Thermostat: %s\n", payload)
    case "home/security/door":
        fmt.Printf("🚪 Door: %s\n", payload)
    default:
        fmt.Printf("📬 %s: %s\n", topic, payload)
    }
})

client.Publish("home/lights/living-room", 0, false, "ON")
client.Publish("home/thermostat", 1, false, "22°C")
```

### Example 3: Multiple Subscriptions

```go
// Subscribe to different patterns
client.Subscribe("sensors/temperature/#", 0, handleTemp)
client.Subscribe("sensors/humidity/#", 0, handleHumidity)
client.Subscribe("alerts/+", 1, handleAlerts)
```

## Docker Setup

```bash
# Start Mosquitto broker
docker run -d -p 1883:1883 eclipse-mosquitto
```

```go
// Connect from Docker container to host
cfg := mqttconfig.MQTTConfig{
    Broker:   "tcp://host.docker.internal:1883",
    ClientID: "docker-app",
}
```

## Testing

```bash
# Subscribe using mosquitto CLI
mosquitto_sub -h localhost -t "device/#" -v

# Publish using mosquitto CLI
mosquitto_pub -h localhost -t "device/sensors" -m "test"
```

## Best Practices

1. ✅ **Unique Client IDs** - Avoid connection conflicts
2. ✅ **Error Handling** - Check errors from Publish/Subscribe
3. ✅ **Graceful Shutdown** - Always call Disconnect()
4. ✅ **Specific Wildcards** - Use narrow patterns to reduce traffic
5. ✅ **Appropriate QoS** - Balance reliability vs performance
6. ✅ **Topic Routing** - Use switch statements for wildcard handlers

## Troubleshooting

### Connection Failed
- Check broker is running
- Verify broker URL: `tcp://localhost:1883`
- Check firewall settings
- Verify credentials

### No Messages
- Ensure subscription before publishing
- Check topic names (case-sensitive)
- Verify QoS compatibility
- Test with mosquitto CLI tools

### Test Connection
```go
if !client.IsConnected() {
    log.Fatal("Not connected to MQTT broker")
}
```

## Files

- `mqtt.go` - Main client implementation
- `example.go` - Usage examples
- `USAGE.md` - Detailed usage guide
- `README.md` - This file

## License

This package uses the Eclipse Paho MQTT Go client library.