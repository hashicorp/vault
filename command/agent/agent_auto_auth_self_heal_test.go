// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	agentConfig "github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/template"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	tokenfile "github.com/hashicorp/vault/command/agentproxyshared/auth/token-file"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	"github.com/stretchr/testify/require"
)

// TestAutoAuthSelfHealing_TokenFileAuth_SinkOutput tests that
// if the token is revoked, Auto Auth is re-triggered and a valid new token
// is written to a sink, and the template is correctly rendered with the new token
func TestAutoAuthSelfHealing_TokenFileAuth_SinkOutput(t *testing.T) {
	// Unset the environment variable so that agent picks up the right test cluster address
	t.Setenv(api.EnvVaultAddress, "")

	cluster := minimal.NewTestSoloCluster(t, nil)
	logger := corehelpers.NewTestLogger(t)
	serverClient := cluster.Cores[0].Client

	// Create token
	secret, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	require.NotNil(t, secret)
	require.NotNil(t, secret.Auth)
	require.NotEmpty(t, secret.Auth.ClientToken)
	token := secret.Auth.ClientToken

	// Write token to the auto-auth token file
	pathVaultToken := makeTempFile(t, "token-file", token)

	// Give us some leeway of 3 errors 1 from each of: auth handler, sink server template server.
	errCh := make(chan error, 3)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// Create auth handler
	am, err := tokenfile.NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path": pathVaultToken,
		},
	})
	require.NoError(t, err)

	// Create sink file
	pathSinkFile := makeTempFile(t, "sink-file", "")
	require.NoError(t, err)

	ahConfig := &auth.AuthHandlerConfig{
		Logger:                       logger.Named("auth.handler"),
		Client:                       serverClient,
		EnableExecTokenCh:            true,
		EnableTemplateTokenCh:        true,
		EnableReauthOnNewCredentials: true,
		ExitOnError:                  false,
	}
	ah := auth.NewAuthHandler(ahConfig)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": pathSinkFile,
		},
	}
	fs, err := file.NewFileSink(config)
	require.NoError(t, err)
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: serverClient,
	})
	go func() {
		errCh <- ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config}, ah.AuthInProgress)
	}()

	// Create template server
	sc := &template.ServerConfig{
		Logger: logger.Named("template.server"),
		AgentConfig: &agentConfig.Config{
			Vault: &agentConfig.Vault{
				Address:       serverClient.Address(),
				TLSSkipVerify: true,
			},
			TemplateConfig: &agentConfig.TemplateConfig{
				StaticSecretRenderInt: 1 * time.Second,
			},
			AutoAuth: &agentConfig.AutoAuth{
				Sinks: []*agentConfig.Sink{
					{
						Type: "file",
						Config: map[string]interface{}{
							"path": pathSinkFile,
						},
					},
				},
			},
			ExitAfterAuth: false,
		},
		LogLevel:      hclog.Trace,
		LogWriter:     hclog.DefaultOutput,
		ExitAfterAuth: false,
	}

	pathTemplateOutput := makeTempFile(t, "template-output", "")
	require.NoError(t, err)
	templateTest := &ctconfig.TemplateConfig{
		Contents:    pointerutil.StringPtr(`{{ with secret "auth/token/lookup-self" }}{{ .Data.id }}{{ end }}`),
		Destination: pointerutil.StringPtr(pathTemplateOutput),
	}
	templatesToRender := []*ctconfig.TemplateConfig{templateTest}

	server := template.NewServer(sc)
	go func() {
		errCh <- server.Run(ctx, ah.TemplateTokenCh, templatesToRender, ah.AuthInProgress, ah.InvalidToken)
	}()

	// Send token to template channel, and wait for the template to render
	ah.TemplateTokenCh <- token
	err = waitForFileContent(t, pathTemplateOutput, token)

	// Revoke Token
	err = serverClient.Auth().Token().RevokeOrphan(token)
	require.NoError(t, err)

	// Create new token
	tokenSecret, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{})
	require.NoError(t, err)
	require.NotNil(t, tokenSecret)
	require.NotNil(t, tokenSecret.Auth)
	require.NotEmpty(t, tokenSecret.Auth.ClientToken)
	newToken := tokenSecret.Auth.ClientToken

	// Write token to file
	err = os.WriteFile(pathVaultToken, []byte(newToken), 0o600)
	require.NoError(t, err)

	// Wait for auto-auth to complete and verify token has been written to the sink
	// and the template has been re-rendered
	err = waitForFileContent(t, pathSinkFile, newToken)
	require.NoError(t, err)

	err = waitForFileContent(t, pathTemplateOutput, newToken)
	require.NoError(t, err)

	// Calling cancel will stop the 'Run' funcs we started in Goroutines, we should
	// then check that there were no errors in our channel.
	cancel()
	wrapUpTimeout := 5 * time.Second
	for {
		select {
		case <-time.After(wrapUpTimeout):
			t.Fatal("test timed out")
		case err := <-errCh:
			require.NoError(t, err)
		case <-ctx.Done():
			// We can finish the test ourselves
			return
		}
	}
}

