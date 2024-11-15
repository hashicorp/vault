package jsonserialization

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	absser "github.com/microsoft/kiota-abstractions-go/serialization"
)

var buffPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// JsonSerializationWriter implements SerializationWriter for JSON.
type JsonSerializationWriter struct {
	writer                     *bytes.Buffer
	separatorIndices           []int
	onBeforeAssignFieldValues  absser.ParsableAction
	onAfterAssignFieldValues   absser.ParsableAction
	onStartObjectSerialization absser.ParsableWriter
}

// NewJsonSerializationWriter creates a new instance of the JsonSerializationWriter.
func NewJsonSerializationWriter() *JsonSerializationWriter {
	return &JsonSerializationWriter{
		writer:           buffPool.Get().(*bytes.Buffer),
		separatorIndices: make([]int, 0),
	}
}
func (w *JsonSerializationWriter) getWriter() *bytes.Buffer {
	if w.writer == nil {
		panic("The writer has already been closed. Call Reset instead of Close to reuse it or instantiate a new one.")
	}

	return w.writer
}
func (w *JsonSerializationWriter) writeRawValue(value ...string) {
	writer := w.getWriter()

	for _, v := range value {
		writer.WriteString(v)
	}
}
func (w *JsonSerializationWriter) writeStringValue(value string) {
	builder := &strings.Builder{}
	// Allocate at least enough space for the string and quotes. However, it's
	// possible that slightly overallocating may be a better strategy because then
	// it would at least be able to handle a few character escape sequences
	// without another allocation.
	builder.Grow(len(value) + 2)

	// Turning off HTML escaping may not be strictly necessary but it matches with
	// the current behavior. Testing with Exchange mail shows that it will
	// accept and properly interpret data sent with and without HTML escaping
	// enabled when creating emails with body content type HTML and HTML tags in
	// the body content.
	enc := json.NewEncoder(builder)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	enc.Encode(value)

	// Note that builder.String() returns a slice referencing the internal memory
	// of builder. This means it's unsafe to continue holding that reference once
	// this function exits (for example some conditions where a pool was used to
	// reduce strings.Builder allocations). We can use it here directly since
	// writeRawValue calls WriteString on a different buffer which should cause a
	// copy of the contents. If that's changed though this will need updated.
	s := builder.String()
	// Need to trim off the trailing newline the encoder adds.
	w.writeRawValue(s[:len(s)-1])
}
func (w *JsonSerializationWriter) writePropertyName(key string) {
	w.writeRawValue("\"", key, "\":")
}
func (w *JsonSerializationWriter) writePropertySeparator() {
	w.separatorIndices = append(w.separatorIndices, w.getWriter().Len())
	w.writeRawValue(",")
}
func (w *JsonSerializationWriter) writeArrayStart() {
	w.writeRawValue("[")
}
func (w *JsonSerializationWriter) writeArrayEnd() {
	w.writeRawValue("]")
}
func (w *JsonSerializationWriter) writeObjectStart() {
	w.writeRawValue("{")
}
func (w *JsonSerializationWriter) writeObjectEnd() {
	w.writeRawValue("}")
}

