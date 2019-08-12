package util

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault-plugin-auth-pcf/models"
)

const BashTimeFormat = "Mon Jan 2 15:04:05 MST 2006"

// NewPCFClient does some work that's needed every time we use the PCF client,
// namely using cleanhttp and configuring it to match the user conf.
func NewPCFClient(config *models.Configuration) (*cfclient.Client, error) {
	clientConf := &cfclient.Config{
		ApiAddress: config.PCFAPIAddr,
		Username:   config.PCFUsername,
		Password:   config.PCFPassword,
		HttpClient: cleanhttp.DefaultClient(),
	}
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	for _, certificate := range config.PCFAPICertificates {
		if ok := rootCAs.AppendCertsFromPEM([]byte(certificate)); !ok {
			return nil, fmt.Errorf("couldn't append PCF API cert to trust: %s", certificate)
		}
	}
	tlsConfig := &tls.Config{
		RootCAs: rootCAs,
	}
	clientConf.HttpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	return cfclient.NewClient(clientConf)
}
