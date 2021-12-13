package vault

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chrismalek/oktasdk-go/okta"
	jwt "github.com/dgrijalva/jwt-go"
	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	memdb "github.com/hashicorp/go-memdb"
	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

const (
	mfaMethodTypeTOTP    = "totp"
	mfaMethodTypeDuo     = "duo"
	mfaMethodTypeOkta    = "okta"
	mfaMethodTypePingID  = "pingid"
	memDBMFAConfigsTable = "mfa_configs"

	mfaTOTPKeysPrefix = systemBarrierPrefix + "mfa/totpkeys/"

	// mfaConfigPrefix is the storage path prefix for persisting MFA method
	// configs
	mfaConfigPrefix = "mfa/method/"
)

type totpKey struct {
	Key string `json:"key"`
}

func (b *SystemBackend) memDBMFAMethods(ws memdb.WatchSet) (memdb.ResultIterator, error) {
	txn := b.db.Txn(false)

	iter, err := txn.Get(memDBMFAConfigsTable, "name")
	if err != nil {
		return nil, err
	}

	ws.Add(iter.WatchCh())

	return iter, nil
}

// pathMFAMethodsList is used to list all the Roles registered with the backend.
func (b *SystemBackend) pathMFAMethodsList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	ws := memdb.NewWatchSet()
	iter, err := b.memDBMFAMethods(ws)
	if err != nil {
		return nil, errwrap.Wrapf("failed to fetch iterator for MFA methods in memdb: {{err}}", err)
	}

	var methodNames []string
	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		methodNames = append(methodNames, raw.(*mfa.Config).Name)
	}

	return logical.ListResponse(methodNames), nil
}

// TODO: I think in OSS we only have root, so leaving this function as is should not be a problem
// TODO: If needed, wrap the neccessary functions with it
func validateRootNS(f framework.OperationFunc) framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		if ns == nil {
			return nil, namespace.ErrNoNamespace
		}
		if ns.ID != namespace.RootNamespaceID {
			return logical.ErrorResponse("this API path can only be called from the root namespace"), nil
		}
		return f(ctx, req, d)
	}
}

func (b *SystemBackend) handleMFAGenerateRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.handleMFAGenerateCommon(ctx, req, d.Get("name").(string), req.EntityID)
}

func (b *SystemBackend) handleMFAAdminGenerateUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.handleMFAGenerateCommon(ctx, req, d.Get("name").(string), d.Get("entity_id").(string))
}

func (b *SystemBackend) handleMFAAdminDestroyUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	methodName := d.Get("name").(string)
	entityID := d.Get("entity_id").(string)
	var entity *identity.Entity
	var err error

	if b.Core.identityStore == nil {
		return nil, fmt.Errorf("identity store not set up, cannot service totp mfa requests")
	}

	if entityID == "" {
		return logical.ErrorResponse("missing entity ID"), nil
	}

	if methodName == "" {
		return logical.ErrorResponse("missing method name"), nil
	}

	entity, err = b.Core.identityStore.MemDBEntityByID(entityID, true)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to find entity with ID %q: {{err}}", entityID), err)
	}

	if entity == nil {
		return logical.ErrorResponse("invalid entity ID"), nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, namespace.ErrNoNamespace
	}
	if entity.NamespaceID != ns.ID {
		return logical.ErrorResponse("destruction of TOTP keys must operate on entities in the same namespace"), nil
	}

	mConfig, err := b.MemDBMFAConfigByMethodName(methodName, false)
	if err != nil {
		return nil, err
	}

	if mConfig == nil {
		return logical.ErrorResponse(fmt.Sprintf("configuration for method name %q does not exist", methodName)), nil
	}

	if mConfig.ID == "" {
		return nil, fmt.Errorf("configuration for method name %q does not contain an identifier", methodName)
	}

	if entity.MFASecrets != nil {
		delete(entity.MFASecrets, mConfig.ID)
	}

	err = b.Core.identityStore.upsertEntity(ctx, entity, nil, true)
	if err != nil {
		return nil, errwrap.Wrapf("failed to persist MFA secret in entity: {{err}}", err)
	}

	return nil, nil
}

func (b *SystemBackend) handleMFAGenerateCommon(ctx context.Context, req *logical.Request, methodName, entityID string) (*logical.Response, error) {
	var entity *identity.Entity
	var err error

	if b.Core.identityStore == nil {
		return nil, fmt.Errorf("identity store not set up, cannot service totp mfa requests")
	}

	if entityID == "" {
		return logical.ErrorResponse("missing entity ID"), nil
	}

	if methodName == "" {
		return logical.ErrorResponse("missing method name"), nil
	}

	entity, err = b.Core.identityStore.MemDBEntityByID(entityID, false)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to find entity with ID %q: {{err}}", entityID), err)
	}

	if entity == nil {
		return logical.ErrorResponse("invalid entity ID"), nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, namespace.ErrNoNamespace
	}
	if entity.NamespaceID != ns.ID {
		return logical.ErrorResponse("generation of TOTP keys must operate on entities in the same namespace"), nil
	}

	mConfig, err := b.MemDBMFAConfigByMethodName(methodName, false)
	if err != nil {
		return nil, err
	}

	if mConfig == nil {
		return logical.ErrorResponse(fmt.Sprintf("configuration for method name %q does not exist", methodName)), nil
	}

	if mConfig.ID == "" {
		return nil, fmt.Errorf("configuration for method name %q does not contain an identifier", methodName)
	}

	switch mConfig.Type {
	case mfaMethodTypeTOTP:
		return b.handleMFAGenerateTOTP(ctx, mConfig, entity.ID)
	default:
		return logical.ErrorResponse(fmt.Sprintf("generate not available for MFA type %q", mConfig.Type)), nil
	}
}

func (c *Core) PersistTOTPKey(ctx context.Context, methodID, entityID, key string) error {
	ks := &totpKey{
		Key: key,
	}
	val, err := jsonutil.EncodeJSON(ks)
	if err != nil {
		return err
	}
	if c.barrier.Put(ctx, &logical.StorageEntry{
		Key:   fmt.Sprintf("%s%s/%s", mfaTOTPKeysPrefix, methodID, entityID),
		Value: val,
	}); err != nil {
		return err
	}
	return nil
}

