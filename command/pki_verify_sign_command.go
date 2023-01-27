package command

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/ryanuber/columnize"
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

  Verifies whether the listed issuer has signed the listed issued certificate.

  POSSIBLE-ISSUER and POSSIBLE-ISSUED are the fully name-spaced path to
  an issuer certificate, for instance: 'ns1/mount1/issuer/issuerName/json'.

  Returns five fields of information:

    - signature_match: was the key of the issuer used to sign the issued.
    - path_match: the possible issuer appears in the valid certificate chain
	  of the issued.
    - key_id_match: does the key-id of the issuer match the key_id of the
	  subject.
    - subject_match: does the subject name of the issuer match the issuer
	  subject of the issued.
    - trust_match: if someone trusted the parent issuer, is the chain
	  provided sufficient to trust the child issued.

` + c.Flags().Help()
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

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to obtain client: %s", err))
		return 1
	}

	err, results := verifySignBetween(client, issuer, issued)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to run verification: %v", err))
		return pkiRetUsage
	}

	c.outputResults(results, issuer, issued)

	return 0
}

func verifySignBetween(client *api.Client, issuerPath string, issuedPath string) (error, map[string]bool) {
	// Note that this eats warnings

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
	issuerKeyId := issuerCertBundle.Certificate.SubjectKeyId

	// Fetch and Parse the Potential Issued Cert
	issuedCertResp, err := client.Logical().Read(issuedPath)
	if err != nil {
		return fmt.Errorf("error: unable to fetch issuer %v: %w", issuerPath, err), nil
	}
	if len(issuedPath) <= 2 {
		return fmt.Errorf("%v", issuedPath), nil
	}
	caChainRaw := issuedCertResp.Data["ca_chain"]
	if caChainRaw == nil {
		return fmt.Errorf("no ca_chain information on %v", issuedPath), nil
	}
	caChainCast := caChainRaw.([]interface{})
	caChain := make([]string, len(caChainCast))
	for i, cert := range caChainCast {
		caChain[i] = cert.(string)
	}
	issuedCertPem := issuedCertResp.Data["certificate"].(string)
	issuedCertBundle, err := certutil.ParsePEMBundle(issuedCertPem)
	if err != nil {
		return err, nil
	}
	parentKeyId := issuedCertBundle.Certificate.AuthorityKeyId

	// Check the Chain-Match
	rootCertPool := x509.NewCertPool()
	rootCertPool.AddCert(issuerCertBundle.Certificate)
	checkTrustPathOptions := x509.VerifyOptions{
		Roots: rootCertPool,
	}
	trust := false
	trusts, err := issuedCertBundle.Certificate.Verify(checkTrustPathOptions)
	if err != nil && !strings.Contains(err.Error(), "certificate signed by unknown authority") {
		return err, nil
	} else if err == nil {
		for _, chain := range trusts {
			// Output of this Should Only Have One Trust with Chain of Length Two (Child followed by Parent)
			for _, cert := range chain {
				if issuedCertBundle.Certificate.Equal(cert) {
					trust = true
					break
				}
			}
		}
	}

	pathMatch := false
	for _, cert := range caChain {
		if strings.TrimSpace(cert) == strings.TrimSpace(issuerCertPem) { // TODO: Decode into ASN1 to Check
			pathMatch = true
			break
		}
	}

	signatureMatch := false
	err = issuedCertBundle.Certificate.CheckSignatureFrom(issuerCertBundle.Certificate)
	if err == nil {
		signatureMatch = true
	}

	result := map[string]bool{
		// This comparison isn't strictly correct, despite a standard ordering these are sets
		"subject_match":   bytes.Equal(issuerCertBundle.Certificate.RawSubject, issuedCertBundle.Certificate.RawIssuer),
		"path_match":      pathMatch,
		"trust_match":     trust, // TODO: Refactor into a reasonable function
		"key_id_match":    bytes.Equal(parentKeyId, issuerKeyId),
		"signature_match": signatureMatch,
	}

	return nil, result
}

func (c *PKIVerifySignCommand) outputResults(results map[string]bool, potentialParent, potentialChild string) error {
	switch Format(c.UI) {
	case "", "table":
		return c.outputResultsTable(results, potentialParent, potentialChild)
	case "json":
		return c.outputResultsJSON(results)
	case "yaml":
		return c.outputResultsYAML(results)
	default:
		return fmt.Errorf("unknown output format: %v", Format(c.UI))
	}
}

func (c *PKIVerifySignCommand) outputResultsTable(results map[string]bool, potentialParent, potentialChild string) error {
	c.UI.Output("issuer:" + potentialParent)
	c.UI.Output("issued:" + potentialChild + "\n")
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
