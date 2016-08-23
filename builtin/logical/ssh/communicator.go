package ssh

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	log "github.com/mgutz/logxi/v1"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type comm struct {
	client  *ssh.Client
	config  *SSHCommConfig
	conn    net.Conn
	address string
}

// SSHCommConfig is the structure used to configure the SSH communicator.
type SSHCommConfig struct {
	// The configuration of the Go SSH connection
	SSHConfig *ssh.ClientConfig

	// Connection returns a new connection. The current connection
	// in use will be closed as part of the Close method, or in the
	// case an error occurs.
	Connection func() (net.Conn, error)

	// Pty, if true, will request a pty from the remote end.
	Pty bool

	// DisableAgent, if true, will not forward the SSH agent.
	DisableAgent bool

	// Logger for output
	Logger log.Logger
}

// Creates a new communicator implementation over SSH. This takes
// an already existing TCP connection and SSH configuration.
func SSHCommNew(address string, config *SSHCommConfig) (result *comm, err error) {
	// Establish an initial connection and connect
	result = &comm{
		config:  config,
		address: address,
	}

	if err = result.reconnect(); err != nil {
		result = nil
		return
	}

	return
}

func (c *comm) Close() error {
	var err error
	if c.conn != nil {
		err = c.conn.Close()
	}
	c.conn = nil
	c.client = nil
	return err
}

func (c *comm) Upload(path string, input io.Reader, fi *os.FileInfo) error {
	// The target directory and file for talking the SCP protocol
	target_dir := filepath.Dir(path)
	target_file := filepath.Base(path)

	// On windows, filepath.Dir uses backslash separators (ie. "\tmp").
	// This does not work when the target host is unix.  Switch to forward slash
	// which works for unix and windows
	target_dir = filepath.ToSlash(target_dir)

	scpFunc := func(w io.Writer, stdoutR *bufio.Reader) error {
		return scpUploadFile(target_file, input, w, stdoutR, fi)
	}

	return c.scpSession("scp -vt "+target_dir, scpFunc)
}

func (c *comm) NewSession() (session *ssh.Session, err error) {
	if c.client == nil {
		err = errors.New("client not available")
	} else {
		session, err = c.client.NewSession()
	}

	if err != nil {
		c.config.Logger.Error("ssh session open error, attempting reconnect", "error", err)
		if err := c.reconnect(); err != nil {
			c.config.Logger.Error("reconnect attempt failed", "error", err)
			return nil, err
		}

		return c.client.NewSession()
	}

	return session, nil
}

func (c *comm) reconnect() error {
	// Close previous connection.
	if c.conn != nil {
		c.Close()
	}

	var err error
	c.conn, err = c.config.Connection()
	if err != nil {
		// Explicitly set this to the REAL nil. Connection() can return
		// a nil implementation of net.Conn which will make the
		// "if c.conn == nil" check fail above. Read here for more information
		// on this psychotic language feature:
		//
		// http://golang.org/doc/faq#nil_error
		c.conn = nil
		c.config.Logger.Error("reconnection error", "error", err)
		return err
	}

	sshConn, sshChan, req, err := ssh.NewClientConn(c.conn, c.address, c.config.SSHConfig)
	if err != nil {
		c.config.Logger.Error("handshake error", "error", err)
		c.Close()
		return err
	}
	if sshConn != nil {
		c.client = ssh.NewClient(sshConn, sshChan, req)
	}
	c.connectToAgent()

	return nil
}

func (c *comm) connectToAgent() {
	if c.client == nil {
		return
	}

	if c.config.DisableAgent {
		return
	}

	// open connection to the local agent
	socketLocation := os.Getenv("SSH_AUTH_SOCK")
	if socketLocation == "" {
		return
	}
	agentConn, err := net.Dial("unix", socketLocation)
	if err != nil {
		c.config.Logger.Error("could not connect to local agent socket", "socket_path", socketLocation)
		return
	}
	defer agentConn.Close()

	// create agent and add in auth
	forwardingAgent := agent.NewClient(agentConn)
	if forwardingAgent == nil {
		c.config.Logger.Error("could not create agent client")
		return
	}

	// add callback for forwarding agent to SSH config
	// XXX - might want to handle reconnects appending multiple callbacks
	auth := ssh.PublicKeysCallback(forwardingAgent.Signers)
	c.config.SSHConfig.Auth = append(c.config.SSHConfig.Auth, auth)
	agent.ForwardToAgent(c.client, forwardingAgent)

	// Setup a session to request agent forwarding
	session, err := c.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	err = agent.RequestAgentForwarding(session)
	if err != nil {
		c.config.Logger.Error("error requesting agent forwarding", "error", err)
		return
	}
	return
}

