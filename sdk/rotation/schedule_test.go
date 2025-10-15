// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rotation

import (
	"testing"
	"time"
)

func TestParseSchedule(t *testing.T) {
	// Actual schedule-parsing tests are the responsibility of the library,
	// (currently robfig/cron), but here are some cases for our specific functionality
	cases := []struct {
		name      string
		in        string
		shouldErr bool
		location  *time.Location
	}{
		{"force local to utc", "* * * * *", false, time.UTC},
		{"seconds are invalid", "* * * * * *", true, nil},

		// Specifically override this usage
		{"custom timezone", "CRON_TZ=Asia/Tokyo * * * * *", false, time.UTC},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sched, err := DefaultScheduler.Parse(c.in)
			if c.shouldErr { // should-err tests end here no matter what
				if err == nil {
					t.Error("should have errored, but didn't")
				} else {
					return
				}
			}

			if err != nil {
				t.Errorf("got unexpected error: %v", err)
			}

			// check the tz
			if sched.Location.String() != c.location.String() {
				t.Errorf("wrong tz, expected %s, got %s", c.location.String(), sched.Location.String())
			}
		})
	}
}
