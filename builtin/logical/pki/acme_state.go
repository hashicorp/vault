// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-secure-stdlib/nonceutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// How many bytes are in a token. Per RFC 8555 Section
	// 8.3. HTTP Challenge and Section 11.3 Token Entropy:
	//
	// > token (required, string):  A random value that uniquely identifies
	// >   the challenge.  This value MUST have at least 128 bits of entropy.
	tokenBytes = 128 / 8

	// Path Prefixes
	acmePathPrefix       = "acme/"
	acmeAccountPrefix    = acmePathPrefix + "accounts/"
	acmeThumbprintPrefix = acmePathPrefix + "account-thumbprints/"
	acmeValidationPrefix = acmePathPrefix + "validations/"
	acmeEabPrefix        = acmePathPrefix + "eab/"
)

type acmeState struct {
	nonces nonceutil.NonceService

	validator *ACMEChallengeEngine

	configDirty *atomic.Bool
	_config     sync.RWMutex
	config      acmeConfigEntry
}

type acmeThumbprint struct {
	Kid        string `json:"kid"`
	Thumbprint string `json:"-"`
}

func NewACMEState() *acmeState {
	state := &acmeState{
		nonces:      nonceutil.NewNonceService(),
		validator:   NewACMEChallengeEngine(),
		configDirty: new(atomic.Bool),
	}
	// Config hasn't been loaded yet; mark dirty.
	state.configDirty.Store(true)

	return state
}

func (a *acmeState) Initialize(b *backend, sc *storageContext) error {
	// Initialize the nonce service.
	if err := a.nonces.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize the ACME nonce service: %w", err)
	}

	// Load the ACME config.
	_, err := a.getConfigWithUpdate(sc)
	if err != nil {
		return fmt.Errorf("error initializing ACME engine: %w", err)
	}

	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary | consts.ReplicationPerformanceStandby) {
		// It is assumed, that if the node does become the active node later
		// the plugin is re-initialized, so this is safe. It also spares the node
		// from loading the existing queue into memory for no reason.
		b.Logger().Debug("Not on an active node, skipping starting ACME challenge validation engine")
		return nil
	}
	// Kick off our ACME challenge validation engine.
	go a.validator.Run(b, a, sc)

	// All good.
	return nil
}

func (a *acmeState) Shutdown(b *backend) {
	// If we aren't the active node, nothing to shutdown
	if b.System().ReplicationState().HasState(consts.ReplicationDRSecondary | consts.ReplicationPerformanceStandby) {
		return
	}

	a.validator.Closing <- struct{}{}
}

func (a *acmeState) markConfigDirty() {
	a.configDirty.Store(true)
}

func (a *acmeState) reloadConfigIfRequired(sc *storageContext) error {
	if !a.configDirty.Load() {
		return nil
	}

	a._config.Lock()
	defer a._config.Unlock()

	if !a.configDirty.Load() {
		// Someone beat us to grabbing the above write lock and already
		// updated the config.
		return nil
	}

	config, err := getAcmeConfig(sc)
	if err != nil {
		return fmt.Errorf("failed reading ACME config: %w", err)
	}

	a.config = *config
	a.configDirty.Store(false)

	return nil
}

func (a *acmeState) getConfigWithUpdate(sc *storageContext) (*acmeConfigEntry, error) {
	if err := a.reloadConfigIfRequired(sc); err != nil {
		return nil, err
	}

	a._config.RLock()
	defer a._config.RUnlock()

	configCopy := a.config
	return &configCopy, nil
}

func (a *acmeState) getConfigWithForcedUpdate(sc *storageContext) (*acmeConfigEntry, error) {
	a.markConfigDirty()
	return a.getConfigWithUpdate(sc)
}

