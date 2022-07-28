package custommetadata

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
)

// CustomMetadata should be arbitrary user-provided key-value pairs meant to
// provide supplemental information about a resource.
type CustomMetadata map[string]string

// The following constants are used by Validate and are meant to be imposed
// broadly for consistency.
const (
	maxKeys               = 64
	maxKeyLength          = 128
	maxValueLength        = 512
	validationErrorPrefix = "custom_metadata validation failed"
)

// Validate will perform input validation for custom metadata. If the key count
// exceeds maxKeys, the validation will be short-circuited to prevent
// unnecessary (and potentially costly) validation to be run. If the key count
// falls at or below maxKeys, multiple checks will be made per key and value.
// These checks include:
//   - 0 < length of key <= maxKeyLength
//   - 0 < length of value <= maxValueLength
//   - keys and values cannot include unprintable characters
func Validate(cm CustomMetadata) error {
	var errs *multierror.Error

	if keyCount := len(cm); keyCount > maxKeys {
		errs = multierror.Append(errs, fmt.Errorf("%s: payload must contain at most %d keys, provided %d",
			validationErrorPrefix,
			maxKeys,
			keyCount))

		return errs.ErrorOrNil()
	}

	// Perform validation on each key and value and return ALL errors
	for key, value := range cm {
		if keyLen := len(key); 0 == keyLen || keyLen > maxKeyLength {
			errs = multierror.Append(errs, fmt.Errorf("%s: length of key %q is %d but must be 0 < len(key) <= %d",
				validationErrorPrefix,
				key,
				keyLen,
				maxKeyLength))
		}

		if valueLen := len(value); 0 == valueLen || valueLen > maxValueLength {
			errs = multierror.Append(errs, fmt.Errorf("%s: length of value for key %q is %d but must be 0 < len(value) <= %d",
				validationErrorPrefix,
				key,
				valueLen,
				maxValueLength))
		}

		if !strutil.Printable(key) {
			// Include unquoted format (%s) to also include the string without the unprintable
			//  characters visible to allow for easier debug and key identification
			errs = multierror.Append(errs, fmt.Errorf("%s: key %q (%s) contains unprintable characters",
				validationErrorPrefix,
				key,
				key))
		}

		if !strutil.Printable(value) {
			errs = multierror.Append(errs, fmt.Errorf("%s: value for key %q contains unprintable characters",
				validationErrorPrefix,
				key))
		}
	}

	return errs.ErrorOrNil()
}
