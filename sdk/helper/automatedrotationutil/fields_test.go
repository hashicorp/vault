// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automatedrotationutil

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/framework"
)

var schemaMap = map[string]*framework.FieldSchema{
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

func TestParseAutomatedRotationFields(t *testing.T) {
	tests := []struct {
		name           string
		data           *framework.FieldData
		expectedParams *AutomatedRotationParams
		expectedError  string
	}{
		{
			name: "basic",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "*/15 * * * * *",
					"rotation_window":   60,
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * * *",
				RotationWindow:           60,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "mutually-exclusive",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "*/15 * * * * *",
					"rotation_period":   15,
					"rotation_window":   60,
				},
				Schema: schemaMap,
			},
			expectedError: ErrRotationMutuallyExclusiveFields.Error(),
		},
		{
			name: "incompatible-fields",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_period": 15,
					"rotation_window": 60,
				},
				Schema: schemaMap,
			},
			expectedError: "rotation_window does not apply to period",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &AutomatedRotationParams{}
			err := p.ParseAutomatedRotationFields(tt.data)

			if err != nil && !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("ParseAutomatedRotationFields() error = %v, expected %s", err, tt.expectedError)
			}

			if err == nil && !reflect.DeepEqual(tt.expectedParams, p) {
				t.Errorf("ParseAutomatedRotationFields() error comparing params; got %v, expected %v", tt.expectedParams, p)
			}
		})
	}
}

func TestPopulateAutomatedRotationData(t *testing.T) {
	tests := []struct {
		name        string
		inputParams *AutomatedRotationParams
		expected    map[string]interface{}
	}{
		{
			name: "basic",
			expected: map[string]interface{}{
				"rotation_schedule":          "*/15 * * * * *",
				"rotation_window":            60,
				"rotation_period":            0,
				"disable_automated_rotation": false,
			},
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * * *",
				RotationWindow:           60,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := make(map[string]interface{})

			tt.inputParams.PopulateAutomatedRotationData(m)

			if !reflect.DeepEqual(m, tt.expected) {
				t.Errorf("PopulateAutomatedRotationData() error comparing values; got %v, expected %v", m, tt.expected)
			}
		})
	}
}
