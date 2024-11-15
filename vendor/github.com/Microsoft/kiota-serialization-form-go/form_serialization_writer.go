package formserialization

import (
	"encoding/base64"
	"errors"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

// FormSerializationWriter implements SerializationWriter for URI form encoded.
type FormSerializationWriter struct {
	writer                     []string
	depth                      int
	onBeforeAssignFieldValues  absser.ParsableAction
	onAfterAssignFieldValues   absser.ParsableAction
	onStartObjectSerialization absser.ParsableWriter
}

// NewFormSerializationWriter creates a new instance of the FormSerializationWriter.
func NewFormSerializationWriter() *FormSerializationWriter {
	return &FormSerializationWriter{
		writer: make([]string, 0),
	}
}
func (w *FormSerializationWriter) writeRawValue(value string) {
	w.writer = append(w.writer, url.QueryEscape(value))
}
func (w *FormSerializationWriter) writeStringValue(value string) {
	w.writeRawValue(value)
}
func (w *FormSerializationWriter) writePropertyName(key string) {
	w.writer = append(w.writer, url.QueryEscape(key)+"=")
}
func (w *FormSerializationWriter) writePropertySeparator() {
	w.writer = append(w.writer, "&")
}
func (w *FormSerializationWriter) trimLastPropertySeparator() {
	lastIdx := len(w.writer) - 1
	if lastIdx > -1 && w.writer[lastIdx] == "&" {
		w.writer[lastIdx] = ""
	}
}

// WriteStringValue writes a String value to underlying the byte array.
func (w *FormSerializationWriter) WriteStringValue(key string, value *string) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue(*value)
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteBoolValue writes a Bool value to underlying the byte array.
func (w *FormSerializationWriter) WriteBoolValue(key string, value *bool) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeRawValue(strconv.FormatBool(*value))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteByteValue writes a Byte value to underlying the byte array.
func (w *FormSerializationWriter) WriteByteValue(key string, value *byte) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt8Value writes a int8 value to underlying the byte array.
func (w *FormSerializationWriter) WriteInt8Value(key string, value *int8) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt32Value writes a Int32 value to underlying the byte array.
func (w *FormSerializationWriter) WriteInt32Value(key string, value *int32) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt64Value writes a Int64 value to underlying the byte array.
func (w *FormSerializationWriter) WriteInt64Value(key string, value *int64) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeRawValue(strconv.FormatInt(*value, 10))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteFloat32Value writes a Float32 value to underlying the byte array.
func (w *FormSerializationWriter) WriteFloat32Value(key string, value *float32) error {
	if value != nil {
		cast := float64(*value)
		return w.WriteFloat64Value(key, &cast)
	}
	return nil
}

// WriteFloat64Value writes a Float64 value to underlying the byte array.
func (w *FormSerializationWriter) WriteFloat64Value(key string, value *float64) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeRawValue(strconv.FormatFloat(*value, 'f', -1, 64))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteTimeValue writes a Time value to underlying the byte array.
func (w *FormSerializationWriter) WriteTimeValue(key string, value *time.Time) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).Format(time.RFC3339))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteISODurationValue writes a ISODuration value to underlying the byte array.
func (w *FormSerializationWriter) WriteISODurationValue(key string, value *absser.ISODuration) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteTimeOnlyValue writes a TimeOnly value to underlying the byte array.
func (w *FormSerializationWriter) WriteTimeOnlyValue(key string, value *absser.TimeOnly) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteDateOnlyValue writes a DateOnly value to underlying the byte array.
func (w *FormSerializationWriter) WriteDateOnlyValue(key string, value *absser.DateOnly) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteUUIDValue writes a UUID value to underlying the byte array.
func (w *FormSerializationWriter) WriteUUIDValue(key string, value *uuid.UUID) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue((*value).String())
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteByteArrayValue writes a ByteArray value to underlying the byte array.
func (w *FormSerializationWriter) WriteByteArrayValue(key string, value []byte) error {
	if key != "" && value != nil {
		w.writePropertyName(key)
	}
	if value != nil {
		w.writeStringValue(base64.StdEncoding.EncodeToString(value))
	}
	if key != "" && value != nil {
		w.writePropertySeparator()
	}
	return nil
}

