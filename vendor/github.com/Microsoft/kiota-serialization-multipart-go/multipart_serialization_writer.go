package multipartserialization

import (
	"errors"
	"time"

	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/google/uuid"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// MultipartSerializationWriter implements SerializationWriter for URI Multipart encoded.
type MultipartSerializationWriter struct {
	writer                     []byte
	onBeforeAssignFieldValues  absser.ParsableAction
	onAfterAssignFieldValues   absser.ParsableAction
	onStartObjectSerialization absser.ParsableWriter
}

// NewMultipartSerializationWriter creates a new instance of the MultipartSerializationWriter.
func NewMultipartSerializationWriter() *MultipartSerializationWriter {
	return &MultipartSerializationWriter{
		writer: make([]byte, 0),
	}
}

// WriteStringValue writes a String value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteStringValue(key string, value *string) error {
	if key != "" {
		w.WriteByteArrayValue("", []byte(key))
	}
	if value != nil {
		if key != "" {
			w.WriteByteArrayValue("", []byte(": "))
		}
		w.WriteByteArrayValue("", []byte(*value))
	}
	w.WriteByteArrayValue("", []byte("\r\n"))
	return nil
}

// WriteBoolValue writes a Bool value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteBoolValue(key string, value *bool) error {
	return errors.New("serialization of bool value is not supported with MultipartSerializationWriter")
}

// WriteByteValue writes a Byte value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteByteValue(key string, value *byte) error {
	return errors.New("serialization of byte value is not supported with MultipartSerializationWriter")
}

// WriteInt8Value writes a int8 value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteInt8Value(key string, value *int8) error {
	return errors.New("serialization of int8 value is not supported with MultipartSerializationWriter")
}

// WriteInt32Value writes a Int32 value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteInt32Value(key string, value *int32) error {
	return errors.New("serialization of int32 value is not supported with MultipartSerializationWriter")
}

// WriteInt64Value writes a Int64 value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteInt64Value(key string, value *int64) error {
	return errors.New("serialization of int64 value is not supported with MultipartSerializationWriter")
}

// WriteFloat32Value writes a Float32 value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteFloat32Value(key string, value *float32) error {
	return errors.New("serialization of float32 value is not supported with MultipartSerializationWriter")
}

// WriteFloat64Value writes a Float64 value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteFloat64Value(key string, value *float64) error {
	return errors.New("serialization of float64 value is not supported with MultipartSerializationWriter")
}

// WriteTimeValue writes a Time value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteTimeValue(key string, value *time.Time) error {
	return errors.New("serialization of time value is not supported with MultipartSerializationWriter")
}

// WriteISODurationValue writes a ISODuration value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteISODurationValue(key string, value *absser.ISODuration) error {
	return errors.New("serialization of ISODuration value is not supported with MultipartSerializationWriter")
}

// WriteTimeOnlyValue writes a TimeOnly value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteTimeOnlyValue(key string, value *absser.TimeOnly) error {
	return errors.New("serialization of TimeOnly value is not supported with MultipartSerializationWriter")
}

// WriteDateOnlyValue writes a DateOnly value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteDateOnlyValue(key string, value *absser.DateOnly) error {
	return errors.New("serialization of DateOnly value is not supported with MultipartSerializationWriter")
}

// WriteUUIDValue writes a UUID value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteUUIDValue(key string, value *uuid.UUID) error {
	return errors.New("serialization of UUID value is not supported with MultipartSerializationWriter")
}

// WriteByteArrayValue writes a ByteArray value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteByteArrayValue(key string, value []byte) error {
	if value != nil {
		w.writer = append(w.writer, value...)
	}
	return nil
}

