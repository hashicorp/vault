// Copyright (c) 2017-2021 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"context"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/decimal128"
)

const format = "2006-01-02 15:04:05.999999999"

// goTypeToSnowflake translates Go data type to Snowflake data type.
func goTypeToSnowflake(v driver.Value, tsmode snowflakeType) snowflakeType {
	switch t := v.(type) {
	case int64:
		return fixedType
	case float64:
		return realType
	case bool:
		return booleanType
	case string:
		return textType
	case []byte:
		if tsmode == binaryType {
			return binaryType // may be redundant but ensures BINARY type
		}
		if t == nil {
			return nullType // invalid byte array. won't take as BINARY
		}
		if len(t) != 1 {
			return unSupportedType
		}
		if _, err := dataTypeMode(t); err != nil {
			return unSupportedType
		}
		return changeType
	case time.Time:
		return tsmode
	}
	if supportedArrayBind(&driver.NamedValue{Value: v}) {
		return sliceType
	}
	return unSupportedType
}

// snowflakeTypeToGo translates Snowflake data type to Go data type.
func snowflakeTypeToGo(dbtype snowflakeType, scale int64) reflect.Type {
	switch dbtype {
	case fixedType:
		if scale == 0 {
			return reflect.TypeOf(int64(0))
		}
		return reflect.TypeOf(float64(0))
	case realType:
		return reflect.TypeOf(float64(0))
	case textType, variantType, objectType, arrayType:
		return reflect.TypeOf("")
	case dateType, timeType, timestampLtzType, timestampNtzType, timestampTzType:
		return reflect.TypeOf(time.Now())
	case binaryType:
		return reflect.TypeOf([]byte{})
	case booleanType:
		return reflect.TypeOf(true)
	}
	logger.Errorf("unsupported dbtype is specified. %v", dbtype)
	return reflect.TypeOf("")
}

// valueToString converts arbitrary golang type to a string. This is mainly used in binding data with placeholders
// in queries.
func valueToString(v driver.Value, tsmode snowflakeType) (*string, error) {
	logger.Debugf("TYPE: %v, %v", reflect.TypeOf(v), reflect.ValueOf(v))
	if v == nil {
		return nil, nil
	}
	v1 := reflect.ValueOf(v)
	switch v1.Kind() {
	case reflect.Bool:
		s := strconv.FormatBool(v1.Bool())
		return &s, nil
	case reflect.Int64:
		s := strconv.FormatInt(v1.Int(), 10)
		return &s, nil
	case reflect.Float64:
		s := strconv.FormatFloat(v1.Float(), 'g', -1, 32)
		return &s, nil
	case reflect.String:
		s := v1.String()
		return &s, nil
	case reflect.Slice, reflect.Map:
		if v1.IsNil() {
			return nil, nil
		}
		if bd, ok := v.([]byte); ok {
			if tsmode == binaryType {
				s := hex.EncodeToString(bd)
				return &s, nil
			}
		}
		// TODO: is this good enough?
		s := v1.String()
		return &s, nil
	case reflect.Struct:
		if tm, ok := v.(time.Time); ok {
			switch tsmode {
			case dateType:
				_, offset := tm.Zone()
				tm = tm.Add(time.Second * time.Duration(offset))
				s := fmt.Sprintf("%d", tm.Unix()*1000)
				return &s, nil
			case timeType:
				s := fmt.Sprintf("%d",
					(tm.Hour()*3600+tm.Minute()*60+tm.Second())*1e9+tm.Nanosecond())
				return &s, nil
			case timestampNtzType, timestampLtzType:
				s := fmt.Sprintf("%d", tm.UnixNano())
				return &s, nil
			case timestampTzType:
				_, offset := tm.Zone()
				s := fmt.Sprintf("%v %v", tm.UnixNano(), offset/60+1440)
				return &s, nil
			}
		}
	}
	return nil, fmt.Errorf("unsupported type: %v", v1.Kind())
}

