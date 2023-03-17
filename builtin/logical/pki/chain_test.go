package pki

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

// For speed, all keys are ECDSA.
type CBGenerateKey struct {
	Name string
}

func (c CBGenerateKey) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	resp, err := CBWrite(b, s, "keys/generate/exported", map[string]interface{}{
		"name": c.Name,
		"algo": "ec",
		"bits": 256,
	})
	if err != nil {
		t.Fatalf("failed to provision key (%v): %v", c.Name, err)
	}
	knownKeys[c.Name] = resp.Data["private"].(string)
}

// Generate a root.
type CBGenerateRoot struct {
	Key          string
	Existing     bool
	Name         string
	CommonName   string
	ErrorMessage string
}

func (c CBGenerateRoot) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	url := "issuers/generate/root/"
	data := make(map[string]interface{})

	if c.Existing {
		url += "existing"
		data["key_ref"] = c.Key
	} else {
		url += "exported"
		data["key_type"] = "ec"
		data["key_bits"] = 256
		data["key_name"] = c.Key
	}

	data["issuer_name"] = c.Name
	data["common_name"] = c.Name
	if len(c.CommonName) > 0 {
		data["common_name"] = c.CommonName
	}

	resp, err := CBWrite(b, s, url, data)
	if err != nil {
		if len(c.ErrorMessage) > 0 {
			if !strings.Contains(err.Error(), c.ErrorMessage) {
				t.Fatalf("failed to generate root cert for issuer (%v): expected (%v) in error message but got %v", c.Name, c.ErrorMessage, err)
			}
			return
		}
		t.Fatalf("failed to provision issuer (%v): %v / body: %v", c.Name, err, data)
	} else if len(c.ErrorMessage) > 0 {
		t.Fatalf("expected to fail generation of issuer (%v) with error message containing (%v)", c.Name, c.ErrorMessage)
	}

	if !c.Existing {
		knownKeys[c.Key] = resp.Data["private_key"].(string)
	}

	knownCerts[c.Name] = resp.Data["certificate"].(string)

	// Validate key_id matches.
	url = "key/" + c.Key
	resp, err = CBRead(b, s, url)
	if err != nil {
		t.Fatalf("failed to fetch key for name %v: %v", c.Key, err)
	}
	if resp == nil {
		t.Fatalf("failed to fetch key for name %v: nil response", c.Key)
	}

	expectedKeyId := resp.Data["key_id"]

	url = "issuer/" + c.Name
	resp, err = CBRead(b, s, url)
	if err != nil {
		t.Fatalf("failed to fetch issuer for name %v: %v", c.Name, err)
	}
	if resp == nil {
		t.Fatalf("failed to fetch issuer for name %v: nil response", c.Name)
	}

	actualKeyId := resp.Data["key_id"]
	if expectedKeyId != actualKeyId {
		t.Fatalf("expected issuer %v to have key matching %v but got mismatch: %v vs %v", c.Name, c.Key, actualKeyId, expectedKeyId)
	}
}

// Generate an intermediate. Might not really be an intermediate; might be
// a cross-signed cert.
type CBGenerateIntermediate struct {
	Key                string
	Existing           bool
	Name               string
	CommonName         string
	Parent             string
	ImportErrorMessage string
}

func (c CBGenerateIntermediate) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	// Build CSR
	url := "issuers/generate/intermediate/"
	data := make(map[string]interface{})

	if c.Existing {
		url += "existing"
		data["key_ref"] = c.Key
	} else {
		url += "exported"
		data["key_type"] = "ec"
		data["key_bits"] = 256
		data["key_name"] = c.Key
	}

	resp, err := CBWrite(b, s, url, data)
	if err != nil {
		t.Fatalf("failed to generate CSR for issuer (%v): %v / body: %v", c.Name, err, data)
	}

	if !c.Existing {
		knownKeys[c.Key] = resp.Data["private_key"].(string)
	}

	csr := resp.Data["csr"].(string)

	// Sign CSR
	url = fmt.Sprintf("issuer/%s/sign-intermediate", c.Parent)
	data = make(map[string]interface{})
	data["csr"] = csr
	data["common_name"] = c.Name
	if len(c.CommonName) > 0 {
		data["common_name"] = c.CommonName
	}
	resp, err = CBWrite(b, s, url, data)
	if err != nil {
		t.Fatalf("failed to sign CSR for issuer (%v): %v / body: %v", c.Name, err, data)
	}

	knownCerts[c.Name] = strings.TrimSpace(resp.Data["certificate"].(string))

	// Set the signed intermediate
	url = "intermediate/set-signed"
	data = make(map[string]interface{})
	data["certificate"] = knownCerts[c.Name]
	data["issuer_name"] = c.Name

	resp, err = CBWrite(b, s, url, data)
	if err != nil {
		if len(c.ImportErrorMessage) > 0 {
			if !strings.Contains(err.Error(), c.ImportErrorMessage) {
				t.Fatalf("failed to import signed cert for issuer (%v): expected (%v) in error message but got %v", c.Name, c.ImportErrorMessage, err)
			}
			return
		}

		t.Fatalf("failed to import signed cert for issuer (%v): %v / body: %v", c.Name, err, data)
	} else if len(c.ImportErrorMessage) > 0 {
		t.Fatalf("expected to fail import (with error %v) of cert for issuer (%v) but was success: response: %v", c.ImportErrorMessage, c.Name, resp)
	}

	// Update the name since set-signed doesn't actually take an issuer name
	// parameter.
	rawNewCerts := resp.Data["imported_issuers"].([]string)
	if len(rawNewCerts) != 1 {
		t.Fatalf("Expected a single new certificate during import of signed cert for %v: got %v\nresp: %v", c.Name, len(rawNewCerts), resp)
	}

	newCertId := rawNewCerts[0]
	_, err = CBWrite(b, s, "issuer/"+newCertId, map[string]interface{}{
		"issuer_name": c.Name,
	})
	if err != nil {
		t.Fatalf("failed to update name for issuer (%v/%v): %v", c.Name, newCertId, err)
	}

	// Validate key_id matches.
	url = "key/" + c.Key
	resp, err = CBRead(b, s, url)
	if err != nil {
		t.Fatalf("failed to fetch key for name %v: %v", c.Key, err)
	}
	if resp == nil {
		t.Fatalf("failed to fetch key for name %v: nil response", c.Key)
	}

	expectedKeyId := resp.Data["key_id"]

	url = "issuer/" + c.Name
	resp, err = CBRead(b, s, url)
	if err != nil {
		t.Fatalf("failed to fetch issuer for name %v: %v", c.Name, err)
	}
	if resp == nil {
		t.Fatalf("failed to fetch issuer for name %v: nil response", c.Name)
	}

	actualKeyId := resp.Data["key_id"]
	if expectedKeyId != actualKeyId {
		t.Fatalf("expected issuer %v to have key matching %v but got mismatch: %v vs %v", c.Name, c.Key, actualKeyId, expectedKeyId)
	}
}

// Delete an issuer; breaks chains.
type CBDeleteIssuer struct {
	Issuer string
}

func (c CBDeleteIssuer) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	url := fmt.Sprintf("issuer/%v", c.Issuer)
	_, err := CBDelete(b, s, url)
	if err != nil {
		t.Fatalf("failed to delete issuer (%v): %v", c.Issuer, err)
	}

	delete(knownCerts, c.Issuer)
}

// Validate the specified chain exists, by name.
type CBValidateChain struct {
	Chains  map[string][]string
	Aliases map[string]string
}

