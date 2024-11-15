// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signals

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"
)

// SIGNULL is the nil/zero signal.
var SIGNULL os.Signal = syscall.Signal(0)

// ValidSignals is the list of all valid signals. This is built at runtime
// because it is OS-dependent.
var ValidSignals []string
var MonitoredSignals []os.Signal

func init() {
	valid := make([]string, 0, len(SignalLookup))
	monitored := make([]os.Signal, 0, len(SignalLookup))
	for k, v := range SignalLookup {
		valid = append(valid, k)
		monitored = append(monitored, v)
	}
	sort.Strings(valid)
	ValidSignals = valid
	MonitoredSignals = monitored
}

// Parse parses the given string as a signal. If the signal is not found,
// an error is returned.
func Parse(s string) (os.Signal, error) {
	sig, ok := SignalLookup[strings.ToUpper(s)]
	if !ok {
		return nil, fmt.Errorf("invalid signal %q - valid signals are %q",
			s, ValidSignals)
	}
	return sig, nil
}
