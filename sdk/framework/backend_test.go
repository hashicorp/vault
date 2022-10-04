package framework

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func BenchmarkBackendRoute(b *testing.B) {
	patterns := []string{
		"foo",
		"bar/(?P<name>.+?)",
		"baz/(?P<name>what)",
		`aws/policy/(?P<policy>\w)`,
		`aws/(?P<policy>\w)`,
	}

	backend := &Backend{Paths: make([]*Path, 0, len(patterns))}
	for _, p := range patterns {
		backend.Paths = append(backend.Paths, &Path{Pattern: p})
	}

	// Warm any caches
	backend.Route("aws/policy/foo")

	// Reset the timer since we did a lot above
	b.ResetTimer()

	// Run through and route. We do a sanity check of the return value
	for i := 0; i < b.N; i++ {
		if p := backend.Route("aws/policy/foo"); p == nil {
			b.Fatal("p should not be nil")
		}
	}
}

func TestBackend_impl(t *testing.T) {
	var _ logical.Backend = new(Backend)
}

func TestBackendHandleRequestFieldWarnings(t *testing.T) {
	handler := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"an_int":   data.Get("an_int"),
				"a_string": data.Get("a_string"),
				"name":     data.Get("name"),
			},
		}, nil
	}

	backend := &Backend{
		Paths: []*Path{
			{
				Pattern: "foo/bar/(?P<name>.+)",
				Fields: map[string]*FieldSchema{
					"an_int":   {Type: TypeInt},
					"a_string": {Type: TypeString},
					"name":     {Type: TypeString},
				},
				Operations: map[logical.Operation]OperationHandler{
					logical.UpdateOperation: &PathOperation{Callback: handler},
				},
			},
		},
	}
	ctx := context.Background()
	resp, err := backend.HandleRequest(ctx, &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "foo/bar/baz",
		Data: map[string]interface{}{
			"an_int":        10,
			"a_string":      "accepted",
			"unrecognized1": "unrecognized",
			"unrecognized2": 20.2,
			"name":          "noop",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	t.Log(resp.Warnings)
	require.Len(t, resp.Warnings, 2)
	require.True(t, strutil.StrListContains(resp.Warnings, "Endpoint ignored these unrecognized parameters: [unrecognized1 unrecognized2]"))
	require.True(t, strutil.StrListContains(resp.Warnings, "Endpoint replaced the value of these parameters with the values captured from the endpoint's path: [name]"))
}

func TestBackendHandleRequest(t *testing.T) {
	callback := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}
	handler := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"amount": data.Get("amount"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": {Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
			{
				Pattern: "foo/baz/handler",
				Fields: map[string]*FieldSchema{
					"amount": {Type: TypeInt},
				},
				Operations: map[logical.Operation]OperationHandler{
					logical.ReadOperation: &PathOperation{Callback: handler},
				},
			},
			{
				Pattern: "foo/both/handler",
				Fields: map[string]*FieldSchema{
					"amount": {Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
				Operations: map[logical.Operation]OperationHandler{
					logical.ReadOperation: &PathOperation{Callback: handler},
				},
			},
		},
		system: &logical.StaticSystemView{},
	}

	for _, path := range []string{"foo/bar", "foo/baz/handler", "foo/both/handler"} {
		key := "value"
		if strings.Contains(path, "handler") {
			key = "amount"
		}
		resp, err := b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.ReadOperation,
			Path:      path,
			Data:      map[string]interface{}{key: "42"},
		})
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if resp.Data[key] != 42 {
			t.Fatalf("bad: %#v", resp)
		}
	}
}

func TestBackendHandleRequest_Forwarding(t *testing.T) {
	tests := map[string]struct {
		fwdStandby   bool
		fwdSecondary bool
		isLocal      bool
		isStandby    bool
		isSecondary  bool
		expectFwd    bool
		nilSysView   bool
	}{
		"no forward": {
			expectFwd: false,
		},
		"no forward, local restricted": {
			isSecondary:  true,
			fwdSecondary: true,
			isLocal:      true,
			expectFwd:    false,
		},
		"no forward, forwarding not requested": {
			isSecondary: true,
			isStandby:   true,
			expectFwd:   false,
		},
		"forward, secondary": {
			fwdSecondary: true,
			isSecondary:  true,
			expectFwd:    true,
		},
		"forward, standby": {
			fwdStandby: true,
			isStandby:  true,
			expectFwd:  true,
		},
		"no forward, only secondary": {
			fwdSecondary: true,
			isStandby:    true,
			expectFwd:    false,
		},
		"no forward, only standby": {
			fwdStandby:  true,
			isSecondary: true,
			expectFwd:   false,
		},
		"nil system view": {
			nilSysView: true,
			expectFwd:  false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var replState consts.ReplicationState
			if test.isStandby {
				replState.AddState(consts.ReplicationPerformanceStandby)
			}
			if test.isSecondary {
				replState.AddState(consts.ReplicationPerformanceSecondary)
			}

			b := &Backend{
				Paths: []*Path{
					{
						Pattern: "foo",
						Operations: map[logical.Operation]OperationHandler{
							logical.ReadOperation: &PathOperation{
								Callback: func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
									return nil, nil
								},
								ForwardPerformanceSecondary: test.fwdSecondary,
								ForwardPerformanceStandby:   test.fwdStandby,
							},
						},
					},
				},

				system: &logical.StaticSystemView{
					LocalMountVal:       test.isLocal,
					ReplicationStateVal: replState,
				},
			}

			if test.nilSysView {
				b.system = nil
			}

			_, err := b.HandleRequest(context.Background(), &logical.Request{
				Operation: logical.ReadOperation,
				Path:      "foo",
			})

			if !test.expectFwd && err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if test.expectFwd && err != logical.ErrReadOnly {
				t.Fatalf("expected ErrReadOnly, got: %v", err)
			}
		})
	}
}

