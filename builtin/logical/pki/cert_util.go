package pki

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/ryanuber/go-glob"
	"golang.org/x/crypto/cryptobyte"
	cbbasn1 "golang.org/x/crypto/cryptobyte/asn1"
	"golang.org/x/net/idna"
)

type inputBundle struct {
	role    *roleEntry
	req     *logical.Request
	apiData *framework.FieldData
}

type caEntity struct {
	// TODO: Have the associated keyEntry here as well.
	issuer   issuerEntry
	caBundle certutil.CAInfoBundle
}

func (cae caEntity) PrettyIssuerId() string {
	return fmt.Sprintf("[id: '%s' name: '%s']", string(cae.issuer.ID), cae.issuer.Name)
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

	// OIDs for X.509 certificate extensions used below.
	oidExtensionSubjectAltName = []int{2, 5, 29, 17}
)

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
func (sc *storageContext) fetchCAInfo(issuerRef string, usage issuerUsage) (*certutil.CAInfoBundle, error) {
	var issuerId issuerID

	if sc.Backend.useLegacyBundleCaStorage() {
		// We have not completed the migration so attempt to load the bundle from the legacy location
		sc.Backend.Logger().Info("Using legacy CA bundle as PKI migration has not completed.")
		issuerId = legacyBundleShimID
	} else {
		var err error
		issuerId, err = sc.resolveIssuerReference(issuerRef)
		if err != nil {
			// Usually a bad label from the user or mis-configured default.
			return nil, errutil.UserError{Err: err.Error()}
		}
	}

	return sc.fetchCAInfoByIssuerId(issuerId, usage)
}

// fetchCAInfoByIssuerId will fetch the CA info, will return an error if no ca info exists for the given issuerId.
// This does support the loading using the legacyBundleShimID
func (sc *storageContext) fetchCAInfoByIssuerId(issuerId issuerID, usage issuerUsage) (*certutil.CAInfoBundle, error) {
	entity, err := sc.fetchCAEntityByIssuerId(issuerId, usage)
	if err != nil {
		return nil, err
	}
	return &entity.caBundle, nil
}

func (sc *storageContext) fetchCAEntityByIssuerId(issuerId issuerID, usage issuerUsage) (caEntity, error) {
	var entity caEntity
	entry, bundle, err := sc.fetchCertBundleByIssuerId(issuerId, true)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return entity, err
		case errutil.InternalError:
			return entity, err
		default:
			return entity, errutil.InternalError{Err: fmt.Sprintf("error fetching CA info: %v", err)}
		}
	}

	if err := entry.EnsureUsage(usage); err != nil {
		return entity, errutil.InternalError{Err: fmt.Sprintf("error while attempting to use issuer %v: %v", issuerId, err)}
	}

	parsedBundle, err := parseCABundle(sc.Context, sc.Backend, bundle)
	if err != nil {
		return entity, errutil.InternalError{Err: err.Error()}
	}

	if parsedBundle.Certificate == nil {
		return entity, errutil.InternalError{Err: "stored CA information not able to be parsed"}
	}
	if parsedBundle.PrivateKey == nil {
		return entity, errutil.UserError{Err: fmt.Sprintf("unable to fetch corresponding key for issuer %v; unable to use this issuer for signing", issuerId)}
	}

	caInfo := &certutil.CAInfoBundle{
		ParsedCertBundle:     *parsedBundle,
		URLs:                 nil,
		LeafNotAfterBehavior: entry.LeafNotAfterBehavior,
		RevocationSigAlg:     entry.RevocationSigAlg,
	}

	entries, err := entry.GetAIAURLs(sc)
	if err != nil {
		return entity, errutil.InternalError{Err: fmt.Sprintf("unable to fetch URL information: %v", err)}
	}
	caInfo.URLs = entries

	return caEntity{
		issuer:   *entry,
		caBundle: *caInfo,
	}, nil
}

func fetchCertBySerialBigInt(ctx context.Context, b *backend, req *logical.Request, prefix string, serial *big.Int) (*logical.StorageEntry, error) {
	return fetchCertBySerial(ctx, b, req, prefix, serialFromBigInt(serial))
}

// Allows fetching certificates from the backend; it handles the slightly
// separate pathing for CRL, and revoked certificates.
//
// Support for fetching CA certificates was removed, due to the new issuers
// changes.
func fetchCertBySerial(ctx context.Context, b *backend, req *logical.Request, prefix, serial string) (*logical.StorageEntry, error) {
	var path, legacyPath string
	var err error
	var certEntry *logical.StorageEntry

	hyphenSerial := normalizeSerial(serial)
	colonSerial := strings.ReplaceAll(strings.ToLower(serial), "-", ":")

	switch {
	// Revoked goes first as otherwise crl get hardcoded paths which fail if
	// we actually want revocation info
	case strings.HasPrefix(prefix, "revoked/"):
		legacyPath = "revoked/" + colonSerial
		path = "revoked/" + hyphenSerial
	case serial == legacyCRLPath || serial == deltaCRLPath:
		if err = b.crlBuilder.rebuildIfForced(ctx, b, req); err != nil {
			return nil, err
		}
		sc := b.makeStorageContext(ctx, req.Storage)
		path, err = sc.resolveIssuerCRLPath(defaultRef)
		if err != nil {
			return nil, err
		}

		if serial == deltaCRLPath {
			if sc.Backend.useLegacyBundleCaStorage() {
				return nil, fmt.Errorf("refusing to serve delta CRL with legacy CA bundle")
			}

			path += deltaCRLPathSuffix
		}
	default:
		legacyPath = "certs/" + colonSerial
		path = "certs/" + hyphenSerial
	}

	certEntry, err = req.Storage.Get(ctx, path)
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
	certEntry, _ = req.Storage.Get(ctx, legacyPath)
	if certEntry == nil {
		return nil, nil
	}
	if certEntry.Value == nil || len(certEntry.Value) == 0 {
		return nil, errutil.InternalError{Err: fmt.Sprintf("returned certificate bytes for serial %s were empty", serial)}
	}

	// Update old-style paths to new-style paths
	certEntry.Key = path
	certsCounted := b.certsCounted.Load()
	if err = req.Storage.Put(ctx, certEntry); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error saving certificate with serial %s to new location", serial)}
	}
	if err = req.Storage.Delete(ctx, legacyPath); err != nil {
		// If we fail here, we have an extra (copy) of a cert in storage, add to metrics:
		switch {
		case strings.HasPrefix(prefix, "revoked/"):
			b.incrementTotalRevokedCertificatesCount(certsCounted, path)
		default:
			b.incrementTotalCertificatesCount(certsCounted, path)
		}
		return nil, errutil.InternalError{Err: fmt.Sprintf("error deleting certificate with serial %s from old location", serial)}
	}

	return certEntry, nil
}

