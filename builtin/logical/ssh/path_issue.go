package ssh

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/crypto/ssh"
)

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
				Description: "TBD",
				Required:    true,
			},
			"key_bits": {
				Type:        framework.TypeInt,
				Description: "TBD",
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
		HelpSynopsis:    "TBD - HelpSynopsis",
		HelpDescription: "TBD - HelpDescription",
	}
}

func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}

	if role.KeyType != "ca" {
		return logical.ErrorResponse("role key type \"%s\" not allowed to issue key pairs", role.KeyType), nil
	}

	keyType := data.Get("key_type").(string)
	keyBits := data.Get("key_bits").(int)

	keyTypeLengths, keyPresent := role.AllowedUserKeyTypesLengths[keyType]
	if keyPresent {
		var bitsPresent bool
		for _, kb := range keyTypeLengths {
			if keyBits == kb {
				bitsPresent = true
				break
			}
		}
		if !bitsPresent {
			// Also return list of allowed key bits?
			return logical.ErrorResponse("key_bits not in list of allowed values for key_type provided"), nil
		}
	}

	publicKey, privateKey, err := generateSSHKeyPair(b.Backend.GetRandomReader(), keyType, keyBits)
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

	err = b.validateSignedKeyRequirements(userPublicKey, role)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("public_key failed to meet the key requirements: %s", err)), nil
	}

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
		parsedPrincipals, err = b.calculateValidPrincipals(data, req, role, "", role.AllowedDomains, validateValidPrincipalForHosts(role))
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	} else {
		parsedPrincipals, err = b.calculateValidPrincipals(data, req, role, role.DefaultUser, role.AllowedUsers, strutil.StrListContains)
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

	extensions, err := b.calculateExtensions(data, req, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
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
			"serial_number":    strconv.FormatUint(certificate.Serial, 16),
			"certificate":      string(signedSSHCertificate),
			"private_key":      privateKey,
			"private_key_type": keyType,
		},
	}

	return response, nil
}
