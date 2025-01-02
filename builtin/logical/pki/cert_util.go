// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"crypto"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/cryptobyte"
	cbbasn1 "golang.org/x/crypto/cryptobyte/asn1"
)

type inputBundle struct {
	role    *issuing.RoleEntry
	req     *logical.Request
	apiData *framework.FieldData
}

var (
	// labelRegex is a single label from a valid domain name and was extracted
	// from hostnameRegex below for use in leftWildLabelRegex, without any
	// label separators (`.`).
	labelRegex = `([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])`

	// A note on hostnameRegex: although we set the StrictDomainName option
	// when doing the idna conversion, this appears to only affect output, not
	// input, so it will allow e.g. host^123.example.com straight through. So
	// we still need to use this to check the output.
	hostnameRegex = regexp.MustCompile(`^(\*\.)?(` + labelRegex + `\.)*` + labelRegex + `\.?$`)

	// Left Wildcard Label Regex is equivalent to a single domain label
	// component from hostnameRegex above, but with additional wildcard
	// characters added. There are four possibilities here:
	//
	//  1. Entire label is a wildcard,
	//  2. Wildcard exists at the start,
	//  3. Wildcard exists at the end,
	//  4. Wildcard exists in the middle.
	allWildRegex       = `\*`
	startWildRegex     = `\*` + labelRegex
	endWildRegex       = labelRegex + `\*`
	middleWildRegex    = labelRegex + `\*` + labelRegex
	leftWildLabelRegex = regexp.MustCompile(`^(` + allWildRegex + `|` + startWildRegex + `|` + endWildRegex + `|` + middleWildRegex + `)$`)

	// Cloned from https://github.com/golang/go/blob/82c713feb05da594567631972082af2fcba0ee4f/src/crypto/x509/x509.go#L327-L379
	oidSignatureMD2WithRSA      = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 2}
	oidSignatureMD5WithRSA      = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 4}
	oidSignatureSHA1WithRSA     = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 5}
	oidSignatureSHA256WithRSA   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 11}
	oidSignatureSHA384WithRSA   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 12}
	oidSignatureSHA512WithRSA   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 13}
	oidSignatureRSAPSS          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 10}
	oidSignatureDSAWithSHA1     = asn1.ObjectIdentifier{1, 2, 840, 10040, 4, 3}
	oidSignatureDSAWithSHA256   = asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 3, 2}
	oidSignatureECDSAWithSHA1   = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 1}
	oidSignatureECDSAWithSHA256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 2}
	oidSignatureECDSAWithSHA384 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 3}
	oidSignatureECDSAWithSHA512 = asn1.ObjectIdentifier{1, 2, 840, 10045, 4, 3, 4}
	oidSignatureEd25519         = asn1.ObjectIdentifier{1, 3, 101, 112}
	oidISOSignatureSHA1WithRSA  = asn1.ObjectIdentifier{1, 3, 14, 3, 2, 29}

	signatureAlgorithmDetails = []struct {
		algo       x509.SignatureAlgorithm
		name       string
		oid        asn1.ObjectIdentifier
		pubKeyAlgo x509.PublicKeyAlgorithm
		hash       crypto.Hash
	}{
		{x509.MD2WithRSA, "MD2-RSA", oidSignatureMD2WithRSA, x509.RSA, crypto.Hash(0) /* no value for MD2 */},
		{x509.MD5WithRSA, "MD5-RSA", oidSignatureMD5WithRSA, x509.RSA, crypto.MD5},
		{x509.SHA1WithRSA, "SHA1-RSA", oidSignatureSHA1WithRSA, x509.RSA, crypto.SHA1},
		{x509.SHA1WithRSA, "SHA1-RSA", oidISOSignatureSHA1WithRSA, x509.RSA, crypto.SHA1},
		{x509.SHA256WithRSA, "SHA256-RSA", oidSignatureSHA256WithRSA, x509.RSA, crypto.SHA256},
		{x509.SHA384WithRSA, "SHA384-RSA", oidSignatureSHA384WithRSA, x509.RSA, crypto.SHA384},
		{x509.SHA512WithRSA, "SHA512-RSA", oidSignatureSHA512WithRSA, x509.RSA, crypto.SHA512},
		{x509.SHA256WithRSAPSS, "SHA256-RSAPSS", oidSignatureRSAPSS, x509.RSA, crypto.SHA256},
		{x509.SHA384WithRSAPSS, "SHA384-RSAPSS", oidSignatureRSAPSS, x509.RSA, crypto.SHA384},
		{x509.SHA512WithRSAPSS, "SHA512-RSAPSS", oidSignatureRSAPSS, x509.RSA, crypto.SHA512},
		{x509.DSAWithSHA1, "DSA-SHA1", oidSignatureDSAWithSHA1, x509.DSA, crypto.SHA1},
		{x509.DSAWithSHA256, "DSA-SHA256", oidSignatureDSAWithSHA256, x509.DSA, crypto.SHA256},
		{x509.ECDSAWithSHA1, "ECDSA-SHA1", oidSignatureECDSAWithSHA1, x509.ECDSA, crypto.SHA1},
		{x509.ECDSAWithSHA256, "ECDSA-SHA256", oidSignatureECDSAWithSHA256, x509.ECDSA, crypto.SHA256},
		{x509.ECDSAWithSHA384, "ECDSA-SHA384", oidSignatureECDSAWithSHA384, x509.ECDSA, crypto.SHA384},
		{x509.ECDSAWithSHA512, "ECDSA-SHA512", oidSignatureECDSAWithSHA512, x509.ECDSA, crypto.SHA512},
		{x509.PureEd25519, "Ed25519", oidSignatureEd25519, x509.Ed25519, crypto.Hash(0) /* no pre-hashing */},
	}
)

