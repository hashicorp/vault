package textserialization

import (
	testing "testing"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
	assert "github.com/stretchr/testify/assert"
)

func TestTextParseFactoryNodeHonoursInterface(t *testing.T) {
	instance := NewTextParseNodeFactory()
	assert.Implements(t, (*absser.ParseNodeFactory)(nil), instance)
}
