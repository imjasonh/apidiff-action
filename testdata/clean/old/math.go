package math

// Add adds two integers
func Add(a, b int) int {
	return a + b
}

// Subtract subtracts b from a
func Subtract(a, b int) int {
	return a - b
}

// Calculator provides basic math operations
type Calculator struct {
	precision int
}

// NewCalculator creates a new calculator
func NewCalculator() *Calculator {
	return &Calculator{precision: 2}
}

// Compute performs a calculation
func (c *Calculator) Compute(operation string, a, b float64) float64 {
	switch operation {
	case "add":
		return a + b
	case "subtract":
		return a - b
	default:
		return 0
	}
}
