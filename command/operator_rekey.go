package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var _ cli.Command = (*OperatorRekeyCommand)(nil)
var _ cli.CommandAutocomplete = (*OperatorRekeyCommand)(nil)

type OperatorRekeyCommand struct {
	*BaseCommand

	flagCancel       bool
	flagInit         bool
	flagKeyShares    int
	flagKeyThreshold int
	flagNonce        string
	flagPGPKeys      []string
	flagStatus       bool
	flagTarget       string
	flagVerify       bool

	// Backup options
	flagBackup         bool
	flagBackupDelete   bool
	flagBackupRetrieve bool

	testStdin io.Reader // for tests
}

func (c *OperatorRekeyCommand) Synopsis() string {
	return "Generates new unseal keys"
}

func (c *OperatorRekeyCommand) Help() string {
	helpText := `
Usage: vault operator rekey [options] [KEY]

  Generates a new set of unseal keys. This can optionally change the total
  number of key shares or the required threshold of those key shares to
  reconstruct the master key. This operation is zero downtime, but it requires
  the Vault is unsealed and a quorum of existing unseal keys are provided.

  An unseal key may be provided directly on the command line as an argument to
  the command. If key is specified as "-", the command will read from stdin. If
  a TTY is available, the command will prompt for text.

  Initialize a rekey:

      $ vault operator rekey \
          -init \
          -key-shares=15 \
          -key-threshold=9

  Rekey and encrypt the resulting unseal keys with PGP:

      $ vault operator rekey \
          -init \
          -key-shares=3 \
          -key-threshold=2 \
          -pgp-keys="keybase:hashicorp,keybase:jefferai,keybase:sethvargo"

  Store encrypted PGP keys in Vault's core:

      $ vault operator rekey \
          -init \
          -pgp-keys="..." \
          -backup

  Retrieve backed-up unseal keys:

      $ vault operator rekey -backup-retrieve

  Delete backed-up unseal keys:

      $ vault operator rekey -backup-delete

` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *OperatorRekeyCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Common Options")

	f.BoolVar(&BoolVar{
		Name:    "init",
		Target:  &c.flagInit,
		Default: false,
		Usage: "Initialize the rekeying operation. This can only be done if no " +
			"rekeying operation is in progress. Customize the new number of key " +
			"shares and key threshold using the -key-shares and -key-threshold " +
			"flags.",
	})

	f.BoolVar(&BoolVar{
		Name:    "cancel",
		Target:  &c.flagCancel,
		Default: false,
		Usage: "Reset the rekeying progress. This will discard any submitted " +
			"unseal keys or configuration.",
	})

	f.BoolVar(&BoolVar{
		Name:    "status",
		Target:  &c.flagStatus,
		Default: false,
		Usage: "Print the status of the current attempt without providing an " +
			"unseal key.",
	})

	f.IntVar(&IntVar{
		Name:       "key-shares",
		Aliases:    []string{"n"},
		Target:     &c.flagKeyShares,
		Default:    5,
		Completion: complete.PredictAnything,
		Usage: "Number of key shares to split the generated master key into. " +
			"This is the number of \"unseal keys\" to generate.",
	})

	f.IntVar(&IntVar{
		Name:       "key-threshold",
		Aliases:    []string{"t"},
		Target:     &c.flagKeyThreshold,
		Default:    3,
		Completion: complete.PredictAnything,
		Usage: "Number of key shares required to reconstruct the master key. " +
			"This must be less than or equal to -key-shares.",
	})

	f.StringVar(&StringVar{
		Name:       "nonce",
		Target:     &c.flagNonce,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "Nonce value provided at initialization. The same nonce value " +
			"must be provided with each unseal key.",
	})

	f.StringVar(&StringVar{
		Name:       "target",
		Target:     &c.flagTarget,
		Default:    "barrier",
		EnvVar:     "",
		Completion: complete.PredictSet("barrier", "recovery"),
		Usage: "Target for rekeying. \"recovery\" only applies when HSM support " +
			"is enabled.",
	})

	f.BoolVar(&BoolVar{
		Name:    "verify",
		Target:  &c.flagVerify,
		Default: false,
		Usage: "Indicates that the action (-status, -cancel, or providing a key " +
			"share) will be affecting verification for the current rekey " +
			"attempt.",
	})

	f.VarFlag(&VarFlag{
		Name:       "pgp-keys",
		Value:      (*pgpkeys.PubKeyFilesFlag)(&c.flagPGPKeys),
		Completion: complete.PredictAnything,
		Usage: "Comma-separated list of paths to files on disk containing " +
			"public GPG keys OR a comma-separated list of Keybase usernames using " +
			"the format \"keybase:<username>\". When supplied, the generated " +
			"unseal keys will be encrypted and base64-encoded in the order " +
			"specified in this list.",
	})

	f = set.NewFlagSet("Backup Options")

	f.BoolVar(&BoolVar{
		Name:    "backup",
		Target:  &c.flagBackup,
		Default: false,
		Usage: "Store a backup of the current PGP encrypted unseal keys in " +
			"Vault's core. The encrypted values can be recovered in the event of " +
			"failure or discarded after success. See the -backup-delete and " +
			"-backup-retrieve options for more information. This option only " +
			"applies when the existing unseal keys were PGP encrypted.",
	})

	f.BoolVar(&BoolVar{
		Name:    "backup-delete",
		Target:  &c.flagBackupDelete,
		Default: false,
		Usage:   "Delete any stored backup unseal keys.",
	})

	f.BoolVar(&BoolVar{
		Name:    "backup-retrieve",
		Target:  &c.flagBackupRetrieve,
		Default: false,
		Usage: "Retrieve the backed-up unseal keys. This option is only available " +
			"if the PGP keys were provided and the backup has not been deleted.",
	})

	return set
}

