package command

import (
	"github.com/mitchellh/cli"
)

// Deprecation
// TODO: remove in 0.9.0

var _ cli.Command = (*PoliciesDeprecatedCommand)(nil)

type PoliciesDeprecatedCommand struct {
	*BaseCommand
}

func (c *PoliciesDeprecatedCommand) Synopsis() string { return "" }

func (c *PoliciesDeprecatedCommand) Help() string {
	return (&PolicyListCommand{
		BaseCommand: c.BaseCommand,
	}).Help()
}

func (c *PoliciesDeprecatedCommand) Run(args []string) int {
	oargs := args

	f := c.flagSet(FlagSetHTTP)
	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	args = f.Args()

	// Got an arg, this is trying to read a policy
	if len(args) > 0 {
		return (&PolicyReadCommand{
			BaseCommand: &BaseCommand{
				UI:          c.UI,
				client:      c.client,
				tokenHelper: c.tokenHelper,
				flagAddress: c.flagAddress,
			},
		}).Run(oargs)
	}

	// No args, probably ran "vault policies" and we want to translate that to
	// "vault policy list"
	return (&PolicyListCommand{
		BaseCommand: &BaseCommand{
			UI:          c.UI,
			client:      c.client,
			tokenHelper: c.tokenHelper,
			flagAddress: c.flagAddress,
		},
	}).Run(oargs)
}
