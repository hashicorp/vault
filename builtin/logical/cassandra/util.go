package cassandra

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/logical"
)

// SplitSQL is used to split a series of SQL statements
func splitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		clean := strings.TrimSpace(p)
		if len(clean) > 0 {
			out = append(out, clean)
		}
	}
	return out
}

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

	if cfg.TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: cfg.InsecureTLS,
		}

		if len(cfg.Certificate) > 0 || len(cfg.IssuingCA) > 0 {
			if len(cfg.Certificate) > 0 && len(cfg.PrivateKey) == 0 {
				return nil, fmt.Errorf("Found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(cfg.Certificate) > 0 {
				certBundle.Certificate = cfg.Certificate
				certBundle.PrivateKey = cfg.PrivateKey
			}
			if len(cfg.IssuingCA) > 0 {
				certBundle.IssuingCA = cfg.IssuingCA
				tlsConfig.InsecureSkipVerify = false
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, fmt.Errorf("Error parsing certificate bundle: %s", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil {
				return nil, fmt.Errorf("Error getting TLS configuration: %s", err)
			}
		}

		clusterConfig.SslOpts = &gocql.SslOptions{
			Config: *tlsConfig,
		}
	}

	session, err := clusterConfig.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating session: %s", err)
	}

	// Verify the info
	err = session.Query(`LIST USERS`).Exec()
	if err != nil {
		return nil, fmt.Errorf("Error validating connection info: %s", err)
	}

	return session, nil
}
