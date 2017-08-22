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

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/meta"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// SSHCommand is a Command that establishes a SSH connection with target by
// generating a dynamic key
type SSHCommand struct {
	meta.Meta

	// API
	client    *api.Client
	sshClient *api.SSH

	// Common options
	mode       string
	noExec     bool
	format     string
	mountPoint string
	role       string
	username   string
	ip         string
	sshArgs    []string

	// Key options
	strictHostKeyChecking string
	userKnownHostsFile    string

	// SSH CA backend specific options
	publicKeyPath     string
	privateKeyPath    string
	hostKeyMountPoint string
	hostKeyHostnames  string
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

	flags := c.Meta.FlagSet("ssh", meta.FlagSetDefault)

	envOrDefault := func(key string, def string) string {
		if k := os.Getenv(key); k != "" {
			return k
		}
		return def
	}

	expandPath := func(p string) string {
		e, err := homedir.Expand(p)
		if err != nil {
			return p
		}
		return e
	}

	// Common options
	flags.StringVar(&c.mode, "mode", "", "")
	flags.BoolVar(&c.noExec, "no-exec", false, "")
	flags.StringVar(&c.format, "format", "table", "")
	flags.StringVar(&c.mountPoint, "mount-point", "ssh", "")
	flags.StringVar(&c.role, "role", "", "")

	// Key options
	flags.StringVar(&c.strictHostKeyChecking, "strict-host-key-checking",
		envOrDefault("VAULT_SSH_STRICT_HOST_KEY_CHECKING", "ask"), "")
	flags.StringVar(&c.userKnownHostsFile, "user-known-hosts-file",
		envOrDefault("VAULT_SSH_USER_KNOWN_HOSTS_FILE", expandPath("~/.ssh/known_hosts")), "")

	// CA-specific options
	flags.StringVar(&c.publicKeyPath, "public-key-path",
		expandPath("~/.ssh/id_rsa.pub"), "")
	flags.StringVar(&c.privateKeyPath, "private-key-path",
		expandPath("~/.ssh/id_rsa"), "")
	flags.StringVar(&c.hostKeyMountPoint, "host-key-mount-point", "", "")
	flags.StringVar(&c.hostKeyHostnames, "host-key-hostnames", "*", "")

	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
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
	c.client = client
	c.sshClient = client.SSHWithMountPoint(c.mountPoint)

	// Extract the username and IP.
	c.username, c.ip, err = c.userAndIP(args[0])
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing user and IP: %s", err))
		return 1
	}

	// The rest of the args are ssh args
	if len(args) > 1 {
		c.sshArgs = args[1:]
	}

	// Credentials are generated only against a registered role. If user
	// does not specify a role with the SSH command, then lookup API is used
	// to fetch all the roles with which this IP is associated. If there is
	// only one role associated with it, use it to establish the connection.
	//
	// TODO: remove in 0.9.0, convert to validation error
	if c.role == "" {
		c.Ui.Warn("" +
			"WARNING: No -role specified. Use -role to tell Vault which ssh role\n" +
			"to use for authentication. In the future, you will need to tell Vault\n" +
			"which role to use. For now, Vault will attempt to guess based on a\n" +
			"the API response.")

		role, err := c.defaultRole(c.mountPoint, c.ip)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error choosing role: %v", err))
			return 1
		}
		// Print the default role chosen so that user knows the role name
		// if something doesn't work. If the role chosen is not allowed to
		// be used by the user (ACL enforcement), then user should see an
		// error message accordingly.
		c.Ui.Output(fmt.Sprintf("Vault SSH: Role: %q", role))
		c.role = role
	}

	// If no mode was given, perform the old-school lookup. Keep this now for
	// backwards-compatability, but print a warning.
	//
	// TODO: remove in 0.9.0, convert to validation error
	if c.mode == "" {
		c.Ui.Warn("" +
			"WARNING: No -mode specified. Use -mode to tell Vault which ssh\n" +
			"authentication mode to use. In the future, you will need to tell\n" +
			"Vault which mode to use. For now, Vault will attempt to guess based\n" +
			"on the API response. This guess involves creating a temporary\n" +
			"credential, reading its type, and then revoking it. To reduce the\n" +
			"number of API calls and surface area, specify -mode directly.")
		secret, cred, err := c.generateCredential()
		if err != nil {
			// This is _very_ hacky, but is the only sane backwards-compatible way
			// to do this. If the error is "key type unknown", we just assume the
			// type is "ca". In the future, mode will be required as an option.
			if strings.Contains(err.Error(), "key type unknown") {
				c.mode = ssh.KeyTypeCA
			} else {
				c.Ui.Error(fmt.Sprintf("Error getting credential: %s", err))
				return 1
			}
		} else {
			c.mode = cred.KeyType
		}

		// Revoke the secret, since the child functions will generate their own
		// credential. Users wishing to avoid this should specify -mode.
		if secret != nil {
			if err := c.client.Sys().Revoke(secret.LeaseID); err != nil {
				c.Ui.Warn(fmt.Sprintf("Failed to revoke temporary key: %s", err))
			}
		}
	}

	switch strings.ToLower(c.mode) {
	case ssh.KeyTypeCA:
		if err := c.handleTypeCA(); err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
	case ssh.KeyTypeOTP:
		if err := c.handleTypeOTP(); err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
	case ssh.KeyTypeDynamic:
		if err := c.handleTypeDynamic(); err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
	default:
		c.Ui.Error(fmt.Sprintf("Unknown SSH mode: %s", c.mode))
		return 1
	}

	return 0
}

