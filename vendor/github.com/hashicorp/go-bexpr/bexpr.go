// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package bexpr is an implementation of a generic boolean expression evaluator.
// The general goal is to be able to evaluate some expression against some
// arbitrary data and get back a boolean indicating if the data was matched by
// the expression
package bexpr

//go:generate pigeon -o grammar/grammar.go -optimize-parser grammar/grammar.peg
//go:generate goimports -w grammar/grammar.go

import (
	"github.com/hashicorp/go-bexpr/grammar"
	"github.com/mitchellh/pointerstructure"
)

// ValueTransformationHookFn provides a way to translate one reflect.Value to another during
// evaluation by bexpr. This facilitates making Go structures appear in a way
// that matches the expected JSON Pointers used for evaluation. This is
// helpful, for example, when working with protocol buffers' well-known types.
type ValueTransformationHookFn = pointerstructure.ValueTransformationHookFn

type Evaluator struct {
	// The syntax tree
	ast                     grammar.Expression
	tagName                 string
	valueTransformationHook ValueTransformationHookFn
	unknownVal              *interface{}
	expression              string
}

// CreateEvaluator is used to create and configure a new Evaluator, the expression
// will be used by the evaluator when evaluating against any supplied datum.
// The following Option types are supported:
// WithHookFn, WithMaxExpressions, WithTagName, WithUnknownValue.
func CreateEvaluator(expression string, opts ...Option) (*Evaluator, error) {
	parsedOpts := getOpts(opts...)
	var parserOpts []grammar.Option
	if parsedOpts.withMaxExpressions != 0 {
		parserOpts = append(parserOpts, grammar.MaxExpressions(parsedOpts.withMaxExpressions))
	}

	ast, err := grammar.Parse("", []byte(expression), parserOpts...)
	if err != nil {
		return nil, err
	}

	eval := &Evaluator{
		ast:                     ast.(grammar.Expression),
		tagName:                 parsedOpts.withTagName,
		valueTransformationHook: parsedOpts.withHookFn,
		unknownVal:              parsedOpts.withUnknown,
		expression:              expression,
	}

	return eval, nil
}

// Evaluate attempts to match the configured expression against the supplied datum.
// It returns a value indicating if a match was found and any error that occurred.
// If an error is returned, the value indicating a match will be false.
func (eval *Evaluator) Evaluate(datum interface{}) (bool, error) {
	opts := []Option{
		WithTagName(eval.tagName),
		WithHookFn(eval.valueTransformationHook),
	}
	if eval.unknownVal != nil {
		opts = append(opts, WithUnknownValue(*eval.unknownVal))
	}

	return evaluate(eval.ast, datum, opts...)
}

// Expression can be used to return the initial expression used to create the Evaluator.
func (eval *Evaluator) Expression() string {
	return eval.expression
}
