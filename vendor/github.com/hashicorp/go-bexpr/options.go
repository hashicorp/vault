// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bexpr

// getOpts - iterate the inbound Options and return a struct
func getOpts(opt ...Option) options {
	opts := getDefaultOptions()
	for _, o := range opt {
		if o != nil {
			o(&opts)
		}
	}
	return opts
}

// a localVariable can either point to a known value or replace another JSON
// Pointer path
type localVariable struct {
	name  string
	path  []string
	value any
}

// Option - how Options are passed as arguments
type Option func(*options)

// options = how options are represented
type options struct {
	withMaxExpressions uint64
	withTagName        string
	withHookFn         ValueTransformationHookFn
	withUnknown        *interface{}
	withLocalVariables []localVariable
}

func WithMaxExpressions(maxExprCnt uint64) Option {
	return func(o *options) {
		o.withMaxExpressions = maxExprCnt
	}
}

// WithTagName indictes what tag to use instead of the default "bexpr"
func WithTagName(tagName string) Option {
	return func(o *options) {
		o.withTagName = tagName
	}
}

// WithHookFn sets a HookFn to be called on the Go data under evaluation
// and all subfields, indexes, and values recursively.  That makes it
// easier for the JSON Pointer to not match exactly the Go value being
// evaluated (for example, when using protocol buffers' well-known types).
func WithHookFn(fn ValueTransformationHookFn) Option {
	return func(o *options) {
		o.withHookFn = fn
	}
}

// WithUnknownValue sets a value that is used for any unknown keys. Normally,
// bexpr will error on any expressions with unknown keys. This can be set to
// instead use a specificed value whenever an unknown key is found. For example,
// this might be set to the empty string "".
func WithUnknownValue(val interface{}) Option {
	return func(o *options) {
		o.withUnknown = &val
	}
}

// WithLocalVariable add a local variable that can either point to another path
// that will be resolved when the local variable is referenced or to a known
// value that will be used directly.
func WithLocalVariable(name string, path []string, value any) Option {
	return func(o *options) {
		o.withLocalVariables = append(o.withLocalVariables, localVariable{
			name:  name,
			path:  path,
			value: value,
		})
	}
}

func getDefaultOptions() options {
	return options{
		withMaxExpressions: 0,
		withTagName:        "bexpr",
		withUnknown:        nil,
	}
}