func (c CBValidateChain) ChainToPEMs(t testing.TB, parent string, chain []string, knownCerts map[string]string) []string {
	var result []string
	for entryIndex, entry := range chain {
		var chainEntry string
		modifiedEntry := entry
		if entryIndex == 0 && entry == "self" {
			modifiedEntry = parent
		}
		for pattern, replacement := range c.Aliases {
			modifiedEntry = strings.ReplaceAll(modifiedEntry, pattern, replacement)
		}
		for _, issuer := range strings.Split(modifiedEntry, ",") {
			cert, ok := knownCerts[issuer]
			if !ok {
				t.Fatalf("Unknown issuer %v in chain for %v: %v", issuer, parent, chain)
			}

			chainEntry += cert
		}
		result = append(result, chainEntry)
	}

	return result
}

func (c CBValidateChain) FindNameForCert(t testing.TB, cert string, knownCerts map[string]string) string {
	for issuer, known := range knownCerts {
		if strings.TrimSpace(known) == strings.TrimSpace(cert) {
			return issuer
		}
	}

	t.Fatalf("Unable to find cert:\n[%v]\nin known map:\n%v\n", cert, knownCerts)
	return ""
}

func (c CBValidateChain) PrettyChain(t testing.TB, chain []string, knownCerts map[string]string) []string {
	var prettyChain []string
	for _, cert := range chain {
		prettyChain = append(prettyChain, c.FindNameForCert(t, cert, knownCerts))
	}

	return prettyChain
}

func ToCertificate(t testing.TB, cert string) *x509.Certificate {
	block, _ := pem.Decode([]byte(cert))
	if block == nil {
		t.Fatalf("Unable to parse certificate: nil PEM block\n[%v]\n", cert)
	}

	ret, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Unable to parse certificate: %v\n[%v]\n", err, cert)
	}

	return ret
}

func ToCRL(t testing.TB, crl string, issuer *x509.Certificate) *pkix.CertificateList {
	block, _ := pem.Decode([]byte(crl))
	if block == nil {
		t.Fatalf("Unable to parse CRL: nil PEM block\n[%v]\n", crl)
	}

	ret, err := x509.ParseCRL(block.Bytes)
	if err != nil {
		t.Fatalf("Unable to parse CRL: %v\n[%v]\n", err, crl)
	}

	if issuer != nil {
		if err := issuer.CheckCRLSignature(ret); err != nil {
			t.Fatalf("Unable to check CRL signature: %v\n[%v]\n[%v]\n", err, crl, issuer)
		}
	}

	return ret
}

func (c CBValidateChain) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	for issuer, chain := range c.Chains {
		resp, err := CBRead(b, s, "issuer/"+issuer)
		if err != nil {
			t.Fatalf("failed to get chain for issuer (%v): %v", issuer, err)
		}

		rawCurrentChain := resp.Data["ca_chain"].([]string)
		var currentChain []string
		for _, entry := range rawCurrentChain {
			currentChain = append(currentChain, strings.TrimSpace(entry))
		}

		// Ensure the issuer cert is always first.
		if currentChain[0] != knownCerts[issuer] {
			pretty := c.FindNameForCert(t, currentChain[0], knownCerts)
			t.Fatalf("expected certificate at index 0 to be self:\n[%v]\n[pretty: %v]\nis not the issuer's cert:\n[%v]\n[pretty: %v]", currentChain[0], pretty, knownCerts[issuer], issuer)
		}

		// Validate it against the expected chain.
		expectedChain := c.ChainToPEMs(t, issuer, chain, knownCerts)
		if len(currentChain) != len(expectedChain) {
			prettyCurrentChain := c.PrettyChain(t, currentChain, knownCerts)
			t.Fatalf("Lengths of chains for issuer %v mismatched: got %v vs expected %v:\n[%v]\n[pretty: %v]\n[%v]\n[pretty: %v]", issuer, len(currentChain), len(expectedChain), currentChain, prettyCurrentChain, expectedChain, chain)
		}

		for currentIndex, currentCert := range currentChain {
			// Chains might be forked so we may not be able to strictly validate
			// the chain against a single value. Instead, use strings.Contains
			// to validate the current cert is in the list of allowed
			// possibilities.
			if !strings.Contains(expectedChain[currentIndex], currentCert) {
				pretty := c.FindNameForCert(t, currentCert, knownCerts)
				t.Fatalf("chain mismatch at index %v for issuer %v: got cert:\n[%v]\n[pretty: %v]\nbut expected one of\n[%v]\n[pretty: %v]\n", currentIndex, issuer, currentCert, pretty, expectedChain[currentIndex], chain[currentIndex])
			}
		}

		// Due to alternate paths, the above doesn't ensure ensure each cert
		// in the chain is only used once. Validate that now.
		for thisIndex, thisCert := range currentChain {
			for otherIndex, otherCert := range currentChain[thisIndex+1:] {
				if thisCert == otherCert {
					thisPretty := c.FindNameForCert(t, thisCert, knownCerts)
					otherPretty := c.FindNameForCert(t, otherCert, knownCerts)
					otherIndex += thisIndex + 1
					t.Fatalf("cert reused in chain for %v:\n[%v]\n[pretty: %v / index: %v]\n[%v]\n[pretty: %v / index: %v]\n", issuer, thisCert, thisPretty, thisIndex, otherCert, otherPretty, otherIndex)
				}
			}
		}

		// Finally, validate that all certs verify something that came before
		// it. In the linear chain sense, this should strictly mean that the
		// parent comes before the child.
		for thisIndex, thisCertPem := range currentChain[1:] {
			thisIndex += 1 // Absolute index.
			parentCert := ToCertificate(t, thisCertPem)

			// Iterate backwards; prefer the most recent cert to the older
			// certs.
			foundCert := false
			for otherIndex := thisIndex - 1; otherIndex >= 0; otherIndex-- {
				otherCertPem := currentChain[otherIndex]
				childCert := ToCertificate(t, otherCertPem)

				if err := childCert.CheckSignatureFrom(parentCert); err == nil {
					foundCert = true
				}
			}

			if !foundCert {
				pretty := c.FindNameForCert(t, thisCertPem, knownCerts)
				t.Fatalf("malformed test scenario: certificate at chain index %v when validating %v does not validate any previous certificates:\n[%v]\n[pretty: %v]\n", thisIndex, issuer, thisCertPem, pretty)
			}
		}
	}
}

// Update an issuer
type CBUpdateIssuer struct {
	Name    string
	CAChain []string
	Usage   string
}

func (c CBUpdateIssuer) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	url := "issuer/" + c.Name
	data := make(map[string]interface{})
	data["issuer_name"] = c.Name

	resp, err := CBRead(b, s, url)
	if err != nil {
		t.Fatalf("failed to read issuer (%v): %v", c.Name, err)
	}

	if len(c.CAChain) == 1 && c.CAChain[0] == "existing" {
		data["manual_chain"] = resp.Data["manual_chain"]
	} else {
		data["manual_chain"] = c.CAChain
	}

	if c.Usage == "existing" {
		data["usage"] = resp.Data["usage"]
	} else if len(c.Usage) == 0 {
		data["usage"] = "read-only,issuing-certificates,crl-signing"
	} else {
		data["usage"] = c.Usage
	}

	_, err = CBWrite(b, s, url, data)
	if err != nil {
		t.Fatalf("failed to update issuer (%v): %v / body: %v", c.Name, err, data)
	}
}

// Issue a leaf, revoke it, and then validate it appears on the CRL.
type CBIssueLeaf struct {
	Issuer string
	Role   string
}

