package ssh

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

type keySpecs struct {
	Type string
	Bits int
}

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameWithAtRegex("role"),

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathIssue,
			},
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `The desired role with configuration for this request.`,
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: "Specifies the desired key type; must be `rsa`, `ed25519` or `ec`",
				Default:     "rsa",
			},
			"key_bits": {
				Type:        framework.TypeInt,
				Description: "Specifies the number of bits to use for the generated keys.",
				Default:     0,
			},
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `The requested Time To Live for the SSH certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be later than the role max TTL.`,
			},
			"valid_principals": {
				Type:        framework.TypeString,
				Description: `Valid principals, either usernames or hostnames, that the certificate should be signed for.`,
			},
			"cert_type": {
				Type:        framework.TypeString,
				Description: `Type of certificate to be created; either "user" or "host".`,
				Default:     "user",
			},
			"key_id": {
				Type:        framework.TypeString,
				Description: `Key id that the created certificate should have. If not specified, the display name of the token will be used.`,
			},
			"critical_options": {
				Type:        framework.TypeMap,
				Description: `Critical options that the certificate should be signed for.`,
			},
			"extensions": {
				Type:        framework.TypeMap,
				Description: `Extensions that the certificate should be signed for.`,
			},
		},
		HelpSynopsis:    pathIssueHelpSyn,
		HelpDescription: pathIssueHelpDesc,
	}
}

func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the role
	roleName := data.Get("role").(string)
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}

	if role.KeyType != "ca" {
		return logical.ErrorResponse("role key type '%s' not allowed to issue key pairs", role.KeyType), nil
	}

	// Validate and extract key specifications
	keySpecs, err := extractKeySpecs(role, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Issue certificate
	return b.pathIssueCertificate(ctx, req, data, role, keySpecs)
}

func (b *backend) pathIssueCertificate(ctx context.Context, req *logical.Request, data *framework.FieldData, role *sshRole, keySpecs *keySpecs) (*logical.Response, error) {
	publicKey, privateKey, err := generateSSHKeyPair(rand.Reader, keySpecs.Type, keySpecs.Bits)
	if err != nil {
		return nil, err
	}

	if publicKey == "" || privateKey == "" {
		return nil, fmt.Errorf("failed to generate or parse the keys")
	}

	// Sign key
	userPublicKey, err := parsePublicSSHKey(publicKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse public_key as SSH key: %s", err)), nil
	}

	response, err := b.pathSignIssueCertificateHelper(ctx, req, data, role, userPublicKey)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Additional to sign response
	response.Data["private_key"] = privateKey
	response.Data["private_key_type"] = keySpecs.Type

	return response, nil
}

func extractKeySpecs(role *sshRole, data *framework.FieldData) (*keySpecs, error) {
	keyType := data.Get("key_type").(string)
	keyBits := data.Get("key_bits").(int)
	keySpecs := keySpecs{
		Type: keyType,
		Bits: keyBits,
	}

	keyTypeToMapKey := createKeyTypeToMapKey(keyType, keyBits)

	if len(role.AllowedUserKeyTypesLengths) != 0 {
		var keyAllowed bool
		var bitsAllowed bool

	keyTypeAliasesLoop:
		for _, keyTypeAlias := range keyTypeToMapKey[keyType] {
			allowedValues, allowed := role.AllowedUserKeyTypesLengths[keyTypeAlias]
			if !allowed {
				continue
			}
			keyAllowed = true

			for _, value := range allowedValues {
				if value == keyBits {
					bitsAllowed = true
					break keyTypeAliasesLoop
				}
			}
		}

		if !keyAllowed {
			return nil, fmt.Errorf("key_type provided not in allowed_user_key_types")
		}

		if !bitsAllowed {
			return nil, fmt.Errorf("key_bits not in list of allowed values for key_type provided")
		}
	}

	return &keySpecs, nil
}

func (b *backend) pathSignIssueCertificateHelper(ctx context.Context, req *logical.Request, data *framework.FieldData, role *sshRole, publicKey ssh.PublicKey) (*logical.Response, error) {
	// Note that these various functions always return "user errors" so we pass
	// them as 4xx values
	keyID, err := b.calculateKeyID(data, req, role, publicKey)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	certificateType, err := b.calculateCertificateType(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	var parsedPrincipals []string
	if certificateType == ssh.HostCert {
		parsedPrincipals, err = b.calculateValidPrincipals(data, req, role, "", role.AllowedDomains, validateValidPrincipalForHosts(role))
		if err != nil {
			return logical.ErrorResponse(err.Error()), err
		}
	} else {
		parsedPrincipals, err = b.calculateValidPrincipals(data, req, role, role.DefaultUser, role.AllowedUsers, strutil.StrListContains)
		if err != nil {
			return logical.ErrorResponse(err.Error()), err
		}
	}

	ttl, err := b.calculateTTL(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	criticalOptions, err := b.calculateCriticalOptions(data, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	extensions, err := b.calculateExtensions(data, req, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	privateKeyEntry, err := caKey(ctx, req.Storage, caPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA private key: %w", err)
	}
	if privateKeyEntry == nil || privateKeyEntry.Key == "" {
		return nil, fmt.Errorf("failed to read CA private key")
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

const pathIssueHelpSyn = `
Request a signed key pair using a certain role with the provided details.
`

const pathIssueHelpDesc = `
`
