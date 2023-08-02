// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import (
	"time"
)

// RotationOptions is an embeddable struct to capture common lease
// settings between a Secret and Auth
type RotationOptions struct {
	// Specifies the amount of time Vault should wait before rotating the
	// password. The minimum is 5 seconds.
	// Mutually exclusive with Schedule
	Period *time.Duration `json:"rotation_period"`

	// Schedule is the maximum duration that this secret is valid for.
	// Mutually exclusive with Period
	// TODO: custom type to handle chron style schedule?
	Schedule *string `json:"rotation_schedule"`
}