func (c CBIssueLeaf) IssueLeaf(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string, errorMessage string) *logical.Response {
	// Write a role
	url := "roles/" + c.Role
	data := make(map[string]interface{})
	data["allow_localhost"] = true
	data["ttl"] = "200s"
	data["key_type"] = "ec"

	_, err := CBWrite(b, s, url, data)
	if err != nil {
		t.Fatalf("failed to update role (%v): %v / body: %v", c.Role, err, data)
	}

	// Issue the certificate.
	url = "issuer/" + c.Issuer + "/issue/" + c.Role
	data = make(map[string]interface{})
	data["common_name"] = "localhost"

	resp, err := CBWrite(b, s, url, data)
	if err != nil {
		if len(errorMessage) >= 0 {
			if !strings.Contains(err.Error(), errorMessage) {
				t.Fatalf("failed to issue cert (%v via %v): %v / body: %v\nExpected error message: %v", c.Issuer, c.Role, err, data, errorMessage)
			}

			return nil
		}

		t.Fatalf("failed to issue cert (%v via %v): %v / body: %v", c.Issuer, c.Role, err, data)
	}
	if resp == nil {
		t.Fatalf("failed to issue cert (%v via %v): nil response / body: %v", c.Issuer, c.Role, data)
	}

	raw_cert := resp.Data["certificate"].(string)
	cert := ToCertificate(t, raw_cert)
	raw_issuer := resp.Data["issuing_ca"].(string)
	issuer := ToCertificate(t, raw_issuer)

	// Validate issuer and signatures are good.
	if strings.TrimSpace(raw_issuer) != strings.TrimSpace(knownCerts[c.Issuer]) {
		t.Fatalf("signing certificate ended with wrong certificate for issuer %v:\n[%v]\n\nvs\n\n[%v]\n", c.Issuer, raw_issuer, knownCerts[c.Issuer])
	}

	if err := cert.CheckSignatureFrom(issuer); err != nil {
		t.Fatalf("failed to verify signature on issued certificate from %v: %v\n[%v]\n[%v]\n", c.Issuer, err, raw_cert, raw_issuer)
	}

	return resp
}

func (c CBIssueLeaf) RevokeLeaf(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string, issueResponse *logical.Response, hasCRL bool, isDefault bool) {
	api_serial := issueResponse.Data["serial_number"].(string)
	raw_cert := issueResponse.Data["certificate"].(string)
	cert := ToCertificate(t, raw_cert)
	raw_issuer := issueResponse.Data["issuing_ca"].(string)
	issuer := ToCertificate(t, raw_issuer)

	// Revoke the certificate.
	url := "revoke"
	data := make(map[string]interface{})
	data["serial_number"] = api_serial
	resp, err := CBWrite(b, s, url, data)
	if err != nil {
		t.Fatalf("failed to revoke issued certificate (%v) under role %v / issuer %v: %v", api_serial, c.Role, c.Issuer, err)
	}
	if resp == nil {
		t.Fatalf("failed to revoke issued certificate (%v) under role %v / issuer %v: nil response", api_serial, c.Role, c.Issuer)
	}
	if _, ok := resp.Data["revocation_time"]; !ok {
		t.Fatalf("failed to revoke issued certificate (%v) under role %v / issuer %v: expected response parameter revocation_time was missing from response:\n%v", api_serial, c.Role, c.Issuer, resp.Data)
	}

	if !hasCRL && isDefault {
		// Nothing further we can test here. We could re-enable CRL building
		// and check that it works, but that seems like a stretch. Other
		// issuers might be functionally the same as this issuer (and thus
		// this CRL will still be issued), but that requires more work to
		// check and verify.
		return
	}

	// Verify it is on this issuer's CRL.
	url = "issuer/" + c.Issuer + "/crl"
	resp, err = CBRead(b, s, url)
	if err != nil {
		t.Fatalf("failed to fetch CRL for issuer %v: %v", c.Issuer, err)
	}
	if resp == nil {
		t.Fatalf("failed to fetch CRL for issuer %v: nil response", c.Issuer)
	}

	raw_crl := resp.Data["crl"].(string)
	crl := ToCRL(t, raw_crl, issuer)

	foundCert := requireSerialNumberInCRL(nil, crl.TBSCertList, api_serial)
	if !foundCert {
		if !hasCRL && !isDefault {
			// Update the issuer we expect to find this on.
			resp, err := CBRead(b, s, "config/issuers")
			if err != nil {
				t.Fatalf("failed to read default issuer config: %v", err)
			}
			if resp == nil {
				t.Fatalf("failed to read default issuer config: nil response")
			}
			defaultID := resp.Data["default"].(issuerID).String()
			c.Issuer = defaultID
			issuer = nil
		}

		// Verify it is on the default issuer's CRL.
		url = "issuer/" + c.Issuer + "/crl"
		resp, err = CBRead(b, s, url)
		if err != nil {
			t.Fatalf("failed to fetch CRL for issuer %v: %v", c.Issuer, err)
		}
		if resp == nil {
			t.Fatalf("failed to fetch CRL for issuer %v: nil response", c.Issuer)
		}

		raw_crl = resp.Data["crl"].(string)
		crl = ToCRL(t, raw_crl, issuer)

		foundCert = requireSerialNumberInCRL(nil, crl.TBSCertList, api_serial)
	}

	if !foundCert {
		// If CRL building is broken, this is useful for finding which issuer's
		// CRL the revoked cert actually appears on.
		for issuerName := range knownCerts {
			url = "issuer/" + issuerName + "/crl"
			resp, err = CBRead(b, s, url)
			if err != nil {
				t.Fatalf("failed to fetch CRL for issuer %v: %v", issuerName, err)
			}
			if resp == nil {
				t.Fatalf("failed to fetch CRL for issuer %v: nil response", issuerName)
			}

			raw_crl := resp.Data["crl"].(string)
			crl := ToCRL(t, raw_crl, nil)

			for index, revoked := range crl.TBSCertList.RevokedCertificates {
				// t.Logf("[%v] revoked serial number: %v -- vs -- %v", index, revoked.SerialNumber, cert.SerialNumber)
				if revoked.SerialNumber.Cmp(cert.SerialNumber) == 0 {
					t.Logf("found revoked cert at index: %v for unexpected issuer: %v", index, issuerName)
					break
				}
			}
		}

		t.Fatalf("expected to find certificate with serial [%v] on issuer %v's CRL but was missing: %v revoked certs\n\nCRL:\n[%v]\n\nLeaf:\n[%v]\n\nIssuer:\n[%v]\n", api_serial, c.Issuer, len(crl.TBSCertList.RevokedCertificates), raw_crl, raw_cert, raw_issuer)
	}
}

func (c CBIssueLeaf) Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	if len(c.Role) == 0 {
		c.Role = "testing"
	}

	resp, err := CBRead(b, s, "config/issuers")
	if err != nil {
		t.Fatalf("failed to read default issuer config: %v", err)
	}
	if resp == nil {
		t.Fatalf("failed to read default issuer config: nil response")
	}
	defaultID := resp.Data["default"].(issuerID).String()

	resp, err = CBRead(b, s, "issuer/"+c.Issuer)
	if err != nil {
		t.Fatalf("failed to read issuer %v: %v", c.Issuer, err)
	}
	if resp == nil {
		t.Fatalf("failed to read issuer %v: nil response", c.Issuer)
	}
	ourID := resp.Data["issuer_id"].(issuerID).String()
	areDefault := ourID == defaultID

	for _, usage := range []string{"read-only", "crl-signing", "issuing-certificates", "issuing-certificates,crl-signing"} {
		ui := CBUpdateIssuer{
			Name:    c.Issuer,
			CAChain: []string{"existing"},
			Usage:   usage,
		}
		ui.Run(t, b, s, knownKeys, knownCerts)

		ilError := "requested usage issuing-certificates for issuer"
		hasIssuing := strings.Contains(usage, "issuing-certificates")
		if hasIssuing {
			ilError = ""
		}

		hasCRL := strings.Contains(usage, "crl-signing")

		resp := c.IssueLeaf(t, b, s, knownKeys, knownCerts, ilError)
		if resp == nil && !hasIssuing {
			continue
		}

		c.RevokeLeaf(t, b, s, knownKeys, knownCerts, resp, hasCRL, areDefault)
	}
}

