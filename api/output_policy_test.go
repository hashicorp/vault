package api

import (
	"net/http"
	"strings"
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

	config.Address = strings.ReplaceAll(config.Address, "127.0.0.1", "localhost")
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	client.SetToken("foo")

	client.SetOutputPolicy(true)

	testCases := []struct {
		path     string
		expected bool
	}{
		{
			"secret/foo", // not in the openAPI response at all
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
