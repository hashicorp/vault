package command

import (
	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"strings"
)

var (
	_ cli.Command             = (*PKIAddRootCommand)(nil)
	_ cli.CommandAutocomplete = (*PKIAddRootCommand)(nil)
)

type PKIAddRootCommand struct {
	*BaseCommand

	flagMountName           string
	flagMaxLeaseTTL         string
	flagCommonName          string
	flagTTL                 string
	flagIssuingCertURLs     string
	flagCRLDistributionURLs string
}

func (c *PKIAddRootCommand) Synopsis() string {
	return "Creates a new root CA"
}

func (c *PKIAddRootCommand) Help() string {
	helpText := `
Usage: vault pki add-root [ARGS]
` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *PKIAddRootCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP)
	f := set.NewFlagSet("Command Options")

	f.StringVar(&StringVar{
		Name:   "mount",
		Target: &c.flagMountName,
		Usage:  "The name of the mount for the root CA. The name must be unique.",
	})

	f.StringVar(&StringVar{
		Name:   "max-ttl",
		Target: &c.flagMaxLeaseTTL,
		Usage:  "The max TTL to use for the root CA",
	})

	f.StringVar(&StringVar{
		Name:   "common-name",
		Target: &c.flagCommonName,
		Usage:  "The common name for the root certificate",
	})

	f.StringVar(&StringVar{
		Name:   "ttl",
		Target: &c.flagTTL,
		Usage:  "The TTL of the root certificate",
	})

	f.StringVar(&StringVar{
		Name:   "issuing-certificates",
		Target: &c.flagIssuingCertURLs,
		Usage:  "The URLs for the Issuing Certificate",
	})

	f.StringVar(&StringVar{
		Name:   "crl-distribution-points",
		Target: &c.flagCRLDistributionURLs,
		Usage:  "The URLs for the CRL Distribution Points",
	})

	return set
}

func (c *PKIAddRootCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictAnything
}

func (c *PKIAddRootCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *PKIAddRootCommand) Run(args []string) int {
	return 0
}