// Stable ordering
func ensureStableOrderingOfChains(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string) {
	// Start by fetching all chains
	certChains := make(map[string][]string)
	for issuer := range knownCerts {
		resp, err := CBRead(b, s, "issuer/"+issuer)
		if err != nil {
			t.Fatalf("failed to get chain for issuer (%v): %v", issuer, err)
		}

		rawCurrentChain := resp.Data["ca_chain"].([]string)
		var currentChain []string
		for _, entry := range rawCurrentChain {
			currentChain = append(currentChain, strings.TrimSpace(entry))
		}

		certChains[issuer] = currentChain
	}

	// Now, generate a bunch of arbitrary roots and validate the chain is
	// consistent.
	var runs []time.Duration
	for i := 0; i < 10; i++ {
		name := "stable-order-root-" + strconv.Itoa(i)
		step := CBGenerateRoot{
			Key:  name,
			Name: name,
		}
		step.Run(t, b, s, make(map[string]string), make(map[string]string))

		before := time.Now()
		_, err := CBDelete(b, s, "issuer/"+name)
		if err != nil {
			t.Fatalf("failed to delete temporary testing issuer %v: %v", name, err)
		}
		after := time.Now()
		elapsed := after.Sub(before)
		runs = append(runs, elapsed)

		for issuer := range knownCerts {
			resp, err := CBRead(b, s, "issuer/"+issuer)
			if err != nil {
				t.Fatalf("failed to get chain for issuer (%v): %v", issuer, err)
			}

			rawCurrentChain := resp.Data["ca_chain"].([]string)
			for index, entry := range rawCurrentChain {
				if strings.TrimSpace(entry) != certChains[issuer][index] {
					t.Fatalf("iteration %d - chain for issuer %v differed at index %d\n%v\nvs\n%v", i, issuer, index, entry, certChains[issuer][index])
				}
			}
		}
	}

	min := runs[0]
	max := runs[0]
	var avg time.Duration
	for _, run := range runs {
		if run < min {
			min = run
		}

		if run > max {
			max = run
		}

		avg += run
	}
	avg = avg / time.Duration(len(runs))

	t.Logf("Chain building run time (deletion) - min: %v / avg: %v / max: %v - entries: %v", min, avg, max, runs)
}

type CBTestStep interface {
	Run(t testing.TB, b *backend, s logical.Storage, knownKeys map[string]string, knownCerts map[string]string)
}

type CBTestScenario struct {
	Steps []CBTestStep
}

