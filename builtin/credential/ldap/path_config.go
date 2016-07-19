package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-ldap/ldap"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfig(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `config`,
		Fields: map[string]*framework.FieldSchema{
			"url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "ldap URL to connect to (default: ldap://127.0.0.1)",
			},

			"userdn": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "LDAP domain to use for users (eg: ou=People,dc=example,dc=org)",
			},

			"binddn": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "LDAP DN for searching for the user DN (optional)",
			},

			"bindpass": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "LDAP password for searching for the user DN (optional)",
			},

			"groupdn": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "LDAP domain to use for groups (eg: ou=Groups,dc=example,dc=org)",
			},

			"upndomain": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Enables userPrincipalDomain login with [username]@UPNDomain (optional)",
			},

			"userattr": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Attribute used for users (default: cn)",
			},

			"certificate": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded (optional)",
			},

			"discoverdn": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Use anonymous bind to discover the bind DN of a user (optional)",
			},

			"insecure_tls": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Skip LDAP server SSL Certificate verification - VERY insecure (optional)",
			},

			"starttls": &framework.FieldSchema{
				Type:        framework.TypeBool,
				Description: "Issue a StartTLS command after establishing unencrypted connection (optional)",
			},

			"tls_min_version": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Minimum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathConfigRead,
			logical.UpdateOperation: b.pathConfigWrite,
		},

		HelpSynopsis:    pathConfigHelpSyn,
		HelpDescription: pathConfigHelpDesc,
	}
}

func (b *backend) Config(req *logical.Request) (*ConfigEntry, error) {
	entry, err := req.Storage.Get("config")
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var result ConfigEntry
	result.SetDefaults()
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (b *backend) pathConfigRead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(cfg).Map(),
	}, nil
}

func (b *backend) pathConfigWrite(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	cfg := &ConfigEntry{}
	url := d.Get("url").(string)
	if url != "" {
		cfg.Url = strings.ToLower(url)
	}
	userattr := d.Get("userattr").(string)
	if userattr != "" {
		cfg.UserAttr = strings.ToLower(userattr)
	}
	userdn := d.Get("userdn").(string)
	if userdn != "" {
		cfg.UserDN = userdn
	}
	groupdn := d.Get("groupdn").(string)
	if groupdn != "" {
		cfg.GroupDN = groupdn
	}
	upndomain := d.Get("upndomain").(string)
	if groupdn != "" {
		cfg.UPNDomain = upndomain
	}
	certificate := d.Get("certificate").(string)
	if certificate != "" {
		cfg.Certificate = certificate
	}
	insecureTLS := d.Get("insecure_tls").(bool)
	if insecureTLS {
		cfg.InsecureTLS = insecureTLS
	}
	cfg.TLSMinVersion = d.Get("tls_min_version").(string)
	if cfg.TLSMinVersion == "" {
		return logical.ErrorResponse("failed to get 'tls_min_version' value"), nil
	}

	var ok bool
	_, ok = tlsutil.TLSLookup[cfg.TLSMinVersion]
	if !ok {
		return logical.ErrorResponse("invalid 'tls_min_version'"), nil
	}

	startTLS := d.Get("starttls").(bool)
	if startTLS {
		cfg.StartTLS = startTLS
	}
	bindDN := d.Get("binddn").(string)
	if bindDN != "" {
		cfg.BindDN = bindDN
	}
	bindPass := d.Get("bindpass").(string)
	if bindPass != "" {
		cfg.BindPassword = bindPass
	}
	discoverDN := d.Get("discoverdn").(bool)
	if discoverDN {
		cfg.DiscoverDN = discoverDN
	}

	// Try to connect to the LDAP server, to validate the URL configuration
	// We can also check the URL at this stage, as anything else would probably
	// require authentication.
	conn, cerr := cfg.DialLDAP()
	if cerr != nil {
		return logical.ErrorResponse(cerr.Error()), nil
	}
	if conn == nil {
		return logical.ErrorResponse("invalid connection returned from LDAP dial"), nil
	}
	conn.Close()

	entry, err := logical.StorageEntryJSON("config", cfg)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type ConfigEntry struct {
	Url           string `json:"url" structs:"url" mapstructure:"url"`
	UserDN        string `json:"userdn" structs:"userdn" mapstructure:"userdn"`
	GroupDN       string `json:"groupdn" structs:"groupdn" mapstructure:"groupdn"`
	UPNDomain     string `json:"upndomain" structs:"upndomain" mapstructure:"upndomain"`
	UserAttr      string `json:"userattr" structs:"userattr" mapstructure:"userattr"`
	Certificate   string `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	InsecureTLS   bool   `json:"insecure_tls" structs:"insecure_tls" mapstructure:"insecure_tls"`
	StartTLS      bool   `json:"starttls" structs:"starttls" mapstructure:"starttls"`
	BindDN        string `json:"binddn" structs:"binddn" mapstructure:"binddn"`
	BindPassword  string `json:"bindpass" structs:"bindpass" mapstructure:"bindpass"`
	DiscoverDN    bool   `json:"discoverdn" structs:"discoverdn" mapstructure:"discoverdn"`
	TLSMinVersion string `json:"tls_min_version" structs:"tls_min_version" mapstructure:"tls_min_version"`
}

func (c *ConfigEntry) GetTLSConfig(host string) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		ServerName: host,
	}

	if c.TLSMinVersion != "" {
		tlsMinVersion, ok := tlsutil.TLSLookup[c.TLSMinVersion]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_min_version' in config")
		}
		tlsConfig.MinVersion = tlsMinVersion
	}

	if c.InsecureTLS {
		tlsConfig.InsecureSkipVerify = true
	}
	if c.Certificate != "" {
		caPool := x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM([]byte(c.Certificate))
		if !ok {
			return nil, fmt.Errorf("could not append CA certificate")
		}
		tlsConfig.RootCAs = caPool
	}
	return tlsConfig, nil
}

func (c *ConfigEntry) DialLDAP() (*ldap.Conn, error) {

	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, err
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		host = u.Host
	}

	var conn *ldap.Conn
	var tlsConfig *tls.Config
	switch u.Scheme {
	case "ldap":
		if port == "" {
			port = "389"
		}
		conn, err = ldap.Dial("tcp", host+":"+port)
		if c.StartTLS {
			tlsConfig, err = c.GetTLSConfig(host)
			if err != nil {
				break
			}
			err = conn.StartTLS(tlsConfig)
		}
	case "ldaps":
		if port == "" {
			port = "636"
		}
		tlsConfig, err = c.GetTLSConfig(host)
		if err != nil {
			break
		}
		conn, err = ldap.DialTLS("tcp", host+":"+port, tlsConfig)
	default:
		return nil, fmt.Errorf("invalid LDAP scheme")
	}
	if err != nil {
		return nil, fmt.Errorf("cannot connect to LDAP: %v", err)
	}

	return conn, nil
}

func (c *ConfigEntry) SetDefaults() {
	c.Url = "ldap://127.0.0.1"
	c.UserAttr = "cn"
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
