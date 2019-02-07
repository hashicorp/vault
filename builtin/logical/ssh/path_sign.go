package ssh

import (
	"context"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

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

func pathSign(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameRegex("role"),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSign,
		},

		Fields: map[string]*framework.FieldSchema{
			"role": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `The desired role with configuration for this request.`,
			},
			"ttl": &framework.FieldSchema{
				Type: framework.TypeDurationSecond,
				Description: `The requested Time To Live for the SSH certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be later than the role max TTL.`,
			},
			"public_key": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `SSH public key that should be signed.`,
			},
			"valid_principals": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Valid principals, either usernames or hostnames, that the certificate should be signed for.`,
			},
			"cert_type": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Type of certificate to be created; either "user" or "host".`,
				Default:     "user",
			},
			"key_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: `Key id that the created certificate should have. If not specified, the display name of the token will be used.`,
			},
			"critical_options": &framework.FieldSchema{
				Type:        framework.TypeMap,
				Description: `Critical options that the certificate should be signed for.`,
			},
			"extensions": &framework.FieldSchema{
				Type:        framework.TypeMap,
				Description: `Extensions that the certificate should be signed for.`,
			},
		},

		HelpSynopsis:    `Request signing an SSH key using a certain role with the provided details.`,
		HelpDescription: `This path allows SSH keys to be signed according to the policy of the given role.`,
	}
}

func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	return b.pathSignCertificate(ctx, req, data, role)
}

