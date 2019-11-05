package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	ctconfig "github.com/hashicorp/consul-template/config"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
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
	if server.UnblockCh == nil {
		t.Fatal("nil blocking channel returned")
	}
}

func TestServerRun(t *testing.T) {
	// create http test server
	ts := httptest.NewServer(http.HandlerFunc(handleRequest))
	defer ts.Close()
	tmpDir, err := ioutil.TempDir("", "agent-tests")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]struct {
		templates []*ctconfig.TemplateConfig
	}{
		"basic": {
			templates: []*ctconfig.TemplateConfig{
				&ctconfig.TemplateConfig{
					Contents: pointerutil.StringPtr(templateContents),
				},
			},
		},
	}

	// secretRender is a simple struct that represents the secret we render to
	// disk. It's used to unmarshal the file contents and test against
	type secretRender struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Version  string `json:"version"`
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			templateTokenCh := make(chan string, 1)
			for i, template := range tc.templates {
				dstFile := fmt.Sprintf("%s/render_%d.txt", tmpDir, i)
				template.Destination = pointerutil.StringPtr(dstFile)
			}

			ctx, cancelFunc := context.WithCancel(context.Background())
			sc := ServerConfig{
				Logger: logging.NewVaultLogger(hclog.Trace),
				VaultConf: &config.Vault{
					Address: ts.URL,
				},
			}

			var server *Server
			server = NewServer(&sc)
			if ts == nil {
				t.Fatal("nil server returned")
			}
			if server.UnblockCh == nil {
				t.Fatal("nil blocking channel returned")
			}

			go server.Run(ctx, templateTokenCh, tc.templates)

			// send a dummy value to trigger the internal Runner to query for secret
			// info
			templateTokenCh <- "test"

			select {
			case <-ctx.Done():
			case <-server.UnblockCh:
			}

			// cancel to clean things up
			cancelFunc()

			// verify test file exists and has the content we're looking for
			for _, template := range tc.templates {
				if template.Destination == nil {
					t.Fatal("nil template destination")
				}
				content, err := ioutil.ReadFile(*template.Destination)
				if err != nil {
					t.Fatal(err)
				}

				secret := secretRender{}
				if err := json.Unmarshal(content, &secret); err != nil {
					t.Fatal(err)
				}
				if secret.Username != "appuser" || secret.Password != "password" || secret.Version != "3" {
					t.Fatalf("secret didn't match: %#v", secret)
				}
			}
		})
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, jsonResponse)
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
