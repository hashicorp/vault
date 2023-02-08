package command

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"nhooyr.io/websocket"
)

var (
	_ cli.Command             = (*EventsSubscribeCommands)(nil)
	_ cli.CommandAutocomplete = (*EventsSubscribeCommands)(nil)
)

type EventsSubscribeCommands struct {
	*BaseCommand

	flagTimeout time.Duration
}

func (c *EventsSubscribeCommands) Synopsis() string {
	return "Subscribe to events"
}

func (c *EventsSubscribeCommands) Help() string {
	helpText := `
Usage: vault events subscribe [-format=json] [-timeout=XYZs] eventType

  Subscribe to events of the given event type (topic). The events will be
  output to standard out.

  If the format parameter is not specified or is not "json", then the
  output format will be protobuf binary format. The schema can be found in
  sdk/logical/event.proto in the Vault source code. 

  If the format is specified as "json", then the output will be a JSON object
  serialized using the default protobuf JSON serialization format, with one
  line per event received.
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *EventsSubscribeCommands) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	// Common Options
	f := set.NewFlagSet("Common Options")

	f.DurationVar(&DurationVar{
		Name:    "timeout",
		Target:  &c.flagTimeout,
		Default: 0,
		Usage:   `Terminate the websocket after this amount of time. 0 means forever.`,
	})
	return set
}

func (c *EventsSubscribeCommands) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *EventsSubscribeCommands) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *EventsSubscribeCommands) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 1:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 1, got %d)", len(args)))
		return 1
	case len(args) > 1:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 1, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	err = c.subscribeRequest(client, "sys/events/subscribe/"+args[0])
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	return 0
}

func (c *EventsSubscribeCommands) subscribeRequest(client *api.Client, path string) error {
	r := client.NewRequest("GET", "/v1/"+path)
	u := r.URL
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else {
		u.Scheme = "wss"
	}
	if c.flagFormat == "json" {
		q := u.Query()
		q.Set("json", "true")
		u.RawQuery = q.Encode()
	}
	client.AddHeader("X-Vault-Token", client.Token())
	client.AddHeader("X-Vault-Namesapce", client.Namespace())
	ctx := context.Background()
	if c.flagTimeout != 0 {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(ctx, c.flagTimeout)
		defer cancelFunc()
	}
	conn, _, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPClient: client.CloneConfig().HttpClient,
		HTTPHeader: client.Headers(),
	})
	if err != nil {
		c.UI.Error(err.Error())
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	for {
		_, message, err := conn.Read(ctx)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(message)
		if err != nil {
			return err
		}
	}
}