// extractTimestamp extracts the internal timestamp data to epoch time in seconds and milliseconds
func extractTimestamp(srcValue *string) (sec int64, nsec int64, err error) {
	logger.Debugf("SRC: %v", srcValue)
	var i int
	for i = 0; i < len(*srcValue); i++ {
		if (*srcValue)[i] == '.' {
			sec, err = strconv.ParseInt((*srcValue)[0:i], 10, 64)
			if err != nil {
				return 0, 0, err
			}
			break
		}
	}
	if i == len(*srcValue) {
		// no fraction
		sec, err = strconv.ParseInt(*srcValue, 10, 64)
		if err != nil {
			return 0, 0, err
		}
		nsec = 0
	} else {
		s := (*srcValue)[i+1:]
		nsec, err = strconv.ParseInt(s+strings.Repeat("0", 9-len(s)), 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	logger.Infof("sec: %v, nsec: %v", sec, nsec)
	return sec, nsec, nil
}

// stringToValue converts a pointer of string data to an arbitrary golang variable. This is mainly used in fetching
// data.
func stringToValue(dest *driver.Value, srcColumnMeta execResponseRowType, srcValue *string) error {
	if srcValue == nil {
		logger.Debugf("snowflake data type: %v, raw value: nil", srcColumnMeta.Type)
		*dest = nil
		return nil
	}
	logger.Debugf("snowflake data type: %v, raw value: %v", srcColumnMeta.Type, *srcValue)
	switch srcColumnMeta.Type {
	case "text", "fixed", "real", "variant", "object":
		*dest = *srcValue
		return nil
	case "date":
		v, err := strconv.ParseInt(*srcValue, 10, 64)
		if err != nil {
			return err
		}
		*dest = time.Unix(v*86400, 0).UTC()
		return nil
	case "time":
		sec, nsec, err := extractTimestamp(srcValue)
		if err != nil {
			return err
		}
		t0 := time.Time{}
		*dest = t0.Add(time.Duration(sec*1e9 + nsec))
		return nil
	case "timestamp_ntz":
		sec, nsec, err := extractTimestamp(srcValue)
		if err != nil {
			return err
		}
		*dest = time.Unix(sec, nsec).UTC()
		return nil
	case "timestamp_ltz":
		sec, nsec, err := extractTimestamp(srcValue)
		if err != nil {
			return err
		}
		*dest = time.Unix(sec, nsec)
		return nil
	case "timestamp_tz":
		logger.Debugf("tz: %v", *srcValue)

		tm := strings.Split(*srcValue, " ")
		if len(tm) != 2 {
			return &SnowflakeError{
				Number:   ErrInvalidTimestampTz,
				SQLState: SQLStateInvalidDataTimeFormat,
				Message:  fmt.Sprintf("invalid TIMESTAMP_TZ data. The value doesn't consist of two numeric values separated by a space: %v", *srcValue),
			}
		}
		sec, nsec, err := extractTimestamp(&tm[0])
		if err != nil {
			return err
		}
		offset, err := strconv.ParseInt(tm[1], 10, 64)
		if err != nil {
			return &SnowflakeError{
				Number:   ErrInvalidTimestampTz,
				SQLState: SQLStateInvalidDataTimeFormat,
				Message:  fmt.Sprintf("invalid TIMESTAMP_TZ data. The offset value is not integer: %v", tm[1]),
			}
		}
		loc := Location(int(offset) - 1440)
		tt := time.Unix(sec, nsec)
		*dest = tt.In(loc)
		return nil
	case "binary":
		b, err := hex.DecodeString(*srcValue)
		if err != nil {
			return &SnowflakeError{
				Number:   ErrInvalidBinaryHexForm,
				SQLState: SQLStateNumericValueOutOfRange,
				Message:  err.Error(),
			}
		}
		*dest = b
		return nil
	}
	*dest = *srcValue
	return nil
}

var decimalShift = new(big.Int).Exp(big.NewInt(2), big.NewInt(64), nil)

func intToBigFloat(val int64, scale int64) *big.Float {
	f := new(big.Float).SetInt64(val)
	s := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(scale), nil))
	return new(big.Float).Quo(f, s)
}

