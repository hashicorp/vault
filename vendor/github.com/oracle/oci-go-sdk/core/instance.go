// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// Instance A compute host. The image used to launch the instance determines its operating system and other
// software. The shape specified during the launch process determines the number of CPUs and memory
// allocated to the instance. For more information, see
// Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type Instance struct {

	// The availability domain the instance is running in.
	// Example: `Uocm:PHX-AD-1`
	AvailabilityDomain *string `mandatory:"true" json:"availabilityDomain"`

	// The OCID of the compartment that contains the instance.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the instance.
	Id *string `mandatory:"true" json:"id"`

	// The current state of the instance.
	LifecycleState InstanceLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The region that contains the availability domain the instance is running in.
	// For the us-phoenix-1 and us-ashburn-1 regions, `phx` and `iad` are returned, respectively.
	// For all other regions, the full region name is returned.
	// Examples: `phx`, `eu-frankfurt-1`
	Region *string `mandatory:"true" json:"region"`

	// The shape of the instance. The shape determines the number of CPUs and the amount of memory
	// allocated to the instance. You can enumerate all available shapes by calling
	// ListShapes.
	Shape *string `mandatory:"true" json:"shape"`

	// The date and time the instance was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	// Example: `My bare metal instance`
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Additional metadata key/value pairs that you provide. They serve the same purpose and functionality as fields in the 'metadata' object.
	// They are distinguished from 'metadata' fields in that these can be nested JSON objects (whereas 'metadata' fields are string/string maps only).
	ExtendedMetadata map[string]interface{} `mandatory:"false" json:"extendedMetadata"`

	// The name of the fault domain the instance is running in.
	// A fault domain is a grouping of hardware and infrastructure within an availability domain.
	// Each availability domain contains three fault domains. Fault domains let you distribute your
	// instances so that they are not on the same physical hardware within a single availability domain.
	// A hardware failure or Compute hardware maintenance that affects one fault domain does not affect
	// instances in other fault domains.
	// If you do not specify the fault domain, the system selects one for you. To change the fault
	// domain for an instance, terminate it and launch a new instance in the preferred fault domain.
	// Example: `FAULT-DOMAIN-1`
	FaultDomain *string `mandatory:"false" json:"faultDomain"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// Deprecated. Use `sourceDetails` instead.
	ImageId *string `mandatory:"false" json:"imageId"`

	// When a bare metal or virtual machine
	// instance boots, the iPXE firmware that runs on the instance is
	// configured to run an iPXE script to continue the boot process.
	// If you want more control over the boot process, you can provide
	// your own custom iPXE script that will run when the instance boots;
	// however, you should be aware that the same iPXE script will run
	// every time an instance boots; not only after the initial
	// LaunchInstance call.
	// The default iPXE script connects to the instance's local boot
	// volume over iSCSI and performs a network boot. If you use a custom iPXE
	// script and want to network-boot from the instance's local boot volume
	// over iSCSI the same way as the default iPXE script, you should use the
	// following iSCSI IP address: 169.254.0.2, and boot volume IQN:
	// iqn.2015-02.oracle.boot.
	// For more information about the Bring Your Own Image feature of
	// Oracle Cloud Infrastructure, see
	// Bring Your Own Image (https://docs.cloud.oracle.com/Content/Compute/References/bringyourownimage.htm).
	// For more information about iPXE, see http://ipxe.org.
	IpxeScript *string `mandatory:"false" json:"ipxeScript"`

	// Specifies the configuration mode for launching virtual machine (VM) instances. The configuration modes are:
	// * `NATIVE` - VM instances launch with iSCSI boot and VFIO devices. The default value for Oracle-provided images.
	// * `EMULATED` - VM instances launch with emulated devices, such as the E1000 network driver and emulated SCSI disk controller.
	// * `PARAVIRTUALIZED` - VM instances launch with paravirtualized devices using virtio drivers.
	// * `CUSTOM` - VM instances launch with custom configuration settings specified in the `LaunchOptions` parameter.
	LaunchMode InstanceLaunchModeEnum `mandatory:"false" json:"launchMode,omitempty"`

	LaunchOptions *LaunchOptions `mandatory:"false" json:"launchOptions"`

	// Custom metadata that you provide.
	Metadata map[string]string `mandatory:"false" json:"metadata"`

	// Details for creating an instance
	SourceDetails InstanceSourceDetails `mandatory:"false" json:"sourceDetails"`

	AgentConfig *InstanceAgentConfig `mandatory:"false" json:"agentConfig"`

	// The date and time the instance is expected to be stopped / started,  in the format defined by RFC3339.
	// After that time if instance hasn't been rebooted, Oracle will reboot the instance within 24 hours of the due time.
	// Regardless of how the instance was stopped, the flag will be reset to empty as soon as instance reaches Stopped state.
	// Example: `2018-05-25T21:10:29.600Z`
	TimeMaintenanceRebootDue *common.SDKTime `mandatory:"false" json:"timeMaintenanceRebootDue"`
}

func (m Instance) String() string {
	return common.PointerString(m)
}

// UnmarshalJSON unmarshals from json
func (m *Instance) UnmarshalJSON(data []byte) (e error) {
	model := struct {
		DefinedTags              map[string]map[string]interface{} `json:"definedTags"`
		DisplayName              *string                           `json:"displayName"`
		ExtendedMetadata         map[string]interface{}            `json:"extendedMetadata"`
		FaultDomain              *string                           `json:"faultDomain"`
		FreeformTags             map[string]string                 `json:"freeformTags"`
		ImageId                  *string                           `json:"imageId"`
		IpxeScript               *string                           `json:"ipxeScript"`
		LaunchMode               InstanceLaunchModeEnum            `json:"launchMode"`
		LaunchOptions            *LaunchOptions                    `json:"launchOptions"`
		Metadata                 map[string]string                 `json:"metadata"`
		SourceDetails            instancesourcedetails             `json:"sourceDetails"`
		AgentConfig              *InstanceAgentConfig              `json:"agentConfig"`
		TimeMaintenanceRebootDue *common.SDKTime                   `json:"timeMaintenanceRebootDue"`
		AvailabilityDomain       *string                           `json:"availabilityDomain"`
		CompartmentId            *string                           `json:"compartmentId"`
		Id                       *string                           `json:"id"`
		LifecycleState           InstanceLifecycleStateEnum        `json:"lifecycleState"`
		Region                   *string                           `json:"region"`
		Shape                    *string                           `json:"shape"`
		TimeCreated              *common.SDKTime                   `json:"timeCreated"`
	}{}

	e = json.Unmarshal(data, &model)
	if e != nil {
		return
	}
	m.DefinedTags = model.DefinedTags
	m.DisplayName = model.DisplayName
	m.ExtendedMetadata = model.ExtendedMetadata
	m.FaultDomain = model.FaultDomain
	m.FreeformTags = model.FreeformTags
	m.ImageId = model.ImageId
	m.IpxeScript = model.IpxeScript
	m.LaunchMode = model.LaunchMode
	m.LaunchOptions = model.LaunchOptions
	m.Metadata = model.Metadata
	nn, e := model.SourceDetails.UnmarshalPolymorphicJSON(model.SourceDetails.JsonData)
	if e != nil {
		return
	}
	if nn != nil {
		m.SourceDetails = nn.(InstanceSourceDetails)
	} else {
		m.SourceDetails = nil
	}
	m.AgentConfig = model.AgentConfig
	m.TimeMaintenanceRebootDue = model.TimeMaintenanceRebootDue
	m.AvailabilityDomain = model.AvailabilityDomain
	m.CompartmentId = model.CompartmentId
	m.Id = model.Id
	m.LifecycleState = model.LifecycleState
	m.Region = model.Region
	m.Shape = model.Shape
	m.TimeCreated = model.TimeCreated
	return
}

// InstanceLaunchModeEnum Enum with underlying type: string
type InstanceLaunchModeEnum string

// Set of constants representing the allowable values for InstanceLaunchModeEnum
const (
	InstanceLaunchModeNative          InstanceLaunchModeEnum = "NATIVE"
	InstanceLaunchModeEmulated        InstanceLaunchModeEnum = "EMULATED"
	InstanceLaunchModeParavirtualized InstanceLaunchModeEnum = "PARAVIRTUALIZED"
	InstanceLaunchModeCustom          InstanceLaunchModeEnum = "CUSTOM"
)

var mappingInstanceLaunchMode = map[string]InstanceLaunchModeEnum{
	"NATIVE":          InstanceLaunchModeNative,
	"EMULATED":        InstanceLaunchModeEmulated,
	"PARAVIRTUALIZED": InstanceLaunchModeParavirtualized,
	"CUSTOM":          InstanceLaunchModeCustom,
}

// GetInstanceLaunchModeEnumValues Enumerates the set of values for InstanceLaunchModeEnum
func GetInstanceLaunchModeEnumValues() []InstanceLaunchModeEnum {
	values := make([]InstanceLaunchModeEnum, 0)
	for _, v := range mappingInstanceLaunchMode {
		values = append(values, v)
	}
	return values
}

// InstanceLifecycleStateEnum Enum with underlying type: string
type InstanceLifecycleStateEnum string

// Set of constants representing the allowable values for InstanceLifecycleStateEnum
const (
	InstanceLifecycleStateProvisioning  InstanceLifecycleStateEnum = "PROVISIONING"
	InstanceLifecycleStateRunning       InstanceLifecycleStateEnum = "RUNNING"
	InstanceLifecycleStateStarting      InstanceLifecycleStateEnum = "STARTING"
	InstanceLifecycleStateStopping      InstanceLifecycleStateEnum = "STOPPING"
	InstanceLifecycleStateStopped       InstanceLifecycleStateEnum = "STOPPED"
	InstanceLifecycleStateCreatingImage InstanceLifecycleStateEnum = "CREATING_IMAGE"
	InstanceLifecycleStateTerminating   InstanceLifecycleStateEnum = "TERMINATING"
	InstanceLifecycleStateTerminated    InstanceLifecycleStateEnum = "TERMINATED"
)

var mappingInstanceLifecycleState = map[string]InstanceLifecycleStateEnum{
	"PROVISIONING":   InstanceLifecycleStateProvisioning,
	"RUNNING":        InstanceLifecycleStateRunning,
	"STARTING":       InstanceLifecycleStateStarting,
	"STOPPING":       InstanceLifecycleStateStopping,
	"STOPPED":        InstanceLifecycleStateStopped,
	"CREATING_IMAGE": InstanceLifecycleStateCreatingImage,
	"TERMINATING":    InstanceLifecycleStateTerminating,
	"TERMINATED":     InstanceLifecycleStateTerminated,
}

// GetInstanceLifecycleStateEnumValues Enumerates the set of values for InstanceLifecycleStateEnum
func GetInstanceLifecycleStateEnumValues() []InstanceLifecycleStateEnum {
	values := make([]InstanceLifecycleStateEnum, 0)
	for _, v := range mappingInstanceLifecycleState {
		values = append(values, v)
	}
	return values
}
