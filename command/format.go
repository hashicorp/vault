package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

const (
	// hopeDelim is the delimiter to use when splitting columns. We call it a
	// hopeDelim because we hope that it's never contained in a secret.
	hopeDelim = "♨"
)

type FormatOptions struct {
	Format string
}

func OutputSecret(ui cli.Ui, secret *api.Secret) int {
	return outputWithFormat(ui, secret, secret)
}

func OutputList(ui cli.Ui, data interface{}) int {
	switch data.(type) {
	case *api.Secret:
		secret := data.(*api.Secret)
		return outputWithFormat(ui, secret, secret.Data["keys"])
	default:
		return outputWithFormat(ui, nil, data)
	}
}

func OutputData(ui cli.Ui, data interface{}) int {
	return outputWithFormat(ui, nil, data)
}

func outputWithFormat(ui cli.Ui, secret *api.Secret, data interface{}) int {
	format := Format(ui)
	formatter, ok := Formatters[format]
	if !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}

	if err := formatter.Output(ui, secret, data); err != nil {
		ui.Error(fmt.Sprintf("Could not parse output: %s", err.Error()))
		return 1
	}
	return 0
}

type Formatter interface {
	Output(ui cli.Ui, secret *api.Secret, data interface{}) error
}

var Formatters = map[string]Formatter{
	"json":  JsonFormatter{},
	"table": TableFormatter{},
	"yaml":  YamlFormatter{},
	"yml":   YamlFormatter{},
}

func format() string {
	format := os.Getenv(EnvVaultFormat)
	if format == "" {
		format = "table"
	}
	return format
}

func Format(ui cli.Ui) string {
	switch ui.(type) {
	case *VaultUI:
		return ui.(*VaultUI).format
	}
	return format()
}

// An output formatter for json output of an object
type JsonFormatter struct{}

func (j JsonFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	ui.Output(string(b))
	return nil
}

// An output formatter for yaml output format of an object
type YamlFormatter struct {
}

func (y YamlFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	b, err := yaml.Marshal(data)
	if err == nil {
		ui.Output(strings.TrimSpace(string(b)))
	}
	return err
}

// An output formatter for table output of an object
type TableFormatter struct {
}

func (t TableFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	switch data.(type) {
	case *api.Secret:
		return t.OutputSecret(ui, secret)
	case []interface{}:
		return t.OutputList(ui, secret, data)
	case []string:
		return t.OutputList(ui, nil, data)
	default:
		return errors.New("Cannot use the table formatter for this type")
	}
}

func (t TableFormatter) OutputList(ui cli.Ui, secret *api.Secret, data interface{}) error {
	t.printWarnings(ui, secret)

	switch data.(type) {
	case []interface{}:
	case []string:
		ui.Output(tableOutput(data.([]string), nil))
		return nil
	default:
		return errors.New("Error: table formatter cannot output list for this data type")
	}

	list := data.([]interface{})

	if len(list) > 0 {
		keys := make([]string, len(list))
		for i, v := range list {
			typed, ok := v.(string)
			if !ok {
				return fmt.Errorf("Error: %v is not a string", v)
			}
			keys[i] = typed
		}
		sort.Strings(keys)

		// Prepend the header
		keys = append([]string{"Keys"}, keys...)

		ui.Output(tableOutput(keys, &columnize.Config{
			Delim: hopeDelim,
		}))
	}

	return nil
}

// printWarnings prints any warnings in the secret.
func (t TableFormatter) printWarnings(ui cli.Ui, secret *api.Secret) {
	if secret != nil && len(secret.Warnings) > 0 {
		ui.Warn("WARNING! The following warnings were returned from Vault:\n")
		for _, warning := range secret.Warnings {
			ui.Warn(wrapAtLengthWithPadding(fmt.Sprintf("* %s", warning), 2))
		}
		ui.Warn("")
	}
}

