// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/constants"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	storageKeyConfig        = "config/keys"
	storageIssuerConfig     = "config/issuers"
	keyPrefix               = "config/key/"
	issuerPrefix            = "config/issuer/"
	storageLocalCRLConfig   = "crls/config"
	storageUnifiedCRLConfig = "unified-crls/config"

	legacyMigrationBundleLogKey = "config/legacyMigrationBundleLog"
	legacyCertBundlePath        = "config/ca_bundle"
	legacyCertBundleBackupPath  = "config/ca_bundle.bak"
	legacyCRLPath               = "crl"
	deltaCRLPath                = "delta-crl"
	deltaCRLPathSuffix          = "-delta"
	unifiedCRLPath              = "unified-crl"
	unifiedDeltaCRLPath         = "unified-delta-crl"
	unifiedCRLPathPrefix        = "unified-"

	autoTidyConfigPath = "config/auto-tidy"
	clusterConfigPath  = "config/cluster"

	// Used as a quick sanity check for a reference id lookups...
	uuidLength = 36

	maxRolesToScanOnIssuerChange = 100
	maxRolesToFindOnIssuerChange = 10

	latestIssuerVersion = 1
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
	ID             keyID                   `json:"id"`
	Name           string                  `json:"name"`
	PrivateKeyType certutil.PrivateKeyType `json:"private_key_type"`
	PrivateKey     string                  `json:"private_key"`
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
	ReadOnlyUsage    issuerUsage = iota
	IssuanceUsage    issuerUsage = 1 << iota
	CRLSigningUsage  issuerUsage = 1 << iota
	OCSPSigningUsage issuerUsage = 1 << iota

	// When adding a new usage in the future, we'll need to create a usage
	// mask field on the IssuerEntry and handle migrations to a newer mask,
	// inferring a value for the new bits.
	AllIssuerUsages = ReadOnlyUsage | IssuanceUsage | CRLSigningUsage | OCSPSigningUsage
)

var namedIssuerUsages = map[string]issuerUsage{
	"read-only":            ReadOnlyUsage,
	"issuing-certificates": IssuanceUsage,
	"crl-signing":          CRLSigningUsage,
	"ocsp-signing":         OCSPSigningUsage,
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

	// Return the known set of usages in a sorted order to not have Terraform state files flipping
	// saying values are different when it's the same list in a different order.
	keys := make([]string, 0, len(namedIssuerUsages))
	for k := range namedIssuerUsages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		usage := namedIssuerUsages[name]
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
	ID                   issuerID                  `json:"id"`
	Name                 string                    `json:"name"`
	KeyID                keyID                     `json:"key_id"`
	Certificate          string                    `json:"certificate"`
	CAChain              []string                  `json:"ca_chain"`
	ManualChain          []issuerID                `json:"manual_chain"`
	SerialNumber         string                    `json:"serial_number"`
	LeafNotAfterBehavior certutil.NotAfterBehavior `json:"not_after_behavior"`
	Usage                issuerUsage               `json:"usage"`
	RevocationSigAlg     x509.SignatureAlgorithm   `json:"revocation_signature_algorithm"`
	Revoked              bool                      `json:"revoked"`
	RevocationTime       int64                     `json:"revocation_time"`
	RevocationTimeUTC    time.Time                 `json:"revocation_time_utc"`
	AIAURIs              *aiaConfigEntry           `json:"aia_uris,omitempty"`
	LastModified         time.Time                 `json:"last_modified"`
	Version              uint                      `json:"version"`
}

type internalCRLConfigEntry struct {
	IssuerIDCRLMap        map[issuerID]crlID  `json:"issuer_id_crl_map"`
	CRLNumberMap          map[crlID]int64     `json:"crl_number_map"`
	LastCompleteNumberMap map[crlID]int64     `json:"last_complete_number_map"`
	CRLExpirationMap      map[crlID]time.Time `json:"crl_expiration_map"`
	LastModified          time.Time           `json:"last_modified"`
	DeltaLastModified     time.Time           `json:"delta_last_modified"`
	UseGlobalQueue        bool                `json:"cross_cluster_revocation"`
}

