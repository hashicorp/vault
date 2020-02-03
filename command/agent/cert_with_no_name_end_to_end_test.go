package agent

import (
	"context"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/api"
	vaultcert "github.com/hashicorp/vault/builtin/credential/cert"
	"github.com/hashicorp/vault/command/agent/auth"
	agentcert "github.com/hashicorp/vault/command/agent/auth/cert"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/helper/dhutil"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestCertWithNoNAmeEndToEnd(t *testing.T) {
	testCertWithNoNAmeEndToEnd(t, false)
	testCertWithNoNAmeEndToEnd(t, true)
}

func testCertWithNoNAmeEndToEnd(t *testing.T, ahWrapping bool) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger,
		CredentialBackends: map[string]logical.Factory{
			"cert": vaultcert.Factory,
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
	err := client.Sys().EnableAuthWithOptions("cert", &api.EnableAuthOptions{
		Type: "cert",
	})
	if err != nil {
		t.Fatal(err)
	}

	certificatePEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cluster.CACert.Raw})

	_, err = client.Logical().Write("auth/cert/certs/test", map[string]interface{}{
		"name":        "test",
		"certificate": string(certificatePEM),
		"policies":    "default",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Generate encryption params
	pub, pri, err := dhutil.GeneratePublicPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	out := ouf.Name()
	ouf.Close()
	os.Remove(out)
	t.Logf("output: %s", out)

	dhpathf, err := ioutil.TempFile("", "auth.dhpath.test.")
	if err != nil {
		t.Fatal(err)
	}
	dhpath := dhpathf.Name()
	dhpathf.Close()
	os.Remove(dhpath)

	// Write DH public key to file
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

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	am, err := agentcert.NewCertAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.cert"),
		MountPath: "auth/cert",
	})
	if err != nil {
		t.Fatal(err)
	}

	ahConfig := &auth.AuthHandlerConfig{
		Logger:                       logger.Named("auth.handler"),
		Client:                       client,
		EnableReauthOnNewCredentials: true,
	}
	if ahWrapping {
		ahConfig.WrapTTL = 10 * time.Second
	}
	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
	}()

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		AAD:    "foobar",
		DHType: "curve25519",
		DHPath: dhpath,
		Config: map[string]interface{}{
			"path": out,
		},
	}
	if !ahWrapping {
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
	defer cancelFunc()

	cloned, err := client.Clone()
	if err != nil {
		t.Fatal(err)
	}

	checkToken := func() string {
		timeout := time.Now().Add(5 * time.Second)
		for {
			if time.Now().After(timeout) {
				t.Fatal("did not find a written token after timeout")
			}
			val, err := ioutil.ReadFile(out)
			if err == nil {
				os.Remove(out)
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
				case ahWrapping && wrapInfo.CreationPath != "auth/cert/login":
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
	checkToken()
}