// WriteStringValue writes a String value to underlying the byte array.
func (w *JsonSerializationWriter) WriteStringValue(key string, value *string) error {
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
func (w *JsonSerializationWriter) WriteBoolValue(key string, value *bool) error {
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
func (w *JsonSerializationWriter) WriteByteValue(key string, value *byte) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt8Value writes a int8 value to underlying the byte array.
func (w *JsonSerializationWriter) WriteInt8Value(key string, value *int8) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt32Value writes a Int32 value to underlying the byte array.
func (w *JsonSerializationWriter) WriteInt32Value(key string, value *int32) error {
	if value != nil {
		cast := int64(*value)
		return w.WriteInt64Value(key, &cast)
	}
	return nil
}

// WriteInt64Value writes a Int64 value to underlying the byte array.
func (w *JsonSerializationWriter) WriteInt64Value(key string, value *int64) error {
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
func (w *JsonSerializationWriter) WriteFloat32Value(key string, value *float32) error {
	if value != nil {
		cast := float64(*value)
		return w.WriteFloat64Value(key, &cast)
	}
	return nil
}

// WriteFloat64Value writes a Float64 value to underlying the byte array.
func (w *JsonSerializationWriter) WriteFloat64Value(key string, value *float64) error {
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
func (w *JsonSerializationWriter) WriteTimeValue(key string, value *time.Time) error {
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
func (w *JsonSerializationWriter) WriteISODurationValue(key string, value *absser.ISODuration) error {
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
func (w *JsonSerializationWriter) WriteTimeOnlyValue(key string, value *absser.TimeOnly) error {
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
func (w *JsonSerializationWriter) WriteDateOnlyValue(key string, value *absser.DateOnly) error {
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
func (w *JsonSerializationWriter) WriteUUIDValue(key string, value *uuid.UUID) error {
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
func (w *JsonSerializationWriter) WriteByteArrayValue(key string, value []byte) error {
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
func (w *JsonSerializationWriter) WriteObjectValue(key string, item absser.Parsable, additionalValuesToMerge ...absser.Parsable) error {
	additionalValuesLen := len(additionalValuesToMerge)
	if item != nil || additionalValuesLen > 0 {
		untypedNode, isUntypedNode := item.(absser.UntypedNodeable)
		if isUntypedNode {
			switch value := untypedNode.(type) {
			case *absser.UntypedBoolean:
				w.WriteBoolValue(key, value.GetValue())
				return nil
			case *absser.UntypedFloat:
				w.WriteFloat32Value(key, value.GetValue())
				return nil
			case *absser.UntypedDouble:
				w.WriteFloat64Value(key, value.GetValue())
				return nil
			case *absser.UntypedInteger:
				w.WriteInt32Value(key, value.GetValue())
				return nil
			case *absser.UntypedLong:
				w.WriteInt64Value(key, value.GetValue())
				return nil
			case *absser.UntypedNull:
				w.WriteNullValue(key)
				return nil
			case *absser.UntypedString:
				w.WriteStringValue(key, value.GetValue())
				return nil
			case *absser.UntypedObject:
				if key != "" {
					w.writePropertyName(key)
				}
				properties := value.GetValue()
				if properties != nil {
					w.writeObjectStart()
					for objectKey, val := range properties {
						err := w.WriteObjectValue(objectKey, val)
						if err != nil {
							return err
						}
					}
					w.writeObjectEnd()
					if key != "" {
						w.writePropertySeparator()
					}
				}
				return nil
			case *absser.UntypedArray:
				if key != "" {
					w.writePropertyName(key)
				}
				values := value.GetValue()
				if values != nil {
					w.writeArrayStart()
					for _, val := range values {
						err := w.WriteObjectValue("", val)
						if err != nil {
							return err
						}
						w.writePropertySeparator()
					}
					w.writeArrayEnd()
				}
				if key != "" {
					w.writePropertySeparator()
				}
				return nil
			}
		}

		if key != "" {
			w.writePropertyName(key)
		}
		abstractions.InvokeParsableAction(w.GetOnBeforeSerialization(), item)
		_, isComposedTypeWrapper := item.(absser.ComposedTypeWrapper)
		if !isComposedTypeWrapper {
			w.writeObjectStart()
		}
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

		if !isComposedTypeWrapper {
			w.writeObjectEnd()
		}
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfObjectValues writes a collection of Parsable values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfObjectValues(key string, collection []absser.Parsable) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			if item != nil {
				err := w.WriteObjectValue("", item)
				if err != nil {
					return err
				}
			} else {
				err := w.WriteNullValue("")
				if err != nil {
					return err
				}
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfStringValues writes a collection of String values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfStringValues(key string, collection []string) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteStringValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfInt32Values writes a collection of Int32 values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfInt32Values(key string, collection []int32) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteInt32Value("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfInt64Values writes a collection of Int64 values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfInt64Values(key string, collection []int64) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteInt64Value("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfFloat32Values writes a collection of Float32 values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfFloat32Values(key string, collection []float32) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteFloat32Value("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfFloat64Values writes a collection of Float64 values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfFloat64Values(key string, collection []float64) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteFloat64Value("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfTimeValues writes a collection of Time values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfTimeValues(key string, collection []time.Time) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteTimeValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfISODurationValues writes a collection of ISODuration values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfISODurationValues(key string, collection []absser.ISODuration) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteISODurationValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfTimeOnlyValues writes a collection of TimeOnly values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfTimeOnlyValues(key string, collection []absser.TimeOnly) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteTimeOnlyValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfDateOnlyValues writes a collection of DateOnly values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfDateOnlyValues(key string, collection []absser.DateOnly) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteDateOnlyValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfUUIDValues writes a collection of UUID values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfUUIDValues(key string, collection []uuid.UUID) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteUUIDValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfBoolValues writes a collection of Bool values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfBoolValues(key string, collection []bool) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteBoolValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfByteValues writes a collection of Byte values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfByteValues(key string, collection []byte) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteByteValue("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// WriteCollectionOfInt8Values writes a collection of int8 values to underlying the byte array.
func (w *JsonSerializationWriter) WriteCollectionOfInt8Values(key string, collection []int8) error {
	if collection != nil { // empty collections are meaningful
		if key != "" {
			w.writePropertyName(key)
		}
		w.writeArrayStart()
		for _, item := range collection {
			err := w.WriteInt8Value("", &item)
			if err != nil {
				return err
			}
			w.writePropertySeparator()
		}

		w.writeArrayEnd()
		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

// GetSerializedContent returns the resulting byte array from the serialization writer.
func (w *JsonSerializationWriter) GetSerializedContent() ([]byte, error) {
	trimmed := w.getWriter().Bytes()
	buffLen := len(trimmed)

	for i := len(w.separatorIndices) - 1; i >= 0; i-- {
		idx := w.separatorIndices[i]

		if idx == buffLen-1 {
			trimmed = trimmed[:idx]
		} else if trimmed[idx+1] == byte(']') || trimmed[idx+1] == byte('}') {
			trimmed = append(trimmed[:idx], trimmed[idx+1:]...)
		}
	}

	trimmedCopy := make([]byte, len(trimmed))
	copy(trimmedCopy, trimmed)

	return trimmedCopy, nil
}

// WriteAnyValue an unknown value as a parameter.
func (w *JsonSerializationWriter) WriteAnyValue(key string, value interface{}) error {
	if value != nil {
		body, err := json.Marshal(value)
		if err != nil {
			return err
		}
		if key != "" {
			w.writePropertyName(key)
		}

		w.writeRawValue(string(body))

		if key != "" {
			w.writePropertySeparator()
		}
	}
	return nil
}

func (w *JsonSerializationWriter) WriteNullValue(key string) error {
	if key != "" {
		w.writePropertyName(key)
	}

	w.writeRawValue("null")

	if key != "" {
		w.writePropertySeparator()
	}

	return nil
}

func (w *JsonSerializationWriter) GetOnBeforeSerialization() absser.ParsableAction {
	return w.onBeforeAssignFieldValues
}

func (w *JsonSerializationWriter) SetOnBeforeSerialization(action absser.ParsableAction) error {
	w.onBeforeAssignFieldValues = action
	return nil
}

func (w *JsonSerializationWriter) GetOnAfterObjectSerialization() absser.ParsableAction {
	return w.onAfterAssignFieldValues
}

func (w *JsonSerializationWriter) SetOnAfterObjectSerialization(action absser.ParsableAction) error {
	w.onAfterAssignFieldValues = action
	return nil
}

func (w *JsonSerializationWriter) GetOnStartObjectSerialization() absser.ParsableWriter {
	return w.onStartObjectSerialization
}

func (w *JsonSerializationWriter) SetOnStartObjectSerialization(writer absser.ParsableWriter) error {
	w.onStartObjectSerialization = writer
	return nil
}

// WriteAdditionalData writes additional data to underlying the byte array.
func (w *JsonSerializationWriter) WriteAdditionalData(value map[string]interface{}) error {
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
			case absser.UntypedNodeable:
				err = w.WriteObjectValue(key, value)
			default:
				err = w.WriteAnyValue(key, &value)
			}
		}
	}
	return err
}

// Reset sets the internal buffer to empty, allowing the writer to be reused.
func (w *JsonSerializationWriter) Reset() error {
	w.getWriter().Reset()
	w.separatorIndices = w.separatorIndices[:0]
	return nil
}

// Close relases the internal buffer. Subsequent calls to the writer will panic.
func (w *JsonSerializationWriter) Close() error {
	if w.writer == nil {
		return nil
	}

	w.writer.Reset()
	buffPool.Put(w.writer)

	w.writer = nil
	w.separatorIndices = nil

	return nil
}
