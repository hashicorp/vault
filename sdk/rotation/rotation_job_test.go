// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rotation

import (
	"reflect"
	"strings"
	"testing"
)

func TestConfigureRotationJob(t *testing.T) {
	tests := []struct {
		name          string
		req           *RotationJobConfigureRequest
		expected      RotationJob
		expectedError string
	}{
		{
			name: "no rotation params",
			req: &RotationJobConfigureRequest{
				MountPoint:       "aws",
				ReqPath:          "config/root",
				RotationSchedule: "",
				RotationWindow:   60,
				RotationPeriod:   0,
			},
			expectedError: "RotationSchedule or RotationPeriod is required to set up rotation job",
		},
		{
			name: "no mount point",
			req: &RotationJobConfigureRequest{
				MountPoint:       "",
				ReqPath:          "config/root",
				RotationSchedule: "",
				RotationWindow:   60,
				RotationPeriod:   5,
			},
			expectedError: "MountPoint is required",
		},
		{
			name: "no req path",
			req: &RotationJobConfigureRequest{
				MountPoint:       "aws",
				ReqPath:          "",
				RotationSchedule: "",
				RotationWindow:   60,
				RotationPeriod:   5,
			},
			expectedError: "ReqPath is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := ConfigureRotationJob(tt.req)
			if err != nil {
				if tt.expectedError == "" {
					t.Errorf("expected no error but received an error: %s", err)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("TestConfigureRotationJob() error = %v, expected %s", err, tt.expectedError)
				}
			}

			if err == nil && !reflect.DeepEqual(tt.expected, out) {
				t.Errorf("TestConfigureRotationJob() error comparing params; got %v, expected %v", out, tt.expected)
			}
		})
	}
}