// handleTypeCA is used to handle SSH logins using the "CA" key type.
func (c *SSHCommand) handleTypeCA() error {
	// Read the key from disk
	publicKey, err := ioutil.ReadFile(c.publicKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to read public key")
	}

	// Attempt to sign the public key
	secret, err := c.sshClient.SignKey(c.role, map[string]interface{}{
		// WARNING: publicKey is []byte, which is b64 encoded on JSON upload. We
		// have to convert it to a string. SV lost many hours to this...
		"public_key":       string(publicKey),
		"valid_principals": c.username,
		"cert_type":        "user",

		// TODO: let the user configure these. In the interim, if users want to
		// customize these values, they can produce the key themselves.
		"extensions": map[string]string{
			"permit-X11-forwarding":   "",
			"permit-agent-forwarding": "",
			"permit-port-forwarding":  "",
			"permit-pty":              "",
			"permit-user-rc":          "",
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to sign public key")
	}
	if secret == nil || secret.Data == nil {
		return fmt.Errorf("client signing returned empty credentials")
	}

	// Handle no-exec
	if c.noExec {
		// This is hacky, but OutputSecret returns an int, not an error :(
		if i := OutputSecret(c.Ui, c.format, secret); i != 0 {
			return fmt.Errorf("an error occurred outputting the secret")
		}
		return nil
	}

	// Extract public key
	key, ok := secret.Data["signed_key"].(string)
	if !ok {
		return fmt.Errorf("missing signed key")
	}

	// Capture the current value - this could be overwritten later if the user
	// enabled host key signing verification.
	userKnownHostsFile := c.userKnownHostsFile
	strictHostKeyChecking := c.strictHostKeyChecking

	// Handle host key signing verification. If the user specified a mount point,
	// download the public key, trust it with the given domains, and use that
	// instead of the user's regular known_hosts file.
	if c.hostKeyMountPoint != "" {
		secret, err := c.client.Logical().Read(c.hostKeyMountPoint + "/config/ca")
		if err != nil {
			return errors.Wrap(err, "failed to get host signing key")
		}
		if secret == nil || secret.Data == nil {
			return fmt.Errorf("missing host signing key")
		}
		publicKey, ok := secret.Data["public_key"].(string)
		if !ok {
			return fmt.Errorf("host signing key is empty")
		}

		// Write the known_hosts file
		name := fmt.Sprintf("vault_ssh_ca_known_hosts_%s_%s", c.username, c.ip)
		data := fmt.Sprintf("@cert-authority %s %s", c.hostKeyHostnames, publicKey)
		knownHosts, err, closer := c.writeTemporaryFile(name, []byte(data), 0644)
		defer closer()
		if err != nil {
			return errors.Wrap(err, "failed to write host public key")
		}

		// Update the variables
		userKnownHostsFile = knownHosts
		strictHostKeyChecking = "yes"
	}

	// Write the signed public key to disk
	name := fmt.Sprintf("vault_ssh_ca_%s_%s", c.username, c.ip)
	signedPublicKeyPath, err, closer := c.writeTemporaryKey(name, []byte(key))
	defer closer()
	if err != nil {
		return errors.Wrap(err, "failed to write signed public key")
	}

	args := append([]string{
		"-i", c.privateKeyPath,
		"-i", signedPublicKeyPath,
		"-o UserKnownHostsFile=" + userKnownHostsFile,
		"-o StrictHostKeyChecking=" + strictHostKeyChecking,
		c.username + "@" + c.ip,
	}, c.sshArgs...)

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to run ssh command")
	}

	// There is no secret to revoke, since it's a certificate signing

	return nil
}

// handleTypeOTP is used to handle SSH logins using the "otp" key type.
func (c *SSHCommand) handleTypeOTP() error {
	secret, cred, err := c.generateCredential()
	if err != nil {
		return errors.Wrap(err, "failed to generate credential")
	}

	// Handle no-exec
	if c.noExec {
		// This is hacky, but OutputSecret returns an int, not an error :(
		if i := OutputSecret(c.Ui, c.format, secret); i != 0 {
			return fmt.Errorf("an error occurred outputting the secret")
		}
		return nil
	}

	var cmd *exec.Cmd

	// Check if the application 'sshpass' is installed in the client machine.
	// If it is then, use it to automate typing in OTP to the prompt. Unfortunately,
	// it was not possible to automate it without a third-party application, with
	// only the Go libraries.
	// Feel free to try and remove this dependency.
	sshpassPath, err := exec.LookPath("sshpass")
	if err != nil {
		c.Ui.Warn("" +
			"Vault could not locate sshpass. The OTP code for the session will be\n" +
			"displayed below. Enter this code in the SSH password prompt. If you\n" +
			"install sshpass, Vault can automatically perform this step for you.")
		c.Ui.Output("OTP for the session is " + cred.Key)

		args := append([]string{
			"-o UserKnownHostsFile=" + c.userKnownHostsFile,
			"-o StrictHostKeyChecking=" + c.strictHostKeyChecking,
			"-p", cred.Port,
			c.username + "@" + c.ip,
		}, c.sshArgs...)
		cmd = exec.Command("ssh", args...)
	} else {
		args := append([]string{
			"-e", // Read password for SSHPASS environment variable
			"ssh",
			"-o UserKnownHostsFile=" + c.userKnownHostsFile,
			"-o StrictHostKeyChecking=" + c.strictHostKeyChecking,
			"-p", cred.Port,
			c.username + "@" + c.ip,
		}, c.sshArgs...)
		cmd = exec.Command(sshpassPath, args...)
		env := os.Environ()
		env = append(env, fmt.Sprintf("SSHPASS=%s", string(cred.Key)))
		cmd.Env = env
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to run ssh command")
	}

	// Revoke the key if it's longer than expected
	if err := c.client.Sys().Revoke(secret.LeaseID); err != nil {
		return errors.Wrap(err, "failed to revoke key")
	}

	return nil
}

// handleTypeDynamic is used to handle SSH logins using the "dyanmic" key type.
func (c *SSHCommand) handleTypeDynamic() error {
	// Generate the credential
	secret, cred, err := c.generateCredential()
	if err != nil {
		return errors.Wrap(err, "failed to generate credential")
	}

	// Handle no-exec
	if c.noExec {
		// This is hacky, but OutputSecret returns an int, not an error :(
		if i := OutputSecret(c.Ui, c.format, secret); i != 0 {
			return fmt.Errorf("an error occurred outputting the secret")
		}
		return nil
	}

	// Write the dynamic key to disk
	name := fmt.Sprintf("vault_ssh_dynamic_%s_%s", c.username, c.ip)
	keyPath, err, closer := c.writeTemporaryKey(name, []byte(cred.Key))
	defer closer()
	if err != nil {
		return errors.Wrap(err, "failed to save dyanmic key")
	}

	args := append([]string{
		"-i", keyPath,
		"-o UserKnownHostsFile=" + c.userKnownHostsFile,
		"-o StrictHostKeyChecking=" + c.strictHostKeyChecking,
		"-p", cred.Port,
		c.username + "@" + c.ip,
	}, c.sshArgs...)

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return errors.Wrap(err, "failed to run ssh command")
	}

	// Revoke the key if it's longer than expected
	if err := c.client.Sys().Revoke(secret.LeaseID); err != nil {
		return errors.Wrap(err, "failed to revoke key")
	}

	return nil
}