func (t TableFormatter) OutputSecret(ui cli.Ui, secret *api.Secret) error {
	if secret == nil {
		return nil
	}

	t.printWarnings(ui, secret)

	out := make([]string, 0, 8)
	if secret.LeaseDuration > 0 {
		if secret.LeaseID != "" {
			out = append(out, fmt.Sprintf("lease_id %s %s", hopeDelim, secret.LeaseID))
			out = append(out, fmt.Sprintf("lease_duration %s %s", hopeDelim, humanDurationInt(secret.LeaseDuration)))
			out = append(out, fmt.Sprintf("lease_renewable %s %t", hopeDelim, secret.Renewable))
		} else {
			// This is probably the generic secret backend which has leases, but we
			// print them as refresh_interval to reduce confusion.
			out = append(out, fmt.Sprintf("refresh_interval %s %s", hopeDelim, humanDurationInt(secret.LeaseDuration)))
		}
	}

	if secret.Auth != nil {
		out = append(out, fmt.Sprintf("token %s %s", hopeDelim, secret.Auth.ClientToken))
		out = append(out, fmt.Sprintf("token_accessor %s %s", hopeDelim, secret.Auth.Accessor))
		// If the lease duration is 0, it's likely a root token, so output the
		// duration as "infinity" to clear things up.
		if secret.Auth.LeaseDuration == 0 {
			out = append(out, fmt.Sprintf("token_duration %s %s", hopeDelim, "∞"))
		} else {
			out = append(out, fmt.Sprintf("token_duration %s %s", hopeDelim, humanDurationInt(secret.Auth.LeaseDuration)))
		}
		out = append(out, fmt.Sprintf("token_renewable %s %t", hopeDelim, secret.Auth.Renewable))
		out = append(out, fmt.Sprintf("token_policies %s %v", hopeDelim, secret.Auth.Policies))
		for k, v := range secret.Auth.Metadata {
			out = append(out, fmt.Sprintf("token_meta_%s %s %v", k, hopeDelim, v))
		}
	}

	if secret.WrapInfo != nil {
		out = append(out, fmt.Sprintf("wrapping_token: %s %s", hopeDelim, secret.WrapInfo.Token))
		out = append(out, fmt.Sprintf("wrapping_accessor: %s %s", hopeDelim, secret.WrapInfo.Accessor))
		out = append(out, fmt.Sprintf("wrapping_token_ttl: %s %s", hopeDelim, humanDurationInt(secret.WrapInfo.TTL)))
		out = append(out, fmt.Sprintf("wrapping_token_creation_time: %s %s", hopeDelim, secret.WrapInfo.CreationTime.String()))
		out = append(out, fmt.Sprintf("wrapping_token_creation_path: %s %s", hopeDelim, secret.WrapInfo.CreationPath))
		if secret.WrapInfo.WrappedAccessor != "" {
			out = append(out, fmt.Sprintf("wrapped_accessor: %s %s", hopeDelim, secret.WrapInfo.WrappedAccessor))
		}
	}

	if len(secret.Data) > 0 {
		keys := make([]string, 0, len(secret.Data))
		for k := range secret.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			out = append(out, fmt.Sprintf("%s %s %v", k, hopeDelim, secret.Data[k]))
		}
	}

	// If we got this far and still don't have any data, there's nothing to print,
	// sorry.
	if len(out) == 0 {
		return nil
	}

	// Prepend the header
	out = append([]string{"Key" + hopeDelim + "Value"}, out...)

	ui.Output(tableOutput(out, &columnize.Config{
		Delim: hopeDelim,
	}))
	return nil
}

// OutputSealStatus will print *api.SealStatusResponse in the CLI according to the format provided
func OutputSealStatus(ui cli.Ui, client *api.Client, status *api.SealStatusResponse) int {
	switch Format(ui) {
	case "table":
	default:
		return OutputData(ui, status)
	}

	var sealPrefix string
	if status.RecoverySeal {
		sealPrefix = "Recovery "
	}

	out := []string{}
	out = append(out, "Key | Value")
	out = append(out, fmt.Sprintf("%sSeal Type | %s", sealPrefix, status.Type))
	out = append(out, fmt.Sprintf("Sealed | %t", status.Sealed))
	out = append(out, fmt.Sprintf("Total %sShares | %d", sealPrefix, status.N))
	out = append(out, fmt.Sprintf("Threshold | %d", status.T))

	if status.Sealed {
		out = append(out, fmt.Sprintf("Unseal Progress | %d/%d", status.Progress, status.T))
		out = append(out, fmt.Sprintf("Unseal Nonce | %s", status.Nonce))
	}

	out = append(out, fmt.Sprintf("Version | %s", status.Version))

	if status.ClusterName != "" && status.ClusterID != "" {
		out = append(out, fmt.Sprintf("Cluster Name | %s", status.ClusterName))
		out = append(out, fmt.Sprintf("Cluster ID | %s", status.ClusterID))
	}

	// Mask the 'Vault is sealed' error, since this means HA is enabled, but that
	// we cannot query for the leader since we are sealed.
	leaderStatus, err := client.Sys().Leader()
	if err != nil && strings.Contains(err.Error(), "Vault is sealed") {
		leaderStatus = &api.LeaderResponse{HAEnabled: true}
	}

	// Output if HA is enabled
	out = append(out, fmt.Sprintf("HA Enabled | %t", leaderStatus.HAEnabled))
	if leaderStatus.HAEnabled {
		mode := "sealed"
		if !status.Sealed {
			mode = "standby"
			if leaderStatus.IsSelf {
				mode = "active"
			}
		}

		out = append(out, fmt.Sprintf("HA Mode | %s", mode))

		if !status.Sealed {
			out = append(out, fmt.Sprintf("HA Cluster | %s", leaderStatus.LeaderClusterAddress))
		}
	}

	ui.Output(tableOutput(out, nil))
	return 0
}
