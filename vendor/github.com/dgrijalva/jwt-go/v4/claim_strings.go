package jwt

import (
	"encoding/json"
	"reflect"
)

// ClaimStrings is used for parsing claim properties that
// can be either a string or array of strings
type ClaimStrings []string

// ParseClaimStrings is used to produce a ClaimStrings value
// from the various forms it may present during encoding/decodeing
func ParseClaimStrings(value interface{}) (ClaimStrings, error) {
	switch v := value.(type) {
	case string:
		return ClaimStrings{v}, nil
	case []string:
		return ClaimStrings(v), nil
	case []interface{}:
		result := make(ClaimStrings, 0, len(v))
		for i, vv := range v {
			if x, ok := vv.(string); ok {
				result = append(result, x)
			} else {
				return nil, &json.UnsupportedTypeError{Type: reflect.TypeOf(v[i])}
			}
		}
		return result, nil
	case nil:
		return nil, nil
	default:
		return nil, &json.UnsupportedTypeError{Type: reflect.TypeOf(v)}
	}
}

// UnmarshalJSON implements the json package's Unmarshaler interface
func (c *ClaimStrings) UnmarshalJSON(data []byte) error {
	var value interface{}
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}

	*c, err = ParseClaimStrings(value)
	return err
}
