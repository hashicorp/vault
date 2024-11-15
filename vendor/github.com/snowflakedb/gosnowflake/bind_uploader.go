// Copyright (c) 2021-2022 Snowflake Computing Inc. All rights reserved.

package gosnowflake

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	bindStageName            = "SYSTEM$BIND"
	createTemporaryStageStmt = "CREATE OR REPLACE TEMPORARY STAGE " + bindStageName +
		" file_format=" + "(type=csv field_optionally_enclosed_by='\"')"

	// size (in bytes) of max input stream (10MB default) as per JDBC specs
	inputStreamBufferSize = 1024 * 1024 * 10
)

type bindUploader struct {
	ctx            context.Context
	sc             *snowflakeConn
	stagePath      string
	fileCount      int
	arrayBindStage string
}

type bindingSchema struct {
	Typ      string          `json:"type"`
	Nullable bool            `json:"nullable"`
	Fields   []fieldMetadata `json:"fields"`
}

type bindingValue struct {
	value  *string
	format string
	schema *bindingSchema
}

func (bu *bindUploader) upload(bindings []driver.NamedValue) (*execResponse, error) {
	bindingRows, err := bu.buildRowsAsBytes(bindings)
	if err != nil {
		return nil, err
	}
	startIdx, numBytes, rowNum := 0, 0, 0
	bu.fileCount = 0
	var data *execResponse
	for rowNum < len(bindingRows) {
		for numBytes < inputStreamBufferSize && rowNum < len(bindingRows) {
			numBytes += len(bindingRows[rowNum])
			rowNum++
		}
		// concatenate all byte arrays into 1 and put into input stream
		var b bytes.Buffer
		b.Grow(numBytes)
		for i := startIdx; i < rowNum; i++ {
			b.Write(bindingRows[i])
		}

		bu.fileCount++
		data, err = bu.uploadStreamInternal(&b, bu.fileCount, true)
		if err != nil {
			return nil, err
		}
		startIdx = rowNum
		numBytes = 0
	}
	return data, nil
}

func (bu *bindUploader) uploadStreamInternal(
	inputStream *bytes.Buffer,
	dstFileName int,
	compressData bool) (
	*execResponse, error) {
	if err := bu.createStageIfNeeded(); err != nil {
		return nil, err
	}
	stageName := bu.stagePath
	if stageName == "" {
		return nil, (&SnowflakeError{
			Number:  ErrBindUpload,
			Message: "stage name is null",
		}).exceptionTelemetry(bu.sc)
	}

	// use a placeholder for source file
	putCommand := fmt.Sprintf("put 'file:///tmp/placeholder/%v' '%v' overwrite=true", dstFileName, stageName)
	// for Windows queries
	putCommand = strings.ReplaceAll(putCommand, "\\", "\\\\")
	// prepare context for PUT command
	ctx := WithFileStream(bu.ctx, inputStream)
	ctx = WithFileTransferOptions(ctx, &SnowflakeFileTransferOptions{
		compressSourceFromStream: compressData})
	return bu.sc.exec(ctx, putCommand, false, true, false, []driver.NamedValue{})
}

func (bu *bindUploader) createStageIfNeeded() error {
	if bu.arrayBindStage != "" {
		return nil
	}
	data, err := bu.sc.exec(bu.ctx, createTemporaryStageStmt, false, false, false, []driver.NamedValue{})
	if err != nil {
		newThreshold := "0"
		bu.sc.cfg.Params[sessionArrayBindStageThreshold] = &newThreshold
		return err
	}
	if !data.Success {
		code, err := strconv.Atoi(data.Code)
		if err != nil {
			return err
		}
		return (&SnowflakeError{
			Number:   code,
			SQLState: data.Data.SQLState,
			Message:  err.Error(),
			QueryID:  data.Data.QueryID,
		}).exceptionTelemetry(bu.sc)
	}
	bu.arrayBindStage = bindStageName
	return nil
}