func (a *acmeState) writeConfig(sc *storageContext, config *acmeConfigEntry) (*acmeConfigEntry, error) {
	a._config.Lock()
	defer a._config.Unlock()

	if err := sc.setAcmeConfig(config); err != nil {
		a.markConfigDirty()
		return nil, fmt.Errorf("failed writing ACME config: %w", err)
	}

	if config != nil {
		a.config = *config
	} else {
		a.config = defaultAcmeConfig
	}

	return config, nil
}

func generateRandomBase64(srcBytes int) (string, error) {
	data := make([]byte, 21)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (a *acmeState) GetNonce() (string, time.Time, error) {
	return a.nonces.Get()
}

func (a *acmeState) RedeemNonce(nonce string) bool {
	return a.nonces.Redeem(nonce)
}

func (a *acmeState) DoTidyNonces() {
	a.nonces.Tidy()
}

type ACMEAccountStatus string

func (aas ACMEAccountStatus) String() string {
	return string(aas)
}

const (
	AccountStatusValid       ACMEAccountStatus = "valid"
	AccountStatusDeactivated ACMEAccountStatus = "deactivated"
	AccountStatusRevoked     ACMEAccountStatus = "revoked"
)

type acmeAccount struct {
	KeyId                string            `json:"-"`
	Status               ACMEAccountStatus `json:"status"`
	Contact              []string          `json:"contact"`
	TermsOfServiceAgreed bool              `json:"terms-of-service-agreed"`
	Jwk                  []byte            `json:"jwk"`
	AcmeDirectory        string            `json:"acme-directory"`
	AccountCreatedDate   time.Time         `json:"account-created-date"`
	MaxCertExpiry        time.Time         `json:"account-max-cert-expiry"`
	AccountRevokedDate   time.Time         `json:"account-revoked-date"`
	Eab                  *eabType          `json:"eab"`
}

type acmeOrder struct {
	OrderId                 string              `json:"-"`
	AccountId               string              `json:"account-id"`
	Status                  ACMEOrderStatusType `json:"status"`
	Expires                 time.Time           `json:"expires"`
	Identifiers             []*ACMEIdentifier   `json:"identifiers"`
	AuthorizationIds        []string            `json:"authorization-ids"`
	CertificateSerialNumber string              `json:"cert-serial-number"`
	CertificateExpiry       time.Time           `json:"cert-expiry"`
	// The actual issuer UUID that issued the certificate, blank if an order exists but no certificate was issued.
	IssuerId issuing.IssuerID `json:"issuer-id"`
}

func (o acmeOrder) getIdentifierDNSValues() []string {
	var identifiers []string
	for _, value := range o.Identifiers {
		if value.Type == ACMEDNSIdentifier {
			// Here, because of wildcard processing, we need to use the
			// original value provided by the caller rather than the
			// post-modification (trimmed '*.' prefix) value.
			identifiers = append(identifiers, value.OriginalValue)
		}
	}
	return identifiers
}

func (o acmeOrder) getIdentifierIPValues() []net.IP {
	var identifiers []net.IP
	for _, value := range o.Identifiers {
		if value.Type == ACMEIPIdentifier {
			identifiers = append(identifiers, net.ParseIP(value.Value))
		}
	}
	return identifiers
}

func (a *acmeState) CreateAccount(ac *acmeContext, c *jwsCtx, contact []string, termsOfServiceAgreed bool, eab *eabType) (*acmeAccount, error) {
	// Write out the thumbprint value/entry out first, if we get an error mid-way through
	// this is easier to recover from. The new kid with the same existing public key
	// will rewrite the thumbprint entry. This goes in hand with LoadAccountByKey that
	// will return a nil, nil value if the referenced kid in a loaded thumbprint does not
	// exist. This effectively makes this self-healing IF the end-user re-attempts the
	// account creation with the same public key.
	thumbprint, err := c.GetKeyThumbprint()
	if err != nil {
		return nil, fmt.Errorf("failed generating thumbprint: %w", err)
	}

	thumbPrint := &acmeThumbprint{
		Kid:        c.Kid,
		Thumbprint: thumbprint,
	}
	thumbPrintEntry, err := logical.StorageEntryJSON(acmeThumbprintPrefix+thumbprint, thumbPrint)
	if err != nil {
		return nil, fmt.Errorf("error generating account thumbprint entry: %w", err)
	}

	if err = ac.sc.Storage.Put(ac.sc.Context, thumbPrintEntry); err != nil {
		return nil, fmt.Errorf("error writing account thumbprint entry: %w", err)
	}

	// Now write out the main value that the thumbprint points too.
	acct := &acmeAccount{
		KeyId:                c.Kid,
		Contact:              contact,
		TermsOfServiceAgreed: termsOfServiceAgreed,
		Jwk:                  c.Jwk,
		Status:               AccountStatusValid,
		AcmeDirectory:        ac.acmeDirectory,
		AccountCreatedDate:   time.Now(),
		Eab:                  eab,
	}
	json, err := logical.StorageEntryJSON(acmeAccountPrefix+c.Kid, acct)
	if err != nil {
		return nil, fmt.Errorf("error creating account entry: %w", err)
	}

	if err := ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return nil, fmt.Errorf("error writing account entry: %w", err)
	}

	return acct, nil
}

