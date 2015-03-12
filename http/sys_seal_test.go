package http

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSysSealStatus(t *testing.T) {
	core := testCore(t)
	testCoreInit(t, core)
	ln, addr := testServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"sealed":   true,
		"t":        float64(1),
		"n":        float64(1),
		"progress": float64(0),
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
