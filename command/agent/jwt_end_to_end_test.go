package agent

import (
	"io/ioutil"
	"testing"

	hclog "github.com/hashicorp/go-hclog"
	vaultjwt "github.com/hashicorp/vault-plugin-auth-jwt"
	"github.com/hashicorp/vault/command/agent/auth"
	"github.com/hashicorp/vault/command/agent/auth/jwt"
	"github.com/hashicorp/vault/command/agent/sink"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestJWTEndtoEnd(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger,
		CredentialBackends: map[string]logical.Factory{
			"jwt": vaultjwt.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	inf, err := ioutil.TempFile("", "auth.jwt.test.")
	if err != nil {
		t.Fatal(err)
	}
	defer inf.Close()
	ouf, err := ioutil.TempFile("", "auth.tokensink.test.")
	if err != nil {
		t.Fatal(err)
	}
	defer ouf.Close()

	am, err := jwt.NewJWTAuthMethod(&auth.AuthConfig{
		Logger:    logger.Named("auth.jwt"),
		MountPath: "auth/jwt",
		Config: map[string]interface{}{
			"path": inf.Name(),
			"role": "test",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ah := auth.NewAuthHandler(&auth.AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	})
	go ah.Run(am)

	fs, err := sink.NewFileSink(&sink.SinkConfig{
		Logger: logger.Named("sink.file"),
		Config: map[string]interface{}{
			"path": ouf.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ss := sink.NewSinkServer(&sink.SinkConfig{
		Logger: logger.Named("sink.server"),
		Client: client,
	})
	go ss.Run(ah.OutputCh, []sink.Sink{fs})

	close(ah.ShutdownCh)
	<-ah.DoneCh
	close(ss.ShutdownCh)
	<-ss.DoneCh
}
