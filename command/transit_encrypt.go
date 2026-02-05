// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/cli"
	envenc "github.com/hashicorp/vault-envelope-encryption-sdk"
	"github.com/posener/complete"
	"google.golang.org/protobuf/types/known/structpb"
)

const defaultSuffix = ".vee"

var (
	_ cli.Command             = (*TransitEncryptCommand)(nil)
	_ cli.CommandAutocomplete = (*TransitEncryptCommand)(nil)
)

type transitEnvEncCommand struct {
	*BaseCommand
	suffix               string
	output               string
	quiet                bool
	timeWindowKeysPerDay int
	aad                  string
}

func (c *transitEnvEncCommand) setFlags(f *FlagSet) {
	f.BoolVar(&BoolVar{
		Name:       "quiet",
		Aliases:    []string{"q"},
		Usage:      "If set, do not display any output.",
		Target:     &c.quiet,
		Completion: complete.PredictAnything,
	})

	f.StringVar(&StringVar{
		Name:       "output",
		Aliases:    []string{"o"},
		Usage:      "The name of the output file to write to, or '-' for stdout.",
		Target:     &c.output,
		Completion: complete.PredictAnything,
	})
	f.StringVar(&StringVar{
		Name:       "suffix",
		Aliases:    []string{"s"},
		Usage:      "The file suffix to add to encrypted files or to trim from decrypted ones.",
		Target:     &c.suffix,
		Completion: complete.PredictAnything,
	})
	f.StringVar(&StringVar{
		Name: "aad",
		Usage: `optional additional authenticated data.  Will be included in the authentication tag of the ciphertext, but not 
stored with it.  The same data provided at encryption time must be provided at decryption time for decryption to succeed.  
Must be base64 encoded or use "@file" to read from a file.`,
		Target:     &c.aad,
		Completion: complete.PredictAnything,
	})
}

type TransitEncryptCommand struct {
	transitEnvEncCommand
	keyBits     int
	mimeType    string
	metadata    string
	omitKeyData bool
}

// The function pattern expected by handleFiles
type envelopeFileOp func(envenc.KeyProvider, io.Reader, io.Writer, *int64) error

func (c *TransitEncryptCommand) Synopsis() string {
	return "Encrypt files or streams using envelope encryption backed by Vault."
}

