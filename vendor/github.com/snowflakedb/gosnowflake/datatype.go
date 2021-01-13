// Copyright (c) 2017-2019 Snowflake Computing Inc. All right reserved.

package gosnowflake

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
)

const (
	fixedType byte = iota
	realType
	textType
	dateType
	variantType
	timestampLtzType
	timestampNtzType
	timestampTzType
	objectType
	arrayType
	binaryType
	timeType
	booleanType
)

var (
	// DataTypeFixed is a FIXED datatype.
	DataTypeFixed = []byte{fixedType}
	// DataTypeReal is a REAL datatype.
	DataTypeReal = []byte{realType}
	// DataTypeText is a TEXT datatype.
	DataTypeText = []byte{textType}
	// DataTypeDate is a Date datatype.
	DataTypeDate = []byte{dateType}
	// DataTypeVariant is a TEXT datatype.
	DataTypeVariant = []byte{variantType}
	// DataTypeTimestampLtz is a TIMESTAMP_LTZ datatype.
	DataTypeTimestampLtz = []byte{timestampLtzType}
	// DataTypeTimestampNtz is a TIMESTAMP_NTZ datatype.
	DataTypeTimestampNtz = []byte{timestampNtzType}
	// DataTypeTimestampTz is a TIMESTAMP_TZ datatype.
	DataTypeTimestampTz = []byte{timestampTzType}
	// DataTypeObject is a OBJECT datatype.
	DataTypeObject = []byte{objectType}
	// DataTypeArray is a ARRAY datatype.
	DataTypeArray = []byte{arrayType}
	// DataTypeBinary is a BINARY datatype.
	DataTypeBinary = []byte{binaryType}
	// DataTypeTime is a TIME datatype.
	DataTypeTime = []byte{timeType}
	// DataTypeBoolean is a BOOLEAN datatype.
	DataTypeBoolean = []byte{booleanType}
)

// dataTypeMode returns the subsequent data type in a string representation.
func dataTypeMode(v driver.Value) (tsmode string, err error) {
	if bd, ok := v.([]byte); ok {
		switch {
		case bytes.Equal(bd, DataTypeDate):
			tsmode = "DATE"
		case bytes.Equal(bd, DataTypeTime):
			tsmode = "TIME"
		case bytes.Equal(bd, DataTypeTimestampLtz):
			tsmode = "TIMESTAMP_LTZ"
		case bytes.Equal(bd, DataTypeTimestampNtz):
			tsmode = "TIMESTAMP_NTZ"
		case bytes.Equal(bd, DataTypeTimestampTz):
			tsmode = "TIMESTAMP_TZ"
		case bytes.Equal(bd, DataTypeBinary):
			tsmode = "BINARY"
		default:
			return "", fmt.Errorf(errMsgInvalidByteArray, v)
		}
	} else {
		return "", fmt.Errorf(errMsgInvalidByteArray, v)
	}
	return tsmode, nil
}

// SnowflakeParameter includes the columns output from SHOW PARAMETER command.
type SnowflakeParameter struct {
	Key                       string
	Value                     string
	Default                   string
	Level                     string
	Description               string
	SetByUser                 string
	SetInJob                  string
	SetOn                     string
	SetByThreadID             string
	SetByThreadName           string
	SetByClass                string
	ParameterComment          string
	Type                      string
	IsExpired                 string
	ExpiresAt                 string
	SetByControllingParameter string
	ActivateVersion           string
	PartialRollout            string
	Unknown                   string // Reserve for added parameter
}

func populateSnowflakeParameter(colname string, p *SnowflakeParameter) interface{} {
	switch colname {
	case "key":
		return &p.Key
	case "value":
		return &p.Value
	case "default":
		return &p.Default
	case "level":
		return &p.Level
	case "description":
		return &p.Description
	case "set_by_user":
		return &p.SetByUser
	case "set_in_job":
		return &p.SetInJob
	case "set_on":
		return &p.SetOn
	case "set_by_thread_id":
		return &p.SetByThreadID
	case "set_by_thread_name":
		return &p.SetByThreadName
	case "set_by_class":
		return &p.SetByClass
	case "parameter_comment":
		return &p.ParameterComment
	case "type":
		return &p.Type
	case "is_expired":
		return &p.IsExpired
	case "expires_at":
		return &p.ExpiresAt
	case "set_by_controlling_parameter":
		return &p.SetByControllingParameter
	case "activate_version":
		return &p.ActivateVersion
	case "partial_rollout":
		return &p.PartialRollout
	default:
		debugPanicf("unknown type: %v", colname)
		return &p.Unknown
	}
}

// ScanSnowflakeParameter binds SnowflakeParameter variable with an array of column buffer.
func ScanSnowflakeParameter(rows *sql.Rows) (*SnowflakeParameter, error) {
	var err error
	var columns []string
	columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}
	colNum := len(columns)
	p := SnowflakeParameter{}
	cols := make([]interface{}, colNum)
	for i := 0; i < colNum; i++ {
		cols[i] = populateSnowflakeParameter(columns[i], &p)
	}
	err = rows.Scan(cols...)
	return &p, err
}
