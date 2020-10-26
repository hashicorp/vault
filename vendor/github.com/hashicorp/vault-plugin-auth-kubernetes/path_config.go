package kubeauth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/briankassouf/jose/jws"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

var (
	localCACertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	localJWTPath    = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

// pathConfig returns the path configuration for CRUD operations on the backend
// configuration.
func pathConfig(b *kubeAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "config$",
		Fields: map[string]*framework.FieldSchema{
			"kubernetes_host": {
				Type:        framework.TypeString,
				Description: "Host must be a host string, a host:port pair, or a URL to the base of the Kubernetes API server.",
			},

			"kubernetes_ca_cert": {
				Type:        framework.TypeString,
				Description: "PEM encoded CA cert for use by the TLS client used to talk with the API.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Kubernetes CA Certificate",
				},
			},
			"token_reviewer_jwt": {
				Type: framework.TypeString,
				Description: `A service account JWT used to access the
TokenReview API to validate other JWTs during login. If not set
the JWT used for login will be used to access the API.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Token Reviewer JWT",
				},
			},
			"pem_keys": {
				Type: framework.TypeCommaStringSlice,
				Description: `Optional list of PEM-formated public keys or certificates
used to verify the signatures of kubernetes service account
JWTs. If a certificate is given, its public key will be
extracted. Not every installation of Kuberentes exposes these keys.`,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Service account verification keys",
				},
			},
			"issuer": {
				Type:        framework.TypeString,
				Description: "Optional JWT issuer. If no issuer is specified, then this plugin will use kubernetes.io/serviceaccount as the default issuer.",
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "JWT Issuer",
				},
			},
			"disable_iss_validation": {
				Type:        framework.TypeBool,
				Description: "Disable JWT issuer validation. Allows to skip ISS validation.",
				Default:     false,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Disable JWT Issuer Validation",
				},
			},
			"disable_local_ca_jwt": {
				Type:        framework.TypeBool,
				Description: "Disable defaulting to the local CA cert and service account JWT when running in a Kubernetes pod",
				Default:     false,
				DisplayAttrs: &framework.DisplayAttributes{
					Name: "Disable use of local CA and service account JWT",
				},
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathConfigWrite,
			logical.CreateOperation: b.pathConfigWrite,
			logical.ReadOperation:   b.pathConfigRead,
		},

		HelpSynopsis:    confHelpSyn,
		HelpDescription: confHelpDesc,
	}
}

// pathConfigWrite handles create and update commands to the config
func (b *kubeAuthBackend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if config, err := b.config(ctx, req.Storage); err != nil {
		return nil, err
	} else if config == nil {
		return nil, nil
	} else {
		// Create a map of data to be returned
		resp := &logical.Response{
			Data: map[string]interface{}{
				"kubernetes_host":        config.Host,
				"kubernetes_ca_cert":     config.CACert,
				"pem_keys":               config.PEMKeys,
				"issuer":                 config.Issuer,
				"disable_iss_validation": config.DisableISSValidation,
				"disable_local_ca_jwt":   config.DisableLocalCAJwt,
			},
		}

		return resp, nil
	}
}

// pathConfigWrite handles create and update commands to the config
func (b *kubeAuthBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	host := data.Get("kubernetes_host").(string)
	if host == "" {
		return logical.ErrorResponse("no host provided"), nil
	}

	disableLocalJWT := data.Get("disable_local_ca_jwt").(bool)
	localCACert := []byte{}
	localTokenReviewer := []byte{}
	if !disableLocalJWT {
		localCACert, _ = ioutil.ReadFile(localCACertPath)
		localTokenReviewer, _ = ioutil.ReadFile(localJWTPath)
	}
	pemList := data.Get("pem_keys").([]string)
	caCert := data.Get("kubernetes_ca_cert").(string)
	issuer := data.Get("issuer").(string)
	disableIssValidation := data.Get("disable_iss_validation").(bool)
	if len(pemList) == 0 && len(caCert) == 0 {
		if len(localCACert) > 0 {
			caCert = string(localCACert)
		} else {
			return logical.ErrorResponse("one of pem_keys or kubernetes_ca_cert must be set"), nil
		}
	}

	tokenReviewer := data.Get("token_reviewer_jwt").(string)
	if !disableLocalJWT && len(tokenReviewer) == 0 && len(localTokenReviewer) > 0 {
		tokenReviewer = string(localTokenReviewer)
	}

	if len(tokenReviewer) > 0 {
		// Validate it's a JWT
		_, err := jws.ParseJWT([]byte(tokenReviewer))
		if err != nil {
			return nil, err
		}
	}

	config := &kubeConfig{
		PublicKeys:           make([]interface{}, len(pemList)),
		PEMKeys:              pemList,
		Host:                 host,
		CACert:               caCert,
		TokenReviewerJWT:     tokenReviewer,
		Issuer:               issuer,
		DisableISSValidation: disableIssValidation,
		DisableLocalCAJwt:    disableLocalJWT,
	}

	var err error
	for i, pem := range pemList {
		config.PublicKeys[i], err = parsePublicKeyPEM([]byte(pem))
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	entry, err := logical.StorageEntryJSON(configPath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}
	return nil, nil
}

// kubeConfig contains the public key certificate used to verify the signature
// on the service account JWTs
type kubeConfig struct {
	// PublicKeys is the list of public key objects used to verify JWTs
	PublicKeys []interface{} `json:"-"`
	// PEMKeys is the list of public key PEMs used to store the keys
	// in storage.
	PEMKeys []string `json:"pem_keys"`
	// Host is the url string for the kubernetes API
	Host string `json:"host"`
	// CACert is the CA Cert to use to call into the kubernetes API
	CACert string `json:"ca_cert"`
	// TokenReviewJWT is the bearer to use during the TokenReview API call
	TokenReviewerJWT string `json:"token_reviewer_jwt"`
	// Issuer is the claim that specifies who issued the token
	Issuer string `json:"issuer"`
	// DisableISSValidation is optional parameter to allow to skip ISS validation
	DisableISSValidation bool `json:"disable_iss_validation"`
	// DisableLocalJWT is an optional parameter to disable defaulting to using
	// the local CA cert and service account jwt when running in a Kubernetes
	// pod
	DisableLocalCAJwt bool `json:"disable_local_ca_jwt"`
}

// PasrsePublicKeyPEM is used to parse RSA and ECDSA public keys from PEMs
func parsePublicKeyPEM(data []byte) (interface{}, error) {
	block, data := pem.Decode(data)
	if block != nil {
		var rawKey interface{}
		var err error
		if rawKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
			if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
				rawKey = cert.PublicKey
			} else {
				return nil, err
			}
		}

		if rsaPublicKey, ok := rawKey.(*rsa.PublicKey); ok {
			return rsaPublicKey, nil
		}
		if ecPublicKey, ok := rawKey.(*ecdsa.PublicKey); ok {
			return ecPublicKey, nil
		}
	}

	return nil, errors.New("data does not contain any valid RSA or ECDSA public keys")
}

const confHelpSyn = `Configures the JWT Public Key and Kubernetes API information.`
const confHelpDesc = `
The Kubernetes Auth backend validates service account JWTs and verifies their
existence with the Kubernetes TokenReview API. This endpoint configures the
public key used to validate the JWT signature and the necessary information to
access the Kubernetes API.
`