// WriteObjectValue writes a Parsable value to underlying the byte array.
func (w *FormSerializationWriter) WriteObjectValue(key string, item absser.Parsable, additionalValuesToMerge ...absser.Parsable) error {
	if w.depth > 0 {
		return errors.New("nested objects serialization is not supported with FormSerializationWriter")
	}
	w.depth++
	additionalValuesLen := len(additionalValuesToMerge)
	if item != nil || additionalValuesLen > 0 {
		if key != "" {
			w.writePropertyName(key)
		}
		abstractions.InvokeParsableAction(w.GetOnBeforeSerialization(), item)
		if item != nil {
			err := abstractions.InvokeParsableWriter(w.GetOnStartObjectSerialization(), item, w)
			if err != nil {
				return err
			}
			err = item.Serialize(w)

			abstractions.InvokeParsableAction(w.GetOnAfterObjectSerialization(), item)
			if err != nil {
				return err
			}
		}

		for _, additionalValue := range additionalValuesToMerge {
			abstractions.InvokeParsableAction(w.GetOnBeforeSerialization(), additionalValue)
			err := abstractions.InvokeParsableWriter(w.GetOnStartObjectSerialization(), additionalValue, w)
			if err != nil {
				return err
			}
			err = additionalValue.Serialize(w)
			if err != nil {
				return err
			}
			abstractions.InvokeParsableAction(w.GetOnAfterObjectSerialization(), additionalValue)
		}

		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfObjectValues writes a collection of Parsable values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfObjectValues(key string, collection []absser.Parsable) error {
	return errors.New("collections serialization is not supported with FormSerializationWriter")
}

func writeCollectionOfPrimitiveValues[T interface{}](key string, writer func(string, *T) error, collection []T) error {
	if collection != nil && len(collection) > 0 {
		for _, item := range collection {
			err := writer(key, &item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// WriteCollectionOfStringValues writes a collection of String values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfStringValues(key string, collection []string) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteStringValue, collection)
}

// WriteCollectionOfInt32Values writes a collection of Int32 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfInt32Values(key string, collection []int32) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteInt32Value, collection)
}

// WriteCollectionOfInt64Values writes a collection of Int64 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfInt64Values(key string, collection []int64) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteInt64Value, collection)
}

// WriteCollectionOfFloat32Values writes a collection of Float32 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfFloat32Values(key string, collection []float32) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteFloat32Value, collection)
}

// WriteCollectionOfFloat64Values writes a collection of Float64 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfFloat64Values(key string, collection []float64) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteFloat64Value, collection)
}

// WriteCollectionOfTimeValues writes a collection of Time values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfTimeValues(key string, collection []time.Time) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteTimeValue, collection)
}

// WriteCollectionOfISODurationValues writes a collection of ISODuration values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfISODurationValues(key string, collection []absser.ISODuration) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteISODurationValue, collection)
}

// WriteCollectionOfTimeOnlyValues writes a collection of TimeOnly values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfTimeOnlyValues(key string, collection []absser.TimeOnly) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteTimeOnlyValue, collection)
}

// WriteCollectionOfDateOnlyValues writes a collection of DateOnly values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfDateOnlyValues(key string, collection []absser.DateOnly) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteDateOnlyValue, collection)
}

// WriteCollectionOfUUIDValues writes a collection of UUID values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfUUIDValues(key string, collection []uuid.UUID) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteUUIDValue, collection)
}

// WriteCollectionOfBoolValues writes a collection of Bool values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfBoolValues(key string, collection []bool) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteBoolValue, collection)
}

// WriteCollectionOfByteValues writes a collection of Byte values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfByteValues(key string, collection []byte) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteByteValue, collection)
}

// WriteCollectionOfInt8Values writes a collection of int8 values to underlying the byte array.
func (w *FormSerializationWriter) WriteCollectionOfInt8Values(key string, collection []int8) error {
	return writeCollectionOfPrimitiveValues(key, w.WriteInt8Value, collection)
}

// GetSerializedContent returns the resulting byte array from the serialization writer.
func (w *FormSerializationWriter) GetSerializedContent() ([]byte, error) {
	w.trimLastPropertySeparator()
	resultStr := strings.Join(w.writer, "")
	return []byte(resultStr), nil
}

// WriteAnyValue an unknown value as a parameter.
func (w *FormSerializationWriter) WriteAnyValue(key string, value interface{}) error {
	return errors.New("serialization of any value is not supported with FormSerializationWriter")
}

