package command

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

const (
	wrappedResponseLocation = "cubbyhole/response"
)

// UnwrapCommand is a Command that behaves like ReadCommand but specifically
// for unwrapping cubbyhole-wrapped secrets
type UnwrapCommand struct {
	meta.Meta
}

func (c *UnwrapCommand) Run(args []string) int {
	var format string
	var field string
	var err error
	var secret *api.Secret
	var flags *flag.FlagSet
	flags = c.Meta.FlagSet("unwrap", meta.FlagSetDefault)
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&field, "field", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 || len(args[0]) == 0 {
		c.Ui.Error("Unwrap expects one argument: the ID of the wrapping token")
		flags.Usage()
		return 1
	}

	tokenID := args[0]
	_, err = uuid.ParseUUID(tokenID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Given token could not be parsed as a UUID: %s", err))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error initializing client: %s", err))
		return 2
	}

	client.SetToken(tokenID)

	secret, err = c.getUnwrappedResponse(client)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	if secret == nil {
		c.Ui.Error("Secret returned was nil")
		return 1
	}

	// Handle single field output
	if field != "" {
		return PrintRawField(c.Ui, secret, field)
	}

	return OutputSecret(c.Ui, format, secret)
}

// getUnwrappedResponse is a helper to do the actual reading and unwrapping
func (c *UnwrapCommand) getUnwrappedResponse(client *api.Client) (*api.Secret, error) {
	secret, err := client.Logical().Read(wrappedResponseLocation)
	if err != nil {
		return nil, fmt.Errorf("Error reading %s: %s", wrappedResponseLocation, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("No value found at %s", wrappedResponseLocation)
	}
	if secret.Data == nil {
		return nil, fmt.Errorf("\"data\" not found in wrapping response")
	}
	if _, ok := secret.Data["response"]; !ok {
		return nil, fmt.Errorf("\"response\" not found in wrapping response \"data\" map")
	}

	wrappedSecret := new(api.Secret)
	buf := bytes.NewBufferString(secret.Data["response"].(string))
	dec := json.NewDecoder(buf)
	dec.UseNumber()
	if err := dec.Decode(wrappedSecret); err != nil {
		return nil, fmt.Errorf("Error unmarshaling wrapped secret: %s", err)
	}

	return wrappedSecret, nil
}

func (c *UnwrapCommand) Synopsis() string {
	return "Unwrap a wrapped secret"
}

func (c *UnwrapCommand) Help() string {
	helpText := `
Usage: vault unwrap [options] <wrapping token ID>

  Unwrap a wrapped secret.

  Unwraps the data wrapped by the given token ID. The returned result is the
  same as a 'read' operation on a non-wrapped secret.

General Options:
` + meta.GeneralOptionsUsage() + `
Read Options:

  -format=table           The format for output. By default it is a whitespace-
                          delimited table. This can also be json or yaml.

  -field=field            If included, the raw value of the specified field
                          will be output raw to stdout.

`
	return strings.TrimSpace(helpText)
}
