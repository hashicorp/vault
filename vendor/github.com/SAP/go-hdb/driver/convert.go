package driver

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	p "github.com/SAP/go-hdb/driver/internal/protocol"
	"github.com/SAP/go-hdb/driver/internal/protocol/levenshtein"
	"golang.org/x/text/transform"
)

// TODO: test.
func reorderNVArgs(pos int, name string, nvargs []driver.NamedValue) {
	for i := pos; i < len(nvargs); i++ {
		if nvargs[i].Name != "" && nvargs[i].Name == name {
			tmp := nvargs[i]
			for j := i; j > pos; j-- {
				nvargs[j] = nvargs[j-1]
			}
			nvargs[pos] = tmp
		}
	}
}

// This function is take from the database/sql package.
// The Elem() test is not needed bacause the function is only
// called for values implementing the driver.Valuer interface.
func callValuerValue(vr driver.Valuer) (v driver.Value, err error) {
	if rv := reflect.ValueOf(vr); rv.Kind() == reflect.Pointer && rv.IsNil() {
		// && rv.Type().Elem().Implements(valuerReflectType) {
		return nil, nil
	}
	return vr.Value()
}

func convertArg(field *p.ParameterField, arg any, cesu8Encoder transform.Transformer) (any, error) {
	// let fields with own value converter convert themselves first (e.g. NullInt64, ...)
	// .check nested Value converters as well (e.g. sql.Null[T] has driver.Decimal as value)
	for {
		valuer, ok := arg.(driver.Valuer)
		if !ok {
			break
		}
		var err error
		arg, err = callValuerValue(valuer)
		if err != nil {
			return nil, err
		}
	}

	// convert field
	return field.Convert(arg, cesu8Encoder)
}

/*
convertExecArgs
  - all fields need to be input fields
  - out parameters are not supported
  - named parameters are not supported
*/
func convertExecArgs(fields []*p.ParameterField, nvargs []driver.NamedValue, cesu8Encoder transform.Transformer, lobChunkSize int) ([]int, error) {
	numField := len(fields)
	if (len(nvargs) % numField) != 0 {
		return nil, fmt.Errorf("invalid number of arguments %d - multiple of %d expected", len(nvargs), numField)
	}
	numRow := len(nvargs) / numField
	addLobDataRecs := []int{}

	for i := range numRow {
		hasAddLobData := false
		for j, field := range fields {
			nvarg := &nvargs[(i*numField)+j]

			if field.Out() {
				return nil, fmt.Errorf("invalid parameter %s - output not allowed", field)
			}
			if _, ok := nvarg.Value.(sql.Out); ok {
				return nil, fmt.Errorf("invalid argument %v - output not allowed", nvarg)
			}
			if nvarg.Name != "" {
				return nil, fmt.Errorf("invalid argument %s - named parameters not supported", nvarg.Name)
			}
			var err error
			if nvarg.Value, err = convertArg(field, nvarg.Value, cesu8Encoder); err != nil {
				return nil, fmt.Errorf("field %s conversion error - %w", field, err)
			}
			// fetch first lob chunk
			if lobInDescr, ok := nvarg.Value.(*p.LobInDescr); ok {
				if err := lobInDescr.FetchNext(lobChunkSize); err != nil {
					return nil, err
				}
				if !lobInDescr.IsLastData() {
					hasAddLobData = true
				}
			}
		}
		if hasAddLobData || i == numRow-1 {
			addLobDataRecs = append(addLobDataRecs, i)
		}
	}
	return addLobDataRecs, nil
}

