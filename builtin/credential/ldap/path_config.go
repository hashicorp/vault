package ldap

import (
	"context"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/ldaputil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields:  ldaputil.ConfigFields(),

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

/*
 * Construct ConfigEntry struct using stored configuration.
 */
func (b *backend) Config(ctx context.Context, req *logical.Request) (*ldaputil.ConfigEntry, error) {
	// Schema for ConfigEntry
	fd, err := b.getConfigFieldData()
	if err != nil {
		return nil, err
	}

	// Create a new ConfigEntry, filling in defaults where appropriate
	result, err := ldaputil.NewConfigEntry(fd)
	if err != nil {
		return nil, err
	}

	storedConfig, err := req.Storage.Get(ctx, "config")
	if err != nil {
		return nil, err
	}

	if storedConfig == nil {
		// No user overrides, return default configuration
		result.CaseSensitiveNames = new(bool)
		*result.CaseSensitiveNames = false

		return result, nil
	}

	// Deserialize stored configuration.
	// Fields not specified in storedConfig will retain their defaults.
	if err := storedConfig.DecodeJSON(&result); err != nil {
		return nil, err
	}

	var persistNeeded bool
	if result.CaseSensitiveNames == nil {
		// Upgrade from before switching to case-insensitive
		result.CaseSensitiveNames = new(bool)
		*result.CaseSensitiveNames = true
		persistNeeded = true
	}

	if persistNeeded && (b.System().LocalMount() || !b.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary|consts.ReplicationPerformanceStandby)) {
		entry, err := logical.StorageEntryJSON("config", result)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: cfg.PasswordlessMap(),
	}
	return resp, nil
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Build a ConfigEntry struct out of the supplied FieldData
	cfg, err := ldaputil.NewConfigEntry(d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// On write, if not specified, use false. We do this here so upgrade logic
	// works since it calls the same newConfigEntry function
	if cfg.CaseSensitiveNames == nil {
		cfg.CaseSensitiveNames = new(bool)
		*cfg.CaseSensitiveNames = false
	}

	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

/*
 * Returns FieldData describing our ConfigEntry struct schema
 */
func (b *backend) getConfigFieldData() (*framework.FieldData, error) {
	configPath := b.Route("config")

	if configPath == nil {
		return nil, logical.ErrUnsupportedPath
	}

	raw := make(map[string]interface{}, len(configPath.Fields))

	fd := framework.FieldData{
		Raw:    raw,
		Schema: configPath.Fields,
	}

	return &fd, nil
}

const pathConfigHelpSyn = `
Configure the LDAP server to connect to, along with its options.
`

const pathConfigHelpDesc = `
This endpoint allows you to configure the LDAP server to connect to and its
configuration options.

The LDAP URL can use either the "ldap://" or "ldaps://" schema. In the former
case, an unencrypted connection will be made with a default port of 389, unless
the "starttls" parameter is set to true, in which case TLS will be used. In the
latter case, a SSL connection will be established with a default port of 636.

## A NOTE ON ESCAPING

It is up to the administrator to provide properly escaped DNs. This includes
the user DN, bind DN for search, and so on.

The only DN escaping performed by this backend is on usernames given at login
time when they are inserted into the final bind DN, and uses escaping rules
defined in RFC 4514.

Additionally, Active Directory has escaping rules that differ slightly from the
RFC; in particular it requires escaping of '#' regardless of position in the DN
(the RFC only requires it to be escaped when it is the first character), and
'=', which the RFC indicates can be escaped with a backslash, but does not
contain in its set of required escapes. If you are using Active Directory and
these appear in your usernames, please ensure that they are escaped, in
addition to being properly escaped in your configured DNs.

For reference, see https://www.ietf.org/rfc/rfc4514.txt and
http://social.technet.microsoft.com/wiki/contents/articles/5312.active-directory-characters-to-escape.aspx
`
