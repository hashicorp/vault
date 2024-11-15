// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfe

// A private struct we need for unmarshalling
type internalWorkspaceRunTask struct {
	ID               string               `jsonapi:"primary,workspace-tasks"`
	EnforcementLevel TaskEnforcementLevel `jsonapi:"attr,enforcement-level"`
	Stage            Stage                `jsonapi:"attr,stage"`
	Stages           []string             `jsonapi:"attr,stages"`

	RunTask   *RunTask   `jsonapi:"relation,task"`
	Workspace *Workspace `jsonapi:"relation,workspace"`
}

// Due to https://github.com/google/jsonapi/issues/74 we must first unmarshall using map[string]interface{}
// and then perform our own conversion for the Stages
func (irt internalWorkspaceRunTask) ToWorkspaceRunTask() *WorkspaceRunTask {
	obj := WorkspaceRunTask{
		ID:               irt.ID,
		EnforcementLevel: irt.EnforcementLevel,
		Stage:            irt.Stage,
		Stages:           make([]Stage, len(irt.Stages)),
		RunTask:          irt.RunTask,
		Workspace:        irt.Workspace,
	}

	for idx, val := range irt.Stages {
		obj.Stages[idx] = Stage(val)
	}

	return &obj
}

// A private struct we need for unmarshalling
type internalWorkspaceRunTaskList struct {
	*Pagination
	Items []*internalWorkspaceRunTask
}

// Due to https://github.com/google/jsonapi/issues/74 we must first unmarshall using
// the internal WorkspaceRunTask struct and convert that a WorkspaceRunTask
func (irt internalWorkspaceRunTaskList) ToWorkspaceRunTaskList() *WorkspaceRunTaskList {
	obj := WorkspaceRunTaskList{
		Pagination: irt.Pagination,
		Items:      make([]*WorkspaceRunTask, len(irt.Items)),
	}

	for idx, src := range irt.Items {
		if src != nil {
			obj.Items[idx] = src.ToWorkspaceRunTask()
		}
	}

	return &obj
}