func (a *acmeState) UpdateAccount(sc *storageContext, acct *acmeAccount) error {
	json, err := logical.StorageEntryJSON(acmeAccountPrefix+acct.KeyId, acct)
	if err != nil {
		return fmt.Errorf("error creating account entry: %w", err)
	}

	if err := sc.Storage.Put(sc.Context, json); err != nil {
		return fmt.Errorf("error writing account entry: %w", err)
	}

	return nil
}

// LoadAccount will load the account object based on the passed in keyId field value
// otherwise will return an error if the account does not exist.
func (a *acmeState) LoadAccount(ac *acmeContext, keyId string) (*acmeAccount, error) {
	acct, err := a.LoadAccountWithoutDirEnforcement(ac.sc, keyId)
	if err != nil {
		return acct, err
	}

	if acct.AcmeDirectory != ac.acmeDirectory {
		return nil, fmt.Errorf("%w: account part of different ACME directory path", ErrMalformed)
	}

	return acct, nil
}

// LoadAccountWithoutDirEnforcement will load the account object based on the passed in keyId field value,
// but does not enforce the ACME directory path, normally this is used by non ACME specific APIs.
func (a *acmeState) LoadAccountWithoutDirEnforcement(sc *storageContext, keyId string) (*acmeAccount, error) {
	entry, err := sc.Storage.Get(sc.Context, acmeAccountPrefix+keyId)
	if err != nil {
		return nil, fmt.Errorf("error loading account: %w", err)
	}
	if entry == nil {
		return nil, fmt.Errorf("account not found: %w", ErrAccountDoesNotExist)
	}

	var acct acmeAccount
	err = entry.DecodeJSON(&acct)
	if err != nil {
		return nil, fmt.Errorf("error decoding account: %w", err)
	}
	acct.KeyId = keyId
	return &acct, nil
}

// LoadAccountByKey will attempt to load the account based on a key thumbprint. If the thumbprint
// or kid is unknown a nil, nil will be returned.
func (a *acmeState) LoadAccountByKey(ac *acmeContext, keyThumbprint string) (*acmeAccount, error) {
	thumbprintEntry, err := ac.sc.Storage.Get(ac.sc.Context, acmeThumbprintPrefix+keyThumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed loading acme thumbprintEntry for key: %w", err)
	}
	if thumbprintEntry == nil {
		return nil, nil
	}

	var thumbprint acmeThumbprint
	err = thumbprintEntry.DecodeJSON(&thumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed decoding thumbprint entry: %s: %w", keyThumbprint, err)
	}

	if len(thumbprint.Kid) == 0 {
		return nil, fmt.Errorf("empty kid within thumbprint entry: %s", keyThumbprint)
	}

	acct, err := a.LoadAccount(ac, thumbprint.Kid)
	if err != nil {
		// If we fail to lookup the account that the thumbprint entry references, assume a bad
		// write previously occurred in which we managed to write out the thumbprint but failed
		// writing out the main account information.
		if errors.Is(err, ErrAccountDoesNotExist) {
			return nil, nil
		}
		return nil, err
	}

	return acct, nil
}

