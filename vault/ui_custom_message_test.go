// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIsTimeNowBetween verifies the proper functioning of the isTimeNowBetween
// function. The test calculates times that are either 1 or 2 hours away from
// time.Now, so there shouldn't be any timing concerns.
func TestIsTimeNowBetween(t *testing.T) {
	var (
		distantPast   time.Time = time.Now().Add(-2 * time.Hour)
		past          time.Time = time.Now().Add(-1 * time.Hour)
		future        time.Time = time.Now().Add(time.Hour)
		distantFuture time.Time = time.Now().Add(2 * time.Hour)
	)

	testcases := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		expectation bool
	}{
		{
			name:        "is between start and end times",
			startTime:   past,
			endTime:     future,
			expectation: true,
		},
		{
			name:        "both start and end times before",
			startTime:   distantPast,
			endTime:     past,
			expectation: false,
		},
		{
			name:        "both start and end times after",
			startTime:   future,
			endTime:     distantFuture,
			expectation: false,
		},
		// These test cases have the startTime occurring after the endTime
		// the algorithm will always return false in these cases.
		{
			name:        "is between start and end times, reversed",
			startTime:   future,
			endTime:     past,
			expectation: false,
		},
		{
			name:        "both start and end times before, reversed",
			startTime:   past,
			endTime:     distantPast,
			expectation: false,
		},
		{
			name:        "both start and end times after, reversed",
			startTime:   distantFuture,
			endTime:     future,
			expectation: false,
		},
	}

	for _, testcase := range testcases {
		result := isTimeNowBetween(testcase.startTime, testcase.endTime)
		assert.Equal(t, testcase.expectation, result, testcase.name)
	}
}

// TestCustomMessageBarrierView verifies that the
// (*UIConfig).customMessageBarrierView returns the correct logical.Storage
// based on the (*UIConfig).nsBarrierView field and whether a
// namespace.Namespace is exists in the provided context.Context.
func TestCustomMessageBarrierView(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	// Make sure that barrierStorage is returned regardless of Namespace in the
	// context when there's no nsBarrierView set.
	view := testUIConfig.customMessageBarrierView(context.Background())
	assert.Same(t, testUIConfig.barrierStorage, view)

	namespaceCtx := namespace.ContextWithNamespace(context.Background(), &namespace.Namespace{
		ID:   "abc12",
		Path: "dev/",
	})

	view = testUIConfig.customMessageBarrierView(namespaceCtx)
	assert.Same(t, testUIConfig.barrierStorage, view)

	// Setup the nsBarrierView to test with and without a namespace in the
	// context.
	testUIConfig.nsBarrierView = NewBarrierView(&logical.InmemStorage{}, "namespaces/")

	view = testUIConfig.customMessageBarrierView(namespaceCtx)
	assert.NotSame(t, testUIConfig.barrierStorage, view)

	view = testUIConfig.customMessageBarrierView(context.Background())
	assert.Same(t, testUIConfig.barrierStorage, view)
}

