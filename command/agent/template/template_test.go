package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
)

// TestNewServer is a simple test to make sure NewServer returns a Server and
// channel
func TestNewServer(t *testing.T) {
	server := NewServer(&ServerConfig{})
	if server == nil {
		t.Fatal("nil server returned")
	}
}

func newAgentConfig(listeners []*configutil.Listener, enableCache bool) *config.Config {
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
		agentConfig.Cache = &config.Cache{UseAutoAuthToken: true}
	}

	return agentConfig
}

func TestCacheConfigUnix(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:        "unix",
			Address:     "foobar",
			TLSDisable:  true,
			SocketMode:  "configmode",
			SocketUser:  "configuser",
			SocketGroup: "configgroup",
		},
		{
			Type:       "tcp",
			Address:    "127.0.0.1:8300",
			TLSDisable: true,
		},
		{
			Type:        "tcp",
			Address:     "127.0.0.1:8400",
			TLSKeyFile:  "/path/to/cakey.pem",
			TLSCertFile: "/path/to/cacert.pem",
		},
	}

	agentConfig := newAgentConfig(listeners, true)
	serverConfig := ServerConfig{AgentConfig: agentConfig}

	ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.HasPrefix(*ctConfig.Vault.Address, "unix") {
		t.Fatalf("expected unix address, got %s", *ctConfig.Vault.Address)
	}

	expected := "unix:/foobar"
	if *ctConfig.Vault.Address != expected {
		t.Fatalf("expected %s, got %s", expected, *ctConfig.Vault.Address)
	}
}
func TestCacheConfigHTTP(t *testing.T) {
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

	agentConfig := newAgentConfig(listeners, true)
	serverConfig := ServerConfig{AgentConfig: agentConfig}

	ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.HasPrefix(*ctConfig.Vault.Address, "http") {
		t.Fatalf("expected http address, got %s", *ctConfig.Vault.Address)
	}

	expected := "http://127.0.0.1:8300"
	if *ctConfig.Vault.Address != expected {
		t.Fatalf("expected %s, got %s", expected, *ctConfig.Vault.Address)
	}
}

func TestCacheConfigHTTPS(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:        "tcp",
			Address:     "127.0.0.1:8300",
			TLSKeyFile:  "/path/to/cakey.pem",
			TLSCertFile: "/path/to/cacert.pem",
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
			Type:       "tcp",
			Address:    "127.0.0.1:8400",
			TLSDisable: true,
		},
	}

	agentConfig := newAgentConfig(listeners, true)
	serverConfig := ServerConfig{AgentConfig: agentConfig}

	ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.HasPrefix(*ctConfig.Vault.Address, "https") {
		t.Fatalf("expected https address, got %s", *ctConfig.Vault.Address)
	}

	expected := "https://127.0.0.1:8300"
	if *ctConfig.Vault.Address != expected {
		t.Fatalf("expected %s, got %s", expected, *ctConfig.Vault.Address)
	}
}

func TestCacheConfigNoCache(t *testing.T) {
	listeners := []*configutil.Listener{
		{
			Type:        "tcp",
			Address:     "127.0.0.1:8300",
			TLSKeyFile:  "/path/to/cakey.pem",
			TLSCertFile: "/path/to/cacert.pem",
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
			Type:       "tcp",
			Address:    "127.0.0.1:8400",
			TLSDisable: true,
		},
	}

	agentConfig := newAgentConfig(listeners, false)
	serverConfig := ServerConfig{AgentConfig: agentConfig}

	ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.HasPrefix(*ctConfig.Vault.Address, "http") {
		t.Fatalf("expected http address, got %s", *ctConfig.Vault.Address)
	}

	expected := "http://127.0.0.1:1111"
	if *ctConfig.Vault.Address != expected {
		t.Fatalf("expected %s, got %s", expected, *ctConfig.Vault.Address)
	}
}

func TestCacheConfigNoListener(t *testing.T) {
	listeners := []*configutil.Listener{}

	agentConfig := newAgentConfig(listeners, true)
	serverConfig := ServerConfig{AgentConfig: agentConfig}

	ctConfig, err := newRunnerConfig(&serverConfig, ctconfig.TemplateConfigs{})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !strings.HasPrefix(*ctConfig.Vault.Address, "http") {
		t.Fatalf("expected http address, got %s", *ctConfig.Vault.Address)
	}

	expected := "http://127.0.0.1:1111"
	if *ctConfig.Vault.Address != expected {
		t.Fatalf("expected %s, got %s", expected, *ctConfig.Vault.Address)
	}
}

func TestServerRun(t *testing.T) {
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

	ts := httptest.NewServer(mux)
	defer ts.Close()
	tmpDir, err := ioutil.TempDir("", "agent-tests")
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
		templateMap map[string]*templateTest
		expectError bool
	}{
		"simple": {
			templateMap: map[string]*templateTest{
				"render_01": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
			},
			expectError: false,
		},
		"multiple": {
			templateMap: map[string]*templateTest{
				"render_01": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_02": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_03": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_04": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_05": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_06": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
				"render_07": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
			},
			expectError: false,
		},
		"bad secret": {
			templateMap: map[string]*templateTest{
				"render_01": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContentsBad),
					},
				},
			},
			expectError: true,
		},
		"missing key": {
			templateMap: map[string]*templateTest{
				"render_01": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents:      pointerutil.StringPtr(templateContentsMissingKey),
						ErrMissingKey: pointerutil.BoolPtr(true),
					},
				},
			},
			expectError: true,
		},
		"permission denied": {
			templateMap: map[string]*templateTest{
				"render_01": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContentsPermDenied),
					},
				},
			},
			expectError: true,
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
			server.testingLimitRetry = 3

			errCh := make(chan error)
			go func() {
				errCh <- server.Run(ctx, templateTokenCh, templatesToRender)
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
			for _, template := range templatesToRender {
				if template.Destination == nil {
					t.Fatal("nil template destination")
				}
				content, err := ioutil.ReadFile(*template.Destination)
				if err != nil {
					t.Fatal(err)
				}
				fileCount++

				secret := secretRender{}
				if err := json.Unmarshal(content, &secret); err != nil {
					t.Fatal(err)
				}
				if secret.Username != "appuser" || secret.Password != "password" || secret.Version != "3" {
					t.Fatalf("secret didn't match: %#v", secret)
				}
			}
			if fileCount != len(templatesToRender) {
				t.Fatalf("mismatch file to template: (%d) / (%d)", fileCount, len(templatesToRender))
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
