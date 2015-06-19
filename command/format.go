package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

func OutputSecret(ui cli.Ui, format string, secret *api.Secret) int {
	switch format {
	case "json":
		return outputFormatJSON(ui, secret)
	case "table":
		fallthrough
	default:
		return outputFormatTable(ui, secret, true)
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

func outputFormatTable(ui cli.Ui, s *api.Secret, whitespace bool) int {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)
	input = append(input, fmt.Sprintf("Key %s Value", config.Delim))

	if s.LeaseID != "" && s.LeaseDuration > 0 {
		input = append(input, fmt.Sprintf("lease_id %s %s", config.Delim, s.LeaseID))
		input = append(input, fmt.Sprintf(
			"lease_duration %s %d", config.Delim, s.LeaseDuration))
		input = append(input, fmt.Sprintf(
			"lease_renewable %s %s", config.Delim, strconv.FormatBool(s.Renewable)))
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

	for k, v := range s.Data {
		input = append(input, fmt.Sprintf("%s %s %v", k, config.Delim, v))
	}

	ui.Output(columnize.Format(input, config))
	return 0
}
