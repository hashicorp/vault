package gocbcore

import (
	"crypto/tls"
	"crypto/x509"
	"net"
)

type dynTLSConfig struct {
	BaseConfig *tls.Config
	Provider   func() *x509.CertPool
}

func (config dynTLSConfig) Clone() *dynTLSConfig {
	return &dynTLSConfig{
		BaseConfig: config.BaseConfig.Clone(),
		Provider:   config.Provider,
	}
}

func (config dynTLSConfig) MakeForHost(serverName string) (*tls.Config, error) {
	newConfig := config.BaseConfig.Clone()

	if config.Provider != nil {
		rootCAs := config.Provider()
		if rootCAs != nil {
			newConfig.RootCAs = rootCAs
			newConfig.InsecureSkipVerify = false
		} else {
			newConfig.RootCAs = nil
			newConfig.InsecureSkipVerify = true
		}
	}

	newConfig.ServerName = serverName
	return newConfig, nil
}

func (config dynTLSConfig) MakeForAddr(addr string) (*tls.Config, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	return config.MakeForHost(host)
}
