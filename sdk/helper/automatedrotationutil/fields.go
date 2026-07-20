// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package automatedrotationutil

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
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
	RotationWindow time.Duration `json:"rotation_window"`
	// RotationPeriod is an alternate choice for simple time-to-live based rotation timing.
	RotationPeriod time.Duration `json:"rotation_period"`

	// If set, will deregister all registered rotation jobs from the RotationManager for plugin.
	DisableAutomatedRotation bool `json:"disable_automated_rotation"`

	// RotationPolicy defines the policy to use when performing retries.
	RotationPolicy string `json:"rotation_policy"`
}

type RotationInfoResponseParams struct {
	// NextVaultRotation represents the next time Vault is expected to rotate the credential.
	NextVaultRotation time.Time `json:"next_vault_rotation"`

	// LastVaultRotation represents the time the credential was initially onboarded to the RM or last rotated.
	LastVaultRotation time.Time `json:"last_vault_rotation"`

	// RotationID represents the ID of this credential.
	RotationID string `json:"rotation_id"`
}

// ParseAutomatedRotationFields provides common field parsing to embedding structs.
func (p *AutomatedRotationParams) ParseAutomatedRotationFields(d *framework.FieldData) error {
	rotationScheduleRaw, scheduleOk := d.GetOk("rotation_schedule")
	rotationWindowSecondsRaw, windowOk := d.GetOk("rotation_window")
	rotationPeriodSecondsRaw, periodOk := d.GetOk("rotation_period")
	disableRotation, disableRotationOk := d.GetOk("disable_automated_rotation")
	rotationPolicyRaw, policyOk := d.GetOk("rotation_policy")

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

	if policyOk {
		p.RotationPolicy = rotationPolicyRaw.(string)
	}

	return nil
}

type RotationJobOperationResponse struct {
	OperationPerformed string
	Logger             hclog.Logger
	RotationInfo       *rotation.RotationInfo
}

type (
	RegisterErrorHandler   func(error) error
	DeregisterErrorHandler func(error) error
)

type ErrorHandlers struct {
	RegisterErrorHandler   RegisterErrorHandler
	DeregisterErrorHandler DeregisterErrorHandler
}

// HandleRotationJob is a helper method for registering or deregistering rotation jobs.
// This wraps the two individual reegister and deregister methods and determines which to call based
// on the request parameters provided.
// Returns a formatted error based on the operation that failed.
// Use this method unless there is a specific need to handle the individual operations independently
// or in a nonstandard way.
func (p *AutomatedRotationParams) HandleRotationJob(ctx context.Context, b *framework.Backend, fieldData *framework.FieldData, req *logical.Request) (RotationJobOperationResponse, error) {
	return p.HandleRotationJobWithErrorHandlers(ctx, b, fieldData, req, &ErrorHandlers{})
}

// HandleRotationJobWithErrorHandlers is the same as HandleRotationJob but allows callers to provide custom error handlers for errors returned from either the register or deregister operations.
func (p *AutomatedRotationParams) HandleRotationJobWithErrorHandlers(ctx context.Context, b *framework.Backend, fieldData *framework.FieldData, req *logical.Request, errorHandlers *ErrorHandlers) (RotationJobOperationResponse, error) {
	resp := RotationJobOperationResponse{
		Logger: b.Logger(),
	}

	if err := p.ParseAutomatedRotationFields(fieldData); err != nil {
		return resp, err
	}

	if p.ShouldDeregisterRotationJob() {
		resp.OperationPerformed = rotation.PerformedDeregistration
		if err := p.HandleDeregisterRotationJob(ctx, b, req); err != nil {
			if errorHandlers.DeregisterErrorHandler != nil {
				err = errorHandlers.DeregisterErrorHandler(err)
			}
			return resp, fmt.Errorf("failed to deregister rotation job: %w", err)
		}
	} else if p.ShouldRegisterRotationJob() {
		resp.OperationPerformed = rotation.PerformedRegistration
		registerResp, err := p.HandleRegisterRotationJob(ctx, b, req)
		if err != nil {
			if errorHandlers.RegisterErrorHandler != nil {
				err = errorHandlers.RegisterErrorHandler(err)
			}
			return resp, fmt.Errorf("failed to register rotation job: %w", err)
		}
		resp.RotationInfo = registerResp
	}

	return resp, nil
}

