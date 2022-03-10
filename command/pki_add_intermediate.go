package command

import (
	"fmt"
	"github.com/hashicorp/vault/command/pkicli"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"strings"
)

var (
	_ cli.Command             = (*PKIAddIntermediateCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIAddIntermediateCommand)(nil)
)

type PKIAddIntermediateCommand struct {
	*BaseCommand

	flagMount string
	flagRootMount string
	flagCommonName string
}

func (c *PKIAddIntermediateCommand) Synopsis() string {
	return "Generate intermediate certificate"
}

func (c *PKIAddIntermediateCommand) Help() string {
	helpText := `
Usage: vault pki add-intermediate [options] ROOT_MOUNT PATH COMMON_NAME [K=V]

  Configures an intermediate mount and generate the intermediate certificate.
  The intermediate certificate is the one from which all leaf certificates will be generated.
  This intermediate will be signed by the root. Other parameters can be specified in
  K=V format (mirroring vault write).

  Configure an intermediate mount at path pki-int with a specific ttl:
      $ vault pki add-intermediate pki pki-int example.com ttl=48000h

  Configure an intermediate mount at path pki-int with CSR type and common name:
      $ vault pki add-intermediate pki pki-int example.com csr=@example.csr format=pem
  

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddIntermediateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:    "mount",
		Target:  &c.flagMount,
		Default: "pki",
		EnvVar:  "",
		Usage:   "PKI intermediate mount",
	})

	f.StringVar(&StringVar{
		Name:    "root-mount",
		Target:  &c.flagRootMount,
		Default: "pki",
		EnvVar:  "",
		Usage:   "PKI root mount",
	})

	f.StringVar(&StringVar{
		Name:    "common_name",
		Target:  &c.flagCommonName,
		EnvVar:  "",
		Usage:   "Common name",
	})

	return set
}

func (c *PKIAddIntermediateCommand) AutocompleteArgs() complete.Predictor {
	// Return an anything predictor here, similar to `vault write`. We
	// don't know what values are valid for the role and/or common names.
	return complete.PredictAnything
}

func (c *PKIAddIntermediateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIAddIntermediateCommand) Run(args []string) int {

	var mountPath string
	var commonName string
	var rootMounthPath string // optional
	var parameters map[string]interface{}

	ops := pkicli.NewOperations(c.client)
	ops.CreateIntermediate(rootMounthPath, mountPath, /* commonName, */ parameters)

	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) < 3 {
		c.UI.Error(fmt.Sprintf("Not enough argumentssssss (expected 3+, got %d)", len(args)))
		return 1
	}
	//TODO - validate that root certificate exists and is valid, check if path is already in use.
	//Step 1: Mount backend for root and generate root certificate
	root_mount := sanitizePath(args[0])
	root_path := sanitizePath(fmt.Sprintf("sys/mounts/%s", root_mount))
	root_data := map[string]interface{}{
		"path": sanitizePath(root_path),
		"description": root_mount + " root CA",
		"type": "pki",
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	_, err = client.Logical().Write(root_path, root_data)

	if err != nil {
		c.UI.Error(err.Error())
		return 3
	}

	// Step 2: Generate root certificate
	root_path = sanitizePath(fmt.Sprintf("%s/root/generate/internal", root_mount))
	secret, err := client.Logical().Write(root_path, root_data)

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error generating root certificate: %v", err))
		return 4
	}

    // Step 3: Mount the backend and configure intermediate CA
	mount := sanitizePath(args[1])
	commonName := args[2]
	path := sanitizePath(fmt.Sprintf("sys/mounts/%s", mount))

	data, err := parseArgsData(nil, args[3:])
	if err != nil {
		c.UI.Error(fmt.Sprintf("Failed to parse K=V data: %s", err))
		return 1
	}

	if data == nil {
		data = make(map[string]interface{})
	}
    data["path"] = sanitizePath(path)
    data["description"] = mount + " intermediate CA"
    data["type"] = "pki"
	data["common_name"] = commonName

	_, err = client.Logical().Write(path, data)

	if err != nil {
		c.UI.Error(err.Error())
		return 3
	}
	c.UI.Info(fmt.Sprintf("Successfully mounted backend for intermediate CA"))

	// Step 4: Generate intermediate certificate signing request
	path = sanitizePath(fmt.Sprintf("%s/intermediate/generate/internal", mount))
	//c.UI.Error(fmt.Sprintf("data for generating csr: %v", data))
	secret, err = client.Logical().Write(path, data)

	if err != nil {
		c.UI.Error(fmt.Sprintf("Error generating intermediate certificate signing request: %v", err))
		return 4
	}
	csr := secret.Data["csr"]
	c.UI.Info(fmt.Sprintf("Successfully generated intermediate certificate signing request: %v", csr))
	//fmt.Println(csr)

	// Step 5: Sign the intermediate CA
	path = sanitizePath(fmt.Sprintf("pki/root/sign-intermediate"))
	data["csr"] = csr
	data["format"] = "pem_bundle"
	secret, err = client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error signing intermediate CA: %v", err))
		return 5
	}
	cert := secret.Data["certificate"]
	//fmt.Println(secret.Data)

	// Set intermediate CA's signing certificate
	path = sanitizePath(fmt.Sprintf("%s/intermediate/set-signed", mount))
	data["certificate"] = cert
	secret, err = client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting intermediate CA's signing certificate: %v", err))
		return 5
	}
	c.UI.Info(fmt.Sprintf("Successfully set intermediate CA's signing certificate: %v", cert))

	// STep 6: Configure a role
	role := "example-dot-com" // Should we get the role name from parameters?
	path = sanitizePath(fmt.Sprintf("%s/roles/%s", mount, role))
	data["allowed_domains"] = commonName
	data["allowed_subdomains"] = true
	_, err = client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error configuring role: %v", err))
		return 5
	}
	c.UI.Info(fmt.Sprintf("Successfully configured role:%s", role))

	//Step 7: Issue certificates
	path = sanitizePath(fmt.Sprintf("%s/issue/%s", mount, role))
	_, err = client.Logical().Write(path, data)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error issuing leaf certificate: %v", err))
		return 5
	}

	return 0
}
