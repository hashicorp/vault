package cassandra

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/logical"
)

// Query templates a query for us.
func substQuery(tpl string, data map[string]string) string {
	for k, v := range data {
		tpl = strings.Replace(tpl, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	return tpl
}

func createSession(cfg *sessionConfig, s logical.Storage) (*gocql.Session, error) {
	clusterConfig := gocql.NewCluster(strings.Split(cfg.Hosts, ",")...)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	clusterConfig.ProtoVersion = cfg.ProtocolVersion
	if clusterConfig.ProtoVersion == 0 {
		clusterConfig.ProtoVersion = 2
	}

	clusterConfig.Timeout = time.Duration(cfg.ConnectTimeout) * time.Second

	if cfg.TLS {
		var tlsConfig *tls.Config
		if len(cfg.Certificate) > 0 || len(cfg.IssuingCA) > 0 {
			if len(cfg.Certificate) > 0 && len(cfg.PrivateKey) == 0 {
				return nil, fmt.Errorf("found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(cfg.Certificate) > 0 {
				certBundle.Certificate = cfg.Certificate
				certBundle.PrivateKey = cfg.PrivateKey
			}
			if len(cfg.IssuingCA) > 0 {
				certBundle.IssuingCA = cfg.IssuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, errwrap.Wrapf("failed to parse certificate bundle: {{err}}", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("failed to get TLS configuration: tlsConfig: %#v; {{err}}", tlsConfig), err)
			}
			tlsConfig.InsecureSkipVerify = cfg.InsecureTLS

			if cfg.TLSMinVersion != "" {
				var ok bool
				tlsConfig.MinVersion, ok = tlsutil.TLSLookup[cfg.TLSMinVersion]
				if !ok {
					return nil, fmt.Errorf("invalid 'tls_min_version' in config")
				}
			} else {
				// MinVersion was not being set earlier. Reset it to
				// zero to gracefully handle upgrades.
				tlsConfig.MinVersion = 0
			}
		}

		clusterConfig.SslOpts = &gocql.SslOptions{
			Config: tlsConfig,
		}
	}

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, errwrap.Wrapf("error creating session: {{err}}", err)
	}

	// Verify the info
	err = session.Query(`LIST USERS`).Exec()
	if err != nil {
		return nil, errwrap.Wrapf("error validating connection info: {{err}}", err)
	}

	return session, nil
}
