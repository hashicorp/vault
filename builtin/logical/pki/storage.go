package pki

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storageKeyConfig      = "config/keys"
	storageIssuerConfig   = "config/issuers"
	keyPrefix             = "config/key/"
	issuerPrefix          = "config/issuer/"
	storageLocalCRLConfig = "crls/config"

	legacyMigrationBundleLogKey = "config/legacyMigrationBundleLog"
	legacyCertBundlePath        = "config/ca_bundle"
	legacyCRLPath               = "crl"

	// Used as a quick sanity check for a reference id lookups...
	uuidLength = 36
)

type keyID string

func (p keyID) String() string {
	return string(p)
}

type issuerID string

func (p issuerID) String() string {
	return string(p)
}

type crlID string

func (p crlID) String() string {
	return string(p)
}

const (
	IssuerRefNotFound = issuerID("not-found")
	KeyRefNotFound    = keyID("not-found")
)

type keyEntry struct {
	ID             keyID                   `json:"id" structs:"id" mapstructure:"id"`
	Name           string                  `json:"name" structs:"name" mapstructure:"name"`
	PrivateKeyType certutil.PrivateKeyType `json:"private_key_type" structs:"private_key_type" mapstructure:"private_key_type"`
	PrivateKey     string                  `json:"private_key" structs:"private_key" mapstructure:"private_key"`
}

func (e keyEntry) getManagedKeyUUID() (UUIDKey, error) {
	if !e.isManagedPrivateKey() {
		return "", errutil.InternalError{Err: "getManagedKeyId called on a key id %s (%s) "}
	}
	return extractManagedKeyId([]byte(e.PrivateKey))
}

func (e keyEntry) isManagedPrivateKey() bool {
	return e.PrivateKeyType == certutil.ManagedPrivateKey
}

type issuerUsage uint

const (
	ReadOnlyUsage   issuerUsage = iota
	IssuanceUsage   issuerUsage = 1 << iota
	CRLSigningUsage issuerUsage = 1 << iota

	// When adding a new usage in the future, we'll need to create a usage
	// mask field on the IssuerEntry and handle migrations to a newer mask,
	// inferring a value for the new bits.
	AllIssuerUsages issuerUsage = ReadOnlyUsage | IssuanceUsage | CRLSigningUsage
)

var namedIssuerUsages = map[string]issuerUsage{
	"read-only":            ReadOnlyUsage,
	"issuing-certificates": IssuanceUsage,
	"crl-signing":          CRLSigningUsage,
}

func (i *issuerUsage) ToggleUsage(usages ...issuerUsage) {
	for _, usage := range usages {
		*i ^= usage
	}
}

func (i issuerUsage) HasUsage(usage issuerUsage) bool {
	return (i & usage) == usage
}

func (i issuerUsage) Names() string {
	var names []string
	var builtUsage issuerUsage

	for name, usage := range namedIssuerUsages {
		if i.HasUsage(usage) {
			names = append(names, name)
			builtUsage.ToggleUsage(usage)
		}
	}

	if i != builtUsage {
		// Found some unknown usage, we should indicate this in the names.
		names = append(names, fmt.Sprintf("unknown:%v", i^builtUsage))
	}

	return strings.Join(names, ",")
}

func NewIssuerUsageFromNames(names []string) (issuerUsage, error) {
	var result issuerUsage
	for index, name := range names {
		usage, ok := namedIssuerUsages[name]
		if !ok {
			return ReadOnlyUsage, fmt.Errorf("unknown name for usage at index %v: %v", index, name)
		}

		result.ToggleUsage(usage)
	}

	return result, nil
}

