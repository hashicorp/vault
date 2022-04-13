package pki

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// For speed, all keys are ECDSA.
type CBGenerateKey struct {
	Name string
}

func (c CBGenerateKey) Run(t *testing.T, client *api.Client, mount string, knownKeys map[string]string, knownCerts map[string]string) {
	resp, err := client.Logical().Write(mount+"/keys/generate/exported", map[string]interface{}{
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

func (c CBGenerateRoot) Run(t *testing.T, client *api.Client, mount string, knownKeys map[string]string, knownCerts map[string]string) {
	url := mount + "/issuers/generate/root/"
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

	resp, err := client.Logical().Write(url, data)
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

func (c CBGenerateIntermediate) Run(t *testing.T, client *api.Client, mount string, knownKeys map[string]string, knownCerts map[string]string) {
	// Build CSR
	url := mount + "/issuers/generate/intermediate/"
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

	resp, err := client.Logical().Write(url, data)
	if err != nil {
		t.Fatalf("failed to generate CSR for issuer (%v): %v / body: %v", c.Name, err, data)
	}

	if !c.Existing {
		knownKeys[c.Key] = resp.Data["private_key"].(string)
	}

	csr := resp.Data["csr"].(string)

	// Sign CSR
	url = fmt.Sprintf(mount+"/issuers/%s/sign-intermediate", c.Parent)
	data = make(map[string]interface{})
	data["csr"] = csr
	data["common_name"] = c.Name
	if len(c.CommonName) > 0 {
		data["common_name"] = c.CommonName
	}
	resp, err = client.Logical().Write(url, data)
	if err != nil {
		t.Fatalf("failed to sign CSR for issuer (%v): %v / body: %v", c.Name, err, data)
	}

	knownCerts[c.Name] = strings.TrimSpace(resp.Data["certificate"].(string))

	// Set the signed intermediate
	url = mount + "/intermediate/set-signed"
	data = make(map[string]interface{})
	data["certificate"] = knownCerts[c.Name]
	data["issuer_name"] = c.Name

	resp, err = client.Logical().Write(url, data)
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
	rawNewCerts := resp.Data["imported_issuers"].([]interface{})
	if len(rawNewCerts) != 1 {
		t.Fatalf("Expected a single new certificate during import of signed cert for %v: got %v\nresp: %v", c.Name, len(rawNewCerts), resp)
	}

	newCertId := rawNewCerts[0].(string)
	_, err = client.Logical().Write(mount+"/issuer/"+newCertId, map[string]interface{}{
		"issuer_name": c.Name,
	})
	if err != nil {
		t.Fatalf("failed to update name for issuer (%v/%v): %v", c.Name, newCertId, err)
	}
}

// Delete an issuer; breaks chains.
type CBDeleteIssuer struct {
	Issuer string
}

func (c CBDeleteIssuer) Run(t *testing.T, client *api.Client, mount string, knownKeys map[string]string, knownCerts map[string]string) {
	url := fmt.Sprintf(mount+"/issuer/%v", c.Issuer)
	_, err := client.Logical().Delete(url)
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

func (c CBValidateChain) ChainToPEMs(t *testing.T, parent string, chain []string, knownCerts map[string]string) []string {
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

func (c CBValidateChain) FindNameForCert(t *testing.T, cert string, knownCerts map[string]string) string {
	for issuer, known := range knownCerts {
		if strings.TrimSpace(known) == strings.TrimSpace(cert) {
			return issuer
		}
	}

	t.Fatalf("Unable to find cert:\n[%v]\nin known map:\n%v\n", cert, knownCerts)
	return ""
}

func (c CBValidateChain) PrettyChain(t *testing.T, chain []string, knownCerts map[string]string) []string {
	var prettyChain []string
	for _, cert := range chain {
		prettyChain = append(prettyChain, c.FindNameForCert(t, cert, knownCerts))
	}

	return prettyChain
}

func (c CBValidateChain) ToCertificate(t *testing.T, cert string) *x509.Certificate {
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

func (c CBValidateChain) Run(t *testing.T, client *api.Client, mount string, knownKeys map[string]string, knownCerts map[string]string) {
	for issuer, chain := range c.Chains {
		resp, err := client.Logical().Read(mount + "/issuer/" + issuer)
		if err != nil {
			t.Fatalf("failed to get chain for issuer (%v): %v", issuer, err)
		}

		rawCurrentChain := resp.Data["ca_chain"].([]interface{})
		var currentChain []string
		for _, entry := range rawCurrentChain {
			currentChain = append(currentChain, strings.TrimSpace(entry.(string)))
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
			parentCert := c.ToCertificate(t, thisCertPem)

			// Iterate backwards; prefer the most recent cert to the older
			// certs.
			foundCert := false
			for otherIndex := thisIndex - 1; otherIndex >= 0; otherIndex-- {
				otherCertPem := currentChain[otherIndex]
				childCert := c.ToCertificate(t, otherCertPem)

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

type CBTestStep interface {
	Run(t *testing.T, client *api.Client, mount string, knownKeys map[string]string, knownCerts map[string]string)
}

type CBTestScenario struct {
	Steps []CBTestStep
}

func Test_CAChainBuilding(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"pki": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client

	testCases := []CBTestScenario{
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
			},
		},
	}

	for testIndex, testCase := range testCases {
		mount := fmt.Sprintf("pki-test-%v", testIndex)
		mountPKIEndpoint(t, client, mount)
		knownKeys := make(map[string]string)
		knownCerts := make(map[string]string)
		for stepIndex, testStep := range testCase.Steps {
			t.Logf("Running %v / %v", testIndex, stepIndex)
			testStep.Run(t, client, mount, knownKeys, knownCerts)
		}

	}
}