func (c *Core) fetchTOTPKey(ctx context.Context, methodID, entityID string) (string, error) {
	entry, err := c.barrier.Get(ctx, fmt.Sprintf("%s%s/%s", mfaTOTPKeysPrefix, methodID, entityID))
	if err != nil {
		return "", err
	}
	if entry == nil {
		return "", nil
	}

	ks := &totpKey{}
	err = jsonutil.DecodeJSON(entry.Value, ks)
	if err != nil {
		return "", err
	}

	return ks.Key, nil
}

func (b *SystemBackend) handleMFAGenerateTOTP(ctx context.Context, mConfig *mfa.Config, entityID string) (*logical.Response, error) {
	var err error
	var totpConfig *mfa.TOTPConfig

	if b.Core.identityStore == nil {
		return nil, fmt.Errorf("identity store not set up, cannot service totp mfa requests")
	}

	switch mConfig.Config.(type) {
	case *mfa.Config_TOTPConfig:
		totpConfig = mConfig.Config.(*mfa.Config_TOTPConfig).TOTPConfig
	default:
		return logical.ErrorResponse(fmt.Sprintf("unknown MFA config type %q", mConfig.Type)), nil
	}

	b.Core.identityStore.lock.Lock()
	defer b.Core.identityStore.lock.Unlock()

	// Read the entity after acquiring the lock
	entity, err := b.Core.identityStore.MemDBEntityByID(entityID, true)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to find entity with ID %q: {{err}}", entityID), err)
	}

	if entity == nil {
		return logical.ErrorResponse("invalid entity ID"), nil
	}

	if entity.MFASecrets == nil {
		entity.MFASecrets = make(map[string]*mfa.Secret)
	} else {
		_, ok := entity.MFASecrets[mConfig.ID]
		if ok {
			resp := &logical.Response{}
			resp.AddWarning(fmt.Sprintf("Entity already has a secret for MFA method %q", mConfig.Name))
			return resp, nil
		}
	}

	keyObject, err := totplib.Generate(totplib.GenerateOpts{
		Issuer:      totpConfig.Issuer,
		AccountName: entity.ID,
		Period:      uint(totpConfig.Period),
		Digits:      otplib.Digits(totpConfig.Digits),
		Algorithm:   otplib.Algorithm(totpConfig.Algorithm),
		SecretSize:  uint(totpConfig.KeySize),
		Rand:        b.Core.secureRandomReader,
	})
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to generate TOTP key for method name %q: {{err}}", mConfig.Name), err)
	}
	if keyObject == nil {
		return nil, fmt.Errorf("failed to generate TOTP key for method name %q", mConfig.Name)
	}

	totpURL := keyObject.String()

	totpB64Barcode := ""
	if totpConfig.QRSize != 0 {
		barcode, err := keyObject.Image(int(totpConfig.QRSize), int(totpConfig.QRSize))
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate QR code image: {{err}}", err)
		}

		var buff bytes.Buffer
		png.Encode(&buff, barcode)
		totpB64Barcode = base64.StdEncoding.EncodeToString(buff.Bytes())
	}

	if err := b.Core.PersistTOTPKey(ctx, mConfig.ID, entity.ID, keyObject.Secret()); err != nil {
		return nil, errwrap.Wrapf("failed to persist totp key: {{err}}", err)
	}

	entity.MFASecrets[mConfig.ID] = &mfa.Secret{
		MethodName: mConfig.Name,
		Value: &mfa.Secret_TOTPSecret{
			TOTPSecret: &mfa.TOTPSecret{
				Issuer:      totpConfig.Issuer,
				AccountName: entity.ID,
				Period:      uint32(totpConfig.Period),
				Algorithm:   int32(totpConfig.Algorithm),
				Digits:      int32(totpConfig.Digits),
				Skew:        uint32(totpConfig.Skew),
				KeySize:     uint32(totpConfig.KeySize),
			},
		},
	}

	err = b.Core.identityStore.upsertEntity(ctx, entity, nil, true)
	if err != nil {
		return nil, errwrap.Wrapf("failed to persist MFA secret in entity: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"url":     totpURL,
			"barcode": totpB64Barcode,
		},
	}, nil
}

func (b *SystemBackend) handleTOTPConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.handleConfigUpdateCommon(ctx, req, d, mfaMethodTypeTOTP)
}

func (b *SystemBackend) handleOktaConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.handleConfigUpdateCommon(ctx, req, d, mfaMethodTypeOkta)
}

func (b *SystemBackend) handleDuoConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.handleConfigUpdateCommon(ctx, req, d, mfaMethodTypeDuo)
}

func (b *SystemBackend) handlePingIDConfigUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return b.handleConfigUpdateCommon(ctx, req, d, mfaMethodTypePingID)
}

