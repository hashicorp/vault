// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package uicustommessages

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestManagerGetEntryForNamespace verifies the behaviour of the
// (*Manager).getEntryForNamespace method in each of its conditional branches.
// Since the method being tested is a helper used by multiple other methods
// of the Manager struct, this test simplifies their tests by eliminating
// duplicate test cases (e.g. storage get returns an error or contains an
// entry with an invalid value that can't be decoded using JSON).
func TestManagerGetEntryForNamespace(t *testing.T) {
	var (
		testManager = NewManager(nil)

		testNs = namespace.RootNamespace
	)
	for _, testcase := range []struct {
		name           string
		context        context.Context
		storage        logical.Storage
		ns             *namespace.Namespace
		errorAssertion func(assert.TestingT, error, ...any) bool
		entryAssertion func(assert.TestingT, any, ...any) bool
	}{
		{
			name:           "namespace nil",
			context:        context.Background(),
			ns:             nil,
			errorAssertion: assert.Error,
			entryAssertion: assert.Nil,
		},
		{
			name:           "storage get fails",
			context:        context.Background(),
			ns:             testNs,
			storage:        &testingStorage{getFails: true},
			errorAssertion: assert.Error,
			entryAssertion: assert.Nil,
		},
		{
			name:           "entry does not exist",
			context:        context.Background(),
			ns:             testNs,
			errorAssertion: assert.NoError,
			entryAssertion: assert.NotNil,
		},
		{
			name:           "entry exists, invalid",
			context:        context.Background(),
			ns:             testNs,
			storage:        buildStorageWithEntry("sys/config/ui/custom-messages", "}-^"),
			errorAssertion: assert.Error,
			entryAssertion: assert.Nil,
		},
		{
			name:           "entry exists, valid",
			context:        context.Background(),
			ns:             testNs,
			storage:        buildStorageWithEntry("sys/config/ui/custom-messages", `{"messages":{}}`),
			errorAssertion: assert.NoError,
			entryAssertion: assert.NotNil,
		},
	} {
		if testcase.storage != nil {
			testManager.view = testcase.storage
		} else {
			testManager.view = &logical.InmemStorage{}
		}

		entry, err := testManager.getEntryForNamespace(testcase.context, testcase.ns)
		testcase.errorAssertion(t, err, testcase.name)
		testcase.entryAssertion(t, entry, testcase.name)
	}
}

// TestManagerGetEntry verifies the behaviour of the (*Manager).getEntry method
// in each of its conditional branches. Since the method being tested is a
// helper used by multiple other methods of the Manager struct, this test
// simplifies their tests by eliminating duplicate test cases (e.g. error
// retrieving the namespace.Namespace from the context.Context).
func TestManagerGetEntry(t *testing.T) {
	testManager := NewManager(buildStorageWithEntry("root", `{}`))

	entry, err := testManager.getEntry(context.Background())
	assert.Error(t, err)
	assert.Nil(t, entry)

	entry, err = testManager.getEntry(namespace.ContextWithNamespace(context.Background(), &namespace.Namespace{ID: "root"}))
	assert.NoError(t, err)
	assert.NotNil(t, entry)
}

// TestManagerPutEntry verifies the behaviour of the (*Manager).putEntry method
// in each of its conditional branches (except for errors returned from
// json.Marshal since no possible Entry struct can be provided to cause one).
// Since the method being tested is a helper used by multiple other methods of
// the Manager struct, this test simplifies their tests by eliminating duplicate
// test cases (e.g. storage put errors or failure to retrieve namespace from
// the context).
func TestManagerPutEntry(t *testing.T) {
	var (
		testManager = NewManager(nil)
		testEntry   = &Entry{
			Messages: make(map[string]Message),
		}
		nsCtx = namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	)

	for _, testcase := range []struct {
		name           string
		context        context.Context
		storage        logical.Storage
		errorAssertion func(assert.TestingT, error, ...any) bool
	}{
		{
			name:           "fail to extract namespace from context",
			context:        nil,
			errorAssertion: assert.Error,
		},
		{
			name:           "storage put fails",
			context:        nsCtx,
			storage:        &testingStorage{putFails: true},
			errorAssertion: assert.Error,
		},
		{
			name:           "no errors",
			context:        nsCtx,
			errorAssertion: assert.NoError,
		},
	} {
		switch testcase.storage {
		case nil:
			testManager.view = &logical.InmemStorage{}
		default:
			testManager.view = testcase.storage
		}

		testcase.errorAssertion(t, testManager.putEntry(testcase.context, testEntry), testcase.name)
	}
}

