// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	agentAppRole "github.com/hashicorp/vault/command/agentproxyshared/auth/approle"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTokenPreload_UsingAutoAuth(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	// Setup Vault
	if err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	}); err != nil {
		t.Fatal(err)
	}

	// Setup Approle
	_, err := client.Logical().Write("auth/approle/role/test1", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "3s",
		"token_max_ttl":  "10s",
		"policies":       []string{"test-autoauth"},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Logical().Write("auth/approle/role/test1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID1 := resp.Data["secret_id"].(string)

	resp, err = client.Logical().Read("auth/approle/role/test1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID1 := resp.Data["role_id"].(string)

	rolef, err := ioutil.TempFile("", "auth.role-id.test.")
	if err != nil {
		t.Fatal(err)
	}
	role := rolef.Name()
	rolef.Close() // WriteFile doesn't need it open
	defer os.Remove(role)
	t.Logf("input role_id_file_path: %s", role)

	secretf, err := ioutil.TempFile("", "auth.secret-id.test.")
	if err != nil {
		t.Fatal(err)
	}
	secret := secretf.Name()
	secretf.Close()
	defer os.Remove(secret)
	t.Logf("input secret_id_file_path: %s", secret)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	conf := map[string]interface{}{
		"role_id_file_path":   role,
		"secret_id_file_path": secret,
	}

	if err := ioutil.WriteFile(role, []byte(roleID1), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test role 1", "path", role)
	}

	if err := ioutil.WriteFile(secret, []byte(secretID1), 0o600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test secret 1", "path", secret)
	}

	// Setup Preload Token
	tokenRespRaw, err := client.Logical().Write("auth/token/create", map[string]interface{}{
		"ttl":              "10s",
		"explicit-max-ttl": "15s",
		"policies":         []string{""},
	})
	if err != nil {
		t.Fatal(err)
	}

	if tokenRespRaw.Auth == nil || tokenRespRaw.Auth.ClientToken == "" {
		t.Fatal("expected token but got none")
	}
	token := tokenRespRaw.Auth.ClientToken

	am, err := agentAppRole.NewApproleAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.approle"),
		MountPath: "auth/approle",
		Config:    conf,
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
		Token:  token,
	}

	ah := auth.NewAuthHandler(ahConfig)

	tmpFile, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	tokenSinkFileName := tmpFile.Name()
	tmpFile.Close()
	os.Remove(tokenSinkFileName)
	t.Logf("output: %s", tokenSinkFileName)

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": tokenSinkFileName,
		},
		WrapTTL: 10 * time.Second,
	}

	fs, err := file.NewFileSink(config)
	if err != nil {
		t.Fatal(err)
	}
	config.Sink = fs

	ss := sink.NewSinkServer(&sink.SinkServerConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})

	errCh := make(chan error)
	go func() {
		errCh <- ah.Run(ctx, am)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	go func() {
		errCh <- ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config}, ah.AuthInProgress)
	}()
	defer func() {
		select {
		case <-ctx.Done():
		case err := <-errCh:
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	// This has to be after the other defers so it happens first. It allows
	// successful test runs to immediately cancel all of the runner goroutines
	// and unblock any of the blocking defer calls by the runner's DoneCh that
	// comes before this and avoid successful tests from taking the entire
	// timeout duration.
	defer cancel()

	if stat, err := os.Lstat(tokenSinkFileName); err == nil {
		t.Fatalf("expected err but got %s", stat)
	} else if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	// Wait 2 seconds for the env variables to be detected and an auth to be generated.
	time.Sleep(time.Second * 2)

	authToken, err := readToken(tokenSinkFileName)
	if err != nil {
		t.Fatal(err)
	}

	if authToken.Token == "" {
		t.Fatal("expected token but didn't receive it")
	}

	wrappedToken := map[string]interface{}{
		"token": authToken.Token,
	}
	unwrapResp, err := client.Logical().Write("sys/wrapping/unwrap", wrappedToken)
	if err != nil {
		t.Fatalf("error unwrapping token: %s", err)
	}

	sinkToken, ok := unwrapResp.Data["token"].(string)
	if !ok {
		t.Fatal("expected token but didn't receive it")
	}

	if sinkToken != token {
		t.Fatalf("auth token and preload token should be the same: expected: %s, actual: %s", token, sinkToken)
	}
}
