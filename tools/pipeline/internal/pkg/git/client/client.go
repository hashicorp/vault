// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/url"
	"os"
	"os/exec"
	"strings"

	slogctx "github.com/veqryn/slog-context"
)

// Client is the local git client.
type Client struct {
	Token  string
	Host   string
	config map[string]string
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
			// Disable the pager so command output is always machine-readable.
			"core.pager": "",
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

// WithHost sets the git remote host used in credential URL rewrites.
// When not set, defaults to "github.com".
func WithHost(host string) NewClientOpt {
	return func(client *Client) {
		client.Host = host
	}
}

// WithConfig sets additional gitconfig in NewClient()
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
	if c.Token != "" || len(c.config) > 0 {
		res.Env = c.configEnv()
		env = append(env, res.Env...)
	}

	cmd := exec.Command("git", append([]string{subCmd}, opts.Strings()...)...)
	cmd.Env = env
	res.Cmd = cmd.String()
	ctx = slogctx.Append(ctx, slog.String("cmd", cmd.String()))
	slog.Default().DebugContext(ctx, "executing git command")
	var err error
	res.Stdout, err = cmd.Output()
	if err != nil {
		slog.Default().ErrorContext(slogctx.Append(
			ctx,
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
	b.WriteString(string(e.Stdout))
	b.WriteString("\n")
	b.WriteString(string(e.Stderr))
	b.WriteString("\n")

	return b.String()
}

// configEnv builds the GIT_CONFIG_* environment variable slice from the
// client's config map and optional token credential rewrite.
//
// It injects config via environment variables to avoid:
//   - modifying local or global gitconfig
//   - relying on preconfigured gitconfig
//   - requiring a credstore
//   - sensitive values like tokens being passed via flags and thus potentially
//     bleeding into STDOUT
func (c *Client) configEnv() []string {
	env := make(map[string]string, len(c.config))
	maps.Copy(env, c.config)

	if c.Token != "" {
		// NOTE: This basic auth token probably only works with Github right now,
		// which is fine because our pipeline only supports Github. Other SCM repos
		// have different rules around the user in the auth portion of the URL.
		// Github doesn't care what the username is but requires one to be set so
		// we always set it to user.
		host := c.Host
		if host == "" {
			host = "github.com"
		}
		token := url.UserPassword("user", c.Token).String()
		env[fmt.Sprintf("url.https://%s@%s.insteadOf", token, host)] = "https://" + host
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

	return vars
}
