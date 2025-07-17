// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	duoapi "github.com/duosecurity/duo_api_golang"
	"github.com/duosecurity/duo_api_golang/authapi"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/quotas"
	"github.com/mitchellh/mapstructure"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"github.com/patrickmn/go-cache"
	otplib "github.com/pquerna/otp"
	totplib "github.com/pquerna/otp/totp"
)

const (
	mfaMethodTypeTOTP              = "totp"
	mfaMethodTypeDuo               = "duo"
	mfaMethodTypeOkta              = "okta"
	mfaMethodTypePingID            = "pingid"
	memDBLoginMFAConfigsTable      = "login_mfa_configs"
	memDBMFALoginEnforcementsTable = "login_enforcements"
	mfaTOTPKeysPrefix              = systemBarrierPrefix + "mfa/totpkeys/"

	// loginMFAConfigPrefix is the storage prefix for persisting login MFA method
	// configs
	loginMFAConfigPrefix      = "login-mfa/method/"
	mfaLoginEnforcementPrefix = "login-mfa/enforcement/"
)

type totpKey struct {
	Key string `json:"key"`
}

// loginMfaPaths returns the API endpoints to configure the new style
// login MFA. The following paths are supported:
// mfa/method/:mfa_method - management of MFA method IDs, which can be used for configuration
// mfa/login_enforcement/:config_name - configures single or two phase MFA auth
func (b *SystemBackend) loginMFAPaths() []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "mfa/validate",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "mfa",
				OperationVerb:   "validate",
			},

			Fields: map[string]*framework.FieldSchema{
				"mfa_request_id": {
					Type:        framework.TypeString,
					Description: "ID for this MFA request",
					Required:    true,
				},
				"mfa_payload": {
					Type:        framework.TypeMap,
					Description: "A map from MFA method ID to a slice of passcodes or an empty slice if the method does not use passcodes",
					Required:    true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.Core.loginMFABackend.handleMFALoginValidate,
					Responses: map[int][]framework.Response{
						http.StatusOK: {{
							Description: "OK",
						}},
					},
					Summary:                   "Validates the login for the given MFA methods. Upon successful validation, it returns an auth response containing the client token",
					ForwardPerformanceStandby: true,
				},
			},
		},
	}
}

// uuidRegex crafts a regex for use in URL paths, somewhat similar to framework.GenericNameRegex, but only accepting
// UUIDs, and only lowercase UUIDs at that.
// It is currently exclusively used for the method_id parameter for MFA methods.
// Think twice before making use of it in any other context, as restricting the valid input in the URL regex results in
// an "unsupported path" error, given input which does not match the regex, which is a fairly unclear way to report an
// invalid parameter value, unless the person seeing the error has an excellent understanding of Vault URL routing.
func uuidRegex(name string) string {
	return fmt.Sprintf("(?P<%s>[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})", name)
}

type MFABackend struct {
	Core        *Core
	mfaLock     *sync.RWMutex
	db          *memdb.MemDB
	mfaLogger   hclog.Logger
	namespacer  Namespacer
	methodTable string
	usedCodes   *cache.Cache
}

type LoginMFABackend struct {
	*MFABackend
}

func loginMFASchemaFuncs() []func() *memdb.TableSchema {
	return []func() *memdb.TableSchema{
		loginMFAConfigTableSchema,
		loginEnforcementTableSchema,
	}
}

func NewLoginMFABackend(core *Core, logger hclog.Logger) *LoginMFABackend {
	b := NewMFABackend(core, logger, memDBLoginMFAConfigsTable, loginMFASchemaFuncs())
	return &LoginMFABackend{b}
}

func NewMFABackend(core *Core, logger hclog.Logger, prefix string, schemaFuncs []func() *memdb.TableSchema) *MFABackend {
	db, _ := SetupMFAMemDB(schemaFuncs)
	return &MFABackend{
		Core:        core,
		mfaLock:     &sync.RWMutex{},
		db:          db,
		mfaLogger:   logger.Named("mfa"),
		namespacer:  core,
		methodTable: prefix,
	}
}

func SetupMFAMemDB(schemaFuncs []func() *memdb.TableSchema) (*memdb.MemDB, error) {
	mfaSchemas := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema),
	}

	for _, schemaFunc := range schemaFuncs {
		schema := schemaFunc()
		if _, ok := mfaSchemas.Tables[schema.Name]; ok {
			panic(fmt.Sprintf("duplicate table name: %s", schema.Name))
		}
		mfaSchemas.Tables[schema.Name] = schema
	}

	db, err := memdb.NewMemDB(mfaSchemas)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (b *LoginMFABackend) ResetLoginMFAMemDB() error {
	var err error

	db, err := SetupMFAMemDB(loginMFASchemaFuncs())
	if err != nil {
		return err
	}

	b.db = db

	return nil
}

func (i *IdentityStore) handleMFAMethodListGlobal(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keys, configInfo, err := i.mfaBackend.mfaMethodList(ctx, "")
	if err != nil {
		return nil, err
	}

	return logical.ListResponseWithInfo(keys, configInfo), nil
}

func (i *IdentityStore) handleMFAMethodListCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, methodType string) (*logical.Response, error) {
	keys, configInfo, err := i.mfaBackend.mfaMethodList(ctx, methodType)
	if err != nil {
		return nil, err
	}

	return logical.ListResponseWithInfo(keys, configInfo), nil
}

func (i *IdentityStore) handleMFAMethodReadGlobal(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return i.handleMFAMethodReadCommon(ctx, req, d, "")
}

func (i *IdentityStore) handleMFAMethodReadCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, methodType string) (*logical.Response, error) {
	methodID := d.Get("method_id").(string)
	if methodID == "" {
		return logical.ErrorResponse("missing method ID"), nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	respData, err := i.mfaBackend.mfaConfigReadByMethodID(methodID)
	if err != nil {
		return nil, err
	}

	if respData == nil {
		return nil, nil
	}

	mfaNs, err := i.namespacer.NamespaceByID(ctx, respData["namespace_id"].(string))
	if err != nil {
		return nil, err
	}

	// reading the method config either from the same namespace or from the parent or from the child should all work
	if !(ns.ID == mfaNs.ID || mfaNs.HasParent(ns) || ns.HasParent(mfaNs)) {
		return logical.ErrorResponse("request namespace does not match method namespace"), logical.ErrPermissionDenied
	}

	if methodType != "" && respData["type"] != methodType {
		return logical.ErrorResponse("failed to find the method ID under MFA type %s.", methodType), nil
	}

	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) handleMFAMethodWriteCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, methodType string) (*logical.Response, error) {
	var err error
	var mConfig *mfa.Config
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// This method handles both the create and update flow. The method_id field is not part of the field schema for the
	// create flow, so we need to use GetOk in order to not panic.
	var methodID string
	methodIDAsInterface, ok := d.GetOk("method_id")
	if ok {
		methodID = methodIDAsInterface.(string)
	}

	methodName := d.Get("method_name").(string)

	b := i.mfaBackend
	b.mfaLock.Lock()
	defer b.mfaLock.Unlock()

	if methodID != "" {
		mConfig, err = b.MemDBMFAConfigByID(methodID)
		if err != nil {
			return nil, err
		}

		// If methodID is specified, but we didn't find anything, return a 404
		if mConfig == nil {
			return nil, nil
		}
	}

	// check if an MFA method configuration exists with that method name
	if methodName != "" {
		namedMfaConfig, err := b.MemDBMFAConfigByName(ctx, methodName)
		if err != nil {
			return nil, err
		}
		if namedMfaConfig != nil {
			if mConfig == nil {
				mConfig = namedMfaConfig
			} else {
				if mConfig.ID != namedMfaConfig.ID {
					return nil, fmt.Errorf("a login MFA method configuration with the method name %s already exists", methodName)
				}
			}
		}
	}

	if mConfig == nil {
		configID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate an identifier for MFA config: %v", err)
		}
		mConfig = &mfa.Config{
			ID:          configID,
			Type:        methodType,
			NamespaceID: ns.ID,
		}
	}

	// Updating the method config name
	if methodName != "" {
		mConfig.Name = methodName
	}

	mfaNs, err := i.namespacer.NamespaceByID(ctx, mConfig.NamespaceID)
	if err != nil {
		return nil, err
	}

	// this logic assumes that the config namespace and the current
	// namespace should be the same. Note an ancestor of mfaNs is not allowed
	// to create/update methodID
	if ns.ID != mfaNs.ID {
		return logical.ErrorResponse("request namespace does not match method namespace"), nil
	}

	mConfig.Type = methodType
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
		err = parseOktaConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	case mfaMethodTypeDuo:
		err = parseDuoConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	case mfaMethodTypePingID:
		err = parsePingIDConfig(mConfig, d)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	default:
		return logical.ErrorResponse(fmt.Sprintf("unrecognized type %q", methodType)), nil
	}

	// Store the config
	err = b.putMFAConfigByID(ctx, mConfig)
	if err != nil {
		return nil, err
	}

	// Back the config in MemDB
	err = b.MemDBUpsertMFAConfig(ctx, mConfig)
	if err != nil {
		return nil, err
	}

	if methodID == "" {
		return &logical.Response{
			Data: map[string]interface{}{
				"method_id": mConfig.ID,
			},
		}, nil
	} else {
		return nil, nil
	}
}

