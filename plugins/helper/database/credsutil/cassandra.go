package credsutil

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

// CassandraCredentialsProducer implements CredentialsProducer and provides an
// interface for cassandra databases to generate user information.
type CassandraCredentialsProducer struct{}

func (ccp *CassandraCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	username := fmt.Sprintf("vault_%s_%s_%d", displayName, userUUID, time.Now().Unix())
	username = strings.Replace(username, "-", "_", -1)

	return username, nil
}

func (ccp *CassandraCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (ccp *CassandraCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return "", nil
}
