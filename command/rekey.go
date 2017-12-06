package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/meta"
	"github.com/posener/complete"
)

// RekeyCommand is a Command that rekeys the vault.
type RekeyCommand struct {
	meta.Meta

	// Key can be used to pre-seed the key. If it is set, it will not
	// be asked with the `password` helper.
	Key string

	// The nonce for the rekey request to send along
	Nonce string

	// Whether to use the recovery key instead of barrier key, if available
	RecoveryKey bool
}

func (c *RekeyCommand) Run(args []string) int {
	var init, cancel, status, delete, retrieve, backup, recoveryKey bool
	var shares, threshold, storedShares int
	var nonce string
	var pgpKeys pgpkeys.PubKeyFilesFlag
	flags := c.Meta.FlagSet("rekey", meta.FlagSetDefault)
	flags.BoolVar(&init, "init", false, "")
	flags.BoolVar(&cancel, "cancel", false, "")
	flags.BoolVar(&status, "status", false, "")
	flags.BoolVar(&delete, "delete", false, "")
	flags.BoolVar(&retrieve, "retrieve", false, "")
	flags.BoolVar(&backup, "backup", false, "")
	flags.BoolVar(&recoveryKey, "recovery-key", c.RecoveryKey, "")
	flags.IntVar(&shares, "key-shares", 5, "")
	flags.IntVar(&threshold, "key-threshold", 3, "")
	flags.IntVar(&storedShares, "stored-shares", 0, "")
	flags.StringVar(&nonce, "nonce", "", "")
	flags.Var(&pgpKeys, "pgp-keys", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if nonce != "" {
		c.Nonce = nonce
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	// Check if we are running doing any restricted variants
	switch {
	case init:
		return c.initRekey(client, shares, threshold, storedShares, pgpKeys, backup, recoveryKey)
	case cancel:
		return c.cancelRekey(client, recoveryKey)
	case status:
		return c.rekeyStatus(client, recoveryKey)
	case retrieve:
		return c.rekeyRetrieveStored(client, recoveryKey)
	case delete:
		return c.rekeyDeleteStored(client, recoveryKey)
	}

	// Check if the rekey is started
	var rekeyStatus *api.RekeyStatusResponse
	if recoveryKey {
		rekeyStatus, err = client.Sys().RekeyRecoveryKeyStatus()
	} else {
		rekeyStatus, err = client.Sys().RekeyStatus()
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading rekey status: %s", err))
		return 1
	}

	// Start the rekey process if not started
	if !rekeyStatus.Started {
		if recoveryKey {
			rekeyStatus, err = client.Sys().RekeyRecoveryKeyInit(&api.RekeyInitRequest{
				SecretShares:    shares,
				SecretThreshold: threshold,
				PGPKeys:         pgpKeys,
				Backup:          backup,
			})
		} else {
			rekeyStatus, err = client.Sys().RekeyInit(&api.RekeyInitRequest{
				SecretShares:    shares,
				SecretThreshold: threshold,
				PGPKeys:         pgpKeys,
				Backup:          backup,
			})
		}
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error initializing rekey: %s", err))
			return 1
		}
		c.Nonce = rekeyStatus.Nonce
	}

	shares = rekeyStatus.N
	threshold = rekeyStatus.T
	serverNonce := rekeyStatus.Nonce

	// Get the unseal key
	args = flags.Args()
	key := c.Key
	if len(args) > 0 {
		key = args[0]
	}
	if key == "" {
		c.Nonce = serverNonce
		fmt.Printf("Rekey operation nonce: %s\n", serverNonce)
		fmt.Printf("Key (will be hidden): ")
		key, err = password.Read(os.Stdin)
		fmt.Printf("\n")
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error attempting to ask for password. The raw error message\n"+
					"is shown below, but the most common reason for this error is\n"+
					"that you attempted to pipe a value into unseal or you're\n"+
					"executing `vault rekey` from outside of a terminal.\n\n"+
					"You should use `vault rekey` from a terminal for maximum\n"+
					"security. If this isn't an option, the unseal key can be passed\n"+
					"in using the first parameter.\n\n"+
					"Raw error: %s", err))
			return 1
		}
	}

	// Provide the key, this may potentially complete the update
	var result *api.RekeyUpdateResponse
	if recoveryKey {
		result, err = client.Sys().RekeyRecoveryKeyUpdate(strings.TrimSpace(key), c.Nonce)
	} else {
		result, err = client.Sys().RekeyUpdate(strings.TrimSpace(key), c.Nonce)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error attempting rekey update: %s", err))
		return 1
	}

	// If we are not complete, then dump the status
	if !result.Complete {
		return c.rekeyStatus(client, recoveryKey)
	}

	// Space between the key prompt, if any, and the output
	c.Ui.Output("\n")
	// Provide the keys
	var haveB64 bool
	if result.KeysB64 != nil && len(result.KeysB64) == len(result.Keys) {
		haveB64 = true
	}
	for i, key := range result.Keys {
		if len(result.PGPFingerprints) > 0 {
			if haveB64 {
				c.Ui.Output(fmt.Sprintf("Key %d fingerprint: %s; value: %s", i+1, result.PGPFingerprints[i], result.KeysB64[i]))
			} else {
				c.Ui.Output(fmt.Sprintf("Key %d fingerprint: %s; value: %s", i+1, result.PGPFingerprints[i], key))
			}
		} else {
			if haveB64 {
				c.Ui.Output(fmt.Sprintf("Key %d: %s", i+1, result.KeysB64[i]))
			} else {
				c.Ui.Output(fmt.Sprintf("Key %d: %s", i+1, key))
			}
		}
	}

	c.Ui.Output(fmt.Sprintf("\nOperation nonce: %s", result.Nonce))

	if len(result.PGPFingerprints) > 0 && result.Backup {
		c.Ui.Output(fmt.Sprintf(
			"\n" +
				"The encrypted unseal keys have been backed up to \"core/unseal-keys-backup\"\n" +
				"in your physical backend. It is your responsibility to remove these if and\n" +
				"when desired.",
		))
	}

	c.Ui.Output(fmt.Sprintf(
		"\n"+
			"Vault rekeyed with %d keys and a key threshold of %d.\n",
		shares,
		threshold,
	))

	// Print this message if keys are returned
	if len(result.Keys) > 0 {
		c.Ui.Output(fmt.Sprintf(
			"\n"+
				"Please securely distribute the above keys. When the vault is re-sealed,\n"+
				"restarted, or stopped, you must provide at least %d of these keys\n"+
				"to unseal it again.\n\n"+
				"Vault does not store the master key. Without at least %[1]d keys,\n"+
				"your vault will remain permanently sealed.",
			threshold,
		))
	}

	return 0
}