// Given a URI SAN, verify that it is allowed.
func validateURISAN(b *backend, data *inputBundle, uri string) bool {
	valid := false
	for _, allowed := range data.role.AllowedURISANs {
		if data.role.AllowedURISANsTemplate {
			isTemplate, _ := framework.ValidateIdentityTemplate(allowed)
			if isTemplate && data.req.EntityID != "" {
				tmpAllowed, err := framework.PopulateIdentityTemplate(allowed, data.req.EntityID, b.System())
				if err != nil {
					continue
				}
				allowed = tmpAllowed
			}
		}
		validURI := glob.Glob(allowed, uri)
		if validURI {
			valid = true
			break
		}
	}
	return valid
}

// Validates a given common name, ensuring it's either an email or a hostname
// after validating it according to the role parameters, or disables
// validation altogether.
func validateCommonName(b *backend, data *inputBundle, name string) string {
	isDisabled := len(data.role.CNValidations) == 1 && data.role.CNValidations[0] == "disabled"
	if isDisabled {
		return ""
	}

	if validateNames(b, data, []string{name}) != "" {
		return name
	}

	// Validations weren't disabled, but the role lacked CN Validations, so
	// don't restrict types. This case is hit in certain existing tests.
	if len(data.role.CNValidations) == 0 {
		return ""
	}

	// If there's an at in the data, ensure email type validation is allowed.
	// Otherwise, ensure hostname is allowed.
	if strings.Contains(name, "@") {
		var allowsEmails bool
		for _, validation := range data.role.CNValidations {
			if validation == "email" {
				allowsEmails = true
				break
			}
		}
		if !allowsEmails {
			return name
		}
	} else {
		var allowsHostnames bool
		for _, validation := range data.role.CNValidations {
			if validation == "hostname" {
				allowsHostnames = true
				break
			}
		}
		if !allowsHostnames {
			return name
		}
	}

	return ""
}

