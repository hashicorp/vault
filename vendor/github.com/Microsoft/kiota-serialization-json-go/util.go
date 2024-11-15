package jsonserialization

import (
	"fmt"
	"math"
	"reflect"
)

type numericRange struct {
	min          float64
	max          float64
	allowDecimal bool
}

var (
	numericTypeRanges = map[reflect.Kind]numericRange{
		reflect.Int8:    {math.MinInt8, math.MaxInt8, false},
		reflect.Uint8:   {0, math.MaxUint8, false},
		reflect.Int16:   {math.MinInt16, math.MaxInt16, false},
		reflect.Uint16:  {0, math.MaxUint16, false},
		reflect.Int32:   {math.MinInt32, math.MaxInt32, false},
		reflect.Uint32:  {0, math.MaxUint32, false},
		reflect.Int64:   {math.MinInt64, math.MaxInt64, false},
		reflect.Uint64:  {0, math.MaxUint64, false},
		reflect.Float32: {-math.MaxFloat32, math.MaxFloat32, true},
		reflect.Float64: {-math.MaxFloat64, math.MaxFloat64, true},
	}
)

type number interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

// isCompatible checks if the value is compatible with the type tp.
// It intentionally excludes checking if types are pointers to allow for possibility.
func isCompatible(value interface{}, tp reflect.Type) bool {
	// Can't join with lower, number types are always "convertible" just not losslessly.
	if isNumericType(value) && isNumericType(tp) {
		//NOTE: no need to check if number is compatible with another, always yes, just overflows
		//Check if number value is TRULY compatible
		return isCompatibleInt(value, tp)
	}

	return reflect.TypeOf(value).ConvertibleTo(tp)
}

// isNil checks if a value is nil or a nil interface, including nested pointers.
func isNil(a interface{}) bool {
	if a == nil {
		return true
	}
	val := reflect.ValueOf(a)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return val.IsNil()
	}
	return false
}

// as converts the value to the type T.
func as[T any](in interface{}, out T) error {
	// No point in trying anything if already nil
	if isNil(in) {
		return nil
	}

	// Make sure nothing is a pointer
	valValue := reflect.ValueOf(in)
	for valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
		in = valValue.Interface()
	}

	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Pointer || isNil(out) {
		return fmt.Errorf("out is not pointer or is nil")
	}

	nestedOutVal := outVal.Elem()
	// Handle the case where out is a pointer to an interface
	if nestedOutVal.Kind() == reflect.Interface && !nestedOutVal.IsNil() {
		nestedOutVal = nestedOutVal.Elem()
	}

	outType := nestedOutVal.Type()

	if !isCompatible(in, outType) {
		return fmt.Errorf("value '%v' is not compatible with type %T", in, nestedOutVal.Interface())
	}

	outVal.Elem().Set(valValue.Convert(outType))
	return nil
}

// isNumericType checks if the given type is a numeric type.
func isNumericType(in interface{}) bool {

	if in == nil {
		return false
	}

	tp, ok := in.(reflect.Type)
	if !ok {
		tp = reflect.TypeOf(in)
	}

	switch tp.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// isCompatibleInt checks if the given value is compatible with the specified integer type.
// It returns true if the value falls within the valid range for the type and has no decimal places.
// Otherwise, it returns false.
func isCompatibleInt(in interface{}, tp reflect.Type) bool {
	if !isNumericType(in) || !isNumericType(tp) {
		return false
	}

	inFloat := reflect.ValueOf(in).Convert(reflect.TypeOf(float64(0))).Float()
	hasDecimal := hasDecimalPlace(inFloat)

	if rangeInfo, ok := numericTypeRanges[tp.Kind()]; ok {
		if inFloat >= rangeInfo.min && inFloat <= rangeInfo.max {
			return rangeInfo.allowDecimal || !hasDecimal
		}
	}
	return false
}

// hasDecimalPlace checks if the given float64 value has a decimal place.
// It returns true if the fractional part of the value is greater than 0.0 (indicating a decimal).
// Otherwise, it returns false.
func hasDecimalPlace(value float64) bool {
	return value != float64(int64(value))
}
