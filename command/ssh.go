package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type SshCommand struct {
	Meta
}

func (c *SshCommand) Run(args []string) int {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

	log.Printf("Vishal: command.ssh.Run returned! OTK:%#v\n", sshOneTimeKey)
	err = ioutil.WriteFile("./vault_ssh_otk_"+args[0]+".pem", []byte(sshOneTimeKey.Key), 0400)
	//if sshOneTimeKey is empty, fail
	//Establish a session directly from client to the target using the one time key received without making the vault server the middle guy:w
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		log.Printf("ssh binary not found in PATH\n")
	}

	sshEnv := os.Environ()

	sshCmdArgs := []string{"ssh", "-i", "vault_ssh_otk_" + args[0] + ".pem", "vishal@localhost"}
	defer os.Remove("vault_ssh_otk_" + args[0] + ".pem")

	if err := syscall.Exec(sshBinary, sshCmdArgs, sshEnv); err != nil {
		log.Printf("Execution failed: sshCommand: " + err.Error())
	}
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
