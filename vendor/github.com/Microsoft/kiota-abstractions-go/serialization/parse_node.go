package serialization

import (
	"time"

	"github.com/google/uuid"
)

// ParseNode Interface for a deserialization node in a parse tree. This interface provides an abstraction layer over serialization formats, libraries and implementations.
type ParseNode interface {
	// GetChildNode returns a new parse node for the given identifier.
	GetChildNode(index string) (ParseNode, error)
	// GetCollectionOfObjectValues returns the collection of Parsable values from the node.
	GetCollectionOfObjectValues(ctor ParsableFactory) ([]Parsable, error)
	// GetCollectionOfPrimitiveValues returns the collection of primitive values from the node.
	GetCollectionOfPrimitiveValues(targetType string) ([]interface{}, error)
	// GetCollectionOfEnumValues returns the collection of Enum values from the node.
	GetCollectionOfEnumValues(parser EnumFactory) ([]interface{}, error)
	// GetObjectValue returns the Parsable value from the node.
	GetObjectValue(ctor ParsableFactory) (Parsable, error)
	// GetStringValue returns a String value from the nodes.
	GetStringValue() (*string, error)
	// GetBoolValue returns a Bool value from the nodes.
	GetBoolValue() (*bool, error)
	// GetInt8Value returns a int8 value from the nodes.
	GetInt8Value() (*int8, error)
	// GetByteValue returns a Byte value from the nodes.
	GetByteValue() (*byte, error)
	// GetFloat32Value returns a Float32 value from the nodes.
	GetFloat32Value() (*float32, error)
	// GetFloat64Value returns a Float64 value from the nodes.
	GetFloat64Value() (*float64, error)
	// GetInt32Value returns a Int32 value from the nodes.
	GetInt32Value() (*int32, error)
	// GetInt64Value returns a Int64 value from the nodes.
	GetInt64Value() (*int64, error)
	// GetTimeValue returns a Time value from the nodes.
	GetTimeValue() (*time.Time, error)
	// GetISODurationValue returns a ISODuration value from the nodes.
	GetISODurationValue() (*ISODuration, error)
	// GetTimeOnlyValue returns a TimeOnly value from the nodes.
	GetTimeOnlyValue() (*TimeOnly, error)
	// GetDateOnlyValue returns a DateOnly value from the nodes.
	GetDateOnlyValue() (*DateOnly, error)
	// GetUUIDValue returns a UUID value from the nodes.
	GetUUIDValue() (*uuid.UUID, error)
	// GetEnumValue returns a Enum value from the nodes.
	GetEnumValue(parser EnumFactory) (interface{}, error)
	// GetByteArrayValue returns a ByteArray value from the nodes.
	GetByteArrayValue() ([]byte, error)
	// GetRawValue returns the values of the node as an interface of any type.
	GetRawValue() (interface{}, error)
	// GetOnBeforeAssignFieldValues returns a callback invoked before the node is deserialized.
	GetOnBeforeAssignFieldValues() ParsableAction
	// SetOnBeforeAssignFieldValues sets a callback invoked before the node is deserialized.
	SetOnBeforeAssignFieldValues(ParsableAction) error
	// GetOnAfterAssignFieldValues returns a callback invoked after the node is deserialized.
	GetOnAfterAssignFieldValues() ParsableAction
	// SetOnAfterAssignFieldValues sets a callback invoked after the node is deserialized.
	SetOnAfterAssignFieldValues(ParsableAction) error
}