func (b *SystemBackend) handleConfigUpdateCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, methodType string) (*logical.Response, error) {
	var err error
	methodName := d.Get("name").(string)
	if methodName == "" {
		return logical.ErrorResponse("missing method name"), nil
	}

	// Acquire the mfa lock early. If the lock is not acquired here, config
	// read from the memdb below can get invalid before the config update hits
	// the storage and can cause inconsistencies.
	b.mfaLock.Lock()
	defer b.mfaLock.Unlock()

	mConfig, err := b.MemDBMFAConfigByMethodName(methodName, true)
	if err != nil {
		return nil, err
	}

	if mConfig == nil {
		configID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, errwrap.Wrapf("failed to generate an identifier for MFA config: {{err}}", err)
		}
		mConfig = &mfa.Config{
			Name: methodName,
			ID:   configID,
			Type: methodType,
		}
	}

	if mConfig.Type != methodType {
		return logical.ErrorResponse(fmt.Sprintf("method name %q is already in use under type %q", mConfig.Name, mConfig.Type)), nil
	}
	mConfig.Type = methodType

	accessorRaw, ok := d.GetOk("mount_accessor")
	if ok {
		accessor := accessorRaw.(string)
		validMount := b.Core.router.ValidateMountByAccessor(accessor)
		if validMount == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid mount accessor %q", accessor)), nil
		}
		mConfig.MountAccessor = accessor
	}

	usernameRaw, ok := d.GetOk("username_format")
	if ok {
		mConfig.UsernameFormat = usernameRaw.(string)
	}

	switch methodType {
	case mfaMethodTypeTOTP:
		err = parseTOTPConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	case mfaMethodTypeOkta:
		if mConfig.MountAccessor == "" {
			return logical.ErrorResponse(fmt.Sprintf("mfa type %q requires a %q parameter", methodType, "mount_accessor")), nil
		}

		err = parseOktaConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	case mfaMethodTypeDuo:
		if mConfig.MountAccessor == "" {
			return logical.ErrorResponse(fmt.Sprintf("mfa type %q requires a %q parameter", methodType, "mount_accessor")), nil
		}
		err = parseDuoConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	case mfaMethodTypePingID:
		if mConfig.MountAccessor == "" {
			return logical.ErrorResponse(fmt.Sprintf("mfa type %q requires a %q parameter", methodType, "mount_accessor")), nil
		}
		err = parsePingIDConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	default:
		return logical.ErrorResponse(fmt.Sprintf("unrecognized type %q", methodType)), nil
	}

	// Store the config
	err = b.putMFAConfig(ctx, mConfig)
	if err != nil {
		return nil, err
	}

	// Back the config in MemDB
	return nil, b.MemDBUpsertMFAConfig(ctx, mConfig)
}

func parseDuoConfig(mConfig *mfa.Config, d *framework.FieldData) error {
	secretKey := d.Get("secret_key").(string)
	if secretKey == "" {
		return fmt.Errorf("secret_key is empty")
	}

	integrationKey := d.Get("integration_key").(string)
	if integrationKey == "" {
		return fmt.Errorf("integration_key is empty")
	}

	apiHostname := d.Get("api_hostname").(string)
	if apiHostname == "" {
		return fmt.Errorf("api_hostname is empty")
	}

	config := &mfa.DuoConfig{
		SecretKey:      secretKey,
		IntegrationKey: integrationKey,
		APIHostname:    apiHostname,
		PushInfo:       d.Get("push_info").(string),
	}

	mConfig.Config = &mfa.Config_DuoConfig{
		DuoConfig: config,
	}

	return nil
}

func parsePingIDConfig(mConfig *mfa.Config, d *framework.FieldData) error {
	fileString := d.Get("settings_file_base64").(string)
	if fileString == "" {
		return fmt.Errorf("settings_file_base64 is empty")
	}

	fileBytes, err := base64.StdEncoding.DecodeString(fileString)
	if err != nil {
		return err
	}

	config := &mfa.PingIDConfig{}
	for _, line := range strings.Split(string(fileBytes), "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		splitLine := strings.SplitN(line, "=", 2)
		if len(splitLine) != 2 {
			return fmt.Errorf("pingid settings file contains a non-empty non-comment line that is not in key=value format: %q", line)
		}
		switch splitLine[0] {
		case "use_base64_key":
			config.UseBase64Key = splitLine[1]
		case "use_signature":
			result, err := parseutil.ParseBool(splitLine[1])
			if err != nil {
				return errors.New("error parsing use_signature value in pingid settings file")
			}
			config.UseSignature = result
		case "token":
			config.Token = splitLine[1]
		case "idp_url":
			config.IDPURL = splitLine[1]
		case "org_alias":
			config.OrgAlias = splitLine[1]
		case "admin_url":
			config.AdminURL = splitLine[1]
		case "authenticator_url":
			config.AuthenticatorURL = splitLine[1]
		default:
			return fmt.Errorf("unknown key %q in pingid settings file", splitLine[0])
		}
	}

	mConfig.Config = &mfa.Config_PingIDConfig{
		PingIDConfig: config,
	}

	return nil
}

func (b *SystemBackend) handleMFAConfigDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	methodName := d.Get("name").(string)
	if methodName == "" {
		return logical.ErrorResponse("missing method name"), nil
	}

	return nil, b.deleteMFAConfigByMethodName(ctx, methodName)
}

func (b *SystemBackend) handleMFAConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	methodName := d.Get("name").(string)
	if methodName == "" {
		return logical.ErrorResponse("missing method name"), nil
	}

	mConfig, err := b.MemDBMFAConfigByMethodName(methodName, false)
	if err != nil {
		return nil, err
	}
	if mConfig == nil {
		return nil, nil
	}

	respData := make(map[string]interface{})

	switch mConfig.Config.(type) {
	case *mfa.Config_TOTPConfig:
		totpConfig := mConfig.GetTOTPConfig()
		respData["issuer"] = totpConfig.Issuer
		respData["period"] = totpConfig.Period
		respData["digits"] = totpConfig.Digits
		respData["skew"] = totpConfig.Skew
		respData["key_size"] = totpConfig.KeySize
		respData["qr_size"] = totpConfig.QRSize
		respData["algorithm"] = otplib.Algorithm(totpConfig.Algorithm).String()
	case *mfa.Config_OktaConfig:
		oktaConfig := mConfig.GetOktaConfig()
		respData["org_name"] = oktaConfig.OrgName
		if oktaConfig.BaseURL != "" {
			respData["base_url"] = oktaConfig.BaseURL
		} else {
			respData["production"] = oktaConfig.Production
		}
		respData["mount_accessor"] = mConfig.MountAccessor
		respData["username_format"] = mConfig.UsernameFormat
	case *mfa.Config_DuoConfig:
		duoConfig := mConfig.GetDuoConfig()
		respData["api_hostname"] = duoConfig.APIHostname
		respData["pushinfo"] = duoConfig.PushInfo
		respData["mount_accessor"] = mConfig.MountAccessor
		respData["username_format"] = mConfig.UsernameFormat
	case *mfa.Config_PingIDConfig:
		pingConfig := mConfig.GetPingIDConfig()
		respData["use_signature"] = pingConfig.UseSignature
		respData["idp_url"] = pingConfig.IDPURL
		respData["org_alias"] = pingConfig.OrgAlias
		respData["admin_url"] = pingConfig.AdminURL
		respData["authenticator_url"] = pingConfig.AuthenticatorURL
	default:
		return nil, fmt.Errorf("invalid method type %q was persisted, underlying type: %T", mConfig.Type, mConfig.Config)
	}

	respData["type"] = mConfig.Type
	respData["id"] = mConfig.ID
	respData["name"] = mConfig.Name

	return &logical.Response{
		Data: respData,
	}, nil
}

