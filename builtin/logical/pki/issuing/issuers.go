// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	ReadOnlyUsage    IssuerUsage = iota
	IssuanceUsage    IssuerUsage = 1 << iota
	CRLSigningUsage  IssuerUsage = 1 << iota
	OCSPSigningUsage IssuerUsage = 1 << iota
)

const (
	// When adding a new usage in the future, we'll need to create a usage
	// mask field on the IssuerEntry and handle migrations to a newer mask,
	// inferring a value for the new bits.
	AllIssuerUsages = ReadOnlyUsage | IssuanceUsage | CRLSigningUsage | OCSPSigningUsage

	DefaultRef   = "default"
	IssuerPrefix = "config/issuer/"

	// Used as a quick sanity check for a reference id lookups...
	uuidLength = 36

	IssuerRefNotFound   = IssuerID("not-found")
	LatestIssuerVersion = 1

	LegacyCertBundlePath  = "config/ca_bundle"
	LegacyBundleShimID    = IssuerID("legacy-entry-shim-id")
	LegacyBundleShimKeyID = KeyID("legacy-entry-shim-key-id")

	LegacyCRLPath        = "crl"
	DeltaCRLPath         = "delta-crl"
	DeltaCRLPathSuffix   = "-delta"
	UnifiedCRLPath       = "unified-crl"
	UnifiedDeltaCRLPath  = "unified-delta-crl"
	UnifiedCRLPathPrefix = "unified-"
)

type IssuerID string

func (p IssuerID) String() string {
	return string(p)
}

type IssuerUsage uint

var namedIssuerUsages = map[string]IssuerUsage{
	"read-only":            ReadOnlyUsage,
	"issuing-certificates": IssuanceUsage,
	"crl-signing":          CRLSigningUsage,
	"ocsp-signing":         OCSPSigningUsage,
}

func (i *IssuerUsage) ToggleUsage(usages ...IssuerUsage) {
	for _, usage := range usages {
		*i ^= usage
	}
}

func (i IssuerUsage) HasUsage(usage IssuerUsage) bool {
	return (i & usage) == usage
}

func (i IssuerUsage) Names() string {
	var names []string
	var builtUsage IssuerUsage

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

func NewIssuerUsageFromNames(names []string) (IssuerUsage, error) {
	var result IssuerUsage
	for index, name := range names {
		usage, ok := namedIssuerUsages[name]
		if !ok {
			return ReadOnlyUsage, fmt.Errorf("unknown name for usage at index %v: %v", index, name)
		}

		result.ToggleUsage(usage)
	}

	return result, nil
}

type IssuerEntry struct {
	ID                   IssuerID                  `json:"id"`
	Name                 string                    `json:"name"`
	KeyID                KeyID                     `json:"key_id"`
	Certificate          string                    `json:"certificate"`
	CAChain              []string                  `json:"ca_chain"`
	ManualChain          []IssuerID                `json:"manual_chain"`
	SerialNumber         string                    `json:"serial_number"`
	LeafNotAfterBehavior certutil.NotAfterBehavior `json:"not_after_behavior"`
	Usage                IssuerUsage               `json:"usage"`
	RevocationSigAlg     x509.SignatureAlgorithm   `json:"revocation_signature_algorithm"`
	Revoked              bool                      `json:"revoked"`
	RevocationTime       int64                     `json:"revocation_time"`
	RevocationTimeUTC    time.Time                 `json:"revocation_time_utc"`
	AIAURIs              *AiaConfigEntry           `json:"aia_uris,omitempty"`
	LastModified         time.Time                 `json:"last_modified"`
	Version              uint                      `json:"version"`
	entIssuerEntry
}

// GetCertificate returns a x509.Certificate of the CA certificate
// represented by this issuer.
func (i IssuerEntry) GetCertificate() (*x509.Certificate, error) {
	cert, err := parsing.ParseCertificateFromString(i.Certificate)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse certificate from issuer: %s: %v", err.Error(), i.ID)}
	}

	return cert, nil
}

