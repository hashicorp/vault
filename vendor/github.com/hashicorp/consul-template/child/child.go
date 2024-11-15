// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package child

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/consul-template/signals"
)

var (
	// ErrMissingCommand is the error returned when no command is specified
	// to run.
	ErrMissingCommand = errors.New("missing command")

	// ExitCodeOK is the default OK exit code.
	ExitCodeOK = 0

	// ExitCodeError is the default error code returned when the child exits with
	// an error without a more specific code.
	ExitCodeError = 127
)

// Child is a wrapper around a child process which can be used to send signals
// and manage the processes' lifecycle.
type Child struct {
	sync.RWMutex

	stdin          io.Reader
	stdout, stderr io.Writer
	command        string
	args           []string
	env            []string

	timeout time.Duration

	reloadSignal os.Signal

	killSignal  os.Signal
	killTimeout time.Duration

	splay time.Duration

	// cmd is the actual child process under management.
	cmd *exec.Cmd

	// exitCh is the channel where the processes exit will be returned.
	exitCh chan int

	// stopLock is the mutex to lock when stopping. stopCh is the circuit breaker
	// to force-terminate any waiting splays to kill the process now. stopped is
	// a boolean that tells us if we have previously been stopped.
	stopLock sync.RWMutex
	stopCh   chan struct{}
	stopped  bool

	// whether to set process group id or not (default on)
	setpgid bool

	// whether to set process session id or not (default off)
	setsid bool

	// a logger that can be used for messages pertinent to this child process
	logger *log.Logger
}

// NewInput is input to the NewChild function.
type NewInput struct {
	// Stdin is the io.Reader where input will come from. This is sent directly to
	// the child process. Stdout and Stderr represent the io.Writer objects where
	// the child process will send output and errorput.
	Stdin          io.Reader
	Stdout, Stderr io.Writer

	// Command is the name of the command to execute. Args are the list of
	// arguments to pass when starting the command.
	Command string
	Args    []string

	// Timeout is the maximum amount of time to allow the command to execute. If
	// set to 0, the command is permitted to run infinitely.
	Timeout time.Duration

	// Env represents the condition of the child processes' environment
	// variables. Only these environment variables will be given to the child, so
	// it is the responsibility of the caller to include the parent processes
	// environment, if required. This should be in the key=value format.
	Env []string

	// ReloadSignal is the signal to send to reload this process. This value may
	// be nil.
	ReloadSignal os.Signal

	// KillSignal is the signal to send to gracefully kill this process. This
	// value may be nil.
	KillSignal os.Signal

	// KillTimeout is the amount of time to wait for the process to gracefully
	// terminate before force-killing.
	KillTimeout time.Duration

	// Splay is the maximum random amount of time to wait before sending signals.
	// This option helps reduce the thundering herd problem by effectively
	// sleeping for a random amount of time before sending the signal. This
	// prevents multiple processes from all signaling at the same time. This value
	// may be zero (which disables the splay entirely).
	Splay time.Duration

	// Setsid flag, if set to true will create the child processes with their own
	// session ID and Process group ID. Note this overrides Setpgid below.
	Setsid bool

	// Setpgid flag, if set to true will create the child processes with their
	// own Process group ID. If set the child processes to have their own PGID
	// but they will have the same SID as that of their parent
	Setpgid bool

	// an optional logger that can be used for messages pertinent to the child process
	Logger *log.Logger
}

// New creates a new child process for management with high-level APIs for
// sending signals to the child process, restarting the child process, and
// gracefully terminating the child process.
func New(i *NewInput) (*Child, error) {
	if i == nil {
		i = new(NewInput)
	}

	if len(i.Command) == 0 {
		return nil, ErrMissingCommand
	}

	if i.Logger == nil {
		i.Logger = log.Default()
	}

	child := &Child{
		stdin:        i.Stdin,
		stdout:       i.Stdout,
		stderr:       i.Stderr,
		command:      i.Command,
		args:         i.Args,
		env:          i.Env,
		timeout:      i.Timeout,
		reloadSignal: i.ReloadSignal,
		killSignal:   i.KillSignal,
		killTimeout:  i.KillTimeout,
		splay:        i.Splay,
		stopCh:       make(chan struct{}, 1),
		setpgid:      i.Setpgid && !i.Setsid,
		setsid:       i.Setsid,
		logger:       i.Logger,
	}

	return child, nil
}

