package example

import "fmt"

// Greeter provides greeting functionality
type Greeter struct {
	Name     string
	Language string // Added field - compatible change
}

// Greet returns a greeting message
// Breaking change: added parameter
func (g *Greeter) Greet(formal bool) string {
	if formal {
		return fmt.Sprintf("Good day, %s!", g.Name)
	}
	return fmt.Sprintf("Hello, %s!", g.Name)
}

// Add adds two integers
func Add(a, b int) int {
	return a + b
}

// Multiply multiplies two integers (new function - compatible)
func Multiply(a, b int) int {
	return a * b
}

// Config holds configuration
type Config struct {
	Host    string
	Port    int
	Timeout int // Added field - compatible change
}

// Status represents service status
type Status int

const (
	StatusUnknown Status = iota
	StatusReady
	StatusError
	StatusStarting // Added constant - compatible change
)
