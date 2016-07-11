package duration

import (
	"strconv"
	"strings"
	"time"
)

func ParseDurationSecond(inp string) (int, error) {
	var result int
	// Look for a suffix otherwise its a plain second value
	if strings.HasSuffix(inp, "s") || strings.HasSuffix(inp, "m") || strings.HasSuffix(inp, "h") {
		dur, err := time.ParseDuration(inp)
		if err != nil {
			return result, err
		}
		result = int(dur.Seconds())
	} else {
		// Plain integer
		val, err := strconv.ParseInt(inp, 10, 64)
		if err != nil {
			return result, err
		}
		result = int(val)
	}

	return result, nil
}