func (c *OperatorRekeyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *OperatorRekeyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OperatorRekeyCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0-1, got %d)", len(args)))
		return 1
	}

	// Create the client
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	switch {
	case c.flagBackupDelete:
		return c.backupDelete(client)
	case c.flagBackupRetrieve:
		return c.backupRetrieve(client)
	case c.flagCancel:
		return c.cancel(client)
	case c.flagInit:
		return c.init(client)
	case c.flagStatus:
		return c.status(client)
	default:
		// If there are no other flags, prompt for an unseal key.
		key := ""
		if len(args) > 0 {
			key = strings.TrimSpace(args[0])
		}
		return c.provide(client, key)
	}
}

// init starts the rekey process.
func (c *OperatorRekeyCommand) init(client *api.Client) int {
	// Handle the different API requests
	var fn func(*api.RekeyInitRequest) (*api.RekeyStatusResponse, error)
	switch strings.ToLower(strings.TrimSpace(c.flagTarget)) {
	case "barrier":
		fn = client.Sys().RekeyInit
	case "recovery", "hsm":
		fn = client.Sys().RekeyRecoveryKeyInit
	default:
		c.UI.Error(fmt.Sprintf("Unknown target: %s", c.flagTarget))
		return 1
	}

	// Make the request
	status, err := fn(&api.RekeyInitRequest{
		SecretShares:        c.flagKeyShares,
		SecretThreshold:     c.flagKeyThreshold,
		PGPKeys:             c.flagPGPKeys,
		Backup:              c.flagBackup,
		RequireVerification: c.flagVerify,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error initializing rekey: %s", err))
		return 2
	}

	// Print warnings about recovery, etc.
	if len(c.flagPGPKeys) == 0 {
		if Format(c.UI) == "table" {
			c.UI.Warn(wrapAtLength(
				"WARNING! If you lose the keys after they are returned, there is no " +
					"recovery. Consider canceling this operation and re-initializing " +
					"with the -pgp-keys flag to protect the returned unseal keys along " +
					"with -backup to allow recovery of the encrypted keys in case of " +
					"emergency. You can delete the stored keys later using the -delete " +
					"flag."))
			c.UI.Output("")
		}
	}
	if len(c.flagPGPKeys) > 0 && !c.flagBackup {
		if Format(c.UI) == "table" {
			c.UI.Warn(wrapAtLength(
				"WARNING! You are using PGP keys for encrypted the resulting unseal " +
					"keys, but you did not enable the option to backup the keys to " +
					"Vault's core. If you lose the encrypted keys after they are " +
					"returned, you will not be able to recover them. Consider canceling " +
					"this operation and re-running with -backup to allow recovery of the " +
					"encrypted unseal keys in case of emergency. You can delete the " +
					"stored keys later using the -delete flag."))
			c.UI.Output("")
		}
	}

	// Provide the current status
	return c.printStatus(status)
}

// cancel is used to abort the rekey process.
func (c *OperatorRekeyCommand) cancel(client *api.Client) int {
	// Handle the different API requests
	var fn func() error
	switch strings.ToLower(strings.TrimSpace(c.flagTarget)) {
	case "barrier":
		fn = client.Sys().RekeyCancel
		if c.flagVerify {
			fn = client.Sys().RekeyVerificationCancel
		}
	case "recovery", "hsm":
		fn = client.Sys().RekeyRecoveryKeyCancel
		if c.flagVerify {
			fn = client.Sys().RekeyRecoveryKeyVerificationCancel
		}

	default:
		c.UI.Error(fmt.Sprintf("Unknown target: %s", c.flagTarget))
		return 1
	}

	// Make the request
	if err := fn(); err != nil {
		c.UI.Error(fmt.Sprintf("Error canceling rekey: %s", err))
		return 2
	}

	c.UI.Output("Success! Canceled rekeying (if it was started)")
	return 0
}

// provide prompts the user for the seal key and posts it to the update root
// endpoint. If this is the last unseal, this function outputs it.
func (c *OperatorRekeyCommand) provide(client *api.Client, key string) int {
	var statusFn func() (interface{}, error)
	var updateFn func(string, string) (interface{}, error)

	switch strings.ToLower(strings.TrimSpace(c.flagTarget)) {
	case "barrier":
		statusFn = func() (interface{}, error) {
			return client.Sys().RekeyStatus()
		}
		updateFn = func(s1 string, s2 string) (interface{}, error) {
			return client.Sys().RekeyUpdate(s1, s2)
		}
		if c.flagVerify {
			statusFn = func() (interface{}, error) {
				return client.Sys().RekeyVerificationStatus()
			}
			updateFn = func(s1 string, s2 string) (interface{}, error) {
				return client.Sys().RekeyVerificationUpdate(s1, s2)
			}
		}
	case "recovery", "hsm":
		statusFn = func() (interface{}, error) {
			return client.Sys().RekeyRecoveryKeyStatus()
		}
		updateFn = func(s1 string, s2 string) (interface{}, error) {
			return client.Sys().RekeyRecoveryKeyUpdate(s1, s2)
		}
		if c.flagVerify {
			statusFn = func() (interface{}, error) {
				return client.Sys().RekeyRecoveryKeyVerificationStatus()
			}
			updateFn = func(s1 string, s2 string) (interface{}, error) {
				return client.Sys().RekeyRecoveryKeyVerificationUpdate(s1, s2)
			}
		}
	default:
		c.UI.Error(fmt.Sprintf("Unknown target: %s", c.flagTarget))
		return 1
	}

	status, err := statusFn()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error getting rekey status: %s", err))
		return 2
	}

	var started bool
	var nonce string

	switch status.(type) {
	case *api.RekeyStatusResponse:
		stat := status.(*api.RekeyStatusResponse)
		started = stat.Started
		nonce = stat.Nonce
	case *api.RekeyVerificationStatusResponse:
		stat := status.(*api.RekeyVerificationStatusResponse)
		started = stat.Started
		nonce = stat.Nonce
	default:
		c.UI.Error("Unknown status type")
		return 1
	}

	// Verify a root token generation is in progress. If there is not one in
	// progress, return an error instructing the user to start one.
	if !started {
		c.UI.Error(wrapAtLength(
			"No rekey is in progress. Start a rekey process by running " +
				"\"vault operator rekey -init\"."))
		return 1
	}

	switch key {
	case "-": // Read from stdin
		nonce = c.flagNonce

		// Pull our fake stdin if needed
		stdin := (io.Reader)(os.Stdin)
		if c.testStdin != nil {
			stdin = c.testStdin
		}

		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdin); err != nil {
			c.UI.Error(fmt.Sprintf("Failed to read from stdin: %s", err))
			return 1
		}

		key = buf.String()
	case "": // Prompt using the tty
		// Nonce value is not required if we are prompting via the terminal
		w := getWriterFromUI(c.UI)
		fmt.Fprintf(w, "Rekey operation nonce: %s\n", nonce)
		fmt.Fprintf(w, "Unseal Key (will be hidden): ")
		key, err = password.Read(os.Stdin)
		fmt.Fprintf(w, "\n")
		if err != nil {
			if err == password.ErrInterrupted {
				c.UI.Error("user canceled")
				return 1
			}

			c.UI.Error(wrapAtLength(fmt.Sprintf("An error occurred attempting to "+
				"ask for the unseal key. The raw error message is shown below, but "+
				"usually this is because you attempted to pipe a value into the "+
				"command or you are executing outside of a terminal (tty). If you "+
				"want to pipe the value, pass \"-\" as the argument to read from "+
				"stdin. The raw error was: %s", err)))
			return 1
		}
	default: // Supplied directly as an arg
		nonce = c.flagNonce
	}

	// Trim any whitespace from they key, especially since we might have
	// prompted the user for it.
	key = strings.TrimSpace(key)

	// Verify we have a nonce value
	if nonce == "" {
		c.UI.Error("Missing nonce value: specify it via the -nonce flag")
		return 1
	}

	// Provide the key, this may potentially complete the update
	resp, err := updateFn(key, nonce)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error posting unseal key: %s", err))
		return 2
	}

	var complete bool
	var mightContainUnsealKeys bool

	switch resp.(type) {
	case *api.RekeyUpdateResponse:
		complete = resp.(*api.RekeyUpdateResponse).Complete
		mightContainUnsealKeys = true
	case *api.RekeyVerificationUpdateResponse:
		complete = resp.(*api.RekeyVerificationUpdateResponse).Complete
	default:
		c.UI.Error("Unknown update response type")
		return 1
	}

	if !complete {
		return c.status(client)
	}

	if mightContainUnsealKeys {
		return c.printUnsealKeys(client, status.(*api.RekeyStatusResponse),
			resp.(*api.RekeyUpdateResponse))
	}

	c.UI.Output(wrapAtLength("Rekey verification successful. The rekey operation is complete and the new keys are now active."))
	return 0
}