// transpose the columns to rows and write them to a list of bytes
func (bu *bindUploader) buildRowsAsBytes(columns []driver.NamedValue) ([][]byte, error) {
	numColumns := len(columns)
	if columns[0].Value == nil {
		return nil, (&SnowflakeError{
			Number:  ErrBindSerialization,
			Message: "no binds found in the first column",
		}).exceptionTelemetry(bu.sc)
	}

	_, column := snowflakeArrayToString(&columns[0], true)
	numRows := len(column)
	csvRows := make([][]byte, 0)
	rows := make([][]interface{}, 0)
	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		rows = append(rows, make([]interface{}, numColumns))
	}

	for rowIdx := 0; rowIdx < numRows; rowIdx++ {
		if column[rowIdx] == nil {
			rows[rowIdx][0] = column[rowIdx]
		} else {
			rows[rowIdx][0] = *column[rowIdx]
		}
	}
	for colIdx := 1; colIdx < numColumns; colIdx++ {
		_, column = snowflakeArrayToString(&columns[colIdx], true)
		iNumRows := len(column)
		if iNumRows != numRows {
			return nil, (&SnowflakeError{
				Number:      ErrBindSerialization,
				Message:     errMsgBindColumnMismatch,
				MessageArgs: []interface{}{colIdx, iNumRows, numRows},
			}).exceptionTelemetry(bu.sc)
		}
		for rowIdx := 0; rowIdx < numRows; rowIdx++ {
			// length of column = number of rows
			if column[rowIdx] == nil {
				rows[rowIdx][colIdx] = column[rowIdx]
			} else {
				rows[rowIdx][colIdx] = *column[rowIdx]
			}
		}
	}
	for _, row := range rows {
		csvRows = append(csvRows, bu.createCSVRecord(row))
	}
	return csvRows, nil
}

func (bu *bindUploader) createCSVRecord(data []interface{}) []byte {
	var b strings.Builder
	b.Grow(1024)
	for i := 0; i < len(data); i++ {
		if i > 0 {
			b.WriteString(",")
		}
		value, ok := data[i].(string)
		if ok {
			b.WriteString(escapeForCSV(value))
		} else if !reflect.ValueOf(data[i]).IsNil() {
			logger.WithContext(bu.ctx).Debugf("Cannot convert value to string in createCSVRecord. value: %v", data[i])
		}
	}
	b.WriteString("\n")
	return []byte(b.String())
}

func (sc *snowflakeConn) processBindings(
	ctx context.Context,
	bindings []driver.NamedValue,
	describeOnly bool,
	requestID UUID,
	req *execRequest) error {
	arrayBindThreshold := sc.getArrayBindStageThreshold()
	numBinds := arrayBindValueCount(bindings)
	if 0 < arrayBindThreshold && arrayBindThreshold <= numBinds && !describeOnly && isArrayBind(bindings) {
		uploader := bindUploader{
			sc:        sc,
			ctx:       ctx,
			stagePath: "@" + bindStageName + "/" + requestID.String(),
		}
		_, err := uploader.upload(bindings)
		if err != nil {
			return err
		}
		req.Bindings = nil
		req.BindStage = uploader.stagePath
	} else {
		var err error
		req.Bindings, err = getBindValues(bindings, sc.cfg.Params)
		if err != nil {
			return err
		}
		req.BindStage = ""
	}
	return nil
}

