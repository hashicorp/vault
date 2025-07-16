// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package template

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	sync "sync/atomic"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/command/agent/internal/ctmanager"
	"github.com/hashicorp/vault/command/agentproxyshared"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/internalshared/listenerutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/test/bufconn"
)

func newRunnerConfig(s *ServerConfig, configs ctconfig.TemplateConfigs) (*ctconfig.Config, error) {
	managerCfg := ctmanager.ManagerConfig{
		AgentConfig: s.AgentConfig,
	}
	cfg, err := ctmanager.NewConfig(managerCfg, configs)
	return cfg, err
}

// TestNewServer is a simple test to make sure NewServer returns a Server and
// channel
func TestNewServer(t *testing.T) {
	server := NewServer(&ServerConfig{})
	if server == nil {
		t.Fatal("nil server returned")
	}
}

func newAgentConfig(listeners []*configutil.Listener, enableCache, enablePersisentCache bool) *config.Config {
	agentConfig := &config.Config{
		SharedConfig: &configutil.SharedConfig{
			PidFile:   "./pidfile",
			Listeners: listeners,
		},
		AutoAuth: &config.AutoAuth{
			Method: &config.Method{
				Type:      "aws",
				MountPath: "auth/aws",
				Config: map[string]interface{}{
					"role": "foobar",
				},
			},
			Sinks: []*config.Sink{
				{
					Type:   "file",
					DHType: "curve25519",
					DHPath: "/tmp/file-foo-dhpath",
					AAD:    "foobar",
					Config: map[string]interface{}{
						"path": "/tmp/file-foo",
					},
				},
			},
		},
		Vault: &config.Vault{
			Address:          "http://127.0.0.1:1111",
			CACert:           "config_ca_cert",
			CAPath:           "config_ca_path",
			TLSSkipVerifyRaw: interface{}("true"),
			TLSSkipVerify:    true,
			ClientCert:       "config_client_cert",
			ClientKey:        "config_client_key",
		},
	}
	if enableCache {
		agentConfig.Cache = &config.Cache{
			UseAutoAuthToken: true,
		}
	}

	if enablePersisentCache {
		agentConfig.Cache.Persist = &agentproxyshared.PersistConfig{Type: "kubernetes"}
	}

	return agentConfig
}

func TestCacheConfig(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:       "tcp",
			Address:    "127.0.0.1:8300",
			TLSDisable: true,
		},
		{
			Type:        "unix",
			Address:     "foobar",
			TLSDisable:  true,
			SocketMode:  "configmode",
			SocketUser:  "configuser",
			SocketGroup: "configgroup",
		},
		{
			Type:        "tcp",
			Address:     "127.0.0.1:8400",
			TLSKeyFile:  "/path/to/cakey.pem",
			TLSCertFile: "/path/to/cacert.pem",
		},
	}

	cases := map[string]struct {
		cacheEnabled           bool
		persistentCacheEnabled bool
		setDialer              bool
		expectedErr            string
		expectCustomDialer     bool
	}{
		"persistent_cache": {
			cacheEnabled:           true,
			persistentCacheEnabled: true,
			setDialer:              true,
			expectedErr:            "",
			expectCustomDialer:     true,
		},
		"memory_cache": {
			cacheEnabled:           true,
			persistentCacheEnabled: false,
			setDialer:              true,
			expectedErr:            "",
			expectCustomDialer:     true,
		},
		"no_cache": {
			cacheEnabled:           false,
			persistentCacheEnabled: false,
			setDialer:              false,
			expectedErr:            "",
			expectCustomDialer:     false,
		},
		"cache_no_dialer": {
			cacheEnabled:           true,
			persistentCacheEnabled: false,
			setDialer:              false,
			expectedErr:            "missing in-process dialer configuration",
			expectCustomDialer:     false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			agentConfig := newAgentConfig(listeners, tc.cacheEnabled, tc.persistentCacheEnabled)
			if tc.setDialer && tc.cacheEnabled {
				bListener := bufconn.Listen(1024 * 1024)
				defer bListener.Close()
				agentConfig.Cache.InProcDialer = listenerutil.NewBufConnWrapper(bListener)
			}
			serverConfig := ServerConfig{AgentConfig: agentConfig}

			ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
			if len(tc.expectedErr) > 0 {
				require.Error(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, ctConfig)
			assert.Equal(t, tc.expectCustomDialer, ctConfig.Vault.Transport.CustomDialer != nil)

			if tc.expectCustomDialer {
				assert.Equal(t, "http://127.0.0.1:8200", *ctConfig.Vault.Address)
			} else {
				assert.Equal(t, "http://127.0.0.1:1111", *ctConfig.Vault.Address)
			}
		})
	}
}

