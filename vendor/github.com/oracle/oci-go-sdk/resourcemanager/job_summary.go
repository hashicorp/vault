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

// JobSummary Returns a listing of all of the specified job's properties and their values.
type JobSummary struct {

	// The job's OCID.
	Id *string `mandatory:"false" json:"id"`

	// OCID of the stack that is associated with the specified job.
	StackId *string `mandatory:"false" json:"stackId"`

	// OCID of the compartment where the stack of the associated job resides.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The job's display name.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The type of job executing
	Operation JobOperationEnum `mandatory:"false" json:"operation,omitempty"`

	ApplyJobPlanResolution *ApplyJobPlanResolution `mandatory:"false" json:"applyJobPlanResolution"`

	// The plan job OCID that was used (if this was an APPLY job and not auto approved).
	ResolvedPlanJobId *string `mandatory:"false" json:"resolvedPlanJobId"`

	// The date and time the job was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The date and time the job succeeded or failed.
	TimeFinished *common.SDKTime `mandatory:"false" json:"timeFinished"`

	// Current state of the specified job. Allowed values are:
	// - ACCEPTED
	// - IN_PROGRESS
	// - FAILED
	// - SUCCEEDED
	// - CANCELING
	// - CANCELED
	LifecycleState JobLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// Free-form tags associated with this resource. Each tag is a key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m JobSummary) String() string {
	return common.PointerString(m)
}
