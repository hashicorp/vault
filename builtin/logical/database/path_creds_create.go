package database

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	v5 "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathCredsCreate(b *databaseBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "creds/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathCredsCreateRead(),
			},

			HelpSynopsis:    pathCredsCreateReadHelpSyn,
			HelpDescription: pathCredsCreateReadHelpDesc,
		},
		{
			Pattern: "static-creds/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the static role.",
				},
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: b.pathStaticCredsRead(),
			},

			HelpSynopsis:    pathStaticCredsReadHelpSyn,
			HelpDescription: pathStaticCredsReadHelpDesc,
		},
	}
}

func (b *databaseBackend) pathCredsCreateRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		// Get the role
		role, err := b.Role(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", name)), nil
		}

		dbConfig, err := b.DatabaseConfig(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		// If role name isn't in the database's allowed roles, send back a
		// permission denied.
		if !strutil.StrListContains(dbConfig.AllowedRoles, "*") && !strutil.StrListContainsGlob(dbConfig.AllowedRoles, name) {
			return nil, fmt.Errorf("%q is not an allowed role", name)
		}

		// Get the Database object
		dbi, err := b.GetConnection(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		dbi.RLock()
		defer dbi.RUnlock()

		ttl, _, err := framework.CalculateTTL(b.System(), 0, role.DefaultTTL, 0, role.MaxTTL, 0, time.Time{})
		if err != nil {
			return nil, err
		}
		expiration := time.Now().Add(ttl)
		// Adding a small buffer since the TTL will be calculated again after this call
		// to ensure the database credential does not expire before the lease
		expiration = expiration.Add(5 * time.Second)

		newUserReq := v5.NewUserRequest{
			UsernameConfig: v5.UsernameMetadata{
				DisplayName: req.DisplayName,
				RoleName:    name,
			},
			Statements: v5.Statements{
				Commands: role.Statements.Creation,
			},
			RollbackStatements: v5.Statements{
				Commands: role.Statements.Rollback,
			},
			Expiration: expiration,
		}

		// Set the credential based on the role's credential type.
		// TODO: can we move this into a credential library? credsutil is unused now.
		var credential interface{}
		switch role.CredentialType {
		case v5.CredentialTypePassword:
			password, err := dbi.database.GeneratePassword(ctx, b.System(), dbConfig.PasswordPolicy)
			if err != nil {
				b.CloseIfShutdown(dbi, err)
				return nil, fmt.Errorf("unable to generate password: %w", err)
			}

			credential = password
			newUserReq.Password = password
			newUserReq.CredentialType = v5.CredentialTypePassword
		case v5.CredentialTypeRSA2048PrivateKey:
			public, private, err := b.generateRSAKeypair()
			if err != nil {
				b.CloseIfShutdown(dbi, err)
				return nil, fmt.Errorf("unable to generate keypair: %w", err)
			}

			credential = string(private)
			newUserReq.PublicKey = string(public)
			newUserReq.CredentialType = v5.CredentialTypeRSA2048PrivateKey
		case v5.CredentialTypeClientCertificate:
			cb, subject, err := b.generateClientCertificate(role, expiration)
			if err != nil {
				b.CloseIfShutdown(dbi, err)
				return nil, fmt.Errorf("unable to generate client certificate: %w", err)
			}

			credential = cb
			newUserReq.UsernameConfig.Subject = subject
			newUserReq.CredentialType = v5.CredentialTypeClientCertificate
		}

		// Overwriting the password in the event this is a legacy database
		// plugin and the provided password is ignored
		newUserResp, password, err := dbi.database.NewUser(ctx, newUserReq)
		if err != nil {
			b.CloseIfShutdown(dbi, err)
			return nil, err
		}

		// Database plugins using the v4 interface generate and return the password.
		// If the returned password is not equal to the credential, then we need to
		// set the credential to the password returned from the v4 database plugin.
		if role.CredentialType == v5.CredentialTypePassword && password != credential {
			credential = password
		}

		respData := map[string]interface{}{
			"username": newUserResp.Username,
		}

		// TODO: another type for client certificate w/ password?
		switch role.CredentialType {
		case v5.CredentialTypePassword:
			respData["password"] = credential
		case v5.CredentialTypeRSA2048PrivateKey:
			respData["private_key"] = credential
		case v5.CredentialTypeClientCertificate:
			cb, ok := credential.(*certutil.CertBundle)
			if !ok {
				return nil, errors.New("type assertion failed on client cert bundle")
			}

			respData["certificate"] = cb.Certificate
			respData["private_key"] = cb.PrivateKey
			respData["private_key_type"] = cb.PrivateKeyType
		}

		internal := map[string]interface{}{
			"username":              newUserResp.Username,
			"role":                  name,
			"db_name":               role.DBName,
			"revocation_statements": role.Statements.Revocation,
		}
		resp := b.Secret(SecretCredsType).Response(respData, internal)
		resp.Secret.TTL = role.DefaultTTL
		resp.Secret.MaxTTL = role.MaxTTL
		return resp, nil
	}
}

