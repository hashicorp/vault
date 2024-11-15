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

var NoStructuredDataError = errors.New("text does not support structured data")
var OnlyOneValue = errors.New("text serialization writer can only write one value")

// TextSerializationWriter implements SerializationWriter for JSON.
type TextSerializationWriter struct {
	writer                     []string
	onBeforeAssignFieldValues  absser.ParsableAction
	onAfterAssignFieldValues   absser.ParsableAction
	onStartObjectSerialization absser.ParsableWriter
}

// NewTextSerializationWriter creates a new instance of the TextSerializationWriter.
func NewTextSerializationWriter() *TextSerializationWriter {
	return &TextSerializationWriter{
		writer: make([]string, 0),
	}
}
func (w *TextSerializationWriter) writeStringValue(key string, value string) error {
	if key != "" {
		return NoStructuredDataError
	}
	if value != "" {
		if len(w.writer) > 0 {
			return OnlyOneValue
		}
		w.writer = append(w.writer, value)
	}
	return nil
}

// WriteStringValue writes a String value to underlying the byte array.
func (w *TextSerializationWriter) WriteStringValue(key string, value *string) error {
	if value != nil {
		return w.writeStringValue(key, *value)
	}
	return nil
}

// WriteBoolValue writes a Bool value to underlying the byte array.
func (w *TextSerializationWriter) WriteBoolValue(key string, value *bool) error {
	if value != nil {
		return w.writeStringValue(key, strconv.FormatBool(*value))
	}
	return nil
}

// WriteByteValue writes a Byte value to underlying the byte array.
func (w *TextSerializationWriter) WriteByteValue(key string, value *byte) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt8Value writes a int8 value to underlying the byte array.
func (w *TextSerializationWriter) WriteInt8Value(key string, value *int8) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt32Value writes a Int32 value to underlying the byte array.
func (w *TextSerializationWriter) WriteInt32Value(key string, value *int32) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt64Value writes a Int64 value to underlying the byte array.
func (w *TextSerializationWriter) WriteInt64Value(key string, value *int64) error {
	if value != nil {
		return w.writeStringValue(key, strconv.FormatInt(*value, 10))
	}
	return nil
}

// WriteFloat32Value writes a Float32 value to underlying the byte array.
func (w *TextSerializationWriter) WriteFloat32Value(key string, value *float32) error {
	if value != nil {
		cast := float64(*value)
		return w.WriteFloat64Value(key, &cast)
	}
	return nil
}

// WriteFloat64Value writes a Float64 value to underlying the byte array.
func (w *TextSerializationWriter) WriteFloat64Value(key string, value *float64) error {
	if value != nil {
		return w.writeStringValue(key, strconv.FormatFloat(*value, 'f', -1, 64))
	}
	return nil
}

// WriteTimeValue writes a Time value to underlying the byte array.
func (w *TextSerializationWriter) WriteTimeValue(key string, value *time.Time) error {
	if value != nil {
		return w.writeStringValue(key, (*value).String())
	}
	return nil
}

// WriteISODurationValue writes a ISODuration value to underlying the byte array.
func (w *TextSerializationWriter) WriteISODurationValue(key string, value *absser.ISODuration) error {
	if value != nil {
		return w.writeStringValue(key, (*value).String())
	}
	return nil
}

// WriteTimeOnlyValue writes a TimeOnly value to underlying the byte array.
func (w *TextSerializationWriter) WriteTimeOnlyValue(key string, value *absser.TimeOnly) error {
	if value != nil {
		return w.writeStringValue(key, (*value).String())
	}
	return nil
}

// WriteDateOnlyValue writes a DateOnly value to underlying the byte array.
func (w *TextSerializationWriter) WriteDateOnlyValue(key string, value *absser.DateOnly) error {
	if value != nil {
		return w.writeStringValue(key, (*value).String())
	}
	return nil
}

// WriteUUIDValue writes a UUID value to underlying the byte array.
func (w *TextSerializationWriter) WriteUUIDValue(key string, value *uuid.UUID) error {
	if value != nil {
		return w.writeStringValue(key, (*value).String())
	}
	return nil
}

// WriteByteArrayValue writes a ByteArray value to underlying the byte array.
func (w *TextSerializationWriter) WriteByteArrayValue(key string, value []byte) error {
	if value != nil {
		return w.writeStringValue(key, base64.StdEncoding.EncodeToString(value))
	}
	return nil
}