// WriteAdditionalData writes additional data to underlying the byte array.
func (w *FormSerializationWriter) WriteAdditionalData(value map[string]interface{}) error {
	var err error
	if len(value) != 0 {
		for key, input := range value {
			switch value := input.(type) {
			case absser.Parsable:
				err = w.WriteObjectValue(key, value)
			case []absser.Parsable:
				err = w.WriteCollectionOfObjectValues(key, value)
			case []string:
				err = w.WriteCollectionOfStringValues(key, value)
			case []bool:
				err = w.WriteCollectionOfBoolValues(key, value)
			case []byte:
				err = w.WriteCollectionOfByteValues(key, value)
			case []int8:
				err = w.WriteCollectionOfInt8Values(key, value)
			case []int32:
				err = w.WriteCollectionOfInt32Values(key, value)
			case []int64:
				err = w.WriteCollectionOfInt64Values(key, value)
			case []float32:
				err = w.WriteCollectionOfFloat32Values(key, value)
			case []float64:
				err = w.WriteCollectionOfFloat64Values(key, value)
			case []uuid.UUID:
				err = w.WriteCollectionOfUUIDValues(key, value)
			case []time.Time:
				err = w.WriteCollectionOfTimeValues(key, value)
			case []absser.ISODuration:
				err = w.WriteCollectionOfISODurationValues(key, value)
			case []absser.TimeOnly:
				err = w.WriteCollectionOfTimeOnlyValues(key, value)
			case []absser.DateOnly:
				err = w.WriteCollectionOfDateOnlyValues(key, value)
			case *string:
				err = w.WriteStringValue(key, value)
			case string:
				err = w.WriteStringValue(key, &value)
			case *bool:
				err = w.WriteBoolValue(key, value)
			case bool:
				err = w.WriteBoolValue(key, &value)
			case *byte:
				err = w.WriteByteValue(key, value)
			case byte:
				err = w.WriteByteValue(key, &value)
			case *int8:
				err = w.WriteInt8Value(key, value)
			case int8:
				err = w.WriteInt8Value(key, &value)
			case *int32:
				err = w.WriteInt32Value(key, value)
			case int32:
				err = w.WriteInt32Value(key, &value)
			case *int64:
				err = w.WriteInt64Value(key, value)
			case int64:
				err = w.WriteInt64Value(key, &value)
			case *float32:
				err = w.WriteFloat32Value(key, value)
			case float32:
				err = w.WriteFloat32Value(key, &value)
			case *float64:
				err = w.WriteFloat64Value(key, value)
			case float64:
				err = w.WriteFloat64Value(key, &value)
			case *uuid.UUID:
				err = w.WriteUUIDValue(key, value)
			case uuid.UUID:
				err = w.WriteUUIDValue(key, &value)
			case *time.Time:
				err = w.WriteTimeValue(key, value)
			case time.Time:
				err = w.WriteTimeValue(key, &value)
			case *absser.ISODuration:
				err = w.WriteISODurationValue(key, value)
			case absser.ISODuration:
				err = w.WriteISODurationValue(key, &value)
			case *absser.TimeOnly:
				err = w.WriteTimeOnlyValue(key, value)
			case absser.TimeOnly:
				err = w.WriteTimeOnlyValue(key, &value)
			case *absser.DateOnly:
				err = w.WriteDateOnlyValue(key, value)
			case absser.DateOnly:
				err = w.WriteDateOnlyValue(key, &value)
			}
		}
	}
	return err
}

// Close clears the internal buffer.
func (w *FormSerializationWriter) Close() error {
	w.writer = make([]string, 0)
	return nil
}

func (w *FormSerializationWriter) WriteNullValue(key string) error {
	if key != "" {
		w.writePropertyName(key)
	}

	w.writeRawValue("null")

	if key != "" {
		w.writePropertySeparator()
	}

	return nil
}

func (w *FormSerializationWriter) GetOnBeforeSerialization() absser.ParsableAction {
	return w.onBeforeAssignFieldValues
}

func (w *FormSerializationWriter) SetOnBeforeSerialization(action absser.ParsableAction) error {
	w.onBeforeAssignFieldValues = action
	return nil
}

func (w *FormSerializationWriter) GetOnAfterObjectSerialization() absser.ParsableAction {
	return w.onAfterAssignFieldValues
}

func (w *FormSerializationWriter) SetOnAfterObjectSerialization(action absser.ParsableAction) error {
	w.onAfterAssignFieldValues = action
	return nil
}

func (w *FormSerializationWriter) GetOnStartObjectSerialization() absser.ParsableWriter {
	return w.onStartObjectSerialization
}

func (w *FormSerializationWriter) SetOnStartObjectSerialization(writer absser.ParsableWriter) error {
	w.onStartObjectSerialization = writer
	return nil
}
