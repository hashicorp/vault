// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
)

func TestCopy_auth(t *testing.T) {
	// Make a non-pointer one so that it can't be modified directly
	expected := logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			TTL: 1 * time.Hour,
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

// TestSalter is a structure that implements the Salter interface in a trivial
// manner.
type TestSalter struct{}

// Salt returns a salt.Salt pointer based on dummy data stored in an in-memory
// storage instance.
func (*TestSalter) Salt(ctx context.Context) (*salt.Salt, error) {
	inmemStorage := &logical.InmemStorage{}
	inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})

	return salt.NewSalt(context.Background(), inmemStorage, &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
	})
}

func TestHashString(t *testing.T) {
	salter := &TestSalter{}

	out, err := HashString(context.Background(), salter, "foo")
	if err != nil {
		t.Fatalf("Error instantiating salt: %s", err)
	}
	if out != "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a" {
		t.Fatalf("err: HashString output did not match expected")
	}
}

func TestHashAuth(t *testing.T) {
	cases := []struct {
		Input        *logical.Auth
		Output       *logical.Auth
		HMACAccessor bool
	}{
		{
			&logical.Auth{ClientToken: "foo"},
			&logical.Auth{ClientToken: "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a"},
			false,
		},
		{
			&logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL: 1 * time.Hour,
				},

				ClientToken: "foo",
			},
			&logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL: 1 * time.Hour,
				},

				ClientToken: "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a",
			},
			false,
		},
	}

	inmemStorage := &logical.InmemStorage{}
	inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	salter := &TestSalter{}
	for _, tc := range cases {
		input := fmt.Sprintf("%#v", tc.Input)
		out, err := HashAuth(context.Background(), salter, tc.Input, tc.HMACAccessor)
		if err != nil {
			t.Fatalf("err: %s\n\n%s", err, input)
		}
		if !reflect.DeepEqual(out, tc.Output) {
			t.Fatalf("bad:\nInput:\n%s\nOutput:\n%#v\nExpected output:\n%#v", input, out, tc.Output)
		}
	}
}

type testOptMarshaler struct {
	S string
	I int
}

func (o *testOptMarshaler) MarshalJSONWithOptions(options *logical.MarshalOptions) ([]byte, error) {
	return json.Marshal(&testOptMarshaler{S: options.ValueHasher(o.S), I: o.I})
}

var _ logical.OptMarshaler = &testOptMarshaler{}

func TestHashRequest(t *testing.T) {
	cases := []struct {
		Input           *logical.Request
		Output          *logical.Request
		NonHMACDataKeys []string
		HMACAccessor    bool
	}{
		{
			&logical.Request{
				Data: map[string]interface{}{
					"foo":              "bar",
					"baz":              "foobar",
					"private_key_type": certutil.PrivateKeyType("rsa"),
					"om":               &testOptMarshaler{S: "bar", I: 1},
				},
			},
			&logical.Request{
				Data: map[string]interface{}{
					"foo":              "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
					"baz":              "foobar",
					"private_key_type": "hmac-sha256:995230dca56fffd310ff591aa404aab52b2abb41703c787cfa829eceb4595bf1",
					"om":               json.RawMessage(`{"S":"hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317","I":1}`),
				},
			},
			[]string{"baz"},
			false,
		},
	}

	inmemStorage := &logical.InmemStorage{}
	inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	salter := &TestSalter{}
	for _, tc := range cases {
		input := fmt.Sprintf("%#v", tc.Input)
		out, err := HashRequest(context.Background(), salter, tc.Input, tc.HMACAccessor, tc.NonHMACDataKeys)
		if err != nil {
			t.Fatalf("err: %s\n\n%s", err, input)
		}
		if diff := deep.Equal(out, tc.Output); len(diff) > 0 {
			t.Fatalf("bad:\nInput:\n%s\nDiff:\n%#v", input, diff)
		}
	}
}

func TestHashResponse(t *testing.T) {
	now := time.Now()

	cases := []struct {
		Input           *logical.Response
		Output          *logical.Response
		NonHMACDataKeys []string
		HMACAccessor    bool
	}{
		{
			&logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
					"baz": "foobar",
					// Responses can contain time values, so test that with
					// a known fixed value.
					"bar": now,
					"om":  &testOptMarshaler{S: "bar", I: 1},
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
					"baz": "foobar",
					"bar": now.Format(time.RFC3339Nano),
					"om":  json.RawMessage(`{"S":"hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317","I":1}`),
				},
				WrapInfo: &wrapping.ResponseWrapInfo{
					TTL:             60,
					Token:           "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
					Accessor:        "hmac-sha256:7c9c6fe666d0af73b3ebcfbfabe6885015558213208e6635ba104047b22f6390",
					CreationTime:    now,
					WrappedAccessor: "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
				},
			},
			[]string{"baz"},
			true,
		},
	}

	inmemStorage := &logical.InmemStorage{}
	inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	salter := &TestSalter{}
	for _, tc := range cases {
		input := fmt.Sprintf("%#v", tc.Input)
		out, err := HashResponse(context.Background(), salter, tc.Input, tc.HMACAccessor, tc.NonHMACDataKeys, false)
		if err != nil {
			t.Fatalf("err: %s\n\n%s", err, input)
		}
		if diff := deep.Equal(out, tc.Output); len(diff) > 0 {
			t.Fatalf("bad:\nInput:\n%s\nDiff:\n%#v", input, diff)
		}
	}
}

func TestHashWalker(t *testing.T) {
	replaceText := "foo"

	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
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
		err := HashStructure(tc.Input, func(string) string {
			return replaceText
		}, nil)
		if err != nil {
			t.Fatalf("err: %s\n\n%#v", err, tc.Input)
		}
		if !reflect.DeepEqual(tc.Input, tc.Output) {
			t.Fatalf("bad:\n\n%#v\n\n%#v", tc.Input, tc.Output)
		}
	}
}

func TestHashWalker_TimeStructs(t *testing.T) {
	replaceText := "bar"

	now := time.Now()
	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
	}{
		// Should not touch map keys of type time.Time.
		{
			map[string]interface{}{
				"hello": map[time.Time]struct{}{
					now: {},
				},
			},
			map[string]interface{}{
				"hello": map[time.Time]struct{}{
					now: {},
				},
			},
		},
		// Should handle map values of type time.Time.
		{
			map[string]interface{}{
				"hello": now,
			},
			map[string]interface{}{
				"hello": now.Format(time.RFC3339Nano),
			},
		},
		// Should handle slice values of type time.Time.
		{
			map[string]interface{}{
				"hello": []interface{}{"foo", now, "foo2"},
			},
			map[string]interface{}{
				"hello": []interface{}{"foobar", now.Format(time.RFC3339Nano), "foo2bar"},
			},
		},
	}

	for _, tc := range cases {
		err := HashStructure(tc.Input, func(s string) string {
			return s + replaceText
		}, nil)
		if err != nil {
			t.Fatalf("err: %v\n\n%#v", err, tc.Input)
		}
		if !reflect.DeepEqual(tc.Input, tc.Output) {
			t.Fatalf("bad:\n\n%#v\n\n%#v", tc.Input, tc.Output)
		}
	}
}
