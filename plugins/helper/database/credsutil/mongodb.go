package credsutil

import (
	"fmt"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

// MongoDBCredentialsProducer implements CredentialsProducer and provides an
// interface for databases to generate user information.
type MongoDBCredentialsProducer struct{}

func (cp *MongoDBCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	username := fmt.Sprintf("vault-%s-%s", displayName, userUUID)

	return username, nil
}

func (cp *MongoDBCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (cp *MongoDBCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return "", nil
}