// Given a set of requested names for a certificate, verifies that all of them
// match the various toggles set in the role for controlling issuance.
// If one does not pass, it is returned in the string argument.
func validateNames(b *backend, data *inputBundle, names []string) string {
	for _, name := range names {
		// Previously, reducedName was called sanitizedName but this made
		// little sense under the previous interpretation of wildcards,
		// leading to two bugs in this implementation. We presently call it
		// "reduced" to indicate that it is still untrusted input (potentially
		// different from the bare Common Name entry we're validating), it
		// might have been modified such as by the removal of wildcard labels
		// or the email prefix.
		reducedName := name
		emailDomain := reducedName
		wildcardLabel := ""
		isEmail := false
		isWildcard := false

		// If it has an @, assume it is an email address and separate out the
		// user from the hostname portion so that we can act on the hostname.
		// Note that this matches behavior from the alt_names parameter. If it
		// ends up being problematic for users, I guess that could be separated
		// into dns_names and email_names in the future to be explicit, but I
		// don't think this is likely.
		if strings.Contains(reducedName, "@") {
			splitEmail := strings.Split(reducedName, "@")
			if len(splitEmail) != 2 {
				return name
			}
			reducedName = splitEmail[1]
			emailDomain = splitEmail[1]
			isEmail = true
		}

		// Per RFC 6125 Section 6.4.3, and explicitly contradicting the earlier
		// RFC 2818 which no modern client will validate against, there are two
		// main types of wildcards, each with a single wildcard specifier (`*`,
		// functionally different from the `*` used as a glob from the
		// AllowGlobDomains parsing path) in the left-most label:
		//
		//  1. Entire label is a single wildcard character (most common and
		//     well-supported),
		//  2. Part of the label contains a single wildcard character (e.g. per
		///    RFC 6125: baz*.example.net, *baz.example.net, or b*z.example.net).
		//
		// We permit issuance of both but not the older RFC 2818 style under
		// the new AllowWildcardCertificates option. However, anything with a
		// glob character is technically a wildcard.
		if strings.Contains(reducedName, "*") {
			// Regardless of later rejections below, this common name contains
			// a wildcard character and is thus technically a wildcard name.
			isWildcard = true

			// Additionally, if AllowWildcardCertificates is explicitly
			// forbidden, it takes precedence over AllowAnyName, thus we should
			// reject the name now.
			//
			// We expect the role to have been correctly migrated but guard for
			// safety.
			if data.role.AllowWildcardCertificates != nil && !*data.role.AllowWildcardCertificates {
				return name
			}

			if strings.Count(reducedName, "*") > 1 {
				// As mentioned above, only one wildcard character is permitted
				// under RFC 6125 semantics.
				return name
			}

			// Split the Common Name into two parts: a left-most label and the
			// remaining segments (if present).
			splitLabels := strings.SplitN(reducedName, ".", 2)
			if len(splitLabels) != 2 {
				// We've been given a single-part domain name that consists
				// entirely of a wildcard. This is a little tricky to handle,
				// but EnforceHostnames validates both the wildcard-containing
				// label and the reduced name, but _only_ the latter if it is
				// non-empty. This allows us to still validate the only label
				// component matches hostname expectations still.
				wildcardLabel = splitLabels[0]
				reducedName = ""
			} else {
				// We have a (at least) two label domain name. But before we can
				// update our names, we need to validate the wildcard ended up
				// in the segment we expected it to. While this is (kinda)
				// validated under EnforceHostnames's leftWildLabelRegex, we
				// still need to validate it in the non-enforced mode.
				//
				// By validated assumption above, we know there's strictly one
				// wildcard in this domain so we only need to check the wildcard
				// label or the reduced name (as one is equivalent to the other).
				// Because we later assume reducedName _lacks_ wildcard segments,
				// we validate that.
				wildcardLabel = splitLabels[0]
				reducedName = splitLabels[1]
				if strings.Contains(reducedName, "*") {
					return name
				}
			}
		}

		// Email addresses using wildcard domain names do not make sense
		// in a Common Name field.
		if isEmail && isWildcard {
			return name
		}

		// AllowAnyName is checked after this because EnforceHostnames still
		// applies when allowing any name. Also, we check the reduced name to
		// ensure that we are not either checking a full email address or a
		// wildcard prefix.
		if data.role.EnforceHostnames {
			if reducedName != "" {
				// See note above about splitLabels having only one segment
				// and setting reducedName to the empty string.
				p := idna.New(
					idna.StrictDomainName(true),
					idna.VerifyDNSLength(true),
				)
				converted, err := p.ToASCII(reducedName)
				if err != nil {
					return name
				}
				if !hostnameRegex.MatchString(converted) {
					return name
				}
			}

			// When a wildcard is specified, we additionally need to validate
			// the label with the wildcard is correctly formed.
			if isWildcard && !leftWildLabelRegex.MatchString(wildcardLabel) {
				return name
			}
		}

		// Self-explanatory, but validations from EnforceHostnames and
		// AllowWildcardCertificates take precedence.
		if data.role.AllowAnyName {
			continue
		}

		// The following blocks all work the same basic way:
		// 1) If a role allows a certain class of base (localhost, token
		// display name, role-configured domains), perform further tests
		//
		// 2) If there is a perfect match on either the sanitized name or it's an
		// email address with a perfect match on the hostname portion, allow it
		//
		// 3) If subdomains are allowed, we check based on the sanitized name;
		// note that if not a wildcard, will be equivalent to the email domain
		// for email checks, and we already checked above for both a wildcard
		// and email address being present in the same name
		// 3a) First we check for a non-wildcard subdomain, as in <name>.<base>
		// 3b) Then we check if it's a wildcard and the base domain is a match
		//
		// Variances are noted in-line

		if data.role.AllowLocalhost {
			if reducedName == "localhost" ||
				reducedName == "localdomain" ||
				(isEmail && emailDomain == "localhost") ||
				(isEmail && emailDomain == "localdomain") {
				continue
			}

			if data.role.AllowSubdomains {
				// It is possible, if unlikely, to have a subdomain of "localhost"
				if strings.HasSuffix(reducedName, ".localhost") ||
					(isWildcard && reducedName == "localhost") {
					continue
				}

				// A subdomain of "localdomain" is also not entirely uncommon
				if strings.HasSuffix(reducedName, ".localdomain") ||
					(isWildcard && reducedName == "localdomain") {
					continue
				}
			}
		}

		if data.role.AllowTokenDisplayName {
			if name == data.req.DisplayName {
				continue
			}

			if data.role.AllowSubdomains {
				if isEmail {
					// If it's an email address, we need to parse the token
					// display name in order to do a proper comparison of the
					// subdomain
					if strings.Contains(data.req.DisplayName, "@") {
						splitDisplay := strings.Split(data.req.DisplayName, "@")
						if len(splitDisplay) == 2 {
							// Compare the sanitized name against the hostname
							// portion of the email address in the broken
							// display name
							if strings.HasSuffix(reducedName, "."+splitDisplay[1]) {
								continue
							}
						}
					}
				}

				if strings.HasSuffix(reducedName, "."+data.req.DisplayName) ||
					(isWildcard && reducedName == data.req.DisplayName) {
					continue
				}
			}
		}

		if len(data.role.AllowedDomains) > 0 {
			valid := false
			for _, currDomain := range data.role.AllowedDomains {
				// If there is, say, a trailing comma, ignore it
				if currDomain == "" {
					continue
				}

				if data.role.AllowedDomainsTemplate {
					isTemplate, _ := framework.ValidateIdentityTemplate(currDomain)
					if isTemplate && data.req.EntityID != "" {
						tmpCurrDomain, err := framework.PopulateIdentityTemplate(currDomain, data.req.EntityID, b.System())
						if err != nil {
							continue
						}

						currDomain = tmpCurrDomain
					}
				}

				// First, allow an exact match of the base domain if that role flag
				// is enabled
				if data.role.AllowBareDomains &&
					(strings.EqualFold(name, currDomain) ||
						(isEmail && strings.EqualFold(emailDomain, currDomain))) {
					valid = true
					break
				}

				if data.role.AllowSubdomains {
					if strings.HasSuffix(reducedName, "."+currDomain) ||
						(isWildcard && strings.EqualFold(reducedName, currDomain)) {
						valid = true
						break
					}
				}

				if data.role.AllowGlobDomains &&
					strings.Contains(currDomain, "*") &&
					glob.Glob(currDomain, name) {
					valid = true
					break
				}
			}

			if valid {
				continue
			}
		}

		return name
	}

	return ""
}