// TestListCustomMessages verifies the filtering logic in the
// (*UIConfig).ListCustomMessages method.
func TestListCustomMessages(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	// Check that with no filtering and no messages in the storage, no error is
	// returned and an empty slice of UICustomMessageEntry.
	results, err := testUIConfig.ListCustomMessages(context.Background(), ListUICustomMessagesFilters{})
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 0, len(results))

	// Load some custom messages in storage.
	commonMessage := "the message"
	commonTitle := "the title"
	distantPastTime := time.Now().Add(-24 * time.Hour)
	pastTime := time.Now().Add(-1 * time.Hour)
	futureTime := time.Now().Add(time.Hour)
	distantFutureTime := time.Now().Add(24 * time.Hour)

	customMessages := []UICustomMessageEntry{
		{
			Id:          mustGenerateUUID(t),
			Title:       commonTitle,
			Message:     commonMessage,
			MessageType: "banner",
			StartTime:   futureTime,
			EndTime:     distantFutureTime,
		},
		{
			Id:          mustGenerateUUID(t),
			Title:       commonTitle,
			Message:     commonMessage,
			MessageType: "banner",
			StartTime:   pastTime,
			EndTime:     futureTime,
		},
		{
			Id:          mustGenerateUUID(t),
			Title:       commonTitle,
			Message:     commonMessage,
			MessageType: "modal",
			StartTime:   distantPastTime,
			EndTime:     pastTime,
		},
		{
			Id:          mustGenerateUUID(t),
			Title:       commonTitle,
			Message:     commonMessage,
			MessageType: "modal",
			StartTime:   pastTime,
			EndTime:     futureTime,
		},
		{
			Id:            mustGenerateUUID(t),
			Title:         commonTitle,
			Message:       commonMessage,
			Authenticated: true,
			MessageType:   "banner",
			StartTime:     futureTime,
			EndTime:       distantFutureTime,
		},
		{
			Id:            mustGenerateUUID(t),
			Title:         commonTitle,
			Message:       commonMessage,
			Authenticated: true,
			MessageType:   "banner",
			StartTime:     pastTime,
			EndTime:       futureTime,
		},
		{
			Id:            mustGenerateUUID(t),
			Title:         commonTitle,
			Message:       commonMessage,
			Authenticated: true,
			MessageType:   "modal",
			StartTime:     distantPastTime,
			EndTime:       pastTime,
		},
		{
			Id:            mustGenerateUUID(t),
			Title:         commonTitle,
			Message:       commonMessage,
			Authenticated: true,
			MessageType:   "modal",
			StartTime:     pastTime,
			EndTime:       futureTime,
		},
	}

	for _, elem := range customMessages {
		storeTestUICustomMessage(t, testUIConfig.barrierStorage, elem)
	}

	// Test all of the different combinations of filters (including omissions)
	trueBool := true
	falseBool := false

	bannerString := "banner"
	modalString := "modal"

	for _, authenticated := range []*bool{nil, &falseBool, &trueBool} {
		for _, messageType := range []*string{nil, &bannerString, &modalString} {
			for _, active := range []*bool{nil, &falseBool, &trueBool} {
				filters := ListUICustomMessagesFilters{
					authenticated: authenticated,
					messageType:   messageType,
					active:        active,
				}

				results, err := testUIConfig.ListCustomMessages(context.Background(), filters)
				assert.NoError(t, err)
				assert.NotNil(t, results)
				// It's impossible to filter out every custom message because
				// each possible combination of all possible values for the 3
				// filtered properties exist in the sample set of messages.
				assert.NotEmpty(t, results)

				// This loop makes sure that only messages that satisfy the
				// filter criteria are returned. When the a filter criteria
				// is nil, then results are not filtered by that property, so
				// there's no need to check the results for that property.
				for _, result := range results {
					if authenticated != nil {
						assert.Equal(t, *authenticated, result.Authenticated)
					}

					if messageType != nil {
						assert.Equal(t, *messageType, result.MessageType)
					}

					if active != nil {
						assert.Equal(t, *active, result.active)
					}
				}
			}
		}
	}

	// Error testing
	//  There's 2 ways that the ListCustomMessages function can return an error:
	//  1. If the List function for the logical.Storage returns an error, or
	//  2. If the get function for the logical.Storage returns an error
	//
	// Both scenarios are simulated here.
	testUIConfig.barrierStorage = &testingStorage{
		listFails: true,
	}

	results, err = testUIConfig.ListCustomMessages(context.Background(), ListUICustomMessagesFilters{})
	assert.Error(t, err)
	assert.Nil(t, results)

	testUIConfig.barrierStorage = &testingStorage{
		getFails: true,
	}

	results, err = testUIConfig.ListCustomMessages(context.Background(), ListUICustomMessagesFilters{})
	assert.Error(t, err)
	assert.Nil(t, results)
}

// TestReadCustomMessage verifies that the (*UIConfig).ReadCustomMessage method
// behaves appropriatly based on different circumstances.
func TestReadCustomMessage(t *testing.T) {
	// Setup a UIConfig for testing with 2 sample custom messages stored in it.
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	entry := UICustomMessageEntry{
		Title:         "title",
		MessageType:   "banner",
		Authenticated: true,
		StartTime:     time.Now().Add(-1 * time.Hour),
		EndTime:       time.Now().Add(time.Hour),
	}

	for i := 0; i < 2; i++ {
		entry.Id = fmt.Sprintf("00%d", i)
		entry.Message = fmt.Sprintf("message %d", i)

		storeTestUICustomMessage(t, testUIConfig.barrierStorage, entry)
	}

	// Check that one of those messages can be read successfully
	result, err := testUIConfig.ReadCustomMessage(context.Background(), "001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "message 1", result.Message)
	assert.True(t, result.active)

	// Check that reading a non-existant message produces the correct errors.
	result, err = testUIConfig.ReadCustomMessage(context.Background(), "555")
	assert.NoError(t, err)
	assert.Nil(t, result)

	// Check that error due to JSON decode failure is returned
	require.NoError(t, testUIConfig.barrierStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   fmt.Sprintf("%s/bad", UICustomMessageKey),
		Value: []byte("_Not}JSO\"N"),
	}))

	result, err = testUIConfig.ReadCustomMessage(context.Background(), "bad")
	assert.Error(t, err)
	assert.Nil(t, result)

	// Check that when an error is return from the (logical.Storage).Get, it is
	// returned.
	testUIConfig.barrierStorage = &testingStorage{
		getFails: true,
	}

	result, err = testUIConfig.ReadCustomMessage(context.Background(), "001")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestCreateCustomMessage verifies that the (*UIConfig).CreateCustomMessage