func (i *IdentityStore) handleMFAMethodDeleteCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, methodType string) (*logical.Response, error) {
	methodID := d.Get("method_id").(string)
	if methodID == "" {
		return logical.ErrorResponse("missing method ID"), nil
	}
	return nil, i.mfaBackend.deleteMFAConfigByMethodID(ctx, methodID, methodType, memDBLoginMFAConfigsTable, loginMFAConfigPrefix)
}

func (i *IdentityStore) handleLoginMFAGenerateUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return i.handleLoginMFAGenerateCommon(ctx, req, d.Get("method_id").(string), req.EntityID)
}

func (i *IdentityStore) handleLoginMFAAdminGenerateUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return i.handleLoginMFAGenerateCommon(ctx, req, d.Get("method_id").(string), d.Get("entity_id").(string))
}

func (i *IdentityStore) handleLoginMFAGenerateCommon(ctx context.Context, req *logical.Request, methodID, entityID string) (*logical.Response, error) {
	if methodID == "" {
		return logical.ErrorResponse("missing method ID"), nil
	}

	if entityID == "" {
		return logical.ErrorResponse("missing entityID"), nil
	}

	mConfig, err := i.mfaBackend.MemDBMFAConfigByID(methodID)
	if err != nil {
		return nil, err
	}
	if mConfig == nil {
		return logical.ErrorResponse(fmt.Sprintf("configuration for method ID %q does not exist", methodID)), nil
	}
	if mConfig.ID == "" {
		return nil, fmt.Errorf("configuration for method ID %q does not contain an identifier", methodID)
	}

	entity, err := i.MemDBEntityByID(entityID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to find entity with ID %q: error: %w", entityID, err)
	}

	if entity == nil {
		return logical.ErrorResponse("invalid entity ID"), nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return logical.ErrorResponse("failed to retrieve the namespace"), nil
	}
	if ns.ID != entity.NamespaceID {
		return logical.ErrorResponse("entity namespace ID does not match the current namespace ID"), nil
	}

	entityNS, err := i.namespacer.NamespaceByID(ctx, entity.NamespaceID)
	if err != nil {
		return logical.ErrorResponse("entity namespace not found"), nil
	}

	configNS, err := i.namespacer.NamespaceByID(ctx, mConfig.NamespaceID)
	if err != nil {
		return logical.ErrorResponse("methodID namespace not found"), nil
	}

	if configNS.ID != entityNS.ID && !entityNS.HasParent(configNS) {
		return logical.ErrorResponse(fmt.Sprintf("entity namespace %s outside of the config namespace %s", entityNS.Path, configNS.Path)), nil
	}

	switch mConfig.Type {
	case mfaMethodTypeTOTP:
		return i.mfaBackend.handleMFAGenerateTOTP(ctx, mConfig, entityID)
	default:
		return logical.ErrorResponse(fmt.Sprintf("generate not available for MFA type %q", mConfig.Type)), nil
	}
}

func (i *IdentityStore) handleLoginMFAAdminDestroyUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var entity *identity.Entity
	var err error

	methodID := d.Get("method_id").(string)
	if methodID == "" {
		return logical.ErrorResponse("missing method ID"), nil
	}

	entityID := d.Get("entity_id").(string)
	if entityID == "" {
		return logical.ErrorResponse("missing entity ID"), nil
	}

	entity, err = i.MemDBEntityByID(entityID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to find entity with ID %q: error: %w", entityID, err)
	}

	if entity == nil {
		return logical.ErrorResponse("invalid entity ID"), nil
	}

	mConfig, err := i.mfaBackend.MemDBMFAConfigByID(methodID)
	if err != nil {
		return nil, err
	}

	if mConfig == nil {
		return logical.ErrorResponse(fmt.Sprintf("configuration for method ID %q does not exist", methodID)), nil
	}

	if mConfig.ID == "" {
		return nil, fmt.Errorf("configuration for method ID %q does not contain an identifier", methodID)
	}

	if mConfig.Type != mfaMethodTypeTOTP {
		return nil, fmt.Errorf("method ID does not match TOTP type")
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return logical.ErrorResponse("failed to retrieve the namespace"), nil
	}
	if ns.ID != entity.NamespaceID {
		return logical.ErrorResponse("entity namespace ID does not match the current namespace ID"), nil
	}

	entityNS, err := i.namespacer.NamespaceByID(ctx, entity.NamespaceID)
	if err != nil {
		return logical.ErrorResponse("entity namespace not found"), nil
	}

	configNS, err := i.namespacer.NamespaceByID(ctx, mConfig.NamespaceID)
	if err != nil {
		return logical.ErrorResponse("methodID namespace not found"), nil
	}

	if configNS.ID != entityNS.ID && !entityNS.HasParent(configNS) {
		return logical.ErrorResponse(fmt.Sprintf("entity namespace %s outside of the current namespace %s", entityNS.Path, ns.Path)), nil
	}

	// destroying the secret on the entity
	if entity.MFASecrets != nil {
		delete(entity.MFASecrets, mConfig.ID)
	}

	err = i.upsertEntity(ctx, entity, nil, true)
	if err != nil {
		return nil, fmt.Errorf("failed to persist MFA secret in entity, error: %w", err)
	}

	return nil, nil
}

// loadMFAMethodConfigs loads MFA method configs for login MFA
func (b *LoginMFABackend) loadMFAMethodConfigs(ctx context.Context, ns *namespace.Namespace) error {
	b.mfaLogger.Trace("loading login MFA configurations")
	barrierView, err := b.Core.barrierViewForNamespace(ns.ID)
	if err != nil {
		return fmt.Errorf("error getting namespace view, namespaceid %s, error %w", ns.ID, err)
	}
	existing, err := barrierView.List(ctx, loginMFAConfigPrefix)
	if err != nil {
		return fmt.Errorf("failed to list MFA configurations for namespace path %s and prefix %s: %w", ns.Path, loginMFAConfigPrefix, err)
	}
	b.mfaLogger.Trace("methods collected", "num_existing", len(existing))

	for _, key := range existing {
		b.mfaLogger.Trace("loading method", "method", key)

		// Read the config from storage
		mConfig, err := b.getMFAConfig(ctx, loginMFAConfigPrefix+key, barrierView)
		if err != nil {
			return err
		}

		if mConfig == nil {
			b.mfaLogger.Trace("failed to find the config related to a method", "namespace", ns.Path, "prefix", loginMFAConfigPrefix, "method", key)
			continue
		}

		// Load the config in MemDB
		err = b.MemDBUpsertMFAConfig(ctx, mConfig)
		if err != nil {
			return fmt.Errorf("failed to load configuration ID %s prefix %s in MemDB: %w", mConfig.ID, loginMFAConfigPrefix, err)
		}
	}

	b.mfaLogger.Trace("configurations restored", "namespace", ns.Path, "prefix", loginMFAConfigPrefix)

	return nil
}

// loadMFAEnforcementConfigs loads MFA method configs for login MFA
func (b *LoginMFABackend) loadMFAEnforcementConfigs(ctx context.Context, ns *namespace.Namespace) ([]*mfa.MFAEnforcementConfig, error) {
	b.mfaLogger.Trace("loading login MFA enforcement configurations")
	barrierView, err := b.Core.barrierViewForNamespace(ns.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting namespace view, namespaceid %s, error %w", ns.ID, err)
	}
	existing, err := barrierView.List(ctx, mfaLoginEnforcementPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list MFA enforcement configurations for namespace %s with prefix %s: %w", ns.Path, mfaLoginEnforcementPrefix, err)
	}
	b.mfaLogger.Trace("enforcements configs collected", "num_existing", len(existing))

	eConfigs := make([]*mfa.MFAEnforcementConfig, 0)
	for _, key := range existing {
		b.mfaLogger.Trace("loading enforcement", "config", key)

		// Read the config from storage
		mConfig, err := b.getMFALoginEnforcementConfig(ctx, mfaLoginEnforcementPrefix+key, barrierView)
		if err != nil {
			return nil, err
		}

		if mConfig == nil {
			b.mfaLogger.Trace("failed to find an enforcement config", "namespace", ns.Path, "prefix", mfaLoginEnforcementPrefix, "config", key)
			continue
		}

		// Load the config in MemDB
		err = b.MemDBUpsertMFALoginEnforcementConfig(ctx, mConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load enforcement configuration ID %s with prefix %s in MemDB: %w", mConfig.ID, mfaLoginEnforcementPrefix, err)
		}

		eConfigs = append(eConfigs, mConfig)
	}

	b.mfaLogger.Trace("enforcement configurations restored", "namespace", ns.Path, "prefix", mfaLoginEnforcementPrefix)

	return eConfigs, nil
}

