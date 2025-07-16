// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package agent

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	vaultaws "github.com/hashicorp/vault/builtin/credential/aws"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
	agentaws "github.com/hashicorp/vault/command/agentproxyshared/auth/aws"
	"github.com/hashicorp/vault/command/agentproxyshared/sink"
	"github.com/hashicorp/vault/command/agentproxyshared/sink/file"
	"github.com/hashicorp/vault/helper/testhelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	// These are the access key and secret that should be used when calling "AssumeRole"
	// for the given AWS_TEST_ROLE_ARN.
	envVarAwsTestAccessKey = "AWS_TEST_ACCESS_KEY"
	envVarAwsTestSecretKey = "AWS_TEST_SECRET_KEY"
	envVarAwsTestRoleArn   = "AWS_TEST_ROLE_ARN"

	// The AWS SDK doesn't export its standard env vars so they're captured here.
	// These are used for the duration of the test to make sure the agent is able to
	// pick up creds from the env.
	//
	// To run this test, do not set these. Only the above ones need to be set.
	envVarAwsAccessKey    = "AWS_ACCESS_KEY_ID"
	envVarAwsSecretKey    = "AWS_SECRET_ACCESS_KEY"
	envVarAwsSessionToken = "AWS_SESSION_TOKEN"
)

func TestAWSEndToEnd(t *testing.T) {
	if !runAcceptanceTests {
		t.SkipNow()
	}

	// Ensure each cred is populated.
	credNames := []string{
		envVarAwsTestAccessKey,
		envVarAwsTestSecretKey,
		envVarAwsTestRoleArn,
	}
	testhelpers.SkipUnlessEnvVarsSet(t, credNames)

	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"aws": vaultaws.Factory,
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
	if err := client.Sys().EnableAuthWithOptions("aws", &api.EnableAuthOptions{
		Type: "aws",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := client.Logical().Write("auth/aws/role/test", map[string]interface{}{
		"auth_type": "iam",
		"policies":  "default",
		// Retain thru the account number of the given arn and wildcard the rest.
		"bound_iam_principal_arn": os.Getenv(envVarAwsTestRoleArn)[:25] + "*",
	}); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// We're going to feed aws auth creds via env variables.
	if err := setAwsEnvCreds(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := unsetAwsEnvCreds(); err != nil {
			t.Fatal(err)
		}
	}()

	am, err := agentaws.NewAWSAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.aws"),
		MountPath: "auth/aws",
		Config: map[string]interface{}{
			"role":                     "test",
			"type":                     "iam",
			"credential_poll_interval": 1,
		},
	})
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

func setAwsEnvCreds() error {
	cfg := &aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv(envVarAwsTestAccessKey), os.Getenv(envVarAwsTestSecretKey), ""),
	}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return err
	}
	client := sts.New(sess)

	roleArn := os.Getenv(envVarAwsTestRoleArn)
	uid, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	input := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &uid,
	}
	output, err := client.AssumeRole(input)
	if err != nil {
		return err
	}

	if err := os.Setenv(envVarAwsAccessKey, *output.Credentials.AccessKeyId); err != nil {
		return err
	}
	if err := os.Setenv(envVarAwsSecretKey, *output.Credentials.SecretAccessKey); err != nil {
		return err
	}
	return os.Setenv(envVarAwsSessionToken, *output.Credentials.SessionToken)
}

func unsetAwsEnvCreds() error {
	if err := os.Unsetenv(envVarAwsAccessKey); err != nil {
		return err
	}
	if err := os.Unsetenv(envVarAwsSecretKey); err != nil {
		return err
	}
	return os.Unsetenv(envVarAwsSessionToken)
}
