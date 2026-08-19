package mqtt

// Config holds MQTT client configuration options.
type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
	// KeepAlive seconds
	KeepAlive int
}
