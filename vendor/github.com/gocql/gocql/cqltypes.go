// Copyright (c) 2012 The gocql Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gocql

type Duration struct {
	Months      int32
	Days        int32
	Nanoseconds int64
}