func doesPublicKeyAlgoMatchSignatureAlgo(pubKey x509.PublicKeyAlgorithm, algo x509.SignatureAlgorithm) bool {
	for _, detail := range signatureAlgorithmDetails {
		if detail.algo == algo {
			return pubKey == detail.pubKeyAlgo
		}
	}

	return false
}

func getFormat(data *framework.FieldData) string {
	format := data.Get("format").(string)
	switch format {
	case "pem":
	case "der":
	case "pem_bundle":
	default:
		format = ""
	}
	return format
}

// fetchCAInfo will fetch the CA info, will return an error if no ca info exists, this does NOT support
// loading using the legacyBundleShimID and should be used with care. This should be called only once
// within the request path otherwise you run the risk of a race condition with the issuer migration on perf-secondaries.
func (sc *storageContext) fetchCAInfo(issuerRef string, usage issuing.IssuerUsage) (*certutil.CAInfoBundle, error) {
	bundle, _, err := sc.fetchCAInfoWithIssuer(issuerRef, usage)
	return bundle, err
}

func (sc *storageContext) fetchCAInfoWithIssuer(issuerRef string, usage issuing.IssuerUsage) (*certutil.CAInfoBundle, issuing.IssuerID, error) {
	var issuerId issuing.IssuerID

	if sc.UseLegacyBundleCaStorage() {
		// We have not completed the migration so attempt to load the bundle from the legacy location
		sc.Logger().Info("Using legacy CA bundle as PKI migration has not completed.")
		issuerId = legacyBundleShimID
	} else {
		var err error
		issuerId, err = sc.resolveIssuerReference(issuerRef)
		if err != nil {
			// Usually a bad label from the user or mis-configured default.
			return nil, issuing.IssuerRefNotFound, errutil.UserError{Err: err.Error()}
		}
	}

	bundle, err := sc.fetchCAInfoByIssuerId(issuerId, usage)
	if err != nil {
		return nil, issuing.IssuerRefNotFound, err
	}

	return bundle, issuerId, nil
}

// fetchCAInfoByIssuerId will fetch the CA info, will return an error if no ca info exists for the given issuerId.
// This does support the loading using the legacyBundleShimID
func (sc *storageContext) fetchCAInfoByIssuerId(issuerId issuing.IssuerID, usage issuing.IssuerUsage) (*certutil.CAInfoBundle, error) {
	return issuing.FetchCAInfoByIssuerId(sc.Context, sc.Storage, sc.GetPkiManagedView(), issuerId, usage)
}

func fetchCertBySerialBigInt(sc *storageContext, prefix string, serial *big.Int) (*logical.StorageEntry, error) {
	return fetchCertBySerial(sc, prefix, serialFromBigInt(serial))
}

