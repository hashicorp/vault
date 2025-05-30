// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package git

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/url"
	"os"
	"os/exec"
	oexec "os/exec"
	"strings"
	"sync"

	slogctx "github.com/veqryn/slog-context"
)

// Client is the local git client.
type Client struct {
	Token   string
	envOnce sync.Once
	envVal  []string
	config  map[string]string
}

// OptStringer is an interface that all sub-command configuration options must
// implement.
type OptStringer interface {
	String() string
	Strings() []string
}

// ExecResponse is the response from the client running a sub-command with Exec()
type ExecResponse struct {
	Cmd    string
	Env    []string
	Stdout []byte
	Stderr []byte
}

// NewClientOpt is a NewClient() functional option
type NewClientOpt func(*Client)

// NewClient takes variable options and returns a default Client.
func NewClient(opts ...NewClientOpt) *Client {
	client := &Client{
		config: map[string]string{
			"core.pager": "",
			"user.name":  "hc-github-team-secure-vault-core",
			"user.email": "github-team-secure-vault-core@hashicorp.com",
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// WithToken sets the Token in NewClient()
func WithToken(token string) NewClientOpt {
	return func(client *Client) {
		client.Token = token
	}
}

// WithToken sets additional gitconfig in NewClient()
func WithConfig(config map[string]string) NewClientOpt {
	return func(client *Client) {
		maps.Copy(client.config, config)
	}
}

// WithLoadTokenFromEnv sets the Token from known env vars in NewClient()
func WithLoadTokenFromEnv() NewClientOpt {
	return func(client *Client) {
		if token, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
			client.Token = token
			return
		}
		if token, ok := os.LookupEnv("GH_TOKEN"); ok {
			client.Token = token
			return
		}
	}
}

// Exec executes a git sub-command.
func (c *Client) Exec(ctx context.Context, subCmd string, opts OptStringer) (*ExecResponse, error) {
	env := os.Environ()
	res := &ExecResponse{Env: os.Environ()}
	if c.Token != "" {
		res.Env = c.configEnv()
		env = append(env, res.Env...)
	}

	cmd := oexec.Command("git", append([]string{subCmd}, opts.Strings()...)...)
	cmd.Env = env
	res.Cmd = cmd.String()
	ctx = slogctx.Append(ctx, slog.String("cmd", cmd.String()))
	slog.Default().DebugContext(ctx, "executing git command")
	var err error
	res.Stdout, err = cmd.Output()
	if err != nil {
		slog.Default().ErrorContext(slogctx.Append(ctx,
			slog.String("error", err.Error()),
		), "executing git command failed")
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			res.Stderr = exitErr.Stderr
		}
	}

	return res, err
}

// String returns the ExecResponse command and output as a string
func (e *ExecResponse) String() string {
	if e == nil {
		return ""
	}

	b := strings.Builder{}
	b.WriteString(e.Cmd)
	b.WriteString("\n")
	for _, line := range strings.Split(string(e.Stdout), "\n") {
		b.WriteString(line)
	}
	b.WriteString("\n")
	for _, line := range strings.Split(string(e.Stderr), "\n") {
		b.WriteString(line)
	}
	b.WriteString("\n")

	return b.String()
}

// configEnv creates a slice of all git configuration as environment variables
// to avoid:
//   - modifying local or global gitconfig
//   - relying on preconfigured gitconfig
//   - requiring a credstore
//   - sensitive values like tokens being passed via flags and thus potentially
//     bleeding into STDOUT
//
// As this is relatively expensive it's only done once and cached so subsequent
// requests can reuse the same configuration.
func (c *Client) configEnv() []string {
	c.envOnce.Do(func() {
		env := c.config

		if c.Token != "" {
			// NOTE: This basic auth token probably only works with Github right now,
			// which is fine because our pipeline only supports Github. Other SCM repos
			// have different rules around the user in the auth portion of the URL.
			// Github doesn't care what the username is but requires one to be set so
			// we always set it to user.
			token := url.UserPassword("user", c.Token).String()
			env[fmt.Sprintf("url.https://%s@github.com.insteadOf", token)] = "https://github.com"
		}

		vars := []string{fmt.Sprintf("GIT_CONFIG_COUNT=%d", len(env))}
		count := 0
		for k, v := range env {
			vars = append(
				vars,
				fmt.Sprintf("GIT_CONFIG_KEY_%d=%s", count, k),
				fmt.Sprintf("GIT_CONFIG_VALUE_%d=%s", count, v),
			)
			count++
		}

		c.envVal = vars
	})

	return c.envVal
}
