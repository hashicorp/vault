package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/pgpkeys"
)

// InitCommand is a Command that initializes a new Vault server.
type InitCommand struct {
	Meta
}

func (c *InitCommand) Run(args []string) int {
	var threshold, shares int
	var pgpKeys pgpkeys.PubKeyFilesFlag
	var idempotent bool
	flags := c.Meta.FlagSet("init", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	flags.IntVar(&shares, "key-shares", 5, "")
	flags.IntVar(&threshold, "key-threshold", 3, "")
	flags.Var(&pgpKeys, "pgp-keys", "")
	flags.BoolVar(&idempotent, "idempotent", false, "")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 1
	}

	resp, err := client.Sys().Init(&api.InitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
		PGPKeys:         pgpKeys,
		Idempotent:      idempotent,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing Vault: %s", err))
		return 1
	}

	if idempotent {
		switch {
		case (resp.Keys == nil || len(resp.Keys) == 0) && resp.RootToken == "":
			c.Ui.Output("Vault is already initialized")
			return 0
		case (resp.Keys == nil || len(resp.Keys) == 0) && resp.RootToken != "":
			c.Ui.Output("Response from Vault could not be understood")
			return 1
		case resp.Keys != nil && len(resp.Keys) > 0 && resp.RootToken == "":
			c.Ui.Output("Response from Vault could not be understood")
			return 1
		default:
			// Do nothing; we have both keys and a root token back so this is a
			// fresh init
		}
	}

	for i, key := range resp.Keys {
		c.Ui.Output(fmt.Sprintf("Key %d: %s", i+1, key))
	}

	c.Ui.Output(fmt.Sprintf("Initial Root Token: %s", resp.RootToken))

	c.Ui.Output(fmt.Sprintf(
		"\n"+
			"Vault initialized with %d keys and a key threshold of %d. Please\n"+
			"securely distribute the above keys. When the Vault is re-sealed,\n"+
			"restarted, or stopped, you must provide at least %d of these keys\n"+
			"to unseal it again.\n\n"+
			"Vault does not store the master key. Without at least %d keys,\n"+
			"your Vault will remain permanently sealed.",
		shares,
		threshold,
		threshold,
		threshold,
	))

	return 0
}

func (c *InitCommand) Synopsis() string {
	return "Initialize a new Vault server"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: vault init [options]

  Initialize a new Vault server.

  This command connects to a Vault server and initializes it for the
  first time. This sets up the initial set of master keys and sets up the
  backend data store structure.

  This command can't be called on an already-initialized Vault.

General Options:

  ` + generalOptionsUsage() + `

Init Options:

  -key-shares=5           The number of key shares to split the master key
                          into.

  -key-threshold=3        The number of key shares required to reconstruct
                          the master key.

  -pgp-keys               If provided, must be a comma-separated list of
                          files on disk containing binary- or base64-format
                          public PGP keys, or Keybase usernames specified as
                          "keybase:<username>". The number of given entries
                          must match 'key-shares'. The output unseal keys will
                          encrypted and hex-encoded, in order, with the given
                          public keys.  If you want to use them with the 'vault
                          unseal' command, you will need to hex decode and
                          decrypt; this will be the plaintext unseal key.

  -idempotent             If set, the command will be idempotent; issuing an
                          init on an already-initialized Vault will return
                          success rather than failure so long as the key shares,
                          threshold, and any PGP keys are the same.
`
	return strings.TrimSpace(helpText)
}
