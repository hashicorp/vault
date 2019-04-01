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
	"syscall"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/posener/complete"
)

var _ cli.Command = (*SSHCommand)(nil)
var _ cli.CommandAutocomplete = (*SSHCommand)(nil)

type SSHCommand struct {
	*BaseCommand

	// Common SSH options
	flagMode                  string
	flagRole                  string
	flagNoExec                bool
	flagMountPoint            string
	flagStrictHostKeyChecking string
	flagSSHExecutable         string
	flagUserKnownHostsFile    string

	// SSH CA Mode options
	flagPublicKeyPath     string
	flagPrivateKeyPath    string
	flagHostKeyMountPoint string
	flagHostKeyHostnames  string
	flagValidPrincipals   string
}

func (c *SSHCommand) Synopsis() string {
	return "Initiate an SSH session"
}

func (c *SSHCommand) Help() string {
	helpText := `
Usage: vault ssh [options] username@ip [ssh options]

  Establishes an SSH connection with the target machine.

  This command uses one of the SSH secrets engines to authenticate and
  automatically establish an SSH connection to a host. This operation requires
  that the SSH secrets engine is mounted and configured.

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

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *SSHCommand) Flags() *FlagSets {
	set := c.flagSet(FlagSetHTTP | FlagSetOutputField | FlagSetOutputFormat)

	f := set.NewFlagSet("SSH Options")

	// TODO: doc field?

	// General
	f.StringVar(&StringVar{
		Name:       "mode",
		Target:     &c.flagMode,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictSet("ca", "dynamic", "otp"),
		Usage:      "Name of the authentication mode (ca, dynamic, otp).",
	})

	f.StringVar(&StringVar{
		Name:       "role",
		Target:     &c.flagRole,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage:      "Name of the role to use to generate the key.",
	})

	f.BoolVar(&BoolVar{
		Name:       "no-exec",
		Target:     &c.flagNoExec,
		Default:    false,
		EnvVar:     "",
		Completion: complete.PredictNothing,
		Usage: "Print the generated credentials, but do not establish a " +
			"connection.",
	})

	f.StringVar(&StringVar{
		Name:       "mount-point",
		Target:     &c.flagMountPoint,
		Default:    "ssh/",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage:      "Mount point to the SSH secrets engine.",
	})

	f.StringVar(&StringVar{
		Name:       "strict-host-key-checking",
		Target:     &c.flagStrictHostKeyChecking,
		Default:    "ask",
		EnvVar:     "VAULT_SSH_STRICT_HOST_KEY_CHECKING",
		Completion: complete.PredictSet("ask", "no", "yes"),
		Usage: "Value to use for the SSH configuration option " +
			"\"StrictHostKeyChecking\".",
	})

	f.StringVar(&StringVar{
		Name:       "user-known-hosts-file",
		Target:     &c.flagUserKnownHostsFile,
		Default:    "",
		EnvVar:     "VAULT_SSH_USER_KNOWN_HOSTS_FILE",
		Completion: complete.PredictFiles("*"),
		Usage: "Value to use for the SSH configuration option " +
			"\"UserKnownHostsFile\".",
	})

	// SSH CA
	f = set.NewFlagSet("CA Mode Options")

	f.StringVar(&StringVar{
		Name:       "public-key-path",
		Target:     &c.flagPublicKeyPath,
		Default:    "~/.ssh/id_rsa.pub",
		EnvVar:     "",
		Completion: complete.PredictFiles("*"),
		Usage:      "Path to the SSH public key to send to Vault for signing.",
	})

	f.StringVar(&StringVar{
		Name:       "private-key-path",
		Target:     &c.flagPrivateKeyPath,
		Default:    "~/.ssh/id_rsa",
		EnvVar:     "",
		Completion: complete.PredictFiles("*"),
		Usage: "Path to the SSH private key to use for authentication. This must " +
			"be the corresponding private key to -public-key-path.",
	})

	f.StringVar(&StringVar{
		Name:       "host-key-mount-point",
		Target:     &c.flagHostKeyMountPoint,
		Default:    "",
		EnvVar:     "VAULT_SSH_HOST_KEY_MOUNT_POINT",
		Completion: complete.PredictAnything,
		Usage: "Mount point to the SSH secrets engine where host keys are signed. " +
			"When given a value, Vault will generate a custom \"known_hosts\" file " +
			"with delegation to the CA at the provided mount point to verify the " +
			"SSH connection's host keys against the provided CA. By default, host " +
			"keys are validated against the user's local \"known_hosts\" file. " +
			"This flag forces strict key host checking and ignores a custom user " +
			"known hosts file.",
	})

	f.StringVar(&StringVar{
		Name:       "host-key-hostnames",
		Target:     &c.flagHostKeyHostnames,
		Default:    "*",
		EnvVar:     "VAULT_SSH_HOST_KEY_HOSTNAMES",
		Completion: complete.PredictAnything,
		Usage: "List of hostnames to delegate for the CA. The default value " +
			"allows all domains and IPs. This is specified as a comma-separated " +
			"list of values.",
	})

	f.StringVar(&StringVar{
		Name:       "valid-principals",
		Target:     &c.flagValidPrincipals,
		Default:    "",
		EnvVar:     "",
		Completion: complete.PredictAnything,
		Usage: "List of valid principal names to include in the generated " +
			"user certificate. This is specified as a comma-separated list of values.",
	})

	f.StringVar(&StringVar{
		Name:       "ssh-executable",
		Target:     &c.flagSSHExecutable,
		Default:    "ssh",
		EnvVar:     "VAULT_SSH_EXECUTABLE",
		Completion: complete.PredictAnything,
		Usage:      "Path to the SSH executable to use when connecting to the host",
	})

	return set
}

func (c *SSHCommand) AutocompleteArgs() complete.Predictor {
	return nil
}

func (c *SSHCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

// Structure to hold the fields returned when asked for a credential from SSH
// secrets engine.
type SSHCredentialResp struct {
	KeyType  string `mapstructure:"key_type"`
	Key      string `mapstructure:"key"`
	Username string `mapstructure:"username"`
	IP       string `mapstructure:"ip"`
	Port     string `mapstructure:"port"`
}

func (c *SSHCommand) Run(args []string) int {
	f := c.Flags()

	if err := f.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Use homedir to expand any relative paths such as ~/.ssh
	c.flagUserKnownHostsFile = expandPath(c.flagUserKnownHostsFile)
	c.flagPublicKeyPath = expandPath(c.flagPublicKeyPath)
	c.flagPrivateKeyPath = expandPath(c.flagPrivateKeyPath)

	args = f.Args()
	if len(args) < 1 {
		c.UI.Error(fmt.Sprintf("Not enough arguments, (expected 1-n, got %d)", len(args)))
		return 1
	}

	// Extract the hostname, username and port from the ssh command
	hostname, username, port, err := c.parseSSHCommand(args)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing the ssh command: %q", err))
		return 1
	}

	// Use the current user if no user was specified in the ssh command
	if username == "" {
		u, err := user.Current()
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error getting the current user: %q", err))
			return 1
		}
		username = u.Username
	}

	ip, err := c.resolveHostname(hostname)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error resolving the ssh hostname: %q", err))
		return 1
	}

	// Set the client in the command
	_, err = c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Credentials are generated only against a registered role. If user
	// does not specify a role with the SSH command, then lookup API is used
	// to fetch all the roles with which this IP is associated. If there is
	// only one role associated with it, use it to establish the connection.
	//
	// TODO: remove in 0.9.0, convert to validation error
	if c.flagRole == "" {
		c.UI.Warn(wrapAtLength(
			"WARNING: No -role specified. Use -role to tell Vault which ssh role " +
				"to use for authentication. In the future, you will need to tell " +
				"Vault which role to use. For now, Vault will attempt to guess based " +
				"on the API response. This will be removed in the Vault 1.1."))

		role, err := c.defaultRole(c.flagMountPoint, ip)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error choosing role: %v", err))
			return 1
		}
		// Print the default role chosen so that user knows the role name
		// if something doesn't work. If the role chosen is not allowed to
		// be used by the user (ACL enforcement), then user should see an
		// error message accordingly.
		c.UI.Output(fmt.Sprintf("Vault SSH: Role: %q", role))
		c.flagRole = role
	}

	// If no mode was given, perform the old-school lookup. Keep this now for
	// backwards-compatibility, but print a warning.
	//
	// TODO: remove in 0.9.0, convert to validation error
	if c.flagMode == "" {
		c.UI.Warn(wrapAtLength(
			"WARNING: No -mode specified. Use -mode to tell Vault which ssh " +
				"authentication mode to use. In the future, you will need to tell " +
				"Vault which mode to use. For now, Vault will attempt to guess based " +
				"on the API response. This guess involves creating a temporary " +
				"credential, reading its type, and then revoking it. To reduce the " +
				"number of API calls and surface area, specify -mode directly. This " +
				"will be removed in Vault 1.1."))
		secret, cred, err := c.generateCredential(username, ip)
		if err != nil {
			// This is _very_ hacky, but is the only sane backwards-compatible way
			// to do this. If the error is "key type unknown", we just assume the
			// type is "ca". In the future, mode will be required as an option.
			if strings.Contains(err.Error(), "key type unknown") {
				c.flagMode = ssh.KeyTypeCA
			} else {
				c.UI.Error(fmt.Sprintf("Error getting credential: %s", err))
				return 1
			}
		} else {
			c.flagMode = cred.KeyType
		}

		// Revoke the secret, since the child functions will generate their own
		// credential. Users wishing to avoid this should specify -mode.
		if secret != nil {
			if err := c.client.Sys().Revoke(secret.LeaseID); err != nil {
				c.UI.Warn(fmt.Sprintf("Failed to revoke temporary key: %s", err))
			}
		}
	}

	switch strings.ToLower(c.flagMode) {
	case ssh.KeyTypeCA:
		return c.handleTypeCA(username, ip, port, args)
	case ssh.KeyTypeOTP:
		return c.handleTypeOTP(username, ip, port, args)
	case ssh.KeyTypeDynamic:
		return c.handleTypeDynamic(username, ip, port, args)
	default:
		c.UI.Error(fmt.Sprintf("Unknown SSH mode: %s", c.flagMode))
		return 1
	}
}

// handleTypeCA is used to handle SSH logins using the "CA" key type.
func (c *SSHCommand) handleTypeCA(username, ip, port string, sshArgs []string) int {
	// Read the key from disk
	publicKey, err := ioutil.ReadFile(c.flagPublicKeyPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to read public key %s: %s",
			c.flagPublicKeyPath, err))
		return 1
	}

	sshClient := c.client.SSHWithMountPoint(c.flagMountPoint)

	var principals = username
	if c.flagValidPrincipals != "" {
		principals = c.flagValidPrincipals
	}

	// Attempt to sign the public key
	secret, err := sshClient.SignKey(c.flagRole, map[string]interface{}{
		// WARNING: publicKey is []byte, which is b64 encoded on JSON upload. We
		// have to convert it to a string. SV lost many hours to this...
		"public_key":       string(publicKey),
		"valid_principals": principals,
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
		c.UI.Error(fmt.Sprintf("failed to sign public key %s: %s",
			c.flagPublicKeyPath, err))
		return 2
	}
	if secret == nil || secret.Data == nil {
		c.UI.Error("missing signed key")
		return 2
	}

	// Handle no-exec
	if c.flagNoExec {
		if c.flagField != "" {
			return PrintRawField(c.UI, secret, c.flagField)
		}
		return OutputSecret(c.UI, secret)
	}

	// Extract public key
	key, ok := secret.Data["signed_key"].(string)
	if !ok || key == "" {
		c.UI.Error("signed key is empty")
		return 2
	}

	// Capture the current value - this could be overwritten later if the user
	// enabled host key signing verification.
	userKnownHostsFile := c.flagUserKnownHostsFile
	strictHostKeyChecking := c.flagStrictHostKeyChecking

	// Handle host key signing verification. If the user specified a mount point,
	// download the public key, trust it with the given domains, and use that
	// instead of the user's regular known_hosts file.
	if c.flagHostKeyMountPoint != "" {
		secret, err := c.client.Logical().Read(c.flagHostKeyMountPoint + "/config/ca")
		if err != nil {
			c.UI.Error(fmt.Sprintf("failed to get host signing key: %s", err))
			return 2
		}
		if secret == nil || secret.Data == nil {
			c.UI.Error("missing host signing key")
			return 2
		}
		publicKey, ok := secret.Data["public_key"].(string)
		if !ok || publicKey == "" {
			c.UI.Error("host signing key is empty")
			return 2
		}

		// Write the known_hosts file
		name := fmt.Sprintf("vault_ssh_ca_known_hosts_%s_%s", username, ip)
		data := fmt.Sprintf("@cert-authority %s %s", c.flagHostKeyHostnames, publicKey)
		knownHosts, err, closer := c.writeTemporaryFile(name, []byte(data), 0644)
		defer closer()
		if err != nil {
			c.UI.Error(fmt.Sprintf("failed to write host public key: %s", err))
			return 1
		}

		// Update the variables
		userKnownHostsFile = knownHosts
		strictHostKeyChecking = "yes"
	}

	// Write the signed public key to disk
	name := fmt.Sprintf("vault_ssh_ca_%s_%s", username, ip)
	signedPublicKeyPath, err, closer := c.writeTemporaryKey(name, []byte(key))
	defer closer()
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to write signed public key: %s", err))
		return 2
	}

	args := []string{
		"-i", c.flagPrivateKeyPath,
		"-i", signedPublicKeyPath,
		"-o StrictHostKeyChecking=" + strictHostKeyChecking,
	}

	if userKnownHostsFile != "" {
		args = append(args,
			"-o UserKnownHostsFile="+userKnownHostsFile,
		)
	}

	// Add extra user defined ssh arguments
	args = append(args, sshArgs...)

	cmd := exec.Command(c.flagSSHExecutable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitCode := 2

		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Success() {
				return 0
			}
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			}
		}

		c.UI.Error(fmt.Sprintf("failed to run ssh command: %s", err))
		return exitCode
	}

	// There is no secret to revoke, since it's a certificate signing
	return 0
}

// handleTypeOTP is used to handle SSH logins using the "otp" key type.
func (c *SSHCommand) handleTypeOTP(username, ip, port string, sshArgs []string) int {
	secret, cred, err := c.generateCredential(username, ip)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to generate credential: %s", err))
		return 2
	}

	// Handle no-exec
	if c.flagNoExec {
		if c.flagField != "" {
			return PrintRawField(c.UI, secret, c.flagField)
		}
		return OutputSecret(c.UI, secret)
	}

	var cmd *exec.Cmd

	// Check if the application 'sshpass' is installed in the client machine. If
	// it is then, use it to automate typing in OTP to the prompt. Unfortunately,
	// it was not possible to automate it without a third-party application, with
	// only the Go libraries. Feel free to try and remove this dependency.
	args := make([]string, 0)
	env := os.Environ()
	sshCmd := c.flagSSHExecutable

	sshpassPath, err := exec.LookPath("sshpass")
	if err != nil {
		// No sshpass available so using normal ssh client
		c.UI.Warn(wrapAtLength(
			"Vault could not locate \"sshpass\". The OTP code for the session is " +
				"displayed below. Enter this code in the SSH password prompt. If you " +
				"install sshpass, Vault can automatically perform this step for you."))
		c.UI.Output("OTP for the session is: " + cred.Key)
	} else {
		// sshpass is available so lets use it instead
		sshCmd = sshpassPath
		args = append(args,
			"-e", // Read password for SSHPASS environment variable
			c.flagSSHExecutable,
		)
		env = append(env, fmt.Sprintf("SSHPASS=%s", string(cred.Key)))
	}

	// Only harcode the knownhostsfile path if it has been set
	if c.flagUserKnownHostsFile != "" {
		args = append(args,
			"-o UserKnownHostsFile="+c.flagUserKnownHostsFile,
		)
	}

	// If a port wasn't specified in the ssh arguments lets use the port we got back from vault
	if port == "" {
		args = append(args, "-p", cred.Port)
	}

	args = append(args,
		"-o StrictHostKeyChecking="+c.flagStrictHostKeyChecking,
	)

	// Add the rest of the ssh args appended by the user
	args = append(args, sshArgs...)

	cmd = exec.Command(sshCmd, args...)
	cmd.Env = env

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitCode := 2

		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Success() {
				return 0
			}
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			}
		}

		c.UI.Error(fmt.Sprintf("failed to run ssh command: %s", err))
		return exitCode
	}

	// Revoke the key if it's longer than expected
	if err := c.client.Sys().Revoke(secret.LeaseID); err != nil {
		c.UI.Error(fmt.Sprintf("failed to revoke key: %s", err))
		return 2
	}

	return 0
}

// handleTypeDynamic is used to handle SSH logins using the "dyanmic" key type.
func (c *SSHCommand) handleTypeDynamic(username, ip, port string, sshArgs []string) int {
	// Generate the credential
	secret, cred, err := c.generateCredential(username, ip)
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to generate credential: %s", err))
		return 2
	}

	// Handle no-exec
	if c.flagNoExec {
		if c.flagField != "" {
			return PrintRawField(c.UI, secret, c.flagField)
		}
		return OutputSecret(c.UI, secret)
	}

	// Write the dynamic key to disk
	name := fmt.Sprintf("vault_ssh_dynamic_%s_%s", username, ip)
	keyPath, err, closer := c.writeTemporaryKey(name, []byte(cred.Key))
	defer closer()
	if err != nil {
		c.UI.Error(fmt.Sprintf("failed to write dynamic key: %s", err))
		return 1
	}

	args := make([]string, 0)
	// If a port wasn't specified in the ssh arguments lets use the port we got back from vault
	if port == "" {
		args = append(args, "-p", cred.Port)
	}

	args = append(args,
		"-i", keyPath,
		"-o UserKnownHostsFile="+c.flagUserKnownHostsFile,
		"-o StrictHostKeyChecking="+c.flagStrictHostKeyChecking,
	)

	// Add extra user defined ssh arguments
	args = append(args, sshArgs...)

	cmd := exec.Command(c.flagSSHExecutable, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		exitCode := 2

		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.Success() {
				return 0
			}
			if ws, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = ws.ExitStatus()
			}
		}

		c.UI.Error(fmt.Sprintf("failed to run ssh command: %s", err))
		return exitCode
	}

	// Revoke the key if it's longer than expected
	if err := c.client.Sys().Revoke(secret.LeaseID); err != nil {
		c.UI.Error(fmt.Sprintf("failed to revoke key: %s", err))
		return 2
	}

	return 0
}

// generateCredential generates a credential for the given role and returns the
// decoded secret data.
func (c *SSHCommand) generateCredential(username, ip string) (*api.Secret, *SSHCredentialResp, error) {
	sshClient := c.client.SSHWithMountPoint(c.flagMountPoint)

	// Attempt to generate the credential.
	secret, err := sshClient.Credential(c.flagRole, map[string]interface{}{
		"username": username,
		"ip":       ip,
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

	if err := f.Close(); err != nil {
		return "", errors.Wrap(err, "closing temporary key"), closer
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
	secret, err := c.client.Logical().Write(mountPoint+"/lookup", data)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("error finding roles for IP %q: {{err}}", ip), err)
	}
	if secret == nil || secret.Data == nil {
		return "", errwrap.Wrapf(fmt.Sprintf("error finding roles for IP %q: {{err}}", ip), err)
	}

	if secret.Data["roles"] == nil {
		return "", fmt.Errorf("no matching roles found for IP %q", ip)
	}

	if len(secret.Data["roles"].([]interface{})) == 1 {
		return secret.Data["roles"].([]interface{})[0].(string), nil
	} else {
		var roleNames string
		for _, item := range secret.Data["roles"].([]interface{}) {
			roleNames += item.(string) + ", "
		}
		roleNames = strings.TrimRight(roleNames, ", ")
		return "", fmt.Errorf("Roles: %q. "+`
			Multiple roles are registered for this IP.
			Select a role using '-role' option.
			Note that all roles may not be permitted, based on ACLs.`, roleNames)
	}
}

func (c *SSHCommand) isSingleSSHArg(arg string) bool {
	// list of single SSH arguments is taken from
	// https://github.com/openssh/openssh-portable/blob/28013759f09ed3ebf7e8335e83a62936bd7a7f47/ssh.c#L204
	singleArgs := []string{
		"4", "6", "A", "a", "C", "f", "G", "g", "K", "k", "M", "N", "n", "q",
		"s", "T", "t", "V", "v", "X", "x", "Y", "y",
	}

	// We want to get the first character after the dash. This is so args like -vvv are picked up as just being -v
	flag := string(arg[1])

	for _, a := range singleArgs {
		if flag == a {
			return true
		}
	}
	return false
}

// Finds the hostname, username (optional) and port (optional) from any valid ssh command
// Supports usrname@hostname but also specifying valid ssh flags like -o User=username,
// -o Port=2222 and -p 2222 anywhere in the command
func (c *SSHCommand) parseSSHCommand(args []string) (hostname string, username string, port string, err error) {
	lastArg := ""
	for _, i := range args {
		arg := lastArg
		lastArg = ""

		// If -p has been specified then this is our ssh port
		if arg == "-p" {
			port = i
			continue
		}

		// this is an ssh option, lets see if User or Port have been set and use it
		if arg == "-o" {
			split := strings.Split(i, "=")
			key := split[0]
			// Incase the value contains = signs we want to get all of them
			value := strings.Join(split[1:], " ")

			if key == "User" {
				// Don't overwrite the user if it is already set by username@hostname
				// This matches the behaviour for how regular ssh reponds when both are specified
				if username == "" {
					username = value
				}
			}

			if key == "Port" {
				// Don't overwrite the port if it is already set by -p
				// This matches the behaviour for how regular ssh reponds when both are specified
				if port == "" {
					port = value
				}
			}
			continue
		}

		// This isn't an ssh argument that we care about. Lets keep on parsing the command
		if arg != "" {
			continue
		}

		// If this is an ssh argument with a value we want to look at it in the next loop
		if strings.HasPrefix(i, "-") {
			// If this isn't a single SSH arg we want to store the flag to we can look at the value next loop
			if !c.isSingleSSHArg(i) {
				lastArg = i
			}
			continue
		}

		// If we have gotten this far it means this is a bare argument
		// The first bare argument is the hostname
		// The second bare argument is the command to run on the remote host

		// If the hostname hasn't been set yet than it means we have found the first bare argument
		if hostname == "" {
			if strings.Contains(i, "@") {
				split := strings.Split(i, "@")
				username = split[0]
				hostname = split[1]
			} else {
				hostname = i
			}
			continue
		} else {
			// The second bare argument is the command to run on the remote host.
			// We need to break out and stop parsing arguments now
			break
		}

	}
	if hostname == "" {
		return "", "", "", errors.Wrap(
			err,
			fmt.Sprintf("failed to find a hostname in ssh command %q", strings.Join(args, " ")),
		)
	}
	return hostname, username, port, nil
}

func (c *SSHCommand) resolveHostname(hostname string) (ip string, err error) {
	// Resolving domain names to IP address on the client side.
	// Vault only deals with IP addresses.
	ipAddr, err := net.ResolveIPAddr("ip", hostname)
	if err != nil {
		return "", errors.Wrap(err, "failed to resolve IP address")
	}
	ip = ipAddr.String()
	return ip, nil
}