type keyConfigEntry struct {
	DefaultKeyId keyID `json:"default"`
}

type issuerConfigEntry struct {
	// This new fetchedDefault field allows us to detect if the default
	// issuer was modified, in turn dispatching the timestamp updater
	// if necessary.
	fetchedDefault             issuerID `json:"-"`
	DefaultIssuerId            issuerID `json:"default"`
	DefaultFollowsLatestIssuer bool     `json:"default_follows_latest_issuer"`
}

type clusterConfigEntry struct {
	Path    string `json:"path"`
	AIAPath string `json:"aia_path"`
}

type aiaConfigEntry struct {
	IssuingCertificates   []string `json:"issuing_certificates"`
	CRLDistributionPoints []string `json:"crl_distribution_points"`
	OCSPServers           []string `json:"ocsp_servers"`
	EnableTemplating      bool     `json:"enable_templating"`
}

func (c *aiaConfigEntry) toURLEntries(sc *storageContext, issuer issuerID) (*certutil.URLEntries, error) {
	if len(c.IssuingCertificates) == 0 && len(c.CRLDistributionPoints) == 0 && len(c.OCSPServers) == 0 {
		return &certutil.URLEntries{}, nil
	}

	result := certutil.URLEntries{
		IssuingCertificates:   c.IssuingCertificates[:],
		CRLDistributionPoints: c.CRLDistributionPoints[:],
		OCSPServers:           c.OCSPServers[:],
	}

	if c.EnableTemplating {
		cfg, err := sc.getClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("error fetching cluster-local address config: %w", err)
		}

		for name, source := range map[string]*[]string{
			"issuing_certificates":    &result.IssuingCertificates,
			"crl_distribution_points": &result.CRLDistributionPoints,
			"ocsp_servers":            &result.OCSPServers,
		} {
			templated := make([]string, len(*source))
			for index, uri := range *source {
				if strings.Contains(uri, "{{cluster_path}}") && len(cfg.Path) == 0 {
					return nil, fmt.Errorf("unable to template AIA URLs as we lack local cluster address information (path)")
				}
				if strings.Contains(uri, "{{cluster_aia_path}}") && len(cfg.AIAPath) == 0 {
					return nil, fmt.Errorf("unable to template AIA URLs as we lack local cluster address information (aia_path)")
				}
				if strings.Contains(uri, "{{issuer_id}}") && len(issuer) == 0 {
					// Elide issuer AIA info as we lack an issuer_id.
					return nil, fmt.Errorf("unable to template AIA URLs as we lack an issuer_id for this operation")
				}

				uri = strings.ReplaceAll(uri, "{{cluster_path}}", cfg.Path)
				uri = strings.ReplaceAll(uri, "{{cluster_aia_path}}", cfg.AIAPath)
				uri = strings.ReplaceAll(uri, "{{issuer_id}}", issuer.String())
				templated[index] = uri
			}

			if uri := validateURLs(templated); uri != "" {
				return nil, fmt.Errorf("error validating templated %v; invalid URI: %v", name, uri)
			}

			*source = templated
		}
	}

	return &result, nil
}

type storageContext struct {
	Context context.Context
	Storage logical.Storage
	Backend *backend
}

func (b *backend) makeStorageContext(ctx context.Context, s logical.Storage) *storageContext {
	return &storageContext{
		Context: ctx,
		Storage: s,
		Backend: b,
	}
}

func (sc *storageContext) listKeys() ([]keyID, error) {
	strList, err := sc.Storage.List(sc.Context, keyPrefix)
	if err != nil {
		return nil, err
	}

	keyIds := make([]keyID, 0, len(strList))
	for _, entry := range strList {
		keyIds = append(keyIds, keyID(entry))
	}

	return keyIds, nil
}

