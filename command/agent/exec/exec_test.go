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
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

func fakeVaultServer() *httptest.Server {
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

// TestServer_Run tests various scenarios of using vault agent as a process supervisor
func TestServer_Run(t *testing.T) {
	fakeVault := fakeVaultServer()
	defer fakeVault.Close()

	testCases := map[string]struct {
		envTemplates                  []*ctconfig.TemplateConfig
		testAppArgs                   []string
		testAppStopSignal             os.Signal
		testAppPort                   int
		simulatedShutdownWaitDuration time.Duration
		expected                      map[string]string
		expectedTestDuration          time.Duration
		expectedError                 error
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
			expectedTestDuration: 15 * time.Second,
			expectedError:        nil,
		},

		"test_app_exits_early": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:          []string{"--stop-after", "1s"},
			testAppStopSignal:    syscall.SIGTERM,
			testAppPort:          34002,
			expectedTestDuration: 15 * time.Second,
			expectedError:        &ProcessExitError{0},
		},

		"test_app_exits_early_non_zero": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:          []string{"--stop-after", "1s", "--exit-code", "5"},
			testAppStopSignal:    syscall.SIGTERM,
			testAppPort:          34003,
			expectedTestDuration: 15 * time.Second,
			expectedError:        &ProcessExitError{5},
		},

		"send_sigterm_expect_test_app_exit": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:                   []string{"--stop-after", "30s", "--sleep-after-stop-signal", "1s"},
			testAppStopSignal:             syscall.SIGTERM,
			testAppPort:                   34004,
			simulatedShutdownWaitDuration: 3 * time.Second,
			expectedTestDuration:          15 * time.Second,
		},

		"send_sigusr1_expect_test_app_exit": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:                   []string{"--stop-after", "30s", "--sleep-after-stop-signal", "1s", "--use-sigusr1"},
			testAppStopSignal:             syscall.SIGUSR1,
			testAppPort:                   34005,
			simulatedShutdownWaitDuration: 3 * time.Second,
			expectedTestDuration:          15 * time.Second,
		},

		"test_app_ignores_stop_signal": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:                   []string{"--stop-after", "60s", "--ignore-stop-signal"},
			testAppStopSignal:             syscall.SIGTERM,
			testAppPort:                   34006,
			simulatedShutdownWaitDuration: 32 * time.Second, // the test app should be stopped immediately after 30s
			expectedTestDuration:          45 * time.Second,
		},
	}

	goBinary, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("could not find go binary on path: %s", err)
	}

	// we must build a test binary since 'go run' does not propagate signals correctly
	testAppBinary := filepath.Join(os.TempDir(), "test-app")

	if err := exec.Command(goBinary, "build", "-o", testAppBinary, "./test-app").Run(); err != nil {
		t.Fatalf("could not build the test application: %s", err)
	}
	defer func() {
		if err := os.Remove(testAppBinary); err != nil {
			t.Fatalf("could not remove %q test application: %s", testAppBinary, err)
		}
	}()

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, cancelContextFunc := context.WithTimeout(context.Background(), testCase.expectedTestDuration)
			defer cancelContextFunc()

			testAppCommand := []string{
				testAppBinary,
				"--port",
				strconv.Itoa(testCase.testAppPort),
			}

			execServer := NewServer(&ServerConfig{
				Logger: logging.NewVaultLogger(hclog.Trace),
				AgentConfig: &config.Config{
					Vault: &config.Vault{
						Address: fakeVault.URL,
						Retry: &config.Retry{
							NumRetries: 3,
						},
					},
					Exec: &config.ExecConfig{
						RestartOnSecretChanges: "always",
						Command:                append(testAppCommand, testCase.testAppArgs...),
						RestartStopSignal:      testCase.testAppStopSignal,
					},
					EnvTemplates: testCase.envTemplates,
				},
				LogLevel:  hclog.Trace,
				LogWriter: hclog.DefaultOutput,
			})

			// start the exec server
			var (
				execServerErrCh   = make(chan error)
				execServerTokenCh = make(chan string, 1)
			)
			go func() {
				execServerErrCh <- execServer.Run(ctx, execServerTokenCh)
			}()

			// send a dummy token to kick off the server
			execServerTokenCh <- "my-token"

			testAppHealthCheckCh := make(chan error)
			testAppAddress := fmt.Sprintf("http://localhost:%d", testCase.testAppPort)

			// check to make sure the test app is running
			if testCase.expectedError == nil {
				go func() {
					time.Sleep(3 * time.Second)
					_, err := http.Head(testAppAddress)
					testAppHealthCheckCh <- err
				}()
			}

			time.Sleep(5 * time.Second)

			select {
			case <-ctx.Done():
				t.Fatal("timeout reached before templates were rendered")

			case err := <-execServerErrCh:
				if testCase.expectedError == nil && err != nil {
					t.Fatalf("exec server did not expect an error, got: %v", err)
				}

				if errors.Is(err, testCase.expectedError) {
					t.Fatalf("exec server expected error %v; got %v", testCase.expectedError, err)
				}

				t.Log("exec server exited without an error")
				return

			case err := <-testAppHealthCheckCh:
				if testCase.expectedError == nil && err != nil {
					t.Fatalf("test app could not be started")
				}

				t.Log("test app started successfully")
			}

			// simulate a shutdown of agent, which, in turn stops the test app
			if testCase.simulatedShutdownWaitDuration > 0 {
				cancelContextFunc()

				time.Sleep(testCase.simulatedShutdownWaitDuration)

				// check if the test app is still alive
				if _, err := http.Head(testAppAddress); err == nil {
					t.Fatalf("the test app is still alive %v after a simulated shutdown!", testCase.simulatedShutdownWaitDuration)
				}

				return
			}

			// verify the environment variables
			resp, err := http.Get(testAppAddress)
			if err != nil {
				t.Fatalf("error making request to test app: %s", err)
			}
			defer resp.Body.Close()

			decoder := json.NewDecoder(resp.Body)
			var response struct {
				EnvironmentVariables map[string]string `json:"environment_variables"`
				ProcessID            int               `json:"process_id"`
			}
			if err := decoder.Decode(&response); err != nil {
				t.Fatalf("unable to parse response from test app: %s", err)
			}

			for key, expectedValue := range testCase.expected {
				actualValue, ok := response.EnvironmentVariables[key]
				if !ok {
					t.Fatalf("expected the test app to return %q environment variable", key)
				}
				if expectedValue != actualValue {
					t.Fatalf("expected environment variable %s to have a value of %q but it has a value of %q", key, expectedValue, actualValue)
				}
			}
		})
	}
}