func TestCacheConfigNoListener(t *testing.T) {
	listeners := []*configutil.Listener{}

	agentConfig := newAgentConfig(listeners, true, true)
	bListener := bufconn.Listen(1024 * 1024)
	defer bListener.Close()
	agentConfig.Cache.InProcDialer = listenerutil.NewBufConnWrapper(bListener)
	serverConfig := ServerConfig{AgentConfig: agentConfig}

	ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Equal(t, "http://127.0.0.1:8200", *ctConfig.Vault.Address)
	assert.NotNil(t, ctConfig.Vault.Transport.CustomDialer)
}

func createHttpTestServer() *httptest.Server {
	// create http test server
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/kv/myapp/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, jsonResponse)
	})
	mux.HandleFunc("/v1/kv/myapp/config-bad", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, `{"errors":[]}`)
	})
	mux.HandleFunc("/v1/kv/myapp/perm-denied", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		fmt.Fprintln(w, `{"errors":["1 error occurred:\n\t* permission denied\n\n"]}`)
	})

	return httptest.NewServer(mux)
}

func TestServerRun(t *testing.T) {
	ts := createHttpTestServer()
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "agent-tests")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// secretRender is a simple struct that represents the secret we render to
	// disk. It's used to unmarshal the file contents and test against
	type secretRender struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Version  string `json:"version"`
	}

	type templateTest struct {
		template *ctconfig.TemplateConfig
	}

	testCases := map[string]struct {
		templateMap        map[string]*templateTest
		expectedValues     *secretRender
		expectError        bool
		exitOnRetryFailure bool
	}{
		"simple": {
			templateMap: map[string]*templateTest{
				"render_01": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
			},
			expectError:        false,
			exitOnRetryFailure: false,
		},
		"multiple": {
			templateMap: map[string]*templateTest{
				"render_01": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_02": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_03": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_04": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_05": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_06": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_07": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
			},
			expectError:        false,
			exitOnRetryFailure: false,
		},
		"bad secret": {
			templateMap: map[string]*templateTest{
				"render_01": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContentsBad),
					},
				},
			},
			expectError:        true,
			exitOnRetryFailure: true,
		},
		"missing key": {
			templateMap: map[string]*templateTest{
				"render_01": {
					template: &ctconfig.TemplateConfig{
						Contents:      pointerutil.StringPtr(templateContentsMissingKey),
						ErrMissingKey: pointerutil.BoolPtr(true),
					},
				},
			},
			expectError:        true,
			exitOnRetryFailure: true,
		},
		"permission denied": {
			templateMap: map[string]*templateTest{
				"render_01": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContentsPermDenied),
					},
				},
			},
			expectError:        true,
			exitOnRetryFailure: true,
		},
		"with sprig functions": {
			templateMap: map[string]*templateTest{
				"render_01": {
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContentsWithSprigFunctions),
					},
				},
			},
			expectedValues: &secretRender{
				Username: "APPUSER",
				Password: "passphrase",
				Version:  "3",
			},
			expectError:        false,
			exitOnRetryFailure: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			templateTokenCh := make(chan string, 1)
			var templatesToRender []*ctconfig.TemplateConfig
			for fileName, templateTest := range tc.templateMap {
				dstFile := fmt.Sprintf("%s/%s", tmpDir, fileName)
				templateTest.template.Destination = pointerutil.StringPtr(dstFile)
				templatesToRender = append(templatesToRender, templateTest.template)
			}

			ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
			sc := ServerConfig{
				Logger: logging.NewVaultLogger(hclog.Trace),
				AgentConfig: &config.Config{
					Vault: &config.Vault{
						Address: ts.URL,
						Retry: &config.Retry{
							NumRetries: 3,
						},
					},
					TemplateConfig: &config.TemplateConfig{
						ExitOnRetryFailure: tc.exitOnRetryFailure,
					},
				},
				LogLevel:      hclog.Trace,
				LogWriter:     hclog.DefaultOutput,
				ExitAfterAuth: true,
			}

			var server *Server
			server = NewServer(&sc)
			if ts == nil {
				t.Fatal("nil server returned")
			}

			errCh := make(chan error)
			serverErrCh := make(chan error, 1)
			go func() {
				errCh <- server.Run(ctx, templateTokenCh, templatesToRender, &sync.Bool{}, serverErrCh)
			}()

			// send a dummy value to trigger the internal Runner to query for secret
			// info
			templateTokenCh <- "test"

			select {
			case <-ctx.Done():
				t.Fatal("timeout reached before templates were rendered")
			case err := <-errCh:
				if err != nil && !tc.expectError {
					t.Fatalf("did not expect error, got: %v", err)
				}
				if err != nil && tc.expectError {
					t.Logf("received expected error: %v", err)
					return
				}
			}

			// verify test file exists and has the content we're looking for
			var fileCount int
			var errs []string
			for _, template := range templatesToRender {
				if template.Destination == nil {
					t.Fatal("nil template destination")
				}
				content, err := os.ReadFile(*template.Destination)
				if err != nil {
					errs = append(errs, err.Error())
					continue
				}
				fileCount++

				secret := secretRender{}
				if err := json.Unmarshal(content, &secret); err != nil {
					t.Fatal(err)
				}
				var expectedValues secretRender
				if tc.expectedValues != nil {
					expectedValues = *tc.expectedValues
				} else {
					expectedValues = secretRender{
						Username: "appuser",
						Password: "password",
						Version:  "3",
					}
				}
				if secret != expectedValues {
					t.Fatalf("secret didn't match, expected: %#v, got: %#v", expectedValues, secret)
				}
			}
			if len(errs) != 0 {
				t.Fatalf("Failed to find the expected files. Expected %d, got %d\n\t%s", len(templatesToRender), fileCount, strings.Join(errs, "\n\t"))
			}
		})
	}
}

