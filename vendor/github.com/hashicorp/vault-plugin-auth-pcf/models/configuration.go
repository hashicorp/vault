package models

import (
	"crypto/x509"
	"errors"
	"fmt"
)

// NewConfiguration is the way a Configuration is intended to be obtained. It ensures the
// given certificates are valid and prepares a CA certificate pool to be used for client
// certificate verification.
func NewConfiguration(certificates []string, pcfAPIAddr, pcfUsername, pcfPassword string) (*Configuration, error) {
	config := &Configuration{
		Certificates: certificates,
		PCFAPIAddr:   pcfAPIAddr,
		PCFUsername:  pcfUsername,
		PCFPassword:  pcfPassword,
	}
	pool := x509.NewCertPool()
	for _, certificate := range certificates {
		if ok := pool.AppendCertsFromPEM([]byte(certificate)); !ok {
			return nil, fmt.Errorf("couldn't append CA certificate: %s", certificate)
		}
	}
	config.verifyOpts = &x509.VerifyOptions{Roots: pool}
	return config, nil
}

// Configuration is not intended to by directly instantiated; please use NewConfiguration.
type Configuration struct {
	// Certificates are the CA certificates that should be used for verifying client certificates.
	Certificates []string `json:"certificates"`

	// PCFAPIAddr is the address of PCF's API, ex: "https://api.dev.cfdev.sh" or "http://127.0.0.1:33671"
	PCFAPIAddr string `json:"pcf_api_addr"`

	// The username for the PCF API.
	PCFUsername string `json:"pcf_username"`

	// The password for the PCF API.
	PCFPassword string `json:"pcf_password"`

	// verifyOpts is intentionally lower-cased so it won't be stored in JSON.
	// Instead, this struct is expected to be created from NewConfiguration
	// so that it'll populate this field.
	verifyOpts *x509.VerifyOptions
}

// VerifyOpts returns the options that can be used for verifying client certificates,
// including the CA certificate pool.
func (c *Configuration) VerifyOpts() (x509.VerifyOptions, error) {
	if c.verifyOpts == nil {
		return x509.VerifyOptions{}, errors.New("verify options are unset")
	}
	return *c.verifyOpts, nil
}
