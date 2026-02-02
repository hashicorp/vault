// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package rotation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

// this string is technically wrong - the schedule is mismatched from the schedule string.
const oldSched = "{\"schedule\":{\"Second\":1,\"Minute\":10376293541461622783,\"Hour\":9223372036871553023,\"Dom\":9223372041149743102,\"Month\":9223372036854783998,\"Dow\":9223372036854775935,\"Location\":{}},\"rotation_window\":60000000000,\"rotation_schedule\":\"0 0 * * *\",\"rotation_period\":0,\"next_vault_rotation\":\"2025-10-21T13:01:06.70935-04:00\",\"last_vault_rotation\":\"0001-01-01T00:00:00Z\"}"

func TestMarshalSchedule(t *testing.T) {
	cases := []struct {
		name string
		in   *RotationSchedule
		out  string
	}{
		{
			"basic",
			makeTestSchedule("* * * * *", 60, 0),
			makeTestScheduleString("* * * * *", 60, 0),
		},
		{
			"period",
			makeTestSchedule("", 0, 100),
			makeTestScheduleString("", 0, 100),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bt, err := json.Marshal(c.in)
			if err != nil {
				t.FailNow()
			}
			if string(bt) != c.out {
				t.Fatalf("marshal output didn't match, expected %q, got %q", c.out, string(bt))
			}
		})
	}
}

func TestUnmarshalSchedule(t *testing.T) {
	cases := []struct {
		name string
		in   string
		out  *RotationSchedule
	}{
		{
			"basic",
			makeTestScheduleString("* * * * *", 60, 0),
			makeTestSchedule("* * * * *", 60, 0),
		},
		{
			"backwards",
			oldSched,
			makeTestSchedule("0 0 * * *", 60, 0),
		},
		{
			"time zone",
			makeTestScheduleString("CRON_TZ=Asia/Tokyo 1/2 * * * *", 45, 0),
			makeTestSchedule("CRON_TZ=Asia/Tokyo 1/2 * * * *", 45, 0),
		},
		{
			"period",
			makeTestScheduleString("", 0, 60),
			makeTestSchedule("", 0, 60),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sched := &RotationSchedule{}
			err := json.Unmarshal([]byte(c.in), sched)
			if err != nil {
				t.Fatalf("got an unexpected error: %s", err)
			}

			if !reflect.DeepEqual(sched.Schedule, c.out.Schedule) {
				t.Fatal("schedule mismatch")
			}
		})
	}
}

// TestThereAndBack validates that a RotationSchedule comes back the same after being json-ed
func TestThereAndBackSchedule(t *testing.T) {
	cases := []struct {
		name  string
		sched *RotationSchedule
	}{
		{
			"basic",
			makeTestSchedule("* * * * *", 60, 0),
		},
		{
			"period",
			makeTestSchedule("", 0, 24601),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			js, err := json.Marshal(c.sched)
			if err != nil {
				t.Fatalf("couldn't marshal json: %s", err)
			}
			out := &RotationSchedule{}
			err = json.Unmarshal(js, out)
			if err != nil {
				t.Fatalf("couldn't unmarshal json: %s", err)
			}

			if !reflect.DeepEqual(out, c.sched) {
				t.Fatalf("input and output non-equal\nin:  %v\nout: %v", c.sched, out)
			}
		})
	}
}

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

// makeTestSchedule is a test utility to make a RotationSchedule. If it is given something weird, it will
// panic.
func makeTestSchedule(schedule string, window, period int) *RotationSchedule {
	var cronSched *cron.SpecSchedule
	var err error

	if schedule == "" {
		cronSched = nil
	} else {
		cronSched, err = DefaultScheduler.Parse(schedule)
		// eat the error so we can use the function directly in struct construction in test cases.
		if err != nil {
			panic(err)
		}
	}

	return &RotationSchedule{
		RotationWindow:    time.Duration(window) * time.Second,
		RotationPeriod:    time.Duration(period) * time.Second,
		Schedule:          cronSched,
		RotationSchedule:  schedule,
		LastVaultRotation: time.Time{},
		NextVaultRotation: time.Time{},
	}
}

// makeTestScheduleString is a test helper that creates a json string that matches our expected json output
// when marshaling a RotationSchedule.
func makeTestScheduleString(schedule string, window, period int) string {
	return fmt.Sprintf(`{"rotation_window":%d,"rotation_schedule":"%s","rotation_period":%d,"next_vault_rotation":"0001-01-01T00:00:00Z","last_vault_rotation":"0001-01-01T00:00:00Z"}`,
		time.Duration(window)*time.Second,
		schedule,
		time.Duration(period)*time.Second,
	)
}
