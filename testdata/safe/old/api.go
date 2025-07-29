package api

// Config represents configuration
type Config struct {
	Host string
	Port int
}

// NewConfig creates a new config
func NewConfig(host string, port int) *Config {
	return &Config{
		Host: host,
		Port: port,
	}
}