func TestBackendHandleRequest_badwrite(t *testing.T) {
	callback := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value").(bool),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": {Type: TypeBool},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.UpdateOperation: callback,
				},
			},
		},
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "foo/bar",
		Data:      map[string]interface{}{"value": "3false3"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !strings.Contains(resp.Data["error"].(string), "Field validation failed") {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_404(t *testing.T) {
	callback := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			{
				Pattern: `foo/bar`,
				Fields: map[string]*FieldSchema{
					"value": {Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "foo/baz",
		Data:      map[string]interface{}{"value": "84"},
	})
	if err != logical.ErrUnsupportedPath {
		t.Fatalf("err: %s", err)
	}
}

func TestBackendHandleRequest_help(t *testing.T) {
	b := &Backend{
		Paths: []*Path{
			{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": {Type: TypeInt},
				},
				HelpSynopsis:    "foo",
				HelpDescription: "bar",
			},
		},
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.HelpOperation,
		Path:      "foo/bar",
		Data:      map[string]interface{}{"value": "42"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data["help"] == nil {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_helpRoot(t *testing.T) {
	b := &Backend{
		Help: "42",
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.HelpOperation,
		Path:      "",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data["help"] == nil {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_renewAuth(t *testing.T) {
	b := &Backend{}

	resp, err := b.HandleRequest(context.Background(), logical.RenewAuthRequest("/foo", &logical.Auth{}, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_renewAuthCallback(t *testing.T) {
	called := new(uint32)
	callback := func(context.Context, *logical.Request, *FieldData) (*logical.Response, error) {
		atomic.AddUint32(called, 1)
		return nil, nil
	}

	b := &Backend{
		AuthRenew: callback,
	}

	_, err := b.HandleRequest(context.Background(), logical.RenewAuthRequest("/foo", &logical.Auth{}, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_renew(t *testing.T) {
	called := new(uint32)
	callback := func(context.Context, *logical.Request, *FieldData) (*logical.Response, error) {
		atomic.AddUint32(called, 1)
		return nil, nil
	}

	secret := &Secret{
		Type:  "foo",
		Renew: callback,
	}
	b := &Backend{
		Secrets: []*Secret{secret},
	}

	_, err := b.HandleRequest(context.Background(), logical.RenewRequest("/foo", secret.Response(nil, nil).Secret, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_revoke(t *testing.T) {
	called := new(uint32)
	callback := func(context.Context, *logical.Request, *FieldData) (*logical.Response, error) {
		atomic.AddUint32(called, 1)
		return nil, nil
	}

	secret := &Secret{
		Type:   "foo",
		Revoke: callback,
	}
	b := &Backend{
		Secrets: []*Secret{secret},
	}

	_, err := b.HandleRequest(context.Background(), logical.RevokeRequest("/foo", secret.Response(nil, nil).Secret, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_rollback(t *testing.T) {
	called := new(uint32)
	callback := func(_ context.Context, req *logical.Request, kind string, data interface{}) error {
		if data == "foo" {
			atomic.AddUint32(called, 1)
		}
		return nil
	}

	b := &Backend{
		WALRollback:       callback,
		WALRollbackMinAge: 1 * time.Millisecond,
	}

	storage := new(logical.InmemStorage)
	if _, err := PutWAL(context.Background(), storage, "kind", "foo"); err != nil {
		t.Fatalf("err: %s", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_rollbackMinAge(t *testing.T) {
	called := new(uint32)
	callback := func(_ context.Context, req *logical.Request, kind string, data interface{}) error {
		if data == "foo" {
			atomic.AddUint32(called, 1)
		}
		return nil
	}

	b := &Backend{
		WALRollback:       callback,
		WALRollbackMinAge: 5 * time.Second,
	}

	storage := new(logical.InmemStorage)
	if _, err := PutWAL(context.Background(), storage, "kind", "foo"); err != nil {
		t.Fatalf("err: %s", err)
	}

	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(called); v != 0 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_unsupportedOperation(t *testing.T) {
	callback := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			{
				Pattern: `foo/bar`,
				Fields: map[string]*FieldSchema{
					"value": {Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "foo/bar",
		Data:      map[string]interface{}{"value": "84"},
	})
	if err != logical.ErrUnsupportedOperation {
		t.Fatalf("err: %s", err)
	}
}

func TestBackendHandleRequest_urlPriority(t *testing.T) {
	callback := func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			{
				Pattern: `foo/(?P<value>\d+)`,
				Fields: map[string]*FieldSchema{
					"value": {Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "foo/42",
		Data:      map[string]interface{}{"value": "84"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data["value"] != 42 {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendRoute(t *testing.T) {
	cases := map[string]struct {
		Patterns []string
		Path     string
		Match    string
	}{
		"no match": {
			[]string{"foo"},
			"bar",
			"",
		},

		"exact": {
			[]string{"foo"},
			"foo",
			"^foo$",
		},

		"regexp": {
			[]string{"fo+"},
			"foo",
			"^fo+$",
		},

		"anchor-start": {
			[]string{"bar"},
			"foobar",
			"",
		},

		"anchor-end": {
			[]string{"bar"},
			"barfoo",
			"",
		},

		"anchor-ambiguous": {
			[]string{"mounts", "sys/mounts"},
			"sys/mounts",
			"^sys/mounts$",
		},
	}

	for n, tc := range cases {
		paths := make([]*Path, len(tc.Patterns))
		for i, pattern := range tc.Patterns {
			paths[i] = &Path{Pattern: pattern}
		}

		b := &Backend{Paths: paths}
		result := b.Route(tc.Path)
		match := ""
		if result != nil {
			match = result.Pattern
		}

		if match != tc.Match {
			t.Fatalf("bad: %s\n\nExpected: %s\nGot: %s",
				n, tc.Match, match)
		}
	}
}

func TestBackendSecret(t *testing.T) {
	cases := map[string]struct {
		Secrets []*Secret
		Search  string
		Match   bool
	}{
		"no match": {
			[]*Secret{{Type: "foo"}},
			"bar",
			false,
		},

		"match": {
			[]*Secret{{Type: "foo"}},
			"foo",
			true,
		},
	}

	for n, tc := range cases {
		b := &Backend{Secrets: tc.Secrets}
		result := b.Secret(tc.Search)
		if tc.Match != (result != nil) {
			t.Fatalf("bad: %s\n\nExpected match: %v", n, tc.Match)
		}
		if result != nil && result.Type != tc.Search {
			t.Fatalf("bad: %s\n\nExpected matching type: %#v", n, result)
		}
	}
}

func TestFieldSchemaDefaultOrZero(t *testing.T) {
	cases := map[string]struct {
		Schema *FieldSchema
		Value  interface{}
	}{
		"default set": {
			&FieldSchema{Type: TypeString, Default: "foo"},
			"foo",
		},

		"default not set": {
			&FieldSchema{Type: TypeString},
			"",
		},

		"default duration set": {
			&FieldSchema{Type: TypeDurationSecond, Default: 60},
			60,
		},

		"default duration int64": {
			&FieldSchema{Type: TypeDurationSecond, Default: int64(60)},
			60,
		},

		"default duration string": {
			&FieldSchema{Type: TypeDurationSecond, Default: "60s"},
			60,
		},

		"illegal default duration string": {
			&FieldSchema{Type: TypeDurationSecond, Default: "h1"},
			0,
		},

		"default duration time.Duration": {
			&FieldSchema{Type: TypeDurationSecond, Default: 60 * time.Second},
			60,
		},

		"default duration not set": {
			&FieldSchema{Type: TypeDurationSecond},
			0,
		},

		"default signed positive duration set": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: 60},
			60,
		},

		"default signed positive duration int64": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: int64(60)},
			60,
		},

		"default signed positive duration string": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: "60s"},
			60,
		},

		"illegal default signed duration string": {
			&FieldSchema{Type: TypeDurationSecond, Default: "-h1"},
			0,
		},

		"default signed positive duration time.Duration": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: 60 * time.Second},
			60,
		},

		"default signed negative duration set": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: -60},
			-60,
		},

		"default signed negative duration int64": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: int64(-60)},
			-60,
		},

		"default signed negative duration string": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: "-60s"},
			-60,
		},

		"default signed negative duration time.Duration": {
			&FieldSchema{Type: TypeSignedDurationSecond, Default: -60 * time.Second},
			-60,
		},

		"default signed negative duration not set": {
			&FieldSchema{Type: TypeSignedDurationSecond},
			0,
		},
		"default header not set": {
			&FieldSchema{Type: TypeHeader},
			http.Header{},
		},
	}

	for name, tc := range cases {
		actual := tc.Schema.DefaultOrZero()
		if !reflect.DeepEqual(actual, tc.Value) {
			t.Errorf("bad: %s\n\nExpected: %#v\nGot: %#v",
				name, tc.Value, actual)
		}
	}
}

func TestInitializeBackend(t *testing.T) {
	var inited bool
	backend := &Backend{InitializeFunc: func(context.Context, *logical.InitializationRequest) error {
		inited = true
		return nil
	}}

	backend.Initialize(nil, &logical.InitializationRequest{Storage: nil})

	if !inited {
		t.Fatal("backend should be open")
	}
}
