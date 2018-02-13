package command

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/mitchellh/cli"
)

// DefaultTokenHelper returns the token helper that is configured for Vault.
func DefaultTokenHelper() (token.TokenHelper, error) {
	config, err := LoadConfig("")
	if err != nil {
		return nil, err
	}

	path := config.TokenHelper
	if path == "" {
		return &token.InternalTokenHelper{}, nil
	}

	path, err = token.ExternalTokenHelperPath(path)
	if err != nil {
		return nil, err
	}
	return &token.ExternalTokenHelper{BinaryPath: path}, nil
}

// RawField extracts the raw field from the given data and returns it as a
// string for printing purposes.
func RawField(secret *api.Secret, field string) (string, bool) {
	var val interface{}
	switch {
	case secret.Auth != nil:
		switch field {
		case "token":
			val = secret.Auth.ClientToken
		case "token_accessor":
			val = secret.Auth.Accessor
		case "token_duration":
			val = secret.Auth.LeaseDuration
		case "token_renewable":
			val = secret.Auth.Renewable
		case "token_policies":
			val = secret.Auth.Policies
		default:
			val = secret.Data[field]
		}

	case secret.WrapInfo != nil:
		switch field {
		case "wrapping_token":
			val = secret.WrapInfo.Token
		case "wrapping_accessor":
			val = secret.WrapInfo.Accessor
		case "wrapping_token_ttl":
			val = secret.WrapInfo.TTL
		case "wrapping_token_creation_time":
			val = secret.WrapInfo.CreationTime.Format(time.RFC3339Nano)
		case "wrapping_token_creation_path":
			val = secret.WrapInfo.CreationPath
		case "wrapped_accessor":
			val = secret.WrapInfo.WrappedAccessor
		default:
			val = secret.Data[field]
		}

	default:
		switch field {
		case "refresh_interval":
			val = secret.LeaseDuration
		default:
			val = secret.Data[field]
		}
	}

	str := fmt.Sprintf("%v", val)
	return str, val != nil
}

// PrintRawField prints raw field from the secret.
func PrintRawField(ui cli.Ui, secret *api.Secret, field string) int {
	str, ok := RawField(secret, field)
	if !ok {
		ui.Error(fmt.Sprintf("Field %q not present in secret", field))
		return 1
	}

	return PrintRaw(ui, str)
}

// PrintRaw prints a raw value to the terminal. If the process is being "piped"
// to something else, the "raw" value is printed without a newline character.
// Otherwise the value is printed as normal.
func PrintRaw(ui cli.Ui, str string) int {
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		ui.Output(str)
	} else {
		// The cli.Ui prints a CR, which is not wanted since the user probably wants
		// just the raw value.
		w := getWriterFromUI(ui)
		fmt.Fprintf(w, str)
	}
	return 0
}

// getWriterFromUI accepts a cli.Ui and returns the underlying io.Writer by
// unwrapping as many wrapped Uis as necessary. If there is an unknown UI
// type, this falls back to os.Stdout.
func getWriterFromUI(ui cli.Ui) io.Writer {
	switch t := ui.(type) {
	case *VaultUI:
		return getWriterFromUI(t.Ui)
	case *cli.BasicUi:
		return t.Writer
	case *cli.ColoredUi:
		return getWriterFromUI(t.Ui)
	case *cli.ConcurrentUi:
		return getWriterFromUI(t.Ui)
	case *cli.MockUi:
		return t.OutputWriter
	default:
		return os.Stdout
	}
}
