// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	nshelper "github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
	"github.com/stretchr/testify/require"
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

// testSalter is a structure that implements the Salter interface in a trivial
// manner.
type testSalter struct{}

// Salt returns a salt.Salt pointer based on dummy data stored in an in-memory
// storage instance.
func (*testSalter) Salt(ctx context.Context) (*salt.Salt, error) {
	inmemStorage := &logical.InmemStorage{}
	err := inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	if err != nil {
		return nil, err
	}

	return salt.NewSalt(context.Background(), inmemStorage, &salt.Config{
		HMAC:     sha256.New,
		HMACType: "hmac-sha256",
	})
}

func TestHashString(t *testing.T) {
	salter := &testSalter{}

	out, err := hashString(context.Background(), salter, "foo")
	if err != nil {
		t.Fatalf("Error instantiating salt: %s", err)
	}
	if out != "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a" {
		t.Fatalf("err: hashString output did not match expected")
	}
}

func TestHashAuth(t *testing.T) {
	cases := map[string]struct {
		Input        *logical.Auth
		Output       *auth
		HMACAccessor bool
	}{
		"no-accessor-hmac": {
			&logical.Auth{
				ClientToken: "foo",
				Accessor:    "very-accessible",
				LeaseOptions: logical.LeaseOptions{
					TTL: 1 * time.Hour,
				},
				TokenType: logical.TokenTypeService,
			},
			&auth{
				ClientToken:   "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a",
				Accessor:      "very-accessible",
				TokenTTL:      3600,
				TokenType:     "service",
				RemainingUses: 5,
			},
			false,
		},
		"accessor-hmac": {
			&logical.Auth{
				Accessor:    "very-accessible",
				ClientToken: "foo",
				LeaseOptions: logical.LeaseOptions{
					TTL: 1 * time.Hour,
				},
				TokenType: logical.TokenTypeBatch,
			},
			&auth{
				ClientToken:   "hmac-sha256:08ba357e274f528065766c770a639abf6809b39ccfd37c2a3157c7f51954da0a",
				Accessor:      "hmac-sha256:5d6d7c8da5b699ace193ea453bbf77082a8aaca42a474436509487d646a7c0af",
				TokenTTL:      3600,
				TokenType:     "batch",
				RemainingUses: 5,
			},
			true,
		},
	}

	inmemStorage := &logical.InmemStorage{}
	err := inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	require.NoError(t, err)
	salter := &testSalter{}
	for _, tc := range cases {
		auditAuth, err := newAuth(tc.Input, 5)
		require.NoError(t, err)
		err = hashAuth(context.Background(), salter, auditAuth, tc.HMACAccessor)
		require.NoError(t, err)
		require.Equal(t, tc.Output, auditAuth)
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
		Output          *request
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
			&request{
				Data: map[string]interface{}{
					"foo":              "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
					"baz":              "foobar",
					"private_key_type": "hmac-sha256:995230dca56fffd310ff591aa404aab52b2abb41703c787cfa829eceb4595bf1",
					"om":               json.RawMessage(`{"S":"hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317","I":1}`),
				},
				Namespace: &namespace{
					ID:   nshelper.RootNamespace.ID,
					Path: nshelper.RootNamespace.Path,
				},
			},
			[]string{"baz"},
			false,
		},
	}

	inmemStorage := &logical.InmemStorage{}
	err := inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	require.NoError(t, err)
	salter := &testSalter{}
	for _, tc := range cases {
		auditReq, err := newRequest(tc.Input, nshelper.RootNamespace)
		require.NoError(t, err)
		err = hashRequest(context.Background(), salter, auditReq, tc.HMACAccessor, tc.NonHMACDataKeys)
		require.NoError(t, err)
		require.Equal(t, tc.Output, auditReq)
	}
}

func TestHashResponse(t *testing.T) {
	now := time.Now()

	resp := &logical.Response{
		Data: map[string]interface{}{
			"foo": "bar",
			"baz": "foobar",
			// Responses can contain time values, so test that with a known fixed value.
			"bar": now,
			"om":  &testOptMarshaler{S: "bar", I: 1},
		},
		WrapInfo: &wrapping.ResponseWrapInfo{
			TTL:             1 * time.Minute,
			Token:           "bar",
			Accessor:        "flimflam",
			CreationTime:    now,
			WrappedAccessor: "bar",
		},
		Auth: &logical.Auth{
			ClientToken: "hvs.QWERTY-T1q5lEjIWux1Tjx-VGqAYJdd4FZtbp1wpD5Ym9pGh4KHGh2cy5TSjRndGoxaU44NzNscm5MSlRLQXZ0ZGg",
			Accessor:    "ABClk9ZNLGOCuTrOEIAooJG3",
			TokenType:   logical.TokenTypeService,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       3,
				MaxTTL:    5,
				Renewable: false,
				Increment: 1,
				IssueTime: now,
			},
			InternalData: map[string]any{
				"foo": "bar",
			},
			LeaseID: "abc",
		},
	}

	req := &logical.Request{MountPoint: "/foo/bar"}
	req.SetMountClass("kv")
	req.SetMountIsExternalPlugin(true)
	req.SetMountRunningVersion("123")
	req.SetMountRunningSha256("256-256!")

	nonHMACDataKeys := []string{"baz"}

	expected := &response{
		Auth: &auth{
			Accessor:    "hmac-sha256:253184715b2d5a6c3a2fc7afe0d2294085f5e886a1275ca735646a6f23be2587",
			ClientToken: "hmac-sha256:2ce541100a8bcd687e8ec7712c8bb4c975a8d8599c02d98945e63ecd413bf0f3",
			TokenType:   "service",
		},
		Data: map[string]interface{}{
			"foo": "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
			"baz": "foobar",
			"bar": now.Format(time.RFC3339Nano),
			"om":  json.RawMessage(`{"S":"hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317","I":1}`),
		},
		WrapInfo: &responseWrapInfo{
			TTL:             60,
			Token:           "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
			Accessor:        "hmac-sha256:7c9c6fe666d0af73b3ebcfbfabe6885015558213208e6635ba104047b22f6390",
			CreationTime:    now.UTC().Format(time.RFC3339Nano),
			WrappedAccessor: "hmac-sha256:f9320baf0249169e73850cd6156ded0106e2bb6ad8cab01b7bbbebe6d1065317",
		},
		MountClass:            "kv",
		MountIsExternalPlugin: true,
		MountPoint:            "/foo/bar",
		MountRunningVersion:   "123",
		MountRunningSha256:    "256-256!",
		Secret: &secret{
			LeaseID: "abc",
		},
	}

	inmemStorage := &logical.InmemStorage{}
	err := inmemStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   "salt",
		Value: []byte("foo"),
	})
	require.NoError(t, err)
	salter := &testSalter{}
	auditResp, err := newResponse(resp, req, false)
	require.NoError(t, err)
	err = hashResponse(context.Background(), salter, auditResp, true, nonHMACDataKeys)
	require.NoError(t, err)
	require.Equal(t, expected, auditResp)
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
		err := hashStructure(tc.Input, func(string) string {
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
		err := hashStructure(tc.Input, func(s string) string {
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
