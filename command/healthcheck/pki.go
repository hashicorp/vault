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
		return nil, fmt.Errorf("invalid certificate: %w", err)
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
		return nil, fmt.Errorf("invalid CRL: %w", err)
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
			return true, nil, nil, fmt.Errorf("unable to parse issuer %v's certificate: %w", issuer, err)
		}

		issuerRet.ParsedCache["certificate"] = cert
	}

	return false, issuerRet, issuerRet.ParsedCache["certificate"].(*x509.Certificate), nil
}

func pkiFetchIssuerEntry(e *Executor, issuer string, versionError func()) (bool, *PathFetch, map[string]interface{}, error) {
	issuerRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/issuer/"+issuer)
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
			return true, nil, nil, fmt.Errorf("unable to parse issuer %v's certificate: %w", issuer, err)
		}

		issuerRet.ParsedCache["certificate"] = cert
	}

	var data map[string]interface{} = nil
	if issuerRet.Secret != nil && len(issuerRet.Secret.Data) > 0 {
		data = issuerRet.Secret.Data
	}

	return false, issuerRet, data, nil
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
			return true, nil, nil, fmt.Errorf("unable to parse issuer %v's %v: %w", issuer, name, err)
		}
		crlRet.ParsedCache["crl"] = crl
	}

	return false, crlRet, crlRet.ParsedCache["crl"].(*x509.RevocationList), nil
}

func pkiFetchKeyEntry(e *Executor, key string, versionError func()) (bool, *PathFetch, map[string]interface{}, error) {
	keyRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/key/"+key)
	if err != nil {
		return true, nil, nil, err
	}

	if !keyRet.IsSecretOK() {
		if keyRet.IsUnsupportedPathError() {
			versionError()
		}
		return true, nil, nil, nil
	}

	var data map[string]interface{} = nil
	if keyRet.Secret != nil && len(keyRet.Secret.Data) > 0 {
		data = keyRet.Secret.Data
	}

	return false, keyRet, data, nil
}

func pkiFetchLeaves(e *Executor, versionError func()) (bool, *PathFetch, []string, error) {
	leavesRet, err := e.FetchIfNotFetched(logical.ListOperation, "/{{mount}}/certs")
	if err != nil {
		return true, nil, nil, err
	}

	if !leavesRet.IsSecretOK() {
		if leavesRet.IsUnsupportedPathError() {
			versionError()
		}

		return true, nil, nil, nil
	}

	if len(leavesRet.ParsedCache) == 0 {
		var leaves []string
		for _, rawSerial := range leavesRet.Secret.Data["keys"].([]interface{}) {
			leaves = append(leaves, rawSerial.(string))
		}
		leavesRet.ParsedCache["leaves"] = leaves
	}

	return false, leavesRet, leavesRet.ParsedCache["leaves"].([]string), nil
}

func pkiFetchLeaf(e *Executor, serial string, versionError func()) (bool, *PathFetch, *x509.Certificate, error) {
	leafRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/cert/"+serial)
	if err != nil {
		return true, nil, nil, err
	}

	if !leafRet.IsSecretOK() {
		if leafRet.IsUnsupportedPathError() {
			versionError()
		}
		return true, nil, nil, nil
	}

	if len(leafRet.ParsedCache) == 0 {
		cert, err := parsePEMCert(leafRet.Secret.Data["certificate"].(string))
		if err != nil {
			return true, nil, nil, fmt.Errorf("unable to parse leaf %v's certificate: %w", serial, err)
		}

		leafRet.ParsedCache["certificate"] = cert
	}

	return false, leafRet, leafRet.ParsedCache["certificate"].(*x509.Certificate), nil
}
