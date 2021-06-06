package parseutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	sockaddr "github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/mitchellh/mapstructure"
)

var validCapacityString = regexp.MustCompile("^[\t ]*([0-9]+)[\t ]?([kmgtKMGT][iI]?[bB])?[\t ]*$")

// ParseCapacityString parses a capacity string and returns the number of bytes it represents.
// Capacity strings are things like 5gib or 10MB. Supported prefixes are kb, kib, mb, mib, gb,
// gib, tb, tib, which are not case sensitive. If no prefix is present, the number is assumed
// to be in bytes already.
func ParseCapacityString(in interface{}) (uint64, error) {
	var cap uint64

	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}

	switch inp := in.(type) {
	case nil:
		// return default of zero
	case string:
		if inp == "" {
			return cap, nil
		}

		matches := validCapacityString.FindStringSubmatch(inp)

		// no sub-groups means we couldn't parse it
		if len(matches) <= 1 {
			return cap, errors.New("could not parse capacity from input")
		}

		var multiplier uint64 = 1
		switch strings.ToLower(matches[2]) {
		case "kb":
			multiplier = 1000
		case "kib":
			multiplier = 1024
		case "mb":
			multiplier = 1000 * 1000
		case "mib":
			multiplier = 1024 * 1024
		case "gb":
			multiplier = 1000 * 1000 * 1000
		case "gib":
			multiplier = 1024 * 1024 * 1024
		case "tb":
			multiplier = 1000 * 1000 * 1000 * 1000
		case "tib":
			multiplier = 1024 * 1024 * 1024 * 1024
		}

		size, err := strconv.ParseUint(matches[1], 10, 64)
		if err != nil {
			return cap, err
		}

		cap = size * multiplier
	case int:
		cap = uint64(inp)
	case int32:
		cap = uint64(inp)
	case int64:
		cap = uint64(inp)
	case uint:
		cap = uint64(inp)
	case uint32:
		cap = uint64(inp)
	case uint64:
		cap = uint64(inp)
	case float32:
		cap = uint64(inp)
	case float64:
		cap = uint64(inp)
	default:
		return cap, errors.New("could not parse capacity from input")
	}

	return cap, nil
}

func ParseDurationSecond(in interface{}) (time.Duration, error) {
	var dur time.Duration
	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}
	switch inp := in.(type) {
	case nil:
		// return default of zero
	case string:
		if inp == "" {
			return dur, nil
		}
		var err error
		// Look for a suffix otherwise its a plain second value
		if strings.HasSuffix(inp, "s") || strings.HasSuffix(inp, "m") || strings.HasSuffix(inp, "h") || strings.HasSuffix(inp, "ms") {
			dur, err = time.ParseDuration(inp)
			if err != nil {
				return dur, err
			}
		} else {
			// Plain integer
			secs, err := strconv.ParseInt(inp, 10, 64)
			if err != nil {
				return dur, err
			}
			dur = time.Duration(secs) * time.Second
		}
	case int:
		dur = time.Duration(inp) * time.Second
	case int32:
		dur = time.Duration(inp) * time.Second
	case int64:
		dur = time.Duration(inp) * time.Second
	case uint:
		dur = time.Duration(inp) * time.Second
	case uint32:
		dur = time.Duration(inp) * time.Second
	case uint64:
		dur = time.Duration(inp) * time.Second
	case float32:
		dur = time.Duration(inp) * time.Second
	case float64:
		dur = time.Duration(inp) * time.Second
	case time.Duration:
		dur = inp
	default:
		return 0, errors.New("could not parse duration from input")
	}

	return dur, nil
}

func ParseAbsoluteTime(in interface{}) (time.Time, error) {
	var t time.Time
	switch inp := in.(type) {
	case nil:
		// return default of zero
		return t, nil
	case string:
		// Allow RFC3339 with nanoseconds, or without,
		// or an epoch time as an integer.
		var err error
		t, err = time.Parse(time.RFC3339Nano, inp)
		if err == nil {
			break
		}
		t, err = time.Parse(time.RFC3339, inp)
		if err == nil {
			break
		}
		epochTime, err := strconv.ParseInt(inp, 10, 64)
		if err == nil {
			t = time.Unix(epochTime, 0)
			break
		}
		return t, errors.New("could not parse string as date and time")
	case json.Number:
		epochTime, err := inp.Int64()
		if err != nil {
			return t, err
		}
		t = time.Unix(epochTime, 0)
	case int:
		t = time.Unix(int64(inp), 0)
	case int32:
		t = time.Unix(int64(inp), 0)
	case int64:
		t = time.Unix(inp, 0)
	case uint:
		t = time.Unix(int64(inp), 0)
	case uint32:
		t = time.Unix(int64(inp), 0)
	case uint64:
		t = time.Unix(int64(inp), 0)
	default:
		return t, errors.New("could not parse time from input type")
	}
	return t, nil
}

func ParseInt(in interface{}) (int64, error) {
	var ret int64
	jsonIn, ok := in.(json.Number)
	if ok {
		in = jsonIn.String()
	}
	switch in.(type) {
	case string:
		inp := in.(string)
		if inp == "" {
			return 0, nil
		}
		var err error
		left, err := strconv.ParseInt(inp, 10, 64)
		if err != nil {
			return ret, err
		}
		ret = left
	case int:
		ret = int64(in.(int))
	case int32:
		ret = int64(in.(int32))
	case int64:
		ret = in.(int64)
	case uint:
		ret = int64(in.(uint))
	case uint32:
		ret = int64(in.(uint32))
	case uint64:
		ret = int64(in.(uint64))
	default:
		return 0, errors.New("could not parse value from input")
	}

	return ret, nil
}

func ParseBool(in interface{}) (bool, error) {
	var result bool
	if err := mapstructure.WeakDecode(in, &result); err != nil {
		return false, err
	}
	return result, nil
}

func ParseString(in interface{}) (string, error) {
	var result string
	if err := mapstructure.WeakDecode(in, &result); err != nil {
		return "", err
	}
	return result, nil
}

func ParseCommaStringSlice(in interface{}) ([]string, error) {
	rawString, ok := in.(string)
	if ok && rawString == "" {
		return []string{}, nil
	}
	var result []string
	config := &mapstructure.DecoderConfig{
		Result:           &result,
		WeaklyTypedInput: true,
		DecodeHook:       mapstructure.StringToSliceHookFunc(","),
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(in); err != nil {
		return nil, err
	}
	return strutil.TrimStrings(result), nil
}

func ParseAddrs(addrs interface{}) ([]*sockaddr.SockAddrMarshaler, error) {
	out := make([]*sockaddr.SockAddrMarshaler, 0)
	stringAddrs := make([]string, 0)

	switch addrs.(type) {
	case string:
		stringAddrs = strutil.ParseArbitraryStringSlice(addrs.(string), ",")
		if len(stringAddrs) == 0 {
			return nil, fmt.Errorf("unable to parse addresses from %v", addrs)
		}

	case []string:
		stringAddrs = addrs.([]string)

	case []interface{}:
		for _, v := range addrs.([]interface{}) {
			stringAddr, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("error parsing %v as string", v)
			}
			stringAddrs = append(stringAddrs, stringAddr)
		}

	default:
		return nil, fmt.Errorf("unknown address input type %T", addrs)
	}

	for _, addr := range stringAddrs {
		sa, err := sockaddr.NewSockAddr(addr)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("error parsing address %q: {{err}}", addr), err)
		}
		out = append(out, &sockaddr.SockAddrMarshaler{
			SockAddr: sa,
		})
	}

	return out, nil
}
