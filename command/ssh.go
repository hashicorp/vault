package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/meta"
	"github.com/mitchellh/mapstructure"
)

// SSHCommand is a Command that establishes a SSH connection
// with target by generating a dynamic key
type SSHCommand struct {
	meta.Meta
}

// Structure to hold the fields returned when asked for a credential from SSHh backend.
type SSHCredentialResp struct {
	KeyType  string `mapstructure:"key_type"`
	Key      string `mapstructure:"key"`
	Username string `mapstructure:"username"`
	IP       string `mapstructure:"ip"`
	Port     string `mapstructure:"port"`
}

func (c *SSHCommand) Run(args []string) int {
	var role, mountPoint, format, userKnownHostsFile, strictHostKeyChecking string
	var noExec bool
	var sshCmdArgs []string
	flags := c.Meta.FlagSet("ssh", meta.FlagSetDefault)
	flags.StringVar(&strictHostKeyChecking, "strict-host-key-checking", "", "")
	flags.StringVar(&userKnownHostsFile, "user-known-hosts-file", "", "")
	flags.StringVar(&format, "format", "table", "")
	flags.StringVar(&role, "role", "", "")
	flags.StringVar(&mountPoint, "mount-point", "ssh", "")
	flags.BoolVar(&noExec, "no-exec", false, "")

	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	// If the flag is already set then it takes the precedence. If the flag is not
	// set, try setting it from env var.
	if os.Getenv("VAULT_SSH_STRICT_HOST_KEY_CHECKING") != "" && strictHostKeyChecking == "" {
		strictHostKeyChecking = os.Getenv("VAULT_SSH_STRICT_HOST_KEY_CHECKING")
	}
	// Assign default value if both flag and env var are not set
	if strictHostKeyChecking == "" {
		strictHostKeyChecking = "ask"
	}

	// If the flag is already set then it takes the precedence. If the flag is not
	// set, try setting it from env var.
	if os.Getenv("VAULT_SSH_USER_KNOWN_HOSTS_FILE") != "" && userKnownHostsFile == "" {
		userKnownHostsFile = os.Getenv("VAULT_SSH_USER_KNOWN_HOSTS_FILE")
	}
	// Assign default value if both flag and env var are not set
	if userKnownHostsFile == "" {
		userKnownHostsFile = "~/.ssh/known_hosts"
	}

	args = flags.Args()
	if len(args) < 1 {
		c.Ui.Error("ssh expects at least one argument")
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %v", err))
		return 1
	}

	// split the parameter username@ip
	input := strings.Split(args[0], "@")
	var username string
	var ipAddr string

	// If only IP is mentioned and username is skipped, assume username to
	// be the current username. Vault SSH role's default username could have
	// been used, but in order to retain the consistency with SSH command,
	// current username is employed.
	if len(input) == 1 {
		u, err := user.Current()
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error fetching username: %v", err))
			return 1
		}
		username = u.Username
		ipAddr = input[0]
	} else if len(input) == 2 {
		username = input[0]
		ipAddr = input[1]
	} else {
		c.Ui.Error(fmt.Sprintf("Invalid parameter: %q", args[0]))
		return 1
	}

	// Resolving domain names to IP address on the client side.
	// Vault only deals with IP addresses.
	ip, err := net.ResolveIPAddr("ip", ipAddr)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error resolving IP Address: %v", err))
		return 1
	}

	// Credentials are generated only against a registered role. If user
	// does not specify a role with the SSH command, then lookup API is used
	// to fetch all the roles with which this IP is associated. If there is
	// only one role associated with it, use it to establish the connection.
	if role == "" {
		role, err = c.defaultRole(mountPoint, ip.String())
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error choosing role: %v", err))
			return 1
		}
		// Print the default role chosen so that user knows the role name
		// if something doesn't work. If the role chosen is not allowed to
		// be used by the user (ACL enforcement), then user should see an
		// error message accordingly.
		c.Ui.Output(fmt.Sprintf("Vault SSH: Role: %q", role))
	}

	data := map[string]interface{}{
		"username": username,
		"ip":       ip.String(),
	}

	keySecret, err := client.SSHWithMountPoint(mountPoint).Credential(role, data)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting key for SSH session: %v", err))
		return 1
	}

	// if no-exec was chosen, just print out the secret and return.
	if noExec {
		return OutputSecret(c.Ui, format, keySecret)
	}

	// Port comes back as a json.Number which mapstructure doesn't like, so convert it
	if keySecret.Data["port"] != nil {
		keySecret.Data["port"] = keySecret.Data["port"].(json.Number).String()
	}
	var resp SSHCredentialResp
	if err := mapstructure.Decode(keySecret.Data, &resp); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing the credential response: %v", err))
		return 1
	}

	if resp.KeyType == ssh.KeyTypeDynamic {
		if len(resp.Key) == 0 {
			c.Ui.Error(fmt.Sprintf("Invalid key"))
			return 1
		}
		sshDynamicKeyFile, err := ioutil.TempFile("", fmt.Sprintf("vault_ssh_%s_%s_", username, ip.String()))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error creating temporary file: %v", err))
			return 1
		}

		// Ensure that we delete the temporary file
		defer os.Remove(sshDynamicKeyFile.Name())

		if err = ioutil.WriteFile(sshDynamicKeyFile.Name(),
			[]byte(resp.Key), 0600); err != nil {
			c.Ui.Error(fmt.Sprintf("Error storing the dynamic key into the temporary file: %v", err))
			return 1
		}
		sshCmdArgs = append(sshCmdArgs, []string{"-i", sshDynamicKeyFile.Name()}...)

	} else if resp.KeyType == ssh.KeyTypeOTP {
		// Check if the application 'sshpass' is installed in the client machine.
		// If it is then, use it to automate typing in OTP to the prompt. Unfortunately,
		// it was not possible to automate it without a third-party application, with
		// only the Go libraries.
		// Feel free to try and remove this dependency.
		sshpassPath, err := exec.LookPath("sshpass")
		if err == nil {
			sshCmdArgs = append(sshCmdArgs, []string{"-p", string(resp.Key), "ssh", "-o UserKnownHostsFile=" + userKnownHostsFile, "-o StrictHostKeyChecking=" + strictHostKeyChecking, "-p", resp.Port, username + "@" + ip.String()}...)
			if len(args) > 1 {
				sshCmdArgs = append(sshCmdArgs, args[1:]...)
			}
			sshCmd := exec.Command(sshpassPath, sshCmdArgs...)
			sshCmd.Stdin = os.Stdin
			sshCmd.Stdout = os.Stdout
			err = sshCmd.Run()
			if err != nil {
				c.Ui.Error(fmt.Sprintf("Failed to establish SSH connection: %q", err))
			}
			return 0
		}
		c.Ui.Output("OTP for the session is " + resp.Key)
		c.Ui.Output("[Note: Install 'sshpass' to automate typing in OTP]")
	}
	sshCmdArgs = append(sshCmdArgs, []string{"-o UserKnownHostsFile=" + userKnownHostsFile, "-o StrictHostKeyChecking=" + strictHostKeyChecking, "-p", resp.Port, username + "@" + ip.String()}...)
	if len(args) > 1 {
		sshCmdArgs = append(sshCmdArgs, args[1:]...)
	}

	sshCmd := exec.Command("ssh", sshCmdArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout

	// Running the command as a separate command. The reason for using exec.Command instead
	// of using crypto/ssh package is that, this way, user can have the same feeling of
	// connecting to remote hosts with the ssh command. Package crypto/ssh did not have a way
	// to establish an independent session like this.
	err = sshCmd.Run()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error while running ssh command: %q", err))
	}

	// If the session established was longer than the lease expiry, the secret
	// might have been revoked already. If not, then revoke it. Since the key
	// file is deleted and since user doesn't know the credential anymore, there
	// is not point in Vault maintaining this secret anymore. Every time the command
	// is run, a fresh credential is generated anyways.
	err = client.Sys().Revoke(keySecret.LeaseID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error revoking the key: %q", err))
	}

	return 0
}

