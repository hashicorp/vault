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

func (scg *sqlCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	// Generate the username, password and expiration. PG limits user to 63 characters
	if scg.displayNameLen > 0 && len(displayName) > scg.displayNameLen {
		displayName = displayName[:scg.displayNameLen]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	username := fmt.Sprintf("%s-%s", displayName, userUUID)
	if scg.usernameLen > 0 && len(username) > scg.usernameLen {
		username = username[:scg.usernameLen]
	}

	return username, nil
}

func (scg *sqlCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (scg *sqlCredentialsProducer) GenerateExpiration(ttl time.Duration) string {
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