// TestGetNamespacesToSearch verifies the behaviour of the getNamespacesToSearch
// function in each of its conditional branches. Since the function being tested
// is a helper used by the (*Manager).FindMessages method, this test simplifies
// its test by making more focused assertions here without all of the additional
// context (e.g. checking that the list contains 1 element and that it's equal
// to namespace.RootNamespace).
func TestGetNamespacesToSearch(t *testing.T) {
	list, err := getNamespacesToSearch(context.Background(), FindFilter{})
	assert.Error(t, err)
	assert.Nil(t, list)

	list, err = getNamespacesToSearch(namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace), FindFilter{})
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 1, len(list))
	assert.Equal(t, namespace.RootNamespace, list[0])
}

// TestStorageKeyForNamespace verifies that the storageKeyForNamespace function
// returns sys/config/ui/custom-messages when the provided namespace is the root
// namespace, otherwise it returns
// namespaces/<ns.id>/sys/config/ui/custom-messages.
func TestStorageKeyForNamespace(t *testing.T) {
	// Check for root namespace
	assert.Equal(t, "sys/config/ui/custom-messages", storageKeyForNamespace(*namespace.RootNamespace))

	// Check for a non-root namespace
	assert.Equal(t, "namespaces/test/sys/config/ui/custom-messages", storageKeyForNamespace(namespace.Namespace{ID: "test", Path: "test/"}))
}

// TestManagerFindMessages verifies the behaviour of the (*Manager).FindMessages
// method in each of its conditional branches.
func TestManagerFindMessages(t *testing.T) {
	var (
		testManager = NewManager(nil)

		nsCtx = namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

		// invalidEntryStorage = buildStorageWithEntry("root", "}^-")
		// validEntryStorage   = buildStorageWithEntry("root", `{"messages":{"abc":{"id":"abc"}}}`)
	)

	for _, testcase := range []struct {
		name            string
		context         context.Context
		filters         FindFilter
		storage         logical.Storage
		errorAssertion  func(assert.TestingT, error, ...any) bool
		resultAssertion func(assert.TestingT, any, ...any) bool
	}{
		{
			name:            "no namespaces to search",
			context:         nil,
			errorAssertion:  assert.Error,
			resultAssertion: assert.Nil,
		},
		{
			name:            "fail to get entry for namespace",
			context:         nsCtx,
			storage:         &testingStorage{getFails: true},
			errorAssertion:  assert.Error,
			resultAssertion: assert.Nil,
		},
		{
			name:            "valid storageEntry",
			context:         nsCtx,
			storage:         buildStorageWithEntry("root", `{"messages":{}}`),
			errorAssertion:  assert.NoError,
			resultAssertion: assert.NotNil,
		},
	} {
		switch testcase.storage {
		case nil:
			testManager.view = &logical.InmemStorage{}
		default:
			testManager.view = testcase.storage
		}

		messages, err := testManager.FindMessages(testcase.context, testcase.filters)

		testcase.errorAssertion(t, err, testcase.name)
		testcase.resultAssertion(t, messages, testcase.name)
	}
}

