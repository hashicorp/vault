package api

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseSecret(t *testing.T) {
	raw := strings.TrimSpace(`
{
	"vault_id": "foo",
	"renewable": true,
	"lease_duration": 10,
	"lease_duration_max": 100,
	"key": "value"
}`)

	secret, err := ParseSecret(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Secret{
		VaultId:          "foo",
		Renewable:        true,
		LeaseDuration:    10,
		LeaseDurationMax: 100,
		Data: map[string]interface{}{
			"key": "value",
		},
	}
	if !reflect.DeepEqual(secret, expected) {
		t.Fatalf("bad: %#v", secret)
	}
}
