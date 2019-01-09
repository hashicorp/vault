package influxdb

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/plugins/helper/database/connutil"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/mitchellh/mapstructure"
)

// influxdbConnectionProducer implements ConnectionProducer and provides an
// interface for influxdb databases to make connections.
type influxdbConnectionProducer struct {
	Host              string      `json:"host" structs:"host" mapstructure:"host"`
	Username          string      `json:"username" structs:"username" mapstructure:"username"`
	Password          string      `json:"password" structs:"password" mapstructure:"password"`
	Port              string      `json:"port" structs:"port" mapstructure:"port"` //default to 8086
	TLS               bool        `json:"tls" structs:"tls" mapstructure:"tls"`
	InsecureTLS       bool        `json:"insecure_tls" structs:"insecure_tls" mapstructure:"insecure_tls"`
	ConnectTimeoutRaw interface{} `json:"connect_timeout" structs:"connect_timeout" mapstructure:"connect_timeout"`
	TLSMinVersion     string      `json:"tls_min_version" structs:"tls_min_version" mapstructure:"tls_min_version"`
	PemBundle         string      `json:"pem_bundle" structs:"pem_bundle" mapstructure:"pem_bundle"`
	PemJSON           string      `json:"pem_json" structs:"pem_json" mapstructure:"pem_json"`

	connectTimeout time.Duration
	certificate    string
	privateKey     string
	issuingCA      string
	rawConfig      map[string]interface{}

	Initialized bool
	Type        string
	client      influx.Client
	sync.Mutex
}

func (i *influxdbConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := i.Init(ctx, conf, verifyConnection)
	return err
}

func (i *influxdbConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	i.Lock()
	defer i.Unlock()

	i.rawConfig = conf

	err := mapstructure.WeakDecode(conf, i)
	if err != nil {
		return nil, err
	}

	if i.ConnectTimeoutRaw == nil {
		i.ConnectTimeoutRaw = "0s"
	}
	if i.Port == "" {
		i.Port = "8086"
	}
	i.connectTimeout, err = parseutil.ParseDurationSecond(i.ConnectTimeoutRaw)
	if err != nil {
		return nil, errwrap.Wrapf("invalid connect_timeout: {{err}}", err)
	}

	switch {
	case len(i.Host) == 0:
		return nil, fmt.Errorf("host cannot be empty")
	case len(i.Username) == 0:
		return nil, fmt.Errorf("username cannot be empty")
	case len(i.Password) == 0:
		return nil, fmt.Errorf("password cannot be empty")
	}

	var certBundle *certutil.CertBundle
	var parsedCertBundle *certutil.ParsedCertBundle
	switch {
	case len(i.PemJSON) != 0:
		parsedCertBundle, err = certutil.ParsePKIJSON([]byte(i.PemJSON))
		if err != nil {
			return nil, errwrap.Wrapf("could not parse given JSON; it must be in the format of the output of the PKI backend certificate issuing command: {{err}}", err)
		}
		certBundle, err = parsedCertBundle.ToCertBundle()
		if err != nil {
			return nil, errwrap.Wrapf("Error marshaling PEM information: {{err}}", err)
		}
		i.certificate = certBundle.Certificate
		i.privateKey = certBundle.PrivateKey
		i.issuingCA = certBundle.IssuingCA
		i.TLS = true

	case len(i.PemBundle) != 0:
		parsedCertBundle, err = certutil.ParsePEMBundle(i.PemBundle)
		if err != nil {
			return nil, errwrap.Wrapf("Error parsing the given PEM information: {{err}}", err)
		}
		certBundle, err = parsedCertBundle.ToCertBundle()
		if err != nil {
			return nil, errwrap.Wrapf("Error marshaling PEM information: {{err}}", err)
		}
		i.certificate = certBundle.Certificate
		i.privateKey = certBundle.PrivateKey
		i.issuingCA = certBundle.IssuingCA
		i.TLS = true
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	i.Initialized = true

	if verifyConnection {
		if _, err := i.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}
	}

	return conf, nil
}

