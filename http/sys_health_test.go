package http

import (
	"io/ioutil"

	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysHealth_get(t *testing.T) {
	core, _, root := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/health")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"initialized": true,
		"sealed":      false,
		"standby":     false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["server_time_utc"] = actual["server_time_utc"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}

	core.Seal(root)

	resp, err = http.Get(addr + "/v1/sys/health")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"initialized": true,
		"sealed":      true,
		"standby":     false,
	}
	testResponseStatus(t, resp, 500)
	testResponseBody(t, resp, &actual)
	expected["server_time_utc"] = actual["server_time_utc"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestSysHealth_customcodes(t *testing.T) {
	core, _, root := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	queryurl, err := url.Parse(addr + "/v1/sys/health?sealedcode=503&activecode=202")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := http.Get(queryurl.String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"initialized": true,
		"sealed":      false,
		"standby":     false,
	}
	testResponseStatus(t, resp, 202)
	testResponseBody(t, resp, &actual)

	expected["server_time_utc"] = actual["server_time_utc"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}

	core.Seal(root)

	queryurl, err = url.Parse(addr + "/v1/sys/health?sealedcode=503&activecode=202")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err = http.Get(queryurl.String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"initialized": true,
		"sealed":      true,
		"standby":     false,
	}
	testResponseStatus(t, resp, 503)
	testResponseBody(t, resp, &actual)
	expected["server_time_utc"] = actual["server_time_utc"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestSysHealth_head(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	testData := []struct {
		uri  string
		code int
	}{
		{"", 200},
		{"?activecode=503", 503},
		{"?activecode=notacode", 400},
	}

	for _, tt := range testData {
		queryurl, err := url.Parse(addr + "/v1/sys/health" + tt.uri)
		if err != nil {
			t.Fatalf("err on %v: %s", queryurl, err)
		}
		resp, err := http.Head(queryurl.String())
		if err != nil {
			t.Fatalf("err on %v: %s", queryurl, err)
		}

		if resp.StatusCode != tt.code {
			t.Fatalf("HEAD %v expected code %d, got %d.", queryurl, tt.code, resp.StatusCode)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("err on %v: %s", queryurl, err)
		}
		if len(data) > 0 {
			t.Fatalf("HEAD %v expected no body, received \"%v\".", queryurl, data)
		}
	}
}
