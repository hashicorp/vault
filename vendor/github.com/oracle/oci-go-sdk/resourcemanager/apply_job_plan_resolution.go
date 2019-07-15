// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Resource Manager API
//
// API for the Resource Manager service. Use this API to install, configure, and manage resources via the "infrastructure-as-code" model. For more information, see Overview of Resource Manager (https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm).
//

package resourcemanager

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ApplyJobPlanResolution Specifies which plan job provides an execution plan for input to the apply or destroy job.
// You can set only one of the three job properties. For destroy jobs, only `isAutoApproved` is permitted.
type ApplyJobPlanResolution struct {

	// OCID that specifies the most recently executed plan job.
	PlanJobId *string `mandatory:"false" json:"planJobId"`

	// Specifies whether to use the OCID of the most recently run plan job.
	// `True` if using the latest job OCID. Must be a plan job that completed successfully.
	IsUseLatestJobId *bool `mandatory:"false" json:"isUseLatestJobId"`

	// Specifies whether to use the configuration directly, without reference to a Plan job.
	// `True` if using the configuration directly. Note that it is not necessary
	// for a Plan job to have run successfully.
	IsAutoApproved *bool `mandatory:"false" json:"isAutoApproved"`
}

func (m ApplyJobPlanResolution) String() string {
	return common.PointerString(m)
}
