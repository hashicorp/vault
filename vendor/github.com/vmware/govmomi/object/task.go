/*
Copyright (c) 2015-2024 VMware, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package object

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/task"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/progress"
	"github.com/vmware/govmomi/vim25/types"
)

// Task is a convenience wrapper around task.Task that keeps a reference to
// the client that was used to create it. This allows users to call the Wait()
// function with only a context parameter, instead of a context parameter, a
// soap.RoundTripper, and reference to the root property collector.
type Task struct {
	Common
}

func NewTask(c *vim25.Client, ref types.ManagedObjectReference) *Task {
	t := Task{
		Common: NewCommon(c, ref),
	}

	return &t
}

// Wait waits for a task to complete.
// NOTE: This method create a thread-safe PropertyCollector instance per-call, so it is thread safe.
// The downside of this approach is the additional resource usage on the vCenter side for each call.
func (t *Task) Wait(ctx context.Context) error {
	_, err := t.WaitForResult(ctx, nil)
	return err
}

// WaitForResult wait for a task to complete.
// NOTE: This method create a thread-safe PropertyCollector instance per-call, so it is thread safe.
// The downside of this approach is the additional resource usage on the vCenter side for each call.
func (t *Task) WaitForResult(ctx context.Context, s ...progress.Sinker) (taskInfo *types.TaskInfo, result error) {
	var pr progress.Sinker
	if len(s) == 1 {
		pr = s[0]
	}
	p, err := property.DefaultCollector(t.c).Create(ctx)
	if err != nil {
		return nil, err
	}

	// Attempt to destroy the collector using the background context, as the
	// specified context may have timed out or have been canceled.
	defer func() {
		if err := p.Destroy(context.Background()); err != nil {
			if result == nil {
				result = err
			} else {
				result = fmt.Errorf(
					"destroy property collector failed with %s after failing to wait for updates: %w",
					err,
					result)
			}
		}
	}()

	return task.WaitEx(ctx, t.Reference(), p, pr)
}

// WaitEx waits for a task to complete.
// NOTE: This method use the same PropertyCollector instance in each call, thus reducing resource usage on the vCenter side.
// The downside of this approach is that this method is not thread safe.
func (t *Task) WaitEx(ctx context.Context) error {
	_, err := t.WaitForResultEx(ctx, nil)
	return err
}

// WaitForResultEx waits for a task to complete.
// NOTE: This method use the same PropertyCollector instance in each call, thus reducing resource usage on the vCenter side.
// The downside of this approach is that this method is not thread safe.
func (t *Task) WaitForResultEx(ctx context.Context, s ...progress.Sinker) (*types.TaskInfo, error) {
	var pr progress.Sinker
	if len(s) == 1 {
		pr = s[0]
	}
	p := property.DefaultCollector(t.c)
	return task.WaitEx(ctx, t.Reference(), p, pr)
}

func (t *Task) Cancel(ctx context.Context) error {
	_, err := methods.CancelTask(ctx, t.Client(), &types.CancelTask{
		This: t.Reference(),
	})

	return err
}

// SetState sets task state and optionally sets results or fault, as appropriate for state.
func (t *Task) SetState(ctx context.Context, state types.TaskInfoState, result types.AnyType, fault *types.LocalizedMethodFault) error {
	req := types.SetTaskState{
		This:   t.Reference(),
		State:  state,
		Result: result,
		Fault:  fault,
	}
	_, err := methods.SetTaskState(ctx, t.Common.Client(), &req)
	return err
}

// SetDescription updates task description to describe the current phase of the task.
func (t *Task) SetDescription(ctx context.Context, description types.LocalizableMessage) error {
	req := types.SetTaskDescription{
		This:        t.Reference(),
		Description: description,
	}
	_, err := methods.SetTaskDescription(ctx, t.Common.Client(), &req)
	return err
}

// UpdateProgress Sets percentage done for this task and recalculates overall percentage done.
// If a percentDone value of less than zero or greater than 100 is specified,
// a value of zero or 100 respectively is used.
func (t *Task) UpdateProgress(ctx context.Context, percentDone int) error {
	req := types.UpdateProgress{
		This:        t.Reference(),
		PercentDone: int32(percentDone),
	}
	_, err := methods.UpdateProgress(ctx, t.Common.Client(), &req)
	return err
}