// ExitCh returns the current exit channel for this child process. This channel
// may change if the process is restarted, so implementers must not cache this
// value. To prevent a potential race condition, this function must be invoked
// strictly after child.Start().
func (c *Child) ExitCh() <-chan int {
	c.RLock()
	defer c.RUnlock()
	return c.exitCh
}

// Pid returns the pid of the child process. If no child process exists, 0 is
// returned.
func (c *Child) Pid() int {
	c.RLock()
	defer c.RUnlock()
	return c.pid()
}

// Command returns the human-formatted command with arguments.
func (c *Child) Command() string {
	list := append([]string{c.command}, c.args...)
	return strings.Join(list, " ")
}

// Start starts and begins execution of the child process. A buffered channel
// is returned which is where the command's exit code will be returned upon
// exit. Any errors that occur prior to starting the command will be returned
// as the second error argument, but any errors returned by the command after
// execution will be returned as a non-zero value over the exit code channel.
func (c *Child) Start() error {
	c.logger.Printf("[INFO] (child) spawning: %s", c.Command())
	c.Lock()
	defer c.Unlock()
	return c.start()
}

// Signal sends the signal to the child process, returning any errors that
// occur.
func (c *Child) Signal(s os.Signal) error {
	c.logger.Printf("[INFO] (child) receiving signal %q", s.String())
	c.RLock()
	defer c.RUnlock()
	switch s {
	case c.reloadSignal:
		return c.reload()
	case c.killSignal:
		c.kill(true)
		return nil
	default:
		return c.signal(s)
	}
}

// Reload sends the reload signal to the child process and does not wait for a
// response. If no reload signal was provided, the process is restarted and
// replaces the process attached to this Child.
func (c *Child) Reload() error {
	if c.reloadSignal == nil {
		c.logger.Printf("[INFO] (child) restarting process")

		// Take a full lock because start is going to replace the process. We also
		// want to make sure that no other routines attempt to send reload signals
		// during this transition.
		c.Lock()
		defer c.Unlock()

		c.kill(false)
		return c.start()
	}
	c.logger.Printf("[INFO] (child) reloading process")

	// We only need a read lock here because neither the process nor the exit
	// channel are changing.
	c.RLock()
	defer c.RUnlock()

	return c.reload()
}

// Kill sends the kill signal to the child process and waits for successful
// termination. If no kill signal is defined, the process is killed with the
// most aggressive kill signal. If the process does not gracefully stop within
// the provided KillTimeout, the process is force-killed. If a splay was
// provided, this function will sleep for a random period of time between 0 and
// the provided splay value to reduce the thundering herd problem. This function
// does not return any errors because it guarantees the process will be dead by
// the return of the function call.
func (c *Child) Kill() {
	c.logger.Printf("[INFO] (child) killing process")
	c.Lock()
	defer c.Unlock()
	c.kill(false)
}

// Stop behaves almost identical to Kill except it suppresses future processes
// from being started by this child and it prevents the killing of the child
// process from sending its value back up the exit channel. This is useful
// when doing a graceful shutdown of an application.
func (c *Child) Stop() {
	c.internalStop(false)
}

// StopImmediately behaves almost identical to Stop except it does not wait
// for any random splay if configured. This is used for performing a fast
// shutdown of consul-template and its children when a kill signal is received.
func (c *Child) StopImmediately() {
	c.internalStop(true)
}

func (c *Child) internalStop(immediately bool) {
	c.logger.Printf("[INFO] (child) stopping process")

	c.Lock()
	defer c.Unlock()

	c.stopLock.Lock()
	defer c.stopLock.Unlock()
	if c.stopped {
		c.logger.Printf("[WARN] (child) already stopped")
		return
	}
	c.kill(immediately)
	close(c.stopCh)
	c.stopped = true
}

