package command

import (
	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/cli"
)

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	VersionInfo *version.VersionInfo
	Ui          cli.Ui
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(_ []string) int {
	out := c.VersionInfo.FullVersionNumber(true)
	if version.CgoEnabled {
		out += " (cgo)"
	}
	c.Ui.Output(out)
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the Vault version"
}