var chainBuildingTestCases = []CBTestScenario{
	{
		// This test builds up two cliques lined by a cycle, dropping into
		// a single intermediate.
		Steps: []CBTestStep{
			// Create a reissued certificate using the same key. These
			// should validate themselves.
			CBGenerateRoot{
				Key:        "key-root-old",
				Name:       "root-old-a",
				CommonName: "root-old",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a": {"self"},
				},
			},
			// After adding the second root using the same key and common
			// name, there should now be two certs in each chain.
			CBGenerateRoot{
				Key:        "key-root-old",
				Existing:   true,
				Name:       "root-old-b",
				CommonName: "root-old",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a": {"self", "root-old-b"},
					"root-old-b": {"self", "root-old-a"},
				},
			},
			// After adding a third root, there are now two possibilities for
			// each later chain entry.
			CBGenerateRoot{
				Key:        "key-root-old",
				Existing:   true,
				Name:       "root-old-c",
				CommonName: "root-old",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a": {"self", "root-old-bc", "root-old-bc"},
					"root-old-b": {"self", "root-old-ac", "root-old-ac"},
					"root-old-c": {"self", "root-old-ab", "root-old-ab"},
				},
				Aliases: map[string]string{
					"root-old-ac": "root-old-a,root-old-c",
					"root-old-ab": "root-old-a,root-old-b",
					"root-old-bc": "root-old-b,root-old-c",
				},
			},
			// If we generate an unrelated issuer, it shouldn't affect either
			// chain.
			CBGenerateRoot{
				Key:        "key-root-new",
				Name:       "root-new-a",
				CommonName: "root-new",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a": {"self", "root-old-bc", "root-old-bc"},
					"root-old-b": {"self", "root-old-ac", "root-old-ac"},
					"root-old-c": {"self", "root-old-ab", "root-old-ab"},
					"root-new-a": {"self"},
				},
				Aliases: map[string]string{
					"root-old-ac": "root-old-a,root-old-c",
					"root-old-ab": "root-old-a,root-old-b",
					"root-old-bc": "root-old-b,root-old-c",
				},
			},
			// Reissuing this new root should form another clique.
			CBGenerateRoot{
				Key:        "key-root-new",
				Existing:   true,
				Name:       "root-new-b",
				CommonName: "root-new",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a": {"self", "root-old-bc", "root-old-bc"},
					"root-old-b": {"self", "root-old-ac", "root-old-ac"},
					"root-old-c": {"self", "root-old-ab", "root-old-ab"},
					"root-new-a": {"self", "root-new-b"},
					"root-new-b": {"self", "root-new-a"},
				},
				Aliases: map[string]string{
					"root-old-ac": "root-old-a,root-old-c",
					"root-old-ab": "root-old-a,root-old-b",
					"root-old-bc": "root-old-b,root-old-c",
				},
			},
			// Generating a cross-signed cert from old->new should result
			// in all old clique certs showing up in the new root's paths.
			// This does not form a cycle.
			CBGenerateIntermediate{
				// In order to validate the existing root-new clique, we
				// have to reuse the key and common name here for
				// cross-signing.
				Key:        "key-root-new",
				Existing:   true,
				Name:       "cross-old-new",
				CommonName: "root-new",
				// Which old issuer is used here doesn't matter as they have
				// the same CN and key.
				Parent: "root-old-a",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a":    {"self", "root-old-bc", "root-old-bc"},
					"root-old-b":    {"self", "root-old-ac", "root-old-ac"},
					"root-old-c":    {"self", "root-old-ab", "root-old-ab"},
					"cross-old-new": {"self", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-a":    {"self", "root-new-b", "cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-b":    {"self", "root-new-a", "cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
				},
				Aliases: map[string]string{
					"root-old-ac":  "root-old-a,root-old-c",
					"root-old-ab":  "root-old-a,root-old-b",
					"root-old-bc":  "root-old-b,root-old-c",
					"root-old-abc": "root-old-a,root-old-b,root-old-c",
				},
			},
			// If we create a new intermediate off of the root-new, we should
			// simply add to the existing chain.
			CBGenerateIntermediate{
				Key:    "key-inter-a-root-new",
				Name:   "inter-a-root-new",
				Parent: "root-new-a",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a":    {"self", "root-old-bc", "root-old-bc"},
					"root-old-b":    {"self", "root-old-ac", "root-old-ac"},
					"root-old-c":    {"self", "root-old-ab", "root-old-ab"},
					"cross-old-new": {"self", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-a":    {"self", "root-new-b", "cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-b":    {"self", "root-new-a", "cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					// If we find cross-old-new first, the old clique will be ahead
					// of the new clique; otherwise the new clique will appear first.
					"inter-a-root-new": {"self", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle"},
				},
				Aliases: map[string]string{
					"root-old-ac":  "root-old-a,root-old-c",
					"root-old-ab":  "root-old-a,root-old-b",
					"root-old-bc":  "root-old-b,root-old-c",
					"root-old-abc": "root-old-a,root-old-b,root-old-c",
					"full-cycle":   "root-old-a,root-old-b,root-old-c,cross-old-new,root-new-a,root-new-b",
				},
			},
			// Now, if we cross-sign back from new to old, we should
			// form cycle with multiple reissued cliques. This means
			// all nodes will have the same chain.
			CBGenerateIntermediate{
				// In order to validate the existing root-old clique, we
				// have to reuse the key and common name here for
				// cross-signing.
				Key:        "key-root-old",
				Existing:   true,
				Name:       "cross-new-old",
				CommonName: "root-old",
				// Which new issuer is used here doesn't matter as they have
				// the same CN and key.
				Parent: "root-new-a",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a":       {"self", "root-old-bc", "root-old-bc", "both-cross-old-new", "both-cross-old-new", "root-new-ab", "root-new-ab"},
					"root-old-b":       {"self", "root-old-ac", "root-old-ac", "both-cross-old-new", "both-cross-old-new", "root-new-ab", "root-new-ab"},
					"root-old-c":       {"self", "root-old-ab", "root-old-ab", "both-cross-old-new", "both-cross-old-new", "root-new-ab", "root-new-ab"},
					"cross-old-new":    {"self", "cross-new-old", "both-cliques", "both-cliques", "both-cliques", "both-cliques", "both-cliques"},
					"cross-new-old":    {"self", "cross-old-new", "both-cliques", "both-cliques", "both-cliques", "both-cliques", "both-cliques"},
					"root-new-a":       {"self", "root-new-b", "both-cross-old-new", "both-cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-b":       {"self", "root-new-a", "both-cross-old-new", "both-cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"inter-a-root-new": {"self", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle"},
				},
				Aliases: map[string]string{
					"root-old-ac":        "root-old-a,root-old-c",
					"root-old-ab":        "root-old-a,root-old-b",
					"root-old-bc":        "root-old-b,root-old-c",
					"root-old-abc":       "root-old-a,root-old-b,root-old-c",
					"root-new-ab":        "root-new-a,root-new-b",
					"both-cross-old-new": "cross-old-new,cross-new-old",
					"both-cliques":       "root-old-a,root-old-b,root-old-c,root-new-a,root-new-b",
					"full-cycle":         "root-old-a,root-old-b,root-old-c,cross-old-new,cross-new-old,root-new-a,root-new-b",
				},
			},
			// Update each old root to only include itself.
			CBUpdateIssuer{
				Name:    "root-old-a",
				CAChain: []string{"root-old-a"},
			},
			CBUpdateIssuer{
				Name:    "root-old-b",
				CAChain: []string{"root-old-b"},
			},
			CBUpdateIssuer{
				Name:    "root-old-c",
				CAChain: []string{"root-old-c"},
			},
			// Step 19
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a":       {"self"},
					"root-old-b":       {"self"},
					"root-old-c":       {"self"},
					"cross-old-new":    {"self", "cross-new-old", "both-cliques", "both-cliques", "both-cliques", "both-cliques", "both-cliques"},
					"cross-new-old":    {"self", "cross-old-new", "both-cliques", "both-cliques", "both-cliques", "both-cliques", "both-cliques"},
					"root-new-a":       {"self", "root-new-b", "both-cross-old-new", "both-cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-b":       {"self", "root-new-a", "both-cross-old-new", "both-cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"inter-a-root-new": {"self", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle"},
				},
				Aliases: map[string]string{
					"root-old-ac":        "root-old-a,root-old-c",
					"root-old-ab":        "root-old-a,root-old-b",
					"root-old-bc":        "root-old-b,root-old-c",
					"root-old-abc":       "root-old-a,root-old-b,root-old-c",
					"root-new-ab":        "root-new-a,root-new-b",
					"both-cross-old-new": "cross-old-new,cross-new-old",
					"both-cliques":       "root-old-a,root-old-b,root-old-c,root-new-a,root-new-b",
					"full-cycle":         "root-old-a,root-old-b,root-old-c,cross-old-new,cross-new-old,root-new-a,root-new-b",
				},
			},
			// Reset the old roots; should get the original chains back.
			CBUpdateIssuer{
				Name: "root-old-a",
			},
			CBUpdateIssuer{
				Name: "root-old-b",
			},
			CBUpdateIssuer{
				Name: "root-old-c",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-old-a":       {"self", "root-old-bc", "root-old-bc", "both-cross-old-new", "both-cross-old-new", "root-new-ab", "root-new-ab"},
					"root-old-b":       {"self", "root-old-ac", "root-old-ac", "both-cross-old-new", "both-cross-old-new", "root-new-ab", "root-new-ab"},
					"root-old-c":       {"self", "root-old-ab", "root-old-ab", "both-cross-old-new", "both-cross-old-new", "root-new-ab", "root-new-ab"},
					"cross-old-new":    {"self", "cross-new-old", "both-cliques", "both-cliques", "both-cliques", "both-cliques", "both-cliques"},
					"cross-new-old":    {"self", "cross-old-new", "both-cliques", "both-cliques", "both-cliques", "both-cliques", "both-cliques"},
					"root-new-a":       {"self", "root-new-b", "both-cross-old-new", "both-cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"root-new-b":       {"self", "root-new-a", "both-cross-old-new", "both-cross-old-new", "root-old-abc", "root-old-abc", "root-old-abc"},
					"inter-a-root-new": {"self", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle", "full-cycle"},
				},
				Aliases: map[string]string{
					"root-old-ac":        "root-old-a,root-old-c",
					"root-old-ab":        "root-old-a,root-old-b",
					"root-old-bc":        "root-old-b,root-old-c",
					"root-old-abc":       "root-old-a,root-old-b,root-old-c",
					"root-new-ab":        "root-new-a,root-new-b",
					"both-cross-old-new": "cross-old-new,cross-new-old",
					"both-cliques":       "root-old-a,root-old-b,root-old-c,root-new-a,root-new-b",
					"full-cycle":         "root-old-a,root-old-b,root-old-c,cross-old-new,cross-new-old,root-new-a,root-new-b",
				},
			},
			CBIssueLeaf{Issuer: "root-old-a"},
			CBIssueLeaf{Issuer: "root-old-b"},
			CBIssueLeaf{Issuer: "root-old-c"},
			CBIssueLeaf{Issuer: "cross-old-new"},
			CBIssueLeaf{Issuer: "cross-new-old"},
			CBIssueLeaf{Issuer: "root-new-a"},
			CBIssueLeaf{Issuer: "root-new-b"},
			CBIssueLeaf{Issuer: "inter-a-root-new"},
		},
	},
	{
		// Here we're testing our chain capacity. First we'll create a
		// bunch of unique roots to form a cycle of length 10.
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:        "key-root-a",
				Name:       "root-a",
				CommonName: "root-a",
			},
			CBGenerateRoot{
				Key:        "key-root-b",
				Name:       "root-b",
				CommonName: "root-b",
			},
			CBGenerateRoot{
				Key:        "key-root-c",
				Name:       "root-c",
				CommonName: "root-c",
			},
			CBGenerateRoot{
				Key:        "key-root-d",
				Name:       "root-d",
				CommonName: "root-d",
			},
			CBGenerateRoot{
				Key:        "key-root-e",
				Name:       "root-e",
				CommonName: "root-e",
			},
			// They should all be disjoint to start.
			CBValidateChain{
				Chains: map[string][]string{
					"root-a": {"self"},
					"root-b": {"self"},
					"root-c": {"self"},
					"root-d": {"self"},
					"root-e": {"self"},
				},
			},
			// Start the cross-signing chains. These are all linear, so there's
			// no error expected; they're just long.
			CBGenerateIntermediate{
				Key:        "key-root-b",
				Existing:   true,
				Name:       "cross-a-b",
				CommonName: "root-b",
				Parent:     "root-a",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-a":    {"self"},
					"cross-a-b": {"self", "root-a"},
					"root-b":    {"self", "cross-a-b", "root-a"},
					"root-c":    {"self"},
					"root-d":    {"self"},
					"root-e":    {"self"},
				},
			},
			CBGenerateIntermediate{
				Key:        "key-root-c",
				Existing:   true,
				Name:       "cross-b-c",
				CommonName: "root-c",
				Parent:     "root-b",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-a":    {"self"},
					"cross-a-b": {"self", "root-a"},
					"root-b":    {"self", "cross-a-b", "root-a"},
					"cross-b-c": {"self", "b-or-cross", "b-chained-cross", "b-chained-cross"},
					"root-c":    {"self", "cross-b-c", "b-or-cross", "b-chained-cross", "b-chained-cross"},
					"root-d":    {"self"},
					"root-e":    {"self"},
				},
				Aliases: map[string]string{
					"b-or-cross":      "root-b,cross-a-b",
					"b-chained-cross": "root-b,cross-a-b,root-a",
				},
			},
			CBGenerateIntermediate{
				Key:        "key-root-d",
				Existing:   true,
				Name:       "cross-c-d",
				CommonName: "root-d",
				Parent:     "root-c",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-a":    {"self"},
					"cross-a-b": {"self", "root-a"},
					"root-b":    {"self", "cross-a-b", "root-a"},
					"cross-b-c": {"self", "b-or-cross", "b-chained-cross", "b-chained-cross"},
					"root-c":    {"self", "cross-b-c", "b-or-cross", "b-chained-cross", "b-chained-cross"},
					"cross-c-d": {"self", "c-or-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross"},
					"root-d":    {"self", "cross-c-d", "c-or-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross"},
					"root-e":    {"self"},
				},
				Aliases: map[string]string{
					"b-or-cross":      "root-b,cross-a-b",
					"b-chained-cross": "root-b,cross-a-b,root-a",
					"c-or-cross":      "root-c,cross-b-c",
					"c-chained-cross": "root-c,cross-b-c,root-b,cross-a-b,root-a",
				},
			},
			CBGenerateIntermediate{
				Key:        "key-root-e",
				Existing:   true,
				Name:       "cross-d-e",
				CommonName: "root-e",
				Parent:     "root-d",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-a":    {"self"},
					"cross-a-b": {"self", "root-a"},
					"root-b":    {"self", "cross-a-b", "root-a"},
					"cross-b-c": {"self", "b-or-cross", "b-chained-cross", "b-chained-cross"},
					"root-c":    {"self", "cross-b-c", "b-or-cross", "b-chained-cross", "b-chained-cross"},
					"cross-c-d": {"self", "c-or-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross"},
					"root-d":    {"self", "cross-c-d", "c-or-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross", "c-chained-cross"},
					"cross-d-e": {"self", "d-or-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross"},
					"root-e":    {"self", "cross-d-e", "d-or-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross", "d-chained-cross"},
				},
				Aliases: map[string]string{
					"b-or-cross":      "root-b,cross-a-b",
					"b-chained-cross": "root-b,cross-a-b,root-a",
					"c-or-cross":      "root-c,cross-b-c",
					"c-chained-cross": "root-c,cross-b-c,root-b,cross-a-b,root-a",
					"d-or-cross":      "root-d,cross-c-d",
					"d-chained-cross": "root-d,cross-c-d,root-c,cross-b-c,root-b,cross-a-b,root-a",
				},
			},
			CBIssueLeaf{Issuer: "root-a"},
			CBIssueLeaf{Issuer: "cross-a-b"},
			CBIssueLeaf{Issuer: "root-b"},
			CBIssueLeaf{Issuer: "cross-b-c"},
			CBIssueLeaf{Issuer: "root-c"},
			CBIssueLeaf{Issuer: "cross-c-d"},
			CBIssueLeaf{Issuer: "root-d"},
			CBIssueLeaf{Issuer: "cross-d-e"},
			CBIssueLeaf{Issuer: "root-e"},
			// Importing the new e->a cross fails because the cycle
			// it builds is too long.
			CBGenerateIntermediate{
				Key:                "key-root-a",
				Existing:           true,
				Name:               "cross-e-a",
				CommonName:         "root-a",
				Parent:             "root-e",
				ImportErrorMessage: "exceeds max size",
			},
			// Deleting any root and one of its crosses (either a->b or b->c)
			// should fix this.
			CBDeleteIssuer{"root-b"},
			CBDeleteIssuer{"cross-b-c"},
			// Importing the new e->a cross fails because the cycle
			// it builds is too long.
			CBGenerateIntermediate{
				Key:        "key-root-a",
				Existing:   true,
				Name:       "cross-e-a",
				CommonName: "root-a",
				Parent:     "root-e",
			},
			CBIssueLeaf{Issuer: "root-a"},
			CBIssueLeaf{Issuer: "cross-a-b"},
			CBIssueLeaf{Issuer: "root-c"},
			CBIssueLeaf{Issuer: "cross-c-d"},
			CBIssueLeaf{Issuer: "root-d"},
			CBIssueLeaf{Issuer: "cross-d-e"},
			CBIssueLeaf{Issuer: "root-e"},
			CBIssueLeaf{Issuer: "cross-e-a"},
		},
	},
	{
		// Here we're testing our clique capacity. First we'll create a
		// bunch of unique roots to form a cycle of length 10.
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:        "key-root",
				Name:       "root-a",
				CommonName: "root",
			},
			CBGenerateRoot{
				Key:        "key-root",
				Existing:   true,
				Name:       "root-b",
				CommonName: "root",
			},
			CBGenerateRoot{
				Key:        "key-root",
				Existing:   true,
				Name:       "root-c",
				CommonName: "root",
			},
			CBGenerateRoot{
				Key:        "key-root",
				Existing:   true,
				Name:       "root-d",
				CommonName: "root",
			},
			CBGenerateRoot{
				Key:        "key-root",
				Existing:   true,
				Name:       "root-e",
				CommonName: "root",
			},
			CBGenerateRoot{
				Key:        "key-root",
				Existing:   true,
				Name:       "root-f",
				CommonName: "root",
			},
			CBIssueLeaf{Issuer: "root-a"},
			CBIssueLeaf{Issuer: "root-b"},
			CBIssueLeaf{Issuer: "root-c"},
			CBIssueLeaf{Issuer: "root-d"},
			CBIssueLeaf{Issuer: "root-e"},
			CBIssueLeaf{Issuer: "root-f"},
			// Seventh reissuance fails.
			CBGenerateRoot{
				Key:          "key-root",
				Existing:     true,
				Name:         "root-g",
				CommonName:   "root",
				ErrorMessage: "excessively reissued certificate",
			},
			// Deleting one and trying again should succeed.
			CBDeleteIssuer{"root-a"},
			CBGenerateRoot{
				Key:        "key-root",
				Existing:   true,
				Name:       "root-g",
				CommonName: "root",
			},
			CBIssueLeaf{Issuer: "root-b"},
			CBIssueLeaf{Issuer: "root-c"},
			CBIssueLeaf{Issuer: "root-d"},
			CBIssueLeaf{Issuer: "root-e"},
			CBIssueLeaf{Issuer: "root-f"},
			CBIssueLeaf{Issuer: "root-g"},
		},
	},
	{
		// There's one more pathological case here: we have a cycle
		// which validates a clique/cycle via cross-signing. We call
		// the parent cycle new roots and the child cycle/clique the
		// old roots.
		Steps: []CBTestStep{
			// New Cycle
			CBGenerateRoot{
				Key:  "key-root-new-a",
				Name: "root-new-a",
			},
			CBGenerateRoot{
				Key:  "key-root-new-b",
				Name: "root-new-b",
			},
			CBGenerateIntermediate{
				Key:        "key-root-new-b",
				Existing:   true,
				Name:       "cross-root-new-b-sig-a",
				CommonName: "root-new-b",
				Parent:     "root-new-a",
			},
			CBGenerateIntermediate{
				Key:        "key-root-new-a",
				Existing:   true,
				Name:       "cross-root-new-a-sig-b",
				CommonName: "root-new-a",
				Parent:     "root-new-b",
			},
			// Old Cycle + Clique
			CBGenerateRoot{
				Key:  "key-root-old-a",
				Name: "root-old-a",
			},
			CBGenerateRoot{
				Key:        "key-root-old-a",
				Existing:   true,
				Name:       "root-old-a-reissued",
				CommonName: "root-old-a",
			},
			CBGenerateRoot{
				Key:  "key-root-old-b",
				Name: "root-old-b",
			},
			CBGenerateRoot{
				Key:        "key-root-old-b",
				Existing:   true,
				Name:       "root-old-b-reissued",
				CommonName: "root-old-b",
			},
			CBGenerateIntermediate{
				Key:        "key-root-old-b",
				Existing:   true,
				Name:       "cross-root-old-b-sig-a",
				CommonName: "root-old-b",
				Parent:     "root-old-a",
			},
			CBGenerateIntermediate{
				Key:        "key-root-old-a",
				Existing:   true,
				Name:       "cross-root-old-a-sig-b",
				CommonName: "root-old-a",
				Parent:     "root-old-b",
			},
			// Validate the chains are separate before linking them.
			CBValidateChain{
				Chains: map[string][]string{
					// New stuff
					"root-new-a":             {"self", "cross-root-new-a-sig-b", "root-new-b-or-cross", "root-new-b-or-cross"},
					"root-new-b":             {"self", "cross-root-new-b-sig-a", "root-new-a-or-cross", "root-new-a-or-cross"},
					"cross-root-new-b-sig-a": {"self", "any-root-new", "any-root-new", "any-root-new"},
					"cross-root-new-a-sig-b": {"self", "any-root-new", "any-root-new", "any-root-new"},

					// Old stuff
					"root-old-a":             {"self", "root-old-a-reissued", "cross-root-old-a-sig-b", "cross-root-old-b-sig-a", "both-root-old-b", "both-root-old-b"},
					"root-old-a-reissued":    {"self", "root-old-a", "cross-root-old-a-sig-b", "cross-root-old-b-sig-a", "both-root-old-b", "both-root-old-b"},
					"root-old-b":             {"self", "root-old-b-reissued", "cross-root-old-b-sig-a", "cross-root-old-a-sig-b", "both-root-old-a", "both-root-old-a"},
					"root-old-b-reissued":    {"self", "root-old-b", "cross-root-old-b-sig-a", "cross-root-old-a-sig-b", "both-root-old-a", "both-root-old-a"},
					"cross-root-old-b-sig-a": {"self", "all-root-old", "all-root-old", "all-root-old", "all-root-old", "all-root-old"},
					"cross-root-old-a-sig-b": {"self", "all-root-old", "all-root-old", "all-root-old", "all-root-old", "all-root-old"},
				},
				Aliases: map[string]string{
					"root-new-a-or-cross": "root-new-a,cross-root-new-a-sig-b",
					"root-new-b-or-cross": "root-new-b,cross-root-new-b-sig-a",
					"both-root-new":       "root-new-a,root-new-b",
					"any-root-new":        "root-new-a,cross-root-new-a-sig-b,root-new-b,cross-root-new-b-sig-a",
					"both-root-old-a":     "root-old-a,root-old-a-reissued",
					"both-root-old-b":     "root-old-b,root-old-b-reissued",
					"all-root-old":        "root-old-a,root-old-a-reissued,root-old-b,root-old-b-reissued,cross-root-old-b-sig-a,cross-root-old-a-sig-b",
				},
			},
			// Finally, generate an intermediate to link new->old. We
			// link root-new-a into root-old-a.
			CBGenerateIntermediate{
				Key:        "key-root-old-a",
				Existing:   true,
				Name:       "cross-root-old-a-sig-root-new-a",
				CommonName: "root-old-a",
				Parent:     "root-new-a",
			},
			CBValidateChain{
				Chains: map[string][]string{
					// New stuff should be unchanged.
					"root-new-a":             {"self", "cross-root-new-a-sig-b", "root-new-b-or-cross", "root-new-b-or-cross"},
					"root-new-b":             {"self", "cross-root-new-b-sig-a", "root-new-a-or-cross", "root-new-a-or-cross"},
					"cross-root-new-b-sig-a": {"self", "any-root-new", "any-root-new", "any-root-new"},
					"cross-root-new-a-sig-b": {"self", "any-root-new", "any-root-new", "any-root-new"},

					// Old stuff
					"root-old-a":             {"self", "root-old-a-reissued", "cross-root-old-a-sig-b", "cross-root-old-b-sig-a", "both-root-old-b", "both-root-old-b", "cross-root-old-a-sig-root-new-a", "any-root-new", "any-root-new", "any-root-new", "any-root-new"},
					"root-old-a-reissued":    {"self", "root-old-a", "cross-root-old-a-sig-b", "cross-root-old-b-sig-a", "both-root-old-b", "both-root-old-b", "cross-root-old-a-sig-root-new-a", "any-root-new", "any-root-new", "any-root-new", "any-root-new"},
					"root-old-b":             {"self", "root-old-b-reissued", "cross-root-old-b-sig-a", "cross-root-old-a-sig-b", "both-root-old-a", "both-root-old-a", "cross-root-old-a-sig-root-new-a", "any-root-new", "any-root-new", "any-root-new", "any-root-new"},
					"root-old-b-reissued":    {"self", "root-old-b", "cross-root-old-b-sig-a", "cross-root-old-a-sig-b", "both-root-old-a", "both-root-old-a", "cross-root-old-a-sig-root-new-a", "any-root-new", "any-root-new", "any-root-new", "any-root-new"},
					"cross-root-old-b-sig-a": {"self", "all-root-old", "all-root-old", "all-root-old", "all-root-old", "all-root-old", "cross-root-old-a-sig-root-new-a", "any-root-new", "any-root-new", "any-root-new", "any-root-new"},
					"cross-root-old-a-sig-b": {"self", "all-root-old", "all-root-old", "all-root-old", "all-root-old", "all-root-old", "cross-root-old-a-sig-root-new-a", "any-root-new", "any-root-new", "any-root-new", "any-root-new"},

					// Link
					"cross-root-old-a-sig-root-new-a": {"self", "root-new-a-or-cross", "any-root-new", "any-root-new", "any-root-new"},
				},
				Aliases: map[string]string{
					"root-new-a-or-cross": "root-new-a,cross-root-new-a-sig-b",
					"root-new-b-or-cross": "root-new-b,cross-root-new-b-sig-a",
					"both-root-new":       "root-new-a,root-new-b",
					"any-root-new":        "root-new-a,cross-root-new-a-sig-b,root-new-b,cross-root-new-b-sig-a",
					"both-root-old-a":     "root-old-a,root-old-a-reissued",
					"both-root-old-b":     "root-old-b,root-old-b-reissued",
					"all-root-old":        "root-old-a,root-old-a-reissued,root-old-b,root-old-b-reissued,cross-root-old-b-sig-a,cross-root-old-a-sig-b",
				},
			},
			CBIssueLeaf{Issuer: "root-new-a"},
			CBIssueLeaf{Issuer: "root-new-b"},
			CBIssueLeaf{Issuer: "cross-root-new-b-sig-a"},
			CBIssueLeaf{Issuer: "cross-root-new-a-sig-b"},
			CBIssueLeaf{Issuer: "root-old-a"},
			CBIssueLeaf{Issuer: "root-old-a-reissued"},
			CBIssueLeaf{Issuer: "root-old-b"},
			CBIssueLeaf{Issuer: "root-old-b-reissued"},
			CBIssueLeaf{Issuer: "cross-root-old-b-sig-a"},
			CBIssueLeaf{Issuer: "cross-root-old-a-sig-b"},
			CBIssueLeaf{Issuer: "cross-root-old-a-sig-root-new-a"},
		},
	},
	{
		// Test a dual-root of trust chaining example with different
		// lengths of chains.
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:  "key-root-new",
				Name: "root-new",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-new",
				Name:   "inter-new",
				Parent: "root-new",
			},
			CBGenerateRoot{
				Key:  "key-root-old",
				Name: "root-old",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-old-a",
				Name:   "inter-old-a",
				Parent: "root-old",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-old-b",
				Name:   "inter-old-b",
				Parent: "inter-old-a",
			},
			// Now generate a cross-signed intermediate to merge these
			// two chains.
			CBGenerateIntermediate{
				Key:        "key-cross-old-new",
				Name:       "cross-old-new-signed-new",
				CommonName: "cross-old-new",
				Parent:     "inter-new",
			},
			CBGenerateIntermediate{
				Key:        "key-cross-old-new",
				Existing:   true,
				Name:       "cross-old-new-signed-old",
				CommonName: "cross-old-new",
				Parent:     "inter-old-b",
			},
			CBGenerateIntermediate{
				Key:    "key-leaf-inter",
				Name:   "leaf-inter",
				Parent: "cross-old-new-signed-new",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-new":                 {"self"},
					"inter-new":                {"self", "root-new"},
					"cross-old-new-signed-new": {"self", "inter-new", "root-new"},
					"root-old":                 {"self"},
					"inter-old-a":              {"self", "root-old"},
					"inter-old-b":              {"self", "inter-old-a", "root-old"},
					"cross-old-new-signed-old": {"self", "inter-old-b", "inter-old-a", "root-old"},
					"leaf-inter":               {"self", "either-cross", "one-intermediate", "other-inter-or-root", "everything-else", "everything-else", "everything-else", "everything-else"},
				},
				Aliases: map[string]string{
					"either-cross":        "cross-old-new-signed-new,cross-old-new-signed-old",
					"one-intermediate":    "inter-new,inter-old-b",
					"other-inter-or-root": "root-new,inter-old-a",
					"everything-else":     "cross-old-new-signed-new,cross-old-new-signed-old,inter-new,inter-old-b,root-new,inter-old-a,root-old",
				},
			},
			CBIssueLeaf{Issuer: "root-new"},
			CBIssueLeaf{Issuer: "inter-new"},
			CBIssueLeaf{Issuer: "root-old"},
			CBIssueLeaf{Issuer: "inter-old-a"},
			CBIssueLeaf{Issuer: "inter-old-b"},
			CBIssueLeaf{Issuer: "cross-old-new-signed-new"},
			CBIssueLeaf{Issuer: "cross-old-new-signed-old"},
			CBIssueLeaf{Issuer: "leaf-inter"},
		},
	},
	{
		// Test just a single root.
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:  "key-root",
				Name: "root",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root": {"self"},
				},
			},
			CBIssueLeaf{Issuer: "root"},
		},
	},
	{
		// Test root + intermediate.
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:  "key-root",
				Name: "root",
			},
			CBGenerateIntermediate{
				Key:    "key-inter",
				Name:   "inter",
				Parent: "root",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root":  {"self"},
					"inter": {"self", "root"},
				},
			},
			CBIssueLeaf{Issuer: "root"},
			CBIssueLeaf{Issuer: "inter"},
		},
	},
	{
		// Test root + intermediate, twice (simulating rotation without
		// chaining).
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:  "key-root-a",
				Name: "root-a",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-a",
				Name:   "inter-a",
				Parent: "root-a",
			},
			CBGenerateRoot{
				Key:  "key-root-b",
				Name: "root-b",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-b",
				Name:   "inter-b",
				Parent: "root-b",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-a":  {"self"},
					"inter-a": {"self", "root-a"},
					"root-b":  {"self"},
					"inter-b": {"self", "root-b"},
				},
			},
			CBIssueLeaf{Issuer: "root-a"},
			CBIssueLeaf{Issuer: "inter-a"},
			CBIssueLeaf{Issuer: "root-b"},
			CBIssueLeaf{Issuer: "inter-b"},
		},
	},
	{
		// Test root + intermediate, twice, chained a->b.
		Steps: []CBTestStep{
			CBGenerateRoot{
				Key:  "key-root-a",
				Name: "root-a",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-a",
				Name:   "inter-a",
				Parent: "root-a",
			},
			CBGenerateRoot{
				Key:  "key-root-b",
				Name: "root-b",
			},
			CBGenerateIntermediate{
				Key:    "key-inter-b",
				Name:   "inter-b",
				Parent: "root-b",
			},
			CBGenerateIntermediate{
				Key:        "key-root-b",
				Existing:   true,
				Name:       "cross-a-b",
				CommonName: "root-b",
				Parent:     "root-a",
			},
			CBValidateChain{
				Chains: map[string][]string{
					"root-a":    {"self"},
					"inter-a":   {"self", "root-a"},
					"root-b":    {"self", "cross-a-b", "root-a"},
					"inter-b":   {"self", "root-b", "cross-a-b", "root-a"},
					"cross-a-b": {"self", "root-a"},
				},
			},
			CBIssueLeaf{Issuer: "root-a"},
			CBIssueLeaf{Issuer: "inter-a"},
			CBIssueLeaf{Issuer: "root-b"},
			CBIssueLeaf{Issuer: "inter-b"},
			CBIssueLeaf{Issuer: "cross-a-b"},
		},
	},
}