func parseTOTPConfig(mConfig *mfa.Config, d *framework.FieldData) error {
	if mConfig == nil {
		return fmt.Errorf("config is nil")
	}

	if d == nil {
		return fmt.Errorf("field data is nil")
	}

	algorithm := d.Get("algorithm").(string)
	var keyAlgorithm otplib.Algorithm
	switch algorithm {
	case "SHA1":
		keyAlgorithm = otplib.AlgorithmSHA1
	case "SHA256":
		keyAlgorithm = otplib.AlgorithmSHA256
	case "SHA512":
		keyAlgorithm = otplib.AlgorithmSHA512
	default:
		return fmt.Errorf("unrecognized algorithm")
	}

	digits := d.Get("digits").(int)
	var keyDigits otplib.Digits
	switch digits {
	case 6:
		keyDigits = otplib.DigitsSix
	case 8:
		keyDigits = otplib.DigitsEight
	default:
		return fmt.Errorf("digits can only be 6 or 8")
	}

	period := d.Get("period").(int)
	if period <= 0 {
		return fmt.Errorf("period must be greater than zero")
	}

	skew := d.Get("skew").(int)
	switch skew {
	case 0:
	case 1:
	default:
		return fmt.Errorf("skew must be 0 or 1")
	}

	keySize := d.Get("key_size").(int)
	if keySize <= 0 {
		return fmt.Errorf("key_size must be greater than zero")
	}

	issuer := d.Get("issuer").(string)
	if issuer == "" {
		return fmt.Errorf("issuer must be set")
	}

	config := &mfa.TOTPConfig{
		Issuer:    issuer,
		Period:    uint32(period),
		Algorithm: int32(keyAlgorithm),
		Digits:    int32(keyDigits),
		Skew:      uint32(skew),
		KeySize:   uint32(keySize),
		QRSize:    int32(d.Get("qr_size").(int)),
	}
	mConfig.Config = &mfa.Config_TOTPConfig{
		TOTPConfig: config,
	}

	return nil
}

func parseOktaConfig(mConfig *mfa.Config, d *framework.FieldData) error {
	if mConfig == nil {
		return errors.New("config is nil")
	}

	if d == nil {
		return errors.New("field data is nil")
	}

	oktaConfig := &mfa.OktaConfig{}

	orgName := d.Get("org_name").(string)
	if orgName == "" {
		return errors.New("org_name must be set")
	}
	oktaConfig.OrgName = orgName

	token := d.Get("api_token").(string)
	if token == "" {
		return errors.New("api_token must be set")
	}
	oktaConfig.APIToken = token

	productionRaw, productionOk := d.GetOk("production")
	if productionOk {
		oktaConfig.Production = productionRaw.(bool)
	} else {
		oktaConfig.Production = true
	}

	baseURLRaw, ok := d.GetOk("base_url")
	if ok {
		oktaConfig.BaseURL = baseURLRaw.(string)
	} else {
		// Only set if not using legacy production flag
		if !productionOk {
			oktaConfig.BaseURL = "okta.com"
		}
	}

	primaryEmailOnly := d.Get("primary_email").(bool)
	if primaryEmailOnly {
		oktaConfig.PrimaryEmail = true
	}

	_, err := url.Parse(fmt.Sprintf("https://%s,%s", oktaConfig.OrgName, oktaConfig.BaseURL))
	if err != nil {
		return errwrap.Wrapf("error parsing given base_url: {{err}}", err)
	}

	mConfig.Config = &mfa.Config_OktaConfig{
		OktaConfig: oktaConfig,
	}

	return nil
}

func (c *Core) validateMFA(ctx context.Context, methodName string, entity *identity.Entity, req *logical.Request) (retErr error) {
	// If the calling token is not tied to an entity, and if the path specified
	// a MFA requirement, then we fail the request
	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	defer func() {
		if retErr != nil {
			c.systemBackend.mfaLogger.Error("validation failed", "method", methodName, "error", retErr)
		} else {
			c.systemBackend.mfaLogger.Debug("validation successful", "method", methodName, "path", req.Path, "entity_id", entity.ID)
		}
	}()

	// Get the configuration for the MFA method set in system backend
	mConfig, err := c.systemBackend.MemDBMFAConfigByMethodName(methodName, false)
	if err != nil {
		return fmt.Errorf("failed to read MFA configuration")
	}

	if mConfig == nil {
		return fmt.Errorf("MFA method configuration not present")
	}

	var alias *identity.Alias
	var finalUsername string
	switch mConfig.Type {
	case mfaMethodTypeDuo, mfaMethodTypeOkta, mfaMethodTypePingID:
		for _, entry := range entity.Aliases {
			if mConfig.MountAccessor == entry.MountAccessor {
				alias = entry
				break
			}
		}
		if alias == nil {
			return fmt.Errorf("could not find alias in entity matching the MFA's mount accessor")
		}
		finalUsername = formatUsername(mConfig.UsernameFormat, alias, entity)
	}

	switch mConfig.Type {
	case mfaMethodTypeTOTP:
		// Get the MFA secret data required to validate the supplied credentials
		if entity.MFASecrets == nil {
			return fmt.Errorf("MFA secret for method name %q not present in entity %q", mConfig.Name, entity.ID)
		}
		entityMFASecret := entity.MFASecrets[mConfig.ID]
		if entityMFASecret == nil {
			return fmt.Errorf("MFA secret for method name %q not present in entity %q", mConfig.Name, entity.ID)
		}

		// Extract the MFA credentials supplied via the request headers
		headerCreds := req.MFACreds[mConfig.Name]
		if headerCreds == nil {
			return fmt.Errorf("MFA credentials not supplied")
		}

		return c.validateTOTP(ctx, headerCreds, entityMFASecret, mConfig.ID, entity.ID)

	case mfaMethodTypeOkta:
		return c.validateOkta(ctx, mConfig, finalUsername)

	case mfaMethodTypeDuo:
		return c.validateDuo(ctx, req.MFACreds[mConfig.Name], mConfig, finalUsername, req)

	case mfaMethodTypePingID:
		return c.validatePingID(ctx, mConfig, finalUsername)

	default:
		return fmt.Errorf("unrecognized MFA type %q", mConfig.Type)
	}
}

