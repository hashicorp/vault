package audit

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/salt"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
)

func TestCopy_auth(t *testing.T) {
	// Make a non-pointer one so that it can't be modified directly
	expected := logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			TTL:       1 * time.Hour,
			IssueTime: time.Now(),
		},

		ClientToken: "foo",
	}
	auth := expected

	// Copy it
	dup, err := copystructure.Copy(&auth)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check equality
	auth2 := dup.(*logical.Auth)
	if !reflect.DeepEqual(*auth2, expected) {
		t.Fatalf("bad:\n\n%#v\n\n%#v", *auth2, expected)
	}
}

func TestCopy_request(t *testing.T) {
	// Make a non-pointer one so that it can't be modified directly
	expected := logical.Request{
		Data: map[string]interface{}{
			"foo": "bar",
		},
		WrapInfo: &logical.RequestWrapInfo{
			TTL: 60 * time.Second,
		},
	}
	arg := expected

	// Copy it
	dup, err := copystructure.Copy(&arg)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check equality
	arg2 := dup.(*logical.Request)
	if !reflect.DeepEqual(*arg2, expected) {
		t.Fatalf("bad:\n\n%#v\n\n%#v", *arg2, expected)
	}
}

func TestCopy_response(t *testing.T) {
	// Make a non-pointer one so that it can't be modified directly
	expected := logical.Response{
		Data: map[string]interface{}{
			"foo": "bar",
		},
		WrapInfo: &wrapping.ResponseWrapInfo{
			TTL:             60,
			Token:           "foo",
			CreationTime:    time.Now(),
			WrappedAccessor: "abcd1234",
		},
	}
	arg := expected

	// Copy it
	dup, err := copystructure.Copy(&arg)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Check equality
	arg2 := dup.(*logical.Response)
	if !reflect.DeepEqual(*arg2, expected) {
		t.Fatalf("bad:\n\n%#v\n\n%#v", *arg2, expected)
	}
}

func TestHashString(t *testing.T) {
	inmemStorage := &logical.InmemStorage{}
	inmemStorage.Put(&logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	localSalt, err := salt.NewSalt(inmemStorage, &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
	})
	if err != nil {
		t.Fatalf("Error instantiating salt: %s", err)
	}
	out := HashString(localSalt, "foo")
	if out != "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a" {
		t.Fatalf("err: HashString output did not match expected")
	}
}

func TestHash(t *testing.T) {
	now := time.Now()

	cases := []struct {
		Input  interface{}
		Output interface{}
	}{
		{
			&logical.Auth{ClientToken: "foo"},
			&logical.Auth{ClientToken: "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a"},
		},
		{
			&logical.Request{
				Data: map[string]interface{}{
					"foo":              "bar",
					"private_key_type": certutil.PrivateKeyType("rsa"),
				},
			},
			&logical.Request{
				Data: map[string]interface{}{
					"foo":              "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
					"private_key_type": "hmac-sha256:995230dca56fffd310ff591aa404aab52b2abb41703c787cfa829eceb4595bf1",
				},
			},
		},
		{
			&logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",

					// Responses can contain time values, so test that with
					// a known fixed value.
					"bar": now,
				},
				WrapInfo: &wrapping.ResponseWrapInfo{
					TTL:             60,
					Token:           "bar",
					Accessor:        "flimflam",
					CreationTime:    now,
					WrappedAccessor: "bar",
				},
			},
			&logical.Response{
				Data: map[string]interface{}{
					"foo": "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
					"bar": now.Format(time.RFC3339Nano),
				},
				WrapInfo: &wrapping.ResponseWrapInfo{
					TTL:             60,
					Token:           "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
					Accessor:        "hmac-sha256:7c9c6fe666d0af73b3ebcfbfabe6885015558213208e6635ba104047b22f6390",
					CreationTime:    now,
					WrappedAccessor: "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
				},
			},
		},
		{
			"foo",
			"foo",
		},
		{
			&logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       1 * time.Hour,
					IssueTime: now,
				},

				ClientToken: "foo",
			},
			&logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       1 * time.Hour,
					IssueTime: now,
				},

				ClientToken: "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a",
			},
		},
	}

	inmemStorage := &logical.InmemStorage{}
	inmemStorage.Put(&logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	localSalt, err := salt.NewSalt(inmemStorage, &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
	})
	if err != nil {
		t.Fatalf("Error instantiating salt: %s", err)
	}
	for _, tc := range cases {
		input := fmt.Sprintf("%#v", tc.Input)
		if err := Hash(localSalt, tc.Input); err != nil {
			t.Fatalf("err: %s\n\n%s", err, input)
		}
		if _, ok := tc.Input.(*logical.Response); ok {
			if !reflect.DeepEqual(tc.Input.(*logical.Response).WrapInfo, tc.Output.(*logical.Response).WrapInfo) {
				t.Fatalf("bad:\nInput:\n%s\nTest case input:\n%#v\nTest case output\n%#v", input, tc.Input.(*logical.Response).WrapInfo, tc.Output.(*logical.Response).WrapInfo)
			}
		}
		if !reflect.DeepEqual(tc.Input, tc.Output) {
			t.Fatalf("bad:\nInput:\n%s\nTest case input:\n%#v\nTest case output\n%#v", input, tc.Input, tc.Output)
		}
	}
}

func TestHashWalker(t *testing.T) {
	replaceText := "foo"

	cases := []struct {
		Input  interface{}
		Output interface{}
	}{
		{
			map[string]interface{}{
				"hello": "foo",
			},
			map[string]interface{}{
				"hello": replaceText,
			},
		},

		{
			map[string]interface{}{
				"hello": []interface{}{"world"},
			},
			map[string]interface{}{
				"hello": []interface{}{replaceText},
			},
		},
	}

	for _, tc := range cases {
		output, err := HashStructure(tc.Input, func(string) string {
			return replaceText
		})
		if err != nil {
			t.Fatalf("err: %s\n\n%#v", err, tc.Input)
		}
		if !reflect.DeepEqual(output, tc.Output) {
			t.Fatalf("bad:\n\n%#v\n\n%#v", tc.Input, output)
		}
	}
}
