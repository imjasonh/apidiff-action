package example

import "fmt"

// Greeter provides greeting functionality
type Greeter struct {
	Name string
}

// Greet returns a greeting message
func (g *Greeter) Greet() string {
	return fmt.Sprintf("Hello, %s!", g.Name)
}

// Add adds two integers
func Add(a, b int) int {
	return a + b
}

// Config holds configuration
type Config struct {
	Host string
	Port int
}

// Status represents service status
type Status int

const (
	StatusUnknown Status = iota
	StatusReady
	StatusError
)
