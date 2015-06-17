package command

import (
	"fmt"
	"log"
	"strings"
)

type SshCommand struct {
	Meta
}

func (c *SshCommand) Run(args []string) int {
	log.Printf("Vishal: SshCommand.Run: args:%#v len(args):%d\n", args, len(args))
	flags := c.Meta.FlagSet("ssh", FlagSetDefault)
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 2
	}
	//if len(args) < 3, fail
	log.Printf("Vishal: sshCommand.Run: args[0]: %#v\n", args[0])
	sshOneTimeKey, err := client.Sys().Ssh(args[0])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting one-time-key for establishing SSH session", err))
		return 2
	}

	log.Printf("Vishal: client.Sys().Ssh() returned! OTK:%#v\n", sshOneTimeKey)
	//if sshOneTimeKey is empty, fail
	//Establish a session directly from client to the target using the one time key received without making the vault server the middle guy:w
	return 0
}

func (c *SshCommand) Synopsis() string {
	return "Initiate a SSH session"
}

func (c *SshCommand) Help() string {
	helpText := `
	SshCommand Help String
	`
	return strings.TrimSpace(helpText)
}
