// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// The verification logic is adapted from
// https://github.com/google/certificate-transparency-go/blob/v1.3.1/x509/verify.go

package x509verify

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

// InvalidReason describes the reason a certificate is invalid.
type InvalidReason int

const (
	// NotAuthorizedToSign results when a certificate is signed by another
	// which isn't marked as a CA certificate.
	NotAuthorizedToSign InvalidReason = iota
	// Expired results when a certificate has expired, based on the time
	// given in the VerifyOptions.
	Expired
	// CANotAuthorizedForThisName results when an intermediate or root
	// certificate has a name constraint which doesn't permit a DNS or
	// other name (including IP address) in the leaf certificate.
	CANotAuthorizedForThisName
	// TooManyIntermediates results when a path length constraint is
	// violated.
	TooManyIntermediates
	// IncompatibleUsage results when the certificate's key usage indicates
	// that it may only be used for a different purpose.
	IncompatibleUsage
	// NameMismatch results when the subject name of a parent certificate
	// does not match the issuer name in the child.
	NameMismatch
	// NameConstraintsWithoutSANs results when a leaf certificate doesn't
	// contain a Subject Alternative Name extension, but a CA certificate
	// contains name constraints, and the Common Name can be interpreted as
	// a hostname.
	//
	// You can avoid this error by setting the experimental GODEBUG environment
	// variable to "x509ignoreCN=1", disabling Common Name matching entirely.
	// This behavior might become the default in the future.
	NameConstraintsWithoutSANs
	// UnconstrainedName results when a CA certificate contains permitted
	// name constraints, but leaf certificate contains a name of an
	// unsupported or unconstrained type.
	UnconstrainedName
	// TooManyConstraints results when the number of comparison operations
	// needed to check a certificate exceeds the limit set by
	// VerifyOptions.MaxConstraintComparisions. This limit exists to
	// prevent pathological certificates can consuming excessive amounts of
	// CPU time to verify.
	TooManyConstraints
	// CANotAuthorizedForExtKeyUsage results when an intermediate or root
	// certificate does not permit a requested extended key usage.
	CANotAuthorizedForExtKeyUsage
)

// CertificateInvalidError results when an odd error occurs. Users of this
// library probably want to handle all these errors uniformly.
type CertificateInvalidError struct {
	Cert   *x509.Certificate
	Reason InvalidReason
	Detail string
}

func (e CertificateInvalidError) Error() string {
	switch e.Reason {
	case NotAuthorizedToSign:
		return "x509: certificate is not authorized to sign other certificates"
	case Expired:
		return "x509: certificate has expired or is not yet valid: " + e.Detail
	case CANotAuthorizedForThisName:
		return "x509: a root or intermediate certificate is not authorized to sign for this name: " + e.Detail
	case CANotAuthorizedForExtKeyUsage:
		return "x509: a root or intermediate certificate is not authorized for an extended key usage: " + e.Detail
	case TooManyIntermediates:
		return "x509: too many intermediates for path length constraint"
	case IncompatibleUsage:
		return "x509: certificate specifies an incompatible key usage"
	case NameMismatch:
		return "x509: issuer name does not match subject from issuing certificate"
	case NameConstraintsWithoutSANs:
		return "x509: issuer has name constraints but leaf doesn't have a SAN extension"
	case UnconstrainedName:
		return "x509: issuer has name constraints but leaf contains unknown or unconstrained name: " + e.Detail
	case TooManyConstraints:
		return "x509: too many constraints were attempted before giving up"
	}
	return "x509: unknown error"
}

// UnhandledCriticalExtension is returned when a certificate contains a
// critical extension that was not handled.
type UnhandledCriticalExtension struct {
	ID asn1.ObjectIdentifier
}

func (h UnhandledCriticalExtension) Error() string {
	return fmt.Sprintf("x509: unhandled critical extension (%v)", h.ID)
}

// UnknownAuthorityError results when the certificate issuer is unknown
type UnknownAuthorityError struct {
	Cert *x509.Certificate
	// hintErr contains an error that may be helpful in determining why an
	// authority wasn't found.
	hintErr error
	// hintCert contains a possible authority certificate that was rejected
	// because of the error in hintErr.
	hintCert *x509.Certificate
}

