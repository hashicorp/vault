// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials/providers"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/hashicorp/go-uuid"
	vaultalicloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	agentalicloud "github.com/hashicorp/vault/command/agentproxyshared/auth/alicloud"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	envVarAlicloudAccessKey = "ALICLOUD_TEST_ACCESS_KEY"
	envVarAlicloudSecretKey = "ALICLOUD_TEST_SECRET_KEY"
	envVarAlicloudRoleArn   = "ALICLOUD_TEST_ROLE_ARN"
)

func TestAliCloudEndToEnd(t *testing.T) {
	if !runAcceptanceTests {
		t.SkipNow()
	}

	// Ensure each cred is populated.
	credNames := []string{
		envVarAlicloudAccessKey,
		envVarAlicloudSecretKey,
		envVarAlicloudRoleArn,
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"alicloud": vaultalicloud.Factory,
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
	if err := client.Sys().EnableAuthWithOptions("alicloud", &api.EnableAuthOptions{
		Type: "alicloud",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Logical().Write("auth/alicloud/role/test", map[string]interface{}{
		"arn": os.Getenv(envVarAlicloudRoleArn),
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// We're going to feed alicloud auth creds via env variables.
	if err := setAliCloudEnvCreds(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := unsetAliCloudEnvCreds(); err != nil {
			t.Fatal(err)
		}
	}()

	am, err := agentalicloud.NewAliCloudAuthMethod(&auth.AuthConfig{
		Logger:    cluster.Logger.Named("auth.alicloud"),
		MountPath: "auth/alicloud",
		Config: map[string]interface{}{
			"role":                     "test",
			"region":                   "us-west-1",
			"credential_poll_interval": 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger: cluster.Logger.Named("auth.handler"),
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
		Logger: cluster.Logger.Named("sink.file"),
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
		Logger: cluster.Logger.Named("sink.server"),
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

func setAliCloudEnvCreds() error {
	config := sdk.NewConfig()
	config.Scheme = "https"
	client, err := sts.NewClientWithOptions("us-west-1", config, credentials.NewAccessKeyCredential(os.Getenv(envVarAlicloudAccessKey), os.Getenv(envVarAlicloudSecretKey)))
	if err != nil {
		return err
	}
	roleSessionName, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	assumeRoleReq := sts.CreateAssumeRoleRequest()
	assumeRoleReq.RoleArn = os.Getenv(envVarAlicloudRoleArn)
	assumeRoleReq.RoleSessionName = strings.ReplaceAll(roleSessionName, "-", "")
	assumeRoleResp, err := client.AssumeRole(assumeRoleReq)
	if err != nil {
		return err
	}

	if err := os.Setenv(providers.EnvVarAccessKeyID, assumeRoleResp.Credentials.AccessKeyId); err != nil {
		return err
	}
	if err := os.Setenv(providers.EnvVarAccessKeySecret, assumeRoleResp.Credentials.AccessKeySecret); err != nil {
		return err
	}
	return os.Setenv(providers.EnvVarAccessKeyStsToken, assumeRoleResp.Credentials.SecurityToken)
}

func unsetAliCloudEnvCreds() error {
	if err := os.Unsetenv(providers.EnvVarAccessKeyID); err != nil {
		return err
	}
	if err := os.Unsetenv(providers.EnvVarAccessKeySecret); err != nil {
		return err
	}
	return os.Unsetenv(providers.EnvVarAccessKeyStsToken)
}
