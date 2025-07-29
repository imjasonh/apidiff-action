package api

// Config represents configuration
type Config struct {
	Host    string
	Port    int
	Timeout int // Added field (safe)
}

// NewConfig creates a new config
func NewConfig(host string, port int) *Config {
	return &Config{
		Host: host,
		Port: port,
	}
}

// NewConfigWithTimeout creates a new config with timeout (new function - safe)
func NewConfigWithTimeout(host string, port int, timeout int) *Config {
	return &Config{
		Host:    host,
		Port:    port,
		Timeout: timeout,
	}
}

// DefaultTimeout is the default timeout value (new constant - safe)
const DefaultTimeout = 30
