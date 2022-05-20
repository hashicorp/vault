package pki

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/errutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathListKeys(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "keys/?$",

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback:                    b.pathListKeysHandler,
				ForwardPerformanceStandby:   false,
				ForwardPerformanceSecondary: false,
			},
		},

		HelpSynopsis:    pathListKeysHelpSyn,
		HelpDescription: pathListKeysHelpDesc,
	}
}

const (
	pathListKeysHelpSyn  = `Fetch a list of all issuer keys`
	pathListKeysHelpDesc = `This endpoint allows listing of known backing keys, returning
their identifier and their name (if set).`
)

func (b *backend) pathListKeysHandler(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not list keys until migration has completed"), nil
	}

	var responseKeys []string
	responseInfo := make(map[string]interface{})

	entries, err := listKeys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	config, err := getKeysConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	for _, identifier := range entries {
		key, err := fetchKeyById(ctx, req.Storage, identifier)
		if err != nil {
			return nil, err
		}

		responseKeys = append(responseKeys, string(identifier))
		responseInfo[string(identifier)] = map[string]interface{}{
			keyNameParam: key.Name,
			"is_default": identifier == config.DefaultKeyId,
		}

	}
	return logical.ListResponseWithInfo(responseKeys, responseInfo), nil
}

func pathKey(b *backend) *framework.Path {
	pattern := "key/" + framework.GenericNameRegex(keyRefParam)
	return buildPathKey(b, pattern)
}

func buildPathKey(b *backend, pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,

		Fields: map[string]*framework.FieldSchema{
			keyRefParam: {
				Type:        framework.TypeString,
				Description: `Reference to key; either "default" for the configured default key, an identifier of a key, or the name assigned to the key.`,
				Default:     defaultRef,
			},
			keyNameParam: {
				Type:        framework.TypeString,
				Description: `Human-readable name for this key.`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.pathGetKeyHandler,
				ForwardPerformanceStandby:   false,
				ForwardPerformanceSecondary: false,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathUpdateKeyHandler,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback:                    b.pathDeleteKeyHandler,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathKeysHelpSyn,
		HelpDescription: pathKeysHelpDesc,
	}
}

const (
	pathKeysHelpSyn  = `Fetch a single issuer key`
	pathKeysHelpDesc = `This allows fetching information associated with the underlying key.

:ref can be either the literal value "default", in which case /config/keys
will be consulted for the present default key, an identifier of a key,
or its assigned name value.

Writing to /key/:ref allows updating of the name field associated with
the certificate.
`
)

func (b *backend) pathGetKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not get keys until migration has completed"), nil
	}

	keyRef := data.Get(keyRefParam).(string)
	if len(keyRef) == 0 {
		return logical.ErrorResponse("missing key reference"), nil
	}

	keyId, err := resolveKeyReference(ctx, req.Storage, keyRef)
	if err != nil {
		return nil, err
	}
	if keyId == "" {
		return logical.ErrorResponse("unable to resolve key id for reference" + keyRef), nil
	}

	key, err := fetchKeyById(ctx, req.Storage, keyId)
	if err != nil {
		return nil, err
	}

	respData := map[string]interface{}{
		keyIdParam:   key.ID,
		keyNameParam: key.Name,
		keyTypeParam: string(key.PrivateKeyType),
	}

	if key.isManagedPrivateKey() {
		managedKeyUUID, err := key.getManagedKeyUUID()
		if err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("failed extracting managed key uuid from key id %s (%s): %v", key.ID, key.Name, err)}
		}

		keyInfo, err := getManagedKeyInfo(ctx, b, managedKeyUUID)
		if err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("failed fetching managed key info from key id %s (%s): %v", key.ID, key.Name, err)}
		}

		// To remain consistent across the api responses (mainly generate root/intermediate calls), return the actual
		// type of key, not that it is a managed key.
		respData[keyTypeParam] = string(keyInfo.keyType)
		respData[managedKeyIdArg] = string(keyInfo.uuid)
		respData[managedKeyNameArg] = string(keyInfo.name)
	}

	return &logical.Response{Data: respData}, nil
}

func (b *backend) pathUpdateKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating keys here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not update keys until migration has completed"), nil
	}

	keyRef := data.Get(keyRefParam).(string)
	if len(keyRef) == 0 {
		return logical.ErrorResponse("missing key reference"), nil
	}

	keyId, err := resolveKeyReference(ctx, req.Storage, keyRef)
	if err != nil {
		return nil, err
	}
	if keyId == "" {
		return logical.ErrorResponse("unable to resolve key id for reference" + keyRef), nil
	}

	key, err := fetchKeyById(ctx, req.Storage, keyId)
	if err != nil {
		return nil, err
	}

	newName := data.Get(keyNameParam).(string)
	if len(newName) > 0 && !nameMatcher.MatchString(newName) {
		return logical.ErrorResponse("new key name outside of valid character limits"), nil
	}

	if newName != key.Name {
		key.Name = newName

		err := writeKey(ctx, req.Storage, *key)
		if err != nil {
			return nil, err
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			keyIdParam:   key.ID,
			keyNameParam: key.Name,
			keyTypeParam: key.PrivateKeyType,
		},
	}

	if len(newName) == 0 {
		resp.AddWarning("Name successfully deleted, you will now need to reference this key by it's Id: " + string(key.ID))
	}

	return resp, nil
}

func (b *backend) pathDeleteKeyHandler(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Since we're planning on updating issuers here, grab the lock so we've
	// got a consistent view.
	b.issuersLock.Lock()
	defer b.issuersLock.Unlock()

	if b.useLegacyBundleCaStorage() {
		return logical.ErrorResponse("Can not delete keys until migration has completed"), nil
	}

	keyRef := data.Get(keyRefParam).(string)
	if len(keyRef) == 0 {
		return logical.ErrorResponse("missing key reference"), nil
	}

	keyId, err := resolveKeyReference(ctx, req.Storage, keyRef)
	if err != nil {
		if keyId == KeyRefNotFound {
			// We failed to lookup the key, we should ignore any error here and reply as if it was deleted.
			return nil, nil
		}
		return nil, err
	}

	keyInUse, issuerId, err := isKeyInUse(keyId.String(), ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if keyInUse {
		return logical.ErrorResponse(fmt.Sprintf("Failed to Delete Key.  Key in Use by Issuer: %s", issuerId)), nil
	}

	wasDefault, err := deleteKey(ctx, req.Storage, keyId)
	if err != nil {
		return nil, err
	}

	var response *logical.Response
	if wasDefault {
		msg := fmt.Sprintf("Deleted key %v (via key_ref %v); this was configured as the default key. Operations without an explicit key will not work until a new default is configured.", string(keyId), keyRef)
		b.Logger().Error(msg)
		response = &logical.Response{}
		response.AddWarning(msg)
	}

	return response, nil
}
