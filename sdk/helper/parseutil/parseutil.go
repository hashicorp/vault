// DEPRECATED: this has been moved to go-secure-stdlib and will be removed
package parseutil

import (
	"time"

	extparseutil "github.com/hashicorp/go-secure-stdlib/parseutil"
	sockaddr "github.com/hashicorp/go-sockaddr"
)

func ParseCapacityString(in any) (uint64, error) {
	return extparseutil.ParseCapacityString(in)
}

func ParseDurationSecond(in any) (time.Duration, error) {
	return extparseutil.ParseDurationSecond(in)
}

func ParseAbsoluteTime(in any) (time.Time, error) {
	return extparseutil.ParseAbsoluteTime(in)
}

func ParseInt(in any) (int64, error) {
	return extparseutil.ParseInt(in)
}

func ParseBool(in any) (bool, error) {
	return extparseutil.ParseBool(in)
}

func ParseString(in any) (string, error) {
	return extparseutil.ParseString(in)
}

func ParseCommaStringSlice(in any) ([]string, error) {
	return extparseutil.ParseCommaStringSlice(in)
}

func ParseAddrs(addrs any) ([]*sockaddr.SockAddrMarshaler, error) {
	return extparseutil.ParseAddrs(addrs)
}
