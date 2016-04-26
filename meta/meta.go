package meta

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
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

// Meta contains the meta-options and functionality that nearly every
// Vault command inherits.
type Meta struct {
	ClientToken string
	Ui          cli.Ui

	// The things below can be set, but aren't common
	ForceAddress string // Address to force for API clients

	// These are set by the command line flags.
	flagAddress    string
	flagCACert     string
	flagCAPath     string
	flagClientCert string
	flagClientKey  string
	flagInsecure   bool

	// Queried if no token can be found
	TokenHelper TokenHelperFunc
}

// Client returns the API client to a Vault server given the configured
// flag settings for this command.
func (m *Meta) Client() (*api.Client, error) {
	config := api.DefaultConfig()

	err := config.ReadEnvironment()
	if err != nil {
		return nil, errwrap.Wrapf("error reading environment: {{err}}", err)
	}

	if m.flagAddress != "" {
		config.Address = m.flagAddress
	}
	if m.ForceAddress != "" {
		config.Address = m.ForceAddress
	}
	// If we need custom TLS configuration, then set it
	if m.flagCACert != "" || m.flagCAPath != "" || m.flagClientCert != "" || m.flagClientKey != "" || m.flagInsecure {
		// We may have set items from the environment so start with the
		// existing TLS config
		tlsConfig := config.HttpClient.Transport.(*http.Transport).TLSClientConfig

		var certPool *x509.CertPool
		var err error
		if m.flagCACert != "" {
			certPool, err = api.LoadCACert(m.flagCACert)
		} else if m.flagCAPath != "" {
			certPool, err = api.LoadCAPath(m.flagCAPath)
		}
		if err != nil {
			return nil, errwrap.Wrapf("Error setting up CA path: {{err}}", err)
		}

		if certPool != nil {
			tlsConfig.RootCAs = certPool
		}
		if m.flagInsecure {
			tlsConfig.InsecureSkipVerify = true
		}

		if m.flagClientCert != "" && m.flagClientKey != "" {
			tlsCert, err := tls.LoadX509KeyPair(m.flagClientCert, m.flagClientKey)
			if err != nil {
				return nil, err
			}
			tlsConfig.Certificates = []tls.Certificate{tlsCert}
		} else if m.flagClientCert != "" || m.flagClientKey != "" {
			return nil, fmt.Errorf("Both client cert and client key must be provided")
		}
	}

	// Build the client
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

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
		f.BoolVar(&m.flagInsecure, "insecure", false, "")
		f.BoolVar(&m.flagInsecure, "tls-skip-verify", false, "")
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

func (m *Meta) loadCACert(path string) (*x509.CertPool, error) {
	certs, err := m.loadCertFromPEM(path)
	if err != nil {
		return nil, fmt.Errorf("Error loading %s: %s", path, err)
	}

	result := x509.NewCertPool()
	for _, cert := range certs {
		result.AddCert(cert)
	}

	return result, nil
}

func (m *Meta) loadCAPath(path string) (*x509.CertPool, error) {
	result := x509.NewCertPool()
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		certs, err := m.loadCertFromPEM(path)
		if err != nil {
			return fmt.Errorf("Error loading %s: %s", path, err)
		}

		for _, cert := range certs {
			result.AddCert(cert)
		}
		return nil
	}

	return result, filepath.Walk(path, fn)
}

func (m *Meta) loadCertFromPEM(path string) ([]*x509.Certificate, error) {
	pemCerts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	certs := make([]*x509.Certificate, 0, 5)
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		certs = append(certs, cert)
	}

	return certs, nil
}

// GeneralOptionsUsage returns the usage documenation for commonly
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
                          -ca-cert and -ca-path are specified, -ca-path is used.
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
	return general
}
