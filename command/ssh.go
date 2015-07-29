package command

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/hashicorp/vault/builtin/logical/ssh"
)

// SSHCommand is a Command that establishes a SSH connection
// with target by generating a dynamic key
type SSHCommand struct {
	Meta
}

func (c *SSHCommand) Run(args []string) int {
	var role, port, path string
	var noExec bool
	var sshCmdArgs []string
	var sshDynamicKeyFileName string
	flags := c.Meta.FlagSet("ssh", FlagSetDefault)
	flags.StringVar(&role, "role", "", "")
	flags.StringVar(&port, "port", "22", "")
	flags.StringVar(&path, "path", "ssh", "")
	flags.BoolVar(&noExec, "no-exec", false, "")

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
	var username string
	var ipAddr string
	if len(input) == 1 {
		u, err := user.Current()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error fetching username: '%s'", err))
		}
		username = u.Username
		ipAddr = input[0]
	} else if len(input) == 2 {
		username = input[0]
		ipAddr = input[1]
	} else {
		c.Ui.Error(fmt.Sprintf("Invalid parameter: %s", args[0]))
		return 2
	}

	ip, err := net.ResolveIPAddr("ip", ipAddr)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error resolving IP Address: %s", err))
		return 2
	}

	if role == "" {
		role, err = c.defaultRole(path, ip.String())
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error setting default role: '%s'", err))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Vault SSH: Role:'%s'\n", role))
	}

	data := map[string]interface{}{
		"username": username,
		"ip":       ip.String(),
	}

	keySecret, err := client.SSH(path).Credential(role, data)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting key for SSH session:%s", err))
		return 2
	}

	if noExec {
		c.Ui.Output(fmt.Sprintf("IP:%s\nUsername: %s\nKey:%s\n", ip.String(), username, keySecret.Data["key"]))
		return 0
	}

	if keySecret.Data["key_type"].(string) == ssh.KeyTypeDynamic {
		sshDynamicKey := string(keySecret.Data["key"].(string))
		if len(sshDynamicKey) == 0 {
			c.Ui.Error(fmt.Sprintf("Invalid key"))
			return 2
		}
		sshDynamicKeyFileName = fmt.Sprintf("vault_temp_file_%s_%s", username, ip.String())
		err = ioutil.WriteFile(sshDynamicKeyFileName, []byte(sshDynamicKey), 0600)
		sshCmdArgs = append(sshCmdArgs, []string{"-i", sshDynamicKeyFileName}...)

	} else if keySecret.Data["key_type"].(string) == ssh.KeyTypeOTP {
		c.Ui.Output(fmt.Sprintf("OTP for the session is %s\n", string(keySecret.Data["key"].(string))))
	} else {
		c.Ui.Error("Error creating key")
	}
	sshCmdArgs = append(sshCmdArgs, []string{"-p", port}...)
	sshCmdArgs = append(sshCmdArgs, args...)

	sshCmd := exec.Command("ssh", sshCmdArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout

	err = sshCmd.Run()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error while running ssh command:%s", err))
	}

	if keySecret.Data["key_type"].(string) == ssh.KeyTypeDynamic {
		err = os.Remove(sshDynamicKeyFileName)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error deleting key file: %s", err))
		}
	}

	err = client.Sys().Revoke(keySecret.LeaseID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error revoking the key: %s", err))
	}

	return 0
}

// If user did not provide the role with which SSH connection has
// to be established and if there is only one role associated with
// the IP, it is used by default.
func (c *SSHCommand) defaultRole(path, ip string) (string, error) {
	data := map[string]interface{}{
		"ip": ip,
	}
	client, err := c.Client()
	if err != nil {
		return "", err
	}
	secret, err := client.Logical().Write(path+"/lookup", data)
	if err != nil {
		return "", fmt.Errorf("Error finding roles for IP '%s':%s", ip, err)

	}
	if secret == nil {
		return "", fmt.Errorf("Error finding roles for IP '%s':%s", ip, err)
	}

	if secret.Data["roles"] == nil {
		return "", fmt.Errorf("No matching roles found for IP '%s'", ip)
	}

	if len(secret.Data["roles"].([]interface{})) == 1 {
		return secret.Data["roles"].([]interface{})[0].(string), nil
	} else {
		var roleNames string
		for _, item := range secret.Data["roles"].([]interface{}) {
			roleNames += item.(string) + ", "
		}
		roleNames = strings.TrimRight(roleNames, ", ")
		return "", fmt.Errorf("IP '%s' has multiple roles.\nSelect a role using '-role' option.\nPossible roles: [%s]\nNote that all roles may not be permitted, based on ACLs.", ip, roleNames)
	}
}

func (c *SSHCommand) Synopsis() string {
	return "Initiate a SSH session"
}

func (c *SSHCommand) Help() string {
	helpText := `
Usage: vault ssh [options] username@ip

  Establishes an SSH connection with the target machine.

  This command generates a key and uses it to establish an SSH
  connection with the target machine. This operation requires
  that SSH backend is mounted and at least one 'role' be registed
  with vault at priori.

  For setting up SSH backends with one-time-passwords, installation
  of agent in target machines is required. 
  See [https://github.com/hashicorp/vault-ssh-agent]

General Options:

  ` + generalOptionsUsage() + `

SSH Options:

  -role                 Role to be used to create the key.
  			Each IP is associated with a role. To see the associated
			roles with IP, use "lookup" endpoint. If you are certain that
			there is only one role associated with the IP, you can
			skip mentioning the role. It will be chosen by default.
			If there are no roles associated with the IP, register
			the CIDR block of that IP using the "roles/" endpoint.

  -port                 Port number to use for SSH connection. This defaults to port 22.

  -no-exec		Shows the credentials but does not establish connection.

  -path			Mount point of SSH backend. If the backend is mounted at
  			'ssh', which is the default as well, this parameter can
			be skipped.
`
	return strings.TrimSpace(helpText)
}
