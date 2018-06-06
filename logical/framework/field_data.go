package framework

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/mitchellh/mapstructure"
)

// FieldData is the structure passed to the callback to handle a path
// containing the populated parameters for fields. This should be used
// instead of the raw (*vault.Request).Data to access data in a type-safe
// way.
type FieldData struct {
	Raw    map[string]interface{}
	Schema map[string]*FieldSchema
}

// Validate cycles through raw data and validate conversions in
// the schema, so we don't get an error/panic later when
// trying to get data out.  Data not in the schema is not
// an error at this point, so we don't worry about it.
func (d *FieldData) Validate() error {
	for field, value := range d.Raw {

		schema, ok := d.Schema[field]
		if !ok {
			continue
		}

		switch schema.Type {
		case TypeBool, TypeInt, TypeMap, TypeDurationSecond, TypeString, TypeLowerCaseString,
			TypeNameString, TypeSlice, TypeStringSlice, TypeCommaStringSlice,
			TypeKVPairs, TypeCommaIntSlice:
			_, _, err := d.getPrimitive(field, schema)
			if err != nil {
				return errwrap.Wrapf(fmt.Sprintf("error converting input %v for field %q: {{err}}", value, field), err)
			}
		default:
			return fmt.Errorf("unknown field type %q for field %q", schema.Type, field)
		}
	}

	return nil
}

// Get gets the value for the given field. If the key is an invalid field,
// FieldData will panic. If you want a safer version of this method, use
// GetOk. If the field k is not set, the default value (if set) will be
// returned, otherwise the zero value will be returned.
func (d *FieldData) Get(k string) interface{} {
	schema, ok := d.Schema[k]
	if !ok {
		panic(fmt.Sprintf("field %s not in the schema", k))
	}

	value, ok := d.GetOk(k)
	if !ok {
		value = schema.DefaultOrZero()
	}

	return value
}

// GetDefaultOrZero gets the default value set on the schema for the given
// field. If there is no default value set, the zero value of the type
// will be returned.
func (d *FieldData) GetDefaultOrZero(k string) interface{} {
	schema, ok := d.Schema[k]
	if !ok {
		panic(fmt.Sprintf("field %s not in the schema", k))
	}

	return schema.DefaultOrZero()
}

// GetFirst gets the value for the given field names, in order from first
// to last. This can be useful for fields with a current name, and one or
// more deprecated names. The second return value will be false if the keys
// are invalid or the keys are not set at all.
func (d *FieldData) GetFirst(k ...string) (interface{}, bool) {
	for _, v := range k {
		if result, ok := d.GetOk(v); ok {
			return result, ok
		}
	}
	return nil, false
}

// GetOk gets the value for the given field. The second return value
// will be false if the key is invalid or the key is not set at all.
func (d *FieldData) GetOk(k string) (interface{}, bool) {
	schema, ok := d.Schema[k]
	if !ok {
		return nil, false
	}

	result, ok, err := d.GetOkErr(k)
	if err != nil {
		panic(fmt.Sprintf("error reading %s: %s", k, err))
	}

	if ok && result == nil {
		result = schema.DefaultOrZero()
	}

	return result, ok
}

// GetOkErr is the most conservative of all the Get methods. It returns
// whether key is set or not, but also an error value. The error value is
// non-nil if the field doesn't exist or there was an error parsing the
// field value.
func (d *FieldData) GetOkErr(k string) (interface{}, bool, error) {
	schema, ok := d.Schema[k]
	if !ok {
		return nil, false, fmt.Errorf("unknown field: %q", k)
	}

	switch schema.Type {
	case TypeBool, TypeInt, TypeMap, TypeDurationSecond, TypeString, TypeLowerCaseString,
		TypeNameString, TypeSlice, TypeStringSlice, TypeCommaStringSlice,
		TypeKVPairs, TypeCommaIntSlice:
		return d.getPrimitive(k, schema)
	default:
		return nil, false,
			fmt.Errorf("unknown field type %q for field %q", schema.Type, k)
	}
}

func (d *FieldData) getPrimitive(k string, schema *FieldSchema) (interface{}, bool, error) {
	raw, ok := d.Raw[k]
	if !ok {
		return nil, false, nil
	}

	switch schema.Type {
	case TypeBool:
		var result bool
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return result, true, nil

	case TypeInt:
		var result int
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return result, true, nil

	case TypeString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return result, true, nil

	case TypeLowerCaseString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return strings.ToLower(result), true, nil

	case TypeNameString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		matched, err := regexp.MatchString("^\\w(([\\w-.]+)?\\w)?$", result)
		if err != nil {
			return nil, true, err
		}
		if !matched {
			return nil, true, errors.New("field does not match the formatting rules")
		}
		return result, true, nil

	case TypeMap:
		var result map[string]interface{}
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return result, true, nil

	case TypeDurationSecond:
		var result int
		switch inp := raw.(type) {
		case nil:
			return nil, false, nil
		case int:
			result = inp
		case int32:
			result = int(inp)
		case int64:
			result = int(inp)
		case uint:
			result = int(inp)
		case uint32:
			result = int(inp)
		case uint64:
			result = int(inp)
		case float32:
			result = int(inp)
		case float64:
			result = int(inp)
		case string:
			dur, err := parseutil.ParseDurationSecond(inp)
			if err != nil {
				return nil, true, err
			}
			result = int(dur.Seconds())
		case json.Number:
			valInt64, err := inp.Int64()
			if err != nil {
				return nil, true, err
			}
			result = int(valInt64)
		default:
			return nil, false, fmt.Errorf("invalid input '%v'", raw)
		}
		return result, true, nil

	case TypeCommaIntSlice:
		var result []int
		config := &mapstructure.DecoderConfig{
			Result:           &result,
			WeaklyTypedInput: true,
			DecodeHook:       mapstructure.StringToSliceHookFunc(","),
		}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return nil, true, err
		}
		if err := decoder.Decode(raw); err != nil {
			return nil, true, err
		}
		return result, true, nil

	case TypeSlice:
		var result []interface{}
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return result, true, nil

	case TypeStringSlice:
		var result []string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, true, err
		}
		return strutil.TrimStrings(result), true, nil

	case TypeCommaStringSlice:
		res, err := parseutil.ParseCommaStringSlice(raw)
		if err != nil {
			return nil, false, err
		}
		return res, true, nil

	case TypeKVPairs:
		// First try to parse this as a map
		var mapResult map[string]string
		if err := mapstructure.WeakDecode(raw, &mapResult); err == nil {
			return mapResult, true, nil
		}

		// If map parse fails, parse as a string list of = delimited pairs
		var listResult []string
		if err := mapstructure.WeakDecode(raw, &listResult); err != nil {
			return nil, true, err
		}

		result := make(map[string]string, len(listResult))
		for _, keyPair := range listResult {
			keyPairSlice := strings.SplitN(keyPair, "=", 2)
			if len(keyPairSlice) != 2 || keyPairSlice[0] == "" {
				return nil, false, fmt.Errorf("invalid key pair %q", keyPair)
			}
			result[keyPairSlice[0]] = keyPairSlice[1]
		}
		return result, true, nil

	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}
}
