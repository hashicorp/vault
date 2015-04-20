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

func TestRequestResetJSONBody(t *testing.T) {
	var r Request
	raw := map[string]interface{}{"foo": "bar"}
	if err := r.SetJSONBody(raw); err != nil {
		t.Fatalf("err: %s", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r.Body); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := r.ResetJSONBody(); err != nil {
		t.Fatalf("err: %s", err)
	}

	var buf2 bytes.Buffer
	if _, err := io.Copy(&buf2, r.Body); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := `{"foo":"bar"}`
	actual := strings.TrimSpace(buf2.String())
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}

	if int64(len(buf2.String())) != r.BodySize {
		t.Fatalf("bad: %d", len(actual))
	}
}
