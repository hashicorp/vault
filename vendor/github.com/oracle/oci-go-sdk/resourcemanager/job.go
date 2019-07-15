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

// Job Jobs perform the actions that are defined in your configuration. There are three job types
// - **Plan job**. A plan job takes your Terraform configuration, parses it, and creates an execution plan.
// - **Apply job**. The apply job takes your execution plan, applies it to the associated stack, then executes
// the configuration's instructions.
// - **Destroy job**. To clean up the infrastructure controlled by the stack, you run a destroy job.
// A destroy job does not delete the stack or associated job resources,
// but instead releases the resources managed by the stack.
type Job struct {

	// The job's OCID.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the stack that is associated with the job.
	StackId *string `mandatory:"false" json:"stackId"`

	// The OCID of the compartment in which the job's associated stack resides.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The job's display name.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The type of job executing.
	Operation JobOperationEnum `mandatory:"false" json:"operation,omitempty"`

	ApplyJobPlanResolution *ApplyJobPlanResolution `mandatory:"false" json:"applyJobPlanResolution"`

	// The plan job OCID that was used (if this was an apply job and was not auto-approved).
	ResolvedPlanJobId *string `mandatory:"false" json:"resolvedPlanJobId"`

	// The date and time at which the job was created.
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// The date and time at which the job stopped running, irrespective of whether the job ran successfully.
	TimeFinished *common.SDKTime `mandatory:"false" json:"timeFinished"`

	LifecycleState JobLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	FailureDetails *FailureDetails `mandatory:"false" json:"failureDetails"`

	// The file path to the directory within the configuration from which the job runs.
	WorkingDirectory *string `mandatory:"false" json:"workingDirectory"`

	// Terraform variables associated with this resource.
	// Maximum number of variables supported is 100.
	// The maximum size of each variable, including both name and value, is 4096 bytes.
	// Example: `{"CompartmentId": "compartment-id-value"}`
	Variables map[string]string `mandatory:"false" json:"variables"`

	// Free-form tags associated with this resource. Each tag is a key-value pair with no predefined name, type, or namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Defined tags for this resource. Each key is predefined and scoped to a namespace.
	// For more information, see Resource Tags (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`
}

func (m Job) String() string {
	return common.PointerString(m)
}

// JobOperationEnum Enum with underlying type: string
type JobOperationEnum string

// Set of constants representing the allowable values for JobOperationEnum
const (
	JobOperationPlan    JobOperationEnum = "PLAN"
	JobOperationApply   JobOperationEnum = "APPLY"
	JobOperationDestroy JobOperationEnum = "DESTROY"
)

var mappingJobOperation = map[string]JobOperationEnum{
	"PLAN":    JobOperationPlan,
	"APPLY":   JobOperationApply,
	"DESTROY": JobOperationDestroy,
}

// GetJobOperationEnumValues Enumerates the set of values for JobOperationEnum
func GetJobOperationEnumValues() []JobOperationEnum {
	values := make([]JobOperationEnum, 0)
	for _, v := range mappingJobOperation {
		values = append(values, v)
	}
	return values
}

// JobLifecycleStateEnum Enum with underlying type: string
type JobLifecycleStateEnum string

// Set of constants representing the allowable values for JobLifecycleStateEnum
const (
	JobLifecycleStateAccepted   JobLifecycleStateEnum = "ACCEPTED"
	JobLifecycleStateInProgress JobLifecycleStateEnum = "IN_PROGRESS"
	JobLifecycleStateFailed     JobLifecycleStateEnum = "FAILED"
	JobLifecycleStateSucceeded  JobLifecycleStateEnum = "SUCCEEDED"
	JobLifecycleStateCanceling  JobLifecycleStateEnum = "CANCELING"
	JobLifecycleStateCanceled   JobLifecycleStateEnum = "CANCELED"
)

var mappingJobLifecycleState = map[string]JobLifecycleStateEnum{
	"ACCEPTED":    JobLifecycleStateAccepted,
	"IN_PROGRESS": JobLifecycleStateInProgress,
	"FAILED":      JobLifecycleStateFailed,
	"SUCCEEDED":   JobLifecycleStateSucceeded,
	"CANCELING":   JobLifecycleStateCanceling,
	"CANCELED":    JobLifecycleStateCanceled,
}

// GetJobLifecycleStateEnumValues Enumerates the set of values for JobLifecycleStateEnum
func GetJobLifecycleStateEnumValues() []JobLifecycleStateEnum {
	values := make([]JobLifecycleStateEnum, 0)
	for _, v := range mappingJobLifecycleState {
		values = append(values, v)
	}
	return values
}