// generateCredential generates a credential for the given role and returns the
// decoded secret data.
func (c *SSHCommand) generateCredential() (*api.Secret, *SSHCredentialResp, error) {
	// Attempt to generate the credential.
	secret, err := c.sshClient.Credential(c.role, map[string]interface{}{
		"username": c.username,
		"ip":       c.ip,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get credentials")
	}
	if secret == nil || secret.Data == nil {
		return nil, nil, fmt.Errorf("vault returned empty credentials")
	}

	// Port comes back as a json.Number which mapstructure doesn't like, so
	// convert it
	if d, ok := secret.Data["port"].(json.Number); ok {
		secret.Data["port"] = d.String()
	}

	// Use mapstructure to decode the response
	var resp SSHCredentialResp
	if err := mapstructure.Decode(secret.Data, &resp); err != nil {
		return nil, nil, errors.Wrap(err, "failed to decode credential")
	}

	// Check for an empty key response
	if len(resp.Key) == 0 {
		return nil, nil, fmt.Errorf("vault returned an invalid key")
	}

	return secret, &resp, nil
}

// writeTemporaryFile writes a file to a temp location with the given data and
// file permissions.
func (c *SSHCommand) writeTemporaryFile(name string, data []byte, perms os.FileMode) (string, error, func() error) {
	// default closer to prevent panic
	closer := func() error { return nil }

	f, err := ioutil.TempFile("", name)
	if err != nil {
		return "", errors.Wrap(err, "creating temporary file"), closer
	}

	closer = func() error { return os.Remove(f.Name()) }

	if err := ioutil.WriteFile(f.Name(), data, perms); err != nil {
		return "", errors.Wrap(err, "writing temporary key"), closer
	}

	return f.Name(), nil, closer
}

// writeTemporaryKey writes the key to a temporary file and returns the path.
// The caller should defer the closer to cleanup the key.
func (c *SSHCommand) writeTemporaryKey(name string, data []byte) (string, error, func() error) {
	return c.writeTemporaryFile(name, data, 0600)
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
	if secret == nil || secret.Data == nil {
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

// userAndIP takes an argument in the format foo@1.2.3.4 and separates the IP
// and user parts, returning any errors.
func (c *SSHCommand) userAndIP(s string) (string, string, error) {
	// split the parameter username@ip
	input := strings.Split(s, "@")
	var username, address string

	// If only IP is mentioned and username is skipped, assume username to
	// be the current username. Vault SSH role's default username could have
	// been used, but in order to retain the consistency with SSH command,
	// current username is employed.
	switch len(input) {
	case 1:
		u, err := user.Current()
		if err != nil {
			return "", "", errors.Wrap(err, "failed to fetch current user")
		}
		username, address = u.Username, input[0]
	case 2:
		username, address = input[0], input[1]
	default:
		return "", "", fmt.Errorf("invalid arguments: %q", s)
	}

	// Resolving domain names to IP address on the client side.
	// Vault only deals with IP addresses.
	ipAddr, err := net.ResolveIPAddr("ip", address)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to resolve IP address")
	}
	ip := ipAddr.String()

	return username, ip, nil
}

func (c *SSHCommand) Synopsis() string {
	return "Initiate an SSH session"
}

func (c *SSHCommand) Help() string {
	helpText := `
Usage: vault ssh [options] username@ip [ssh options]

  Establishes an SSH connection with the target machine.

  This command uses one of the SSH authentication backends to authenticate and
  automatically establish an SSH connection to a host. This operation requires
  that the SSH backend is mounted and configured.

  SSH using the OTP mode (requires sshpass for full automation):

      $ vault ssh -mode=otp -role=my-role user@1.2.3.4

  SSH using the CA mode:

      $ vault ssh -mode=ca -role=my-role user@1.2.3.4

  SSH using CA mode with host key verification:

      $ vault ssh \
          -mode=ca \
          -role=my-role \
          -host-key-mount-point=host-signer \
          -host-key-hostnames=example.com \
          user@example.com

  For the full list of options and arguments, please see the documentation.

General Options:
` + meta.GeneralOptionsUsage() + `
SSH Options:

  -role            Role to be used to create the key. Each IP is associated with
                   a role. To see the associated roles with IP, use "lookup"
                   endpoint. If you are certain that there is only one role
                   associated with the IP, you can skip mentioning the role. It
                   will be chosen by default. If there are no roles associated
                   with the IP, register the CIDR block of that IP using the
                   "roles/" endpoint.

  -no-exec         Shows the credentials but does not establish connection.

  -mount-point     Mount point of SSH backend. If the backend is mounted at
                   "ssh" (default), this parameter can be skipped.

  -format          If the "no-exec" option is enabled, the credentials will be
                   printed out and SSH connection will not be established. The
                   format of the output can be "json" or "table" (default).

  -strict-host-key-checking   This option corresponds to "StrictHostKeyChecking"
                   of SSH configuration. If "sshpass" is employed to enable
                   automated login, then if host key is not "known" to the
                   client, "vault ssh" command will fail. Set this option to
                   "no" to bypass the host key checking. Defaults to "ask".
                   Can also be specified with the
                   "VAULT_SSH_STRICT_HOST_KEY_CHECKING" environment variable.

  -user-known-hosts-file   This option corresponds to "UserKnownHostsFile" of
                   SSH configuration. Assigns the file to use for storing the
                   host keys. If this option is set to "/dev/null" along with
                   "-strict-host-key-checking=no", both warnings and host key
                   checking can be avoided while establishing the connection.
                   Defaults to "~/.ssh/known_hosts". Can also be specified with
                   "VAULT_SSH_USER_KNOWN_HOSTS_FILE" environment variable.

CA Mode Options:

  - public-key-path=<path>
      The path to the public key to send to Vault for signing. The default value
      is ~/.ssh/id_rsa.pub.

  - private-key-path=<path>
      The path to the private key to use for authentication. This must be the
      corresponding private key to -public-key-path. The default value is
      ~/.ssh/id_rsa.

  - host-key-mount-point=<string>
      The mount point to the SSH backend where host keys are signed. When given
      a value, Vault will generate a custom known_hosts file with delegation to
      the CA at the provided mount point and verify the SSH connection's host
      keys against the provided CA. By default, this command uses the users's
      existing known_hosts file. When this flag is set, this command will force
      strict host key checking and will override any values provided for a
      custom -user-known-hosts-file.

  - host-key-hostnames=<string>
      The list of hostnames to delegate for this certificate authority. By
      default, this is "*", which allows all domains and IPs. To restrict
      validation to a series of hostnames, specify them as comma-separated
      values here.
`
	return strings.TrimSpace(helpText)
}