func (b *databaseBackend) pathStaticCredsRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		name := data.Get("name").(string)

		role, err := b.StaticRole(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if role == nil {
			return logical.ErrorResponse("unknown role: %s", name), nil
		}

		dbConfig, err := b.DatabaseConfig(ctx, req.Storage, role.DBName)
		if err != nil {
			return nil, err
		}

		// If role name isn't in the database's allowed roles, send back a
		// permission denied.
		if !strutil.StrListContains(dbConfig.AllowedRoles, "*") && !strutil.StrListContainsGlob(dbConfig.AllowedRoles, name) {
			return nil, fmt.Errorf("%q is not an allowed role", name)
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"username":            role.StaticAccount.Username,
				"password":            role.StaticAccount.Password,
				"ttl":                 role.StaticAccount.PasswordTTL().Seconds(),
				"rotation_period":     role.StaticAccount.RotationPeriod.Seconds(),
				"last_vault_rotation": role.StaticAccount.LastVaultRotation,
			},
		}, nil
	}
}

// generateRSAKeypair returns PEM encodings of an RSA public and private key pair.
// The first return value is the PEM encoding of the PKIX marshalled public key.
// The second return value is the PEM encoding of the PKCS #8 marshalled private key.
func (b *databaseBackend) generateRSAKeypair() ([]byte, []byte, error) {
	// TODO: parameterize bits on role
	key, err := rsa.GenerateKey(b.GetRandomReader(), 2048)
	if err != nil {
		return nil, nil, err
	}

	// Marshal the public key to PKIX, ASN.1 DER form.
	public, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return nil, nil, err
	}

	// Marshal the private key to PKCS #8, ASN.1 DER form
	private, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	// Create PEM blocks
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: public,
	}
	privateBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: private,
	}

	// TODO: x509.EncryptPEMBlock is deprecated, so not sure how we'll support
	//       password encrypted private keys quite yet.

	return pem.EncodeToMemory(publicBlock), pem.EncodeToMemory(privateBlock), nil
}

func (b *databaseBackend) generateClientCertificate(role *roleEntry, expiration time.Time) (*certutil.CertBundle, string, error) {
	parsedCABundle, err := role.CertBundle.ToParsedCertBundle()
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate client certificate: %w", err)
	}

	caSign := &certutil.CAInfoBundle{
		ParsedCertBundle: *parsedCABundle,
		URLs: &certutil.URLEntries{
			IssuingCertificates:   []string{},
			CRLDistributionPoints: []string{},
			OCSPServers:           []string{},
		},
	}

	subject := pkix.Name{
		CommonName: role.CommonName,
		// Additional subject DN options intentionally omitted for now
	}

	creation := &certutil.CreationBundle{
		Params: &certutil.CreationParameters{
			Subject: subject,
			// Additional SAN-related fields intentionally omitted for now
			// SANs below:
			// DNSNames:                      strutil.RemoveDuplicates(dnsNames, false),
			// EmailAddresses:                strutil.RemoveDuplicates(emailAddresses, false),
			// IPAddresses:                   ipAddresses,
			// URIs:                          URIs,
			// OtherSANs:                     otherSANs,
			KeyType:       "rsa", // TODO: parameterize on role
			KeyBits:       2048,  // TODO: parameterize on role
			SignatureBits: 256,   // TODO: parameterize on role
			NotAfter:      expiration,
			KeyUsage:      x509.KeyUsageDigitalSignature,
			ExtKeyUsage:   certutil.ClientAuthExtKeyUsage,
			// ExtKeyUsageOIDs:               data.role.ExtKeyUsageOIDs,
			// PolicyIdentifiers:             data.role.PolicyIdentifiers,
			BasicConstraintsValidForNonCA: true,
			NotBeforeDuration:             0,
			URLs: &certutil.URLEntries{
				IssuingCertificates:   []string{},
				CRLDistributionPoints: []string{},
				OCSPServers:           []string{},
			},
		},
		SigningBundle: caSign,
	}

	parsedClientBundle, err := certutil.CreateCertificateWithRandomSource(creation, b.GetRandomReader())
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate client certificate: %w", err)
	}

	cb, err := parsedClientBundle.ToCertBundle()
	if err != nil {
		return nil, "", fmt.Errorf("error converting raw cert bundle to cert bundle: %w", err)
	}

	return cb, subject.String(), nil
}

const pathCredsCreateReadHelpSyn = `
Request database credentials for a certain role.
`

const pathCredsCreateReadHelpDesc = `
This path reads database credentials for a certain role. The
database credentials will be generated on demand and will be automatically
revoked when the lease is up.
`

const pathStaticCredsReadHelpSyn = `
Request database credentials for a certain static role. These credentials are
rotated periodically.
`

const pathStaticCredsReadHelpDesc = `
This path reads database credentials for a certain static role. The database
credentials are rotated periodically according to their configuration, and will
return the same password until they are rotated.
`