// status is used just to fetch and dump the status.
func (c *OperatorRekeyCommand) status(client *api.Client) int {
	// Handle the different API requests
	var fn func() (interface{}, error)
	switch strings.ToLower(strings.TrimSpace(c.flagTarget)) {
	case "barrier":
		fn = func() (interface{}, error) {
			return client.Sys().RekeyStatus()
		}
		if c.flagVerify {
			fn = func() (interface{}, error) {
				return client.Sys().RekeyVerificationStatus()
			}
		}
	case "recovery", "hsm":
		fn = func() (interface{}, error) {
			return client.Sys().RekeyRecoveryKeyStatus()
		}
		if c.flagVerify {
			fn = func() (interface{}, error) {
				return client.Sys().RekeyRecoveryKeyVerificationStatus()
			}
		}
	default:
		c.UI.Error(fmt.Sprintf("Unknown target: %s", c.flagTarget))
		return 1
	}

	// Make the request
	status, err := fn()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error reading rekey status: %s", err))
		return 2
	}

	return c.printStatus(status)
}

// backupRetrieve retrieves the stored backup keys.
func (c *OperatorRekeyCommand) backupRetrieve(client *api.Client) int {
	// Handle the different API requests
	var fn func() (*api.RekeyRetrieveResponse, error)
	switch strings.ToLower(strings.TrimSpace(c.flagTarget)) {
	case "barrier":
		fn = client.Sys().RekeyRetrieveBackup
	case "recovery", "hsm":
		fn = client.Sys().RekeyRetrieveRecoveryBackup
	default:
		c.UI.Error(fmt.Sprintf("Unknown target: %s", c.flagTarget))
		return 1
	}

	// Make the request
	storedKeys, err := fn()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error retrieving rekey stored keys: %s", err))
		return 2
	}

	secret := &api.Secret{
		Data: structs.New(storedKeys).Map(),
	}

	return OutputSecret(c.UI, secret)
}