// works as expected in all circumstances. This includes errors in the
// underlying logical.Storage, reaching the maximum count of custom messages.
func TestCreateCustomMessage(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	testUIConfig.nsBarrierView = NewBarrierView(testUIConfig.barrierStorage, "namespaces/")

	entry := UICustomMessageEntry{
		Title:         "title",
		Message:       "message",
		Authenticated: true,
		MessageType:   "modal",
		StartTime:     time.Now().Add(-1 * time.Hour),
		EndTime:       time.Now().Add(time.Hour),
	}

	// Create a custom message
	result, err := testUIConfig.CreateCustomMessage(context.Background(), entry)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.True(t, result.active)

	// Create more custom message to have the maximum number of custom messages.
	for i := 0; i < MaximumCustomMessageCountPerNamespace-1; i++ {
		result, err = testUIConfig.CreateCustomMessage(context.Background(), entry)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	}

	// We should now have the maximum number of custom messages in storage,
	// try to create one more to get an error.
	result, err = testUIConfig.CreateCustomMessage(context.Background(), entry)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "maximum")
	assert.Nil(t, result)

	// Make sure that the time validation check occurs before the maximum
	// custom messages check (by not only looking for an error, but checking
	// that it's the failure error for the validateStartAndEndTimes function).
	entryWithInvalidTimes := entry
	entryWithInvalidTimes.StartTime = entry.EndTime
	entryWithInvalidTimes.EndTime = entry.StartTime
	result, err = testUIConfig.CreateCustomMessage(context.Background(), entryWithInvalidTimes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "occur before")
	assert.Nil(t, result)

	// Make sure that a custom message can be still be created in a different
	// namespace.
	ns := &namespace.Namespace{
		ID:   "abc",
		Path: "def/",
	}
	result, err = testUIConfig.CreateCustomMessage(namespace.ContextWithNamespace(context.Background(), ns), entry)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)

	// Check that error is returned when an error occurs in the underlying
	// logical.Storage.
	testUIConfig.barrierStorage = &testingStorage{
		listFails: true,
	}

	result, err = testUIConfig.CreateCustomMessage(context.Background(), entry)
	assert.Error(t, err)
	assert.Nil(t, result)

	testUIConfig.barrierStorage = &testingStorage{
		putFails: true,
	}

	result, err = testUIConfig.CreateCustomMessage(context.Background(), entry)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestDeleteCustomMessage verifies that the (*UIConfig).DeleteCustomMessage
// method behaves correctly in all expected circumstances.
func TestDeleteCustomMessage(t *testing.T) {
	// Setup UIConfig with a custom message
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	entry := UICustomMessageEntry{
		Id:          "id",
		Title:       "title",
		Message:     "message",
		MessageType: "banner",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(time.Hour),
	}

	storeTestUICustomMessage(t, testUIConfig.barrierStorage, entry)

	// Check that the custom message can be deleted and no error is returned.
	// Then count the remaining custom messages.
	assert.NoError(t, testUIConfig.DeleteCustomMessage(context.Background(), "id"))

	count, err := testUIConfig.countCustomMessagesInternal(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// Check that deleting a non-existing custom message, doesn't return an
	// error.
	assert.NoError(t, testUIConfig.DeleteCustomMessage(context.Background(), "id"))

	// Check that an error is returned when there's an error in the underlying
	// logical.Storage.
	testUIConfig.barrierStorage = &testingStorage{
		deleteFails: true,
	}

	assert.Error(t, testUIConfig.DeleteCustomMessage(context.Background(), "id"))
}

// TestUpdateCustomMessage verifies that the (*UIConfig).UpdateCustomMessage
// method behaves correctly in all expected circumstances.
func TestUpdateCustomMessage(t *testing.T) {
	// Setup a UIConfig with a sample custom message
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	entry := UICustomMessageEntry{
		Id:          "id",
		Title:       "title",
		Message:     "message",
		MessageType: "banner",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(time.Hour),
	}

	storeTestUICustomMessage(t, testUIConfig.barrierStorage, entry)

	// Update the custom message entry and pass it to the UpdateCustomMessage
	// function.
	entry.Message = "Updated message"

	result, err := testUIConfig.UpdateCustomMessage(context.Background(), entry)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify that the change made it into the storage
	underlyingEntry, err := testUIConfig.retrieveCustomMessageInternal(context.Background(), "id")
	require.NoError(t, err)
	require.NotNil(t, underlyingEntry)
	assert.Equal(t, "Updated message", underlyingEntry.Message)

	// Make sure that modifying the returned *UICustomMessageEntry does nothing
	// to what's stored in the logical.Storage.
	result.MessageType = "modal"
	underlyingEntry, err = testUIConfig.retrieveCustomMessageInternal(context.Background(), "id")
	require.NoError(t, err)
	require.NotNil(t, underlyingEntry)
	assert.Equal(t, "banner", underlyingEntry.MessageType)

	// Check that updating an entry to a state that consists of invalid times
	// results in an error.
	entryWithInvalidTimes := entry
	entryWithInvalidTimes.StartTime = entry.EndTime
	entryWithInvalidTimes.EndTime = entry.StartTime
	result, err = testUIConfig.UpdateCustomMessage(context.Background(), entryWithInvalidTimes)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "occur before")
	assert.Nil(t, result)

	testUIConfig.barrierStorage = &testingStorage{
		putFails: true,
	}

	// Check that an error is returned if an error occurred in the underlying
	// logical.Storage.
	result, err = testUIConfig.UpdateCustomMessage(context.Background(), entry)
	assert.Error(t, err)
	assert.Nil(t, result)
}

