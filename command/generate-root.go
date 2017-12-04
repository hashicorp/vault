package command

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/password"
	"github.com/hashicorp/vault/helper/pgpkeys"
	"github.com/hashicorp/vault/helper/xor"
	"github.com/hashicorp/vault/meta"
	"github.com/posener/complete"
)

// GenerateRootCommand is a Command that generates a new root token.
type GenerateRootCommand struct {
	meta.Meta

	// Key can be used to pre-seed the key. If it is set, it will not
	// be asked with the `password` helper.
	Key string

	// The nonce for the rekey request to send along
	Nonce string
}

func (c *GenerateRootCommand) Run(args []string) int {
	var init, cancel, status, genotp, drToken bool
	var nonce, decode, otp, pgpKey string
	var pgpKeyArr pgpkeys.PubKeyFilesFlag
	flags := c.Meta.FlagSet("generate-root", meta.FlagSetDefault)
	flags.BoolVar(&init, "init", false, "")
	flags.BoolVar(&drToken, "dr-token", false, "")
	flags.BoolVar(&cancel, "cancel", false, "")
	flags.BoolVar(&status, "status", false, "")
	flags.BoolVar(&genotp, "genotp", false, "")
	flags.StringVar(&decode, "decode", "", "")
	flags.StringVar(&otp, "otp", "", "")
	flags.StringVar(&nonce, "nonce", "", "")
	flags.Var(&pgpKeyArr, "pgp-key", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	if genotp {
		buf := make([]byte, 16)
		readLen, err := rand.Read(buf)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading random bytes: %s", err))
			return 1
		}
		if readLen != 16 {
			c.Ui.Error(fmt.Sprintf("Read %d bytes when we should have read 16", readLen))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("OTP: %s", base64.StdEncoding.EncodeToString(buf)))
		return 0
	}

	if len(decode) > 0 {
		if len(otp) == 0 {
			c.Ui.Error("Both the value to decode and the OTP must be passed in")
			return 1
		}
		return c.decode(decode, otp)
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	// Check if the root generation is started
	f := client.Sys().GenerateRootStatus
	if drToken {
		f = client.Sys().GenerateDROperationTokenStatus
	}
	rootGenerationStatus, err := f()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading root generation status: %s", err))
		return 1
	}

	// If we are initing, or if we are not started but are not running a
	// special function, check otp and pgpkey
	checkOtpPgp := false
	switch {
	case init:
		checkOtpPgp = true
	case cancel:
	case status:
	case genotp:
	case len(decode) != 0:
	case rootGenerationStatus.Started:
	default:
		checkOtpPgp = true
	}
	if checkOtpPgp {
		switch {
		case len(otp) == 0 && (pgpKeyArr == nil || len(pgpKeyArr) == 0):
			c.Ui.Error(c.Help())
			return 1
		case len(otp) != 0 && pgpKeyArr != nil && len(pgpKeyArr) != 0:
			c.Ui.Error(c.Help())
			return 1
		case len(otp) != 0:
			err := c.verifyOTP(otp)
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Error verifying the provided OTP: %s", err))
				return 1
			}
		case pgpKeyArr != nil:
			if len(pgpKeyArr) != 1 {
				c.Ui.Error("Could not parse PGP key")
				return 1
			}
			if len(pgpKeyArr[0]) == 0 {
				c.Ui.Error("Got an empty PGP key")
				return 1
			}
			pgpKey = pgpKeyArr[0]
		default:
			panic("unreachable case")
		}
	}

	if nonce != "" {
		c.Nonce = nonce
	}

	// Check if we are running doing any restricted variants
	switch {
	case init:
		return c.initGenerateRoot(client, otp, pgpKey, drToken)
	case cancel:
		return c.cancelGenerateRoot(client, drToken)
	case status:
		return c.rootGenerationStatus(client, drToken)
	}

	// Start the root generation process if not started
	if !rootGenerationStatus.Started {
		f := client.Sys().GenerateRootInit
		if drToken {
			f = client.Sys().GenerateDROperationTokenInit
		}
		rootGenerationStatus, err = f(otp, pgpKey)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error initializing root generation: %s", err))
			return 1
		}
		c.Nonce = rootGenerationStatus.Nonce
	}

	serverNonce := rootGenerationStatus.Nonce

	// Get the unseal key
	args = flags.Args()
	key := c.Key
	if len(args) > 0 {
		key = args[0]
	}
	if key == "" {
		c.Nonce = serverNonce
		fmt.Printf("Root generation operation nonce: %s\n", serverNonce)
		fmt.Printf("Key (will be hidden): ")
		key, err = password.Read(os.Stdin)
		fmt.Printf("\n")
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error attempting to ask for password. The raw error message\n"+
					"is shown below, but the most common reason for this error is\n"+
					"that you attempted to pipe a value into unseal or you're\n"+
					"executing `vault generate-root` from outside of a terminal.\n\n"+
					"You should use `vault generate-root` from a terminal for maximum\n"+
					"security. If this isn't an option, the unseal key can be passed\n"+
					"in using the first parameter.\n\n"+
					"Raw error: %s", err))
			return 1
		}
	}

	// Provide the key, this may potentially complete the update
	{
		f := client.Sys().GenerateRootUpdate
		if drToken {
			f = client.Sys().GenerateDROperationTokenUpdate
		}
		statusResp, err := f(strings.TrimSpace(key), c.Nonce)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error attempting generate-root update: %s", err))
			return 1
		}

		c.dumpStatus(statusResp)
	}
	return 0
}

