// Copyright 2021 Couchbase
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocbcore

import (
	"errors"
	"fmt"
	"strconv"
)

// From Java impl:
// ${Mutation.CAS} is written by kvengine with 'macroToString(htonll(info.cas))'.  Discussed this with KV team and,
// though there is consensus that this is off (htonll is definitely wrong, and a string is an odd choice), there are
// clients (SyncGateway) that consume the current string, so it can't be changed.  Note that only little-endian
// servers are supported for Couchbase, so the 8 byte long inside the string will always be little-endian ordered.
//
// Looks like: "0x000058a71dd25c15"
// Want:        0x155CD21DA7580000   (1539336197457313792 in base10, an epoch time in millionths of a second)
func parseCASToMilliseconds(in string) (int64, error) {
	if len(in) < 18 {
		logWarnf("Invalid mutation cas value seen in cleanup: %s", in)
		return 0, errors.New("invalid cas value provided")
	}
	offsetIndex := 2 // for the initial "0x"
	result := int64(0)

	for octetIndex := 7; octetIndex >= 0; octetIndex-- {
		char1 := in[offsetIndex+(octetIndex*2)]
		char2 := in[offsetIndex+(octetIndex*2)+1]

		octet1 := int64(0)
		octet2 := int64(0)

		if char1 >= 'a' && char1 <= 'f' {
			octet1 = int64(char1 - 'a' + 10)
		} else if char1 >= 'A' && char1 <= 'F' {
			octet1 = int64(char1 - 'A' + 10)
		} else if char1 >= '0' && char1 <= '9' {
			octet1 = int64(char1 - '0')
		} else {
			return 0, fmt.Errorf("could not parse CAS: %s", in)
		}

		if char2 >= 'a' && char2 <= 'f' {
			octet2 = int64(char2 - 'a' + 10)
		} else if char2 >= 'A' && char2 <= 'F' {
			octet2 = int64(char2 - 'A' + 10)
		} else if char2 >= '0' && char2 <= '9' {
			octet2 = int64(char2 - '0')
		} else {
			return 0, fmt.Errorf("could not parse CAS: %s", in)
		}

		result |= octet1 << ((octetIndex * 8) + 4)
		result |= octet2 << (octetIndex * 8)
	}

	// It's in nanoseconds, let's return milliseconds.
	return result / 1000000, nil
}

func parseHLCToSeconds(hlc jsonHLC) (int64, error) {
	return strconv.ParseInt(hlc.NowSecs, 10, 64)
}