func formatUsername(format string, alias *identity.Alias, entity *identity.Entity) string {
	if format == "" {
		return alias.Name
	}

	username := format
	username = strings.Replace(username, "{{alias.name}}", alias.Name, -1)
	username = strings.Replace(username, "{{entity.name}}", entity.Name, -1)
	for k, v := range alias.Metadata {
		username = strings.Replace(username, fmt.Sprintf("{{alias.metadata.%s}}", k), v, -1)
	}
	for k, v := range entity.Metadata {
		username = strings.Replace(username, fmt.Sprintf("{{entity.metadata.%s}}", k), v, -1)
	}
	return username
}

func (c *Core) validateDuo(ctx context.Context, creds []string, mConfig *mfa.Config, username string, req *logical.Request) error {
	duoConfig := mConfig.GetDuoConfig()
	if duoConfig == nil {
		return fmt.Errorf("failed to get Duo configuration for method %q", mConfig.Name)
	}

	passcode := ""
	for _, cred := range creds {
		if strings.HasPrefix(cred, "passcode") {
			splits := strings.SplitN(cred, "=", 2)
			if len(splits) != 2 {
				return fmt.Errorf("invalid credential %q", cred)
			}
			if splits[0] == "passcode" {
				passcode = splits[1]
			}
		}
	}

	client := duoapi.NewDuoApi(
		duoConfig.IntegrationKey,
		duoConfig.SecretKey,
		duoConfig.APIHostname,
		duoConfig.PushInfo,
	)

	authClient := authapi.NewAuthApi(*client)
	check, err := authClient.Check()
	if err != nil {
		return err
	}
	if check == nil {
		return errors.New("Duo api check returned nil, possibly bad integration key")
	}
	var message string
	var messageDetail string
	if check.StatResult.Message != nil {
		message = *check.StatResult.Message
	}
	if check.StatResult.Message_Detail != nil {
		messageDetail = *check.StatResult.Message_Detail
	}
	if check.StatResult.Stat != "OK" {
		return fmt.Errorf("check against Duo failed; message (if given): %q; message detail (if given): %q", message, messageDetail)
	}

	preauth, err := authClient.Preauth(authapi.PreauthUsername(username), authapi.PreauthIpAddr(req.Connection.RemoteAddr))
	if err != nil {
		return errwrap.Wrapf("failed to perform Duo preauth: {{err}}", err)
	}
	if preauth == nil {
		return fmt.Errorf("failed to perform Duo preauth")
	}
	if preauth.StatResult.Stat != "OK" {
		return fmt.Errorf("failed to perform Duo preauth: %q - %q", *preauth.StatResult.Message, *preauth.StatResult.Message_Detail)
	}

	switch preauth.Response.Result {
	case "allow":
		return nil
	case "deny":
		return fmt.Errorf(preauth.Response.Status_Msg)
	case "enroll":
		return fmt.Errorf(fmt.Sprintf("%q - %q", preauth.Response.Status_Msg, preauth.Response.Enroll_Portal_Url))
	case "auth":
		break
	default:
		return fmt.Errorf("invalid response from Duo preauth: %q", preauth.Response.Result)
	}

	options := []func(*url.Values){}
	factor := "push"
	if passcode != "" {
		factor = "passcode"
		options = append(options, authapi.AuthPasscode(passcode))
	} else {
		options = append(options, authapi.AuthDevice("auto"))
		if duoConfig.PushInfo != "" {
			options = append(options, authapi.AuthPushinfo(duoConfig.PushInfo))
		}
	}

	options = append(options, authapi.AuthUsername(username))
	options = append(options, authapi.AuthAsync())

	result, err := authClient.Auth(factor, options...)
	if err != nil {
		return errwrap.Wrapf("failed to authenticate with Duo: {{err}}", err)
	}
	if result.StatResult.Stat != "OK" {
		return fmt.Errorf("failed to authenticate with Duo: %q - %q", *result.StatResult.Message, *result.StatResult.Message_Detail)
	}
	if result.Response.Txid == "" {
		return fmt.Errorf("failed to get transaction ID for Duo authentication")
	}

	for {
		// AuthStatus does the long polling until there is a status update. So
		// there is no need to wait for a second before we invoke this API.
		statusResult, err := authClient.AuthStatus(result.Response.Txid)
		if err != nil {
			return errwrap.Wrapf("failed to get authentication status from Duo: {{err}}", err)
		}
		if statusResult == nil {
			return errwrap.Wrapf("failed to get authentication status from Duo: {{err}}", err)
		}
		if statusResult.StatResult.Stat != "OK" {
			return fmt.Errorf("failed to get authentication status from Duo: %q - %q", *statusResult.StatResult.Message, *statusResult.StatResult.Message_Detail)
		}

		switch statusResult.Response.Result {
		case "deny":
			return fmt.Errorf("duo authentication failed: %q", statusResult.Response.Status_Msg)
		case "allow":
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("duo push verification operation canceled")
		case <-time.After(time.Second):
		}
	}
}