func (a *acmeState) LoadJWK(ac *acmeContext, keyId string) ([]byte, error) {
	key, err := a.LoadAccount(ac, keyId)
	if err != nil {
		return nil, err
	}

	if len(key.Jwk) == 0 {
		return nil, fmt.Errorf("malformed key entry lacks JWK")
	}

	return key.Jwk, nil
}

func (a *acmeState) LoadAuthorization(ac *acmeContext, userCtx *jwsCtx, authId string) (*ACMEAuthorization, error) {
	if authId == "" {
		return nil, fmt.Errorf("malformed authorization identifier")
	}

	authorizationPath := getAuthorizationPath(userCtx.Kid, authId)

	authz, err := loadAuthorizationAtPath(ac.sc, authorizationPath)
	if err != nil {
		return nil, err
	}

	if userCtx.Kid != authz.AccountId {
		return nil, ErrUnauthorized
	}

	return authz, nil
}

func loadAuthorizationAtPath(sc *storageContext, authorizationPath string) (*ACMEAuthorization, error) {
	entry, err := sc.Storage.Get(sc.Context, authorizationPath)
	if err != nil {
		return nil, fmt.Errorf("error loading authorization: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("authorization does not exist: %w", ErrMalformed)
	}

	var authz ACMEAuthorization
	err = entry.DecodeJSON(&authz)
	if err != nil {
		return nil, fmt.Errorf("error decoding authorization: %w", err)
	}

	return &authz, nil
}

func (a *acmeState) SaveAuthorization(ac *acmeContext, authz *ACMEAuthorization) error {
	path := getAuthorizationPath(authz.AccountId, authz.Id)
	return saveAuthorizationAtPath(ac.sc, path, authz)
}

func saveAuthorizationAtPath(sc *storageContext, path string, authz *ACMEAuthorization) error {
	if authz.Id == "" {
		return fmt.Errorf("invalid authorization, missing id")
	}

	if authz.AccountId == "" {
		return fmt.Errorf("invalid authorization, missing account id")
	}

	json, err := logical.StorageEntryJSON(path, authz)
	if err != nil {
		return fmt.Errorf("error creating authorization entry: %w", err)
	}

	if err = sc.Storage.Put(sc.Context, json); err != nil {
		return fmt.Errorf("error writing authorization entry: %w", err)
	}

	return nil
}

func (a *acmeState) ParseRequestParams(ac *acmeContext, req *logical.Request, data *framework.FieldData) (*jwsCtx, map[string]interface{}, error) {
	var c jwsCtx
	var m map[string]interface{}

	// Parse the key out.
	rawJWKBase64, ok := data.GetOk("protected")
	if !ok {
		return nil, nil, fmt.Errorf("missing required field 'protected': %w", ErrMalformed)
	}
	jwkBase64 := rawJWKBase64.(string)

	jwkBytes, err := base64.RawURLEncoding.DecodeString(jwkBase64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to base64 parse 'protected': %s: %w", err, ErrMalformed)
	}
	if err = c.UnmarshalOuterJwsJson(a, ac, jwkBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to json unmarshal 'protected': %w", err)
	}

	// Since we already parsed the header to verify the JWS context, we
	// should read and redeem the nonce here too, to avoid doing any extra
	// work if it is invalid.
	if !a.RedeemNonce(c.Nonce) {
		return nil, nil, fmt.Errorf("invalid or reused nonce: %w", ErrBadNonce)
	}

	// If the path is incorrect, reject the request.
	//
	// See RFC 8555 Section 6.4. Request URL Integrity:
	//
	// > As noted in Section 6.2, all ACME request objects carry a "url"
	// > header parameter in their protected header. ... On receiving such
	// > an object in an HTTP request, the server MUST compare the "url"
	// > header parameter to the request URL.  If the two do not match,
	// > then the server MUST reject the request as unauthorized.
	if len(c.Url) == 0 {
		return nil, nil, fmt.Errorf("missing required parameter 'url' in 'protected': %w", ErrMalformed)
	}
	if ac.clusterUrl.JoinPath(req.Path).String() != c.Url {
		return nil, nil, fmt.Errorf("invalid value for 'url' in 'protected': got '%v' expected '%v': %w", c.Url, ac.clusterUrl.JoinPath(req.Path).String(), ErrUnauthorized)
	}

	rawPayloadBase64, ok := data.GetOk("payload")
	if !ok {
		return nil, nil, fmt.Errorf("missing required field 'payload': %w", ErrMalformed)
	}
	payloadBase64 := rawPayloadBase64.(string)

	rawSignatureBase64, ok := data.GetOk("signature")
	if !ok {
		return nil, nil, fmt.Errorf("missing required field 'signature': %w", ErrMalformed)
	}
	signatureBase64 := rawSignatureBase64.(string)

	// go-jose only seems to support compact signature encodings.
	compactSig := fmt.Sprintf("%v.%v.%v", jwkBase64, payloadBase64, signatureBase64)
	m, err = c.VerifyJWS(compactSig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to verify signature: %w", err)
	}

	return &c, m, nil
}

func (a *acmeState) LoadOrder(ac *acmeContext, userCtx *jwsCtx, orderId string) (*acmeOrder, error) {
	path := getOrderPath(userCtx.Kid, orderId)
	entry, err := ac.sc.Storage.Get(ac.sc.Context, path)
	if err != nil {
		return nil, fmt.Errorf("error loading order: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("order does not exist: %w", ErrMalformed)
	}

	var order acmeOrder
	err = entry.DecodeJSON(&order)
	if err != nil {
		return nil, fmt.Errorf("error decoding order: %w", err)
	}

	if userCtx.Kid != order.AccountId {
		return nil, ErrUnauthorized
	}

	order.OrderId = orderId

	return &order, nil
}

// LoadAccountOrders will load all orders for a given account ID, this should be used by the
// management interface only, not through any of the ACME APIs.
func (a *acmeState) LoadAccountOrders(sc *storageContext, accountId string) ([]*acmeOrder, error) {
	orderIds, err := a.ListOrderIds(sc, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed listing order ids for account id %s: %w", accountId, err)
	}

	var orders []*acmeOrder
	for _, orderId := range orderIds {
		order, err := a.LoadOrder(&acmeContext{sc: sc}, &jwsCtx{Kid: accountId}, orderId)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (a *acmeState) SaveOrder(ac *acmeContext, order *acmeOrder) error {
	if order.OrderId == "" {
		return fmt.Errorf("invalid order, missing order id")
	}

	if order.AccountId == "" {
		return fmt.Errorf("invalid order, missing account id")
	}
	path := getOrderPath(order.AccountId, order.OrderId)
	json, err := logical.StorageEntryJSON(path, order)
	if err != nil {
		return fmt.Errorf("error serializing order entry: %w", err)
	}

	if err = ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing order entry: %w", err)
	}

	return nil
}

func (a *acmeState) ListOrderIds(sc *storageContext, accountId string) ([]string, error) {
	accountOrderPrefixPath := acmeAccountPrefix + accountId + "/orders/"

	rawOrderIds, err := sc.Storage.List(sc.Context, accountOrderPrefixPath)
	if err != nil {
		return nil, fmt.Errorf("failed listing order ids for account %s: %w", accountId, err)
	}

	return filterDirEntries(rawOrderIds), nil
}

type acmeCertEntry struct {
	Serial  string `json:"-"`
	Account string `json:"-"`
	Order   string `json:"order"`
}

func (a *acmeState) TrackIssuedCert(ac *acmeContext, accountId string, serial string, orderId string) error {
	path := getAcmeSerialToAccountTrackerPath(accountId, serial)
	entry := acmeCertEntry{
		Order: orderId,
	}

	json, err := logical.StorageEntryJSON(path, &entry)
	if err != nil {
		return fmt.Errorf("error serializing acme cert entry: %w", err)
	}

	if err = ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing acme cert entry: %w", err)
	}

	return nil
}

func (a *acmeState) GetIssuedCert(ac *acmeContext, accountId string, serial string) (*acmeCertEntry, error) {
	path := acmeAccountPrefix + accountId + "/certs/" + normalizeSerial(serial)

	entry, err := ac.sc.Storage.Get(ac.sc.Context, path)
	if err != nil {
		return nil, fmt.Errorf("error loading acme cert entry: %w", err)
	}

	if entry == nil {
		return nil, fmt.Errorf("no certificate with this serial was issued for this account")
	}

	var cert acmeCertEntry
	err = entry.DecodeJSON(&cert)
	if err != nil {
		return nil, fmt.Errorf("error decoding acme cert entry: %w", err)
	}

	cert.Serial = denormalizeSerial(serial)
	cert.Account = accountId

	return &cert, nil
}

func (a *acmeState) SaveEab(sc *storageContext, eab *eabType) error {
	json, err := logical.StorageEntryJSON(path.Join(acmeEabPrefix, eab.KeyID), eab)
	if err != nil {
		return err
	}
	return sc.Storage.Put(sc.Context, json)
}

func (a *acmeState) LoadEab(sc *storageContext, eabKid string) (*eabType, error) {
	rawEntry, err := sc.Storage.Get(sc.Context, path.Join(acmeEabPrefix, eabKid))
	if err != nil {
		return nil, err
	}
	if rawEntry == nil {
		return nil, fmt.Errorf("%w: no eab found for kid %s", ErrStorageItemNotFound, eabKid)
	}

	var eab eabType
	err = rawEntry.DecodeJSON(&eab)
	if err != nil {
		return nil, err
	}

	eab.KeyID = eabKid
	return &eab, nil
}

func (a *acmeState) DeleteEab(sc *storageContext, eabKid string) (bool, error) {
	rawEntry, err := sc.Storage.Get(sc.Context, path.Join(acmeEabPrefix, eabKid))
	if err != nil {
		return false, err
	}
	if rawEntry == nil {
		return false, nil
	}

	err = sc.Storage.Delete(sc.Context, path.Join(acmeEabPrefix, eabKid))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *acmeState) ListEabIds(sc *storageContext) ([]string, error) {
	entries, err := sc.Storage.List(sc.Context, acmeEabPrefix)
	if err != nil {
		return nil, err
	}
	ids := filterDirEntries(entries)

	return ids, nil
}

func (a *acmeState) ListAccountIds(sc *storageContext) ([]string, error) {
	entries, err := sc.Storage.List(sc.Context, acmeAccountPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed listing ACME account prefix directory %s: %w", acmeAccountPrefix, err)
	}

	return filterDirEntries(entries), nil
}

func getAcmeSerialToAccountTrackerPath(accountId string, serial string) string {
	return acmeAccountPrefix + accountId + "/certs/" + normalizeSerial(serial)
}

func getAuthorizationPath(accountId string, authId string) string {
	return acmeAccountPrefix + accountId + "/authorizations/" + authId
}

func getOrderPath(accountId string, orderId string) string {
	return acmeAccountPrefix + accountId + "/orders/" + orderId
}

func getACMEToken() (string, error) {
	return generateRandomBase64(tokenBytes)
}
