package agent

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	log "github.com/hashicorp/go-hclog"
	credPCF "github.com/hashicorp/vault-plugin-auth-pcf"
	"github.com/hashicorp/vault-plugin-auth-pcf/testing/certificates"
	pcfAPI "github.com/hashicorp/vault-plugin-auth-pcf/testing/pcf"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
	agentpcf "github.com/hashicorp/vault/command/agent/auth/pcf"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/command/agent/sink/file"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestPCFEndToEnd(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"pcf": credPCF.Factory,
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
	if err := client.Sys().EnableAuthWithOptions("pcf", &api.EnableAuthOptions{
		Type: "pcf",
	}); err != nil {
		t.Fatal(err)
	}

	testIPAddress := "127.0.0.1"

	// Generate some valid certs that look like the ones we get from PCF.
	testPCFCerts, err := certificates.Generate(pcfAPI.FoundServiceGUID, pcfAPI.FoundOrgGUID, pcfAPI.FoundSpaceGUID, pcfAPI.FoundAppGUID, testIPAddress)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := testPCFCerts.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	// Start a mock server representing their API.
	mockPCFAPI := pcfAPI.MockServer(false)
	defer mockPCFAPI.Close()

	// Configure a CA certificate like a Vault operator would in setting up PCF.
	if _, err := client.Logical().Write("auth/pcf/config", map[string]interface{}{
		"identity_ca_certificates": testPCFCerts.CACertificate,
		"pcf_api_addr":             mockPCFAPI.URL,
		"pcf_username":             pcfAPI.AuthUsername,
		"pcf_password":             pcfAPI.AuthPassword,
	}); err != nil {
		t.Fatal(err)
	}

	// Configure a role to be used for logging in, another thing a Vault operator would do.
	if _, err := client.Logical().Write("auth/pcf/roles/test-role", map[string]interface{}{
		"bound_instance_ids":     pcfAPI.FoundServiceGUID,
		"bound_organization_ids": pcfAPI.FoundOrgGUID,
		"bound_space_ids":        pcfAPI.FoundSpaceGUID,
		"bound_application_ids":  pcfAPI.FoundAppGUID,
	}); err != nil {
		t.Fatal(err)
	}

	os.Setenv(credPCF.EnvVarInstanceCertificate, testPCFCerts.PathToInstanceCertificate)
	os.Setenv(credPCF.EnvVarInstanceKey, testPCFCerts.PathToInstanceKey)

	ctx, cancelFunc := context.WithCancel(context.Background())
	timer := time.AfterFunc(30*time.Second, func() {
		cancelFunc()
	})
	defer timer.Stop()

	am, err := agentpcf.NewPCFAuthMethod(&auth.AuthConfig{
		MountPath: "auth/pcf",
		Config: map[string]interface{}{
			"role": "test-role",
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
