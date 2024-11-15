// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package bexpr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-bexpr/grammar"
	"github.com/mitchellh/pointerstructure"
)

var byteSliceTyp reflect.Type = reflect.TypeOf([]byte{})

func primitiveEqualityFn(kind reflect.Kind) func(first interface{}, second reflect.Value) bool {
	switch kind {
	case reflect.Bool:
		return doEqualBool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return doEqualInt64
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return doEqualUint64
	case reflect.Float32:
		return doEqualFloat32
	case reflect.Float64:
		return doEqualFloat64
	case reflect.String:
		return doEqualString
	default:
		return nil
	}
}

func doEqualBool(first interface{}, second reflect.Value) bool {
	return first.(bool) == second.Bool()
}

func doEqualInt64(first interface{}, second reflect.Value) bool {
	return first.(int64) == second.Int()
}

func doEqualUint64(first interface{}, second reflect.Value) bool {
	return first.(uint64) == second.Uint()
}

func doEqualFloat32(first interface{}, second reflect.Value) bool {
	return first.(float32) == float32(second.Float())
}

func doEqualFloat64(first interface{}, second reflect.Value) bool {
	return first.(float64) == second.Float()
}

func doEqualString(first interface{}, second reflect.Value) bool {
	return first.(string) == second.String()
}

// Get rid of 0 to many levels of pointers to get at the real type
func derefType(rtype reflect.Type) reflect.Type {
	for rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}
	return rtype
}

func doMatchMatches(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	if !value.Type().ConvertibleTo(byteSliceTyp) {
		return false, fmt.Errorf("Value of type %s is not convertible to []byte", value.Type())
	}

	var re *regexp.Regexp
	var ok bool
	if expression.Value.Converted != nil {
		re, ok = expression.Value.Converted.(*regexp.Regexp)
	}
	if !ok || re == nil {
		var err error
		re, err = regexp.Compile(expression.Value.Raw)
		if err != nil {
			return false, fmt.Errorf("Failed to compile regular expression %q: %v", expression.Value.Raw, err)
		}
		expression.Value.Converted = re
	}

	return re.Match(value.Convert(byteSliceTyp).Interface().([]byte)), nil
}

func doMatchEqual(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	eqFn := primitiveEqualityFn(value.Kind())
	if eqFn == nil {
		return false, errors.New("unable to find suitable primitive comparison function for matching")
	}
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}
	return eqFn(matchValue, value), nil
}

func doMatchIn(expression *grammar.MatchExpression, value reflect.Value) (bool, error) {
	matchValue, err := getMatchExprValue(expression, value.Kind())
	if err != nil {
		return false, fmt.Errorf("error getting match value in expression: %w", err)
	}

	switch kind := value.Kind(); kind {
	case reflect.Map:
		found := value.MapIndex(reflect.ValueOf(matchValue))
		return found.IsValid(), nil

	case reflect.Slice, reflect.Array:
		itemType := derefType(value.Type().Elem())
		kind := itemType.Kind()
		switch kind {
		case reflect.Interface:
			// If it's an interface, that is, the type was []interface{}, we
			// have to treat each element individually, checking each element's
			// type/kind and rederiving the match value.
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i).Elem()
				itemType := derefType(item.Type())
				kind := itemType.Kind()
				// We need to special case errors here. The reason is that in an
				// interface slice there can be a mix/match of types, but the
				// coerce functions expect a certain type. So the expression
				// passed in might be `"true" in "/my/slice"` but the value it's
				// checking against might be an integer, thus it will try to
				// coerce "true" to an integer and fail. However, all of the
				// functions use strconv which has a specific error type for
				// syntax errors, so as a special case in this situation, don't
				// error on a strconv.ErrSyntax, just continue on to the next
				// element.
				matchValue, err = getMatchExprValue(expression, kind)
				if err != nil {
					if errors.Is(err, strconv.ErrSyntax) {
						continue
					}
					return false, errors.New(`error getting interface slice match value in expression`)
				}
				eqFn := primitiveEqualityFn(kind)
				if eqFn == nil {
					return false, fmt.Errorf(`unable to find suitable primitive comparison function for "in" comparison in interface slice: %s`, kind)
				}
				// the value will be the correct type as we verified the itemType
				if eqFn(matchValue, reflect.Indirect(item)) {
					return true, nil
				}
			}
			return false, nil

		default:
			// Otherwise it's a concrete type and we can essentially cache the
			// answers. First we need to re-derive the match value for equality
			// assertion.
			matchValue, err = getMatchExprValue(expression, kind)
			if err != nil {
				return false, fmt.Errorf("error getting match value in expression: %w", err)
			}
			eqFn := primitiveEqualityFn(kind)
			if eqFn == nil {
				return false, errors.New(`unable to find suitable primitive comparison function for "in" comparison`)
			}
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i)
				// the value will be the correct type as we verified the itemType
				if eqFn(matchValue, reflect.Indirect(item)) {
					return true, nil
				}
			}
			return false, nil
		}

	case reflect.String:
		return strings.Contains(value.String(), matchValue.(string)), nil

	default:
		return false, fmt.Errorf("Cannot perform in/contains operations on type %s for selector: %q", kind, expression.Selector)
	}
}

