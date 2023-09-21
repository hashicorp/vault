package schedule

import (
	"fmt"

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
	return sched, nil
}

func (d *DefaultSchedule) ValidateRotationWindow(s int) error {
	if s < minRotationWindowSeconds {
		return fmt.Errorf("rotation_window must be %d seconds or more", minRotationWindowSeconds)
	}
	return nil
}
