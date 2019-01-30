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
	hclog "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"
	vaultalicloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	agentalicloud "github.com/hashicorp/vault/command/agent/auth/alicloud"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
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

	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger,
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

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

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
		Logger:    logger.Named("auth.alicloud"),
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
		Logger: logger.Named("auth.handler"),
		Client: client,
	}

	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
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
	go ss.Run(ctx, ah.OutputCh, []*sink.SinkConfig{config})
	defer func() {
		<-ss.DoneCh
	}()

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
	assumeRoleReq.RoleSessionName = strings.Replace(roleSessionName, "-", "", -1)
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
