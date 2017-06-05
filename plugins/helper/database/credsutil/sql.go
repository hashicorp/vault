package credsutil

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
)

// SQLCredentialsProducer implements CredentialsProducer and provides a generic credentials producer for most sql database types.
type SQLCredentialsProducer struct {
	DisplayNameLen int
	RoleNameLen    int
	UsernameLen    int
	Separator      string
}

func (scp *SQLCredentialsProducer) GenerateUsername(config dbplugin.UsernameConfig) (string, error) {
	displayName := config.DisplayName
	if scp.DisplayNameLen > 0 && len(displayName) > scp.DisplayNameLen {
		displayName = displayName[:scp.DisplayNameLen]
	}
	roleName := config.RoleName
	if scp.RoleNameLen > 0 && len(roleName) > scp.RoleNameLen {
		roleName = roleName[:scp.RoleNameLen]
	}

	userUUID, err := RandomAlphaNumericOfLen(20)
	if err != nil {
		return "", err
	}

	username := strings.Join([]string{"v", displayName, roleName, string(userUUID), fmt.Sprint(time.Now().UTC().Unix())}, scp.Separator)
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