func decimalToBigInt(num decimal128.Num) *big.Int {
	high := new(big.Int).SetInt64(num.HighBits())
	low := new(big.Int).SetUint64(num.LowBits())
	return new(big.Int).Add(new(big.Int).Mul(high, decimalShift), low)
}

func decimalToBigFloat(num decimal128.Num, scale int64) *big.Float {
	f := new(big.Float).SetInt(decimalToBigInt(num))
	s := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(scale), nil))
	return new(big.Float).Quo(f, s)
}

func stringIntToDecimal(src string) (decimal128.Num, bool) {
	b, ok := new(big.Int).SetString(src, 10)
	if !ok {
		return decimal128.Num{}, ok
	}
	var high, low big.Int
	high.QuoRem(b, decimalShift, &low)
	return decimal128.New(high.Int64(), low.Uint64()), ok
}

func stringFloatToDecimal(src string, scale int64) (decimal128.Num, bool) {
	b, ok := new(big.Float).SetString(src)
	if !ok {
		return decimal128.Num{}, ok
	}
	s := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(scale), nil))
	n := new(big.Float).Mul(b, s)
	if !n.IsInt() {
		return decimal128.Num{}, false
	}
	var high, low, z big.Int
	n.Int(&z)
	high.QuoRem(&z, decimalShift, &low)
	return decimal128.New(high.Int64(), low.Uint64()), ok
}

