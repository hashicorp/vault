package cacheboltdb

const (
	// SecretLeaseType - Bucket/type for leases with secret info
	SecretLeaseType = "secret-lease"

	// AuthLeaseType - Bucket/type for leases with auth info
	AuthLeaseType = "auth-lease"

	// TokenType - Bucket/type for auto-auth tokens
	TokenType = "token"

	// AutoAuthToken - key for the latest auto-auth token
	AutoAuthToken = "auto-auth-token"
)

// Storage interface for persistent storage
type Storage interface {
	// Set saves an Index item in the persistent storage, with a string key,
	// []byte value, and type of Index
	Set(string, []byte, string) error

	// Delete an Index item from the persistent storage
	Delete(id string) error

	// GetByType - retrieve a list of serialized Index's by type
	GetByType(string) ([][]byte, error)

	// GetAutoAuthToken - retrieve the latest auto-auth token if present
	GetAutoAuthToken() ([]byte, error)

	// Close the persistent storage
	Close() error

	// Clear the persistent storage
	Clear() error
}