// WriteObjectValue writes a Parsable value to underlying the byte array.
func (w *TextSerializationWriter) WriteObjectValue(key string, item absser.Parsable, additionalValuesToMerge ...absser.Parsable) error {
	return NoStructuredDataError
}

// WriteCollectionOfObjectValues writes a collection of Parsable values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfObjectValues(key string, collection []absser.Parsable) error {
	return NoStructuredDataError
}

// WriteCollectionOfStringValues writes a collection of String values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfStringValues(key string, collection []string) error {
	return NoStructuredDataError
}

// WriteCollectionOfInt32Values writes a collection of Int32 values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfInt32Values(key string, collection []int32) error {
	return NoStructuredDataError
}

// WriteCollectionOfInt64Values writes a collection of Int64 values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfInt64Values(key string, collection []int64) error {
	return NoStructuredDataError
}

// WriteCollectionOfFloat32Values writes a collection of Float32 values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfFloat32Values(key string, collection []float32) error {
	return NoStructuredDataError
}

// WriteCollectionOfFloat64Values writes a collection of Float64 values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfFloat64Values(key string, collection []float64) error {
	return NoStructuredDataError
}

// WriteCollectionOfTimeValues writes a collection of Time values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfTimeValues(key string, collection []time.Time) error {
	return NoStructuredDataError
}

// WriteCollectionOfISODurationValues writes a collection of ISODuration values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfISODurationValues(key string, collection []absser.ISODuration) error {
	return NoStructuredDataError
}

// WriteCollectionOfTimeOnlyValues writes a collection of TimeOnly values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfTimeOnlyValues(key string, collection []absser.TimeOnly) error {
	return NoStructuredDataError
}

// WriteCollectionOfDateOnlyValues writes a collection of DateOnly values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfDateOnlyValues(key string, collection []absser.DateOnly) error {
	return NoStructuredDataError
}

// WriteCollectionOfUUIDValues writes a collection of UUID values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfUUIDValues(key string, collection []uuid.UUID) error {
	return NoStructuredDataError
}

// WriteCollectionOfBoolValues writes a collection of Bool values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfBoolValues(key string, collection []bool) error {
	return NoStructuredDataError
}

// WriteCollectionOfByteValues writes a collection of Byte values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfByteValues(key string, collection []byte) error {
	return NoStructuredDataError
}

// WriteCollectionOfInt8Values writes a collection of int8 values to underlying the byte array.
func (w *TextSerializationWriter) WriteCollectionOfInt8Values(key string, collection []int8) error {
	return NoStructuredDataError
}

// GetSerializedContent returns the resulting byte array from the serialization writer.
func (w *TextSerializationWriter) GetSerializedContent() ([]byte, error) {
	resultStr := strings.Join(w.writer, "")
	return []byte(resultStr), nil
}

// WriteAdditionalData writes additional data to underlying the byte array.
func (w *TextSerializationWriter) WriteAdditionalData(value map[string]interface{}) error {
	return NoStructuredDataError
}

// WriteAnyValue an unknown value as a parameter.
func (w *TextSerializationWriter) WriteAnyValue(key string, value interface{}) error {
	return NoStructuredDataError
}

// Close clears the internal buffer.
func (w *TextSerializationWriter) Close() error {
	return nil
}

func (w *TextSerializationWriter) WriteNullValue(key string) error {
	return NoStructuredDataError
}

func (w *TextSerializationWriter) GetOnBeforeSerialization() absser.ParsableAction {
	return w.onBeforeAssignFieldValues
}

func (w *TextSerializationWriter) SetOnBeforeSerialization(action absser.ParsableAction) error {
	w.onBeforeAssignFieldValues = action
	return nil
}

func (w *TextSerializationWriter) GetOnAfterObjectSerialization() absser.ParsableAction {
	return w.onAfterAssignFieldValues
}

func (w *TextSerializationWriter) SetOnAfterObjectSerialization(action absser.ParsableAction) error {
	w.onAfterAssignFieldValues = action
	return nil
}

func (w *TextSerializationWriter) GetOnStartObjectSerialization() absser.ParsableWriter {
	return w.onStartObjectSerialization
}

func (w *TextSerializationWriter) SetOnStartObjectSerialization(writer absser.ParsableWriter) error {
	w.onStartObjectSerialization = writer
	return nil
}
