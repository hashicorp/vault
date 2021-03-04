package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"log"

	"github.com/Shopify/sarama"
)

func producerFromConfig(config map[string]string) (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	cfg.ClientID = "vault"
	cfg.Producer.Return.Successes = true

	tlsConfig, err := tlsConfig(config)
	if err != nil {
		return nil, err
	}

	tlsDisabled, tlsDisabledSet := config["tls_disabled"]

	if tlsDisabledSet && tlsDisabled == "true" {
		log.Println("[WARN] TLS disabled")
	} else {
		cfg.Net.TLS.Enable = true
		cfg.Net.TLS.Config = tlsConfig
	}

	return cfg, nil
}

func tlsConfig(config map[string]string) (*tls.Config, error) {
	tlsConfig := &tls.Config{}

	caCert, caOk := config["ca_cert"]
	clientCert, cCertOk := config["client_cert"]
	clientPrivateKey, cPKOk := config["client_private_key"]

	if cCertOk && cPKOk {
		clientCert, err := keyPair(clientCert, clientPrivateKey)
		if err != nil {
			return tlsConfig, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}

	if caOk {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(caCert))
		tlsConfig.RootCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}

	return tlsConfig, nil
}

func keyPair(cert, key string) (tls.Certificate, error) {
	return tls.X509KeyPair([]byte(cert), []byte(key))
}