// TestNewServerLogLevels tests that the server can be started with any log
// level.
func TestNewServerLogLevels(t *testing.T) {
	ts := createHttpTestServer()
	defer ts.Close()

	tmpDir, err := os.MkdirTemp("", "agent-tests")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	levels := []hclog.Level{hclog.NoLevel, hclog.Trace, hclog.Debug, hclog.Info, hclog.Warn, hclog.Error}
	for _, level := range levels {
		name := fmt.Sprintf("log_%s", level)
		t.Run(name, func(t *testing.T) {
			server := NewServer(&ServerConfig{
				Logger:        logging.NewVaultLogger(level),
				LogWriter:     hclog.DefaultOutput,
				LogLevel:      level,
				ExitAfterAuth: true,
				AgentConfig: &config.Config{
					Vault: &config.Vault{
						Address: ts.URL,
					},
				},
			})
			if server == nil {
				t.Fatal("nil server returned")
			}
			defer server.Stop()

			templateTokenCh := make(chan string, 1)

			templateTest := &ctconfig.TemplateConfig{
				Contents: pointerutil.StringPtr(templateContents),
			}
			dstFile := fmt.Sprintf("%s/%s", tmpDir, name)
			templateTest.Destination = pointerutil.StringPtr(dstFile)
			templatesToRender := []*ctconfig.TemplateConfig{templateTest}

			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			errCh := make(chan error)
			serverErrCh := make(chan error, 1)
			go func() {
				errCh <- server.Run(ctx, templateTokenCh, templatesToRender, &sync.Bool{}, serverErrCh)
			}()

			// send a dummy value to trigger auth so the server will exit
			templateTokenCh <- "test"

			select {
			case <-ctx.Done():
				t.Fatal("timeout reached before templates were rendered")
			case err := <-errCh:
				if err != nil {
					t.Fatalf("did not expect error, got: %v", err)
				}
			}
		})
	}
}

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

var templateContents = `
{{ with secret "kv/myapp/config"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
{{ if .Data.metadata.version}}"version":"{{ .Data.metadata.version }}"{{ end }}
}
{{ end }}
`

var templateContentsMissingKey = `
{{ with secret "kv/myapp/config"}}
{
{{ if .Data.data.foo}}"foo":"{{ .Data.data.foo}}"{{ end }}
}
{{ end }}
`

var templateContentsBad = `
{{ with secret "kv/myapp/config-bad"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
{{ if .Data.metadata.version}}"version":"{{ .Data.metadata.version }}"{{ end }}
}
{{ end }}
`

var templateContentsPermDenied = `
{{ with secret "kv/myapp/perm-denied"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username}}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password }}",{{ end }}
{{ if .Data.metadata.version}}"version":"{{ .Data.metadata.version }}"{{ end }}
}
{{ end }}
`

var templateContentsWithSprigFunctions = `
{{ with secret "kv/myapp/config"}}
{
{{ if .Data.data.username}}"username":"{{ .Data.data.username | sprig_upper }}",{{ end }}
{{ if .Data.data.password }}"password":"{{ .Data.data.password | sprig_replace "word" "phrase" }}",{{ end }}
{{ if .Data.metadata.version}}"version":"{{ .Data.metadata.version }}"{{ end }}
}
{{ end }}
`
