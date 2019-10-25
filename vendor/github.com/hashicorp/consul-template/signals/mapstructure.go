package signals

import (
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// StringToSignalFunc parses a string as a signal based on the signal lookup
// table. If the user supplied an empty string or nil, a special "nil signal"
// is returned. Clients should check for this value and set the response back
// nil after mapstructure finishes parsing.
func StringToSignalFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t.String() != "os.Signal" {
			return data, nil
		}

		if data == nil || data.(string) == "" {
			return SIGNIL, nil
		}

		return Parse(data.(string))
	}
}
