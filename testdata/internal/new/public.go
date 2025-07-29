package example

// PublicFunc is a public function
// BREAKING: Added required parameter
func PublicFunc(s string, uppercase bool) string {
	if uppercase {
		return s
	}
	return s
}

// PublicType is a public type
type PublicType struct {
	Name  string
	Email string // Added field - compatible
}
