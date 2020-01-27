package client

import (
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"

	"github.com/hashicorp/vault/sdk/helper/certutil"
)

const (
	// These environment variables aren't set by default.
	// Vault may read them in if set through these environment variables.
	// Example here:
	// https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/
	EnvVarKubernetesNamespace = "VAULT_NAMESPACE"
	EnvVarKubernetesPodName   = "VAULT_POD_NAME"

	// The service host and port environment variables are
	// set by default inside a Kubernetes environment.
	EnvVarKubernetesServiceHost = "KUBERNETES_SERVICE_HOST"
	EnvVarKubernetesServicePort = "KUBERNETES_SERVICE_PORT"
)

var (
	// These are presented as variables so they can be updated
	// to point at test fixtures if needed. They aren't passed
	// into inClusterConfig to avoid dependency injection.
	Scheme     = "https://"
	TokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	RootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

// inClusterConfig returns a config object which uses the service account
// kubernetes gives to services. It's intended for clients that expect to be
// running inside a service running on kubernetes. It will return ErrNotInCluster
// if called from a process not running in a kubernetes environment.
func inClusterConfig() (*Config, error) {
	host, port := os.Getenv(EnvVarKubernetesServiceHost), os.Getenv(EnvVarKubernetesServicePort)
	if len(host) == 0 || len(port) == 0 {
		return nil, ErrNotInCluster
	}

	token, err := ioutil.ReadFile(TokenFile)
	if err != nil {
		return nil, err
	}

	pool, err := certutil.NewCertPool(RootCAFile)
	if err != nil {
		return nil, err
	}
	return &Config{
		Host:            Scheme + net.JoinHostPort(host, port),
		CACertPool:      pool,
		BearerToken:     string(token),
		BearerTokenFile: TokenFile,
	}, nil
}

type Config struct {
	Host            string
	BearerToken     string
	BearerTokenFile string
	CACertPool      *x509.CertPool
}