func (c *GenerateRootCommand) verifyOTP(otp string) error {
	if len(otp) == 0 {
		return fmt.Errorf("No OTP passed in")
	}
	otpBytes, err := base64.StdEncoding.DecodeString(otp)
	if err != nil {
		return fmt.Errorf("Error decoding base64 OTP value: %s", err)
	}
	if otpBytes == nil || len(otpBytes) != 16 {
		return fmt.Errorf("Decoded OTP value is invalid or wrong length")
	}

	return nil
}

func (c *GenerateRootCommand) decode(encodedVal, otp string) int {
	tokenBytes, err := xor.XORBase64(encodedVal, otp)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	token, err := uuid.FormatUUID(tokenBytes)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error formatting base64 token value: %v", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Root token: %s", token))

	return 0
}

// initGenerateRoot is used to start the generation process
func (c *GenerateRootCommand) initGenerateRoot(client *api.Client, otp string, pgpKey string, drToken bool) int {
	// Start the rekey
	f := client.Sys().GenerateRootInit
	if drToken {
		f = client.Sys().GenerateDROperationTokenInit
	}

	status, err := f(otp, pgpKey)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing root generation: %s", err))
		return 1
	}

	c.dumpStatus(status)

	return 0
}

// cancelGenerateRoot is used to abort the generation process
func (c *GenerateRootCommand) cancelGenerateRoot(client *api.Client, drToken bool) int {
	f := client.Sys().GenerateRootCancel
	if drToken {
		f = client.Sys().GenerateDROperationTokenCancel
	}
	err := f()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to cancel root generation: %s", err))
		return 1
	}
	c.Ui.Output("Root generation canceled.")
	return 0
}

// rootGenerationStatus is used just to fetch and dump the status
func (c *GenerateRootCommand) rootGenerationStatus(client *api.Client, drToken bool) int {
	// Check the status
	f := client.Sys().GenerateRootStatus
	if drToken {
		f = client.Sys().GenerateDROperationTokenStatus
	}
	status, err := f()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading root generation status: %s", err))
		return 1
	}

	c.dumpStatus(status)

	return 0
}

// dumpStatus dumps the status to output
func (c *GenerateRootCommand) dumpStatus(status *api.GenerateRootStatusResponse) {
	// Dump the status
	statString := fmt.Sprintf(
		"Nonce: %s\n"+
			"Started: %v\n"+
			"Generate Root Progress: %d\n"+
			"Required Keys: %d\n"+
			"Complete: %t",
		status.Nonce,
		status.Started,
		status.Progress,
		status.Required,
		status.Complete,
	)
	if len(status.PGPFingerprint) > 0 {
		statString = fmt.Sprintf("%s\nPGP Fingerprint: %s", statString, status.PGPFingerprint)
	}
	if len(status.EncodedRootToken) > 0 {
		statString = fmt.Sprintf("%s\n\nEncoded root token: %s", statString, status.EncodedRootToken)
	} else if len(status.EncodedToken) > 0 {
		statString = fmt.Sprintf("%s\n\nEncoded token: %s", statString, status.EncodedToken)
	}
	c.Ui.Output(statString)
}

func (c *GenerateRootCommand) Synopsis() string {
	return "Generates a new root token"
}

func (c *GenerateRootCommand) Help() string {
	helpText := `
Usage: vault generate-root [options] [key]

  'generate-root' is used to create a new root token.

  Root generation can only be done when the vault is already unsealed. The
  operation is done online, but requires that a threshold of the current unseal
  keys be provided.

  One (and only one) of the following must be provided when initializing the
  root generation attempt:

  1) A 16-byte, base64-encoded One Time Password (OTP) provided in the '-otp'
  flag; the token is XOR'd with this value before it is returned once the final
  unseal key has been provided. The '-decode' operation can be used with this
  value and the OTP to output the final token value. The '-genotp' flag can be
  used to generate a suitable value.

  or

  2) A file containing a PGP key (binary or base64-encoded) or a Keybase.io
  username in the format of "keybase:<username>" in the '-pgp-key' flag. The
  final token value will be encrypted with this public key and base64-encoded.

General Options:
` + meta.GeneralOptionsUsage() + `
Generate Root Options:

  -init                   Initialize the root generation attempt. This can only
                          be done if no generation is already initiated.

  -cancel                 Reset the root generation process by throwing away
                          prior unseal keys and the configuration.

  -status                 Prints the status of the current attempt. This can be
                          used to see the status without attempting to provide
                          an unseal key.

  -decode=abcd            Decodes and outputs the generated root token. The OTP
                          used at '-init' time must be provided in the '-otp'
                          parameter.

  -genotp                 Returns a high-quality OTP suitable for passing into
                          the '-init' method.

  -otp=abcd               The base64-encoded 16-byte OTP for use with the
                          '-init' or '-decode' methods.

  -pgp-key                A file on disk containing a binary- or base64-format
                          public PGP key, or a Keybase username specified as
                          "keybase:<username>". The output root token will be
                          encrypted and base64-encoded, in order, with the given
                          public key.

  -nonce=abcd             The nonce provided at initialization time. This same
                          nonce value must be provided with each unseal key. If
                          the unseal key is not being passed in via the command
                          line the nonce parameter is not required, and will
                          instead be displayed with the key prompt.

  -dr-token               Generate a Disaster Recovery operation token. This flag
                          should be set on '-init', '-cancel', and every time a 
                          key is provided to specify the type of token to generate.
`
	return strings.TrimSpace(helpText)
}

func (c *GenerateRootCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *GenerateRootCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{
		"-init":    complete.PredictNothing,
		"-cancel":  complete.PredictNothing,
		"-status":  complete.PredictNothing,
		"-decode":  complete.PredictNothing,
		"-genotp":  complete.PredictNothing,
		"-otp":     complete.PredictNothing,
		"-pgp-key": complete.PredictNothing,
		"-nonce":   complete.PredictNothing,
	}
}
