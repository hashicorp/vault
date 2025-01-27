// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

var ErrStorageItemNotFound = errors.New("storage item not found")

const (
	storageKeyConfig        = issuing.StorageKeyConfig
	storageIssuerConfig     = issuing.StorageIssuerConfig
	keyPrefix               = issuing.KeyPrefix
	issuerPrefix            = issuing.IssuerPrefix
	storageLocalCRLConfig   = issuing.StorageLocalCRLConfig
	storageUnifiedCRLConfig = issuing.StorageUnifiedCRLConfig

	legacyMigrationBundleLogKey = "config/legacyMigrationBundleLog"
	legacyCertBundlePath        = issuing.LegacyCertBundlePath
	legacyCertBundleBackupPath  = "config/ca_bundle.bak"

	legacyCRLPath        = issuing.LegacyCRLPath
	deltaCRLPath         = issuing.DeltaCRLPath
	deltaCRLPathSuffix   = issuing.DeltaCRLPathSuffix
	unifiedCRLPath       = issuing.UnifiedCRLPath
	unifiedDeltaCRLPath  = issuing.UnifiedDeltaCRLPath
	unifiedCRLPathPrefix = issuing.UnifiedCRLPathPrefix

	autoTidyConfigPath = "config/auto-tidy"
	clusterConfigPath  = "config/cluster"

	autoTidyLastRunPath = "config/auto-tidy-last-run"

	maxRolesToScanOnIssuerChange = 100
	maxRolesToFindOnIssuerChange = 10
)

func ToURLEntries(sc *storageContext, issuer issuing.IssuerID, c *issuing.AiaConfigEntry) (*certutil.URLEntries, error) {
	return issuing.ToURLEntries(sc.Context, sc.Storage, issuer, c)
}

type storageContext struct {
	Context context.Context
	Storage logical.Storage
	Backend *backend
}

var _ pki_backend.StorageContext = (*storageContext)(nil)

func (b *backend) makeStorageContext(ctx context.Context, s logical.Storage) *storageContext {
	return &storageContext{
		Context: ctx,
		Storage: s,
		Backend: b,
	}
}

func (sc *storageContext) WithFreshTimeout(timeout time.Duration) (*storageContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &storageContext{
		Context: ctx,
		Storage: sc.Storage,
		Backend: sc.Backend,
	}, cancel
}

func (sc *storageContext) GetContext() context.Context {
	return sc.Context
}

func (sc *storageContext) GetStorage() logical.Storage {
	return sc.Storage
}

func (sc *storageContext) Logger() hclog.Logger {
	return sc.Backend.Logger()
}

func (sc *storageContext) System() logical.SystemView {
	return sc.Backend.System()
}

func (sc *storageContext) CrlBuilder() pki_backend.CrlBuilderType {
	return sc.Backend.CrlBuilder()
}

func (sc *storageContext) GetUnifiedTransferStatus() *UnifiedTransferStatus {
	return sc.Backend.GetUnifiedTransferStatus()
}

func (sc *storageContext) GetPkiManagedView() managed_key.PkiManagedKeyView {
	return sc.Backend
}

func (sc *storageContext) GetCertificateCounter() issuing.CertificateCounter {
	return sc.Backend.GetCertificateCounter()
}

func (sc *storageContext) UseLegacyBundleCaStorage() bool {
	return sc.Backend.UseLegacyBundleCaStorage()
}

func (sc *storageContext) GetRevokeStorageLock() *sync.RWMutex {
	return sc.Backend.GetRevokeStorageLock()
}

func (sc *storageContext) GetRole(name string) (*issuing.RoleEntry, error) {
	return sc.Backend.GetRole(sc.Context, sc.Storage, name)
}

func (sc *storageContext) listKeys() ([]issuing.KeyID, error) {
	return issuing.ListKeys(sc.Context, sc.Storage)
}

func (sc *storageContext) fetchKeyById(keyId issuing.KeyID) (*issuing.KeyEntry, error) {
	return issuing.FetchKeyById(sc.Context, sc.Storage, keyId)
}

func (sc *storageContext) writeKey(key issuing.KeyEntry) error {
	return issuing.WriteKey(sc.Context, sc.Storage, key)
}

func (sc *storageContext) deleteKey(id issuing.KeyID) (bool, error) {
	return issuing.DeleteKey(sc.Context, sc.Storage, id)
}

