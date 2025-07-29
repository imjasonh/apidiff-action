package internal

// Helper is an internal helper function
func Helper(s string) string {
	return s + "!"
}

// InternalType is an internal type
type InternalType struct {
	Value int
}

// Process does internal processing
func (i *InternalType) Process() int {
	return i.Value * 2
}
