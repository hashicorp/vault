package pki

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_calcRandomStartupDelayer(t *testing.T) {
	t.Parallel()
	type args struct {
		min time.Duration
		max time.Duration
	}
	staticNow := time.Now()
	nowFunc := func() time.Time { return staticNow }
	tests := []struct {
		name      string
		args      args
		exactTime bool
	}{
		{name: "disable", args: args{0, 0}, exactTime: true},
		{name: "lower min", args: args{1 * time.Millisecond, maxStartDelay}, exactTime: false},
		{name: "higher max", args: args{minStartDelay, 24 * time.Hour}, exactTime: false},
		{name: "both match min", args: args{minStartDelay, minStartDelay}, exactTime: false},
		{name: "both match max", args: args{maxStartDelay, maxStartDelay}, exactTime: false},
		{name: "both under min", args: args{1 * time.Millisecond, 24 * time.Millisecond}, exactTime: false},
		{name: "both higher than max", args: args{1 * time.Hour, 24 * time.Hour}, exactTime: false},
	}

	// Make sure our constants are what we expect
	assert.Equal(t, 1*time.Minute, minStartDelay)
	assert.Equal(t, 15*time.Minute, maxStartDelay)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delayedTime := _calcRandomStartupDelayer(tt.args.min, tt.args.max, nowFunc)
			if tt.exactTime {
				assert.Equalf(t, staticNow, delayedTime, "calcRandomStartupDelayer(%v, %v)", tt.args.min, tt.args.max)
			} else {
				// Returned time should be between min/max
				assert.LessOrEqual(t, delayedTime, staticNow.Add(maxStartDelay))
				assert.GreaterOrEqual(t, delayedTime, staticNow.Add(minStartDelay))
			}
		})
	}
}
