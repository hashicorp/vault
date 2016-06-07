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
	},
	"warnings": [
		"a warning!"
	],
	"wrap_info": {
		"token": "token",
		"ttl": 60,
		"creation_time": 100000
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
		Warnings: []string{
			"a warning!",
		},
		WrapInfo: &SecretWrapInfo{
			Token:        "token",
			TTL:          60,
			CreationTime: int64(100000),
		},
	}
	if !reflect.DeepEqual(secret, expected) {
		t.Fatalf("bad: %#v %#v", secret, expected)
	}
}