func (b *backend) pathSignCertificate(ctx context.Context, req *logical.Request, data *framework.FieldData, role *sshRole) (*logical.Response, error) {
	publicKey := data.Get("public_key").(string)
	if publicKey == "" {
		return logical.ErrorResponse("missing public_key"), nil
	}

	userPublicKey, err := parsePublicSSHKey(publicKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse public_key as SSH key: %s", err)), nil
	}

	err = b.validateSignedKeyRequirements(userPublicKey, role)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("public_key failed to meet the key requirements: %s", err)), nil
	}

	// Note that these various functions always return "user errors" so we pass
	// them as 4xx values
	keyID, err := b.calculateKeyID(data, req, role, userPublicKey)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	certificateType, err := b.calculateCertificateType(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	var parsedPrincipals []string
	if certificateType == ssh.HostCert {
		parsedPrincipals, err = b.calculateValidPrincipals(data, "", role.AllowedDomains, validateValidPrincipalForHosts(role))
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	} else {
		parsedPrincipals, err = b.calculateValidPrincipals(data, role.DefaultUser, role.AllowedUsers, strutil.StrListContains)
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

	extensions, err := b.calculateExtensions(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	privateKeyEntry, err := caKey(ctx, req.Storage, caPrivateKey)
	if err != nil {
		return nil, errwrap.Wrapf("failed to read CA private key: {{err}}", err)
	}
	if privateKeyEntry == nil || privateKeyEntry.Key == "" {
		return nil, fmt.Errorf("failed to read CA private key")
	}

	signer, err := ssh.ParsePrivateKey([]byte(privateKeyEntry.Key))
	if err != nil {
		return nil, errwrap.Wrapf("failed to parse stored CA private key: {{err}}", err)
	}

	cBundle := creationBundle{
		KeyID:           keyID,
		PublicKey:       userPublicKey,
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
		return nil, fmt.Errorf("error marshaling signed certificate")
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			"serial_number": strconv.FormatUint(certificate.Serial, 16),
			"signed_key":    string(signedSSHCertificate),
		},
	}

	return response, nil
}

func (b *backend) calculateValidPrincipals(data *framework.FieldData, defaultPrincipal, principalsAllowedByRole string, validatePrincipal func([]string, string) bool) ([]string, error) {
	validPrincipals := ""
	validPrincipalsRaw, ok := data.GetOk("valid_principals")
	if ok {
		validPrincipals = validPrincipalsRaw.(string)
	} else {
		validPrincipals = defaultPrincipal
	}

	parsedPrincipals := strutil.RemoveDuplicates(strutil.ParseStringSlice(validPrincipals, ","), false)
	allowedPrincipals := strutil.RemoveDuplicates(strutil.ParseStringSlice(principalsAllowedByRole, ","), false)
	switch {
	case len(parsedPrincipals) == 0:
		// There is nothing to process
		return nil, nil
	case len(allowedPrincipals) == 0:
		// User has requested principals to be set, but role is not configured
		// with any principals
		return nil, fmt.Errorf("role is not configured to allow any principles")
	default:
		// Role was explicitly configured to allow any principal.
		if principalsAllowedByRole == "*" {
			return parsedPrincipals, nil
		}

		for _, principal := range parsedPrincipals {
			if !validatePrincipal(allowedPrincipals, principal) {
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

func (b *backend) calculateExtensions(data *framework.FieldData, role *sshRole) (map[string]string, error) {
	unparsedExtensions := data.Get("extensions").(map[string]interface{})
	if len(unparsedExtensions) == 0 {
		return role.DefaultExtensions, nil
	}

	extensions := convertMapToStringValue(unparsedExtensions)

	if role.AllowedExtensions != "" {
		notAllowed := []string{}
		allowedExtensions := strings.Split(role.AllowedExtensions, ",")

		for extension := range extensions {
			if !strutil.StrListContains(allowedExtensions, extension) {
				notAllowed = append(notAllowed, extension)
			}
		}

		if len(notAllowed) != 0 {
			return nil, fmt.Errorf("extensions %v are not on allowed list", notAllowed)
		}
	}

	return extensions, nil
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
	if len(role.AllowedUserKeyLengths) != 0 {
		var kstr string
		var kbits int

		switch k := publickey.(type) {
		case ssh.CryptoPublicKey:
			ff := k.CryptoPublicKey()
			switch k := ff.(type) {
			case *rsa.PublicKey:
				kstr = "rsa"
				kbits = k.N.BitLen()
			case *dsa.PublicKey:
				kstr = "dsa"
				kbits = k.Parameters.P.BitLen()
			case *ecdsa.PublicKey:
				kstr = "ecdsa"
				kbits = k.Curve.Params().BitSize
			case ed25519.PublicKey:
				kstr = "ed25519"
			default:
				return fmt.Errorf("public key type of %s is not allowed", kstr)
			}
		default:
			return fmt.Errorf("pubkey not suitable for crypto (expected ssh.CryptoPublicKey but found %T)", k)
		}

		if value, ok := role.AllowedUserKeyLengths[kstr]; ok {
			var pass bool
			switch kstr {
			case "rsa":
				if kbits == value {
					pass = true
				}
			case "dsa":
				if kbits == value {
					pass = true
				}
			case "ecdsa":
				if kbits == value {
					pass = true
				}
			case "ed25519":
				// ed25519 public keys are always 256 bits in length,
				// so there is no need to inspect their value
				pass = true
			}

			if !pass {
				return fmt.Errorf("key is of an invalid size: %v", kbits)
			}

		} else {
			return fmt.Errorf("key type of %s is not allowed", kstr)
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

	certificate := &ssh.Certificate{
		Serial:          serialNumber.Uint64(),
		Key:             b.PublicKey,
		KeyId:           b.KeyID,
		ValidPrincipals: b.ValidPrincipals,
		ValidAfter:      uint64(now.Add(-30 * time.Second).In(time.UTC).Unix()),
		ValidBefore:     uint64(now.Add(b.TTL).In(time.UTC).Unix()),
		CertType:        b.CertificateType,
		Permissions: ssh.Permissions{
			CriticalOptions: b.CriticalOptions,
			Extensions:      b.Extensions,
		},
	}

	err = certificate.SignCert(rand.Reader, b.Signer)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signed SSH key")
	}

	return certificate, nil
}
