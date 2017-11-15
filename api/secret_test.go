package api

import (
	"reflect"
	"strings"
	"testing"
	"time"
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
		"accessor": "accessor",
		"ttl": 60,
		"creation_time": "2016-06-07T15:52:10-04:00",
		"wrapped_accessor": "abcd1234"
	}
}`)

	rawTime, _ := time.Parse(time.RFC3339, "2016-06-07T15:52:10-04:00")

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
			Token:           "token",
			Accessor:        "accessor",
			TTL:             60,
			CreationTime:    rawTime,
			WrappedAccessor: "abcd1234",
		},
	}
	if !reflect.DeepEqual(secret, expected) {
		t.Fatalf("bad:\ngot\n%#v\nexpected\n%#v\n", secret, expected)
	}
}
