package credsutil

import "time"

// CredentialsProducer can be used as an embeded interface in the Database
// definition. It implements the methods for generating user information for a
// particular database type and is used in all the builtin database types.
type CredentialsProducer interface {
	GenerateUsername(displayName string) (string, error)
	GeneratePassword() (string, error)
	GenerateExpiration(ttl time.Time) (string, error)
}