func (e UnknownAuthorityError) Error() string {
	s := "x509: certificate signed by unknown authority"
	if e.hintErr != nil {
		certName := e.hintCert.Subject.CommonName
		if len(certName) == 0 {
			if len(e.hintCert.Subject.Organization) > 0 {
				certName = e.hintCert.Subject.Organization[0]
			} else {
				certName = "serial:" + e.hintCert.SerialNumber.String()
			}
		}
		s += fmt.Sprintf(" (possibly because of %q while trying to verify candidate authority certificate %q)", e.hintErr, certName)
	}
	return s
}

// SystemRootsError results when we fail to load the system root certificates.
type SystemRootsError struct {
	Err error
}

func (se SystemRootsError) Error() string {
	msg := "x509: failed to load system roots and no roots provided"
	if se.Err != nil {
		return msg + "; " + se.Err.Error()
	}
	return msg
}

// errNotParsed is returned when a certificate without ASN.1 contents is
// verified. Platform-specific verification needs the ASN.1 contents.
var errNotParsed = errors.New("x509: missing ASN.1 contents; use ParseCertificate")

// VerifyOptions contains parameters for certificate chain verification.
// It mirrors the extended options exposed by the Google Certificate Transparency
// x509 fork, rewritten to use only standard-library types.
type VerifyOptions struct {
	// DNSName, if non-empty, is the hostname to verify the leaf certificate against.
	DNSName string

	// Intermediates is an optional pool of certificates that are not trust
	// anchors but can be used to form the chain.
	Intermediates *CertPool

	// Roots is the pool of trusted root certificates. Callers must provide this;
	// x509verify does not currently support using the host's system root pool.
	Roots *CertPool

	// CurrentTime is used for expiry checks. If zero, time.Now() is used.
	CurrentTime time.Time

	// DisableTimeChecks, when true, skips NotBefore/NotAfter validation.
	DisableTimeChecks bool

	// DisableCriticalExtensionChecks, when true, skips rejection of
	// certificates containing unhandled critical extensions.
	DisableCriticalExtensionChecks bool

	// DisableNameChecks, when true, skips issuer/subject name mismatch
	// checks during chain building.
	DisableNameChecks bool

	// DisableEKUChecks, when true, skips Extended Key Usage enforcement
	// down the certificate chain.
	DisableEKUChecks bool

	// DisablePathLenChecks, when true, skips path length constraint
	// validation.
	DisablePathLenChecks bool

	// DisableNameConstraintChecks, when true, skips name constraint
	// enforcement on intermediate and root certificates.
	DisableNameConstraintChecks bool

	// KeyUsages specifies which Extended Key Usage values are acceptable for
	// the leaf certificate. An empty list defaults to ExtKeyUsageServerAuth.
	// To accept any key usage, include x509.ExtKeyUsageAny.
	KeyUsages []x509.ExtKeyUsage

	// MaxConstraintComparisions limits the number of comparisons performed
	// when checking name constraints. Zero means a sensible default (250000).
	MaxConstraintComparisions int
}

const (
	leafCertificate         = 0
	intermediateCertificate = 1
	rootCertificate         = 2
)

// Well-known OIDs used during verification.
var (
	oidExtensionSubjectAltName  = asn1.ObjectIdentifier{2, 5, 29, 17}
	oidExtensionNameConstraints = asn1.ObjectIdentifier{2, 5, 29, 30}
)

const (
	nameTypeEmail = 1
	nameTypeDNS   = 2
	nameTypeURI   = 6
	nameTypeIP    = 7
)

// oidInExtensions reports whether an extension with the given OID exists.
func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool {
	for _, e := range extensions {
		if e.Id.Equal(oid) {
			return true
		}
	}
	return false
}

func hasSANExtension(c *x509.Certificate) bool {
	return oidInExtensions(oidExtensionSubjectAltName, c.Extensions)
}

func hasNameConstraints(c *x509.Certificate) bool {
	return oidInExtensions(oidExtensionNameConstraints, c.Extensions)
}

func getSANExtension(c *x509.Certificate) []byte {
	for _, e := range c.Extensions {
		if e.Id.Equal(oidExtensionSubjectAltName) {
			return e.Value
		}
	}
	return nil
}

// rfc2821Mailbox represents an email address broken into local and domain parts.
type rfc2821Mailbox struct {
	local, domain string
}