// Arrow Interface (Column) converter. This is called when Arrow chunks are
// downloaded to convert to the corresponding row type.
func arrowToValue(
	destcol *[]snowflakeValue,
	srcColumnMeta execResponseRowType,
	srcValue array.Interface,
	higherPrecision bool) error {
	data := srcValue.Data()
	var err error
	if len(*destcol) != srcValue.Data().Len() {
		err = fmt.Errorf("array interface length mismatch")
	}
	logger.Debugf("snowflake data type: %v, arrow data type: %v", srcColumnMeta.Type, srcValue.DataType())

	switch getSnowflakeType(strings.ToUpper(srcColumnMeta.Type)) {
	case fixedType:
		switch srcValue.DataType().ID() {
		case arrow.DECIMAL:
			for i, num := range array.NewDecimal128Data(data).Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						if higherPrecision {
							(*destcol)[i] = decimalToBigInt(num)
						} else {
							(*destcol)[i] = decimalToBigInt(num).String()
						}
					} else {
						f := decimalToBigFloat(num, srcColumnMeta.Scale)
						if higherPrecision {
							(*destcol)[i] = f
						} else {
							(*destcol)[i] = fmt.Sprintf("%f", f)
						}
					}
				}
			}
		case arrow.INT64:
			for i, val := range array.NewInt64Data(data).Int64Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						if higherPrecision {
							(*destcol)[i] = val
						} else {
							(*destcol)[i] = fmt.Sprintf("%d", val)
						}
					} else {
						if higherPrecision {
							f := intToBigFloat(val, srcColumnMeta.Scale)
							(*destcol)[i] = f
						} else {
							(*destcol)[i] = fmt.Sprintf("%f", float64(val)/math.Pow10(int(srcColumnMeta.Scale)))
						}
					}
				}
			}
		case arrow.INT32:
			for i, val := range array.NewInt32Data(data).Int32Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						if higherPrecision {
							(*destcol)[i] = int64(val)
						} else {
							(*destcol)[i] = fmt.Sprintf("%d", val)
						}
					} else {
						if higherPrecision {
							f := intToBigFloat(int64(val), srcColumnMeta.Scale)
							(*destcol)[i] = f
						} else {
							(*destcol)[i] = fmt.Sprintf("%f", float64(val)/math.Pow10(int(srcColumnMeta.Scale)))
						}
					}
				}
			}
		case arrow.INT16:
			for i, val := range array.NewInt16Data(data).Int16Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						if higherPrecision {
							(*destcol)[i] = int64(val)
						} else {
							(*destcol)[i] = fmt.Sprintf("%d", val)
						}
					} else {
						if higherPrecision {
							f := intToBigFloat(int64(val), srcColumnMeta.Scale)
							(*destcol)[i] = f
						} else {
							(*destcol)[i] = fmt.Sprintf("%f", float64(val)/math.Pow10(int(srcColumnMeta.Scale)))
						}
					}
				}
			}
		case arrow.INT8:
			for i, val := range array.NewInt8Data(data).Int8Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						if higherPrecision {
							(*destcol)[i] = int64(val)
						} else {
							(*destcol)[i] = fmt.Sprintf("%d", val)
						}
					} else {
						if higherPrecision {
							f := intToBigFloat(int64(val), srcColumnMeta.Scale)
							(*destcol)[i] = f
						} else {
							(*destcol)[i] = fmt.Sprintf("%f", float64(val)/math.Pow10(int(srcColumnMeta.Scale)))
						}
					}
				}
			}
		}
		return err
	case booleanType:
		boolData := array.NewBooleanData(data)
		for i := range *destcol {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = boolData.Value(i)
			}
		}
		return err
	case realType:
		for i, flt64 := range array.NewFloat64Data(data).Float64Values() {
			if !srcValue.IsNull(i) {
				if higherPrecision {
					(*destcol)[i] = flt64
				} else {
					(*destcol)[i] = fmt.Sprintf("%f", flt64)
				}
			}
		}
		return err
	case textType, arrayType, variantType, objectType:
		strings := array.NewStringData(data)
		for i := range *destcol {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = strings.Value(i)
			}
		}
		return err
	case binaryType:
		binaryData := array.NewBinaryData(data)
		for i := range *destcol {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = binaryData.Value(i)
			}
		}
		return err
	case dateType:
		for i, date32 := range array.NewDate32Data(data).Date32Values() {
			if !srcValue.IsNull(i) {
				t0 := time.Unix(int64(date32)*86400, 0).UTC()
				(*destcol)[i] = t0
			}
		}
		return err
	case timeType:
		if srcValue.DataType().ID() == arrow.INT64 {
			for i, i64 := range array.NewInt64Data(data).Int64Values() {
				if !srcValue.IsNull(i) {
					t0 := time.Time{}
					(*destcol)[i] = t0.Add(time.Duration(i64))
				}
			}
		} else {
			for i, i32 := range array.NewInt32Data(data).Int32Values() {
				if !srcValue.IsNull(i) {
					t0 := time.Time{}
					(*destcol)[i] = t0.Add(time.Duration(int64(i32) * int64(math.Pow10(9-int(srcColumnMeta.Scale)))))
				}
			}
		}
		return err
	case timestampNtzType:
		if srcValue.DataType().ID() == arrow.STRUCT {
			structData := array.NewStructData(data)
			epoch := array.NewInt64Data(structData.Field(0).Data()).Int64Values()
			fraction := array.NewInt32Data(structData.Field(1).Data()).Int32Values()
			for i := range *destcol {
				if !srcValue.IsNull(i) {
					(*destcol)[i] = time.Unix(epoch[i], int64(fraction[i])).UTC()
				}
			}
		} else {
			for i, t := range array.NewInt64Data(data).Int64Values() {
				if !srcValue.IsNull(i) {
					(*destcol)[i] = time.Unix(0, t*int64(math.Pow10(9-int(srcColumnMeta.Scale)))).UTC()
				}
			}
		}
		return err
	case timestampLtzType:
		if srcValue.DataType().ID() == arrow.STRUCT {
			structData := array.NewStructData(data)
			epoch := array.NewInt64Data(structData.Field(0).Data()).Int64Values()
			fraction := array.NewInt32Data(structData.Field(1).Data()).Int32Values()
			for i := range *destcol {
				if !srcValue.IsNull(i) {
					(*destcol)[i] = time.Unix(epoch[i], int64(fraction[i]))
				}
			}
		} else {
			for i, t := range array.NewInt64Data(data).Int64Values() {
				if !srcValue.IsNull(i) {
					q := t / int64(math.Pow10(int(srcColumnMeta.Scale)))
					r := t % int64(math.Pow10(int(srcColumnMeta.Scale)))
					(*destcol)[i] = time.Unix(q, r)
				}
			}
		}
		return err
	case timestampTzType:
		structData := array.NewStructData(data)
		if structData.NumField() == 2 {
			epoch := array.NewInt64Data(structData.Field(0).Data()).Int64Values()
			timezone := array.NewInt32Data(structData.Field(1).Data()).Int32Values()
			for i := range *destcol {
				if !srcValue.IsNull(i) {
					loc := Location(int(timezone[i]) - 1440)
					tt := time.Unix(epoch[i], 0)
					(*destcol)[i] = tt.In(loc)
				}
			}
		} else {
			epoch := array.NewInt64Data(structData.Field(0).Data()).Int64Values()
			fraction := array.NewInt32Data(structData.Field(1).Data()).Int32Values()
			timezone := array.NewInt32Data(structData.Field(2).Data()).Int32Values()
			for i := range *destcol {
				if !srcValue.IsNull(i) {
					loc := Location(int(timezone[i]) - 1440)
					tt := time.Unix(epoch[i], int64(fraction[i]))
					(*destcol)[i] = tt.In(loc)
				}
			}
		}
		return err
	}

	return fmt.Errorf("unsupported data type")
}