func doMatchIsEmpty(matcher *grammar.MatchExpression, value reflect.Value) (bool, error) {
	// NOTE: see preconditions in evaluategrammar.MatchExpressionRecurse
	return value.Len() == 0, nil
}

func getMatchExprValue(expression *grammar.MatchExpression, rvalue reflect.Kind) (interface{}, error) {
	if expression.Value == nil {
		return nil, nil
	}

	switch rvalue {
	case reflect.Bool:
		return CoerceBool(expression.Value.Raw)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return CoerceInt64(expression.Value.Raw)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return CoerceUint64(expression.Value.Raw)

	case reflect.Float32:
		return CoerceFloat32(expression.Value.Raw)

	case reflect.Float64:
		return CoerceFloat64(expression.Value.Raw)

	default:
		return expression.Value.Raw, nil
	}
}

// evaluateNotPresent is called after a pointerstructure.ErrNotFound is
// encountered during evaluation.
//
// Returns true if the Selector Path's parent is a map as the missing key may
// be handled by the MatchOperator's NotPresentDisposition method.
//
// Returns false if the Selector Path has a length of 1, or if the parent of
// the Selector's Path is not a map, a pointerstructure.ErrNotFound error is
// returned.
func evaluateNotPresent(ptr pointerstructure.Pointer, datum interface{}) bool {
	if len(ptr.Parts) < 2 {
		return false
	}

	// Pop the missing leaf part of the path
	ptr.Parts = ptr.Parts[0 : len(ptr.Parts)-1]

	val, _ := ptr.Get(datum)
	return reflect.ValueOf(val).Kind() == reflect.Map
}

// getValue resolves path to the value it references by first looking into the
// the local variables, then into the global datum state if it does not.
//
// When the path points to a local variable we have multiple cases we have to
// take care of, in some constructions like
//
//	all Slice as item { item != "forbidden" }
//
// `item` is actually an alias to "/Slice/0", "/Slice/1", etc. In that case we
// compute the full path because we tracked what each of them points to.
//
// In some other cases like
//
//	all Map as key { key != "forbidden" }
//
// `key` has no equivalent JSON Pointer. In that case we kept track of the the
// concrete value instead of the path and we return it directly.
func getValue(datum interface{}, path []string, opt ...Option) (interface{}, bool, error) {
	opts := getOpts(opt...)
	if len(path) != 0 && len(opts.withLocalVariables) > 0 {
		for i := len(opts.withLocalVariables) - 1; i >= 0; i-- {
			name := path[0]
			lv := opts.withLocalVariables[i]
			if name == lv.name {
				if len(lv.path) == 0 {
					// This local variable is a key or an index and we know its
					// value without having to call pointerstructure, we stop
					// here.
					if len(path) > 1 {
						first := pointerstructure.Pointer{Parts: []string{name}}
						full := pointerstructure.Pointer{Parts: path}
						return nil, false, fmt.Errorf("%s references a %T so %s is invalid", first.String(), lv.value, full.String())
					}
					return lv.value, true, nil
				} else {
					// This local variable references another value, we prepend the
					// path of the selector it replaces and continue searching
					prefix := append([]string(nil), lv.path...)
					path = append(prefix, path[1:]...)
				}
			}
		}
	}

	// This is not a local variable, we use pointerstructure to look for it
	// in the global datum
	ptr := pointerstructure.Pointer{
		Parts: path,
		Config: pointerstructure.Config{
			TagName:                 opts.withTagName,
			ValueTransformationHook: opts.withHookFn,
		},
	}
	val, err := ptr.Get(datum)
	if err != nil {
		if errors.Is(err, pointerstructure.ErrNotFound) {
			// Prefer the withUnknown option if set, otherwise defer to NotPresent
			// disposition
			switch {
			case opts.withUnknown != nil:
				err = nil
				val = *opts.withUnknown
			case evaluateNotPresent(ptr, datum):
				return nil, false, nil
			}
		}

		if err != nil {
			return false, false, fmt.Errorf("error finding value in datum: %w", err)
		}
	}

	return val, true, nil
}

