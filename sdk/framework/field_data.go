package framework

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
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
		case TypeBool, TypeInt, TypeInt64, TypeMap, TypeDurationSecond, TypeSignedDurationSecond, TypeString,
			TypeLowerCaseString, TypeNameString, TypeSlice, TypeStringSlice, TypeCommaStringSlice,
			TypeKVPairs, TypeCommaIntSlice, TypeHeader, TypeFloat, TypeTime:
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

	// If the value can't be decoded, use the zero or default value for the field
	// type
	value, ok := d.GetOk(k)
	if !ok || value == nil {
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

// GetOk gets the value for the given field. The second return value will be
// false if the key is invalid or the key is not set at all. If the field k is
// set and the decoded value is nil, the default or zero value
// will be returned instead.
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
	case TypeBool, TypeInt, TypeInt64, TypeMap, TypeDurationSecond, TypeSignedDurationSecond, TypeString,
		TypeLowerCaseString, TypeNameString, TypeSlice, TypeStringSlice, TypeCommaStringSlice,
		TypeKVPairs, TypeCommaIntSlice, TypeHeader, TypeFloat, TypeTime:
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

	switch t := schema.Type; t {
	case TypeBool:
		var result bool
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return result, true, nil

	case TypeInt:
		var result int
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return result, true, nil

	case TypeInt64:
		var result int64
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return result, true, nil

	case TypeFloat:
		var result float64
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return result, true, nil

	case TypeString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return result, true, nil

	case TypeLowerCaseString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return strings.ToLower(result), true, nil

	case TypeNameString:
		var result string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		matched, err := regexp.MatchString("^\\w(([\\w-.]+)?\\w)?$", result)
		if err != nil {
			return nil, false, err
		}
		if !matched {
			return nil, false, errors.New("field does not match the formatting rules")
		}
		return result, true, nil

	case TypeMap:
		var result map[string]interface{}
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		return result, true, nil

	case TypeDurationSecond, TypeSignedDurationSecond:
		var result int
		switch inp := raw.(type) {
		case nil:
			return nil, false, nil
		default:
			dur, err := parseutil.ParseDurationSecond(inp)
			if err != nil {
				return nil, false, err
			}
			result = int(dur.Seconds())
		}
		if t == TypeDurationSecond && result < 0 {
			return nil, false, fmt.Errorf("cannot provide negative value '%d'", result)
		}
		return result, true, nil

	case TypeTime:
		switch inp := raw.(type) {
		case nil:
			// Handle nil interface{} as a non-error case
			return nil, false, nil
		default:
			time, err := parseutil.ParseAbsoluteTime(inp)
			if err != nil {
				return nil, false, err
			}
			return time.UTC(), true, nil
		}

	case TypeCommaIntSlice:
		var result []int
		config := &mapstructure.DecoderConfig{
			Result:           &result,
			WeaklyTypedInput: true,
			DecodeHook:       mapstructure.StringToSliceHookFunc(","),
		}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return nil, false, err
		}
		if err := decoder.Decode(raw); err != nil {
			return nil, false, err
		}
		if len(result) == 0 {
			return make([]int, 0), true, nil
		}
		return result, true, nil

	case TypeSlice:
		var result []interface{}
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		if len(result) == 0 {
			return make([]interface{}, 0), true, nil
		}
		return result, true, nil

	case TypeStringSlice:
		rawString, ok := raw.(string)
		if ok && rawString == "" {
			return []string{}, true, nil
		}

		var result []string
		if err := mapstructure.WeakDecode(raw, &result); err != nil {
			return nil, false, err
		}
		if len(result) == 0 {
			return make([]string, 0), true, nil
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
			return nil, false, err
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

	case TypeHeader:
		/*

			There are multiple ways a header could be provided:

			1.	As a map[string]interface{} that resolves to a map[string]string or map[string][]string, or a mix of both
				because that's permitted for headers.
				This mainly comes from the API.

			2.	As a string...
				a. That contains JSON that originally was JSON, but then was base64 encoded.
				b. That contains JSON, ex. `{"content-type":"text/json","accept":["encoding/json"]}`.
				This mainly comes from the API and is used to save space while sending in the header.

			3.	As an array of strings that contains comma-delimited key-value pairs associated via a colon,
				ex: `content-type:text/json`,`accept:encoding/json`.
				This mainly comes from the CLI.

			We go through these sequentially below.

		*/
		result := http.Header{}

		toHeader := func(resultMap map[string]interface{}) (http.Header, error) {
			header := http.Header{}
			for headerKey, headerValGroup := range resultMap {
				switch typedHeader := headerValGroup.(type) {
				case string:
					header.Add(headerKey, typedHeader)
				case []string:
					for _, headerVal := range typedHeader {
						header.Add(headerKey, headerVal)
					}
				case json.Number:
					header.Add(headerKey, typedHeader.String())
				case []interface{}:
					for _, headerVal := range typedHeader {
						switch typedHeader := headerVal.(type) {
						case string:
							header.Add(headerKey, typedHeader)
						case json.Number:
							header.Add(headerKey, typedHeader.String())
						default:
							// All header values should already be strings when they're being sent in.
							// Even numbers and booleans will be treated as strings.
							return nil, fmt.Errorf("received non-string value for header key:%s, val:%s", headerKey, headerValGroup)
						}
					}
				default:
					return nil, fmt.Errorf("unrecognized type for %s", headerValGroup)
				}
			}
			return header, nil
		}

		resultMap := make(map[string]interface{})

		// 1. Are we getting a map from the API?
		if err := mapstructure.WeakDecode(raw, &resultMap); err == nil {
			result, err = toHeader(resultMap)
			if err != nil {
				return nil, false, err
			}
			return result, true, nil
		}

		// 2. Are we getting a JSON string?
		if headerStr, ok := raw.(string); ok {
			// a. Is it base64 encoded?
			headerBytes, err := base64.StdEncoding.DecodeString(headerStr)
			if err != nil {
				// b. It's not base64 encoded, it's a straight-out JSON string.
				headerBytes = []byte(headerStr)
			}
			if err := jsonutil.DecodeJSON(headerBytes, &resultMap); err != nil {
				return nil, false, err
			}
			result, err = toHeader(resultMap)
			if err != nil {
				return nil, false, err
			}
			return result, true, nil
		}

		// 3. Are we getting an array of fields like "content-type:encoding/json" from the CLI?
		var keyPairs []interface{}
		if err := mapstructure.WeakDecode(raw, &keyPairs); err == nil {
			for _, keyPairIfc := range keyPairs {
				keyPair, ok := keyPairIfc.(string)
				if !ok {
					return nil, false, fmt.Errorf("invalid key pair %q", keyPair)
				}
				keyPairSlice := strings.SplitN(keyPair, ":", 2)
				if len(keyPairSlice) != 2 || keyPairSlice[0] == "" {
					return nil, false, fmt.Errorf("invalid key pair %q", keyPair)
				}
				result.Add(keyPairSlice[0], keyPairSlice[1])
			}
			return result, true, nil
		}
		return nil, false, fmt.Errorf("%s not provided an expected format", raw)

	default:
		panic(fmt.Sprintf("Unknown type: %s", schema.Type))
	}
}