// fetchCertBySerial allows fetching certificates from the backend; it handles the slightly
// separate pathing for CRL, and revoked certificates.
//
// Support for fetching CA certificates was removed, due to the new issuers
// changes.
func fetchCertBySerial(sc pki_backend.StorageContext, prefix, serial string) (*logical.StorageEntry, error) {
	var path, legacyPath string
	var err error
	var certEntry *logical.StorageEntry

	hyphenSerial := parsing.NormalizeSerialForStorage(serial)
	colonSerial := strings.ReplaceAll(strings.ToLower(serial), "-", ":")

	switch {
	// Revoked goes first as otherwise crl get hardcoded paths which fail if
	// we actually want revocation info
	case strings.HasPrefix(prefix, "revoked/"):
		legacyPath = "revoked/" + colonSerial
		path = "revoked/" + hyphenSerial
	case serial == issuing.LegacyCRLPath || serial == issuing.DeltaCRLPath || serial == issuing.UnifiedCRLPath || serial == issuing.UnifiedDeltaCRLPath:
		warnings, err := sc.CrlBuilder().RebuildIfForced(sc)
		if err != nil {
			return nil, err
		}
		if len(warnings) > 0 {
			msg := "During rebuild of CRL for cert fetch, got the following warnings:"
			for index, warning := range warnings {
				msg = fmt.Sprintf("%v\n %d. %v", msg, index+1, warning)
			}
			sc.Logger().Warn(msg)
		}

		unified := serial == issuing.UnifiedCRLPath || serial == issuing.UnifiedDeltaCRLPath
		path, err = issuing.ResolveIssuerCRLPath(sc.GetContext(), sc.GetStorage(), sc.UseLegacyBundleCaStorage(), issuing.DefaultRef, unified)
		if err != nil {
			return nil, err
		}

		if serial == issuing.DeltaCRLPath || serial == issuing.UnifiedDeltaCRLPath {
			if sc.UseLegacyBundleCaStorage() {
				return nil, fmt.Errorf("refusing to serve delta CRL with legacy CA bundle")
			}

			path += issuing.DeltaCRLPathSuffix
		}
	default:
		legacyPath = issuing.PathCerts + colonSerial
		path = issuing.PathCerts + hyphenSerial
	}

	certEntry, err = sc.GetStorage().Get(sc.GetContext(), path)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error fetching certificate %s: %s", serial, err)}
	}
	if certEntry != nil {
		if certEntry.Value == nil || len(certEntry.Value) == 0 {
			return nil, errutil.InternalError{Err: fmt.Sprintf("returned certificate bytes for serial %s were empty", serial)}
		}
		return certEntry, nil
	}

	// If legacyPath is unset, it's going to be a CA or CRL; return immediately
	if legacyPath == "" {
		return nil, nil
	}

	// Retrieve the old-style path.  We disregard errors here because they
	// always manifest on Windows, and thus the initial check for a revoked
	// cert fails would return an error when the cert isn't revoked, preventing
	// the happy path from working.
	certEntry, _ = sc.GetStorage().Get(sc.GetContext(), legacyPath)
	if certEntry == nil {
		return nil, nil
	}
	if certEntry.Value == nil || len(certEntry.Value) == 0 {
		return nil, errutil.InternalError{Err: fmt.Sprintf("returned certificate bytes for serial %s were empty", serial)}
	}

	// Update old-style paths to new-style paths
	certEntry.Key = path
	certCounter := sc.GetCertificateCounter()
	certsCounted := certCounter.IsInitialized()
	if err = sc.GetStorage().Put(sc.GetContext(), certEntry); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error saving certificate with serial %s to new location: %s", serial, err)}
	}
	if err = sc.GetStorage().Delete(sc.GetContext(), legacyPath); err != nil {
		// If we fail here, we have an extra (copy) of a cert in storage, add to metrics:
		switch {
		case strings.HasPrefix(prefix, "revoked/"):
			certCounter.IncrementTotalRevokedCertificatesCount(certsCounted, path)
		default:
			certCounter.IncrementTotalCertificatesCount(certsCounted, path)
		}
		return nil, errutil.InternalError{Err: fmt.Sprintf("error deleting certificate with serial %s from old location: %s", serial, err)}
	}

	return certEntry, nil
}

