/*
Copyright (c) 2024-2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fault

import (
	"reflect"

	"github.com/vmware/govmomi/vim25/types"
)

// As finds the first fault in the error's tree that matches target, and if one
// is found, sets the target to that fault value and returns the fault's
// localized message and true. Otherwise, false is returned.
//
// The tree is inspected according to the object type. If the object implements
// Golang's error interface, the Unwrap() error or Unwrap() []error methods are
// repeatedly checked for additional errors. If the object implements GoVmomi's
// BaseMethodFault or HasLocalizedMethodFault interfaces, the object is checked
// for an underlying FaultCause. When err wraps multiple errors or faults, err
// is examined followed by a depth-first traversal of its children.
//
// An error matches target if the error's concrete value is assignable to the
// value pointed to by target, or if the error has a method
// AsFault(BaseMethodFault) (string, bool) such that AsFault(BaseMethodFault)
// returns true. In the latter case, the AsFault method is responsible for
// setting target.
//
// An error type might provide an AsFault method so it can be treated as if it
// were a different error type.
//
// This function panics if err does not implement error, types.BaseMethodFault,
// types.HasLocalizedMethodFault, Fault() types.BaseMethodFault, or if target is
// not a pointer.
func As(err, target any) (localizedMessage string, okay bool) {
	if err == nil {
		return
	}
	if target == nil {
		panic("fault: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("fault: target must be a non-nil pointer")
	}
	targetType := typ.Elem()
	if targetType.Kind() != reflect.Interface &&
		!targetType.Implements(baseMethodFaultType) {
		panic("fault: *target must be interface or implement BaseMethodFault")
	}
	if !as(err, target, val, targetType, &localizedMessage) {
		return "", false
	}
	return localizedMessage, true
}

func as(
	err,
	target any,
	targetVal reflect.Value,
	targetType reflect.Type,
	localizedMsg *string) bool {

	for {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			targetVal.Elem().Set(reflect.ValueOf(err))
			return true
		}
		if tErr, ok := err.(hasAsFault); ok {
			if msg, ok := tErr.AsFault(target); ok {
				*localizedMsg = msg
				return true
			}
			return false
		}
		switch tErr := err.(type) {
		case types.HasLocalizedMethodFault:
			if fault := tErr.GetLocalizedMethodFault(); fault != nil {
				*localizedMsg = fault.LocalizedMessage
				if fault.Fault != nil {
					return as(
						fault.Fault,
						target,
						targetVal,
						targetType,
						localizedMsg)
				}
			}
			return false
		case types.BaseMethodFault:
			if fault := tErr.GetMethodFault(); fault != nil {
				if fault.FaultCause != nil {
					*localizedMsg = fault.FaultCause.LocalizedMessage
					return as(
						fault.FaultCause,
						target,
						targetVal,
						targetType,
						localizedMsg)
				}
			}
			return false
		case hasFault:
			if fault := tErr.Fault(); fault != nil {
				return as(fault, target, targetVal, targetType, localizedMsg)
			}
			return false
		case unwrappableError:
			if err = tErr.Unwrap(); err == nil {
				return false
			}
		case unwrappableErrorSlice:
			for _, err := range tErr.Unwrap() {
				if err == nil {
					continue
				}
				return as(err, target, targetVal, targetType, localizedMsg)
			}
			return false
		default:
			return false
		}
	}
}

// Is reports whether any fault in err's tree matches target.
//
// The tree is inspected according to the object type. If the object implements
// Golang's error interface, the Unwrap() error or Unwrap() []error methods are
// repeatedly checked for additional errors. If the object implements GoVmomi's
// BaseMethodFault or HasLocalizedMethodFault interfaces, the object is checked
// for an underlying FaultCause. When err wraps multiple errors or faults, err
// is examined followed by a depth-first traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method IsFault(BaseMethodFault) bool such that
// IsFault(BaseMethodFault) returns true.
//
// An error type might provide an IsFault method so it can be treated as
// equivalent to an existing fault. For example, if MyFault defines:
//
//	func (m MyFault) IsFault(target BaseMethodFault) bool {
//		return target == &types.NotSupported{}
//	}
//
// then IsFault(MyError{}, &types.NotSupported{}) returns true. An IsFault
// method should only shallowly compare err and the target and not unwrap
// either.
func Is(err any, target types.BaseMethodFault) bool {
	if target == nil {
		return err == target
	}
	isComparable := reflect.TypeOf(target).Comparable()
	return is(err, target, isComparable)
}

func is(err any, target types.BaseMethodFault, targetComparable bool) bool {
	for {
		if targetComparable && err == target {
			return true
		}
		if tErr, ok := err.(hasIsFault); ok && tErr.IsFault(target) {
			return true
		}
		switch tErr := err.(type) {
		case types.HasLocalizedMethodFault:
			fault := tErr.GetLocalizedMethodFault()
			if fault == nil {
				return false
			}
			err = fault.Fault
		case types.BaseMethodFault:
			if reflect.ValueOf(err).Type() == reflect.ValueOf(target).Type() {
				return true
			}
			fault := tErr.GetMethodFault()
			if fault == nil {
				return false
			}
			err = fault.FaultCause
		case hasFault:
			if err = tErr.Fault(); err == nil {
				return false
			}
		case unwrappableError:
			if err = tErr.Unwrap(); err == nil {
				return false
			}
		case unwrappableErrorSlice:
			for _, err := range tErr.Unwrap() {
				if is(err, target, targetComparable) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

// OnFaultFn is called for every fault encountered when inspecting an error
// or fault for a fault tree. The In function returns when the entire tree is
// inspected or the OnFaultFn returns true.
type OnFaultFn func(
	fault types.BaseMethodFault,
	localizedMessage string,
	localizableMessages []types.LocalizableMessage) bool

// In invokes onFaultFn for each fault in err's tree.
//
// The tree is inspected according to the object type. If the object implements
// Golang's error interface, the Unwrap() error or Unwrap() []error methods are
// repeatedly checked for additional errors. If the object implements GoVmomi's
// BaseMethodFault or HasLocalizedMethodFault interfaces, the object is checked
// for an underlying FaultCause. When err wraps multiple errors or faults, err
// is examined followed by a depth-first traversal of its children.
//
// This function panics if err does not implement error, types.BaseMethodFault,
// types.HasLocalizedMethodFault, Fault() types.BaseMethodFault, or if onFaultFn
// is nil.
func In(err any, onFaultFn OnFaultFn) {
	if onFaultFn == nil {
		panic("fault: onFaultFn must not be nil")
	}
	switch tErr := err.(type) {
	case types.HasLocalizedMethodFault:
		inFault(tErr.GetLocalizedMethodFault(), onFaultFn)
	case types.BaseMethodFault:
		inFault(&types.LocalizedMethodFault{Fault: tErr}, onFaultFn)
	case hasFault:
		if fault := tErr.Fault(); fault != nil {
			inFault(&types.LocalizedMethodFault{Fault: fault}, onFaultFn)
		}
	case unwrappableError:
		In(tErr.Unwrap(), onFaultFn)
	case unwrappableErrorSlice:
		for _, uErr := range tErr.Unwrap() {
			if uErr == nil {
				continue
			}
			In(uErr, onFaultFn)
		}
	case error:
		// No-op
	default:
		panic("fault: err must implement error, types.BaseMethodFault, or " +
			"types.HasLocalizedMethodFault")
	}
}

func inFault(
	localizedMethodFault *types.LocalizedMethodFault,
	onFaultFn OnFaultFn) {

	if localizedMethodFault == nil {
		return
	}

	fault := localizedMethodFault.Fault
	if fault == nil {
		return
	}

	var (
		faultCause    *types.LocalizedMethodFault
		faultMessages []types.LocalizableMessage
	)

	if methodFault := fault.GetMethodFault(); methodFault != nil {
		faultCause = methodFault.FaultCause
		faultMessages = methodFault.FaultMessage
	}

	if onFaultFn(fault, localizedMethodFault.LocalizedMessage, faultMessages) {
		return
	}

	// Check the fault's children.
	inFault(faultCause, onFaultFn)
}

type hasFault interface {
	Fault() types.BaseMethodFault
}

type hasAsFault interface {
	AsFault(target any) (string, bool)
}

type hasIsFault interface {
	IsFault(target types.BaseMethodFault) bool
}

type unwrappableError interface {
	Unwrap() error
}

type unwrappableErrorSlice interface {
	Unwrap() []error
}

var baseMethodFaultType = reflect.TypeOf((*types.BaseMethodFault)(nil)).Elem()
