package ssh

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/crypto/ssh"
	"log"
	"strconv"
	"strings"
	"time"
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
		return nil, errutil.UserError{Err: "\"public_key\" is empty"}
	}

	keyParts := strings.Split(publicKey, " ")
	if len(keyParts) > 1 {
		// Someone has sent the 'full' public key rather than just the base64 encoded part that the ssh library wants
		publicKey = keyParts[1]
	}

	decodedKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, errutil.UserError{Err: "Unable to decode \"public_key\" as SSH key"}
	}

	userPublicKey, err := ssh.ParsePublicKey([]byte(decodedKey))
	if err != nil {
		log.Printf("Failed to parse key: %s", err)
		return nil, errutil.UserError{Err: "Unable to parse \"public_key\" as SSH key"}
	}

	keyId := data.Get("key_id").(string)
	if keyId == "" {
		keyId = req.DisplayName
	}

	certificateType, err := b.calculateCertificateType(data, role)
	if err != nil {
		return nil, err
	}

	parsedPrincipals, err := b.calculateValidPrincipals(data, role, certificateType == ssh.HostCert)
	if err != nil {
		return nil, err
	}

	ttl, err := b.calculateTtl(data, role)
	if err != nil {
		return nil, err
	}

	criticalOptions, err := b.calculateCriticalOptions(data, role)
	if err != nil {
		return nil, err
	}

	extensions, err := b.calculateExtensions(data, role)
	if err != nil {
		return nil, err
	}

	storedBundle, err := req.Storage.Get("config/ssh_certificate_bundle")
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch local CA certificate/key: %v", err)}
	}
	if storedBundle == nil {
		return nil, errutil.UserError{Err: "backend must be configured with a CA certificate/key"}
	}

	var bundle signingBundle
	if err := storedBundle.DecodeJSON(&bundle); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode local CA certificate/key: %v", err)}
	}

	signingBundle := createSigningBundle(keyId, userPublicKey, parsedPrincipals, ttl, certificateType, bundle, role, criticalOptions, extensions)

	certificate, err := signingBundle.sign()
	if err != nil {
		return nil, err
	}

	signedSSHCertificate := string(ssh.MarshalAuthorizedKey(certificate))

	response := b.Secret(SecretCertsType).Response(
		map[string]interface{}{
			"serial_number": strconv.FormatUint(certificate.Serial, 16),
			"signed_key":    signedSSHCertificate,
		},
		map[string]interface{}{
			"serial_number": strconv.FormatUint(certificate.Serial, 16),
			"signed_key":    signedSSHCertificate,
		})

	return response, nil
}

func (b *backend) calculateValidPrincipals(data *framework.FieldData, role *sshRole, isHostCertificate bool) ([]string, error) {
	validPrincipals := data.Get("valid_principals").(string)
	if validPrincipals == "" {
		if role.DefaultUser != "" {
			return []string {role.DefaultUser}, nil
		}
		if role.AllowedUsers == "" {
			return []string{}, nil
		}

		return nil, errutil.UserError{Err: `"valid_principals" value required by role`}
	}

	parsedPrincipals := strings.Split(validPrincipals, ",")
	if role.AllowedUsers == "" {
		return nil, errutil.UserError{Err: `"valid_principals" not in allowed list`}
	}

	// Role was explicitly configured to allow any principal.
	if role.AllowedUsers == "*" {
		return parsedPrincipals, nil
	}

	allowedPrincipals := strings.Split(role.AllowedUsers, ",")
	for _, principal := range parsedPrincipals {
		if !validateValidPrincipal(principal, allowedPrincipals, role, isHostCertificate) {
			return nil, errutil.UserError{Err: fmt.Sprintf(`%v is not a valid value for "valid_principals"`, principal)}
		}
	}

	return parsedPrincipals, nil
}

func validateValidPrincipal(validPrincipal string, allowedPrincipals []string, role *sshRole, isHostCertificate bool) bool {

	if isHostCertificate {
		for _, allowedPrincipal := range allowedPrincipals {
			if allowedPrincipal == validPrincipal && role.AllowBareDomains {
				return true
			}
			if role.AllowSubdomains && strings.HasSuffix(validPrincipal, "." + allowedPrincipal) {
				return true
			}
		}

		return false
	} else {

		return contains(allowedPrincipals, validPrincipal)
	}
}

func (b *backend) calculateCertificateType(data *framework.FieldData, role *sshRole) (uint32, error) {
	requestedCertificateType := data.Get("cert_type").(string)

	var certificateType uint32
	switch requestedCertificateType {
	case "user":
		if !role.AllowUserCertificates {
			return 0, errutil.UserError{Err: `"cert_type" 'user' is not allowed by role`}
		}
		certificateType = ssh.UserCert
		break
	case "host":
		if !role.AllowHostCertificates {
			return 0, errutil.UserError{Err: `"cert_type" 'host' is not allowed by role`}
		}
		certificateType = ssh.HostCert
		break
	default:
		return 0, errutil.UserError{Err: "\"cert_type\" must be either 'user' or 'host'"}
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
			if !contains(allowedCriticalOptions, option) {
				notAllowedOptions = append(notAllowedOptions, option)
			}
		}

		if len(notAllowedOptions) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf("Critical options not on allowed list: %v", notAllowedOptions)}
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
			if !contains(allowedExtensions, extension) {
				notAllowed = append(notAllowed, extension)
			}
		}

		if len(notAllowed) != 0 {
			return nil, errutil.UserError{Err: fmt.Sprintf("Extensions not on allowed list: %v", notAllowed)}
		}
	}

	return extensions, nil
}

func (b *backend) calculateTtl(data *framework.FieldData, role *sshRole) (time.Duration, error) {

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
		ttl, err = time.ParseDuration(ttlField)
		if err != nil {
			return time.Nanosecond, errutil.UserError{Err: fmt.Sprintf("invalid requested ttl: %s", err)}
		}
	}

	if len(role.MaxTTL) == 0 {
		maxTTL = b.System().MaxLeaseTTL()
	} else {
		var err error
		maxTTL, err = time.ParseDuration(role.MaxTTL)
		if err != nil {
			return time.Nanosecond, errutil.UserError{Err: fmt.Sprintf("invalid ttl: %s", err)}
		}
	}

	if ttl > maxTTL {
		// Don't error if they were using system defaults, only error if
		// they specifically chose a bad TTL
		if len(ttlField) == 0 {
			ttl = maxTTL
		} else {
			return time.Nanosecond, errutil.UserError{Err: fmt.Sprintf("ttl is larger than maximum allowed (%d)", maxTTL/time.Second)}
		}
	}

	return ttl, nil
}

func createSigningBundle(keyId string, userPublicKey ssh.PublicKey, validPrincipals []string, duration time.Duration, certificateType uint32, sshCertificateBundle signingBundle,
	role *sshRole, criticalOptions, extensions map[string]string) creationBundle {

	return creationBundle{
		KeyId:           keyId,
		PublicKey:       userPublicKey,
		SigningBundle:   sshCertificateBundle,
		ValidPrincipals: validPrincipals,
		TTL:             duration,
		CertificateType: certificateType,
		Role:            role,
		criticalOptions: criticalOptions,
		extensions:      extensions,
	}
}

func (b *creationBundle) sign() (*ssh.Certificate, error) {

	signingKey, err := b.SigningBundle.toSSHSigningKey()
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

func (b *signingBundle) toSSHSigningKey() (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(b.Certificate))
}