func (c *Core) validateOkta(ctx context.Context, mConfig *mfa.Config, username string) error {
	oktaConfig := mConfig.GetOktaConfig()
	if oktaConfig == nil {
		return fmt.Errorf("failed to get Okta configuration for method %q", mConfig.Name)
	}

	var client *okta.Client
	if oktaConfig.BaseURL != "" {
		var err error
		client, err = okta.NewClientWithDomain(cleanhttp.DefaultClient(), oktaConfig.OrgName, oktaConfig.BaseURL, oktaConfig.APIToken)
		if err != nil {
			return errwrap.Wrapf("error getting Okta client: {{err}}", err)
		}
	} else {
		client = okta.NewClient(cleanhttp.DefaultClient(), oktaConfig.OrgName, oktaConfig.APIToken, oktaConfig.Production)
	}

	var filterOpts *okta.UserListFilterOptions
	if oktaConfig.PrimaryEmail {
		filterOpts = &okta.UserListFilterOptions{
			EmailEqualTo: username,
		}
	} else {
		filterOpts = &okta.UserListFilterOptions{
			LoginEqualTo: username,
		}
	}

	users, _, err := client.Users.ListWithFilter(filterOpts)
	if err != nil {
		return err
	}
	switch {
	case len(users) == 0:
		return fmt.Errorf("no users found for e-mail address")
	case len(users) > 1:
		return fmt.Errorf("more than one user found for e-mail address")
	}

	user := &users[0]

	_, err = client.Users.PopulateMFAFactors(user)
	if err != nil {
		return err
	}

	if len(user.MFAFactors) == 0 {
		return fmt.Errorf("no MFA factors found for user")
	}

	var factorID string
	for _, factor := range user.MFAFactors {
		if factor.FactorType == "push" {
			factorID = factor.ID
			break
		}
	}

	if factorID == "" {
		return fmt.Errorf("no push-type MFA factor found for user")
	}

	type pollInfo struct {
		ValidationURL string `json:"href"`
	}

	type pushLinks struct {
		Poll pollInfo `json:"poll"`
	}

	type pushResult struct {
		Expiration   time.Time `json:"expiresAt"`
		FactorResult string    `json:"factorResult"`
		Links        pushLinks `json:"_links"`
	}

	req, err := client.NewRequest("POST", fmt.Sprintf("users/%s/factors/%s/verify", user.ID, factorID), nil)
	if err != nil {
		return err
	}

	var result pushResult
	_, err = client.Do(req, &result)
	if err != nil {
		return err
	}

	if result.FactorResult != "WAITING" {
		return fmt.Errorf("expected WAITING status for push status, got %q", result.FactorResult)
	}

	for {
		req, err := client.NewRequest("GET", result.Links.Poll.ValidationURL, nil)
		if err != nil {
			return err
		}
		var result pushResult
		_, err = client.Do(req, &result)
		if err != nil {
			return err
		}
		switch result.FactorResult {
		case "WAITING":
		case "SUCCESS":
			return nil
		case "REJECTED":
			return fmt.Errorf("push verification explicitly rejected")
		case "TIMEOUT":
			return fmt.Errorf("push verification timed out")
		default:
			return fmt.Errorf("unknown status code")
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("push verification operation canceled")
		case <-time.After(time.Second):
		}
	}
}

func (c *Core) validatePingID(ctx context.Context, mConfig *mfa.Config, username string) error {
	pingConfig := mConfig.GetPingIDConfig()
	if pingConfig == nil {
		return fmt.Errorf("failed to get PingID configuration for method %q", mConfig.Name)
	}

	signingKey, err := base64.StdEncoding.DecodeString(pingConfig.UseBase64Key)
	if err != nil {
		return errwrap.Wrapf("failed decoding pingid signing key: {{err}}", err)
	}

	client := cleanhttp.DefaultClient()

	createRequest := func(reqPath string, reqBody map[string]interface{}) (*http.Request, error) {
		// Construct the token
		token := &jwt.Token{
			Method: jwt.SigningMethodHS256,
			Header: map[string]interface{}{
				"alg":       "HS256",
				"org_alias": pingConfig.OrgAlias,
				"token":     pingConfig.Token,
			},
			Claims: jwt.MapClaims{
				"reqHeader": map[string]interface{}{
					"locale":    "en",
					"orgAlias":  pingConfig.OrgAlias,
					"secretKey": pingConfig.Token,
					"timestamp": time.Now().Format("2006-01-02  15:04:05.000"),
					"version":   "4.9",
				},
				"reqBody": reqBody,
			},
		}
		signedToken, err := token.SignedString(signingKey)
		if err != nil {
			return nil, errwrap.Wrapf("failed signing pingid request token: {{err}}", err)
		}

		// Construct the URL
		if !strings.HasPrefix(reqPath, "/") {
			reqPath = "/" + reqPath
		}
		reqURL, err := url.Parse(pingConfig.IDPURL + reqPath)
		if err != nil {
			return nil, errwrap.Wrapf("failed to parse pingid request url: {{err}}", err)
		}

		// Construct the request; WithContext is done here since it's a shallow
		// copy
		req := &http.Request{}
		req = req.WithContext(ctx)
		req.Method = "POST"
		req.URL = reqURL
		req.Body = ioutil.NopCloser(bytes.NewBufferString(signedToken))
		if req.Header == nil {
			req.Header = make(http.Header)
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}

	do := func(req *http.Request) (*jwt.Token, error) {
		// Run the request and fetch the response
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return nil, fmt.Errorf("nil response from pingid")
		}
		if resp.Body == nil {
			return nil, fmt.Errorf("nil body in pingid response")
		}
		bodyBytes := bytes.NewBuffer(nil)
		_, err = bodyBytes.ReadFrom(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, errwrap.Wrapf("error reading pingid response: {{err}}", err)
		}

		// Parse the body, which is a JWT. Ensure that it's using HMAC signing
		// and return the signing key in the func for validation
		token, err := jwt.Parse(bodyBytes.String(), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method %q from pingid response", token.Header["alg"])
			}
			return signingKey, nil
		})
		if err != nil {
			return nil, errwrap.Wrapf("error parsing pingid response: {{err}}", err)
		}

		// Check if parameters are as expected
		if _, ok := token.Header["token"]; !ok {
			return nil, fmt.Errorf("%q header not found", "token")
		}
		if headerTokenStr, ok := token.Header["token"].(string); !ok || headerTokenStr != pingConfig.Token {
			return nil, fmt.Errorf("invalid token in ping response")
		}

		// validate org alias
		// This was originally 'org_alias', but it appears to now be returned as 'orgAlias'. Official
		// Ping docs are sparse on the header details. We now prefer orgAlias but will still handle
		// org_alias.
		oa := token.Header["orgAlias"]
		if oa == nil {
			if oa = token.Header["org_alias"]; oa == nil {
				return nil, fmt.Errorf("neither orgAlias nor org_alias headers were found")
			}
		}

		if headerOrgAliasStr, ok := oa.(string); !ok || headerOrgAliasStr != pingConfig.OrgAlias {
			return nil, fmt.Errorf("invalid org_alias in ping response")
		}
		return token, nil
	}

	type deviceDetails struct {
		PushEnabled bool  `mapstructure:"pushEnabled"`
		DeviceID    int64 `mapstructure:"deviceId"`
	}

	type respBody struct {
		SessionID   string          `mapstructure:"sessionId"`
		ErrorID     int64           `mapstructure:"errorId"`
		ErrorMsg    string          `mapstructure:"errorMsg"`
		UserDevices []deviceDetails `mapstructure:"userDevices"`
	}

	type apiResponse struct {
		ResponseBody respBody `mapstructure:"responseBody"`
	}

	/*
		// Normally we don't leave in commented code, however:
		// Explicitly setting the device ID didn't work even when the device was
		// push enabled (said the application was not installed on the device), so
		// instead trigger default behavior, which does work, even when there's
		// only one device and the deviceid matched :-/
		// We're leaving this here because if we support other types we'll likely
		// still need it, and if we get device ID selection working we'll want it.

		req, err := createRequest("rest/4/startauthentication/do", map[string]interface{}{
			"spAlias":  "web",
			"userName": username,
		})
		if err != nil {
			return err
		}
		token, err := do(req)
		if err != nil {
			return err
		}

		// We get back a map from the JWT library so use mapstructure
		var startResp apiResponse
		err = mapstructure.Decode(token.Claims, &startResp)
		if err != nil {
			return err
		}

		// Look for at least one push-enabled method
		body := startResp.ResponseBody
		var foundPush bool
		switch {
		case body.ErrorID != 30007:
			return fmt.Errorf("only pingid push authentication is currently supported")

		case len(body.UserDevices) == 0:
			return fmt.Errorf("no user mfa devices returned from pingid")

		default:
			for _, dev := range body.UserDevices {
				if dev.PushEnabled {
					foundPush = true
					break
				}
			}

			if !foundPush {
				return fmt.Errorf("no push enabled device id found from pingid")
			}
		}
	*/
	req, err := createRequest("rest/4/authonline/do", map[string]interface{}{
		"spAlias":  "web",
		"userName": username,
		"authType": "CONFIRM",
	})
	if err != nil {
		return err
	}
	token, err := do(req)
	if err != nil {
		return err
	}

	// Ensure a success response
	var authResp apiResponse
	err = mapstructure.Decode(token.Claims, &authResp)
	if err != nil {
		return err
	}

	if authResp.ResponseBody.ErrorID != 200 {
		return errors.New(authResp.ResponseBody.ErrorMsg)
	}

	return nil
}

