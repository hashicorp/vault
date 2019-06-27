package credsutil

import (
	"context"
	"time"

	"fmt"

	"github.com/hashicorp/vault/sdk/database/dbplugin"
	"github.com/hashicorp/vault/sdk/helper/base62"
)

// CredentialsProducer can be used as an embedded interface in the Database
// definition. It implements the methods for generating user information for a
// particular database type and is used in all the builtin database types.
type CredentialsProducer interface {
	GenerateCredentials(context.Context) (string, error)
	GenerateUsername(dbplugin.UsernameConfig) (string, error)
	GeneratePassword() (string, error)
	GenerateExpiration(time.Time) (string, error)
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

	randomStr, err := base62.Random(length - len(prefix))
	if err != nil {
		return "", err
	}

	return prefix + randomStr, nil
}