// initRekey is used to start the rekey process
func (c *RekeyCommand) initRekey(client *api.Client,
	shares, threshold, storedShares int,
	pgpKeys pgpkeys.PubKeyFilesFlag,
	backup, recoveryKey bool) int {
	// Start the rekey
	request := &api.RekeyInitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
		StoredShares:    storedShares,
		PGPKeys:         pgpKeys,
		Backup:          backup,
	}
	var status *api.RekeyStatusResponse
	var err error
	if recoveryKey {
		status, err = client.Sys().RekeyRecoveryKeyInit(request)
	} else {
		status, err = client.Sys().RekeyInit(request)
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing rekey: %s", err))
		return 1
	}

	if pgpKeys == nil || len(pgpKeys) == 0 {
		c.Ui.Output(`
WARNING: If you lose the keys after they are returned to you, there is no
recovery. Consider using the '-pgp-keys' option to protect the returned unseal
keys along with '-backup=true' to allow recovery of the encrypted keys in case
of emergency. They can easily be deleted at a later time with
'vault rekey -delete'.
`)
	}

	if pgpKeys != nil && len(pgpKeys) > 0 && !backup {
		c.Ui.Output(`
WARNING: You are using PGP keys for encryption, but have not set the option to
back up the new unseal keys to physical storage. If you lose the keys after
they are returned to you, there is no recovery. Consider setting '-backup=true'
to allow recovery of the encrypted keys in case of emergency. They can easily
be deleted at a later time with 'vault rekey -delete'.
`)
	}

	// Provide the current status
	return c.dumpRekeyStatus(status)
}

// cancelRekey is used to abort the rekey process
func (c *RekeyCommand) cancelRekey(client *api.Client, recovery bool) int {
	var err error
	if recovery {
		err = client.Sys().RekeyRecoveryKeyCancel()
	} else {
		err = client.Sys().RekeyCancel()
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to cancel rekey: %s", err))
		return 1
	}
	c.Ui.Output("Rekey canceled.")
	return 0
}

// rekeyStatus is used just to fetch and dump the status
func (c *RekeyCommand) rekeyStatus(client *api.Client, recovery bool) int {
	// Check the status
	var status *api.RekeyStatusResponse
	var err error
	if recovery {
		status, err = client.Sys().RekeyRecoveryKeyStatus()
	} else {
		status, err = client.Sys().RekeyStatus()
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading rekey status: %s", err))
		return 1
	}

	return c.dumpRekeyStatus(status)
}