func (sc *storageContext) importKey(keyValue string, keyName string, keyType certutil.PrivateKeyType) (*issuing.KeyEntry, bool, error) {
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
	knownKeys, err := sc.listKeys()
	if err != nil {
		return nil, false, err
	}

	// Get our public key from the current inbound key, to compare against all the other keys.
	var pkForImportingKey crypto.PublicKey
	if keyType == certutil.ManagedPrivateKey {
		pkForImportingKey, err = managed_key.GetPublicKeyFromKeyBytes(sc.Context, sc.Backend, []byte(keyValue))
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
		existingKey, err := sc.fetchKeyById(identifier)
		if err != nil {
			return nil, false, err
		}
		areEqual, err := comparePublicKey(sc, existingKey, pkForImportingKey)
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
	var result issuing.KeyEntry
	result.ID = genKeyId()
	result.Name = keyName
	result.PrivateKey = keyValue
	result.PrivateKeyType = keyType

	// Finally, we can write the key to storage.
	if err := sc.writeKey(result); err != nil {
		return nil, false, err
	}

	// Before we return below, we need to iterate over _all_ issuers and see if
	// one of them has a missing KeyId link, and if so, point it back to
	// ourselves. We fetch the list of issuers up front, even when don't need
	// it, to give ourselves a better chance of succeeding below.
	knownIssuers, err := sc.listIssuers()
	if err != nil {
		return nil, false, err
	}

	issuerDefaultSet, err := sc.isDefaultIssuerSet()
	if err != nil {
		return nil, false, err
	}

	// Now, for each issuer, try and compute the issuer<->key link if missing.
	for _, identifier := range knownIssuers {
		existingIssuer, err := sc.fetchIssuerById(identifier)
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
			if err := sc.writeIssuer(existingIssuer); err != nil {
				return nil, false, err
			}

			// If there was no prior default value set and/or we had no known
			// issuers when we started, set this issuer as default.
			if !issuerDefaultSet {
				err = sc.updateDefaultIssuerId(existingIssuer.ID)
				if err != nil {
					return nil, false, err
				}
				issuerDefaultSet = true
			}
		}
	}

	// If there was no prior default value set and/or we had no known
	// keys when we started, set this key as default.
	keyDefaultSet, err := sc.isDefaultKeySet()
	if err != nil {
		return nil, false, err
	}
	if len(knownKeys) == 0 || !keyDefaultSet {
		if err = sc.updateDefaultKeyId(result.ID); err != nil {
			return nil, false, err
		}
	}

	// All done; return our new key reference.
	return &result, false, nil
}

func GetAIAURLs(sc *storageContext, i *issuing.IssuerEntry) (*certutil.URLEntries, error) {
	return issuing.GetAIAURLs(sc.Context, sc.Storage, i)
}

func (sc *storageContext) listIssuers() ([]issuing.IssuerID, error) {
	return issuing.ListIssuers(sc.Context, sc.Storage)
}

func (sc *storageContext) resolveKeyReference(reference string) (issuing.KeyID, error) {
	return issuing.ResolveKeyReference(sc.Context, sc.Storage, reference)
}

// fetchIssuerById returns an IssuerEntry based on issuerId, if none found an error is returned.
func (sc *storageContext) fetchIssuerById(issuerId issuing.IssuerID) (*issuing.IssuerEntry, error) {
	return issuing.FetchIssuerById(sc.Context, sc.Storage, issuerId)
}

func (sc *storageContext) writeIssuer(issuer *issuing.IssuerEntry) error {
	return issuing.WriteIssuer(sc.Context, sc.Storage, issuer)
}

func (sc *storageContext) deleteIssuer(id issuing.IssuerID) (bool, error) {
	return issuing.DeleteIssuer(sc.Context, sc.Storage, id)
}

