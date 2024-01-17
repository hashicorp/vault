// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"bytes"
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
	// The client itself does nothing directly with these variables, it's
	// up to the caller. However, they live here so they'll be consistently
	// named should the client ever be reused.
	// We generally recommend preferring environmental settings over configured
	// ones, allowing settings from the Downward API to override hard-coded
	// ones.
	EnvVarKubernetesNamespace = "VAULT_K8S_NAMESPACE"
	EnvVarKubernetesPodName   = "VAULT_K8S_POD_NAME"

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
// inClusterConfig is based on this:
// https://github.com/kubernetes/client-go/blob/a56922badea0f2a91771411eaa1173c9e9243908/rest/config.go#L451
func inClusterConfig() (*Config, error) {
	host, port := os.Getenv(EnvVarKubernetesServiceHost), os.Getenv(EnvVarKubernetesServicePort)
	if len(host) == 0 || len(port) == 0 {
		return nil, ErrNotInCluster
	}

	token, err := ioutil.ReadFile(TokenFile)
	if err != nil {
		return nil, err
	}

	caBytes, err := ioutil.ReadFile(RootCAFile)
	if err != nil {
		return nil, err
	}
	pool, err := certutil.NewCertPool(bytes.NewReader(caBytes))
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

// This config is based on the one returned here:
// https://github.com/kubernetes/client-go/blob/a56922badea0f2a91771411eaa1173c9e9243908/rest/config.go#L451
// It is pared down to the absolute minimum fields used by this code.
// The CACertPool is promoted to the top level from being originally on the TLSClientConfig
// because it is the only parameter of the TLSClientConfig used by this code.
// Also, it made more sense to simply reuse the pool rather than holding raw values
// and parsing it repeatedly.
type Config struct {
	CACertPool *x509.CertPool

	// Host must be a host string, a host:port pair, or a URL to the base of the apiserver.
	// If a URL is given then the (optional) Path of that URL represents a prefix that must
	// be appended to all request URIs used to access the apiserver. This allows a frontend
	// proxy to easily relocate all of the apiserver endpoints.
	Host string

	// Server requires Bearer authentication. This client will not attempt to use
	// refresh tokens for an OAuth2 flow.
	BearerToken string

	// Path to a file containing a BearerToken.
	// If set, checks for a new token in the case of authorization errors.
	BearerTokenFile string
}