// parseRFC2821Mailbox parses an email address into local and domain parts,
// based on the ABNF for a “Mailbox” from RFC 2821. According to RFC 5280,
// Section 4.2.1.6 that's correct for an rfc822Name from a certificate: “The
// format of an rfc822Name is a "Mailbox" as defined in RFC 2821, Section 4.1.2”.
func parseRFC2821Mailbox(in string) (mailbox rfc2821Mailbox, ok bool) {
	if len(in) == 0 {
		return mailbox, false
	}

	localPartBytes := make([]byte, 0, len(in)/2)

	if in[0] == '"' {
		// Quoted-string = DQUOTE *qcontent DQUOTE
		// non-whitespace-control = %d1-8 / %d11 / %d12 / %d14-31 / %d127
		// qcontent = qtext / quoted-pair
		// qtext = non-whitespace-control /
		//         %d33 / %d35-91 / %d93-126
		// quoted-pair = ("\" text) / obs-qp
		// text = %d1-9 / %d11 / %d12 / %d14-127 / obs-text
		//
		// (Names beginning with “obs-” are the obsolete syntax from RFC 2822,
		// Section 4. Since it has been 16 years, we no longer accept that.)
		in = in[1:]
	QuotedString:
		for {
			if len(in) == 0 {
				return mailbox, false
			}
			c := in[0]
			in = in[1:]

			switch {
			case c == '"':
				break QuotedString

			case c == '\\':
				// quoted-pair
				if len(in) == 0 {
					return mailbox, false
				}
				if in[0] == 11 || in[0] == 12 ||
					(1 <= in[0] && in[0] <= 9) ||
					(14 <= in[0] && in[0] <= 127) {
					localPartBytes = append(localPartBytes, in[0])
					in = in[1:]
				} else {
					return mailbox, false
				}

			case c == 11 ||
				c == 12 ||
				// Space (char 32) is not allowed based on the
				// BNF, but RFC 3696 gives an example that
				// assumes that it is. Several “verified”
				// errata continue to argue about this point.
				// We choose to accept it.
				c == 32 ||
				c == 33 ||
				c == 127 ||
				(1 <= c && c <= 8) ||
				(14 <= c && c <= 31) ||
				(35 <= c && c <= 91) ||
				(93 <= c && c <= 126):
				// qtext
				localPartBytes = append(localPartBytes, c)

			default:
				return mailbox, false
			}
		}
	} else {
		// Atom ("." Atom)*
	NextChar:
		for len(in) > 0 {
			// atext from RFC 2822, Section 3.2.4
			c := in[0]

			switch {
			case c == '\\':
				// Examples given in RFC 3696 suggest that
				// escaped characters can appear outside of a
				// quoted string. Several “verified” errata
				// continue to argue the point. We choose to
				// accept it.
				in = in[1:]
				if len(in) == 0 {
					return mailbox, false
				}
				fallthrough
			case ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') ||
				('A' <= c && c <= 'Z') ||
				c == '!' || c == '#' || c == '$' || c == '%' ||
				c == '&' || c == '\'' || c == '*' || c == '+' ||
				c == '-' || c == '/' || c == '=' || c == '?' ||
				c == '^' || c == '_' || c == '`' || c == '{' ||
				c == '|' || c == '}' || c == '~' || c == '.':
				localPartBytes = append(localPartBytes, in[0])
				in = in[1:]
			default:
				break NextChar
			}
		}
		if len(localPartBytes) == 0 {
			return mailbox, false
		}

		// From RFC 3696, Section 3:
		// “period (".") may also appear, but may not be used to start
		// or end the local part, nor may two or more consecutive
		// periods appear.”
		twoDots := []byte{'.', '.'}
		if localPartBytes[0] == '.' ||
			localPartBytes[len(localPartBytes)-1] == '.' ||
			bytes.Contains(localPartBytes, twoDots) {
			return mailbox, false
		}
	}

	if len(in) == 0 || in[0] != '@' {
		return mailbox, false
	}
	in = in[1:]

	// The RFC species a format for domains, but that's known to be
	// violated in practice so we accept that anything after an '@' is the
	// domain part.
	if _, ok := domainToReverseLabels(in); !ok {
		return mailbox, false
	}

	mailbox.local = string(localPartBytes)
	mailbox.domain = in
	return mailbox, true
}