type issuerEntry struct {
	ID                   issuerID                  `json:"id" structs:"id" mapstructure:"id"`
	Name                 string                    `json:"name" structs:"name" mapstructure:"name"`
	KeyID                keyID                     `json:"key_id" structs:"key_id" mapstructure:"key_id"`
	Certificate          string                    `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	CAChain              []string                  `json:"ca_chain" structs:"ca_chain" mapstructure:"ca_chain"`
	ManualChain          []issuerID                `json:"manual_chain" structs:"manual_chain" mapstructure:"manual_chain"`
	SerialNumber         string                    `json:"serial_number" structs:"serial_number" mapstructure:"serial_number"`
	LeafNotAfterBehavior certutil.NotAfterBehavior `json:"not_after_behavior" structs:"not_after_behavior" mapstructure:"not_after_behavior"`
	Usage                issuerUsage               `json:"usage" structs:"usage" mapstructure:"usage"`
}

type localCRLConfigEntry struct {
	IssuerIDCRLMap map[issuerID]crlID `json:"issuer_id_crl_map" structs:"issuer_id_crl_map" mapstructure:"issuer_id_crl_map"`
	CRLNumberMap   map[crlID]int64    `json:"crl_number_map" structs:"crl_number_map" mapstructure:"crl_number_map"`
}

type keyConfigEntry struct {
	DefaultKeyId keyID `json:"default" structs:"default" mapstructure:"default"`
}

type issuerConfigEntry struct {
	DefaultIssuerId issuerID `json:"default" structs:"default" mapstructure:"default"`
}

func listKeys(ctx context.Context, s logical.Storage) ([]keyID, error) {
	strList, err := s.List(ctx, keyPrefix)
	if err != nil {
		return nil, err
	}

	keyIds := make([]keyID, 0, len(strList))
	for _, entry := range strList {
		keyIds = append(keyIds, keyID(entry))
	}

	return keyIds, nil
}

func fetchKeyById(ctx context.Context, s logical.Storage, keyId keyID) (*keyEntry, error) {
	if len(keyId) == 0 {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: empty key identifier")}
	}

	entry, err := s.Get(ctx, keyPrefix+keyId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: %v", err)}
	}
	if entry == nil {
		// FIXME: Dedicated/specific error for this?
		return nil, errutil.UserError{Err: fmt.Sprintf("pki key id %s does not exist", keyId.String())}
	}

	var key keyEntry
	if err := entry.DecodeJSON(&key); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki key with id %s: %v", keyId.String(), err)}
	}

	return &key, nil
}

func writeKey(ctx context.Context, s logical.Storage, key keyEntry) error {
	keyId := key.ID

	json, err := logical.StorageEntryJSON(keyPrefix+keyId.String(), key)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func deleteKey(ctx context.Context, s logical.Storage, id keyID) (bool, error) {
	config, err := getKeysConfig(ctx, s)
	if err != nil {
		return false, err
	}

	wasDefault := false
	if config.DefaultKeyId == id {
		wasDefault = true
		config.DefaultKeyId = keyID("")
		if err := setKeysConfig(ctx, s, config); err != nil {
			return wasDefault, err
		}
	}

	return wasDefault, s.Delete(ctx, keyPrefix+id.String())
}

func importKey(ctx context.Context, b *backend, s logical.Storage, keyValue string, keyName string, keyType certutil.PrivateKeyType) (*keyEntry, bool, error) {
	// importKey imports the specified PEM-format key (from keyValue) into
	// the new PKI storage format. The first return field is a reference to
	// the new key; the second is whether or not the key already existed
	// during import (in which case, *key points to the existing key reference
	// and identifier); the last return field is whether or not an error
	// occurred.
	//
	// Normalize whitespace before beginning.  See note in importIssuer as to
	// why we do this.
	keyValue = strings.TrimSpace(keyValue) + "\n"
	//
	// Before we can import a known key, we first need to know if the key
	// exists in storage already. This means iterating through all known
	// keys and comparing their private value against this value.
	knownKeys, err := listKeys(ctx, s)
	if err != nil {
		return nil, false, err
	}

	// Get our public key from the current inbound key, to compare against all the other keys.
	var pkForImportingKey crypto.PublicKey
	if keyType == certutil.ManagedPrivateKey {
		managedKeyUUID, err := extractManagedKeyId([]byte(keyValue))
		if err != nil {
			return nil, false, errutil.InternalError{Err: fmt.Sprintf("failed extracting managed key uuid from key: %v", err)}
		}
		pkForImportingKey, err = getManagedKeyPublicKey(ctx, b, managedKeyUUID)
		if err != nil {
			return nil, false, err
		}
	} else {
		pkForImportingKey, err = getPublicKeyFromBytes([]byte(keyValue))
		if err != nil {
			return nil, false, err
		}
	}

	foundExistingKeyWithName := false
	for _, identifier := range knownKeys {
		existingKey, err := fetchKeyById(ctx, s, identifier)
		if err != nil {
			return nil, false, err
		}
		areEqual, err := comparePublicKey(ctx, b, existingKey, pkForImportingKey)
		if err != nil {
			return nil, false, err
		}

		if areEqual {
			// Here, we don't need to stitch together the issuer entries,
			// because the last run should've done that for us (or, when
			// importing an issuer).
			return existingKey, true, nil
		}

		// Allow us to find an existing matching key with a different name before erroring out
		if keyName != "" && existingKey.Name == keyName {
			foundExistingKeyWithName = true
		}
	}

	// Another key with a different value is using the keyName so reject this request.
	if foundExistingKeyWithName {
		return nil, false, errutil.UserError{Err: fmt.Sprintf("an existing key is using the requested key name value: %s", keyName)}
	}

	// Haven't found a key, so we've gotta create it and write it into storage.
	var result keyEntry
	result.ID = genKeyId()
	result.Name = keyName
	result.PrivateKey = keyValue
	result.PrivateKeyType = keyType

	// Finally, we can write the key to storage.
	if err := writeKey(ctx, s, result); err != nil {
		return nil, false, err
	}

	// Before we return below, we need to iterate over _all_ issuers and see if
	// one of them has a missing KeyId link, and if so, point it back to
	// ourselves. We fetch the list of issuers up front, even when don't need
	// it, to give ourselves a better chance of succeeding below.
	knownIssuers, err := listIssuers(ctx, s)
	if err != nil {
		return nil, false, err
	}

	// Now, for each issuer, try and compute the issuer<->key link if missing.
	for _, identifier := range knownIssuers {
		existingIssuer, err := fetchIssuerById(ctx, s, identifier)
		if err != nil {
			return nil, false, err
		}

		// If the KeyID value is already present, we can skip it.
		if len(existingIssuer.KeyID) > 0 {
			continue
		}

		// Otherwise, compare public values. Note that there might be multiple
		// certificates (e.g., cross-signed) with the same key.

		cert, err := existingIssuer.GetCertificate()
		if err != nil {
			// Malformed issuer.
			return nil, false, err
		}

		equal, err := certutil.ComparePublicKeysAndType(cert.PublicKey, pkForImportingKey)
		if err != nil {
			return nil, false, err
		}

		if equal {
			// These public keys are equal, so this key entry must be the
			// corresponding private key to this issuer; update it accordingly.
			existingIssuer.KeyID = result.ID
			if err := writeIssuer(ctx, s, existingIssuer); err != nil {
				return nil, false, err
			}
		}
	}

	// If there was no prior default value set and/or we had no known
	// keys when we started, set this key as default.
	keyDefaultSet, err := isDefaultKeySet(ctx, s)
	if err != nil {
		return nil, false, err
	}
	if len(knownKeys) == 0 || !keyDefaultSet {
		if err = updateDefaultKeyId(ctx, s, result.ID); err != nil {
			return nil, false, err
		}
	}

	// All done; return our new key reference.
	return &result, false, nil
}

func (i issuerEntry) GetCertificate() (*x509.Certificate, error) {
	cert, err := parseCertificateFromBytes([]byte(i.Certificate))
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse certificate from issuer: %s: %v", err.Error(), i.ID)}
	}

	return cert, nil
}

func (i issuerEntry) EnsureUsage(usage issuerUsage) error {
	// We want to spit out a nice error message about missing usages.
	if i.Usage.HasUsage(usage) {
		return nil
	}

	issuerRef := fmt.Sprintf("id:%v", i.ID)
	if len(i.Name) > 0 {
		issuerRef = fmt.Sprintf("%v / name:%v", issuerRef, i.Name)
	}

	// These usages differ at some point in time. We've gotta find the first
	// usage that differs and return a logical-sounding error message around
	// that difference.
	for name, candidate := range namedIssuerUsages {
		if usage.HasUsage(candidate) && !i.Usage.HasUsage(candidate) {
			return fmt.Errorf("requested usage %v for issuer [%v] but only had usage %v", name, issuerRef, i.Usage.Names())
		}
	}

	// Maybe we have an unnamed usage that's requested.
	return fmt.Errorf("unknown delta between usages: %v -> %v / for issuer [%v]", usage.Names(), i.Usage.Names(), issuerRef)
}

func listIssuers(ctx context.Context, s logical.Storage) ([]issuerID, error) {
	strList, err := s.List(ctx, issuerPrefix)
	if err != nil {
		return nil, err
	}

	issuerIds := make([]issuerID, 0, len(strList))
	for _, entry := range strList {
		issuerIds = append(issuerIds, issuerID(entry))
	}

	return issuerIds, nil
}

func resolveKeyReference(ctx context.Context, s logical.Storage, reference string) (keyID, error) {
	if reference == defaultRef {
		// Handle fetching the default key.
		config, err := getKeysConfig(ctx, s)
		if err != nil {
			return keyID("config-error"), err
		}
		if len(config.DefaultKeyId) == 0 {
			return KeyRefNotFound, fmt.Errorf("no default key currently configured")
		}

		return config.DefaultKeyId, nil
	}

	// Lookup by a direct get first to see if our reference is an ID, this is quick and cached.
	if len(reference) == uuidLength {
		entry, err := s.Get(ctx, keyPrefix+reference)
		if err != nil {
			return keyID("key-read"), err
		}
		if entry != nil {
			return keyID(reference), nil
		}
	}

	// ... than to pull all keys from storage.
	keys, err := listKeys(ctx, s)
	if err != nil {
		return keyID("list-error"), err
	}
	for _, keyId := range keys {
		key, err := fetchKeyById(ctx, s, keyId)
		if err != nil {
			return keyID("key-read"), err
		}

		if key.Name == reference {
			return key.ID, nil
		}
	}

	// Otherwise, we must not have found the key.
	return KeyRefNotFound, errutil.UserError{Err: fmt.Sprintf("unable to find PKI key for reference: %v", reference)}
}

// fetchIssuerById returns an issuerEntry based on issuerId, if none found an error is returned.
func fetchIssuerById(ctx context.Context, s logical.Storage, issuerId issuerID) (*issuerEntry, error) {
	if len(issuerId) == 0 {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: empty issuer identifier")}
	}

	entry, err := s.Get(ctx, issuerPrefix+issuerId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
	}
	if entry == nil {
		// FIXME: Dedicated/specific error for this?
		return nil, errutil.UserError{Err: fmt.Sprintf("pki issuer id %s does not exist", issuerId.String())}
	}

	var issuer issuerEntry
	if err := entry.DecodeJSON(&issuer); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki issuer with id %s: %v", issuerId.String(), err)}
	}

	return &issuer, nil
}

func writeIssuer(ctx context.Context, s logical.Storage, issuer *issuerEntry) error {
	issuerId := issuer.ID

	json, err := logical.StorageEntryJSON(issuerPrefix+issuerId.String(), issuer)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func deleteIssuer(ctx context.Context, s logical.Storage, id issuerID) (bool, error) {
	config, err := getIssuersConfig(ctx, s)
	if err != nil {
		return false, err
	}

	wasDefault := false
	if config.DefaultIssuerId == id {
		wasDefault = true
		config.DefaultIssuerId = issuerID("")
		if err := setIssuersConfig(ctx, s, config); err != nil {
			return wasDefault, err
		}
	}

	return wasDefault, s.Delete(ctx, issuerPrefix+id.String())
}

func importIssuer(ctx context.Context, b *backend, s logical.Storage, certValue string, issuerName string) (*issuerEntry, bool, error) {
	// importIssuers imports the specified PEM-format certificate (from
	// certValue) into the new PKI storage format. The first return field is a
	// reference to the new issuer; the second is whether or not the issuer
	// already existed during import (in which case, *issuer points to the
	// existing issuer reference and identifier); the last return field is
	// whether or not an error occurred.

	// Before we begin, we need to ensure the PEM formatted certificate looks
	// good. Restricting to "just" `CERTIFICATE` entries is a little
	// restrictive, as it could be a `X509 CERTIFICATE` entry or a custom
	// value wrapping an actual DER cert. So validating the contents of the
	// PEM header is out of the question (and validating the contents of the
	// PEM block is left to our GetCertificate call below).
	//
	// However, we should trim all leading and trailing spaces and add a
	// single new line. This allows callers to blindly concatenate PEM
	// blobs from the API and get roughly what they'd expect.
	//
	// Discussed further in #11960 and RFC 7468.
	certValue = strings.TrimSpace(certValue) + "\n"

	// Extracting the certificate is necessary for two reasons: first, it lets
	// us fetch the serial number; second, for the public key comparison with
	// known keys.
	issuerCert, err := parseCertificateFromBytes([]byte(certValue))
	if err != nil {
		return nil, false, err
	}

	// Ensure this certificate is a usable as a CA certificate.
	if !issuerCert.BasicConstraintsValid || !issuerCert.IsCA {
		return nil, false, errutil.UserError{Err: "Refusing to import non-CA certificate"}
	}

	// Before we can import a known issuer, we first need to know if the issuer
	// exists in storage already. This means iterating through all known
	// issuers and comparing their private value against this value.
	knownIssuers, err := listIssuers(ctx, s)
	if err != nil {
		return nil, false, err
	}

	foundExistingIssuerWithName := false
	for _, identifier := range knownIssuers {
		existingIssuer, err := fetchIssuerById(ctx, s, identifier)
		if err != nil {
			return nil, false, err
		}
		existingIssuerCert, err := existingIssuer.GetCertificate()
		if err != nil {
			return nil, false, err
		}
		if areCertificatesEqual(existingIssuerCert, issuerCert) {
			// Here, we don't need to stitch together the key entries,
			// because the last run should've done that for us (or, when
			// importing a key).
			return existingIssuer, true, nil
		}

		// Allow us to find an existing matching issuer with a different name before erroring out
		if issuerName != "" && existingIssuer.Name == issuerName {
			foundExistingIssuerWithName = true
		}
	}

	if foundExistingIssuerWithName {
		return nil, false, errutil.UserError{Err: fmt.Sprintf("another issuer is using the requested name: %s", issuerName)}
	}

	// Haven't found an issuer, so we've gotta create it and write it into
	// storage.
	var result issuerEntry
	result.ID = genIssuerId()
	result.Name = issuerName
	result.Certificate = certValue
	result.LeafNotAfterBehavior = certutil.ErrNotAfterBehavior
	result.Usage.ToggleUsage(IssuanceUsage, CRLSigningUsage)

	// We shouldn't add CSRs or multiple certificates in this
	countCertificates := strings.Count(result.Certificate, "-BEGIN ")
	if countCertificates != 1 {
		return nil, false, fmt.Errorf("bad issuer: potentially multiple PEM blobs in one certificate storage entry:\n%v", result.Certificate)
	}

	result.SerialNumber = strings.TrimSpace(certutil.GetHexFormatted(issuerCert.SerialNumber.Bytes(), ":"))

	// Before we return below, we need to iterate over _all_ keys and see if
	// one of them a public key matching this certificate, and if so, update our
	// link accordingly. We fetch the list of keys up front, even may not need
	// it, to give ourselves a better chance of succeeding below.
	knownKeys, err := listKeys(ctx, s)
	if err != nil {
		return nil, false, err
	}

	// Now, for each key, try and compute the issuer<->key link. We delay
	// writing issuer to storage as we won't need to update the key, only
	// the issuer.
	for _, identifier := range knownKeys {
		existingKey, err := fetchKeyById(ctx, s, identifier)
		if err != nil {
			return nil, false, err
		}

		equal, err := comparePublicKey(ctx, b, existingKey, issuerCert.PublicKey)
		if err != nil {
			return nil, false, err
		}

		if equal {
			result.KeyID = existingKey.ID
			// Here, there's exactly one stored key with the same public key
			// as us, per guarantees in importKey; as we're importing an
			// issuer, there's no other keys or issuers we'd need to read or
			// update, so exit.
			break
		}
	}

	// Finally, rebuild the chains. In this process, because the provided
	// reference issuer is non-nil, we'll save this issuer to storage.
	if err := rebuildIssuersChains(ctx, s, &result); err != nil {
		return nil, false, err
	}

	// If there was no prior default value set and/or we had no known
	// issuers when we started, set this issuer as default.
	issuerDefaultSet, err := isDefaultIssuerSet(ctx, s)
	if err != nil {
		return nil, false, err
	}
	if len(knownIssuers) == 0 || !issuerDefaultSet {
		if err = updateDefaultIssuerId(ctx, s, result.ID); err != nil {
			return nil, false, err
		}
	}

	// All done; return our new key reference.
	return &result, false, nil
}

func areCertificatesEqual(cert1 *x509.Certificate, cert2 *x509.Certificate) bool {
	return bytes.Compare(cert1.Raw, cert2.Raw) == 0
}

func setLocalCRLConfig(ctx context.Context, s logical.Storage, mapping *localCRLConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageLocalCRLConfig, mapping)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getLocalCRLConfig(ctx context.Context, s logical.Storage) (*localCRLConfigEntry, error) {
	entry, err := s.Get(ctx, storageLocalCRLConfig)
	if err != nil {
		return nil, err
	}

	mapping := &localCRLConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(mapping); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode cluster-local CRL configuration: %v", err)}
		}
	}

	if len(mapping.IssuerIDCRLMap) == 0 {
		mapping.IssuerIDCRLMap = make(map[issuerID]crlID)
	}

	if len(mapping.CRLNumberMap) == 0 {
		mapping.CRLNumberMap = make(map[crlID]int64)
	}

	return mapping, nil
}

func setKeysConfig(ctx context.Context, s logical.Storage, config *keyConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageKeyConfig, config)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getKeysConfig(ctx context.Context, s logical.Storage) (*keyConfigEntry, error) {
	entry, err := s.Get(ctx, storageKeyConfig)
	if err != nil {
		return nil, err
	}

	keyConfig := &keyConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(keyConfig); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode key configuration: %v", err)}
		}
	}

	return keyConfig, nil
}

func setIssuersConfig(ctx context.Context, s logical.Storage, config *issuerConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageIssuerConfig, config)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func getIssuersConfig(ctx context.Context, s logical.Storage) (*issuerConfigEntry, error) {
	entry, err := s.Get(ctx, storageIssuerConfig)
	if err != nil {
		return nil, err
	}

	issuerConfig := &issuerConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(issuerConfig); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode issuer configuration: %v", err)}
		}
	}

	return issuerConfig, nil
}

// Lookup within storage the value of reference, assuming the string is a reference to an issuer entry,
// returning the converted issuerID or an error if not found. This method will not properly resolve the
// special legacyBundleShimID value as we do not want to confuse our special value and a user-provided name of the
// same value.
func resolveIssuerReference(ctx context.Context, s logical.Storage, reference string) (issuerID, error) {
	if reference == defaultRef {
		// Handle fetching the default issuer.
		config, err := getIssuersConfig(ctx, s)
		if err != nil {
			return issuerID("config-error"), err
		}
		if len(config.DefaultIssuerId) == 0 {
			return IssuerRefNotFound, fmt.Errorf("no default issuer currently configured")
		}

		return config.DefaultIssuerId, nil
	}

	// Lookup by a direct get first to see if our reference is an ID, this is quick and cached.
	if len(reference) == uuidLength {
		entry, err := s.Get(ctx, issuerPrefix+reference)
		if err != nil {
			return issuerID("issuer-read"), err
		}
		if entry != nil {
			return issuerID(reference), nil
		}
	}

	// ... than to pull all issuers from storage.
	issuers, err := listIssuers(ctx, s)
	if err != nil {
		return issuerID("list-error"), err
	}

	for _, issuerId := range issuers {
		issuer, err := fetchIssuerById(ctx, s, issuerId)
		if err != nil {
			return issuerID("issuer-read"), err
		}

		if issuer.Name == reference {
			return issuer.ID, nil
		}
	}

	// Otherwise, we must not have found the issuer.
	return IssuerRefNotFound, errutil.UserError{Err: fmt.Sprintf("unable to find PKI issuer for reference: %v", reference)}
}

func resolveIssuerCRLPath(ctx context.Context, b *backend, s logical.Storage, reference string) (string, error) {
	if b.useLegacyBundleCaStorage() {
		return legacyCRLPath, nil
	}

	issuer, err := resolveIssuerReference(ctx, s, reference)
	if err != nil {
		return legacyCRLPath, err
	}

	crlConfig, err := getLocalCRLConfig(ctx, s)
	if err != nil {
		return legacyCRLPath, err
	}

	if crlId, ok := crlConfig.IssuerIDCRLMap[issuer]; ok && len(crlId) > 0 {
		return fmt.Sprintf("crls/%v", crlId), nil
	}

	return legacyCRLPath, fmt.Errorf("unable to find CRL for issuer: id:%v/ref:%v", issuer, reference)
}

// Builds a certutil.CertBundle from the specified issuer identifier,
// optionally loading the key or not. This method supports loading legacy
// bundles using the legacyBundleShimID issuerId, and if no entry is found will return an error.
func fetchCertBundleByIssuerId(ctx context.Context, s logical.Storage, id issuerID, loadKey bool) (*issuerEntry, *certutil.CertBundle, error) {
	if id == legacyBundleShimID {
		// We have not completed the migration, or started a request in legacy mode, so
		// attempt to load the bundle from the legacy location
		issuer, bundle, err := getLegacyCertBundle(ctx, s)
		if err != nil {
			return nil, nil, err
		}
		if issuer == nil || bundle == nil {
			return nil, nil, errutil.UserError{Err: "no legacy cert bundle exists"}
		}

		return issuer, bundle, err
	}

	issuer, err := fetchIssuerById(ctx, s, id)
	if err != nil {
		return nil, nil, err
	}

	var bundle certutil.CertBundle
	bundle.Certificate = issuer.Certificate
	bundle.CAChain = issuer.CAChain
	bundle.SerialNumber = issuer.SerialNumber

	// Fetch the key if it exists. Sometimes we don't need the key immediately.
	if loadKey && issuer.KeyID != keyID("") {
		key, err := fetchKeyById(ctx, s, issuer.KeyID)
		if err != nil {
			return nil, nil, err
		}

		bundle.PrivateKeyType = key.PrivateKeyType
		bundle.PrivateKey = key.PrivateKey
	}

	return issuer, &bundle, nil
}

func writeCaBundle(ctx context.Context, b *backend, s logical.Storage, caBundle *certutil.CertBundle, issuerName string, keyName string) (*issuerEntry, *keyEntry, error) {
	myKey, _, err := importKey(ctx, b, s, caBundle.PrivateKey, keyName, caBundle.PrivateKeyType)
	if err != nil {
		return nil, nil, err
	}

	myIssuer, _, err := importIssuer(ctx, b, s, caBundle.Certificate, issuerName)
	if err != nil {
		return nil, nil, err
	}

	for _, cert := range caBundle.CAChain {
		if _, _, err = importIssuer(ctx, b, s, cert, ""); err != nil {
			return nil, nil, err
		}
	}

	return myIssuer, myKey, nil
}

func genIssuerId() issuerID {
	return issuerID(genUuid())
}

func genKeyId() keyID {
	return keyID(genUuid())
}

func genCRLId() crlID {
	return crlID(genUuid())
}

func genUuid() string {
	aUuid, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return aUuid
}

func isKeyInUse(keyId string, ctx context.Context, s logical.Storage) (inUse bool, issuerId string, err error) {
	knownIssuers, err := listIssuers(ctx, s)
	if err != nil {
		return true, "", err
	}

	for _, issuerId := range knownIssuers {
		issuerEntry, err := fetchIssuerById(ctx, s, issuerId)
		if err != nil {
			return true, issuerId.String(), errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
		}
		if issuerEntry == nil {
			return true, issuerId.String(), errutil.InternalError{Err: fmt.Sprintf("Issuer listed: %s does not exist", issuerId.String())}
		}
		if issuerEntry.KeyID.String() == keyId {
			return true, issuerId.String(), nil
		}
	}

	return false, "", nil
}
