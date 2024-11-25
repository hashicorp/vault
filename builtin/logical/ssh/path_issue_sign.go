// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

var containsTemplateRegex = regexp.MustCompile(`{{.+?}}`)

var ecCurveBitsToAlgoName = map[int]string{
	256: ssh.KeyAlgoECDSA256,
	384: ssh.KeyAlgoECDSA384,
	521: ssh.KeyAlgoECDSA521,
}

// If the algorithm is not found, it could be that we have a curve
// that we haven't added a constant for yet. But they could allow it
// (assuming x/crypto/ssh can parse it) via setting a ec: <keyBits>
// mapping rather than using a named SSH key type, so erring out here
// isn't advisable.

type creationBundle struct {
	KeyID           string
	ValidPrincipals []string
	PublicKey       ssh.PublicKey
	CertificateType uint32
	TTL             time.Duration
	Signer          ssh.Signer
	Role            *sshRole
	CriticalOptions map[string]string
	Extensions      map[string]string
}

func (b *backend) pathSignIssueCertificateHelper(ctx context.Context, req *logical.Request, data *framework.FieldData, role *sshRole, publicKey ssh.PublicKey) (*logical.Response, error) {
	// Note that these various functions always return "user errors" so we pass
	// them as 4xx values
	keyID, err := b.calculateKeyID(data, req, role, publicKey)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	certificateType, err := b.calculateCertificateType(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	var parsedPrincipals []string
	if certificateType == ssh.HostCert {
		parsedPrincipals, err = b.calculateValidPrincipals(data, req, role, "", role.AllowedDomains, role.AllowedDomainsTemplate, validateValidPrincipalForHosts(role))
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	} else {
		defaultPrincipal := role.DefaultUser
		if role.DefaultUserTemplate {
			defaultPrincipal, err = b.renderPrincipal(role.DefaultUser, req)
			if err != nil {
				return nil, err
			}
		}
		parsedPrincipals, err = b.calculateValidPrincipals(data, req, role, defaultPrincipal, role.AllowedUsers, role.AllowedUsersTemplate, strutil.StrListContains)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	ttl, err := b.calculateTTL(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	criticalOptions, err := b.calculateCriticalOptions(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	extensions, addExtTemplatingWarning, err := b.calculateExtensions(data, req, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	privateKeyEntry, err := caKey(ctx, req.Storage, caPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA private key: %w", err)
	}
	if privateKeyEntry == nil || privateKeyEntry.Key == "" {
		return nil, errors.New("failed to read CA private key")
	}

	signer, err := ssh.ParsePrivateKey([]byte(privateKeyEntry.Key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse stored CA private key: %w", err)
	}

	cBundle := creationBundle{
		KeyID:           keyID,
		PublicKey:       publicKey,
		Signer:          signer,
		ValidPrincipals: parsedPrincipals,
		TTL:             ttl,
		CertificateType: certificateType,
		Role:            role,
		CriticalOptions: criticalOptions,
		Extensions:      extensions,
	}

	certificate, err := cBundle.sign()
	if err != nil {
		return nil, err
	}

	signedSSHCertificate := ssh.MarshalAuthorizedKey(certificate)
	if len(signedSSHCertificate) == 0 {
		return nil, errors.New("error marshaling signed certificate")
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			"serial_number": strconv.FormatUint(certificate.Serial, 16),
			"signed_key":    string(signedSSHCertificate),
		},
	}

	if addExtTemplatingWarning {
		response.AddWarning("default_extension templating enabled with at least one extension requiring identity templating. However, this request lacked identity entity information, causing one or more extensions to be skipped from the generated certificate.")
	}

	return response, nil
}

func (b *backend) renderPrincipal(principal string, req *logical.Request) (string, error) {
	// Look for templating markers {{ .* }}
	matched := containsTemplateRegex.MatchString(principal)
	if matched {
		if req.EntityID != "" {
			// Retrieve principal based on template + entityID from request.
			renderedPrincipal, err := framework.PopulateIdentityTemplate(principal, req.EntityID, b.System())
			if err != nil {
				return "", fmt.Errorf("template '%s' could not be rendered -> %s", principal, err)
			}
			return renderedPrincipal, nil
		}
	}
	// Static principal
	return principal, nil
}

func (b *backend) calculateValidPrincipals(data *framework.FieldData, req *logical.Request, role *sshRole, defaultPrincipal, principalsAllowedByRole string, enableTemplating bool, validatePrincipal func([]string, string) bool) ([]string, error) {
	validPrincipals := ""
	validPrincipalsRaw, ok := data.GetOk("valid_principals")
	if ok {
		validPrincipals = validPrincipalsRaw.(string)
	} else {
		validPrincipals = defaultPrincipal
	}

	parsedPrincipals := strutil.RemoveDuplicates(strutil.ParseStringSlice(validPrincipals, ","), false)
	// Build list of allowed Principals from template and static principalsAllowedByRole
	var allowedPrincipals []string
	if enableTemplating {
		rendered, err := b.renderPrincipal(principalsAllowedByRole, req)
		if err != nil {
			return nil, err
		}
		allowedPrincipals = strutil.RemoveDuplicates(strutil.ParseStringSlice(rendered, ","), false)
	} else {
		allowedPrincipals = strutil.RemoveDuplicates(strutil.ParseStringSlice(principalsAllowedByRole, ","), false)
	}

	if len(parsedPrincipals) == 0 && defaultPrincipal != "" {
		// defaultPrincipal will either be the defaultUser or a rendered defaultUserTemplate
		parsedPrincipals = []string{defaultPrincipal}
	}

	switch {
	case len(parsedPrincipals) == 0:
		if role.AllowEmptyPrincipals {
			// There is nothing to process
			return nil, nil
		} else {
			return nil, fmt.Errorf("empty valid principals not allowed by role")
		}
	case len(allowedPrincipals) == 0:
		// User has requested principals to be set, but role is not configured
		// with any principals
		return nil, fmt.Errorf("role is not configured to allow any principals")
	default:
		// Role was explicitly configured to allow any principal.
		if principalsAllowedByRole == "*" {
			return parsedPrincipals, nil
		}

		for _, principal := range parsedPrincipals {
			if !validatePrincipal(strutil.RemoveDuplicates(allowedPrincipals, false), principal) {
				return nil, fmt.Errorf("%v is not a valid value for valid_principals", principal)
			}
		}
		return parsedPrincipals, nil
	}
}

func validateValidPrincipalForHosts(role *sshRole) func([]string, string) bool {
	return func(allowedPrincipals []string, validPrincipal string) bool {
		for _, allowedPrincipal := range allowedPrincipals {
			if allowedPrincipal == validPrincipal && role.AllowBareDomains {
				return true
			}
			if role.AllowSubdomains && strings.HasSuffix(validPrincipal, "."+allowedPrincipal) {
				return true
			}
		}

		return false
	}
}

func (b *backend) calculateCertificateType(data *framework.FieldData, role *sshRole) (uint32, error) {
	requestedCertificateType := data.Get("cert_type").(string)

	var certificateType uint32
	switch requestedCertificateType {
	case "user":
		if !role.AllowUserCertificates {
			return 0, errors.New("cert_type 'user' is not allowed by role")
		}
		certificateType = ssh.UserCert
	case "host":
		if !role.AllowHostCertificates {
			return 0, errors.New("cert_type 'host' is not allowed by role")
		}
		certificateType = ssh.HostCert
	default:
		return 0, errors.New("cert_type must be either 'user' or 'host'")
	}

	return certificateType, nil
}

func (b *backend) calculateKeyID(data *framework.FieldData, req *logical.Request, role *sshRole, pubKey ssh.PublicKey) (string, error) {
	reqID := data.Get("key_id").(string)

	if reqID != "" {
		if !role.AllowUserKeyIDs {
			return "", fmt.Errorf("setting key_id is not allowed by role")
		}
		return reqID, nil
	}

	keyIDFormat := "vault-{{token_display_name}}-{{public_key_hash}}"
	if req.DisplayName == "" {
		keyIDFormat = "vault-{{public_key_hash}}"
	}

	if role.KeyIDFormat != "" {
		keyIDFormat = role.KeyIDFormat
	}

	keyID := substQuery(keyIDFormat, map[string]string{
		"token_display_name": req.DisplayName,
		"role_name":          data.Get("role").(string),
		"public_key_hash":    fmt.Sprintf("%x", sha256.Sum256(pubKey.Marshal())),
	})

	return keyID, nil
}

func (b *backend) calculateCriticalOptions(data *framework.FieldData, role *sshRole) (map[string]string, error) {
	unparsedCriticalOptions := data.Get("critical_options").(map[string]interface{})
	if len(unparsedCriticalOptions) == 0 {
		return role.DefaultCriticalOptions, nil
	}

	criticalOptions := convertMapToStringValue(unparsedCriticalOptions)

	if role.AllowedCriticalOptions != "" {
		notAllowedOptions := []string{}
		allowedCriticalOptions := strings.Split(role.AllowedCriticalOptions, ",")

		for option := range criticalOptions {
			if !strutil.StrListContains(allowedCriticalOptions, option) {
				notAllowedOptions = append(notAllowedOptions, option)
			}
		}

		if len(notAllowedOptions) != 0 {
			return nil, fmt.Errorf("critical options not on allowed list: %v", notAllowedOptions)
		}
	}

	return criticalOptions, nil
}

func (b *backend) calculateExtensions(data *framework.FieldData, req *logical.Request, role *sshRole) (map[string]string, bool, error) {
	unparsedExtensions := data.Get("extensions").(map[string]interface{})
	extensions := make(map[string]string)

	if len(unparsedExtensions) > 0 {
		extensions := convertMapToStringValue(unparsedExtensions)
		if role.AllowedExtensions == "*" {
			// Allowed extensions was configured to allow all
			return extensions, false, nil
		}

		notAllowed := []string{}
		allowedExtensions := strings.Split(role.AllowedExtensions, ",")
		for extensionKey := range extensions {
			if !strutil.StrListContains(allowedExtensions, extensionKey) {
				notAllowed = append(notAllowed, extensionKey)
			}
		}

		if len(notAllowed) != 0 {
			return nil, false, fmt.Errorf("extensions %v are not on allowed list", notAllowed)
		}
		return extensions, false, nil
	}

	haveMissingEntityInfoWithTemplatedExt := false

	if role.DefaultExtensionsTemplate {
		for extensionKey, extensionValue := range role.DefaultExtensions {
			// Look for templating markers {{ .* }}
			matched := containsTemplateRegex.MatchString(extensionValue)
			if matched {
				if req.EntityID != "" {
					// Retrieve extension value based on template + entityID from request.
					templateExtensionValue, err := framework.PopulateIdentityTemplate(extensionValue, req.EntityID, b.System())
					if err == nil {
						// Template returned an extension value that we can use
						extensions[extensionKey] = templateExtensionValue
					} else {
						return nil, false, fmt.Errorf("template '%s' could not be rendered -> %s", extensionValue, err)
					}
				} else {
					haveMissingEntityInfoWithTemplatedExt = true
				}
			} else {
				// Static extension value or err template
				extensions[extensionKey] = extensionValue
			}
		}
	} else {
		extensions = role.DefaultExtensions
	}

	return extensions, haveMissingEntityInfoWithTemplatedExt, nil
}

func (b *backend) calculateTTL(data *framework.FieldData, role *sshRole) (time.Duration, error) {
	var ttl, maxTTL time.Duration
	var err error

	ttlRaw, specifiedTTL := data.GetOk("ttl")
	if specifiedTTL {
		ttl = time.Duration(ttlRaw.(int)) * time.Second
	} else {
		ttl, err = parseutil.ParseDurationSecond(role.TTL)
		if err != nil {
			return 0, err
		}
	}
	if ttl == 0 {
		ttl = b.System().DefaultLeaseTTL()
	}

	maxTTL, err = parseutil.ParseDurationSecond(role.MaxTTL)
	if err != nil {
		return 0, err
	}
	if maxTTL == 0 {
		maxTTL = b.System().MaxLeaseTTL()
	}

	if ttl > maxTTL {
		// Don't error if they were using system defaults, only error if
		// they specifically chose a bad TTL
		if !specifiedTTL {
			ttl = maxTTL
		} else {
			return 0, fmt.Errorf("ttl is larger than maximum allowed %d", maxTTL/time.Second)
		}
	}

	return ttl, nil
}

func (b *backend) validateSignedKeyRequirements(publickey ssh.PublicKey, role *sshRole) error {
	if len(role.AllowedUserKeyTypesLengths) != 0 {
		var keyType string
		var keyBits int

		switch k := publickey.(type) {
		case ssh.CryptoPublicKey:
			ff := k.CryptoPublicKey()
			switch k := ff.(type) {
			case *rsa.PublicKey:
				keyType = "rsa"
				keyBits = k.N.BitLen()
			case *dsa.PublicKey:
				keyType = "dsa"
				keyBits = k.Parameters.P.BitLen()
			case *ecdsa.PublicKey:
				keyType = "ecdsa"
				keyBits = k.Curve.Params().BitSize
			case ed25519.PublicKey:
				keyType = "ed25519"
			default:
				return fmt.Errorf("public key type of %s is not allowed", keyType)
			}
		default:
			return fmt.Errorf("pubkey not suitable for crypto (expected ssh.CryptoPublicKey but found %T)", k)
		}

		keyTypeToMapKey := createKeyTypeToMapKey(keyType, keyBits)

		var present bool
		var pass bool
		for _, kstr := range keyTypeToMapKey[keyType] {
			allowed_values, ok := role.AllowedUserKeyTypesLengths[kstr]
			if !ok {
				continue
			}

			present = true

			for _, value := range allowed_values {
				if keyType == "rsa" || keyType == "dsa" {
					// Regardless of map naming, we always need to validate the
					// bit length of RSA and DSA keys. Use the keyType flag to
					if keyBits == value {
						pass = true
					}
				} else if kstr == "ec" || kstr == "ecdsa" {
					// If the map string is "ecdsa", we have to validate the keyBits
					// are a match for an allowed value, meaning that our curve
					// is allowed. This isn't necessary when a named curve (e.g.
					// ssh.KeyAlgoECDSA256) is allowed (and hence kstr is that),
					// because keyBits is already specified in the kstr. Thus,
					// we have conditioned around kstr and not keyType (like with
					// rsa or dsa).
					if keyBits == value {
						pass = true
					}
				} else {
					// We get here in two cases: we have a algo-named EC key
					// matching a format specifier in the key map (e.g., a P-256
					// key with a KeyAlgoECDSA256 entry in the map) or we have a
					// ed25519 key (which is always allowed).
					pass = true
				}
			}
		}

		if !present {
			return fmt.Errorf("key of type %s is not allowed", keyType)
		}

		if !pass {
			return fmt.Errorf("key is of an invalid size: %v", keyBits)
		}
	}
	return nil
}

func (b *creationBundle) sign() (retCert *ssh.Certificate, retErr error) {
	defer func() {
		if r := recover(); r != nil {
			errMsg, ok := r.(string)
			if ok {
				retCert = nil
				retErr = errors.New(errMsg)
			}
		}
	}()

	serialNumber, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	sshAlgorithmSigner, ok := b.Signer.(ssh.AlgorithmSigner)
	if !ok {
		return nil, fmt.Errorf("failed to generate signed SSH key: signer is not an AlgorithmSigner")
	}

	// prepare certificate for signing
	nonce := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate signed SSH key: error generating random nonce: %w", err)
	}
	certificate := &ssh.Certificate{
		Serial:          serialNumber.Uint64(),
		Key:             b.PublicKey,
		KeyId:           b.KeyID,
		ValidPrincipals: b.ValidPrincipals,
		ValidAfter:      uint64(now.Add(-b.Role.NotBeforeDuration).In(time.UTC).Unix()),
		ValidBefore:     uint64(now.Add(b.TTL).In(time.UTC).Unix()),
		CertType:        b.CertificateType,
		Permissions: ssh.Permissions{
			CriticalOptions: b.CriticalOptions,
			Extensions:      b.Extensions,
		},
		Nonce:        nonce,
		SignatureKey: sshAlgorithmSigner.PublicKey(),
	}

	// get bytes to sign; this is based on Certificate.bytesForSigning() from the go ssh lib
	out := certificate.Marshal()
	// Drop trailing signature length.
	certificateBytes := out[:len(out)-4]

	algo := b.Role.AlgorithmSigner

	// Handle the new default algorithm selection process correctly.
	if algo == DefaultAlgorithmSigner && sshAlgorithmSigner.PublicKey().Type() == ssh.KeyAlgoRSA {
		algo = ssh.SigAlgoRSASHA2256
	} else if algo == DefaultAlgorithmSigner {
		algo = ""
	}

	sig, err := sshAlgorithmSigner.SignWithAlgorithm(rand.Reader, certificateBytes, algo)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signed SSH key: sign error: %w", err)
	}

	certificate.Signature = sig

	return certificate, nil
}

func createKeyTypeToMapKey(keyType string, keyBits int) map[string][]string {
	keyTypeToMapKey := map[string][]string{
		"rsa":     {"rsa", ssh.KeyAlgoRSA},
		"dsa":     {"dsa", ssh.KeyAlgoDSA},
		"ecdsa":   {"ecdsa", "ec"},
		"ed25519": {"ed25519", ssh.KeyAlgoED25519},
	}

	if keyType == "ecdsa" {
		if algo, ok := ecCurveBitsToAlgoName[keyBits]; ok {
			keyTypeToMapKey[keyType] = append(keyTypeToMapKey[keyType], algo)
		}
	}

	return keyTypeToMapKey
}
