package ssh

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/duration"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/ssh"
)

type signingBundle struct {
	Certificate string `json:"certificate" structs:"certificate" mapstructure:"certificate"`
}

type creationBundle struct {
	KeyId           string
	ValidPrincipals []string
	PublicKey       ssh.PublicKey
	CertificateType uint32
	TTL             time.Duration
	SigningBundle   signingBundle
	Role            *sshRole
	criticalOptions map[string]string
	extensions      map[string]string
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
				Type: framework.TypeString,
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

func (b *backend) pathSign(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	return b.pathSignCertificate(req, data, role)
}

func (b *backend) pathSignCertificate(req *logical.Request, data *framework.FieldData, role *sshRole) (*logical.Response, error) {
	publicKey := data.Get("public_key").(string)
	if publicKey == "" {
		return logical.ErrorResponse("missing public_key"), nil
	}

	userPublicKey, err := parsePublicSSHKey(publicKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("unable to decode \"public_key\" as SSH key: %s", err)), nil
	}

	keyId := data.Get("key_id").(string)
	if keyId == "" {
		keyId = req.DisplayName
	}

	// Note that these various functions always return "user errors" so we pass
	// them as 4xx values
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

	storedBundle, err := req.Storage.Get("config/ca_bundle")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch local CA certificate/key: %v", err)
	}
	if storedBundle == nil {
		return logical.ErrorResponse("backend must be configured with a CA certificate/key"), nil
	}

	var bundle signingBundle
	if err := storedBundle.DecodeJSON(&bundle); err != nil {
		return nil, fmt.Errorf("unable to decode local CA certificate/key: %v", err)
	}

	signingBundle := creationBundle{
		KeyId:           keyId,
		PublicKey:       userPublicKey,
		SigningBundle:   bundle,
		ValidPrincipals: parsedPrincipals,
		TTL:             ttl,
		CertificateType: certificateType,
		Role:            role,
		criticalOptions: criticalOptions,
		extensions:      extensions,
	}

	certificate, err := signingBundle.sign()
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
	if principalsAllowedByRole == "" {
		return nil, fmt.Errorf(`"role is not configured to allow any principles`)
	}

	validPrincipals := data.Get("valid_principals").(string)
	if validPrincipals == "" {
		if defaultPrincipal != "" {
			return []string{defaultPrincipal}, nil
		}

		return nil, fmt.Errorf(`"valid_principals" not supplied and no default set in the role`)
	}

	parsedPrincipals := strings.Split(validPrincipals, ",")

	// Role was explicitly configured to allow any principal.
	if principalsAllowedByRole == "*" {
		return parsedPrincipals, nil
	}

	allowedPrincipals := strings.Split(principalsAllowedByRole, ",")
	for _, principal := range parsedPrincipals {
		if !validatePrincipal(allowedPrincipals, principal) {
			return nil, fmt.Errorf(`%v is not a valid value for "valid_principals"`, principal)
		}
	}

	return parsedPrincipals, nil
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
			return 0, errors.New(`"cert_type" 'user' is not allowed by role`)
		}
		certificateType = ssh.UserCert
	case "host":
		if !role.AllowHostCertificates {
			return 0, errors.New(`"cert_type" 'host' is not allowed by role`)
		}
		certificateType = ssh.HostCert
	default:
		return 0, errors.New(`"cert_type" must be either 'user' or 'host'`)
	}

	return certificateType, nil
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
			return nil, fmt.Errorf("Critical options not on allowed list: %v", notAllowedOptions)
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
			return nil, fmt.Errorf("Extensions not on allowed list: %v", notAllowed)
		}
	}

	return extensions, nil
}

func (b *backend) calculateTTL(data *framework.FieldData, role *sshRole) (time.Duration, error) {

	var ttl, maxTTL time.Duration
	var ttlField string
	ttlFieldInt, ok := data.GetOk("ttl")
	if !ok {
		ttlField = role.TTL
	} else {
		ttlField = ttlFieldInt.(string)
	}

	if len(ttlField) == 0 {
		ttl = b.System().DefaultLeaseTTL()
	} else {
		var err error
		ttl, err = duration.ParseDurationSecond(ttlField)
		if err != nil {
			return 0, fmt.Errorf("invalid requested ttl: %s", err)
		}
	}

	if len(role.MaxTTL) == 0 {
		maxTTL = b.System().MaxLeaseTTL()
	} else {
		var err error
		maxTTL, err = duration.ParseDurationSecond(role.MaxTTL)
		if err != nil {
			return 0, fmt.Errorf("invalid requested max ttl: %s", err)
		}
	}

	if ttl > maxTTL {
		// Don't error if they were using system defaults, only error if
		// they specifically chose a bad TTL
		if len(ttlField) == 0 {
			ttl = maxTTL
		} else {
			return 0, fmt.Errorf("ttl is larger than maximum allowed (%d)", maxTTL/time.Second)
		}
	}

	return ttl, nil
}

func (b *creationBundle) sign() (*ssh.Certificate, error) {
	signingKey, err := ssh.ParsePrivateKey([]byte(b.SigningBundle.Certificate))
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("stored SSH signing key cannot be parsed: %v", err)}
	}

	serialNumber, err := certutil.GenerateSerialNumber()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	certificate := &ssh.Certificate{
		Serial:          serialNumber.Uint64(),
		Key:             b.PublicKey,
		KeyId:           b.KeyId,
		ValidPrincipals: b.ValidPrincipals,
		ValidAfter:      uint64(now.Add(-30 * time.Second).In(time.UTC).Unix()),
		ValidBefore:     uint64(now.Add(b.TTL).In(time.UTC).Unix()),
		CertType:        b.CertificateType,
		Permissions: ssh.Permissions{
			CriticalOptions: b.criticalOptions,
			Extensions:      b.extensions,
		},
	}

	err = certificate.SignCert(rand.Reader, signingKey)
	if err != nil {
		return nil, errutil.InternalError{Err: "Failed to generate signed SSH key"}
	}

	return certificate, nil
}
