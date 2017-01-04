package dbs

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

type CredentialsProducer interface {
	GenerateUsername(displayName string) (string, error)
	GeneratePassword() (string, error)
	GenerateExpiration(ttl time.Duration) string
}

// sqlCredentialsProducer impliments CredentialsProducer and provides a generic credentials producer for most sql database types.
type sqlCredentialsProducer struct {
	displayNameLen int
	usernameLen    int
}

func (scp *sqlCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	// Generate the username, password and expiration. PG limits user to 63 characters
	if scp.displayNameLen > 0 && len(displayName) > scp.displayNameLen {
		displayName = displayName[:scp.displayNameLen]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	username := fmt.Sprintf("%s-%s", displayName, userUUID)
	if scp.usernameLen > 0 && len(username) > scp.usernameLen {
		username = username[:scp.usernameLen]
	}

	return username, nil
}

func (scp *sqlCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (scp *sqlCredentialsProducer) GenerateExpiration(ttl time.Duration) string {
	return time.Now().
		Add(ttl).
		Format("2006-01-02 15:04:05-0700")
}

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

func (ccp *cassandraCredentialsProducer) GenerateExpiration(ttl time.Duration) string {
	return ""
}
