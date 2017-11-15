package command

import (
	"fmt"
	"os"
	"reflect"
	"time"

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

func PrintRawField(ui cli.Ui, secret *api.Secret, field string) int {
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

	if val != nil {
		// c.Ui.Output() prints a CR character which in this case is
		// not desired. Since Vault CLI currently only uses BasicUi,
		// which writes to standard output, os.Stdout is used here to
		// directly print the message. If mitchellh/cli exposes method
		// to print without CR, this check needs to be removed.
		if reflect.TypeOf(ui).String() == "*cli.BasicUi" {
			fmt.Fprintf(os.Stdout, "%v", val)
		} else {
			ui.Output(fmt.Sprintf("%v", val))
		}
		return 0
	} else {
		ui.Error(fmt.Sprintf(
			"Field %s not present in secret", field))
		return 1
	}
}
