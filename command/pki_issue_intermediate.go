package command

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/posener/complete"
	"io"
	"os"
	paths "path"
	"strings"
)

type PKIIssueCACommand struct {
	*BaseCommand

	flagConfig          string
	flagReturnIndicator string
	flagDefaultDisabled bool
	flagList            bool

	flagKeyStorageSource string
	flagNewIssuerName    string
}

func (c *PKIIssueCACommand) Synopsis() string {
	return "Given a Parent Certificate, and a List of Generation Parameters, Creates an Issue on a Specified Moount"
}

func (c *PKIIssueCACommand) Help() string {
	helpText := `
Usage: vault pki issue PARENT CHILD_MOUNT options
`
	return strings.TrimSpace(helpText)
}

func (c *PKIIssueCACommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:       "type",
		Target:     &c.flagKeyStorageSource,
		Default:    "internal",
		EnvVar:     "",
		Usage:      `Options are “existing” - to use an existing key inside vault, “internal” - to generate a new key inside vault, or “kms” - to link to an external key.  Exported keys are not available through this API.`,
		Completion: complete.PredictSet("internal", "existing", "kms"),
	})

	f.StringVar(&StringVar{
		Name:    "issuer_name",
		Target:  &c.flagNewIssuerName,
		Default: "",
		EnvVar:  "",
		Usage:   `If present, the newly created issuer will be given this name`,
	})

	return set
}

func (c *PKIIssueCACommand) Run(args []string) int {
	f := c.Flags()
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	if len(args) < 3 {
		c.UI.Error("Not enough arguments expected parent issuer and child-mount location and some key_value argument")
		return 1
	}

	parentMountIssuer := sanitizePath(args[0])            // /pki/issuer/default
	_, parentIssuerName := paths.Split(parentMountIssuer) // TODO: Use this in order to name the issuer when we import it into the intermediate mount

	client, err := c.Client()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to obtain client: %w", err))
		return 1
	}

	_, err = client.Logical().Read(parentMountIssuer)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Unable to access parent issuer %v: %v", parentMountIssuer, err))
	}

	intermediateMount := sanitizePath(args[1])

	// Pull our fake stdin if needed
	stdin := (io.Reader)(os.Stdin)

	data, err := parseArgsData(stdin, args[2:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	csrResp, err := client.Logical().Write(intermediateMount+"/intermediate/generate/"+c.flagKeyStorageSource, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failled to Generate Intermediate CSR on %v: %v", intermediateMount, err))
		return 1
	}

	csrPem := csrResp.Data["csr"].(string)
	data["csr"] = csrPem

	rootResp, err := client.Logical().Write(parentMountIssuer+"/sign-intermediate", data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error Signing Intermiate On %v", err))
		return 1
	}

	// Import Certificate
	certificate := rootResp.Data["certificate"].(string)
	err = importIssuerWithName(client, intermediateMount, certificate, c.flagNewIssuerName)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error Importing Into %v Newly Created Issuer %v: %v", intermediateMount, certificate, err))
		return 1
	}

	// Import Issuing Certificate
	issuingCa := rootResp.Data["issuing_ca"].(string)
	err = importIssuerWithName(client, intermediateMount, issuingCa, parentIssuerName)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error Importing Into %v Newly Created Issuer %v: %v", intermediateMount, certificate, err))
		return 1
	}

	// Import CA_Chain (just in case there's more information)
	caChain := rootResp.Data["ca_chain"].([]interface{})
	pemBundle := ""
	for _, cert := range caChain {
		pemBundle += cert.(string) + "\n"
	}
	importData := map[string]interface{}{
		"pem_bundle": pemBundle,
	}
	finalResp, err := client.Logical().Write(intermediateMount+"/issuers/import/cert", importData)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error Importing CaChain into %v: %v", intermediateMount, err))
		return 1
	}

	c.UI.Output(fmt.Sprintf("%v", finalResp))

	return 0
}

func importIssuerWithName(client *api.Client, mount string, bundle string, name string) error {
	importData := map[string]interface{}{
		"pem_bundle": bundle,
	}
	writeResp, err := client.Logical().Write(mount+"/issuers/import/cert", importData)
	if err != nil {
		return err
	}
	mapping := writeResp.Data["mapping"].(map[string]interface{})
	if len(mapping) > 1 {
		return fmt.Errorf("multiple issuers returned, while expected one, got %v", writeResp)
	}
	if name != "" && name != "default" {
		issuerUUID := ""
		for issuerId, _ := range mapping {
			issuerUUID = issuerId
		}
		nameReq := map[string]interface{}{
			"issuer_name": name,
		}
		ctx := context.Background()
		_, err = client.Logical().JSONMergePatch(ctx, mount+"/issuer/"+issuerUUID, nameReq)
		if err != nil {
			return fmt.Errorf("error naming issuer %v to %v: %v", issuerUUID, name, err)
		}
	}
	return nil
}