type (
	intArray          []int
	int32Array        []int32
	int64Array        []int64
	float64Array      []float64
	float32Array      []float32
	boolArray         []bool
	stringArray       []string
	byteArray         [][]byte
	timestampNtzArray []time.Time
	timestampLtzArray []time.Time
	timestampTzArray  []time.Time
	dateArray         []time.Time
	timeArray         []time.Time
)

// Array takes in a column of a row to be inserted via array binding, bulk or
// otherwise, and converts it into a native snowflake type for binding
func Array(a interface{}, typ ...snowflakeType) interface{} {
	switch t := a.(type) {
	case []int:
		return (*intArray)(&t)
	case []int32:
		return (*int32Array)(&t)
	case []int64:
		return (*int64Array)(&t)
	case []float64:
		return (*float64Array)(&t)
	case []float32:
		return (*float32Array)(&t)
	case []bool:
		return (*boolArray)(&t)
	case []string:
		return (*stringArray)(&t)
	case [][]byte:
		return (*byteArray)(&t)
	case []time.Time:
		if len(typ) < 1 {
			return a
		}
		switch typ[0] {
		case timestampNtzType:
			return (*timestampNtzArray)(&t)
		case timestampLtzType:
			return (*timestampLtzArray)(&t)
		case timestampTzType:
			return (*timestampTzArray)(&t)
		case dateType:
			return (*dateArray)(&t)
		case timeType:
			return (*timeArray)(&t)
		default:
			return a
		}

	case *[]int:
		return (*intArray)(t)
	case *[]int32:
		return (*int32Array)(t)
	case *[]int64:
		return (*int64Array)(t)
	case *[]float64:
		return (*float64Array)(t)
	case *[]float32:
		return (*float32Array)(t)
	case *[]bool:
		return (*boolArray)(t)
	case *[]string:
		return (*stringArray)(t)
	case *[][]byte:
		return (*byteArray)(t)
	case *[]time.Time:
		if len(typ) < 1 {
			return a
		}
		switch typ[0] {
		case timestampNtzType:
			return (*timestampNtzArray)(t)
		case timestampLtzType:
			return (*timestampLtzArray)(t)
		case timestampTzType:
			return (*timestampTzArray)(t)
		case dateType:
			return (*dateArray)(t)
		case timeType:
			return (*timeArray)(t)
		default:
			return a
		}
	default:
		return a
	}
}

