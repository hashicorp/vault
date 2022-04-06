package api

import (
	"net/http"
	"testing"
)

func TestIsSudoPath(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(
			`{
				"paths": {
					"/sudo/path": {
						"x-vault-sudo": true
					},
					"/not/a/sudo/path": {
						"x-vault-sudo": false
					},
					"not/a/sudo/path/either": {},
					"/sudo/path/with/template/{name}": {
						"x-vault-sudo": true
					},
					"/sudo/path/with/nested/template/{prefix}/too": {
						"x-vault-sudo": true
					},
					"/sudo/path/with/multiple/templates/{type}/{header}": {
						"x-vault-sudo": true
					}
				}
			}`,
		))
	}

	config, ln := testHTTPServer(t, http.HandlerFunc(handler))
	defer ln.Close()

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	client.SetToken("foo")

	testCases := []struct {
		path     string
		expected bool
	}{
		{
			"not/in/openapi/response",
			false,
		},
		{
			"sudo/path",
			true,
		},
		{
			"not/a/sudo/path",
			false,
		},
		{
			"not/a/sudo/path/either",
			false,
		},
		{
			"sudo/path/with/template/foo",
			true,
		},
		{
			"sudo/path/with/nested/template/foo/too",
			true,
		},
		{
			"sudo/path/with/multiple/templates/foo/bar",
			true,
		},
	}

	for _, tc := range testCases {
		res, err := isSudoPath(client, tc.path)
		if err != nil {
			t.Fatalf("error checking if path is sudo: %v", err)
		}
		if res != tc.expected {
			t.Fatalf("expected isSudoPath to return %v for path %s but it returned %v", tc.expected, tc.path, res)
		}
	}
}