// validateOtherSANs checks if the values requested are allowed. If an OID
// isn't allowed, it will be returned as the first string. If a value isn't
// allowed, it will be returned as the second string. Empty strings + error
// means everything is okay.
func validateOtherSANs(data *inputBundle, requested map[string][]string) (string, string, error) {
	if len(data.role.AllowedOtherSANs) == 1 && data.role.AllowedOtherSANs[0] == "*" {
		// Anything is allowed
		return "", "", nil
	}

	allowed, err := parseOtherSANs(data.role.AllowedOtherSANs)
	if err != nil {
		return "", "", fmt.Errorf("error parsing role's allowed SANs: %w", err)
	}
	for oid, names := range requested {
		for _, name := range names {
			allowedNames, ok := allowed[oid]
			if !ok {
				return oid, "", nil
			}

			valid := false
			for _, allowedName := range allowedNames {
				if glob.Glob(allowedName, name) {
					valid = true
					break
				}
			}

			if !valid {
				return oid, name, nil
			}
		}
	}

	return "", "", nil
}

func parseOtherSANs(others []string) (map[string][]string, error) {
	result := map[string][]string{}
	for _, other := range others {
		splitOther := strings.SplitN(other, ";", 2)
		if len(splitOther) != 2 {
			return nil, fmt.Errorf("expected a semicolon in other SAN %q", other)
		}
		splitType := strings.SplitN(splitOther[1], ":", 2)
		if len(splitType) != 2 {
			return nil, fmt.Errorf("expected a colon in other SAN %q", other)
		}
		switch {
		case strings.EqualFold(splitType[0], "utf8"):
		case strings.EqualFold(splitType[0], "utf-8"):
		default:
			return nil, fmt.Errorf("only utf8 other SANs are supported; found non-supported type in other SAN %q", other)
		}
		result[splitOther[0]] = append(result[splitOther[0]], splitType[1])
	}

	return result, nil
}

func validateSerialNumber(data *inputBundle, serialNumber string) string {
	valid := false
	if len(data.role.AllowedSerialNumbers) > 0 {
		for _, currSerialNumber := range data.role.AllowedSerialNumbers {
			if currSerialNumber == "" {
				continue
			}

			if (strings.Contains(currSerialNumber, "*") &&
				glob.Glob(currSerialNumber, serialNumber)) ||
				currSerialNumber == serialNumber {
				valid = true
				break
			}
		}
	}
	if !valid {
		return serialNumber
	} else {
		return ""
	}
}

