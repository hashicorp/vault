package credsutil

import (
	"fmt"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

// SQLCredentialsProducer implements CredentialsProducer and provides a generic credentials producer for most sql database types.
type SQLCredentialsProducer struct {
	DisplayNameLen int
	UsernameLen    int
}

func (scp *SQLCredentialsProducer) GenerateUsername(displayName string) (string, error) {
	if scp.DisplayNameLen > 0 && len(displayName) > scp.DisplayNameLen {
		displayName = displayName[:scp.DisplayNameLen]
	}
	userUUID, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	username := fmt.Sprintf("v-%s-%s", displayName, userUUID)
	if scp.UsernameLen > 0 && len(username) > scp.UsernameLen {
		username = username[:scp.UsernameLen]
	}

	return username, nil
}

func (scp *SQLCredentialsProducer) GeneratePassword() (string, error) {
	password, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}

	return password, nil
}

func (scp *SQLCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return ttl.Format("2006-01-02 15:04:05-0700"), nil
}
