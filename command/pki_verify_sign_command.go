package command

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/ryanuber/columnize"
	"strconv"
	"strings"
)

type PKIVerifySignCommand struct {
	*BaseCommand

	flagConfig          string
	flagReturnIndicator string
	flagDefaultDisabled bool
	flagList            bool
}

func (c *PKIVerifySignCommand) Synopsis() string {
	return "Check Whether One Certificate Validates Another Specified Certificate"
}

func (c *PKIVerifySignCommand) Help() string {
	helpText := `
Usage: vault pki verify-sign POSSIBLE-ISSUER POSSIBLE-ISSUED
Returns four fields of information:
- signature_match: was the key of the issuer used to sign the issued
- path_match: the possible issuer appears in the valid certificate chain of the issued
- key_id_match: does the key-id of the issuer match the key_id of the subject
- subject_match: does the subject name of the issuer match the issuer subject of the issued
`
	return strings.TrimSpace(helpText)
}

func (c *PKIVerifySignCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	return set
}

func (c *PKIVerifySignCommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	if len(args) < 2 {
		if len(args) == 0 {
			c.UI.Error("Not enough arguments (expected potential issuer and issued, got nothing)")
		} else {
			c.UI.Error("Not enough arguments (expected both potential issuer and issued, got only one)")
		}
		return 1
	} else if len(args) > 2 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected only potential issuer and issued, got %d arguments)", len(args)))
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				c.UI.Warn(fmt.Sprintf("Options (%v) must be specified before positional arguments (%v)", arg, args[0]))
				break
			}
		}
		return 1
	}

	issuer := sanitizePath(args[0])
	issued := sanitizePath(args[1])

	err, results := verifySignBetween(c.client, issuer, issued)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to run verification: %v", err))
		return pkiRetUsage
	}

	c.outputResults(results)

	return 0
}

func verifySignBetween(client *api.Client, issuerPath string, issuedPath string) (error, map[string]bool) {
	// TODO: Stop Eating Warnings Here

	// Fetch and Parse the Potential Issuer:
	issuerResp, err := client.Logical().Read(issuerPath)
	if err != nil {
		return fmt.Errorf("error: unable to fetch issuer %v: %w", issuerPath, err), nil
	}
	issuerCertPem := issuerResp.Data["certificate"].(string)
	issuerCertBundle, err := certutil.ParsePEMBundle(issuerCertPem)
	if err != nil {
		return err, nil
	}
	issuerName := issuerCertBundle.Certificate.Subject
	issuerKeyId := issuerCertBundle.Certificate.SubjectKeyId

	// Fetch and Parse the Potential Issued Cert
	issuedCertResp, err := client.Logical().Read(issuedPath)
	if err != nil {
		return fmt.Errorf("error: unable to fetch issuer %v: %w", issuerPath, err), nil
	}
	caChain := issuedCertResp.Data["ca_chain"].(string)
	issuedCertPem := issuedCertResp.Data["certificate"].(string)
	issuedCertBundle, err := certutil.ParsePEMBundle(issuedCertPem)
	if err != nil {
		return err, nil
	}
	parentIssuerName := issuerCertBundle.Certificate.Issuer
	parentKeyId := issuerCertBundle.Certificate.AuthorityKeyId

	// Check the Chain-Match
	rootCertPool := x509.NewCertPool() // TODO: Check if it matters to use Root Only Here (There's also an Intermediate Cert Pool)
	rootCertPool.AddCert(issuerCertBundle.Certificate)
	checkTrustPathOptions := x509.VerifyOptions{
		Roots: rootCertPool,
	}
	trusts, err := issuedCertBundle.Certificate.Verify(checkTrustPathOptions)
	trust := false
	for _, chain := range trusts {
		for _, cert := range chain {
			if issuedCertBundle.Certificate.Equal(cert) {
				trust = true
			}
		}
	}

	result := map[string]bool{
		"subject_match":   parentIssuerName.String() == issuerName.String(), // TODO: No Equals Defined on Name, Check This
		"path_match":      strings.Contains(caChain, issuerCertPem),         // TODO: Check Trimming
		"trust_match":     trust,                                            // TODO: Refactor into a reasonable function
		"key_id_match":    isKeyIDEqual(parentKeyId, issuerKeyId),
		"signature_match": false, // TODO: Checking the Signature Has to Be Done Per-Algorithm
	}

	return nil, result
}

func isKeyIDEqual(first, second []byte) bool {
	// TODO: Check if Trimming Makes Sense Here - Eg. Could this Be Padded?
	// TODO: There has to be a library that does this
	if len(first) != len(second) {
		return false
	}
	for i, byteOfFirst := range first {
		if byteOfFirst != second[i] {
			return false
		}
	}
	return true
}

func (c *PKIVerifySignCommand) outputResults(results map[string]bool) error {
	switch Format(c.UI) {
	case "", "table":
		return c.outputResultsTable(results)
	case "json":
		return c.outputResultsJSON(results)
	case "yaml":
		return c.outputResultsYAML(results)
	default:
		return fmt.Errorf("unknown output format: %v", Format(c.UI))
	}
}

func (c *PKIVerifySignCommand) outputResultsTable(results map[string]bool) error {
	data := []string{"field" + hopeDelim + "value"}
	for field, finding := range results {
		row := field + hopeDelim + strconv.FormatBool(finding)
		data = append(data, row)
	}
	c.UI.Output(tableOutput(data, &columnize.Config{
		Delim: hopeDelim,
	}))
	c.UI.Output("\n")

	return nil
}

func (c *PKIVerifySignCommand) outputResultsJSON(results map[string]bool) error {
	bytes, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	c.UI.Output(string(bytes))
	return nil
}

func (c *PKIVerifySignCommand) outputResultsYAML(results map[string]bool) error {
	bytes, err := yaml.Marshal(results)
	if err != nil {
		return err
	}

	c.UI.Output(string(bytes))
	return nil
}