// If user did not provide the role with which SSH connection has
// to be established and if there is only one role associated with
// the IP, it is used by default.
func (c *SSHCommand) defaultRole(mountPoint, ip string) (string, error) {
	data := map[string]interface{}{
		"ip": ip,
	}
	client, err := c.Client()
	if err != nil {
		return "", err
	}
	secret, err := client.Logical().Write(mountPoint+"/lookup", data)
	if err != nil {
		return "", fmt.Errorf("Error finding roles for IP %q: %q", ip, err)

	}
	if secret == nil {
		return "", fmt.Errorf("Error finding roles for IP %q: %q", ip, err)
	}

	if secret.Data["roles"] == nil {
		return "", fmt.Errorf("No matching roles found for IP %q", ip)
	}

	if len(secret.Data["roles"].([]interface{})) == 1 {
		return secret.Data["roles"].([]interface{})[0].(string), nil
	} else {
		var roleNames string
		for _, item := range secret.Data["roles"].([]interface{}) {
			roleNames += item.(string) + ", "
		}
		roleNames = strings.TrimRight(roleNames, ", ")
		return "", fmt.Errorf("Roles:%q. "+`
		Multiple roles are registered for this IP.
		Select a role using '-role' option.
		Note that all roles may not be permitted, based on ACLs.`, roleNames)
	}
}