// Test_NoAutoAuthSelfHealing_BadPolicy tests that auto auth
// is not re-triggered if a token with incorrect policy access
// is used to render a template
func Test_NoAutoAuthSelfHealing_BadPolicy(t *testing.T) {
	// Unset the environment variable so that agent picks up the right test cluster address
	t.Setenv(api.EnvVaultAddress, "")

	policyName := "kv-access"

	cluster := minimal.NewTestSoloCluster(t, nil)
	logger := corehelpers.NewTestLogger(t)
	serverClient := cluster.Cores[0].Client

	// Write a policy with correct access to the secrets
	err := serverClient.Sys().PutPolicy(policyName, `
path "/kv/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}
path "/secret/*" {
	capabilities = ["create", "read", "update", "delete", "list"]
}`)
	require.NoError(t, err)

	// Create a token without enough policy access to the kv secrets
	secret, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"default"},
	})
	require.NoError(t, err)
	require.NotNil(t, secret)
	require.NotNil(t, secret.Auth)
	require.NotEmpty(t, secret.Auth.ClientToken)
	require.Len(t, secret.Auth.Policies, 1)
	require.Contains(t, secret.Auth.Policies, "default")
	token := secret.Auth.ClientToken

	// Write token to vault-token file
	pathVaultToken := makeTempFile(t, "vault-token", token)

	// Give us some leeway of 3 errors 1 from each of: auth handler, sink server template server.
	errCh := make(chan error, 3)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// Create auth handler
	am, err := tokenfile.NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path": pathVaultToken,
		},
	})
	require.NoError(t, err)

	ahConfig := &auth.AuthHandlerConfig{
		Logger:                       logger.Named("auth.handler"),
		Client:                       serverClient,
		EnableExecTokenCh:            true,
		EnableReauthOnNewCredentials: true,
		ExitOnError:                  false,
	}
	ah := auth.NewAuthHandler(ahConfig)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()

	// Create sink file
	pathSinkFile := makeTempFile(t, "sink-file", "")

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": pathSinkFile,
		},
	}
	fs, err := file.NewFileSink(config)
	require.NoError(t, err)
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: serverClient,
	})
	go func() {
		errCh <- ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config}, ah.AuthInProgress)
	}()

	// Create template server
	sc := template.ServerConfig{
		Logger: logger.Named("template.server"),
		AgentConfig: &agentConfig.Config{
			Vault: &agentConfig.Vault{
				Address:       serverClient.Address(),
				TLSSkipVerify: true,
			},
			TemplateConfig: &agentConfig.TemplateConfig{
				StaticSecretRenderInt: 1 * time.Second,
			},
			// Need to create at least one sink output so that it does not exit after rendering
			AutoAuth: &agentConfig.AutoAuth{
				Sinks: []*agentConfig.Sink{
					{
						Type: "file",
						Config: map[string]interface{}{
							"path": pathSinkFile,
						},
					},
				},
			},
			ExitAfterAuth: false,
		},
		LogLevel:      hclog.Trace,
		LogWriter:     hclog.DefaultOutput,
		ExitAfterAuth: false,
	}

	pathTemplateDestination := makeTempFile(t, "kv-data", "")
	templateTest := &ctconfig.TemplateConfig{
		Contents:    pointerutil.StringPtr(`"{{ with secret "secret/data/otherapp" }}{{ .Data.data.username }}{{ end }}"`),
		Destination: pointerutil.StringPtr(pathTemplateDestination),
	}
	templatesToRender := []*ctconfig.TemplateConfig{templateTest}

	server := template.NewServer(&sc)
	go func() {
		errCh <- server.Run(ctx, ah.TemplateTokenCh, templatesToRender, ah.AuthInProgress, ah.InvalidToken)
	}()

	// Send token to the template channel
	ah.TemplateTokenCh <- token

	// Create new token with the correct policy access
	tokenSecret, err := serverClient.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{policyName},
	})
	require.NoError(t, err)
	require.NotNil(t, tokenSecret)
	require.NotNil(t, tokenSecret.Auth)
	require.NotEmpty(t, tokenSecret.Auth.ClientToken)
	require.Len(t, tokenSecret.Auth.Policies, 2)
	require.Contains(t, tokenSecret.Auth.Policies, "default")
	require.Contains(t, tokenSecret.Auth.Policies, policyName)
	newToken := tokenSecret.Auth.ClientToken

	// Write new token to token file (where Agent would re-auto-auth from if
	// it were triggered)
	err = os.WriteFile(pathVaultToken, []byte(newToken), 0o600)
	require.NoError(t, err)

	// Wait for any potential *incorrect* re-triggers of auto auth
	time.Sleep(time.Second * 3)

	// Auto auth should not have been re-triggered because of just a permission denied error
	// Verify that the new token has NOT been written to the token sink
	tokenInSink, err := os.ReadFile(pathSinkFile)
	require.NoError(t, err)
	require.Equal(t, token, string(tokenInSink))

	// Validate that the template still hasn't been rendered.
	templateContent, err := os.ReadFile(pathTemplateDestination)
	require.NoError(t, err)
	require.Equal(t, "", string(templateContent))

	cancel()
	wrapUpTimeout := 5 * time.Second
	for {
		select {
		case <-time.After(wrapUpTimeout):
			t.Fatal("test timed out")
		case err := <-errCh:
			require.NoError(t, err)
		case <-ctx.Done():
			// We can finish the test ourselves
			return
		}
	}
}