func generateCert(sc *storageContext,
	input *inputBundle,
	caSign *certutil.CAInfoBundle,
	isCA bool,
	randomSource io.Reader) (*certutil.ParsedCertBundle, []string, error,
) {
	ctx := sc.Context
	b := sc.Backend

	if input.role == nil {
		return nil, nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	if input.role.KeyType == "rsa" && input.role.KeyBits < 2048 {
		return nil, nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
	}

	data, warnings, err := generateCreationBundle(b, input, caSign, nil)
	if err != nil {
		return nil, nil, err
	}
	if data.Params == nil {
		return nil, nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	if isCA {
		data.Params.IsCA = isCA
		data.Params.PermittedDNSDomains = input.apiData.Get("permitted_dns_domains").([]string)

		if data.SigningBundle == nil {
			// Generating a self-signed root certificate. Since we have no
			// issuer entry yet, we default to the global URLs.
			entries, err := getGlobalAIAURLs(ctx, sc.Storage)
			if err != nil {
				return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch URL information: %v", err)}
			}
			data.Params.URLs = entries

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

// N.B.: This is only meant to be used for generating intermediate CAs.
// It skips some sanity checks.
func generateIntermediateCSR(sc *storageContext, input *inputBundle, randomSource io.Reader) (*certutil.ParsedCSRBundle, []string, error) {
	b := sc.Backend

	creation, warnings, err := generateCreationBundle(b, input, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	if creation.Params == nil {
		return nil, nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	addBasicConstraints := input.apiData != nil && input.apiData.Get("add_basic_constraints").(bool)
	parsedBundle, err := generateCSRBundle(sc, input, creation, addBasicConstraints, randomSource)
	if err != nil {
		return nil, nil, err
	}

	return parsedBundle, warnings, nil
}

func signCert(b *backend,
	data *inputBundle,
	caSign *certutil.CAInfoBundle,
	isCA bool,
	useCSRValues bool) (*certutil.ParsedCertBundle, []string, error,
) {
	if data.role == nil {
		return nil, nil, errutil.InternalError{Err: "no role found in data bundle"}
	}

	csrString := data.apiData.Get("csr").(string)
	if csrString == "" {
		return nil, nil, errutil.UserError{Err: "\"csr\" is empty"}
	}

	pemBlock, _ := pem.Decode([]byte(csrString))
	if pemBlock == nil {
		return nil, nil, errutil.UserError{Err: "csr contains no data"}
	}
	csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		return nil, nil, errutil.UserError{Err: fmt.Sprintf("certificate request could not be parsed: %v", err)}
	}

	if csr.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm || csr.PublicKey == nil {
		return nil, nil, errutil.UserError{Err: "Refusing to sign CSR with empty PublicKey. This usually means the SubjectPublicKeyInfo field has an OID not recognized by Go, such as 1.2.840.113549.1.1.10 for rsaPSS."}
	}

	// This switch validates that the CSR key type matches the role and sets
	// the value in the actualKeyType/actualKeyBits values.
	actualKeyType := ""
	actualKeyBits := 0

	switch data.role.KeyType {
	case "rsa":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.RSA {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				data.role.KeyType)}
		}

		pubKey, ok := csr.PublicKey.(*rsa.PublicKey)
		if !ok {
			return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		actualKeyType = "rsa"
		actualKeyBits = pubKey.N.BitLen()
	case "ec":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.ECDSA {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				data.role.KeyType)}
		}
		pubKey, ok := csr.PublicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		actualKeyType = "ec"
		actualKeyBits = pubKey.Params().BitSize
	case "ed25519":
		// Verify that the key matches the role type
		if csr.PublicKeyAlgorithm != x509.Ed25519 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires keys of type %s",
				data.role.KeyType)}
		}

		_, ok := csr.PublicKey.(ed25519.PublicKey)
		if !ok {
			return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
		}

		actualKeyType = "ed25519"
		actualKeyBits = 0
	case "any":
		// We need to compute the actual key type and key bits, to correctly
		// validate minimums and SignatureBits below.
		switch csr.PublicKeyAlgorithm {
		case x509.RSA:
			pubKey, ok := csr.PublicKey.(*rsa.PublicKey)
			if !ok {
				return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
			}
			if pubKey.N.BitLen() < 2048 {
				return nil, nil, errutil.UserError{Err: "RSA keys < 2048 bits are unsafe and not supported"}
			}

			actualKeyType = "rsa"
			actualKeyBits = pubKey.N.BitLen()
		case x509.ECDSA:
			pubKey, ok := csr.PublicKey.(*ecdsa.PublicKey)
			if !ok {
				return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
			}

			actualKeyType = "ec"
			actualKeyBits = pubKey.Params().BitSize
		case x509.Ed25519:
			_, ok := csr.PublicKey.(ed25519.PublicKey)
			if !ok {
				return nil, nil, errutil.UserError{Err: "could not parse CSR's public key"}
			}

			actualKeyType = "ed25519"
			actualKeyBits = 0
		default:
			return nil, nil, errutil.UserError{Err: "Unknown key type in CSR: " + csr.PublicKeyAlgorithm.String()}
		}
	default:
		return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unsupported key type value: %s", data.role.KeyType)}
	}

	// Before validating key lengths, update our KeyBits/SignatureBits based
	// on the actual CSR key type.
	if data.role.KeyType == "any" {
		// We update the value of KeyBits and SignatureBits here (from the
		// role), using the specified key type. This allows us to convert
		// the default value (0) for SignatureBits and KeyBits to a
		// meaningful value.
		//
		// We ignore the role's original KeyBits value if the KeyType is any
		// as legacy (pre-1.10) roles had default values that made sense only
		// for RSA keys (key_bits=2048) and the older code paths ignored the role value
		// set for KeyBits when KeyType was set to any. This also enforces the
		// docs saying when key_type=any, we only enforce our specified minimums
		// for signing operations
		if data.role.KeyBits, data.role.SignatureBits, err = certutil.ValidateDefaultOrValueKeyTypeSignatureLength(
			actualKeyType, 0, data.role.SignatureBits); err != nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unknown internal error updating default values: %v", err)}
		}

		// We're using the KeyBits field as a minimum value below, and P-224 is safe
		// and a previously allowed value. However, the above call defaults
		// to P-256 as that's a saner default than P-224 (w.r.t. generation), so
		// override it here to allow 224 as the smallest size we permit.
		if actualKeyType == "ec" {
			data.role.KeyBits = 224
		}
	}

	// At this point, data.role.KeyBits and data.role.SignatureBits should both
	// be non-zero, for RSA and ECDSA keys. Validate the actualKeyBits based on
	// the role's values. If the KeyType was any, and KeyBits was set to 0,
	// KeyBits should be updated to 2048 unless some other value was chosen
	// explicitly.
	//
	// This validation needs to occur regardless of the role's key type, so
	// that we always validate both RSA and ECDSA key sizes.
	if actualKeyType == "rsa" {
		if actualKeyBits < data.role.KeyBits {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				data.role.KeyBits, actualKeyBits)}
		}

		if actualKeyBits < 2048 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"Vault requires a minimum of a 2048-bit key, but CSR's key is %d bits",
				actualKeyBits)}
		}
	} else if actualKeyType == "ec" {
		if actualKeyBits < data.role.KeyBits {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"role requires a minimum of a %d-bit key, but CSR's key is %d bits",
				data.role.KeyBits,
				actualKeyBits)}
		}
	}

	creation, warnings, err := generateCreationBundle(b, data, caSign, csr)
	if err != nil {
		return nil, nil, err
	}
	if creation.Params == nil {
		return nil, nil, errutil.InternalError{Err: "nil parameters received from parameter bundle generation"}
	}

	creation.Params.IsCA = isCA
	creation.Params.UseCSRValues = useCSRValues

	if isCA {
		creation.Params.PermittedDNSDomains = data.apiData.Get("permitted_dns_domains").([]string)
	}

	parsedBundle, err := certutil.SignCertificate(creation)
	if err != nil {
		return nil, nil, err
	}

	return parsedBundle, warnings, nil
}

// otherNameRaw describes a name related to a certificate which is not in one
// of the standard name formats. RFC 5280, 4.2.1.6:
//
//	OtherName ::= SEQUENCE {
//	     type-id    OBJECT IDENTIFIER,
//	     value      [0] EXPLICIT ANY DEFINED BY type-id }
type otherNameRaw struct {
	TypeID asn1.ObjectIdentifier
	Value  asn1.RawValue
}

type otherNameUtf8 struct {
	oid   string
	value string
}

// ExtractUTF8String returns the UTF8 string contained in the Value, or an error
// if none is present.
func (oraw *otherNameRaw) extractUTF8String() (*otherNameUtf8, error) {
	svalue := cryptobyte.String(oraw.Value.Bytes)
	var outTag cbbasn1.Tag
	var val cryptobyte.String
	read := svalue.ReadAnyASN1(&val, &outTag)

	if read && outTag == asn1.TagUTF8String {
		return &otherNameUtf8{oid: oraw.TypeID.String(), value: string(val)}, nil
	}
	return nil, fmt.Errorf("no UTF-8 string found in OtherName")
}