// GetFullCaChain returns a slice of x509.Certificate values of this issuer full ca chain,
// which starts with the CA certificate represented by this issuer followed by the entire CA chain
func (i IssuerEntry) GetFullCaChain() ([]*x509.Certificate, error) {
	var chains []*x509.Certificate
	issuerCert, err := i.GetCertificate()
	if err != nil {
		return nil, err
	}

	chains = append(chains, issuerCert)

	for rangeI, chainVal := range i.CAChain {
		parsedChainVal, err := parsing.ParseCertificateFromString(chainVal)
		if err != nil {
			return nil, fmt.Errorf("error parsing issuer %s ca chain index value [%d]: %w", i.ID, rangeI, err)
		}

		if bytes.Equal(parsedChainVal.Raw, issuerCert.Raw) {
			continue
		}
		chains = append(chains, parsedChainVal)
	}

	return chains, nil
}

func (i IssuerEntry) EnsureUsage(usage IssuerUsage) error {
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

func (i IssuerEntry) CanMaybeSignWithAlgo(algo x509.SignatureAlgorithm) error {
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

// ResolveAndFetchIssuerForIssuance takes a name or uuid referencing an issuer, loads the issuer
// and validates that we have the associated private key and is allowed to perform issuance operations.
func ResolveAndFetchIssuerForIssuance(ctx context.Context, s logical.Storage, issuerName string) (*IssuerEntry, error) {
	if len(issuerName) == 0 {
		return nil, fmt.Errorf("unable to fetch pki issuer: empty issuer name")
	}
	issuerId, err := ResolveIssuerReference(ctx, s, issuerName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve issuer %s: %w", issuerName, err)
	}

	issuer, err := FetchIssuerById(ctx, s, issuerId)
	if err != nil {
		return nil, fmt.Errorf("failed to load issuer %s: %w", issuerName, err)
	}

	if issuer.Usage.HasUsage(IssuanceUsage) && len(issuer.KeyID) > 0 {
		return issuer, nil
	}

	return nil, fmt.Errorf("issuer %s missing proper issuance usage or doesn't have associated key", issuerName)
}

func ResolveIssuerReference(ctx context.Context, s logical.Storage, reference string) (IssuerID, error) {
	if reference == DefaultRef {
		// Handle fetching the default issuer.
		config, err := GetIssuersConfig(ctx, s)
		if err != nil {
			return IssuerID("config-error"), err
		}
		if len(config.DefaultIssuerId) == 0 {
			return IssuerRefNotFound, fmt.Errorf("no default issuer currently configured")
		}

		return config.DefaultIssuerId, nil
	}

	// Lookup by a direct get first to see if our reference is an ID, this is quick and cached.
	if len(reference) == uuidLength {
		entry, err := s.Get(ctx, IssuerPrefix+reference)
		if err != nil {
			return IssuerID("issuer-read"), err
		}
		if entry != nil {
			return IssuerID(reference), nil
		}
	}

	// ... than to pull all issuers from storage.
	issuers, err := ListIssuers(ctx, s)
	if err != nil {
		return IssuerID("list-error"), err
	}

	for _, issuerId := range issuers {
		issuer, err := FetchIssuerById(ctx, s, issuerId)
		if err != nil {
			return IssuerID("issuer-read"), err
		}

		if issuer.Name == reference {
			return issuer.ID, nil
		}
	}

	// Otherwise, we must not have found the issuer.
	return IssuerRefNotFound, errutil.UserError{Err: fmt.Sprintf("unable to find PKI issuer for reference: %v", reference)}
}

func ListIssuers(ctx context.Context, s logical.Storage) ([]IssuerID, error) {
	strList, err := s.List(ctx, IssuerPrefix)
	if err != nil {
		return nil, err
	}

	issuerIds := make([]IssuerID, 0, len(strList))
	for _, entry := range strList {
		issuerIds = append(issuerIds, IssuerID(entry))
	}

	return issuerIds, nil
}

// FetchIssuerById returns an IssuerEntry based on issuerId, if none found an error is returned.
func FetchIssuerById(ctx context.Context, s logical.Storage, issuerId IssuerID) (*IssuerEntry, error) {
	if len(issuerId) == 0 {
		return nil, errutil.InternalError{Err: "unable to fetch pki issuer: empty issuer identifier"}
	}

	entry, err := s.Get(ctx, IssuerPrefix+issuerId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki issuer: %v", err)}
	}
	if entry == nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("pki issuer id %s does not exist", issuerId.String())}
	}

	var issuer IssuerEntry
	if err := entry.DecodeJSON(&issuer); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki issuer with id %s: %v", issuerId.String(), err)}
	}

	return upgradeIssuerIfRequired(&issuer), nil
}