func (c *RekeyCommand) dumpRekeyStatus(status *api.RekeyStatusResponse) int {
	// Dump the status
	statString := fmt.Sprintf(
		"Nonce: %s\n"+
			"Started: %t\n"+
			"Key Shares: %d\n"+
			"Key Threshold: %d\n"+
			"Rekey Progress: %d\n"+
			"Required Keys: %d",
		status.Nonce,
		status.Started,
		status.N,
		status.T,
		status.Progress,
		status.Required,
	)
	if len(status.PGPFingerprints) != 0 {
		statString = fmt.Sprintf("%s\nPGP Key Fingerprints: %s", statString, status.PGPFingerprints)
		statString = fmt.Sprintf("%s\nBackup Storage: %t", statString, status.Backup)
	}
	c.Ui.Output(statString)
	return 0
}

func (c *RekeyCommand) rekeyRetrieveStored(client *api.Client, recovery bool) int {
	var storedKeys *api.RekeyRetrieveResponse
	var err error
	if recovery {
		storedKeys, err = client.Sys().RekeyRetrieveRecoveryBackup()
	} else {
		storedKeys, err = client.Sys().RekeyRetrieveBackup()
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error retrieving stored keys: %s", err))
		return 1
	}

	secret := &api.Secret{
		Data: structs.New(storedKeys).Map(),
	}

	return OutputSecret(c.Ui, "table", secret)
}

func (c *RekeyCommand) rekeyDeleteStored(client *api.Client, recovery bool) int {
	var err error
	if recovery {
		err = client.Sys().RekeyDeleteRecoveryBackup()
	} else {
		err = client.Sys().RekeyDeleteBackup()
	}
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to delete stored keys: %s", err))
		return 1
	}
	c.Ui.Output("Stored keys deleted.")
	return 0
}

func (c *RekeyCommand) Synopsis() string {
	return "Rekeys Vault to generate new unseal keys"
}

func (c *RekeyCommand) Help() string {
	helpText := `
Usage: vault rekey [options] [key]

  Rekey is used to change the unseal keys. This can be done to generate
  a new set of unseal keys or to change the number of shares and the
  required threshold.

  Rekey can only be done when the vault is already unsealed. The operation
  is done online, but requires that a threshold of the current unseal
  keys be provided.

General Options:
` + meta.GeneralOptionsUsage() + `
Rekey Options:

  -init                   Initialize the rekey operation by setting the desired
                          number of shares and the key threshold. This can only be
                          done if no rekey is already initiated.

  -cancel                 Reset the rekey process by throwing away
                          prior keys and the rekey configuration.

  -status                 Prints the status of the current rekey operation.
                          This can be used to see the status without attempting
                          to provide an unseal key.

  -retrieve               Retrieve backed-up keys. Only available if the PGP keys
                          were provided and the backup has not been deleted.

  -delete                 Delete any backed-up keys.

  -key-shares=5           The number of key shares to split the master key
                          into.

  -key-threshold=3        The number of key shares required to reconstruct
                          the master key.

  -nonce=abcd             The nonce provided at rekey initialization time. This
                          same nonce value must be provided with each unseal
                          key. If the unseal key is not being passed in via the
                          the command line the nonce parameter is not required,
                          and will instead be displayed with the key prompt.

  -pgp-keys               If provided, must be a comma-separated list of
                          files on disk containing binary- or base64-format
                          public PGP keys, or Keybase usernames specified as
                          "keybase:<username>". The number of given entries
                          must match 'key-shares'. The output unseal keys will
                          be encrypted and base64-encoded, in order, with the
                          given public keys.  If you want to use them with the
                          'vault unseal' command, you will need to base64-decode
                          and decrypt; this will be the plaintext unseal key.

  -backup=false           If true, and if the key shares are PGP-encrypted, a
                          plaintext backup of the PGP-encrypted keys will be
                          stored at "core/unseal-keys-backup" in your physical
                          storage. You can retrieve or delete them via the
                          'sys/rekey/backup' endpoint.

  -recovery-key=false     Whether to rekey the recovery key instead of the
                          barrier key. Only used with Vault HSM.
`
	return strings.TrimSpace(helpText)
}

func (c *RekeyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RekeyCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-init":          complete.PredictNothing,
		"-cancel":        complete.PredictNothing,
		"-status":        complete.PredictNothing,
		"-retrieve":      complete.PredictNothing,
		"-delete":        complete.PredictNothing,
		"-key-shares":    complete.PredictNothing,
		"-key-threshold": complete.PredictNothing,
		"-nonce":         complete.PredictNothing,
		"-pgp-keys":      complete.PredictNothing,
		"-backup":        complete.PredictNothing,
		"-recovery-key":  complete.PredictNothing,
	}
}
