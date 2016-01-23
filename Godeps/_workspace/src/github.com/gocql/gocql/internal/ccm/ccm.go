// +build ccm

package ccm

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func execCmd(args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command("ccm", args...)
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = &bytes.Buffer{}
	if err := cmd.Run(); err != nil {
		return nil, errors.New(cmd.Stderr.(*bytes.Buffer).String())
	}

	return stdout, nil
}

func AllUp() error {
	status, err := Status()
	if err != nil {
		return err
	}

	for _, host := range status {
		if !host.State.IsUp() {
			if err := NodeUp(host.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

func NodeUp(node string) error {
	_, err := execCmd(node, "start", "--wait-for-binary-proto", "--wait-other-notice")
	return err
}

func NodeDown(node string) error {
	_, err := execCmd(node, "stop")
	return err
}

type Host struct {
	State NodeState
	Addr  string
	Name  string
}

type NodeState int

func (n NodeState) String() string {
	if n == NodeStateUp {
		return "UP"
	} else if n == NodeStateDown {
		return "DOWN"
	} else {
		return fmt.Sprintf("UNKNOWN_STATE_%d", n)
	}
}

func (n NodeState) IsUp() bool {
	return n == NodeStateUp
}

const (
	NodeStateUp NodeState = iota
	NodeStateDown
)

func Status() (map[string]Host, error) {
	// TODO: parse into struct o maniuplate
	out, err := execCmd("status", "-v")
	if err != nil {
		return nil, err
	}

	const (
		stateCluster = iota
		stateCommas
		stateNode
		stateOption
	)

	nodes := make(map[string]Host)
	// didnt really want to write a full state machine parser
	state := stateCluster
	sc := bufio.NewScanner(out)

	var host Host

	for sc.Scan() {
		switch state {
		case stateCluster:
			text := sc.Text()
			if !strings.HasPrefix(text, "Cluster:") {
				return nil, fmt.Errorf("expected 'Cluster:' got %q", text)
			}
			state = stateCommas
		case stateCommas:
			text := sc.Text()
			if !strings.HasPrefix(text, "-") {
				return nil, fmt.Errorf("expected commas got %q", text)
			}
			state = stateNode
		case stateNode:
			// assume nodes start with node
			text := sc.Text()
			if !strings.HasPrefix(text, "node") {
				return nil, fmt.Errorf("expected 'node' got %q", text)
			}
			line := strings.Split(text, ":")
			host.Name = line[0]

			nodeState := strings.TrimSpace(line[1])
			switch nodeState {
			case "UP":
				host.State = NodeStateUp
			case "DOWN":
				host.State = NodeStateDown
			default:
				return nil, fmt.Errorf("unknown node state from ccm: %q", nodeState)
			}

			state = stateOption
		case stateOption:
			text := sc.Text()
			if text == "" {
				state = stateNode
				nodes[host.Name] = host
				host = Host{}
				continue
			}

			line := strings.Split(strings.TrimSpace(text), "=")
			k, v := line[0], line[1]
			if k == "binary" {
				// could check errors
				// ('127.0.0.1', 9042)
				v = v[2:] // (''
				if i := strings.IndexByte(v, '\''); i < 0 {
					return nil, fmt.Errorf("invalid binary v=%q", v)
				} else {
					host.Addr = v[:i]
					// dont need port
				}
			}
		default:
			return nil, fmt.Errorf("unexpected state: %q", state)
		}
	}

	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("unable to parse ccm status: %v", err)
	}

	return nodes, nil
}