// TestManagerCreateMessage verifies the behaviour of the
// (*Manager).CreateMessage method in each of its conditional branches.
func TestManagerCreateMessage(t *testing.T) {
	var (
		testManager = NewManager(nil)

		validMessageTpl = Message{
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
			Message:   "created message",
			Type:      "banner",
		}
		invalidMessageTpl = Message{
			StartTime: time.Now().Add(time.Hour),
			EndTime:   time.Now(),
		}

		nsCtx = namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	)

	for _, testcase := range []struct {
		name             string
		storage          logical.Storage
		message          Message
		errorAssertion   func(assert.TestingT, error, ...any) bool
		messageAssertion func(assert.TestingT, any, ...any) bool
	}{
		{
			name:             "storage fail to get entry",
			storage:          &testingStorage{getFails: true},
			message:          Message{},
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:             "entry fail to create message",
			message:          invalidMessageTpl,
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:             "storage fail to put entry",
			storage:          &testingStorage{putFails: true, getResponseValue: `{"messages":{}}`},
			message:          validMessageTpl,
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:             "message created",
			message:          validMessageTpl,
			errorAssertion:   assert.NoError,
			messageAssertion: assert.NotNil,
		},
	} {
		switch testcase.storage {
		case nil:
			testManager.view = &logical.InmemStorage{}
		default:
			testManager.view = testcase.storage
		}

		message, err := testManager.CreateMessage(nsCtx, testcase.message)
		testcase.errorAssertion(t, err, testcase.name)
		testcase.messageAssertion(t, message, testcase.name)

		if message != nil {
			assert.NotEmpty(t, message.ID, testcase.name)

			entry, err := testManager.getEntry(nsCtx)
			require.NoError(t, err, testcase.name)
			require.NotNil(t, entry, testcase.name)

			entryMessage, ok := entry.Messages[message.ID]
			assert.True(t, ok, testcase.name)
			assert.Equal(t, testcase.message.Message, entryMessage.Message, testcase.name)
			assert.Equal(t, message.ID, entryMessage.ID, testcase.name)
		}
	}
}

// TestManagerReadMessage verifies the behaviour of the (*Manager).ReadMessage
// method in each of its conditional branches.
func TestManagerReadMessage(t *testing.T) {
	var (
		testManager = NewManager(nil)

		nsCtx = namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	)

	for _, testcase := range []struct {
		name             string
		storage          logical.Storage
		messageID        string
		errorAssertion   func(assert.TestingT, error, ...any) bool
		errorIsAssertion func(assert.TestingT, error, error, ...any) bool
		messageAssertion func(assert.TestingT, any, ...any) bool
	}{
		{
			name:             "storage get fails",
			storage:          &testingStorage{getFails: true},
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:             "message does not exist",
			storage:          &logical.InmemStorage{},
			errorAssertion:   assert.Error,
			errorIsAssertion: assert.ErrorIs,
			messageAssertion: assert.Nil,
		},
		{
			name:             "message exists",
			storage:          buildStorageWithEntry("sys/config/ui/custom-messages", `{"messages":{"abc":{"id":"abc"}}}`),
			messageID:        "abc",
			errorAssertion:   assert.NoError,
			messageAssertion: assert.NotNil,
		},
	} {
		testManager.view = testcase.storage

		message, err := testManager.ReadMessage(nsCtx, testcase.messageID)
		testcase.errorAssertion(t, err, testcase.name)
		if testcase.errorIsAssertion != nil {
			testcase.errorIsAssertion(t, err, logical.ErrNotFound, testcase.name)
		}
		testcase.messageAssertion(t, message, testcase.name)
	}
}

// TestManagerUpdateMessage verifies the behaviour of the
// (*Manager).UpdateMessage method in each of its conditional branches.
func TestManagerUpdateMessage(t *testing.T) {
	var (
		testManager = NewManager(nil)

		nsCtx = namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	)

	for _, testcase := range []struct {
		name             string
		message          Message
		storage          logical.Storage
		errorAssertion   func(assert.TestingT, error, ...any) bool
		messageAssertion func(assert.TestingT, any, ...any) bool
	}{
		{
			name:             "storage get fails",
			storage:          &testingStorage{getFails: true},
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:    "updating to invalid times",
			storage: buildStorageWithEntry("sys/config/ui/custom-messages", `{"messages":{"abc":{"id":"abc"}}}`),
			message: Message{
				ID:        "abc",
				StartTime: time.Now().Add(time.Hour),
				EndTime:   time.Now().Add(-1 * time.Hour),
			},
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:    "updating non-existant message",
			storage: &logical.InmemStorage{},
			message: Message{
				ID:        "abc",
				StartTime: time.Now().Add(-1 * time.Hour),
				EndTime:   time.Now().Add(3 * time.Hour),
			},
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:    "storage put fails",
			storage: &testingStorage{putFails: true, getResponseValue: `{"messages":{"abc":{"id":"abc"}}}`},
			message: Message{
				ID:        "abc",
				StartTime: time.Now().Add(-5 * time.Hour),
				EndTime:   time.Now().Add(time.Hour),
				Type:      "banner",
			},
			errorAssertion:   assert.Error,
			messageAssertion: assert.Nil,
		},
		{
			name:    "message updated",
			storage: buildStorageWithEntry("sys/config/ui/custom-messages", `{"messages":{"abc":{"id":"abc"}}}`),
			message: Message{
				ID:        "abc",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
				Message:   "updated value",
				Type:      "banner",
			},
			errorAssertion:   assert.NoError,
			messageAssertion: assert.NotNil,
		},
	} {
		testManager.view = testcase.storage

		message, err := testManager.UpdateMessage(nsCtx, testcase.message)
		testcase.errorAssertion(t, err, testcase.name)
		testcase.messageAssertion(t, message, testcase.name)

		if message != nil {
			entry, err := testManager.getEntry(nsCtx)
			require.NoError(t, err, testcase.name)
			require.NotNil(t, entry, testcase.name)

			entryMessage, ok := entry.Messages[testcase.message.ID]
			assert.True(t, ok, testcase.name)
			assert.Equal(t, testcase.message.Message, entryMessage.Message, testcase.name)
		}
	}
}

