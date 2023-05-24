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
	"strconv"
	"syscall"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

func dummyVaultServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kv/my-app/creds", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
                "request_id": "8af096e9-518c-7351-eff5-5ba20554b21f",
                "lease_id": "",
                "renewable": false,
                "lease_duration": 0,
                "data": {
                    "data": {
                        "password": "s3cr3t",
                        "user": "app-user"
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
            }`)
	})

	return httptest.NewServer(mux)
}

func processErrorCodeChecker(expectedExitCode int) func(t *testing.T, err error) {
	return func(t *testing.T, err error) {
		var processExitError *ProcessExitError
		if errors.As(err, &processExitError) {
			if processExitError.ExitCode != expectedExitCode {
				t.Fatalf("expected exit code %d, got %d", expectedExitCode, processExitError.ExitCode)
			}
		} else {
			t.Fatalf("expected error of type ProcessExitError")
		}
	}
}

func TestServer_Run(t *testing.T) {
	vault := dummyVaultServer()
	defer vault.Close()

	testCases := map[string]struct {
		envTemplates      []*ctconfig.TemplateConfig
		checkError        func(*testing.T, error)
		testAppArgs       []string
		testAppStopSignal os.Signal
		testAppPort       int
		expected          map[string]string
		expectedExit      bool
	}{
		"simple": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}, {
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.password }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_PASSWORD"),
			}},
			testAppArgs:       []string{"--stop-after", "10s"},
			testAppStopSignal: syscall.SIGTERM,
			testAppPort:       34001,
			expected: map[string]string{
				"MY_USER":     "app-user",
				"MY_PASSWORD": "s3cr3t",
			},
			expectedExit: false,
		},
		"exits_early_success": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:       []string{"--stop-after", "2s"},
			testAppStopSignal: syscall.SIGTERM,
			testAppPort:       34002,
			expected: map[string]string{
				"MY_USER": "app-user",
			},
			expectedExit: true,
			checkError:   processErrorCodeChecker(0),
		},
		"exits_early_non_zero": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:       []string{"--stop-after", "2s", "--exit-code", "5"},
			testAppStopSignal: syscall.SIGTERM,
			testAppPort:       34003,
			expected: map[string]string{
				"MY_USER": "app-user",
			},
			expectedExit: true,
			checkError:   processErrorCodeChecker(5),
		},
	}

	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("could not find go binary on path: %s", err)
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			exampleAppUrl := fmt.Sprintf("http://localhost:%d", testCase.testAppPort)
			baseCmdArgs := []string{
				goBin,
				"run",
				"./test-app",
				"--port",
				strconv.Itoa(testCase.testAppPort),
			}

			execServer := NewServer(&ServerConfig{
				Logger: logging.NewVaultLogger(hclog.Trace),
				AgentConfig: &config.Config{
					Vault: &config.Vault{
						Address: vault.URL,
						Retry: &config.Retry{
							NumRetries: 3,
						},
					},
					Exec: &config.ExecConfig{
						RestartOnSecretChanges: "always",
						Command:                append(baseCmdArgs, testCase.testAppArgs...),
						RestartStopSignal:      testCase.testAppStopSignal,
					},
					EnvTemplates: testCase.envTemplates,
				},
				LogLevel:  hclog.Trace,
				LogWriter: hclog.DefaultOutput,
			})

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
				errCh <- execServer.Run(ctx, templateTokenCh)
			}()

			// send a dummy value to kick off the server
			templateTokenCh <- "my-token"

			// check to make sure the app is running
			if !testCase.expectedExit {
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

			time.Sleep(5 * time.Second)

			select {
			case <-ctx.Done():
				t.Fatal("timeout reached before templates were rendered")
			case err := <-errCh:
				if err != nil && !testCase.expectedExit {
					t.Fatalf("did not expect error, got: %v", err)
				}
				if err != nil && testCase.expectedExit {
					t.Logf("received expected error: %v", err)
					if testCase.checkError != nil {
						testCase.checkError(t, err)
					}
					return
				}
			case <-appStartedCh:
				t.Log("app has started")
			case <-appStartErrCh:
				if !testCase.expectedExit {
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

			for key, expectedValue := range testCase.expected {
				actualValue, ok := response.EnvVars[key]
				if !ok {
					t.Fatalf("expected the test app to return %q environment variable", key)
				}
				if expectedValue != actualValue {
					t.Fatalf("expected environment variable %s to have a value of %q but it has a value of %q", key, expectedValue, actualValue)
				}
			}

			// explicitly cancel
			cancel()
		})
	}
}