func (b *LoginMFABackend) loginMFAMethodExistenceCheck(eConfig *mfa.MFAEnforcementConfig) error {
	var aggErr *multierror.Error
	for _, confID := range eConfig.MFAMethodIDs {
		config, memErr := b.MemDBMFAConfigByID(confID)
		if memErr != nil {
			aggErr = multierror.Append(aggErr, memErr)
			return aggErr.ErrorOrNil()
		}
		if config == nil {
			aggErr = multierror.Append(aggErr, fmt.Errorf("found an MFA method ID in enforcement config, but failed to find the MFA method config method ID %s", confID))
		}
	}

	return aggErr.ErrorOrNil()
}

// sanitizeMFACredsWithLoginEnforcementMethodIDs updates the MFACred map
// looping through the matched login enforcement configurations, and
// replacing MFA method names with MFA method IDs
func (b *LoginMFABackend) sanitizeMFACredsWithLoginEnforcementMethodIDs(ctx context.Context, mfaCredsMap logical.MFACreds, mfaMethodIDs []string) (logical.MFACreds, error) {
	sanitizedMfaCreds := make(logical.MFACreds, 0)
	var multiError *multierror.Error
	for _, methodID := range mfaMethodIDs {
		val, ok := mfaCredsMap[methodID]
		if ok {
			sanitizedMfaCreds[methodID] = val
			continue
		}
		mConfig, err := b.MemDBMFAConfigByID(methodID)
		if err != nil {
			return nil, err
		}
		if mConfig == nil {
			multiError = multierror.Append(multiError, fmt.Errorf("failed to find MFA config for method ID %s", methodID))
			continue
		}

		// method name in the MFACredsMap should be the method full name,
		// i.e., namespacePath+name. This is because, a user in a child
		// namespace can reference an MFA method ID in a parent namespace
		configNS, err := NamespaceByID(ctx, mConfig.NamespaceID, b.Core)
		if err != nil {
			return nil, err
		}
		if configNS != nil {
			val, ok = mfaCredsMap[configNS.Path+mConfig.Name]
			if ok {
				sanitizedMfaCreds[mConfig.ID] = val
			} else {
				multiError = multierror.Append(multiError, fmt.Errorf("failed to find MFA credentials associated with an MFA method ID %v, method name %v", methodID, configNS.Path+mConfig.Name))
			}
		} else {
			multiError = multierror.Append(multiError, fmt.Errorf("failed to find the namespace associated with an MFA method ID %v", mConfig.ID))
		}
	}

	// we don't need to find every MFA method identifiers in the MFA header
	// So, don't return errors if that is the case.
	if len(sanitizedMfaCreds) > 0 {
		return sanitizedMfaCreds, nil
	}

	return sanitizedMfaCreds, multiError
}

func (b *LoginMFABackend) handleMFALoginValidate(ctx context.Context, req *logical.Request, d *framework.FieldData) (retResp *logical.Response, retErr error) {
	// mfaReqID is the ID of the login request
	mfaReqID := d.Get("mfa_request_id").(string)
	if mfaReqID == "" {
		return logical.ErrorResponse("missing request ID"), nil
	}

	// a map of methodID to passcode
	mfaPayload := d.Get("mfa_payload")
	if mfaPayload == nil {
		return logical.ErrorResponse("missing mfa payload"), nil
	}

	var mfaCreds logical.MFACreds
	err := mapstructure.Decode(mfaPayload, &mfaCreds)
	if err != nil {
		return logical.ErrorResponse("invalid mfa payload"), nil
	}

	// getting the cached response Auth. We should note that the entry is
	// removed from the queue, and if any error happens before the validation
	// and creating a token succeed, we need to push the entry back to the queue.
	cachedResponseAuth, err := b.Core.PopMFAResponseAuthByID(mfaReqID)
	if err != nil || cachedResponseAuth == nil {
		return logical.ErrorResponse("invalid request ID"), nil
	}
	defer func() {
		// Only if retErr is NOT nil, then push back the valid entry
		if retErr == nil {
			return
		}
		pushErr := b.Core.SaveMFAResponseAuth(cachedResponseAuth)
		if pushErr != nil {
			retErr = multierror.Append(retErr, pushErr)
		}
	}()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("MFA validation failed. Namespace not found. error: %v", err)
	}

	if ns.ID != cachedResponseAuth.RequestNSID {
		return nil, fmt.Errorf("original request was issued in a different namesapce %v, current namespace is %v", cachedResponseAuth.RequestNSPath, ns.Path)
	}

	entity, _, err := b.Core.fetchEntityAndDerivedPolicies(ctx, ns, cachedResponseAuth.CachedAuth.EntityID, true)
	if err != nil || entity == nil {
		return nil, fmt.Errorf("MFA validation failed. entity not found: %v", err)
	}

	// finding the MFAEnforcement config that matches our ns. ns could be root as well
	matchedMfaEnforcementList, err := b.Core.buildMFAEnforcementConfigList(ctx, entity, cachedResponseAuth.RequestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find MFAEnforcement configuration")
	}

	if len(matchedMfaEnforcementList) == 0 {
		return nil, fmt.Errorf("found nil or empty MFAEnforcement configuration")
	}

	for _, eConfig := range matchedMfaEnforcementList {
		err = b.Core.validateLoginMFA(ctx, eConfig, entity, req.Connection.RemoteAddr, mfaCreds)
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to satisfy enforcement %s. error: %s", eConfig.Name, err.Error())), logical.ErrPermissionDenied
		}
	}

	// MFA validation has passed. Let's generate the token
	resp, err := b.Core.LoginMFACreateToken(ctx, cachedResponseAuth.RequestPath, cachedResponseAuth.CachedAuth, req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to create a token. error: %v", err)
	}

	return resp, nil
}

func (c *Core) teardownLoginMFA() error {
	if !c.IsDRSecondary() {
		// Clear any cached auth response
		c.mfaResponseAuthQueueLock.Lock()
		c.mfaResponseAuthQueue = nil
		c.mfaResponseAuthQueueLock.Unlock()

		c.loginMFABackend.usedCodes = nil

		if err := c.loginMFABackend.ResetLoginMFAMemDB(); err != nil {
			return err
		}
	}
	return nil
}

// LoginMFACreateToken creates a token after the login MFA is validated.
// It also applies the lease quotas on the original login request path.
func (c *Core) LoginMFACreateToken(ctx context.Context, reqPath string, cachedAuth *logical.Auth, loginRequestData map[string]interface{}) (*logical.Response, error) {
	auth := cachedAuth
	resp := &logical.Response{
		Auth: auth,
	}

	// Determine the source of the login
	mountPoint := c.router.MatchingMount(ctx, reqPath)

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("namespace not found: %w", err)
	}

	var role string
	if reqRole := ctx.Value(logical.CtxKeyRequestRole{}); reqRole != nil {
		role = reqRole.(string)
	}

	// The request successfully authenticated itself. Run the quota checks on
	// the original login request path before creating the token.
	quotaResp, quotaErr := c.applyLeaseCountQuota(ctx, &quotas.Request{
		Path:          reqPath,
		MountPath:     strings.TrimPrefix(mountPoint, ns.Path),
		Role:          role,
		NamespacePath: ns.Path,
	})

	if quotaErr != nil {
		c.logger.Error("failed to apply quota", "path", reqPath, "error", quotaErr)
		return nil, quotaErr
	}

	if !quotaResp.Allowed {
		if c.logger.IsTrace() {
			c.logger.Trace("request rejected due to lease count quota violation", "request_path", reqPath)
		}

		return nil, fmt.Errorf("request path %q: %w", reqPath, quotas.ErrLeaseCountQuotaExceeded)
	}

	// note that we don't need to handle the error for the following function right away.
	// The function takes the response as in input variable and modify it. So, the returned
	// arguments are resp and err.
	leaseGenerated, resp, err := c.LoginCreateToken(ctx, ns, reqPath, mountPoint, role, resp)

	if quotaResp.Access != nil {
		quotaAckErr := c.ackLeaseQuota(quotaResp.Access, leaseGenerated)
		if quotaAckErr != nil {
			err = multierror.Append(err, quotaAckErr)
		}
	}

	return resp, err
}

func (i *IdentityStore) handleMFALoginEnforcementList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keys, configInfo, err := i.mfaBackend.mfaLoginEnforcementList(ctx)
	if err != nil {
		return nil, err
	}

	return logical.ListResponseWithInfo(keys, configInfo), nil
}

