// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"
	"reflect"
	"strconv"
)

// formatterConfig is used to provide basic configuration to a formatter.
// Use newFormatterConfig to initialize the formatterConfig struct.
type formatterConfig struct {
	formatterConfigEnt

	raw          bool
	hmacAccessor bool

	// Vault lacks pagination in its APIs. As a result, certain list operations can return **very** large responses.
	// The user's chosen audit sinks may experience difficulty consuming audit records that swell to tens of megabytes
	// of JSON. The responses of list operations are typically not very interesting, as they are mostly lists of keys,
	// or, even when they include a "key_info" field, are not returning confidential information. They become even less
	// interesting once HMAC-ed by the audit system.
	//
	// Some example Vault "list" operations that are prone to becoming very large in an active Vault installation are:
	//   auth/token/accessors/
	//   identity/entity/id/
	//   identity/entity-alias/id/
	//   pki/certs/
	//
	// This option exists to provide such users with the option to have response data elided from audit logs, only when
	// the operation type is "list". For added safety, the elision only applies to the "keys" and "key_info" fields
	// within the response data - these are conventionally the only fields present in a list response - see
	// logical.ListResponse, and logical.ListResponseWithInfo. However, other fields are technically possible if a
	// plugin author writes unusual code, and these will be preserved in the audit log even with this option enabled.
	// The elision replaces the values of the "keys" and "key_info" fields with an integer count of the number of
	// entries. This allows even the elided audit logs to still be useful for answering questions like
	// "Was any data returned?" or "How many records were listed?".
	elideListResponses bool

	// This should only ever be used in a testing context
	omitTime bool

	// The required/target format for the event (supported: jsonFormat and jsonxFormat).
	requiredFormat format

	// headerFormatter specifies the formatter used for headers that existing in any incoming audit request.
	headerFormatter HeaderFormatter

	// prefix specifies a prefix that should be prepended to any formatted request or response before serialization.
	prefix string
}

// newFormatterConfig creates the configuration required by a formatter node using the config map supplied to the factory.
func newFormatterConfig(headerFormatter HeaderFormatter, config map[string]string) (formatterConfig, error) {
	if headerFormatter == nil || reflect.ValueOf(headerFormatter).IsNil() {
		return formatterConfig{}, fmt.Errorf("header formatter is required: %w", ErrInvalidParameter)
	}

	var opt []option

	if format, ok := config[optionFormat]; ok {
		if !isValidFormat(format) {
			return formatterConfig{}, fmt.Errorf("unsupported %q: %w", optionFormat, ErrExternalOptions)
		}

		opt = append(opt, withFormat(format))
	}

	// Check if hashing of accessor is disabled
	if hmacAccessorRaw, ok := config[optionHMACAccessor]; ok {
		v, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return formatterConfig{}, fmt.Errorf("unable to parse %q: %w", optionHMACAccessor, ErrExternalOptions)
		}
		opt = append(opt, withHMACAccessor(v))
	}

	// Check if raw logging is enabled
	if raw, ok := config[optionLogRaw]; ok {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return formatterConfig{}, fmt.Errorf("unable to parse %q: %w", optionLogRaw, ErrExternalOptions)
		}
		opt = append(opt, withRaw(v))
	}

	if elideListResponsesRaw, ok := config[optionElideListResponses]; ok {
		v, err := strconv.ParseBool(elideListResponsesRaw)
		if err != nil {
			return formatterConfig{}, fmt.Errorf("unable to parse %q: %w", optionElideListResponses, ErrExternalOptions)
		}
		opt = append(opt, withElision(v))
	}

	if prefix, ok := config[optionPrefix]; ok {
		opt = append(opt, withPrefix(prefix))
	}

	opts, err := getOpts(opt...)
	if err != nil {
		return formatterConfig{}, err
	}

	fmtCfgEnt, err := newFormatterConfigEnt(config)
	if err != nil {
		return formatterConfig{}, err
	}

	return formatterConfig{
		formatterConfigEnt: fmtCfgEnt,
		headerFormatter:    headerFormatter,
		elideListResponses: opts.withElision,
		hmacAccessor:       opts.withHMACAccessor,
		omitTime:           opts.withOmitTime, // This must be set in code after creation.
		prefix:             opts.withPrefix,
		raw:                opts.withRaw,
		requiredFormat:     opts.withFormat,
	}, nil
}