func (c *Child) start() error {
	cmd := exec.Command(c.command, c.args...)
	cmd.Stdin = c.stdin
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	cmd.Env = c.env
	setSysProcAttr(cmd, c.setpgid, c.setsid)
	if err := cmd.Start(); err != nil {
		return err
	}
	c.cmd = cmd

	// Create a new exitCh so that previously invoked commands (if any) don't
	// cause us to exit, and start a goroutine to wait for that process to end.
	exitCh := make(chan int, 1)
	go func() {
		var code int
		err := cmd.Wait()
		if err == nil {
			code = ExitCodeOK
		} else {
			code = ExitCodeError
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					code = status.ExitStatus()
				}
			}
		}

		// If the child is in the process of killing, do not send a response back
		// down the exit channel.
		c.stopLock.RLock()
		defer c.stopLock.RUnlock()
		if !c.stopped {
			select {
			case <-c.stopCh:
			case exitCh <- code:
			}
		}

		close(exitCh)
	}()

	c.exitCh = exitCh

	// If a timeout was given, start the timer to wait for the child to exit
	if c.timeout != 0 {
		select {
		case code := <-exitCh:
			if code != 0 {
				return fmt.Errorf(
					"command exited with a non-zero exit status:\n"+
						"\n    %s\n\n"+
						"Please ensure the command exits with a zero (0) "+
						"exit status to indicate success",
					c.Command(),
				)
			}
		case <-time.After(c.timeout):
			// Force-kill the process
			c.stopLock.Lock()
			defer c.stopLock.Unlock()
			if c.cmd != nil && c.cmd.Process != nil {
				c.signal(os.Kill)
			}

			return fmt.Errorf(
				"command did not exit within time limit (%q):\n    %s",
				c.timeout, c.Command())
		}
	}

	return nil
}

func (c *Child) pid() int {
	if !c.running() {
		return 0
	}
	return c.cmd.Process.Pid
}

func (c *Child) signal(s os.Signal) error {
	if !c.running() {
		return nil
	}

	sig, ok := s.(syscall.Signal)
	switch {
	case !ok:
		return fmt.Errorf("bad signal: %s", s)
	case sig == signals.SIGNULL:
		// skip on SIGNULL (ie. no signal)
		return nil
	}

	pid := c.cmd.Process.Pid
	if c.setpgid {
		// kill takes negative pid to indicate that you want to use gpid
		pid = -(pid)
	}
	// cross platform way to signal process/process group
	if p, err := os.FindProcess(pid); err != nil {
		return err
	} else {
		return p.Signal(sig)
	}
}

func (c *Child) reload() error {
	select {
	case <-c.stopCh:
	case <-c.randomSplay():
	}

	return c.signal(c.reloadSignal)
}

// kill sends the signal to kill the process using the configured signal
// if set, else the default system signal
func (c *Child) kill(immediately bool) {
	if !c.running() {
		c.logger.Printf("[DEBUG] (child) Kill() called but process dead; not waiting for splay.")
		return
	} else if immediately {
		c.logger.Printf("[DEBUG] (child) Kill() called but performing immediate shutdown; not waiting for splay.")
	} else {
		select {
		case <-c.stopCh:
		case <-c.randomSplay():
		}
	}

	var exited bool
	defer func() {
		if !exited {
			c.signal(os.Kill)
		}
		c.cmd = nil
	}()

	if c.killSignal == nil {
		return
	}

	if err := c.signal(c.killSignal); err != nil {
		c.logger.Printf("[ERR] (child) Kill failed: %s", err)
		if processNotFoundErr(err) {
			exited = true // checked in defer
		}
		return
	}

	killCh := make(chan struct{}, 1)
	go func() {
		defer close(killCh)
		c.cmd.Process.Wait()
	}()

	select {
	case <-c.stopCh:
	case <-killCh:
		exited = true
	case <-time.After(c.killTimeout):
	}
}

func (c *Child) running() bool {
	select {
	case <-c.exitCh:
		return false
	default:
	}
	return c.cmd != nil && c.cmd.Process != nil
}

func (c *Child) randomSplay() <-chan time.Time {
	if c.splay == 0 {
		return time.After(0)
	}

	ns := c.splay.Nanoseconds()
	offset := rand.Int63n(ns)
	t := time.Duration(offset)

	c.logger.Printf("[DEBUG] (child) waiting %.2fs for random splay", t.Seconds())

	return time.After(t)
}