/*
_convertQueryArgs
  - all fields need to be input fields
  - out parameters are not supported
  - named parameters are not supported
*/
func convertQueryArgs(fields []*p.ParameterField, nvargs []driver.NamedValue, cesu8Encoder transform.Transformer, lobChunkSize int) error {
	if len(nvargs) != len(fields) {
		return fmt.Errorf("invalid number of arguments %d - %d expected", len(nvargs), len(fields))
	}

	for i, field := range fields {
		nvarg := &nvargs[i]
		if field.Out() {
			return fmt.Errorf("invalid parameter %s - output not allowed", field)
		}
		if _, ok := nvarg.Value.(sql.Out); ok {
			return fmt.Errorf("invalid argument %v - output not allowed", nvarg)
		}
		if nvarg.Name != "" {
			return fmt.Errorf("invalid argument %s - named parameters not supported", nvarg.Name)
		}
		var err error
		if nvarg.Value, err = convertArg(field, nvarg.Value, cesu8Encoder); err != nil {
			return fmt.Errorf("field %s conversion error - %w", field, err)
		}
		// fetch first lob chunk
		if lobInDescr, ok := nvarg.Value.(*p.LobInDescr); ok {
			if err := lobInDescr.FetchNext(lobChunkSize); err != nil {
				return err
			}
		}
	}
	return nil
}

// convertCallArgs
// - fields could be input or output fields
// - number of args needs to be equal to number of fields
// - named parameters are supported

type callArgs struct {
	inFields, outFields []*p.ParameterField
	inArgs, outArgs     []driver.NamedValue
}

func newCallArgs() *callArgs {
	return &callArgs{
		inFields:  []*p.ParameterField{},
		outFields: []*p.ParameterField{},
		inArgs:    []driver.NamedValue{},
		outArgs:   []driver.NamedValue{},
	}
}

func convertCallArgs(fields []*p.ParameterField, nvargs []driver.NamedValue, cesu8Encoder transform.Transformer, lobChunkSize int) (*callArgs, error) {
	callArgs := newCallArgs()

	if len(nvargs) < len(fields) { // number of fields needs to match number of args or be greater (add table output args)
		return nil, fmt.Errorf("invalid number of arguments %d - %d expected", len(nvargs), len(fields))
	}

	prmnvargs := nvargs[:len(fields)]

	for i, field := range fields {
		reorderNVArgs(i, field.Name(), prmnvargs)

		nvarg := &prmnvargs[i]

		if nvarg.Name != "" && nvarg.Name != field.Name() {
			return nil, fmt.Errorf("invalid argument name %s - did you mean %s?",
				nvarg.Name,
				levenshtein.MinString(fields, func(field *p.ParameterField) string { return field.Name() }, nvarg.Name, false),
			)
		}

		out, isOut := nvarg.Value.(sql.Out)

		var err error
		if field.In() {
			if isOut {
				if !out.In {
					return nil, fmt.Errorf("argument field %s mismatch - use in argument with out field", field)
				}
				if out.Dest, err = convertArg(field, out.Dest, cesu8Encoder); err != nil {
					return nil, fmt.Errorf("field %s conversion error - %w", field, err)
				}
			} else {
				if nvarg.Value, err = convertArg(field, nvarg.Value, cesu8Encoder); err != nil {
					return nil, fmt.Errorf("field %s conversion error - %w", field, err)
				}
			}
			// fetch first lob chunk
			if lobInDescr, ok := nvarg.Value.(*p.LobInDescr); ok {
				if err := lobInDescr.FetchNext(lobChunkSize); err != nil {
					return nil, err
				}
			}
			callArgs.inArgs = append(callArgs.inArgs, *nvarg)
			callArgs.inFields = append(callArgs.inFields, field)
		}

		if field.Out() {
			if !isOut {
				return nil, fmt.Errorf("argument field %s mismatch - use out argument with non-out field", field)
			}
			if _, ok := out.Dest.(*sql.Rows); ok {
				return nil, fmt.Errorf("invalid output parameter type %T", out.Dest)
			}
			callArgs.outArgs = append(callArgs.outArgs, *nvarg)
			callArgs.outFields = append(callArgs.outFields, field)
		}
	}

	// table output args
	for i := len(fields); i < len(nvargs); i++ {
		nvarg := &nvargs[i]
		out, ok := nvarg.Value.(sql.Out)
		if !ok {
			return nil, fmt.Errorf("invalid parameter type %T at %d - output parameter expected", nvarg.Value, i)
		}
		if _, ok := out.Dest.(*sql.Rows); !ok {
			return nil, fmt.Errorf("invalid output parameter %T at %d - sql.Rows expected", out.Dest, i)
		}
		callArgs.outArgs = append(callArgs.outArgs, *nvarg)
	}
	return callArgs, nil
}
