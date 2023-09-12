// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	"github.com/hashicorp/go-retryablehttp"

	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

func fakeVaultServer(t *testing.T) *httptest.Server {
	t.Helper()

	firstRequest := true

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kv/my-app/creds", func(w http.ResponseWriter, r *http.Request) {
		// change the password on the second request to re-render the template
		var password string

		if firstRequest {
			password = "s3cr3t"
		} else {
			password = "s3cr3t-two"
		}

		firstRequest = false

		fmt.Fprintf(w, `{
                "request_id": "8af096e9-518c-7351-eff5-5ba20554b21f",
                "lease_id": "",
                "renewable": false,
                "lease_duration": 0,
                "data": {
                    "data": {
                        "password": "%s",
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
            }`,
			password,
		)
	})

	return httptest.NewServer(mux)
}

// TestExecServer_Run tests various scenarios of using vault agent as a process
// supervisor. At its core is a sample application referred to as 'test app',
// compiled from ./test-app/main.go. Each test case verifies that the test app
// is started and/or stopped correctly by exec.Server.Run(). There are 3
// high-level scenarios we want to test for:
//
//  1. test app is started and is injected with environment variables
//  2. test app exits early (either with zero or non-zero extit code)
//  3. test app needs to be stopped (and restarted) by exec.Server
func TestExecServer_Run(t *testing.T) {
	// we must build a test-app binary since 'go run' does not propagate signals correctly
	goBinary, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("could not find go binary on path: %s", err)
	}

	testAppBinary := filepath.Join(os.TempDir(), "test-app")

	if err := exec.Command(goBinary, "build", "-o", testAppBinary, "./test-app").Run(); err != nil {
		t.Fatalf("could not build the test application: %s", err)
	}
	defer func() {
		if err := os.Remove(testAppBinary); err != nil {
			t.Fatalf("could not remove %q test application: %s", testAppBinary, err)
		}
	}()

	testCases := map[string]struct {
		// skip this test case
		skip       bool
		skipReason string

		// inputs to the exec server
		envTemplates               []*ctconfig.TemplateConfig
		staticSecretRenderInterval time.Duration

		// test app parameters
		testAppArgs       []string
		testAppStopSignal os.Signal
		testAppPort       int

		// simulate a shutdown of agent, which, in turn stops the test app
		simulateShutdown             bool
		simulateShutdownWaitDuration time.Duration

		// expected results
		expected             map[string]string
		expectedTestDuration time.Duration
		expectedError        error
	}{
		"ensure_environment_variables_are_injected": {
			skip: true,
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

		"password_changes_test_app_should_restart": {
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}, {
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.password }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_PASSWORD"),
			}},
			staticSecretRenderInterval: 5 * time.Second,
			testAppArgs:                []string{"--stop-after", "15s", "--sleep-after-stop-signal", "0s"},
			testAppStopSignal:          syscall.SIGTERM,
			testAppPort:                34002,
			expected: map[string]string{
				"MY_USER":     "app-user",
				"MY_PASSWORD": "s3cr3t-two",
			},
			expectedTestDuration: 15 * time.Second,
			expectedError:        nil,
		},

		"test_app_exits_early": {
			skip: true,
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:          []string{"--stop-after", "1s"},
			testAppStopSignal:    syscall.SIGTERM,
			testAppPort:          34003,
			expectedTestDuration: 15 * time.Second,
			expectedError:        &ProcessExitError{0},
		},

		"test_app_exits_early_non_zero": {
			skip: true,
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:          []string{"--stop-after", "1s", "--exit-code", "5"},
			testAppStopSignal:    syscall.SIGTERM,
			testAppPort:          34004,
			expectedTestDuration: 15 * time.Second,
			expectedError:        &ProcessExitError{5},
		},

		"send_sigterm_expect_test_app_exit": {
			skip: true,
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:                  []string{"--stop-after", "30s", "--sleep-after-stop-signal", "1s"},
			testAppStopSignal:            syscall.SIGTERM,
			testAppPort:                  34005,
			simulateShutdown:             true,
			simulateShutdownWaitDuration: 3 * time.Second,
			expectedTestDuration:         15 * time.Second,
			expectedError:                nil,
		},

		"send_sigusr1_expect_test_app_exit": {
			skip: true,
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:                  []string{"--stop-after", "30s", "--sleep-after-stop-signal", "1s", "--use-sigusr1"},
			testAppStopSignal:            syscall.SIGUSR1,
			testAppPort:                  34006,
			simulateShutdown:             true,
			simulateShutdownWaitDuration: 3 * time.Second,
			expectedTestDuration:         15 * time.Second,
			expectedError:                nil,
		},

		"test_app_ignores_stop_signal": {
			skip:       true,
			skipReason: "This test currently fails with 'go test -race' (see hashicorp/consul-template/issues/1753).",
			envTemplates: []*ctconfig.TemplateConfig{{
				Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
				MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
			}},
			testAppArgs:                  []string{"--stop-after", "60s", "--sleep-after-stop-signal", "60s"},
			testAppStopSignal:            syscall.SIGTERM,
			testAppPort:                  34007,
			simulateShutdown:             true,
			simulateShutdownWaitDuration: 32 * time.Second, // the test app should be stopped immediately after 30s
			expectedTestDuration:         45 * time.Second,
			expectedError:                nil,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if testCase.skip {
				t.Skip(testCase.skipReason)
			}

			t.Logf("test case %s: begin", name)
			defer t.Logf("test case %s: end", name)

			fakeVault := fakeVaultServer(t)
			defer fakeVault.Close()

			ctx, cancelContextFunc := context.WithTimeout(context.Background(), testCase.expectedTestDuration)
			defer cancelContextFunc()

			testAppCommand := []string{
				testAppBinary,
				"--port",
				strconv.Itoa(testCase.testAppPort),
			}

			execServer, err := NewServer(&ServerConfig{
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
					TemplateConfig: &config.TemplateConfig{
						ExitOnRetryFailure:    true,
						StaticSecretRenderInt: testCase.staticSecretRenderInterval,
					},
				},
				LogLevel:  hclog.Trace,
				LogWriter: hclog.DefaultOutput,
			})
			if err != nil {
				t.Fatalf("could not create exec server: %q", err)
			}

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

			// ensure the test app is running after 3 seconds
			var (
				testAppAddr      = fmt.Sprintf("http://localhost:%d", testCase.testAppPort)
				testAppStartedCh = make(chan error)
			)
			if testCase.expectedError == nil {
				time.AfterFunc(500*time.Millisecond, func() {
					_, err := retryablehttp.Head(testAppAddr)
					testAppStartedCh <- err
				})
			}

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

			case err := <-testAppStartedCh:
				if testCase.expectedError == nil && err != nil {
					t.Fatalf("test app could not be started")
				}

				t.Log("test app started successfully")
			}

			// expect the test app to restart after staticSecretRenderInterval + debounce timer due to a password change
			if testCase.staticSecretRenderInterval != 0 {
				t.Logf("sleeping for %v to wait for application restart", testCase.staticSecretRenderInterval+5*time.Second)
				time.Sleep(testCase.staticSecretRenderInterval + 5*time.Second)
			}

			// simulate a shutdown of agent, which, in turn stops the test app
			if testCase.simulateShutdown {
				cancelContextFunc()

				time.Sleep(testCase.simulateShutdownWaitDuration)

				// check if the test app is still alive
				if _, err := http.Head(testAppAddr); err == nil {
					t.Fatalf("the test app is still alive %v after a simulated shutdown!", testCase.simulateShutdownWaitDuration)
				}

				return
			}

			// verify the environment variables
			t.Logf("verifying test-app's environment variables")

			resp, err := retryablehttp.Get(testAppAddr)
			if err != nil {
				t.Fatalf("error making request to the test app: %s", err)
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

func TestExecServer_LogFiles(t *testing.T) {
	goBinary, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("could not find go binary on path: %s", err)
	}

	testAppBinary := filepath.Join(os.TempDir(), "test-app")

	if err := exec.Command(goBinary, "build", "-o", testAppBinary, "./test-app").Run(); err != nil {
		t.Fatalf("could not build the test application: %s", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(testAppBinary); err != nil {
			t.Fatalf("could not remove %q test application: %s", testAppBinary, err)
		}
	})

	tempStderr := filepath.Join(os.TempDir(), "vault-exec-test.stderr.log")
	t.Cleanup(func() {
		_ = os.Remove(tempStderr)
	})

	testCases := map[string]struct {
		testAppPort   int
		stderrFile    string
		expectedError error
	}{
		"can_log_to_file": {
			testAppPort: 34001,
			stderrFile:  tempStderr,
		},
		"cant_open_file": {
			testAppPort:   34002,
			stderrFile:    "/file/does/not/exist",
			expectedError: os.ErrNotExist,
		},
	}

	for tcName, testCase := range testCases {
		t.Run(tcName, func(t *testing.T) {
			fakeVault := fakeVaultServer(t)
			defer fakeVault.Close()

			testAppCommand := []string{
				testAppBinary,
				"--port",
				strconv.Itoa(testCase.testAppPort),
				"--stop-after",
				"60s",
			}

			execServer, err := NewServer(&ServerConfig{
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
						Command:                testAppCommand,
						ChildProcessStderr:     testCase.stderrFile,
					},
					EnvTemplates: []*ctconfig.TemplateConfig{{
						Contents:                 pointerutil.StringPtr(`{{ with secret "kv/my-app/creds" }}{{ .Data.data.user }}{{ end }}`),
						MapToEnvironmentVariable: pointerutil.StringPtr("MY_USER"),
					}},
					TemplateConfig: &config.TemplateConfig{
						ExitOnRetryFailure:    true,
						StaticSecretRenderInt: 5 * time.Second,
					},
				},
				LogLevel:  hclog.Trace,
				LogWriter: hclog.DefaultOutput,
			})
			if err != nil {
				if testCase.expectedError != nil {
					if errors.Is(err, testCase.expectedError) {
						t.Log("test passes! caught expected err")
						return
					} else {
						t.Fatalf("caught error %q did not match expected error %q", err, testCase.expectedError)
					}
				}
				t.Fatalf("could not create exec server: %q", err)
			}

			// replace the tempfile created with one owned by this test
			var stdoutBuffer bytes.Buffer
			execServer.childProcessStderr = noopCloser{&stdoutBuffer}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

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

			// ensure the test app is running after 500ms
			var (
				testAppAddr      = fmt.Sprintf("http://localhost:%d", testCase.testAppPort)
				testAppStartedCh = make(chan error)
			)
			time.AfterFunc(500*time.Millisecond, func() {
				_, err := retryablehttp.Head(testAppAddr)
				testAppStartedCh <- err
			})

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

			case <-testAppStartedCh:
				t.Log("test app started successfully")
			}

			// let the app run a bit
			time.Sleep(5 * time.Second)
			// stop the app
			cancel()
			// wait for app to stop
			time.Sleep(5 * time.Second)

			if stdoutBuffer.Len() == 0 {
				t.Fatalf("stdout log file does not have any data!")
			}
		})
	}
}

type noopCloser struct {
	io.Writer
}

func (_ noopCloser) Close() error {
	return nil
}
