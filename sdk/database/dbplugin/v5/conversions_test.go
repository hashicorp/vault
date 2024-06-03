// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/hashicorp/vault/sdk/database/dbplugin/v5/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConversionsHaveAllFields(t *testing.T) {
	t.Run("initReqToProto", func(t *testing.T) {
		req := InitializeRequest{
			Config: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": "baz",
				},
			},
			VerifyConnection: true,
		}

		protoReq, err := initReqToProto(req)
		if err != nil {
			t.Fatalf("Failed to convert request to proto request: %s", err)
		}

		values := getAllGetterValues(protoReq)
		if len(values) == 0 {
			// Probably a test failure - the protos used in these tests should have Get functions on them
			t.Fatalf("No values found from Get functions!")
		}

		for _, gtr := range values {
			err := assertAllFieldsSet(fmt.Sprintf("InitializeRequest.%s", gtr.name), gtr.value)
			if err != nil {
				t.Fatalf("%s", err)
			}
		}
	})

	t.Run("newUserReqToProto", func(t *testing.T) {
		req := NewUserRequest{
			UsernameConfig: UsernameMetadata{
				DisplayName: "dispName",
				RoleName:    "roleName",
			},
			Statements: Statements{
				Commands: []string{
					"statement",
				},
			},
			RollbackStatements: Statements{
				Commands: []string{
					"rollback_statement",
				},
			},
			CredentialType: CredentialTypeRSAPrivateKey,
			PublicKey:      []byte("-----BEGIN PUBLIC KEY-----"),
			Password:       "password",
			Subject:        "subject",
			Expiration:     time.Now(),
		}

		protoReq, err := newUserReqToProto(req)
		if err != nil {
			t.Fatalf("Failed to convert request to proto request: %s", err)
		}

		values := getAllGetterValues(protoReq)
		if len(values) == 0 {
			// Probably a test failure - the protos used in these tests should have Get functions on them
			t.Fatalf("No values found from Get functions!")
		}

		for _, gtr := range values {
			err := assertAllFieldsSet(fmt.Sprintf("NewUserRequest.%s", gtr.name), gtr.value)
			if err != nil {
				t.Fatalf("%s", err)
			}
		}
	})

	t.Run("updateUserReqToProto", func(t *testing.T) {
		req := UpdateUserRequest{
			Username:       "username",
			CredentialType: CredentialTypeRSAPrivateKey,
			Password: &ChangePassword{
				NewPassword: "newpassword",
				Statements: Statements{
					Commands: []string{
						"statement",
					},
				},
			},
			PublicKey: &ChangePublicKey{
				NewPublicKey: []byte("-----BEGIN PUBLIC KEY-----"),
				Statements: Statements{
					Commands: []string{
						"statement",
					},
				},
			},
			Expiration: &ChangeExpiration{
				NewExpiration: time.Now(),
				Statements: Statements{
					Commands: []string{
						"statement",
					},
				},
			},
		}

		protoReq, err := updateUserReqToProto(req)
		if err != nil {
			t.Fatalf("Failed to convert request to proto request: %s", err)
		}

		values := getAllGetterValues(protoReq)
		if len(values) == 0 {
			// Probably a test failure - the protos used in these tests should have Get functions on them
			t.Fatalf("No values found from Get functions!")
		}

		for _, gtr := range values {
			err := assertAllFieldsSet(fmt.Sprintf("UpdateUserRequest.%s", gtr.name), gtr.value)
			if err != nil {
				t.Fatalf("%s", err)
			}
		}
	})

	t.Run("deleteUserReqToProto", func(t *testing.T) {
		req := DeleteUserRequest{
			Username: "username",
			Statements: Statements{
				Commands: []string{
					"statement",
				},
			},
		}

		protoReq, err := deleteUserReqToProto(req)
		if err != nil {
			t.Fatalf("Failed to convert request to proto request: %s", err)
		}

		values := getAllGetterValues(protoReq)
		if len(values) == 0 {
			// Probably a test failure - the protos used in these tests should have Get functions on them
			t.Fatalf("No values found from Get functions!")
		}

		for _, gtr := range values {
			err := assertAllFieldsSet(fmt.Sprintf("DeleteUserRequest.%s", gtr.name), gtr.value)
			if err != nil {
				t.Fatalf("%s", err)
			}
		}
	})

	t.Run("getUpdateUserRequest", func(t *testing.T) {
		req := &proto.UpdateUserRequest{
			Username:       "username",
			CredentialType: int32(CredentialTypeRSAPrivateKey),
			Password: &proto.ChangePassword{
				NewPassword: "newpass",
				Statements: &proto.Statements{
					Commands: []string{
						"statement",
					},
				},
			},
			PublicKey: &proto.ChangePublicKey{
				NewPublicKey: []byte("-----BEGIN PUBLIC KEY-----"),
				Statements: &proto.Statements{
					Commands: []string{
						"statement",
					},
				},
			},
			Expiration: &proto.ChangeExpiration{
				NewExpiration: timestamppb.Now(),
				Statements: &proto.Statements{
					Commands: []string{
						"statement",
					},
				},
			},
		}

		protoReq, err := getUpdateUserRequest(req)
		if err != nil {
			t.Fatalf("Failed to convert request to proto request: %s", err)
		}

		err = assertAllFieldsSet("proto.UpdateUserRequest", protoReq)
		if err != nil {
			t.Fatalf("%s", err)
		}
	})
}

