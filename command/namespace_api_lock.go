package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
)

var (
	_ cli.Command             = (*NamespaceAPILockCommand)(nil)
	_ cli.CommandAutocomplete = (*NamespaceAPILockCommand)(nil)
)

type NamespaceAPILockCommand struct {
	*BaseCommand
}

func (c *NamespaceAPILockCommand) Synopsis() string {
	return "Lock the API for particular namespaces"
}

func (c *NamespaceAPILockCommand) Help() string {
	helpText := `
Usage: vault namespace lock PATH

	Lock the current namespace, and all descendants:

		$ vault namespace lock

	Lock a child namespace, and all of its descendants (e.g. ns1/ns2/):

		$ vault namespace lock ns1/ns2

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *NamespaceAPILockCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP | FlagSetOutputFormat)
}

func (c *NamespaceAPILockCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultNamespaces()
}

func (c *NamespaceAPILockCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *NamespaceAPILockCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()
	if len(args) > 1 {
		c.UI.Error(fmt.Sprintf("Too many arguments (expected 0 or 1, got %d)", len(args)))
		return 1
	}

	// current namespace is already encoded in the :client:
	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 2
	}

	optionalChildNSPath := ""
	if len(args) == 1 {
		optionalChildNSPath = fmt.Sprintf("/%s", namespace.Canonicalize(args[0]))
	}

	resp, err := client.Logical().Write(fmt.Sprintf("sys/namespaces/api-lock/lock%s", optionalChildNSPath), nil)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error locking namespace: %v", err))
		return 2
	}

	return OutputSecret(c.UI, resp)
}
