package credsutil

import (
	"time"

	"fmt"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/helper/keysutil"
)

// CredentialsProducer can be used as an embeded interface in the Database
// definition. It implements the methods for generating user information for a
// particular database type and is used in all the builtin database types.
type CredentialsProducer interface {
	GenerateUsername(usernameConfig dbplugin.UsernameConfig) (string, error)
	GeneratePassword() (string, error)
	GenerateExpiration(ttl time.Time) (string, error)
}

const (
	reqStr    = `A1a-`
	minStrLen = 10
)

// RandomAlphaNumeric returns a random string of characters [A-Za-z0-9-]
// of the provided length. The string generated takes up to 4 characters
// of space that are predefined and prepended to ensure password
// character requirements. It also requires a min length of 10 characters.
func RandomAlphaNumeric(length int, prependA1a bool) (string, error) {
	if length < minStrLen {
		return "", fmt.Errorf("minimum length of %d is required", minStrLen)
	}

	var prefix string
	if prependA1a {
		prefix = reqStr
	}

	buf, err := uuid.GenerateRandomBytes(length - len(prefix))
	if err != nil {
		return "", err
	}

	output := (prefix + keysutil.Base62Encode(buf))[:length]

	return output, nil
}
