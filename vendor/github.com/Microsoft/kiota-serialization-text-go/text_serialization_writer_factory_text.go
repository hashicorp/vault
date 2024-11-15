package textserialization

import (
	testing "testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

func TestSerializationWriterFactoryHonoursInterface(t *testing.T) {
	instance := NewTextSerializationWriterFactory()
	assert.Implements(t, (*absser.SerializationWriterFactory)(nil), instance)
}
