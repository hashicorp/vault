// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/cli"
	envenc "github.com/hashicorp/vault-envelope-encryption-sdk"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TransitDecryptCommand)(nil)
	_ cli.CommandAutocomplete = (*TransitDecryptCommand)(nil)
)

type TransitDecryptCommand struct {
	transitEnvEncCommand
}

func (c *TransitDecryptCommand) Synopsis() string {
	return "Decrypt files or streams using envelope encryption backed by Vault."
}

func (c *TransitDecryptCommand) Help() string {
	helpText := `
Usage: vault transit envelope decrypt transit-key-path [options...] [filenames...]

  Using a Transit encryption key, decrypt one or more files or a stream from STDIN
  with envelope encryption.  In this mode, the client requests Vault to decrypt an Encrypted Data Key
  attached to the file/stream to retrieve its Data Encryption Key (DEK), which is then
  used to decrypt the data.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransitDecryptCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputFormat)

	f := set.NewFlagSet("Command Options")
	c.transitEnvEncCommand.setFlags(f)
	f.IntVar(&IntVar{
		Name:       "time-window-keys-per-day",
		Usage:      "For time keyed window key management, how many keyed time periods subdivide a day.  0 or not specified to use per session data keys.",
		Target:     &c.timeWindowKeysPerDay,
		Completion: complete.PredictAnything,
	})

	return set
}

func (c *TransitDecryptCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitDecryptCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitDecryptCommand) Run(args []string) int {
	return c.envelopeDecrypt(c.BaseCommand, c.Flags(), args)
}

// error codes: 1: user error, 2: internal computation error, 3: remote api call error
func (tc *TransitDecryptCommand) envelopeDecrypt(c *BaseCommand, flags *FlagSets, args []string) int {
	// Parse and validate the arguments.
	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.UI.Error(fmt.Sprintf("Incorrect argument count (expected 1+, got %d). Wanted transit key path (:mount:/keys/:key_name:)", len(args)))
		return 1
	}

	keyPath := args[0]
	targets := args[1:]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	if tc.output != "" {
		if tc.suffix != "" {
			c.UI.Error("Cannot specify suffix and output at the same time.")
			return 1
		}
		if len(targets) > 1 {
			c.UI.Error("Cannot specify output with more than one input.")
			return 1
		}
	}

	aad, err := readFlagBytes(tc.aad, "aad")
	if err != nil {
		c.UI.Error(fmt.Sprintf("aad could not be read: %v", err))
		return 1
	}

	path, key, err := transitKeyAndPath(keyPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not parse transit path and key: %v", err))
		return 1
	}

	kp, err := envenc.NewTransitKeyProvider(envenc.ProviderConfig{
		Client:    client,
		CacheSize: len(targets),
		KeyName:   key,
		Backend:   path,
	})
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not initialize transit key provider: %v", err))
		return 2
	}

	// For encrypt, this is straightforward.  For decrypt, I imagined building the Header channel in the command, and passing a function with the channel
	// in the closure of envelopeDecrypt, something like:
	// func envelopeDecrypt(ch chan *Header) envelopeOpFunc {
	//   return func(...) { NewDecrypting, etc }
	// }
	// Then when decryption happens, the outer function can read the channel and display the header, while handleFiles doesn't care

	var headerOut chan *envenc.Header
	var wg sync.WaitGroup
	if !tc.quiet {
		wg.Add(1)
		headerOut = make(chan *envenc.Header)
		go func() {
			defer wg.Done()
			for header := range headerOut {
				OutputData(c.UI, header.Map())
			}
		}()
	}
	rc := handleFiles(&tc.transitEnvEncCommand, targets, kp, envelopeDecrypt(headerOut, aad), true)

	if !tc.quiet {
		close(headerOut)
		wg.Wait()
	}
	return rc
}

func envelopeDecrypt(headerOut chan *envenc.Header, aad []byte) envelopeFileOp {
	return func(kp envenc.KeyProvider, in io.Reader, out io.Writer, length *int64) error {
		opts := []envenc.Option{}
		if headerOut != nil {
			opts = append(opts, envenc.WithHeaderOutChan(headerOut))
		}
		if len(aad) > 0 {
			opts = append(opts, envenc.WithAad(aad))
		}
		if length != nil {
			opts = append(opts, envenc.WithLength(length))
		}
		decIn, err := envenc.NewDecryptingReader(kp, in, opts...)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, decIn)
		return err
	}
}
