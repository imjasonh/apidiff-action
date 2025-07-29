package api

import "time"

// User represents a user in the system
type User struct {
	ID        string
	Name      string
	Email     string    // Added field - compatible change
	CreatedAt time.Time // Added field - compatible change
}

// GetUser retrieves a user by ID
func GetUser(id string) (*User, error) {
	// Mock implementation
	return &User{ID: id, Name: "John Doe"}, nil
}

// GetUsers retrieves all users (new function - compatible)
func GetUsers() ([]*User, error) {
	return []*User{}, nil
}

// UpdateUser updates a user (new function - compatible)
func UpdateUser(user *User) error {
	return nil
}

// Status represents the service status
type Status string

const (
	StatusOK      Status = "ok"
	StatusError   Status = "error"
	StatusPending Status = "pending" // Added constant - compatible change
	StatusUnknown Status = "unknown" // Added constant - compatible change
)

// Service provides API functionality
type Service struct {
	endpoint string
	timeout  time.Duration // Added field - compatible change
}

// NewService creates a new service
func NewService(endpoint string) *Service {
	return &Service{endpoint: endpoint}
}

// NewServiceWithTimeout creates a new service with timeout (new function - compatible)
func NewServiceWithTimeout(endpoint string, timeout time.Duration) *Service {
	return &Service{endpoint: endpoint, timeout: timeout}
}

// Call makes an API call
func (s *Service) Call(method string) error {
	// Mock implementation
	return nil
}

// CallWithContext makes an API call with context (new method - compatible)
func (s *Service) CallWithContext(method string, timeout time.Duration) error {
	return nil
}

// GetStatus returns the service status (new method - compatible)
func (s *Service) GetStatus() Status {
	return StatusOK
}

// ErrorCode represents error codes (new type - compatible)
type ErrorCode int

const (
	ErrNone    ErrorCode = 0 // New constant - compatible
	ErrInvalid ErrorCode = 1 // New constant - compatible
	ErrTimeout ErrorCode = 2 // New constant - compatible
)