func (i *influxdbConnectionProducer) Connection(_ context.Context) (interface{}, error) {
	if !i.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	// If we already have a DB, return it
	if i.client != nil {
		return i.client, nil
	}

	cli, err := i.createClient()
	if err != nil {
		return nil, err
	}

	//  Store the session in backend for reuse
	i.client = cli

	return cli, nil
}

func (i *influxdbConnectionProducer) Close() error {
	// Grab the write lock
	i.Lock()
	defer i.Unlock()

	if i.client != nil {
		i.client.Close()
	}

	i.client = nil

	return nil
}

func (i *influxdbConnectionProducer) createClient() (influx.Client, error) {
	clientConfig := influx.HTTPConfig{
		Addr:      fmt.Sprintf("http://%s:%s", i.Host, i.Port),
		Username:  i.Username,
		Password:  i.Password,
		UserAgent: "vault-influxdb-plugin",
		Timeout:   i.connectTimeout,
	}

	if i.TLS {
		var tlsConfig *tls.Config
		if len(i.certificate) > 0 || len(i.issuingCA) > 0 {
			if len(i.certificate) > 0 && len(i.privateKey) == 0 {
				return nil, fmt.Errorf("found certificate for TLS authentication but no private key")
			}

			certBundle := &certutil.CertBundle{}
			if len(i.certificate) > 0 {
				certBundle.Certificate = i.certificate
				certBundle.PrivateKey = i.privateKey
			}
			if len(i.issuingCA) > 0 {
				certBundle.IssuingCA = i.issuingCA
			}

			parsedCertBundle, err := certBundle.ToParsedCertBundle()
			if err != nil {
				return nil, errwrap.Wrapf("failed to parse certificate bundle: {{err}}", err)
			}

			tlsConfig, err = parsedCertBundle.GetTLSConfig(certutil.TLSClient)
			if err != nil || tlsConfig == nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("failed to get TLS configuration: tlsConfig:%#v err:{{err}}", tlsConfig), err)
			}
			tlsConfig.InsecureSkipVerify = i.InsecureTLS

			if i.TLSMinVersion != "" {
				var ok bool
				tlsConfig.MinVersion, ok = tlsutil.TLSLookup[i.TLSMinVersion]
				if !ok {
					return nil, fmt.Errorf("invalid 'tls_min_version' in config")
				}
			} else {
				// MinVersion was not being set earlier. Reset it to
				// zero to gracefully handle upgrades.
				tlsConfig.MinVersion = 0
			}
		}
		clientConfig.TLSConfig = tlsConfig
		clientConfig.Addr = fmt.Sprintf("https://%s:%s", i.Host, i.Port)
	}

	cli, err := influx.NewHTTPClient(clientConfig)
	if err != nil {
		return nil, errwrap.Wrapf("error creating client: {{err}}", err)
	}

	// Checking server status
	_, _, err = cli.Ping(i.connectTimeout)
	if err != nil {
		return nil, errwrap.Wrapf("error checking cluster status: {{err}}", err)
	}

	// verifying infos about the connection
	isAdmin, err := isUserAdmin(cli, i.Username)
	if err != nil {
		return nil, errwrap.Wrapf("error getting if provided username is admin: {{err}}", err)
	}
	if !isAdmin {
		return nil, fmt.Errorf("the provided user is not an admin of the influxDB server")
	}

	return cli, nil
}

func (i *influxdbConnectionProducer) secretValues() map[string]interface{} {
	return map[string]interface{}{
		i.Password:  "[password]",
		i.PemBundle: "[pem_bundle]",
		i.PemJSON:   "[pem_json]",
	}
}

func isUserAdmin(cli influx.Client, user string) (bool, error) {
	q := influx.NewQuery("SHOW USERS", "", "")
	response, err := cli.Query(q)
	if err != nil {
		return false, err
	}
	if response.Error() != nil {
		return false, response.Error()
	}
	for _, res := range response.Results {
		for _, serie := range res.Series {
			for _, val := range serie.Values {
				if val[0].(string) == user && val[1].(bool) == true {
					return true, nil
				}
			}
		}
	}
	return false, fmt.Errorf("the provided username is not a valid user in the influxdb")
}
