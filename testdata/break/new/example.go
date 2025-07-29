package example

import "context"

// Greeter handles greetings
type Greeter struct {
	Name     string
	Language string // Added field (safe)
}

// Greet returns a greeting
// Breaking change: added parameter
func (g *Greeter) Greet(formal bool) string {
	if formal {
		return "Greetings, " + g.Name
	}
	return "Hello, " + g.Name
}

// Process handles data processing
// Breaking change: added context parameter
func Process(ctx context.Context, data string) (string, error) {
	return data, nil
}

// MaxRetries was removed (breaking change)
