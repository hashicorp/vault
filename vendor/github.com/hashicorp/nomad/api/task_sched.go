// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

type TaskSchedule struct {
	Cron *TaskScheduleCron `hcl:"cron,block"`
}

type TaskScheduleCron struct {
	Start    string `hcl:"start,optional"`
	End      string `hcl:"end,optional"`
	Timezone string `hcl:"timezone,optional"`
}
