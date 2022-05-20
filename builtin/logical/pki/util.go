package pki

import (
	"context"
	"crypto"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/vault/sdk/helper/certutil"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault/sdk/framework"

	"github.com/hashicorp/vault/sdk/helper/errutil"
)

const (
	managedKeyNameArg = "managed_key_name"
	managedKeyIdArg   = "managed_key_id"
	defaultRef        = "default"
)

var (
	nameMatcher        = regexp.MustCompile("^" + framework.GenericNameRegex(issuerRefParam) + "$")
	errIssuerNameInUse = errutil.UserError{Err: "issuer name already in use"}
	errKeyNameInUse    = errutil.UserError{Err: "key name already in use"}
)

func normalizeSerial(serial string) string {
	return strings.Replace(strings.ToLower(serial), ":", "-", -1)
}

func denormalizeSerial(serial string) string {
	return strings.Replace(strings.ToLower(serial), "-", ":", -1)
}

func kmsRequested(input *inputBundle) bool {
	return kmsRequestedFromFieldData(input.apiData)
}

func kmsRequestedFromFieldData(data *framework.FieldData) bool {
	exportedStr, ok := data.GetOk("exported")
	if !ok {
		return false
	}
	return exportedStr.(string) == "kms"
}

func existingKeyRequested(input *inputBundle) bool {
	return existingKeyRequestedFromFieldData(input.apiData)
}

func existingKeyRequestedFromFieldData(data *framework.FieldData) bool {
	exportedStr, ok := data.GetOk("exported")
	if !ok {
		return false
	}
	return exportedStr.(string) == "existing"
}

type managedKeyId interface {
	String() string
}

type (
	UUIDKey string
	NameKey string
)

func (u UUIDKey) String() string {
	return string(u)
}

func (n NameKey) String() string {
	return string(n)
}

type managedKeyInfo struct {
	publicKey crypto.PublicKey
	keyType   certutil.PrivateKeyType
	name      NameKey
	uuid      UUIDKey
}

// getManagedKeyId returns a NameKey or a UUIDKey, whichever was specified in the
// request API data.
func getManagedKeyId(data *framework.FieldData) (managedKeyId, error) {
	name, UUID, err := getManagedKeyNameOrUUID(data)
	if err != nil {
		return nil, err
	}

	var keyId managedKeyId = NameKey(name)
	if len(UUID) > 0 {
		keyId = UUIDKey(UUID)
	}

	return keyId, nil
}

func getKeyRefWithErr(data *framework.FieldData) (string, error) {
	keyRef := getKeyRef(data)

	if len(keyRef) == 0 {
		return "", errutil.UserError{Err: fmt.Sprintf("missing argument key_ref for existing type")}
	}

	return keyRef, nil
}

func getManagedKeyNameOrUUID(data *framework.FieldData) (name string, UUID string, err error) {
	getApiData := func(argName string) (string, error) {
		arg, ok := data.GetOk(argName)
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

func getIssuerName(ctx context.Context, s logical.Storage, data *framework.FieldData) (string, error) {
	issuerName := ""
	issuerNameIface, ok := data.GetOk("issuer_name")
	if ok {
		issuerName = strings.TrimSpace(issuerNameIface.(string))

		if strings.ToLower(issuerName) == defaultRef {
			return issuerName, errutil.UserError{Err: "reserved keyword 'default' can not be used as issuer name"}
		}

		if !nameMatcher.MatchString(issuerName) {
			return issuerName, errutil.UserError{Err: "issuer name contained invalid characters"}
		}
		issuerId, err := resolveIssuerReference(ctx, s, issuerName)
		if err == nil {
			return issuerName, errIssuerNameInUse
		}

		if err != nil && issuerId != IssuerRefNotFound {
			return issuerName, errutil.InternalError{Err: err.Error()}
		}
	}
	return issuerName, nil
}

func getKeyName(ctx context.Context, s logical.Storage, data *framework.FieldData) (string, error) {
	keyName := ""
	keyNameIface, ok := data.GetOk(keyNameParam)
	if ok {
		keyName = strings.TrimSpace(keyNameIface.(string))

		if strings.ToLower(keyName) == defaultRef {
			return "", errutil.UserError{Err: "reserved keyword 'default' can not be used as key name"}
		}

		if !nameMatcher.MatchString(keyName) {
			return "", errutil.UserError{Err: "key name contained invalid characters"}
		}
		keyId, err := resolveKeyReference(ctx, s, keyName)
		if err == nil {
			return "", errKeyNameInUse
		}

		if err != nil && keyId != KeyRefNotFound {
			return "", errutil.InternalError{Err: err.Error()}
		}
	}
	return keyName, nil
}

func getIssuerRef(data *framework.FieldData) string {
	return extractRef(data, issuerRefParam)
}

func getKeyRef(data *framework.FieldData) string {
	return extractRef(data, keyRefParam)
}

func extractRef(data *framework.FieldData, paramName string) string {
	value := strings.TrimSpace(data.Get(paramName).(string))
	if strings.EqualFold(value, defaultRef) {
		return defaultRef
	}
	return value
}