// backupDelete deletes the stored backup keys.
func (c *OperatorRekeyCommand) backupDelete(client *api.Client) int {
	// Handle the different API requests
	var fn func() error
	switch strings.ToLower(strings.TrimSpace(c.flagTarget)) {
	case "barrier":
		fn = client.Sys().RekeyDeleteBackup
	case "recovery", "hsm":
		fn = client.Sys().RekeyDeleteRecoveryBackup
	default:
		c.UI.Error(fmt.Sprintf("Unknown target: %s", c.flagTarget))
		return 1
	}

	// Make the request
	if err := fn(); err != nil {
		c.UI.Error(fmt.Sprintf("Error deleting rekey stored keys: %s", err))
		return 2
	}

	c.UI.Output("Success! Delete stored keys (if they existed)")
	return 0
}

// printStatus dumps the status to output
func (c *OperatorRekeyCommand) printStatus(in interface{}) int {
	out := []string{}
	out = append(out, "Key | Value")

	switch in.(type) {
	case *api.RekeyStatusResponse:
		status := in.(*api.RekeyStatusResponse)
		out = append(out, fmt.Sprintf("Nonce | %s", status.Nonce))
		out = append(out, fmt.Sprintf("Started | %t", status.Started))
		if status.Started {
			if status.Progress == status.Required {
				out = append(out, fmt.Sprintf("Rekey Progress | %d/%d (verification in progress)", status.Progress, status.Required))
			} else {
				out = append(out, fmt.Sprintf("Rekey Progress | %d/%d", status.Progress, status.Required))
			}
			out = append(out, fmt.Sprintf("New Shares | %d", status.N))
			out = append(out, fmt.Sprintf("New Threshold | %d", status.T))
			out = append(out, fmt.Sprintf("Verification Required | %t", status.VerificationRequired))
			if status.VerificationNonce != "" {
				out = append(out, fmt.Sprintf("Verification Nonce | %s", status.VerificationNonce))
			}
		}
		if len(status.PGPFingerprints) > 0 {
			out = append(out, fmt.Sprintf("PGP Fingerprints | %s", status.PGPFingerprints))
			out = append(out, fmt.Sprintf("Backup | %t", status.Backup))
		}
	case *api.RekeyVerificationStatusResponse:
		status := in.(*api.RekeyVerificationStatusResponse)
		out = append(out, fmt.Sprintf("Started | %t", status.Started))
		out = append(out, fmt.Sprintf("New Shares | %d", status.N))
		out = append(out, fmt.Sprintf("New Threshold | %d", status.T))
		out = append(out, fmt.Sprintf("Verification Nonce | %s", status.Nonce))
		out = append(out, fmt.Sprintf("Verification Progress | %d/%d", status.Progress, status.T))
	default:
		c.UI.Error("Unknown status type")
		return 1
	}

	switch Format(c.UI) {
	case "table":
		c.UI.Output(tableOutput(out, nil))
		return 0
	default:
		return OutputData(c.UI, in)
	}
}