type getter struct {
	name  string
	value interface{}
}

func getAllGetterValues(value interface{}) (values []getter) {
	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if !strings.HasPrefix(method.Name, "Get") {
			continue
		}
		valMethod := val.Method(i)
		resp := valMethod.Call(nil)
		getVal := resp[0].Interface()
		gtr := getter{
			name:  strings.TrimPrefix(method.Name, "Get"),
			value: getVal,
		}
		values = append(values, gtr)
	}
	return values
}

// Ensures the assertion works properly
func TestAssertAllFieldsSet(t *testing.T) {
	type testCase struct {
		value     interface{}
		expectErr bool
	}

	tests := map[string]testCase{
		"zero int": {
			value:     0,
			expectErr: true,
		},
		"non-zero int": {
			value:     1,
			expectErr: false,
		},
		"zero float64": {
			value:     0.0,
			expectErr: true,
		},
		"non-zero float64": {
			value:     1.0,
			expectErr: false,
		},
		"empty string": {
			value:     "",
			expectErr: true,
		},
		"true boolean": {
			value:     true,
			expectErr: false,
		},
		"false boolean": { // False is an exception to the "is zero" rule
			value:     false,
			expectErr: false,
		},
		"blank struct": {
			value:     struct{}{},
			expectErr: true,
		},
		"non-blank but empty struct": {
			value: struct {
				str string
			}{
				str: "",
			},
			expectErr: true,
		},
		"non-empty string": {
			value:     "foo",
			expectErr: false,
		},
		"non-empty struct": {
			value: struct {
				str string
			}{
				str: "foo",
			},
			expectErr: false,
		},
		"empty nested struct": {
			value: struct {
				Str       string
				Substruct struct {
					Substr string
				}
			}{
				Str: "foo",
				Substruct: struct {
					Substr string
				}{}, // Empty sub-field
			},
			expectErr: true,
		},
		"filled nested struct": {
			value: struct {
				str       string
				substruct struct {
					substr string
				}
			}{
				str: "foo",
				substruct: struct {
					substr string
				}{
					substr: "sub-foo",
				},
			},
			expectErr: false,
		},
		"nil map": {
			value:     map[string]string(nil),
			expectErr: true,
		},
		"empty map": {
			value:     map[string]string{},
			expectErr: true,
		},
		"filled map": {
			value: map[string]string{
				"foo": "bar",
				"int": "42",
			},
			expectErr: false,
		},
		"map with empty string value": {
			value: map[string]string{
				"foo": "",
			},
			expectErr: true,
		},
		"nested map with empty string value": {
			value: map[string]interface{}{
				"bar": "baz",
				"foo": map[string]interface{}{
					"subfoo": "",
				},
			},
			expectErr: true,
		},
		"nil slice": {
			value:     []string(nil),
			expectErr: true,
		},
		"empty slice": {
			value:     []string{},
			expectErr: true,
		},
		"filled slice": {
			value: []string{
				"foo",
			},
			expectErr: false,
		},
		"slice with empty string value": {
			value: []string{
				"",
			},
			expectErr: true,
		},
		"empty structpb": {
			value:     newStructPb(t, map[string]interface{}{}),
			expectErr: true,
		},
		"filled structpb": {
			value: newStructPb(t, map[string]interface{}{
				"foo": "bar",
				"int": 42,
			}),
			expectErr: false,
		},

		"pointer to zero int": {
			value:     intPtr(0),
			expectErr: true,
		},
		"pointer to non-zero int": {
			value:     intPtr(1),
			expectErr: false,
		},
		"pointer to zero float64": {
			value:     float64Ptr(0.0),
			expectErr: true,
		},
		"pointer to non-zero float64": {
			value:     float64Ptr(1.0),
			expectErr: false,
		},
		"pointer to nil string": {
			value:     new(string),
			expectErr: true,
		},
		"pointer to non-nil string": {
			value:     strPtr("foo"),
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := assertAllFieldsSet("", test.value)
			if test.expectErr && err == nil {
				t.Fatalf("err expected, got nil")
			}
			if !test.expectErr && err != nil {
				t.Fatalf("no error expected, got: %s", err)
			}
		})
	}
}