// WriteObjectValue writes a Parsable value to underlying the byte array.
func (w *MultipartSerializationWriter) WriteObjectValue(key string, item absser.Parsable, additionalValuesToMerge ...absser.Parsable) error {
	if item != nil {
		abstractions.InvokeParsableAction(w.GetOnBeforeSerialization(), item)
		err := abstractions.InvokeParsableWriter(w.GetOnStartObjectSerialization(), item, w)
		if err != nil {
			return err
		}
		if _, ok := item.(abstractions.MultipartBody); !ok {
			return errors.New("only the serialization of multipart bodies is supported with MultipartSerializationWriter")
		}
		err = item.Serialize(w)

		abstractions.InvokeParsableAction(w.GetOnAfterObjectSerialization(), item)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteCollectionOfObjectValues writes a collection of Parsable values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfObjectValues(key string, collection []absser.Parsable) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfStringValues writes a collection of String values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfStringValues(key string, collection []string) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfInt32Values writes a collection of Int32 values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfInt32Values(key string, collection []int32) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfInt64Values writes a collection of Int64 values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfInt64Values(key string, collection []int64) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfFloat32Values writes a collection of Float32 values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfFloat32Values(key string, collection []float32) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfFloat64Values writes a collection of Float64 values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfFloat64Values(key string, collection []float64) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfTimeValues writes a collection of Time values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfTimeValues(key string, collection []time.Time) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfISODurationValues writes a collection of ISODuration values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfISODurationValues(key string, collection []absser.ISODuration) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfTimeOnlyValues writes a collection of TimeOnly values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfTimeOnlyValues(key string, collection []absser.TimeOnly) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfDateOnlyValues writes a collection of DateOnly values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfDateOnlyValues(key string, collection []absser.DateOnly) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfUUIDValues writes a collection of UUID values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfUUIDValues(key string, collection []uuid.UUID) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfBoolValues writes a collection of Bool values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfBoolValues(key string, collection []bool) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfByteValues writes a collection of Byte values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfByteValues(key string, collection []byte) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// WriteCollectionOfInt8Values writes a collection of int8 values to underlying the byte array.
func (w *MultipartSerializationWriter) WriteCollectionOfInt8Values(key string, collection []int8) error {
	return errors.New("collections serialization is not supported with MultipartSerializationWriter")
}

// GetSerializedContent returns the resulting byte array from the serialization writer.
func (w *MultipartSerializationWriter) GetSerializedContent() ([]byte, error) {
	return w.writer, nil
}

// WriteAnyValue an unknown value as a parameter.
func (w *MultipartSerializationWriter) WriteAnyValue(key string, value interface{}) error {
	return errors.New("serialization of any value is not supported with MultipartSerializationWriter")
}

// WriteAdditionalData writes additional data to underlying the byte array.
func (w *MultipartSerializationWriter) WriteAdditionalData(value map[string]interface{}) error {
	return errors.New("serialization of additional data is not supported with MultipartSerializationWriter")
}

// Close clears the internal buffer.
func (w *MultipartSerializationWriter) Close() error {
	w.writer = make([]byte, 0)
	return nil
}

func (w *MultipartSerializationWriter) WriteNullValue(key string) error {
	return errors.New("serialization of null value is not supported with MultipartSerializationWriter")
}

func (w *MultipartSerializationWriter) GetOnBeforeSerialization() absser.ParsableAction {
	return w.onBeforeAssignFieldValues
}

func (w *MultipartSerializationWriter) SetOnBeforeSerialization(action absser.ParsableAction) error {
	w.onBeforeAssignFieldValues = action
	return nil
}

func (w *MultipartSerializationWriter) GetOnAfterObjectSerialization() absser.ParsableAction {
	return w.onAfterAssignFieldValues
}

func (w *MultipartSerializationWriter) SetOnAfterObjectSerialization(action absser.ParsableAction) error {
	w.onAfterAssignFieldValues = action
	return nil
}

func (w *MultipartSerializationWriter) GetOnStartObjectSerialization() absser.ParsableWriter {
	return w.onStartObjectSerialization
}

func (w *MultipartSerializationWriter) SetOnStartObjectSerialization(writer absser.ParsableWriter) error {
	w.onStartObjectSerialization = writer
	return nil
}
