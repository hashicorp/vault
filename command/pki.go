package command

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"

	"github.com/mitchellh/cli"
)

var _ cli.Command = (*PKICommand)(nil)

type PKICommand struct {
	*BaseCommand
}

func (c *PKICommand) Synopsis() string {
	return "Interact with PKI Secret Engines"
}

func (c *PKICommand) Help() string {
	helpText := `
Usage: vault pki <subcommand> [options] [args]

  This command groups subcommands for interacting with Vault's PKI Secrets
  Engine. Operators can manage PKI mounts and roles.

  To test role based issuance:

       $ vault pki role-test -mount=pki-int server-role example.com

  Please see the individual subcommand help for detailed usage information.
`

	return strings.TrimSpace(helpText)
}

func (c *PKICommand) Run(args []string) int {
	return cli.RunResultHelp
}

func (c *PKICommand) testCreateRoot() int {
	params := pkiCreateRootParameters{
		path: "pki-root",
		maxLeaseTTL: "24h",
	}
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
	}
	ops := &pkiOps{client: client}
	_, err = ops.createRoot(params)
	if err != nil {
		c.UI.Error(err.Error())
	}

	return 0
}

type pkiOperations interface {
	createRoot(params pkiCreateRootParameters) (*pkiCreateRootResponse, error)
	createIntermediate(params pkiCreateIntermediateParameters) (*pkiCreateIntermediateResponse, error)
}

type pkiCreateRootParameters struct {
	path        string
	maxLeaseTTL string
	commonName  string
	ttl         string
	// etc.
}

type pkiCreateRootResponse struct {
	certificate string
}

type pkiCreateIntermediateParameters struct {
	path        string
	maxLeaseTTL string
	commonName  string
	ttl         string
	// etc.
}

type pkiCreateIntermediateResponse struct {
	// TODO
}

var _ pkiOperations = (*pkiOps)(nil)

type pkiOps struct {
	client *api.Client
}

func (p pkiOps) createRoot(params pkiCreateRootParameters) (*pkiCreateRootResponse, error) {

	err := p.secretsEnable(params.path, params.path + " root CA", params.maxLeaseTTL)
	if err != nil {
		return nil, err
	}

	s, err := p.rootGenerate(params)
	cert := s.Data["certificate"]
	if err != nil {
		return nil, err
	}

	err = p.configUrls(params)
	if err != nil {
		// FIXME(victorr): should not really return nil here
		return nil, err
	}

	//fmt.Println(cert)
	r := &pkiCreateRootResponse{
		certificate: cert.(string),
	}
	return r, nil
}

func (p pkiOps) createIntermediate(params pkiCreateIntermediateParameters) (*pkiCreateIntermediateResponse, error) {
	panic("implement me")
}

func (p pkiOps) secretsEnable(mountPath string, desrcription string, maxLeaseTTL string) error {
	// https://www.vaultproject.io/api-docs/system/mounts#enable-secrets-engine
	data := map[string]interface{}{
		"path": sanitizePath(mountPath),
		"description": desrcription,
		"type": "pki",
	}
	if maxLeaseTTL != "" {
		data["config"] = map[string]interface{} {
			"max_lease_ttl": maxLeaseTTL,
		}
	}
	path := sanitizePath(fmt.Sprintf("sys/mounts/%s", mountPath))
	_, err := p.client.Logical().Write(path, data)

	return err
}

// Returns a secret with keys: certificate, expiration (number), issuing_ca, serial_number.
func (p pkiOps) rootGenerate(params pkiCreateRootParameters) (*api.Secret, error) {
	// https://www.vaultproject.io/api/secret/pki#generate-root

	data := map[string]interface{}{
		"common_name": params.commonName,
		"ttl": params.ttl,
	}
	path := sanitizePath(fmt.Sprintf("%s/root/generate/internal", params.path))

	return p.client.Logical().Write(path, data)
}

func (p pkiOps) configUrls(params pkiCreateRootParameters) error {
	// https://www.vaultproject.io/api/secret/pki#set-urls

	// TODO
	return nil
}