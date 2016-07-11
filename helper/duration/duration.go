package duration

import (
	"strconv"
	"strings"
	"time"
)

func ParseDurationSecond(inp string) (time.Duration, error) {
	var err error
	var dur time.Duration
	// Look for a suffix otherwise its a plain second value
	if strings.HasSuffix(inp, "s") || strings.HasSuffix(inp, "m") || strings.HasSuffix(inp, "h") {
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

	return dur, nil
}
