package pkicli

import (
	"github.com/hashicorp/vault/api"
)

type Operations interface {
	// CreateRoot mounts a PKI backend and creates a root CA.
	// The parameterMap accepts:
	// - max_lease_ttl, to tune the backend
	// - all the parameters for https://www.vaultproject.io/api/secret/pki#generate-root
	// - all the parameters for https://www.vaultproject.io/api/secret/pki#set-urls
	CreateRoot(mountPath string, parameterMap mapStringAny) (*pkiCreateRootResponse, error)

	// CreateIntermediate mounts a PKI backend, and creates an intermediate CA signed by the CA at signingMountPath
	// The parameterMap accepts:
	// - max_lease_ttl, to tune the backend
	// - all the parameters for https://www.vaultproject.io/api/secret/pki#generate-intermediate
	// - all the parameters for https://www.vaultproject.io/api-docs/secret/pki#sign-intermediate
	// - all the parameters for https://www.vaultproject.io/api-docs/secret/pki#set-signed-intermediate
	CreateIntermediate(signingMountPath, mountPath string, parameterMap mapStringAny) (*pkiCreateIntermediateResponse, error)

}

type pkiCreateRootResponse struct {
	cert string
}

type pkiCreateIntermediateResponse struct {
	csr     string
	certPem string
}

var _ Operations = (*pkiOps)(nil)

type pkiOps struct {
	client *api.Client
}

func NewOperations(client *api.Client) Operations {
	return &pkiOps{client: client}
}

func (p pkiOps) CreateRoot(mountPath string, parameterMap mapStringAny) (*pkiCreateRootResponse, error) {

	params := newParams(parameterMap)
	params.put("_mount", mountPath)
	params.put("_ca_type", "root")
	params.putDefault("description", mountPath+" root CA")

	// 1. Enable the backend
	err := p.enableBackend(params)
	if err != nil {
		return nil, err
	}

	// 2. Generate the root CA
	generateResp, err := p.generateCA(params)
	if err != nil {
		return nil, err
	}

	// 3. Set the config URLs
	err = p.configUrls(params)
	if err != nil {
		return nil, err
	}

	return &pkiCreateRootResponse{
		cert: generateResp.cert,
	}, nil
}

func (p pkiOps) CreateIntermediate(rootMountPath, mountPath string, parameterMap mapStringAny) (*pkiCreateIntermediateResponse, error) {
	params := newParams(parameterMap)
	params.put("_mount", mountPath)

	// 1. Enable the backend
	enableParams := params.clone()
	enableParams.putDefault("description", mountPath+" intermediate CA")
	err := p.enableBackend(enableParams)
	if err != nil {
		return nil, err
	}

	// 2. Generate the intermediate CA
	generateParams := params.clone()
	generateParams.put("_ca_type", "intermediate")
	generateResp, err := p.generateCA(generateParams)
	if err != nil {
		return nil, err
	}

	// 3. Set the config URLs
	err = p.configUrls(params)
	if err != nil {
		return nil, err
	}

	if rootMountPath == "" {
		return &pkiCreateIntermediateResponse{
			csr: generateResp.csr,
		}, nil
	}

	// 4. Sign the intermediate CSR
	signParams := params.clone()
	signParams.put("_mount", rootMountPath)
	signParams.put("csr", generateResp.csr)
	signResp, err := p.signIntermediate(signParams)
	if err != nil {
		return nil, err
	}

	// 5. Set the signed certificate in the intermediate
	setSignedParams := params.clone()
	setSignedParams.put("certificate", signResp.certPem)
	err = p.setSigned(setSignedParams)
	if err != nil {
		return nil, err
	}

	return &pkiCreateIntermediateResponse{
		certPem: signResp.certPem,
	}, nil
}


