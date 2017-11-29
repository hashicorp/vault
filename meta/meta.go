package meta

import (
	"bufio"
	"flag"
	"io"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/helper/flag-slice"
	"github.com/mitchellh/cli"
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

type TokenHelperFunc func() (token.TokenHelper, error)

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetServer  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetServer
)

var (
	additionalOptionsUsage = func() string {
		return `
  -wrap-ttl=""            Indicates that the response should be wrapped in a
                          cubbyhole token with the requested TTL. The response
                          can be fetched by calling the "sys/wrapping/unwrap"
                          endpoint, passing in the wrapping token's ID. This
                          is a numeric string with an optional suffix
                          "s", "m", or "h"; if no suffix is specified it will
                          be parsed as seconds. May also be specified via
                          VAULT_WRAP_TTL.

  -policy-override        Indicates that any soft-mandatory Sentinel policies
                          be overridden.
`
	}
)

// Meta contains the meta-options and functionality that nearly every
// Vault command inherits.
type Meta struct {
	ClientToken string
	Ui          cli.Ui

	// The things below can be set, but aren't common
	ForceAddress string // Address to force for API clients

	// These are set by the command line flags.
	flagAddress        string
	flagCACert         string
	flagCAPath         string
	flagClientCert     string
	flagClientKey      string
	flagWrapTTL        string
	flagInsecure       bool
	flagMFA            []string
	flagPolicyOverride bool

	// Queried if no token can be found
	TokenHelper TokenHelperFunc
}

func (m *Meta) DefaultWrappingLookupFunc(operation, path string) string {
	if m.flagWrapTTL != "" {
		return m.flagWrapTTL
	}

	return api.DefaultWrappingLookupFunc(operation, path)
}

// Client returns the API client to a Vault server given the configured
// flag settings for this command.
func (m *Meta) Client() (*api.Client, error) {
	config := api.DefaultConfig()

	if m.flagAddress != "" {
		config.Address = m.flagAddress
	}
	if m.ForceAddress != "" {
		config.Address = m.ForceAddress
	}
	// If we need custom TLS configuration, then set it
	if m.flagCACert != "" || m.flagCAPath != "" || m.flagClientCert != "" || m.flagClientKey != "" || m.flagInsecure {
		t := &api.TLSConfig{
			CACert:        m.flagCACert,
			CAPath:        m.flagCAPath,
			ClientCert:    m.flagClientCert,
			ClientKey:     m.flagClientKey,
			TLSServerName: "",
			Insecure:      m.flagInsecure,
		}
		if err := config.ConfigureTLS(t); err != nil {
			return nil, err
		}
	}

	// Build the client
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetWrappingLookupFunc(m.DefaultWrappingLookupFunc)

	var mfaCreds []string

	// Extract the MFA credentials from environment variable first
	if os.Getenv(api.EnvVaultMFA) != "" {
		mfaCreds = []string{os.Getenv(api.EnvVaultMFA)}
	}

	// If CLI MFA flags were supplied, prefer that over environment variable
	if len(m.flagMFA) != 0 {
		mfaCreds = m.flagMFA
	}

	client.SetMFACreds(mfaCreds)

	client.SetPolicyOverride(m.flagPolicyOverride)

	// If we have a token directly, then set that
	token := m.ClientToken

	// Try to set the token to what is already stored
	if token == "" {
		token = client.Token()
	}

	// If we don't have a token, check the token helper
	if token == "" {
		if m.TokenHelper != nil {
			// If we have a token, then set that
			tokenHelper, err := m.TokenHelper()
			if err != nil {
				return nil, err
			}
			token, err = tokenHelper.Get()
			if err != nil {
				return nil, err
			}
		}
	}

	// Set the token
	if token != "" {
		client.SetToken(token)
	}

	return client, nil
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	// FlagSetServer tells us to enable the settings for selecting
	// the server information.
	if fs&FlagSetServer != 0 {
		f.StringVar(&m.flagAddress, "address", "", "")
		f.StringVar(&m.flagCACert, "ca-cert", "", "")
		f.StringVar(&m.flagCAPath, "ca-path", "", "")
		f.StringVar(&m.flagClientCert, "client-cert", "", "")
		f.StringVar(&m.flagClientKey, "client-key", "", "")
		f.StringVar(&m.flagWrapTTL, "wrap-ttl", "", "")
		f.BoolVar(&m.flagInsecure, "insecure", false, "")
		f.BoolVar(&m.flagInsecure, "tls-skip-verify", false, "")
		f.BoolVar(&m.flagPolicyOverride, "policy-override", false, "")
		f.Var((*sliceflag.StringFlag)(&m.flagMFA), "mfa", "")
	}

	// Create an io.Writer that writes to our Ui properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.Ui.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

// GeneralOptionsUsage returns the usage documentation for commonly
// available options
func GeneralOptionsUsage() string {
	general := `
  -address=addr           The address of the Vault server.
                          Overrides the VAULT_ADDR environment variable if set.

  -ca-cert=path           Path to a PEM encoded CA cert file to use to
                          verify the Vault server SSL certificate.
                          Overrides the VAULT_CACERT environment variable if set.

  -ca-path=path           Path to a directory of PEM encoded CA cert files
                          to verify the Vault server SSL certificate. If both
                          -ca-cert and -ca-path are specified, -ca-cert is used.
                          Overrides the VAULT_CAPATH environment variable if set.

  -client-cert=path       Path to a PEM encoded client certificate for TLS
                          authentication to the Vault server. Must also specify
                          -client-key. Overrides the VAULT_CLIENT_CERT
                          environment variable if set.

  -client-key=path        Path to an unencrypted PEM encoded private key
                          matching the client certificate from -client-cert.
                          Overrides the VAULT_CLIENT_KEY environment variable
                          if set.

  -tls-skip-verify        Do not verify TLS certificate. This is highly
                          not recommended. Verification will also be skipped
                          if VAULT_SKIP_VERIFY is set.
`

	general += additionalOptionsUsage()
	return general
}
