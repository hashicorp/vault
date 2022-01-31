package pki

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/errutil"
)

const (
	managedKeyNameArg = "managed_key_name"
	managedKeyIdArg = "managed_key_id"
)

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

func getManagedKeyNameOrUUID(input *inputBundle) (name string, UUID string, err error) {
	getApiData := func(argName string) (string, error) {
		arg, ok := input.apiData.GetOk(argName)
		if !ok {
			return "", nil
		}

		argValue, ok := arg.(string)
		if !ok {
			return "", errutil.UserError{Err: fmt.Sprintf("invalid type for argument %s", argName)}
		}

		return strings.TrimSpace(argValue), nil
	}

	keyName, err := getApiData(managedKeyNameArg)
	keyUUID, err2 := getApiData(managedKeyIdArg)
	switch {
	case err != nil:
		return "", "", err
	case err2 != nil:
		return "", "", err2
	case len(keyName) == 0 && len(keyUUID) == 0:
		return "", "", errutil.UserError{Err: fmt.Sprintf("missing argument %s or %s", managedKeyNameArg, managedKeyIdArg)}
	case len(keyName) > 0 && len(keyUUID) > 0:
		return "", "", errutil.UserError{Err: fmt.Sprintf("only one argument of %s or %s should be specified", managedKeyNameArg, managedKeyIdArg)}
	}

	return keyName, keyUUID, nil
}