func (sc *storageContext) fetchKeyById(keyId keyID) (*keyEntry, error) {
	if len(keyId) == 0 {
		return nil, errutil.InternalError{Err: "unable to fetch pki key: empty key identifier"}
	}

	entry, err := sc.Storage.Get(sc.Context, keyPrefix+keyId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: %v", err)}
	}
	if entry == nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("pki key id %s does not exist", keyId.String())}
	}

	var key keyEntry
	if err := entry.DecodeJSON(&key); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki key with id %s: %v", keyId.String(), err)}
	}

	return &key, nil
}

func (sc *storageContext) writeKey(key keyEntry) error {
	keyId := key.ID

	json, err := logical.StorageEntryJSON(keyPrefix+keyId.String(), key)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, json)
}

func (sc *storageContext) deleteKey(id keyID) (bool, error) {
	config, err := sc.getKeysConfig()
	if err != nil {
		return false, err
	}

	wasDefault := false
	if config.DefaultKeyId == id {
		wasDefault = true
		config.DefaultKeyId = keyID("")
		if err := sc.setKeysConfig(config); err != nil {
			return wasDefault, err
		}
	}

	return wasDefault, sc.Storage.Delete(sc.Context, keyPrefix+id.String())
}

func (sc *storageContext) importKey(keyValue string, keyName string, keyType certutil.PrivateKeyType) (*keyEntry, bool, error) {
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
		managedKeyUUID, err := extractManagedKeyId([]byte(keyValue))
		if err != nil {
			return nil, false, errutil.InternalError{Err: fmt.Sprintf("failed extracting managed key uuid from key: %v", err)}
		}
		pkForImportingKey, err = getManagedKeyPublicKey(sc.Context, sc.Backend, managedKeyUUID)
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
	var result keyEntry
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

func (i issuerEntry) CanMaybeSignWithAlgo(algo x509.SignatureAlgorithm) error {
	// Hack: Go isn't kind enough expose its lovely signatureAlgorithmDetails
	// informational struct for our usage. However, we don't want to actually
	// fetch the private key and attempt a signature with this algo (as we'll
	// mint new, previously unsigned material in the process that could maybe
	// be potentially abused if it leaks).
	//
	// So...
	//
	// ...we maintain our own mapping of cert.PKI<->sigAlgos. Notably, we
	// exclude DSA support as the PKI engine has never supported DSA keys.
	if algo == x509.UnknownSignatureAlgorithm {
		// Special cased to indicate upgrade and letting Go automatically
		// chose the correct value.
		return nil
	}

	cert, err := i.GetCertificate()
	if err != nil {
		return fmt.Errorf("unable to parse issuer's potential signature algorithm types: %w", err)
	}

	switch cert.PublicKeyAlgorithm {
	case x509.RSA:
		switch algo {
		case x509.SHA256WithRSA, x509.SHA384WithRSA, x509.SHA512WithRSA,
			x509.SHA256WithRSAPSS, x509.SHA384WithRSAPSS,
			x509.SHA512WithRSAPSS:
			return nil
		}
	case x509.ECDSA:
		switch algo {
		case x509.ECDSAWithSHA256, x509.ECDSAWithSHA384, x509.ECDSAWithSHA512:
			return nil
		}
	case x509.Ed25519:
		switch algo {
		case x509.PureEd25519:
			return nil
		}
	}

	return fmt.Errorf("unable to use issuer of type %v to sign with %v key type", cert.PublicKeyAlgorithm.String(), algo.String())
}

func (i issuerEntry) GetAIAURLs(sc *storageContext) (*certutil.URLEntries, error) {
	// Default to the per-issuer AIA URLs.
	entries := i.AIAURIs

	// If none are set (either due to a nil entry or because no URLs have
	// been provided), fall back to the global AIA URL config.
	if entries == nil || (len(entries.IssuingCertificates) == 0 && len(entries.CRLDistributionPoints) == 0 && len(entries.OCSPServers) == 0) {
		var err error

		entries, err = getGlobalAIAURLs(sc.Context, sc.Storage)
		if err != nil {
			return nil, err
		}
	}

	if entries == nil {
		return &certutil.URLEntries{}, nil
	}

	return entries.toURLEntries(sc, i.ID)
}

func (sc *storageContext) listIssuers() ([]issuerID, error) {
	strList, err := sc.Storage.List(sc.Context, issuerPrefix)
	if err != nil {
		return nil, err
	}

	issuerIds := make([]issuerID, 0, len(strList))
	for _, entry := range strList {
		issuerIds = append(issuerIds, issuerID(entry))
	}

	return issuerIds, nil
}

func (sc *storageContext) resolveKeyReference(reference string) (keyID, error) {
	if reference == defaultRef {
		// Handle fetching the default key.
		config, err := sc.getKeysConfig()
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
		entry, err := sc.Storage.Get(sc.Context, keyPrefix+reference)
		if err != nil {
			return keyID("key-read"), err
		}
		if entry != nil {
			return keyID(reference), nil
		}
	}

	// ... than to pull all keys from storage.
	keys, err := sc.listKeys()
	if err != nil {
		return keyID("list-error"), err
	}
	for _, keyId := range keys {
		key, err := sc.fetchKeyById(keyId)
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
func (sc *storageContext) fetchIssuerById(issuerId issuerID) (*issuerEntry, error) {
	if len(issuerId) == 0 {
		return nil, errutil.InternalError{Err: "unable to fetch pki issuer: empty issuer identifier"}
	}

	entry, err := sc.Storage.Get(sc.Context, issuerPrefix+issuerId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
	}
	if entry == nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("pki issuer id %s does not exist", issuerId.String())}
	}

	var issuer issuerEntry
	if err := entry.DecodeJSON(&issuer); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki issuer with id %s: %v", issuerId.String(), err)}
	}

	return sc.upgradeIssuerIfRequired(&issuer), nil
}

