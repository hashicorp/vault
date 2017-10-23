package cassandra

import (
	"crypto/tls"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"

	"github.com/gocql/gocql"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
)

// cassandraConnectionProducer implements ConnectionProducer and provides an
// interface for cassandra databases to make connections.
type cassandraConnectionProducer struct {
	Hosts             string      `json:"hosts" structs:"hosts" mapstructure:"hosts"`
	Port              int         `json:"port" structs:"port" mapstructure:"port"`
	Username          string      `json:"username" structs:"username" mapstructure:"username"`
	Password          string      `json:"password" structs:"password" mapstructure:"password"`
	TLS               bool        `json:"tls" structs:"tls" mapstructure:"tls"`
	InsecureTLS       bool        `json:"insecure_tls" structs:"insecure_tls" mapstructure:"insecure_tls"`
	ProtocolVersion   int         `json:"protocol_version" structs:"protocol_version" mapstructure:"protocol_version"`
	ConnectTimeoutRaw interface{} `json:"connect_timeout" structs:"connect_timeout" mapstructure:"connect_timeout"`
	TLSMinVersion     string      `json:"tls_min_version" structs:"tls_min_version" mapstructure:"tls_min_version"`
	Consistency       string      `json:"consistency" structs:"consistency" mapstructure:"consistency"`
	PemBundle         string      `json:"pem_bundle" structs:"pem_bundle" mapstructure:"pem_bundle"`
	PemJSON           string      `json:"pem_json" structs:"pem_json" mapstructure:"pem_json"`

	connectTimeout time.Duration
	certificate    string
	privateKey     string
	issuingCA      string

	Initialized bool
	Type        string
	session     *gocql.Session
	sync.Mutex
}

func (c *cassandraConnectionProducer) Initialize(conf map[string]interface{}, verifyConnection bool) error {
	c.Lock()
	defer c.Unlock()

	err := mapstructure.WeakDecode(conf, c)
	if err != nil {
		return err
	}

	if c.ConnectTimeoutRaw == nil {
		c.ConnectTimeoutRaw = "0s"
	}
	c.connectTimeout, err = parseutil.ParseDurationSecond(c.ConnectTimeoutRaw)
	if err != nil {
		return fmt.Errorf("invalid connect_timeout: %s", err)
	}

	switch {
	case len(c.Hosts) == 0:
		return fmt.Errorf("hosts cannot be empty")
	case len(c.Username) == 0:
		return fmt.Errorf("username cannot be empty")
	case len(c.Password) == 0:
		return fmt.Errorf("password cannot be empty")
	}

	var certBundle *certutil.CertBundle
	var parsedCertBundle *certutil.ParsedCertBundle
	switch {
	case len(c.PemJSON) != 0:
		parsedCertBundle, err = certutil.ParsePKIJSON([]byte(c.PemJSON))
		if err != nil {
			return fmt.Errorf("could not parse given JSON; it must be in the format of the output of the PKI backend certificate issuing command: %s", err)
		}
		certBundle, err = parsedCertBundle.ToCertBundle()
		if err != nil {
			return fmt.Errorf("Error marshaling PEM information: %s", err)
		}
		c.certificate = certBundle.Certificate
		c.privateKey = certBundle.PrivateKey
		c.issuingCA = certBundle.IssuingCA
		c.TLS = true

	case len(c.PemBundle) != 0:
		parsedCertBundle, err = certutil.ParsePEMBundle(c.PemBundle)
		if err != nil {
			return fmt.Errorf("Error parsing the given PEM information: %s", err)
		}
		certBundle, err = parsedCertBundle.ToCertBundle()
		if err != nil {
			return fmt.Errorf("Error marshaling PEM information: %s", err)
		}
		c.certificate = certBundle.Certificate
		c.privateKey = certBundle.PrivateKey
		c.issuingCA = certBundle.IssuingCA
		c.TLS = true
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(); err != nil {
			return fmt.Errorf("error verifying connection: %s", err)
		}
	}

	return nil
}

func (c *cassandraConnectionProducer) Connection() (interface{}, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, return it
	if c.session != nil {
		return c.session, nil
	}

	session, err := c.createSession()
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	c.session = session

	return session, nil
}

func (c *cassandraConnectionProducer) Close() error {
	// Grab the write lock
	c.Lock()
	defer c.Unlock()

	if c.session != nil {
		c.session.Close()
	}

	c.session = nil

	return nil
}

func (c *cassandraConnectionProducer) createSession() (*gocql.Session, error) {
	hosts := strings.Split(c.Hosts, ",")
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Authenticator = gocql.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	}

	if c.Port != 0 {
		clusterConfig.Port = c.Port
	}

	clusterConfig.ProtoVersion = c.ProtocolVersion
	if clusterConfig.ProtoVersion == 0 {
		clusterConfig.ProtoVersion = 2
	}

	clusterConfig.Timeout = c.connectTimeout
	if c.TLS {
		var tlsConfig *tls.Config
		if len(c.certificate) > 0 || len(c.issuingCA) > 0 {
			if len(c.certificate) > 0 && len(c.privateKey) == 0 {
				return nil, fmt.Errorf("found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(c.certificate) > 0 {
				certBundle.Certificate = c.certificate
				certBundle.PrivateKey = c.privateKey
			}
			if len(c.issuingCA) > 0 {
				certBundle.IssuingCA = c.issuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate bundle: %s", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, fmt.Errorf("failed to get TLS configuration: tlsConfig:%#v err:%v", tlsConfig, err)
			}
			tlsConfig.InsecureSkipVerify = c.InsecureTLS

			if c.TLSMinVersion != "" {
				var ok bool
				tlsConfig.MinVersion, ok = tlsutil.TLSLookup[c.TLSMinVersion]
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
		return nil, fmt.Errorf("error creating session: %s", err)
	}

	// Set consistency
	if c.Consistency != "" {
		consistencyValue, err := gocql.ParseConsistencyWrapper(c.Consistency)
		if err != nil {
			return nil, err
		}

		session.SetConsistency(consistencyValue)
	}

	// Verify the info
	err = session.Query(`LIST ALL`).Exec()
	if err != nil {
		return nil, fmt.Errorf("error validating connection info: %s", err)
	}

	return session, nil
}
