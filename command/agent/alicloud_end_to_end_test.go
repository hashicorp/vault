package agent

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/hashicorp/go-uuid"
	vaultalicloud "github.com/hashicorp/vault-plugin-auth-alicloud"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	agentalicloud "github.com/hashicorp/vault/command/agent/auth/alicloud"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/helper/dhutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

const (
	envVarRunAccTests = "VAULT_ACC"
	envVarAccessKey   = "ALICLOUD_TEST_ACCESS_KEY"
	envVarSecretKey   = "ALICLOUD_TEST_SECRET_KEY"
	envVarRoleArn     = "ALICLOUD_TEST_ROLE_ARN"
)

var runAcceptanceTests = os.Getenv(envVarRunAccTests) == "1"

func TestAliCloudEndToEnd(t *testing.T) {
	if !runAcceptanceTests {
		t.SkipNow()
	}

	testAliCloudEndToEnd(t, false)
	testAliCloudEndToEnd(t, true)
}

func testAliCloudEndToEnd(t *testing.T, ahWrapping bool) {
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
		"arn": os.Getenv(envVarRoleArn),
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
		// Let's make sure we unset these when the test is done
		os.Unsetenv(providers.EnvVarAccessKeyID)
		os.Unsetenv(providers.EnvVarAccessKeySecret)
		os.Unsetenv(providers.EnvVarAccessKeyStsToken)
	}()

	am, err := agentalicloud.NewAliCloudAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.alicloud"),
		MountPath: "auth/alicloud",
		Config: map[string]interface{}{
			"role":                    "test",
			"region":                  "us-west-1",
			"cred_check_freq_seconds": 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	}
	if ahWrapping {
		ahConfig.WrapTTL = 10 * time.Second
	}
	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
		fmt.Print("ah.DoneCh closed") // TODO stripme
	}()

	// Set up the token sink, which is where generated tokens should be sent.
	pub, pri, err := dhutil.GeneratePublicPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	tokenSinkFile := ouf.Name()
	ouf.Close() // TODO it confuses me that the outfile is being closed here, and the tokenSinkFile is being removed - the whole flow here confuses me
	os.Remove(tokenSinkFile)
	t.Logf("output: %s", tokenSinkFile)

	dhpathf, err := ioutil.TempFile("", "auth.dhpath.test.") // TODO what is this?
	if err != nil {
		t.Fatal(err)
	}
	dhpath := dhpathf.Name()
	dhpathf.Close()
	os.Remove(dhpath)

	mPubKey, err := jsonutil.EncodeJSON(&dhutil.PublicKeyInfo{
		Curve25519PublicKey: pub,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(dhpath, mPubKey, 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote dh param file", "path", dhpath)
	}

	// TODO need to read what this config is all about
	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		AAD:    "foobar",
		DHType: "curve25519",
		DHPath: dhpath,
		Config: map[string]interface{}{
			"path": tokenSinkFile,
		},
	}
	if !ahWrapping {
		// TODO what is this?
		config.WrapTTL = 10 * time.Second
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

	// This has to be after the other defers so it happens first
	// TODO is this being overwritten by the earlier cancelFunc? naming clash?
	defer cancelFunc()

	if stat, err := os.Lstat(tokenSinkFile); err == nil {
		t.Fatalf("expected err but got %s", stat)
	} else if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	// TODO why are we doing this?
	cloned, err := client.Clone()
	if err != nil {
		t.Fatal(err)
	}

	// Rotate the env vars
	if err := setAliCloudEnvCreds(); err != nil {
		t.Fatal(err)
	}

	checkToken := func() string {
		timeout := time.Now().Add(5 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("did not find a written token after timeout")
			}
			val, err := ioutil.ReadFile(tokenSinkFile)
			if err == nil {
				os.Remove(tokenSinkFile)
				if len(val) == 0 {
					t.Fatal("written token was empty")
				}

				// First decrypt it
				resp := new(dhutil.Envelope)
				if err := jsonutil.DecodeJSON(val, resp); err != nil {
					continue
				}

				aesKey, err := dhutil.GenerateSharedKey(pri, resp.Curve25519PublicKey)
				if err != nil {
					t.Fatal(err)
				}
				if len(aesKey) == 0 {
					t.Fatal("got empty aes key")
				}

				val, err = dhutil.DecryptAES(aesKey, resp.EncryptedPayload, resp.Nonce, []byte("foobar"))
				if err != nil {
					t.Fatalf("error: %v\nresp: %v", err, string(val))
				}

				// Now unwrap it
				wrapInfo := new(api.SecretWrapInfo)
				if err := jsonutil.DecodeJSON(val, wrapInfo); err != nil {
					t.Fatal(err)
				}
				switch {
				case wrapInfo.TTL != 10:
					t.Fatalf("bad wrap info: %v", wrapInfo.TTL)
				case !ahWrapping && wrapInfo.CreationPath != "sys/wrapping/wrap":
					t.Fatalf("bad wrap path: %v", wrapInfo.CreationPath)
				case ahWrapping && wrapInfo.CreationPath != "auth/alicloud/login":
					t.Fatalf("bad wrap path: %v", wrapInfo.CreationPath)
				case wrapInfo.Token == "":
					t.Fatal("wrap token is empty")
				}
				cloned.SetToken(wrapInfo.Token)
				secret, err := cloned.Logical().Unwrap("")
				if err != nil {
					t.Fatal(err)
				}
				if ahWrapping {
					switch {
					case secret.Auth == nil:
						t.Fatal("unwrap secret auth is nil")
					case secret.Auth.ClientToken == "":
						t.Fatal("unwrap token is nil")
					}
					return secret.Auth.ClientToken
				} else {
					switch {
					case secret.Data == nil:
						t.Fatal("unwrap secret data is nil")
					case secret.Data["token"] == nil:
						t.Fatal("unwrap token is nil")
					}
					return secret.Data["token"].(string)
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
	origToken := checkToken()

	// We only check this if the renewer is actually renewing for us
	if !ahWrapping {
		// Period of 3 seconds, so should still be alive after 7
		timeout := time.Now().Add(7 * time.Second)
		cloned.SetToken(origToken)
		for {
			if time.Now().After(timeout) {
				break
			}
			secret, err := cloned.Auth().Token().LookupSelf()
			if err != nil {
				t.Fatal(err)
			}
			ttl, err := secret.Data["ttl"].(json.Number).Int64()
			if err != nil {
				t.Fatal(err)
			}
			if ttl > 3 {
				t.Fatalf("unexpected ttl: %v", secret.Data["ttl"])
			}
		}
	}

	// Rotate the env vars again
	if err := setAliCloudEnvCreds(); err != nil {
		t.Fatal(err)
	}

	newToken := checkToken()
	if newToken == origToken {
		t.Fatal("found same token written")
	}

	if !ahWrapping {
		// Repeat the period test. At the end the old token should have expired and
		// the new token should still be alive after 7
		timeout := time.Now().Add(7 * time.Second)
		cloned.SetToken(newToken)
		for {
			if time.Now().After(timeout) {
				break
			}
			secret, err := cloned.Auth().Token().LookupSelf()
			if err != nil {
				t.Fatal(err)
			}
			ttl, err := secret.Data["ttl"].(json.Number).Int64()
			if err != nil {
				t.Fatal(err)
			}
			if ttl > 3 {
				t.Fatalf("unexpected ttl: %v", secret.Data["ttl"])
			}
		}

		cloned.SetToken(origToken)
		_, err = cloned.Auth().Token().LookupSelf()
		if err == nil {
			t.Fatal("expected error")
		}
	}
}

func setAliCloudEnvCreds() error {
	config := sdk.NewConfig()
	config.Scheme = "https"
	client, err := sts.NewClientWithOptions("us-west-1", config, credentials.NewAccessKeyCredential(os.Getenv(envVarAccessKey), os.Getenv(envVarSecretKey)))
	if err != nil {
		return err
	}
	roleSessionName, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	assumeRoleReq := sts.CreateAssumeRoleRequest()
	assumeRoleReq.RoleArn = os.Getenv(envVarRoleArn)
	assumeRoleReq.RoleSessionName = strings.Replace(roleSessionName, "-", "", -1)
	assumeRoleResp, err := client.AssumeRole(assumeRoleReq)
	if err != nil {
		return err
	}

	os.Setenv(providers.EnvVarAccessKeyID, assumeRoleResp.Credentials.AccessKeyId)
	os.Setenv(providers.EnvVarAccessKeySecret, assumeRoleResp.Credentials.AccessKeySecret)
	os.Setenv(providers.EnvVarAccessKeyStsToken, assumeRoleResp.Credentials.SecurityToken)
	return nil
}
