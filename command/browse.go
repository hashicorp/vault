package command

import (
	"fmt"
	"github.com/posener/complete"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/nsf/termbox-go"
)

var _ cli.Command = (*BrowseCommand)(nil)

// back buffer width and height
var bbw, bbh int

// back buffer
var backbuf []termbox.Cell

type BrowseCommand struct {
	*BaseCommand
}

func (c *BrowseCommand) Synopsis() string {
	return "Interact with Vault's Key-Value storage"
}

func (c *BrowseCommand) Help() string {
	helpText := `
Usage: vault browse

  This command has subcommands for interacting with Vault's key-value
  store. Here are some simple examples, and more detailed examples are
  available in the subcommands or the documentation.

  Create or update the key named "foo" in the "secret" mount with the value
  "bar=baz":

      $ vault browse

  Read this value back:

      $ vault browse

  Get metadata for the key:

      $ vault browse
	  
  Get a specific version of the key:

      $ vault browse

  Please see the individual subcommand help for detailed usage information.
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *BrowseCommand) Flags() *FlagSets {
	return c.flagSet(FlagSetHTTP)
}

func (c *BrowseCommand) AutocompleteArgs() complete.Predictor {
	return c.PredictVaultFolders()
}

func (c *BrowseCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *BrowseCommand) Run(args []string) int {
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

	path := ensureTrailingSlash(sanitizePath(args[0]))

	tui, err := NewTerminalUI(client, path)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not initialize terminal UI at %s: %s", path, err.Error()))
		return 2
	}

	err = termbox.Init()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not initialize termbox %s", err.Error()))
		return 2
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc)
	reallocBackBuffer(termbox.Size())

	tui.draw()

mainloop:
	for {

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc {
				break mainloop
			}

			tui.handleInput(ev.Key)
		case termbox.EventResize:
			reallocBackBuffer(ev.Width, ev.Height)
		}

		tui.draw()
	}

	return 0
}

func reallocBackBuffer(w, h int) {
	bbw, bbh = w, h
	backbuf = make([]termbox.Cell, w*h)
}
