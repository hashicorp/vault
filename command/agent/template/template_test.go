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
	}{
		"simple": {
			templateMap: map[string]*templateTest{
				"render_01": &templateTest{
					template: &ctconfig.TemplateConfig{
						Contents: pointerutil.StringPtr(templateContents),
					},
				},
			},
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

			ctx := context.Background()
			sc := ServerConfig{
				Logger: logging.NewVaultLogger(hclog.Trace),
				VaultConf: &config.Vault{
					Address: ts.URL,
				},
				ExitAfterAuth: true,
			}

			var server *Server
			server = NewServer(&sc)
			if ts == nil {
				t.Fatal("nil server returned")
			}

			go server.Run(ctx, templateTokenCh, templatesToRender)

			// send a dummy value to trigger the internal Runner to query for secret
			// info
			templateTokenCh <- "test"
			<-server.DoneCh

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