// waitForFileContent waits for the file at filePath to exist and contain fileContent
// or it will return in an error. Waits for five seconds, with 100ms intervals.
// Returns nil if content became the same, or non-nil if it didn't.
func waitForFileContent(t *testing.T, filePath, expectedContent string) error {
	t.Helper()

	var err error
	tick := time.Tick(100 * time.Millisecond)
	timeout := time.After(5 * time.Second)
	// We need to wait for the files to be updated...
	for {
		select {
		case <-timeout:
			return fmt.Errorf("timed out waiting for file content, last error: %w", err)
		case <-tick:
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return err
		}

		stringContent := string(content)
		if stringContent != expectedContent {
			err = fmt.Errorf("content not yet the same, expectedContent=%s, content=%s", expectedContent, stringContent)
			continue
		}

		return nil
	}
}

// makeTempFile creates a temp file with the specified name, populates it with the
// supplied contents and closes it. The path to the file is returned, also the file
// will be automatically removed when the test which created it, finishes.
func makeTempFile(t *testing.T, name, contents string) string {
	t.Helper()

	f, err := os.Create(filepath.Join(t.TempDir(), name))
	require.NoError(t, err)
	path := f.Name()

	_, err = f.WriteString(contents)
	require.NoError(t, err)

	err = f.Close()
	require.NoError(t, err)

	return path
}
