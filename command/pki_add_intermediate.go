package command

import (
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"strings"
)

var (
	_ cli.Command = (*PKIAddIntermediateCommand)(nil)
	_ cli.Command = (*PKIAddIntermediateCommand)(nil)
)

type PKIAddIntermediateCommand struct {
	*BaseCommand

	flagRootMount   string
	flagMountName   string
	flagMaxLeaseTTL string
	flagCommonName  string
}

func (c *PKIAddIntermediateCommand) Synopsis() string {
	return "Creates a new parseIntermediateArgs CA"
}

func (c *PKIAddIntermediateCommand) Help() string {
	helpText := `
Usage: vault pki add-intermediate [ARGS]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddIntermediateCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "root-mount",
		Target: &c.flagRootMount,
		Usage:  "The name of the root mount to use for signing the parseIntermediateArgs certificate",
	})

	f.StringVar(&StringVar{
		Name:   "mount",
		Target: &c.flagMountName,
		Usage:  "The name of the mount for the root CA. The name must be unique.",
	})

	f.StringVar(&StringVar{
		Name:   "max-ttl",
		Target: &c.flagMaxLeaseTTL,
		Usage:  "The max TTL to use for parseIntermediateArgs CA",
	})

	f.StringVar(&StringVar{
		Name:   "common-name",
		Target: &c.flagCommonName,
		Usage:  "The common name for the parseIntermediateArgs CA",
	})

	return set
}

func (c *PKIAddIntermediateCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *PKIAddIntermediateCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIAddIntermediateCommand) Run(args []string) int {
	return 0
}
