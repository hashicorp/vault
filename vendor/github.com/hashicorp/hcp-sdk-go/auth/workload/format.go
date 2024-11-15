// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workload

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonpointer"
)

const (
	// FormatTypeValue indicates that the value itself contains the access_token
	FormatTypeValue = "value"

	// FormatTypeJSON indicates that the response is a JSON payload that
	// contains the access_token.
	FormatTypeJSON = "json"
)

// CredentialFormat configures how to extract the credential from the source
// value. It supports either treating the entire response as the value or
// extracting a particular field from a JSON response.
type CredentialFormat struct {
	// Type is either "text" or "json". When not provided "text" type is assumed.
	Type string `json:"format_type,omitempty"`

	// SubjectCredentialPointer is a JSON pointer that indicates how to access
	// the subject credential.
	SubjectCredentialPointer string `json:"subject_cred_pointer,omitempty"`
}

// Validate validates the format configuration.
func (cf CredentialFormat) Validate() error {
	credType := cf.Type
	if credType == "" {
		credType = FormatTypeValue
	}

	if credType != FormatTypeValue && credType != FormatTypeJSON {
		return fmt.Errorf("format type must either be %q or %q", FormatTypeValue, FormatTypeJSON)
	}

	if credType == FormatTypeValue && cf.SubjectCredentialPointer != "" {
		return fmt.Errorf("subject credential pointer must not be set with format type %q", FormatTypeValue)
	}

	if credType == FormatTypeJSON {
		if cf.SubjectCredentialPointer == "" {
			return fmt.Errorf("subject credential pointer must be set with format type %q", FormatTypeJSON)
		}

		if _, err := gojsonpointer.NewJsonPointer(cf.SubjectCredentialPointer); err != nil {
			return fmt.Errorf("subject credential pointer is invalid: %v", err)
		}
	}

	return nil
}

// get extracts the subject token from the passed response value based on the
// CredentialFormat configuration.
func (cf CredentialFormat) get(value []byte) (string, error) {
	if cf.Type == FormatTypeValue || cf.Type == "" {
		return string(value), nil
	}

	// Unmarshal the JSON value
	jsonData := make(map[string]interface{})
	if err := json.Unmarshal(value, &jsonData); err != nil {
		return "", fmt.Errorf("failed to unmarshal json value: %v", err)
	}

	jsonp, err := gojsonpointer.NewJsonPointer(cf.SubjectCredentialPointer)
	if err != nil {
		return "", fmt.Errorf("subject credential pointer is invalid: %v", err)
	}

	// Retrieve the access token using the JSON pointer
	val, _, err := jsonp.Get(jsonData)
	if err != nil {
		return "", err
	}

	cred, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("credential must be a string; got %T", val)
	}

	return cred, nil
}
