package healthcheck

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

func pkiFetchIssuers(e *Executor, versionError func()) (bool, *PathFetch, []string, error) {
	issuersRet, err := e.FetchIfNotFetched(logical.ListOperation, "/{{mount}}/issuers")
	if err != nil {
		return true, nil, nil, err
	}

	if !issuersRet.IsSecretOK() {
		if issuersRet.IsUnsupportedPathError() {
			versionError()
		}

		return true, nil, nil, nil
	}

	if len(issuersRet.ParsedCache) == 0 {
		var issuers []string
		for _, rawIssuerId := range issuersRet.Secret.Data["keys"].([]interface{}) {
			issuers = append(issuers, rawIssuerId.(string))
		}
		issuersRet.ParsedCache["issuers"] = issuers
	}

	return false, issuersRet, issuersRet.ParsedCache["issuers"].([]string), nil
}

func parsePEM(contents string) ([]byte, error) {
	// Need to parse out the issuer from its PEM format.
	pemBlock, _ := pem.Decode([]byte(contents))
	if pemBlock == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}

	return pemBlock.Bytes, nil
}

func parsePEMCert(contents string) (*x509.Certificate, error) {
	parsed, err := parsePEM(contents)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(parsed)
	if err != nil {
		return nil, fmt.Errorf("invalid certificate: %v", err)
	}

	return cert, nil
}

func parsePEMCRL(contents string) (*x509.RevocationList, error) {
	parsed, err := parsePEM(contents)
	if err != nil {
		return nil, err
	}

	crl, err := x509.ParseRevocationList(parsed)
	if err != nil {
		return nil, fmt.Errorf("invalid CRL: %v", err)
	}

	return crl, nil
}

func pkiFetchIssuer(e *Executor, issuer string, versionError func()) (bool, *PathFetch, *x509.Certificate, error) {
	issuerRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/issuer/"+issuer+"/json")
	if err != nil {
		return true, nil, nil, err
	}

	if !issuerRet.IsSecretOK() {
		if issuerRet.IsUnsupportedPathError() {
			versionError()
		}
		return true, nil, nil, nil
	}

	if len(issuerRet.ParsedCache) == 0 {
		cert, err := parsePEMCert(issuerRet.Secret.Data["certificate"].(string))
		if err != nil {
			return true, nil, nil, fmt.Errorf("unable to parse issuer %v's certificate: %v", issuer, err)
		}

		issuerRet.ParsedCache["certificate"] = cert
	}

	return false, issuerRet, issuerRet.ParsedCache["certificate"].(*x509.Certificate), nil
}

func pkiFetchIssuerCRL(e *Executor, issuer string, delta bool, versionError func()) (bool, *PathFetch, *x509.RevocationList, error) {
	path := "/{{mount}}/issuer/" + issuer + "/crl"
	name := "CRL"
	if delta {
		path += "/delta"
		name = "Delta CRL"
	}

	crlRet, err := e.FetchIfNotFetched(logical.ReadOperation, path)
	if err != nil {
		return true, nil, nil, err
	}

	if !crlRet.IsSecretOK() {
		if crlRet.IsUnsupportedPathError() {
			versionError()
		}
		return true, nil, nil, nil
	}

	if len(crlRet.ParsedCache) == 0 {
		crl, err := parsePEMCRL(crlRet.Secret.Data["crl"].(string))
		if err != nil {
			return true, nil, nil, fmt.Errorf("unable to parse issuer %v's %v: %v", issuer, name, err)
		}
		crlRet.ParsedCache["crl"] = crl
	}

	return false, crlRet, crlRet.ParsedCache["crl"].(*x509.RevocationList), nil
}