func (i *IdentityStore) handleMFALoginEnforcementRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	respData, err := i.mfaBackend.mfaLoginEnforcementConfigByNameAndNamespace(name, ns.ID)
	if err != nil {
		return nil, err
	}

	if respData == nil {
		return nil, nil
	}

	// The config is readable only from the same namespace
	if ns.ID != respData["namespace_id"].(string) {
		return logical.ErrorResponse("request namespace does not match method namespace"), nil
	}

	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) handleMFALoginEnforcementUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var err error
	var eConfig *mfa.MFAEnforcementConfig

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing enforcement name"), nil
	}

	b := i.mfaBackend
	b.mfaLock.Lock()
	defer b.mfaLock.Unlock()

	eConfig, err = b.MemDBMFALoginEnforcementConfigByNameAndNamespace(name, ns.ID)
	if err != nil {
		return nil, err
	}

	if eConfig == nil {
		configID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate an identifier for MFA login enforcement config: %w", err)
		}
		eConfig = &mfa.MFAEnforcementConfig{
			Name:        name,
			NamespaceID: ns.ID,
			ID:          configID,
		}
	}

	mfaMethodIds, ok := d.GetOk("mfa_method_ids")
	if !ok {
		return logical.ErrorResponse("missing method ids"), nil
	}

	for _, mmid := range mfaMethodIds.([]string) {
		// make sure this method id actually exists
		config, err := b.mfaConfigReadByMethodID(mmid)
		if err != nil {
			return nil, err
		}
		if config == nil {
			return logical.ErrorResponse("one of the provided method ids doesn't exist"), nil
		}

		mfaNs, err := i.namespacer.NamespaceByID(ctx, config["namespace_id"].(string))
		if err != nil {
			return logical.ErrorResponse("failed to retrieve config namespace"), nil
		}

		if ns.ID != mfaNs.ID && !ns.HasParent(mfaNs) {
			return logical.ErrorResponse("one of the provided method ids is in an incompatible namespace and can't be used"), nil
		}
	}
	eConfig.MFAMethodIDs = mfaMethodIds.([]string)

	oneOfLastFour := false
	authMethodAccessors, ok := d.GetOk("auth_method_accessors")
	if ok {
		for _, accessor := range authMethodAccessors.([]string) {
			found, err := b.validateAuthEntriesForAccessorOrType(ctx, ns, func(entry *MountEntry) bool {
				return accessor == entry.Accessor
			})
			if err != nil {
				return nil, err
			}
			if !found {
				return logical.ErrorResponse("one of the auth method accessors provided is invalid"), nil
			}
		}
		eConfig.AuthMethodAccessors = authMethodAccessors.([]string)
		oneOfLastFour = true
	}

	authMethodTypes, ok := d.GetOk("auth_method_types")
	if ok {
		for _, authType := range authMethodTypes.([]string) {
			found, err := b.validateAuthEntriesForAccessorOrType(ctx, ns, func(entry *MountEntry) bool {
				return authType == entry.Type
			})
			if err != nil {
				return nil, err
			}
			if !found {
				return logical.ErrorResponse("one of the auth method types provided is invalid"), nil
			}
		}
		eConfig.AuthMethodTypes = authMethodTypes.([]string)
		oneOfLastFour = true
	}

	identityGroupIds, ok := d.GetOk("identity_group_ids")
	if ok {
		for _, groupId := range identityGroupIds.([]string) {
			group, err := i.MemDBGroupByID(groupId, true)
			if err != nil {
				return nil, err
			}
			if group == nil {
				return logical.ErrorResponse("one of the provided group ids doesn't exist"), nil
			}
		}
		eConfig.IdentityGroupIds = identityGroupIds.([]string)
		oneOfLastFour = true
	}

	identityEntityIds, ok := d.GetOk("identity_entity_ids")
	if ok {
		for _, entityId := range identityEntityIds.([]string) {
			entity, err := i.MemDBEntityByID(entityId, true)
			if err != nil {
				return nil, err
			}
			if entity == nil {
				return logical.ErrorResponse("one of the provided entity ids doesn't exist"), nil
			}
		}
		eConfig.IdentityEntityIDs = identityEntityIds.([]string)
		oneOfLastFour = true
	}

	if !oneOfLastFour {
		return logical.ErrorResponse("One of auth_method_accessors, auth_method_types, identity_group_ids, identity_entity_ids must be specified"), nil
	}

	// Store the config
	err = b.putMFALoginEnforcementConfig(ctx, eConfig)
	if err != nil {
		return nil, err
	}

	// Back the config in MemDB
	return nil, b.MemDBUpsertMFALoginEnforcementConfig(ctx, eConfig)
}

func (i *IdentityStore) handleMFALoginEnforcementDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	return nil, i.mfaBackend.deleteMFALoginEnforcementConfigByNameAndNamespace(ctx, name, ns.ID)
}

func (b *LoginMFABackend) validateAuthEntriesForAccessorOrType(ctx context.Context, ns *namespace.Namespace, validFunc func(entry *MountEntry) bool) (bool, error) {
	b.Core.authLock.RLock()
	defer b.Core.authLock.RUnlock()

	for _, entry := range b.Core.auth.Entries {
		// only check auth methods in the current namespace
		if entry.Namespace().ID != ns.ID {
			continue
		}

		cont, err := b.Core.checkReplicatedFiltering(ctx, entry, credentialRoutePrefix)
		if err != nil {
			return false, err
		}
		if cont {
			continue
		}

		if validFunc(entry) {
			return true, nil
		}
	}

	return false, nil
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

func (b *MFABackend) handleMFAGenerateTOTP(ctx context.Context, mConfig *mfa.Config, entityID string) (*logical.Response, error) {
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
		UsePasscode:    d.Get("use_passcode").(bool),
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

func (b *LoginMFABackend) mfaConfigReadByMethodID(id string) (map[string]interface{}, error) {
	mConfig, err := b.MemDBMFAConfigByID(id)
	if err != nil {
		return nil, err
	}
	if mConfig == nil {
		return nil, nil
	}

	return b.mfaConfigToMap(mConfig)
}

func (b *LoginMFABackend) mfaMethodList(ctx context.Context, methodType string) ([]string, map[string]interface{}, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	ws := memdb.NewWatchSet()
	txn := b.db.Txn(false)

	var iter memdb.ResultIterator
	switch {
	case methodType == "":
		// get all the configs
		iter, err = txn.Get(b.methodTable, "id")
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch iterator for login mfa method configs in memdb: %w", err)
		}
	default:
		// get all the configs for the given type
		iter, err = txn.Get(b.methodTable, "type", methodType)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch iterator for login mfa method configs in memdb: %w", err)
		}
	}

	ws.Add(iter.WatchCh())

	var keys []string
	configInfo := map[string]interface{}{}

	for {
		// check for timeouts
		select {
		case <-ctx.Done():
			return keys, configInfo, nil
		default:
			break
		}

		raw := iter.Next()
		if raw == nil {
			break
		}
		config := raw.(*mfa.Config)

		// return this config if it's in the same ns as the request ns OR it's in a parent ns of the request ns
		mfaNs, err := b.namespacer.NamespaceByID(ctx, config.NamespaceID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch namespace: %w", err)
		}

		// the namespaces have to match, or the config namespace needs to be a parent of the request namespace
		if !(ns.ID == mfaNs.ID || ns.HasParent(mfaNs)) {
			continue
		}

		keys = append(keys, config.ID)
		configInfoEntry, err := b.mfaConfigToMap(config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert config to map: %w", err)
		}
		configInfo[config.ID] = configInfoEntry
	}

	return keys, configInfo, nil
}

func (b *LoginMFABackend) mfaLoginEnforcementList(ctx context.Context) ([]string, map[string]interface{}, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	ws := memdb.NewWatchSet()
	txn := b.db.Txn(false)

	// get all the login enforcements in our namespace
	iter, err := txn.Get(memDBMFALoginEnforcementsTable, "namespace", ns.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch iterator for login enforcement configs in memdb: %w", err)
	}

	ws.Add(iter.WatchCh())

	var keys []string
	enforcementInfo := map[string]interface{}{}

	for {
		// check for timeouts
		select {
		case <-ctx.Done():
			return keys, enforcementInfo, nil
		default:
			break
		}

		raw := iter.Next()
		if raw == nil {
			break
		}
		config := raw.(*mfa.MFAEnforcementConfig)
		keys = append(keys, config.Name)
		configInfoEntry, err := b.mfaLoginEnforcementConfigToMap(config)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert enforcement to map: %w", err)
		}
		enforcementInfo[config.Name] = configInfoEntry
	}

	return keys, enforcementInfo, nil
}

func (b *LoginMFABackend) mfaLoginEnforcementConfigByNameAndNamespace(name, namespaceId string) (map[string]interface{}, error) {
	eConfig, err := b.MemDBMFALoginEnforcementConfigByNameAndNamespace(name, namespaceId)
	if err != nil {
		return nil, err
	}
	if eConfig == nil {
		return nil, nil
	}
	return b.mfaLoginEnforcementConfigToMap(eConfig)
}

