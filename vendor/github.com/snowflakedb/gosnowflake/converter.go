// Copyright (c) 2017-2020 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"github.com/apache/arrow/go/arrow"
	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/decimal128"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// goTypeToSnowflake translates Go data type to Snowflake data type.
func goTypeToSnowflake(v driver.Value, tsmode string) string {
	switch v := v.(type) {
	case int64:
		return "FIXED"
	case float64:
		return "REAL"
	case bool:
		return "BOOLEAN"
	case string:
		return "TEXT"
	case []byte:
		if tsmode == "BINARY" {
			return "BINARY" // may be redundant but ensures BINARY type
		}
		if v == nil || len(v) != 1 {
			return "TEXT" // invalid byte array. won't take as BINARY
		}
		_, err := dataTypeMode(v)
		if err != nil {
			return "TEXT" // not supported dataType
		}
		return "CHANGE_TYPE"
	case []int, []int64, []float64, []bool, []string:
		return "ARRAY"
	case time.Time:
		return tsmode
	}
	return "TEXT"
}

// snowflakeTypeToGo translates Snowflake data type to Go data type.
func snowflakeTypeToGo(dbtype string, scale int64) reflect.Type {
	switch dbtype {
	case "fixed":
		if scale == 0 {
			return reflect.TypeOf(int64(0))
		}
		return reflect.TypeOf(float64(0))
	case "real":
		return reflect.TypeOf(float64(0))
	case "text", "variant", "object", "array":
		return reflect.TypeOf("")
	case "date", "time", "timestamp_ltz", "timestamp_ntz", "timestamp_tz":
		return reflect.TypeOf(time.Now())
	case "binary":
		return reflect.TypeOf([]byte{})
	case "boolean":
		return reflect.TypeOf(true)
	}
	logger.Errorf("unsupported dbtype is specified. %v", dbtype)
	return reflect.TypeOf("")
}

// valueToString converts arbitrary golang type to a string. This is mainly used in binding data with placeholders
// in queries.
func valueToString(v driver.Value, tsmode string) (*string, error) {
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
			if tsmode == "BINARY" {
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
			case "DATE":
				_, offset := tm.Zone()
				tm = tm.Add(time.Second * time.Duration(offset))
				s := fmt.Sprintf("%d", tm.Unix()*1000)
				return &s, nil
			case "TIME":
				s := fmt.Sprintf("%d",
					(tm.Hour()*3600+tm.Minute()*60+tm.Second())*1e9+tm.Nanosecond())
				return &s, nil
			case "TIMESTAMP_NTZ", "TIMESTAMP_LTZ":
				s := fmt.Sprintf("%d", tm.UnixNano())
				return &s, nil
			case "TIMESTAMP_TZ":
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

func arrayToString(v driver.Value) (string, []string) {
	var t string
	var arr []string
	switch a := v.(type) {
	case []int:
		t = "FIXED"
		for _, x := range a {
			arr = append(arr, strconv.Itoa(x))
		}
	case []int64:
		t = "FIXED"
		for _, x := range a {
			arr = append(arr, strconv.Itoa(int(x)))
		}
	case []float64:
		t = "REAL"
		for _, x := range a {
			arr = append(arr, fmt.Sprintf("%g", x))
		}
	case []bool:
		t = "BOOLEAN"
		for _, x := range a {
			arr = append(arr, strconv.FormatBool(x))
		}
	case []string:
		t = "TEXT"
		arr = a
	}
	return t, arr
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

// Arrow Interface (Column) converter. This is called when Arrow chunks are downloaded to convert to the corresponding
// row type.
func arrowToValue(destcol *[]snowflakeValue, srcColumnMeta execResponseRowType, srcValue array.Interface) error {
	data := srcValue.Data()
	var err error
	if len(*destcol) != srcValue.Data().Len() {
		err = fmt.Errorf("array interface length mismatch")
	}
	logger.Debugf("snowflake data type: %v, arrow data type: %v", srcColumnMeta.Type, srcValue.DataType())

	switch strings.ToUpper(srcColumnMeta.Type) {
	case "FIXED":
		switch srcValue.DataType().ID() {
		case arrow.DECIMAL:
			for i, num := range array.NewDecimal128Data(data).Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						(*destcol)[i] = decimalToBigInt(num)
					} else {
						(*destcol)[i] = decimalToBigFloat(num, srcColumnMeta.Scale)
					}
				}
			}
		case arrow.INT64:
			for i, val := range array.NewInt64Data(data).Int64Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						(*destcol)[i] = val
					} else {
						f := intToBigFloat(val, srcColumnMeta.Scale)
						(*destcol)[i] = f
					}
				}
			}
		case arrow.INT32:
			for i, val := range array.NewInt32Data(data).Int32Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						(*destcol)[i] = int64(val)
					} else {
						f := intToBigFloat(int64(val), srcColumnMeta.Scale)
						(*destcol)[i] = f
					}
				}
			}
		case arrow.INT16:
			for i, val := range array.NewInt16Data(data).Int16Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						(*destcol)[i] = int64(val)
					} else {
						f := intToBigFloat(int64(val), srcColumnMeta.Scale)
						(*destcol)[i] = f
					}
				}
			}
		case arrow.INT8:
			for i, val := range array.NewInt8Data(data).Int8Values() {
				if !srcValue.IsNull(i) {
					if srcColumnMeta.Scale == 0 {
						(*destcol)[i] = int64(val)
					} else {
						f := intToBigFloat(int64(val), srcColumnMeta.Scale)
						(*destcol)[i] = f
					}
				}
			}
		}
		return err
	case "BOOLEAN":
		boolData := array.NewBooleanData(data)
		for i := range *destcol {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = boolData.Value(i)
			}
		}
		return err
	case "REAL":
		for i, float64 := range array.NewFloat64Data(data).Float64Values() {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = float64
			}
		}
		return err
	case "TEXT", "ARRAY", "VARIANT", "OBJECT":
		strings := array.NewStringData(data)
		for i := range *destcol {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = strings.Value(i)
			}
		}
		return err
	case "BINARY":
		binaryData := array.NewBinaryData(data)
		for i := range *destcol {
			if !srcValue.IsNull(i) {
				(*destcol)[i] = binaryData.Value(i)
			}
		}
		return err
	case "DATE":
		for i, date32 := range array.NewDate32Data(data).Date32Values() {
			if !srcValue.IsNull(i) {
				t0 := time.Unix(int64(date32)*86400, 0).UTC()
				(*destcol)[i] = t0
			}
		}
		return err
	case "TIME":
		if srcValue.DataType().ID() == arrow.INT64 {
			for i, int64 := range array.NewInt64Data(data).Int64Values() {
				if !srcValue.IsNull(i) {
					t0 := time.Time{}
					(*destcol)[i] = t0.Add(time.Duration(int64))
				}
			}
		} else {
			for i, int32 := range array.NewInt32Data(data).Int32Values() {
				if !srcValue.IsNull(i) {
					t0 := time.Time{}
					(*destcol)[i] = t0.Add(time.Duration(int64(int32) * int64(math.Pow10(9-int(srcColumnMeta.Scale)))))
				}
			}
		}
		return err
	case "TIMESTAMP_NTZ":
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
	case "TIMESTAMP_LTZ":
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
	case "TIMESTAMP_TZ":
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

	err = fmt.Errorf("unsupported data type")
	return err
}
