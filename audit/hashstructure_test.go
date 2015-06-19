package audit

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/copystructure"
)

func TestCopy_auth(t *testing.T) {
	// Make a non-pointer one so that it can't be modified directly
	expected := logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			Lease:      1 * time.Hour,
			LeaseIssue: time.Now().UTC(),
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

func TestHash(t *testing.T) {
	now := time.Now().UTC()

	cases := []struct {
		Input  interface{}
		Output interface{}
	}{
		{
			&logical.Auth{ClientToken: "foo"},
			&logical.Auth{ClientToken: "sha1:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"},
		},
		{
			&logical.Request{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			&logical.Request{
				Data: map[string]interface{}{
					"foo": "sha1:62cdb7020ff920e5aa642c3d4066950dd1f01f4d",
				},
			},
		},
		{
			&logical.Response{
				Data: map[string]interface{}{
					"foo": "bar",
				},
			},
			&logical.Response{
				Data: map[string]interface{}{
					"foo": "sha1:62cdb7020ff920e5aa642c3d4066950dd1f01f4d",
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
					Lease:      1 * time.Hour,
					LeaseIssue: now,
				},

				ClientToken: "foo",
			},
			&logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					Lease:      1 * time.Hour,
					LeaseIssue: now,
				},

				ClientToken: "sha1:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33",
			},
		},
	}

	for _, tc := range cases {
		input := fmt.Sprintf("%#v", tc.Input)
		if err := Hash(tc.Input); err != nil {
			t.Fatalf("err: %s\n\n%s", err, input)
		}
		if !reflect.DeepEqual(tc.Input, tc.Output) {
			t.Fatalf("bad:\n\n%s\n\n%#v\n\n%#v", input, tc.Input, tc.Output)
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
		output, err := HashStructure(tc.Input, func(string) (string, error) {
			return replaceText, nil
		})
		if err != nil {
			t.Fatalf("err: %s\n\n%#v", err, tc.Input)
		}
		if !reflect.DeepEqual(output, tc.Output) {
			t.Fatalf("bad:\n\n%#v\n\n%#v", tc.Input, output)
		}
	}
}

func TestHashSHA1(t *testing.T) {
	fn := HashSHA1("")
	result, err := fn("foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if result != "sha1:0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33" {
		t.Fatalf("bad: %#v", result)
	}
}
