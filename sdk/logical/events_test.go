package logical

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeSender struct {
	captured *EventData
}

func (f *fakeSender) SendEvent(ctx context.Context, eventType EventType, event *EventData) error {
	f.captured = event
	return nil
}

// TestSendEventWithOddParametersAddsExtraMetadata tests that an extra parameter is added to the metadata
// with a special key to note that it was extra.
func TestSendEventWithOddParametersAddsExtraMetadata(t *testing.T) {
	sender := &fakeSender{}
	// 0 or 2 arguments are okay
	err := SendEvent(context.Background(), sender, "foo")
	if err != nil {
		t.Fatal(err)
	}
	m := sender.captured.Metadata.AsMap()
	assert.NotContains(t, m, extraMetadataArgument)
	err = SendEvent(context.Background(), sender, "foo", "bar", "baz")
	if err != nil {
		t.Fatal(err)
	}
	m = sender.captured.Metadata.AsMap()
	assert.NotContains(t, m, extraMetadataArgument)

	// 1 or 3 arguments should give result in extraMetadataArgument in metadata
	err = SendEvent(context.Background(), sender, "foo", "extra")
	if err != nil {
		t.Fatal(err)
	}
	m = sender.captured.Metadata.AsMap()
	assert.Contains(t, m, extraMetadataArgument)
	assert.Equal(t, "extra", m[extraMetadataArgument])

	err = SendEvent(context.Background(), sender, "foo", "bar", "baz", "extra")
	if err != nil {
		t.Fatal(err)
	}
	m = sender.captured.Metadata.AsMap()
	assert.Contains(t, m, extraMetadataArgument)
	assert.Equal(t, "extra", m[extraMetadataArgument])
}