func (b *LoginMFABackend) mfaLoginEnforcementConfigToMap(eConfig *mfa.MFAEnforcementConfig) (map[string]interface{}, error) {
	resp := make(map[string]interface{})
	resp["name"] = eConfig.Name
	ns, err := b.namespacer.NamespaceByID(context.Background(), eConfig.NamespaceID)
	if err != nil {
		return nil, err
	}
	if ns != nil {
		resp["namespace_path"] = ns.Path
	}
	resp["namespace_id"] = eConfig.NamespaceID
	resp["mfa_method_ids"] = append([]string{}, eConfig.MFAMethodIDs...)
	resp["auth_method_accessors"] = append([]string{}, eConfig.AuthMethodAccessors...)
	resp["auth_method_types"] = append([]string{}, eConfig.AuthMethodTypes...)
	resp["identity_group_ids"] = append([]string{}, eConfig.IdentityGroupIds...)
	resp["identity_entity_ids"] = append([]string{}, eConfig.IdentityEntityIDs...)
	resp["id"] = eConfig.ID
	return resp, nil
}

func (b *MFABackend) mfaConfigToMap(mConfig *mfa.Config) (map[string]interface{}, error) {
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
		respData["max_validation_attempts"] = totpConfig.MaxValidationAttempts
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
		respData["use_passcode"] = duoConfig.UsePasscode
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
	respData["namespace_id"] = mConfig.NamespaceID
	ns, err := b.namespacer.NamespaceByID(context.Background(), mConfig.NamespaceID)
	if err != nil {
		return nil, err
	}
	if ns != nil {
		respData["namespace_path"] = ns.Path
	}

	return respData, nil
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

	maxValidationAttempt := d.Get("max_validation_attempts").(int)
	if maxValidationAttempt < 0 {
		return fmt.Errorf("max_validation_attempts must be greater than zero")
	}
	if maxValidationAttempt == 0 {
		maxValidationAttempt = defaultMaxTOTPValidateAttempts
	}

	config := &mfa.TOTPConfig{
		Issuer:                issuer,
		Period:                uint32(period),
		Algorithm:             int32(keyAlgorithm),
		Digits:                int32(keyDigits),
		Skew:                  uint32(skew),
		KeySize:               uint32(keySize),
		QRSize:                int32(d.Get("qr_size").(int)),
		MaxValidationAttempts: uint32(maxValidationAttempt),
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

func (c *Core) validateLoginMFA(ctx context.Context, eConfig *mfa.MFAEnforcementConfig, entity *identity.Entity, requestConnRemoteAddr string, mfaCredsMap logical.MFACreds) error {
	sanitizedMfaCreds, err := c.loginMFABackend.sanitizeMFACredsWithLoginEnforcementMethodIDs(ctx, mfaCredsMap, eConfig.MFAMethodIDs)
	if err != nil {
		return fmt.Errorf("failed to sanitize MFA creds, %w", err)
	}
	if len(sanitizedMfaCreds) == 0 && len(eConfig.MFAMethodIDs) > 0 {
		return fmt.Errorf("login MFA validation failed for methodID: %v", eConfig.MFAMethodIDs)
	}

	var retErr error
	for _, methodID := range eConfig.MFAMethodIDs {
		// as configID is the same as methodID, and methodID is unique, we can
		// use it to retrieve the MFACreds
		mfaCreds, ok := sanitizedMfaCreds[methodID]
		if !ok || mfaCreds == nil {
			continue
		}

		err := c.validateLoginMFAInternal(ctx, methodID, entity, requestConnRemoteAddr, mfaCreds)
		if err != nil {
			retErr = multierror.Append(retErr, err)
			continue
		}
		return nil
	}

	return multierror.Append(retErr, fmt.Errorf("login MFA validation failed for methodID: %v", eConfig.MFAMethodIDs))
}

func (c *Core) validateLoginMFAInternal(ctx context.Context, methodID string, entity *identity.Entity, reqConnectionRemoteAddress string, mfaCreds []string) (retErr error) {
	if entity == nil {
		return fmt.Errorf("entity is nil")
	}

	// Get the configuration for the MFA method set in system backend
	mConfig, err := c.loginMFABackend.MemDBMFAConfigByID(methodID)
	if err != nil {
		return fmt.Errorf("failed to read MFA configuration")
	}

	if mConfig == nil {
		return fmt.Errorf("MFA method configuration not present")
	}

	var finalUsername string
	switch mConfig.Type {
	case mfaMethodTypeDuo, mfaMethodTypeOkta, mfaMethodTypePingID:
		if mConfig.UsernameFormat == "" {
			finalUsername = entity.Name
		} else {
			directGroups, inheritedGroups, err := c.identityStore.groupsByEntityID(entity.ID)
			if err != nil {
				return fmt.Errorf("failed to fetch group memberships: %w", err)
			}
			groups := append(directGroups, inheritedGroups...)

			_, finalUsername, err = identitytpl.PopulateString(identitytpl.PopulateStringInput{
				Mode:        identitytpl.ACLTemplating,
				String:      mConfig.UsernameFormat,
				Entity:      identity.ToSDKEntity(entity),
				Groups:      identity.ToSDKGroups(groups),
				NamespaceID: entity.NamespaceID,
			})
			if err != nil {
				return err
			}
		}
	}

	mfaFactors, err := parseMfaFactors(mfaCreds)
	if err != nil {
		return fmt.Errorf("failed to parse MFA factor, %w", err)
	}

	switch mConfig.Type {
	case mfaMethodTypeTOTP:
		// Get the MFA secret data required to validate the supplied credentials
		if entity.MFASecrets == nil {
			return fmt.Errorf("MFA secret for method ID %q not present in entity %q", mConfig.ID, entity.ID)
		}
		entityMFASecret := entity.MFASecrets[mConfig.ID]
		if entityMFASecret == nil {
			return fmt.Errorf("MFA secret for method name %q not present in entity %q", mConfig.Name, entity.ID)
		}

		return c.validateTOTP(ctx, mfaFactors, entityMFASecret, mConfig.ID, entity.ID, c.loginMFABackend.usedCodes, mConfig.GetTOTPConfig().MaxValidationAttempts)

	case mfaMethodTypeOkta:
		return c.validateOkta(ctx, mConfig, finalUsername)

	case mfaMethodTypeDuo:
		return c.validateDuo(ctx, mfaFactors, mConfig, finalUsername, reqConnectionRemoteAddress)

	case mfaMethodTypePingID:
		return c.validatePingID(ctx, mConfig, finalUsername)

	default:
		return fmt.Errorf("unrecognized MFA type %q", mConfig.Type)
	}
}

func (c *Core) buildMFAEnforcementConfigList(ctx context.Context, entity *identity.Entity, reqPath string) ([]*mfa.MFAEnforcementConfig, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace from context. %s, %v", "error", err)
	}

	eConfigIter, err := c.loginMFABackend.MemDBMFALoginEnforcementConfigIterator()
	if err != nil {
		return nil, err
	}

	me := c.router.MatchingMountEntry(ctx, reqPath)
	if me == nil {
		return nil, fmt.Errorf("failed to find matching mount entry for path %v", reqPath)
	}

	var matchedMfaEnforcementConfig []*mfa.MFAEnforcementConfig
	// finding the MFAEnforcement config that matches our ns. ns could be root as well
ECONFIG_LOOP:
	for eConfigRaw := eConfigIter.Next(); eConfigRaw != nil; eConfigRaw = eConfigIter.Next() {
		eConfig := eConfigRaw.(*mfa.MFAEnforcementConfig)

		// check if this config's ns applies to current req,
		// i.e. is it the req's ns or an ancestor of req's ns?
		eConfigNS, err := c.NamespaceByID(ctx, eConfig.NamespaceID)
		if err != nil {
			return nil, fmt.Errorf("failed to find the MFAEnforcementConfig namespace: %w", err)
		}

		if eConfig == nil || eConfigNS == nil || (eConfigNS.ID != ns.ID && !ns.HasParent(eConfigNS)) {
			continue
		}

		// if entity is nil, an MFAEnforcementConfig could still be configured
		// having mount type/accessor
		if entity != nil {
			if entity.NamespaceID != ns.ID {
				return nil, fmt.Errorf("entity namespace ID is different than the current ns ID")
			}

			// Check if entityID is in the MFAEnforcement config
			if strutil.StrListContains(eConfig.IdentityEntityIDs, entity.ID) {
				matchedMfaEnforcementConfig = append(matchedMfaEnforcementConfig, eConfig)
				continue
			}

			// Retrieve entity groups
			directGroups, inheritedGroups, err := c.identityStore.groupsByEntityID(entity.ID)
			if err != nil {
				return nil, fmt.Errorf("error on retrieving groups by entityID in MFA")
			}
			for _, g := range directGroups {
				if strutil.StrListContains(eConfig.IdentityGroupIds, g.ID) {
					matchedMfaEnforcementConfig = append(matchedMfaEnforcementConfig, eConfig)
					continue ECONFIG_LOOP
				}
			}
			for _, g := range inheritedGroups {
				if strutil.StrListContains(eConfig.IdentityGroupIds, g.ID) {
					matchedMfaEnforcementConfig = append(matchedMfaEnforcementConfig, eConfig)
					continue ECONFIG_LOOP
				}
			}
		}

		for _, acc := range eConfig.AuthMethodAccessors {
			if me != nil && me.Accessor == acc {
				matchedMfaEnforcementConfig = append(matchedMfaEnforcementConfig, eConfig)
				continue ECONFIG_LOOP
			}
		}

		for _, authT := range eConfig.AuthMethodTypes {
			if me != nil && me.Type == authT {
				matchedMfaEnforcementConfig = append(matchedMfaEnforcementConfig, eConfig)
				continue ECONFIG_LOOP
			}
		}
	}

	return matchedMfaEnforcementConfig, nil
}

func formatUsername(format string, alias *identity.Alias, entity *identity.Entity) string {
	if format == "" {
		return alias.Name
	}

	username := format
	username = strings.ReplaceAll(username, "{{alias.name}}", alias.Name)
	username = strings.ReplaceAll(username, "{{entity.name}}", entity.Name)
	for k, v := range alias.Metadata {
		username = strings.ReplaceAll(username, fmt.Sprintf("{{alias.metadata.%s}}", k), v)
	}
	for k, v := range entity.Metadata {
		username = strings.ReplaceAll(username, fmt.Sprintf("{{entity.metadata.%s}}", k), v)
	}
	return username
}

type MFAFactor struct {
	passcode string
}

func parseMfaFactors(creds []string) (*MFAFactor, error) {
	mfaFactor := &MFAFactor{}

	for _, cred := range creds {
		switch {
		case cred == "": // for the case of push notification
			continue
		case strings.HasPrefix(cred, "passcode="):
			if mfaFactor.passcode != "" {
				return nil, fmt.Errorf("found multiple passcodes for the same MFA method")
			}

			splits := strings.SplitN(cred, "=", 2)
			if splits[1] == "" {
				return nil, fmt.Errorf("invalid passcode")
			}

			mfaFactor.passcode = splits[1]
		case strings.Contains(cred, "="):
			return nil, fmt.Errorf("found an invalid MFA cred: %v", cred)
		default:
			// a non-empty cred that does not match the above
			// means it is a passcode
			if mfaFactor.passcode != "" {
				return nil, fmt.Errorf("found multiple passcodes for the same MFA method")
			}
			mfaFactor.passcode = cred
		}
	}

	if mfaFactor.passcode == "" {
		return nil, nil
	}

	return mfaFactor, nil
}

func (c *Core) validateDuo(ctx context.Context, mfaFactors *MFAFactor, mConfig *mfa.Config, username, reqConnectionRemoteAddr string) error {
	duoConfig := mConfig.GetDuoConfig()
	if duoConfig == nil {
		return fmt.Errorf("failed to get Duo configuration for method %q", mConfig.Name)
	}

	var passcode string
	if mfaFactors != nil {
		passcode = mfaFactors.passcode
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

	preauth, err := authClient.Preauth(authapi.PreauthUsername(username), authapi.PreauthIpAddr(reqConnectionRemoteAddr))
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
		return errors.New(preauth.Response.Status_Msg)
	case "enroll":
		return fmt.Errorf("%q - %q", preauth.Response.Status_Msg, preauth.Response.Enroll_Portal_Url)
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

	options = append(options, authapi.AuthIpAddr(reqConnectionRemoteAddr))
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
		timer := time.NewTimer(time.Second)

		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("duo push verification operation canceled")
		case <-timer.C:
		}
	}
}

