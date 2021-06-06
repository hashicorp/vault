package signals

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// SIGNIL is the nil signal.
var SIGNIL os.Signal = new(NilSignal)

// ValidSignals is the list of all valid signals. This is built at runtime
// because it is OS-dependent.
var ValidSignals []string

func init() {
	valid := make([]string, 0, len(SignalLookup))
	for k := range SignalLookup {
		valid = append(valid, k)
	}
	sort.Strings(valid)
	ValidSignals = valid
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
