// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

// A private struct we need for unmarshalling
type internalRunTask struct {
	ID          string                 `jsonapi:"primary,tasks"`
	Name        string                 `jsonapi:"attr,name"`
	URL         string                 `jsonapi:"attr,url"`
	Description string                 `jsonapi:"attr,description"`
	Category    string                 `jsonapi:"attr,category"`
	HMACKey     *string                `jsonapi:"attr,hmac-key,omitempty"`
	Enabled     bool                   `jsonapi:"attr,enabled"`
	RawGlobal   map[string]interface{} `jsonapi:"attr,global-configuration,omitempty"`

	Organization      *Organization               `jsonapi:"relation,organization"`
	WorkspaceRunTasks []*internalWorkspaceRunTask `jsonapi:"relation,workspace-tasks"`
}

// Due to https://github.com/google/jsonapi/issues/74 we must first unmarshall using map[string]interface{}
// and then perform our own conversion from the map into a GlobalRunTask struct
func (irt internalRunTask) ToRunTask() *RunTask {
	obj := RunTask{
		ID:          irt.ID,
		Name:        irt.Name,
		URL:         irt.URL,
		Description: irt.Description,
		Category:    irt.Category,
		HMACKey:     irt.HMACKey,
		Enabled:     irt.Enabled,

		Organization: irt.Organization,
	}

	// Convert the WorkspaceRunTasks
	workspaceTasks := make([]*WorkspaceRunTask, len(irt.WorkspaceRunTasks))
	for idx, rawTask := range irt.WorkspaceRunTasks {
		if rawTask != nil {
			workspaceTasks[idx] = rawTask.ToWorkspaceRunTask()
		}
	}
	obj.WorkspaceRunTasks = workspaceTasks

	// Check if the global configuration exists
	if val, ok := irt.RawGlobal["enabled"]; !ok {
		// The enabled property is required so we can assume now that the
		// global configuration was not supplied
		return &obj
	} else if boolVal, ok := val.(bool); !ok {
		// The enabled property exists but it is invalid (Couldn't cast to boolean)
		// so assume the global configuration was not supplied
		return &obj
	} else {
		obj.Global = &GlobalRunTask{
			Enabled: boolVal,
		}
	}

	// Global Enforcement Level
	if val, ok := irt.RawGlobal["enforcement-level"]; ok {
		if stringVal, ok := val.(string); ok {
			obj.Global.EnforcementLevel = TaskEnforcementLevel(stringVal)
		}
	}

	// Global Stages
	if val, ok := irt.RawGlobal["stages"]; ok {
		if stringsVal, ok := val.([]interface{}); ok {
			obj.Global.Stages = make([]Stage, len(stringsVal))
			for idx, stageName := range stringsVal {
				if stringVal, ok := stageName.(string); ok {
					obj.Global.Stages[idx] = Stage(stringVal)
				}
			}
		}
	}

	return &obj
}

// A private struct we need for unmarshalling
type internalRunTaskList struct {
	*Pagination
	Items []*internalRunTask
}

// Due to https://github.com/google/jsonapi/issues/74 we must first unmarshall using
// the internal RunTask struct and convert that a RunTask
func (irt internalRunTaskList) ToRunTaskList() *RunTaskList {
	obj := RunTaskList{
		Pagination: irt.Pagination,
		Items:      make([]*RunTask, len(irt.Items)),
	}

	for idx, src := range irt.Items {
		if src != nil {
			obj.Items[idx] = src.ToRunTask()
		}
	}

	return &obj
}
