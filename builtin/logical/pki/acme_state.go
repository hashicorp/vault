package pki

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// How long nonces are considered valid.
	nonceExpiry = 15 * time.Minute

	// Path Prefixes // TODO, these are the things to tidy
	acmePathPrefix          = "acme/"
	acmeAccountPrefix       = acmePathPrefix + "accounts/"
	acmeThumbprintPrefix    = acmePathPrefix + "account-thumbprints/"
	acmeAuthroizationPrefix = acmePathPrefix + "authz/" // TODO: These contain the challenges, which will need to be cleaned up by tidy (need to wait on orders)
)

type acmeState struct {
	nextExpiry *atomic.Int64
	nonces     *sync.Map // map[string]time.Time
}

type acmeThumbprint struct {
	Kid        string `json:"kid"`
	Thumbprint string `json:"-"`
}

func NewACMEState() *acmeState {
	return &acmeState{
		nextExpiry: new(atomic.Int64),
		nonces:     new(sync.Map),
	}
}

func generateNonce() (string, error) {
	data := make([]byte, 21)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func (a *acmeState) GetNonce() (string, time.Time, error) {
	now := time.Now()
	nonce, err := generateNonce()
	if err != nil {
		return "", now, err
	}

	then := now.Add(nonceExpiry)
	a.nonces.Store(nonce, then)

	nextExpiry := a.nextExpiry.Load()
	next := time.Unix(nextExpiry, 0)
	if now.After(next) || then.Before(next) {
		a.nextExpiry.Store(then.Unix())
	}

	return nonce, then, nil
}

func (a *acmeState) RedeemNonce(nonce string) bool {
	rawTimeout, present := a.nonces.LoadAndDelete(nonce)
	if !present {
		return false
	}

	timeout := rawTimeout.(time.Time)
	if time.Now().After(timeout) {
		return false
	}

	return true
}

func (a *acmeState) DoTidyNonces() {
	now := time.Now()
	expiry := a.nextExpiry.Load()
	then := time.Unix(expiry, 0)

	if expiry == 0 || now.After(then) {
		a.TidyNonces()
	}
}

func (a *acmeState) TidyNonces() {
	now := time.Now()
	nextRun := now.Add(nonceExpiry)

	a.nonces.Range(func(key, value any) bool {
		timeout := value.(time.Time)
		if now.After(timeout) {
			a.nonces.Delete(key)
		}

		if timeout.Before(nextRun) {
			nextRun = timeout
		}

		return false /* don't quit looping */
	})

	a.nextExpiry.Store(nextRun.Unix())
}

type ACMEAccountStatus string

func (aas ACMEAccountStatus) String() string {
	return string(aas)
}

const (
	StatusValid       ACMEAccountStatus = "valid"
	StatusDeactivated ACMEAccountStatus = "deactivated"
	StatusRevoked     ACMEAccountStatus = "revoked"
)

type acmeAccount struct {
	KeyId                string            `json:"-"`
	Status               ACMEAccountStatus `json:"status"`
	Contact              []string          `json:"contact"`
	TermsOfServiceAgreed bool              `json:"termsOfServiceAgreed"`
	Jwk                  []byte            `json:"jwk"`
	LastAccessDate       time.Time         `json:"last_access_date"`
	MarkedAsRevokedDate  *time.Time        `json:"marked_as_revoked_date"`
}

func (a *acmeState) CreateAccount(ac *acmeContext, c *jwsCtx, contact []string, termsOfServiceAgreed bool) (*acmeAccount, error) {
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
		Status:               StatusValid,
		LastAccessDate:       time.Now(),
		MarkedAsRevokedDate:  nil,
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

func (a *acmeState) UpdateAccount(ac *acmeContext, acct *acmeAccount) error {
	json, err := logical.StorageEntryJSON(acmeAccountPrefix+acct.KeyId, acct)
	if err != nil {
		return fmt.Errorf("error creating account entry: %w", err)
	}

	if err := ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing account entry: %w", err)
	}

	return nil
}

// LoadAccount will load the account object based on the passed in keyId field value
// otherwise will return an error if the account does not exist.
func (a *acmeState) LoadAccount(ac *acmeContext, keyId string) (*acmeAccount, error) {
	entry, err := ac.sc.Storage.Get(ac.sc.Context, acmeAccountPrefix+keyId)
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

func (a *acmeState) DeleteAccountByThumbprint(ac *acmeContext, keyThumbprint string) error {
	thumbprintEntry, err := ac.sc.Storage.Get(ac.sc.Context, acmeThumbprintPrefix+keyThumbprint)
	if err != nil {
		return fmt.Errorf("error retrieving thumbprint entry %v, unable to find corresponding account entry: %w", keyThumbprint, err)
	}
	if thumbprintEntry == nil {
		return fmt.Errorf("empty thumbprint entry %v, unable to find corresponding account entry", keyThumbprint)
	}

	var thumbprint acmeThumbprint
	err = thumbprintEntry.DecodeJSON(&thumbprint)
	if err != nil {
		return fmt.Errorf("unable to decode thumbprint entry %v to find account entry: %w", keyThumbprint, err)
	}

	if len(thumbprint.Kid) == 0 {
		return fmt.Errorf("unable to find account entry: empty kid within thumbprint entry: %s", keyThumbprint)
	}

	// We Delete the Account Associated with this Thumbprint:
	err = ac.sc.Storage.Delete(ac.sc.Context, acmeAccountPrefix+thumbprint.Kid)
	if err != nil {
		return err
	}

	// Now we delete the Thumbprint Associated with the Account:
	err = ac.sc.Storage.Delete(ac.sc.Context, acmeThumbprintPrefix+keyThumbprint)
	if err != nil {
		return err
	}

	return nil
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

func (a *acmeState) TidyAcmeAccount(ac *acmeContext, keyThumbprint string, markRevokedIfLastAccessedBefore time.Time, deleteIfMarkedRevokedBefore time.Time) (wasRevoked bool, wasDeleted bool, err error) {
	acmeAccount, err := a.LoadAccountByKey(ac, keyThumbprint)
	if err != nil { // Possibly Already Deleted
		return false, false, err
	}
	if acmeAccount.MarkedAsRevokedDate != nil && !acmeAccount.MarkedAsRevokedDate.IsZero() && acmeAccount.MarkedAsRevokedDate.Before(deleteIfMarkedRevokedBefore) {
		err := a.DeleteAccountByThumbprint(ac, keyThumbprint)
		if err != nil {
			return false, false, err
		}
		return false, true, nil
	}
	if acmeAccount.LastAccessDate.Before(markRevokedIfLastAccessedBefore) {
		acmeAccount.MarkedAsRevokedDate = new(time.Time)
		*acmeAccount.MarkedAsRevokedDate = time.Now()
		acmeAccount.Status = StatusRevoked
		err := a.UpdateAccount(ac, acmeAccount)
		if err != nil {
			return false, false, err
		}
		return true, false, nil
	}

	return false, false, nil
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

	entry, err := ac.sc.Storage.Get(ac.sc.Context, acmeAuthroizationPrefix+authId)
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

	if userCtx.Kid != authz.AccountId {
		return nil, ErrUnauthorized
	}

	return &authz, nil
}

func (a *acmeState) SaveAuthorization(ac *acmeContext, authz *ACMEAuthorization) error {
	if authz.Id == "" {
		return fmt.Errorf("invalid authorization, missing id")
	}

	json, err := logical.StorageEntryJSON(acmeAuthroizationPrefix+authz.Id, authz)
	if err != nil {
		return fmt.Errorf("error creating authorization entry: %w", err)
	}

	if err := ac.sc.Storage.Put(ac.sc.Context, json); err != nil {
		return fmt.Errorf("error writing authorization entry: %w", err)
	}

	return nil
}

func (a *acmeState) ParseRequestParams(ac *acmeContext, data *framework.FieldData) (*jwsCtx, map[string]interface{}, error) {
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
	if err = c.UnmarshalJSON(a, ac, jwkBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to json unmarshal 'protected': %w", err)
	}

	// Since we already parsed the header to verify the JWS context, we
	// should read and redeem the nonce here too, to avoid doing any extra
	// work if it is invalid.
	if !a.RedeemNonce(c.Nonce) {
		return nil, nil, fmt.Errorf("invalid or reused nonce: %w", ErrBadNonce)
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
