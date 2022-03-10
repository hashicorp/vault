package pkicli

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

func (p pkiOps) enableBackend(params *Params) error {
	//mountPath string, description string, maxLeaseTTL string) error {
	// https://www.vaultproject.io/api-docs/system/mounts#enable-secrets-engine

	params = params.clone()
	mount := params.pop("_mount")
	maxLeaseTTL := params.popDefault("max_lease_ttl", "")
	if err := params.error(); err != nil {
		return err
	}

	params.put("type", "pki")

	if maxLeaseTTL != "" {
		params.put("config", mapStringAny{
			"max_lease_ttl": maxLeaseTTL,
		})
	}
	_, err := p.write(params, "sys/mounts/%s", mount)

	return err
}

type generateCAResponse struct {
	secret *api.Secret
	cert string
	csr string
}

// Returns a secret with keys: certificate, expiration (number), issuing_ca, serial_number.
func (p pkiOps) generateCA(params *Params) (*generateCAResponse, error) {
	// https://www.vaultproject.io/api/secret/pki#generate-root
	// https://www.vaultproject.io/api/secret/pki#generate-intermediate

	params = params.clone()
	caType := params.pop("_ca_type")
	kind := params.popDefault("type", "internal")
	if err := params.error(); err != nil {
		return nil, err
	}

	// write to, for example, mypki/root/generate/internal
	secret, err := p.writeMount(params, "%s/generate/%s", caType, kind)
	if err != nil {
		return nil, err
	}

	var cert, csr string
	if v, ok := secret.Data["certificate"]; ok {
		cert = v.(string)
	}
	if v, ok := secret.Data["csr"]; ok {
		csr = v.(string)
	}
	return &generateCAResponse{
		secret: secret,
		cert: cert,
		csr: csr,
		// expiration json.Number
		// issuing_ca string
		// serial_number string
	}, nil
}

type signIntermediateResponse struct {
	secret *api.Secret
	certPem string
}

func (p pkiOps) signIntermediate(params *Params) (*signIntermediateResponse, error) {
	// https://www.vaultproject.io/api-docs/secret/pki#sign-intermediate

	params = params.clone()
	params.require("common_name", "csr")
	params.put("format", "pem_bundle")
	if err := params.error(); err != nil {
		return nil, err
	}

	secret, err := p.writeMount(params, "root/sign-intermediate")
	if err != nil {
		return nil, err
	}

	return &signIntermediateResponse{
		secret: secret,
		certPem: secret.Data["certificate"].(string),
	}, nil
}

func (p pkiOps) setSigned(params *Params) error {
	params = params.clone()
	params.require("certificate")
	if err := params.error(); err != nil {
		return err
	}

	_, err := p.writeMount(params, "intermediate/set-signed")

	return err
}

func (p pkiOps) configUrls(params *Params) error {
	// https://www.vaultproject.io/api/secret/pki#set-urls

	if !params.hasAny("issuing_certificates", "crl_distribution_points", "ocsp_servers") {
		return nil
	}
	params = params.clone()
	_, err := p.writeMount(params, "config/urls")

	return err
}

func (p pkiOps) readCAChainPEM(mountPath string) (string, error) {
	params := newParams(mapStringAny{"_mount": mountPath})
	s, err := p.readMount(params, "cert/ca_chain")

	return s.Data["ca_chain"].(string), err
}

func (p pkiOps) writeMount(params *Params, format string, a ...interface{}) (*api.Secret, error) {
	writePath, err := mountPath(params, format, a...)
	if err != nil {
		return nil, err
	}

	return p.write(params, writePath)
}

func mountPath(params *Params, format string, a ...interface{}) (string, error) {
	mount := params.pop("_mount")
	if err := params.error(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", mount, fmt.Sprintf(format, a...)), nil
}

func (p pkiOps) write(params *Params, format string, a ...interface{}) (*api.Secret, error) {
	writePath := fmt.Sprintf(format, a...)
	return p.client.Logical().Write(writePath, params.data())
}

func (p pkiOps) readMount(params *Params, format string, a ...interface{}) (*api.Secret, error) {
	readPath, err := mountPath(params, format, a...)
	if err != nil {
		return nil, err
	}

	return p.client.Logical().Read(readPath)
}