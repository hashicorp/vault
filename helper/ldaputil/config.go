package ldaputil

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/helper/tlsutil"
)

type ConfigEntry struct {
	Url           string `json:"url"`
	UserDN        string `json:"userdn"`
	GroupDN       string `json:"groupdn"`
	GroupFilter   string `json:"groupfilter"`
	GroupAttr     string `json:"groupattr"`
	UPNDomain     string `json:"upndomain"`
	UserAttr      string `json:"userattr"`
	Certificate   string `json:"certificate"`
	InsecureTLS   bool   `json:"insecure_tls"`
	StartTLS      bool   `json:"starttls"`
	BindDN        string `json:"binddn"`
	BindPassword  string `json:"bindpass"`
	DenyNullBind  bool   `json:"deny_null_bind"`
	DiscoverDN    bool   `json:"discoverdn"`
	TLSMinVersion string `json:"tls_min_version"`
	TLSMaxVersion string `json:"tls_max_version"`

	// This json tag deviates from snake case because there was a past issue
	// where the tag was being ignored, causing it to be jsonified as "CaseSensitiveNames".
	// To continue reading in users' previously stored values,
	// we chose to carry that forward.
	CaseSensitiveNames *bool `json:"CaseSensitiveNames,omitempty"`
}

func (c *ConfigEntry) Validate() error {
	if len(c.Url) == 0 {
		return errors.New("at least one url must be provided")
	}
	// Note: This logic is driven by the logic in GetUserBindDN.
	// If updating this, please also update the logic there.
	if !c.DiscoverDN && (c.BindDN == "" || c.BindPassword == "") && c.UPNDomain == "" && c.UserDN == "" {
		return errors.New("cannot derive UserBindDN")
	}
	tlsMinVersion, ok := tlsutil.TLSLookup[c.TLSMinVersion]
	if !ok {
		return errors.New("invalid 'tls_min_version' in config")
	}
	tlsMaxVersion, ok := tlsutil.TLSLookup[c.TLSMaxVersion]
	if !ok {
		return errors.New("invalid 'tls_max_version' in config")
	}
	if tlsMaxVersion < tlsMinVersion {
		return errors.New("'tls_max_version' must be greater than or equal to 'tls_min_version'")
	}
	if c.Certificate != "" {
		block, _ := pem.Decode([]byte(c.Certificate))
		if block == nil || block.Type != "CERTIFICATE" {
			return errors.New("failed to decode PEM block in the certificate")
		}
		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate %s", err.Error())
		}
	}
	return nil
}
