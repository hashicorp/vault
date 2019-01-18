package agent

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	"github.com/hashicorp/vault/command/agent/auth"
	agentapprole "github.com/hashicorp/vault/command/agent/auth/approle"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestAppRoleEndToEnd(t *testing.T) {
	t.Run("preserve_secret_id_file", func(t *testing.T) {
		testAppRoleEndToEnd(t, false)
	})

	t.Run("remove_secret_id_file", func(t *testing.T) {
		testAppRoleEndToEnd(t, true)
	})
}

func testAppRoleEndToEnd(t *testing.T, removeSecretIDFile bool) {
	var err error
	logger := logging.NewVaultLogger(log.Trace)
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test1", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "3s",
		"token_max_ttl":  "10s",
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

	_, err = client.Logical().Write("auth/approle/role/test2", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "3s",
		"token_max_ttl":  "10s",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = client.Logical().Write("auth/approle/role/test2/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID2 := resp.Data["secret_id"].(string)

	resp, err = client.Logical().Read("auth/approle/role/test2/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID2 := resp.Data["role_id"].(string)

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

	// We close these right away because we're just basically testing
	// permissions and finding a usable file name
	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	out := ouf.Name()
	ouf.Close()
	os.Remove(out)
	t.Logf("output: %s", out)

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	conf := map[string]interface{}{
		"role_id_file_path":   role,
		"secret_id_file_path": secret,
	}
	if !removeSecretIDFile {
		conf["remove_secret_id_file_after_reading"] = removeSecretIDFile
	}

	am, err := agentapprole.NewApproleAuthMethod(&auth.AuthConfig{
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
	}
	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
	}()

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": out,
		},
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

	// Check that no sink file exists
	_, err = os.Lstat(out)
	if err == nil {
		t.Fatal("expected err")
	}
	if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	if err := ioutil.WriteFile(role, []byte(roleID1), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test role 1", "path", role)
	}

	if err := ioutil.WriteFile(secret, []byte(secretID1), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test secret 1", "path", secret)
	}

	checkToken := func() string {
		timeout := time.Now().Add(10 * time.Second)
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

				_, err = os.Stat(secret)
				switch {
				case removeSecretIDFile && err == nil:
					t.Fatal("secret file exists but was supposed to be removed")
				case !removeSecretIDFile && err != nil:
					t.Fatal("secret ID file does not exist but was not supposed to be removed")
				}

				client.SetToken(string(val))
				secret, err := client.Auth().Token().LookupSelf()
				if err != nil {
					t.Fatal(err)
				}
				return secret.Data["entity_id"].(string)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
	origEntity := checkToken()

	// Make sure it gets renewed
	timeout := time.Now().Add(4 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		secret, err := client.Auth().Token().LookupSelf()
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

	// Write new values
	if err := ioutil.WriteFile(role, []byte(roleID2), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test role 2", "path", role)
	}

	if err := ioutil.WriteFile(secret, []byte(secretID2), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test secret 2", "path", secret)
	}

	newEntity := checkToken()
	if newEntity == origEntity {
		t.Fatal("found same entity")
	}

	timeout = time.Now().Add(4 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		secret, err := client.Auth().Token().LookupSelf()
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

func TestAppRoleWithWrapping(t *testing.T) {
	var err error
	logger := logging.NewVaultLogger(log.Trace)
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client
	origToken := client.Token()

	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test1", map[string]interface{}{
		"bind_secret_id": "true",
		"token_ttl":      "3s",
		"token_max_ttl":  "10s",
	})
	if err != nil {
		t.Fatal(err)
	}

	client.SetWrappingLookupFunc(func(operation, path string) string {
		if path == "auth/approle/role/test1/secret-id" {
			return "10s"
		}
		return ""
	})

	resp, err := client.Logical().Write("auth/approle/role/test1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID1 := resp.WrapInfo.Token

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

	// We close these right away because we're just basically testing
	// permissions and finding a usable file name
	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	out := ouf.Name()
	ouf.Close()
	os.Remove(out)
	t.Logf("output: %s", out)

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	conf := map[string]interface{}{
		"role_id_file_path":                   role,
		"secret_id_file_path":                 secret,
		"secret_id_response_wrapping_path":    "auth/approle/role/test1/secret-id",
		"remove_secret_id_file_after_reading": true,
	}

	am, err := agentapprole.NewApproleAuthMethod(&auth.AuthConfig{
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
	}
	ah := auth.NewAuthHandler(ahConfig)
	go ah.Run(ctx, am)
	defer func() {
		<-ah.DoneCh
	}()

	config := &sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": out,
		},
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

	// Check that no sink file exists
	_, err = os.Lstat(out)
	if err == nil {
		t.Fatal("expected err")
	}
	if !os.IsNotExist(err) {
		t.Fatal("expected notexist err")
	}

	if err := ioutil.WriteFile(role, []byte(roleID1), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test role 1", "path", role)
	}

	if err := ioutil.WriteFile(secret, []byte(secretID1), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test secret 1", "path", secret)
	}

	checkToken := func() string {
		timeout := time.Now().Add(10 * time.Second)
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

				if _, err := os.Stat(secret); err == nil {
					t.Fatal("secret ID file does not exist but was not supposed to be removed")
				}

				client.SetToken(string(val))
				secret, err := client.Auth().Token().LookupSelf()
				if err != nil {
					t.Fatal(err)
				}
				return secret.Data["entity_id"].(string)
			}
			time.Sleep(250 * time.Millisecond)
		}
	}
	origEntity := checkToken()

	// Make sure it gets renewed
	timeout := time.Now().Add(4 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		secret, err := client.Auth().Token().LookupSelf()
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

	// Write new values
	client.SetToken(origToken)
	resp, err = client.Logical().Write("auth/approle/role/test1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID2 := resp.WrapInfo.Token
	if err := ioutil.WriteFile(secret, []byte(secretID2), 0600); err != nil {
		t.Fatal(err)
	} else {
		logger.Trace("wrote test secret 2", "path", secret)
	}

	newEntity := checkToken()
	if newEntity != origEntity {
		t.Fatal("did not find same entity")
	}

	timeout = time.Now().Add(4 * time.Second)
	for {
		if time.Now().After(timeout) {
			break
		}
		secret, err := client.Auth().Token().LookupSelf()
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