func (sc *storageContext) upgradeIssuerIfRequired(issuer *issuerEntry) *issuerEntry {
	// *NOTE*: Don't attempt to write out the issuer here as it may cause ErrReadOnly that will direct the
	// request all the way up to the primary cluster which would be horrible for local cluster operations such
	// as generating a leaf cert or a revoke.
	// Also even though we could tell if we are the primary cluster's active node, we can't tell if we have the
	// a full rw issuer lock, so it might not be safe to write.
	if issuer.Version == latestIssuerVersion {
		return issuer
	}

	if issuer.Version == 0 {
		// Upgrade at this step requires interrogating the certificate itself;
		// if this decode fails, it indicates internal problems and the
		// request will subsequently fail elsewhere. However, decoding this
		// certificate is mildly expensive, so we only do it in the event of
		// a Version 0 certificate.
		cert, err := issuer.GetCertificate()
		if err != nil {
			return issuer
		}

		hadCRL := issuer.Usage.HasUsage(CRLSigningUsage)
		// Remove CRL signing usage if it exists on the issuer but doesn't
		// exist in the KU of the x509 certificate.
		if hadCRL && (cert.KeyUsage&x509.KeyUsageCRLSign) == 0 {
			issuer.Usage.ToggleUsage(OCSPSigningUsage)
		}

		// Handle our new OCSPSigning usage flag for earlier versions. If we
		// had it (prior to removing it in this upgrade), we'll add the OCSP
		// flag since EKUs don't matter.
		if hadCRL && !issuer.Usage.HasUsage(OCSPSigningUsage) {
			issuer.Usage.ToggleUsage(OCSPSigningUsage)
		}
	}

	issuer.Version = latestIssuerVersion
	return issuer
}

func (sc *storageContext) writeIssuer(issuer *issuerEntry) error {
	issuerId := issuer.ID
	if issuer.LastModified.IsZero() {
		issuer.LastModified = time.Now().UTC()
	}

	json, err := logical.StorageEntryJSON(issuerPrefix+issuerId.String(), issuer)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, json)
}

