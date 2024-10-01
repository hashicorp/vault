// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDatabaseErrorSanitizerMiddleware(t *testing.T) {
	type testCase struct {
		inputErr    error
		secretsFunc func() map[string]string

		expectedError error
	}

	tests := map[string]testCase{
		"nil error": {
			inputErr:      nil,
			expectedError: nil,
		},
		"url error": {
			inputErr:      new(url.Error),
			expectedError: errors.New("unable to parse connection url"),
		},
		"nil secrets func": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			expectedError: errors.New("here is my password: iofsd9473tg"),
		},
		"secrets with empty string": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			secretsFunc:   secretFunc(t, "", ""),
			expectedError: errors.New("here is my password: iofsd9473tg"),
		},
		"secrets that do not match": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			secretsFunc:   secretFunc(t, "asdf", "<redacted>"),
			expectedError: errors.New("here is my password: iofsd9473tg"),
		},
		"secrets that do match": {
			inputErr:      errors.New("here is my password: iofsd9473tg"),
			secretsFunc:   secretFunc(t, "iofsd9473tg", "<redacted>"),
			expectedError: errors.New("here is my password: <redacted>"),
		},
		"multiple secrets": {
			inputErr: errors.New("here is my password: iofsd9473tg"),
			secretsFunc: secretFunc(t,
				"iofsd9473tg", "<redacted>",
				"password", "<this was the word password>",
			),
			expectedError: errors.New("here is my <this was the word password>: <redacted>"),
		},
		"gRPC status error": {
			inputErr:      status.Error(codes.InvalidArgument, "an error with a password iofsd9473tg"),
			secretsFunc:   secretFunc(t, "iofsd9473tg", "<redacted>"),
			expectedError: status.Errorf(codes.InvalidArgument, "an error with a password <redacted>"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db := fakeDatabase{}
			mw := NewDatabaseErrorSanitizerMiddleware(db, test.secretsFunc)

			actualErr := mw.sanitize(test.inputErr)
			if !reflect.DeepEqual(actualErr, test.expectedError) {
				t.Fatalf("Actual error: %s\nExpected error: %s", actualErr, test.expectedError)
			}
		})
	}

	t.Run("Initialize", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				initErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.Initialize(context.Background(), InitializeRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 1)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("NewUser", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				newUserErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.NewUser(context.Background(), NewUserRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 1)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				updateUserErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.UpdateUser(context.Background(), UpdateUserRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 1)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				deleteUserErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.DeleteUser(context.Background(), DeleteUserRequest{})
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 1)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Type", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				typeErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		_, err := mw.Type()
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 1)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Close", func(t *testing.T) {
		db := &recordingDatabase{
			next: fakeDatabase{
				closeErr: errors.New("password: iofsd9473tg with some stuff after it"),
			},
		}
		mw := DatabaseErrorSanitizerMiddleware{
			next:      db,
			secretsFn: secretFunc(t, "iofsd9473tg", "<redacted>"),
		}

		expectedErr := errors.New("password: <redacted> with some stuff after it")

		err := mw.Close()
		if !reflect.DeepEqual(err, expectedErr) {
			t.Fatalf("Actual err: %s\n Expected err: %s", err, expectedErr)
		}

		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 1)
	})
}

func secretFunc(t *testing.T, vals ...string) func() map[string]string {
	t.Helper()
	if len(vals)%2 != 0 {
		t.Fatalf("Test configuration error: secretFunc must be called with an even number of values")
	}

	m := map[string]string{}

	for i := 0; i < len(vals); i += 2 {
		key := vals[i]
		m[key] = vals[i+1]
	}

	return func() map[string]string {
		return m
	}
}

func TestTracingMiddleware(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.Initialize(context.Background(), InitializeRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 1)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("NewUser", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.NewUser(context.Background(), NewUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 1)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.UpdateUser(context.Background(), UpdateUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 1)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.DeleteUser(context.Background(), DeleteUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 1)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Type", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		_, err := mw.Type()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 1)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Close", func(t *testing.T) {
		db := &recordingDatabase{}
		logger := hclog.NewNullLogger()
		mw := databaseTracingMiddleware{
			next:   db,
			logger: logger,
		}
		err := mw.Close()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 1)
	})
}

func TestMetricsMiddleware(t *testing.T) {
	t.Run("Initialize", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.Initialize(context.Background(), InitializeRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 1)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("NewUser", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.NewUser(context.Background(), NewUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 1)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.UpdateUser(context.Background(), UpdateUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 1)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.DeleteUser(context.Background(), DeleteUserRequest{})
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 1)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Type", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		_, err := mw.Type()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 1)
		assertEquals(t, db.closeCalls, 0)
	})

	t.Run("Close", func(t *testing.T) {
		db := &recordingDatabase{}
		mw := databaseMetricsMiddleware{
			next:    db,
			typeStr: "metrics",
		}
		err := mw.Close()
		if err != nil {
			t.Fatalf("Expected no error, but got: %s", err)
		}
		assertEquals(t, db.initializeCalls, 0)
		assertEquals(t, db.newUserCalls, 0)
		assertEquals(t, db.updateUserCalls, 0)
		assertEquals(t, db.deleteUserCalls, 0)
		assertEquals(t, db.typeCalls, 0)
		assertEquals(t, db.closeCalls, 1)
	})
}

func assertEquals(t *testing.T, actual, expected int) {
	t.Helper()
	if actual != expected {
		t.Fatalf("Actual: %d Expected: %d", actual, expected)
	}
}
