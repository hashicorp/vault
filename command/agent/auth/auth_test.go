package auth

import (
	"context"
	"net"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
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

// fakeTestMethod is meant to fake the minimum amount needed to get to the point
// where testing connection failed attempts is possible
type fakeTestMethod struct {
	// embed a userpassTestMethod to inherit NewCreds, CredSuccess, and Shutdown
	// methods
	*userpassTestMethod
}

func newFakeTestMethod(t *testing.T, client *api.Client) AuthMethod {
	return &fakeTestMethod{}
}

func (u *fakeTestMethod) Authenticate(_ context.Context, client *api.Client) (string, map[string]interface{}, error) {
	// simply need to return something here
	return "auth/userpass/login/foo", map[string]interface{}{
		"password": "bar",
	}, nil
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
		case <-ah.TemplateTokenCh:
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

/// TestConnectionAttempts verifies the max-connection-attempts flag is fully
//supported
func TestConnectionAttempts(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)

	// start a network listener, intentionally don't listen for anything. This
	// simulates an unresponsive Vault cluster
	testCases := map[string]int{
		"unset": 0,
		"low":   1,
		"high":  15,
	}

	for name, maxAttempts := range testCases {
		t.Run(name, func(t *testing.T) {
			l, err := net.Listen("tcp", ":1234")
			if err != nil {
				t.Fatal(err)
			}

			// Use local host with the same port as the net listener. We can't use
			// l.Addr().String() because it's not an http address, and then the listener
			// will never get connection attempts from Vault
			client, err := api.NewClient(nil)
			if err != nil {
				t.Fatal(err)
			}
			client.SetAddress("http://127.0.0.1:1234")
			client.SetMaxRetries(0)

			ctx, cancelFunc := context.WithCancel(context.Background())

			// attempts tracks the connection attempts
			var attempts int
			go func(ctx context.Context) {
				for {
					conn, _ := l.Accept()
					select {
					case <-ctx.Done():
						return
					default:
					}

					conn.Close()
					attempts++
				}
			}(ctx)

			ah := NewAuthHandler(&AuthHandlerConfig{
				Logger: logger.Named("auth.handler"),
				Client: client,
			})

			am := newFakeTestMethod(t, client)
			go ah.RunWithMaxAttempts(ctx, am, maxAttempts)

			// all the tests should complete well within this time. We use a timer to
			// ensure that if something does break, we don't hang forever
			timeout := maxAttempts * 5
			if timeout == 0 {
				timeout = 5
			}
			stopTime := time.Now().Add(time.Duration(timeout) * time.Second)

			var closed bool
			// Consume tokens so we don't block
		consumption:
			for {
				select {
				case <-ah.OutputCh:
				case <-ah.TemplateTokenCh:
				// Nothing
				case <-time.After(stopTime.Sub(time.Now())):
					// if test max is not zero and we've reached the timeout, something is
					// wrong
					if !closed {
						cancelFunc()
						closed = true
					}
					if maxAttempts != 0 {
						t.Fatalf("test timeout. Expected: %d, actual: %d", maxAttempts, attempts)
					}
					break consumption
				case <-ah.DoneCh:
					break consumption
				}
			}

			if !closed {
				cancelFunc()
				closed = true
			}

			if maxAttempts != 0 {
				if attempts != maxAttempts {
					t.Fatalf("connection attempts did not match expected. Expected: %d, actual: %d", maxAttempts, attempts)
				}
			}

			l.Close()
		})
	}
}