func (o otherNameUtf8) String() string {
	return fmt.Sprintf("%s;%s:%s", o.oid, "UTF-8", o.value)
}

func getOtherSANsFromX509Extensions(exts []pkix.Extension) ([]otherNameUtf8, error) {
	var ret []otherNameUtf8
	for _, ext := range exts {
		if !ext.Id.Equal(oidExtensionSubjectAltName) {
			continue
		}
		err := forEachSAN(ext.Value, func(tag int, data []byte) error {
			if tag != 0 {
				return nil
			}

			var other otherNameRaw
			_, err := asn1.UnmarshalWithParams(data, &other, "tag:0")
			if err != nil {
				return fmt.Errorf("could not parse requested other SAN: %w", err)
			}
			val, err := other.extractUTF8String()
			if err != nil {
				return err
			}
			ret = append(ret, *val)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func forEachSAN(extension []byte, callback func(tag int, data []byte) error) error {
	// RFC 5280, 4.2.1.6

	// SubjectAltName ::= GeneralNames
	//
	// GeneralNames ::= SEQUENCE SIZE (1..MAX) OF GeneralName
	//
	// GeneralName ::= CHOICE {
	//      otherName                       [0]     OtherName,
	//      rfc822Name                      [1]     IA5String,
	//      dNSName                         [2]     IA5String,
	//      x400Address                     [3]     ORAddress,
	//      directoryName                   [4]     Name,
	//      ediPartyName                    [5]     EDIPartyName,
	//      uniformResourceIdentifier       [6]     IA5String,
	//      iPAddress                       [7]     OCTET STRING,
	//      registeredID                    [8]     OBJECT IDENTIFIER }
	var seq asn1.RawValue
	rest, err := asn1.Unmarshal(extension, &seq)
	if err != nil {
		return err
	} else if len(rest) != 0 {
		return fmt.Errorf("x509: trailing data after X.509 extension")
	}
	if !seq.IsCompound || seq.Tag != 16 || seq.Class != 0 {
		return asn1.StructuralError{Msg: "bad SAN sequence"}
	}

	rest = seq.Bytes
	for len(rest) > 0 {
		var v asn1.RawValue
		rest, err = asn1.Unmarshal(rest, &v)
		if err != nil {
			return err
		}

		if err := callback(v.Tag, v.FullBytes); err != nil {
			return err
		}
	}

	return nil
}

// generateCreationBundle is a shared function that reads parameters supplied
// from the various endpoints and generates a CreationParameters with the
// parameters that can be used to issue or sign
func generateCreationBundle(b *backend, data *inputBundle, caSign *certutil.CAInfoBundle, csr *x509.CertificateRequest) (*certutil.CreationBundle, []string, error) {
	// Read in names -- CN, DNS and email addresses
	var cn string
	var ridSerialNumber string
	var warnings []string
	dnsNames := []string{}
	emailAddresses := []string{}
	{
		if csr != nil && data.role.UseCSRCommonName {
			cn = csr.Subject.CommonName
		}
		if cn == "" {
			cn = data.apiData.Get("common_name").(string)
			if cn == "" && data.role.RequireCN {
				return nil, nil, errutil.UserError{Err: `the common_name field is required, or must be provided in a CSR with "use_csr_common_name" set to true, unless "require_cn" is set to false`}
			}
		}

		ridSerialNumber = data.apiData.Get("serial_number").(string)

		// only take serial number from CSR if one was not supplied via API
		if ridSerialNumber == "" && csr != nil {
			ridSerialNumber = csr.Subject.SerialNumber
		}

		if csr != nil && data.role.UseCSRSANs {
			dnsNames = csr.DNSNames
			emailAddresses = csr.EmailAddresses
		}

		if cn != "" && !data.apiData.Get("exclude_cn_from_sans").(bool) {
			if strings.Contains(cn, "@") {
				// Note: emails are not disallowed if the role's email protection
				// flag is false, because they may well be included for
				// informational purposes; it is up to the verifying party to
				// ensure that email addresses in a subject alternate name can be
				// used for the purpose for which they are presented
				emailAddresses = append(emailAddresses, cn)
			} else {
				// Only add to dnsNames if it's actually a DNS name but convert
				// idn first
				p := idna.New(
					idna.StrictDomainName(true),
					idna.VerifyDNSLength(true),
				)
				converted, err := p.ToASCII(cn)
				if err != nil {
					return nil, nil, errutil.UserError{Err: err.Error()}
				}
				if hostnameRegex.MatchString(converted) {
					dnsNames = append(dnsNames, converted)
				}
			}
		}

		if csr == nil || !data.role.UseCSRSANs {
			cnAltRaw, ok := data.apiData.GetOk("alt_names")
			if ok {
				cnAlt := strutil.ParseDedupAndSortStrings(cnAltRaw.(string), ",")
				for _, v := range cnAlt {
					if strings.Contains(v, "@") {
						emailAddresses = append(emailAddresses, v)
					} else {
						// Only add to dnsNames if it's actually a DNS name but
						// convert idn first
						p := idna.New(
							idna.StrictDomainName(true),
							idna.VerifyDNSLength(true),
						)
						converted, err := p.ToASCII(v)
						if err != nil {
							return nil, nil, errutil.UserError{Err: err.Error()}
						}
						if hostnameRegex.MatchString(converted) {
							dnsNames = append(dnsNames, converted)
						}
					}
				}
			}
		}

		// Check the CN. This ensures that the CN is checked even if it's
		// excluded from SANs.
		if cn != "" {
			badName := validateCommonName(b, data, cn)
			if len(badName) != 0 {
				return nil, nil, errutil.UserError{Err: fmt.Sprintf(
					"common name %s not allowed by this role", badName)}
			}
		}

		if ridSerialNumber != "" {
			badName := validateSerialNumber(data, ridSerialNumber)
			if len(badName) != 0 {
				return nil, nil, errutil.UserError{Err: fmt.Sprintf(
					"serial_number %s not allowed by this role", badName)}
			}
		}

		// Check for bad email and/or DNS names
		badName := validateNames(b, data, dnsNames)
		if len(badName) != 0 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"subject alternate name %s not allowed by this role", badName)}
		}

		badName = validateNames(b, data, emailAddresses)
		if len(badName) != 0 {
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"email address %s not allowed by this role", badName)}
		}
	}

	// otherSANsInput has the same format as the other_sans HTTP param in the
	// Vault PKI API: it is a list of strings of the form <oid>;<type>:<value>
	// where <type> must be UTF8/UTF-8.
	var otherSANsInput []string
	// otherSANs is the output of parseOtherSANs(otherSANsInput): its keys are
	// the <oid> value, its values are of the form [<type>, <value>]
	var otherSANs map[string][]string
	if sans := data.apiData.Get("other_sans").([]string); len(sans) > 0 {
		otherSANsInput = sans
	}
	if data.role.UseCSRSANs && csr != nil && len(csr.Extensions) > 0 {
		others, err := getOtherSANsFromX509Extensions(csr.Extensions)
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Errorf("could not parse requested other SAN: %w", err).Error()}
		}
		for _, other := range others {
			otherSANsInput = append(otherSANsInput, other.String())
		}
	}
	if len(otherSANsInput) > 0 {
		requested, err := parseOtherSANs(otherSANsInput)
		if err != nil {
			return nil, nil, errutil.UserError{Err: fmt.Errorf("could not parse requested other SAN: %w", err).Error()}
		}
		badOID, badName, err := validateOtherSANs(data, requested)
		switch {
		case err != nil:
			return nil, nil, errutil.UserError{Err: err.Error()}
		case len(badName) > 0:
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"other SAN %s not allowed for OID %s by this role", badName, badOID)}
		case len(badOID) > 0:
			return nil, nil, errutil.UserError{Err: fmt.Sprintf(
				"other SAN OID %s not allowed by this role", badOID)}
		default:
			otherSANs = requested
		}
	}

	// Get and verify any IP SANs
	ipAddresses := []net.IP{}
	{
		if csr != nil && data.role.UseCSRSANs {
			if len(csr.IPAddresses) > 0 {
				if !data.role.AllowIPSANs {
					return nil, nil, errutil.UserError{Err: "IP Subject Alternative Names are not allowed in this role, but was provided some via CSR"}
				}
				ipAddresses = csr.IPAddresses
			}
		} else {
			ipAlt := data.apiData.Get("ip_sans").([]string)
			if len(ipAlt) > 0 {
				if !data.role.AllowIPSANs {
					return nil, nil, errutil.UserError{Err: fmt.Sprintf(
						"IP Subject Alternative Names are not allowed in this role, but was provided %s", ipAlt)}
				}
				for _, v := range ipAlt {
					parsedIP := net.ParseIP(v)
					if parsedIP == nil {
						return nil, nil, errutil.UserError{Err: fmt.Sprintf(
							"the value %q is not a valid IP address", v)}
					}
					ipAddresses = append(ipAddresses, parsedIP)
				}
			}
		}
	}

	URIs := []*url.URL{}
	{
		if csr != nil && data.role.UseCSRSANs {
			if len(csr.URIs) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return nil, nil, errutil.UserError{
						Err: "URI Subject Alternative Names are not allowed in this role, but were provided via CSR",
					}
				}

				// validate uri sans
				for _, uri := range csr.URIs {
					valid := validateURISAN(b, data, uri.String())
					if !valid {
						return nil, nil, errutil.UserError{
							Err: "URI Subject Alternative Names were provided via CSR which are not valid for this role",
						}
					}

					URIs = append(URIs, uri)
				}
			}
		} else {
			uriAlt := data.apiData.Get("uri_sans").([]string)
			if len(uriAlt) > 0 {
				if len(data.role.AllowedURISANs) == 0 {
					return nil, nil, errutil.UserError{
						Err: "URI Subject Alternative Names are not allowed in this role, but were provided via the API",
					}
				}

				for _, uri := range uriAlt {
					valid := validateURISAN(b, data, uri)
					if !valid {
						return nil, nil, errutil.UserError{
							Err: "URI Subject Alternative Names were provided via the API which are not valid for this role",
						}
					}

					parsedURI, err := url.Parse(uri)
					if parsedURI == nil || err != nil {
						return nil, nil, errutil.UserError{
							Err: fmt.Sprintf(
								"the provided URI Subject Alternative Name %q is not a valid URI", uri),
						}
					}

					URIs = append(URIs, parsedURI)
				}
			}
		}
	}

	// Most of these could also be RemoveDuplicateStable, or even
	// leave duplicates in, but OU is the one most likely to be duplicated.
	subject := pkix.Name{
		CommonName:         cn,
		SerialNumber:       ridSerialNumber,
		Country:            strutil.RemoveDuplicatesStable(data.role.Country, false),
		Organization:       strutil.RemoveDuplicatesStable(data.role.Organization, false),
		OrganizationalUnit: strutil.RemoveDuplicatesStable(data.role.OU, false),
		Locality:           strutil.RemoveDuplicatesStable(data.role.Locality, false),
		Province:           strutil.RemoveDuplicatesStable(data.role.Province, false),
		StreetAddress:      strutil.RemoveDuplicatesStable(data.role.StreetAddress, false),
		PostalCode:         strutil.RemoveDuplicatesStable(data.role.PostalCode, false),
	}

	// Get the TTL and verify it against the max allowed
	var ttl time.Duration
	var maxTTL time.Duration
	var notAfter time.Time
	var err error
	{
		ttl = time.Duration(data.apiData.Get("ttl").(int)) * time.Second
		notAfterAlt := data.role.NotAfter
		if notAfterAlt == "" {
			notAfterAltRaw, ok := data.apiData.GetOk("not_after")
			if ok {
				notAfterAlt = notAfterAltRaw.(string)
			}

		}
		if ttl > 0 && notAfterAlt != "" {
			return nil, nil, errutil.UserError{
				Err: "Either ttl or not_after should be provided. Both should not be provided in the same request.",
			}
		}

		if ttl == 0 && data.role.TTL > 0 {
			ttl = data.role.TTL
		}

		if data.role.MaxTTL > 0 {
			maxTTL = data.role.MaxTTL
		}

		if ttl == 0 {
			ttl = b.System().DefaultLeaseTTL()
		}
		if maxTTL == 0 {
			maxTTL = b.System().MaxLeaseTTL()
		}
		if ttl > maxTTL {
			warnings = append(warnings, fmt.Sprintf("TTL %q is longer than permitted maxTTL %q, so maxTTL is being used", ttl, maxTTL))
			ttl = maxTTL
		}

		if notAfterAlt != "" {
			notAfter, err = time.Parse(time.RFC3339, notAfterAlt)
			if err != nil {
				return nil, nil, errutil.UserError{Err: err.Error()}
			}
		} else {
			notAfter = time.Now().Add(ttl)
		}
		if caSign != nil && notAfter.After(caSign.Certificate.NotAfter) {
			// If it's not self-signed, verify that the issued certificate
			// won't be valid past the lifetime of the CA certificate, and
			// act accordingly. This is dependent based on the issuer's
			// LeafNotAfterBehavior argument.
			switch caSign.LeafNotAfterBehavior {
			case certutil.PermitNotAfterBehavior:
				// Explicitly do nothing.
			case certutil.TruncateNotAfterBehavior:
				notAfter = caSign.Certificate.NotAfter
			case certutil.ErrNotAfterBehavior:
				fallthrough
			default:
				return nil, nil, errutil.UserError{Err: fmt.Sprintf(
					"cannot satisfy request, as TTL would result in notAfter %s that is beyond the expiration of the CA certificate at %s", notAfter.Format(time.RFC3339Nano), caSign.Certificate.NotAfter.Format(time.RFC3339Nano))}
			}
		}
	}

	// Parse SKID from the request for cross-signing.
	var skid []byte
	{
		if rawSKIDValue, ok := data.apiData.GetOk("skid"); ok {
			// Handle removing common separators to make copy/paste from tool
			// output easier. Chromium uses space, OpenSSL uses colons, and at
			// one point, Vault had preferred dash as a separator for hex
			// strings.
			var err error
			skidValue := rawSKIDValue.(string)
			for _, separator := range []string{":", "-", " "} {
				skidValue = strings.ReplaceAll(skidValue, separator, "")
			}

			skid, err = hex.DecodeString(skidValue)
			if err != nil {
				return nil, nil, errutil.UserError{Err: fmt.Sprintf("cannot parse requested SKID value as hex: %v", err)}
			}
		}
	}

	creation := &certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			Subject:                       subject,
			DNSNames:                      strutil.RemoveDuplicates(dnsNames, false),
			EmailAddresses:                strutil.RemoveDuplicates(emailAddresses, false),
			IPAddresses:                   ipAddresses,
			URIs:                          URIs,
			OtherSANs:                     otherSANs,
			KeyType:                       data.role.KeyType,
			KeyBits:                       data.role.KeyBits,
			SignatureBits:                 data.role.SignatureBits,
			UsePSS:                        data.role.UsePSS,
			NotAfter:                      notAfter,
			KeyUsage:                      x509.KeyUsage(parseKeyUsages(data.role.KeyUsage)),
			ExtKeyUsage:                   parseExtKeyUsages(data.role),
			ExtKeyUsageOIDs:               data.role.ExtKeyUsageOIDs,
			PolicyIdentifiers:             data.role.PolicyIdentifiers,
			BasicConstraintsValidForNonCA: data.role.BasicConstraintsValidForNonCA,
			NotBeforeDuration:             data.role.NotBeforeDuration,
			ForceAppendCaChain:            caSign != nil,
			SKID:                          skid,
		},
		SigningBundle: caSign,
		CSR:           csr,
	}

	// Don't deal with URLs or max path length if it's self-signed, as these
	// normally come from the signing bundle
	if caSign == nil {
		return creation, warnings, nil
	}

	// This will have been read in from the getGlobalAIAURLs function
	creation.Params.URLs = caSign.URLs

	// If the max path length in the role is not nil, it was specified at
	// generation time with the max_path_length parameter; otherwise derive it
	// from the signing certificate
	if data.role.MaxPathLength != nil {
		creation.Params.MaxPathLength = *data.role.MaxPathLength
	} else {
		switch {
		case caSign.Certificate.MaxPathLen < 0:
			creation.Params.MaxPathLength = -1
		case caSign.Certificate.MaxPathLen == 0 &&
			caSign.Certificate.MaxPathLenZero:
			// The signing function will ensure that we do not issue a CA cert
			creation.Params.MaxPathLength = 0
		default:
			// If this takes it to zero, we handle this case later if
			// necessary
			creation.Params.MaxPathLength = caSign.Certificate.MaxPathLen - 1
		}
	}

	return creation, warnings, nil
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
		Id: asn1.ObjectIdentifier(oidExtensionSubjectAltName),
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
	block, extra := pem.Decode(certBytes)
	if block == nil {
		return nil, errors.New("unable to parse certificate: invalid PEM")
	}
	if len(strings.TrimSpace(string(extra))) > 0 {
		return nil, errors.New("unable to parse certificate: trailing PEM data")
	}

	return x509.ParseCertificate(block.Bytes)
}
