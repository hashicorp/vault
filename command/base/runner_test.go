// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// TestRunCustomWithCommandsFn_FactoryError checks that help and output flags don't
// hide command setup errors or change how we handle successful and unknown commands.
func TestRunCustomWithCommandsFn_FactoryError(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		args []string
		code int
	}{
		{"plain", []string{"unavailable"}, 1},
		{"help", []string{"unavailable", "--help"}, 1},
		{"curl", []string{"unavailable", "-output-curl-string"}, 1},
		{"policy", []string{"unavailable", "-output-policy"}, 1},
		{"success", []string{"success"}, 0},
		{"unknown", []string{"unknown"}, 127},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var stdout, stderr bytes.Buffer
			factoryErr := errors.New("factory rejected command")
			success := &cli.MockCommand{HelpText: "command help", SynopsisText: "success"}
			initFn := func(cli.Ui, cli.Ui, *RunOptions) map[string]cli.CommandFactory {
				return map[string]cli.CommandFactory{
					"unavailable": func() (cli.Command, error) { return nil, factoryErr },
					"success":     func() (cli.Command, error) { return success, nil },
				}
			}
			code := RunCustomWithCommandsFn(tc.args, &RunOptions{
				Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr, BinaryName: "vault",
			}, initFn, "unavailable")
			require.Equal(t, tc.code, code)
			switch tc.name {
			case "success":
				require.True(t, success.RunCalled)
				require.Empty(t, stdout.String()+stderr.String())
			case "unknown":
				require.False(t, success.RunCalled)
				require.Contains(t, stderr.String(), "Usage:")
				require.NotContains(t, stderr.String(), factoryErr.Error())
			default:
				require.False(t, success.RunCalled)
				require.Empty(t, stdout.String())
				require.Equal(t, fmt.Sprintf("Error executing CLI: %s\n", factoryErr), stderr.String())
			}
		})
	}
}

// TestGenerateCurlString_Freshness checks that we don't hide command failures
// with curl output from a previous request. New requests should still produce output.
func TestGenerateCurlString_Freshness(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		marker     string
		exitCode   int
		wantCode   int
		wantOutput bool
		wantError  string
	}{
		{"nil", "nil", 1, 1, false, "refused\nUnable to generate cURL string from command\n"},
		{"cleared", "cleared", 1, 1, false, "refused\nUnable to generate cURL string from command\n"},
		{"stale", "stale", 1, 1, false, "refused\nUnable to generate cURL string from command\n"},
		{"stale_nonzero", "stale", 2, 2, false, "refused\nUnable to generate cURL string from command\n"},
		{"fresh_nonzero", "fresh", 2, 0, true, ""},
		{"first_nonzero", "first", 2, 0, true, ""},
		{"success", "fresh", 0, 1, false, "Could not generate cURL command"},
		{"unknown_nil", "nil", 127, 127, false, ""},
		{"unknown_stale", "stale", 127, 127, false, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			request, err := retryablehttp.NewRequest(http.MethodGet, "https://example.test/v1/example", nil)
			require.NoError(t, err)
			var previous, current *api.OutputStringError
			if tc.marker != "nil" && tc.marker != "first" {
				previous = &api.OutputStringError{Request: request}
			}
			switch tc.marker {
			case "stale":
				current = previous
			case "fresh", "first":
				current = &api.OutputStringError{Request: request}
			}
			var stdout, stderr bytes.Buffer
			code := generateCurlString(tc.exitCode, &RunOptions{Stdout: &stdout, Stderr: &stderr},
				bytes.NewBufferString("refused\n"), previous, current)
			require.Equal(t, tc.wantCode, code)
			require.Equal(t, tc.wantError, stderr.String())
			if tc.wantOutput {
				require.Equal(t, "curl https://example.test/v1/example\n", stdout.String())
			} else {
				require.Empty(t, stdout.String())
			}
		})
	}
}

// TestGeneratePolicy_Freshness checks that we don't hide command failures by reusing
// an old policy. We should still generate a policy when this command makes a request.
func TestGeneratePolicy_Freshness(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name       string
		marker     string
		exitCode   int
		wantCode   int
		wantOutput bool
		wantError  string
	}{
		{"nil", "nil", 1, 1, false, "refused\nUnable to generate policy from command\n"},
		{"cleared", "cleared", 1, 1, false, "refused\nUnable to generate policy from command\n"},
		{"stale", "stale", 1, 1, false, "refused\nUnable to generate policy from command\n"},
		{"stale_nonzero", "stale", 2, 2, false, "refused\nUnable to generate policy from command\n"},
		{"fresh_nonzero", "fresh", 2, 0, true, ""},
		{"first_nonzero", "first", 2, 0, true, ""},
		{"success", "fresh", 0, 1, false, "Could not generate policy"},
		{"unknown_nil", "nil", 127, 127, false, ""},
		{"unknown_stale", "stale", 127, 127, false, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var previous, current *api.OutputPolicyError
			if tc.marker != "nil" && tc.marker != "first" {
				previous = &api.OutputPolicyError{}
			}
			switch tc.marker {
			case "stale":
				current = previous
			case "fresh", "first":
				current = &api.OutputPolicyError{}
			}
			var stdout, stderr bytes.Buffer
			code := generatePolicy(tc.exitCode, &RunOptions{Stdout: &stdout, Stderr: &stderr},
				bytes.NewBufferString("refused\n"), previous, current)
			require.Equal(t, tc.wantCode, code)
			require.Equal(t, tc.wantError, stderr.String())
			if tc.wantOutput {
				require.Equal(t, "path \"\" {\n  capabilities = [\"read\"]\n}\n", stdout.String())
			} else {
				require.Empty(t, stdout.String())
			}
		})
	}
}
