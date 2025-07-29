package math

// Add adds two integers
func Add(a, b int) int {
	// Implementation changed but API is the same
	result := a + b
	return result
}

// Subtract subtracts b from a
func Subtract(a, b int) int {
	// Refactored implementation
	diff := a - b
	return diff
}

// Calculator provides basic math operations
type Calculator struct {
	precision int
}

// NewCalculator creates a new calculator
func NewCalculator() *Calculator {
	// Changed internal default but API is the same
	return &Calculator{precision: 4}
}

// Compute performs a calculation
func (c *Calculator) Compute(operation string, a, b float64) float64 {
	// Refactored with better performance
	var result float64
	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	default:
		result = 0
	}
	return result
}
