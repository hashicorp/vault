package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*TokenListCommand)(nil)
	_ cli.CommandAutocomplete = (*TokenListCommand)(nil)
)

type TokenListCommand struct {
	*BaseCommand
}

func (c *TokenListCommand) Synopsis() string {
	return "List information about each token"
}

func (c *TokenListCommand) Help() string {
	helpText := `
Usage: vault token list

	Displays information about each token.

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *TokenListCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *TokenListCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *TokenListCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *TokenListCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 0 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0, got %d)", len(args)))
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	tokens, err := client.Auth().Token().ListTokens()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error listing tokens: %s", err))
		return 2
	}

	switch Format(c.UI) {
	case "table":
		c.UI.Output(tableOutput(c.simpleMounts(tokens), nil))
		return 0
	default:
		return OutputData(c.UI, tokens)
	}
}

func (c *TokenListCommand) simpleMounts(tokens []*api.TokenInfo) []string {
	out := []string{"Accessor | Display Name | Role | TTL | Policies"}
	for _, token := range tokens {
		tokenTTL, err := parseutil.ParseDurationSecond(token.TTL)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing TTL: %s", err))
		}
		out = append(
			out,
			fmt.Sprintf(
				"%s | %s | %s | %s | %s",
				token.Accessor,
				token.DisplayName,
				token.Role,
				tokenTTL,
				token.Policies,
			),
		)
	}

	return out
}