// Given a URI SAN, verify that it is allowed.
func validateURISAN(b *backend, data *inputBundle, uri string) bool {
	entityInfo := issuing.NewEntityInfoFromReq(data.req)
	return issuing.ValidateURISAN(b.System(), data.role, entityInfo, uri)
}

// Validates a given common name, ensuring it's either an email or a hostname
// after validating it according to the role parameters, or disables
// validation altogether.
func validateCommonName(b *backend, data *inputBundle, name string) string {
	entityInfo := issuing.NewEntityInfoFromReq(data.req)
	return issuing.ValidateCommonName(b.System(), data.role, entityInfo, name)
}

func isWildcardDomain(name string) bool {
	return issuing.IsWildcardDomain(name)
}

func validateWildcardDomain(name string) (string, string, error) {
	return issuing.ValidateWildcardDomain(name)
}

// Given a set of requested names for a certificate, verifies that all of them
// match the various toggles set in the role for controlling issuance.
// If one does not pass, it is returned in the string argument.
func validateNames(b *backend, data *inputBundle, names []string) string {
	entityInfo := issuing.NewEntityInfoFromReq(data.req)
	return issuing.ValidateNames(b.System(), data.role, entityInfo, names)
}

// validateOtherSANs checks if the values requested are allowed. If an OID
// isn't allowed, it will be returned as the first string. If a value isn't
// allowed, it will be returned as the second string. Empty strings + error
// means everything is okay.
func validateOtherSANs(data *inputBundle, requested map[string][]string) (string, string, error) {
	return issuing.ValidateOtherSANs(data.role, requested)
}

func parseOtherSANs(others []string) (map[string][]string, error) {
	return issuing.ParseOtherSANs(others)
}

// Returns bool stating whether the given UserId is Valid
func validateUserId(data *inputBundle, userId string) bool {
	return issuing.ValidateUserId(data.role, userId)
}

func validateSerialNumber(data *inputBundle, serialNumber string) string {
	return issuing.ValidateSerialNumber(data.role, serialNumber)
}

func generateCert(sc *storageContext,
	input *inputBundle,
	caSign *certutil.CAInfoBundle,
	isCA bool,
	randomSource io.Reader) (*certutil.ParsedCertBundle, []string, error,
) {
	ctx := sc.Context

	if input.role == nil {
		return nil, nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	if input.role.KeyType == "rsa" && input.role.KeyBits < 2048 {
		return nil, nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
	}

	data, warnings, err := generateCreationBundle(sc.System(), input, caSign, nil)
	if err != nil {
		return nil, nil, err
	}
	if data.Params == nil {
		return nil, nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	if isCA {
		data.Params.IsCA = isCA
		data.Params.PermittedDNSDomains = input.apiData.Get("permitted_dns_domains").([]string)
		data.Params.ExcludedDNSDomains = input.apiData.Get("excluded_dns_domains").([]string)
		data.Params.PermittedIPRanges, err = convertIpRanges(input.apiData.Get("permitted_ip_ranges").([]string))
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf("invalid permitted_ip_ranges value: %s", err)}
		}
		data.Params.ExcludedIPRanges, err = convertIpRanges(input.apiData.Get("excluded_ip_ranges").([]string))
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf("invalid excluded_ip_ranges value: %s", err)}
		}
		data.Params.PermittedEmailAddresses = input.apiData.Get("permitted_email_addresses").([]string)
		data.Params.ExcludedEmailAddresses = input.apiData.Get("excluded_email_addresses").([]string)
		data.Params.PermittedURIDomains = input.apiData.Get("permitted_uri_domains").([]string)
		data.Params.ExcludedURIDomains = input.apiData.Get("excluded_uri_domains").([]string)

		if data.SigningBundle == nil {
			// Generating a self-signed root certificate. Since we have no
			// issuer entry yet, we default to the global URLs.
			entries, err := getGlobalAIAURLs(ctx, sc.Storage)
			if err != nil {
				return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch AIA URL information: %v", err)}
			}

			uris, err := ToURLEntries(sc, issuing.IssuerID(""), entries)
			if err != nil {
				// When generating root issuers, don't err on missing issuer
				// ID; there is little value in including AIA info on a root,
				// as this info would point back to itself; though RFC 5280 is
				// a touch vague on this point, this seems to be consensus
				// from public CAs such as DigiCert Global Root G3, ISRG Root
				// X1, and others.
				//
				// This is a UX bug if we do err here, as it requires AIA
				// templating to not include issuer id (a best practice for
				// child certs issued from root and intermediate mounts
				// however), and setting this before root generation (or, on
				// root renewal) could cause problems.
				if _, nonEmptyIssuerErr := ToURLEntries(sc, issuing.IssuerID("empty-issuer-id"), entries); nonEmptyIssuerErr != nil {
					return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse AIA URL information: %v\nUsing templated AIA URL's {{issuer_id}} field when generating root certificates is not supported.", err)}
				}

				uris = &certutil.URLEntries{}

				msg := "When generating root CA, found global AIA configuration with issuer_id template unsuitable for root generation. This AIA configuration has been ignored. To include AIA on this root CA, set the global AIA configuration to not include issuer_id and instead to refer to a static issuer name."
				warnings = append(warnings, msg)
			}

			data.Params.URLs = uris

			if input.role.MaxPathLength == nil {
				data.Params.MaxPathLength = -1
			} else {
				data.Params.MaxPathLength = *input.role.MaxPathLength
			}
		}
	}

	parsedBundle, err := generateCABundle(sc, input, data, randomSource)
	if err != nil {
		return nil, nil, err
	}

	return parsedBundle, warnings, nil
}

