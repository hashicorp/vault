// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubesecrets

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	configPath        = "config"
	localCACertPath   = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	localJWTPath      = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	k8sServiceHostEnv = "KUBERNETES_SERVICE_HOST"
	k8sServicePortEnv = "KUBERNETES_SERVICE_PORT_HTTPS"
)

// kubeConfig contains the public key certificate used to verify the signature
// on the service account JWTs
type kubeConfig struct {
	// Host is the url string for the kubernetes API
	Host string `json:"kubernetes_host"`

	// CACert is the CA Cert to use to call into the kubernetes API
	CACert string `json:"kubernetes_ca_cert"`

	// ServiceAccountJwt is the bearer token to use when authenticating to the
	// kubernetes API
	ServiceAccountJwt string `json:"service_account_jwt"`

	// DisableLocalJWT is an optional parameter to disable defaulting to using
	// the local CA cert and service account jwt when running in a Kubernetes
	// pod
	DisableLocalCAJwt bool `json:"disable_local_ca_jwt"`
}

func (b *backend) pathConfig() *framework.Path {
	return &framework.Path{
		Pattern: configPath,
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKubernetes,
		},
		Fields: map[string]*framework.FieldSchema{
			"disable_local_ca_jwt": {
				Type:        framework.TypeBool,
				Description: "Disable defaulting to the local CA certificate and service account JWT when running in a Kubernetes pod.",
				Default:     false,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Disable use of local CA and service account JWT",
				},
			},
			"kubernetes_ca_cert": {
				Type:        framework.TypeString,
				Description: "PEM encoded CA certificate to use to verify the Kubernetes API server certificate. Defaults to the local pod's CA if found.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Kubernetes CA Certificate",
				},
			},
			"kubernetes_host": {
				Type:        framework.TypeString,
				Description: "Kubernetes API URL to connect to. Defaults to https://$KUBERNETES_SERVICE_HOST:KUBERNETES_SERVICE_PORT if those environment variables are set.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Kubernetes API URL",
				},
			},
			"service_account_jwt": {
				Type:        framework.TypeString,
				Description: "The JSON web token of the service account used by the secret engine to manage Kubernetes credentials. Defaults to the local pod's JWT if found.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Kubernetes API JWT",
				},
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "configure",
				},
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "configuration",
				},
			},
		},
		HelpSynopsis: "Configure the Kubernetes secret engine plugin.",
		HelpDescription: "This path configures the Kubernetes secret engine plugin. See the documentation for the " +
			"plugin specified for a full list of accepted connection details.",
	}
}

// pathConfigWrite handles create and update commands to the config
func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if config, err := getConfig(ctx, req.Storage); err != nil {
		return nil, err
	} else if config == nil {
		return nil, nil
	} else {
		// Create a map of data to be returned. Note that these reflect just the
		// values that the user set, not what the defaults will be if they
		// aren't set (see configWithDynamicValues() for those defaults). And
		// the service account jwt is omitted as sensitive data.
		resp := &logical.Response{
			Data: map[string]interface{}{
				"disable_local_ca_jwt": config.DisableLocalCAJwt,
				"kubernetes_ca_cert":   config.CACert,
				"kubernetes_host":      config.Host,
			},
		}

		return resp, nil
	}
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &kubeConfig{}
	}

	if host, ok := data.GetOk("kubernetes_host"); ok {
		config.Host = host.(string)
	} else if _, err := getK8sURLFromEnv(); err != nil {
		return nil, errors.New("kubernetes_host was unset and could not be determined from environment variables")
	}
	if disableLocalJWT, ok := data.GetOk("disable_local_ca_jwt"); ok {
		config.DisableLocalCAJwt = disableLocalJWT.(bool)
	}
	if caCert, ok := data.GetOk("kubernetes_ca_cert"); ok {
		config.CACert = caCert.(string)
	}
	if serviceAccountJWT, ok := data.GetOk("service_account_jwt"); ok {
		config.ServiceAccountJwt = serviceAccountJWT.(string)
	}

	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	// reset the client so the next invocation will pick up the new configuration
	b.reset()

	return nil, nil
}

func (b *backend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, configPath)

	if err == nil {
		b.reset()
	}

	return nil, err
}

// configWithDynamicValues fetches the kubeConfig from storage and sets any
// runtime defaults for host, local token, and local CA certificate.
func (b *backend) configWithDynamicValues(ctx context.Context, s logical.Storage) (*kubeConfig, error) {
	config, err := getConfig(ctx, s)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("could not load backend configuration")
	}

	// If host is blank, default to reading from env
	if config.Host == "" {
		config.Host, err = getK8sURLFromEnv()
		if err != nil {
			return nil, errors.New("kubernetes_host was unset and could not determine it from environment variables")
		}
	}

	// Nothing more to do if loading local CA cert and JWT token is disabled.
	if config.DisableLocalCAJwt {
		return config, nil
	}

	// Read local JWT token unless it was not stored in config.
	if config.ServiceAccountJwt == "" {
		jwtBytes, err := b.localSATokenReader.ReadFile()
		if err != nil {
			// Ignore error: make best effort trying to load local JWT,
			// otherwise the JWT submitted in login payload will be used.
			b.Logger().Debug("failed to read local service account token, will use client token", "error", err)
		}
		config.ServiceAccountJwt = string(jwtBytes)
	}

	// Read local CA cert unless it was stored in config.
	if config.CACert == "" {
		caBytes, err := b.localCACertReader.ReadFile()
		if err != nil {
			return nil, err
		}
		config.CACert = string(caBytes)
	}

	return config, nil
}

func getConfig(ctx context.Context, s logical.Storage) (*kubeConfig, error) {
	entry, err := s.Get(ctx, configPath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := new(kubeConfig)
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	// return the config, we are done
	return config, nil
}

func getK8sURLFromEnv() (string, error) {
	host := os.Getenv(k8sServiceHostEnv)
	port := os.Getenv(k8sServicePortEnv)
	if host == "" || port == "" {
		return "", fmt.Errorf("failed to find k8s API host variables %q and %q in env", k8sServiceHostEnv, k8sServicePortEnv)
	}
	return fmt.Sprintf("https://%s:%s", host, port), nil
}
