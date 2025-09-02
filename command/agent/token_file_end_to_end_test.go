// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	token_file "github.com/hashicorp/vault/command/agentproxyshared/auth/token-file"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault"
)

func TestTokenFileEndToEnd(t *testing.T) {
	logger := logging.NewVaultLogger(log.Trace)
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	secret, err := client.Auth().Token().Create(nil)
	if err != nil || secret == nil {
		t.Fatal(err)
	}

	tokenFile, err := os.Create(filepath.Join(t.TempDir(), "token_file"))
	if err != nil {
		t.Fatal(err)
	}
	tokenFileName := tokenFile.Name()
	tokenFile.Close() // WriteFile doesn't need it open
	os.WriteFile(tokenFileName, []byte(secret.Auth.ClientToken), 0o666)
	defer os.Remove(tokenFileName)

	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	am, err := token_file.NewTokenFileAuthMethod(&auth.AuthConfig{
		Logger: logger.Named("auth.method"),
		Config: map[string]interface{}{
			"token_file_path": tokenFileName,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ah := auth.NewAuthHandler(ahConfig)
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

	// We close these right away because we're just basically testing
	// permissions and finding a usable file name
	sinkFile, err := os.Create(filepath.Join(t.TempDir(), "auth.tokensink.test."))
	if err != nil {
		t.Fatal(err)
	}
	tokenSinkFileName := sinkFile.Name()
	sinkFile.Close()
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

	// This has to be after the other defers, so it happens first. It allows
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

	token, err := readToken(tokenSinkFileName)
	if err != nil {
		t.Fatal(err)
	}

	if token.Token == "" {
		t.Fatal("expected token but didn't receive it")
	}

	_, err = os.Stat(tokenFileName)
	if err != nil {
		t.Fatal("Token file removed")
	}
}
