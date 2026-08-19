package mqtt

import "context"

// Light represents a simplified Hue light state used by the HTTP API.
type Light struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	On         bool   `json:"on"`
	Brightness int    `json:"bri"`
	Hue        int    `json:"hue"`
	Sat        int    `json:"sat"`
}

// MessageHandler processes incoming MQTT messages for a topic.
type MessageHandler func(topic string, payload []byte)

// Unsubscribe function returned by Subscribe.
type UnsubscribeFunc func()

// Client is a minimal MQTT client abstraction used by the HTTP API.
type Client interface {
	Connect(ctx context.Context) error
	Disconnect()
	// GetLights returns a map of light id to Light. Implementations may maintain
	// an in-memory cache populated from telemetry topics.
	GetLights(ctx context.Context) (map[string]Light, error)
	PublishCommand(ctx context.Context, topic string, payload []byte) error
	Subscribe(topic string, handler MessageHandler) (UnsubscribeFunc, error)
}
