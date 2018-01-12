package api

import (
	"time"

	"github.com/hashicorp/nomad/helper"
	"github.com/hashicorp/nomad/helper/uuid"
)

func MockJob() *Job {
	job := &Job{
		Region:      helper.StringToPtr("global"),
		ID:          helper.StringToPtr(uuid.Generate()),
		Name:        helper.StringToPtr("my-job"),
		Type:        helper.StringToPtr("service"),
		Priority:    helper.IntToPtr(50),
		AllAtOnce:   helper.BoolToPtr(false),
		Datacenters: []string{"dc1"},
		Constraints: []*Constraint{
			{
				LTarget: "${attr.kernel.name}",
				RTarget: "linux",
				Operand: "=",
			},
		},
		TaskGroups: []*TaskGroup{
			{
				Name:  helper.StringToPtr("web"),
				Count: helper.IntToPtr(10),
				EphemeralDisk: &EphemeralDisk{
					SizeMB: helper.IntToPtr(150),
				},
				RestartPolicy: &RestartPolicy{
					Attempts: helper.IntToPtr(3),
					Interval: helper.TimeToPtr(10 * time.Minute),
					Delay:    helper.TimeToPtr(1 * time.Minute),
					Mode:     helper.StringToPtr("delay"),
				},
				Tasks: []*Task{
					{
						Name:   "web",
						Driver: "exec",
						Config: map[string]interface{}{
							"command": "/bin/date",
						},
						Env: map[string]string{
							"FOO": "bar",
						},
						Services: []*Service{
							{
								Name:      "${TASK}-frontend",
								PortLabel: "http",
								Tags:      []string{"pci:${meta.pci-dss}", "datacenter:${node.datacenter}"},
								Checks: []ServiceCheck{
									{
										Name:     "check-table",
										Type:     "script",
										Command:  "/usr/local/check-table-${meta.database}",
										Args:     []string{"${meta.version}"},
										Interval: 30 * time.Second,
										Timeout:  5 * time.Second,
									},
								},
							},
							{
								Name:      "${TASK}-admin",
								PortLabel: "admin",
							},
						},
						LogConfig: DefaultLogConfig(),
						Resources: &Resources{
							CPU:      helper.IntToPtr(500),
							MemoryMB: helper.IntToPtr(256),
							Networks: []*NetworkResource{
								{
									MBits:        helper.IntToPtr(50),
									DynamicPorts: []Port{{Label: "http"}, {Label: "admin"}},
								},
							},
						},
						Meta: map[string]string{
							"foo": "bar",
						},
					},
				},
				Meta: map[string]string{
					"elb_check_type":     "http",
					"elb_check_interval": "30s",
					"elb_check_min":      "3",
				},
			},
		},
		Meta: map[string]string{
			"owner": "armon",
		},
	}
	job.Canonicalize()
	return job
}

func MockPeriodicJob() *Job {
	j := MockJob()
	j.Type = helper.StringToPtr("batch")
	j.Periodic = &PeriodicConfig{
		Enabled:  helper.BoolToPtr(true),
		SpecType: helper.StringToPtr("cron"),
		Spec:     helper.StringToPtr("*/30 * * * *"),
	}
	return j
}
