package driver

import (
	_ "embed" // embed stats configuration
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

//go:embed statscfg.json
var statsCfgRaw []byte

var statsCfg struct {
	TimeUnit        string    `json:"timeUnit"`
	SQLTimeTexts    []string  `json:"sqlTimeTexts"`
	TimeUpperBounds []float64 `json:"timeUpperBounds"`
}

// time unit map (see go package time format.go).
var timeUnitMap = map[string]uint64{
	"ns": uint64(time.Nanosecond),
	"us": uint64(time.Microsecond),
	"µs": uint64(time.Microsecond), // U+00B5 = micro symbol
	"μs": uint64(time.Microsecond), // U+03BC = Greek letter mu
	"ms": uint64(time.Millisecond),
	"s":  uint64(time.Second),
	"m":  uint64(time.Minute),
	"h":  uint64(time.Hour),
}

func loadStatsCfg() error {

	if err := json.Unmarshal(statsCfgRaw, &statsCfg); err != nil {
		return fmt.Errorf("invalid statscfg.json file: %w", err)
	}

	if len(statsCfg.SQLTimeTexts) != int(numSQLTime) {
		return fmt.Errorf("invalid number of statscfg.json sqlTimeTexts %d - expected %d", len(statsCfg.SQLTimeTexts), numSQLTime)
	}
	if len(statsCfg.TimeUpperBounds) == 0 {
		return fmt.Errorf("number of statscfg.json timeUpperBounds needs to be greater than %d", 0)
	}

	if _, ok := timeUnitMap[statsCfg.TimeUnit]; !ok {
		return fmt.Errorf("invalid time unit in statscfg.json %s", statsCfg.TimeUnit)
	}

	// sort and dedup timeBuckets
	slices.Sort(statsCfg.TimeUpperBounds)
	statsCfg.TimeUpperBounds = slices.Compact(statsCfg.TimeUpperBounds)

	return nil
}