func WriteIssuer(ctx context.Context, s logical.Storage, issuer *IssuerEntry) error {
	issuerId := issuer.ID
	if issuer.LastModified.IsZero() {
		issuer.LastModified = time.Now().UTC()
	}

	json, err := logical.StorageEntryJSON(IssuerPrefix+issuerId.String(), issuer)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func DeleteIssuer(ctx context.Context, s logical.Storage, id IssuerID) (bool, error) {
	config, err := GetIssuersConfig(ctx, s)
	if err != nil {
		return false, err
	}

	wasDefault := false
	if config.DefaultIssuerId == id {
		wasDefault = true
		// Overwrite the fetched default issuer as we're going to remove this
		// entry.
		config.fetchedDefault = IssuerID("")
		config.DefaultIssuerId = IssuerID("")
		if err := SetIssuersConfig(ctx, s, config); err != nil {
			return wasDefault, err
		}
	}

	return wasDefault, s.Delete(ctx, IssuerPrefix+id.String())
}

func upgradeIssuerIfRequired(issuer *IssuerEntry) *IssuerEntry {
	// *NOTE*: Don't attempt to write out the issuer here as it may cause ErrReadOnly that will direct the
	// request all the way up to the primary cluster which would be horrible for local cluster operations such
	// as generating a leaf cert or a revoke.
	// Also even though we could tell if we are the primary cluster's active node, we can't tell if we have the
	// a full rw issuer lock, so it might not be safe to write.
	if issuer.Version == LatestIssuerVersion {
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
			issuer.Usage.ToggleUsage(CRLSigningUsage)
		}

		// Handle our new OCSPSigning usage flag for earlier versions. If we
		// had it (prior to removing it in this upgrade), we'll add the OCSP
		// flag since EKUs don't matter.
		if hadCRL && !issuer.Usage.HasUsage(OCSPSigningUsage) {
			issuer.Usage.ToggleUsage(OCSPSigningUsage)
		}
	}

	issuer.Version = LatestIssuerVersion
	return issuer
}

// FetchCAInfoByIssuerId will fetch the CA info, will return an error if no ca info exists for the given issuerId.
// This does support the loading using the legacyBundleShimID
func FetchCAInfoByIssuerId(ctx context.Context, s logical.Storage, mkv managed_key.PkiManagedKeyView, issuerId IssuerID, usage IssuerUsage) (*certutil.CAInfoBundle, error) {
	entry, bundle, err := FetchCertBundleByIssuerId(ctx, s, issuerId, true)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return nil, err
		case errutil.InternalError:
			return nil, err
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf("error fetching CA info: %v", err)}
		}
	}

	if err = entry.EnsureUsage(usage); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error while attempting to use issuer %v: %v", issuerId, err)}
	}

	parsedBundle, err := ParseCABundle(ctx, mkv, bundle)
	if err != nil {
		return nil, errutil.InternalError{Err: err.Error()}
	}

	if parsedBundle.Certificate == nil {
		return nil, errutil.InternalError{Err: "stored CA information not able to be parsed"}
	}
	if parsedBundle.PrivateKey == nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("unable to fetch corresponding key for issuer %v; unable to use this issuer for signing", issuerId)}
	}

	caInfo := &certutil.CAInfoBundle{
		ParsedCertBundle:     *parsedBundle,
		URLs:                 nil,
		LeafNotAfterBehavior: entry.LeafNotAfterBehavior,
		RevocationSigAlg:     entry.RevocationSigAlg,
	}

	entries, err := GetAIAURLs(ctx, s, entry)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch AIA URL information: %v", err)}
	}
	caInfo.URLs = entries

	return caInfo, nil
}