func (sc *storageContext) deleteIssuer(id issuerID) (bool, error) {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return false, err
	}

	wasDefault := false
	if config.DefaultIssuerId == id {
		wasDefault = true
		// Overwrite the fetched default issuer as we're going to remove this
		// entry.
		config.fetchedDefault = issuerID("")
		config.DefaultIssuerId = issuerID("")
		if err := sc.setIssuersConfig(config); err != nil {
			return wasDefault, err
		}
	}

	return wasDefault, sc.Storage.Delete(sc.Context, issuerPrefix+id.String())
}

func (sc *storageContext) importIssuer(certValue string, issuerName string) (*issuerEntry, bool, error) {
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
	var result issuerEntry
	result.ID = genIssuerId()
	result.Name = issuerName
	result.Certificate = certValue
	result.LeafNotAfterBehavior = certutil.ErrNotAfterBehavior
	result.Usage.ToggleUsage(AllIssuerUsages)
	result.Version = latestIssuerVersion

	// If we lack relevant bits for CRL, prohibit it from being set
	// on the usage side.
	if (issuerCert.KeyUsage&x509.KeyUsageCRLSign) == 0 && result.Usage.HasUsage(CRLSigningUsage) {
		result.Usage.ToggleUsage(CRLSigningUsage)
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

func (sc *storageContext) _setInternalCRLConfig(mapping *internalCRLConfigEntry, path string) error {
	json, err := logical.StorageEntryJSON(path, mapping)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, json)
}

func (sc *storageContext) setLocalCRLConfig(mapping *internalCRLConfigEntry) error {
	return sc._setInternalCRLConfig(mapping, storageLocalCRLConfig)
}

func (sc *storageContext) setUnifiedCRLConfig(mapping *internalCRLConfigEntry) error {
	return sc._setInternalCRLConfig(mapping, storageUnifiedCRLConfig)
}

func (sc *storageContext) _getInternalCRLConfig(path string) (*internalCRLConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, path)
	if err != nil {
		return nil, err
	}

	mapping := &internalCRLConfigEntry{}
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

	if len(mapping.LastCompleteNumberMap) == 0 {
		mapping.LastCompleteNumberMap = make(map[crlID]int64)

		// Since this might not exist on migration, we want to guess as
		// to the last full CRL number was. This was likely the last
		// value from CRLNumberMap if it existed, since we're just adding
		// the mapping here in this block.
		//
		// After the next full CRL build, we will have set this value
		// correctly, so it doesn't really matter in the long term if
		// we're off here.
		for id, number := range mapping.CRLNumberMap {
			// Decrement by one, since CRLNumberMap is the future number,
			// not the last built number.
			mapping.LastCompleteNumberMap[id] = number - 1
		}
	}

	if len(mapping.CRLExpirationMap) == 0 {
		mapping.CRLExpirationMap = make(map[crlID]time.Time)
	}

	return mapping, nil
}

func (sc *storageContext) getLocalCRLConfig() (*internalCRLConfigEntry, error) {
	return sc._getInternalCRLConfig(storageLocalCRLConfig)
}

func (sc *storageContext) getUnifiedCRLConfig() (*internalCRLConfigEntry, error) {
	return sc._getInternalCRLConfig(storageUnifiedCRLConfig)
}

func (sc *storageContext) setKeysConfig(config *keyConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageKeyConfig, config)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, json)
}

func (sc *storageContext) getKeysConfig() (*keyConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, storageKeyConfig)
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

func (sc *storageContext) setIssuersConfig(config *issuerConfigEntry) error {
	json, err := logical.StorageEntryJSON(storageIssuerConfig, config)
	if err != nil {
		return err
	}

	if err := sc.Storage.Put(sc.Context, json); err != nil {
		return err
	}

	if err := sc.changeDefaultIssuerTimestamps(config.fetchedDefault, config.DefaultIssuerId); err != nil {
		return err
	}

	return nil
}

func (sc *storageContext) getIssuersConfig() (*issuerConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, storageIssuerConfig)
	if err != nil {
		return nil, err
	}

	issuerConfig := &issuerConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(issuerConfig); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode issuer configuration: %v", err)}
		}
	}
	issuerConfig.fetchedDefault = issuerConfig.DefaultIssuerId

	return issuerConfig, nil
}