// snowflakeArrayToString converts the array binding to snowflake's native
// string type. The string value differs whether it's directly bound or
// uploaded via stream.
func snowflakeArrayToString(nv *driver.NamedValue, stream bool) (snowflakeType, []*string) {
	var t snowflakeType
	var arr []*string
	switch reflect.TypeOf(nv.Value) {
	case reflect.TypeOf(&intArray{}):
		t = fixedType
		a := nv.Value.(*intArray)
		for _, x := range *a {
			v := strconv.Itoa(x)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&int64Array{}):
		t = fixedType
		a := nv.Value.(*int64Array)
		for _, x := range *a {
			v := strconv.FormatInt(x, 10)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&int32Array{}):
		t = fixedType
		a := nv.Value.(*int32Array)
		for _, x := range *a {
			v := strconv.Itoa(int(x))
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&float64Array{}):
		t = realType
		a := nv.Value.(*float64Array)
		for _, x := range *a {
			v := fmt.Sprintf("%g", x)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&float32Array{}):
		t = realType
		a := nv.Value.(*float32Array)
		for _, x := range *a {
			v := fmt.Sprintf("%g", x)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&boolArray{}):
		t = booleanType
		a := nv.Value.(*boolArray)
		for _, x := range *a {
			v := strconv.FormatBool(x)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&stringArray{}):
		t = textType
		a := nv.Value.(*stringArray)
		for _, x := range *a {
			v := x // necessary for address to be not overwritten
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&byteArray{}):
		t = binaryType
		a := nv.Value.(*byteArray)
		for _, x := range *a {
			v := hex.EncodeToString(x)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&timestampNtzArray{}):
		t = timestampNtzType
		a := nv.Value.(*timestampNtzArray)
		for _, x := range *a {
			v := strconv.FormatInt(x.UnixNano(), 10)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&timestampLtzArray{}):
		t = timestampLtzType
		a := nv.Value.(*timestampLtzArray)
		for _, x := range *a {
			v := strconv.FormatInt(x.UnixNano(), 10)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&timestampTzArray{}):
		t = timestampTzType
		a := nv.Value.(*timestampTzArray)
		for _, x := range *a {
			var v string
			if stream {
				v = x.Format(format)
			} else {
				_, offset := x.Zone()
				v = fmt.Sprintf("%v %v", x.UnixNano(), offset/60+1440)
			}
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&dateArray{}):
		t = dateType
		a := nv.Value.(*dateArray)
		for _, x := range *a {
			_, offset := x.Zone()
			x = x.Add(time.Second * time.Duration(offset))
			v := fmt.Sprintf("%d", x.Unix()*1000)
			arr = append(arr, &v)
		}
	case reflect.TypeOf(&timeArray{}):
		t = timeType
		a := nv.Value.(*timeArray)
		for _, x := range *a {
			var v string
			if stream {
				v = x.Format(format[11:19])
			} else {
				h, m, s := x.Clock()
				tm := int64(h)*int64(time.Hour) + int64(m)*int64(time.Minute) + int64(s)*int64(time.Second) + int64(x.Nanosecond())
				v = strconv.FormatInt(tm, 10)
			}
			arr = append(arr, &v)
		}
	default:
		return unSupportedType, nil
	}
	return t, arr
}

func higherPrecisionEnabled(ctx context.Context) bool {
	v := ctx.Value(enableHigherPrecision)
	if v == nil {
		return false
	}
	d, ok := v.(bool)
	return ok && d
}