func ParseCABundle(ctx context.Context, mkv managed_key.PkiManagedKeyView, bundle *certutil.CertBundle) (*certutil.ParsedCertBundle, error) {
	if bundle.PrivateKeyType == certutil.ManagedPrivateKey {
		return managed_key.ParseManagedKeyCABundle(ctx, mkv, bundle)
	}
	return bundle.ToParsedCertBundle()
}

// FetchCertBundleByIssuerId builds a certutil.CertBundle from the specified issuer identifier,
// optionally loading the key or not. This method supports loading legacy
// bundles using the legacyBundleShimID issuerId, and if no entry is found will return an error.
func FetchCertBundleByIssuerId(ctx context.Context, s logical.Storage, id IssuerID, loadKey bool) (*IssuerEntry, *certutil.CertBundle, error) {
	if id == LegacyBundleShimID {
		// We have not completed the migration, or started a request in legacy mode, so
		// attempt to load the bundle from the legacy location
		issuer, bundle, err := GetLegacyCertBundle(ctx, s)
		if err != nil {
			return nil, nil, err
		}
		if issuer == nil || bundle == nil {
			return nil, nil, errutil.UserError{Err: "no legacy cert bundle exists"}
		}

		return issuer, bundle, err
	}

	issuer, err := FetchIssuerById(ctx, s, id)
	if err != nil {
		return nil, nil, err
	}

	var bundle certutil.CertBundle
	bundle.Certificate = issuer.Certificate
	bundle.CAChain = issuer.CAChain
	bundle.SerialNumber = issuer.SerialNumber

	// Fetch the key if it exists. Sometimes we don't need the key immediately.
	if loadKey && issuer.KeyID != KeyID("") {
		key, err := FetchKeyById(ctx, s, issuer.KeyID)
		if err != nil {
			return nil, nil, err
		}

		bundle.PrivateKeyType = key.PrivateKeyType
		bundle.PrivateKey = key.PrivateKey
	}

	return issuer, &bundle, nil
}

func GetLegacyCertBundle(ctx context.Context, s logical.Storage) (*IssuerEntry, *certutil.CertBundle, error) {
	entry, err := s.Get(ctx, LegacyCertBundlePath)
	if err != nil {
		return nil, nil, err
	}

	if entry == nil {
		return nil, nil, nil
	}

	cb := &certutil.CertBundle{}
	err = entry.DecodeJSON(cb)
	if err != nil {
		return nil, nil, err
	}

	// Fake a storage entry with backwards compatibility in mind.
	issuer := &IssuerEntry{
		ID:                   LegacyBundleShimID,
		KeyID:                LegacyBundleShimKeyID,
		Name:                 "legacy-entry-shim",
		Certificate:          cb.Certificate,
		CAChain:              cb.CAChain,
		SerialNumber:         cb.SerialNumber,
		LeafNotAfterBehavior: certutil.ErrNotAfterBehavior,
	}
	issuer.Usage.ToggleUsage(AllIssuerUsages)

	return issuer, cb, nil
}

func ResolveIssuerCRLPath(ctx context.Context, storage logical.Storage, useLegacyBundleCaStorage bool, reference string, unified bool) (string, error) {
	if useLegacyBundleCaStorage {
		return "crl", nil
	}

	issuer, err := ResolveIssuerReference(ctx, storage, reference)
	if err != nil {
		return "crl", err
	}

	var crlConfig *InternalCRLConfigEntry
	if unified {
		crlConfig, err = GetUnifiedCRLConfig(ctx, storage)
		if err != nil {
			return "crl", err
		}
	} else {
		crlConfig, err = GetLocalCRLConfig(ctx, storage)
		if err != nil {
			return "crl", err
		}
	}

	if crlId, ok := crlConfig.IssuerIDCRLMap[issuer]; ok && len(crlId) > 0 {
		path := fmt.Sprintf("crls/%v", crlId)
		if unified {
			path = ("unified-") + path
		}

		return path, nil
	}

	return "crl", fmt.Errorf("unable to find CRL for issuer: id:%v/ref:%v", issuer, reference)
}