// Lookup within storage the value of reference, assuming the string is a reference to an issuer entry,
// returning the converted issuerID or an error if not found. This method will not properly resolve the
// special legacyBundleShimID value as we do not want to confuse our special value and a user-provided name of the
// same value.
func (sc *storageContext) resolveIssuerReference(reference string) (issuerID, error) {
	if reference == defaultRef {
		// Handle fetching the default issuer.
		config, err := sc.getIssuersConfig()
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
		entry, err := sc.Storage.Get(sc.Context, issuerPrefix+reference)
		if err != nil {
			return issuerID("issuer-read"), err
		}
		if entry != nil {
			return issuerID(reference), nil
		}
	}

	// ... than to pull all issuers from storage.
	issuers, err := sc.listIssuers()
	if err != nil {
		return issuerID("list-error"), err
	}

	for _, issuerId := range issuers {
		issuer, err := sc.fetchIssuerById(issuerId)
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

func (sc *storageContext) resolveIssuerCRLPath(reference string, unified bool) (string, error) {
	if sc.Backend.useLegacyBundleCaStorage() {
		return legacyCRLPath, nil
	}

	issuer, err := sc.resolveIssuerReference(reference)
	if err != nil {
		return legacyCRLPath, err
	}

	configPath := storageLocalCRLConfig
	if unified {
		configPath = storageUnifiedCRLConfig
	}

	crlConfig, err := sc._getInternalCRLConfig(configPath)
	if err != nil {
		return legacyCRLPath, err
	}

	if crlId, ok := crlConfig.IssuerIDCRLMap[issuer]; ok && len(crlId) > 0 {
		path := fmt.Sprintf("crls/%v", crlId)
		if unified {
			path = unifiedCRLPathPrefix + path
		}

		return path, nil
	}

	return legacyCRLPath, fmt.Errorf("unable to find CRL for issuer: id:%v/ref:%v", issuer, reference)
}

// Builds a certutil.CertBundle from the specified issuer identifier,
// optionally loading the key or not. This method supports loading legacy
// bundles using the legacyBundleShimID issuerId, and if no entry is found will return an error.
func (sc *storageContext) fetchCertBundleByIssuerId(id issuerID, loadKey bool) (*issuerEntry, *certutil.CertBundle, error) {
	if id == legacyBundleShimID {
		// We have not completed the migration, or started a request in legacy mode, so
		// attempt to load the bundle from the legacy location
		issuer, bundle, err := getLegacyCertBundle(sc.Context, sc.Storage)
		if err != nil {
			return nil, nil, err
		}
		if issuer == nil || bundle == nil {
			return nil, nil, errutil.UserError{Err: "no legacy cert bundle exists"}
		}

		return issuer, bundle, err
	}

	issuer, err := sc.fetchIssuerById(id)
	if err != nil {
		return nil, nil, err
	}

	var bundle certutil.CertBundle
	bundle.Certificate = issuer.Certificate
	bundle.CAChain = issuer.CAChain
	bundle.SerialNumber = issuer.SerialNumber

	// Fetch the key if it exists. Sometimes we don't need the key immediately.
	if loadKey && issuer.KeyID != keyID("") {
		key, err := sc.fetchKeyById(issuer.KeyID)
		if err != nil {
			return nil, nil, err
		}

		bundle.PrivateKeyType = key.PrivateKeyType
		bundle.PrivateKey = key.PrivateKey
	}

	return issuer, &bundle, nil
}

func (sc *storageContext) writeCaBundle(caBundle *certutil.CertBundle, issuerName string, keyName string) (*issuerEntry, *keyEntry, error) {
	myKey, _, err := sc.importKey(caBundle.PrivateKey, keyName, caBundle.PrivateKeyType)
	if err != nil {
		return nil, nil, err
	}

	// We may have existing mounts that only contained a key with no certificate yet as a signed CSR
	// was never setup within the mount.
	if caBundle.Certificate == "" {
		return &issuerEntry{}, myKey, nil
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
			var role roleEntry
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

func (sc *storageContext) getRevocationConfig() (*crlConfig, error) {
	entry, err := sc.Storage.Get(sc.Context, "config/crl")
	if err != nil {
		return nil, err
	}

	var result crlConfig
	if entry == nil {
		result = defaultCrlConfig
		return &result, nil
	}

	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	if result.Version == 0 {
		// Automatically update existing configurations.
		result.OcspDisable = defaultCrlConfig.OcspDisable
		result.OcspExpiry = defaultCrlConfig.OcspExpiry
		result.AutoRebuild = defaultCrlConfig.AutoRebuild
		result.AutoRebuildGracePeriod = defaultCrlConfig.AutoRebuildGracePeriod
		result.Version = 1
	}
	if result.Version == 1 {
		if result.DeltaRebuildInterval == "" {
			result.DeltaRebuildInterval = defaultCrlConfig.DeltaRebuildInterval
		}
		result.Version = 2
	}

	// Depending on client version, it's possible that the expiry is unset.
	// This sets the default value to prevent issues in downstream code.
	if result.Expiry == "" {
		result.Expiry = defaultCrlConfig.Expiry
	}

	if !constants.IsEnterprise && (result.UnifiedCRLOnExistingPaths || result.UnifiedCRL || result.UseGlobalQueue) {
		// An end user must have had Enterprise, enabled the unified config args and then downgraded to OSS.
		sc.Backend.Logger().Warn("Not running Vault Enterprise, " +
			"disabling unified_crl, unified_crl_on_existing_paths and cross_cluster_revocation config flags.")
		result.UnifiedCRLOnExistingPaths = false
		result.UnifiedCRL = false
		result.UseGlobalQueue = false
	}

	return &result, nil
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

	sc.Backend.publishCertCountMetrics.Store(config.PublishMetrics)

	// To Potentially Disable Certificate Counting
	if config.MaintainCount == false {
		certCountWasEnabled := sc.Backend.certCountEnabled.Swap(config.MaintainCount)
		if certCountWasEnabled {
			sc.Backend.certsCounted.Store(true)
			sc.Backend.certCountError = "Cert Count is Disabled: enable via Tidy Config maintain_stored_certificate_counts"
			sc.Backend.possibleDoubleCountedSerials = nil        // This won't stop a list operation, but will stop an expensive clean-up during initialize
			sc.Backend.possibleDoubleCountedRevokedSerials = nil // This won't stop a list operation, but will stop an expensive clean-up during initialize
			sc.Backend.certCount.Store(0)
			sc.Backend.revokedCertCount.Store(0)
		}
	} else { // To Potentially Enable Certificate Counting
		if sc.Backend.certCountEnabled.Load() == false {
			// We haven't written "re-enable certificate counts" outside the initialize function
			// Any call derived call to do so is likely to time out on ~2 million certs
			sc.Backend.certCountError = "Certificate Counting Has Not Been Initialized, re-initialize this mount"
		}
	}

	return nil
}

func (sc *storageContext) listRevokedCerts() ([]string, error) {
	list, err := sc.Storage.List(sc.Context, revokedPath)
	if err != nil {
		return nil, fmt.Errorf("failed listing revoked certs: %w", err)
	}

	return list, err
}

func (sc *storageContext) getClusterConfig() (*clusterConfigEntry, error) {
	entry, err := sc.Storage.Get(sc.Context, clusterConfigPath)
	if err != nil {
		return nil, err
	}

	var result clusterConfigEntry
	if entry == nil {
		return &result, nil
	}

	if err = entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (sc *storageContext) writeClusterConfig(config *clusterConfigEntry) error {
	entry, err := logical.StorageEntryJSON(clusterConfigPath, config)
	if err != nil {
		return err
	}

	return sc.Storage.Put(sc.Context, entry)
}

func (sc *storageContext) fetchRevocationInfo(serial string) (*revocationInfo, error) {
	var revInfo *revocationInfo
	revEntry, err := fetchCertBySerial(sc, revokedPath, serial)
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
