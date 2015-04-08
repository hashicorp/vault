package api

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseSecret(t *testing.T) {
	raw := strings.TrimSpace(`
{
	"lease_id": "foo",
	"renewable": true,
	"lease_duration": 10,
	"data": {
		"key": "value"
	}
}`)

	secret, err := ParseSecret(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Secret{
		LeaseID:       "foo",
		Renewable:     true,
		LeaseDuration: 10,
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	if !reflect.DeepEqual(secret, expected) {
		t.Fatalf("bad: %#v %#v", secret, expected)
	}
}
