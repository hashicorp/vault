// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	// Minimum allowed value for rotation_window
	minRotationWindowSeconds = 3600
	parseOptions             = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow
)

// RootSchedule holds the parsed and unparsed versions of the schedule, along with the projected next rotation time.
type RootSchedule struct {
	Schedule          *cron.SpecSchedule `json:"schedule"`
	RotationWindow    time.Duration      `json:"rotation_window"` // seconds of window
	RotationSchedule  string             `json:"rotation_schedule"`
	NextVaultRotation time.Time          `json:"next_vault_rotation"`
}

type Scheduler interface {
	Parse(rotationSchedule string) (*cron.SpecSchedule, error)
	ValidateRotationWindow(s int) error
	NextRotationTimeFromInput(next *RootSchedule, input time.Time) time.Time
	IsInsideRotationWindow(rotation *RootSchedule, t time.Time) bool
	ShouldRotate(rotation *RootSchedule, priority int64, t time.Time) bool
	NextRotationTime(next *RootSchedule) time.Time
	SetNextVaultRotation(next *RootSchedule, t time.Time)
}

var DefaultScheduler Scheduler = &DefaultSchedule{}

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
	return sched, nil
}

func (d *DefaultSchedule) ValidateRotationWindow(s int) error {
	if s < minRotationWindowSeconds {
		return fmt.Errorf("rotation_window must be %d seconds or more", minRotationWindowSeconds)
	}
	return nil
}

// NextRotationTime calculates the next scheduled rotation
func (d *DefaultSchedule) NextRotationTime(next *RootSchedule) time.Time {
	return next.Schedule.Next(time.Now())
}

// NextRotationTimeFromInput calculates and returns the next rotation time based on the provided  schedule and input time
func (d *DefaultSchedule) NextRotationTimeFromInput(next *RootSchedule, input time.Time) time.Time {
	return next.Schedule.Next(input)
}

// IsInsideRotationWindow checks if the current time is before the calculated end of the rotation window,
// to make sure that t time is within the specified rotation window
// It returns true if rotation window is not specified
func (d *DefaultSchedule) IsInsideRotationWindow(rotation *RootSchedule, t time.Time) bool {
	if rotation.RotationWindow != 0 {
		return t.Before(rotation.NextVaultRotation.Add(rotation.RotationWindow))
	}
	return true
}

// ShouldRotate checks if the rotation should occur based on  priority, current time, and rotation window
// It returns true if the priority is less than or equal to the current time and the current time is within the rotation window
func (d *DefaultSchedule) ShouldRotate(rotation *RootSchedule, priority int64, t time.Time) bool {
	return priority <= t.Unix() && d.IsInsideRotationWindow(rotation, t)
}

// SetNextVaultRotation calculates the next rotation time of a given schedule based on the time.
func (d *DefaultSchedule) SetNextVaultRotation(next *RootSchedule, t time.Time) {
	next.NextVaultRotation = next.Schedule.Next(t)
}
