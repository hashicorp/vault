// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*EventsSubscriptionListCommand)(nil)
	_ cli.Command             = (*EventsSubscriptionGetCommand)(nil)
	_ cli.Command             = (*EventsSubscriptionCreateCommand)(nil)
	_ cli.Command             = (*EventsSubscriptionDeleteCommand)(nil)
	_ cli.CommandAutocomplete = (*EventsSubscriptionListCommand)(nil)
	_ cli.CommandAutocomplete = (*EventsSubscriptionGetCommand)(nil)
	_ cli.CommandAutocomplete = (*EventsSubscriptionCreateCommand)(nil)
	_ cli.CommandAutocomplete = (*EventsSubscriptionDeleteCommand)(nil)
)

type EventsSubscriptionListCommand struct {
	*BaseCommand
}

func (c *EventsSubscriptionListCommand) Synopsis() string {
	return "List subscriptions"
}

func (c *EventsSubscriptionListCommand) Help() string {
	helpText := `
Usage: vault events subscriptions list

  List all event subscriptions.
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *EventsSubscriptionListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *EventsSubscriptionListCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *EventsSubscriptionListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *EventsSubscriptionListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	secret, err := client.Logical().Read("sys/events/subscriptions")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	if secret == nil || secret.Data == nil || secret.Data["subscriptions"] == nil {
		return 0
	}
	switch Format(c.UI) {
	case "table":
		var out []string
		out = append(out, "ID | Plugin")
		for _, a := range secret.Data["subscriptions"].([]any) {
			m := a.(map[string]any)
			out = append(out, fmt.Sprintf("%s | %s", m["id"], m["plugin"]))
		}
		c.UI.Output(tableOutput(out, nil))
		return 0
	default:
		return OutputData(c.UI, secret.Data)
	}
}

type EventsSubscriptionGetCommand struct {
	*BaseCommand
}

func (c *EventsSubscriptionGetCommand) Synopsis() string {
	return "Get subscription"
}

func (c *EventsSubscriptionGetCommand) Help() string {
	helpText := `
Usage: vault events subscriptions get plugin id

  Get details about the specified event subscription.
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *EventsSubscriptionGetCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
}

func (c *EventsSubscriptionGetCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *EventsSubscriptionGetCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *EventsSubscriptionGetCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}
	plugin := args[0]
	id := args[1]

	secret, err := client.Logical().Read("sys/events/subscriptions/" + plugin + "/" + id)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	if secret == nil || secret.Data == nil {
		return 0
	}
	return OutputData(c.UI, secret.Data)
}

type EventsSubscriptionCreateCommand struct {
	*BaseCommand

	bexprFilter string
	config      map[string]string
}

func (c *EventsSubscriptionCreateCommand) Synopsis() string {
	return "Create an external, push subscription to events"
}

func (c *EventsSubscriptionCreateCommand) Help() string {
	helpText := `
Usage: vault events subscriptions create -config=a=b [-config=c=d] [-filter=filterExpression] plugin eventType

  Subscribe to events of the given event type (topic), which may be a glob
  pattern (with "*" treated as a wildcard). The events will routed to the specified plugin, which
  will be created with the config options.
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *EventsSubscriptionCreateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)
	f := set.NewFlagSet("Subscribe Options")
	f.StringVar(&StringVar{
		Name: "filter",
		Usage: `A boolean expression to use to filter events. Only events matching
                the filter will be subscribed to. This is applied after any filtering
                by event type or namespace.`,
		Default: "",
		Target:  &c.bexprFilter,
	})
	f.StringMapVar(&StringMapVar{
		Name: "config",
		Usage: `Map of key-value pairs to pass to the plugin. See plugin
                documentation for specific values.`,
		Default: map[string]string{},
		Target:  &c.config,
	})
	return set
}

func (c *EventsSubscriptionCreateCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *EventsSubscriptionCreateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *EventsSubscriptionCreateCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}
	plugin := args[0]
	eventType := args[1]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	input := map[string]interface{}{
		"config":     c.config,
		"plugin":     plugin,
		"event_type": eventType,
	}
	if c.bexprFilter != "" {
		input["filter"] = c.bexprFilter
	}
	secret, err := client.Logical().Write("sys/events/subscriptions", input)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	return OutputData(c.UI, secret.Data)
}

type EventsSubscriptionDeleteCommand struct {
	*BaseCommand
}

func (c *EventsSubscriptionDeleteCommand) Synopsis() string {
	return "Delete a subscription (unsubscribe)"
}

func (c *EventsSubscriptionDeleteCommand) Help() string {
	helpText := `
Usage: vault events subscriptions delete plugin id

  Delete (unsubscribe) the specified event subscription.
` + c.Flags().Help()
	return strings.TrimSpace(helpText)
}

func (c *EventsSubscriptionDeleteCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *EventsSubscriptionDeleteCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *EventsSubscriptionDeleteCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *EventsSubscriptionDeleteCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	switch {
	case len(args) < 2:
		c.UI.Error(fmt.Sprintf("Not enough arguments (expected 2, got %d)", len(args)))
		return 1
	case len(args) > 2:
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 2, got %d)", len(args)))
		return 1
	}
	plugin := args[0]
	id := args[1]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}
	_, err = client.Logical().Delete("sys/events/subscriptions/" + plugin + "/" + id)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	return 0
}