// convertIpRanges parses each string in the input slice as an IP network. Input
// strings are expected to be in the CIDR notation of IP address and prefix length
// like "192.0.2.0/24" or "2001:db8::/32", as defined in RFC 4632 and RFC 4291.
func convertIpRanges(ipRanges []string) ([]*net.IPNet, error) {
	var ret []*net.IPNet
	for _, ipRange := range ipRanges {
		_, ipnet, err := net.ParseCIDR(ipRange)
		if err != nil {
			return nil, fmt.Errorf("error parsing IP range %q: %w", ipRange, err)
		}
		ret = append(ret, ipnet)
	}
	return ret, nil
}

// N.B.: This is only meant to be used for generating intermediate CAs.
// It skips some sanity checks.
func generateIntermediateCSR(sc *storageContext, input *inputBundle, randomSource io.Reader) (*certutil.ParsedCSRBundle, []string, error) {
	creation, warnings, err := generateCreationBundle(sc.System(), input, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	if creation.Params == nil {
		return nil, nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	_, exists := input.apiData.GetOk("key_usage")
	if !exists {
		creation.Params.KeyUsage = 0
	}
	addBasicConstraints := input.apiData != nil && input.apiData.Get("add_basic_constraints").(bool)
	parsedBundle, err := generateCSRBundle(sc, input, creation, addBasicConstraints, randomSource)
	if err != nil {
		return nil, nil, err
	}

	return parsedBundle, warnings, nil
}

func NewSignCertInputFromDataFields(data *framework.FieldData, isCA bool, useCSRValues bool) SignCertInputFromDataFields {
	certBundle := NewCreationBundleInputFromFieldData(data)
	return SignCertInputFromDataFields{
		CreationBundleInputFromFieldData: certBundle,
		data:                             data,
		isCA:                             isCA,
		useCSRValues:                     useCSRValues,
	}
}

type SignCertInputFromDataFields struct {
	CreationBundleInputFromFieldData
	data         *framework.FieldData
	isCA         bool
	useCSRValues bool
}

var _ issuing.SignCertInput = SignCertInputFromDataFields{}

func (i SignCertInputFromDataFields) GetCSR() (*x509.CertificateRequest, error) {
	csrString := i.data.Get("csr").(string)
	if csrString == "" {
		return nil, errutil.UserError{Err: "\"csr\" is empty"}
	}

	pemBlock, _ := pem.Decode([]byte(csrString))
	if pemBlock == nil {
		return nil, errutil.UserError{Err: "csr contains no data"}
	}
	csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("certificate request could not be parsed: %v", err)}
	}

	return csr, nil
}

func (i SignCertInputFromDataFields) IsCA() bool {
	return i.isCA
}

func (i SignCertInputFromDataFields) UseCSRValues() bool {
	return i.useCSRValues
}

func (i SignCertInputFromDataFields) GetPermittedDomains() []string {
	return i.data.Get("permitted_dns_domains").([]string)
}

func (i SignCertInputFromDataFields) GetExcludedDomains() []string {
	return i.data.Get("excluded_dns_domains").([]string)
}

