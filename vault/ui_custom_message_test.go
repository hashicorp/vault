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

func TestIsTimeNowBetween(t *testing.T) {
	testcases := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		expectation bool
	}{
		{
			name:        "is between start and end times",
			startTime:   time.Now().Add(-1 * time.Hour),
			endTime:     time.Now().Add(time.Hour),
			expectation: true,
		},
		{
			name:        "both start and end times before",
			startTime:   time.Now().Add(-2 * time.Hour),
			endTime:     time.Now().Add(-1 * time.Hour),
			expectation: false,
		},
		{
			name:        "both start and end times after",
			startTime:   time.Now().Add(time.Hour),
			endTime:     time.Now().Add(2 * time.Hour),
			expectation: false,
		},
		{
			name:        "is between start and end times, reversed",
			startTime:   time.Now().Add(time.Hour),
			endTime:     time.Now().Add(-1 * time.Hour),
			expectation: false,
		},
		{
			name:        "both start and end times before, reversed",
			startTime:   time.Now().Add(-1 * time.Hour),
			endTime:     time.Now().Add(-2 * time.Hour),
			expectation: false,
		},
		{
			name:        "both start and end times after, reversed",
			startTime:   time.Now().Add(2 * time.Hour),
			endTime:     time.Now().Add(time.Hour),
			expectation: false,
		},
	}

	for _, testcase := range testcases {
		result := isTimeNowBetween(testcase.startTime, testcase.endTime)
		assert.Equal(t, testcase.expectation, result, testcase.name)
	}
}

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

	// Make sure that barrierStorage is returned when no namespace in context.
	testUIConfig.nsBarrierView = NewBarrierView(&logical.InmemStorage{}, "namespaces/")

	view = testUIConfig.customMessageBarrierView(context.Background())
	assert.Same(t, testUIConfig.barrierStorage, view)

	// Make sure that a sub view is returned when Namespace is in the context
	// and nsBarrierView is set.
	view = testUIConfig.customMessageBarrierView(namespaceCtx)
	assert.NotSame(t, testUIConfig.barrierStorage, view)
	assert.NotSame(t, testUIConfig.nsBarrierView, view)
}

func TestListCustomMessages(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	// Happy path testing first
	results, err := testUIConfig.ListCustomMessages(context.Background(), ListUICustomMessagesFilters{})
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 0, len(results))

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

	// Non-happy path testing
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

func TestReadCustomMessage(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	// Happy path testing
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

		require.NoError(t, testUIConfig.saveCustomMessage(context.Background(), entry))
	}

	result, err := testUIConfig.ReadCustomMessage(context.Background(), "001")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "message 1", result.Message)
	assert.True(t, result.active)

	// Check the error paths

	// Setup a JSON decode failure
	require.NoError(t, testUIConfig.barrierStorage.Put(context.Background(), &logical.StorageEntry{
		Key:   fmt.Sprintf("%s/bad", UICustomMessageKey),
		Value: []byte("_Not}JSO\"N"),
	}))

	result, err = testUIConfig.ReadCustomMessage(context.Background(), "bad")
	assert.Error(t, err)
	assert.Nil(t, result)

	// Setup a failure when Get is called on the logical.Storage.
	testUIConfig.barrierStorage = &testingStorage{
		getFails: true,
	}

	result, err = testUIConfig.ReadCustomMessage(context.Background(), "001")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCreateCustomMessage(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	entry := UICustomMessageEntry{
		Title:         "title",
		Message:       "message",
		Authenticated: true,
		MessageType:   "modal",
		StartTime:     time.Now().Add(-1 * time.Hour),
		EndTime:       time.Now().Add(time.Hour),
	}

	// Happy path first!
	result, err := testUIConfig.CreateCustomMessage(context.Background(), entry)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Id)
	assert.True(t, result.active)

	// Non-happy path test cases
	// Exceeding maximum custom message count
	for i := 0; i < MaximumCustomMessageCount-1; i++ {
		result, err = testUIConfig.CreateCustomMessage(context.Background(), entry)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	}

	// We should now have the maximum number of custom messages in storage,
	// try to create one more to get an error.
	result, err = testUIConfig.CreateCustomMessage(context.Background(), entry)
	assert.Error(t, err)
	assert.Nil(t, result)

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

func TestDeleteCustomMessage(t *testing.T) {
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

	require.NoError(t, testUIConfig.saveCustomMessage(context.Background(), entry))

	assert.NoError(t, testUIConfig.DeleteCustomMessage(context.Background(), "id"))

	count, err := testUIConfig.countCustomMessages(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	testUIConfig.barrierStorage = &testingStorage{
		deleteFails: true,
	}

	assert.Error(t, testUIConfig.DeleteCustomMessage(context.Background(), "id"))
}

func TestUpdateCustomMessage(t *testing.T) {
	testUIConfig := &UIConfig{
		barrierStorage: &logical.InmemStorage{},
	}

	// Store a custom message entry in storage
	entry := UICustomMessageEntry{
		Id:          "id",
		Title:       "title",
		Message:     "message",
		MessageType: "banner",
		StartTime:   time.Now().Add(-1 * time.Hour),
		EndTime:     time.Now().Add(time.Hour),
	}

	testUIConfig.saveCustomMessage(context.Background(), entry)

	// Update the custom message entry and pass it to the UpdateCustomMessage
	// function.
	entry.Message = "Updated message"

	result, err := testUIConfig.UpdateCustomMessage(context.Background(), entry)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Verify that the change made it into the storage
	underlyingEntry, err := testUIConfig.retrieveCustomMessage(context.Background(), "id")
	require.NoError(t, err)
	require.NotNil(t, underlyingEntry)
	assert.Equal(t, "Updated message", underlyingEntry.Message)

	// Make sure that modifying the returned *UICustomMessageEntry does nothing
	// to what's stored in the logical.Storage.
	result.MessageType = "modal"
	underlyingEntry, err = testUIConfig.retrieveCustomMessage(context.Background(), "id")
	require.NoError(t, err)
	require.NotNil(t, underlyingEntry)
	assert.Equal(t, "banner", underlyingEntry.MessageType)

	testUIConfig.barrierStorage = &testingStorage{
		putFails: true,
	}

	// Error path
	result, err = testUIConfig.UpdateCustomMessage(context.Background(), entry)
	assert.Error(t, err)
	assert.Nil(t, result)
}

type testingStorage struct {
	listFails   bool
	getFails    bool
	deleteFails bool
	putFails    bool
}

func (s *testingStorage) List(_ context.Context, _ string) ([]string, error) {
	if s.listFails {
		return nil, errors.New("failure")
	}

	return []string{"dummy"}, nil
}

func (s *testingStorage) Get(_ context.Context, _ string) (*logical.StorageEntry, error) {
	if s.getFails {
		return nil, errors.New("failure")
	}

	return &logical.StorageEntry{
		Key:   "dummy",
		Value: []byte("{}"),
	}, nil
}

func (s *testingStorage) Delete(_ context.Context, _ string) error {
	if s.deleteFails {
		return errors.New("failure")
	}

	return nil
}

func (s *testingStorage) Put(_ context.Context, _ *logical.StorageEntry) error {
	if s.putFails {
		return errors.New("failure")
	}

	return nil
}

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

func mustGenerateUUID(t *testing.T) string {
	result, err := uuid.GenerateUUID()
	require.NoError(t, err)

	return result
}
