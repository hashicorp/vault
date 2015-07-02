package command

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type SSHCommand struct {
	Meta
}

func (c *SSHCommand) Run(args []string) int {
	var role string
	flags := c.Meta.FlagSet("ssh", FlagSetDefault)
	flags.StringVar(&role, "role", "", "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}
	args = flags.Args()
	if len(args) < 1 {
		c.Ui.Error("ssh expects at least one argument")
		return 2
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 2
	}
	input := strings.Split(args[0], "@")
	username := input[0]
	ip, err := net.ResolveIPAddr("ip4", input[1])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error resolving IP Address: %s", err))
		return 2
	}

	if role == "" {
		data := map[string]interface{}{
			"ip": ip.String(),
		}
		secret, err := client.Logical().Write("ssh/lookup", data)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error finding roles for IP:%s Error:%s", ip.String(), err))
			return 1
		}

		if secret.Data["roles"] == nil {
			c.Ui.Error(fmt.Sprintf("IP '%s' not registered under any role", ip.String()))
			return 1
		}

		if len(secret.Data["roles"].([]interface{})) == 1 {
			role = secret.Data["roles"].([]interface{})[0].(string)
			c.Ui.Output(fmt.Sprintf("Using role[%s]\n", role))
		} else {
			c.Ui.Error(fmt.Sprintf("Multiple roles for IP '%s'. Select one of '%s' using '-role' option", ip, secret.Data["roles"]))
			return 1
		}
	}

	data := map[string]interface{}{
		"username": username,
		"ip":       ip.String(),
	}
	keySecret, err := client.SSH().KeyCreate(role, data)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting key for SSH session:%s", err))
		return 2
	}

	sshDynamicKey := string(keySecret.Data["key"].(string))
	if len(sshDynamicKey) == 0 {
		c.Ui.Error(fmt.Sprintf("Invalid key"))
		return 2
	}
	sshDynamicKeyFileName := fmt.Sprintf("vault_ssh_key_%s_%s.pem", username, ip.String())
	err = ioutil.WriteFile(sshDynamicKeyFileName, []byte(sshDynamicKey), 0600)
	sshBinary, err := exec.LookPath("ssh")
	if err != nil {
		c.Ui.Error("ssh binary not found in PATH\n")
		return 2
	}

	sshEnv := os.Environ()

	sshCmdArgs := []string{"ssh", "-i", sshDynamicKeyFileName, args[0]}

	if err := syscall.Exec(sshBinary, sshCmdArgs, sshEnv); err != nil {
		c.Ui.Error(fmt.Sprintf("Could not launch 'ssh' binary:'%s", err))
		return 2
	}
	return 0
}

type OneTimeKey struct {
	Key string
}

func (c *SSHCommand) Synopsis() string {
	return "Initiate a SSH session"
}

func (c *SSHCommand) Help() string {
	helpText := `
	SSHCommand Help String
	`
	return strings.TrimSpace(helpText)
}
