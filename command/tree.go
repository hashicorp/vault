package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/meta"
)

// TreeCommand is a Command that shows keys in tree format
type TreeCommand struct {
	meta.Meta
}

func (c *TreeCommand) Run(args []string) int {
	var flags *flag.FlagSet
	flags = c.Meta.FlagSet("list", meta.FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	args = flags.Args()
	if len(args) != 1 || len(args[0]) == 0 {
		c.Ui.Error("tree expects one argument")
		flags.Usage()
		return 1
	}

	path := args[0]
	if path[0] == '/' {
		path = path[1:]
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

  c.ShowKeys(path, 0)

  return 0
}

func (c *TreeCommand) ShowKeys(path string, level int) int {
	var err error
	var secret *api.Secret

  client, err := c.Client()
  if err != nil {
    c.Ui.Error(fmt.Sprintf(
      "Error initializing client: %s", err))
    return 2
  }

	secret, err = client.Logical().List(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error reading %s: %s", path, err))
		return 1
	}
	if secret == nil {
		// c.Ui.Error(fmt.Sprintf(
		// 	"No value found at %s", path))
		return 1
	}

  if keys, ok := secret.Data["keys"].([]interface{}); ok {
    for _, k := range keys {
      key := k.(string)
      if level > 0 {
        for i := 0; i < level - 1; i++ {
          fmt.Printf("    ")
        }
        fmt.Printf("├── ")
      }
      fmt.Printf("%s\n", key)

      c.ShowKeys(path + key, level + 1)
    }
  }

  return 0
}

func (c *TreeCommand) Synopsis() string {
	return "Show data or secrets in Vault in tree format"
}

func (c *TreeCommand) Help() string {
	helpText :=
		`
Usage: vault tree [options] path

  Show data from Vault in tree format.

  Retrieve a listing of available data. The data returned, if any, is backend-
  and endpoint-specific.

General Options:
` + meta.GeneralOptionsUsage()
	return strings.TrimSpace(helpText)
}
