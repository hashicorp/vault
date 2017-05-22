package mongodb

import (
	"fmt"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

// mongoDBCredentialsProducer implements CredentialsProducer and provides an
// interface for databases to generate user information.
type mongoDBCredentialsProducer struct{}

func (cp *mongoDBCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	username := fmt.Sprintf("vault-%s-%s", displayName, userUUID)

	return username, nil
}

func (cp *mongoDBCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (cp *mongoDBCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return "", nil
}
