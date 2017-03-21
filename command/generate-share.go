package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/meta"
)

// GenerateShareCommand is a Command that generates a new master key share.
type GenerateShareCommand struct {
	meta.Meta

	// Key can be used to pre-seed the key. If it is set, it will not
	// be asked with the `password` helper.
	Key string
}

func (c *GenerateShareCommand) Run(args []string) int {
	var init, cancel, status bool
	var pgpKey string
	var pgpKeyArr pgpkeys.PubKeyFilesFlag
	flags := c.Meta.FlagSet("generate-share", meta.FlagSetDefault)
	flags.BoolVar(&init, "init", false, "")
	flags.BoolVar(&cancel, "cancel", false, "")
	flags.BoolVar(&status, "status", false, "")
	flags.Var(&pgpKeyArr, "pgp-key", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	// Check if the root generation is started
	shareGenerationStatus, err := client.Sys().GenerateShareStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading share generation status: %s", err))
		return 1
	}

	// If we are initing, or if we are not started but are not running a
	// special function, check pgpkey
	checkPgp := false
	switch {
	case init:
		checkPgp = true
	case cancel:
	case status:
	case shareGenerationStatus.Started:
	default:
		checkPgp = true
	}
	if checkPgp {
		switch {
		case (pgpKeyArr == nil || len(pgpKeyArr) == 0):
			c.Ui.Error(c.Help())
			return 1
		case pgpKeyArr != nil:
			if len(pgpKeyArr[0]) == 0 {
				c.Ui.Error("Got an empty PGP key")
				return 1
			}
			if len(pgpKeyArr) != 1 {
				c.Ui.Error("Could not parse PGP key")
				return 1
			}
			pgpKey = pgpKeyArr[0]
		default:
			panic("unreachable case")
		}
	}

	// Check if we are running doing any restricted variants
	switch {
	case init:
		return c.initGenerateShare(client, pgpKey)
	case cancel:
		return c.cancelGenerateShare(client)
	case status:
		return c.shareGenerationStatus(client)
	}

	// Start the share generation process if not started
	if !shareGenerationStatus.Started {
		shareGenerationStatus, err = client.Sys().GenerateShareInit(pgpKey)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error initializing root generation: %s", err))
			return 1
		}
	}

	// Get the unseal key
	args = flags.Args()
	key := c.Key
	if len(args) > 0 {
		key = args[0]
	}
	if key == "" {
		fmt.Printf("Key (will be hidden): ")
		key, err = password.Read(os.Stdin)
		fmt.Printf("\n")
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error attempting to ask for password. The raw error message\n"+
					"is shown below, but the most common reason for this error is\n"+
					"that you attempted to pipe a value into unseal or you're\n"+
					"executing `vault generate-share` from outside of a terminal.\n\n"+
					"You should use `vault generate-share` from a terminal for maximum\n"+
					"security. If this isn't an option, the unseal key can be passed\n"+
					"in using the first parameter.\n\n"+
					"Raw error: %s", err))
			return 1
		}
	}

	// Provide the key, this may potentially complete the update
	statusResp, err := client.Sys().GenerateShareUpdate(strings.TrimSpace(key))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error attempting generate-share update: %s", err))
		return 1
	}

	c.dumpStatus(statusResp)

	return 0
}

// initGenerateShare is used to start the generation process
func (c *GenerateShareCommand) initGenerateShare(client *api.Client, pgpKey string) int {
	// Start the rekey
	status, err := client.Sys().GenerateShareInit(pgpKey)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing share generation: %s", err))
		return 1
	}

	c.dumpStatus(status)

	return 0
}

// cancelGenerateShare is used to abort the generation process
func (c *GenerateShareCommand) cancelGenerateShare(client *api.Client) int {
	err := client.Sys().GenerateShareCancel()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to cancel share generation: %s", err))
		return 1
	}
	c.Ui.Output("Share generation canceled.")
	return 0
}

// shareGenerationStatus is used just to fetch and dump the status
func (c *GenerateShareCommand) shareGenerationStatus(client *api.Client) int {
	// Check the status
	status, err := client.Sys().GenerateShareStatus()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading share generation status: %s", err))
		return 1
	}

	c.dumpStatus(status)

	return 0
}

// dumpStatus dumps the status to output
func (c *GenerateShareCommand) dumpStatus(status *api.GenerateShareStatusResponse) {
	// Dump the status
	statString := fmt.Sprintf(
		"Started: %v\n"+
			"Generate Share Progress: %d\n"+
			"Required Keys: %d\n"+
			"Complete: %t",
		status.Started,
		status.Progress,
		status.Required,
		status.Complete,
	)
	if len(status.PGPFingerprint) > 0 {
		statString = fmt.Sprintf("%s\nPGP Fingerprint: %s", statString, status.PGPFingerprint)
	}
	if len(status.Key) > 0 {
		statString = fmt.Sprintf("%s\nShare: %s", statString, status.Key)
	}

	c.Ui.Output(statString)
}

func (c *GenerateShareCommand) Synopsis() string {
	return "Generates a new master key share"
}

func (c *GenerateShareCommand) Help() string {
	helpText := `
Usage: vault generate-share [options] [key]

  'generate-share' is used to create a new master key share.

  Share generation can only be done when the Vault is already unsealed. The
  operation is done online, but requires that a threshold of the current unseal
  keys be provided.

  The following must be provided at attempt initialization time:

  A file containing a PGP key (binary or base64-encoded) or a Keybase.io
  username in the format of "keybase:<username>" in the '-pgp-key' flag. The
  final share value will be encrypted with this public key and base64-encoded.

General Options:
` + meta.GeneralOptionsUsage() + `
Generate Share Options:

  -init                   Initialize the share generation attempt. This can only
                          be done if no generation is already initiated.

  -cancel                 Reset the share generation process by throwing away
                          prior unseal keys and the configuration.

  -status                 Prints the status of the current attempt. This can be
                          used to see the status without attempting to provide
                          an unseal key.

  -pgp-key                A file on disk containing a binary- or base64-format
                          public PGP key, or a Keybase username specified as
                          "keybase:<username>". The output root token will be
                          encrypted and base64-encoded, in order, with the given
                          public key.
`
	return strings.TrimSpace(helpText)
}
