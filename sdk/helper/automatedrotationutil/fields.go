// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automatedrotationutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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
	RotationWindow int `json:"rotation_window"`
	// RotationPeriod is an alternate choice for simple time-to-live based rotation timing.
	RotationPeriod int `json:"rotation_period"`

	// If set, will deregister all registered rotation jobs from the RotationManager for plugin.
	DisableAutomatedRotation bool `json:"disable_automated_rotation"`
}

// ParseAutomatedRotationFields provides common field parsing to embedding structs.
func (p *AutomatedRotationParams) ParseAutomatedRotationFields(d *framework.FieldData) error {
	rotationScheduleRaw, scheduleOk := d.GetOk("rotation_schedule")
	rotationWindowRaw, windowOk := d.GetOk("rotation_window")
	rotationPeriodRaw, periodOk := d.GetOk("rotation_period")

	if scheduleOk {
		if periodOk {
			return ErrRotationMutuallyExclusiveFields
		}
		p.RotationSchedule = rotationScheduleRaw.(string)
	}

	if windowOk {
		if periodOk {
			return fmt.Errorf("rotation_window does not apply to period")
		}
		p.RotationWindow = rotationWindowRaw.(int)
	}

	if periodOk {
		p.RotationPeriod = rotationPeriodRaw.(int)
	}

	if (scheduleOk && !windowOk) || (windowOk && !scheduleOk) {
		return fmt.Errorf("must include both schedule and window")
	}

	_, err := rotation.DefaultScheduler.Parse(p.RotationSchedule)
	if err != nil {
		return fmt.Errorf("failed to parse provided rotation_schedule: %w", err)
	}

	p.DisableAutomatedRotation = d.Get("disable_automated_rotation").(bool)

	return nil
}

// PopulateAutomatedRotationData adds PluginIdentityTokenParams info into the given map.
func (p *AutomatedRotationParams) PopulateAutomatedRotationData(m map[string]interface{}) {
	m["rotation_schedule"] = p.RotationSchedule
	m["rotation_window"] = p.RotationWindow
	m["rotation_period"] = p.RotationPeriod
	m["disable_automated_rotation"] = p.DisableAutomatedRotation
}

func (p *AutomatedRotationParams) ShouldRegisterRotationJob() bool {
	return p.RotationSchedule != "" || p.RotationPeriod != 0
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
			Type:        framework.TypeInt,
			Description: "Specifies the amount of time in which the rotation is allowed to occur starting from a given rotation_schedule",
		},
		"rotation_period": {
			Type:        framework.TypeInt,
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

// HandleRotationJobOperation is a helper method wrapper to the two individual helper methods
// that register and deregister a rotation job each, and returns a standardized error response
// with a warning in the event of an unexpected error.
// revertStorageFunc is a function that should revert the original config back in storege before
// any fields were updated.
func (p *AutomatedRotationParams) HandleRotationJobOperation(ctx context.Context, b *framework.Backend, req *logical.Request, revertStorageFunc func() error) *logical.Response {
	if p.DisableAutomatedRotation {
		// Disable automated rotation and deregister credentials if required.
		if err := p.HandleDeregisterRotationJob(ctx, b, req); err != nil {
			resp := logical.ErrorResponse("error deregistering rotation job but config was successfully updated: %s", err)
			resp.AddWarning("config was successfully updated despite failing to disable automated rotation")

			if err := revertStorageFunc(); err != nil {
				return resp
			}

			return nil
		}
	} else {
		// Register the rotation job if it's required.
		if err := p.HandleRegisterRotationJob(ctx, b, req); err != nil {
			resp := logical.ErrorResponse("error registering rotation job but config was successfully updated: %s", err)
			resp.AddWarning("config was successfully updated despite failing to enable automated rotation")

			if err := revertStorageFunc(); err != nil {
				return resp
			}

			return nil
		}
	}

	return nil
}

// HandleRegisterRotationJob is a helper method to register rotation jobs from a plugin.
// It is up to callers to perform cleanup or storage backouts if necessary.
func (p *AutomatedRotationParams) HandleRegisterRotationJob(ctx context.Context, b *framework.Backend, req *logical.Request) error {
	// Now that the root config is set up, register the rotation job if it's required.
	if p.ShouldRegisterRotationJob() {
		cfgReq := &rotation.RotationJobConfigureRequest{
			MountType:        req.MountType,
			ReqPath:          req.Path,
			RotationSchedule: p.RotationSchedule,
			RotationWindow:   p.RotationWindow,
			RotationPeriod:   p.RotationPeriod,
		}

		b.Logger().Debug("Registering rotation job", "mount", req.MountPoint+req.Path)
		if _, err := b.System().RegisterRotationJob(ctx, cfgReq); err != nil {
			return err
		}
	}

	return nil
}

// HandleDeregisterRotationJob is a helper method to deregister rotation jobs from a plugin.
// It is up to callers to perform cleanup or storage backouts if necessary.
func (p *AutomatedRotationParams) HandleDeregisterRotationJob(ctx context.Context, b *framework.Backend, req *logical.Request) error {
	// Disable Automated Rotation and Deregister credentials if required
	deregisterReq := &rotation.RotationJobDeregisterRequest{
		MountType: req.MountType,
		ReqPath:   req.Path,
	}

	b.Logger().Debug("Deregistering rotation job", "mount", req.MountPoint+req.Path)
	if err := b.System().DeregisterRotationJob(ctx, deregisterReq); err != nil {
		return err
	}

	return nil
}
