package internal

// Helper is an internal helper function
// BREAKING: Changed return type (but should be ignored)
func Helper(s string) (string, error) {
	return s + "!", nil
}

// InternalType is an internal type
// BREAKING: Removed field (but should be ignored)
type InternalType struct {
	NewValue string
}

// Process does internal processing
// BREAKING: Changed signature (but should be ignored)
func (i *InternalType) Process(multiplier int) string {
	return i.NewValue
}
