package disk

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// DefaultPath is the default path where the Vault token is stored.
const DefaultPath = "~/.vault-token"

type Command struct {
	Path string
}

func (c *Command) Run(args []string) int {
	var path string
	pathDefault := DefaultPath
	if c.Path != "" {
		pathDefault = c.Path
	}

	f := flag.NewFlagSet("token-disk", flag.ContinueOnError)
	f.StringVar(&path, "path", pathDefault, "")
	f.Usage = func() { fmt.Fprintf(os.Stderr, c.Help()+"\n") }
	if err := f.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "\n%s\n", err)
		return 1
	}

	path, err := homedir.Expand(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error expanding directory: %s\n", err)
		return 1
	}

	args = f.Args()
	switch args[0] {
	case "get":
		f, err := os.Open(path)
		if os.IsNotExist(err) {
			return 0
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
		defer f.Close()

		if _, err := io.Copy(os.Stdout, f); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	case "store":
		f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
		defer f.Close()

		if _, err := io.Copy(f, os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	case "erase":
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return 1
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown subcommand: %s\n", args[0])
		return 1
	}

	return 0
}

func (c *Command) Synopsis() string {
	return "Stores Vault tokens on disk"
}

func (c *Command) Help() string {
	helpText := `
Usage: vault token-disk [options] [operation]

  Vault token helper (see vault config "token_helper") that writes
  authenticated tokens to disk unencrypted.

Options:

  -path=path      Path to store the token.

`
	return strings.TrimSpace(helpText)
}
