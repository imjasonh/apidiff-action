package example

// Greeter handles greetings
type Greeter struct {
	Name string
}

// Greet returns a greeting
func (g *Greeter) Greet() string {
	return "Hello, " + g.Name
}

// Process handles data processing
func Process(data string) (string, error) {
	return data, nil
}

// MaxRetries is the maximum number of retries
const MaxRetries = 3