func (i SignCertInputFromDataFields) GetPermittedIpRanges() ([]*net.IPNet, error) {
	return convertIpRanges(i.data.Get("permitted_ip_ranges").([]string))
}

func (i SignCertInputFromDataFields) GetExcludedIpRanges() ([]*net.IPNet, error) {
	return convertIpRanges(i.data.Get("excluded_ip_ranges").([]string))
}

func (i SignCertInputFromDataFields) GetPermittedEmailAddresses() []string {
	return i.data.Get("permitted_email_addresses").([]string)
}

func (i SignCertInputFromDataFields) GetExcludedEmailAddresses() []string {
	return i.data.Get("excluded_email_addresses").([]string)
}

func (i SignCertInputFromDataFields) GetPermittedUriDomains() []string {
	return i.data.Get("permitted_uri_domains").([]string)
}

func (i SignCertInputFromDataFields) GetExcludedUriDomains() []string {
	return i.data.Get("excluded_uri_domains").([]string)
}

func (i SignCertInputFromDataFields) IgnoreCSRSignature() bool {
	return false
}

func signCert(sysView logical.SystemView, data *inputBundle, caSign *certutil.CAInfoBundle, isCA bool, useCSRValues bool) (*certutil.ParsedCertBundle, []string, error) {
	if data.role == nil {
		return nil, nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	entityInfo := issuing.NewEntityInfoFromReq(data.req)
	signCertInput := NewSignCertInputFromDataFields(data.apiData, isCA, useCSRValues)

	return issuing.SignCert(sysView, data.role, entityInfo, caSign, signCertInput)
}

func getOtherSANsFromX509Extensions(exts []pkix.Extension) ([]certutil.OtherNameUtf8, error) {
	return certutil.GetOtherSANsFromX509Extensions(exts)
}

var _ issuing.CreationBundleInput = CreationBundleInputFromFieldData{}

func NewCreationBundleInputFromFieldData(data *framework.FieldData) CreationBundleInputFromFieldData {
	certNotAfter := NewCertNotAfterInputFromFieldData(data)
	return CreationBundleInputFromFieldData{
		CertNotAfterInputFromFieldData: certNotAfter,
		data:                           data,
	}
}

type CreationBundleInputFromFieldData struct {
	CertNotAfterInputFromFieldData
	data *framework.FieldData
}

func (cb CreationBundleInputFromFieldData) IgnoreCSRSignature() bool {
	return false
}

func (cb CreationBundleInputFromFieldData) GetCommonName() string {
	return cb.data.Get("common_name").(string)
}

func (cb CreationBundleInputFromFieldData) GetSerialNumber() string {
	return cb.data.Get("serial_number").(string)
}

func (cb CreationBundleInputFromFieldData) GetExcludeCnFromSans() bool {
	return cb.data.Get("exclude_cn_from_sans").(bool)
}

func (cb CreationBundleInputFromFieldData) GetOptionalAltNames() (interface{}, bool) {
	return cb.data.GetOk("alt_names")
}

func (cb CreationBundleInputFromFieldData) GetOtherSans() []string {
	return cb.data.Get("other_sans").([]string)
}

func (cb CreationBundleInputFromFieldData) GetIpSans() []string {
	return cb.data.Get("ip_sans").([]string)
}

func (cb CreationBundleInputFromFieldData) GetURISans() []string {
	return cb.data.Get("uri_sans").([]string)
}

func (cb CreationBundleInputFromFieldData) GetOptionalSkid() (interface{}, bool) {
	return cb.data.GetOk("skid")
}

func (cb CreationBundleInputFromFieldData) IsUserIdInSchema() (interface{}, bool) {
	val, present := cb.data.Schema["user_ids"]
	return val, present
}

func (cb CreationBundleInputFromFieldData) GetUserIds() []string {
	return cb.data.Get("user_ids").([]string)
}

// generateCreationBundle is a shared function that reads parameters supplied
// from the various endpoints and generates a CreationParameters with the
// parameters that can be used to issue or sign
func generateCreationBundle(sysView logical.SystemView, data *inputBundle, caSign *certutil.CAInfoBundle, csr *x509.CertificateRequest) (*certutil.CreationBundle, []string, error) {
	entityInfo := issuing.NewEntityInfoFromReq(data.req)
	creationBundleInput := NewCreationBundleInputFromFieldData(data.apiData)

	return issuing.GenerateCreationBundle(sysView, data.role, entityInfo, creationBundleInput, caSign, csr)
}

// getCertificateNotAfter compute a certificate's NotAfter date based on the mount ttl, role, signing bundle and input
// api data being sent. Returns a NotAfter time, a set of warnings or an error.
func getCertificateNotAfter(sysView logical.SystemView, data *inputBundle, caSign *certutil.CAInfoBundle) (time.Time, []string, error) {
	input := NewCertNotAfterInputFromFieldData(data.apiData)
	return issuing.GetCertificateNotAfter(sysView, data.role, input, caSign)
}

// applyIssuerLeafNotAfterBehavior resets a certificate's notAfter time or errors out based on the
// issuer's notAfter date along with the LeafNotAfterBehavior configuration
func applyIssuerLeafNotAfterBehavior(caSign *certutil.CAInfoBundle, notAfter time.Time) (time.Time, error) {
	return issuing.ApplyIssuerLeafNotAfterBehavior(caSign, notAfter)
}

func convertRespToPKCS8(resp *logical.Response) error {
	privRaw, ok := resp.Data["private_key"]
	if !ok {
		return nil
	}
	priv, ok := privRaw.(string)
	if !ok {
		return fmt.Errorf("error converting response to pkcs8: could not parse original value as string")
	}

	privKeyTypeRaw, ok := resp.Data["private_key_type"]
	if !ok {
		return fmt.Errorf("error converting response to pkcs8: %q not found in response", "private_key_type")
	}
	privKeyType, ok := privKeyTypeRaw.(certutil.PrivateKeyType)
	if !ok {
		return fmt.Errorf("error converting response to pkcs8: could not parse original type value as string")
	}

	var keyData []byte
	var pemUsed bool
	var err error
	var signer crypto.Signer

	block, _ := pem.Decode([]byte(priv))
	if block == nil {
		keyData, err = base64.StdEncoding.DecodeString(priv)
		if err != nil {
			return fmt.Errorf("error converting response to pkcs8: error decoding original value: %w", err)
		}
	} else {
		keyData = block.Bytes
		pemUsed = true
	}

	switch privKeyType {
	case certutil.RSAPrivateKey:
		signer, err = x509.ParsePKCS1PrivateKey(keyData)
	case certutil.ECPrivateKey:
		signer, err = x509.ParseECPrivateKey(keyData)
	case certutil.Ed25519PrivateKey:
		k, err := x509.ParsePKCS8PrivateKey(keyData)
		if err != nil {
			return fmt.Errorf("error converting response to pkcs8: error parsing previous key: %w", err)
		}
		signer = k.(crypto.Signer)
	default:
		return fmt.Errorf("unknown private key type %q", privKeyType)
	}
	if err != nil {
		return fmt.Errorf("error converting response to pkcs8: error parsing previous key: %w", err)
	}

	keyData, err = x509.MarshalPKCS8PrivateKey(signer)
	if err != nil {
		return fmt.Errorf("error converting response to pkcs8: error marshaling pkcs8 key: %w", err)
	}

	if pemUsed {
		block.Type = "PRIVATE KEY"
		block.Bytes = keyData
		resp.Data["private_key"] = strings.TrimSpace(string(pem.EncodeToMemory(block)))
	} else {
		resp.Data["private_key"] = base64.StdEncoding.EncodeToString(keyData)
	}

	return nil
}

func handleOtherCSRSANs(in *x509.CertificateRequest, sans map[string][]string) error {
	certTemplate := &x509.Certificate{
		DNSNames:       in.DNSNames,
		IPAddresses:    in.IPAddresses,
		EmailAddresses: in.EmailAddresses,
		URIs:           in.URIs,
	}
	if err := handleOtherSANs(certTemplate, sans); err != nil {
		return err
	}
	if len(certTemplate.ExtraExtensions) > 0 {
		in.ExtraExtensions = append(in.ExtraExtensions, certTemplate.ExtraExtensions...)
	}
	return nil
}

func handleOtherSANs(in *x509.Certificate, sans map[string][]string) error {
	// If other SANs is empty we return which causes normal Go stdlib parsing
	// of the other SAN types
	if len(sans) == 0 {
		return nil
	}

	var rawValues []asn1.RawValue

	// We need to generate an IMPLICIT sequence for compatibility with OpenSSL
	// -- it's an open question what the default for RFC 5280 actually is, see
	// https://github.com/openssl/openssl/issues/5091 -- so we have to use
	// cryptobyte because using the asn1 package's marshaling always produces
	// an EXPLICIT sequence. Note that asn1 is way too magical according to
	// agl, and cryptobyte is modeled after the CBB/CBS bits that agl put into
	// boringssl.
	for oid, vals := range sans {
		for _, val := range vals {
			var b cryptobyte.Builder
			oidStr, err := stringToOid(oid)
			if err != nil {
				return err
			}
			b.AddASN1ObjectIdentifier(oidStr)
			b.AddASN1(cbbasn1.Tag(0).ContextSpecific().Constructed(), func(b *cryptobyte.Builder) {
				b.AddASN1(cbbasn1.UTF8String, func(b *cryptobyte.Builder) {
					b.AddBytes([]byte(val))
				})
			})
			m, err := b.Bytes()
			if err != nil {
				return err
			}
			rawValues = append(rawValues, asn1.RawValue{Tag: 0, Class: 2, IsCompound: true, Bytes: m})
		}
	}

	// If other SANs is empty we return which causes normal Go stdlib parsing
	// of the other SAN types
	if len(rawValues) == 0 {
		return nil
	}

	// Append any existing SANs, sans marshalling
	rawValues = append(rawValues, marshalSANs(in.DNSNames, in.EmailAddresses, in.IPAddresses, in.URIs)...)

	// Marshal and add to ExtraExtensions
	ext := pkix.Extension{
		// This is the defined OID for subjectAltName
		Id: certutil.OidExtensionSubjectAltName,
	}
	var err error
	ext.Value, err = asn1.Marshal(rawValues)
	if err != nil {
		return err
	}
	in.ExtraExtensions = append(in.ExtraExtensions, ext)

	return nil
}

// Note: Taken from the Go source code since it's not public, and used in the
// modified function below (which also uses these consts upstream)
const (
	nameTypeOther = 0
	nameTypeEmail = 1
	nameTypeDNS   = 2
	nameTypeURI   = 6
	nameTypeIP    = 7
)

// Note: Taken from the Go source code since it's not public, plus changed to not marshal
// marshalSANs marshals a list of addresses into the contents of an X.509
// SubjectAlternativeName extension.
func marshalSANs(dnsNames, emailAddresses []string, ipAddresses []net.IP, uris []*url.URL) []asn1.RawValue {
	var rawValues []asn1.RawValue
	for _, name := range dnsNames {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeDNS, Class: 2, Bytes: []byte(name)})
	}
	for _, email := range emailAddresses {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeEmail, Class: 2, Bytes: []byte(email)})
	}
	for _, rawIP := range ipAddresses {
		// If possible, we always want to encode IPv4 addresses in 4 bytes.
		ip := rawIP.To4()
		if ip == nil {
			ip = rawIP
		}
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeIP, Class: 2, Bytes: ip})
	}
	for _, uri := range uris {
		rawValues = append(rawValues, asn1.RawValue{Tag: nameTypeURI, Class: 2, Bytes: []byte(uri.String())})
	}
	return rawValues
}

func stringToOid(in string) (asn1.ObjectIdentifier, error) {
	split := strings.Split(in, ".")
	ret := make(asn1.ObjectIdentifier, 0, len(split))
	for _, v := range split {
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func parseCertificateFromBytes(certBytes []byte) (*x509.Certificate, error) {
	return parsing.ParseCertificateFromBytes(certBytes)
}

func NewCertNotAfterInputFromFieldData(data *framework.FieldData) CertNotAfterInputFromFieldData {
	return CertNotAfterInputFromFieldData{data: data}
}

var _ issuing.CertNotAfterInput = CertNotAfterInputFromFieldData{}

type CertNotAfterInputFromFieldData struct {
	data *framework.FieldData
}

func (i CertNotAfterInputFromFieldData) GetTTL() int {
	return i.data.Get("ttl").(int)
}

func (i CertNotAfterInputFromFieldData) GetOptionalNotAfter() (interface{}, bool) {
	return i.data.GetOk("not_after")
}