// TestManagerDeleteMessage verifies the behaviour of the
// (*Manager).DeleteMessage method in each of its conditional branches.
func TestManagerDeleteMessage(t *testing.T) {
	var (
		testManager = NewManager(nil)

		nsCtx = namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)
	)

	for _, testcase := range []struct {
		name           string
		messageID      string
		storage        logical.Storage
		errorAssertion func(assert.TestingT, error, ...any) bool
		checkStorage   bool
	}{
		{
			name:           "storage get fails",
			storage:        &testingStorage{getFails: true},
			errorAssertion: assert.Error,
		},
		{
			name:           "storage put fails",
			storage:        &testingStorage{putFails: true, getResponseValue: `{"messages":{}}`},
			errorAssertion: assert.Error,
		},
		{
			name:           "message deleted",
			storage:        buildStorageWithEntry("root", `{"messages":{"abc":{"id":"abc"}}}`),
			messageID:      "abc",
			errorAssertion: assert.NoError,
			checkStorage:   true,
		},
	} {
		testManager.view = testcase.storage
		testcase.errorAssertion(t, testManager.DeleteMessage(nsCtx, testcase.messageID), testcase.name)

		if testcase.checkStorage {
			entry, err := testManager.getEntry(nsCtx)
			require.NoError(t, err, testcase.name)
			require.NotNil(t, entry, testcase.name)

			assert.NotContains(t, entry.Messages, testcase.messageID)
		}
	}
}

// buildStorageWithEntry is a helper function that returns a logical.Storage
// with a logical.StorageEntry built using the key and value arguments stored
// in it.
func buildStorageWithEntry(key, value string) logical.Storage {
	storage := &logical.InmemStorage{}

	entry := &logical.StorageEntry{
		Key:   key,
		Value: []byte(value),
	}

	storage.Put(context.Background(), entry)

	return storage
}

// testingStorage is a struct that implements the logical.Storage interface
// that's used to simulate errors occurring when interacting with the interface.
// Each of the methods will return an error if their correspond *Fails field is
// set to true. Otherwise, the *ResponseValue fields can be used to specify
// what the corresponding method should return.
type testingStorage struct {
	listFails   bool
	getFails    bool
	deleteFails bool
	putFails    bool

	listResponseValue []string
	getResponseValue  string
}

// List fails if s.listFails is true, otherwise it returns s.listResponseValue.
func (s *testingStorage) List(_ context.Context, key string) ([]string, error) {
	if s.listFails {
		return nil, errors.New("failure")
	}

	if s.listResponseValue == nil {
		s.listResponseValue = make([]string, 0)
	}

	return s.listResponseValue, nil
}

// Get fails if s.getFails is true, otherwise it returns s.getResponseValue.
func (s *testingStorage) Get(_ context.Context, key string) (*logical.StorageEntry, error) {
	if s.getFails {
		return nil, errors.New("failure")
	}

	return &logical.StorageEntry{
		Key:   key,
		Value: []byte(s.getResponseValue),
	}, nil
}

// Get fails if s.deleteFails is true, otherwise nothing happens.
func (s *testingStorage) Delete(_ context.Context, _ string) error {
	if s.deleteFails {
		return errors.New("failure")
	}

	return nil
}

// Put fails if s.putFails is true, otherwise nothing happens.
func (s *testingStorage) Put(_ context.Context, _ *logical.StorageEntry) error {
	if s.putFails {
		return errors.New("failure")
	}

	return nil
}
