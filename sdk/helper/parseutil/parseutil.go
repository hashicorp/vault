// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package parseutil

import (
	"time"

	extparseutil "github.com/hashicorp/go-secure-stdlib/parseutil"
	sockaddr "github.com/hashicorp/go-sockaddr"
)

func ParseCapacityString(in interface{}) (uint64, error) {
	return extparseutil.ParseCapacityString(in)
}

func ParseDurationSecond(in interface{}) (time.Duration, error) {
	return extparseutil.ParseDurationSecond(in)
}

func ParseAbsoluteTime(in interface{}) (time.Time, error) {
	return extparseutil.ParseAbsoluteTime(in)
}

func ParseInt(in interface{}) (int64, error) {
	return extparseutil.ParseInt(in)
}

func ParseBool(in interface{}) (bool, error) {
	return extparseutil.ParseBool(in)
}

func ParseString(in interface{}) (string, error) {
	return extparseutil.ParseString(in)
}

func ParseCommaStringSlice(in interface{}) ([]string, error) {
	return extparseutil.ParseCommaStringSlice(in)
}

func ParseAddrs(addrs interface{}) ([]*sockaddr.SockAddrMarshaler, error) {
	return extparseutil.ParseAddrs(addrs)
}