func (c *comm) scpSession(scpCommand string, f func(io.Writer, *bufio.Reader) error) error {
	session, err := c.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// Get a pipe to stdin so that we can send data down
	stdinW, err := session.StdinPipe()
	if err != nil {
		return err
	}

	// We only want to close once, so we nil w after we close it,
	// and only close in the defer if it hasn't been closed already.
	defer func() {
		if stdinW != nil {
			stdinW.Close()
		}
	}()

	// Get a pipe to stdout so that we can get responses back
	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	stdoutR := bufio.NewReader(stdoutPipe)

	// Set stderr to a bytes buffer
	stderr := new(bytes.Buffer)
	session.Stderr = stderr

	// Start the sink mode on the other side
	if err := session.Start(scpCommand); err != nil {
		return err
	}

	// Call our callback that executes in the context of SCP. We ignore
	// EOF errors if they occur because it usually means that SCP prematurely
	// ended on the other side.
	if err := f(stdinW, stdoutR); err != nil && err != io.EOF {
		return err
	}

	// Close the stdin, which sends an EOF, and then set w to nil so that
	// our defer func doesn't close it again since that is unsafe with
	// the Go SSH package.
	stdinW.Close()
	stdinW = nil

	// Wait for the SCP connection to close, meaning it has consumed all
	// our data and has completed. Or has errored.
	err = session.Wait()
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			// Otherwise, we have an ExitErorr, meaning we can just read
			// the exit status
			c.config.Logger.Error("got non-zero exit status", "exit_status", exitErr.ExitStatus())

			// If we exited with status 127, it means SCP isn't available.
			// Return a more descriptive error for that.
			if exitErr.ExitStatus() == 127 {
				return errors.New(
					"SCP failed to start. This usually means that SCP is not\n" +
						"properly installed on the remote system.")
			}
		}

		return err
	}
	return nil
}

// checkSCPStatus checks that a prior command sent to SCP completed
// successfully. If it did not complete successfully, an error will
// be returned.
func checkSCPStatus(r *bufio.Reader) error {
	code, err := r.ReadByte()
	if err != nil {
		return err
	}

	if code != 0 {
		// Treat any non-zero (really 1 and 2) as fatal errors
		message, _, err := r.ReadLine()
		if err != nil {
			return fmt.Errorf("Error reading error message: %s", err)
		}

		return errors.New(string(message))
	}

	return nil
}

func scpUploadFile(dst string, src io.Reader, w io.Writer, r *bufio.Reader, fi *os.FileInfo) error {
	var mode os.FileMode
	var size int64

	if fi != nil && (*fi).Mode().IsRegular() {
		mode = (*fi).Mode().Perm()
		size = (*fi).Size()
	} else {
		// Create a temporary file where we can copy the contents of the src
		// so that we can determine the length, since SCP is length-prefixed.
		tf, err := ioutil.TempFile("", "vault-ssh-upload")
		if err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}
		defer os.Remove(tf.Name())
		defer tf.Close()

		mode = 0644

		if _, err := io.Copy(tf, src); err != nil {
			return err
		}

		// Sync the file so that the contents are definitely on disk, then
		// read the length of it.
		if err := tf.Sync(); err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}

		// Seek the file to the beginning so we can re-read all of it
		if _, err := tf.Seek(0, 0); err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}

		tfi, err := tf.Stat()
		if err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}

		size = tfi.Size()
		src = tf
	}

	// Start the protocol
	perms := fmt.Sprintf("C%04o", mode)

	fmt.Fprintln(w, perms, size, dst)
	if err := checkSCPStatus(r); err != nil {
		return err
	}

	if _, err := io.CopyN(w, src, size); err != nil {
		return err
	}

	fmt.Fprint(w, "\x00")
	if err := checkSCPStatus(r); err != nil {
		return err
	}

	return nil
}
