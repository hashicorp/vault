// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automatedrotationutil

import (
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
)

var (
	ErrRotationMutuallyExclusiveFields = errors.New("mutually exclusive fields rotation_schedule and rotation_ttl were both specified; only one of them can be provided")
	ErrRotationManagerUnsupported      = errors.New("rotation manager capabilities not supported in Vault community edition")
)

// AutomatedRotationParams contains a set of common parameters that plugins
// can use for setting automated credential rotation.
type AutomatedRotationParams struct {
	// RotationSchedule is the CRON-style rotation schedule.
	RotationSchedule string `json:"rotation_schedule"`
	// RotationWindow specifies the amount of time in which the rotation is allowed to occur starting from a given rotation_schedule.
	RotationWindow int `json:"rotation_window"`
	// RotationTTL is an alternate choice for simple time-to-live based rotation timing.
	RotationTTL int `json:"rotation_ttl"`

	// RotationID is the unique ID of the registered rotation job.
	// Used by the plugin to track the rotation.
	RotationID string `json:"rotation_id"`
}

// ParseAutomatedRotationFields provides common field parsing to embedding structs.
func (p *AutomatedRotationParams) ParseAutomatedRotationFields(d *framework.FieldData) error {
	rotationScheduleRaw, scheduleOk := d.GetOk("rotation_schedule")
	rotationWindowRaw, windowOk := d.GetOk("rotation_window")
	rotationTTLRaw, ttlOk := d.GetOk("rotation_ttl")

	if scheduleOk {
		if ttlOk {
			return ErrRotationMutuallyExclusiveFields
		}
		p.RotationSchedule = rotationScheduleRaw.(string)
	}

	if windowOk {
		if ttlOk {
			return fmt.Errorf("rotation_window does not apply to ttl")
		}
		p.RotationWindow = rotationWindowRaw.(int)
	}

	if ttlOk {
		p.RotationTTL = rotationTTLRaw.(int)
	}

	if (scheduleOk && !windowOk) || (windowOk && !scheduleOk) {
		return fmt.Errorf("must include both schedule and window")
	}

	return nil
}

// PopulateAutomatedRotationData adds PluginIdentityTokenParams info into the given map.
func (p *AutomatedRotationParams) PopulateAutomatedRotationData(m map[string]interface{}) {
	m["rotation_schedule"] = p.RotationSchedule
	m["rotation_window"] = p.RotationWindow
	m["rotation_ttl"] = p.RotationTTL
	m["rotation_id"] = p.RotationID
}

// AddAutomatedRotationFields adds plugin identity token fields to the given
// field schema map.
func AddAutomatedRotationFields(m map[string]*framework.FieldSchema) {
	fields := map[string]*framework.FieldSchema{
		"rotation_schedule": {
			Type:        framework.TypeString,
			Description: "CRON-style string that will define the schedule on which rotations should occur. Mutually exclusive with TTL",
		},
		"rotation_window": {
			Type:        framework.TypeInt,
			Description: "Specifies the amount of time in which the rotation is allowed to occur starting from a given rotation_schedule",
		},
		"rotation_ttl": {
			Type:        framework.TypeInt,
			Description: "TTL for automatic credential rotation of the given username. Mutually exclusive with rotation_schedule",
		},
		"rotation_id": {
			Type:        framework.TypeInt,
			Description: "Unique ID of the registered rotation job",
		},
	}

	for name, schema := range fields {
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("adding field %q would overwrite existing field", name))
		}
		m[name] = schema
	}
}