func (c *TransitEncryptCommand) Help() string {
	helpText := `
Usage: vault transit envelope encrypt transit-key-path [options...] [filenames...]

  Using a Transit encryption key, encrypt one or more files or a stream from STDIN
  with envelope encryption.  In this mode, Vault provides a Data Encryption Key (DEK)
  and an Encrypted Data Key (EDK), which is the DEK encrypted by the key in Vault Transit.
  It then performs encryption of the target files or stream client side, with no practical
  limit to size, and without plaintext transiting to Vault.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TransitEncryptCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)

	f := set.NewFlagSet("Command Options")
	c.transitEnvEncCommand.setFlags(f)

	f.IntVar(&IntVar{
		Name:       "key-bits",
		Usage:      "The number of bits for the AES Data Encryption Key (DEK)",
		Target:     &c.keyBits,
		Completion: complete.PredictAnything,
	})

	f.IntVar(&IntVar{
		Name:       "time-window-keys-per-day",
		Usage:      "For time keyed window key management, how many keyed time periods subdivide a day.  0 or not specified to use per session data keys.",
		Target:     &c.timeWindowKeysPerDay,
		Completion: complete.PredictAnything,
	})
	// TODO: Mime type auto-detection?
	f.StringVar(&StringVar{
		Name:       "mime-type",
		Usage:      "Optionally, the mime-type of the plaintext file or stream.",
		Target:     &c.mimeType,
		Completion: complete.PredictAnything,
	})
	// TODO: Mime type auto-detection?
	f.StringVar(&StringVar{
		Name:       "metadata",
		Usage:      "optional arbitrary metadata that will be authenticated at decrypt time.  Must be base64 encoded or use \"@file\" to read from a file.",
		Target:     &c.metadata,
		Completion: complete.PredictAnything,
	})

	f.BoolVar(&BoolVar{
		Name:       "omit-key-data",
		Usage:      "Do not add to the header information about the transit key's namespace, path, and name",
		Target:     &c.omitKeyData,
		Completion: complete.PredictAnything,
	})

	return set
}

func (c *TransitEncryptCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TransitEncryptCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TransitEncryptCommand) Run(args []string) int {
	return c.EnvelopeEncrypt(c.BaseCommand, transitKeyPath, c.Flags(), args)
}

func transitKeyAndPath(s string) (path string, key string, err error) {
	parts := keyPath.FindStringSubmatch(s)
	if len(parts) != 3 {
		return "", "", errors.New("expected transit path and key name in the form :path:/keys/:name:")
	}
	path = parts[1]
	keyName := parts[2]

	return path, keyName, nil
}

// error codes: 1: user error, 2: internal computation error, 3: remote api call error
func (tc *TransitEncryptCommand) EnvelopeEncrypt(c *BaseCommand, pathFunc TransitKeyFunc, flags *FlagSets, args []string) int {
	// Parse and validate the arguments.
	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = flags.Args()
	if len(args) < 1 {
		c.UI.Error(fmt.Sprintf("Incorrect argument count (expected 2+, got %d). Wanted transit key path (:mount:/keys/:key_name:)", len(args)))
		return 1
	}

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
		if len(args[1:]) > 1 {
			c.UI.Error("Cannot specify output with more than one input.")
			return 1
		}
	}

	var metadata map[string]any
	if len(tc.metadata) > 0 {
		metadataBytes, err := readFlagBytes(tc.metadata, "metadata")
		if len(metadataBytes) > 0 {
			metadata = make(map[string]any)
			err = json.Unmarshal(metadataBytes, &metadata)
			if err != nil {
				c.UI.Error(fmt.Sprintf("metadata field could not be JSON parsed: %v", err))
				return 1
			}
		}
	}

	var aad []byte
	if len(tc.aad) > 0 {
		aad, err = readFlagBytes(tc.aad, "aad")
		if err != nil {
			c.UI.Error(fmt.Sprintf("aad could not be read: %v", err))
			return 1
		}
	}

	path, key, err := transitKeyAndPath(args[0])
	if err != nil {
		c.UI.Error(fmt.Sprintf("could not parse transit path and key: %v", err))
		return 1
	}

	kp, err := envenc.NewTransitKeyProvider(envenc.ProviderConfig{
		Client:    client,
		CacheSize: len(args[1:]),
		KeyName:   key,
		Backend:   path,
		KeyBits:   tc.keyBits,
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

	var keyData *envenc.KeyData
	ns := client.Namespace()
	keyData = &envenc.KeyData{
		MountPath: &path,
		KeyName:   &key,
	}
	if len(ns) > 0 {
		keyData.Namespace = &ns
	}

	ef, err := envelopeEncrypt(keyData, tc.mimeType, metadata, aad, tc.omitKeyData)
	if err != nil {
		c.UI.Error(fmt.Sprintf("error setting up encryption context: %v", err))
		return 3
	}
	return handleFiles(&tc.transitEnvEncCommand, args[1:], kp, ef, false)
}

func readFlagBytes(flagVal string, flagName string) (data []byte, err error) {
	if len(flagVal) > 0 {
		if flagVal[0] == '@' {
			data, err := os.ReadFile(flagVal[1:])
			if err != nil {
				return nil, fmt.Errorf("%s data file could not be read: %w", flagName, err)
			}
			return data, nil
		}
		data, err = base64.StdEncoding.DecodeString(flagVal)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func handleFiles(c *transitEnvEncCommand, i []string, kp envenc.KeyProvider, op envelopeFileOp, stripSuffix bool) int {
	for _, file := range i {
		errInt := func() int {
			var length *int64
			// So defers are preformed per file
			var out io.WriteCloser
			df, in, length, err := openInput(file, c.BaseCommand)
			if err != nil {
				c.UI.Error(err.Error())
				return 2
			}
			if df != nil {
				defer df()
			}

			var cleanup func()
			if c.output == "-" {
				out = os.Stdout
				cleanup = func() {}
			} else {
				var outFilename string
				if len(c.output) > 0 {
					outFilename = c.output
				} else {
					var suffix string
					if c.suffix == "" {
						suffix = defaultSuffix
					} else {
						suffix = c.suffix
					}
					if stripSuffix {
						if strings.HasSuffix(file, suffix) {
							outFilename = file[:len(file)-len(suffix)]
						}
					} else {
						outFilename = file + suffix
					}
				}

				// Avoid clobbering the input by using the same file as output
				fi2, err := os.Stat(outFilename)
				if err == nil {
					fi1, err := os.Stat(file)
					if err != nil {
						c.UI.Error(fmt.Sprintf("error stating input file: %v", err))
						return 6
					}

					if os.SameFile(fi1, fi2) {
						c.UI.Error("cannot output to the same file as input")
						return 7
					}
				} else if !os.IsNotExist(err) {
					c.UI.Error(fmt.Sprintf("error stating output file: %v", err))
					return 6
				}

				file, err := os.Create(outFilename)
				if err != nil {
					c.UI.Error(err.Error())
					return 3
				}
				cleanup = func() {
					if err := os.Remove(outFilename); err != nil {
						c.UI.Error(fmt.Sprintf("error cleaning up output file: %v", err))
					}
				}
				defer file.Close()
				out = file
			}
			bufferedOut := bufio.NewWriter(out)
			err = op(kp, in, bufferedOut, length)
			if err != nil {
				c.UI.Error(fmt.Sprintf("error processing stream: %v", err))
				cleanup()
				return 5
			}
			err = bufferedOut.Flush()
			if err != nil {
				c.UI.Error(fmt.Sprintf("error flushing stream: %v", err))
				cleanup()
				return 5
			}
			err = out.Close()
			if err != nil {
				c.UI.Error(fmt.Sprintf("error closing output: %v", err))
				cleanup()
			}
			return 0
		}()
		if errInt > 0 {
			return errInt
		}
	}
	return 0
}

func openInput(file string, c *BaseCommand) (defers func() error, in io.Reader, length *int64, err error) {
	var l int64
	var df func() error
	if file == "-" {
		return nil, os.Stdin, nil, nil
	} else {
		fileIn, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			return nil, nil, nil, err
		}
		df = fileIn.Close
		fi, err := fileIn.Stat()
		if err != nil {
			return nil, nil, nil, err
		}
		l = fi.Size()
		in = fileIn
	}
	in = bufio.NewReader(in)
	return df, in, &l, nil
}

func envelopeEncrypt(keyData *envenc.KeyData, mimeType string, metadata map[string]any, aad []byte, omitKeyData bool) (envelopeFileOp, error) {
	var md *structpb.Struct
	if metadata != nil {
		var err error
		md, err = structpb.NewStruct(metadata)
		if err != nil {
			return nil, err
		}
	}

	return func(kp envenc.KeyProvider, in io.Reader, out io.Writer, length *int64) error {
		header := envenc.NewHeader()
		v1 := header.GetV1()
		v1.MimeType = mimeType
		v1.Metadata = md
		if keyData != nil {
			v1.KeyData = keyData
		}

		opts := []envenc.Option{envenc.WithHeader(header)}
		if length != nil {
			opts = append(opts, envenc.WithLength(length))
		}
		if len(aad) > 0 {
			opts = append(opts, envenc.WithAad(aad))
		}
		if omitKeyData {
			opts = append(opts, envenc.WithOmitKeyData(true))
		}
		encOut, err := envenc.NewEncryptingWriter(kp, out, opts...)
		if err != nil {
			return err
		}
		_, err = io.Copy(encOut, in)
		if err != nil {
			encOut.Close()
			return err
		}
		err = encOut.Close()
		if err != nil {
			return fmt.Errorf("error closing output : %w", err)
		}
		return nil
	}, nil
}

var _ cli.Command = (*TransitEnvelopeCommand)(nil)

type TransitEnvelopeCommand struct {
	*BaseCommand
}

func (c *TransitEnvelopeCommand) Synopsis() string {
	return "Encrypt and decrypt files and streams using envelope encryption."
}

func (c *TransitEnvelopeCommand) Help() string {
	helpText := `
Usage: vault transit envelope encrypt|decrypt options... files...

  This command has subcommands for interacting with Vault's Transit Envelope
  Encryption.  This allows client side encryption and decryption of arbitrary 
  sized files and streams.

  An example of encrypting and decrypting a file using this feature:

  $ vault transit envelope encrypt transit/keys/kek myfile.txt

  $ vault transit envelope decrypt transit/keys/key myfile.txt.vee

  Please see the individual subcommand help for detailed usage information.
`
	return strings.TrimSpace(helpText)
}

func (c *TransitEnvelopeCommand) Run(args []string) int {
	return cli.RunResultHelp
}
