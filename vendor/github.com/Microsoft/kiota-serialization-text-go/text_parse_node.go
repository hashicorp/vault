// Package textserialization is the default Kiota serialization implementation for text.
package textserialization

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// TextParseNode is a ParseNode implementation for JSON.
type TextParseNode struct {
	value                     string
	onBeforeAssignFieldValues absser.ParsableAction
	onAfterAssignFieldValues  absser.ParsableAction
}

// NewTextParseNode creates a new TextParseNode.
func NewTextParseNode(content []byte) (*TextParseNode, error) {
	if len(content) == 0 {
		return nil, errors.New("content is empty")
	}
	value, err := loadTextTree(content)
	return value, err
}
func loadTextTree(content []byte) (*TextParseNode, error) {
	return &TextParseNode{
		value: string(content),
	}, nil
}

// GetChildNode returns a new parse node for the given identifier.
func (n *TextParseNode) GetChildNode(index string) (absser.ParseNode, error) {
	return nil, NoStructuredDataError
}

// GetObjectValue returns the Parsable value from the node.
func (n *TextParseNode) GetObjectValue(ctor absser.ParsableFactory) (absser.Parsable, error) {
	return nil, NoStructuredDataError
}

// GetCollectionOfObjectValues returns the collection of Parsable values from the node.
func (n *TextParseNode) GetCollectionOfObjectValues(ctor absser.ParsableFactory) ([]absser.Parsable, error) {
	return nil, NoStructuredDataError
}

// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
func (n *TextParseNode) GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error) {
	return nil, NoStructuredDataError
}

// GetCollectionOfEnumValues returns the collection of Enum values from the node.
func (n *TextParseNode) GetCollectionOfEnumValues(parser absser.EnumFactory) ([]interface{}, error) {
	return nil, NoStructuredDataError
}

// GetStringValue returns a String value from the nodes.
func (n *TextParseNode) GetStringValue() (*string, error) {
	if n == nil {
		return nil, nil
	}
	val := strings.Trim(n.value, "\"")
	return &val, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *TextParseNode) GetBoolValue() (*bool, error) {
	if n == nil {
		return nil, nil
	}
	val, err := strconv.ParseBool(n.value)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

// GetInt8Value returns a int8 value from the nodes.
func (n *TextParseNode) GetInt8Value() (*int8, error) {
	if n == nil {
		return nil, nil
	}
	val, err := strconv.ParseInt(n.value, 0, 8)
	if err != nil {
		return nil, err
	}
	cast := int8(val)
	return &cast, nil
}

// GetBoolValue returns a Bool value from the nodes.
func (n *TextParseNode) GetByteValue() (*byte, error) {
	if n == nil {
		return nil, nil
	}
	val, err := strconv.ParseInt(n.value, 0, 8)
	if err != nil {
		return nil, err
	}
	cast := uint8(val)
	return &cast, nil
}

// GetFloat32Value returns a Float32 value from the nodes.
func (n *TextParseNode) GetFloat32Value() (*float32, error) {
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	cast := float32(*v)
	return &cast, nil
}

// GetFloat64Value returns a Float64 value from the nodes.
func (n *TextParseNode) GetFloat64Value() (*float64, error) {
	if n == nil {
		return nil, nil
	}
	val, err := strconv.ParseFloat(n.value, 0)
	if err != nil {
		return nil, err
	}
	cast := float64(val)
	return &cast, nil
}

// GetInt32Value returns a Int32 value from the nodes.
func (n *TextParseNode) GetInt32Value() (*int32, error) {
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	cast := int32(*v)
	return &cast, nil
}

// GetInt64Value returns a Int64 value from the nodes.
func (n *TextParseNode) GetInt64Value() (*int64, error) {
	v, err := n.GetFloat64Value()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	cast := int64(*v)
	return &cast, nil
}

// GetTimeValue returns a Time value from the nodes.
func (n *TextParseNode) GetTimeValue() (*time.Time, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *v)
	return &parsed, err
}

// GetISODurationValue returns a ISODuration value from the nodes.
func (n *TextParseNode) GetISODurationValue() (*absser.ISODuration, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseISODuration(*v)
}

// GetTimeOnlyValue returns a TimeOnly value from the nodes.
func (n *TextParseNode) GetTimeOnlyValue() (*absser.TimeOnly, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseTimeOnly(*v)
}

// GetDateOnlyValue returns a DateOnly value from the nodes.
func (n *TextParseNode) GetDateOnlyValue() (*absser.DateOnly, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return absser.ParseDateOnly(*v)
}

// GetUUIDValue returns a UUID value from the nodes.
func (n *TextParseNode) GetUUIDValue() (*uuid.UUID, error) {
	v, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	parsed, err := uuid.Parse(*v)
	return &parsed, err
}

// GetEnumValue returns a Enum value from the nodes.
func (n *TextParseNode) GetEnumValue(parser absser.EnumFactory) (interface{}, error) {
	if parser == nil {
		return nil, errors.New("parser is nil")
	}
	s, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return parser(*s)
}

// GetByteArrayValue returns a ByteArray value from the nodes.
func (n *TextParseNode) GetByteArrayValue() ([]byte, error) {
	s, err := n.GetStringValue()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(*s)
}

// GetRawValue returns a ByteArray value from the nodes.
func (n *TextParseNode) GetRawValue() (interface{}, error) {
	return n.value, nil
}

// GetOnBeforeAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) GetOnBeforeAssignFieldValues() absser.ParsableAction {
	return n.onBeforeAssignFieldValues
}

// SetOnBeforeAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) SetOnBeforeAssignFieldValues(action absser.ParsableAction) error {
	n.onBeforeAssignFieldValues = action
	return nil
}

// GetOnAfterAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) GetOnAfterAssignFieldValues() absser.ParsableAction {
	return n.onAfterAssignFieldValues
}

// SetOnAfterAssignFieldValues returns a ByteArray value from the nodes.
func (n *TextParseNode) SetOnAfterAssignFieldValues(action absser.ParsableAction) error {
	n.onAfterAssignFieldValues = action
	return nil
}