func (c *SSHCommand) Synopsis() string {
	return "Initiate an SSH session"
}

func (c *SSHCommand) Help() string {
	helpText := `
Usage: vault ssh [options] username@ip

  Establishes an SSH connection with the target machine.

  This command generates a key and uses it to establish an SSH
  connection with the target machine. This operation requires
  that the SSH backend is mounted and at least one 'role' is
  registered with Vault beforehand.

  For setting up SSH backends with one-time-passwords, installation
  of vault-ssh-helper or a compatible agent on target machines
  is required. See [https://github.com/hashicorp/vault-ssh-agent].

General Options:
` + meta.GeneralOptionsUsage() + `
SSH Options:

	-role				Role to be used to create the key.
					Each IP is associated with a role. To see the associated
					roles with IP, use "lookup" endpoint. If you are certain
					that there is only one role associated with the IP, you can
					skip mentioning the role. It will be chosen by default.  If
					there are no roles associated with the IP, register the
					CIDR block of that IP using the "roles/" endpoint.

	-no-exec			Shows the credentials but does not establish connection.

	-mount-point			Mount point of SSH backend. If the backend is mounted at
					'ssh', which is the default as well, this parameter can be
					skipped.

	-format				If no-exec option is enabled, then the credentials will be
					printed out and SSH connection will not be established. The
					format of the output can be 'json' or 'table'. JSON output
					is useful when writing scripts. Default is 'table'.

	-strict-host-key-checking	This option corresponds to StrictHostKeyChecking of SSH configuration.
					If 'sshpass' is employed to enable automated login, then if host key
					is not "known" to the client, 'vault ssh' command will fail. Set this
					option to "no" to bypass the host key checking. Defaults to "ask".
					Can also be specified with VAULT_SSH_STRICT_HOST_KEY_CHECKING environment
					variable.

	-user-known-hosts-file		This option corresponds to UserKnownHostsFile of SSH configuration.
					Assigns the file to use for storing the host keys. If this option is
					set to "/dev/null" along with "-strict-host-key-checking=no", both
					warnings and host key checking can be avoided while establishing the
					connection. Defaults to "~/.ssh/known_hosts". Can also be specified
					with VAULT_SSH_USER_KNOWN_HOSTS_FILE environment variable.
`
	return strings.TrimSpace(helpText)
}
