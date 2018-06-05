package credsutil

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
)

const (
	NoneLength int = -1
)

// SQLCredentialsProducer implements CredentialsProducer and provides a generic credentials producer for most sql database types.
type SQLCredentialsProducer struct {
	DisplayNameLen int
	RoleNameLen    int
	UsernameLen    int
	Separator      string
}

func (scp *SQLCredentialsProducer) GenerateUsername(config dbplugin.UsernameConfig) (string, error) {
	username := "v"

	displayName := config.DisplayName
	if scp.DisplayNameLen > 0 && len(displayName) > scp.DisplayNameLen {
		displayName = displayName[:scp.DisplayNameLen]
	} else if scp.DisplayNameLen == NoneLength {
		displayName = ""
	}

	if len(displayName) > 0 {
		username = fmt.Sprintf("%s%s%s", username, scp.Separator, displayName)
	}

	roleName := config.RoleName
	if scp.RoleNameLen > 0 && len(roleName) > scp.RoleNameLen {
		roleName = roleName[:scp.RoleNameLen]
	} else if scp.RoleNameLen == NoneLength {
		roleName = ""
	}

	if len(roleName) > 0 {
		username = fmt.Sprintf("%s%s%s", username, scp.Separator, roleName)
	}

	userUUID, err := RandomAlphaNumeric(20, false)
	if err != nil {
		return "", err
	}

	username = fmt.Sprintf("%s%s%s", username, scp.Separator, userUUID)
	username = fmt.Sprintf("%s%s%s", username, scp.Separator, fmt.Sprint(time.Now().Unix()))
	if scp.UsernameLen > 0 && len(username) > scp.UsernameLen {
		username = username[:scp.UsernameLen]
	}

	return username, nil
}

func (scp *SQLCredentialsProducer) GeneratePassword() (string, error) {
	password, err := RandomAlphaNumeric(20, true)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (scp *SQLCredentialsProducer) GenerateExpiration(ttl time.Time) (string, error) {
	return ttl.Format("2006-01-02 15:04:05-0700"), nil
}
