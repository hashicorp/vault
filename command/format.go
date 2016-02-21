package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

func OutputSecret(ui cli.Ui, format string, secret *api.Secret) int {
	switch format {
	case "json":
		return outputFormatJSON(ui, secret)
	case "yaml":
		return outputFormatYAML(ui, secret)
	case "table":
		return outputFormatTable(ui, secret, true)
	default:
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}
}

func OutputList(ui cli.Ui, format string, secret *api.Secret) int {
	switch format {
	case "json":
		return outputFormatJSONList(ui, secret)
	case "yaml":
		return outputFormatYAMLList(ui, secret)
	case "table":
		return outputFormatTableList(ui, secret)
	default:
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}
}

func outputFormatJSON(ui cli.Ui, s *api.Secret) int {
	b, err := json.Marshal(s)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error formatting secret: %s", err))
		return 1
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	ui.Output(out.String())
	return 0
}

func outputFormatJSONList(ui cli.Ui, s *api.Secret) int {
	b, err := json.Marshal(s.Data["keys"])
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error formatting keys: %s", err))
		return 1
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	ui.Output(out.String())
	return 0
}

func outputFormatYAML(ui cli.Ui, s *api.Secret) int {
	b, err := yaml.Marshal(s)
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error formatting secret: %s", err))
		return 1
	}

	ui.Output(strings.TrimSpace(string(b)))
	return 0
}

func outputFormatYAMLList(ui cli.Ui, s *api.Secret) int {
	b, err := yaml.Marshal(s.Data["keys"])
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error formatting secret: %s", err))
		return 1
	}

	ui.Output(strings.TrimSpace(string(b)))
	return 0
}

func outputFormatTable(ui cli.Ui, s *api.Secret, whitespace bool) int {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)

	input = append(input, fmt.Sprintf("Key %s Value", config.Delim))

	if s.LeaseDuration > 0 {
		if s.LeaseID != "" {
			input = append(input, fmt.Sprintf("lease_id %s %s", config.Delim, s.LeaseID))
		}
		input = append(input, fmt.Sprintf(
			"lease_duration %s %d", config.Delim, s.LeaseDuration))
		if s.LeaseID != "" {
			input = append(input, fmt.Sprintf(
				"lease_renewable %s %s", config.Delim, strconv.FormatBool(s.Renewable)))
		}
	}

	if s.Auth != nil {
		input = append(input, fmt.Sprintf("token %s %s", config.Delim, s.Auth.ClientToken))
		input = append(input, fmt.Sprintf("token_duration %s %d", config.Delim, s.Auth.LeaseDuration))
		input = append(input, fmt.Sprintf("token_renewable %s %v", config.Delim, s.Auth.Renewable))
		input = append(input, fmt.Sprintf("token_policies %s %v", config.Delim, s.Auth.Policies))
		for k, v := range s.Auth.Metadata {
			input = append(input, fmt.Sprintf("token_meta_%s %s %#v", k, config.Delim, v))
		}
	}

	keys := make([]string, 0, len(s.Data))
	for k := range s.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		input = append(input, fmt.Sprintf("%s %s %v", k, config.Delim, s.Data[k]))
	}

	if len(s.Warnings) != 0 {
		input = append(input, "")
		input = append(input, "The following warnings were returned from the Vault server:")
		for _, warning := range s.Warnings {
			input = append(input, fmt.Sprintf("* %s", warning))
		}
	}

	ui.Output(columnize.Format(input, config))
	return 0
}

func outputFormatTableList(ui cli.Ui, s *api.Secret) int {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)

	input = append(input, "Keys")

	keys := make([]string, 0, len(s.Data["keys"].([]interface{})))
	for _, k := range s.Data["keys"].([]interface{}) {
		keys = append(keys, k.(string))
	}
	sort.Strings(keys)

	for _, k := range keys {
		input = append(input, fmt.Sprintf("%s", k))
	}

	if len(s.Warnings) != 0 {
		input = append(input, "")
		for _, warning := range s.Warnings {
			input = append(input, fmt.Sprintf("* %s", warning))
		}
	}

	ui.Output(columnize.Format(input, config))
	return 0
}
