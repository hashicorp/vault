package cassandra

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

// cassandraCredentialsProducer implements CredentialsProducer and provides an
// interface for cassandra databases to generate user information.
type cassandraCredentialsProducer struct{}

func (ccp *cassandraCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	username := fmt.Sprintf("vault_%s_%s_%d", displayName, userUUID, time.Now().Unix())
	username = strings.Replace(username, "-", "_", -1)

	return username, nil
}

func (ccp *cassandraCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (ccp *cassandraCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return "", nil
}
