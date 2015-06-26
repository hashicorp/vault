package command

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
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
	if len(args) < 1 {
		c.Ui.Error(fmt.Sprintf("Insufficient arguments"))
		return 2
	}
	log.Printf("Vishal: sshCommand.Run: args[0]: %#v\n", args[0])
	input := strings.Split(args[0], "@")
	username := input[0]
	ipAddr, err := net.ResolveIPAddr("ip4", input[1])
	log.Printf("Vishal: ssh.Ssh ipAddr_resolved: %#v\n", ipAddr.String())
	data := map[string]interface{}{
		"username": username,
		"ip":       ipAddr.String(),
	}

	keySecret, err := client.Ssh().KeyCreate(data)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting key for establishing SSH session", err))
		return 2
	}
	sshOneTimeKey := string(keySecret.Data["key"].(string))
	log.Printf("Vishal: command.ssh.Run returned! len(key):%d\n", len(sshOneTimeKey))
	ag := strings.Split(args[0], "@")
	sshOtkFileName := "vault_ssh_otk_" + ag[0] + "_" + ag[1] + ".pem"
	err = ioutil.WriteFile(sshOtkFileName, []byte(sshOneTimeKey), 0400)
	//if sshOneTimeKey is empty, fail
	//Establish a session directly from client to the target using the one time key received without making the vault server the middle guy:w
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		log.Printf("ssh binary not found in PATH\n")
	}

	sshEnv := os.Environ()

	sshNew := "ssh -i " + sshOtkFileName + " " + args[0]
	log.Printf("Vishal: sshNew:%#v\n", sshNew)
	sshCmdArgs := []string{"ssh", "-i", sshOtkFileName, args[0]}
	//defer os.Remove("vault_ssh_otk_" + args[0] + ".pem")

	if err := syscall.Exec(sshBinary, sshCmdArgs, sshEnv); err != nil {
		log.Printf("Execution failed: sshCommand: " + err.Error())
	}
	return 0
}

type OneTimeKey struct {
	Key string
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
