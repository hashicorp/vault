// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logical

import "context"

// ObservationRecorder records an observation.
type ObservationRecorder interface {
	RecordObservationFromPlugin(ctx context.Context, observationType string, data map[string]interface{}) error
}
