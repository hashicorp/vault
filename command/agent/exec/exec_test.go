package exec

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ctconfig "github.com/hashicorp/consul-template/config"
	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/pointerutil"
	"github.com/hashicorp/vault/vault"
)

func testVaultServer(t *testing.T) (*api.Client, func()) {
	t.Helper()

	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
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

func TestServer_Run(t *testing.T) {
	testServer := createHttpTestServer()
	defer testServer.Close()

	testCases := map[string]struct {
		templates []*ctconfig.TemplateConfig
	}{
		"simple": {
			templates: []*ctconfig.TemplateConfig{
				{
					Contents:                 pointerutil.StringPtr(`{{ with secret "kv/myapp/config"}}{{.Data.data.username}}{{end}}`),
					MapToEnvironmentVariable: pointerutil.StringPtr("MY_USERNAME"),
				},
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {

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