func (sc *storageContext) importIssuer(certValue string, issuerName string) (*issuing.IssuerEntry, bool, error) {
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

	// Ensure this certificate has a parsed public key. Otherwise, we've
	// likely been given a bad certificate.
	if issuerCert.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm || issuerCert.PublicKey == nil {
		return nil, false, errutil.UserError{Err: "Refusing to import CA certificate with empty PublicKey. This usually means the SubjectPublicKeyInfo field has an OID not recognized by Go, such as 1.2.840.113549.1.1.10 for rsaPSS."}
	}

	// Before we can import a known issuer, we first need to know if the issuer
	// exists in storage already. This means iterating through all known
	// issuers and comparing their private value against this value.
	knownIssuers, err := sc.listIssuers()
	if err != nil {
		return nil, false, err
	}

	foundExistingIssuerWithName := false
	for _, identifier := range knownIssuers {
		existingIssuer, err := sc.fetchIssuerById(identifier)
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
	var result issuing.IssuerEntry
	result.ID = genIssuerId()
	result.Name = issuerName
	result.Certificate = certValue
	result.LeafNotAfterBehavior = certutil.ErrNotAfterBehavior
	result.Usage.ToggleUsage(issuing.AllIssuerUsages)
	result.Version = issuing.LatestIssuerVersion

	// If we lack relevant bits for CRL, prohibit it from being set
	// on the usage side.
	if (issuerCert.KeyUsage&x509.KeyUsageCRLSign) == 0 && result.Usage.HasUsage(issuing.CRLSigningUsage) {
		result.Usage.ToggleUsage(issuing.CRLSigningUsage)
	}

	// We shouldn't add CSRs or multiple certificates in this
	countCertificates := strings.Count(result.Certificate, "-BEGIN ")
	if countCertificates != 1 {
		return nil, false, fmt.Errorf("bad issuer: potentially multiple PEM blobs in one certificate storage entry:\n%v", result.Certificate)
	}

	result.SerialNumber = serialFromCert(issuerCert)

	// Before we return below, we need to iterate over _all_ keys and see if
	// one of them a public key matching this certificate, and if so, update our
	// link accordingly. We fetch the list of keys up front, even may not need
	// it, to give ourselves a better chance of succeeding below.
	knownKeys, err := sc.listKeys()
	if err != nil {
		return nil, false, err
	}

	// Now, for each key, try and compute the issuer<->key link. We delay
	// writing issuer to storage as we won't need to update the key, only
	// the issuer.
	for _, identifier := range knownKeys {
		existingKey, err := sc.fetchKeyById(identifier)
		if err != nil {
			return nil, false, err
		}

		equal, err := comparePublicKey(sc, existingKey, issuerCert.PublicKey)
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
	if err := sc.rebuildIssuersChains(&result); err != nil {
		return nil, false, err
	}

	// If there was no prior default value set and/or we had no known
	// issuers when we started, set this issuer as default.
	issuerDefaultSet, err := sc.isDefaultIssuerSet()
	if err != nil {
		return nil, false, err
	}
	if (len(knownIssuers) == 0 || !issuerDefaultSet) && len(result.KeyID) != 0 {
		if err = sc.updateDefaultIssuerId(result.ID); err != nil {
			return nil, false, err
		}
	}

	// All done; return our new key reference.
	return &result, false, nil
}

func areCertificatesEqual(cert1 *x509.Certificate, cert2 *x509.Certificate) bool {
	return bytes.Equal(cert1.Raw, cert2.Raw)
}

func (sc *storageContext) setLocalCRLConfig(mapping *issuing.InternalCRLConfigEntry) error {
	return issuing.SetLocalCRLConfig(sc.Context, sc.Storage, mapping)
}

func (sc *storageContext) setUnifiedCRLConfig(mapping *issuing.InternalCRLConfigEntry) error {
	return issuing.SetUnifiedCRLConfig(sc.Context, sc.Storage, mapping)
}

func (sc *storageContext) getLocalCRLConfig() (*issuing.InternalCRLConfigEntry, error) {
	return issuing.GetLocalCRLConfig(sc.Context, sc.Storage)
}

func (sc *storageContext) getUnifiedCRLConfig() (*issuing.InternalCRLConfigEntry, error) {
	return issuing.GetUnifiedCRLConfig(sc.Context, sc.Storage)
}

func (sc *storageContext) setKeysConfig(config *issuing.KeyConfigEntry) error {
	return issuing.SetKeysConfig(sc.Context, sc.Storage, config)
}

func (sc *storageContext) getKeysConfig() (*issuing.KeyConfigEntry, error) {
	return issuing.GetKeysConfig(sc.Context, sc.Storage)
}

func (sc *storageContext) setIssuersConfig(config *issuing.IssuerConfigEntry) error {
	return issuing.SetIssuersConfig(sc.Context, sc.Storage, config)
}

func (sc *storageContext) getIssuersConfig() (*issuing.IssuerConfigEntry, error) {
	return issuing.GetIssuersConfig(sc.Context, sc.Storage)
}

// Lookup within storage the value of reference, assuming the string is a reference to an issuer entry,
// returning the converted IssuerID or an error if not found. This method will not properly resolve the
// special legacyBundleShimID value as we do not want to confuse our special value and a user-provided name of the
// same value.
func (sc *storageContext) resolveIssuerReference(reference string) (issuing.IssuerID, error) {
	return issuing.ResolveIssuerReference(sc.Context, sc.Storage, reference)
}

// Builds a certutil.CertBundle from the specified issuer identifier,
// optionally loading the key or not. This method supports loading legacy
// bundles using the legacyBundleShimID issuerId, and if no entry is found will return an error.
func (sc *storageContext) fetchCertBundleByIssuerId(id issuing.IssuerID, loadKey bool) (*issuing.IssuerEntry, *certutil.CertBundle, error) {
	return issuing.FetchCertBundleByIssuerId(sc.Context, sc.Storage, id, loadKey)
}

func (sc *storageContext) writeCaBundle(caBundle *certutil.CertBundle, issuerName string, keyName string) (*issuing.IssuerEntry, *issuing.KeyEntry, error) {
	myKey, _, err := sc.importKey(caBundle.PrivateKey, keyName, caBundle.PrivateKeyType)
	if err != nil {
		return nil, nil, err
	}

	// We may have existing mounts that only contained a key with no certificate yet as a signed CSR
	// was never setup within the mount.
	if caBundle.Certificate == "" {
		return &issuing.IssuerEntry{}, myKey, nil
	}

	myIssuer, _, err := sc.importIssuer(caBundle.Certificate, issuerName)
	if err != nil {
		return nil, nil, err
	}

	for _, cert := range caBundle.CAChain {
		if _, _, err = sc.importIssuer(cert, ""); err != nil {
			return nil, nil, err
		}
	}

	return myIssuer, myKey, nil
}

func genIssuerId() issuing.IssuerID {
	return issuing.IssuerID(genUuid())
}

func genKeyId() issuing.KeyID {
	return issuing.KeyID(genUuid())
}

func genCRLId() issuing.CrlID {
	return issuing.CrlID(genUuid())
}

func genUuid() string {
	aUuid, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return aUuid
}

func (sc *storageContext) isKeyInUse(keyId string) (inUse bool, issuerId string, err error) {
	knownIssuers, err := sc.listIssuers()
	if err != nil {
		return true, "", err
	}

	for _, issuerId := range knownIssuers {
		issuerEntry, err := sc.fetchIssuerById(issuerId)
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

func (sc *storageContext) checkForRolesReferencing(issuerId string) (timeout bool, inUseBy int32, err error) {
	roleEntries, err := sc.Storage.List(sc.Context, "role/")
	if err != nil {
		return false, 0, err
	}

	inUseBy = 0
	checkedRoles := 0

	for _, roleName := range roleEntries {
		entry, err := sc.Storage.Get(sc.Context, "role/"+roleName)
		if err != nil {
			return false, 0, err
		}
		if entry != nil { // If nil, someone deleted an entry since we haven't taken a lock here so just continue
			var role issuing.RoleEntry
			err = entry.DecodeJSON(&role)
			if err != nil {
				return false, inUseBy, err
			}
			if role.Issuer == issuerId {
				inUseBy = inUseBy + 1
				if inUseBy >= maxRolesToFindOnIssuerChange {
					return true, inUseBy, nil
				}
			}
		}
		checkedRoles = checkedRoles + 1
		if checkedRoles >= maxRolesToScanOnIssuerChange {
			return true, inUseBy, nil
		}
	}

	return false, inUseBy, nil
}

func (sc *storageContext) getRevocationConfig() (*pki_backend.CrlConfig, error) {
	entry, err := sc.Storage.Get(sc.Context, "config/crl")
	if err != nil {
		return nil, err
	}

	var result pki_backend.CrlConfig
	if entry == nil {
		result = pki_backend.DefaultCrlConfig
		return &result, nil
	}

	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.Version == 0 {
		// Automatically update existing configurations.
		result.OcspDisable = pki_backend.DefaultCrlConfig.OcspDisable
		result.OcspExpiry = pki_backend.DefaultCrlConfig.OcspExpiry
		result.AutoRebuild = pki_backend.DefaultCrlConfig.AutoRebuild
		result.AutoRebuildGracePeriod = pki_backend.DefaultCrlConfig.AutoRebuildGracePeriod
		result.Version = 1
	}
	if result.Version == 1 {
		if result.DeltaRebuildInterval == "" {
			result.DeltaRebuildInterval = pki_backend.DefaultCrlConfig.DeltaRebuildInterval
		}
		result.Version = 2
	}

	// Depending on client version, it's possible that the expiry is unset.
	// This sets the default value to prevent issues in downstream code.
	if result.Expiry == "" {
		result.Expiry = pki_backend.DefaultCrlConfig.Expiry
	}

	isLocalMount := sc.System().LocalMount()
	if (!constants.IsEnterprise || isLocalMount) && (result.UnifiedCRLOnExistingPaths || result.UnifiedCRL || result.UseGlobalQueue) {
		// An end user must have had Enterprise, enabled the unified config args and then downgraded to OSS.
		sc.Logger().Warn("Not running Vault Enterprise or using a local mount, " +
			"disabling unified_crl, unified_crl_on_existing_paths and cross_cluster_revocation config flags.")
		result.UnifiedCRLOnExistingPaths = false
		result.UnifiedCRL = false
		result.UseGlobalQueue = false
	}

	return &result, nil
}

func (sc *storageContext) setRevocationConfig(config *pki_backend.CrlConfig) error {
	entry, err := logical.StorageEntryJSON("config/crl", config)
	if err != nil {
		return fmt.Errorf("failed building storage entry JSON: %w", err)
	}

	err = sc.Storage.Put(sc.Context, entry)
	if err != nil {
		return fmt.Errorf("failed writing storage entry: %w", err)
	}

	return nil
}

func (sc *storageContext) getAutoTidyConfig() (*tidyConfig, error) {
	entry, err := sc.Storage.Get(sc.Context, autoTidyConfigPath)
	if err != nil {
		return nil, err
	}

	var result tidyConfig
	if entry == nil {
		result = defaultTidyConfig
		return &result, nil
	}

	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.IssuerSafetyBuffer == 0 {
		result.IssuerSafetyBuffer = defaultTidyConfig.IssuerSafetyBuffer
	}

	if result.MinStartupBackoff == 0 {
		result.MinStartupBackoff = defaultTidyConfig.MinStartupBackoff
	}

	if result.MaxStartupBackoff == 0 {
		result.MaxStartupBackoff = defaultTidyConfig.MaxStartupBackoff
	}

	return &result, nil
}

func (sc *storageContext) writeAutoTidyConfig(config *tidyConfig) error {
	entry, err := logical.StorageEntryJSON(autoTidyConfigPath, config)
	if err != nil {
		return err
	}

	err = sc.Storage.Put(sc.Context, entry)
	if err != nil {
		return err
	}

	certCounter := sc.Backend.GetCertificateCounter()
	certCounter.ReconfigureWithTidyConfig(config)

	return nil
}

func (sc *storageContext) listRevokedCerts() ([]string, error) {
	list, err := sc.Storage.List(sc.Context, revokedPath)
	if err != nil {
		return nil, fmt.Errorf("failed listing revoked certs: %w", err)
	}

	return list, err
}

func (sc *storageContext) getClusterConfig() (*issuing.ClusterConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, clusterConfigPath)
	if err != nil {
		return nil, err
	}

	var result issuing.ClusterConfigEntry
	if entry == nil {
		return &result, nil
	}

	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (sc *storageContext) writeClusterConfig(config *issuing.ClusterConfigEntry) error {
	entry, err := logical.StorageEntryJSON(clusterConfigPath, config)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, entry)
}

// tidyLastRun Track the various pieces of information around tidy on a specific cluster
type tidyLastRun struct {
	LastRunTime time.Time
}

func (sc *storageContext) getAutoTidyLastRun() (time.Time, error) {
	entry, err := sc.Storage.Get(sc.Context, autoTidyLastRunPath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed getting auto tidy last run: %w", err)
	}
	if entry == nil {
		return time.Time{}, nil
	}

	var result tidyLastRun
	if err = entry.DecodeJSON(&result); err != nil {
		return time.Time{}, fmt.Errorf("failed parsing auto tidy last run: %w", err)
	}
	return result.LastRunTime, nil
}

func (sc *storageContext) writeAutoTidyLastRun(lastRunTime time.Time) error {
	lastRun := tidyLastRun{LastRunTime: lastRunTime}
	entry, err := logical.StorageEntryJSON(autoTidyLastRunPath, lastRun)
	if err != nil {
		return fmt.Errorf("failed generating json for auto tidy last run: %w", err)
	}

	if err := sc.Storage.Put(sc.Context, entry); err != nil {
		return fmt.Errorf("failed writing auto tidy last run: %w", err)
	}

	return nil
}

func fetchRevocationInfo(sc pki_backend.StorageContext, serial string) (*revocation.RevocationInfo, error) {
	var revInfo *revocation.RevocationInfo
	revEntry, err := fetchCertBySerial(sc, revocation.RevokedPath, serial)
	if err != nil {
		return nil, err
	}
	if revEntry != nil {
		err = revEntry.DecodeJSON(&revInfo)
		if err != nil {
			return nil, fmt.Errorf("error decoding existing revocation info: %w", err)
		}
	}

	return revInfo, nil
}

// filterDirEntries filters out directory entries from a list of entries normally from a List operation.
func filterDirEntries(entries []string) []string {
	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		if strings.HasSuffix(entry, "/") {
			continue
		}
		ids = append(ids, entry)
	}
	return ids
}
