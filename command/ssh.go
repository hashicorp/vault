package command

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/vault/api"
)

// SSHCommand is a Command that establishes a SSH connection with target by generating a dynamic key
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
	ip, err := net.ResolveIPAddr("ip", input[1])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error resolving IP Address: %s", err))
		return 2
	}

	if role == "" {
		role, err = setDefaultRole(client, ip.String())
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error setting default role: %s", err.Error()))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Using role[%s]\n", role))
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
	sshDynamicKeyFileName := fmt.Sprintf("vault_temp_file_%s_%s", username, ip.String())
	err = ioutil.WriteFile(sshDynamicKeyFileName, []byte(sshDynamicKey), 0600)

	cmd := exec.Command("ssh", "-i", sshDynamicKeyFileName, args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error while running ssh command:%s", err))
	}

	err = os.Remove(sshDynamicKeyFileName)
	if err != nil {
		c.Ui.Error("Error cleaning up") // Intentionally not mentioning the exact error
	}

	err = client.SSH().KeyRevoke(keySecret.LeaseID)
	if err != nil {
		c.Ui.Error("Error cleaning up") // Intentionally not mentioning the exact error
	}

	return 0
}

func setDefaultRole(client *api.Client, ip string) (string, error) {
	data := map[string]interface{}{
		"ip": ip,
	}
	secret, err := client.Logical().Write("ssh/lookup", data)
	if err != nil {
		return "", fmt.Errorf("Error finding roles for IP '%s':%s", ip, err)

	}
	if secret == nil {
		return "", fmt.Errorf("Error finding roles for IP '%s':%s", ip, err)
	}

	if secret.Data["roles"] == nil {
		return "", fmt.Errorf("IP '%s' not registered under any role", ip)
	}

	if len(secret.Data["roles"].([]interface{})) == 1 {
		return secret.Data["roles"].([]interface{})[0].(string), nil
	} else {
		return "", fmt.Errorf("Multiple roles for IP '%s'. Select one of '%s' using '-role' option", ip, secret.Data["roles"])
	}
}

func (c *SSHCommand) Synopsis() string {
	return "Initiate a SSH session"
}

func (c *SSHCommand) Help() string {
	helpText := `
Usage: vault ssh [options] username@ip

  Establishes an SSH connection with the target machine.

  This command generates a dynamic key and uses it to establish an
  SSH connection with the target machine. This operation requires
  that SSH backend is mounted and at least one 'role' be registed
  with vault at priori.

General Options:

  ` + generalOptionsUsage() + `

SSH Options:

  -role			Mention the role to be used to create dynamic key.
  			Each IP is associated with a role. To see the associated
			roles with IP, use "lookup" endpoint. If you are certain that
			there is only one role associated with the IP, you can
			skip mentioning the role. It will be chosen by default.
			If there are no roless associated with the IP, register
			the CIDR block of that IP using the "roles/" endpoint.
`
	return strings.TrimSpace(helpText)
}
