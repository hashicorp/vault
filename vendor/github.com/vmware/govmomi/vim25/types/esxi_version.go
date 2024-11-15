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

package types

import (
	"fmt"
	"regexp"
	"strconv"
)

// ESXiVersion is an ESXi version.
type ESXiVersion uint8

const (
	esxiVersionBegin ESXiVersion = iota

	ESXi2000
	ESXi3000
	ESXi4000
	ESXi5000
	ESXi5100
	ESXi5500
	ESXi6000
	ESXi6500
	ESXi6700
	ESXi6720
	ESXi7000
	ESXi7010
	ESXi7020
	ESXi8000
	ESXi8010
	ESXi8020

	esxiVersionEnd
)

// HardwareVersion returns the maximum hardware version supported by this
// version of ESXi, per https://kb.vmware.com/s/article/1003746.
func (ev ESXiVersion) HardwareVersion() HardwareVersion {
	switch ev {
	case ESXi2000:
		return VMX3
	case ESXi3000:
		return VMX4
	case ESXi4000:
		return VMX7
	case ESXi5000:
		return VMX8
	case ESXi5100:
		return VMX9
	case ESXi5500:
		return VMX10
	case ESXi6000:
		return VMX11
	case ESXi6500:
		return VMX13
	case ESXi6700:
		return VMX14
	case ESXi6720:
		return VMX15
	case ESXi7000:
		return VMX17
	case ESXi7010:
		return VMX18
	case ESXi7020:
		return VMX19
	case ESXi8000, ESXi8010:
		return VMX20
	case ESXi8020:
		return VMX21
	}
	return 0
}

// IsHardwareVersionSupported returns true if the provided hardware version is
// supported by the given version of ESXi.
func (ev ESXiVersion) IsHardwareVersionSupported(hv HardwareVersion) bool {
	return hv <= ev.HardwareVersion()
}

func (ev ESXiVersion) IsValid() bool {
	return ev.String() != ""
}

func (ev ESXiVersion) String() string {
	switch ev {
	case ESXi2000:
		return "2"
	case ESXi3000:
		return "3"
	case ESXi4000:
		return "4"
	case ESXi5000:
		return "5.0"
	case ESXi5100:
		return "5.1"
	case ESXi5500:
		return "5.5"
	case ESXi6000:
		return "6.0"
	case ESXi6500:
		return "6.5"
	case ESXi6700:
		return "6.7"
	case ESXi6720:
		return "6.7.2"
	case ESXi7000:
		return "7.0"
	case ESXi7010:
		return "7.0.1"
	case ESXi7020:
		return "7.0.2"
	case ESXi8000:
		return "8.0"
	case ESXi8010:
		return "8.0.1"
	case ESXi8020:
		return "8.0.2"
	}
	return ""
}

func (ev ESXiVersion) MarshalText() ([]byte, error) {
	return []byte(ev.String()), nil
}

func (ev *ESXiVersion) UnmarshalText(text []byte) error {
	v, err := ParseESXiVersion(string(text))
	if err != nil {
		return err
	}
	*ev = v
	return nil
}

// MustParseESXiVersion parses the provided string into an ESXi version.
func MustParseESXiVersion(s string) ESXiVersion {
	v, err := ParseESXiVersion(s)
	if err != nil {
		panic(err)
	}
	return v
}

var esxiRe = regexp.MustCompile(`(?i)^v?(\d)(?:\.(\d))?(?:\.(\d))?(?:\s*u(\d))?$`)

// ParseESXiVersion parses the provided string into an ESXi version.
func ParseESXiVersion(s string) (ESXiVersion, error) {
	if m := esxiRe.FindStringSubmatch(s); len(m) > 0 {
		var (
			major  int64
			minor  int64
			patch  int64
			update int64
		)

		major, _ = strconv.ParseInt(m[1], 0, 0)
		if len(m) > 2 {
			minor, _ = strconv.ParseInt(m[2], 0, 0)
		}
		if len(m) > 3 {
			patch, _ = strconv.ParseInt(m[3], 0, 0)
		}
		if len(m) > 4 {
			update, _ = strconv.ParseInt(m[4], 0, 0)
		}

		switch {
		case major == 2 && minor == 0 && patch == 0 && update == 0:
			return ESXi2000, nil
		case major == 3 && minor == 0 && patch == 0 && update == 0:
			return ESXi3000, nil
		case major == 4 && minor == 0 && patch == 0 && update == 0:
			return ESXi4000, nil
		case major == 5 && minor == 0 && patch == 0 && update == 0:
			return ESXi5000, nil
		case major == 5 && minor == 1 && patch == 0 && update == 0:
			return ESXi5100, nil
		case major == 5 && minor == 5 && patch == 0 && update == 0:
			return ESXi5500, nil
		case major == 6 && minor == 0 && patch == 0 && update == 0:
			return ESXi6000, nil
		case major == 6 && minor == 5 && patch == 0 && update == 0:
			return ESXi6500, nil
		case major == 6 && minor == 7 && patch == 0 && update == 0:
			return ESXi6700, nil
		case major == 6 && minor == 7 && patch == 2 && update == 0,
			major == 6 && minor == 7 && patch == 0 && update == 2:
			return ESXi6720, nil
		case major == 7 && minor == 0 && patch == 0 && update == 0:
			return ESXi7000, nil
		case major == 7 && minor == 0 && patch == 1 && update == 0,
			major == 7 && minor == 0 && patch == 0 && update == 1:
			return ESXi7010, nil
		case major == 7 && minor == 0 && patch == 2 && update == 0,
			major == 7 && minor == 0 && patch == 0 && update == 2:
			return ESXi7020, nil
		case major == 8 && minor == 0 && patch == 0 && update == 0:
			return ESXi8000, nil
		case major == 8 && minor == 0 && patch == 1 && update == 0,
			major == 8 && minor == 0 && patch == 0 && update == 1:
			return ESXi8010, nil
		case major == 8 && minor == 0 && patch == 2 && update == 0,
			major == 8 && minor == 0 && patch == 0 && update == 2:
			return ESXi8020, nil
		}
	}

	return 0, fmt.Errorf("invalid version: %q", s)
}

// GetESXiVersions returns a list of ESXi versions.
func GetESXiVersions() []ESXiVersion {
	dst := make([]ESXiVersion, esxiVersionEnd-1)
	for i := esxiVersionBegin + 1; i < esxiVersionEnd; i++ {
		dst[i-1] = i
	}
	return dst
}
