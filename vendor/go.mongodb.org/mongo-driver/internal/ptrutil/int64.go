// Copyright (C) MongoDB, Inc. 2023-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package ptrutil

// CompareInt64 is a piecewise function with the following return conditions:
//
// (1)  2, ptr1 != nil AND ptr2 == nil
// (2)  1, *ptr1 > *ptr2
// (3)  0, ptr1 == ptr2 or *ptr1 == *ptr2
// (4) -1, *ptr1 < *ptr2
// (5) -2, ptr1 == nil AND ptr2 != nil
func CompareInt64(ptr1, ptr2 *int64) int {
	if ptr1 == ptr2 {
		// This will catch the double nil or same-pointer cases.
		return 0
	}

	if ptr1 == nil && ptr2 != nil {
		return -2
	}

	if ptr1 != nil && ptr2 == nil {
		return 2
	}

	if *ptr1 > *ptr2 {
		return 1
	}

	if *ptr1 < *ptr2 {
		return -1
	}

	return 0
}
