package opts

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ValidateEnv validates an environment variable and returns it.
// If no value is specified, it obtains its value from the current environment
//
// As on ParseEnvFile and related to #16585, environment variable names
// are not validated, and it's up to the application inside the container
// to validate them or not.
//
// The only validation here is to check if name is empty, per #25099
func ValidateEnv(val string) (string, error) {
	k, _, hasValue := strings.Cut(val, "=")
	if k == "" {
		return "", errors.New("invalid environment variable: " + val)
	}
	if hasValue {
		// val contains a "=" (but value may be an empty string)
		return val, nil
	}
	if envVal, ok := os.LookupEnv(k); ok {
		return k + "=" + envVal, nil
	}
	return val, nil
}
