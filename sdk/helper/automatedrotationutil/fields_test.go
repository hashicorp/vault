// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automatedrotationutil

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
)

var schemaMap = map[string]*framework.FieldSchema{
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

func TestParseAutomatedRotationFields(t *testing.T) {
	tests := []struct {
		name           string
		data           *framework.FieldData
		expectedParams *AutomatedRotationParams
		initialParams  *AutomatedRotationParams
		expectedError  string
	}{
		{
			name: "basic",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "*/15 * * * *",
					"rotation_window":   60,
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * *",
				RotationWindow:           60 * time.Second,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "mutually-exclusive",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "*/15 * * * *",
					"rotation_period":   15,
					"rotation_window":   60,
				},
				Schema: schemaMap,
			},
			expectedError: ErrRotationMutuallyExclusiveFields.Error(),
		},
		{
			name: "bad-schedule",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "random string",
				},
				Schema: schemaMap,
			},
			expectedError: "failed to parse provided rotation_schedule",
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
		{
			name: "window-without-schedule",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_window": 60,
				},
				Schema: schemaMap,
			},
			expectedError: "cannot use rotation_window without rotation_schedule",
		},
		{
			name: "rotation-period-duration-seconds",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_period": "2m",
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           0,
				RotationPeriod:           120 * time.Second,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "rotation-window-duration-seconds",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_window":   "12h",
					"rotation_schedule": "* */2 * * *",
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "* */2 * * *",
				RotationWindow:           12 * time.Hour,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "period-and-window-ok",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_window": 0,
					"rotation_period": 10,
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           0,
				RotationPeriod:           10 * time.Second,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "period-and-window-ok-strings",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "* */2 * * *",
					"rotation_window":   "5h",
					"rotation_period":   "",
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "* */2 * * *",
				RotationWindow:           5 * time.Hour,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "period-and-schedule-ok",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "",
					"rotation_window":   "",
					"rotation_period":   "2m",
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           0,
				RotationPeriod:           2 * time.Minute,
				DisableAutomatedRotation: false,
			},
		},
		{
			name: "zero-out-schedule-and-window-set-period",
			data: &framework.FieldData{
				Raw: map[string]interface{}{
					"rotation_schedule": "",
					"rotation_window":   "",
					"rotation_period":   "2m",
				},
				Schema: schemaMap,
			},
			expectedParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           0,
				RotationPeriod:           2 * time.Minute,
				DisableAutomatedRotation: false,
			},
			initialParams: &AutomatedRotationParams{
				RotationSchedule:         "*/1 * * * *",
				RotationWindow:           30 * time.Second,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &AutomatedRotationParams{}
			if tt.initialParams != nil {
				p = tt.initialParams
			}
			err := p.ParseAutomatedRotationFields(tt.data)
			if err != nil {
				if tt.expectedError == "" {
					t.Errorf("expected no error but received an error: %s", err)
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("ParseAutomatedRotationFields() error = %v, expected %s", err, tt.expectedError)
				}
			}

			if err == nil && !reflect.DeepEqual(tt.expectedParams, p) {
				t.Errorf("ParseAutomatedRotationFields() error comparing params; got %v, expected %v", p, tt.expectedParams)
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
				"rotation_schedule":          "*/15 * * * *",
				"rotation_window":            time.Duration(60).Seconds(),
				"rotation_period":            time.Duration(0).Seconds(),
				"disable_automated_rotation": false,
			},
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * *",
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

func TestShouldRegisterRotationJob(t *testing.T) {
	tests := []struct {
		name        string
		inputParams *AutomatedRotationParams
		expected    bool
	}{
		{
			name:     "false",
			expected: false,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           60,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name:     "true-schedule",
			expected: true,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * *",
				RotationWindow:           0,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name:     "true-period",
			expected: true,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           0,
				RotationPeriod:           5,
				DisableAutomatedRotation: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.inputParams.ShouldRegisterRotationJob()
			if out != tt.expected {
				t.Errorf("ShouldRegisterRotationJob() output = %v, expected: %v", out, tt.expected)
			}
		})
	}
}

func TestShouldDeregisterRotationJob(t *testing.T) {
	tests := []struct {
		name        string
		inputParams *AutomatedRotationParams
		expected    bool
	}{
		{
			name:     "zero out schedule and period",
			expected: true,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           60,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name:     "false-schedule set",
			expected: false,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * *",
				RotationWindow:           0,
				RotationPeriod:           0,
				DisableAutomatedRotation: false,
			},
		},
		{
			name:     "false-period set",
			expected: false,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "",
				RotationWindow:           0,
				RotationPeriod:           5,
				DisableAutomatedRotation: false,
			},
		},
		{
			name:     "true-disable set",
			expected: true,
			inputParams: &AutomatedRotationParams{
				RotationSchedule:         "*/15 * * * *",
				RotationWindow:           0,
				RotationPeriod:           0,
				DisableAutomatedRotation: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.inputParams.ShouldDeregisterRotationJob()
			if out != tt.expected {
				t.Errorf("ShouldRegisterRotationJob() output = %v, expected: %v", out, tt.expected)
			}
		})
	}
}
