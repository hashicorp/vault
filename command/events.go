// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

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

	namespaces  []string
	bexprFilter string
}

func (c *EventsSubscribeCommands) Synopsis() string {
	return "Subscribe to events"
}

func (c *EventsSubscribeCommands) Help() string {
	helpText := `
Usage: vault events subscribe [-namespaces=ns1] [-timeout=XYZs] [-filter=filterExpression] eventType

  Subscribe to events of the given event type (topic), which may be a glob
  pattern (with "*"" treated as a wildcard). The events will be sent to
  standard out.

  The output will be a JSON object serialized using the default protobuf
  JSON serialization format, with one line per event received.
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *EventsSubscribeCommands) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Subscribe Options")
	f.StringVar(&StringVar{
		Name: "filter",
		Usage: `A boolean expression to use to filter events. Only events matching
                the filter will be subscribed to. This is applied after any filtering
                by event type or namespace.`,
		Default: "",
		Target:  &c.bexprFilter,
	})
	f.StringSliceVar(&StringSliceVar{
		Name: "namespaces",
		Usage: `Specifies one or more patterns of additional child namespaces
                to subscribe to. The namespace of the request is automatically
                prepended, so specifying 'ns2' when the request is in the 'ns1'
                namespace will result in subscribing to 'ns1/ns2', in addition to
                'ns1'. Patterns can include "*" characters to indicate
                wildcards. The default is to subscribe only to the request's
                namespace.`,
		Default: []string{},
		Target:  &c.namespaces,
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

// cleanNamespace removes leading and trailing space and /'s from the namespace path.
func cleanNamespace(ns string) string {
	ns = strings.TrimSpace(ns)
	ns = strings.Trim(ns, "/")
	return ns
}

// cleanNamespaces removes leading and trailing space and /'s from the namespace paths.
func cleanNamespaces(namespaces []string) []string {
	cleaned := make([]string, len(namespaces))
	for i, ns := range namespaces {
		cleaned[i] = cleanNamespace(ns)
	}
	return cleaned
}

func (c *EventsSubscribeCommands) subscribeRequest(client *api.Client, path string) error {
	r := client.NewRequest("GET", "/v1/"+path)
	u := r.URL
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else {
		u.Scheme = "wss"
	}
	q := u.Query()
	q.Set("json", "true")
	if len(c.namespaces) > 0 {
		q["namespaces"] = cleanNamespaces(c.namespaces)
	}
	bexprFilter := strings.TrimSpace(c.bexprFilter)
	if bexprFilter != "" {
		q.Set("filter", bexprFilter)
	}
	u.RawQuery = q.Encode()
	client.AddHeader("X-Vault-Token", client.Token())
	client.AddHeader("X-Vault-Namespace", client.Namespace())
	ctx := context.Background()

	// Follow redirects in case our request if our request is forwarded to the leader.
	url := u.String()
	var conn *websocket.Conn
	var err error
	for attempt := 0; attempt < 10; attempt++ {
		var resp *http.Response
		conn, resp, err = websocket.Dial(ctx, url, &websocket.DialOptions{
			HTTPClient: client.CloneConfig().HttpClient,
			HTTPHeader: client.Headers(),
		})

		if err == nil {
			break
		}

		switch {
		case resp == nil:
			return err
		case resp.StatusCode == http.StatusTemporaryRedirect:
			url = resp.Header.Get("Location")
			continue
		case resp.StatusCode == http.StatusNotFound:
			return errors.New("events endpoint not found; check `vault read sys/experiments` to see if an events experiment is available but disabled")
		default:
			return err
		}
	}

	if conn == nil {
		return fmt.Errorf("too many redirects")
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
