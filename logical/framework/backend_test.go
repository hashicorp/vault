package framework

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
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

func TestBackendHandleRequest(t *testing.T) {
	callback := func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	resp, err := b.HandleRequest(&logical.Request{
		Operation: logical.ReadOperation,
		Path:      "foo/bar",
		Data:      map[string]interface{}{"value": "42"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp.Data["value"] != 42 {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_badwrite(t *testing.T) {
	callback := func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value").(bool),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeBool},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.UpdateOperation: callback,
				},
			},
		},
	}

	_, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "foo/bar",
		Data:      map[string]interface{}{"value": "3false3"},
	})

	if err == nil {
		t.Fatalf("should have thrown a conversion error")
	}

}

func TestBackendHandleRequest_404(t *testing.T) {
	callback := func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: `foo/bar`,
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	_, err := b.HandleRequest(&logical.Request{
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
			&Path{
				Pattern: "foo/bar",
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				HelpSynopsis:    "foo",
				HelpDescription: "bar",
			},
		},
	}

	resp, err := b.HandleRequest(&logical.Request{
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

	resp, err := b.HandleRequest(&logical.Request{
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

	resp, err := b.HandleRequest(logical.RenewAuthRequest(
		"/foo", &logical.Auth{}, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}
}

func TestBackendHandleRequest_renewAuthCallback(t *testing.T) {
	var called uint32
	callback := func(*logical.Request, *FieldData) (*logical.Response, error) {
		atomic.AddUint32(&called, 1)
		return nil, nil
	}

	b := &Backend{
		AuthRenew: callback,
	}

	_, err := b.HandleRequest(logical.RenewAuthRequest(
		"/foo", &logical.Auth{}, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(&called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}
func TestBackendHandleRequest_renew(t *testing.T) {
	var called uint32
	callback := func(*logical.Request, *FieldData) (*logical.Response, error) {
		atomic.AddUint32(&called, 1)
		return nil, nil
	}

	secret := &Secret{
		Type:  "foo",
		Renew: callback,
	}
	b := &Backend{
		Secrets: []*Secret{secret},
	}

	_, err := b.HandleRequest(logical.RenewRequest(
		"/foo", secret.Response(nil, nil).Secret, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(&called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_renewExtend(t *testing.T) {
	sysView := logical.StaticSystemView{
		DefaultLeaseTTLVal: 5 * time.Minute,
		MaxLeaseTTLVal:     30 * time.Hour,
	}

	secret := &Secret{
		Type:            "foo",
		Renew:           LeaseExtend(0, 0, sysView),
		DefaultDuration: 5 * time.Minute,
	}
	b := &Backend{
		Secrets: []*Secret{secret},
	}

	req := logical.RenewRequest("/foo", secret.Response(nil, nil).Secret, nil)
	req.Secret.IssueTime = time.Now()
	req.Secret.Increment = 1 * time.Hour
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp == nil || resp.Secret == nil {
		t.Fatal("should have secret")
	}

	if resp.Secret.TTL < 59*time.Minute || resp.Secret.TTL > 61*time.Minute {
		t.Fatalf("bad: %s", resp.Secret.TTL)
	}
}

func TestBackendHandleRequest_revoke(t *testing.T) {
	var called uint32
	callback := func(*logical.Request, *FieldData) (*logical.Response, error) {
		atomic.AddUint32(&called, 1)
		return nil, nil
	}

	secret := &Secret{
		Type:   "foo",
		Revoke: callback,
	}
	b := &Backend{
		Secrets: []*Secret{secret},
	}

	_, err := b.HandleRequest(logical.RevokeRequest(
		"/foo", secret.Response(nil, nil).Secret, nil))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(&called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_rollback(t *testing.T) {
	var called uint32
	callback := func(req *logical.Request, kind string, data interface{}) error {
		if data == "foo" {
			atomic.AddUint32(&called, 1)
		}

		return nil
	}

	b := &Backend{
		WALRollback:       callback,
		WALRollbackMinAge: 1 * time.Millisecond,
	}

	storage := new(logical.InmemStorage)
	if _, err := PutWAL(storage, "kind", "foo"); err != nil {
		t.Fatalf("err: %s", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err := b.HandleRequest(&logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(&called); v != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_rollbackMinAge(t *testing.T) {
	var called uint32
	callback := func(req *logical.Request, kind string, data interface{}) error {
		if data == "foo" {
			atomic.AddUint32(&called, 1)
		}

		return nil
	}

	b := &Backend{
		WALRollback:       callback,
		WALRollbackMinAge: 5 * time.Second,
	}

	storage := new(logical.InmemStorage)
	if _, err := PutWAL(storage, "kind", "foo"); err != nil {
		t.Fatalf("err: %s", err)
	}

	_, err := b.HandleRequest(&logical.Request{
		Operation: logical.RollbackOperation,
		Path:      "",
		Storage:   storage,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := atomic.LoadUint32(&called); v != 0 {
		t.Fatalf("bad: %#v", v)
	}
}

func TestBackendHandleRequest_unsupportedOperation(t *testing.T) {
	callback := func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: `foo/bar`,
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	_, err := b.HandleRequest(&logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "foo/bar",
		Data:      map[string]interface{}{"value": "84"},
	})
	if err != logical.ErrUnsupportedOperation {
		t.Fatalf("err: %s", err)
	}
}

func TestBackendHandleRequest_urlPriority(t *testing.T) {
	callback := func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		return &logical.Response{
			Data: map[string]interface{}{
				"value": data.Get("value"),
			},
		}, nil
	}

	b := &Backend{
		Paths: []*Path{
			&Path{
				Pattern: `foo/(?P<value>\d+)`,
				Fields: map[string]*FieldSchema{
					"value": &FieldSchema{Type: TypeInt},
				},
				Callbacks: map[logical.Operation]OperationFunc{
					logical.ReadOperation: callback,
				},
			},
		},
	}

	resp, err := b.HandleRequest(&logical.Request{
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
			[]*Secret{&Secret{Type: "foo"}},
			"bar",
			false,
		},

		"match": {
			[]*Secret{&Secret{Type: "foo"}},
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

		"default duration not set": {
			&FieldSchema{Type: TypeDurationSecond},
			0,
		},
	}

	for name, tc := range cases {
		actual := tc.Schema.DefaultOrZero()
		if !reflect.DeepEqual(actual, tc.Value) {
			t.Fatalf("bad: %s\n\nExpected: %#v\nGot: %#v",
				name, tc.Value, actual)
		}
	}
}
