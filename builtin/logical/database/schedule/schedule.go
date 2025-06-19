// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	// Minimum allowed value for rotation_window
	minRotationWindowSeconds = 3600
	parseOptions             = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
)

type Scheduler interface {
	Parse(string) (*cron.SpecSchedule, error)
	ValidateRotationWindow(int) error
}

var _ Scheduler = &DefaultSchedule{}

type DefaultSchedule struct{}

func (d *DefaultSchedule) Parse(rotationSchedule string) (*cron.SpecSchedule, error) {
	parser := cron.NewParser(parseOptions)
	schedule, err := parser.Parse(rotationSchedule)
	if err != nil {
		return nil, err
	}
	sched, ok := schedule.(*cron.SpecSchedule)
	if !ok {
		return nil, fmt.Errorf("invalid rotation schedule")
	}
	// Force the location to be UTC instead of the local timezone
	sched.Location = time.UTC
	return sched, nil
}

func (d *DefaultSchedule) ValidateRotationWindow(s int) error {
	if s < minRotationWindowSeconds {
		return fmt.Errorf("rotation_window must be %d seconds or more", minRotationWindowSeconds)
	}
	return nil
}