func assertAllFieldsSet(name string, val interface{}) error {
	if val == nil {
		return fmt.Errorf("value is nil")
	}

	rVal := reflect.ValueOf(val)
	return assertAllFieldsSetValue(name, rVal)
}

func assertAllFieldsSetValue(name string, rVal reflect.Value) error {
	// All booleans are allowed - we don't have a way of differentiating between
	// and intentional false and a missing false
	if rVal.Kind() == reflect.Bool {
		return nil
	}

	// Primitives fall through here
	if rVal.IsZero() {
		return fmt.Errorf("%s is zero", name)
	}

	switch rVal.Kind() {
	case reflect.Ptr, reflect.Interface:
		return assertAllFieldsSetValue(name, rVal.Elem())
	case reflect.Struct:
		return assertAllFieldsSetStruct(name, rVal)
	case reflect.Map:
		if rVal.Len() == 0 {
			return fmt.Errorf("%s (map type) is empty", name)
		}

		iter := rVal.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()

			err := assertAllFieldsSetValue(fmt.Sprintf("%s[%s]", name, k), v)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		if rVal.Len() == 0 {
			return fmt.Errorf("%s (slice type) is empty", name)
		}
		for i := 0; i < rVal.Len(); i++ {
			sliceVal := rVal.Index(i)
			err := assertAllFieldsSetValue(fmt.Sprintf("%s[%d]", name, i), sliceVal)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func assertAllFieldsSetStruct(name string, rVal reflect.Value) error {
	switch rVal.Type() {
	case reflect.TypeOf(timestamppb.Timestamp{}):
		ts := rVal.Interface().(timestamppb.Timestamp)
		if ts.AsTime().IsZero() {
			return fmt.Errorf("%s is zero", name)
		}
		return nil
	default:
		for i := 0; i < rVal.NumField(); i++ {
			field := rVal.Field(i)
			fieldName := rVal.Type().Field(i)

			// Skip fields that aren't exported
			if unicode.IsLower([]rune(fieldName.Name)[0]) {
				continue
			}

			err := assertAllFieldsSetValue(fmt.Sprintf("%s.%s", name, fieldName.Name), field)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func strPtr(str string) *string {
	return &str
}

func newStructPb(t *testing.T, m map[string]interface{}) *structpb.Struct {
	t.Helper()

	s, err := structpb.NewStruct(m)
	if err != nil {
		t.Fatalf("Failed to convert map to struct: %s", err)
	}
	return s
}
