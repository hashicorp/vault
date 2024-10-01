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
	vaultoci "github.com/hashicorp/vault-plugin-auth-oci"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	agentoci "github.com/hashicorp/vault/command/agentproxyshared/auth/oci"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	envVarOCITestTenancyOCID    = "OCI_TEST_TENANCY_OCID"
	envVarOCITestUserOCID       = "OCI_TEST_USER_OCID"
	envVarOCITestFingerprint    = "OCI_TEST_FINGERPRINT"
	envVarOCITestPrivateKeyPath = "OCI_TEST_PRIVATE_KEY_PATH"
	envVAROCITestOCIDList       = "OCI_TEST_OCID_LIST"

	// The OCI SDK doesn't export its standard env vars so they're captured here.
	// These are used for the duration of the test to make sure the agent is able to
	// pick up creds from the env.
	//
	// To run this test, do not set these. Only the above ones need to be set.
	envVarOCITenancyOCID    = "OCI_tenancy_ocid"
	envVarOCIUserOCID       = "OCI_user_ocid"
	envVarOCIFingerprint    = "OCI_fingerprint"
	envVarOCIPrivateKeyPath = "OCI_private_key_path"
)

func TestOCIEndToEnd(t *testing.T) {
	if !runAcceptanceTests {
		t.SkipNow()
	}

	// Ensure each cred is populated.
	credNames := []string{
		envVarOCITestTenancyOCID,
		envVarOCITestUserOCID,
		envVarOCITestFingerprint,
		envVarOCITestPrivateKeyPath,
		envVAROCITestOCIDList,
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"oci": vaultoci.Factory,
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
	if err := client.Sys().EnableAuthWithOptions("oci", &api.EnableAuthOptions{
		Type: "oci",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Logical().Write("auth/oci/config", map[string]interface{}{
		"home_tenancy_id": os.Getenv(envVarOCITestTenancyOCID),
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Logical().Write("auth/oci/role/test", map[string]interface{}{
		"ocid_list": os.Getenv(envVAROCITestOCIDList),
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// We're going to feed oci auth creds via env variables.
	if err := setOCIEnvCreds(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := unsetOCIEnvCreds(); err != nil {
			t.Fatal(err)
		}
	}()

	vaultAddr := "http://" + cluster.Cores[0].Listeners[0].Addr().String()

	am, err := agentoci.NewOCIAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.oci"),
		MountPath: "auth/oci",
		Config: map[string]interface{}{
			"type": "apikey",
			"role": "test",
		},
	}, vaultAddr)
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
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

	token, err := readToken(tokenSinkFileName)
	if err != nil {
		t.Fatal(err)
	}

	if token.Token == "" {
		t.Fatal("expected token but didn't receive it")
	}
}

func setOCIEnvCreds() error {
	if err := os.Setenv(envVarOCITenancyOCID, os.Getenv(envVarOCITestTenancyOCID)); err != nil {
		return err
	}
	if err := os.Setenv(envVarOCIUserOCID, os.Getenv(envVarOCITestUserOCID)); err != nil {
		return err
	}
	if err := os.Setenv(envVarOCIFingerprint, os.Getenv(envVarOCITestFingerprint)); err != nil {
		return err
	}
	return os.Setenv(envVarOCIPrivateKeyPath, os.Getenv(envVarOCITestPrivateKeyPath))
}

func unsetOCIEnvCreds() error {
	if err := os.Unsetenv(envVarOCITenancyOCID); err != nil {
		return err
	}
	if err := os.Unsetenv(envVarOCIUserOCID); err != nil {
		return err
	}
	if err := os.Unsetenv(envVarOCIFingerprint); err != nil {
		return err
	}
	return os.Unsetenv(envVarOCIPrivateKeyPath)
}