// domainToReverseLabels converts a textual domain name like foo.example.com to
// the list of labels in reverse order, e.g. ["com", "example", "foo"].
func domainToReverseLabels(domain string) (reverseLabels []string, ok bool) {
	for len(domain) > 0 {
		if i := strings.LastIndexByte(domain, '.'); i == -1 {
			reverseLabels = append(reverseLabels, domain)
			domain = ""
		} else {
			reverseLabels = append(reverseLabels, domain[i+1:])
			domain = domain[:i]
		}
	}

	if len(reverseLabels) > 0 && len(reverseLabels[0]) == 0 {
		// An empty label at the end indicates an absolute value.
		return nil, false
	}

	for _, label := range reverseLabels {
		if len(label) == 0 {
			// Empty labels are otherwise invalid.
			return nil, false
		}

		for _, c := range label {
			if c < 33 || c > 126 {
				// Invalid character.
				return nil, false
			}
		}
	}

	return reverseLabels, true
}

func matchEmailConstraint(mailbox rfc2821Mailbox, constraint string) (bool, error) {
	// If the constraint contains an @, then it specifies an exact mailbox
	// name.
	if strings.Contains(constraint, "@") {
		constraintMailbox, ok := parseRFC2821Mailbox(constraint)
		if !ok {
			return false, fmt.Errorf("x509: internal error: cannot parse constraint %q", constraint)
		}
		return mailbox.local == constraintMailbox.local && strings.EqualFold(mailbox.domain, constraintMailbox.domain), nil
	}

	// Otherwise the constraint is like a DNS constraint of the domain part
	// of the mailbox.
	return matchDomainConstraint(mailbox.domain, constraint)
}

func matchURIConstraint(uri *url.URL, constraint string) (bool, error) {
	// From RFC 5280, Section 4.2.1.10:
	// “a uniformResourceIdentifier that does not include an authority
	// component with a host name specified as a fully qualified domain
	// name (e.g., if the URI either does not include an authority
	// component or includes an authority component in which the host name
	// is specified as an IP address), then the application MUST reject the
	// certificate.”

	host := uri.Host
	if len(host) == 0 {
		return false, fmt.Errorf("URI with empty host (%q) cannot be matched against constraints", uri.String())
	}

	if strings.Contains(host, ":") && !strings.HasSuffix(host, "]") {
		var err error
		host, _, err = net.SplitHostPort(uri.Host)
		if err != nil {
			return false, err
		}
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") || net.ParseIP(host) != nil {
		return false, fmt.Errorf("URI with IP (%q) cannot be matched against constraints", uri.String())
	}

	return matchDomainConstraint(host, constraint)
}

func matchIPConstraint(ip net.IP, constraint *net.IPNet) (bool, error) {
	if len(ip) != len(constraint.IP) {
		return false, nil
	}

	for i := range ip {
		if mask := constraint.Mask[i]; ip[i]&mask != constraint.IP[i]&mask {
			return false, nil
		}
	}

	return true, nil
}

func matchDomainConstraint(domain, constraint string) (bool, error) {
	// The meaning of zero length constraints is not specified, but this
	// code follows NSS and accepts them as matching everything.
	if len(constraint) == 0 {
		return true, nil
	}

	domainLabels, ok := domainToReverseLabels(domain)
	if !ok {
		return false, fmt.Errorf("x509: internal error: cannot parse domain %q", domain)
	}

	// RFC 5280 says that a leading period in a domain name means that at
	// least one label must be prepended, but only for URI and email
	// constraints, not DNS constraints. The code also supports that
	// behaviour for DNS constraints.

	mustHaveSubdomains := false
	if constraint[0] == '.' {
		mustHaveSubdomains = true
		constraint = constraint[1:]
	}

	constraintLabels, ok := domainToReverseLabels(constraint)
	if !ok {
		return false, fmt.Errorf("x509: internal error: cannot parse domain %q", constraint)
	}

	if len(domainLabels) < len(constraintLabels) ||
		(mustHaveSubdomains && len(domainLabels) == len(constraintLabels)) {
		return false, nil
	}

	for i, constraintLabel := range constraintLabels {
		if !strings.EqualFold(constraintLabel, domainLabels[i]) {
			return false, nil
		}
	}

	return true, nil
}

// toLowerCaseASCII returns a lower-case ASCII version of in.
func toLowerCaseASCII(in string) string {
	isAlreadyLower := true
	for _, c := range in {
		if c == utf8.RuneError {
			isAlreadyLower = false
			break
		}
		if 'A' <= c && c <= 'Z' {
			isAlreadyLower = false
			break
		}
	}
	if isAlreadyLower {
		return in
	}
	out := []byte(in)
	for i, c := range out {
		if 'A' <= c && c <= 'Z' {
			out[i] += 'a' - 'A'
		}
	}
	return string(out)
}