// testingStorage is a struct with methods that satisfy the logical.Storage
// interface. It can be programmed to unconditionally return errors for any
// of its methods.
type testingStorage struct {
	listFails   bool
	getFails    bool
	deleteFails bool
	putFails    bool
}

// List returns a single key, dummy, unless the receiver has been configured to
// return an error, in which case it returns nil and a made up error.
func (s *testingStorage) List(_ context.Context, _ string) ([]string, error) {
	if s.listFails {
		return nil, errors.New("failure")
	}

	return []string{"dummy"}, nil
}

// Get returns a dummy logical.StorageEntry unless the receiver has been
// configured to return an error, in which case it returns nil and a made up
// error.
func (s *testingStorage) Get(_ context.Context, _ string) (*logical.StorageEntry, error) {
	if s.getFails {
		return nil, errors.New("failure")
	}

	return &logical.StorageEntry{
		Key:   "dummy",
		Value: []byte("{}"),
	}, nil
}

// Delete returns nothing, unless the receiver has been configured to return an
// error, in which case it returns a made up error.
func (s *testingStorage) Delete(_ context.Context, _ string) error {
	if s.deleteFails {
		return errors.New("failure")
	}

	return nil
}

// Put returns nothing, unless the receiver has been configured to return an
// error, in which case it returns a made up error.
func (s *testingStorage) Put(_ context.Context, _ *logical.StorageEntry) error {
	if s.putFails {
		return errors.New("failure")
	}

	return nil
}

// storeTestUICustomMessage takes care storing a UICustomMessageEntry into
// the provided logical.Storage without any errors.
func storeTestUICustomMessage(t *testing.T, storage logical.Storage, customMessage UICustomMessageEntry) {
	storageEntryValue, err := jsonutil.EncodeJSON(&customMessage)
	require.NoError(t, err)
	require.NotNil(t, storageEntryValue)
	require.NotEmpty(t, storageEntryValue)

	storageEntry := &logical.StorageEntry{
		Key:   fmt.Sprintf("%s/%s", UICustomMessageKey, customMessage.Id),
		Value: storageEntryValue,
	}

	err = storage.Put(context.Background(), storageEntry)
	require.NoError(t, err)
}

// mustGenerateUUID calls the uuid.GenerateUUID method and fails the current
// test if an error occurs.
func mustGenerateUUID(t *testing.T) string {
	result, err := uuid.GenerateUUID()
	require.NoError(t, err)

	return result
}

// TestValidateStartAndEndTimes verifies that the logic in the
// validateStartAndEndTimes function is correct.
func TestValidateStartAndEndTimes(t *testing.T) {
	var (
		timeNow        = time.Now()
		timeNowPlusOne = timeNow.Add(time.Second)
	)

	assert.NoError(t, validateStartAndEndTimes(timeNow, timeNowPlusOne))
	assert.Error(t, validateStartAndEndTimes(timeNow, timeNow))
	assert.Error(t, validateStartAndEndTimes(timeNowPlusOne, timeNow))
}
