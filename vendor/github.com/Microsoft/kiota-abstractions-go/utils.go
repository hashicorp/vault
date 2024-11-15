package abstractions

import (
	"github.com/google/uuid"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	"time"
)

// SetValue receives a source function and applies the results to the setter
//
// source is any function that produces (*T, error)
// setter recipient function of the result of the source if no error is produces
func SetValue[T interface{}](source func() (*T, error), setter func(t *T)) error {
	val, err := source()
	if err != nil {
		return err
	}
	if val != nil {
		setter(val)
	}
	return nil
}

// SetObjectValueFromSource takes a generic source with a discriminator receiver, reads value and writes it to a setter
//
// `source func() (*T, error)` is a generic getter with possible error response.
// `setter func(t *T)` generic function that can write a value from the source
func SetObjectValueFromSource[T interface{}](source func(ctor serialization.ParsableFactory) (serialization.Parsable, error), ctor serialization.ParsableFactory, setter func(t T)) error {
	val, err := source(ctor)
	if err != nil {
		return err
	}
	if val != nil {
		res := (val).(T)
		setter(res)
	}
	return nil
}

// SetCollectionValue is a utility function that receives a collection that can be cast to Parsable and a function that expects the results
//
// source is any function that receives a `ParsableFactory` and returns a slice of Parsable or error
// ctor is a ParsableFactory
// setter is a recipient of the function results
func SetCollectionValue[T interface{}](source func(ctor serialization.ParsableFactory) ([]serialization.Parsable, error), ctor serialization.ParsableFactory, setter func(t []T)) error {
	val, err := source(ctor)
	if err != nil {
		return err
	}
	if val != nil {
		res := CollectionCast[T](val)
		setter(res)
	}
	return nil
}

// CollectionApply applies an operation to every element of the slice and returns a result of the modified collection
//
//	is a slice of all the elementents to be mutated
//
// mutator applies an operation to the collection and returns a response of type `R`
func CollectionApply[T any, R interface{}](collection []T, mutator func(t T) R) []R {
	cast := make([]R, len(collection))
	for i, v := range collection {
		cast[i] = mutator(v)
	}
	return cast
}

// SetReferencedEnumValue is a utility function that receives an enum getter , EnumFactory and applies a de-referenced result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetReferencedEnumValue[T interface{}](source func(parser serialization.EnumFactory) (interface{}, error), parser serialization.EnumFactory, setter func(t *T)) error {
	val, err := source(parser)
	if err != nil {
		return err
	}
	if val != nil {
		res := (val).(*T)
		setter(res)
	}
	return nil
}

// SetCollectionOfReferencedEnumValue is a utility function that receives an enum collection source , EnumFactory and applies a de-referenced result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetCollectionOfReferencedEnumValue[T interface{}](source func(parser serialization.EnumFactory) ([]interface{}, error), parser serialization.EnumFactory, setter func(t []*T)) error {
	val, err := source(parser)
	if err != nil {
		return err
	}
	if val != nil {
		res := CollectionApply(val, func(v interface{}) *T { return (v).(*T) })
		setter(res)
	}
	return nil
}

// SetCollectionOfPrimitiveValue is a utility function that receives a collection of primitives , targetType and applies the result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// targetType is a string representing the type of result
// setter is a recipient of the function results
func SetCollectionOfPrimitiveValue[T interface{}](source func(targetType string) ([]interface{}, error), targetType string, setter func(t []T)) error {
	val, err := source(targetType)
	if err != nil {
		return err
	}
	if val != nil {
		res := CollectionCast[T](val)
		setter(res)
	}
	return nil
}

// SetCollectionOfReferencedPrimitiveValue is a utility function that receives a collection of primitives , targetType and applies the re-referenced result of the factory to a setter
//
// source is any function that receives a `EnumFactory` and returns an interface or error
// parser is an EnumFactory
// setter is a recipient of the function results
func SetCollectionOfReferencedPrimitiveValue[T interface{}](source func(targetType string) ([]interface{}, error), targetType string, setter func(t []T)) error {
	val, err := source(targetType)
	if err != nil {
		return err
	}
	if val != nil {
		res := CollectionValueCast[T](val)
		setter(res)
	}
	return nil
}

func p[T interface{}](t T) *T {
	return &t
}

// GetValueOrDefault Converts a Pointer to a value or returns a default value
func GetValueOrDefault[T interface{}](source func() *T, defaultValue T) T {
	result := source()
	if result != nil {
		return *result
	} else {
		return defaultValue
	}
}

// SetCollectionOfObjectValues returns an objects collection prototype for deserialization
func SetCollectionOfObjectValues[T interface{}](ctor serialization.ParsableFactory, setter func(t []T)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfObjectValues(ctor)
		if err != nil {
			return err
		}
		if val != nil {
			res := CollectionCast[T](val)
			setter(res)
		}
		return nil
	}
}

// SetCollectionOfPrimitiveValues returns a primitive's collection prototype for deserialization
func SetCollectionOfPrimitiveValues[T interface{}](targetType string, setter func(t []T)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfPrimitiveValues(targetType)
		if err != nil {
			return err
		}
		if val != nil {
			res := CollectionValueCast[T](val)
			setter(res)
		}
		return nil
	}
}

