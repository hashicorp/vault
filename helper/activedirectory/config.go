package activedirectory

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	DefaultTLSMinVersion = "tls12"
	DefaultTLSMaxVersion = "tls12"
)

func NewConfiguration(logger hclog.Logger, fieldData *framework.FieldData) (*Configuration, error) {

	certificate, err := getValidatedCertificate(fieldData)
	if err != nil {
		return nil, err
	}

	tlsMinVersion, err := getTLSMinVersion(fieldData)
	if err != nil {
		return nil, err
	}

	tlsMaxVersion, err := getTLSMaxVersion(fieldData)
	if err != nil {
		return nil, err
	}

	conf := &Configuration{
		Certificate:   certificate,
		InsecureTLS:   fieldData.Get("insecure_tls").(bool),
		Password:      fieldData.Get("password").(string),
		StartTLS:      getStartTLS(fieldData),
		TLSMinVersion: tlsMinVersion,
		TLSMaxVersion: tlsMaxVersion,
		URL:           fieldData.Get("url").(string),
		Username:      fieldData.Get("username").(string),
		logger:        logger,
	}

	if err := conf.validate(); err != nil {
		return nil, err
	}

	return conf, nil
}

type Configuration struct {
	Certificate   string `json:"certificate" structs:"certificate" mapstructure:"certificate"`
	InsecureTLS   bool   `json:"insecure_tls" structs:"insecuretls" mapstructure:"insecuretls"`
	Password      string `json:"password" structs:"password" mapstructure:"password"`
	StartTLS      bool   `json:"starttls" structs:"starttls" mapstructure:"starttls"`
	TLSMinVersion uint16 `json:"tlsminversion" structs:"tlsminversion" mapstructure:"tlsminversion"`
	TLSMaxVersion uint16 `json:"tlsmaxversion" structs:"tlsmaxversion" mapstructure:"tlsmaxversion"`
	URL           string `json:"url" structs:"url" mapstructure:"certificate"`
	Username      string `json:"username" structs:"username" mapstructure:"username"`

	// *tlsConfig objects aren't jsonable, so we must avoid storing them and instead generate them on the fly
	tlsConfigs map[*url.URL]*tls.Config
	logger     hclog.Logger
}

func (c *Configuration) validate() error {
	if c.TLSMinVersion < c.TLSMaxVersion {
		return fmt.Errorf("'tls_max_version' must be greater than or equal to 'tls_min_version'")
	}
	return nil
}

func (c *Configuration) GetTLSConfigs() (map[*url.URL]*tls.Config, error) {
	if len(c.tlsConfigs) <= 0 {
		configs, err := c.getTLSConfigs()
		if err != nil {
			return nil, err
		}
		c.tlsConfigs = configs
	}
	return c.tlsConfigs, nil
}

func (c *Configuration) getTLSConfigs() (map[*url.URL]*tls.Config, error) {

	confUrls := strings.ToLower(c.URL)
	urls := strings.Split(confUrls, ",")

	tlsConfigs := make(map[*url.URL]*tls.Config)
	for _, uut := range urls {

		u, err := url.Parse(uut)
		if err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn(fmt.Sprintf("unable to parse %s: %s, ignoring", uut, err.Error()))
			}
			continue
		}

		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			// err intentionally ignored
			// fall back to using the parsed url's host
			host = u.Host
		}

		tlsConfig := &tls.Config{
			ServerName:         host,
			MinVersion:         c.TLSMinVersion,
			MaxVersion:         c.TLSMaxVersion,
			InsecureSkipVerify: c.InsecureTLS,
		}

		if c.Certificate != "" {
			caPool := x509.NewCertPool()
			ok := caPool.AppendCertsFromPEM([]byte(c.Certificate))
			if !ok {
				// this probably won't succeed on further attempts, so return
				return nil, fmt.Errorf("could not append CA certificate")
			}
			tlsConfig.RootCAs = caPool
		}

		tlsConfigs[u] = tlsConfig
	}

	return tlsConfigs, nil
}

func getValidatedCertificate(fieldData *framework.FieldData) (string, error) {

	confCertificate := fieldData.Get("certificate").(string)
	if confCertificate == "" {
		// no certificate was provided
		return "", nil
	}

	block, _ := pem.Decode([]byte(confCertificate))
	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode PEM block in the certificate")
	}

	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate %s", err.Error())
	}

	return confCertificate, nil
}

func getStartTLS(fieldData *framework.FieldData) bool {

	startTLSIfc, ok := fieldData.GetOk("starttls")
	if !ok {
		return true
	}

	confStartTLS, ok := startTLSIfc.(bool)
	if !ok {
		return true
	}

	return confStartTLS
}

func getTLSMinVersion(fieldData *framework.FieldData) (uint16, error) {

	confTLSMinVersion := fieldData.Get("tls_min_version").(string)
	if confTLSMinVersion == "" {
		confTLSMinVersion = DefaultTLSMinVersion
	}

	tlsMinVersion, ok := tlsutil.TLSLookup[confTLSMinVersion]
	if !ok {
		return 0, fmt.Errorf("invalid 'tls_min_version' in config")
	}

	return tlsMinVersion, nil
}

func getTLSMaxVersion(fieldData *framework.FieldData) (uint16, error) {

	confTLSMaxVersion := fieldData.Get("tls_max_version").(string)
	if confTLSMaxVersion == "" {
		confTLSMaxVersion = DefaultTLSMaxVersion
	}

	tlsMaxVersion, ok := tlsutil.TLSLookup[confTLSMaxVersion]
	if !ok {
		return 0, fmt.Errorf("invalid 'tls_max_version' in config")
	}

	return tlsMaxVersion, nil
}
