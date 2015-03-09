package api

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRequestSetJSONBody(t *testing.T) {
	var r Request
	raw := map[string]interface{}{"foo": "bar"}
	if err := r.SetJSONBody(raw); err != nil {
		t.Fatalf("err: %s", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := `{"foo":"bar"}`
	actual := strings.TrimSpace(buf.String())
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}

	if int64(len(buf.String())) != r.BodySize {
		t.Fatalf("bad: %d", len(actual))
	}
}
