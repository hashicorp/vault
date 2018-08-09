package api

import (
	"strings"
	"testing"
)

func TestRequestSetJSONBody(t *testing.T) {
	var r Request
	raw := map[string]interface{}{"foo": "bar"}
	if err := r.SetJSONBody(raw); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := `{"foo":"bar"}`
	actual := strings.TrimSpace(string(r.BodyBytes))
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestRequestResetJSONBody(t *testing.T) {
	var r Request
	raw := map[string]interface{}{"foo": "bar"}
	if err := r.SetJSONBody(raw); err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := r.ResetJSONBody(); err != nil {
		t.Fatalf("err: %s", err)
	}

	buf := make([]byte, len(r.BodyBytes))
	copy(buf, r.BodyBytes)

	expected := `{"foo":"bar"}`
	actual := strings.TrimSpace(string(buf))
	if actual != expected {
		t.Fatalf("bad: actual %s, expected %s", actual, expected)
	}
}
