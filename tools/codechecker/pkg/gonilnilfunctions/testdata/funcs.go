// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package testdata

func ReturnReturnOkay() (any, error) {
	var i interface{}
	return i, nil
}

func OneGoodOneBad() (any, error) { // want "Function OneGoodOneBad can return an error, and has a statement that returns only nils"
	var i interface{}
	if true {
		return i, nil
	}
	return nil, nil
}

func OneBadOneGood() (any, error) { // want "Function OneBadOneGood can return an error, and has a statement that returns only nils"
	var i interface{}
	if true {
		return nil, nil
	}
	return i, nil
}

func EmptyFunc() {}

func TwoNilNils() (any, error) { // want "Function TwoNilNils can return an error, and has a statement that returns only nils"
	if true {
		return nil, nil
	}
	return nil, nil
}

// ThreeResults should not fail, as while it returns nil, nil, nil, it has three results, not two.
func ThreeResults() (any, any, error) {
	return nil, nil, nil
}

func TwoArgsNoError() (any, any) {
	return nil, nil
}

func NestedReturn() (any, error) { // want "Function NestedReturn can return an error, and has a statement that returns only nils"
	{
		{
			{
				return nil, nil
			}
		}
	}
}

func NestedForReturn() (any, error) { // want "Function NestedForReturn can return an error, and has a statement that returns only nils"
	for {
		for i := 0; i < 100; i++ {
			{
				return nil, nil
			}
		}
	}
}

func AnyErrorNilNil() (any, error) { // want "Function AnyErrorNilNil can return an error, and has a statement that returns only nils"
	return nil, nil
}

// Skipped should be skipped because of the following line:
// ignore-nil-nil-function-check
func Skipped() (any, error) {
	return nil, nil
}