func (c *Core) validateTOTP(ctx context.Context, creds []string, entityMethodSecret *mfa.Secret, configID, entityID string) error {
	if len(creds) == 0 {
		return fmt.Errorf("missing TOTP passcode")
	}

	if len(creds) > 1 {
		return fmt.Errorf("more than one TOTP passcode supplied")
	}

	totpSecret := entityMethodSecret.GetTOTPSecret()
	if totpSecret == nil {
		return fmt.Errorf("entity does not contain the TOTP secret")
	}

	key, err := c.fetchTOTPKey(ctx, configID, entityID)
	if err != nil {
		return errwrap.Wrapf("error fetching TOTP key: {{err}}", err)
	}

	if key == "" {
		return fmt.Errorf("empty key for entity's TOTP secret")
	}

	validateOpts := totplib.ValidateOpts{
		Period:    uint(totpSecret.Period),
		Skew:      uint(totpSecret.Skew),
		Digits:    otplib.Digits(int(totpSecret.Digits)),
		Algorithm: otplib.Algorithm(int(totpSecret.Algorithm)),
	}

	valid, err := totplib.ValidateCustom(creds[0], key, time.Now(), validateOpts)
	if err != nil && err != otplib.ErrValidateInputInvalidLength {
		return errwrap.Wrapf("failed to validate TOTP passcode: {{err}}", err)
	}

	if !valid {
		return fmt.Errorf("failed to validate TOTP passcode")
	}

	return nil
}

func mfaConfigTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: memDBMFAConfigsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"name": {
				Name:   "name",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

func (b *SystemBackend) loadMFAConfigs(ctx context.Context) error {
	// Accumulate existing entities
	b.mfaLogger.Debug("loading configurations")
	existing, err := b.Core.systemBarrierView.List(ctx, mfaConfigPrefix)
	if err != nil {
		return errwrap.Wrapf("failed to list MFA configurations: {{err}}", err)
	}
	b.mfaLogger.Debug("configurations collected", "num_existing", len(existing))

	for _, key := range existing {
		if b.logger.IsTrace() {
			b.mfaLogger.Trace("loading configuration", "method", key)
		}

		// Read the config from storage
		mConfig, err := b.getMFAConfig(ctx, mfaConfigPrefix+key)
		if err != nil {
			return err
		}

		if mConfig == nil {
			continue
		}

		// Load the config in MemDB
		err = b.MemDBUpsertMFAConfig(ctx, mConfig)
		if err != nil {
			return errwrap.Wrapf("failed to load configuration in MemDB: {{err}}", err)
		}
	}

	b.mfaLogger.Info("configurations restored")

	return nil
}

func (b *SystemBackend) MemDBUpsertMFAConfig(ctx context.Context, mConfig *mfa.Config) error {
	txn := b.db.Txn(true)
	defer txn.Abort()

	err := b.MemDBUpsertMFAConfigInTxn(txn, mConfig)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (b *SystemBackend) MemDBUpsertMFAConfigInTxn(txn *memdb.Txn, mConfig *mfa.Config) error {
	if txn == nil {
		return fmt.Errorf("nil txn")
	}

	if mConfig == nil {
		return fmt.Errorf("config is nil")
	}

	mConfigRaw, err := txn.First(memDBMFAConfigsTable, "id", mConfig.ID)
	if err != nil {
		return errwrap.Wrapf("failed to lookup MFA config from MemDB using id: {{err}}", err)
	}

	if mConfigRaw != nil {
		err = txn.Delete(memDBMFAConfigsTable, mConfigRaw)
		if err != nil {
			return errwrap.Wrapf("failed to delete MFA config from MemDB: {{err}}", err)
		}
	}

	if err := txn.Insert(memDBMFAConfigsTable, mConfig); err != nil {
		return errwrap.Wrapf("failed to update MFA config into MemDB: {{err}}", err)
	}

	return nil
}

func (b *SystemBackend) MemDBMFAConfigByMethodNameInTxn(txn *memdb.Txn, methodName string, clone bool) (*mfa.Config, error) {
	if methodName == "" {
		return nil, fmt.Errorf("missing method name")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	mConfigRaw, err := txn.First(memDBMFAConfigsTable, "name", methodName)
	if err != nil {
		return nil, errwrap.Wrapf("failed to fetch MFA config from memdb using method name: {{err}}", err)
	}

	if mConfigRaw == nil {
		return nil, nil
	}

	mConfig, ok := mConfigRaw.(*mfa.Config)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched MFA config")
	}

	if clone {
		return mConfig.Clone()
	}

	return mConfig, nil
}

func (b *SystemBackend) MemDBMFAConfigByMethodName(methodName string, clone bool) (*mfa.Config, error) {
	if methodName == "" {
		return nil, fmt.Errorf("missing method name")
	}

	txn := b.db.Txn(false)

	return b.MemDBMFAConfigByMethodNameInTxn(txn, methodName, clone)
}

func (b *SystemBackend) MemDBMFAConfigByIDInTxn(txn *memdb.Txn, mConfigID string, clone bool) (*mfa.Config, error) {
	if mConfigID == "" {
		return nil, fmt.Errorf("missing config id")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	mConfigRaw, err := txn.First(memDBMFAConfigsTable, "id", mConfigID)
	if err != nil {
		return nil, errwrap.Wrapf("failed to fetch MFA config from memdb using id: {{err}}", err)
	}

	if mConfigRaw == nil {
		return nil, nil
	}

	mConfig, ok := mConfigRaw.(*mfa.Config)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched MFA config")
	}

	if clone {
		return mConfig.Clone()
	}

	return mConfig, nil
}

// TODO: should I remove this? seems not used anywhere
func (b *SystemBackend) MemDBMFAConfigByID(mConfigID string, clone bool) (*mfa.Config, error) {
	if mConfigID == "" {
		return nil, fmt.Errorf("missing config id")
	}

	txn := b.db.Txn(false)

	return b.MemDBMFAConfigByIDInTxn(txn, mConfigID, clone)
}

func (b *SystemBackend) deleteMFAConfigByMethodName(ctx context.Context, methodName string) error {
	var err error

	if methodName == "" {
		return fmt.Errorf("missing method name")
	}

	b.mfaLock.Lock()
	defer b.mfaLock.Unlock()

	// Delete the config from storage
	entryIndex := mfaConfigPrefix + methodName
	err = b.Core.systemBarrierView.Delete(ctx, entryIndex)
	if err != nil {
		return err
	}

	// Create a MemDB transaction to delete config
	txn := b.db.Txn(true)
	defer txn.Abort()

	mConfig, err := b.MemDBMFAConfigByMethodNameInTxn(txn, methodName, false)
	if err != nil {
		return err
	}

	if mConfig == nil {
		return nil
	}

	if mConfig.Type == "totp" && mConfig.ID != "" {
		// This is best effort; if they end up hanging around it's okay, they're encrypted anyways
		if err := logical.ClearView(ctx, NewBarrierView(b.Core.barrier, fmt.Sprintf("%s%s", mfaTOTPKeysPrefix, mConfig.ID))); err != nil {
			b.mfaLogger.Warn("unable to clear TOTP keys", "method", mConfig.Name, "error", err)
		}
	}

	// Delete the config from MemDB
	err = b.MemDBDeleteMFAConfigByMethodNameInTxn(txn, methodName)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (b *SystemBackend) MemDBDeleteMFAConfigByMethodNameInTxn(txn *memdb.Txn, methodName string) error {
	if methodName == "" {
		return nil
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	mConfig, err := b.MemDBMFAConfigByMethodNameInTxn(txn, methodName, false)
	if err != nil {
		return err
	}

	if mConfig == nil {
		return nil
	}

	err = txn.Delete(memDBMFAConfigsTable, mConfig)
	if err != nil {
		return errwrap.Wrapf("failed to delete MFA config from memdb: {{err}}", err)
	}

	return nil
}

// TODO: should this be remove as it is not used anywhere?
func (b *SystemBackend) MemDBDeleteMFAConfigByIDInTxn(txn *memdb.Txn, configID string) error {
	if configID == "" {
		return nil
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	mConfig, err := b.MemDBMFAConfigByIDInTxn(txn, configID, false)
	if err != nil {
		return err
	}

	if mConfig == nil {
		return nil
	}

	err = txn.Delete(memDBMFAConfigsTable, mConfig)
	if err != nil {
		return errwrap.Wrapf("failed to delete MFA config from memdb: {{err}}", err)
	}

	return nil
}

func (b *SystemBackend) putMFAConfig(ctx context.Context, mConfig *mfa.Config) error {
	entryIndex := mfaConfigPrefix + mConfig.Name

	marshaledEntry, err := proto.Marshal(mConfig)
	if err != nil {
		return err
	}

	return b.Core.systemBarrierView.Put(ctx, &logical.StorageEntry{
		Key:   entryIndex,
		Value: marshaledEntry,
	})
}

func (b *SystemBackend) getMFAConfig(ctx context.Context, key string) (*mfa.Config, error) {
	if !strings.HasPrefix(key, mfaConfigPrefix) {
		return nil, nil
	}

	entry, err := b.Core.systemBarrierView.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var mConfig mfa.Config
	err = proto.Unmarshal(entry.Value, &mConfig)
	if err != nil {
		return nil, err
	}

	return &mConfig, nil
}
