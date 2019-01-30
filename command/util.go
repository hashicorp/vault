package command

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/config"
	"github.com/hashicorp/vault/command/token"
	"github.com/mitchellh/cli"
)

// DefaultTokenHelper returns the token helper that is configured for Vault.
func DefaultTokenHelper() (token.TokenHelper, error) {
	return config.DefaultTokenHelper()
}

// RawField extracts the raw field from the given data and returns it as a
// string for printing purposes.
func RawField(secret *api.Secret, field string) interface{} {
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
			val = secret.Auth.TokenPolicies
		case "identity_policies":
			val = secret.Auth.IdentityPolicies
		case "policies":
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
		case "lease_duration":
			val = secret.LeaseDuration
		case "lease_id":
			val = secret.LeaseID
		case "request_id":
			val = secret.RequestID
		case "renewable":
			val = secret.Renewable
		case "refresh_interval":
			val = secret.LeaseDuration
		case "data":
			var ok bool
			val, ok = secret.Data["data"]
			if !ok {
				val = secret.Data
			}
		default:
			val = secret.Data[field]
		}
	}

	return val
}

// PrintRawField prints raw field from the secret.
func PrintRawField(ui cli.Ui, data interface{}, field string) int {
	var val interface{}
	switch data.(type) {
	case *api.Secret:
		val = RawField(data.(*api.Secret), field)
	case map[string]interface{}:
		val = data.(map[string]interface{})[field]
	}

	if val == nil {
		ui.Error(fmt.Sprintf("Field %q not present in secret", field))
		return 1
	}

	format := Format(ui)
	if format == "" || format == "table" {
		return PrintRaw(ui, fmt.Sprintf("%v", val))
	}

	// Handle specific format flags as best as possible
	formatter, ok := Formatters[format]
	if !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	b, err := formatter.Format(val)
	if err != nil {
		ui.Error(fmt.Sprintf("Error formatting output: %s", err))
		return 1
	}

	return PrintRaw(ui, string(b))
}

// PrintRaw prints a raw value to the terminal. If the process is being "piped"
// to something else, the "raw" value is printed without a newline character.
// Otherwise the value is printed as normal.
func PrintRaw(ui cli.Ui, str string) int {
	if !color.NoColor {
		ui.Output(str)
	} else {
		// The cli.Ui prints a CR, which is not wanted since the user probably wants
		// just the raw value.
		w := getWriterFromUI(ui)
		fmt.Fprint(w, str)
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