func (c *Core) validateOkta(ctx context.Context, mConfig *mfa.Config, username string) error {
	oktaConfig := mConfig.GetOktaConfig()
	if oktaConfig == nil {
		return fmt.Errorf("failed to get Okta configuration for method %q", mConfig.Name)
	}

	baseURL := oktaConfig.BaseURL
	if baseURL == "" {
		baseURL = "okta.com"
	}
	orgURL, err := url.Parse(fmt.Sprintf("https://%s.%s", oktaConfig.OrgName, baseURL))
	if err != nil {
		return err
	}

	cfg, err := okta.NewConfiguration(
		okta.WithToken(oktaConfig.APIToken),
		okta.WithOrgUrl(orgURL.String()),
		// Do not use cache or polling MFA will not refresh
		okta.WithCache(false),
	)
	if err != nil {
		return fmt.Errorf("error creating client: %s", err)
	}
	client := okta.NewAPIClient(cfg)

	filterField := "profile.login"
	if oktaConfig.PrimaryEmail {
		filterField = "profile.email"
	}
	filterQuery := fmt.Sprintf("%s eq %q", filterField, username)

	users, _, err := client.UserAPI.ListUsers(client.GetConfig().Context).Filter(filterQuery).Execute()
	if err != nil {
		return err
	}
	switch {
	case len(users) == 0:
		return fmt.Errorf("no users found for e-mail address")
	case len(users) > 1:
		return fmt.Errorf("more than one user found for e-mail address")
	}

	user := users[0]

	factors, _, err := client.UserFactorAPI.ListFactors(ctx, user.GetId()).Execute()
	if err != nil {
		return err
	}

	if len(factors) == 0 {
		return fmt.Errorf("no MFA factors found for user")
	}

	var factorFound bool
	var userFactor *okta.UserFactorPush
	for _, factor := range factors {
		if factor.UserFactorPush != nil {
			userFactor = factor.UserFactorPush
			factorFound = true
			break
		}
	}

	if !factorFound {
		return fmt.Errorf("no push-type MFA factor found for user")
	}

	result, _, err := client.UserFactorAPI.VerifyFactor(ctx, user.GetId(), userFactor.GetId()).Execute()
	if err != nil {
		return err
	}

	if result.GetFactorResult() != "WAITING" {
		return fmt.Errorf("expected WAITING status for push status, got %q", result.GetFactorResult())
	}

	// Parse links to get polling link
	type linksObj struct {
		Poll struct {
			Href string `mapstructure:"href"`
		} `mapstructure:"poll"`
	}
	links := new(linksObj)
	if err := mapstructure.WeakDecode(result.Links, links); err != nil {
		return err
	}
	// Strip the org URL from the fully qualified poll URL
	url, err := url.Parse(strings.Replace(links.Poll.Href, orgURL.String(), "", 1))
	if err != nil {
		return err
	}

	// Okta doesn't return the transactionID as a parameter in the response, but it's encoded in the URL
	// this approach comes from: https://github.com/okta/okta-sdk-golang/issues/300, but it's not ideal.
	// It is, however, what the dotnet library by Okta themselves does.
	txRx := regexp.MustCompile("^.*/transactions/(.*)$")
	matches := txRx.FindStringSubmatch(url.Path)
	if len(matches) != 2 {
		return fmt.Errorf("couldn't determine transaction id from url")
	}
	transactionID := matches[1]

	// poll verifyfactor until termination (e.g., the user responds to the push factor)
	for {
		result, _, err := client.UserFactorAPI.GetFactorTransactionStatus(client.GetConfig().Context, user.GetId(), userFactor.GetId(), transactionID).Execute()
		if err != nil {
			return err
		}

		// the transaction status returns an inner object set based on what the factor status is.
		// the other ones are nil. This is (probably) because the structure of the returned JSON
		// varies based on what the factor status is.
		switch {
		case result.UserFactorPushTransactionWaiting != nil:
		case result.UserFactorPushTransaction != nil:
			return nil
		case result.UserFactorPushTransactionRejected != nil:
			return fmt.Errorf("push verification explicitly rejected")
		case result.UserFactorPushTransactionTimeout != nil:
			return fmt.Errorf("push verification timed out")
		default:
			return fmt.Errorf("unknown status code")
		}
		timer := time.NewTimer(time.Second)

		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("push verification operation canceled")
		case <-timer.C:
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
		req.Body = io.NopCloser(bytes.NewBufferString(signedToken))
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

func (c *Core) validateTOTP(ctx context.Context, mfaFactors *MFAFactor, entityMethodSecret *mfa.Secret, configID, entityID string, usedCodes *cache.Cache, maximumValidationAttempts uint32) error {
	if mfaFactors == nil || mfaFactors.passcode == "" {
		return fmt.Errorf("MFA credentials not supplied")
	}
	passcode := mfaFactors.passcode

	totpSecret := entityMethodSecret.GetTOTPSecret()
	if totpSecret == nil {
		return fmt.Errorf("entity does not contain the TOTP secret")
	}

	usedName := fmt.Sprintf("%s_%s", configID, passcode)

	_, ok := usedCodes.Get(usedName)
	if ok {
		return fmt.Errorf("code already used; new code is available in %v seconds", totpSecret.Period)
	}

	// The duration in which a passcode is stored in cache to enforce
	// rate limit on failed totp passcode validation
	passcodeTTL := time.Duration(int64(time.Second) * int64(totpSecret.Period))

	// Enforcing rate limit per MethodID per EntityID
	rateLimitID := fmt.Sprintf("%s_%s", configID, entityID)

	numAttempts, _ := usedCodes.Get(rateLimitID)
	if numAttempts == nil {
		usedCodes.Set(rateLimitID, uint32(1), passcodeTTL)
	} else {
		num, ok := numAttempts.(uint32)
		if !ok {
			return fmt.Errorf("invalid counter type returned in TOTP usedCode cache")
		}
		if num == maximumValidationAttempts {
			return fmt.Errorf("maximum TOTP validation attempts %d exceeded the allowed attempts %d. Please try again in %v seconds", num+1, maximumValidationAttempts, passcodeTTL)
		}
		err := usedCodes.Increment(rateLimitID, 1)
		if err != nil {
			return fmt.Errorf("failed to increment the TOTP code counter")
		}
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

	valid, err := totplib.ValidateCustom(passcode, key, time.Now(), validateOpts)
	if err != nil && err != otplib.ErrValidateInputInvalidLength {
		return errwrap.Wrapf("failed to validate TOTP passcode: {{err}}", err)
	}

	if !valid {
		return fmt.Errorf("failed to validate TOTP passcode")
	}

	// Take the key skew, add two for behind and in front, and multiply that by
	// the period to cover the full possibility of the validity of the key
	validityPeriod := time.Duration(int64(time.Second) * int64(totpSecret.Period) * int64(2+totpSecret.Skew))

	// Adding the used code to the cache
	err = usedCodes.Add(usedName, nil, validityPeriod)
	if err != nil {
		return fmt.Errorf("error adding code to used cache: %w", err)
	}

	// deleting the cache entry after a successful MFA validation
	usedCodes.Delete(rateLimitID)

	return nil
}

func loginMFAConfigTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: memDBLoginMFAConfigsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"namespace_id": {
				Name:   "namespace_id",
				Unique: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "NamespaceID",
				},
			},
			"type": {
				Name:   "type",
				Unique: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "Type",
				},
			},
			"name": {
				Name:         "name",
				Unique:       true,
				AllowMissing: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
						&memdb.StringFieldIndex{
							Field: "Name",
						},
					},
				},
			},
		},
	}
}

