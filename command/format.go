package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ghodss/yaml"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/ryanuber/columnize"
)

var predictFormat complete.Predictor = complete.PredictSet("json", "yaml")

func OutputSecret(ui cli.Ui, format string, secret *api.Secret) int {
	return outputWithFormat(ui, format, secret, secret)
}

func OutputList(ui cli.Ui, format string, secret *api.Secret) int {
	return outputWithFormat(ui, format, secret, secret.Data["keys"])
}

func outputWithFormat(ui cli.Ui, format string, secret *api.Secret, data interface{}) int {
	formatter, ok := Formatters[strings.ToLower(format)]
	if !ok {
		ui.Error(fmt.Sprintf("Invalid output format: %s", format))
		return 1
	}
	if err := formatter.Output(ui, secret, data); err != nil {
		ui.Error(fmt.Sprintf("Could not output secret: %s", err.Error()))
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

// An output formatter for json output of an object
type JsonFormatter struct {
}

func (j JsonFormatter) Output(ui cli.Ui, secret *api.Secret, data interface{}) error {
	b, err := json.Marshal(data)
	if err == nil {
		var out bytes.Buffer
		json.Indent(&out, b, "", "\t")
		ui.Output(out.String())
	}
	return err
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
	// TODO: this should really use reflection like the other formatters do
	if s, ok := data.(*api.Secret); ok {
		return t.OutputSecret(ui, secret, s)
	}
	if s, ok := data.([]interface{}); ok {
		return t.OutputList(ui, secret, s)
	}
	return errors.New("Cannot use the table formatter for this type")
}

func (t TableFormatter) OutputList(ui cli.Ui, secret *api.Secret, list []interface{}) error {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)

	if len(list) > 0 {
		input = append(input, "Keys")
		input = append(input, "----")

		keys := make([]string, 0, len(list))
		for _, k := range list {
			keys = append(keys, k.(string))
		}
		sort.Strings(keys)

		for _, k := range keys {
			input = append(input, fmt.Sprintf("%s", k))
		}
	}

	tableOutputStr := columnize.Format(input, config)

	// Print the warning separately because the length of first
	// column in the output will be increased by the length of
	// the longest warning string making the output look bad.
	warningsInput := make([]string, 0, 5)
	if len(secret.Warnings) != 0 {
		warningsInput = append(warningsInput, "")
		warningsInput = append(warningsInput, "The following warnings were returned from the Vault server:")
		for _, warning := range secret.Warnings {
			warningsInput = append(warningsInput, fmt.Sprintf("* %s", warning))
		}
	}

	warningsOutputStr := columnize.Format(warningsInput, config)

	ui.Output(fmt.Sprintf("%s\n%s", tableOutputStr, warningsOutputStr))

	return nil
}

func (t TableFormatter) OutputSecret(ui cli.Ui, secret, s *api.Secret) error {
	config := columnize.DefaultConfig()
	config.Delim = "♨"
	config.Glue = "\t"
	config.Prefix = ""

	input := make([]string, 0, 5)

	onceHeader := &sync.Once{}
	headerFunc := func() {
		input = append(input, fmt.Sprintf("Key %s Value", config.Delim))
		input = append(input, fmt.Sprintf("--- %s -----", config.Delim))
	}

	if s.LeaseDuration > 0 {
		onceHeader.Do(headerFunc)
		if s.LeaseID != "" {
			input = append(input, fmt.Sprintf("lease_id %s %s", config.Delim, s.LeaseID))
			input = append(input, fmt.Sprintf(
				"lease_duration %s %s", config.Delim, (time.Second*time.Duration(s.LeaseDuration)).String()))
		} else {
			input = append(input, fmt.Sprintf(
				"refresh_interval %s %s", config.Delim, (time.Second*time.Duration(s.LeaseDuration)).String()))
		}
		if s.LeaseID != "" {
			input = append(input, fmt.Sprintf(
				"lease_renewable %s %s", config.Delim, strconv.FormatBool(s.Renewable)))
		}
	}

	if s.Auth != nil {
		onceHeader.Do(headerFunc)
		input = append(input, fmt.Sprintf("token %s %s", config.Delim, s.Auth.ClientToken))
		input = append(input, fmt.Sprintf("token_accessor %s %s", config.Delim, s.Auth.Accessor))
		input = append(input, fmt.Sprintf("token_duration %s %s", config.Delim, (time.Second*time.Duration(s.Auth.LeaseDuration)).String()))
		input = append(input, fmt.Sprintf("token_renewable %s %v", config.Delim, s.Auth.Renewable))
		input = append(input, fmt.Sprintf("token_policies %s %v", config.Delim, s.Auth.Policies))
		for k, v := range s.Auth.Metadata {
			input = append(input, fmt.Sprintf("token_meta_%s %s %#v", k, config.Delim, v))
		}
	}

	if s.WrapInfo != nil {
		onceHeader.Do(headerFunc)
		input = append(input, fmt.Sprintf("wrapping_token: %s %s", config.Delim, s.WrapInfo.Token))
		input = append(input, fmt.Sprintf("wrapping_accessor: %s %s", config.Delim, s.WrapInfo.Accessor))
		input = append(input, fmt.Sprintf("wrapping_token_ttl: %s %s", config.Delim, (time.Second*time.Duration(s.WrapInfo.TTL)).String()))
		input = append(input, fmt.Sprintf("wrapping_token_creation_time: %s %s", config.Delim, s.WrapInfo.CreationTime.String()))
		input = append(input, fmt.Sprintf("wrapping_token_creation_path: %s %s", config.Delim, s.WrapInfo.CreationPath))
		if s.WrapInfo.WrappedAccessor != "" {
			input = append(input, fmt.Sprintf("wrapped_accessor: %s %s", config.Delim, s.WrapInfo.WrappedAccessor))
		}
	}

	if s.Data != nil && len(s.Data) > 0 {
		onceHeader.Do(headerFunc)
		keys := make([]string, 0, len(s.Data))
		for k := range s.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			input = append(input, fmt.Sprintf("%s %s %v", k, config.Delim, s.Data[k]))
		}
	}

	tableOutputStr := columnize.Format(input, config)

	// Print the warning separately because the length of first
	// column in the output will be increased by the length of
	// the longest warning string making the output look bad.
	warningsInput := make([]string, 0, 5)
	if len(s.Warnings) != 0 {
		warningsInput = append(warningsInput, "")
		warningsInput = append(warningsInput, "The following warnings were returned from the Vault server:")
		for _, warning := range s.Warnings {
			warningsInput = append(warningsInput, fmt.Sprintf("* %s", warning))
		}
	}

	warningsOutputStr := columnize.Format(warningsInput, config)

	ui.Output(fmt.Sprintf("%s\n%s", tableOutputStr, warningsOutputStr))

	return nil
}