func getBindValues(bindings []driver.NamedValue, params map[string]*string) (map[string]execBindParameter, error) {
	tsmode := timestampNtzType
	idx := 1
	var err error
	bindValues := make(map[string]execBindParameter, len(bindings))
	for _, binding := range bindings {
		if tnt, ok := binding.Value.(TypedNullTime); ok {
			tsmode = convertTzTypeToSnowflakeType(tnt.TzType)
			binding.Value = tnt.Time
		}
		t := goTypeToSnowflake(binding.Value, tsmode)
		if t == changeType {
			tsmode, err = dataTypeMode(binding.Value)
			if err != nil {
				return nil, err
			}
		} else {
			var val interface{}
			var bv bindingValue
			if t == sliceType {
				// retrieve array binding data
				t, val = snowflakeArrayToString(&binding, false)
			} else {
				bv, err = valueToString(binding.Value, tsmode, params)
				val = bv.value
				if err != nil {
					return nil, err
				}
			}
			if t == nullType || t == unSupportedType {
				t = textType // if null or not supported, pass to GS as text
			} else if t == nilObjectType || t == mapType || t == nilMapType {
				t = objectType
			} else if t == nilArrayType {
				t = arrayType
			}
			bindValues[bindingName(binding, idx)] = execBindParameter{
				Type:   t.String(),
				Value:  val,
				Format: bv.format,
				Schema: bv.schema,
			}
			idx++
		}
	}
	return bindValues, nil
}

func bindingName(nv driver.NamedValue, idx int) string {
	if nv.Name != "" {
		return nv.Name
	}
	return strconv.Itoa(idx)
}

func arrayBindValueCount(bindValues []driver.NamedValue) int {
	if !isArrayBind(bindValues) {
		return 0
	}
	_, arr := snowflakeArrayToString(&bindValues[0], false)
	return len(bindValues) * len(arr)
}

func isArrayBind(bindings []driver.NamedValue) bool {
	if len(bindings) == 0 {
		return false
	}
	for _, binding := range bindings {
		if supported := supportedArrayBind(&binding); !supported {
			return false
		}
	}
	return true
}

func supportedArrayBind(nv *driver.NamedValue) bool {
	switch reflect.TypeOf(nv.Value) {
	case reflect.TypeOf(&intArray{}), reflect.TypeOf(&int32Array{}),
		reflect.TypeOf(&int64Array{}), reflect.TypeOf(&float64Array{}),
		reflect.TypeOf(&float32Array{}), reflect.TypeOf(&boolArray{}),
		reflect.TypeOf(&stringArray{}), reflect.TypeOf(&byteArray{}),
		reflect.TypeOf(&timestampNtzArray{}), reflect.TypeOf(&timestampLtzArray{}),
		reflect.TypeOf(&timestampTzArray{}), reflect.TypeOf(&dateArray{}),
		reflect.TypeOf(&timeArray{}):
		return true
	case reflect.TypeOf([]uint8{}):
		// internal binding ts mode
		val, ok := nv.Value.([]uint8)
		if !ok {
			return ok
		}
		if len(val) == 0 {
			return true // for null binds
		}
		if fixedType <= snowflakeType(val[0]) && snowflakeType(val[0]) <= unSupportedType {
			return true
		}
		return false
	default:
		// TODO SNOW-176486 variant, object, array

		// Support for bulk array binding insertion using []interface{}
		if isInterfaceArrayBinding(nv.Value) {
			return true
		}
		return false
	}
}

func supportedNullBind(nv *driver.NamedValue) bool {
	switch reflect.TypeOf(nv.Value) {
	case reflect.TypeOf(sql.NullString{}), reflect.TypeOf(sql.NullInt64{}),
		reflect.TypeOf(sql.NullBool{}), reflect.TypeOf(sql.NullFloat64{}), reflect.TypeOf(TypedNullTime{}):
		return true
	}
	return false
}

func supportedStructuredObjectWriterBind(nv *driver.NamedValue) bool {
	if _, ok := nv.Value.(StructuredObjectWriter); ok {
		return true
	}
	_, ok := nv.Value.(reflect.Type)
	return ok
}

func supportedStructuredArrayBind(nv *driver.NamedValue) bool {
	typ := reflect.TypeOf(nv.Value)
	return typ != nil && (typ.Kind() == reflect.Array || typ.Kind() == reflect.Slice)
}

func supportedStructuredMapBind(nv *driver.NamedValue) bool {
	typ := reflect.TypeOf(nv.Value)
	return typ != nil && (typ.Kind() == reflect.Map || typ == reflect.TypeOf(NilMapTypes{}))
}