// checkNameConstraints checks that c permits a child certificate to claim the
// given name, of type nameType. The argument parsedName contains the parsed
// form of name, suitable for passing to the match function. The total number
// of comparisons is tracked in the given count and should not exceed the given
// limit.
func checkNameConstraints(c *x509.Certificate, count *int, maxConstraintComparisons int, nameType, name string, parsedName interface{},
	match func(parsedName, constraint interface{}) (bool, error),
	permitted, excluded interface{},
) error {
	excludedValue := reflect.ValueOf(excluded)
	*count += excludedValue.Len()
	if *count > maxConstraintComparisons {
		return CertificateInvalidError{c, TooManyConstraints, ""}
	}
	for i := 0; i < excludedValue.Len(); i++ {
		constraint := excludedValue.Index(i).Interface()
		match, err := match(parsedName, constraint)
		if err != nil {
			return CertificateInvalidError{c, CANotAuthorizedForThisName, err.Error()}
		}
		if match {
			return CertificateInvalidError{c, CANotAuthorizedForThisName, fmt.Sprintf("%s %q is excluded by constraint %q", nameType, name, constraint)}
		}
	}

	permittedValue := reflect.ValueOf(permitted)
	*count += permittedValue.Len()
	if *count > maxConstraintComparisons {
		return CertificateInvalidError{c, TooManyConstraints, ""}
	}
	ok := true
	for i := 0; i < permittedValue.Len(); i++ {
		constraint := permittedValue.Index(i).Interface()
		var err error
		if ok, err = match(parsedName, constraint); err != nil {
			return CertificateInvalidError{c, CANotAuthorizedForThisName, err.Error()}
		}
		if ok {
			break
		}
	}
	if !ok {
		return CertificateInvalidError{c, CANotAuthorizedForThisName, fmt.Sprintf("%s %q is not permitted by any constraint", nameType, name)}
	}
	return nil
}