func Test_CAChainBuilding(t *testing.T) {
	for testIndex, testCase := range chainBuildingTestCases {
		b, s := createBackendWithStorage(t)

		knownKeys := make(map[string]string)
		knownCerts := make(map[string]string)
		for stepIndex, testStep := range testCase.Steps {
			t.Logf("Running %v / %v", testIndex, stepIndex)
			testStep.Run(t, b, s, knownKeys, knownCerts)
		}

		t.Logf("Checking stable ordering of chains...")
		ensureStableOrderingOfChains(t, b, s, knownKeys, knownCerts)
	}
}

func BenchmarkChainBuilding(benchies *testing.B) {
	for testIndex, testCase := range chainBuildingTestCases {
		name := "test-case-" + strconv.Itoa(testIndex)
		benchies.Run(name, func(bench *testing.B) {
			// Stop the timer as we setup the infra and certs.
			bench.StopTimer()
			bench.ResetTimer()

			b, s := createBackendWithStorage(bench)

			knownKeys := make(map[string]string)
			knownCerts := make(map[string]string)
			for _, testStep := range testCase.Steps {
				testStep.Run(bench, b, s, knownKeys, knownCerts)
			}

			// Run the benchmark.
			ctx := context.Background()
			bench.StartTimer()
			for n := 0; n < bench.N; n++ {
				rebuildIssuersChains(ctx, s, nil)
			}
		})
	}
}
