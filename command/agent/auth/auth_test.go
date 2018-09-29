package auth

import (
	"context"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/logging"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

type userpassTestMethod struct{}

func newUserpassTestMethod(t *testing.T, client *api.Client) AuthMethod {
	err := client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
		Config: api.AuthConfigInput{
			DefaultLeaseTTL: "1s",
			MaxLeaseTTL:     "3s",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	return &userpassTestMethod{}
}

func (u *userpassTestMethod) Authenticate(_ context.Context, client *api.Client) (string, map[string]interface{}, error) {
	_, err := client.Logical().Write("auth/userpass/users/foo", map[string]interface{}{
		"password": "bar",
	})
	if err != nil {
		return "", nil, err
	}
	return "auth/userpass/login/foo", map[string]interface{}{
		"password": "bar",
	}, nil
}

func (u *userpassTestMethod) NewCreds() chan struct{} {
	return nil
}

func (u *userpassTestMethod) CredSuccess() {
}

func (u *userpassTestMethod) Shutdown() {
}

func TestAuthHandler(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger,
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	ctx, cancelFunc := context.WithCancel(context.Background())

	ah := NewAuthHandler(&AuthHandlerConfig{
		Logger: logger.Named("auth.handler"),
		Client: client,
	})

	am := newUserpassTestMethod(t, client)
	go ah.Run(ctx, am)

	// Consume tokens so we don't block
	stopTime := time.Now().Add(5 * time.Second)
	closed := false
consumption:
	for {
		select {
		case <-ah.OutputCh:
			// Nothing
		case <-time.After(stopTime.Sub(time.Now())):
			if !closed {
				cancelFunc()
				closed = true
			}
		case <-ah.DoneCh:
			break consumption
		}
	}
}
