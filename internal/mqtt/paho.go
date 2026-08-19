package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqttp "github.com/eclipse/paho.mqtt.golang"
)

type pahoClient struct {
	cfg    Config
	client mqttp.Client
	cache  map[string]Light
	mu     sync.RWMutex
}

// NewPahoClient creates a new Paho-based MQTT client wrapper.
func NewPahoClient(cfg Config) Client {
	return &pahoClient{
		cfg:   cfg,
		cache: make(map[string]Light),
	}
}

func (p *pahoClient) Connect(ctx context.Context) error {
	opts := mqttp.NewClientOptions()
	opts.AddBroker(p.cfg.Broker)
	opts.SetClientID(p.cfg.ClientID)
	if p.cfg.Username != "" {
		opts.SetUsername(p.cfg.Username)
		opts.SetPassword(p.cfg.Password)
	}
	opts.SetKeepAlive(time.Duration(p.cfg.KeepAlive) * time.Second)
	opts.SetAutoReconnect(true)
	opts.OnConnect = func(c mqttp.Client) {
		// subscribe to light state topics to populate cache
		// wildcard subscription example: bridge/lights/+/state
		_ = c.Subscribe("bridge/lights/+/state", 0, p.messageHandler)
	}

	p.client = mqttp.NewClient(opts)

	// try connect with context cancelation
	done := make(chan error, 1)
	go func() {
		if token := p.client.Connect(); token.Wait() && token.Error() != nil {
			done <- token.Error()
			return
		}
		done <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (p *pahoClient) Disconnect() {
	if p.client != nil && p.client.IsConnected() {
		p.client.Disconnect(250)
	}
}

func (p *pahoClient) GetLights(ctx context.Context) (map[string]Light, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	// return a copy to avoid races
	out := make(map[string]Light, len(p.cache))
	for k, v := range p.cache {
		out[k] = v
	}
	return out, nil
}

func (p *pahoClient) PublishCommand(ctx context.Context, topic string, payload []byte) error {
	if p.client == nil {
		return fmt.Errorf("mqtt client not initialized")
	}
	token := p.client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (p *pahoClient) Subscribe(topic string, handler MessageHandler) (UnsubscribeFunc, error) {
	if p.client == nil {
		return nil, fmt.Errorf("mqtt client not initialized")
	}

	wrapped := func(_ mqttp.Client, msg mqttp.Message) {
		handler(msg.Topic(), msg.Payload())
	}

	token := p.client.Subscribe(topic, 0, wrapped)
	token.Wait()
	if err := token.Error(); err != nil {
		return nil, err
	}

	return func() { p.client.Unsubscribe(topic) }, nil
}

func (p *pahoClient) messageHandler(client mqttp.Client, msg mqttp.Message) {
	// Expect payload to be JSON representing Light (or state)
	var l Light
	if err := json.Unmarshal(msg.Payload(), &l); err != nil {
		return
	}

	// derive id from topic if not present: bridge/lights/{id}/state
	topic := msg.Topic()
	var id string
	fmt.Sscanf(topic, "bridge/lights/%s/state", &id)
	// trim possible trailing parts
	if id != "" {
		// remove trailing slash or anything after id
		for i, c := range id {
			if c == '/' {
				id = id[:i]
				break
			}
		}
		l.ID = id
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache[l.ID] = l
}