// SetCollectionOfEnumValues returns an enum prototype for deserialization
func SetCollectionOfEnumValues[T interface{}](parser serialization.EnumFactory, setter func(t []T)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		val, err := n.GetCollectionOfEnumValues(parser)
		if err != nil {
			return err
		}
		if val != nil {
			res := CollectionValueCast[T](val)
			setter(res)
		}
		return nil
	}
}

// SetObjectValue returns an object prototype for deserialization
func SetObjectValue[T interface{}](ctor serialization.ParsableFactory, setter func(t T)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetObjectValueFromSource(n.GetObjectValue, ctor, setter)
	}
}

// SetStringValue returns a string prototype for deserialization
func SetStringValue(setter func(t *string)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetStringValue, setter)
	}
}

// SetBoolValue returns a boolean prototype for deserialization
func SetBoolValue(setter func(t *bool)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetBoolValue, setter)
	}
}

// SetInt8Value Returns an int8 prototype for deserialization
func SetInt8Value(setter func(t *int8)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetInt8Value, setter)
	}
}

// SetByteValue returns a byte prototype for deserialization
func SetByteValue(setter func(t *byte)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetByteValue, setter)
	}
}

// SetFloat32Value returns a float32 prototype for deserialization
func SetFloat32Value(setter func(t *float32)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetFloat32Value, setter)
	}
}

// SetFloat64Value returns a float64 prototype for deserialization
func SetFloat64Value(setter func(t *float64)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetFloat64Value, setter)
	}
}

// SetInt32Value returns a int32 prototype for deserialization
func SetInt32Value(setter func(t *int32)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetInt32Value, setter)
	}
}

// SetInt64Value returns a int64 prototype for deserialization
func SetInt64Value(setter func(t *int64)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetInt64Value, setter)
	}
}

// SetTimeValue returns a time.Time prototype for deserialization
func SetTimeValue(setter func(t *time.Time)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetTimeValue, setter)
	}
}

// SetISODurationValue returns a ISODuration prototype for deserialization
func SetISODurationValue(setter func(t *serialization.ISODuration)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetISODurationValue, setter)
	}
}

// SetTimeOnlyValue returns a TimeOnly prototype for deserialization
func SetTimeOnlyValue(setter func(t *serialization.TimeOnly)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetTimeOnlyValue, setter)
	}
}

// SetDateOnlyValue returns a DateOnly prototype for deserialization
func SetDateOnlyValue(setter func(t *serialization.DateOnly)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetDateOnlyValue, setter)
	}
}

// SetUUIDValue returns a DateOnly prototype for deserialization
func SetUUIDValue(setter func(t *uuid.UUID)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetValue(n.GetUUIDValue, setter)
	}
}

// SetEnumValue returns a Enum prototype for deserialization
func SetEnumValue[T interface{}](parser serialization.EnumFactory, setter func(t *T)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		return SetReferencedEnumValue(n.GetEnumValue, parser, setter)
	}
}

// SetByteArrayValue returns a []byte prototype for deserialization
func SetByteArrayValue(setter func(t []byte)) serialization.NodeParser {
	return func(n serialization.ParseNode) error {
		val, err := n.GetByteArrayValue()
		if err != nil {
			return err
		}
		if val != nil {
			setter(val)
		}
		return nil
	}
}

// CollectionCast casts a collection of values from any type T to given type R
func CollectionCast[R interface{}, T any](items []T) []R {
	cast := CollectionApply(items, func(v T) R { return any(v).(R) })
	return cast
}

// CollectionValueCast casts a collection of values from any type T to given type R
//
// Value cast can be used to cast memory addresses to the value of the pointer
func CollectionValueCast[R interface{}, T any](items []T) []R {
	cast := CollectionApply(items, func(v T) R { return *(any(v).(*R)) })
	return cast
}

// CollectionStructCast casts a collection of values from any type T to given type R
//
// Value cast can be used to cast memory addresses to the value of the pointer
func CollectionStructCast[R interface{}, T any](items []T) []R {
	cast := CollectionApply(items, func(v T) R {
		temp := v
		return any(&temp).(R)
	})
	return cast
}

// InvokeParsableAction nil safe execution of ParsableAction
func InvokeParsableAction(action serialization.ParsableAction, parsable serialization.Parsable) {
	if action != nil {
		action(parsable)
	}
}

// InvokeParsableWriter executes the ParsableAction in a nil safe way
func InvokeParsableWriter(writer serialization.ParsableWriter, parsable serialization.Parsable, serializer serialization.SerializationWriter) error {
	if writer != nil {
		return writer(parsable, serializer)
	}
	return nil
}

// CopyMap returns a copy of map[string]string
func CopyMap[T comparable, R interface{}](items map[T]R) map[T]R {
	result := make(map[T]R)
	for idx, item := range items {
		result[idx] = item
	}
	return result
}

// CopyStringMap returns a copy of map[string]string
func CopyStringMap(items map[string]string) map[string]string {
	result := make(map[string]string)
	for idx, item := range items {
		result[idx] = item
	}
	return result
}
