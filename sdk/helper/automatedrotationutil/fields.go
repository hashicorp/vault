// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automatedrotationutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/rotation"
)

var (
	ErrRotationMutuallyExclusiveFields = errors.New("mutually exclusive fields rotation_schedule and rotation_period were both specified; only one of them can be provided")
	ErrRotationManagerUnsupported      = errors.New("rotation manager capabilities not supported in Vault community edition")
)

// AutomatedRotationParams contains a set of common parameters that plugins
// can use for setting automated credential rotation.
type AutomatedRotationParams struct {
	// RotationSchedule is the CRON-style rotation schedule.
	RotationSchedule string `json:"rotation_schedule"`
	// RotationWindow specifies the amount of time in which the rotation is allowed to
	// occur starting from a given rotation_schedule.
	RotationWindow time.Duration `json:"rotation_window"`
	// RotationPeriod is an alternate choice for simple time-to-live based rotation timing.
	RotationPeriod time.Duration `json:"rotation_period"`

	// If set, will deregister all registered rotation jobs from the RotationManager for plugin.
	DisableAutomatedRotation bool `json:"disable_automated_rotation"`
}

// ParseAutomatedRotationFields provides common field parsing to embedding structs.
func (p *AutomatedRotationParams) ParseAutomatedRotationFields(d *framework.FieldData) error {
	rotationScheduleRaw, scheduleOk := d.GetOk("rotation_schedule")
	rotationWindowSecondsRaw, windowOk := d.GetOk("rotation_window")
	rotationPeriodSecondsRaw, periodOk := d.GetOk("rotation_period")
	disableRotation, disableRotationOk := d.GetOk("disable_automated_rotation")

	if scheduleOk {
		if periodOk && rotationPeriodSecondsRaw.(int) != 0 && rotationScheduleRaw.(string) != "" {
			return ErrRotationMutuallyExclusiveFields
		}
		p.RotationSchedule = rotationScheduleRaw.(string)

		// parse schedule to ensure it is valid
		if p.RotationSchedule != "" {
			_, err := rotation.DefaultScheduler.Parse(p.RotationSchedule)
			if err != nil {
				return fmt.Errorf("failed to parse provided rotation_schedule: %w", err)
			}
		}
	}

	if windowOk {
		if periodOk && rotationPeriodSecondsRaw.(int) != 0 && rotationWindowSecondsRaw.(int) != 0 {
			return fmt.Errorf("rotation_window does not apply to period")
		}
		rotationWindowSeconds := rotationWindowSecondsRaw.(int)
		p.RotationWindow = time.Duration(rotationWindowSeconds) * time.Second
	}

	if periodOk {
		rotationPeriodSeconds := rotationPeriodSecondsRaw.(int)
		p.RotationPeriod = time.Duration(rotationPeriodSeconds) * time.Second
	}

	if (windowOk && rotationWindowSecondsRaw.(int) != 0) && !scheduleOk {
		return fmt.Errorf("cannot use rotation_window without rotation_schedule")
	}

	if disableRotationOk {
		p.DisableAutomatedRotation = disableRotation.(bool)
	}

	return nil
}

// PopulateAutomatedRotationData adds PluginIdentityTokenParams info into the given map.
func (p *AutomatedRotationParams) PopulateAutomatedRotationData(m map[string]interface{}) {
	m["rotation_schedule"] = p.RotationSchedule
	m["rotation_window"] = p.RotationWindow.Seconds()
	m["rotation_period"] = p.RotationPeriod.Seconds()
	m["disable_automated_rotation"] = p.DisableAutomatedRotation
}

func (p *AutomatedRotationParams) ShouldRegisterRotationJob() bool {
	return p.RotationSchedule != "" || p.RotationPeriod != 0
}

func (p *AutomatedRotationParams) ShouldDeregisterRotationJob() bool {
	return p.DisableAutomatedRotation || (p.RotationSchedule == "" && p.RotationPeriod == 0)
}

// AddAutomatedRotationFields adds plugin identity token fields to the given
// field schema map.
func AddAutomatedRotationFields(m map[string]*framework.FieldSchema) {
	fields := map[string]*framework.FieldSchema{
		"rotation_schedule": {
			Type:        framework.TypeString,
			Description: "CRON-style string that will define the schedule on which rotations should occur. Mutually exclusive with rotation_period",
		},
		"rotation_window": {
			Type:        framework.TypeDurationSecond,
			Description: "Specifies the amount of time in which the rotation is allowed to occur starting from a given rotation_schedule",
		},
		"rotation_period": {
			Type:        framework.TypeDurationSecond,
			Description: "TTL for automatic credential rotation of the given username. Mutually exclusive with rotation_schedule",
		},
		"disable_automated_rotation": {
			Type:        framework.TypeBool,
			Description: "If set to true, will deregister all registered rotation jobs from the RotationManager for the plugin.",
		},
	}

	for name, schema := range fields {
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("adding field %q would overwrite existing field", name))
		}
		m[name] = schema
	}
}