func (c *OperatorRekeyCommand) printUnsealKeys(client *api.Client, status *api.RekeyStatusResponse, resp *api.RekeyUpdateResponse) int {
	switch Format(c.UI) {
	case "table":
	default:
		return OutputData(c.UI, resp)
	}

	// Space between the key prompt, if any, and the output
	c.UI.Output("")

	// Provide the keys
	var haveB64 bool
	if resp.KeysB64 != nil && len(resp.KeysB64) == len(resp.Keys) {
		haveB64 = true
	}
	for i, key := range resp.Keys {
		if len(resp.PGPFingerprints) > 0 {
			if haveB64 {
				c.UI.Output(fmt.Sprintf("Key %d fingerprint: %s; value: %s", i+1, resp.PGPFingerprints[i], resp.KeysB64[i]))
			} else {
				c.UI.Output(fmt.Sprintf("Key %d fingerprint: %s; value: %s", i+1, resp.PGPFingerprints[i], key))
			}
		} else {
			if haveB64 {
				c.UI.Output(fmt.Sprintf("Key %d: %s", i+1, resp.KeysB64[i]))
			} else {
				c.UI.Output(fmt.Sprintf("Key %d: %s", i+1, key))
			}
		}
	}

	c.UI.Output("")
	c.UI.Output(fmt.Sprintf("Operation nonce: %s", resp.Nonce))

	if len(resp.PGPFingerprints) > 0 && resp.Backup {
		c.UI.Output("")
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"The encrypted unseal keys are backed up to \"core/unseal-keys-backup\"" +
				"in the storage backend. Remove these keys at any time using " +
				"\"vault operator rekey -backup-delete\". Vault does not automatically " +
				"remove these keys.",
		)))
	}

	switch status.VerificationRequired {
	case false:
		c.UI.Output("")
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Vault rekeyed with %d key shares and a key threshold of %d. Please "+
				"securely distribute the key shares printed above. When Vault is "+
				"re-sealed, restarted, or stopped, you must supply at least %d of "+
				"these keys to unseal it before it can start servicing requests.",
			status.N,
			status.T,
			status.T)))
	default:
		c.UI.Output("")
		c.UI.Output(wrapAtLength(fmt.Sprintf(
			"Vault has created a new key, split into %d key shares and a key threshold "+
				"of %d. These will not be active until after verification is complete. "+
				"Please securely distribute the key shares printed above. When Vault "+
				"is re-sealed, restarted, or stopped, you must supply at least %d of "+
				"these keys to unseal it before it can start servicing requests.",
			status.N,
			status.T,
			status.T)))
		c.UI.Output("")
		c.UI.Warn(wrapAtLength(
			"Again, these key shares are _not_ valid until verification is performed. " +
				"Do not lose or discard your current key shares until after verification " +
				"is complete or you will be unable to unseal Vault. If you cancel the " +
				"rekey process or seal Vault before verification is complete the new " +
				"shares will be discarded and the current shares will remain valid.",
		))
		c.UI.Output("")
		c.UI.Warn(wrapAtLength(
			"The current verification status, including initial nonce, is shown below.",
		))
		c.UI.Output("")

		c.flagVerify = true
		return c.status(client)
	}

	return 0
}
