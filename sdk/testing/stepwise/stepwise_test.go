package stepwise

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
)

type testRun struct {
	mt        *mockT
	expectedT *mockT
	testCase  Case
}

func TestStepwise_SkipIfNotAcc(t *testing.T) {
	if err := os.Setenv(TestEnvVar, ""); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Setenv(TestEnvVar, "1")
	skipCase := Case{
		Environment: new(mockEnvironment),
		Steps:       []Step{Step{}},
	}

	mt := new(mockT)
	et := mockT{
		SkipCalled: true,
	}

	testRun := testRun{
		mt: mt,
		expectedT: &mockT{
			SkipCalled: true,
		},
		testCase: skipCase,
	}

	Run(mt, testRun.testCase)

	if mt.SkipCalled != et.SkipCalled {
		t.Fatalf("expected SkipCalled (%t), got (%t)", et.SkipCalled, mt.SkipCalled)
	}
}

func TestStepwise_Run(t *testing.T) {
	basicCase := func() Case {
		// envOptions := stepwise.EnvironmentOptions{
		// 	Name:            "transit2",
		// 	PluginType:      stepwise.PluginTypeSecrets,
		// 	PluginName:      "transit",
		// 	MountPathPrefix: "transit_temp",
		// }
		return Case{
			Environment: new(mockEnvironment),
			Steps: []Step{
				Step{
					Operation: ListOperation,
					Path:      "keys",
					Check: func(resp *api.Secret, err error) error {
						return nil
					},
				},
				// testAccStepwiseWritePolicy(t, "test", true),
			},
		}
	}

	testRuns := map[string]testRun{
		"basic": {
			mt:        new(mockT),
			expectedT: new(mockT),
			testCase:  basicCase(),
		},
	}

	for name, tr := range testRuns {
		t.Run(name, func(t *testing.T) {
			Run(tr.mt, tr.testCase)
		})
	}
}

func TestStepwise_makeRequest(t *testing.T) {
	me := new(mockEnvironment)
	me.Setup()
	mt := new(mockT)

	type testRequest struct {
		Operation         StepOperation
		Path              string
		ExpectedRequestID string
		ExpectErr         bool
	}
	testRequests := map[string]testRequest{
		"list": testRequest{
			Operation:         ListOperation,
			Path:              "keys",
			ExpectedRequestID: "list-request",
		},
		"read": testRequest{
			Operation:         ReadOperation,
			Path:              "keys/name",
			ExpectedRequestID: "read-request",
		},
		"error": testRequest{
			Operation: ReadOperation,
			Path:      "error",
			ExpectErr: true,
		},
	}

	for name, tc := range testRequests {
		t.Run(name, func(t *testing.T) {
			step := Step{
				Operation: tc.Operation,
				Path:      tc.Path,
			}

			secret, err := makeRequest(mt, me, step)
			if err != nil && !tc.ExpectErr {
				t.Fatalf("unexpected error: %s", err)
			}
			if err == nil && tc.ExpectErr {
				t.Fatal("expected error but got none:")
			}

			if err != nil && tc.ExpectErr {
				return
			}
			if secret.RequestID != tc.ExpectedRequestID {
				t.Fatalf("expected (%s), got (%s)", tc.ExpectedRequestID, secret.RequestID)
			}
		})
	}
}

type mockEnvironment struct {
	ts     *httptest.Server
	client *api.Client
	l      sync.Mutex
}

// Setup creates the mock environment, establishing a test HTTP server
func (m *mockEnvironment) Setup() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/test/keys", func(w http.ResponseWriter, req *http.Request) {
		listResp := api.Secret{
			RequestID: "list-request",
			LeaseID:   "lease-id",
		}
		out, err := jsonutil.EncodeJSON(listResp)
		if err != nil {
			panic(err)
		}
		w.Write(out)
	})
	mux.HandleFunc("/v1/test/keys/name", func(w http.ResponseWriter, req *http.Request) {
		readResp := api.Secret{
			RequestID: "read-request",
			LeaseID:   "lease-id",
		}
		out, err := jsonutil.EncodeJSON(readResp)
		if err != nil {
			panic(err)
		}
		w.Write(out)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		fmt.Fprintf(w, "{}")
	})
	m.ts = httptest.NewServer(mux)

	return nil
}

// Client creates a Vault API client configured to the mock environment's test
// server
func (m *mockEnvironment) Client() (*api.Client, error) {
	m.l.Lock()
	defer m.l.Unlock()
	// this shouldn't be needed but being defensive
	if m.ts == nil {
		if err := m.Setup(); err != nil {
			return nil, err
		}
	}

	if m.client != nil {
		return m.client, nil
	}

	cfg := api.Config{
		Address:    m.ts.URL,
		HttpClient: cleanhttp.DefaultPooledClient(),
		Timeout:    time.Second * 5,
		MaxRetries: 2,
	}

	return api.NewClient(&cfg)
}

func (m *mockEnvironment) Teardown() error {
	m.ts.Close()
	return nil
}
func (m *mockEnvironment) Name() string { return "" }
func (m *mockEnvironment) ExpandPath(path string) string {
	return "/test/" + path
}
func (m *mockEnvironment) MountPath() string { return "" }
func (m *mockEnvironment) RootToken() string { return "" }