// isValid performs validity checks on c given that it is a candidate to append
// to the chain in currentChain.
func isValid(c *x509.Certificate, certType int, currentChain []*x509.Certificate, opts *VerifyOptions) error {
	if !opts.DisableCriticalExtensionChecks && len(c.UnhandledCriticalExtensions) > 0 {
		return UnhandledCriticalExtension{ID: c.UnhandledCriticalExtensions[0]}
	}

	if !opts.DisableNameChecks && len(currentChain) > 0 {
		child := currentChain[len(currentChain)-1]
		if !bytes.Equal(child.RawIssuer, c.RawSubject) {
			return CertificateInvalidError{c, NameMismatch, ""}
		}
	}

	if !opts.DisableTimeChecks {
		now := opts.CurrentTime
		if now.IsZero() {
			now = time.Now()
		}
		if now.Before(c.NotBefore) {
			return CertificateInvalidError{
				Cert:   c,
				Reason: Expired,
				Detail: fmt.Sprintf("current time %s is before %s", now.Format(time.RFC3339), c.NotBefore.Format(time.RFC3339)),
			}
		} else if now.After(c.NotAfter) {
			return CertificateInvalidError{
				Cert:   c,
				Reason: Expired,
				Detail: fmt.Sprintf("current time %s is after %s", now.Format(time.RFC3339), c.NotAfter.Format(time.RFC3339)),
			}
		}
	}

	maxConstraintComparisons := opts.MaxConstraintComparisions
	if maxConstraintComparisons == 0 {
		maxConstraintComparisons = 250000
	}
	comparisonCount := 0

	var leaf *x509.Certificate
	if certType == intermediateCertificate || certType == rootCertificate {
		if len(currentChain) == 0 {
			return errors.New("x509: internal error: empty chain when appending CA cert")
		}
		leaf = currentChain[0]
	}

	doNameConstraintCheck := !opts.DisableNameConstraintChecks &&
		(certType == intermediateCertificate || certType == rootCertificate) &&
		hasNameConstraints(c)
	if doNameConstraintCheck && commonNameAsHostname(leaf) {
		// Legacy CN-as-hostname: we cannot safely check name constraints
		// against it, so reject to avoid bypassing the constraint.
		return CertificateInvalidError{c, NameConstraintsWithoutSANs, ""}
	} else if doNameConstraintCheck && hasSANExtension(leaf) {
		err := forEachSAN(getSANExtension(leaf), func(tag int, data []byte) error {
			switch tag {
			case nameTypeEmail:
				name := string(data)
				mailbox, ok := parseRFC2821Mailbox(name)
				if !ok {
					return fmt.Errorf("x509: cannot parse rfc822Name %q", name)
				}
				return checkNameConstraints(c, &comparisonCount, maxConstraintComparisons,
					"email address", name, mailbox,
					func(parsedName, constraint interface{}) (bool, error) {
						return matchEmailConstraint(parsedName.(rfc2821Mailbox), constraint.(string))
					},
					c.PermittedEmailAddresses, c.ExcludedEmailAddresses)

			case nameTypeDNS:
				name := string(data)
				if _, ok := domainToReverseLabels(name); !ok {
					return fmt.Errorf("x509: cannot parse dnsName %q", name)
				}
				return checkNameConstraints(c, &comparisonCount, maxConstraintComparisons,
					"DNS name", name, name,
					func(parsedName, constraint interface{}) (bool, error) {
						return matchDomainConstraint(parsedName.(string), constraint.(string))
					},
					c.PermittedDNSDomains, c.ExcludedDNSDomains)

			case nameTypeURI:
				name := string(data)
				uri, err := url.Parse(name)
				if err != nil {
					return fmt.Errorf("x509: internal error: URI SAN %q failed to parse", name)
				}
				return checkNameConstraints(c, &comparisonCount, maxConstraintComparisons,
					"URI", name, uri,
					func(parsedName, constraint interface{}) (bool, error) {
						return matchURIConstraint(parsedName.(*url.URL), constraint.(string))
					},
					c.PermittedURIDomains, c.ExcludedURIDomains)

			case nameTypeIP:
				ip := net.IP(data)
				if l := len(ip); l != net.IPv4len && l != net.IPv6len {
					return fmt.Errorf("x509: internal error: IP SAN %x failed to parse", data)
				}
				return checkNameConstraints(c, &comparisonCount, maxConstraintComparisons,
					"IP address", ip.String(), ip,
					func(parsedName, constraint interface{}) (bool, error) {
						return matchIPConstraint(parsedName.(net.IP), constraint.(*net.IPNet))
					},
					c.PermittedIPRanges, c.ExcludedIPRanges)
			default:
				// Unknown SAN types are ignored.
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	// KeyUsage status flags are ignored. From Engineering Security, Peter
	// Gutmann: A European government CA marked its signing certificates as
	// being valid for encryption only, but no-one noticed. Another
	// European CA marked its signature keys as not being valid for
	// signatures. A different CA marked its own trusted root certificate
	// as being invalid for certificate signing. Another national CA
	// distributed a certificate to be used to encrypt data for the
	// country’s tax authority that was marked as only being usable for
	// digital signatures but not for encryption. Yet another CA reversed
	// the order of the bit flags in the keyUsage due to confusion over
	// encoding endianness, essentially setting a random keyUsage in
	// certificates that it issued. Another CA created a self-invalidating
	// certificate by adding a certificate policy statement stipulating
	// that the certificate had to be used strictly as specified in the
	// keyUsage, and a keyUsage containing a flag indicating that the RSA
	// encryption key could only be used for Diffie-Hellman key agreement.

	if certType == intermediateCertificate && (!c.BasicConstraintsValid || !c.IsCA) {
		return CertificateInvalidError{c, NotAuthorizedToSign, ""}
	}

	if !opts.DisablePathLenChecks && c.BasicConstraintsValid && c.MaxPathLen >= 0 {
		numIntermediates := len(currentChain) - 1
		if numIntermediates > c.MaxPathLen {
			return CertificateInvalidError{c, TooManyIntermediates, ""}
		}
	}

	return nil
}

// Verify attempts to verify c by building one or more chains from c to a
// certificate in opts.Roots, using certificates in opts.Intermediates if
// needed. If successful, it returns one or more chains where the first element
// is c and the last element is from opts.Roots.
//
// Unlike crypto/x509.Certificate.Verify, this function supports additional
// disable flags in VerifyOptions to bypass checks that are not appropriate for
// all Vault PKI use cases (e.g. DisableTimeChecks, DisableEKUChecks).
func Verify(c *x509.Certificate, opts VerifyOptions) (chains [][]*x509.Certificate, err error) {
	if len(c.Raw) == 0 {
		return nil, errors.New("x509: missing ASN.1 contents; use ParseCertificate")
	}

	if opts.Intermediates != nil {
		for _, intermediate := range opts.Intermediates.certs {
			if len(intermediate.Raw) == 0 {
				return nil, errNotParsed
			}
		}
	}

	if opts.Roots == nil {
		// For simplicity and security, require callers to supply an explicit
		// root pool. Using the system pool through the stdlib would bypass our
		// extended verification options.
		return nil, errors.New("x509: Roots must be provided; system root pool is not supported by x509verify")
	}

	if opts.Intermediates == nil {
		opts.Intermediates = NewCertPool()
	}

	if err = isValid(c, leafCertificate, nil, &opts); err != nil {
		return
	}

	if len(opts.DNSName) > 0 {
		if err = c.VerifyHostname(opts.DNSName); err != nil {
			return
		}
	}

	var candidateChains [][]*x509.Certificate
	if opts.Roots.contains(c) {
		candidateChains = append(candidateChains, []*x509.Certificate{c})
	} else {
		if candidateChains, err = buildChains(c, nil, []*x509.Certificate{c}, nil, &opts); err != nil {
			return nil, err
		}
	}

	keyUsages := opts.KeyUsages
	if len(keyUsages) == 0 {
		keyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	}

	// If any key usage is acceptable then we're done.
	for _, usage := range keyUsages {
		if usage == x509.ExtKeyUsageAny {
			return candidateChains, nil
		}
	}

	for _, candidate := range candidateChains {
		if opts.DisableEKUChecks || checkChainForKeyUsage(candidate, keyUsages) {
			chains = append(chains, candidate)
		}
	}

	if len(chains) == 0 {
		return nil, CertificateInvalidError{c, IncompatibleUsage, ""}
	}

	return chains, nil
}

func appendToFreshChain(chain []*x509.Certificate, cert *x509.Certificate) []*x509.Certificate {
	n := make([]*x509.Certificate, len(chain)+1)
	copy(n, chain)
	n[len(chain)] = cert
	return n
}

// maxChainSignatureChecks is the maximum number of CheckSignatureFrom calls
// that an invocation of buildChains will (tranistively) make. Most chains are
// less than 15 certificates long, so this leaves space for multiple chains and
// for failed checks due to different intermediates having the same Subject.
const maxChainSignatureChecks = 100

func buildChains(c *x509.Certificate, cache map[*x509.Certificate][][]*x509.Certificate, currentChain []*x509.Certificate, sigChecks *int, opts *VerifyOptions) (chains [][]*x509.Certificate, err error) {
	var (
		hintErr  error
		hintCert *x509.Certificate
	)

	considerCandidate := func(certType int, candidate *x509.Certificate) {
		for _, cert := range currentChain {
			if cert.Equal(candidate) {
				return
			}
		}
		if sigChecks == nil {
			sigChecks = new(int)
		}
		*sigChecks++
		if *sigChecks > maxChainSignatureChecks {
			err = errors.New("x509: signature check attempts limit reached while verifying certificate chain")
			return
		}
		if err := c.CheckSignatureFrom(candidate); err != nil {
			if hintErr == nil {
				hintErr = err
				hintCert = candidate
			}
			return
		}
		if err = isValid(candidate, certType, currentChain, opts); err != nil {
			return
		}
		switch certType {
		case rootCertificate:
			chains = append(chains, appendToFreshChain(currentChain, candidate))
		case intermediateCertificate:
			if cache == nil {
				cache = make(map[*x509.Certificate][][]*x509.Certificate)
			}
			childChains, ok := cache[candidate]
			if !ok {
				childChains, err = buildChains(candidate, cache, appendToFreshChain(currentChain, candidate), sigChecks, opts)
				cache[candidate] = childChains
			}
			chains = append(chains, childChains...)
		}
	}

	for _, rootNum := range opts.Roots.findPotentialParents(c) {
		considerCandidate(rootCertificate, opts.Roots.certs[rootNum])
	}
	for _, intermediateNum := range opts.Intermediates.findPotentialParents(c) {
		considerCandidate(intermediateCertificate, opts.Intermediates.certs[intermediateNum])
	}

	if len(chains) > 0 {
		err = nil
	}
	if len(chains) == 0 && err == nil {
		err = UnknownAuthorityError{c, hintErr, hintCert}
	}

	return
}

// validHostname reports whether host is a valid hostname that can be matched or
// matched against according to RFC 6125 2.2, with some leniency to accommodate
// legacy values.
func validHostname(host string) bool {
	host = strings.TrimSuffix(host, ".")
	if len(host) == 0 {
		return false
	}
	for i, part := range strings.Split(host, ".") {
		if part == "" {
			// Empty label.
			return false
		}
		if i == 0 && part == "*" {
			// Only allow full left-most wildcards, as those are the only ones
			// we match, and matching literal '*' characters is probably never
			// the expected behavior.
			continue
		}
		for j, ch := range part {
			if 'a' <= ch && ch <= 'z' {
				continue
			}
			if '0' <= ch && ch <= '9' {
				continue
			}
			if 'A' <= ch && ch <= 'Z' {
				continue
			}
			if ch == '-' && j != 0 {
				continue
			}
			if ch == '_' || ch == ':' {
				// Not valid characters in hostnames, but commonly
				// found in deployments outside the WebPKI.
				continue
			}
			return false
		}
	}
	return true
}

// commonNameAsHostname reports whether the Common Name field should be
// considered the hostname that the certificate is valid for. This is a legacy
// behavior, disabled if the Subject Alt Name extension is present.
//
// It applies the strict validHostname check to the Common Name field, so that
// certificates without SANs can still be validated against CAs with name
// constraints if there is no risk the CN would be matched as a hostname.
// See NameConstraintsWithoutSANs and issue 24151.
func commonNameAsHostname(c *x509.Certificate) bool {
	return !hasSANExtension(c) && validHostname(c.Subject.CommonName)
}

func checkChainForKeyUsage(chain []*x509.Certificate, keyUsages []x509.ExtKeyUsage) bool {
	usages := make([]x509.ExtKeyUsage, len(keyUsages))
	copy(usages, keyUsages)

	if len(chain) == 0 {
		return false
	}
	usagesRemaining := len(usages)

	// We walk down the list and cross out any usages that aren't supported
	// by each certificate. If we cross out all the usages, then the chain
	// is unacceptable.
NextCert:
	for i := len(chain) - 1; i >= 0; i-- {
		cert := chain[i]
		if len(cert.ExtKeyUsage) == 0 && len(cert.UnknownExtKeyUsage) == 0 {
			// The certificate doesn't have any extended key usage specified.
			continue
		}
		for _, usage := range cert.ExtKeyUsage {
			if usage == x509.ExtKeyUsageAny {
				// The certificate is explicitly good for any usage.
				continue NextCert
			}
		}
		const invalidUsage x509.ExtKeyUsage = -1

	NextRequestedUsage:
		for i, requestedUsage := range usages {
			if requestedUsage == invalidUsage {
				continue
			}
			for _, usage := range cert.ExtKeyUsage {
				if requestedUsage == usage {
					continue NextRequestedUsage
				} else if requestedUsage == x509.ExtKeyUsageServerAuth &&
					(usage == x509.ExtKeyUsageNetscapeServerGatedCrypto ||
						usage == x509.ExtKeyUsageMicrosoftServerGatedCrypto) {
					// In order to support COMODO
					// certificate chains, we have to
					// accept Netscape or Microsoft SGC
					// usages as equal to ServerAuth.
					continue NextRequestedUsage
				}
			}
			usages[i] = invalidUsage
			usagesRemaining--
			if usagesRemaining == 0 {
				return false
			}
		}
	}
	return true
}

// forEachSAN calls callback for each SAN value in the DER-encoded extension.
func forEachSAN(extension []byte, callback func(tag int, data []byte) error) error {
	var seq asn1.RawValue
	rest, err := asn1.Unmarshal(extension, &seq)
	if err != nil {
		return err
	}
	if len(rest) != 0 {
		return errors.New("x509: trailing data after X.509 extension")
	}
	if !seq.IsCompound || seq.Tag != asn1.TagSequence || seq.Class != asn1.ClassUniversal {
		return asn1.StructuralError{Msg: "bad SAN sequence"}
	}
	rest = seq.Bytes
	for len(rest) > 0 {
		var v asn1.RawValue
		rest, err = asn1.Unmarshal(rest, &v)
		if err != nil {
			return err
		}
		if err := callback(v.Tag, v.Bytes); err != nil {
			return err
		}
	}
	return nil
}