// HandleStorageErrorAfterRotationJob is a helper method to log and wrap errors from storage operations
// and should only be used after performing rotation job operations.
// The caller is responsible for determining when it is appropriate to call this method.
func (r *RotationJobOperationResponse) HandleStorageErrorAfterRotationJob(req *logical.Request, err error) error {
	if err != nil {
		if r.OperationPerformed != "" {
			// Write to storage failed but a rotation job operation succeeded beforehand, return this as an error.
			msg := "write to storage failed but the rotation operation still succeeded; this may cause next_vault_rotation and last_vault_rotation values to be out-of-sync with the actual rotation schedule"
			r.Logger.Error(msg,
				"operation", r.OperationPerformed, "mount", req.MountPoint, "path", req.Path, "error", err.Error())
			return fmt.Errorf("%s; operation=%s, mount=%s, path=%s, error=%s", msg, r.OperationPerformed, req.MountPoint, req.Path, err)
		} else {
			// There was no rotation operation performed, so just return the storage error.
			return err
		}
	}

	return nil
}

// HandleRegisterRotationJob is a helper method to register rotation jobs from a plugin.
// Use HandleRotationJob or HandleRotationJobWithErrorHandlers instead when possible.
// Returns the raw system error.
// Callers are responsible for validating when to register a RotationJob.
// Callers are responsible for performing cleanup or storage rollbacks if necessary.
func (p *AutomatedRotationParams) HandleRegisterRotationJob(ctx context.Context, b *framework.Backend, req *logical.Request) (*rotation.RotationInfo, error) {
	registerReq := &rotation.RotationJobConfigureRequest{
		MountPoint:       req.MountPoint,
		ReqPath:          req.Path,
		RotationSchedule: p.RotationSchedule,
		RotationWindow:   p.RotationWindow,
		RotationPeriod:   p.RotationPeriod,
	}

	b.Logger().Debug("Registering rotation job", "mount", req.MountPoint, "path", req.Path)
	resp, err := b.System().RegisterRotationJobWithResponse(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// HandleDeregisterRotationJob is a helper method to deregister rotation jobs from a plugin.
// Use HandleRotationJob or HandleRotationJobWithErrorHandlers instead when possible.
// Returns the raw system error.
// Callers are responsible for validating when to deregister a RotationJob.
// Callers are responsible for performing cleanup or storage rollbacks if necessary.
func (p *AutomatedRotationParams) HandleDeregisterRotationJob(ctx context.Context, b *framework.Backend, req *logical.Request) error {
	deregisterReq := &rotation.RotationJobDeregisterRequest{
		MountPoint: req.MountPoint,
		ReqPath:    req.Path,
	}

	b.Logger().Debug("Deregistering rotation job", "mount", req.MountPoint, "path", req.Path)
	if err := b.System().DeregisterRotationJob(ctx, deregisterReq); err != nil {
		return err
	}

	return nil
}

// Use PopulateSetAutomatedRotationData instead, *unless* all these
// fields are necessary to maintain backwards compatibility with the plugin's pre-existing response API.
// PopulateAutomatedRotationData adds AutomatedRotationParams info into the given map.
func (p *AutomatedRotationParams) PopulateAutomatedRotationData(m map[string]interface{}) {
	m["rotation_schedule"] = p.RotationSchedule
	m["rotation_window"] = p.RotationWindow.Seconds()
	m["rotation_period"] = p.RotationPeriod.Seconds()
	m["disable_automated_rotation"] = p.DisableAutomatedRotation
	m["rotation_policy"] = p.RotationPolicy
}

// PopulateSetAutomatedRotationData adds AutomatedRotationParams info into the given map, based
// on which fields were set for rotation. Setting a rotation schedule will not return a rotation
// period, and setting a rotation period will not return a rotation schedule or rotation window.
func (p *AutomatedRotationParams) PopulateSetAutomatedRotationData(m map[string]interface{}) {
	// Always set these even if they are zero values, to avoid confusion.
	m["disable_automated_rotation"] = p.DisableAutomatedRotation
	m["rotation_policy"] = p.RotationPolicy

	// Set both of these if a schedule is set.
	if p.RotationSchedule != "" {
		m["rotation_schedule"] = p.RotationSchedule
		m["rotation_window"] = p.RotationWindow.Seconds()
	}

	// Set this if a period is set.
	if p.RotationPeriod != 0 {
		m["rotation_period"] = p.RotationPeriod.Seconds()
	}
}

// PopulateRotationInfo adds RotationInfoResponseParams info into the given map.
func (p *RotationInfoResponseParams) PopulateRotationInfo(m map[string]interface{}) {
	// Only set last_vault_rotation and next_vault_rotation if they are non-zero
	if !p.LastVaultRotation.IsZero() {
		m["last_vault_rotation"] = p.LastVaultRotation.UTC()
	} else {
		m["last_vault_rotation"] = nil
	}

	if !p.NextVaultRotation.IsZero() {
		m["next_vault_rotation"] = p.NextVaultRotation.UTC()
	} else {
		m["next_vault_rotation"] = nil
	}
}

// SetRotationInfo sets the rotation info. It ensures a consistent format across different uses.
// Plugins should use this when registering credentials or in the RotateCredential callback to keep rotation state up to date.
func (p *RotationInfoResponseParams) SetRotationInfo(r *rotation.RotationInfo) {
	if r != nil {
		// LastVaultRotation is only provided by the RM on rotateCredential requests
		// On a registration, we do not need to set this info on the credential.
		if !r.LastVaultRotation.IsZero() {
			// only set if provided
			// only care about precision up until seconds, drop everything below
			p.LastVaultRotation = r.LastVaultRotation.UTC().Truncate(time.Second)
		}
		p.NextVaultRotation = r.NextVaultRotation.UTC().Truncate(time.Second)
		p.RotationID = r.RotationID
	}
}

// SetLastVaultRotation sets the LastVaultRotation. It ensures a consistent format across different uses.
// Plugins should only use this when manually rotating credentials to keep rotation state up to date.
func (p *RotationInfoResponseParams) SetLastVaultRotation() {
	p.LastVaultRotation = time.Now().UTC().Truncate(time.Second)
}

// GetTTL computes the TTL in seconds until expireTime from now.
// This method should be used by plugins to compute TTL values for static credentials.
func (p *RotationInfoResponseParams) GetTTL() int64 {
	ttl := int64(p.NextVaultRotation.Sub(time.Now()).Seconds())
	// a negative value here means the time has arrived, but the queue hasn't been checked yet. If the queue is checked
	// every <n> seconds, we could get a value as low as -n. This can be a little confusing on the user end, so we clamp
	// the value to zero. To quote another doc, "Users should not trust passwords with a zero ttl, as they are likely
	// in the process of being rotated and will quickly become invalidated."
	if ttl < 0 {
		ttl = 0
	}

	return ttl
}

func (p *AutomatedRotationParams) ShouldRegisterRotationJob() bool {
	return p.HasNonzeroRotationValues()
}

func (p *AutomatedRotationParams) ShouldDeregisterRotationJob() bool {
	return p.DisableAutomatedRotation || (p.RotationSchedule == "" && p.RotationPeriod == 0)
}

// HasNonzeroRotationValues returns true if either of the primary rotation values (RotationSchedule or RotationPeriod)
// are not the zero value.
func (p *AutomatedRotationParams) HasNonzeroRotationValues() bool {
	return p.RotationSchedule != "" || p.RotationPeriod != 0
}

// AddAutomatedRotationFieldsWithGroup adds rotation fields to the given field schema map
// the fields are associated to the provided display attribute group
func AddAutomatedRotationFieldsWithGroup(m map[string]*framework.FieldSchema, group string) {
	fields := map[string]*framework.FieldSchema{
		"rotation_schedule": {
			Type:        framework.TypeString,
			Description: "CRON-style string that will define the schedule on which rotations should occur. Mutually exclusive with rotation_period",
			DisplayAttrs: &framework.DisplayAttributes{
				Group: group,
			},
		},
		"rotation_window": {
			Type:        framework.TypeDurationSecond,
			Description: "Specifies the amount of time in which the rotation is allowed to occur starting from a given rotation_schedule",
			DisplayAttrs: &framework.DisplayAttributes{
				Group: group,
			},
		},
		"rotation_period": {
			Type:        framework.TypeDurationSecond,
			Description: "TTL for automatic credential rotation of the given username. Mutually exclusive with rotation_schedule",
			DisplayAttrs: &framework.DisplayAttributes{
				Group: group,
			},
		},
		"disable_automated_rotation": {
			Type:        framework.TypeBool,
			Description: "If set to true, will deregister all registered rotation jobs from the RotationManager for the plugin.",
			DisplayAttrs: &framework.DisplayAttributes{
				Group: group,
			},
		},
		"rotation_policy": {
			Type:        framework.TypeString,
			Description: "Defines the rotation policy to use when performing automated rotations.",
			DisplayAttrs: &framework.DisplayAttributes{
				Group: group,
			},
		},
	}

	for name, schema := range fields {
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("adding field %q would overwrite existing field", name))
		}
		m[name] = schema
	}
}

// stubbing original function for compatibility
// AddAutomatedRotationFieldsWithGroup should be used directly
// future utils that define fields should include a group parameter
func AddAutomatedRotationFields(m map[string]*framework.FieldSchema) {
	AddAutomatedRotationFieldsWithGroup(m, "default")
}

// Equals returns true if the automated rotation parameters match the other instance.
// Useful for detecting configuration changes after parsing new field data.
func (p *AutomatedRotationParams) Equals(other AutomatedRotationParams) bool {
	return *p == other
}
