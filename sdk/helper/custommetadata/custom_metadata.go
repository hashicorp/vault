// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package custommetadata

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/mitchellh/mapstructure"
)

// The following constants are used by Validate and are meant to be imposed
// broadly for consistency.
const (
	maxKeys               = 64
	maxKeyLength          = 128
	maxValueLength        = 512
	validationErrorPrefix = "custom_metadata validation failed"
)

// Parse is used to effectively convert the TypeMap
// (map[string]interface{}) into a TypeKVPairs (map[string]string)
// which is how custom_metadata is stored. Defining custom_metadata
// as a TypeKVPairs will convert nulls into empty strings. A null,
// however, is essential for a PATCH operation in that it signals
// the handler to remove the field. The filterNils flag should
// only be used during a patch operation.
func Parse(raw map[string]interface{}, filterNils bool) (map[string]string, error) {
	customMetadata := map[string]string{}
	for k, v := range raw {
		if filterNils && v == nil {
			continue
		}

		var s string
		if err := mapstructure.WeakDecode(v, &s); err != nil {
			return nil, err
		}

		customMetadata[k] = s
	}

	return customMetadata, nil
}

// Validate will perform input validation for custom metadata.
// CustomMetadata should be arbitrary user-provided key-value pairs meant to
// provide supplemental information about a resource. If the key count
// exceeds maxKeys, the validation will be short-circuited to prevent
// unnecessary (and potentially costly) validation to be run. If the key count
// falls at or below maxKeys, multiple checks will be made per key and value.
// These checks include:
//   - 0 < length of key <= maxKeyLength
//   - 0 < length of value <= maxValueLength
//   - keys and values cannot include unprintable characters
func Validate(cm map[string]string) error {
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