func evaluateMatchExpression(expression *grammar.MatchExpression, datum interface{}, opt ...Option) (bool, error) {
	val, present, err := getValue(
		datum,
		expression.Selector.Path,
		opt...,
	)
	if err != nil {
		return false, err
	}
	if !present {
		return expression.Operator.NotPresentDisposition(), nil
	}

	if jn, ok := val.(json.Number); ok {
		if jni, err := jn.Int64(); err == nil {
			val = jni
		} else if jnf, err := jn.Float64(); err == nil {
			val = jnf
		} else {
			return false, fmt.Errorf("unable to convert json number %s to int or float", jn)
		}
	}

	rvalue := reflect.Indirect(reflect.ValueOf(val))
	switch expression.Operator {
	case grammar.MatchEqual:
		return doMatchEqual(expression, rvalue)
	case grammar.MatchNotEqual:
		result, err := doMatchEqual(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case grammar.MatchIn:
		return doMatchIn(expression, rvalue)
	case grammar.MatchNotIn:
		result, err := doMatchIn(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case grammar.MatchIsEmpty:
		return doMatchIsEmpty(expression, rvalue)
	case grammar.MatchIsNotEmpty:
		result, err := doMatchIsEmpty(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	case grammar.MatchMatches:
		return doMatchMatches(expression, rvalue)
	case grammar.MatchNotMatches:
		result, err := doMatchMatches(expression, rvalue)
		if err == nil {
			return !result, nil
		}
		return false, err
	default:
		return false, fmt.Errorf("Invalid match operation: %d", expression.Operator)
	}
}

func evaluateCollectionExpression(expression *grammar.CollectionExpression, datum interface{}, opt ...Option) (bool, error) {
	val, present, err := getValue(
		datum,
		expression.Selector.Path,
		opt...,
	)
	if err != nil {
		return false, err
	}
	if !present {
		return expression.Op == grammar.CollectionOpAll, nil
	}

	v := reflect.ValueOf(val)

	var keys []reflect.Value
	if v.Kind() == reflect.Map {
		if v.Type().Key() != reflect.TypeOf("") {
			return false, fmt.Errorf("%s can only iterate over maps indexed with strings", expression.Op)
		}
		keys = v.MapKeys()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		for i := 0; i < v.Len(); i++ {
			innerOpt := append([]Option(nil), opt...)

			if expression.NameBinding.Mode == grammar.CollectionBindIndexAndValue &&
				expression.NameBinding.Index == expression.NameBinding.Value {
				return false, fmt.Errorf("%q cannot be used as a placeholder for both the index and the value", expression.NameBinding.Index)
			}

			if v.Kind() == reflect.Map {
				key := keys[i]
				if expression.NameBinding.Default != "" {
					innerOpt = append(innerOpt, WithLocalVariable(expression.NameBinding.Default, nil, key.Interface()))
				}
				if expression.NameBinding.Index != "" {
					innerOpt = append(innerOpt, WithLocalVariable(expression.NameBinding.Index, nil, key.Interface()))
				}
				if expression.NameBinding.Value != "" {
					path := make([]string, 0, len(expression.Selector.Path)+1)
					path = append(path, expression.Selector.Path...)
					path = append(path, key.Interface().(string))
					innerOpt = append(innerOpt, WithLocalVariable(expression.NameBinding.Value, path, nil))
				}
			} else {
				if expression.NameBinding.Index != "" {
					innerOpt = append(innerOpt, WithLocalVariable(expression.NameBinding.Index, nil, i))
				}

				pathValue := make([]string, 0, len(expression.Selector.Path)+1)
				pathValue = append(pathValue, expression.Selector.Path...)
				pathValue = append(pathValue, fmt.Sprintf("%d", i))
				if expression.NameBinding.Default != "" {
					innerOpt = append(innerOpt, WithLocalVariable(expression.NameBinding.Default, pathValue, nil))
				}
				if expression.NameBinding.Value != "" {
					innerOpt = append(innerOpt, WithLocalVariable(expression.NameBinding.Value, pathValue, nil))
				}
			}

			result, err := evaluate(expression.Inner, datum, innerOpt...)
			if err != nil {
				return false, err
			}
			if (result && expression.Op == grammar.CollectionOpAny) || (!result && expression.Op == grammar.CollectionOpAll) {
				return result, nil
			}
		}

		return expression.Op == grammar.CollectionOpAll, nil

	default:
		return false, fmt.Errorf(`%s is not a list or a map`, expression.Selector.String())
	}
}

func evaluate(ast grammar.Expression, datum interface{}, opt ...Option) (bool, error) {
	switch node := ast.(type) {
	case *grammar.UnaryExpression:
		switch node.Operator {
		case grammar.UnaryOpNot:
			result, err := evaluate(node.Operand, datum, opt...)
			return !result, err
		}
	case *grammar.BinaryExpression:
		switch node.Operator {
		case grammar.BinaryOpAnd:
			result, err := evaluate(node.Left, datum, opt...)
			if err != nil || !result {
				return result, err
			}

			return evaluate(node.Right, datum, opt...)

		case grammar.BinaryOpOr:
			result, err := evaluate(node.Left, datum, opt...)
			if err != nil || result {
				return result, err
			}

			return evaluate(node.Right, datum, opt...)
		}
	case *grammar.MatchExpression:
		return evaluateMatchExpression(node, datum, opt...)
	case *grammar.CollectionExpression:
		return evaluateCollectionExpression(node, datum, opt...)
	}
	return false, fmt.Errorf("Invalid AST node")
}
