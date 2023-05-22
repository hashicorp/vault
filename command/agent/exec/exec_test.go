package exec

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/config"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	"github.com/hashicorp/vault/vault"
)

const (
	exampleAppUrl = "http://localhost:8000"
)

func testVaultServer(t *testing.T) (*api.Client, func()) {
	t.Helper()

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       hclog.NewNullLogger(),
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	// enable kv-v2 backend
	if err := client.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)

	return client, cluster.Cleanup
}

func createHttpTestServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kv/myapp/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, jsonResponse)
	})

	return httptest.NewServer(mux)
}

func processErrorCodeChecker(expectedExitCode int) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		var processExitError *ProcessExitError
		if errors.As(err, &processExitError) {
			if processExitError.ExitCode != expectedExitCode {
				t.Fatalf("expected there to be an exit code of %d, got %d", expectedExitCode, processExitError.ExitCode)
			}
		} else {
			t.Fatalf("expected error of type ProcessExitError")
		}
	}
}

func TestServer_Run(t *testing.T) {
	testServer := createHttpTestServer()
	defer testServer.Close()

	testCases := map[string]struct {
		envTemplates   []*ctconfig.TemplateConfig
		expectedValues map[string]string
		extraAppArgs   []string
		expectError    bool
		checkError     func(*testing.T, error)
		processTime    time.Duration
		stopSignal     os.Signal
	}{
		"simple": {
			envTemplates: []*ctconfig.TemplateConfig{
				{
					Contents:                 pointerutil.StringPtr(`{{ with secret "kv/myapp/config"}}{{.Data.data.username}}{{end}}`),
					MapToEnvironmentVariable: pointerutil.StringPtr("MY_USERNAME"),
				},
				{
					Contents:                 pointerutil.StringPtr(`{{ with secret "kv/myapp/config"}}{{.Data.data.password}}{{end}}`),
					MapToEnvironmentVariable: pointerutil.StringPtr("MY_PASSWORD"),
				},
			},
			expectedValues: map[string]string{
				"MY_USERNAME": "appuser",
				"MY_PASSWORD": "password",
			},
			expectError: false,
			stopSignal:  syscall.SIGTERM,
		},
		"exits_early": {
			envTemplates: []*ctconfig.TemplateConfig{
				{
					Contents:                 pointerutil.StringPtr(`{{ with secret "kv/myapp/config"}}{{.Data.data.username}}{{end}}`),
					MapToEnvironmentVariable: pointerutil.StringPtr("MY_USERNAME"),
				},
			},
			expectedValues: map[string]string{
				"MY_USERNAME": "appuser",
			},
			processTime:  time.Second * 5,
			extraAppArgs: []string{"--stop-after", "2s"},
			expectError:  true,
			stopSignal:   syscall.SIGTERM,
			checkError:   processErrorCodeChecker(0),
		},
	}

	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("could not find go binary on path: %s", err)
	}

	baseCmdArgs := []string{
		goBin,
		"run",
		"./test-app",
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			serverConfig := &ServerConfig{
				Logger: logging.NewVaultLogger(hclog.Trace),
				AgentConfig: &config.Config{
					Vault: &config.Vault{
						Address: testServer.URL,
						Retry: &config.Retry{
							NumRetries: 3,
						},
					},
					Exec: &config.ExecConfig{
						RestartOnSecretChanges: "always",
						Command:                append(baseCmdArgs, testCase.extraAppArgs...),
						RestartKillSignal:      testCase.stopSignal,
					},
					EnvTemplates: testCase.envTemplates,
				},
				LogLevel:  hclog.Trace,
				LogWriter: hclog.DefaultOutput,
			}

			server := NewServer(serverConfig)

			ctx, cancelTimeout := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancelTimeout()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			templateTokenCh := make(chan string, 1)
			errCh := make(chan error)
			appStartedCh := make(chan struct{})
			appStartErrCh := make(chan error)

			// start the exec server
			go func() {
				errCh <- server.Run(ctx, templateTokenCh)
			}()

			// send a dummy value to kick off the server
			templateTokenCh <- "test"

			// check to make sure the app is running
			if !testCase.expectError {
				go func() {
					time.Sleep(3 * time.Second)
					_, err := retryablehttp.Head(exampleAppUrl)
					if err != nil {
						appStartErrCh <- err
					} else {
						appStartedCh <- struct{}{}
					}
				}()
			}

			time.Sleep(testCase.processTime)

			select {
			case <-ctx.Done():
				t.Fatal("timeout reached before templates were rendered")
			case err := <-errCh:
				if err != nil && !testCase.expectError {
					t.Fatalf("did not expect error, got: %v", err)
				}
				if err != nil && testCase.expectError {
					t.Logf("received expected error: %v", err)
					if testCase.checkError != nil {
						testCase.checkError(t, err)
					}
					return
				}
			case <-appStartedCh:
				t.Log("app has started")
			case <-appStartErrCh:
				if !testCase.expectError {
					t.Fatal("app could not be started")
				}
			}

			// we started the server, now call it to see if
			// the environment variables are what they are supposed to be
			res, err := http.Get(exampleAppUrl)
			if err != nil {
				t.Fatalf("error making request to test app: %s", err)
			}
			defer res.Body.Close()

			decoder := json.NewDecoder(res.Body)
			var response struct {
				EnvVars   map[string]string `json:"env_vars"`
				ProcessID int               `json:"process_id"`
			}
			if err := decoder.Decode(&response); err != nil {
				t.Fatalf("unable to parse response from test app: %s", err)
			}

			for key, expectedValue := range testCase.expectedValues {
				actualValue, ok := response.EnvVars[key]
				if !ok {
					t.Fatalf("expected the test app to return %q env var", key)
				}
				if expectedValue != actualValue {
					t.Fatalf("expected env var %s to have a value of %q but it has a value of %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}

// copied from template_test.go
var jsonResponse = `
{
  "request_id": "8af096e9-518c-7351-eff5-5ba20554b21f",
  "lease_id": "",
  "renewable": false,
  "lease_duration": 0,
  "data": {
    "data": {
      "password": "password",
      "username": "appuser"
    },
    "metadata": {
      "created_time": "2019-10-07T22:18:44.233247Z",
      "deletion_time": "",
      "destroyed": false,
      "version": 3
    }
  },
  "wrap_info": null,
  "warnings": null,
  "auth": null
}
`