// turns out every memdb table schema must have an id index
func loginEnforcementTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: memDBMFALoginEnforcementsTable,
		Indexes: map[string]*memdb.IndexSchema{
			"id": {
				Name:   "id",
				Unique: true,
				Indexer: &memdb.StringFieldIndex{
					Field: "ID",
				},
			},
			"namespace": {
				Name:   "namespace",
				Unique: false,
				Indexer: &memdb.StringFieldIndex{
					Field: "NamespaceID",
				},
			},
			"nameAndNamespace": {
				Name:   "nameAndNamespace",
				Unique: true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "Name",
						},
						&memdb.StringFieldIndex{
							Field: "NamespaceID",
						},
					},
				},
			},
		},
	}
}

func (b *MFABackend) MemDBUpsertMFAConfig(ctx context.Context, mConfig *mfa.Config) error {
	txn := b.db.Txn(true)
	defer txn.Abort()

	err := b.MemDBUpsertMFAConfigInTxn(txn, mConfig)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (b *MFABackend) MemDBUpsertMFAConfigInTxn(txn *memdb.Txn, mConfig *mfa.Config) error {
	if txn == nil {
		return fmt.Errorf("nil txn")
	}

	if mConfig == nil {
		return fmt.Errorf("config is nil")
	}

	mConfigRaw, err := txn.First(b.methodTable, "id", mConfig.ID)
	if err != nil {
		return errwrap.Wrapf("failed to lookup MFA config from MemDB using id: {{err}}", err)
	}

	if mConfigRaw != nil {
		err = txn.Delete(b.methodTable, mConfigRaw)
		if err != nil {
			return errwrap.Wrapf("failed to delete MFA config from MemDB: {{err}}", err)
		}
	}

	if err := txn.Insert(b.methodTable, mConfig); err != nil {
		return errwrap.Wrapf("failed to update MFA config into MemDB: {{err}}", err)
	}

	return nil
}

func (b *LoginMFABackend) MemDBUpsertMFALoginEnforcementConfig(ctx context.Context, eConfig *mfa.MFAEnforcementConfig) error {
	if eConfig == nil {
		return fmt.Errorf("config is nil")
	}

	txn := b.db.Txn(true)
	defer txn.Abort()

	eConfigRaw, err := txn.First(memDBMFALoginEnforcementsTable, "nameAndNamespace", eConfig.Name, eConfig.NamespaceID)
	if err != nil {
		return fmt.Errorf("failed to lookup MFA login enforcement config from MemDB using name: %w", err)
	}

	if eConfigRaw != nil {
		err = txn.Delete(memDBMFALoginEnforcementsTable, eConfigRaw)
		if err != nil {
			return fmt.Errorf("failed to delete MFA login enforcement config from MemDB: %w", err)
		}
	}

	if err := txn.Insert(memDBMFALoginEnforcementsTable, eConfig); err != nil {
		return fmt.Errorf("failed to update MFA login enforcement config in MemDB: %w", err)
	}

	txn.Commit()
	return nil
}

func (b *LoginMFABackend) MemDBMFAConfigByIDInTxn(txn *memdb.Txn, mConfigID string) (*mfa.Config, error) {
	if mConfigID == "" {
		return nil, fmt.Errorf("missing config id")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	mConfigRaw, err := txn.First(b.methodTable, "id", mConfigID)
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

	return mConfig.Clone()
}

func (b *LoginMFABackend) MemDBMFAConfigByID(mConfigID string) (*mfa.Config, error) {
	if mConfigID == "" {
		return nil, fmt.Errorf("missing config id")
	}

	txn := b.db.Txn(false)

	return b.MemDBMFAConfigByIDInTxn(txn, mConfigID)
}

func (b *LoginMFABackend) MemDBMFAConfigByNameInTxn(ctx context.Context, txn *memdb.Txn, mConfigName string) (*mfa.Config, error) {
	if mConfigName == "" {
		return nil, fmt.Errorf("missing config name")
	}

	if txn == nil {
		return nil, fmt.Errorf("txn is nil")
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	mConfigRaw, err := txn.First(b.methodTable, "name", ns.ID, mConfigName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MFA config from memdb using name: %w", err)
	}

	if mConfigRaw == nil {
		return nil, nil
	}

	mConfig, ok := mConfigRaw.(*mfa.Config)
	if !ok {
		return nil, fmt.Errorf("failed to declare the type of fetched MFA config")
	}

	return mConfig.Clone()
}

func (b *LoginMFABackend) MemDBMFAConfigByName(ctx context.Context, name string) (*mfa.Config, error) {
	if name == "" {
		return nil, fmt.Errorf("missing config name")
	}

	txn := b.db.Txn(false)

	return b.MemDBMFAConfigByNameInTxn(ctx, txn, name)
}

func (b *LoginMFABackend) MemDBMFALoginEnforcementConfigByNameAndNamespace(name, namespaceId string) (*mfa.MFAEnforcementConfig, error) {
	if name == "" {
		return nil, fmt.Errorf("missing config name")
	}

	txn := b.db.Txn(false)
	defer txn.Abort()

	eConfigRaw, err := txn.First(memDBMFALoginEnforcementsTable, "nameAndNamespace", name, namespaceId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MFA login enforcement config from memdb using name: %w", err)
	}

	if eConfigRaw == nil {
		return nil, nil
	}

	eConfig, ok := eConfigRaw.(*mfa.MFAEnforcementConfig)
	if !ok {
		return nil, fmt.Errorf("invalid type for MFA login enforcement config in memdb")
	}

	return eConfig.Clone()
}

func (b *LoginMFABackend) MemDBMFALoginEnforcementConfigByID(id string) (*mfa.MFAEnforcementConfig, error) {
	if id == "" {
		return nil, fmt.Errorf("missing config id")
	}

	txn := b.db.Txn(false)
	defer txn.Abort()

	eConfigRaw, err := txn.First(memDBMFALoginEnforcementsTable, "id", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MFA login enforcement config from memdb using id: %w", err)
	}

	if eConfigRaw == nil {
		return nil, nil
	}

	eConfig, ok := eConfigRaw.(*mfa.MFAEnforcementConfig)
	if !ok {
		return nil, fmt.Errorf("invalid type for MFA login enforcement config in memdb")
	}

	return eConfig.Clone()
}

func (b *LoginMFABackend) MemDBMFALoginEnforcementConfigIterator() (memdb.ResultIterator, error) {
	txn := b.db.Txn(false)
	defer txn.Abort()

	// List all the MFAEnforcementConfigs
	it, err := txn.Get(memDBMFALoginEnforcementsTable, "id")
	if err != nil {
		return nil, fmt.Errorf("failed to get an iterator over the MFAEnforcementConfig table: %w", err)
	}

	return it, nil
}

func (b *LoginMFABackend) deleteMFALoginEnforcementConfigByNameAndNamespace(ctx context.Context, name, namespaceId string) error {
	var err error

	if name == "" {
		return fmt.Errorf("missing config name")
	}

	b.mfaLock.Lock()
	defer b.mfaLock.Unlock()

	// delete the config from storage
	eConfig, err := b.MemDBMFALoginEnforcementConfigByNameAndNamespace(name, namespaceId)
	if err != nil {
		return err
	}

	if eConfig == nil {
		return nil
	}

	entryIndex := mfaLoginEnforcementPrefix + eConfig.ID
	barrierView, err := b.Core.barrierViewForNamespace(eConfig.NamespaceID)
	if err != nil {
		return err
	}

	err = barrierView.Delete(ctx, entryIndex)
	if err != nil {
		return err
	}

	// create a memdb transaction to delete config
	txn := b.db.Txn(true)
	defer txn.Abort()

	err = txn.Delete(memDBMFALoginEnforcementsTable, eConfig)
	if err != nil {
		return fmt.Errorf("failed to delete MFA login enforcement config from memdb: %w", err)
	}

	txn.Commit()
	return nil
}

func (b *LoginMFABackend) MemDBDeleteMFALoginEnforcementConfigByID(id string) error {
	if id == "" {
		return nil
	}

	txn := b.db.Txn(true)
	defer txn.Abort()

	eConfig, err := b.MemDBMFALoginEnforcementConfigByID(id)
	if err != nil {
		return err
	}

	if eConfig == nil {
		return nil
	}

	err = txn.Delete(memDBMFALoginEnforcementsTable, eConfig)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (b *LoginMFABackend) MemDBDeleteMFALoginEnforcementConfigByNameAndNamespace(name, namespaceId, tableName string) error {
	if name == "" || namespaceId == "" {
		return nil
	}

	txn := b.db.Txn(true)
	defer txn.Abort()

	eConfig, err := b.MemDBMFALoginEnforcementConfigByNameAndNamespace(name, namespaceId)
	if err != nil {
		return err
	}
	if eConfig == nil {
		return nil
	}

	err = txn.Delete(memDBMFALoginEnforcementsTable, eConfig)
	if err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (b *LoginMFABackend) deleteMFAConfigByMethodID(ctx context.Context, configID, methodType, tableName, prefix string) error {
	var err error

	if configID == "" {
		return fmt.Errorf("missing config id")
	}

	b.mfaLock.Lock()
	defer b.mfaLock.Unlock()

	eConfigIter, err := b.MemDBMFALoginEnforcementConfigIterator()
	if err != nil {
		return err
	}

	for eConfigRaw := eConfigIter.Next(); eConfigRaw != nil; eConfigRaw = eConfigIter.Next() {
		eConfig := eConfigRaw.(*mfa.MFAEnforcementConfig)
		if strutil.StrListContains(eConfig.MFAMethodIDs, configID) {
			return fmt.Errorf("methodID is still used by an enforcement configuration with ID: %s", eConfig.ID)
		}
	}

	// Delete the config from storage
	entryIndex := prefix + configID
	err = b.Core.systemBarrierView.Delete(ctx, entryIndex)
	if err != nil {
		return err
	}

	// Create a MemDB transaction to delete config
	txn := b.db.Txn(true)
	defer txn.Abort()

	mConfig, err := b.MemDBMFAConfigByIDInTxn(txn, configID)
	if err != nil {
		return err
	}

	if mConfig == nil {
		return nil
	}

	if mConfig.Type != methodType {
		return fmt.Errorf("method type does not match the MFA config type")
	}

	mfaNs, err := b.Core.NamespaceByID(ctx, mConfig.NamespaceID)
	if err != nil {
		return err
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	// this logic assumes that the config namespace and the current
	// namespace should be the same. Note an ancestor of mfaNs is not allowed
	// to delete methodID
	if ns.ID != mfaNs.ID {
		return fmt.Errorf("request namespace does not match method namespace")
	}

	if mConfig.Type == "totp" && mConfig.ID != "" {
		// This is best effort; if they end up hanging around it's okay, they're encrypted anyways
		if err := logical.ClearView(ctx, NewBarrierView(b.Core.barrier, fmt.Sprintf("%s%s", mfaTOTPKeysPrefix, mConfig.ID))); err != nil {
			b.mfaLogger.Warn("unable to clear TOTP keys", "method", mConfig.Name, "error", err)
		}
	}

	// Delete the config from MemDB
	err = b.MemDBDeleteMFAConfigByIDInTxn(txn, configID)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (b *LoginMFABackend) MemDBDeleteMFAConfigByID(methodId, tableName string) error {
	if methodId == "" {
		return nil
	}

	txn := b.db.Txn(true)
	defer txn.Abort()

	err := b.MemDBDeleteMFAConfigByIDInTxn(txn, methodId)
	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}

func (b *LoginMFABackend) MemDBDeleteMFAConfigByIDInTxn(txn *memdb.Txn, configID string) error {
	if configID == "" {
		return nil
	}

	if txn == nil {
		return fmt.Errorf("txn is nil")
	}

	mConfig, err := b.MemDBMFAConfigByIDInTxn(txn, configID)
	if err != nil {
		return err
	}

	if mConfig == nil {
		return nil
	}

	err = txn.Delete(b.methodTable, mConfig)
	if err != nil {
		return errwrap.Wrapf("failed to delete MFA config from memdb: {{err}}", err)
	}

	return nil
}

func (b *LoginMFABackend) putMFAConfigByID(ctx context.Context, mConfig *mfa.Config) error {
	barrierView, err := b.Core.barrierViewForNamespace(mConfig.NamespaceID)
	if err != nil {
		return err
	}
	return b.putMFAConfigCommon(ctx, mConfig, loginMFAConfigPrefix, mConfig.ID, barrierView)
}

func (b *MFABackend) putMFAConfigCommon(ctx context.Context, mConfig *mfa.Config, prefix, suffix string, barrierView *BarrierView) error {
	entryIndex := prefix + suffix
	marshaledEntry, err := proto.Marshal(mConfig)
	if err != nil {
		return err
	}

	return barrierView.Put(ctx, &logical.StorageEntry{
		Key:   entryIndex,
		Value: marshaledEntry,
	})
}

func (b *MFABackend) getMFAConfig(ctx context.Context, path string, barrierView *BarrierView) (*mfa.Config, error) {
	entry, err := barrierView.Get(ctx, path)
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

func (b *LoginMFABackend) getMFALoginEnforcementConfig(ctx context.Context, path string, barrierView *BarrierView) (*mfa.MFAEnforcementConfig, error) {
	entry, err := barrierView.Get(ctx, path)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var mConfig mfa.MFAEnforcementConfig
	err = proto.Unmarshal(entry.Value, &mConfig)
	if err != nil {
		return nil, err
	}

	return &mConfig, nil
}

func (b *LoginMFABackend) putMFALoginEnforcementConfig(ctx context.Context, eConfig *mfa.MFAEnforcementConfig) error {
	entryIndex := mfaLoginEnforcementPrefix + eConfig.ID
	marshaledEntry, err := proto.Marshal(eConfig)
	if err != nil {
		return err
	}

	barrierView, err := b.Core.barrierViewForNamespace(eConfig.NamespaceID)
	if err != nil {
		return err
	}

	return barrierView.Put(ctx, &logical.StorageEntry{
		Key:   entryIndex,
		Value: marshaledEntry,
	})
}

var mfaHelp = map[string][2]string{
	"methods-list": {
		"Lists all the available MFA methods by their name.",
		"",
	},
	"totp-generate": {
		`Generates a TOTP secret for the given method name on the entity of the
		calling token.`,
		`This endpoint generates an MFA secret based on the
		configuration tied to the method name and stores it in the entity of
		the token making this request.`,
	},
	"totp-admin-generate": {
		`Generates a TOTP secret for the given method name on the given entity.`,
		`This endpoint generates an MFA secret based on the configuration tied
		to the method name and stores it in the entity corresponding to the
		given entity identifier. This endpoint is used to administratively
		generate TOTP secrets on entities.`,
	},
	"totp-admin-destroy": {
		`Deletes the TOTP secret for the given method name on the given entity.`,
		`This endpoint removes the secret belonging to method name from the
		entity regardless of the secret type.`,
	},
	"totp-method": {
		"Defines or updates a TOTP MFA method.",
		"",
	},
	"okta-method": {
		"Defines or updates an Okta MFA method.",
		"",
	},
	"duo-method": {
		"Defines or updates a Duo MFA method.",
		"",
	},
	"pingid-method": {
		"Defines or updates a PingID MFA method.",
		"",
	},
}
