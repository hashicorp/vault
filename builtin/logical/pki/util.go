package pki

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/errutil"
)

const managedKeyNameArg = "managed_key_name"

func normalizeSerial(serial string) string {
	return strings.Replace(strings.ToLower(serial), ":", "-", -1)
}

func denormalizeSerial(serial string) string {
	return strings.Replace(strings.ToLower(serial), "-", ":", -1)
}

func kmsRequested(input *inputBundle) bool {
	exportedStr, ok := input.apiData.GetOk("exported")
	if !ok {
		return false
	}
	return exportedStr.(string) == "kms"
}

func getManagedKeyName(input *inputBundle) (string, error) {
	arg, ok := input.apiData.GetOk(managedKeyNameArg)
	if !ok {
		return "", errutil.UserError{Err: fmt.Sprintf("missing %s argument", managedKeyNameArg)}
	}
	keyName, ok := arg.(string)
	if !ok {
		return "", errutil.UserError{Err: fmt.Sprintf("invalid type for argument %s", managedKeyNameArg)}
	}

	trimmedArg := strings.TrimSpace(keyName)
	if len(trimmedArg) == 0 {
		return "", errutil.UserError{Err: fmt.Sprintf("invalid value for argument %s", managedKeyNameArg)}
	}

	return trimmedArg, nil
}
