/*
Copyright (c) 2014-2024 VMware, Inc. All Rights Reserved.

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

package types

import "reflect"

// These constant strings can be used as parameters in user-specified
// email subject and body templates as well as in scripts.
//
// The action processor
// in VirtualCenter substitutes the run-time values for the parameters.
// For example, an email subject provided by the client could be the string:
// `Alarm - {alarmName} Description:\n{eventDescription}`.
// Or a script action provided could be: `myScript {alarmName}`.
type ActionParameter string

const (
	// The name of the entity where the alarm is triggered.
	ActionParameterTargetName = ActionParameter("targetName")
	// The name of the triggering alarm.
	ActionParameterAlarmName = ActionParameter("alarmName")
	// The status prior to the alarm being triggered.
	ActionParameterOldStatus = ActionParameter("oldStatus")
	// The status after the alarm is triggered.
	ActionParameterNewStatus = ActionParameter("newStatus")
	// A summary of information involved in triggering the alarm.
	ActionParameterTriggeringSummary = ActionParameter("triggeringSummary")
	// A summary of declarations made during the triggering of the alarm.
	ActionParameterDeclaringSummary = ActionParameter("declaringSummary")
	// The event description.
	ActionParameterEventDescription = ActionParameter("eventDescription")
	// The object of the entity where the alarm is associated.
	ActionParameterTarget = ActionParameter("target")
	// The object of the triggering alarm.
	ActionParameterAlarm = ActionParameter("alarm")
)

func (e ActionParameter) Values() []ActionParameter {
	return []ActionParameter{
		ActionParameterTargetName,
		ActionParameterAlarmName,
		ActionParameterOldStatus,
		ActionParameterNewStatus,
		ActionParameterTriggeringSummary,
		ActionParameterDeclaringSummary,
		ActionParameterEventDescription,
		ActionParameterTarget,
		ActionParameterAlarm,
	}
}

func (e ActionParameter) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ActionParameter"] = reflect.TypeOf((*ActionParameter)(nil)).Elem()
}

// Pre-defined constants for possible action types.
//
// Virtual Center
// uses this information to coordinate with the clients.
type ActionType string

const (
	// Migration action type
	ActionTypeMigrationV1 = ActionType("MigrationV1")
	// Virtual machine power action type
	ActionTypeVmPowerV1 = ActionType("VmPowerV1")
	// Host power action type
	ActionTypeHostPowerV1 = ActionType("HostPowerV1")
	// Host entering maintenance mode action type
	ActionTypeHostMaintenanceV1 = ActionType("HostMaintenanceV1")
	// Storage migration action type
	ActionTypeStorageMigrationV1 = ActionType("StorageMigrationV1")
	// Initial placement action for a virtual machine or a virtual disk
	ActionTypeStoragePlacementV1 = ActionType("StoragePlacementV1")
	// Initial placement action for a virtual machine and its virtual disks
	ActionTypePlacementV1 = ActionType("PlacementV1")
	// Host changing infrastructure update ha mode action type.
	ActionTypeHostInfraUpdateHaV1 = ActionType("HostInfraUpdateHaV1")
)

func (e ActionType) Values() []ActionType {
	return []ActionType{
		ActionTypeMigrationV1,
		ActionTypeVmPowerV1,
		ActionTypeHostPowerV1,
		ActionTypeHostMaintenanceV1,
		ActionTypeStorageMigrationV1,
		ActionTypeStoragePlacementV1,
		ActionTypePlacementV1,
		ActionTypeHostInfraUpdateHaV1,
	}
}

func (e ActionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ActionType"] = reflect.TypeOf((*ActionType)(nil)).Elem()
}

// Types of affinities.
type AffinityType string

const (
	AffinityTypeMemory = AffinityType("memory")
	AffinityTypeCpu    = AffinityType("cpu")
)

func (e AffinityType) Values() []AffinityType {
	return []AffinityType{
		AffinityTypeMemory,
		AffinityTypeCpu,
	}
}

func (e AffinityType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AffinityType"] = reflect.TypeOf((*AffinityType)(nil)).Elem()
}

type AgentInstallFailedReason string

const (
	// There is not enough storage space on the host to install the agent.
	AgentInstallFailedReasonNotEnoughSpaceOnDevice = AgentInstallFailedReason("NotEnoughSpaceOnDevice")
	// Failed to initialize the upgrade directory on the host.
	AgentInstallFailedReasonPrepareToUpgradeFailed = AgentInstallFailedReason("PrepareToUpgradeFailed")
	// The agent was installed but is not running.
	AgentInstallFailedReasonAgentNotRunning = AgentInstallFailedReason("AgentNotRunning")
	// The agent was installed but did not respond to requests.
	AgentInstallFailedReasonAgentNotReachable = AgentInstallFailedReason("AgentNotReachable")
	// The agent install took too long.
	AgentInstallFailedReasonInstallTimedout = AgentInstallFailedReason("InstallTimedout")
	// The signature verification for the installer failed.
	AgentInstallFailedReasonSignatureVerificationFailed = AgentInstallFailedReason("SignatureVerificationFailed")
	// Failed to upload the agent installer.
	AgentInstallFailedReasonAgentUploadFailed = AgentInstallFailedReason("AgentUploadFailed")
	// The agent upload took too long.
	AgentInstallFailedReasonAgentUploadTimedout = AgentInstallFailedReason("AgentUploadTimedout")
	// The agent installer failed for an unknown reason.
	AgentInstallFailedReasonUnknownInstallerError = AgentInstallFailedReason("UnknownInstallerError")
)

func (e AgentInstallFailedReason) Values() []AgentInstallFailedReason {
	return []AgentInstallFailedReason{
		AgentInstallFailedReasonNotEnoughSpaceOnDevice,
		AgentInstallFailedReasonPrepareToUpgradeFailed,
		AgentInstallFailedReasonAgentNotRunning,
		AgentInstallFailedReasonAgentNotReachable,
		AgentInstallFailedReasonInstallTimedout,
		AgentInstallFailedReasonSignatureVerificationFailed,
		AgentInstallFailedReasonAgentUploadFailed,
		AgentInstallFailedReasonAgentUploadTimedout,
		AgentInstallFailedReasonUnknownInstallerError,
	}
}

func (e AgentInstallFailedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AgentInstallFailedReason"] = reflect.TypeOf((*AgentInstallFailedReason)(nil)).Elem()
}

// Alarm entity type
type AlarmFilterSpecAlarmTypeByEntity string

const (
	// Alarms on all entity types.
	AlarmFilterSpecAlarmTypeByEntityEntityTypeAll = AlarmFilterSpecAlarmTypeByEntity("entityTypeAll")
	// Host alarms
	AlarmFilterSpecAlarmTypeByEntityEntityTypeHost = AlarmFilterSpecAlarmTypeByEntity("entityTypeHost")
	// VM alarms
	AlarmFilterSpecAlarmTypeByEntityEntityTypeVm = AlarmFilterSpecAlarmTypeByEntity("entityTypeVm")
)

func (e AlarmFilterSpecAlarmTypeByEntity) Values() []AlarmFilterSpecAlarmTypeByEntity {
	return []AlarmFilterSpecAlarmTypeByEntity{
		AlarmFilterSpecAlarmTypeByEntityEntityTypeAll,
		AlarmFilterSpecAlarmTypeByEntityEntityTypeHost,
		AlarmFilterSpecAlarmTypeByEntityEntityTypeVm,
	}
}

func (e AlarmFilterSpecAlarmTypeByEntity) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AlarmFilterSpecAlarmTypeByEntity"] = reflect.TypeOf((*AlarmFilterSpecAlarmTypeByEntity)(nil)).Elem()
}

// Alarm triggering type.
//
// The main divisions are event triggered and
// metric- or state-based alarms.
type AlarmFilterSpecAlarmTypeByTrigger string

const (
	// All alarm types.
	AlarmFilterSpecAlarmTypeByTriggerTriggerTypeAll = AlarmFilterSpecAlarmTypeByTrigger("triggerTypeAll")
	// Event based alarms
	AlarmFilterSpecAlarmTypeByTriggerTriggerTypeEvent = AlarmFilterSpecAlarmTypeByTrigger("triggerTypeEvent")
	// Metric or state alarms
	AlarmFilterSpecAlarmTypeByTriggerTriggerTypeMetric = AlarmFilterSpecAlarmTypeByTrigger("triggerTypeMetric")
)

func (e AlarmFilterSpecAlarmTypeByTrigger) Values() []AlarmFilterSpecAlarmTypeByTrigger {
	return []AlarmFilterSpecAlarmTypeByTrigger{
		AlarmFilterSpecAlarmTypeByTriggerTriggerTypeAll,
		AlarmFilterSpecAlarmTypeByTriggerTriggerTypeEvent,
		AlarmFilterSpecAlarmTypeByTriggerTriggerTypeMetric,
	}
}

func (e AlarmFilterSpecAlarmTypeByTrigger) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AlarmFilterSpecAlarmTypeByTrigger"] = reflect.TypeOf((*AlarmFilterSpecAlarmTypeByTrigger)(nil)).Elem()
}

// Defines the result status values for a validating answer file.
type AnswerFileValidationInfoStatus string

const (
	// Answer File validation was successful.
	AnswerFileValidationInfoStatusSuccess = AnswerFileValidationInfoStatus("success")
	// Answer File validation failed.
	AnswerFileValidationInfoStatusFailed = AnswerFileValidationInfoStatus("failed")
	// Answer File validation failed to generate default.
	AnswerFileValidationInfoStatusFailed_defaults = AnswerFileValidationInfoStatus("failed_defaults")
)

func (e AnswerFileValidationInfoStatus) Values() []AnswerFileValidationInfoStatus {
	return []AnswerFileValidationInfoStatus{
		AnswerFileValidationInfoStatusSuccess,
		AnswerFileValidationInfoStatusFailed,
		AnswerFileValidationInfoStatusFailed_defaults,
	}
}

func (e AnswerFileValidationInfoStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AnswerFileValidationInfoStatus"] = reflect.TypeOf((*AnswerFileValidationInfoStatus)(nil)).Elem()
}

type ApplyHostProfileConfigurationResultStatus string

const (
	// Remediation succeeded.
	ApplyHostProfileConfigurationResultStatusSuccess = ApplyHostProfileConfigurationResultStatus("success")
	// Remediation failed.
	ApplyHostProfileConfigurationResultStatusFailed = ApplyHostProfileConfigurationResultStatus("failed")
	// Remediation succeeded but reboot after remediation failed.
	//
	// May treat this as a warning.
	ApplyHostProfileConfigurationResultStatusReboot_failed = ApplyHostProfileConfigurationResultStatus("reboot_failed")
	// Stateless reboot for remediation failed.
	ApplyHostProfileConfigurationResultStatusStateless_reboot_failed = ApplyHostProfileConfigurationResultStatus("stateless_reboot_failed")
	// Remediation and reboot succeeded but check compliance after reboot
	// failed.
	//
	// May treat this as a warning.
	ApplyHostProfileConfigurationResultStatusCheck_compliance_failed = ApplyHostProfileConfigurationResultStatus("check_compliance_failed")
	// The required state is not satisfied so host profiel apply cannot
	// be done.
	ApplyHostProfileConfigurationResultStatusState_not_satisfied = ApplyHostProfileConfigurationResultStatus("state_not_satisfied")
	// Exit maintenance mode failed.
	ApplyHostProfileConfigurationResultStatusExit_maintenancemode_failed = ApplyHostProfileConfigurationResultStatus("exit_maintenancemode_failed")
	// The remediation was canceled.
	ApplyHostProfileConfigurationResultStatusCanceled = ApplyHostProfileConfigurationResultStatus("canceled")
)

func (e ApplyHostProfileConfigurationResultStatus) Values() []ApplyHostProfileConfigurationResultStatus {
	return []ApplyHostProfileConfigurationResultStatus{
		ApplyHostProfileConfigurationResultStatusSuccess,
		ApplyHostProfileConfigurationResultStatusFailed,
		ApplyHostProfileConfigurationResultStatusReboot_failed,
		ApplyHostProfileConfigurationResultStatusStateless_reboot_failed,
		ApplyHostProfileConfigurationResultStatusCheck_compliance_failed,
		ApplyHostProfileConfigurationResultStatusState_not_satisfied,
		ApplyHostProfileConfigurationResultStatusExit_maintenancemode_failed,
		ApplyHostProfileConfigurationResultStatusCanceled,
	}
}

func (e ApplyHostProfileConfigurationResultStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ApplyHostProfileConfigurationResultStatus"] = reflect.TypeOf((*ApplyHostProfileConfigurationResultStatus)(nil)).Elem()
}

// This list specifies the type of operation being performed on the array.
type ArrayUpdateOperation string

const (
	// indicates an addition to the array.
	ArrayUpdateOperationAdd = ArrayUpdateOperation("add")
	// indicates the removal of an element in the
	// array.
	//
	// In this case the key field must contain the key of the element
	// to be removed.
	ArrayUpdateOperationRemove = ArrayUpdateOperation("remove")
	// indicates changes to an element in the array.
	ArrayUpdateOperationEdit = ArrayUpdateOperation("edit")
)

func (e ArrayUpdateOperation) Values() []ArrayUpdateOperation {
	return []ArrayUpdateOperation{
		ArrayUpdateOperationAdd,
		ArrayUpdateOperationRemove,
		ArrayUpdateOperationEdit,
	}
}

func (e ArrayUpdateOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ArrayUpdateOperation"] = reflect.TypeOf((*ArrayUpdateOperation)(nil)).Elem()
}

type AutoStartAction string

const (
	// No action is taken for this virtual machine.
	//
	// This virtual machine is
	// not a part of the auto-start sequence. This can be used for both auto-start
	// and auto-start settings.
	AutoStartActionNone = AutoStartAction("none")
	// The default system action is taken for this virtual machine when it is next in
	// the auto-start order.
	//
	// This can be used for both auto-start and auto-start
	// settings.
	AutoStartActionSystemDefault = AutoStartAction("systemDefault")
	// This virtual machine is powered on when it is next in the auto-start order.
	AutoStartActionPowerOn = AutoStartAction("powerOn")
	// This virtual machine is powered off when it is next in the auto-stop order.
	//
	// This is the default stopAction.
	AutoStartActionPowerOff = AutoStartAction("powerOff")
	// The guest operating system for a virtual machine is shut down when that
	// virtual machine in next in the auto-stop order.
	AutoStartActionGuestShutdown = AutoStartAction("guestShutdown")
	// This virtual machine is suspended when it is next in the auto-stop order.
	AutoStartActionSuspend = AutoStartAction("suspend")
)

func (e AutoStartAction) Values() []AutoStartAction {
	return []AutoStartAction{
		AutoStartActionNone,
		AutoStartActionSystemDefault,
		AutoStartActionPowerOn,
		AutoStartActionPowerOff,
		AutoStartActionGuestShutdown,
		AutoStartActionSuspend,
	}
}

func (e AutoStartAction) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AutoStartAction"] = reflect.TypeOf((*AutoStartAction)(nil)).Elem()
}

// Determines if the virtual machine should start after receiving a heartbeat,
// ignore heartbeats and start after the startDelay has elapsed, or follow the
// system default before powering on.
//
// When a virtual machine is next in the start
// order, the system either waits a specified period of time for a virtual
// machine to power on or it waits until it receives a successful heartbeat from a
// powered on virtual machine. By default, this is set to no.
type AutoStartWaitHeartbeatSetting string

const (
	// The system waits until receiving a heartbeat before powering on the next
	// machine in the order.
	AutoStartWaitHeartbeatSettingYes = AutoStartWaitHeartbeatSetting("yes")
	// The system does not wait to receive a heartbeat before powering on the next
	// machine in the order.
	//
	// This is the default setting.
	AutoStartWaitHeartbeatSettingNo = AutoStartWaitHeartbeatSetting("no")
	// The system uses the default value to determine whether or not to wait to
	// receive a heartbeat before powering on the next machine in the order.
	AutoStartWaitHeartbeatSettingSystemDefault = AutoStartWaitHeartbeatSetting("systemDefault")
)

func (e AutoStartWaitHeartbeatSetting) Values() []AutoStartWaitHeartbeatSetting {
	return []AutoStartWaitHeartbeatSetting{
		AutoStartWaitHeartbeatSettingYes,
		AutoStartWaitHeartbeatSettingNo,
		AutoStartWaitHeartbeatSettingSystemDefault,
	}
}

func (e AutoStartWaitHeartbeatSetting) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["AutoStartWaitHeartbeatSetting"] = reflect.TypeOf((*AutoStartWaitHeartbeatSetting)(nil)).Elem()
}

// Provisioning type constants.
type BaseConfigInfoDiskFileBackingInfoProvisioningType string

const (
	// Space required for thin-provisioned virtual disk is allocated
	// and zeroed on demand as the space is used.
	BaseConfigInfoDiskFileBackingInfoProvisioningTypeThin = BaseConfigInfoDiskFileBackingInfoProvisioningType("thin")
	// An eager zeroed thick virtual disk has all space allocated and
	// wiped clean of any previous contents on the physical media at
	// creation time.
	//
	// Such virtual disk may take longer time
	// during creation compared to other provisioning formats.
	BaseConfigInfoDiskFileBackingInfoProvisioningTypeEagerZeroedThick = BaseConfigInfoDiskFileBackingInfoProvisioningType("eagerZeroedThick")
	// A thick virtual disk has all space allocated at creation time.
	//
	// This space may contain stale data on the physical media.
	BaseConfigInfoDiskFileBackingInfoProvisioningTypeLazyZeroedThick = BaseConfigInfoDiskFileBackingInfoProvisioningType("lazyZeroedThick")
)

func (e BaseConfigInfoDiskFileBackingInfoProvisioningType) Values() []BaseConfigInfoDiskFileBackingInfoProvisioningType {
	return []BaseConfigInfoDiskFileBackingInfoProvisioningType{
		BaseConfigInfoDiskFileBackingInfoProvisioningTypeThin,
		BaseConfigInfoDiskFileBackingInfoProvisioningTypeEagerZeroedThick,
		BaseConfigInfoDiskFileBackingInfoProvisioningTypeLazyZeroedThick,
	}
}

func (e BaseConfigInfoDiskFileBackingInfoProvisioningType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["BaseConfigInfoDiskFileBackingInfoProvisioningType"] = reflect.TypeOf((*BaseConfigInfoDiskFileBackingInfoProvisioningType)(nil)).Elem()
}

// Enum representing result of batch-APis.
type BatchResultResult string

const (
	BatchResultResultSuccess = BatchResultResult("success")
	BatchResultResultFail    = BatchResultResult("fail")
)

func (e BatchResultResult) Values() []BatchResultResult {
	return []BatchResultResult{
		BatchResultResultSuccess,
		BatchResultResultFail,
	}
}

func (e BatchResultResult) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["BatchResultResult"] = reflect.TypeOf((*BatchResultResult)(nil)).Elem()
}

type CannotEnableVmcpForClusterReason string

const (
	// APD timeout has been disabled on one of the host
	CannotEnableVmcpForClusterReasonAPDTimeoutDisabled = CannotEnableVmcpForClusterReason("APDTimeoutDisabled")
)

func (e CannotEnableVmcpForClusterReason) Values() []CannotEnableVmcpForClusterReason {
	return []CannotEnableVmcpForClusterReason{
		CannotEnableVmcpForClusterReasonAPDTimeoutDisabled,
	}
}

func (e CannotEnableVmcpForClusterReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CannotEnableVmcpForClusterReason"] = reflect.TypeOf((*CannotEnableVmcpForClusterReason)(nil)).Elem()
}

type CannotMoveFaultToleranceVmMoveType string

const (
	// Move out of the resouce pool
	CannotMoveFaultToleranceVmMoveTypeResourcePool = CannotMoveFaultToleranceVmMoveType("resourcePool")
	// Move out of the cluster
	CannotMoveFaultToleranceVmMoveTypeCluster = CannotMoveFaultToleranceVmMoveType("cluster")
)

func (e CannotMoveFaultToleranceVmMoveType) Values() []CannotMoveFaultToleranceVmMoveType {
	return []CannotMoveFaultToleranceVmMoveType{
		CannotMoveFaultToleranceVmMoveTypeResourcePool,
		CannotMoveFaultToleranceVmMoveTypeCluster,
	}
}

func (e CannotMoveFaultToleranceVmMoveType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CannotMoveFaultToleranceVmMoveType"] = reflect.TypeOf((*CannotMoveFaultToleranceVmMoveType)(nil)).Elem()
}

type CannotPowerOffVmInClusterOperation string

const (
	// suspend
	CannotPowerOffVmInClusterOperationSuspend = CannotPowerOffVmInClusterOperation("suspend")
	// power off
	CannotPowerOffVmInClusterOperationPowerOff = CannotPowerOffVmInClusterOperation("powerOff")
	// guest shutdown
	CannotPowerOffVmInClusterOperationGuestShutdown = CannotPowerOffVmInClusterOperation("guestShutdown")
	// guest suspend
	CannotPowerOffVmInClusterOperationGuestSuspend = CannotPowerOffVmInClusterOperation("guestSuspend")
)

func (e CannotPowerOffVmInClusterOperation) Values() []CannotPowerOffVmInClusterOperation {
	return []CannotPowerOffVmInClusterOperation{
		CannotPowerOffVmInClusterOperationSuspend,
		CannotPowerOffVmInClusterOperationPowerOff,
		CannotPowerOffVmInClusterOperationGuestShutdown,
		CannotPowerOffVmInClusterOperationGuestSuspend,
	}
}

func (e CannotPowerOffVmInClusterOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CannotPowerOffVmInClusterOperation"] = reflect.TypeOf((*CannotPowerOffVmInClusterOperation)(nil)).Elem()
}

type CannotUseNetworkReason string

const (
	// Network does not support reservation
	CannotUseNetworkReasonNetworkReservationNotSupported = CannotUseNetworkReason("NetworkReservationNotSupported")
	// Source and destination networks do not have same security policies
	CannotUseNetworkReasonMismatchedNetworkPolicies = CannotUseNetworkReason("MismatchedNetworkPolicies")
	// Source and destination DVS do not have same version or vendor
	CannotUseNetworkReasonMismatchedDvsVersionOrVendor = CannotUseNetworkReason("MismatchedDvsVersionOrVendor")
	// VMotion to unsupported destination network type
	CannotUseNetworkReasonVMotionToUnsupportedNetworkType = CannotUseNetworkReason("VMotionToUnsupportedNetworkType")
	// The network is under maintenance
	CannotUseNetworkReasonNetworkUnderMaintenance = CannotUseNetworkReason("NetworkUnderMaintenance")
	// Source and destination networks do not have same ENS(Enhanced Network Stack) mode
	CannotUseNetworkReasonMismatchedEnsMode = CannotUseNetworkReason("MismatchedEnsMode")
)

func (e CannotUseNetworkReason) Values() []CannotUseNetworkReason {
	return []CannotUseNetworkReason{
		CannotUseNetworkReasonNetworkReservationNotSupported,
		CannotUseNetworkReasonMismatchedNetworkPolicies,
		CannotUseNetworkReasonMismatchedDvsVersionOrVendor,
		CannotUseNetworkReasonVMotionToUnsupportedNetworkType,
		CannotUseNetworkReasonNetworkUnderMaintenance,
		CannotUseNetworkReasonMismatchedEnsMode,
	}
}

func (e CannotUseNetworkReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CannotUseNetworkReason"] = reflect.TypeOf((*CannotUseNetworkReason)(nil)).Elem()
}

// The types of tests which can requested by any of the methods in either
// `VirtualMachineCompatibilityChecker` or `VirtualMachineProvisioningChecker`.
type CheckTestType string

const (
	// Tests that examine only the configuration
	// of the virtual machine and its current host; the destination
	// resource pool and host or cluster are irrelevant.
	CheckTestTypeSourceTests = CheckTestType("sourceTests")
	// Tests that examine both the virtual
	// machine and the destination host or cluster; the destination
	// resource pool is irrelevant.
	//
	// This set excludes tests that fall
	// into the datastoreTests group.
	CheckTestTypeHostTests = CheckTestType("hostTests")
	// Tests that check that the destination resource
	// pool can support the virtual machine if it is powered on.
	//
	// The
	// destination host or cluster is relevant because it will affect the
	// amount of overhead memory required to run the virtual machine.
	CheckTestTypeResourcePoolTests = CheckTestType("resourcePoolTests")
	// Tests that check that the
	// destination host or cluster can see the datastores where the virtual
	// machine's virtual disks are going to be located.
	//
	// The destination
	// resource pool is irrelevant.
	CheckTestTypeDatastoreTests = CheckTestType("datastoreTests")
	// Tests that check that the
	// destination host or cluster can see the networks that the virtual
	// machine's virtual nic devices are going to be connected.
	CheckTestTypeNetworkTests = CheckTestType("networkTests")
)

func (e CheckTestType) Values() []CheckTestType {
	return []CheckTestType{
		CheckTestTypeSourceTests,
		CheckTestTypeHostTests,
		CheckTestTypeResourcePoolTests,
		CheckTestTypeDatastoreTests,
		CheckTestTypeNetworkTests,
	}
}

func (e CheckTestType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CheckTestType"] = reflect.TypeOf((*CheckTestType)(nil)).Elem()
}

// HCIWorkflowState identifies the state of the cluser from the perspective of HCI
// workflow.
//
// The workflow begins with in\_progress mode and can transition
// to 'done' or 'invalid', both of which are terminal states.
type ClusterComputeResourceHCIWorkflowState string

const (
	// Indicates cluster is getting configured or will be configured.
	ClusterComputeResourceHCIWorkflowStateIn_progress = ClusterComputeResourceHCIWorkflowState("in_progress")
	// Indicates cluster configuration is complete.
	ClusterComputeResourceHCIWorkflowStateDone = ClusterComputeResourceHCIWorkflowState("done")
	// Indicates the workflow was abandoned on the cluster before the
	// configuration could complete.
	ClusterComputeResourceHCIWorkflowStateInvalid = ClusterComputeResourceHCIWorkflowState("invalid")
)

func (e ClusterComputeResourceHCIWorkflowState) Values() []ClusterComputeResourceHCIWorkflowState {
	return []ClusterComputeResourceHCIWorkflowState{
		ClusterComputeResourceHCIWorkflowStateIn_progress,
		ClusterComputeResourceHCIWorkflowStateDone,
		ClusterComputeResourceHCIWorkflowStateInvalid,
	}
}

func (e ClusterComputeResourceHCIWorkflowState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterComputeResourceHCIWorkflowState"] = reflect.TypeOf((*ClusterComputeResourceHCIWorkflowState)(nil)).Elem()
}

type ClusterComputeResourceVcsHealthStatus string

const (
	// Indicates vCS health status is normal.
	ClusterComputeResourceVcsHealthStatusHealthy = ClusterComputeResourceVcsHealthStatus("healthy")
	// Indicates only vCS is unhealthy.
	ClusterComputeResourceVcsHealthStatusDegraded = ClusterComputeResourceVcsHealthStatus("degraded")
	// Indicates vCS is unhealthy and other cluster services are impacted.
	ClusterComputeResourceVcsHealthStatusNonhealthy = ClusterComputeResourceVcsHealthStatus("nonhealthy")
)

func (e ClusterComputeResourceVcsHealthStatus) Values() []ClusterComputeResourceVcsHealthStatus {
	return []ClusterComputeResourceVcsHealthStatus{
		ClusterComputeResourceVcsHealthStatusHealthy,
		ClusterComputeResourceVcsHealthStatusDegraded,
		ClusterComputeResourceVcsHealthStatusNonhealthy,
	}
}

func (e ClusterComputeResourceVcsHealthStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterComputeResourceVcsHealthStatus"] = reflect.TypeOf((*ClusterComputeResourceVcsHealthStatus)(nil)).Elem()
	minAPIVersionForType["ClusterComputeResourceVcsHealthStatus"] = "7.0.1.1"
}

type ClusterCryptoConfigInfoCryptoMode string

const (
	// Put each host into the crypto safe state automatically when needed.
	ClusterCryptoConfigInfoCryptoModeOnDemand = ClusterCryptoConfigInfoCryptoMode("onDemand")
	// Put each host into the crypto safe state immediately.
	ClusterCryptoConfigInfoCryptoModeForceEnable = ClusterCryptoConfigInfoCryptoMode("forceEnable")
)

func (e ClusterCryptoConfigInfoCryptoMode) Values() []ClusterCryptoConfigInfoCryptoMode {
	return []ClusterCryptoConfigInfoCryptoMode{
		ClusterCryptoConfigInfoCryptoModeOnDemand,
		ClusterCryptoConfigInfoCryptoModeForceEnable,
	}
}

func (e ClusterCryptoConfigInfoCryptoMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterCryptoConfigInfoCryptoMode"] = reflect.TypeOf((*ClusterCryptoConfigInfoCryptoMode)(nil)).Elem()
}

// The `ClusterDasAamNodeStateDasState_enum` enumerated type defines
// values for host HA configuration and runtime state properties
// (`ClusterDasAamNodeState.configState` and
// `ClusterDasAamNodeState.runtimeState`).
type ClusterDasAamNodeStateDasState string

const (
	// HA has never been enabled on the the host.
	ClusterDasAamNodeStateDasStateUninitialized = ClusterDasAamNodeStateDasState("uninitialized")
	// HA agents have been installed but are not running on the the host.
	ClusterDasAamNodeStateDasStateInitialized = ClusterDasAamNodeStateDasState("initialized")
	// HA configuration is in progress.
	ClusterDasAamNodeStateDasStateConfiguring = ClusterDasAamNodeStateDasState("configuring")
	// HA configuration is being removed.
	ClusterDasAamNodeStateDasStateUnconfiguring = ClusterDasAamNodeStateDasState("unconfiguring")
	// HA agent is running on this host.
	ClusterDasAamNodeStateDasStateRunning = ClusterDasAamNodeStateDasState("running")
	// There is an error condition.
	//
	// This can represent a configuration
	// error or a host agent runtime error.
	ClusterDasAamNodeStateDasStateError = ClusterDasAamNodeStateDasState("error")
	// The HA agent has been shut down.
	ClusterDasAamNodeStateDasStateAgentShutdown = ClusterDasAamNodeStateDasState("agentShutdown")
	// The host is not reachable.
	//
	// This can represent a host failure
	// or an isolated host.
	ClusterDasAamNodeStateDasStateNodeFailed = ClusterDasAamNodeStateDasState("nodeFailed")
)

func (e ClusterDasAamNodeStateDasState) Values() []ClusterDasAamNodeStateDasState {
	return []ClusterDasAamNodeStateDasState{
		ClusterDasAamNodeStateDasStateUninitialized,
		ClusterDasAamNodeStateDasStateInitialized,
		ClusterDasAamNodeStateDasStateConfiguring,
		ClusterDasAamNodeStateDasStateUnconfiguring,
		ClusterDasAamNodeStateDasStateRunning,
		ClusterDasAamNodeStateDasStateError,
		ClusterDasAamNodeStateDasStateAgentShutdown,
		ClusterDasAamNodeStateDasStateNodeFailed,
	}
}

func (e ClusterDasAamNodeStateDasState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasAamNodeStateDasState"] = reflect.TypeOf((*ClusterDasAamNodeStateDasState)(nil)).Elem()
}

// The policy to determine the candidates from which vCenter Server can
// choose heartbeat datastores.
type ClusterDasConfigInfoHBDatastoreCandidate string

const (
	// vCenter Server chooses heartbeat datastores from the set specified
	// by the user (see `ClusterDasConfigInfo.heartbeatDatastore`).
	//
	// More specifically,
	// datastores not included in the set will not be chosen. Note that if
	// `ClusterDasConfigInfo.heartbeatDatastore` is empty, datastore heartbeating will
	// be disabled for HA.
	ClusterDasConfigInfoHBDatastoreCandidateUserSelectedDs = ClusterDasConfigInfoHBDatastoreCandidate("userSelectedDs")
	// vCenter Server chooses heartbeat datastores from all the feasible ones,
	// i.e., the datastores that are accessible to more than one host in
	// the cluster.
	//
	// The choice will be made without giving preference to those
	// specified by the user (see `ClusterDasConfigInfo.heartbeatDatastore`).
	ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDs = ClusterDasConfigInfoHBDatastoreCandidate("allFeasibleDs")
	// vCenter Server chooses heartbeat datastores from all the feasible ones
	// while giving preference to those specified by the user (see `ClusterDasConfigInfo.heartbeatDatastore`).
	//
	// More specifically, the datastores not included in `ClusterDasConfigInfo.heartbeatDatastore` will be
	// chosen if and only if the specified ones are not sufficient.
	ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDsWithUserPreference = ClusterDasConfigInfoHBDatastoreCandidate("allFeasibleDsWithUserPreference")
)

func (e ClusterDasConfigInfoHBDatastoreCandidate) Values() []ClusterDasConfigInfoHBDatastoreCandidate {
	return []ClusterDasConfigInfoHBDatastoreCandidate{
		ClusterDasConfigInfoHBDatastoreCandidateUserSelectedDs,
		ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDs,
		ClusterDasConfigInfoHBDatastoreCandidateAllFeasibleDsWithUserPreference,
	}
}

func (e ClusterDasConfigInfoHBDatastoreCandidate) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasConfigInfoHBDatastoreCandidate"] = reflect.TypeOf((*ClusterDasConfigInfoHBDatastoreCandidate)(nil)).Elem()
}

// Possible states of an HA service.
//
// All services support the
// disabled and enabled states.
type ClusterDasConfigInfoServiceState string

const (
	// HA service is disabled.
	ClusterDasConfigInfoServiceStateDisabled = ClusterDasConfigInfoServiceState("disabled")
	// HA service is enabled.
	ClusterDasConfigInfoServiceStateEnabled = ClusterDasConfigInfoServiceState("enabled")
)

func (e ClusterDasConfigInfoServiceState) Values() []ClusterDasConfigInfoServiceState {
	return []ClusterDasConfigInfoServiceState{
		ClusterDasConfigInfoServiceStateDisabled,
		ClusterDasConfigInfoServiceStateEnabled,
	}
}

func (e ClusterDasConfigInfoServiceState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasConfigInfoServiceState"] = reflect.TypeOf((*ClusterDasConfigInfoServiceState)(nil)).Elem()
}

// The `ClusterDasConfigInfoVmMonitoringState_enum` enum defines values that indicate
// the state of Virtual Machine Health Monitoring.
//
// Health Monitoring
// uses the vmTools (guest) and application agent heartbeat modules.
// You can configure HA to respond to heartbeat failures of either one
// or both modules. You can also disable the HA response to heartbeat failures.
//   - To set the cluster default for health monitoring, use the
//     ClusterConfigSpecEx.dasConfig.`ClusterDasConfigInfo.vmMonitoring` property.
//   - To set health monitoring for a virtual machine, use the
//     ClusterConfigSpecEx.dasVmConfigSpec.info.dasSettings.`ClusterDasVmSettings.vmToolsMonitoringSettings` property.
//   - To retrieve the current state of health monitoring (cluster setting), use the
//     ClusterConfigInfoEx.dasConfig.`ClusterDasConfigInfo.vmMonitoring`
//     property.
//   - To retrieve the current state of health monitoring for a virtual machine, use the
//     ClusterConfigInfoEx.dasVmConfig\[\].dasSettings.vmToolsMonitoringSettings.`ClusterVmToolsMonitoringSettings.vmMonitoring`
//     property.
type ClusterDasConfigInfoVmMonitoringState string

const (
	// Virtual machine health monitoring is disabled.
	//
	// In this state,
	// HA response to guest and application heartbeat failures are disabled.
	ClusterDasConfigInfoVmMonitoringStateVmMonitoringDisabled = ClusterDasConfigInfoVmMonitoringState("vmMonitoringDisabled")
	// HA response to guest heartbeat failure is enabled.
	//
	// To retrieve the guest heartbeat status, use the
	// `VirtualMachine*.*VirtualMachine.guestHeartbeatStatus`
	// property.
	ClusterDasConfigInfoVmMonitoringStateVmMonitoringOnly = ClusterDasConfigInfoVmMonitoringState("vmMonitoringOnly")
	// HA response to both guest and application heartbeat failure is enabled.
	//   - To retrieve the guest heartbeat status, use the
	//     `VirtualMachine*.*VirtualMachine.guestHeartbeatStatus`
	//     property.
	//   - To retrieve the application heartbeat status, use the
	//     `GuestInfo*.*GuestInfo.appHeartbeatStatus`
	//     property.
	ClusterDasConfigInfoVmMonitoringStateVmAndAppMonitoring = ClusterDasConfigInfoVmMonitoringState("vmAndAppMonitoring")
)

func (e ClusterDasConfigInfoVmMonitoringState) Values() []ClusterDasConfigInfoVmMonitoringState {
	return []ClusterDasConfigInfoVmMonitoringState{
		ClusterDasConfigInfoVmMonitoringStateVmMonitoringDisabled,
		ClusterDasConfigInfoVmMonitoringStateVmMonitoringOnly,
		ClusterDasConfigInfoVmMonitoringStateVmAndAppMonitoring,
	}
}

func (e ClusterDasConfigInfoVmMonitoringState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasConfigInfoVmMonitoringState"] = reflect.TypeOf((*ClusterDasConfigInfoVmMonitoringState)(nil)).Elem()
}

// The `ClusterDasFdmAvailabilityState_enum` enumeration describes the
// availability states of hosts in a vSphere HA cluster.
//
// In the HA
// architecture, a agent called the Fault Domain Manager runs on
// each active host. These agents elect a master and the others become
// its slaves. The availability state assigned to a given host is
// determined from information reported by the Fault Domain Manager
// running on the host, by a Fault Domain Manager that has been elected
// master, and by vCenter Server. See `ClusterDasFdmHostState`
// for more information about the vSphere HA architecture.
type ClusterDasFdmAvailabilityState string

const (
	// The Fault Domain Manager for the host has not yet been
	// initialized.
	//
	// Hence the host is not part of a vSphere HA
	// fault domain. This state is reported by vCenter Server or
	// by the host itself.
	ClusterDasFdmAvailabilityStateUninitialized = ClusterDasFdmAvailabilityState("uninitialized")
	// The Fault Domain Manager on the host has been initialized and
	// the host is either waiting to join the existing master or
	// is participating in an election for a new master.
	//
	// This state
	// is reported by vCenter Server or by the host itself.
	ClusterDasFdmAvailabilityStateElection = ClusterDasFdmAvailabilityState("election")
	// The Fault Domain Manager on the host has been elected a
	// master.
	//
	// This state is reported by the the host itself.
	ClusterDasFdmAvailabilityStateMaster = ClusterDasFdmAvailabilityState("master")
	// The normal operating state for a slave host.
	//
	// In this state,
	// the host is exchanging heartbeats with a master over
	// the management network, and is thus connected to it. If
	// there is a management network partition, the slave will be
	// in this state only if it is in the same partition as the master.
	// This state is reported by the master of a slave host.
	ClusterDasFdmAvailabilityStateConnectedToMaster = ClusterDasFdmAvailabilityState("connectedToMaster")
	// A slave host is alive and has management network connectivity, but
	// the management network has been partitioned.
	//
	// This state is reported
	// by masters that are in a partition other than the one containing the
	// slave host; the master in the slave's partition will report the slave state
	// as `connectedToMaster`.
	ClusterDasFdmAvailabilityStateNetworkPartitionedFromMaster = ClusterDasFdmAvailabilityState("networkPartitionedFromMaster")
	// A host is alive but is isolated from the management network.
	//
	// See `ClusterDasVmSettingsIsolationResponse_enum` for the criteria
	// used to determine whether a host is isolated.
	ClusterDasFdmAvailabilityStateNetworkIsolated = ClusterDasFdmAvailabilityState("networkIsolated")
	// The slave host appears to be down.
	//
	// This state is reported by the
	// master of a slave host.
	ClusterDasFdmAvailabilityStateHostDown = ClusterDasFdmAvailabilityState("hostDown")
	// An error occurred when initilizating the Fault Domain Manager
	// on a host due to a problem with installing the
	// agent or configuring it.
	//
	// This condition can often be cleared by
	// reconfiguring HA for the host. This state is reported by vCenter
	// Server.
	ClusterDasFdmAvailabilityStateInitializationError = ClusterDasFdmAvailabilityState("initializationError")
	// An error occurred when unconfiguring the Fault Domain Manager
	// running on a host.
	//
	// In order to clear this condition the host might
	// need to be reconnected to the cluster and reconfigured first.
	// This state is reported by vCenter
	// Server.
	ClusterDasFdmAvailabilityStateUninitializationError = ClusterDasFdmAvailabilityState("uninitializationError")
	// The Fault Domain Manager (FDM) on the host cannot be reached.
	//
	// This
	// state is reported in two unlikely situations.
	//   - First, it is reported by
	//     a master if the host responds to ICMP pings sent by the master over the
	//     management network but the FDM on the host cannot be reached by the master.
	//     This situation will occur if the FDM is unable to run or exit the
	//     uninitialized state.
	//   - Second, it is reported by vCenter Server if it cannot connect to a
	//     master nor the FDM for the host. This situation would occur if all hosts
	//     in the cluster failed but vCenter Server is still running. It may also
	//     occur if all FDMs are unable to run or exit the uninitialized state.
	ClusterDasFdmAvailabilityStateFdmUnreachable = ClusterDasFdmAvailabilityState("fdmUnreachable")
	// Config/Reconfig/upgrade operation has failed in first attempt and
	// a retry of these operations is scheduled.
	//
	// If any of the retry attempts succeed, the state is set to initialized.
	// If all retry attempts fail, the state is set to initializationError.
	// This state is reported by vCenter.
	ClusterDasFdmAvailabilityStateRetry = ClusterDasFdmAvailabilityState("retry")
)

func (e ClusterDasFdmAvailabilityState) Values() []ClusterDasFdmAvailabilityState {
	return []ClusterDasFdmAvailabilityState{
		ClusterDasFdmAvailabilityStateUninitialized,
		ClusterDasFdmAvailabilityStateElection,
		ClusterDasFdmAvailabilityStateMaster,
		ClusterDasFdmAvailabilityStateConnectedToMaster,
		ClusterDasFdmAvailabilityStateNetworkPartitionedFromMaster,
		ClusterDasFdmAvailabilityStateNetworkIsolated,
		ClusterDasFdmAvailabilityStateHostDown,
		ClusterDasFdmAvailabilityStateInitializationError,
		ClusterDasFdmAvailabilityStateUninitializationError,
		ClusterDasFdmAvailabilityStateFdmUnreachable,
		ClusterDasFdmAvailabilityStateRetry,
	}
}

func (e ClusterDasFdmAvailabilityState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasFdmAvailabilityState"] = reflect.TypeOf((*ClusterDasFdmAvailabilityState)(nil)).Elem()
	minAPIVersionForEnumValue["ClusterDasFdmAvailabilityState"] = map[string]string{
		"retry": "8.0.0.0",
	}
}

// The `ClusterDasVmSettingsIsolationResponse_enum` enum defines
// values that indicate whether or not the virtual machine should be
// powered off if a host determines that it is isolated from the rest of
// the cluster.
//
// Host network isolation occurs when a host is still running but it can no
// longer communicate with other hosts in the cluster and it cannot ping
// the configured isolation address(es). When the HA agent on a host loses
// contact with the other hosts, it will ping the isolation addresses. If
// the pings fail, the host will declare itself isolated.
//
// Once the HA agent declares the host isolated, it will initiate the
// isolation response workflow after a 30 second delay. You can use the FDM
// advanced option fdm.isolationPolicyDelaySec to increase the delay. For
// each virtual machine, the HA agent attempts to determine if a master is
// responsible for restarting the virtual machine. If it cannot make the
// determination, or there is a master that is responsible, the agent will
// apply the configured isolation response. This workflow will continue
// until the configuration policy, has been applied to all virtual
// machines, the agent reconnects to another HA agent in the cluster, or
// the isolation address pings start succeeding. If there is a master agent
// in the cluster, it will attempt to restart the virtual machines that
// were powered off during isolation.
//
// By default, the isolated host leaves its virtual machines powered on.
// You can override the isolation response default with a cluster-wide
// setting (`ClusterDasConfigInfo.defaultVmSettings`)
// or a virtual machine setting
// (`ClusterDasVmSettings.isolationResponse`).
//   - All isolation response values are valid for the
//     `ClusterDasVmSettings.isolationResponse`
//     property specified in a single virtual machine HA configuration.
//   - All values except for <code>clusterIsolationResponse</code> are valid
//     for the cluster-wide default HA configuration for virtual machines
//     (`ClusterDasConfigInfo.defaultVmSettings`).
//
// If you ensure that your network infrastructure is sufficiently redundant
// and that at least one network path is available at all times, host network
// isolation should be a rare occurrence.
type ClusterDasVmSettingsIsolationResponse string

const (
	// Do not power off the virtual machine in the event of a host network
	// isolation.
	ClusterDasVmSettingsIsolationResponseNone = ClusterDasVmSettingsIsolationResponse("none")
	// Power off the virtual machine in the event of a host network
	// isolation.
	ClusterDasVmSettingsIsolationResponsePowerOff = ClusterDasVmSettingsIsolationResponse("powerOff")
	// Shut down the virtual machine guest operating system in the event of
	// a host network isolation.
	//
	// If the guest operating system fails to
	// shutdown within five minutes, HA will initiate a forced power off.
	//
	// When you use the shutdown isolation response, failover can take
	// longer (compared to the
	// `powerOff`
	// response) because the virtual machine cannot fail over until it is
	// shutdown.
	ClusterDasVmSettingsIsolationResponseShutdown = ClusterDasVmSettingsIsolationResponse("shutdown")
	// Use the default isolation response defined for the cluster
	// that contains this virtual machine.
	ClusterDasVmSettingsIsolationResponseClusterIsolationResponse = ClusterDasVmSettingsIsolationResponse("clusterIsolationResponse")
)

func (e ClusterDasVmSettingsIsolationResponse) Values() []ClusterDasVmSettingsIsolationResponse {
	return []ClusterDasVmSettingsIsolationResponse{
		ClusterDasVmSettingsIsolationResponseNone,
		ClusterDasVmSettingsIsolationResponsePowerOff,
		ClusterDasVmSettingsIsolationResponseShutdown,
		ClusterDasVmSettingsIsolationResponseClusterIsolationResponse,
	}
}

func (e ClusterDasVmSettingsIsolationResponse) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasVmSettingsIsolationResponse"] = reflect.TypeOf((*ClusterDasVmSettingsIsolationResponse)(nil)).Elem()
}

// The `ClusterDasVmSettingsRestartPriority_enum` enum defines
// virtual machine restart priority values to resolve resource contention.
//
// The priority determines the preference that HA gives to a virtual
// machine if sufficient capacity is not available to power on all failed
// virtual machines. For example, high priority virtual machines on a host
// get preference over low priority virtual machines.
//
// All priority values are valid for the restart priority specified in a
// single virtual machine HA configuration (`ClusterDasVmConfigInfo.dasSettings`).
// All values except for <code>clusterRestartPriority</code> are valid for
// the cluster-wide default HA configuration for virtual machines
// (`ClusterDasConfigInfo.defaultVmSettings`).
type ClusterDasVmSettingsRestartPriority string

const (
	// vSphere HA is disabled for this virtual machine.
	ClusterDasVmSettingsRestartPriorityDisabled = ClusterDasVmSettingsRestartPriority("disabled")
	// Virtual machines with this priority have the lowest chance of
	// powering on after a failure if there is insufficient capacity on
	// hosts to meet all virtual machine needs.
	ClusterDasVmSettingsRestartPriorityLowest = ClusterDasVmSettingsRestartPriority("lowest")
	// Virtual machines with this priority have a lower chance of powering
	// on after a failure if there is insufficient capacity on hosts to meet
	// all virtual machine needs.
	ClusterDasVmSettingsRestartPriorityLow = ClusterDasVmSettingsRestartPriority("low")
	// Virtual machines with this priority have an intermediate chance of
	// powering on after a failure if there is insufficient capacity on
	// hosts to meet all virtual machine needs.
	ClusterDasVmSettingsRestartPriorityMedium = ClusterDasVmSettingsRestartPriority("medium")
	// Virtual machines with this priority have a higher chance of powering
	// on after a failure if there is insufficient capacity on hosts to meet
	// all virtual machine needs.
	ClusterDasVmSettingsRestartPriorityHigh = ClusterDasVmSettingsRestartPriority("high")
	// Virtual machines with this priority have the highest chance of
	// powering on after a failure if there is insufficient capacity on
	// hosts to meet all virtual machine needs.
	ClusterDasVmSettingsRestartPriorityHighest = ClusterDasVmSettingsRestartPriority("highest")
	// Virtual machines with this priority use the default restart
	// priority defined for the cluster that contains this virtual machine.
	ClusterDasVmSettingsRestartPriorityClusterRestartPriority = ClusterDasVmSettingsRestartPriority("clusterRestartPriority")
)

func (e ClusterDasVmSettingsRestartPriority) Values() []ClusterDasVmSettingsRestartPriority {
	return []ClusterDasVmSettingsRestartPriority{
		ClusterDasVmSettingsRestartPriorityDisabled,
		ClusterDasVmSettingsRestartPriorityLowest,
		ClusterDasVmSettingsRestartPriorityLow,
		ClusterDasVmSettingsRestartPriorityMedium,
		ClusterDasVmSettingsRestartPriorityHigh,
		ClusterDasVmSettingsRestartPriorityHighest,
		ClusterDasVmSettingsRestartPriorityClusterRestartPriority,
	}
}

func (e ClusterDasVmSettingsRestartPriority) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterDasVmSettingsRestartPriority"] = reflect.TypeOf((*ClusterDasVmSettingsRestartPriority)(nil)).Elem()
}

// Describes the operation type of the action.
//
// enterexitQuarantine suggests
// that the host is only exiting the quarantine state (i.e. not the
// maintenance mode).
type ClusterHostInfraUpdateHaModeActionOperationType string

const (
	ClusterHostInfraUpdateHaModeActionOperationTypeEnterQuarantine  = ClusterHostInfraUpdateHaModeActionOperationType("enterQuarantine")
	ClusterHostInfraUpdateHaModeActionOperationTypeExitQuarantine   = ClusterHostInfraUpdateHaModeActionOperationType("exitQuarantine")
	ClusterHostInfraUpdateHaModeActionOperationTypeEnterMaintenance = ClusterHostInfraUpdateHaModeActionOperationType("enterMaintenance")
)

func (e ClusterHostInfraUpdateHaModeActionOperationType) Values() []ClusterHostInfraUpdateHaModeActionOperationType {
	return []ClusterHostInfraUpdateHaModeActionOperationType{
		ClusterHostInfraUpdateHaModeActionOperationTypeEnterQuarantine,
		ClusterHostInfraUpdateHaModeActionOperationTypeExitQuarantine,
		ClusterHostInfraUpdateHaModeActionOperationTypeEnterMaintenance,
	}
}

func (e ClusterHostInfraUpdateHaModeActionOperationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterHostInfraUpdateHaModeActionOperationType"] = reflect.TypeOf((*ClusterHostInfraUpdateHaModeActionOperationType)(nil)).Elem()
}

type ClusterInfraUpdateHaConfigInfoBehaviorType string

const (
	// With this behavior configured, the proposed DRS recommendations
	// require manual approval before they are executed.
	ClusterInfraUpdateHaConfigInfoBehaviorTypeManual = ClusterInfraUpdateHaConfigInfoBehaviorType("Manual")
	// With this behavior configured, the proposed DRS recommendations are
	// executed immediately.
	ClusterInfraUpdateHaConfigInfoBehaviorTypeAutomated = ClusterInfraUpdateHaConfigInfoBehaviorType("Automated")
)

func (e ClusterInfraUpdateHaConfigInfoBehaviorType) Values() []ClusterInfraUpdateHaConfigInfoBehaviorType {
	return []ClusterInfraUpdateHaConfigInfoBehaviorType{
		ClusterInfraUpdateHaConfigInfoBehaviorTypeManual,
		ClusterInfraUpdateHaConfigInfoBehaviorTypeAutomated,
	}
}

func (e ClusterInfraUpdateHaConfigInfoBehaviorType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterInfraUpdateHaConfigInfoBehaviorType"] = reflect.TypeOf((*ClusterInfraUpdateHaConfigInfoBehaviorType)(nil)).Elem()
}

type ClusterInfraUpdateHaConfigInfoRemediationType string

const (
	// With this behavior configured, a degraded host will be recommended
	// to be placed in Quarantine Mode.
	ClusterInfraUpdateHaConfigInfoRemediationTypeQuarantineMode = ClusterInfraUpdateHaConfigInfoRemediationType("QuarantineMode")
	// With this behavior configured, a degraded host will be recommended
	// to be placed in Maintenance Mode.
	ClusterInfraUpdateHaConfigInfoRemediationTypeMaintenanceMode = ClusterInfraUpdateHaConfigInfoRemediationType("MaintenanceMode")
)

func (e ClusterInfraUpdateHaConfigInfoRemediationType) Values() []ClusterInfraUpdateHaConfigInfoRemediationType {
	return []ClusterInfraUpdateHaConfigInfoRemediationType{
		ClusterInfraUpdateHaConfigInfoRemediationTypeQuarantineMode,
		ClusterInfraUpdateHaConfigInfoRemediationTypeMaintenanceMode,
	}
}

func (e ClusterInfraUpdateHaConfigInfoRemediationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterInfraUpdateHaConfigInfoRemediationType"] = reflect.TypeOf((*ClusterInfraUpdateHaConfigInfoRemediationType)(nil)).Elem()
}

// Defines the options for a Datacenter::powerOnVm() invocation.
type ClusterPowerOnVmOption string

const (
	// Override the DRS automation level.
	//
	// Value type: `DrsBehavior_enum`
	// Default value: current behavior
	ClusterPowerOnVmOptionOverrideAutomationLevel = ClusterPowerOnVmOption("OverrideAutomationLevel")
	// Reserve resources for the powering-on VMs throughout the
	// power-on session.
	//
	// When this option is set to true, the server
	// will return at most one recommended host per manual VM, and
	// the VM's reservations are held on the recommended host until
	// the VM is actually powered on (either by applying the
	// recommendation or by a power-on request on the VM), or until
	// the recommendation is cancelled, or until the recommendation
	// expires. The expiration time is currently set to 10
	// minutes. This option does not have an effect on automatic VMs
	// since their recommendations are executed immediately. This
	// option is effective on DRS clusters only.
	// Value type: boolean
	// Default value: false
	ClusterPowerOnVmOptionReserveResources = ClusterPowerOnVmOption("ReserveResources")
)

func (e ClusterPowerOnVmOption) Values() []ClusterPowerOnVmOption {
	return []ClusterPowerOnVmOption{
		ClusterPowerOnVmOptionOverrideAutomationLevel,
		ClusterPowerOnVmOptionReserveResources,
	}
}

func (e ClusterPowerOnVmOption) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterPowerOnVmOption"] = reflect.TypeOf((*ClusterPowerOnVmOption)(nil)).Elem()
}

// Type of services for which Profile can be requested for
type ClusterProfileServiceType string

const (
	// Distributed Resource Scheduling
	ClusterProfileServiceTypeDRS = ClusterProfileServiceType("DRS")
	// High Availability
	ClusterProfileServiceTypeHA = ClusterProfileServiceType("HA")
	// Distributed Power Management
	ClusterProfileServiceTypeDPM = ClusterProfileServiceType("DPM")
	// Fault tolerance
	ClusterProfileServiceTypeFT = ClusterProfileServiceType("FT")
)

func (e ClusterProfileServiceType) Values() []ClusterProfileServiceType {
	return []ClusterProfileServiceType{
		ClusterProfileServiceTypeDRS,
		ClusterProfileServiceTypeHA,
		ClusterProfileServiceTypeDPM,
		ClusterProfileServiceTypeFT,
	}
}

func (e ClusterProfileServiceType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterProfileServiceType"] = reflect.TypeOf((*ClusterProfileServiceType)(nil)).Elem()
}

type ClusterSystemVMsConfigInfoDeploymentMode string

const (
	// System VMs are fully managed by the system.
	ClusterSystemVMsConfigInfoDeploymentModeSYSTEM_MANAGED = ClusterSystemVMsConfigInfoDeploymentMode("SYSTEM_MANAGED")
	// System VMs are absent on the managed entity.
	ClusterSystemVMsConfigInfoDeploymentModeABSENT = ClusterSystemVMsConfigInfoDeploymentMode("ABSENT")
)

func (e ClusterSystemVMsConfigInfoDeploymentMode) Values() []ClusterSystemVMsConfigInfoDeploymentMode {
	return []ClusterSystemVMsConfigInfoDeploymentMode{
		ClusterSystemVMsConfigInfoDeploymentModeSYSTEM_MANAGED,
		ClusterSystemVMsConfigInfoDeploymentModeABSENT,
	}
}

func (e ClusterSystemVMsConfigInfoDeploymentMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterSystemVMsConfigInfoDeploymentMode"] = reflect.TypeOf((*ClusterSystemVMsConfigInfoDeploymentMode)(nil)).Elem()
	minAPIVersionForType["ClusterSystemVMsConfigInfoDeploymentMode"] = "8.0.2.0"
}

// The VM policy settings that determine the response to
// storage failures.
type ClusterVmComponentProtectionSettingsStorageVmReaction string

const (
	// VM Component Protection service will not monitor or react to
	// the component failure.
	//
	// This setting does not affect other vSphere
	// HA services such as Host Monitoring or VM Health Monitoring.
	ClusterVmComponentProtectionSettingsStorageVmReactionDisabled = ClusterVmComponentProtectionSettingsStorageVmReaction("disabled")
	// VM Component Protection service will monitor component failures but
	// will not restart an affected VM.
	//
	// Rather it will notify users about
	// the component failures. This setting does not affect other vSphere HA
	// services such as Host Monitoring or VM Health Monitoring.
	ClusterVmComponentProtectionSettingsStorageVmReactionWarning = ClusterVmComponentProtectionSettingsStorageVmReaction("warning")
	// VM Component Protection service protects VMs conservatively.
	//
	// With this
	// setting, when the service can't determine that capacity is available to
	// restart a VM, it will favor keeping the VM running.
	ClusterVmComponentProtectionSettingsStorageVmReactionRestartConservative = ClusterVmComponentProtectionSettingsStorageVmReaction("restartConservative")
	// VM Component Protection service protects VMs aggressively.
	//
	// With this setting,
	// the service will terminate an affected VM even if it can't determine that
	// capacity exists to restart the VM.
	ClusterVmComponentProtectionSettingsStorageVmReactionRestartAggressive = ClusterVmComponentProtectionSettingsStorageVmReaction("restartAggressive")
	// VM will use the cluster default setting.
	//
	// This option is only meaningful for
	// per-VM settings.
	ClusterVmComponentProtectionSettingsStorageVmReactionClusterDefault = ClusterVmComponentProtectionSettingsStorageVmReaction("clusterDefault")
)

func (e ClusterVmComponentProtectionSettingsStorageVmReaction) Values() []ClusterVmComponentProtectionSettingsStorageVmReaction {
	return []ClusterVmComponentProtectionSettingsStorageVmReaction{
		ClusterVmComponentProtectionSettingsStorageVmReactionDisabled,
		ClusterVmComponentProtectionSettingsStorageVmReactionWarning,
		ClusterVmComponentProtectionSettingsStorageVmReactionRestartConservative,
		ClusterVmComponentProtectionSettingsStorageVmReactionRestartAggressive,
		ClusterVmComponentProtectionSettingsStorageVmReactionClusterDefault,
	}
}

func (e ClusterVmComponentProtectionSettingsStorageVmReaction) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterVmComponentProtectionSettingsStorageVmReaction"] = reflect.TypeOf((*ClusterVmComponentProtectionSettingsStorageVmReaction)(nil)).Elem()
}

// If an APD condition clears after an APD timeout condition has been declared and before
// VM Component Protection service terminated the VM, the guestOS and application may
// no longer be operational.
//
// VM Component Protection may be configured to reset the
// VM (`VirtualMachine.ResetVM_Task`) to restore the service of guest applications.
type ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared string

const (
	// VM Component Protection service will not react after APD condition is cleared.
	ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedNone = ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared("none")
	// VM Component Protection service will reset the VM after APD condition is cleared.
	//
	// Note this only applies if the subject VM is still powered on.
	ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedReset = ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared("reset")
	// VM will use the cluster default setting.
	//
	// This option is only meaningful for
	// per-VM settings.
	ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedUseClusterDefault = ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared("useClusterDefault")
)

func (e ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared) Values() []ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared {
	return []ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared{
		ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedNone,
		ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedReset,
		ClusterVmComponentProtectionSettingsVmReactionOnAPDClearedUseClusterDefault,
	}
}

func (e ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared"] = reflect.TypeOf((*ClusterVmComponentProtectionSettingsVmReactionOnAPDCleared)(nil)).Elem()
}

// Condition for VM's readiness
type ClusterVmReadinessReadyCondition string

const (
	// No ready condition specified.
	//
	// In case of vSphere HA, higher restart priority VMs are still
	// placed before lower priority VMs.
	ClusterVmReadinessReadyConditionNone = ClusterVmReadinessReadyCondition("none")
	// VM is powered on.
	ClusterVmReadinessReadyConditionPoweredOn = ClusterVmReadinessReadyCondition("poweredOn")
	// VM guest operating system is up and responding normally (VM tools
	// heartbeat status is green).
	ClusterVmReadinessReadyConditionGuestHbStatusGreen = ClusterVmReadinessReadyCondition("guestHbStatusGreen")
	// An application running inside the VM is responding normally.
	//
	// To enable Application Monitoring, you must first obtain the
	// appropriate SDK (or be using an application that supports VMware
	// Application Monitoring) and use it to set up customized heartbeats
	// for the applications you want to monitor.
	// See `ClusterDasConfigInfo.vmMonitoring`.
	ClusterVmReadinessReadyConditionAppHbStatusGreen = ClusterVmReadinessReadyCondition("appHbStatusGreen")
	// VM will use the cluster default setting.
	//
	// This option is only
	// meaningful for per-VM settings.
	ClusterVmReadinessReadyConditionUseClusterDefault = ClusterVmReadinessReadyCondition("useClusterDefault")
)

func (e ClusterVmReadinessReadyCondition) Values() []ClusterVmReadinessReadyCondition {
	return []ClusterVmReadinessReadyCondition{
		ClusterVmReadinessReadyConditionNone,
		ClusterVmReadinessReadyConditionPoweredOn,
		ClusterVmReadinessReadyConditionGuestHbStatusGreen,
		ClusterVmReadinessReadyConditionAppHbStatusGreen,
		ClusterVmReadinessReadyConditionUseClusterDefault,
	}
}

func (e ClusterVmReadinessReadyCondition) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ClusterVmReadinessReadyCondition"] = reflect.TypeOf((*ClusterVmReadinessReadyCondition)(nil)).Elem()
}

type ComplianceResultStatus string

const (
	// Entity is in Compliance
	ComplianceResultStatusCompliant = ComplianceResultStatus("compliant")
	// Entity is out of Compliance
	ComplianceResultStatusNonCompliant = ComplianceResultStatus("nonCompliant")
	// Compliance status of the entity is not known
	ComplianceResultStatusUnknown = ComplianceResultStatus("unknown")
	// Compliance check on this host is running.
	ComplianceResultStatusRunning = ComplianceResultStatus("running")
)

func (e ComplianceResultStatus) Values() []ComplianceResultStatus {
	return []ComplianceResultStatus{
		ComplianceResultStatusCompliant,
		ComplianceResultStatusNonCompliant,
		ComplianceResultStatusUnknown,
		ComplianceResultStatusRunning,
	}
}

func (e ComplianceResultStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ComplianceResultStatus"] = reflect.TypeOf((*ComplianceResultStatus)(nil)).Elem()
}

// The SPBM(Storage Policy Based Management) license state for a host
type ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState string

const (
	// The host is licensed
	ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseStateLicensed = ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState("licensed")
	// The host is not licensed
	ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseStateUnlicensed = ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState("unlicensed")
	// The host license information is unknown, this could happen if the
	// host is not in a available state
	ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseStateUnknown = ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState("unknown")
)

func (e ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState) Values() []ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState {
	return []ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState{
		ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseStateLicensed,
		ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseStateUnlicensed,
		ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseStateUnknown,
	}
}

func (e ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState"] = reflect.TypeOf((*ComputeResourceHostSPBMLicenseInfoHostSPBMLicenseState)(nil)).Elem()
}

// Config spec operation type.
type ConfigSpecOperation string

const (
	// Indicates the addition of an element to the configuration.
	ConfigSpecOperationAdd = ConfigSpecOperation("add")
	// Indicates the change of an element in the configuration.
	ConfigSpecOperationEdit = ConfigSpecOperation("edit")
	// Indicates the removal of an element in the configuration.
	ConfigSpecOperationRemove = ConfigSpecOperation("remove")
)

func (e ConfigSpecOperation) Values() []ConfigSpecOperation {
	return []ConfigSpecOperation{
		ConfigSpecOperationAdd,
		ConfigSpecOperationEdit,
		ConfigSpecOperationRemove,
	}
}

func (e ConfigSpecOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ConfigSpecOperation"] = reflect.TypeOf((*ConfigSpecOperation)(nil)).Elem()
}

type CryptoManagerHostKeyManagementType string

const (
	CryptoManagerHostKeyManagementTypeUnknown  = CryptoManagerHostKeyManagementType("unknown")
	CryptoManagerHostKeyManagementTypeInternal = CryptoManagerHostKeyManagementType("internal")
	CryptoManagerHostKeyManagementTypeExternal = CryptoManagerHostKeyManagementType("external")
)

func (e CryptoManagerHostKeyManagementType) Values() []CryptoManagerHostKeyManagementType {
	return []CryptoManagerHostKeyManagementType{
		CryptoManagerHostKeyManagementTypeUnknown,
		CryptoManagerHostKeyManagementTypeInternal,
		CryptoManagerHostKeyManagementTypeExternal,
	}
}

func (e CryptoManagerHostKeyManagementType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CryptoManagerHostKeyManagementType"] = reflect.TypeOf((*CryptoManagerHostKeyManagementType)(nil)).Elem()
	minAPIVersionForType["CryptoManagerHostKeyManagementType"] = "8.0.1.0"
}

type CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason string

const (
	// Key not found in VC cache and does not specify a provider
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateMissingInCache = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateMissingInCache")
	// Key provider is invalid
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateClusterInvalid = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateClusterInvalid")
	// Can not reach the key provider
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateClusterUnreachable = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateClusterUnreachable")
	// Key not found in KMS
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateMissingInKMS = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateMissingInKMS")
	// Key not active or enabled
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateNotActiveOrEnabled = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateNotActiveOrEnabled")
	// Key is managed by Trust Authority
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateManagedByTrustAuthority = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateManagedByTrustAuthority")
	// Key is managed by Native Key Provider
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateManagedByNKP = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("KeyStateManagedByNKP")
	// No permission to access key provider
	CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonNoPermissionToAccessKeyProvider = CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason("NoPermissionToAccessKeyProvider")
)

func (e CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason) Values() []CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason {
	return []CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason{
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateMissingInCache,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateClusterInvalid,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateClusterUnreachable,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateMissingInKMS,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateNotActiveOrEnabled,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateManagedByTrustAuthority,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonKeyStateManagedByNKP,
		CryptoManagerKmipCryptoKeyStatusKeyUnavailableReasonNoPermissionToAccessKeyProvider,
	}
}

func (e CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason"] = reflect.TypeOf((*CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason)(nil)).Elem()
	minAPIVersionForEnumValue["CryptoManagerKmipCryptoKeyStatusKeyUnavailableReason"] = map[string]string{
		"KeyStateManagedByNKP":            "8.0.3.0",
		"NoPermissionToAccessKeyProvider": "8.0.3.0",
	}
}

type CustomizationFailedReasonCode string

const (
	// The user defined script is disabled during customization
	CustomizationFailedReasonCodeUserDefinedScriptDisabled = CustomizationFailedReasonCode("userDefinedScriptDisabled")
	// The guest customization is disabled by VMware Tools
	CustomizationFailedReasonCodeCustomizationDisabled = CustomizationFailedReasonCode("customizationDisabled")
	// The cloud-init version is too old to support cloud-init raw data
	CustomizationFailedReasonCodeRawDataIsNotSupported = CustomizationFailedReasonCode("rawDataIsNotSupported")
	// The cloud-init meta data is not valid format
	CustomizationFailedReasonCodeWrongMetadataFormat = CustomizationFailedReasonCode("wrongMetadataFormat")
)

func (e CustomizationFailedReasonCode) Values() []CustomizationFailedReasonCode {
	return []CustomizationFailedReasonCode{
		CustomizationFailedReasonCodeUserDefinedScriptDisabled,
		CustomizationFailedReasonCodeCustomizationDisabled,
		CustomizationFailedReasonCodeRawDataIsNotSupported,
		CustomizationFailedReasonCodeWrongMetadataFormat,
	}
}

func (e CustomizationFailedReasonCode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CustomizationFailedReasonCode"] = reflect.TypeOf((*CustomizationFailedReasonCode)(nil)).Elem()
	minAPIVersionForEnumValue["CustomizationFailedReasonCode"] = map[string]string{
		"customizationDisabled": "7.0.1.0",
		"rawDataIsNotSupported": "7.0.3.0",
		"wrongMetadataFormat":   "7.0.3.0",
	}
}

// Enumeration of AutoMode values.
type CustomizationLicenseDataMode string

const (
	// Indicates that client access licenses have been purchased for the server,
	// allowing a certain number of concurrent connections to the VirtualCenter
	// server.
	CustomizationLicenseDataModePerServer = CustomizationLicenseDataMode("perServer")
	// Indicates that a client access license has been purchased for each computer
	// that accesses the VirtualCenter server.
	CustomizationLicenseDataModePerSeat = CustomizationLicenseDataMode("perSeat")
)

func (e CustomizationLicenseDataMode) Values() []CustomizationLicenseDataMode {
	return []CustomizationLicenseDataMode{
		CustomizationLicenseDataModePerServer,
		CustomizationLicenseDataModePerSeat,
	}
}

func (e CustomizationLicenseDataMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CustomizationLicenseDataMode"] = reflect.TypeOf((*CustomizationLicenseDataMode)(nil)).Elem()
}

// NetBIOS setting for Windows.
type CustomizationNetBIOSMode string

const (
	// DHCP server decides whether or not to use NetBIOS.
	CustomizationNetBIOSModeEnableNetBIOSViaDhcp = CustomizationNetBIOSMode("enableNetBIOSViaDhcp")
	// Always use NetBIOS.
	CustomizationNetBIOSModeEnableNetBIOS = CustomizationNetBIOSMode("enableNetBIOS")
	// Never use NetBIOS.
	CustomizationNetBIOSModeDisableNetBIOS = CustomizationNetBIOSMode("disableNetBIOS")
)

func (e CustomizationNetBIOSMode) Values() []CustomizationNetBIOSMode {
	return []CustomizationNetBIOSMode{
		CustomizationNetBIOSModeEnableNetBIOSViaDhcp,
		CustomizationNetBIOSModeEnableNetBIOS,
		CustomizationNetBIOSModeDisableNetBIOS,
	}
}

func (e CustomizationNetBIOSMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CustomizationNetBIOSMode"] = reflect.TypeOf((*CustomizationNetBIOSMode)(nil)).Elem()
}

// A enum constant specifying what should be done to the guest vm after running
// sysprep.
type CustomizationSysprepRebootOption string

const (
	// Reboot the machine after running sysprep.
	//
	// This will cause values
	// specified in the sysprep.xml to be applied immediately.
	CustomizationSysprepRebootOptionReboot = CustomizationSysprepRebootOption("reboot")
	// Take no action.
	//
	// Leave the guest os running after running sysprep. This
	// option can be used to look at values for debugging purposes after
	// running sysprep.
	CustomizationSysprepRebootOptionNoreboot = CustomizationSysprepRebootOption("noreboot")
	// Shutdown the machine after running sysprep.
	//
	// This puts the vm in a
	// sealed state.
	CustomizationSysprepRebootOptionShutdown = CustomizationSysprepRebootOption("shutdown")
)

func (e CustomizationSysprepRebootOption) Values() []CustomizationSysprepRebootOption {
	return []CustomizationSysprepRebootOption{
		CustomizationSysprepRebootOptionReboot,
		CustomizationSysprepRebootOptionNoreboot,
		CustomizationSysprepRebootOptionShutdown,
	}
}

func (e CustomizationSysprepRebootOption) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["CustomizationSysprepRebootOption"] = reflect.TypeOf((*CustomizationSysprepRebootOption)(nil)).Elem()
}

// Set of possible values for
// `DVPortStatus*.*DVPortStatus.vmDirectPathGen2InactiveReasonNetwork`.
type DVPortStatusVmDirectPathGen2InactiveReasonNetwork string

const (
	// The switch for which this port is defined does not support VMDirectPath Gen 2.
	//
	// See
	// `DVSFeatureCapability*.*DVSFeatureCapability.vmDirectPathGen2Supported`.
	DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptIncompatibleDvs = DVPortStatusVmDirectPathGen2InactiveReasonNetwork("portNptIncompatibleDvs")
	// None of the physical NICs used as uplinks for this port support
	// VMDirectPath Gen 2.
	//
	// See also `PhysicalNic.vmDirectPathGen2Supported`.
	DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptNoCompatibleNics = DVPortStatusVmDirectPathGen2InactiveReasonNetwork("portNptNoCompatibleNics")
	// At least some of the physical NICs used as uplinks for this port
	// support VMDirectPath Gen 2, but all available network-passthrough
	// resources are in use by other ports.
	DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptNoVirtualFunctionsAvailable = DVPortStatusVmDirectPathGen2InactiveReasonNetwork("portNptNoVirtualFunctionsAvailable")
	// VMDirectPath Gen 2 has been explicitly disabled for this port.
	DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptDisabledForPort = DVPortStatusVmDirectPathGen2InactiveReasonNetwork("portNptDisabledForPort")
)

func (e DVPortStatusVmDirectPathGen2InactiveReasonNetwork) Values() []DVPortStatusVmDirectPathGen2InactiveReasonNetwork {
	return []DVPortStatusVmDirectPathGen2InactiveReasonNetwork{
		DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptIncompatibleDvs,
		DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptNoCompatibleNics,
		DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptNoVirtualFunctionsAvailable,
		DVPortStatusVmDirectPathGen2InactiveReasonNetworkPortNptDisabledForPort,
	}
}

func (e DVPortStatusVmDirectPathGen2InactiveReasonNetwork) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DVPortStatusVmDirectPathGen2InactiveReasonNetwork"] = reflect.TypeOf((*DVPortStatusVmDirectPathGen2InactiveReasonNetwork)(nil)).Elem()
}

// Set of possible values for
// `DVPortStatus*.*DVPortStatus.vmDirectPathGen2InactiveReasonOther`.
type DVPortStatusVmDirectPathGen2InactiveReasonOther string

const (
	// The host for which this port is defined does not support VMDirectPath Gen 2.
	//
	// See `HostCapability*.*HostCapability.vmDirectPathGen2Supported`
	DVPortStatusVmDirectPathGen2InactiveReasonOtherPortNptIncompatibleHost = DVPortStatusVmDirectPathGen2InactiveReasonOther("portNptIncompatibleHost")
	// Configuration or state of the port's connectee prevents
	// VMDirectPath Gen 2.
	//
	// See
	// `VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeState.vmDirectPathGen2InactiveReasonVm`
	// and/or
	// `VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeState.vmDirectPathGen2InactiveReasonExtended`
	// in the appropriate element of the RuntimeInfo.device array of the
	// virtual machine connected to this port.
	DVPortStatusVmDirectPathGen2InactiveReasonOtherPortNptIncompatibleConnectee = DVPortStatusVmDirectPathGen2InactiveReasonOther("portNptIncompatibleConnectee")
)

func (e DVPortStatusVmDirectPathGen2InactiveReasonOther) Values() []DVPortStatusVmDirectPathGen2InactiveReasonOther {
	return []DVPortStatusVmDirectPathGen2InactiveReasonOther{
		DVPortStatusVmDirectPathGen2InactiveReasonOtherPortNptIncompatibleHost,
		DVPortStatusVmDirectPathGen2InactiveReasonOtherPortNptIncompatibleConnectee,
	}
}

func (e DVPortStatusVmDirectPathGen2InactiveReasonOther) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DVPortStatusVmDirectPathGen2InactiveReasonOther"] = reflect.TypeOf((*DVPortStatusVmDirectPathGen2InactiveReasonOther)(nil)).Elem()
}

type DVSFilterSpecLinkConfig string

const (
	// The port link state: blocked.
	DVSFilterSpecLinkConfigBlocked = DVSFilterSpecLinkConfig("blocked")
	// The port link state: unblocked.
	DVSFilterSpecLinkConfigUnblocked = DVSFilterSpecLinkConfig("unblocked")
)

func (e DVSFilterSpecLinkConfig) Values() []DVSFilterSpecLinkConfig {
	return []DVSFilterSpecLinkConfig{
		DVSFilterSpecLinkConfigBlocked,
		DVSFilterSpecLinkConfigUnblocked,
	}
}

func (e DVSFilterSpecLinkConfig) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DVSFilterSpecLinkConfig"] = reflect.TypeOf((*DVSFilterSpecLinkConfig)(nil)).Elem()
}

type DVSFilterSpecLinkState string

const (
	// The port link state: down.
	DVSFilterSpecLinkStateDown = DVSFilterSpecLinkState("down")
	// The port link state: up.
	DVSFilterSpecLinkStateUp = DVSFilterSpecLinkState("up")
)

func (e DVSFilterSpecLinkState) Values() []DVSFilterSpecLinkState {
	return []DVSFilterSpecLinkState{
		DVSFilterSpecLinkStateDown,
		DVSFilterSpecLinkStateUp,
	}
}

func (e DVSFilterSpecLinkState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DVSFilterSpecLinkState"] = reflect.TypeOf((*DVSFilterSpecLinkState)(nil)).Elem()
}

type DVSMacLimitPolicyType string

const (
	DVSMacLimitPolicyTypeAllow = DVSMacLimitPolicyType("allow")
	DVSMacLimitPolicyTypeDrop  = DVSMacLimitPolicyType("drop")
)

func (e DVSMacLimitPolicyType) Values() []DVSMacLimitPolicyType {
	return []DVSMacLimitPolicyType{
		DVSMacLimitPolicyTypeAllow,
		DVSMacLimitPolicyTypeDrop,
	}
}

func (e DVSMacLimitPolicyType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DVSMacLimitPolicyType"] = reflect.TypeOf((*DVSMacLimitPolicyType)(nil)).Elem()
}

type DasConfigFaultDasConfigFaultReason string

const (
	// There is a problem with the host network configuration.
	DasConfigFaultDasConfigFaultReasonHostNetworkMisconfiguration = DasConfigFaultDasConfigFaultReason("HostNetworkMisconfiguration")
	// There is a problem with the host configuration.
	DasConfigFaultDasConfigFaultReasonHostMisconfiguration = DasConfigFaultDasConfigFaultReason("HostMisconfiguration")
	// The privileges were insuffient for the operation.
	DasConfigFaultDasConfigFaultReasonInsufficientPrivileges = DasConfigFaultDasConfigFaultReason("InsufficientPrivileges")
	// There was no running primary agent available to contact.
	//
	// Check that your other hosts don't have HA errors
	DasConfigFaultDasConfigFaultReasonNoPrimaryAgentAvailable = DasConfigFaultDasConfigFaultReason("NoPrimaryAgentAvailable")
	// The HA configuration failed for other reasons.
	DasConfigFaultDasConfigFaultReasonOther = DasConfigFaultDasConfigFaultReason("Other")
	// No datastores defined for this host
	DasConfigFaultDasConfigFaultReasonNoDatastoresConfigured = DasConfigFaultDasConfigFaultReason("NoDatastoresConfigured")
	// Failure to create config vvol
	DasConfigFaultDasConfigFaultReasonCreateConfigVvolFailed = DasConfigFaultDasConfigFaultReason("CreateConfigVvolFailed")
	// Host in vSAN cluster does not support vSAN.
	DasConfigFaultDasConfigFaultReasonVSanNotSupportedOnHost = DasConfigFaultDasConfigFaultReason("VSanNotSupportedOnHost")
	// There is a problem with the cluster network configuration.
	DasConfigFaultDasConfigFaultReasonDasNetworkMisconfiguration = DasConfigFaultDasConfigFaultReason("DasNetworkMisconfiguration")
	// Setting desired imageSpec in Personality Manager failed
	DasConfigFaultDasConfigFaultReasonSetDesiredImageSpecFailed = DasConfigFaultDasConfigFaultReason("SetDesiredImageSpecFailed")
	// The ApplyHA call to Personality Manager failed
	DasConfigFaultDasConfigFaultReasonApplyHAVibsOnClusterFailed = DasConfigFaultDasConfigFaultReason("ApplyHAVibsOnClusterFailed")
)

func (e DasConfigFaultDasConfigFaultReason) Values() []DasConfigFaultDasConfigFaultReason {
	return []DasConfigFaultDasConfigFaultReason{
		DasConfigFaultDasConfigFaultReasonHostNetworkMisconfiguration,
		DasConfigFaultDasConfigFaultReasonHostMisconfiguration,
		DasConfigFaultDasConfigFaultReasonInsufficientPrivileges,
		DasConfigFaultDasConfigFaultReasonNoPrimaryAgentAvailable,
		DasConfigFaultDasConfigFaultReasonOther,
		DasConfigFaultDasConfigFaultReasonNoDatastoresConfigured,
		DasConfigFaultDasConfigFaultReasonCreateConfigVvolFailed,
		DasConfigFaultDasConfigFaultReasonVSanNotSupportedOnHost,
		DasConfigFaultDasConfigFaultReasonDasNetworkMisconfiguration,
		DasConfigFaultDasConfigFaultReasonSetDesiredImageSpecFailed,
		DasConfigFaultDasConfigFaultReasonApplyHAVibsOnClusterFailed,
	}
}

func (e DasConfigFaultDasConfigFaultReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DasConfigFaultDasConfigFaultReason"] = reflect.TypeOf((*DasConfigFaultDasConfigFaultReason)(nil)).Elem()
}

// Deprecated as of VI API 2.5, use `ClusterDasVmSettingsRestartPriority_enum`.
//
// The priority of the virtual machine determines the preference
// given to it if sufficient capacity is not available to power
// on all failed virtual machines.
//
// For example, high priority
// virtual machines on a host get preference over low priority
// virtual machines.
type DasVmPriority string

const (
	// vSphere HA is disabled for this virtual machine.
	DasVmPriorityDisabled = DasVmPriority("disabled")
	// Virtual machines with this priority have a lower chance of powering on after a
	// failure if there is insufficient capacity on hosts to meet all virtual machine
	// needs.
	DasVmPriorityLow = DasVmPriority("low")
	// Virtual machines with this priority have an intermediate chance of powering
	// on after a failure if there is insufficient capacity on hosts to meet all
	// virtual machine needs.
	DasVmPriorityMedium = DasVmPriority("medium")
	// Virtual machines with this priority have a higher chance of powering on after a
	// failure if there is insufficient capacity on hosts to meet all virtual machine
	// needs.
	DasVmPriorityHigh = DasVmPriority("high")
)

func (e DasVmPriority) Values() []DasVmPriority {
	return []DasVmPriority{
		DasVmPriorityDisabled,
		DasVmPriorityLow,
		DasVmPriorityMedium,
		DasVmPriorityHigh,
	}
}

func (e DasVmPriority) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DasVmPriority"] = reflect.TypeOf((*DasVmPriority)(nil)).Elem()
}

type DatastoreAccessible string

const (
	// Is accessible
	DatastoreAccessibleTrue = DatastoreAccessible("True")
	// Is not accessible
	DatastoreAccessibleFalse = DatastoreAccessible("False")
)

func (e DatastoreAccessible) Values() []DatastoreAccessible {
	return []DatastoreAccessible{
		DatastoreAccessibleTrue,
		DatastoreAccessibleFalse,
	}
}

func (e DatastoreAccessible) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DatastoreAccessible"] = reflect.TypeOf((*DatastoreAccessible)(nil)).Elem()
}

// Defines the current maintenance mode state of the datastore.
type DatastoreSummaryMaintenanceModeState string

const (
	// Default state.
	DatastoreSummaryMaintenanceModeStateNormal = DatastoreSummaryMaintenanceModeState("normal")
	// Started entering maintenance mode, but not finished.
	//
	// This could happen when waiting for user input or for
	// long-running vmotions to complete.
	DatastoreSummaryMaintenanceModeStateEnteringMaintenance = DatastoreSummaryMaintenanceModeState("enteringMaintenance")
	// Successfully entered maintenance mode.
	DatastoreSummaryMaintenanceModeStateInMaintenance = DatastoreSummaryMaintenanceModeState("inMaintenance")
)

func (e DatastoreSummaryMaintenanceModeState) Values() []DatastoreSummaryMaintenanceModeState {
	return []DatastoreSummaryMaintenanceModeState{
		DatastoreSummaryMaintenanceModeStateNormal,
		DatastoreSummaryMaintenanceModeStateEnteringMaintenance,
		DatastoreSummaryMaintenanceModeStateInMaintenance,
	}
}

func (e DatastoreSummaryMaintenanceModeState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DatastoreSummaryMaintenanceModeState"] = reflect.TypeOf((*DatastoreSummaryMaintenanceModeState)(nil)).Elem()
}

type DayOfWeek string

const (
	DayOfWeekSunday    = DayOfWeek("sunday")
	DayOfWeekMonday    = DayOfWeek("monday")
	DayOfWeekTuesday   = DayOfWeek("tuesday")
	DayOfWeekWednesday = DayOfWeek("wednesday")
	DayOfWeekThursday  = DayOfWeek("thursday")
	DayOfWeekFriday    = DayOfWeek("friday")
	DayOfWeekSaturday  = DayOfWeek("saturday")
)

func (e DayOfWeek) Values() []DayOfWeek {
	return []DayOfWeek{
		DayOfWeekSunday,
		DayOfWeekMonday,
		DayOfWeekTuesday,
		DayOfWeekWednesday,
		DayOfWeekThursday,
		DayOfWeekFriday,
		DayOfWeekSaturday,
	}
}

func (e DayOfWeek) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DayOfWeek"] = reflect.TypeOf((*DayOfWeek)(nil)).Elem()
}

// Reasons why a virtual device would not be supported on a host.
type DeviceNotSupportedReason string

const (
	// The host does not support this virtual device at all.
	DeviceNotSupportedReasonHost = DeviceNotSupportedReason("host")
	// The device is supported by the host in general, but not for
	// the specific guest OS the virtual machine is using.
	DeviceNotSupportedReasonGuest = DeviceNotSupportedReason("guest")
	// The device is supported by the host and guest OS, but not for
	// the vSphere Fault Tolerance.
	DeviceNotSupportedReasonFt = DeviceNotSupportedReason("ft")
)

func (e DeviceNotSupportedReason) Values() []DeviceNotSupportedReason {
	return []DeviceNotSupportedReason{
		DeviceNotSupportedReasonHost,
		DeviceNotSupportedReasonGuest,
		DeviceNotSupportedReasonFt,
	}
}

func (e DeviceNotSupportedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DeviceNotSupportedReason"] = reflect.TypeOf((*DeviceNotSupportedReason)(nil)).Elem()
	minAPIVersionForEnumValue["DeviceNotSupportedReason"] = map[string]string{
		"ft": "8.0.3.0",
	}
}

// The list of Device Protocols.
type DeviceProtocol string

const (
	DeviceProtocolNVMe = DeviceProtocol("NVMe")
	DeviceProtocolSCSI = DeviceProtocol("SCSI")
)

func (e DeviceProtocol) Values() []DeviceProtocol {
	return []DeviceProtocol{
		DeviceProtocolNVMe,
		DeviceProtocolSCSI,
	}
}

func (e DeviceProtocol) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DeviceProtocol"] = reflect.TypeOf((*DeviceProtocol)(nil)).Elem()
	minAPIVersionForType["DeviceProtocol"] = "8.0.1.0"
}

// Pre-defined constants for possible creators of log files.
type DiagnosticManagerLogCreator string

const (
	// VirtualCenter service
	DiagnosticManagerLogCreatorVpxd = DiagnosticManagerLogCreator("vpxd")
	// VirtualCenter agent
	DiagnosticManagerLogCreatorVpxa = DiagnosticManagerLogCreator("vpxa")
	// Host agent
	DiagnosticManagerLogCreatorHostd = DiagnosticManagerLogCreator("hostd")
	// Host server agent
	DiagnosticManagerLogCreatorServerd = DiagnosticManagerLogCreator("serverd")
	// Installation
	DiagnosticManagerLogCreatorInstall = DiagnosticManagerLogCreator("install")
	// Virtual infrastructure client
	DiagnosticManagerLogCreatorVpxClient = DiagnosticManagerLogCreator("vpxClient")
	// System Record Log
	DiagnosticManagerLogCreatorRecordLog = DiagnosticManagerLogCreator("recordLog")
)

func (e DiagnosticManagerLogCreator) Values() []DiagnosticManagerLogCreator {
	return []DiagnosticManagerLogCreator{
		DiagnosticManagerLogCreatorVpxd,
		DiagnosticManagerLogCreatorVpxa,
		DiagnosticManagerLogCreatorHostd,
		DiagnosticManagerLogCreatorServerd,
		DiagnosticManagerLogCreatorInstall,
		DiagnosticManagerLogCreatorVpxClient,
		DiagnosticManagerLogCreatorRecordLog,
	}
}

func (e DiagnosticManagerLogCreator) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DiagnosticManagerLogCreator"] = reflect.TypeOf((*DiagnosticManagerLogCreator)(nil)).Elem()
}

// Constants for defined formats.
//
// For more information, see the comment for the format property.
type DiagnosticManagerLogFormat string

const (
	// A standard ASCII-based line-based log file.
	DiagnosticManagerLogFormatPlain = DiagnosticManagerLogFormat("plain")
)

func (e DiagnosticManagerLogFormat) Values() []DiagnosticManagerLogFormat {
	return []DiagnosticManagerLogFormat{
		DiagnosticManagerLogFormatPlain,
	}
}

func (e DiagnosticManagerLogFormat) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DiagnosticManagerLogFormat"] = reflect.TypeOf((*DiagnosticManagerLogFormat)(nil)).Elem()
}

// Type of partition indicating the type of storage on which the partition
// resides.
//
// If the diagnostic partition is local only, it will only need
// one slot. If the diagnostic partition is on shared storage, it could
// be used by multiple hosts. As a result, it will need multiple slots.
type DiagnosticPartitionStorageType string

const (
	DiagnosticPartitionStorageTypeDirectAttached  = DiagnosticPartitionStorageType("directAttached")
	DiagnosticPartitionStorageTypeNetworkAttached = DiagnosticPartitionStorageType("networkAttached")
)

func (e DiagnosticPartitionStorageType) Values() []DiagnosticPartitionStorageType {
	return []DiagnosticPartitionStorageType{
		DiagnosticPartitionStorageTypeDirectAttached,
		DiagnosticPartitionStorageTypeNetworkAttached,
	}
}

func (e DiagnosticPartitionStorageType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DiagnosticPartitionStorageType"] = reflect.TypeOf((*DiagnosticPartitionStorageType)(nil)).Elem()
}

// The type of diagnostic partition.
//
// Private diagnostic partition has one
// slot, so can only be used by one host. Shared diagnostic parititon
// needs multiple slots so to be usable by multiple hosts.
type DiagnosticPartitionType string

const (
	DiagnosticPartitionTypeSingleHost = DiagnosticPartitionType("singleHost")
	DiagnosticPartitionTypeMultiHost  = DiagnosticPartitionType("multiHost")
)

func (e DiagnosticPartitionType) Values() []DiagnosticPartitionType {
	return []DiagnosticPartitionType{
		DiagnosticPartitionTypeSingleHost,
		DiagnosticPartitionTypeMultiHost,
	}
}

func (e DiagnosticPartitionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DiagnosticPartitionType"] = reflect.TypeOf((*DiagnosticPartitionType)(nil)).Elem()
}

// The disallowed change type.
type DisallowedChangeByServiceDisallowedChange string

const (
	// Online extend disk operation.
	DisallowedChangeByServiceDisallowedChangeHotExtendDisk = DisallowedChangeByServiceDisallowedChange("hotExtendDisk")
)

func (e DisallowedChangeByServiceDisallowedChange) Values() []DisallowedChangeByServiceDisallowedChange {
	return []DisallowedChangeByServiceDisallowedChange{
		DisallowedChangeByServiceDisallowedChangeHotExtendDisk,
	}
}

func (e DisallowedChangeByServiceDisallowedChange) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DisallowedChangeByServiceDisallowedChange"] = reflect.TypeOf((*DisallowedChangeByServiceDisallowedChange)(nil)).Elem()
}

// The `DistributedVirtualPortgroupBackingType_enum` enum defines
// the distributed virtual portgroup backing type.
type DistributedVirtualPortgroupBackingType string

const (
	// The portgroup is created by vCenter.
	DistributedVirtualPortgroupBackingTypeStandard = DistributedVirtualPortgroupBackingType("standard")
	// The portgroup is created by NSX manager.
	//
	// For NSX backing type, We only support ephemeral portgroup type.
	// If `DistributedVirtualPortgroupPortgroupType_enum` is
	// ephemeral, A `DistributedVirtualPort` will be
	// dynamicly created by NSX when the virtual machine is reconfigured
	// to connect to the portgroup.
	DistributedVirtualPortgroupBackingTypeNsx = DistributedVirtualPortgroupBackingType("nsx")
)

func (e DistributedVirtualPortgroupBackingType) Values() []DistributedVirtualPortgroupBackingType {
	return []DistributedVirtualPortgroupBackingType{
		DistributedVirtualPortgroupBackingTypeStandard,
		DistributedVirtualPortgroupBackingTypeNsx,
	}
}

func (e DistributedVirtualPortgroupBackingType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualPortgroupBackingType"] = reflect.TypeOf((*DistributedVirtualPortgroupBackingType)(nil)).Elem()
}

// The meta tag names recognizable in the
// `DVPortgroupConfigInfo.portNameFormat` string.
type DistributedVirtualPortgroupMetaTagName string

const (
	// This tag will be expanded to the name of the switch.
	DistributedVirtualPortgroupMetaTagNameDvsName = DistributedVirtualPortgroupMetaTagName("dvsName")
	// This tag will be expanded to the name of the portgroup.
	DistributedVirtualPortgroupMetaTagNamePortgroupName = DistributedVirtualPortgroupMetaTagName("portgroupName")
	// This tag will be expanded to the current index of the port.
	DistributedVirtualPortgroupMetaTagNamePortIndex = DistributedVirtualPortgroupMetaTagName("portIndex")
)

func (e DistributedVirtualPortgroupMetaTagName) Values() []DistributedVirtualPortgroupMetaTagName {
	return []DistributedVirtualPortgroupMetaTagName{
		DistributedVirtualPortgroupMetaTagNameDvsName,
		DistributedVirtualPortgroupMetaTagNamePortgroupName,
		DistributedVirtualPortgroupMetaTagNamePortIndex,
	}
}

func (e DistributedVirtualPortgroupMetaTagName) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualPortgroupMetaTagName"] = reflect.TypeOf((*DistributedVirtualPortgroupMetaTagName)(nil)).Elem()
}

// The `DistributedVirtualPortgroupPortgroupType_enum` enum defines
// the distributed virtual portgroup types
// (`DistributedVirtualPortgroup*.*DistributedVirtualPortgroup.config*.*DVPortgroupConfigInfo.type`).
//
// Early binding specifies a static set of ports that are created
// when you create the distributed virtual portgroup. An ephemeral portgroup uses dynamic
// ports that are created when you power on a virtual machine.
type DistributedVirtualPortgroupPortgroupType string

const (
	// A free `DistributedVirtualPort` will be selected and assigned to
	// a `VirtualMachine` when the virtual machine is reconfigured to
	// connect to the portgroup.
	DistributedVirtualPortgroupPortgroupTypeEarlyBinding = DistributedVirtualPortgroupPortgroupType("earlyBinding")
	// Deprecated as of vSphere API 5.0.
	//
	// A free `DistributedVirtualPort` will be selected and
	// assigned to a `VirtualMachine` when the virtual machine is
	// powered on.
	DistributedVirtualPortgroupPortgroupTypeLateBinding = DistributedVirtualPortgroupPortgroupType("lateBinding")
	// A `DistributedVirtualPort` will be created and assigned to a
	// `VirtualMachine` when the virtual machine is powered on, and will
	// be deleted when the virtual machine is powered off.
	//
	// An ephemeral portgroup has
	// no limit on the number of ports that can be a part of this portgroup.
	// In cases where the vCenter Server is unavailable the host can
	// create conflict ports in this portgroup to be used by a virtual machine
	// at power on.
	DistributedVirtualPortgroupPortgroupTypeEphemeral = DistributedVirtualPortgroupPortgroupType("ephemeral")
)

func (e DistributedVirtualPortgroupPortgroupType) Values() []DistributedVirtualPortgroupPortgroupType {
	return []DistributedVirtualPortgroupPortgroupType{
		DistributedVirtualPortgroupPortgroupTypeEarlyBinding,
		DistributedVirtualPortgroupPortgroupTypeLateBinding,
		DistributedVirtualPortgroupPortgroupTypeEphemeral,
	}
}

func (e DistributedVirtualPortgroupPortgroupType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualPortgroupPortgroupType"] = reflect.TypeOf((*DistributedVirtualPortgroupPortgroupType)(nil)).Elem()
}

// List of possible host infrastructure traffic classes
type DistributedVirtualSwitchHostInfrastructureTrafficClass string

const (
	// Management Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassManagement = DistributedVirtualSwitchHostInfrastructureTrafficClass("management")
	// Fault Tolerance (FT) Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassFaultTolerance = DistributedVirtualSwitchHostInfrastructureTrafficClass("faultTolerance")
	// vMotion Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassVmotion = DistributedVirtualSwitchHostInfrastructureTrafficClass("vmotion")
	// Virtual Machine Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassVirtualMachine = DistributedVirtualSwitchHostInfrastructureTrafficClass("virtualMachine")
	// iSCSI Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassISCSI = DistributedVirtualSwitchHostInfrastructureTrafficClass("iSCSI")
	// NFS Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassNfs = DistributedVirtualSwitchHostInfrastructureTrafficClass("nfs")
	// vSphere Replication (VR) Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassHbr = DistributedVirtualSwitchHostInfrastructureTrafficClass("hbr")
	// vSphere Storage Area Network Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassVsan = DistributedVirtualSwitchHostInfrastructureTrafficClass("vsan")
	// vSphere Data Protection - Backup Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassVdp = DistributedVirtualSwitchHostInfrastructureTrafficClass("vdp")
	// vSphere Backup NFC Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassBackupNfc = DistributedVirtualSwitchHostInfrastructureTrafficClass("backupNfc")
	// vSphere NVMETCP Traffic
	DistributedVirtualSwitchHostInfrastructureTrafficClassNvmetcp = DistributedVirtualSwitchHostInfrastructureTrafficClass("nvmetcp")
)

func (e DistributedVirtualSwitchHostInfrastructureTrafficClass) Values() []DistributedVirtualSwitchHostInfrastructureTrafficClass {
	return []DistributedVirtualSwitchHostInfrastructureTrafficClass{
		DistributedVirtualSwitchHostInfrastructureTrafficClassManagement,
		DistributedVirtualSwitchHostInfrastructureTrafficClassFaultTolerance,
		DistributedVirtualSwitchHostInfrastructureTrafficClassVmotion,
		DistributedVirtualSwitchHostInfrastructureTrafficClassVirtualMachine,
		DistributedVirtualSwitchHostInfrastructureTrafficClassISCSI,
		DistributedVirtualSwitchHostInfrastructureTrafficClassNfs,
		DistributedVirtualSwitchHostInfrastructureTrafficClassHbr,
		DistributedVirtualSwitchHostInfrastructureTrafficClassVsan,
		DistributedVirtualSwitchHostInfrastructureTrafficClassVdp,
		DistributedVirtualSwitchHostInfrastructureTrafficClassBackupNfc,
		DistributedVirtualSwitchHostInfrastructureTrafficClassNvmetcp,
	}
}

func (e DistributedVirtualSwitchHostInfrastructureTrafficClass) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchHostInfrastructureTrafficClass"] = reflect.TypeOf((*DistributedVirtualSwitchHostInfrastructureTrafficClass)(nil)).Elem()
	minAPIVersionForEnumValue["DistributedVirtualSwitchHostInfrastructureTrafficClass"] = map[string]string{
		"backupNfc": "7.0.1.0",
		"nvmetcp":   "7.0.3.0",
	}
}

// Describes the state of the host proxy switch.
type DistributedVirtualSwitchHostMemberHostComponentState string

const (
	// The host proxy switch is up and running.
	DistributedVirtualSwitchHostMemberHostComponentStateUp = DistributedVirtualSwitchHostMemberHostComponentState("up")
	// The host proxy switch is waiting to be initialized.
	DistributedVirtualSwitchHostMemberHostComponentStatePending = DistributedVirtualSwitchHostMemberHostComponentState("pending")
	// The proxy switch configuration is not the same as the
	// distributed virtual switch configuration in the vCenter Server.
	DistributedVirtualSwitchHostMemberHostComponentStateOutOfSync = DistributedVirtualSwitchHostMemberHostComponentState("outOfSync")
	// The host requires attention.
	DistributedVirtualSwitchHostMemberHostComponentStateWarning = DistributedVirtualSwitchHostMemberHostComponentState("warning")
	// The host is disconnected or it is not responding.
	DistributedVirtualSwitchHostMemberHostComponentStateDisconnected = DistributedVirtualSwitchHostMemberHostComponentState("disconnected")
	// The host proxy is down.
	DistributedVirtualSwitchHostMemberHostComponentStateDown = DistributedVirtualSwitchHostMemberHostComponentState("down")
)

func (e DistributedVirtualSwitchHostMemberHostComponentState) Values() []DistributedVirtualSwitchHostMemberHostComponentState {
	return []DistributedVirtualSwitchHostMemberHostComponentState{
		DistributedVirtualSwitchHostMemberHostComponentStateUp,
		DistributedVirtualSwitchHostMemberHostComponentStatePending,
		DistributedVirtualSwitchHostMemberHostComponentStateOutOfSync,
		DistributedVirtualSwitchHostMemberHostComponentStateWarning,
		DistributedVirtualSwitchHostMemberHostComponentStateDisconnected,
		DistributedVirtualSwitchHostMemberHostComponentStateDown,
	}
}

func (e DistributedVirtualSwitchHostMemberHostComponentState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchHostMemberHostComponentState"] = reflect.TypeOf((*DistributedVirtualSwitchHostMemberHostComponentState)(nil)).Elem()
}

// Describe the runtime state of the uplink.
type DistributedVirtualSwitchHostMemberHostUplinkStateState string

const (
	DistributedVirtualSwitchHostMemberHostUplinkStateStateActive  = DistributedVirtualSwitchHostMemberHostUplinkStateState("active")
	DistributedVirtualSwitchHostMemberHostUplinkStateStateStandby = DistributedVirtualSwitchHostMemberHostUplinkStateState("standby")
)

func (e DistributedVirtualSwitchHostMemberHostUplinkStateState) Values() []DistributedVirtualSwitchHostMemberHostUplinkStateState {
	return []DistributedVirtualSwitchHostMemberHostUplinkStateState{
		DistributedVirtualSwitchHostMemberHostUplinkStateStateActive,
		DistributedVirtualSwitchHostMemberHostUplinkStateStateStandby,
	}
}

func (e DistributedVirtualSwitchHostMemberHostUplinkStateState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchHostMemberHostUplinkStateState"] = reflect.TypeOf((*DistributedVirtualSwitchHostMemberHostUplinkStateState)(nil)).Elem()
}

// Transport zone type.
type DistributedVirtualSwitchHostMemberTransportZoneType string

const (
	// VLAN based networking
	DistributedVirtualSwitchHostMemberTransportZoneTypeVlan = DistributedVirtualSwitchHostMemberTransportZoneType("vlan")
	// VXLAN based networking
	DistributedVirtualSwitchHostMemberTransportZoneTypeOverlay = DistributedVirtualSwitchHostMemberTransportZoneType("overlay")
)

func (e DistributedVirtualSwitchHostMemberTransportZoneType) Values() []DistributedVirtualSwitchHostMemberTransportZoneType {
	return []DistributedVirtualSwitchHostMemberTransportZoneType{
		DistributedVirtualSwitchHostMemberTransportZoneTypeVlan,
		DistributedVirtualSwitchHostMemberTransportZoneTypeOverlay,
	}
}

func (e DistributedVirtualSwitchHostMemberTransportZoneType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchHostMemberTransportZoneType"] = reflect.TypeOf((*DistributedVirtualSwitchHostMemberTransportZoneType)(nil)).Elem()
}

// Network resource control version types.
type DistributedVirtualSwitchNetworkResourceControlVersion string

const (
	// Network Resource Control API version 2
	DistributedVirtualSwitchNetworkResourceControlVersionVersion2 = DistributedVirtualSwitchNetworkResourceControlVersion("version2")
	// Network Resource Control API version 3
	DistributedVirtualSwitchNetworkResourceControlVersionVersion3 = DistributedVirtualSwitchNetworkResourceControlVersion("version3")
)

func (e DistributedVirtualSwitchNetworkResourceControlVersion) Values() []DistributedVirtualSwitchNetworkResourceControlVersion {
	return []DistributedVirtualSwitchNetworkResourceControlVersion{
		DistributedVirtualSwitchNetworkResourceControlVersionVersion2,
		DistributedVirtualSwitchNetworkResourceControlVersionVersion3,
	}
}

func (e DistributedVirtualSwitchNetworkResourceControlVersion) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchNetworkResourceControlVersion"] = reflect.TypeOf((*DistributedVirtualSwitchNetworkResourceControlVersion)(nil)).Elem()
}

// List of possible teaming modes supported by the vNetwork Distributed
// Switch.
//
// The different policy modes define the way traffic is routed
// through the different uplink ports in a team.
type DistributedVirtualSwitchNicTeamingPolicyMode string

const (
	// Routing based on IP hash
	DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_ip = DistributedVirtualSwitchNicTeamingPolicyMode("loadbalance_ip")
	// Route based on source MAC hash
	DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_srcmac = DistributedVirtualSwitchNicTeamingPolicyMode("loadbalance_srcmac")
	// Route based on the source of the port ID
	DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_srcid = DistributedVirtualSwitchNicTeamingPolicyMode("loadbalance_srcid")
	// Use explicit failover order
	DistributedVirtualSwitchNicTeamingPolicyModeFailover_explicit = DistributedVirtualSwitchNicTeamingPolicyMode("failover_explicit")
	// Routing based by dynamically balancing traffic through the NICs
	// in a team.
	//
	// This is the recommended teaming policy when the
	// network I/O control feature is enabled for the vNetwork
	// Distributed Switch.
	DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_loadbased = DistributedVirtualSwitchNicTeamingPolicyMode("loadbalance_loadbased")
)

func (e DistributedVirtualSwitchNicTeamingPolicyMode) Values() []DistributedVirtualSwitchNicTeamingPolicyMode {
	return []DistributedVirtualSwitchNicTeamingPolicyMode{
		DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_ip,
		DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_srcmac,
		DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_srcid,
		DistributedVirtualSwitchNicTeamingPolicyModeFailover_explicit,
		DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_loadbased,
	}
}

func (e DistributedVirtualSwitchNicTeamingPolicyMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchNicTeamingPolicyMode"] = reflect.TypeOf((*DistributedVirtualSwitchNicTeamingPolicyMode)(nil)).Elem()
}

// The connectee types.
type DistributedVirtualSwitchPortConnecteeConnecteeType string

const (
	// The port connects to a Physical NIC.
	DistributedVirtualSwitchPortConnecteeConnecteeTypePnic = DistributedVirtualSwitchPortConnecteeConnecteeType("pnic")
	// The port connects to a Virtual NIC in a Virtual Machine.
	DistributedVirtualSwitchPortConnecteeConnecteeTypeVmVnic = DistributedVirtualSwitchPortConnecteeConnecteeType("vmVnic")
	// The port connects to a console Virtual NIC on a host.
	DistributedVirtualSwitchPortConnecteeConnecteeTypeHostConsoleVnic = DistributedVirtualSwitchPortConnecteeConnecteeType("hostConsoleVnic")
	// The port connects to a VMkernel Virtual NIC on a host.
	DistributedVirtualSwitchPortConnecteeConnecteeTypeHostVmkVnic = DistributedVirtualSwitchPortConnecteeConnecteeType("hostVmkVnic")
	// The port connects to a Virtual NIC in a System CRX VM.
	DistributedVirtualSwitchPortConnecteeConnecteeTypeSystemCrxVnic = DistributedVirtualSwitchPortConnecteeConnecteeType("systemCrxVnic")
)

func (e DistributedVirtualSwitchPortConnecteeConnecteeType) Values() []DistributedVirtualSwitchPortConnecteeConnecteeType {
	return []DistributedVirtualSwitchPortConnecteeConnecteeType{
		DistributedVirtualSwitchPortConnecteeConnecteeTypePnic,
		DistributedVirtualSwitchPortConnecteeConnecteeTypeVmVnic,
		DistributedVirtualSwitchPortConnecteeConnecteeTypeHostConsoleVnic,
		DistributedVirtualSwitchPortConnecteeConnecteeTypeHostVmkVnic,
		DistributedVirtualSwitchPortConnecteeConnecteeTypeSystemCrxVnic,
	}
}

func (e DistributedVirtualSwitchPortConnecteeConnecteeType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchPortConnecteeConnecteeType"] = reflect.TypeOf((*DistributedVirtualSwitchPortConnecteeConnecteeType)(nil)).Elem()
	minAPIVersionForEnumValue["DistributedVirtualSwitchPortConnecteeConnecteeType"] = map[string]string{
		"systemCrxVnic": "8.0.1.0",
	}
}

// The product spec operation types.
type DistributedVirtualSwitchProductSpecOperationType string

const (
	// Push the switch's host component of the specified product info to the
	// host members of the switch at a fixed location known by the host.
	DistributedVirtualSwitchProductSpecOperationTypePreInstall = DistributedVirtualSwitchProductSpecOperationType("preInstall")
	// Change the switch implementation to use the specified one.
	//
	// If the
	// property values in the specified product info are different from
	// those of the corresponding properties in the switch's product info,
	// a host component preinstall and switch upgrade will be performed.
	DistributedVirtualSwitchProductSpecOperationTypeUpgrade = DistributedVirtualSwitchProductSpecOperationType("upgrade")
	// Set the product information for an available switch upgrade that
	// would be done by the switch implementation.
	//
	// This operation will post
	// a config issue on the switch to signal the availability of an upgrade.
	// This operation is applicable only in the case when switch policy
	// `DVSPolicy.autoUpgradeAllowed`
	// is set to false.
	DistributedVirtualSwitchProductSpecOperationTypeNotifyAvailableUpgrade = DistributedVirtualSwitchProductSpecOperationType("notifyAvailableUpgrade")
	// If productSpec is set to be same as that in the
	// `DvsUpgradeAvailableEvent` configIssue, the switch
	// implementation will proceed with the upgrade.
	//
	// To reject or stop the
	// upgrade, leave the productSpec unset. If productSpec is set but does not
	// match that in `DvsUpgradeAvailableEvent` configIssue,
	// a fault will be raised.
	// This operation is applicable only in the case when switch policy
	// `DVSPolicy.autoUpgradeAllowed`
	// is set to false.
	DistributedVirtualSwitchProductSpecOperationTypeProceedWithUpgrade = DistributedVirtualSwitchProductSpecOperationType("proceedWithUpgrade")
	// Update the bundle URL and ID information.
	//
	// If other properties in
	// the specified product info differ from the
	// corresponding properties of the switch's product info, a fault will
	// be thrown. Updating the bundle ID will result in installing the new host
	// component identified by the bundle ID.
	DistributedVirtualSwitchProductSpecOperationTypeUpdateBundleInfo = DistributedVirtualSwitchProductSpecOperationType("updateBundleInfo")
)

func (e DistributedVirtualSwitchProductSpecOperationType) Values() []DistributedVirtualSwitchProductSpecOperationType {
	return []DistributedVirtualSwitchProductSpecOperationType{
		DistributedVirtualSwitchProductSpecOperationTypePreInstall,
		DistributedVirtualSwitchProductSpecOperationTypeUpgrade,
		DistributedVirtualSwitchProductSpecOperationTypeNotifyAvailableUpgrade,
		DistributedVirtualSwitchProductSpecOperationTypeProceedWithUpgrade,
		DistributedVirtualSwitchProductSpecOperationTypeUpdateBundleInfo,
	}
}

func (e DistributedVirtualSwitchProductSpecOperationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DistributedVirtualSwitchProductSpecOperationType"] = reflect.TypeOf((*DistributedVirtualSwitchProductSpecOperationType)(nil)).Elem()
}

type DpmBehavior string

const (
	// Specifies that VirtualCenter should generate recommendations
	// for host power operations, but should not execute the
	// recommendations automatically.
	DpmBehaviorManual = DpmBehavior("manual")
	// Specifies that VirtualCenter should generate recommendations
	// for host power operations, and should execute the
	// recommendations automatically.
	DpmBehaviorAutomated = DpmBehavior("automated")
)

func (e DpmBehavior) Values() []DpmBehavior {
	return []DpmBehavior{
		DpmBehaviorManual,
		DpmBehaviorAutomated,
	}
}

func (e DpmBehavior) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DpmBehavior"] = reflect.TypeOf((*DpmBehavior)(nil)).Elem()
}

type DrsBehavior string

const (
	// Specifies that VirtualCenter should generate recommendations for
	// virtual machine migration and for placement with a host,
	// but should not implement the recommendations automatically.
	DrsBehaviorManual = DrsBehavior("manual")
	// Specifies that VirtualCenter should generate recommendations for
	// virtual machine migration and for placement with a host,
	// but should automatically implement only the placement at power on.
	DrsBehaviorPartiallyAutomated = DrsBehavior("partiallyAutomated")
	// Specifies that VirtualCenter should automate both the migration
	// of virtual machines and their placement with a host at power on.
	DrsBehaviorFullyAutomated = DrsBehavior("fullyAutomated")
)

func (e DrsBehavior) Values() []DrsBehavior {
	return []DrsBehavior{
		DrsBehaviorManual,
		DrsBehaviorPartiallyAutomated,
		DrsBehaviorFullyAutomated,
	}
}

func (e DrsBehavior) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DrsBehavior"] = reflect.TypeOf((*DrsBehavior)(nil)).Elem()
}

// Correlation state as computed by storageRM
// module on host.
type DrsInjectorWorkloadCorrelationState string

const (
	DrsInjectorWorkloadCorrelationStateCorrelated   = DrsInjectorWorkloadCorrelationState("Correlated")
	DrsInjectorWorkloadCorrelationStateUncorrelated = DrsInjectorWorkloadCorrelationState("Uncorrelated")
)

func (e DrsInjectorWorkloadCorrelationState) Values() []DrsInjectorWorkloadCorrelationState {
	return []DrsInjectorWorkloadCorrelationState{
		DrsInjectorWorkloadCorrelationStateCorrelated,
		DrsInjectorWorkloadCorrelationStateUncorrelated,
	}
}

func (e DrsInjectorWorkloadCorrelationState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DrsInjectorWorkloadCorrelationState"] = reflect.TypeOf((*DrsInjectorWorkloadCorrelationState)(nil)).Elem()
}

// Deprecated as of VI API 2.5 use `RecommendationReasonCode_enum`.
//
// List of defined migration reason codes:
type DrsRecommendationReasonCode string

const (
	// Balance average CPU utilization.
	DrsRecommendationReasonCodeFairnessCpuAvg = DrsRecommendationReasonCode("fairnessCpuAvg")
	// Balance average memory utilization.
	DrsRecommendationReasonCodeFairnessMemAvg = DrsRecommendationReasonCode("fairnessMemAvg")
	// Fulfill affinity rule.
	DrsRecommendationReasonCodeJointAffin = DrsRecommendationReasonCode("jointAffin")
	// Fulfill anti-affinity rule.
	DrsRecommendationReasonCodeAntiAffin = DrsRecommendationReasonCode("antiAffin")
	// Host entering maintenance mode.
	DrsRecommendationReasonCodeHostMaint = DrsRecommendationReasonCode("hostMaint")
)

func (e DrsRecommendationReasonCode) Values() []DrsRecommendationReasonCode {
	return []DrsRecommendationReasonCode{
		DrsRecommendationReasonCodeFairnessCpuAvg,
		DrsRecommendationReasonCodeFairnessMemAvg,
		DrsRecommendationReasonCodeJointAffin,
		DrsRecommendationReasonCodeAntiAffin,
		DrsRecommendationReasonCodeHostMaint,
	}
}

func (e DrsRecommendationReasonCode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DrsRecommendationReasonCode"] = reflect.TypeOf((*DrsRecommendationReasonCode)(nil)).Elem()
}

// The port blocked/unblocked state.
type DvsEventPortBlockState string

const (
	// The dvs port is in unset state
	DvsEventPortBlockStateUnset = DvsEventPortBlockState("unset")
	// The dvs port is in blocked state
	DvsEventPortBlockStateBlocked = DvsEventPortBlockState("blocked")
	// The dvs port is in unblocked state
	DvsEventPortBlockStateUnblocked = DvsEventPortBlockState("unblocked")
	// The dvs port is in unknown state
	DvsEventPortBlockStateUnknown = DvsEventPortBlockState("unknown")
)

func (e DvsEventPortBlockState) Values() []DvsEventPortBlockState {
	return []DvsEventPortBlockState{
		DvsEventPortBlockStateUnset,
		DvsEventPortBlockStateBlocked,
		DvsEventPortBlockStateUnblocked,
		DvsEventPortBlockStateUnknown,
	}
}

func (e DvsEventPortBlockState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DvsEventPortBlockState"] = reflect.TypeOf((*DvsEventPortBlockState)(nil)).Elem()
}

// Network Filter on Failure Type.
//
// It specifies whether all the
// packets will be allowed or all the packets will be denied when
// Filter fails to configure.
type DvsFilterOnFailure string

const (
	// Allows all the packets when the Filter fails to configure.
	DvsFilterOnFailureFailOpen = DvsFilterOnFailure("failOpen")
	// Denies all the packets when the Filter fails to configure.
	DvsFilterOnFailureFailClosed = DvsFilterOnFailure("failClosed")
)

func (e DvsFilterOnFailure) Values() []DvsFilterOnFailure {
	return []DvsFilterOnFailure{
		DvsFilterOnFailureFailOpen,
		DvsFilterOnFailureFailClosed,
	}
}

func (e DvsFilterOnFailure) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DvsFilterOnFailure"] = reflect.TypeOf((*DvsFilterOnFailure)(nil)).Elem()
}

// Network Traffic Rule direction types.
//
// It specifies whether rule
// needs to be applied for packets which are incoming/outgoing or both.
type DvsNetworkRuleDirectionType string

const (
	// This specifies that the network rule has to be applied only for
	// incoming packets.
	DvsNetworkRuleDirectionTypeIncomingPackets = DvsNetworkRuleDirectionType("incomingPackets")
	// This specifies that the network rule has to be applied only for
	// outgoing packets.
	DvsNetworkRuleDirectionTypeOutgoingPackets = DvsNetworkRuleDirectionType("outgoingPackets")
	// This specifies that the network rule has to be applied only for
	// both incoming and outgoing packets.
	DvsNetworkRuleDirectionTypeBoth = DvsNetworkRuleDirectionType("both")
)

func (e DvsNetworkRuleDirectionType) Values() []DvsNetworkRuleDirectionType {
	return []DvsNetworkRuleDirectionType{
		DvsNetworkRuleDirectionTypeIncomingPackets,
		DvsNetworkRuleDirectionTypeOutgoingPackets,
		DvsNetworkRuleDirectionTypeBoth,
	}
}

func (e DvsNetworkRuleDirectionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["DvsNetworkRuleDirectionType"] = reflect.TypeOf((*DvsNetworkRuleDirectionType)(nil)).Elem()
}

// The `EntityImportType_enum` enum defines the import type for a
// `DistributedVirtualSwitchManager*.*DistributedVirtualSwitchManager.DVSManagerImportEntity_Task`
// operation.
type EntityImportType string

const (
	// Create the entity with new identifiers.
	//
	// Specify the
	// `EntityBackupConfig*.*EntityBackupConfig.name` and
	// `EntityBackupConfig*.*EntityBackupConfig.container`
	// properties.
	//
	// The Server ignores any value for the
	// `EntityBackupConfig*.*EntityBackupConfig.key`
	// property.
	EntityImportTypeCreateEntityWithNewIdentifier = EntityImportType("createEntityWithNewIdentifier")
	// Recreate the entities with the original identifiers of the entity from which backup was created.
	//
	// The Server throws an exception if an entity with the same identifier already exists.
	// This option will also add the host members to the `DistributedVirtualSwitch` and will
	// try to get the virtual machine networking back with the same `DistributedVirtualPortgroup`.
	// Specify a `Folder` as the
	// `EntityBackupConfig*.*EntityBackupConfig.container`
	// for `EntityBackupConfig*.*EntityBackupConfig.entityType`
	// "distributedVirtualSwitch".
	//
	// The Server ignores any values for the
	// `EntityBackupConfig*.*EntityBackupConfig.key` and
	// `EntityBackupConfig*.*EntityBackupConfig.name`
	// properties.
	EntityImportTypeCreateEntityWithOriginalIdentifier = EntityImportType("createEntityWithOriginalIdentifier")
	// Apply the configuration specified in the
	// `EntityBackupConfig*.*EntityBackupConfig.configBlob`
	// property to the entity specified in the
	// `EntityBackupConfig*.*EntityBackupConfig.entityType` and
	// `EntityBackupConfig*.*EntityBackupConfig.key`
	// properties.
	//
	// If you specify
	// `EntityBackupConfig*.*EntityBackupConfig.name`,
	// the Server uses the specified name to rename the entity.
	//
	// The Server ignores any value for the
	// `EntityBackupConfig*.*EntityBackupConfig.container`
	// property.
	EntityImportTypeApplyToEntitySpecified = EntityImportType("applyToEntitySpecified")
)

func (e EntityImportType) Values() []EntityImportType {
	return []EntityImportType{
		EntityImportTypeCreateEntityWithNewIdentifier,
		EntityImportTypeCreateEntityWithOriginalIdentifier,
		EntityImportTypeApplyToEntitySpecified,
	}
}

func (e EntityImportType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["EntityImportType"] = reflect.TypeOf((*EntityImportType)(nil)).Elem()
}

// The `EntityType_enum` enum identifies
// the type of entity that was exported
// (`DistributedVirtualSwitchManager.DVSManagerExportEntity_Task`).
type EntityType string

const (
	// Indicates the exported entity is a `DistributedVirtualSwitch`.
	EntityTypeDistributedVirtualSwitch = EntityType("distributedVirtualSwitch")
	// Indicates the exported entity is a `DistributedVirtualPortgroup`.
	EntityTypeDistributedVirtualPortgroup = EntityType("distributedVirtualPortgroup")
)

func (e EntityType) Values() []EntityType {
	return []EntityType{
		EntityTypeDistributedVirtualSwitch,
		EntityTypeDistributedVirtualPortgroup,
	}
}

func (e EntityType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["EntityType"] = reflect.TypeOf((*EntityType)(nil)).Elem()
}

// Basic Comparison operators
type EventAlarmExpressionComparisonOperator string

const (
	// attribute equals specified value
	EventAlarmExpressionComparisonOperatorEquals = EventAlarmExpressionComparisonOperator("equals")
	// attribute does not equal specified value
	EventAlarmExpressionComparisonOperatorNotEqualTo = EventAlarmExpressionComparisonOperator("notEqualTo")
	// attribute starts with specified value
	EventAlarmExpressionComparisonOperatorStartsWith = EventAlarmExpressionComparisonOperator("startsWith")
	// attribute does not start with specified value
	EventAlarmExpressionComparisonOperatorDoesNotStartWith = EventAlarmExpressionComparisonOperator("doesNotStartWith")
	// attribute ends with specified value
	EventAlarmExpressionComparisonOperatorEndsWith = EventAlarmExpressionComparisonOperator("endsWith")
	// attribute does not end with specified value
	EventAlarmExpressionComparisonOperatorDoesNotEndWith = EventAlarmExpressionComparisonOperator("doesNotEndWith")
)

func (e EventAlarmExpressionComparisonOperator) Values() []EventAlarmExpressionComparisonOperator {
	return []EventAlarmExpressionComparisonOperator{
		EventAlarmExpressionComparisonOperatorEquals,
		EventAlarmExpressionComparisonOperatorNotEqualTo,
		EventAlarmExpressionComparisonOperatorStartsWith,
		EventAlarmExpressionComparisonOperatorDoesNotStartWith,
		EventAlarmExpressionComparisonOperatorEndsWith,
		EventAlarmExpressionComparisonOperatorDoesNotEndWith,
	}
}

func (e EventAlarmExpressionComparisonOperator) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["EventAlarmExpressionComparisonOperator"] = reflect.TypeOf((*EventAlarmExpressionComparisonOperator)(nil)).Elem()
}

type EventCategory string

const (
	// Returns informational events.
	EventCategoryInfo = EventCategory("info")
	// Returns warning events.
	EventCategoryWarning = EventCategory("warning")
	// Returns error events.
	EventCategoryError = EventCategory("error")
	// Returns events pertaining to users.
	EventCategoryUser = EventCategory("user")
)

func (e EventCategory) Values() []EventCategory {
	return []EventCategory{
		EventCategoryInfo,
		EventCategoryWarning,
		EventCategoryError,
		EventCategoryUser,
	}
}

func (e EventCategory) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["EventCategory"] = reflect.TypeOf((*EventCategory)(nil)).Elem()
}

// Severity level constants.
type EventEventSeverity string

const (
	// Something that must be corrected
	EventEventSeverityError = EventEventSeverity("error")
	// Should be corrected, but the system can continue operating normally
	EventEventSeverityWarning = EventEventSeverity("warning")
	// An informational message
	EventEventSeverityInfo = EventEventSeverity("info")
	// A user-related message
	EventEventSeverityUser = EventEventSeverity("user")
)

func (e EventEventSeverity) Values() []EventEventSeverity {
	return []EventEventSeverity{
		EventEventSeverityError,
		EventEventSeverityWarning,
		EventEventSeverityInfo,
		EventEventSeverityUser,
	}
}

func (e EventEventSeverity) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["EventEventSeverity"] = reflect.TypeOf((*EventEventSeverity)(nil)).Elem()
}

// This option specifies how to select events based on child relationships
// in the inventory hierarchy.
//
// If a managed entity has children, their events
// can be retrieved with this filter option.
type EventFilterSpecRecursionOption string

const (
	// Returns events that pertain only to the specified managed entity,
	// and not its children.
	EventFilterSpecRecursionOptionSelf = EventFilterSpecRecursionOption("self")
	// Returns events pertaining to child entities only.
	//
	// Excludes
	// events pertaining to the specified managed entity itself.
	EventFilterSpecRecursionOptionChildren = EventFilterSpecRecursionOption("children")
	// Returns events pertaining either to the specified managed entity
	// or to its child entities.
	EventFilterSpecRecursionOptionAll = EventFilterSpecRecursionOption("all")
)

func (e EventFilterSpecRecursionOption) Values() []EventFilterSpecRecursionOption {
	return []EventFilterSpecRecursionOption{
		EventFilterSpecRecursionOptionSelf,
		EventFilterSpecRecursionOptionChildren,
		EventFilterSpecRecursionOptionAll,
	}
}

func (e EventFilterSpecRecursionOption) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["EventFilterSpecRecursionOption"] = reflect.TypeOf((*EventFilterSpecRecursionOption)(nil)).Elem()
}

// The operating mode of the adapter.
type FibreChannelPortType string

const (
	FibreChannelPortTypeFabric       = FibreChannelPortType("fabric")
	FibreChannelPortTypeLoop         = FibreChannelPortType("loop")
	FibreChannelPortTypePointToPoint = FibreChannelPortType("pointToPoint")
	FibreChannelPortTypeUnknown      = FibreChannelPortType("unknown")
)

func (e FibreChannelPortType) Values() []FibreChannelPortType {
	return []FibreChannelPortType{
		FibreChannelPortTypeFabric,
		FibreChannelPortTypeLoop,
		FibreChannelPortTypePointToPoint,
		FibreChannelPortTypeUnknown,
	}
}

func (e FibreChannelPortType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["FibreChannelPortType"] = reflect.TypeOf((*FibreChannelPortType)(nil)).Elem()
}

// Status of volume's support for vStorage hardware acceleration.
//
// The ESX Server determines the status based on the capabilities
// of the devices that support the file system volume.
// When a host boots, the support status is unknown.
// As the ESX host attempts hardware-accelerated operations,
// it determines whether the storage device supports hardware
// acceleration and sets the `HostFileSystemMountInfo.vStorageSupport`
// property accordingly.
type FileSystemMountInfoVStorageSupportStatus string

const (
	// Storage device supports hardware acceleration.
	//
	// The ESX host will use the feature to offload certain
	// storage-related operations to the device.
	FileSystemMountInfoVStorageSupportStatusVStorageSupported = FileSystemMountInfoVStorageSupportStatus("vStorageSupported")
	// Storage device does not support hardware acceleration.
	//
	// The ESX host will handle all storage-related operations.
	FileSystemMountInfoVStorageSupportStatusVStorageUnsupported = FileSystemMountInfoVStorageSupportStatus("vStorageUnsupported")
	// Initial support status value.
	FileSystemMountInfoVStorageSupportStatusVStorageUnknown = FileSystemMountInfoVStorageSupportStatus("vStorageUnknown")
)

func (e FileSystemMountInfoVStorageSupportStatus) Values() []FileSystemMountInfoVStorageSupportStatus {
	return []FileSystemMountInfoVStorageSupportStatus{
		FileSystemMountInfoVStorageSupportStatusVStorageSupported,
		FileSystemMountInfoVStorageSupportStatusVStorageUnsupported,
		FileSystemMountInfoVStorageSupportStatusVStorageUnknown,
	}
}

func (e FileSystemMountInfoVStorageSupportStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["FileSystemMountInfoVStorageSupportStatus"] = reflect.TypeOf((*FileSystemMountInfoVStorageSupportStatus)(nil)).Elem()
}

type FolderDesiredHostState string

const (
	// Add host in maintenance mode.
	FolderDesiredHostStateMaintenance = FolderDesiredHostState("maintenance")
	// Add host in non-maintenance mode.
	FolderDesiredHostStateNon_maintenance = FolderDesiredHostState("non_maintenance")
)

func (e FolderDesiredHostState) Values() []FolderDesiredHostState {
	return []FolderDesiredHostState{
		FolderDesiredHostStateMaintenance,
		FolderDesiredHostStateNon_maintenance,
	}
}

func (e FolderDesiredHostState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["FolderDesiredHostState"] = reflect.TypeOf((*FolderDesiredHostState)(nil)).Elem()
}

// HostSelectionType defines how the host was selected
type FtIssuesOnHostHostSelectionType string

const (
	// The host was specified by the user
	FtIssuesOnHostHostSelectionTypeUser = FtIssuesOnHostHostSelectionType("user")
	// The host was selected by Virtual Center
	FtIssuesOnHostHostSelectionTypeVc = FtIssuesOnHostHostSelectionType("vc")
	// The host was selected by DRS
	FtIssuesOnHostHostSelectionTypeDrs = FtIssuesOnHostHostSelectionType("drs")
)

func (e FtIssuesOnHostHostSelectionType) Values() []FtIssuesOnHostHostSelectionType {
	return []FtIssuesOnHostHostSelectionType{
		FtIssuesOnHostHostSelectionTypeUser,
		FtIssuesOnHostHostSelectionTypeVc,
		FtIssuesOnHostHostSelectionTypeDrs,
	}
}

func (e FtIssuesOnHostHostSelectionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["FtIssuesOnHostHostSelectionType"] = reflect.TypeOf((*FtIssuesOnHostHostSelectionType)(nil)).Elem()
}

type GuestFileType string

const (
	// Regular files, and on Posix filesystems, unix domain sockets
	// and devices.
	GuestFileTypeFile = GuestFileType("file")
	// directory
	GuestFileTypeDirectory = GuestFileType("directory")
	// symbolic link
	GuestFileTypeSymlink = GuestFileType("symlink")
)

func (e GuestFileType) Values() []GuestFileType {
	return []GuestFileType{
		GuestFileTypeFile,
		GuestFileTypeDirectory,
		GuestFileTypeSymlink,
	}
}

func (e GuestFileType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestFileType"] = reflect.TypeOf((*GuestFileType)(nil)).Elem()
}

// Application state type.
type GuestInfoAppStateType string

const (
	// The application state wasn't set from the guest by the application agent.
	//
	// This is the default.
	GuestInfoAppStateTypeNone = GuestInfoAppStateType("none")
	// The guest's application agent declared its state as normal and doesn't
	// require any action
	GuestInfoAppStateTypeAppStateOk = GuestInfoAppStateType("appStateOk")
	// Guest's application agent asks for immediate reset
	GuestInfoAppStateTypeAppStateNeedReset = GuestInfoAppStateType("appStateNeedReset")
)

func (e GuestInfoAppStateType) Values() []GuestInfoAppStateType {
	return []GuestInfoAppStateType{
		GuestInfoAppStateTypeNone,
		GuestInfoAppStateTypeAppStateOk,
		GuestInfoAppStateTypeAppStateNeedReset,
	}
}

func (e GuestInfoAppStateType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestInfoAppStateType"] = reflect.TypeOf((*GuestInfoAppStateType)(nil)).Elem()
}

type GuestInfoCustomizationStatus string

const (
	// No guest customizationSpec has been applied for the VM
	GuestInfoCustomizationStatusTOOLSDEPLOYPKG_IDLE = GuestInfoCustomizationStatus("TOOLSDEPLOYPKG_IDLE")
	// The guest customizationSpec has been applied for the VM,
	// but the customization process has not yet started inside the guest OS
	GuestInfoCustomizationStatusTOOLSDEPLOYPKG_PENDING = GuestInfoCustomizationStatus("TOOLSDEPLOYPKG_PENDING")
	// The customization process is currently running inside the guest OS
	GuestInfoCustomizationStatusTOOLSDEPLOYPKG_RUNNING = GuestInfoCustomizationStatus("TOOLSDEPLOYPKG_RUNNING")
	// The customization process has completed successfully inside the
	// guest OS
	GuestInfoCustomizationStatusTOOLSDEPLOYPKG_SUCCEEDED = GuestInfoCustomizationStatus("TOOLSDEPLOYPKG_SUCCEEDED")
	// The customizatio process has failed inside the guest OS
	GuestInfoCustomizationStatusTOOLSDEPLOYPKG_FAILED = GuestInfoCustomizationStatus("TOOLSDEPLOYPKG_FAILED")
)

func (e GuestInfoCustomizationStatus) Values() []GuestInfoCustomizationStatus {
	return []GuestInfoCustomizationStatus{
		GuestInfoCustomizationStatusTOOLSDEPLOYPKG_IDLE,
		GuestInfoCustomizationStatusTOOLSDEPLOYPKG_PENDING,
		GuestInfoCustomizationStatusTOOLSDEPLOYPKG_RUNNING,
		GuestInfoCustomizationStatusTOOLSDEPLOYPKG_SUCCEEDED,
		GuestInfoCustomizationStatusTOOLSDEPLOYPKG_FAILED,
	}
}

func (e GuestInfoCustomizationStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestInfoCustomizationStatus"] = reflect.TypeOf((*GuestInfoCustomizationStatus)(nil)).Elem()
	minAPIVersionForType["GuestInfoCustomizationStatus"] = "7.0.2.0"
}

// Firmware types
type GuestOsDescriptorFirmwareType string

const (
	// BIOS firmware
	GuestOsDescriptorFirmwareTypeBios = GuestOsDescriptorFirmwareType("bios")
	// Extensible Firmware Interface
	GuestOsDescriptorFirmwareTypeEfi = GuestOsDescriptorFirmwareType("efi")
)

func (e GuestOsDescriptorFirmwareType) Values() []GuestOsDescriptorFirmwareType {
	return []GuestOsDescriptorFirmwareType{
		GuestOsDescriptorFirmwareTypeBios,
		GuestOsDescriptorFirmwareTypeEfi,
	}
}

func (e GuestOsDescriptorFirmwareType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestOsDescriptorFirmwareType"] = reflect.TypeOf((*GuestOsDescriptorFirmwareType)(nil)).Elem()
}

// Guest OS support level
type GuestOsDescriptorSupportLevel string

const (
	// This operating system is not supported,
	// but may be supported in the future.
	GuestOsDescriptorSupportLevelExperimental = GuestOsDescriptorSupportLevel("experimental")
	// This operating system is not fully supported,
	// but may have been supported in the past.
	GuestOsDescriptorSupportLevelLegacy = GuestOsDescriptorSupportLevel("legacy")
	// No longer supported.
	GuestOsDescriptorSupportLevelTerminated = GuestOsDescriptorSupportLevel("terminated")
	// Fully supported.
	GuestOsDescriptorSupportLevelSupported = GuestOsDescriptorSupportLevel("supported")
	// This operating system is not supported.
	GuestOsDescriptorSupportLevelUnsupported = GuestOsDescriptorSupportLevel("unsupported")
	// Support for this operating system will be terminated in the future.
	//
	// Please migrate to using a different operating system.
	GuestOsDescriptorSupportLevelDeprecated = GuestOsDescriptorSupportLevel("deprecated")
	// This operating system may not be supported yet,
	// please check VMware compatibility guide.
	GuestOsDescriptorSupportLevelTechPreview = GuestOsDescriptorSupportLevel("techPreview")
)

func (e GuestOsDescriptorSupportLevel) Values() []GuestOsDescriptorSupportLevel {
	return []GuestOsDescriptorSupportLevel{
		GuestOsDescriptorSupportLevelExperimental,
		GuestOsDescriptorSupportLevelLegacy,
		GuestOsDescriptorSupportLevelTerminated,
		GuestOsDescriptorSupportLevelSupported,
		GuestOsDescriptorSupportLevelUnsupported,
		GuestOsDescriptorSupportLevelDeprecated,
		GuestOsDescriptorSupportLevelTechPreview,
	}
}

func (e GuestOsDescriptorSupportLevel) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestOsDescriptorSupportLevel"] = reflect.TypeOf((*GuestOsDescriptorSupportLevel)(nil)).Elem()
}

// End guest quiesce phase error types.
type GuestQuiesceEndGuestQuiesceError string

const (
	// Fail the end phase of guest quiesce creation.
	GuestQuiesceEndGuestQuiesceErrorFailure = GuestQuiesceEndGuestQuiesceError("failure")
)

func (e GuestQuiesceEndGuestQuiesceError) Values() []GuestQuiesceEndGuestQuiesceError {
	return []GuestQuiesceEndGuestQuiesceError{
		GuestQuiesceEndGuestQuiesceErrorFailure,
	}
}

func (e GuestQuiesceEndGuestQuiesceError) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestQuiesceEndGuestQuiesceError"] = reflect.TypeOf((*GuestQuiesceEndGuestQuiesceError)(nil)).Elem()
}

// This describes the bitness (32-bit or 64-bit) of a registry view in a
// Windows OS that supports WOW64.
//
// WOW64 (short for Windows 32-bit on Windows 64-bit) is the x86 emulator
// that allows 32-bit Windows-based applications to run seamlessly on
// 64-bit Windows. Please refer to these MSDN sites for more details:
// http://msdn.microsoft.com/en-us/library/aa384249(v=vs.85).aspx and
// http://msdn.microsoft.com/en-us/library/aa384253(v=vs.85).aspx
type GuestRegKeyWowSpec string

const (
	// Access the key from the native view of the
	// Registry (32-bit on 32-bit versions of Windows,
	// 64-bit on 64-bit versions of Windows).
	GuestRegKeyWowSpecWOWNative = GuestRegKeyWowSpec("WOWNative")
	// Access the key from the 32-bit view of the Registry.
	GuestRegKeyWowSpecWOW32 = GuestRegKeyWowSpec("WOW32")
	// Access the key from the 64-bit view of the Registry.
	GuestRegKeyWowSpecWOW64 = GuestRegKeyWowSpec("WOW64")
)

func (e GuestRegKeyWowSpec) Values() []GuestRegKeyWowSpec {
	return []GuestRegKeyWowSpec{
		GuestRegKeyWowSpecWOWNative,
		GuestRegKeyWowSpecWOW32,
		GuestRegKeyWowSpecWOW64,
	}
}

func (e GuestRegKeyWowSpec) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["GuestRegKeyWowSpec"] = reflect.TypeOf((*GuestRegKeyWowSpec)(nil)).Elem()
}

type HealthUpdateInfoComponentType string

const (
	HealthUpdateInfoComponentTypeMemory  = HealthUpdateInfoComponentType("Memory")
	HealthUpdateInfoComponentTypePower   = HealthUpdateInfoComponentType("Power")
	HealthUpdateInfoComponentTypeFan     = HealthUpdateInfoComponentType("Fan")
	HealthUpdateInfoComponentTypeNetwork = HealthUpdateInfoComponentType("Network")
	HealthUpdateInfoComponentTypeStorage = HealthUpdateInfoComponentType("Storage")
)

func (e HealthUpdateInfoComponentType) Values() []HealthUpdateInfoComponentType {
	return []HealthUpdateInfoComponentType{
		HealthUpdateInfoComponentTypeMemory,
		HealthUpdateInfoComponentTypePower,
		HealthUpdateInfoComponentTypeFan,
		HealthUpdateInfoComponentTypeNetwork,
		HealthUpdateInfoComponentTypeStorage,
	}
}

func (e HealthUpdateInfoComponentType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HealthUpdateInfoComponentType"] = reflect.TypeOf((*HealthUpdateInfoComponentType)(nil)).Elem()
}

// Defines different access modes that a user may have on the host for
// direct host connections.
//
// The assumption here is that when the host is managed by vCenter,
// we don't need fine-grained control on local user permissions like the
// interface provided by `AuthorizationManager`.
type HostAccessMode string

const (
	// Indicates that the user has no explicitly defined permissions or roles.
	//
	// This is used when we want to remove all permissions for some user.
	//
	// Note that this is not the same as `accessNoAccess`.
	HostAccessModeAccessNone = HostAccessMode("accessNone")
	// Describes a propagating Admin role on the root inventory object
	// (root folder) on the host, and no other non-Admin role on any other
	// object.
	//
	// The same permissions are needed to login to local or remote
	// shell (ESXiShell or SSH).
	HostAccessModeAccessAdmin = HostAccessMode("accessAdmin")
	// Describes a propagating NoAccess role on the root inventory object
	// (root folder) on the host, and no other roles.
	//
	// Even if the user has another (redundant) NoAccess role on some other
	// inventory object, then the access mode for this user will be
	// classified as `accessOther`.
	//
	// This mode may be used to restrict a specific user account without
	// restricting the access mode for the group to which the user belongs.
	HostAccessModeAccessNoAccess = HostAccessMode("accessNoAccess")
	// Describes a propagating ReadOnly role on the root inventory object
	// (root folder) on the host, and no other roles.
	//
	// Even if the user has another (redundant) ReadOnly role on some other
	// inventory object, then the access mode for this user will be
	// `accessOther`.
	HostAccessModeAccessReadOnly = HostAccessMode("accessReadOnly")
	// Describes a combination of one or more roles/permissions which are
	// none of the above.
	HostAccessModeAccessOther = HostAccessMode("accessOther")
)

func (e HostAccessMode) Values() []HostAccessMode {
	return []HostAccessMode{
		HostAccessModeAccessNone,
		HostAccessModeAccessAdmin,
		HostAccessModeAccessNoAccess,
		HostAccessModeAccessReadOnly,
		HostAccessModeAccessOther,
	}
}

func (e HostAccessMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostAccessMode"] = reflect.TypeOf((*HostAccessMode)(nil)).Elem()
}

type HostActiveDirectoryAuthenticationCertificateDigest string

const (
	HostActiveDirectoryAuthenticationCertificateDigestSHA1 = HostActiveDirectoryAuthenticationCertificateDigest("SHA1")
)

func (e HostActiveDirectoryAuthenticationCertificateDigest) Values() []HostActiveDirectoryAuthenticationCertificateDigest {
	return []HostActiveDirectoryAuthenticationCertificateDigest{
		HostActiveDirectoryAuthenticationCertificateDigestSHA1,
	}
}

func (e HostActiveDirectoryAuthenticationCertificateDigest) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostActiveDirectoryAuthenticationCertificateDigest"] = reflect.TypeOf((*HostActiveDirectoryAuthenticationCertificateDigest)(nil)).Elem()
}

type HostActiveDirectoryInfoDomainMembershipStatus string

const (
	// The Active Directory integration provider does not support
	// domain trust checks.
	HostActiveDirectoryInfoDomainMembershipStatusUnknown = HostActiveDirectoryInfoDomainMembershipStatus("unknown")
	// No problems with the domain membership.
	HostActiveDirectoryInfoDomainMembershipStatusOk = HostActiveDirectoryInfoDomainMembershipStatus("ok")
	// The host thinks it's part of a domain,
	// but no domain controllers could be reached to confirm.
	HostActiveDirectoryInfoDomainMembershipStatusNoServers = HostActiveDirectoryInfoDomainMembershipStatus("noServers")
	// The client side of the trust relationship is broken.
	HostActiveDirectoryInfoDomainMembershipStatusClientTrustBroken = HostActiveDirectoryInfoDomainMembershipStatus("clientTrustBroken")
	// The server side of the trust relationship is broken
	// (or bad machine password).
	HostActiveDirectoryInfoDomainMembershipStatusServerTrustBroken = HostActiveDirectoryInfoDomainMembershipStatus("serverTrustBroken")
	// Unexpected domain controller responded.
	HostActiveDirectoryInfoDomainMembershipStatusInconsistentTrust = HostActiveDirectoryInfoDomainMembershipStatus("inconsistentTrust")
	// There's some problem with the domain membership.
	HostActiveDirectoryInfoDomainMembershipStatusOtherProblem = HostActiveDirectoryInfoDomainMembershipStatus("otherProblem")
)

func (e HostActiveDirectoryInfoDomainMembershipStatus) Values() []HostActiveDirectoryInfoDomainMembershipStatus {
	return []HostActiveDirectoryInfoDomainMembershipStatus{
		HostActiveDirectoryInfoDomainMembershipStatusUnknown,
		HostActiveDirectoryInfoDomainMembershipStatusOk,
		HostActiveDirectoryInfoDomainMembershipStatusNoServers,
		HostActiveDirectoryInfoDomainMembershipStatusClientTrustBroken,
		HostActiveDirectoryInfoDomainMembershipStatusServerTrustBroken,
		HostActiveDirectoryInfoDomainMembershipStatusInconsistentTrust,
		HostActiveDirectoryInfoDomainMembershipStatusOtherProblem,
	}
}

func (e HostActiveDirectoryInfoDomainMembershipStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostActiveDirectoryInfoDomainMembershipStatus"] = reflect.TypeOf((*HostActiveDirectoryInfoDomainMembershipStatus)(nil)).Elem()
}

type HostBIOSInfoFirmwareType string

const (
	HostBIOSInfoFirmwareTypeBIOS = HostBIOSInfoFirmwareType("BIOS")
	HostBIOSInfoFirmwareTypeUEFI = HostBIOSInfoFirmwareType("UEFI")
)

func (e HostBIOSInfoFirmwareType) Values() []HostBIOSInfoFirmwareType {
	return []HostBIOSInfoFirmwareType{
		HostBIOSInfoFirmwareTypeBIOS,
		HostBIOSInfoFirmwareTypeUEFI,
	}
}

func (e HostBIOSInfoFirmwareType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostBIOSInfoFirmwareType"] = reflect.TypeOf((*HostBIOSInfoFirmwareType)(nil)).Elem()
	minAPIVersionForType["HostBIOSInfoFirmwareType"] = "8.0.2.0"
}

// Deprecated as of vSphere API 7.0, use
// `VmFaultToleranceConfigIssueReasonForIssue_enum`.
//
// Set of possible values for
// `HostCapability.ftCompatibilityIssues`
type HostCapabilityFtUnsupportedReason string

const (
	// No VMotion license
	HostCapabilityFtUnsupportedReasonVMotionNotLicensed = HostCapabilityFtUnsupportedReason("vMotionNotLicensed")
	// VMotion nic is not configured on the host
	HostCapabilityFtUnsupportedReasonMissingVMotionNic = HostCapabilityFtUnsupportedReason("missingVMotionNic")
	// FT logging nic is not configured on the host
	HostCapabilityFtUnsupportedReasonMissingFTLoggingNic = HostCapabilityFtUnsupportedReason("missingFTLoggingNic")
	// Host does not have proper FT license
	HostCapabilityFtUnsupportedReasonFtNotLicensed = HostCapabilityFtUnsupportedReason("ftNotLicensed")
	// Host does not have HA agent running properly
	HostCapabilityFtUnsupportedReasonHaAgentIssue = HostCapabilityFtUnsupportedReason("haAgentIssue")
	// VMware product installed on the host does not support
	// fault tolerance
	HostCapabilityFtUnsupportedReasonUnsupportedProduct = HostCapabilityFtUnsupportedReason("unsupportedProduct")
	// Host CPU does not support hardware virtualization
	HostCapabilityFtUnsupportedReasonCpuHvUnsupported = HostCapabilityFtUnsupportedReason("cpuHvUnsupported")
	// Host CPU does not support hardware MMU virtualization
	HostCapabilityFtUnsupportedReasonCpuHwmmuUnsupported = HostCapabilityFtUnsupportedReason("cpuHwmmuUnsupported")
	// Host CPU is compatible for replay-based FT, but hardware
	// virtualization has been disabled in the BIOS.
	HostCapabilityFtUnsupportedReasonCpuHvDisabled = HostCapabilityFtUnsupportedReason("cpuHvDisabled")
)

func (e HostCapabilityFtUnsupportedReason) Values() []HostCapabilityFtUnsupportedReason {
	return []HostCapabilityFtUnsupportedReason{
		HostCapabilityFtUnsupportedReasonVMotionNotLicensed,
		HostCapabilityFtUnsupportedReasonMissingVMotionNic,
		HostCapabilityFtUnsupportedReasonMissingFTLoggingNic,
		HostCapabilityFtUnsupportedReasonFtNotLicensed,
		HostCapabilityFtUnsupportedReasonHaAgentIssue,
		HostCapabilityFtUnsupportedReasonUnsupportedProduct,
		HostCapabilityFtUnsupportedReasonCpuHvUnsupported,
		HostCapabilityFtUnsupportedReasonCpuHwmmuUnsupported,
		HostCapabilityFtUnsupportedReasonCpuHvDisabled,
	}
}

func (e HostCapabilityFtUnsupportedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCapabilityFtUnsupportedReason"] = reflect.TypeOf((*HostCapabilityFtUnsupportedReason)(nil)).Elem()
}

// Set of VMFS unmap API version.
type HostCapabilityUnmapMethodSupported string

const (
	// only the unmap priority is supported
	HostCapabilityUnmapMethodSupportedPriority = HostCapabilityUnmapMethodSupported("priority")
	// the unmap bandwidth can be set as a fixed value
	HostCapabilityUnmapMethodSupportedFixed = HostCapabilityUnmapMethodSupported("fixed")
	// the unmap bandwidth can be set as a range, where the actual
	// bandwidth will be dynamically throttled by the backened
	HostCapabilityUnmapMethodSupportedDynamic = HostCapabilityUnmapMethodSupported("dynamic")
)

func (e HostCapabilityUnmapMethodSupported) Values() []HostCapabilityUnmapMethodSupported {
	return []HostCapabilityUnmapMethodSupported{
		HostCapabilityUnmapMethodSupportedPriority,
		HostCapabilityUnmapMethodSupportedFixed,
		HostCapabilityUnmapMethodSupportedDynamic,
	}
}

func (e HostCapabilityUnmapMethodSupported) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCapabilityUnmapMethodSupported"] = reflect.TypeOf((*HostCapabilityUnmapMethodSupported)(nil)).Elem()
}

// Set of possible values for `HostCapability.vmDirectPathGen2UnsupportedReason`.
type HostCapabilityVmDirectPathGen2UnsupportedReason string

const (
	// The host software does not support VMDirectPath Gen 2.
	HostCapabilityVmDirectPathGen2UnsupportedReasonHostNptIncompatibleProduct = HostCapabilityVmDirectPathGen2UnsupportedReason("hostNptIncompatibleProduct")
	// The host hardware does not support VMDirectPath Gen 2.
	//
	// Note that
	// this is a general capability for the host and is independent of
	// support by a given physical NIC.
	HostCapabilityVmDirectPathGen2UnsupportedReasonHostNptIncompatibleHardware = HostCapabilityVmDirectPathGen2UnsupportedReason("hostNptIncompatibleHardware")
	// The host is configured to disable VMDirectPath Gen 2.
	HostCapabilityVmDirectPathGen2UnsupportedReasonHostNptDisabled = HostCapabilityVmDirectPathGen2UnsupportedReason("hostNptDisabled")
)

func (e HostCapabilityVmDirectPathGen2UnsupportedReason) Values() []HostCapabilityVmDirectPathGen2UnsupportedReason {
	return []HostCapabilityVmDirectPathGen2UnsupportedReason{
		HostCapabilityVmDirectPathGen2UnsupportedReasonHostNptIncompatibleProduct,
		HostCapabilityVmDirectPathGen2UnsupportedReasonHostNptIncompatibleHardware,
		HostCapabilityVmDirectPathGen2UnsupportedReasonHostNptDisabled,
	}
}

func (e HostCapabilityVmDirectPathGen2UnsupportedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCapabilityVmDirectPathGen2UnsupportedReason"] = reflect.TypeOf((*HostCapabilityVmDirectPathGen2UnsupportedReason)(nil)).Elem()
}

// The status of a given certificate as computed per the soft and the hard
// thresholds in vCenter Server.
//
// There are two different thresholds for the host certificate
// expirations; a soft threshold (which constitutes of two phases) and a
// hard threshold.
//
// Soft Threshold:
//
// Phase One: vCenter Server will publish an event at
// this time to let the user know about the status, but, no alarms or
// warnings are raised.
//
// Phase Two: During this phase, vCenter Server will publish an event and
// indicate the certificate status as expiring in the UI.
//
// Hard Threshold:
//
// vCenter Server will publish an alarm and indicate via the UI that the
// certificate expiration is imminent.
type HostCertificateManagerCertificateInfoCertificateStatus string

const (
	// The certificate status is unknown.
	HostCertificateManagerCertificateInfoCertificateStatusUnknown = HostCertificateManagerCertificateInfoCertificateStatus("unknown")
	// The certificate has expired.
	HostCertificateManagerCertificateInfoCertificateStatusExpired = HostCertificateManagerCertificateInfoCertificateStatus("expired")
	// The certificate is expiring shortly.
	//
	// (soft threshold - 1)
	HostCertificateManagerCertificateInfoCertificateStatusExpiring = HostCertificateManagerCertificateInfoCertificateStatus("expiring")
	// The certificate is expiring shortly.
	//
	// (soft threshold - 2)
	HostCertificateManagerCertificateInfoCertificateStatusExpiringShortly = HostCertificateManagerCertificateInfoCertificateStatus("expiringShortly")
	// The certificate expiration is imminent.
	//
	// (hard threshold)
	HostCertificateManagerCertificateInfoCertificateStatusExpirationImminent = HostCertificateManagerCertificateInfoCertificateStatus("expirationImminent")
	// The certificate is good.
	HostCertificateManagerCertificateInfoCertificateStatusGood = HostCertificateManagerCertificateInfoCertificateStatus("good")
)

func (e HostCertificateManagerCertificateInfoCertificateStatus) Values() []HostCertificateManagerCertificateInfoCertificateStatus {
	return []HostCertificateManagerCertificateInfoCertificateStatus{
		HostCertificateManagerCertificateInfoCertificateStatusUnknown,
		HostCertificateManagerCertificateInfoCertificateStatusExpired,
		HostCertificateManagerCertificateInfoCertificateStatusExpiring,
		HostCertificateManagerCertificateInfoCertificateStatusExpiringShortly,
		HostCertificateManagerCertificateInfoCertificateStatusExpirationImminent,
		HostCertificateManagerCertificateInfoCertificateStatusGood,
	}
}

func (e HostCertificateManagerCertificateInfoCertificateStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCertificateManagerCertificateInfoCertificateStatus"] = reflect.TypeOf((*HostCertificateManagerCertificateInfoCertificateStatus)(nil)).Elem()
}

type HostCertificateManagerCertificateKind string

const (
	// Machine certificate of the Host
	HostCertificateManagerCertificateKindMachine = HostCertificateManagerCertificateKind("Machine")
	// VASA Client certificate used for communication with VASA Provider
	HostCertificateManagerCertificateKindVASAClient = HostCertificateManagerCertificateKind("VASAClient")
)

func (e HostCertificateManagerCertificateKind) Values() []HostCertificateManagerCertificateKind {
	return []HostCertificateManagerCertificateKind{
		HostCertificateManagerCertificateKindMachine,
		HostCertificateManagerCertificateKindVASAClient,
	}
}

func (e HostCertificateManagerCertificateKind) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCertificateManagerCertificateKind"] = reflect.TypeOf((*HostCertificateManagerCertificateKind)(nil)).Elem()
	minAPIVersionForType["HostCertificateManagerCertificateKind"] = "8.0.1.0"
}

// This is a global mode on a configuration specification indicating
// whether the structure represents the desired state or the set of
// operations to apply on the managed object.
type HostConfigChangeMode string

const (
	// Indicates that the structure represents the
	// set of operations to apply on the managed object.
	HostConfigChangeModeModify = HostConfigChangeMode("modify")
	// Indicates that the structure represents the
	// desired state of the managed object.
	HostConfigChangeModeReplace = HostConfigChangeMode("replace")
)

func (e HostConfigChangeMode) Values() []HostConfigChangeMode {
	return []HostConfigChangeMode{
		HostConfigChangeModeModify,
		HostConfigChangeModeReplace,
	}
}

func (e HostConfigChangeMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostConfigChangeMode"] = reflect.TypeOf((*HostConfigChangeMode)(nil)).Elem()
}

// This list indicates the operation that should be performed for an
// entity.
type HostConfigChangeOperation string

const (
	// Indicates the addition of an entity to the configuration.
	HostConfigChangeOperationAdd = HostConfigChangeOperation("add")
	// Indicates the removal of an entity from the configuration.
	HostConfigChangeOperationRemove = HostConfigChangeOperation("remove")
	// Indicates changes on the entity.
	//
	// The entity must exist or a
	// `NotFound` error will be thrown.
	HostConfigChangeOperationEdit = HostConfigChangeOperation("edit")
	// Indicates that an entity will be ignored: it won't be added when it
	// doesn't exist, or removed/changed when it exists.
	HostConfigChangeOperationIgnore = HostConfigChangeOperation("ignore")
)

func (e HostConfigChangeOperation) Values() []HostConfigChangeOperation {
	return []HostConfigChangeOperation{
		HostConfigChangeOperationAdd,
		HostConfigChangeOperationRemove,
		HostConfigChangeOperationEdit,
		HostConfigChangeOperationIgnore,
	}
}

func (e HostConfigChangeOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostConfigChangeOperation"] = reflect.TypeOf((*HostConfigChangeOperation)(nil)).Elem()
}

type HostCpuPackageVendor string

const (
	HostCpuPackageVendorUnknown = HostCpuPackageVendor("unknown")
	HostCpuPackageVendorIntel   = HostCpuPackageVendor("intel")
	HostCpuPackageVendorAmd     = HostCpuPackageVendor("amd")
	HostCpuPackageVendorHygon   = HostCpuPackageVendor("hygon")
)

func (e HostCpuPackageVendor) Values() []HostCpuPackageVendor {
	return []HostCpuPackageVendor{
		HostCpuPackageVendorUnknown,
		HostCpuPackageVendorIntel,
		HostCpuPackageVendorAmd,
		HostCpuPackageVendorHygon,
	}
}

func (e HostCpuPackageVendor) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCpuPackageVendor"] = reflect.TypeOf((*HostCpuPackageVendor)(nil)).Elem()
}

// Possible values for Current CPU power management policy
type HostCpuPowerManagementInfoPolicyType string

const (
	HostCpuPowerManagementInfoPolicyTypeOff           = HostCpuPowerManagementInfoPolicyType("off")
	HostCpuPowerManagementInfoPolicyTypeStaticPolicy  = HostCpuPowerManagementInfoPolicyType("staticPolicy")
	HostCpuPowerManagementInfoPolicyTypeDynamicPolicy = HostCpuPowerManagementInfoPolicyType("dynamicPolicy")
)

func (e HostCpuPowerManagementInfoPolicyType) Values() []HostCpuPowerManagementInfoPolicyType {
	return []HostCpuPowerManagementInfoPolicyType{
		HostCpuPowerManagementInfoPolicyTypeOff,
		HostCpuPowerManagementInfoPolicyTypeStaticPolicy,
		HostCpuPowerManagementInfoPolicyTypeDynamicPolicy,
	}
}

func (e HostCpuPowerManagementInfoPolicyType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCpuPowerManagementInfoPolicyType"] = reflect.TypeOf((*HostCpuPowerManagementInfoPolicyType)(nil)).Elem()
}

type HostCpuSchedulerInfoCpuSchedulerPolicyInfo string

const (
	// The CPU scheduler on this host is running without any modifications
	// or mitigations.
	HostCpuSchedulerInfoCpuSchedulerPolicyInfoSystemDefault = HostCpuSchedulerInfoCpuSchedulerPolicyInfo("systemDefault")
	// The CPU scheduler on this host is using only one hyperthread per
	// core to mitigate a security vulnerability.
	HostCpuSchedulerInfoCpuSchedulerPolicyInfoScav1 = HostCpuSchedulerInfoCpuSchedulerPolicyInfo("scav1")
	// The CPU scheduler on this host is using hyperthreads, with
	// Side-Channel aware scheduling to mitigate a security vulnerability.
	HostCpuSchedulerInfoCpuSchedulerPolicyInfoScav2 = HostCpuSchedulerInfoCpuSchedulerPolicyInfo("scav2")
)

func (e HostCpuSchedulerInfoCpuSchedulerPolicyInfo) Values() []HostCpuSchedulerInfoCpuSchedulerPolicyInfo {
	return []HostCpuSchedulerInfoCpuSchedulerPolicyInfo{
		HostCpuSchedulerInfoCpuSchedulerPolicyInfoSystemDefault,
		HostCpuSchedulerInfoCpuSchedulerPolicyInfoScav1,
		HostCpuSchedulerInfoCpuSchedulerPolicyInfoScav2,
	}
}

func (e HostCpuSchedulerInfoCpuSchedulerPolicyInfo) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCpuSchedulerInfoCpuSchedulerPolicyInfo"] = reflect.TypeOf((*HostCpuSchedulerInfoCpuSchedulerPolicyInfo)(nil)).Elem()
	minAPIVersionForType["HostCpuSchedulerInfoCpuSchedulerPolicyInfo"] = "8.0.3.0"
}

// Defines a host's encryption state
type HostCryptoState string

const (
	// The host is not safe for receiving sensitive material.
	HostCryptoStateIncapable = HostCryptoState("incapable")
	// The host is prepared for receiving sensitive material
	// but does not have a host key set yet.
	HostCryptoStatePrepared = HostCryptoState("prepared")
	// The host is crypto safe and has a host key set.
	HostCryptoStateSafe = HostCryptoState("safe")
	// The host is explicitly crypto disabled and pending reboot to be
	// applied.
	//
	// When host is in this state, creating encrypted virtual
	// machines is not allowed, but still need a reboot to totally clean
	// up and enter incapable state.
	HostCryptoStatePendingIncapable = HostCryptoState("pendingIncapable")
)

func (e HostCryptoState) Values() []HostCryptoState {
	return []HostCryptoState{
		HostCryptoStateIncapable,
		HostCryptoStatePrepared,
		HostCryptoStateSafe,
		HostCryptoStatePendingIncapable,
	}
}

func (e HostCryptoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostCryptoState"] = reflect.TypeOf((*HostCryptoState)(nil)).Elem()
}

type HostDVSConfigSpecSwitchMode string

const (
	// traditional package processing mode.
	HostDVSConfigSpecSwitchModeNormal = HostDVSConfigSpecSwitchMode("normal")
	// ENS mode which skips packet parsing and flow table lookup.
	HostDVSConfigSpecSwitchModeMux = HostDVSConfigSpecSwitchMode("mux")
)

func (e HostDVSConfigSpecSwitchMode) Values() []HostDVSConfigSpecSwitchMode {
	return []HostDVSConfigSpecSwitchMode{
		HostDVSConfigSpecSwitchModeNormal,
		HostDVSConfigSpecSwitchModeMux,
	}
}

func (e HostDVSConfigSpecSwitchMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDVSConfigSpecSwitchMode"] = reflect.TypeOf((*HostDVSConfigSpecSwitchMode)(nil)).Elem()
	minAPIVersionForType["HostDVSConfigSpecSwitchMode"] = "8.0.0.1"
}

type HostDasErrorEventHostDasErrorReason string

const (
	// Error while configuring/unconfiguring HA
	HostDasErrorEventHostDasErrorReasonConfigFailed = HostDasErrorEventHostDasErrorReason("configFailed")
	// Timeout while communicating with HA agent
	HostDasErrorEventHostDasErrorReasonTimeout = HostDasErrorEventHostDasErrorReason("timeout")
	// HA communication initialization failed
	HostDasErrorEventHostDasErrorReasonCommunicationInitFailed = HostDasErrorEventHostDasErrorReason("communicationInitFailed")
	// Health check script failed
	HostDasErrorEventHostDasErrorReasonHealthCheckScriptFailed = HostDasErrorEventHostDasErrorReason("healthCheckScriptFailed")
	// HA agent has an error
	HostDasErrorEventHostDasErrorReasonAgentFailed = HostDasErrorEventHostDasErrorReason("agentFailed")
	// HA agent was shutdown
	HostDasErrorEventHostDasErrorReasonAgentShutdown = HostDasErrorEventHostDasErrorReason("agentShutdown")
	// HA isolation address unpingable
	HostDasErrorEventHostDasErrorReasonIsolationAddressUnpingable = HostDasErrorEventHostDasErrorReason("isolationAddressUnpingable")
	// Other reason
	HostDasErrorEventHostDasErrorReasonOther = HostDasErrorEventHostDasErrorReason("other")
)

func (e HostDasErrorEventHostDasErrorReason) Values() []HostDasErrorEventHostDasErrorReason {
	return []HostDasErrorEventHostDasErrorReason{
		HostDasErrorEventHostDasErrorReasonConfigFailed,
		HostDasErrorEventHostDasErrorReasonTimeout,
		HostDasErrorEventHostDasErrorReasonCommunicationInitFailed,
		HostDasErrorEventHostDasErrorReasonHealthCheckScriptFailed,
		HostDasErrorEventHostDasErrorReasonAgentFailed,
		HostDasErrorEventHostDasErrorReasonAgentShutdown,
		HostDasErrorEventHostDasErrorReasonIsolationAddressUnpingable,
		HostDasErrorEventHostDasErrorReasonOther,
	}
}

func (e HostDasErrorEventHostDasErrorReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDasErrorEventHostDasErrorReason"] = reflect.TypeOf((*HostDasErrorEventHostDasErrorReason)(nil)).Elem()
}

// Types of time synchronization protocols.
type HostDateTimeInfoProtocol string

const (
	// Network Time Protocol (NTP).
	HostDateTimeInfoProtocolNtp = HostDateTimeInfoProtocol("ntp")
	// Precision Time Protocol (PTP).
	HostDateTimeInfoProtocolPtp = HostDateTimeInfoProtocol("ptp")
)

func (e HostDateTimeInfoProtocol) Values() []HostDateTimeInfoProtocol {
	return []HostDateTimeInfoProtocol{
		HostDateTimeInfoProtocolNtp,
		HostDateTimeInfoProtocolPtp,
	}
}

func (e HostDateTimeInfoProtocol) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDateTimeInfoProtocol"] = reflect.TypeOf((*HostDateTimeInfoProtocol)(nil)).Elem()
}

// The set of digest methods that can be used by TPM to calculate the PCR
// values.
type HostDigestInfoDigestMethodType string

const (
	HostDigestInfoDigestMethodTypeSHA1 = HostDigestInfoDigestMethodType("SHA1")
	// Deprecated as of vSphere API 6.7.
	//
	// MD5.
	HostDigestInfoDigestMethodTypeMD5     = HostDigestInfoDigestMethodType("MD5")
	HostDigestInfoDigestMethodTypeSHA256  = HostDigestInfoDigestMethodType("SHA256")
	HostDigestInfoDigestMethodTypeSHA384  = HostDigestInfoDigestMethodType("SHA384")
	HostDigestInfoDigestMethodTypeSHA512  = HostDigestInfoDigestMethodType("SHA512")
	HostDigestInfoDigestMethodTypeSM3_256 = HostDigestInfoDigestMethodType("SM3_256")
)

func (e HostDigestInfoDigestMethodType) Values() []HostDigestInfoDigestMethodType {
	return []HostDigestInfoDigestMethodType{
		HostDigestInfoDigestMethodTypeSHA1,
		HostDigestInfoDigestMethodTypeMD5,
		HostDigestInfoDigestMethodTypeSHA256,
		HostDigestInfoDigestMethodTypeSHA384,
		HostDigestInfoDigestMethodTypeSHA512,
		HostDigestInfoDigestMethodTypeSM3_256,
	}
}

func (e HostDigestInfoDigestMethodType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDigestInfoDigestMethodType"] = reflect.TypeOf((*HostDigestInfoDigestMethodType)(nil)).Elem()
}

// This enum specifies the supported digest verification settings.
//
// For NVMe over TCP connections, both header and data digests may be
// requested during the process of establishing the connection.
// For details, see:
//   - NVM Express Technical Proposal 8000 - NVMe/TCP Transport,
type HostDigestVerificationSetting string

const (
	// Both header and data digest verification are disabled.
	HostDigestVerificationSettingDigestDisabled = HostDigestVerificationSetting("digestDisabled")
	// Only header digest verification is enabled.
	HostDigestVerificationSettingHeaderOnly = HostDigestVerificationSetting("headerOnly")
	// Only data digest verification is enabled.
	HostDigestVerificationSettingDataOnly = HostDigestVerificationSetting("dataOnly")
	// Both header and data digest verification are enabled.
	HostDigestVerificationSettingHeaderAndData = HostDigestVerificationSetting("headerAndData")
)

func (e HostDigestVerificationSetting) Values() []HostDigestVerificationSetting {
	return []HostDigestVerificationSetting{
		HostDigestVerificationSettingDigestDisabled,
		HostDigestVerificationSettingHeaderOnly,
		HostDigestVerificationSettingDataOnly,
		HostDigestVerificationSettingHeaderAndData,
	}
}

func (e HostDigestVerificationSetting) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDigestVerificationSetting"] = reflect.TypeOf((*HostDigestVerificationSetting)(nil)).Elem()
	minAPIVersionForType["HostDigestVerificationSetting"] = "7.0.3.0"
}

type HostDisconnectedEventReasonCode string

const (
	// Failed to verify SSL thumbprint
	HostDisconnectedEventReasonCodeSslThumbprintVerifyFailed = HostDisconnectedEventReasonCode("sslThumbprintVerifyFailed")
	// License expired for the host
	HostDisconnectedEventReasonCodeLicenseExpired = HostDisconnectedEventReasonCode("licenseExpired")
	// Agent is being upgraded
	HostDisconnectedEventReasonCodeAgentUpgrade = HostDisconnectedEventReasonCode("agentUpgrade")
	// User requested disconnect
	HostDisconnectedEventReasonCodeUserRequest = HostDisconnectedEventReasonCode("userRequest")
	// License not available after host upgrade
	HostDisconnectedEventReasonCodeInsufficientLicenses = HostDisconnectedEventReasonCode("insufficientLicenses")
	// Agent is out of date
	HostDisconnectedEventReasonCodeAgentOutOfDate = HostDisconnectedEventReasonCode("agentOutOfDate")
	// Failed to decrypt password
	HostDisconnectedEventReasonCodePasswordDecryptFailure = HostDisconnectedEventReasonCode("passwordDecryptFailure")
	// Unknown reason
	HostDisconnectedEventReasonCodeUnknown = HostDisconnectedEventReasonCode("unknown")
	// The vRAM capacity of vCenter will be exceeded
	HostDisconnectedEventReasonCodeVcVRAMCapacityExceeded = HostDisconnectedEventReasonCode("vcVRAMCapacityExceeded")
)

func (e HostDisconnectedEventReasonCode) Values() []HostDisconnectedEventReasonCode {
	return []HostDisconnectedEventReasonCode{
		HostDisconnectedEventReasonCodeSslThumbprintVerifyFailed,
		HostDisconnectedEventReasonCodeLicenseExpired,
		HostDisconnectedEventReasonCodeAgentUpgrade,
		HostDisconnectedEventReasonCodeUserRequest,
		HostDisconnectedEventReasonCodeInsufficientLicenses,
		HostDisconnectedEventReasonCodeAgentOutOfDate,
		HostDisconnectedEventReasonCodePasswordDecryptFailure,
		HostDisconnectedEventReasonCodeUnknown,
		HostDisconnectedEventReasonCodeVcVRAMCapacityExceeded,
	}
}

func (e HostDisconnectedEventReasonCode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDisconnectedEventReasonCode"] = reflect.TypeOf((*HostDisconnectedEventReasonCode)(nil)).Elem()
}

// List of partition format types.
//
// This denotes the partition table layout.
type HostDiskPartitionInfoPartitionFormat string

const (
	HostDiskPartitionInfoPartitionFormatGpt     = HostDiskPartitionInfoPartitionFormat("gpt")
	HostDiskPartitionInfoPartitionFormatMbr     = HostDiskPartitionInfoPartitionFormat("mbr")
	HostDiskPartitionInfoPartitionFormatUnknown = HostDiskPartitionInfoPartitionFormat("unknown")
)

func (e HostDiskPartitionInfoPartitionFormat) Values() []HostDiskPartitionInfoPartitionFormat {
	return []HostDiskPartitionInfoPartitionFormat{
		HostDiskPartitionInfoPartitionFormatGpt,
		HostDiskPartitionInfoPartitionFormatMbr,
		HostDiskPartitionInfoPartitionFormatUnknown,
	}
}

func (e HostDiskPartitionInfoPartitionFormat) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDiskPartitionInfoPartitionFormat"] = reflect.TypeOf((*HostDiskPartitionInfoPartitionFormat)(nil)).Elem()
}

// List of symbol partition types
type HostDiskPartitionInfoType string

const (
	HostDiskPartitionInfoTypeNone          = HostDiskPartitionInfoType("none")
	HostDiskPartitionInfoTypeVmfs          = HostDiskPartitionInfoType("vmfs")
	HostDiskPartitionInfoTypeLinuxNative   = HostDiskPartitionInfoType("linuxNative")
	HostDiskPartitionInfoTypeLinuxSwap     = HostDiskPartitionInfoType("linuxSwap")
	HostDiskPartitionInfoTypeExtended      = HostDiskPartitionInfoType("extended")
	HostDiskPartitionInfoTypeNtfs          = HostDiskPartitionInfoType("ntfs")
	HostDiskPartitionInfoTypeVmkDiagnostic = HostDiskPartitionInfoType("vmkDiagnostic")
	HostDiskPartitionInfoTypeVffs          = HostDiskPartitionInfoType("vffs")
)

func (e HostDiskPartitionInfoType) Values() []HostDiskPartitionInfoType {
	return []HostDiskPartitionInfoType{
		HostDiskPartitionInfoTypeNone,
		HostDiskPartitionInfoTypeVmfs,
		HostDiskPartitionInfoTypeLinuxNative,
		HostDiskPartitionInfoTypeLinuxSwap,
		HostDiskPartitionInfoTypeExtended,
		HostDiskPartitionInfoTypeNtfs,
		HostDiskPartitionInfoTypeVmkDiagnostic,
		HostDiskPartitionInfoTypeVffs,
	}
}

func (e HostDiskPartitionInfoType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDiskPartitionInfoType"] = reflect.TypeOf((*HostDiskPartitionInfoType)(nil)).Elem()
}

type HostDistributedVirtualSwitchManagerFailoverReason string

const (
	// The failover is caused by DPU crash.
	HostDistributedVirtualSwitchManagerFailoverReasonCrash = HostDistributedVirtualSwitchManagerFailoverReason("crash")
	// The failover is caused by DPU's vmnic(s) link down.
	HostDistributedVirtualSwitchManagerFailoverReasonLinkDown = HostDistributedVirtualSwitchManagerFailoverReason("linkDown")
	// The failover is triggered by the user.
	HostDistributedVirtualSwitchManagerFailoverReasonUserInitiated = HostDistributedVirtualSwitchManagerFailoverReason("userInitiated")
)

func (e HostDistributedVirtualSwitchManagerFailoverReason) Values() []HostDistributedVirtualSwitchManagerFailoverReason {
	return []HostDistributedVirtualSwitchManagerFailoverReason{
		HostDistributedVirtualSwitchManagerFailoverReasonCrash,
		HostDistributedVirtualSwitchManagerFailoverReasonLinkDown,
		HostDistributedVirtualSwitchManagerFailoverReasonUserInitiated,
	}
}

func (e HostDistributedVirtualSwitchManagerFailoverReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDistributedVirtualSwitchManagerFailoverReason"] = reflect.TypeOf((*HostDistributedVirtualSwitchManagerFailoverReason)(nil)).Elem()
	minAPIVersionForType["HostDistributedVirtualSwitchManagerFailoverReason"] = "8.0.3.0"
}

type HostDistributedVirtualSwitchManagerFailoverStage string

const (
	HostDistributedVirtualSwitchManagerFailoverStageSTAGE_1 = HostDistributedVirtualSwitchManagerFailoverStage("STAGE_1")
)

func (e HostDistributedVirtualSwitchManagerFailoverStage) Values() []HostDistributedVirtualSwitchManagerFailoverStage {
	return []HostDistributedVirtualSwitchManagerFailoverStage{
		HostDistributedVirtualSwitchManagerFailoverStageSTAGE_1,
	}
}

func (e HostDistributedVirtualSwitchManagerFailoverStage) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostDistributedVirtualSwitchManagerFailoverStage"] = reflect.TypeOf((*HostDistributedVirtualSwitchManagerFailoverStage)(nil)).Elem()
	minAPIVersionForType["HostDistributedVirtualSwitchManagerFailoverStage"] = "8.0.3.0"
}

// Set of possible values for
// `HostFeatureVersionInfo.key`, which
// is a unique key that identifies a feature.
type HostFeatureVersionKey string

const (
	// VMware Fault Tolerance feature.
	//
	// For pre-4.1 hosts, the
	// version value reported will be empty in which case
	// `AboutInfo.build` should be used. For all
	// other hosts, the version number reported will be a component-specific
	// version identifier of the form X.Y.Z, where:
	// X refers to host agent Fault Tolerance version number,
	// Y refers to VMX Fault Tolerance version number,
	// Z refers to VMkernal Fault Tolerance version
	HostFeatureVersionKeyFaultTolerance = HostFeatureVersionKey("faultTolerance")
)

func (e HostFeatureVersionKey) Values() []HostFeatureVersionKey {
	return []HostFeatureVersionKey{
		HostFeatureVersionKeyFaultTolerance,
	}
}

func (e HostFeatureVersionKey) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFeatureVersionKey"] = reflect.TypeOf((*HostFeatureVersionKey)(nil)).Elem()
}

// Type of file system volume.
type HostFileSystemVolumeFileSystemType string

const (
	// VMware File System (ESX Server only).
	//
	// If this is set,
	// the type of the file system volume is VMFS.
	HostFileSystemVolumeFileSystemTypeVMFS = HostFileSystemVolumeFileSystemType("VMFS")
	// Network file system v3 linux &amp; esx servers only.
	//
	// If this is
	// set, the type of the file system volume is NFS v3.
	HostFileSystemVolumeFileSystemTypeNFS = HostFileSystemVolumeFileSystemType("NFS")
	// Network file system v4.1 linux &amp; esx servers only.
	//
	// If this is
	// set, the type of the file system volume is NFS v4.1 or later.
	HostFileSystemVolumeFileSystemTypeNFS41 = HostFileSystemVolumeFileSystemType("NFS41")
	// Common Internet File System.
	//
	// If this is set, the type of the
	// file system volume is Common Internet File System.
	HostFileSystemVolumeFileSystemTypeCIFS = HostFileSystemVolumeFileSystemType("CIFS")
	// VSAN File System (ESX Server only).
	HostFileSystemVolumeFileSystemTypeVsan = HostFileSystemVolumeFileSystemType("vsan")
	// vFlash File System (ESX Server only).
	//
	// If this is set, the type of the file system volume is VFFS.
	HostFileSystemVolumeFileSystemTypeVFFS = HostFileSystemVolumeFileSystemType("VFFS")
	// vvol File System (ESX Server only).
	HostFileSystemVolumeFileSystemTypeVVOL = HostFileSystemVolumeFileSystemType("VVOL")
	// Persistent Memory File System (ESX Server only).
	HostFileSystemVolumeFileSystemTypePMEM = HostFileSystemVolumeFileSystemType("PMEM")
	// VSAN direct file system.
	HostFileSystemVolumeFileSystemTypeVsanD = HostFileSystemVolumeFileSystemType("vsanD")
	// Used if the file system is not one of the specified file systems.
	//
	// Used mostly for reporting purposes. The other types are described
	// by the otherType property.
	HostFileSystemVolumeFileSystemTypeOTHER = HostFileSystemVolumeFileSystemType("OTHER")
)

func (e HostFileSystemVolumeFileSystemType) Values() []HostFileSystemVolumeFileSystemType {
	return []HostFileSystemVolumeFileSystemType{
		HostFileSystemVolumeFileSystemTypeVMFS,
		HostFileSystemVolumeFileSystemTypeNFS,
		HostFileSystemVolumeFileSystemTypeNFS41,
		HostFileSystemVolumeFileSystemTypeCIFS,
		HostFileSystemVolumeFileSystemTypeVsan,
		HostFileSystemVolumeFileSystemTypeVFFS,
		HostFileSystemVolumeFileSystemTypeVVOL,
		HostFileSystemVolumeFileSystemTypePMEM,
		HostFileSystemVolumeFileSystemTypeVsanD,
		HostFileSystemVolumeFileSystemTypeOTHER,
	}
}

func (e HostFileSystemVolumeFileSystemType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFileSystemVolumeFileSystemType"] = reflect.TypeOf((*HostFileSystemVolumeFileSystemType)(nil)).Elem()
	minAPIVersionForEnumValue["HostFileSystemVolumeFileSystemType"] = map[string]string{
		"vsanD": "7.0.1.0",
	}
}

// Enumeration of port directions.
type HostFirewallRuleDirection string

const (
	HostFirewallRuleDirectionInbound  = HostFirewallRuleDirection("inbound")
	HostFirewallRuleDirectionOutbound = HostFirewallRuleDirection("outbound")
)

func (e HostFirewallRuleDirection) Values() []HostFirewallRuleDirection {
	return []HostFirewallRuleDirection{
		HostFirewallRuleDirectionInbound,
		HostFirewallRuleDirectionOutbound,
	}
}

func (e HostFirewallRuleDirection) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFirewallRuleDirection"] = reflect.TypeOf((*HostFirewallRuleDirection)(nil)).Elem()
}

// Enumeration of port types.
type HostFirewallRulePortType string

const (
	HostFirewallRulePortTypeSrc = HostFirewallRulePortType("src")
	HostFirewallRulePortTypeDst = HostFirewallRulePortType("dst")
)

func (e HostFirewallRulePortType) Values() []HostFirewallRulePortType {
	return []HostFirewallRulePortType{
		HostFirewallRulePortTypeSrc,
		HostFirewallRulePortTypeDst,
	}
}

func (e HostFirewallRulePortType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFirewallRulePortType"] = reflect.TypeOf((*HostFirewallRulePortType)(nil)).Elem()
}

// Set of valid port protocols.
type HostFirewallRuleProtocol string

const (
	HostFirewallRuleProtocolTcp = HostFirewallRuleProtocol("tcp")
	HostFirewallRuleProtocolUdp = HostFirewallRuleProtocol("udp")
)

func (e HostFirewallRuleProtocol) Values() []HostFirewallRuleProtocol {
	return []HostFirewallRuleProtocol{
		HostFirewallRuleProtocolTcp,
		HostFirewallRuleProtocolUdp,
	}
}

func (e HostFirewallRuleProtocol) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFirewallRuleProtocol"] = reflect.TypeOf((*HostFirewallRuleProtocol)(nil)).Elem()
}

type HostFirewallSystemRuleSetId string

const (
	HostFirewallSystemRuleSetIdFaultTolerance = HostFirewallSystemRuleSetId("faultTolerance")
	HostFirewallSystemRuleSetIdFdm            = HostFirewallSystemRuleSetId("fdm")
	HostFirewallSystemRuleSetIdUpdateManager  = HostFirewallSystemRuleSetId("updateManager")
	HostFirewallSystemRuleSetIdVpxHeartbeats  = HostFirewallSystemRuleSetId("vpxHeartbeats")
)

func (e HostFirewallSystemRuleSetId) Values() []HostFirewallSystemRuleSetId {
	return []HostFirewallSystemRuleSetId{
		HostFirewallSystemRuleSetIdFaultTolerance,
		HostFirewallSystemRuleSetIdFdm,
		HostFirewallSystemRuleSetIdUpdateManager,
		HostFirewallSystemRuleSetIdVpxHeartbeats,
	}
}

func (e HostFirewallSystemRuleSetId) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFirewallSystemRuleSetId"] = reflect.TypeOf((*HostFirewallSystemRuleSetId)(nil)).Elem()
	minAPIVersionForType["HostFirewallSystemRuleSetId"] = "8.0.2.0"
}

type HostFirewallSystemServiceName string

const (
	HostFirewallSystemServiceNameVpxa = HostFirewallSystemServiceName("vpxa")
)

func (e HostFirewallSystemServiceName) Values() []HostFirewallSystemServiceName {
	return []HostFirewallSystemServiceName{
		HostFirewallSystemServiceNameVpxa,
	}
}

func (e HostFirewallSystemServiceName) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFirewallSystemServiceName"] = reflect.TypeOf((*HostFirewallSystemServiceName)(nil)).Elem()
	minAPIVersionForType["HostFirewallSystemServiceName"] = "8.0.2.0"
}

// The vendor definition for type of Field Replaceable Unit (FRU).
type HostFruFruType string

const (
	HostFruFruTypeUndefined = HostFruFruType("undefined")
	HostFruFruTypeBoard     = HostFruFruType("board")
	HostFruFruTypeProduct   = HostFruFruType("product")
)

func (e HostFruFruType) Values() []HostFruFruType {
	return []HostFruFruType{
		HostFruFruTypeUndefined,
		HostFruFruTypeBoard,
		HostFruFruTypeProduct,
	}
}

func (e HostFruFruType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostFruFruType"] = reflect.TypeOf((*HostFruFruType)(nil)).Elem()
}

// Supported values for graphics type.
type HostGraphicsConfigGraphicsType string

const (
	// Shared graphics (ex.
	//
	// virtual shared graphics acceleration).
	HostGraphicsConfigGraphicsTypeShared = HostGraphicsConfigGraphicsType("shared")
	// Shared direct graphics (ex.
	//
	// vendor vGPU shared passthrough).
	HostGraphicsConfigGraphicsTypeSharedDirect = HostGraphicsConfigGraphicsType("sharedDirect")
)

func (e HostGraphicsConfigGraphicsType) Values() []HostGraphicsConfigGraphicsType {
	return []HostGraphicsConfigGraphicsType{
		HostGraphicsConfigGraphicsTypeShared,
		HostGraphicsConfigGraphicsTypeSharedDirect,
	}
}

func (e HostGraphicsConfigGraphicsType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostGraphicsConfigGraphicsType"] = reflect.TypeOf((*HostGraphicsConfigGraphicsType)(nil)).Elem()
}

// Supported values for shared passthrough assignment policy
type HostGraphicsConfigSharedPassthruAssignmentPolicy string

const (
	// Performance policy: assign VM to GPU with fewest VMs.
	HostGraphicsConfigSharedPassthruAssignmentPolicyPerformance = HostGraphicsConfigSharedPassthruAssignmentPolicy("performance")
	// Consolidation policy: group like VMs on GPU until fully loaded.
	HostGraphicsConfigSharedPassthruAssignmentPolicyConsolidation = HostGraphicsConfigSharedPassthruAssignmentPolicy("consolidation")
)

func (e HostGraphicsConfigSharedPassthruAssignmentPolicy) Values() []HostGraphicsConfigSharedPassthruAssignmentPolicy {
	return []HostGraphicsConfigSharedPassthruAssignmentPolicy{
		HostGraphicsConfigSharedPassthruAssignmentPolicyPerformance,
		HostGraphicsConfigSharedPassthruAssignmentPolicyConsolidation,
	}
}

func (e HostGraphicsConfigSharedPassthruAssignmentPolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostGraphicsConfigSharedPassthruAssignmentPolicy"] = reflect.TypeOf((*HostGraphicsConfigSharedPassthruAssignmentPolicy)(nil)).Elem()
}

type HostGraphicsConfigVgpuMode string

const (
	// vGPU time-sliced same size.
	HostGraphicsConfigVgpuModeSameSize = HostGraphicsConfigVgpuMode("sameSize")
	// vGPU time-sliced mixed size.
	HostGraphicsConfigVgpuModeMixedSize = HostGraphicsConfigVgpuMode("mixedSize")
)

func (e HostGraphicsConfigVgpuMode) Values() []HostGraphicsConfigVgpuMode {
	return []HostGraphicsConfigVgpuMode{
		HostGraphicsConfigVgpuModeSameSize,
		HostGraphicsConfigVgpuModeMixedSize,
	}
}

func (e HostGraphicsConfigVgpuMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostGraphicsConfigVgpuMode"] = reflect.TypeOf((*HostGraphicsConfigVgpuMode)(nil)).Elem()
	minAPIVersionForType["HostGraphicsConfigVgpuMode"] = "8.0.3.0"
}

// Possible values for graphics type.
type HostGraphicsInfoGraphicsType string

const (
	// Basic graphics when no host driver is available.
	HostGraphicsInfoGraphicsTypeBasic = HostGraphicsInfoGraphicsType("basic")
	// Shared graphics (ex.
	//
	// virtual shared graphics acceleration).
	HostGraphicsInfoGraphicsTypeShared = HostGraphicsInfoGraphicsType("shared")
	// Direct graphics (ex.
	//
	// passthrough).
	HostGraphicsInfoGraphicsTypeDirect = HostGraphicsInfoGraphicsType("direct")
	// Shared direct graphics (ex.
	//
	// vGPU shared passthrough).
	HostGraphicsInfoGraphicsTypeSharedDirect = HostGraphicsInfoGraphicsType("sharedDirect")
)

func (e HostGraphicsInfoGraphicsType) Values() []HostGraphicsInfoGraphicsType {
	return []HostGraphicsInfoGraphicsType{
		HostGraphicsInfoGraphicsTypeBasic,
		HostGraphicsInfoGraphicsTypeShared,
		HostGraphicsInfoGraphicsTypeDirect,
		HostGraphicsInfoGraphicsTypeSharedDirect,
	}
}

func (e HostGraphicsInfoGraphicsType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostGraphicsInfoGraphicsType"] = reflect.TypeOf((*HostGraphicsInfoGraphicsType)(nil)).Elem()
}

type HostGraphicsInfoVgpuMode string

const (
	// vGPU mode not applicable.
	HostGraphicsInfoVgpuModeNone = HostGraphicsInfoVgpuMode("none")
	// vGPU time-sliced same size.
	HostGraphicsInfoVgpuModeSameSize = HostGraphicsInfoVgpuMode("sameSize")
	// vGPU time-sliced mixed size.
	HostGraphicsInfoVgpuModeMixedSize = HostGraphicsInfoVgpuMode("mixedSize")
	// vGPU multi-instance GPU.
	HostGraphicsInfoVgpuModeMultiInstanceGpu = HostGraphicsInfoVgpuMode("multiInstanceGpu")
)

func (e HostGraphicsInfoVgpuMode) Values() []HostGraphicsInfoVgpuMode {
	return []HostGraphicsInfoVgpuMode{
		HostGraphicsInfoVgpuModeNone,
		HostGraphicsInfoVgpuModeSameSize,
		HostGraphicsInfoVgpuModeMixedSize,
		HostGraphicsInfoVgpuModeMultiInstanceGpu,
	}
}

func (e HostGraphicsInfoVgpuMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostGraphicsInfoVgpuMode"] = reflect.TypeOf((*HostGraphicsInfoVgpuMode)(nil)).Elem()
	minAPIVersionForType["HostGraphicsInfoVgpuMode"] = "8.0.3.0"
}

// The current status of the hardware
type HostHardwareElementStatus string

const (
	// The implementation cannot report on the current status of the
	// physical element
	HostHardwareElementStatusUnknown = HostHardwareElementStatus("Unknown")
	// The physical element is functioning as expected
	HostHardwareElementStatusGreen = HostHardwareElementStatus("Green")
	// All functionality is available but some might be degraded.
	HostHardwareElementStatusYellow = HostHardwareElementStatus("Yellow")
	// The physical element is failing.
	//
	// It is possible that some or all
	// functionalities of this physical element is degraded or not working.
	HostHardwareElementStatusRed = HostHardwareElementStatus("Red")
)

func (e HostHardwareElementStatus) Values() []HostHardwareElementStatus {
	return []HostHardwareElementStatus{
		HostHardwareElementStatusUnknown,
		HostHardwareElementStatusGreen,
		HostHardwareElementStatusYellow,
		HostHardwareElementStatusRed,
	}
}

func (e HostHardwareElementStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostHardwareElementStatus"] = reflect.TypeOf((*HostHardwareElementStatus)(nil)).Elem()
}

type HostHasComponentFailureHostComponentType string

const (
	HostHasComponentFailureHostComponentTypeDatastore = HostHasComponentFailureHostComponentType("Datastore")
)

func (e HostHasComponentFailureHostComponentType) Values() []HostHasComponentFailureHostComponentType {
	return []HostHasComponentFailureHostComponentType{
		HostHasComponentFailureHostComponentTypeDatastore,
	}
}

func (e HostHasComponentFailureHostComponentType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostHasComponentFailureHostComponentType"] = reflect.TypeOf((*HostHasComponentFailureHostComponentType)(nil)).Elem()
}

// Acceptance level definitions
type HostImageAcceptanceLevel string

const (
	// "VMware-certified"
	HostImageAcceptanceLevelVmware_certified = HostImageAcceptanceLevel("vmware_certified")
	// "VMware-accepted"
	HostImageAcceptanceLevelVmware_accepted = HostImageAcceptanceLevel("vmware_accepted")
	// "Partner-supported"
	HostImageAcceptanceLevelPartner = HostImageAcceptanceLevel("partner")
	// "Community-supported"
	HostImageAcceptanceLevelCommunity = HostImageAcceptanceLevel("community")
)

func (e HostImageAcceptanceLevel) Values() []HostImageAcceptanceLevel {
	return []HostImageAcceptanceLevel{
		HostImageAcceptanceLevelVmware_certified,
		HostImageAcceptanceLevelVmware_accepted,
		HostImageAcceptanceLevelPartner,
		HostImageAcceptanceLevelCommunity,
	}
}

func (e HostImageAcceptanceLevel) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostImageAcceptanceLevel"] = reflect.TypeOf((*HostImageAcceptanceLevel)(nil)).Elem()
}

// Reasons why fault tolerance is not supported on the host.
type HostIncompatibleForFaultToleranceReason string

const (
	// The product does not support fault tolerance.
	HostIncompatibleForFaultToleranceReasonProduct = HostIncompatibleForFaultToleranceReason("product")
	// The product supports fault tolerance but the host CPU does not.
	HostIncompatibleForFaultToleranceReasonProcessor = HostIncompatibleForFaultToleranceReason("processor")
)

func (e HostIncompatibleForFaultToleranceReason) Values() []HostIncompatibleForFaultToleranceReason {
	return []HostIncompatibleForFaultToleranceReason{
		HostIncompatibleForFaultToleranceReasonProduct,
		HostIncompatibleForFaultToleranceReasonProcessor,
	}
}

func (e HostIncompatibleForFaultToleranceReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostIncompatibleForFaultToleranceReason"] = reflect.TypeOf((*HostIncompatibleForFaultToleranceReason)(nil)).Elem()
}

// Reasons why record/replay is not supported on a host.
type HostIncompatibleForRecordReplayReason string

const (
	// The product does not support record/replay.
	HostIncompatibleForRecordReplayReasonProduct = HostIncompatibleForRecordReplayReason("product")
	// The product supports record/replay but the host CPU does not.
	HostIncompatibleForRecordReplayReasonProcessor = HostIncompatibleForRecordReplayReason("processor")
)

func (e HostIncompatibleForRecordReplayReason) Values() []HostIncompatibleForRecordReplayReason {
	return []HostIncompatibleForRecordReplayReason{
		HostIncompatibleForRecordReplayReasonProduct,
		HostIncompatibleForRecordReplayReasonProcessor,
	}
}

func (e HostIncompatibleForRecordReplayReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostIncompatibleForRecordReplayReason"] = reflect.TypeOf((*HostIncompatibleForRecordReplayReason)(nil)).Elem()
}

// The type of CHAP authentication setting to use.
//
// prohibited : do not use CHAP.
// preferred : use CHAP if successfully negotiated,
// but allow non-CHAP connections as fallback
// discouraged : use non-CHAP, but allow CHAP connectsion as fallback
// required : use CHAP for connection strictly, and fail if CHAP
// negotiation fails.
// Defaults to preferred on first configuration if unspecified.
type HostInternetScsiHbaChapAuthenticationType string

const (
	HostInternetScsiHbaChapAuthenticationTypeChapProhibited  = HostInternetScsiHbaChapAuthenticationType("chapProhibited")
	HostInternetScsiHbaChapAuthenticationTypeChapDiscouraged = HostInternetScsiHbaChapAuthenticationType("chapDiscouraged")
	HostInternetScsiHbaChapAuthenticationTypeChapPreferred   = HostInternetScsiHbaChapAuthenticationType("chapPreferred")
	HostInternetScsiHbaChapAuthenticationTypeChapRequired    = HostInternetScsiHbaChapAuthenticationType("chapRequired")
)

func (e HostInternetScsiHbaChapAuthenticationType) Values() []HostInternetScsiHbaChapAuthenticationType {
	return []HostInternetScsiHbaChapAuthenticationType{
		HostInternetScsiHbaChapAuthenticationTypeChapProhibited,
		HostInternetScsiHbaChapAuthenticationTypeChapDiscouraged,
		HostInternetScsiHbaChapAuthenticationTypeChapPreferred,
		HostInternetScsiHbaChapAuthenticationTypeChapRequired,
	}
}

func (e HostInternetScsiHbaChapAuthenticationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostInternetScsiHbaChapAuthenticationType"] = reflect.TypeOf((*HostInternetScsiHbaChapAuthenticationType)(nil)).Elem()
}

// The type of integrity checks to use.
//
// The digest setting for header
// and data traffic can be separately configured.
// prohibited : do not use digest.
// preferred : use digest if successfully negotiated, but skip the use
// of digest otherwise.
// discouraged : do not use digest if target allows, otherwise use digest.
// required : use digest strictly, and fail if target does not support
// digest.
// Defaults to preferred on first configuration if unspecified.
type HostInternetScsiHbaDigestType string

const (
	HostInternetScsiHbaDigestTypeDigestProhibited  = HostInternetScsiHbaDigestType("digestProhibited")
	HostInternetScsiHbaDigestTypeDigestDiscouraged = HostInternetScsiHbaDigestType("digestDiscouraged")
	HostInternetScsiHbaDigestTypeDigestPreferred   = HostInternetScsiHbaDigestType("digestPreferred")
	HostInternetScsiHbaDigestTypeDigestRequired    = HostInternetScsiHbaDigestType("digestRequired")
)

func (e HostInternetScsiHbaDigestType) Values() []HostInternetScsiHbaDigestType {
	return []HostInternetScsiHbaDigestType{
		HostInternetScsiHbaDigestTypeDigestProhibited,
		HostInternetScsiHbaDigestTypeDigestDiscouraged,
		HostInternetScsiHbaDigestTypeDigestPreferred,
		HostInternetScsiHbaDigestTypeDigestRequired,
	}
}

func (e HostInternetScsiHbaDigestType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostInternetScsiHbaDigestType"] = reflect.TypeOf((*HostInternetScsiHbaDigestType)(nil)).Elem()
}

// enum listing possible IPv6 address configuration methods.
type HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType string

const (
	// DHCP
	HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeDHCP = HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType("DHCP")
	// Auto configured.
	//
	// Auto configured Link local address and Router Advertisement addresses
	// would be of this type.
	HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeAutoConfigured = HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType("AutoConfigured")
	// Static address.
	//
	// Typically user specified addresses will be static addresses.
	// User can specify link local address. Only Static addresses can be added or removed.
	HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeStatic = HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType("Static")
	// Other or unknown type.
	HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeOther = HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType("Other")
)

func (e HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType) Values() []HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType {
	return []HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType{
		HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeDHCP,
		HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeAutoConfigured,
		HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeStatic,
		HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationTypeOther,
	}
}

func (e HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType"] = reflect.TypeOf((*HostInternetScsiHbaIscsiIpv6AddressAddressConfigurationType)(nil)).Elem()
}

// enum listing IPv6 address operations.
type HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation string

const (
	HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperationAdd    = HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation("add")
	HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperationRemove = HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation("remove")
)

func (e HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation) Values() []HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation {
	return []HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation{
		HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperationAdd,
		HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperationRemove,
	}
}

func (e HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation"] = reflect.TypeOf((*HostInternetScsiHbaIscsiIpv6AddressIPv6AddressOperation)(nil)).Elem()
}

// The binding mode of the adapter.
type HostInternetScsiHbaNetworkBindingSupportType string

const (
	HostInternetScsiHbaNetworkBindingSupportTypeNotsupported = HostInternetScsiHbaNetworkBindingSupportType("notsupported")
	HostInternetScsiHbaNetworkBindingSupportTypeOptional     = HostInternetScsiHbaNetworkBindingSupportType("optional")
	HostInternetScsiHbaNetworkBindingSupportTypeRequired     = HostInternetScsiHbaNetworkBindingSupportType("required")
)

func (e HostInternetScsiHbaNetworkBindingSupportType) Values() []HostInternetScsiHbaNetworkBindingSupportType {
	return []HostInternetScsiHbaNetworkBindingSupportType{
		HostInternetScsiHbaNetworkBindingSupportTypeNotsupported,
		HostInternetScsiHbaNetworkBindingSupportTypeOptional,
		HostInternetScsiHbaNetworkBindingSupportTypeRequired,
	}
}

func (e HostInternetScsiHbaNetworkBindingSupportType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostInternetScsiHbaNetworkBindingSupportType"] = reflect.TypeOf((*HostInternetScsiHbaNetworkBindingSupportType)(nil)).Elem()
}

// The method of discovery of an iScsi target.
//
// staticMethod: static discovery
// sendTargetsMethod: sendtarget discovery
// slpMethod: Service Location Protocol discovery
// isnsMethod: Internet Storage Name Service discovery
// unknownMethod: discovery method not identified by iscsi stack
type HostInternetScsiHbaStaticTargetTargetDiscoveryMethod string

const (
	HostInternetScsiHbaStaticTargetTargetDiscoveryMethodStaticMethod     = HostInternetScsiHbaStaticTargetTargetDiscoveryMethod("staticMethod")
	HostInternetScsiHbaStaticTargetTargetDiscoveryMethodSendTargetMethod = HostInternetScsiHbaStaticTargetTargetDiscoveryMethod("sendTargetMethod")
	HostInternetScsiHbaStaticTargetTargetDiscoveryMethodSlpMethod        = HostInternetScsiHbaStaticTargetTargetDiscoveryMethod("slpMethod")
	HostInternetScsiHbaStaticTargetTargetDiscoveryMethodIsnsMethod       = HostInternetScsiHbaStaticTargetTargetDiscoveryMethod("isnsMethod")
	HostInternetScsiHbaStaticTargetTargetDiscoveryMethodUnknownMethod    = HostInternetScsiHbaStaticTargetTargetDiscoveryMethod("unknownMethod")
)

func (e HostInternetScsiHbaStaticTargetTargetDiscoveryMethod) Values() []HostInternetScsiHbaStaticTargetTargetDiscoveryMethod {
	return []HostInternetScsiHbaStaticTargetTargetDiscoveryMethod{
		HostInternetScsiHbaStaticTargetTargetDiscoveryMethodStaticMethod,
		HostInternetScsiHbaStaticTargetTargetDiscoveryMethodSendTargetMethod,
		HostInternetScsiHbaStaticTargetTargetDiscoveryMethodSlpMethod,
		HostInternetScsiHbaStaticTargetTargetDiscoveryMethodIsnsMethod,
		HostInternetScsiHbaStaticTargetTargetDiscoveryMethodUnknownMethod,
	}
}

func (e HostInternetScsiHbaStaticTargetTargetDiscoveryMethod) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostInternetScsiHbaStaticTargetTargetDiscoveryMethod"] = reflect.TypeOf((*HostInternetScsiHbaStaticTargetTargetDiscoveryMethod)(nil)).Elem()
}

// This specifies how the ipv6 address is configured for the interface.
//
// We follow rfc4293 in defining the values for the configType.
type HostIpConfigIpV6AddressConfigType string

const (
	// Any other type of address configuration other than the below
	// mentioned ones will fall under this category.
	//
	// For e.g., automatic
	// address configuration for the link local address falls under
	// this type.
	HostIpConfigIpV6AddressConfigTypeOther = HostIpConfigIpV6AddressConfigType("other")
	// The address is configured manually.
	HostIpConfigIpV6AddressConfigTypeManual = HostIpConfigIpV6AddressConfigType("manual")
	// The address is configured through dhcp.
	HostIpConfigIpV6AddressConfigTypeDhcp = HostIpConfigIpV6AddressConfigType("dhcp")
	// The address is obtained through stateless autoconfiguration.
	HostIpConfigIpV6AddressConfigTypeLinklayer = HostIpConfigIpV6AddressConfigType("linklayer")
	// The address is chosen by the system at random
	// e.g., an IPv4 address within 169.254/16, or an RFC
	// 3041 privacy address.
	HostIpConfigIpV6AddressConfigTypeRandom = HostIpConfigIpV6AddressConfigType("random")
)

func (e HostIpConfigIpV6AddressConfigType) Values() []HostIpConfigIpV6AddressConfigType {
	return []HostIpConfigIpV6AddressConfigType{
		HostIpConfigIpV6AddressConfigTypeOther,
		HostIpConfigIpV6AddressConfigTypeManual,
		HostIpConfigIpV6AddressConfigTypeDhcp,
		HostIpConfigIpV6AddressConfigTypeLinklayer,
		HostIpConfigIpV6AddressConfigTypeRandom,
	}
}

func (e HostIpConfigIpV6AddressConfigType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostIpConfigIpV6AddressConfigType"] = reflect.TypeOf((*HostIpConfigIpV6AddressConfigType)(nil)).Elem()
}

type HostIpConfigIpV6AddressStatus string

const (
	// Indicates that this is a valid address.
	HostIpConfigIpV6AddressStatusPreferred = HostIpConfigIpV6AddressStatus("preferred")
	// Indicates that this is a valid but deprecated address
	// that should no longer be used as a source address.
	HostIpConfigIpV6AddressStatusDeprecated = HostIpConfigIpV6AddressStatus("deprecated")
	// Indicates that this isn't a valid.
	HostIpConfigIpV6AddressStatusInvalid = HostIpConfigIpV6AddressStatus("invalid")
	// Indicates that the address is not accessible because
	// interface is not operational.
	HostIpConfigIpV6AddressStatusInaccessible = HostIpConfigIpV6AddressStatus("inaccessible")
	// Indicates that the status cannot be determined.
	HostIpConfigIpV6AddressStatusUnknown = HostIpConfigIpV6AddressStatus("unknown")
	// Indicates that the uniqueness of the
	// address on the link is presently being verified.
	HostIpConfigIpV6AddressStatusTentative = HostIpConfigIpV6AddressStatus("tentative")
	// Indicates the address has been determined to be non-unique
	// on the link, this address will not be reachable.
	HostIpConfigIpV6AddressStatusDuplicate = HostIpConfigIpV6AddressStatus("duplicate")
)

func (e HostIpConfigIpV6AddressStatus) Values() []HostIpConfigIpV6AddressStatus {
	return []HostIpConfigIpV6AddressStatus{
		HostIpConfigIpV6AddressStatusPreferred,
		HostIpConfigIpV6AddressStatusDeprecated,
		HostIpConfigIpV6AddressStatusInvalid,
		HostIpConfigIpV6AddressStatusInaccessible,
		HostIpConfigIpV6AddressStatusUnknown,
		HostIpConfigIpV6AddressStatusTentative,
		HostIpConfigIpV6AddressStatusDuplicate,
	}
}

func (e HostIpConfigIpV6AddressStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostIpConfigIpV6AddressStatus"] = reflect.TypeOf((*HostIpConfigIpV6AddressStatus)(nil)).Elem()
}

// Identifiers of currently supported resources.
type HostLicensableResourceKey string

const (
	// Number of CPU packages on this host.
	HostLicensableResourceKeyNumCpuPackages = HostLicensableResourceKey("numCpuPackages")
	// Number of licensable CPU cores/compute-units on this host.
	HostLicensableResourceKeyNumCpuCores = HostLicensableResourceKey("numCpuCores")
	// Total size of memory installed on this host, measured in kilobytes.
	HostLicensableResourceKeyMemorySize = HostLicensableResourceKey("memorySize")
	// Total size of memory configured for VMs on this host, measured in kilobytes.
	HostLicensableResourceKeyMemoryForVms = HostLicensableResourceKey("memoryForVms")
	// Number of VMs already running on this host.
	HostLicensableResourceKeyNumVmsStarted = HostLicensableResourceKey("numVmsStarted")
	// Number of VMs that are currently powering-on, immigrating, etc.
	HostLicensableResourceKeyNumVmsStarting = HostLicensableResourceKey("numVmsStarting")
	// vSAN capacity in TiB on this host.
	HostLicensableResourceKeyVsanCapacity = HostLicensableResourceKey("vsanCapacity")
)

func (e HostLicensableResourceKey) Values() []HostLicensableResourceKey {
	return []HostLicensableResourceKey{
		HostLicensableResourceKeyNumCpuPackages,
		HostLicensableResourceKeyNumCpuCores,
		HostLicensableResourceKeyMemorySize,
		HostLicensableResourceKeyMemoryForVms,
		HostLicensableResourceKeyNumVmsStarted,
		HostLicensableResourceKeyNumVmsStarting,
		HostLicensableResourceKeyVsanCapacity,
	}
}

func (e HostLicensableResourceKey) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostLicensableResourceKey"] = reflect.TypeOf((*HostLicensableResourceKey)(nil)).Elem()
	minAPIVersionForEnumValue["HostLicensableResourceKey"] = map[string]string{
		"vsanCapacity": "8.0.3.0",
	}
}

// Defines the possible states of lockdown mode.
type HostLockdownMode string

const (
	// Indicates that lockdown mode is disabled.
	HostLockdownModeLockdownDisabled = HostLockdownMode("lockdownDisabled")
	// Indicates that lockdown mode is enabled with service DCUI
	// (Direct Console User Interface) running.
	HostLockdownModeLockdownNormal = HostLockdownMode("lockdownNormal")
	// Indicates that lockdown mode is enabled with service DCUI stopped.
	//
	// If the host is in "strict" lockdown mode then no one will be able
	// to exit lockdown mode through DCUI in emergency situations,
	// i.e. when the connection to vCenter server is permanently lost.
	HostLockdownModeLockdownStrict = HostLockdownMode("lockdownStrict")
)

func (e HostLockdownMode) Values() []HostLockdownMode {
	return []HostLockdownMode{
		HostLockdownModeLockdownDisabled,
		HostLockdownModeLockdownNormal,
		HostLockdownModeLockdownStrict,
	}
}

func (e HostLockdownMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostLockdownMode"] = reflect.TypeOf((*HostLockdownMode)(nil)).Elem()
}

// This enum defines the possible types of file types that can be reserved
// or deleted
type HostLowLevelProvisioningManagerFileType string

const (
	HostLowLevelProvisioningManagerFileTypeFile        = HostLowLevelProvisioningManagerFileType("File")
	HostLowLevelProvisioningManagerFileTypeVirtualDisk = HostLowLevelProvisioningManagerFileType("VirtualDisk")
	HostLowLevelProvisioningManagerFileTypeDirectory   = HostLowLevelProvisioningManagerFileType("Directory")
)

func (e HostLowLevelProvisioningManagerFileType) Values() []HostLowLevelProvisioningManagerFileType {
	return []HostLowLevelProvisioningManagerFileType{
		HostLowLevelProvisioningManagerFileTypeFile,
		HostLowLevelProvisioningManagerFileTypeVirtualDisk,
		HostLowLevelProvisioningManagerFileTypeDirectory,
	}
}

func (e HostLowLevelProvisioningManagerFileType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostLowLevelProvisioningManagerFileType"] = reflect.TypeOf((*HostLowLevelProvisioningManagerFileType)(nil)).Elem()
}

// The target of the disk reload.
type HostLowLevelProvisioningManagerReloadTarget string

const (
	// Specifies the reload of the current config of the virtual machine.
	HostLowLevelProvisioningManagerReloadTargetCurrentConfig = HostLowLevelProvisioningManagerReloadTarget("currentConfig")
	// Specifies the reload of the snapshot config of the virtual machine.
	//
	// If the virtual machine has multiple snapshots, all of the snapshot's
	// config will be reloaded.
	HostLowLevelProvisioningManagerReloadTargetSnapshotConfig = HostLowLevelProvisioningManagerReloadTarget("snapshotConfig")
)

func (e HostLowLevelProvisioningManagerReloadTarget) Values() []HostLowLevelProvisioningManagerReloadTarget {
	return []HostLowLevelProvisioningManagerReloadTarget{
		HostLowLevelProvisioningManagerReloadTargetCurrentConfig,
		HostLowLevelProvisioningManagerReloadTargetSnapshotConfig,
	}
}

func (e HostLowLevelProvisioningManagerReloadTarget) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostLowLevelProvisioningManagerReloadTarget"] = reflect.TypeOf((*HostLowLevelProvisioningManagerReloadTarget)(nil)).Elem()
}

type HostMaintenanceSpecPurpose string

const (
	HostMaintenanceSpecPurposeHostUpgrade = HostMaintenanceSpecPurpose("hostUpgrade")
)

func (e HostMaintenanceSpecPurpose) Values() []HostMaintenanceSpecPurpose {
	return []HostMaintenanceSpecPurpose{
		HostMaintenanceSpecPurposeHostUpgrade,
	}
}

func (e HostMaintenanceSpecPurpose) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMaintenanceSpecPurpose"] = reflect.TypeOf((*HostMaintenanceSpecPurpose)(nil)).Elem()
}

// Enumeration of flags pertaining to a memory tier.
//
// Here are some examples of what the flags will look like for various memory
// configurations:
//   - Traditional memory (`noTiering`): The host has a DRAM tier
//     for the main memory and nothing else. The DRAM tier will have the
//     `memoryTier` flag.
//   - App Direct mode (`noTiering`): The host has a DRAM tier
//     and a PMem tier, but the two are independent and unrelated. The PMem tier is
//     non-volatile and is exposed as an NVDIMM device. Applications can decide whether to
//     direct the reads and writes to DRAM or PMem by using the appropriate system call. The
//     DRAM tier will have the `memoryTier` flag and the PMem tier will
//     have the `persistentTier` flag.
//   - Memory mode (`hardwareTiering`): The host has a DRAM tier
//     and a PMem tier, but the DRAM is hidden from applications and is just a cache
//     for the PMem main memory. The PMem tier is volatile, and is abstracted by the hardware
//     layer to look like traditional memory. Applications can read from/write to memory
//     using the traditional memory system calls. The memory controller in the hardware will
//     internally direct those to the DRAM cache first, and on a cache miss redirect them to
//     the PMem main memory. The DRAM tier will have the `cachingTier`
type HostMemoryTierFlags string

const (
	// Flag indicating that the tier is the primary memory tier visible from the
	// host.
	HostMemoryTierFlagsMemoryTier = HostMemoryTierFlags("memoryTier")
	// Flag indicating that the tier is used as non-volatile storage, e.g.
	//
	// PMem in
	// App Direct mode.
	HostMemoryTierFlagsPersistentTier = HostMemoryTierFlags("persistentTier")
	// Flag indicating that the tier is a cache for main memory.
	HostMemoryTierFlagsCachingTier = HostMemoryTierFlags("cachingTier")
	// `**Since:**` vSphere API Release 8.0.3.0
	HostMemoryTierFlagsUnmappableTier = HostMemoryTierFlags("unmappableTier")
)

func (e HostMemoryTierFlags) Values() []HostMemoryTierFlags {
	return []HostMemoryTierFlags{
		HostMemoryTierFlagsMemoryTier,
		HostMemoryTierFlagsPersistentTier,
		HostMemoryTierFlagsCachingTier,
		HostMemoryTierFlagsUnmappableTier,
	}
}

func (e HostMemoryTierFlags) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMemoryTierFlags"] = reflect.TypeOf((*HostMemoryTierFlags)(nil)).Elem()
	minAPIVersionForType["HostMemoryTierFlags"] = "7.0.3.0"
	minAPIVersionForEnumValue["HostMemoryTierFlags"] = map[string]string{
		"unmappableTier": "8.0.3.0",
	}
}

type HostMemoryTierType string

const (
	// Dynamic random-access memory.
	HostMemoryTierTypeDRAM = HostMemoryTierType("DRAM")
	// Persistent memory.
	HostMemoryTierTypePMem = HostMemoryTierType("PMem")
	// NVMe memory.
	HostMemoryTierTypeNVMe = HostMemoryTierType("NVMe")
)

func (e HostMemoryTierType) Values() []HostMemoryTierType {
	return []HostMemoryTierType{
		HostMemoryTierTypeDRAM,
		HostMemoryTierTypePMem,
		HostMemoryTierTypeNVMe,
	}
}

func (e HostMemoryTierType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMemoryTierType"] = reflect.TypeOf((*HostMemoryTierType)(nil)).Elem()
	minAPIVersionForType["HostMemoryTierType"] = "7.0.3.0"
	minAPIVersionForEnumValue["HostMemoryTierType"] = map[string]string{
		"NVMe": "8.0.3.0",
	}
}

type HostMemoryTieringType string

const (
	// The traditional memory configuration without any tiers.
	HostMemoryTieringTypeNoTiering = HostMemoryTieringType("noTiering")
	// The memory configuration where a tier is hardware-controlled and invisible to
	// applications, e.g.
	//
	// Intel's Memory Mode.
	HostMemoryTieringTypeHardwareTiering = HostMemoryTieringType("hardwareTiering")
	// The memory configuration where all memory tiers are managed by software (ESX).
	HostMemoryTieringTypeSoftwareTiering = HostMemoryTieringType("softwareTiering")
)

func (e HostMemoryTieringType) Values() []HostMemoryTieringType {
	return []HostMemoryTieringType{
		HostMemoryTieringTypeNoTiering,
		HostMemoryTieringTypeHardwareTiering,
		HostMemoryTieringTypeSoftwareTiering,
	}
}

func (e HostMemoryTieringType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMemoryTieringType"] = reflect.TypeOf((*HostMemoryTieringType)(nil)).Elem()
	minAPIVersionForType["HostMemoryTieringType"] = "7.0.3.0"
	minAPIVersionForEnumValue["HostMemoryTieringType"] = map[string]string{
		"softwareTiering": "8.0.3.0",
	}
}

// A datastore can become inaccessible due to a number of reasons as
// defined in this enum `HostMountInfoInaccessibleReason_enum`.
//
// The reason for a datastore being inaccessible is reported in
// `HostMountInfo.inaccessibleReason`.
// APD ("All Paths Down") is a condition where a SAN or NFS storage has
// become inaccessible for unknown reasons. It only indicates loss of
// connectivity and does not indicate storage device failure or
// LUN removal (Permanent Device Loss or PDL)
// A difference between APD and PDL is that APD may recover
// in which case all use cases will start to work as before. In case of PDL
// the failed datastore/device is unlikely to recover and hence the device
// path information and data cache will be emptied. If the PDL condition
// recovers, the failed datastores have to be added back to the host. Once
// in PDL a datastore cannot be added back until there are no longer any
// open files on the datastore.
// PDL is not linked to the APD and can happen at any time with or without APD
// preceding. If APD and PDL occur at the same time, APD will be reported first.
// Once (and if) the APD condition clears, PermanentDataLoss will be reported if
// PDL condition still exists.
type HostMountInfoInaccessibleReason string

const (
	// AllPathsDown\_Start value is reported when all paths down state is detected
	HostMountInfoInaccessibleReasonAllPathsDown_Start = HostMountInfoInaccessibleReason("AllPathsDown_Start")
	// After a wait for a system default time (which is user modifiable)
	// to ascertain the state is indeed an APD, AllPathsDown\_Timeout property
	// is reported.
	//
	// The host advanced option used to set timeout period
	// is "/Misc/APDTimeout"
	// After the datastore property is set to AllPathsDown\_Timeout, all data i/o
	// to the datastore will be fast-failed (failed immediately).
	HostMountInfoInaccessibleReasonAllPathsDown_Timeout = HostMountInfoInaccessibleReason("AllPathsDown_Timeout")
	// A PDL condition is reported as PermanentDeviceLoss.
	HostMountInfoInaccessibleReasonPermanentDeviceLoss = HostMountInfoInaccessibleReason("PermanentDeviceLoss")
)

func (e HostMountInfoInaccessibleReason) Values() []HostMountInfoInaccessibleReason {
	return []HostMountInfoInaccessibleReason{
		HostMountInfoInaccessibleReasonAllPathsDown_Start,
		HostMountInfoInaccessibleReasonAllPathsDown_Timeout,
		HostMountInfoInaccessibleReasonPermanentDeviceLoss,
	}
}

func (e HostMountInfoInaccessibleReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMountInfoInaccessibleReason"] = reflect.TypeOf((*HostMountInfoInaccessibleReason)(nil)).Elem()
}

// NFS mount request can be failed due to a number of reasons as
// defined in this enum `HostMountInfoMountFailedReason_enum`.
//
// The reason for the mount failure is reported in
// `HostMountInfo.mountFailedReason`. This is applicable only for those
type HostMountInfoMountFailedReason string

const (
	// Failed to get port or connect.
	//
	// Or MOUNT/FSINFO RPC failed.
	HostMountInfoMountFailedReasonCONNECT_FAILURE = HostMountInfoMountFailedReason("CONNECT_FAILURE")
	// Server doesn't support MOUNT\_PROGRAM/MOUNT\_PROGRAM\_VERSION.
	HostMountInfoMountFailedReasonMOUNT_NOT_SUPPORTED = HostMountInfoMountFailedReason("MOUNT_NOT_SUPPORTED")
	// Server doesn't support NFS\_PROGRAM/NFS\_PROGRAM\_VERSION.
	HostMountInfoMountFailedReasonNFS_NOT_SUPPORTED = HostMountInfoMountFailedReason("NFS_NOT_SUPPORTED")
	// No permission to mount the remote volume or it doesn't exist.
	HostMountInfoMountFailedReasonMOUNT_DENIED = HostMountInfoMountFailedReason("MOUNT_DENIED")
	// Remote path not a directory.
	HostMountInfoMountFailedReasonMOUNT_NOT_DIR = HostMountInfoMountFailedReason("MOUNT_NOT_DIR")
	// Maximum NFS volumes have been mounted.
	HostMountInfoMountFailedReasonVOLUME_LIMIT_EXCEEDED = HostMountInfoMountFailedReason("VOLUME_LIMIT_EXCEEDED")
	// Maximum connections for NFS has been reached.
	HostMountInfoMountFailedReasonCONN_LIMIT_EXCEEDED = HostMountInfoMountFailedReason("CONN_LIMIT_EXCEEDED")
	// Volume already mounted or a different mount exists with same label.
	HostMountInfoMountFailedReasonMOUNT_EXISTS = HostMountInfoMountFailedReason("MOUNT_EXISTS")
	// Any other reason which is not present in above list.
	HostMountInfoMountFailedReasonOTHERS = HostMountInfoMountFailedReason("OTHERS")
)

func (e HostMountInfoMountFailedReason) Values() []HostMountInfoMountFailedReason {
	return []HostMountInfoMountFailedReason{
		HostMountInfoMountFailedReasonCONNECT_FAILURE,
		HostMountInfoMountFailedReasonMOUNT_NOT_SUPPORTED,
		HostMountInfoMountFailedReasonNFS_NOT_SUPPORTED,
		HostMountInfoMountFailedReasonMOUNT_DENIED,
		HostMountInfoMountFailedReasonMOUNT_NOT_DIR,
		HostMountInfoMountFailedReasonVOLUME_LIMIT_EXCEEDED,
		HostMountInfoMountFailedReasonCONN_LIMIT_EXCEEDED,
		HostMountInfoMountFailedReasonMOUNT_EXISTS,
		HostMountInfoMountFailedReasonOTHERS,
	}
}

func (e HostMountInfoMountFailedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMountInfoMountFailedReason"] = reflect.TypeOf((*HostMountInfoMountFailedReason)(nil)).Elem()
	minAPIVersionForType["HostMountInfoMountFailedReason"] = "8.0.0.1"
}

// Defines the access mode of the datastore.
type HostMountMode string

const (
	// The host system has read/write access to the file system.
	HostMountModeReadWrite = HostMountMode("readWrite")
	// The host system has read-only access to the file system.
	HostMountModeReadOnly = HostMountMode("readOnly")
)

func (e HostMountMode) Values() []HostMountMode {
	return []HostMountMode{
		HostMountModeReadWrite,
		HostMountModeReadOnly,
	}
}

func (e HostMountMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostMountMode"] = reflect.TypeOf((*HostMountMode)(nil)).Elem()
}

// Security type supported.
type HostNasVolumeSecurityType string

const (
	// Authentication based on traditional UNIX identifiers (UID and GID).
	//
	// Server trusts the IDs sent by the client for each request and uses them
	// to perform access control. Current implementation only supports
	// AUTH\_SYS with root user.
	HostNasVolumeSecurityTypeAUTH_SYS = HostNasVolumeSecurityType("AUTH_SYS")
	// Ensures RPC header authentication using Kerberos session keys.
	//
	// When
	// this option is enabled, the client uses the information specified in
	// `HostNasVolumeUserInfo` to establish shared keys with the server using
	// Kerberos. These shared keys are used to generate and verify message
	// authentication codes for RPC header of NFS requests and responses,
	// respectively. This method does not secure NFS file data.
	HostNasVolumeSecurityTypeSEC_KRB5 = HostNasVolumeSecurityType("SEC_KRB5")
	// Extends SEC\_KRB5 to generate and verify message authentication codes
	// for the payload of NFS requests and responses respectively.
	//
	// This
	// ensures the integrity of the NFS file data.
	HostNasVolumeSecurityTypeSEC_KRB5I = HostNasVolumeSecurityType("SEC_KRB5I")
)

func (e HostNasVolumeSecurityType) Values() []HostNasVolumeSecurityType {
	return []HostNasVolumeSecurityType{
		HostNasVolumeSecurityTypeAUTH_SYS,
		HostNasVolumeSecurityTypeSEC_KRB5,
		HostNasVolumeSecurityTypeSEC_KRB5I,
	}
}

func (e HostNasVolumeSecurityType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNasVolumeSecurityType"] = reflect.TypeOf((*HostNasVolumeSecurityType)(nil)).Elem()
}

// Define TCP congestion control algorithm used by an instance
type HostNetStackInstanceCongestionControlAlgorithmType string

const (
	// New Reno Algorithm.
	//
	// See http://tools.ietf.org/html/rfc3782 for detail.
	HostNetStackInstanceCongestionControlAlgorithmTypeNewreno = HostNetStackInstanceCongestionControlAlgorithmType("newreno")
	// Cubic Algorithm.
	//
	// See http://tools.ietf.org/id/draft-rhee-tcp-cubic-00.txt for detail.
	HostNetStackInstanceCongestionControlAlgorithmTypeCubic = HostNetStackInstanceCongestionControlAlgorithmType("cubic")
)

func (e HostNetStackInstanceCongestionControlAlgorithmType) Values() []HostNetStackInstanceCongestionControlAlgorithmType {
	return []HostNetStackInstanceCongestionControlAlgorithmType{
		HostNetStackInstanceCongestionControlAlgorithmTypeNewreno,
		HostNetStackInstanceCongestionControlAlgorithmTypeCubic,
	}
}

func (e HostNetStackInstanceCongestionControlAlgorithmType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNetStackInstanceCongestionControlAlgorithmType"] = reflect.TypeOf((*HostNetStackInstanceCongestionControlAlgorithmType)(nil)).Elem()
}

// Define the instance identifier for different traffic type
type HostNetStackInstanceSystemStackKey string

const (
	// The default stack used by applications
	HostNetStackInstanceSystemStackKeyDefaultTcpipStack = HostNetStackInstanceSystemStackKey("defaultTcpipStack")
	// Stack key used for vMotion applications
	HostNetStackInstanceSystemStackKeyVmotion = HostNetStackInstanceSystemStackKey("vmotion")
	// Stack key used for vSphere provisioning NFC traffic
	HostNetStackInstanceSystemStackKeyVSphereProvisioning = HostNetStackInstanceSystemStackKey("vSphereProvisioning")
	// Stack key used for port mirroring
	HostNetStackInstanceSystemStackKeyMirror = HostNetStackInstanceSystemStackKey("mirror")
	// Stack key used for ops applications
	HostNetStackInstanceSystemStackKeyOps = HostNetStackInstanceSystemStackKey("ops")
)

func (e HostNetStackInstanceSystemStackKey) Values() []HostNetStackInstanceSystemStackKey {
	return []HostNetStackInstanceSystemStackKey{
		HostNetStackInstanceSystemStackKeyDefaultTcpipStack,
		HostNetStackInstanceSystemStackKeyVmotion,
		HostNetStackInstanceSystemStackKeyVSphereProvisioning,
		HostNetStackInstanceSystemStackKeyMirror,
		HostNetStackInstanceSystemStackKeyOps,
	}
}

func (e HostNetStackInstanceSystemStackKey) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNetStackInstanceSystemStackKey"] = reflect.TypeOf((*HostNetStackInstanceSystemStackKey)(nil)).Elem()
	minAPIVersionForEnumValue["HostNetStackInstanceSystemStackKey"] = map[string]string{
		"mirror": "8.0.0.1",
		"ops":    "8.0.0.1",
	}
}

// Health state of the numeric sensor as reported by the sensor probes.
//
// Same data reported using command line: esxcli hardware ipmi sdr list
type HostNumericSensorHealthState string

const (
	// The implementation cannot report on the current health state of the
	// physical element
	HostNumericSensorHealthStateUnknown = HostNumericSensorHealthState("unknown")
	// The sensor is operating under normal conditions
	HostNumericSensorHealthStateGreen = HostNumericSensorHealthState("green")
	// The sensor is operating under conditions that are non-critical.
	HostNumericSensorHealthStateYellow = HostNumericSensorHealthState("yellow")
	// The sensor is operating under critical or fatal conditions.
	//
	// This may
	// directly affect the functioning of both the sensor and related
	// components.
	HostNumericSensorHealthStateRed = HostNumericSensorHealthState("red")
)

func (e HostNumericSensorHealthState) Values() []HostNumericSensorHealthState {
	return []HostNumericSensorHealthState{
		HostNumericSensorHealthStateUnknown,
		HostNumericSensorHealthStateGreen,
		HostNumericSensorHealthStateYellow,
		HostNumericSensorHealthStateRed,
	}
}

func (e HostNumericSensorHealthState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNumericSensorHealthState"] = reflect.TypeOf((*HostNumericSensorHealthState)(nil)).Elem()
}

// Sensor Types for specific hardware component are either based on
// class of sensor or what the sensor monitors to allow for grouping
type HostNumericSensorType string

const (
	// Fan sensor
	HostNumericSensorTypeFan = HostNumericSensorType("fan")
	// Power sensor
	HostNumericSensorTypePower = HostNumericSensorType("power")
	// Temperature sensor
	HostNumericSensorTypeTemperature = HostNumericSensorType("temperature")
	// Voltage Sensor
	HostNumericSensorTypeVoltage = HostNumericSensorType("voltage")
	// Other sensor.
	HostNumericSensorTypeOther = HostNumericSensorType("other")
	// Processor sensor.
	HostNumericSensorTypeProcessor = HostNumericSensorType("processor")
	// Memory sensor.
	HostNumericSensorTypeMemory = HostNumericSensorType("memory")
	// disk/storage sensor.
	HostNumericSensorTypeStorage = HostNumericSensorType("storage")
	// system board sensor.
	HostNumericSensorTypeSystemBoard = HostNumericSensorType("systemBoard")
	// Battery sensor.
	HostNumericSensorTypeBattery = HostNumericSensorType("battery")
	// BIOS/firmware related sensor.
	HostNumericSensorTypeBios = HostNumericSensorType("bios")
	// cable related sensor.
	HostNumericSensorTypeCable = HostNumericSensorType("cable")
	// Watchdog related sensor.
	HostNumericSensorTypeWatchdog = HostNumericSensorType("watchdog")
)

func (e HostNumericSensorType) Values() []HostNumericSensorType {
	return []HostNumericSensorType{
		HostNumericSensorTypeFan,
		HostNumericSensorTypePower,
		HostNumericSensorTypeTemperature,
		HostNumericSensorTypeVoltage,
		HostNumericSensorTypeOther,
		HostNumericSensorTypeProcessor,
		HostNumericSensorTypeMemory,
		HostNumericSensorTypeStorage,
		HostNumericSensorTypeSystemBoard,
		HostNumericSensorTypeBattery,
		HostNumericSensorTypeBios,
		HostNumericSensorTypeCable,
		HostNumericSensorTypeWatchdog,
	}
}

func (e HostNumericSensorType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNumericSensorType"] = reflect.TypeOf((*HostNumericSensorType)(nil)).Elem()
}

// This enum represents the supported NVM subsystem types.
type HostNvmeDiscoveryLogSubsystemType string

const (
	// A Discovery service, composed of Discovery controllers.
	HostNvmeDiscoveryLogSubsystemTypeDiscovery = HostNvmeDiscoveryLogSubsystemType("discovery")
	// An NVM subsystem whose controllers may have attached namespaces.
	HostNvmeDiscoveryLogSubsystemTypeNvm = HostNvmeDiscoveryLogSubsystemType("nvm")
)

func (e HostNvmeDiscoveryLogSubsystemType) Values() []HostNvmeDiscoveryLogSubsystemType {
	return []HostNvmeDiscoveryLogSubsystemType{
		HostNvmeDiscoveryLogSubsystemTypeDiscovery,
		HostNvmeDiscoveryLogSubsystemTypeNvm,
	}
}

func (e HostNvmeDiscoveryLogSubsystemType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNvmeDiscoveryLogSubsystemType"] = reflect.TypeOf((*HostNvmeDiscoveryLogSubsystemType)(nil)).Elem()
}

// This enum represents the supported types of transport requirements.
type HostNvmeDiscoveryLogTransportRequirements string

const (
	// A fabric secure channel is required.
	HostNvmeDiscoveryLogTransportRequirementsSecureChannelRequired = HostNvmeDiscoveryLogTransportRequirements("secureChannelRequired")
	// A fabric secure channel is not required.
	HostNvmeDiscoveryLogTransportRequirementsSecureChannelNotRequired = HostNvmeDiscoveryLogTransportRequirements("secureChannelNotRequired")
	// Requirements are not specified
	HostNvmeDiscoveryLogTransportRequirementsRequirementsNotSpecified = HostNvmeDiscoveryLogTransportRequirements("requirementsNotSpecified")
)

func (e HostNvmeDiscoveryLogTransportRequirements) Values() []HostNvmeDiscoveryLogTransportRequirements {
	return []HostNvmeDiscoveryLogTransportRequirements{
		HostNvmeDiscoveryLogTransportRequirementsSecureChannelRequired,
		HostNvmeDiscoveryLogTransportRequirementsSecureChannelNotRequired,
		HostNvmeDiscoveryLogTransportRequirementsRequirementsNotSpecified,
	}
}

func (e HostNvmeDiscoveryLogTransportRequirements) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNvmeDiscoveryLogTransportRequirements"] = reflect.TypeOf((*HostNvmeDiscoveryLogTransportRequirements)(nil)).Elem()
}

// This enum specifies the supported address families for
// NVME over Fabrics.
//
// For details, see:
//   - "NVM Express over Fabrics 1.0", Section 5.3, Figure 34,
//     "Discovery Log Page Entry"
type HostNvmeTransportParametersNvmeAddressFamily string

const (
	// IPv4 address, format specified in IETF RFC 791.
	HostNvmeTransportParametersNvmeAddressFamilyIpv4 = HostNvmeTransportParametersNvmeAddressFamily("ipv4")
	// IPv6 address, format specified in IETF RFC 2373.
	HostNvmeTransportParametersNvmeAddressFamilyIpv6 = HostNvmeTransportParametersNvmeAddressFamily("ipv6")
	// InfiniBand address family.
	HostNvmeTransportParametersNvmeAddressFamilyInfiniBand = HostNvmeTransportParametersNvmeAddressFamily("infiniBand")
	// Fibre Channel address family.
	HostNvmeTransportParametersNvmeAddressFamilyFc = HostNvmeTransportParametersNvmeAddressFamily("fc")
	// Intra-host transport.
	HostNvmeTransportParametersNvmeAddressFamilyLoopback = HostNvmeTransportParametersNvmeAddressFamily("loopback")
	// Unrecognized address family.
	HostNvmeTransportParametersNvmeAddressFamilyUnknown = HostNvmeTransportParametersNvmeAddressFamily("unknown")
)

func (e HostNvmeTransportParametersNvmeAddressFamily) Values() []HostNvmeTransportParametersNvmeAddressFamily {
	return []HostNvmeTransportParametersNvmeAddressFamily{
		HostNvmeTransportParametersNvmeAddressFamilyIpv4,
		HostNvmeTransportParametersNvmeAddressFamilyIpv6,
		HostNvmeTransportParametersNvmeAddressFamilyInfiniBand,
		HostNvmeTransportParametersNvmeAddressFamilyFc,
		HostNvmeTransportParametersNvmeAddressFamilyLoopback,
		HostNvmeTransportParametersNvmeAddressFamilyUnknown,
	}
}

func (e HostNvmeTransportParametersNvmeAddressFamily) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNvmeTransportParametersNvmeAddressFamily"] = reflect.TypeOf((*HostNvmeTransportParametersNvmeAddressFamily)(nil)).Elem()
}

// The set of NVM Express over Fabrics transport types.
//
// For details, see:
//   - "NVM Express over Fabrics 1.0", Section 1.5.1,
//     "Fabrics and Transports".
type HostNvmeTransportType string

const (
	// PCI Express transport type
	HostNvmeTransportTypePcie = HostNvmeTransportType("pcie")
	// Fibre Channel transport type
	HostNvmeTransportTypeFibreChannel = HostNvmeTransportType("fibreChannel")
	// Remote Direct Memory Access transport type
	HostNvmeTransportTypeRdma = HostNvmeTransportType("rdma")
	// Transmission Control Protocol transport type
	HostNvmeTransportTypeTcp = HostNvmeTransportType("tcp")
	// Intra-host transport.
	HostNvmeTransportTypeLoopback = HostNvmeTransportType("loopback")
	// The transport type is not among the currently supported ones.
	HostNvmeTransportTypeUnsupported = HostNvmeTransportType("unsupported")
)

func (e HostNvmeTransportType) Values() []HostNvmeTransportType {
	return []HostNvmeTransportType{
		HostNvmeTransportTypePcie,
		HostNvmeTransportTypeFibreChannel,
		HostNvmeTransportTypeRdma,
		HostNvmeTransportTypeTcp,
		HostNvmeTransportTypeLoopback,
		HostNvmeTransportTypeUnsupported,
	}
}

func (e HostNvmeTransportType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostNvmeTransportType"] = reflect.TypeOf((*HostNvmeTransportType)(nil)).Elem()
	minAPIVersionForEnumValue["HostNvmeTransportType"] = map[string]string{
		"tcp": "7.0.3.0",
	}
}

type HostOpaqueSwitchOpaqueSwitchState string

const (
	// The opaque switch is up and running.
	HostOpaqueSwitchOpaqueSwitchStateUp = HostOpaqueSwitchOpaqueSwitchState("up")
	// The opaque switch requires attention.
	HostOpaqueSwitchOpaqueSwitchStateWarning = HostOpaqueSwitchOpaqueSwitchState("warning")
	// The opaque switch is down.
	HostOpaqueSwitchOpaqueSwitchStateDown = HostOpaqueSwitchOpaqueSwitchState("down")
	// The opaque switch is under upgrade.
	HostOpaqueSwitchOpaqueSwitchStateMaintenance = HostOpaqueSwitchOpaqueSwitchState("maintenance")
)

func (e HostOpaqueSwitchOpaqueSwitchState) Values() []HostOpaqueSwitchOpaqueSwitchState {
	return []HostOpaqueSwitchOpaqueSwitchState{
		HostOpaqueSwitchOpaqueSwitchStateUp,
		HostOpaqueSwitchOpaqueSwitchStateWarning,
		HostOpaqueSwitchOpaqueSwitchStateDown,
		HostOpaqueSwitchOpaqueSwitchStateMaintenance,
	}
}

func (e HostOpaqueSwitchOpaqueSwitchState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostOpaqueSwitchOpaqueSwitchState"] = reflect.TypeOf((*HostOpaqueSwitchOpaqueSwitchState)(nil)).Elem()
}

// The following enum describes some common kinds of partial maintenance modes,
type HostPartialMaintenanceModeId string

const (
	// When the host is in the quick patch partial maintenance mode, it is safe to
	// perform a quick patch.
	//
	// When the host is in this partial maintenance mode, any virtual machines
	// and/or pods placed on it will continue to run but operations which may
	// lead to new workloads starting on the host such as power on or incoming
	// vmotions may be blocked.
	// It is generally unsafe to reboot the host in this state.
	HostPartialMaintenanceModeIdQuickPatchPartialMM = HostPartialMaintenanceModeId("quickPatchPartialMM")
)

func (e HostPartialMaintenanceModeId) Values() []HostPartialMaintenanceModeId {
	return []HostPartialMaintenanceModeId{
		HostPartialMaintenanceModeIdQuickPatchPartialMM,
	}
}

func (e HostPartialMaintenanceModeId) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPartialMaintenanceModeId"] = reflect.TypeOf((*HostPartialMaintenanceModeId)(nil)).Elem()
	minAPIVersionForType["HostPartialMaintenanceModeId"] = "8.0.3.0"
	minAPIVersionForEnumValue["HostPartialMaintenanceModeId"] = map[string]string{
		"quickPatchPartialMM": "8.0.3.0",
	}
}

// The following enum contains the list of possible statuses associated
type HostPartialMaintenanceModeStatus string

const (
	// The host is not in the particular partial maintenance mode.
	HostPartialMaintenanceModeStatusNotInPartialMM = HostPartialMaintenanceModeStatus("notInPartialMM")
	// The host is in the process of entering the particular partial maintenance
	// mode.
	HostPartialMaintenanceModeStatusEnteringPartialMM = HostPartialMaintenanceModeStatus("enteringPartialMM")
	// The host is in the process of exiting the particular partial maintenance
	// mode.
	HostPartialMaintenanceModeStatusExitingPartialMM = HostPartialMaintenanceModeStatus("exitingPartialMM")
	// The host is in the particular partial maintenance mode.
	HostPartialMaintenanceModeStatusInPartialMM = HostPartialMaintenanceModeStatus("inPartialMM")
)

func (e HostPartialMaintenanceModeStatus) Values() []HostPartialMaintenanceModeStatus {
	return []HostPartialMaintenanceModeStatus{
		HostPartialMaintenanceModeStatusNotInPartialMM,
		HostPartialMaintenanceModeStatusEnteringPartialMM,
		HostPartialMaintenanceModeStatusExitingPartialMM,
		HostPartialMaintenanceModeStatusInPartialMM,
	}
}

func (e HostPartialMaintenanceModeStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPartialMaintenanceModeStatus"] = reflect.TypeOf((*HostPartialMaintenanceModeStatus)(nil)).Elem()
	minAPIVersionForType["HostPartialMaintenanceModeStatus"] = "8.0.3.0"
}

// The installation state if the update is installed on the server.
type HostPatchManagerInstallState string

const (
	// The server has been restarted since the update installation.
	HostPatchManagerInstallStateHostRestarted = HostPatchManagerInstallState("hostRestarted")
	// Indicates if the newly installed image is active on the server
	HostPatchManagerInstallStateImageActive = HostPatchManagerInstallState("imageActive")
)

func (e HostPatchManagerInstallState) Values() []HostPatchManagerInstallState {
	return []HostPatchManagerInstallState{
		HostPatchManagerInstallStateHostRestarted,
		HostPatchManagerInstallStateImageActive,
	}
}

func (e HostPatchManagerInstallState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPatchManagerInstallState"] = reflect.TypeOf((*HostPatchManagerInstallState)(nil)).Elem()
}

// The integrity validation status.
type HostPatchManagerIntegrityStatus string

const (
	// The update is successfully validated.
	HostPatchManagerIntegrityStatusValidated = HostPatchManagerIntegrityStatus("validated")
	// The integrity can not be verified since a public key to
	// verify the update cannot be found.
	HostPatchManagerIntegrityStatusKeyNotFound = HostPatchManagerIntegrityStatus("keyNotFound")
	// A public key to verify the update has been revoked.
	HostPatchManagerIntegrityStatusKeyRevoked = HostPatchManagerIntegrityStatus("keyRevoked")
	// A public key to verify the update is expired.
	HostPatchManagerIntegrityStatusKeyExpired = HostPatchManagerIntegrityStatus("keyExpired")
	// A digital signature of the update does not match.
	HostPatchManagerIntegrityStatusDigestMismatch = HostPatchManagerIntegrityStatus("digestMismatch")
	// Not enough signed signatures on the update.
	HostPatchManagerIntegrityStatusNotEnoughSignatures = HostPatchManagerIntegrityStatus("notEnoughSignatures")
	// The integrity validation failed.
	HostPatchManagerIntegrityStatusValidationError = HostPatchManagerIntegrityStatus("validationError")
)

func (e HostPatchManagerIntegrityStatus) Values() []HostPatchManagerIntegrityStatus {
	return []HostPatchManagerIntegrityStatus{
		HostPatchManagerIntegrityStatusValidated,
		HostPatchManagerIntegrityStatusKeyNotFound,
		HostPatchManagerIntegrityStatusKeyRevoked,
		HostPatchManagerIntegrityStatusKeyExpired,
		HostPatchManagerIntegrityStatusDigestMismatch,
		HostPatchManagerIntegrityStatusNotEnoughSignatures,
		HostPatchManagerIntegrityStatusValidationError,
	}
}

func (e HostPatchManagerIntegrityStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPatchManagerIntegrityStatus"] = reflect.TypeOf((*HostPatchManagerIntegrityStatus)(nil)).Elem()
}

// Reasons why an update is not applicable to the ESX host.
type HostPatchManagerReason string

const (
	// The update is made obsolete by other patches installed on the host.
	HostPatchManagerReasonObsoleted = HostPatchManagerReason("obsoleted")
	// The update depends on another update that is neither installed
	// nor in the scanned list of updates.
	HostPatchManagerReasonMissingPatch = HostPatchManagerReason("missingPatch")
	// The update depends on certain libraries or RPMs that are not
	// available.
	HostPatchManagerReasonMissingLib = HostPatchManagerReason("missingLib")
	// The update depends on an update that is not installed but is
	// in the scanned list of updates.
	HostPatchManagerReasonHasDependentPatch = HostPatchManagerReason("hasDependentPatch")
	// The update conflicts with certain updates that are already
	// installed on the host.
	HostPatchManagerReasonConflictPatch = HostPatchManagerReason("conflictPatch")
	// The update conflicts with RPMs or libraries installed on the
	// host.
	HostPatchManagerReasonConflictLib = HostPatchManagerReason("conflictLib")
)

func (e HostPatchManagerReason) Values() []HostPatchManagerReason {
	return []HostPatchManagerReason{
		HostPatchManagerReasonObsoleted,
		HostPatchManagerReasonMissingPatch,
		HostPatchManagerReasonMissingLib,
		HostPatchManagerReasonHasDependentPatch,
		HostPatchManagerReasonConflictPatch,
		HostPatchManagerReasonConflictLib,
	}
}

func (e HostPatchManagerReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPatchManagerReason"] = reflect.TypeOf((*HostPatchManagerReason)(nil)).Elem()
}

type HostPowerOperationType string

const (
	// Power On Operation
	HostPowerOperationTypePowerOn = HostPowerOperationType("powerOn")
	// Power Off Operation.
	//
	// Power off operation puts the host in
	// a state that can be woken up remotely.
	HostPowerOperationTypePowerOff = HostPowerOperationType("powerOff")
)

func (e HostPowerOperationType) Values() []HostPowerOperationType {
	return []HostPowerOperationType{
		HostPowerOperationTypePowerOn,
		HostPowerOperationTypePowerOff,
	}
}

func (e HostPowerOperationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPowerOperationType"] = reflect.TypeOf((*HostPowerOperationType)(nil)).Elem()
}

// The `HostProfileManagerAnswerFileStatus_enum` enum
// defines possible values for answer file status.
type HostProfileManagerAnswerFileStatus string

const (
	// Answer file is valid.
	HostProfileManagerAnswerFileStatusValid = HostProfileManagerAnswerFileStatus("valid")
	// Answer file is not valid.
	//
	// The file is either missing or incomplete.
	//   - To produce an answer file, pass host-specific data (user input) to the
	//     `HostProfileManager*.*HostProfileManager.ApplyHostConfig_Task`
	//     method.
	//   - To produce a complete answer file, call the
	//     `HostProfile*.*HostProfile.ExecuteHostProfile`
	//     method and fill in any missing parameters in the returned
	//     `ProfileExecuteResult*.*ProfileExecuteResult.requireInput`
	//     list. After you execute the profile successfully, you can pass the complete required
	//     input list to the apply method.
	HostProfileManagerAnswerFileStatusInvalid = HostProfileManagerAnswerFileStatus("invalid")
	// Answer file status is not known.
	HostProfileManagerAnswerFileStatusUnknown = HostProfileManagerAnswerFileStatus("unknown")
)

func (e HostProfileManagerAnswerFileStatus) Values() []HostProfileManagerAnswerFileStatus {
	return []HostProfileManagerAnswerFileStatus{
		HostProfileManagerAnswerFileStatusValid,
		HostProfileManagerAnswerFileStatusInvalid,
		HostProfileManagerAnswerFileStatusUnknown,
	}
}

func (e HostProfileManagerAnswerFileStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProfileManagerAnswerFileStatus"] = reflect.TypeOf((*HostProfileManagerAnswerFileStatus)(nil)).Elem()
}

// The composition status class.
type HostProfileManagerCompositionResultResultElementStatus string

const (
	HostProfileManagerCompositionResultResultElementStatusSuccess = HostProfileManagerCompositionResultResultElementStatus("success")
	HostProfileManagerCompositionResultResultElementStatusError   = HostProfileManagerCompositionResultResultElementStatus("error")
)

func (e HostProfileManagerCompositionResultResultElementStatus) Values() []HostProfileManagerCompositionResultResultElementStatus {
	return []HostProfileManagerCompositionResultResultElementStatus{
		HostProfileManagerCompositionResultResultElementStatusSuccess,
		HostProfileManagerCompositionResultResultElementStatusError,
	}
}

func (e HostProfileManagerCompositionResultResultElementStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProfileManagerCompositionResultResultElementStatus"] = reflect.TypeOf((*HostProfileManagerCompositionResultResultElementStatus)(nil)).Elem()
}

// The composition validation status class.
type HostProfileManagerCompositionValidationResultResultElementStatus string

const (
	HostProfileManagerCompositionValidationResultResultElementStatusSuccess = HostProfileManagerCompositionValidationResultResultElementStatus("success")
	HostProfileManagerCompositionValidationResultResultElementStatusError   = HostProfileManagerCompositionValidationResultResultElementStatus("error")
)

func (e HostProfileManagerCompositionValidationResultResultElementStatus) Values() []HostProfileManagerCompositionValidationResultResultElementStatus {
	return []HostProfileManagerCompositionValidationResultResultElementStatus{
		HostProfileManagerCompositionValidationResultResultElementStatusSuccess,
		HostProfileManagerCompositionValidationResultResultElementStatusError,
	}
}

func (e HostProfileManagerCompositionValidationResultResultElementStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProfileManagerCompositionValidationResultResultElementStatus"] = reflect.TypeOf((*HostProfileManagerCompositionValidationResultResultElementStatus)(nil)).Elem()
}

// The `HostProfileManagerTaskListRequirement_enum` enum
// defines possible values for requirements when applying a `HostConfigSpec`
// object returned as part of a <code>generateConfigTaskList</code>
// operation.
type HostProfileManagerTaskListRequirement string

const (
	// The ESXi host must be in maintenance mode before the task list can be
	// applied.
	HostProfileManagerTaskListRequirementMaintenanceModeRequired = HostProfileManagerTaskListRequirement("maintenanceModeRequired")
	// The ESXi host must be rebooted after the task list is applied in order
	// for the new settings in the `HostConfigSpec` to take
	// effect on the host.
	HostProfileManagerTaskListRequirementRebootRequired = HostProfileManagerTaskListRequirement("rebootRequired")
)

func (e HostProfileManagerTaskListRequirement) Values() []HostProfileManagerTaskListRequirement {
	return []HostProfileManagerTaskListRequirement{
		HostProfileManagerTaskListRequirementMaintenanceModeRequired,
		HostProfileManagerTaskListRequirementRebootRequired,
	}
}

func (e HostProfileManagerTaskListRequirement) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProfileManagerTaskListRequirement"] = reflect.TypeOf((*HostProfileManagerTaskListRequirement)(nil)).Elem()
}

// Types of host profile update.
type HostProfileValidationFailureInfoUpdateType string

const (
	// Update host profile from host.
	HostProfileValidationFailureInfoUpdateTypeHostBased = HostProfileValidationFailureInfoUpdateType("HostBased")
	// Import host profile.
	HostProfileValidationFailureInfoUpdateTypeImport = HostProfileValidationFailureInfoUpdateType("Import")
	// Edit host profile.
	HostProfileValidationFailureInfoUpdateTypeEdit = HostProfileValidationFailureInfoUpdateType("Edit")
	// Compose setting from host profile.
	HostProfileValidationFailureInfoUpdateTypeCompose = HostProfileValidationFailureInfoUpdateType("Compose")
)

func (e HostProfileValidationFailureInfoUpdateType) Values() []HostProfileValidationFailureInfoUpdateType {
	return []HostProfileValidationFailureInfoUpdateType{
		HostProfileValidationFailureInfoUpdateTypeHostBased,
		HostProfileValidationFailureInfoUpdateTypeImport,
		HostProfileValidationFailureInfoUpdateTypeEdit,
		HostProfileValidationFailureInfoUpdateTypeCompose,
	}
}

func (e HostProfileValidationFailureInfoUpdateType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProfileValidationFailureInfoUpdateType"] = reflect.TypeOf((*HostProfileValidationFailureInfoUpdateType)(nil)).Elem()
}

// This defines validation state values for host profile.
type HostProfileValidationState string

const (
	HostProfileValidationStateReady   = HostProfileValidationState("Ready")
	HostProfileValidationStateRunning = HostProfileValidationState("Running")
	HostProfileValidationStateFailed  = HostProfileValidationState("Failed")
)

func (e HostProfileValidationState) Values() []HostProfileValidationState {
	return []HostProfileValidationState{
		HostProfileValidationStateReady,
		HostProfileValidationStateRunning,
		HostProfileValidationStateFailed,
	}
}

func (e HostProfileValidationState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProfileValidationState"] = reflect.TypeOf((*HostProfileValidationState)(nil)).Elem()
}

// Deprecated from all vmodl version above @released("6.0").
//
// ProtocolEndpoint Type.
type HostProtocolEndpointPEType string

const (
	HostProtocolEndpointPETypeBlock = HostProtocolEndpointPEType("block")
	HostProtocolEndpointPETypeNas   = HostProtocolEndpointPEType("nas")
)

func (e HostProtocolEndpointPEType) Values() []HostProtocolEndpointPEType {
	return []HostProtocolEndpointPEType{
		HostProtocolEndpointPETypeBlock,
		HostProtocolEndpointPETypeNas,
	}
}

func (e HostProtocolEndpointPEType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProtocolEndpointPEType"] = reflect.TypeOf((*HostProtocolEndpointPEType)(nil)).Elem()
}

// ProtocolEndpoint type.
type HostProtocolEndpointProtocolEndpointType string

const (
	HostProtocolEndpointProtocolEndpointTypeScsi  = HostProtocolEndpointProtocolEndpointType("scsi")
	HostProtocolEndpointProtocolEndpointTypeNfs   = HostProtocolEndpointProtocolEndpointType("nfs")
	HostProtocolEndpointProtocolEndpointTypeNfs4x = HostProtocolEndpointProtocolEndpointType("nfs4x")
)

func (e HostProtocolEndpointProtocolEndpointType) Values() []HostProtocolEndpointProtocolEndpointType {
	return []HostProtocolEndpointProtocolEndpointType{
		HostProtocolEndpointProtocolEndpointTypeScsi,
		HostProtocolEndpointProtocolEndpointTypeNfs,
		HostProtocolEndpointProtocolEndpointTypeNfs4x,
	}
}

func (e HostProtocolEndpointProtocolEndpointType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostProtocolEndpointProtocolEndpointType"] = reflect.TypeOf((*HostProtocolEndpointProtocolEndpointType)(nil)).Elem()
}

type HostPtpConfigDeviceType string

const (
	// No device.
	HostPtpConfigDeviceTypeNone = HostPtpConfigDeviceType("none")
	// Virtual network adapter.
	HostPtpConfigDeviceTypeVirtualNic = HostPtpConfigDeviceType("virtualNic")
	// A network PCI device capable of PTP hardware timestamping,
	// enabled for passthru.
	//
	// See `HostPciPassthruSystem`
	// for information on PCI devices enabled for passthru available
	// on the host.
	HostPtpConfigDeviceTypePciPassthruNic = HostPtpConfigDeviceType("pciPassthruNic")
)

func (e HostPtpConfigDeviceType) Values() []HostPtpConfigDeviceType {
	return []HostPtpConfigDeviceType{
		HostPtpConfigDeviceTypeNone,
		HostPtpConfigDeviceTypeVirtualNic,
		HostPtpConfigDeviceTypePciPassthruNic,
	}
}

func (e HostPtpConfigDeviceType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostPtpConfigDeviceType"] = reflect.TypeOf((*HostPtpConfigDeviceType)(nil)).Elem()
	minAPIVersionForType["HostPtpConfigDeviceType"] = "7.0.3.0"
}

type HostQualifiedNameType string

const (
	// The NVMe Qualified Name (NQN) of this host.
	HostQualifiedNameTypeNvmeQualifiedName = HostQualifiedNameType("nvmeQualifiedName")
	// The NVMe Qualified Name (NQN) of this host used by Vvol.
	HostQualifiedNameTypeVvolNvmeQualifiedName = HostQualifiedNameType("vvolNvmeQualifiedName")
)

func (e HostQualifiedNameType) Values() []HostQualifiedNameType {
	return []HostQualifiedNameType{
		HostQualifiedNameTypeNvmeQualifiedName,
		HostQualifiedNameTypeVvolNvmeQualifiedName,
	}
}

func (e HostQualifiedNameType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostQualifiedNameType"] = reflect.TypeOf((*HostQualifiedNameType)(nil)).Elem()
	minAPIVersionForType["HostQualifiedNameType"] = "7.0.3.0"
	minAPIVersionForEnumValue["HostQualifiedNameType"] = map[string]string{
		"vvolNvmeQualifiedName": "8.0.0.0",
	}
}

// Possible RDMA device connection states.
//
// These correspond
// to possible link states as defined by the
// Infiniband (TM) specification.
//
// Further details can be found in:
//   - "Infiniband (TM) Architecture Specification, Volume 1"
//     section 7.2 "Link states"
type HostRdmaDeviceConnectionState string

const (
	// Connection state unknown.
	//
	// Indicates that the driver returned
	// unexpected or no connection state information.
	HostRdmaDeviceConnectionStateUnknown = HostRdmaDeviceConnectionState("unknown")
	// Device down.
	//
	// Indicates that both the logical link and
	// underlying physical link are down. Packets
	// are discarded.
	HostRdmaDeviceConnectionStateDown = HostRdmaDeviceConnectionState("down")
	// Device initializing.
	//
	// Indicates that the physical link is up, but
	// the logical link is still initializing.
	// Only subnet management and flow control link
	// packets can be received and transmitted.
	HostRdmaDeviceConnectionStateInit = HostRdmaDeviceConnectionState("init")
	// Device armed.
	//
	// Indicates that the physical link is up, but
	// the logical link is not yet fully configured.
	// Packets can be received, but non-SMPs
	// (subnet management packets) to be sent are discarded.
	HostRdmaDeviceConnectionStateArmed = HostRdmaDeviceConnectionState("armed")
	// Device active.
	//
	// Indicates that both the physical and logical
	// link are up. Packets can be transmitted and received.
	HostRdmaDeviceConnectionStateActive = HostRdmaDeviceConnectionState("active")
	// Device in active defer state.
	//
	// Indicates that the logical link was active, but the
	// physical link has suffered a failure. If it recovers
	// within a timeout, the connection state will return to active,
	// otherwise it will move to down.
	HostRdmaDeviceConnectionStateActiveDefer = HostRdmaDeviceConnectionState("activeDefer")
)

func (e HostRdmaDeviceConnectionState) Values() []HostRdmaDeviceConnectionState {
	return []HostRdmaDeviceConnectionState{
		HostRdmaDeviceConnectionStateUnknown,
		HostRdmaDeviceConnectionStateDown,
		HostRdmaDeviceConnectionStateInit,
		HostRdmaDeviceConnectionStateArmed,
		HostRdmaDeviceConnectionStateActive,
		HostRdmaDeviceConnectionStateActiveDefer,
	}
}

func (e HostRdmaDeviceConnectionState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostRdmaDeviceConnectionState"] = reflect.TypeOf((*HostRdmaDeviceConnectionState)(nil)).Elem()
}

// Deprecated as of vSphere API 6.0.
//
// Set of possible values for
// `HostCapability.replayUnsupportedReason` and
// `HostCapability.replayCompatibilityIssues`.
type HostReplayUnsupportedReason string

const (
	HostReplayUnsupportedReasonIncompatibleProduct = HostReplayUnsupportedReason("incompatibleProduct")
	HostReplayUnsupportedReasonIncompatibleCpu     = HostReplayUnsupportedReason("incompatibleCpu")
	HostReplayUnsupportedReasonHvDisabled          = HostReplayUnsupportedReason("hvDisabled")
	HostReplayUnsupportedReasonCpuidLimitSet       = HostReplayUnsupportedReason("cpuidLimitSet")
	HostReplayUnsupportedReasonOldBIOS             = HostReplayUnsupportedReason("oldBIOS")
	HostReplayUnsupportedReasonUnknown             = HostReplayUnsupportedReason("unknown")
)

func (e HostReplayUnsupportedReason) Values() []HostReplayUnsupportedReason {
	return []HostReplayUnsupportedReason{
		HostReplayUnsupportedReasonIncompatibleProduct,
		HostReplayUnsupportedReasonIncompatibleCpu,
		HostReplayUnsupportedReasonHvDisabled,
		HostReplayUnsupportedReasonCpuidLimitSet,
		HostReplayUnsupportedReasonOldBIOS,
		HostReplayUnsupportedReasonUnknown,
	}
}

func (e HostReplayUnsupportedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostReplayUnsupportedReason"] = reflect.TypeOf((*HostReplayUnsupportedReason)(nil)).Elem()
}

// Define the instance state type
type HostRuntimeInfoNetStackInstanceRuntimeInfoState string

const (
	// The instance is deleted or not running
	HostRuntimeInfoNetStackInstanceRuntimeInfoStateInactive = HostRuntimeInfoNetStackInstanceRuntimeInfoState("inactive")
	// The instance is running
	HostRuntimeInfoNetStackInstanceRuntimeInfoStateActive = HostRuntimeInfoNetStackInstanceRuntimeInfoState("active")
	// The instance is in the progress of asynchronous deletion
	HostRuntimeInfoNetStackInstanceRuntimeInfoStateDeactivating = HostRuntimeInfoNetStackInstanceRuntimeInfoState("deactivating")
	// Reserved state for future proofing asynchronous creation
	HostRuntimeInfoNetStackInstanceRuntimeInfoStateActivating = HostRuntimeInfoNetStackInstanceRuntimeInfoState("activating")
)

func (e HostRuntimeInfoNetStackInstanceRuntimeInfoState) Values() []HostRuntimeInfoNetStackInstanceRuntimeInfoState {
	return []HostRuntimeInfoNetStackInstanceRuntimeInfoState{
		HostRuntimeInfoNetStackInstanceRuntimeInfoStateInactive,
		HostRuntimeInfoNetStackInstanceRuntimeInfoStateActive,
		HostRuntimeInfoNetStackInstanceRuntimeInfoStateDeactivating,
		HostRuntimeInfoNetStackInstanceRuntimeInfoStateActivating,
	}
}

func (e HostRuntimeInfoNetStackInstanceRuntimeInfoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostRuntimeInfoNetStackInstanceRuntimeInfoState"] = reflect.TypeOf((*HostRuntimeInfoNetStackInstanceRuntimeInfoState)(nil)).Elem()
}

type HostRuntimeInfoStateEncryptionInfoProtectionMode string

const (
	// Encryption is not protected.
	HostRuntimeInfoStateEncryptionInfoProtectionModeNone = HostRuntimeInfoStateEncryptionInfoProtectionMode("none")
	// Encryption is TPM protected.
	HostRuntimeInfoStateEncryptionInfoProtectionModeTpm = HostRuntimeInfoStateEncryptionInfoProtectionMode("tpm")
)

func (e HostRuntimeInfoStateEncryptionInfoProtectionMode) Values() []HostRuntimeInfoStateEncryptionInfoProtectionMode {
	return []HostRuntimeInfoStateEncryptionInfoProtectionMode{
		HostRuntimeInfoStateEncryptionInfoProtectionModeNone,
		HostRuntimeInfoStateEncryptionInfoProtectionModeTpm,
	}
}

func (e HostRuntimeInfoStateEncryptionInfoProtectionMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostRuntimeInfoStateEncryptionInfoProtectionMode"] = reflect.TypeOf((*HostRuntimeInfoStateEncryptionInfoProtectionMode)(nil)).Elem()
	minAPIVersionForType["HostRuntimeInfoStateEncryptionInfoProtectionMode"] = "7.0.3.0"
}

type HostRuntimeInfoStatelessNvdsMigrationState string

const (
	// The host is ready for NVDS to VDS migration.
	HostRuntimeInfoStatelessNvdsMigrationStateReady = HostRuntimeInfoStatelessNvdsMigrationState("ready")
	// The host does not need NVDS to VDS migration
	HostRuntimeInfoStatelessNvdsMigrationStateNotNeeded = HostRuntimeInfoStatelessNvdsMigrationState("notNeeded")
	// The host is disconnected from VC.
	HostRuntimeInfoStatelessNvdsMigrationStateUnknown = HostRuntimeInfoStatelessNvdsMigrationState("unknown")
)

func (e HostRuntimeInfoStatelessNvdsMigrationState) Values() []HostRuntimeInfoStatelessNvdsMigrationState {
	return []HostRuntimeInfoStatelessNvdsMigrationState{
		HostRuntimeInfoStatelessNvdsMigrationStateReady,
		HostRuntimeInfoStatelessNvdsMigrationStateNotNeeded,
		HostRuntimeInfoStatelessNvdsMigrationStateUnknown,
	}
}

func (e HostRuntimeInfoStatelessNvdsMigrationState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostRuntimeInfoStatelessNvdsMigrationState"] = reflect.TypeOf((*HostRuntimeInfoStatelessNvdsMigrationState)(nil)).Elem()
	minAPIVersionForType["HostRuntimeInfoStatelessNvdsMigrationState"] = "7.0.2.0"
}

// Set of valid service policy strings.
type HostServicePolicy string

const (
	// Service should be started when the host starts up.
	HostServicePolicyOn = HostServicePolicy("on")
	// Service should run if and only if it has open firewall ports.
	HostServicePolicyAutomatic = HostServicePolicy("automatic")
	// Service should not be started when the host starts up.
	HostServicePolicyOff = HostServicePolicy("off")
)

func (e HostServicePolicy) Values() []HostServicePolicy {
	return []HostServicePolicy{
		HostServicePolicyOn,
		HostServicePolicyAutomatic,
		HostServicePolicyOff,
	}
}

func (e HostServicePolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostServicePolicy"] = reflect.TypeOf((*HostServicePolicy)(nil)).Elem()
}

type HostSevInfoSevState string

const (
	HostSevInfoSevStateUninitialized = HostSevInfoSevState("uninitialized")
	HostSevInfoSevStateInitialized   = HostSevInfoSevState("initialized")
	HostSevInfoSevStateWorking       = HostSevInfoSevState("working")
)

func (e HostSevInfoSevState) Values() []HostSevInfoSevState {
	return []HostSevInfoSevState{
		HostSevInfoSevStateUninitialized,
		HostSevInfoSevStateInitialized,
		HostSevInfoSevStateWorking,
	}
}

func (e HostSevInfoSevState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSevInfoSevState"] = reflect.TypeOf((*HostSevInfoSevState)(nil)).Elem()
	minAPIVersionForType["HostSevInfoSevState"] = "7.0.1.0"
}

// Flexible Launch Enclave (FLC) modes.
type HostSgxInfoFlcModes string

const (
	// Flexible Launch Enclave (FLC) is not available on the host.
	//
	// The
	// "launch enclave MSRs" are initialized with Intel's public key hash.
	HostSgxInfoFlcModesOff = HostSgxInfoFlcModes("off")
	// FLC is available and the "launch Enclave MSRs" are locked and
	// initialized with the provided public key hash.
	HostSgxInfoFlcModesLocked = HostSgxInfoFlcModes("locked")
	// FLC is available and the "launch enclave MSRs" are writeable and
	// initialized with Intel's public key hash.
	HostSgxInfoFlcModesUnlocked = HostSgxInfoFlcModes("unlocked")
)

func (e HostSgxInfoFlcModes) Values() []HostSgxInfoFlcModes {
	return []HostSgxInfoFlcModes{
		HostSgxInfoFlcModesOff,
		HostSgxInfoFlcModesLocked,
		HostSgxInfoFlcModesUnlocked,
	}
}

func (e HostSgxInfoFlcModes) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSgxInfoFlcModes"] = reflect.TypeOf((*HostSgxInfoFlcModes)(nil)).Elem()
}

// Host SGX states.
type HostSgxInfoSgxStates string

const (
	// SGX is not present in the CPU.
	HostSgxInfoSgxStatesNotPresent = HostSgxInfoSgxStates("notPresent")
	// SGX is disabled in the BIOS.
	HostSgxInfoSgxStatesDisabledBIOS = HostSgxInfoSgxStates("disabledBIOS")
	// SGX is disabled because CPU erratum CFW101 is present.
	HostSgxInfoSgxStatesDisabledCFW101 = HostSgxInfoSgxStates("disabledCFW101")
	// SGX is disabled due to a mismatch in the SGX capabilities
	// exposed by different CPUs.
	HostSgxInfoSgxStatesDisabledCPUMismatch = HostSgxInfoSgxStates("disabledCPUMismatch")
	// SGX is disabled because the CPU does not support FLC.
	HostSgxInfoSgxStatesDisabledNoFLC = HostSgxInfoSgxStates("disabledNoFLC")
	// SGX is disabled because the host uses NUMA, which is not
	// supported with SGX.
	HostSgxInfoSgxStatesDisabledNUMAUnsup = HostSgxInfoSgxStates("disabledNUMAUnsup")
	// SGX is disabled because the host exceeds the maximum supported
	// number of EPC regions.
	HostSgxInfoSgxStatesDisabledMaxEPCRegs = HostSgxInfoSgxStates("disabledMaxEPCRegs")
	// SGX is enabled.
	HostSgxInfoSgxStatesEnabled = HostSgxInfoSgxStates("enabled")
)

func (e HostSgxInfoSgxStates) Values() []HostSgxInfoSgxStates {
	return []HostSgxInfoSgxStates{
		HostSgxInfoSgxStatesNotPresent,
		HostSgxInfoSgxStatesDisabledBIOS,
		HostSgxInfoSgxStatesDisabledCFW101,
		HostSgxInfoSgxStatesDisabledCPUMismatch,
		HostSgxInfoSgxStatesDisabledNoFLC,
		HostSgxInfoSgxStatesDisabledNUMAUnsup,
		HostSgxInfoSgxStatesDisabledMaxEPCRegs,
		HostSgxInfoSgxStatesEnabled,
	}
}

func (e HostSgxInfoSgxStates) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSgxInfoSgxStates"] = reflect.TypeOf((*HostSgxInfoSgxStates)(nil)).Elem()
}

type HostSgxRegistrationInfoRegistrationStatus string

const (
	// SGX is not available or the host is unisocket.
	HostSgxRegistrationInfoRegistrationStatusNotApplicable = HostSgxRegistrationInfoRegistrationStatus("notApplicable")
	// SGX registration is incomplete.
	HostSgxRegistrationInfoRegistrationStatusIncomplete = HostSgxRegistrationInfoRegistrationStatus("incomplete")
	// SGX registration is complete.
	HostSgxRegistrationInfoRegistrationStatusComplete = HostSgxRegistrationInfoRegistrationStatus("complete")
)

func (e HostSgxRegistrationInfoRegistrationStatus) Values() []HostSgxRegistrationInfoRegistrationStatus {
	return []HostSgxRegistrationInfoRegistrationStatus{
		HostSgxRegistrationInfoRegistrationStatusNotApplicable,
		HostSgxRegistrationInfoRegistrationStatusIncomplete,
		HostSgxRegistrationInfoRegistrationStatusComplete,
	}
}

func (e HostSgxRegistrationInfoRegistrationStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSgxRegistrationInfoRegistrationStatus"] = reflect.TypeOf((*HostSgxRegistrationInfoRegistrationStatus)(nil)).Elem()
	minAPIVersionForType["HostSgxRegistrationInfoRegistrationStatus"] = "8.0.0.1"
}

type HostSgxRegistrationInfoRegistrationType string

const (
	// Indicates that an Initial Platform Establishment
	// or TCB recovery registration is pending.
	HostSgxRegistrationInfoRegistrationTypeManifest = HostSgxRegistrationInfoRegistrationType("manifest")
	// Indicates that new CPU package was added.
	HostSgxRegistrationInfoRegistrationTypeAddPackage = HostSgxRegistrationInfoRegistrationType("addPackage")
)

func (e HostSgxRegistrationInfoRegistrationType) Values() []HostSgxRegistrationInfoRegistrationType {
	return []HostSgxRegistrationInfoRegistrationType{
		HostSgxRegistrationInfoRegistrationTypeManifest,
		HostSgxRegistrationInfoRegistrationTypeAddPackage,
	}
}

func (e HostSgxRegistrationInfoRegistrationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSgxRegistrationInfoRegistrationType"] = reflect.TypeOf((*HostSgxRegistrationInfoRegistrationType)(nil)).Elem()
	minAPIVersionForType["HostSgxRegistrationInfoRegistrationType"] = "8.0.0.1"
}

// SNMP Agent supported capabilities enum
type HostSnmpAgentCapability string

const (
	// Implements test notifications and allows agent configuration
	HostSnmpAgentCapabilityCOMPLETE = HostSnmpAgentCapability("COMPLETE")
	// Implements only test notification capability only
	HostSnmpAgentCapabilityDIAGNOSTICS = HostSnmpAgentCapability("DIAGNOSTICS")
	// Allows for agent configuration only
	HostSnmpAgentCapabilityCONFIGURATION = HostSnmpAgentCapability("CONFIGURATION")
)

func (e HostSnmpAgentCapability) Values() []HostSnmpAgentCapability {
	return []HostSnmpAgentCapability{
		HostSnmpAgentCapabilityCOMPLETE,
		HostSnmpAgentCapabilityDIAGNOSTICS,
		HostSnmpAgentCapabilityCONFIGURATION,
	}
}

func (e HostSnmpAgentCapability) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSnmpAgentCapability"] = reflect.TypeOf((*HostSnmpAgentCapability)(nil)).Elem()
}

// Defines a host's standby mode.
type HostStandbyMode string

const (
	// The host is entering standby mode.
	HostStandbyModeEntering = HostStandbyMode("entering")
	// The host is exiting standby mode.
	HostStandbyModeExiting = HostStandbyMode("exiting")
	// The host is in standby mode.
	HostStandbyModeIn = HostStandbyMode("in")
	// The host is not in standy mode, and it is not
	// in the process of entering/exiting standby mode.
	HostStandbyModeNone = HostStandbyMode("none")
)

func (e HostStandbyMode) Values() []HostStandbyMode {
	return []HostStandbyMode{
		HostStandbyModeEntering,
		HostStandbyModeExiting,
		HostStandbyModeIn,
		HostStandbyModeNone,
	}
}

func (e HostStandbyMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostStandbyMode"] = reflect.TypeOf((*HostStandbyMode)(nil)).Elem()
}

// The set of supported host bus adapter protocols.
type HostStorageProtocol string

const (
	// The Small Computer System Interface (SCSI) protocol.
	HostStorageProtocolScsi = HostStorageProtocol("scsi")
	// The Non-Volatile Memory Express (NVME) protocol.
	HostStorageProtocolNvme = HostStorageProtocol("nvme")
)

func (e HostStorageProtocol) Values() []HostStorageProtocol {
	return []HostStorageProtocol{
		HostStorageProtocolScsi,
		HostStorageProtocolNvme,
	}
}

func (e HostStorageProtocol) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostStorageProtocol"] = reflect.TypeOf((*HostStorageProtocol)(nil)).Elem()
}

// Defines a host's connection state.
type HostSystemConnectionState string

const (
	// Connected to the server.
	//
	// For ESX Server, this is always the setting.
	HostSystemConnectionStateConnected = HostSystemConnectionState("connected")
	// VirtualCenter is not receiving heartbeats from the server.
	//
	// The state
	// automatically changes to connected once heartbeats are received
	// again. This state is typically used to trigger an alarm on the host.
	HostSystemConnectionStateNotResponding = HostSystemConnectionState("notResponding")
	// The user has explicitly taken the host down.
	//
	// VirtualCenter does not expect to
	// receive heartbeats from the host. The next time a heartbeat is received, the
	// host is moved to the connected state again and an event is logged.
	HostSystemConnectionStateDisconnected = HostSystemConnectionState("disconnected")
)

func (e HostSystemConnectionState) Values() []HostSystemConnectionState {
	return []HostSystemConnectionState{
		HostSystemConnectionStateConnected,
		HostSystemConnectionStateNotResponding,
		HostSystemConnectionStateDisconnected,
	}
}

func (e HostSystemConnectionState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSystemConnectionState"] = reflect.TypeOf((*HostSystemConnectionState)(nil)).Elem()
}

type HostSystemIdentificationInfoIdentifier string

const (
	// The Asset tag of the system
	HostSystemIdentificationInfoIdentifierAssetTag = HostSystemIdentificationInfoIdentifier("AssetTag")
	// The Service tag of the system
	HostSystemIdentificationInfoIdentifierServiceTag = HostSystemIdentificationInfoIdentifier("ServiceTag")
	// OEM specific string
	HostSystemIdentificationInfoIdentifierOemSpecificString = HostSystemIdentificationInfoIdentifier("OemSpecificString")
	// The Enclosure Serial Number tag of the system
	HostSystemIdentificationInfoIdentifierEnclosureSerialNumberTag = HostSystemIdentificationInfoIdentifier("EnclosureSerialNumberTag")
	// The Serial Number tag of the system
	HostSystemIdentificationInfoIdentifierSerialNumberTag = HostSystemIdentificationInfoIdentifier("SerialNumberTag")
)

func (e HostSystemIdentificationInfoIdentifier) Values() []HostSystemIdentificationInfoIdentifier {
	return []HostSystemIdentificationInfoIdentifier{
		HostSystemIdentificationInfoIdentifierAssetTag,
		HostSystemIdentificationInfoIdentifierServiceTag,
		HostSystemIdentificationInfoIdentifierOemSpecificString,
		HostSystemIdentificationInfoIdentifierEnclosureSerialNumberTag,
		HostSystemIdentificationInfoIdentifierSerialNumberTag,
	}
}

func (e HostSystemIdentificationInfoIdentifier) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSystemIdentificationInfoIdentifier"] = reflect.TypeOf((*HostSystemIdentificationInfoIdentifier)(nil)).Elem()
}

// Defines a host's power state.
type HostSystemPowerState string

const (
	// The host is powered on.
	//
	// A host that is entering standby mode
	// `entering` is also in this state.
	HostSystemPowerStatePoweredOn = HostSystemPowerState("poweredOn")
	// The host was specifically powered off by the user through
	// VirtualCenter.
	//
	// This state is not a cetain state, because
	// after VirtualCenter issues the command to power off the host,
	// the host might crash, or kill all the processes but fail to
	// power off.
	HostSystemPowerStatePoweredOff = HostSystemPowerState("poweredOff")
	// The host was specifically put in standby mode, either
	// explicitly by the user, or automatically by DPM.
	//
	// This state
	// is not a cetain state, because after VirtualCenter issues the
	// command to put the host in standby state, the host might
	// crash, or kill all the processes but fail to power off. A host
	// that is exiting standby mode `exiting`
	// is also in this state.
	HostSystemPowerStateStandBy = HostSystemPowerState("standBy")
	// If the host is disconnected, or notResponding, we cannot
	// possibly have knowledge of its power state.
	//
	// Hence, the host
	// is marked as unknown.
	HostSystemPowerStateUnknown = HostSystemPowerState("unknown")
)

func (e HostSystemPowerState) Values() []HostSystemPowerState {
	return []HostSystemPowerState{
		HostSystemPowerStatePoweredOn,
		HostSystemPowerStatePoweredOff,
		HostSystemPowerStateStandBy,
		HostSystemPowerStateUnknown,
	}
}

func (e HostSystemPowerState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSystemPowerState"] = reflect.TypeOf((*HostSystemPowerState)(nil)).Elem()
}

// Valid state for host profile remediation.
type HostSystemRemediationStateState string

const (
	// Before precheck remediation and remediation.
	HostSystemRemediationStateStateRemediationReady = HostSystemRemediationStateState("remediationReady")
	// Preecheck remediation is running.
	HostSystemRemediationStateStatePrecheckRemediationRunning = HostSystemRemediationStateState("precheckRemediationRunning")
	// Preecheck remediation succeeded.
	HostSystemRemediationStateStatePrecheckRemediationComplete = HostSystemRemediationStateState("precheckRemediationComplete")
	// Preecheck remediation failed.
	HostSystemRemediationStateStatePrecheckRemediationFailed = HostSystemRemediationStateState("precheckRemediationFailed")
	// Remediation is running.
	HostSystemRemediationStateStateRemediationRunning = HostSystemRemediationStateState("remediationRunning")
	// Remediation failed.
	HostSystemRemediationStateStateRemediationFailed = HostSystemRemediationStateState("remediationFailed")
)

func (e HostSystemRemediationStateState) Values() []HostSystemRemediationStateState {
	return []HostSystemRemediationStateState{
		HostSystemRemediationStateStateRemediationReady,
		HostSystemRemediationStateStatePrecheckRemediationRunning,
		HostSystemRemediationStateStatePrecheckRemediationComplete,
		HostSystemRemediationStateStatePrecheckRemediationFailed,
		HostSystemRemediationStateStateRemediationRunning,
		HostSystemRemediationStateStateRemediationFailed,
	}
}

func (e HostSystemRemediationStateState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostSystemRemediationStateState"] = reflect.TypeOf((*HostSystemRemediationStateState)(nil)).Elem()
}

// Status constants of TPM attestation.
type HostTpmAttestationInfoAcceptanceStatus string

const (
	// TPM attestation failed.
	HostTpmAttestationInfoAcceptanceStatusNotAccepted = HostTpmAttestationInfoAcceptanceStatus("notAccepted")
	// TPM attestation succeeded.
	HostTpmAttestationInfoAcceptanceStatusAccepted = HostTpmAttestationInfoAcceptanceStatus("accepted")
)

func (e HostTpmAttestationInfoAcceptanceStatus) Values() []HostTpmAttestationInfoAcceptanceStatus {
	return []HostTpmAttestationInfoAcceptanceStatus{
		HostTpmAttestationInfoAcceptanceStatusNotAccepted,
		HostTpmAttestationInfoAcceptanceStatusAccepted,
	}
}

func (e HostTpmAttestationInfoAcceptanceStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostTpmAttestationInfoAcceptanceStatus"] = reflect.TypeOf((*HostTpmAttestationInfoAcceptanceStatus)(nil)).Elem()
}

type HostTrustAuthorityAttestationInfoAttestationStatus string

const (
	// Attestation succeeded.
	HostTrustAuthorityAttestationInfoAttestationStatusAttested = HostTrustAuthorityAttestationInfoAttestationStatus("attested")
	// Attestation failed.
	HostTrustAuthorityAttestationInfoAttestationStatusNotAttested = HostTrustAuthorityAttestationInfoAttestationStatus("notAttested")
	// Attestation status is unknown.
	HostTrustAuthorityAttestationInfoAttestationStatusUnknown = HostTrustAuthorityAttestationInfoAttestationStatus("unknown")
)

func (e HostTrustAuthorityAttestationInfoAttestationStatus) Values() []HostTrustAuthorityAttestationInfoAttestationStatus {
	return []HostTrustAuthorityAttestationInfoAttestationStatus{
		HostTrustAuthorityAttestationInfoAttestationStatusAttested,
		HostTrustAuthorityAttestationInfoAttestationStatusNotAttested,
		HostTrustAuthorityAttestationInfoAttestationStatusUnknown,
	}
}

func (e HostTrustAuthorityAttestationInfoAttestationStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostTrustAuthorityAttestationInfoAttestationStatus"] = reflect.TypeOf((*HostTrustAuthorityAttestationInfoAttestationStatus)(nil)).Elem()
	minAPIVersionForType["HostTrustAuthorityAttestationInfoAttestationStatus"] = "7.0.1.0"
}

// Reasons for identifying the disk extent
// as copy of VMFS volume extent.
type HostUnresolvedVmfsExtentUnresolvedReason string

const (
	// The VMFS detected 'diskid' does not match with
	// LVM detected 'diskId'
	HostUnresolvedVmfsExtentUnresolvedReasonDiskIdMismatch = HostUnresolvedVmfsExtentUnresolvedReason("diskIdMismatch")
	// VMFS 'uuid' does not match
	HostUnresolvedVmfsExtentUnresolvedReasonUuidConflict = HostUnresolvedVmfsExtentUnresolvedReason("uuidConflict")
)

func (e HostUnresolvedVmfsExtentUnresolvedReason) Values() []HostUnresolvedVmfsExtentUnresolvedReason {
	return []HostUnresolvedVmfsExtentUnresolvedReason{
		HostUnresolvedVmfsExtentUnresolvedReasonDiskIdMismatch,
		HostUnresolvedVmfsExtentUnresolvedReasonUuidConflict,
	}
}

func (e HostUnresolvedVmfsExtentUnresolvedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostUnresolvedVmfsExtentUnresolvedReason"] = reflect.TypeOf((*HostUnresolvedVmfsExtentUnresolvedReason)(nil)).Elem()
}

type HostUnresolvedVmfsResolutionSpecVmfsUuidResolution string

const (
	// Resignature the Unresolved VMFS volume.
	//
	// In the event the volume to be resignatured contains multiple
	// extents but only a single copy of each extent exists, only the
	// head extent needs to be specified.
	HostUnresolvedVmfsResolutionSpecVmfsUuidResolutionResignature = HostUnresolvedVmfsResolutionSpecVmfsUuidResolution("resignature")
	// Keep the original Uuid of the VMFS volume and mount it
	//
	// In the event the volume to be force mounted contains multiple
	// extents but only a single copy of each extent exists, only the
	// head extent needs to be specified.
	HostUnresolvedVmfsResolutionSpecVmfsUuidResolutionForceMount = HostUnresolvedVmfsResolutionSpecVmfsUuidResolution("forceMount")
)

func (e HostUnresolvedVmfsResolutionSpecVmfsUuidResolution) Values() []HostUnresolvedVmfsResolutionSpecVmfsUuidResolution {
	return []HostUnresolvedVmfsResolutionSpecVmfsUuidResolution{
		HostUnresolvedVmfsResolutionSpecVmfsUuidResolutionResignature,
		HostUnresolvedVmfsResolutionSpecVmfsUuidResolutionForceMount,
	}
}

func (e HostUnresolvedVmfsResolutionSpecVmfsUuidResolution) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostUnresolvedVmfsResolutionSpecVmfsUuidResolution"] = reflect.TypeOf((*HostUnresolvedVmfsResolutionSpecVmfsUuidResolution)(nil)).Elem()
}

type HostVirtualNicManagerNicType string

const (
	// The VirtualNic is used for VMotion.
	HostVirtualNicManagerNicTypeVmotion = HostVirtualNicManagerNicType("vmotion")
	// The VirtualNic is used for Fault Tolerance logging.
	HostVirtualNicManagerNicTypeFaultToleranceLogging = HostVirtualNicManagerNicType("faultToleranceLogging")
	// The VirtualNic is used for vSphere Replication LWD traffic
	// (i.e From the primary host to the VR server).
	HostVirtualNicManagerNicTypeVSphereReplication = HostVirtualNicManagerNicType("vSphereReplication")
	// The VirtualNic is used for vSphere Replication NFC traffic (i.e.
	//
	// From
	// the VR server to the secondary host).
	HostVirtualNicManagerNicTypeVSphereReplicationNFC = HostVirtualNicManagerNicType("vSphereReplicationNFC")
	// The VirtualNic is used for management network traffic .
	//
	// This nicType is available only when the system does not
	// support service console adapters.
	//
	// See also `HostNetCapabilities.usesServiceConsoleNic`.
	HostVirtualNicManagerNicTypeManagement = HostVirtualNicManagerNicType("management")
	// The VirtualNic is used for Virtual SAN data traffic.
	//
	// To enable or disable a VirtualNic for VSAN networking,
	// use `HostVsanSystem.UpdateVsan_Task`.
	//
	// See also `HostVsanSystem`, `HostVsanSystem.UpdateVsan_Task`, `ComputeResource.ReconfigureComputeResource_Task`.
	HostVirtualNicManagerNicTypeVsan = HostVirtualNicManagerNicType("vsan")
	// The VirtualNic is used for vSphere provisioning NFC traffic
	// (i.e.
	//
	// the NFC traffic between ESX hosts as a part of a VC initiated
	// provisioning operations like cold-migrations, clones, snapshot and
	// cold data in hot migration).
	HostVirtualNicManagerNicTypeVSphereProvisioning = HostVirtualNicManagerNicType("vSphereProvisioning")
	// The VirtualNic is used for Virtual SAN witness traffic.
	//
	// Witness traffic vmknic is required for Virtual SAN stretched cluster,
	// to help on communication between Virtual SAN data node and witness
	// node.
	// To enable or disable a VirtualNic for Virtual SAN networking,
	// use `HostVsanSystem.UpdateVsan_Task`.
	//
	// See also `HostVsanSystem`, `HostVsanSystem.UpdateVsan_Task`.
	HostVirtualNicManagerNicTypeVsanWitness = HostVirtualNicManagerNicType("vsanWitness")
	// The VirtualNic is used for vSphere backup NFC traffic
	// (i.e.
	//
	// the NFC traffic between backup appliance and ESX hosts).
	HostVirtualNicManagerNicTypeVSphereBackupNFC = HostVirtualNicManagerNicType("vSphereBackupNFC")
	// The VirtualNic is used for Precision Time Protocol (PTP).
	HostVirtualNicManagerNicTypePtp = HostVirtualNicManagerNicType("ptp")
	// The VirtualNic is used for NVMe over TCP traffic.
	HostVirtualNicManagerNicTypeNvmeTcp = HostVirtualNicManagerNicType("nvmeTcp")
	// The VirtualNic is used for NVMe over RDMA traffic.
	HostVirtualNicManagerNicTypeNvmeRdma = HostVirtualNicManagerNicType("nvmeRdma")
)

func (e HostVirtualNicManagerNicType) Values() []HostVirtualNicManagerNicType {
	return []HostVirtualNicManagerNicType{
		HostVirtualNicManagerNicTypeVmotion,
		HostVirtualNicManagerNicTypeFaultToleranceLogging,
		HostVirtualNicManagerNicTypeVSphereReplication,
		HostVirtualNicManagerNicTypeVSphereReplicationNFC,
		HostVirtualNicManagerNicTypeManagement,
		HostVirtualNicManagerNicTypeVsan,
		HostVirtualNicManagerNicTypeVSphereProvisioning,
		HostVirtualNicManagerNicTypeVsanWitness,
		HostVirtualNicManagerNicTypeVSphereBackupNFC,
		HostVirtualNicManagerNicTypePtp,
		HostVirtualNicManagerNicTypeNvmeTcp,
		HostVirtualNicManagerNicTypeNvmeRdma,
	}
}

func (e HostVirtualNicManagerNicType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostVirtualNicManagerNicType"] = reflect.TypeOf((*HostVirtualNicManagerNicType)(nil)).Elem()
	minAPIVersionForEnumValue["HostVirtualNicManagerNicType"] = map[string]string{
		"nvmeTcp":  "7.0.3.0",
		"nvmeRdma": "7.0.3.0",
	}
}

// Set of possible values for mode field in AccessSpec.
type HostVmciAccessManagerMode string

const (
	// Grant access to specified services in addition to existing services.
	HostVmciAccessManagerModeGrant = HostVmciAccessManagerMode("grant")
	// Replace existing services with specified services.
	HostVmciAccessManagerModeReplace = HostVmciAccessManagerMode("replace")
	// Revoke the specified services.
	HostVmciAccessManagerModeRevoke = HostVmciAccessManagerMode("revoke")
)

func (e HostVmciAccessManagerMode) Values() []HostVmciAccessManagerMode {
	return []HostVmciAccessManagerMode{
		HostVmciAccessManagerModeGrant,
		HostVmciAccessManagerModeReplace,
		HostVmciAccessManagerModeRevoke,
	}
}

func (e HostVmciAccessManagerMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostVmciAccessManagerMode"] = reflect.TypeOf((*HostVmciAccessManagerMode)(nil)).Elem()
}

// VMFS unmap bandwidth policy.
//
// VMFS unmap reclaims unused storage space.
// This specifies the bandwidth policy option of unmaps.
type HostVmfsVolumeUnmapBandwidthPolicy string

const (
	// Unmap bandwidth is a fixed value.
	HostVmfsVolumeUnmapBandwidthPolicyFixed = HostVmfsVolumeUnmapBandwidthPolicy("fixed")
	// Unmaps bandwidth is a dynamic value with lower and upper limits
	HostVmfsVolumeUnmapBandwidthPolicyDynamic = HostVmfsVolumeUnmapBandwidthPolicy("dynamic")
)

func (e HostVmfsVolumeUnmapBandwidthPolicy) Values() []HostVmfsVolumeUnmapBandwidthPolicy {
	return []HostVmfsVolumeUnmapBandwidthPolicy{
		HostVmfsVolumeUnmapBandwidthPolicyFixed,
		HostVmfsVolumeUnmapBandwidthPolicyDynamic,
	}
}

func (e HostVmfsVolumeUnmapBandwidthPolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostVmfsVolumeUnmapBandwidthPolicy"] = reflect.TypeOf((*HostVmfsVolumeUnmapBandwidthPolicy)(nil)).Elem()
}

// VMFS unmap priority.
//
// VMFS unmap reclaims unused storage space.
// This specifies the processing rate of unmaps.
type HostVmfsVolumeUnmapPriority string

const (
	// Unmap is disabled.
	HostVmfsVolumeUnmapPriorityNone = HostVmfsVolumeUnmapPriority("none")
	// Unmaps are processed at low rate.
	HostVmfsVolumeUnmapPriorityLow = HostVmfsVolumeUnmapPriority("low")
)

func (e HostVmfsVolumeUnmapPriority) Values() []HostVmfsVolumeUnmapPriority {
	return []HostVmfsVolumeUnmapPriority{
		HostVmfsVolumeUnmapPriorityNone,
		HostVmfsVolumeUnmapPriorityLow,
	}
}

func (e HostVmfsVolumeUnmapPriority) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HostVmfsVolumeUnmapPriority"] = reflect.TypeOf((*HostVmfsVolumeUnmapPriority)(nil)).Elem()
}

// List of supported algorithms for checksum calculation.
type HttpNfcLeaseManifestEntryChecksumType string

const (
	HttpNfcLeaseManifestEntryChecksumTypeSha1   = HttpNfcLeaseManifestEntryChecksumType("sha1")
	HttpNfcLeaseManifestEntryChecksumTypeSha256 = HttpNfcLeaseManifestEntryChecksumType("sha256")
)

func (e HttpNfcLeaseManifestEntryChecksumType) Values() []HttpNfcLeaseManifestEntryChecksumType {
	return []HttpNfcLeaseManifestEntryChecksumType{
		HttpNfcLeaseManifestEntryChecksumTypeSha1,
		HttpNfcLeaseManifestEntryChecksumTypeSha256,
	}
}

func (e HttpNfcLeaseManifestEntryChecksumType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HttpNfcLeaseManifestEntryChecksumType"] = reflect.TypeOf((*HttpNfcLeaseManifestEntryChecksumType)(nil)).Elem()
}

// List of supported modes by HttpNfcLease
type HttpNfcLeaseMode string

const (
	// Client pushes or downloads individual files from/to
	// each host/url provided by this lease in `HttpNfcLease.info`
	HttpNfcLeaseModePushOrGet = HttpNfcLeaseMode("pushOrGet")
	// Mode where hosts itself pull files from source URLs.
	//
	// See `HttpNfcLease.HttpNfcLeasePullFromUrls_Task`
	HttpNfcLeaseModePull = HttpNfcLeaseMode("pull")
)

func (e HttpNfcLeaseMode) Values() []HttpNfcLeaseMode {
	return []HttpNfcLeaseMode{
		HttpNfcLeaseModePushOrGet,
		HttpNfcLeaseModePull,
	}
}

func (e HttpNfcLeaseMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HttpNfcLeaseMode"] = reflect.TypeOf((*HttpNfcLeaseMode)(nil)).Elem()
}

// List of possible states of a lease.
type HttpNfcLeaseState string

const (
	// When the lease is being initialized.
	HttpNfcLeaseStateInitializing = HttpNfcLeaseState("initializing")
	// When the lease is ready and disks may be transferred.
	HttpNfcLeaseStateReady = HttpNfcLeaseState("ready")
	// When the import/export session is completed, and the lease
	// is no longer held.
	HttpNfcLeaseStateDone = HttpNfcLeaseState("done")
	// When an error has occurred.
	HttpNfcLeaseStateError = HttpNfcLeaseState("error")
)

func (e HttpNfcLeaseState) Values() []HttpNfcLeaseState {
	return []HttpNfcLeaseState{
		HttpNfcLeaseStateInitializing,
		HttpNfcLeaseStateReady,
		HttpNfcLeaseStateDone,
		HttpNfcLeaseStateError,
	}
}

func (e HttpNfcLeaseState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["HttpNfcLeaseState"] = reflect.TypeOf((*HttpNfcLeaseState)(nil)).Elem()
}

type IncompatibleHostForVmReplicationIncompatibleReason string

const (
	// Host does not support the RPO configured for VM replication.
	IncompatibleHostForVmReplicationIncompatibleReasonRpo = IncompatibleHostForVmReplicationIncompatibleReason("rpo")
	// Host does not support network compression configured for VM
	// replication.
	IncompatibleHostForVmReplicationIncompatibleReasonNetCompression = IncompatibleHostForVmReplicationIncompatibleReason("netCompression")
)

func (e IncompatibleHostForVmReplicationIncompatibleReason) Values() []IncompatibleHostForVmReplicationIncompatibleReason {
	return []IncompatibleHostForVmReplicationIncompatibleReason{
		IncompatibleHostForVmReplicationIncompatibleReasonRpo,
		IncompatibleHostForVmReplicationIncompatibleReasonNetCompression,
	}
}

func (e IncompatibleHostForVmReplicationIncompatibleReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["IncompatibleHostForVmReplicationIncompatibleReason"] = reflect.TypeOf((*IncompatibleHostForVmReplicationIncompatibleReason)(nil)).Elem()
}

// The available iSNS discovery methods.
type InternetScsiSnsDiscoveryMethod string

const (
	InternetScsiSnsDiscoveryMethodIsnsStatic = InternetScsiSnsDiscoveryMethod("isnsStatic")
	InternetScsiSnsDiscoveryMethodIsnsDhcp   = InternetScsiSnsDiscoveryMethod("isnsDhcp")
	InternetScsiSnsDiscoveryMethodIsnsSlp    = InternetScsiSnsDiscoveryMethod("isnsSlp")
)

func (e InternetScsiSnsDiscoveryMethod) Values() []InternetScsiSnsDiscoveryMethod {
	return []InternetScsiSnsDiscoveryMethod{
		InternetScsiSnsDiscoveryMethodIsnsStatic,
		InternetScsiSnsDiscoveryMethodIsnsDhcp,
		InternetScsiSnsDiscoveryMethodIsnsSlp,
	}
}

func (e InternetScsiSnsDiscoveryMethod) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["InternetScsiSnsDiscoveryMethod"] = reflect.TypeOf((*InternetScsiSnsDiscoveryMethod)(nil)).Elem()
}

type InvalidDasConfigArgumentEntryForInvalidArgument string

const (
	// Policies for admission control
	InvalidDasConfigArgumentEntryForInvalidArgumentAdmissionControl = InvalidDasConfigArgumentEntryForInvalidArgument("admissionControl")
	// User-specified heartbeat datastores
	InvalidDasConfigArgumentEntryForInvalidArgumentUserHeartbeatDs = InvalidDasConfigArgumentEntryForInvalidArgument("userHeartbeatDs")
	// VM override
	InvalidDasConfigArgumentEntryForInvalidArgumentVmConfig = InvalidDasConfigArgumentEntryForInvalidArgument("vmConfig")
)

func (e InvalidDasConfigArgumentEntryForInvalidArgument) Values() []InvalidDasConfigArgumentEntryForInvalidArgument {
	return []InvalidDasConfigArgumentEntryForInvalidArgument{
		InvalidDasConfigArgumentEntryForInvalidArgumentAdmissionControl,
		InvalidDasConfigArgumentEntryForInvalidArgumentUserHeartbeatDs,
		InvalidDasConfigArgumentEntryForInvalidArgumentVmConfig,
	}
}

func (e InvalidDasConfigArgumentEntryForInvalidArgument) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["InvalidDasConfigArgumentEntryForInvalidArgument"] = reflect.TypeOf((*InvalidDasConfigArgumentEntryForInvalidArgument)(nil)).Elem()
}

type InvalidProfileReferenceHostReason string

const (
	// The associated host and profile version are incompatible.
	InvalidProfileReferenceHostReasonIncompatibleVersion = InvalidProfileReferenceHostReason("incompatibleVersion")
	// There is no reference host associated with the profile.
	InvalidProfileReferenceHostReasonMissingReferenceHost = InvalidProfileReferenceHostReason("missingReferenceHost")
)

func (e InvalidProfileReferenceHostReason) Values() []InvalidProfileReferenceHostReason {
	return []InvalidProfileReferenceHostReason{
		InvalidProfileReferenceHostReasonIncompatibleVersion,
		InvalidProfileReferenceHostReasonMissingReferenceHost,
	}
}

func (e InvalidProfileReferenceHostReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["InvalidProfileReferenceHostReason"] = reflect.TypeOf((*InvalidProfileReferenceHostReason)(nil)).Elem()
}

// Defines the type of operation for an IO Filter.
type IoFilterOperation string

const (
	// Install an IO Filter.
	IoFilterOperationInstall = IoFilterOperation("install")
	// Uninstall an IO Filter.
	IoFilterOperationUninstall = IoFilterOperation("uninstall")
	// Upgrade an IO Filter.
	IoFilterOperationUpgrade = IoFilterOperation("upgrade")
)

func (e IoFilterOperation) Values() []IoFilterOperation {
	return []IoFilterOperation{
		IoFilterOperationInstall,
		IoFilterOperationUninstall,
		IoFilterOperationUpgrade,
	}
}

func (e IoFilterOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["IoFilterOperation"] = reflect.TypeOf((*IoFilterOperation)(nil)).Elem()
}

// Defines the type of an IO Filter.
type IoFilterType string

const (
	// Cache.
	IoFilterTypeCache = IoFilterType("cache")
	// Replication.
	IoFilterTypeReplication = IoFilterType("replication")
	// Encryption.
	IoFilterTypeEncryption = IoFilterType("encryption")
	// Compression.
	IoFilterTypeCompression = IoFilterType("compression")
	// Inspection.
	IoFilterTypeInspection = IoFilterType("inspection")
	// Datastore I/O Control.
	IoFilterTypeDatastoreIoControl = IoFilterType("datastoreIoControl")
	// Data Provider.
	IoFilterTypeDataProvider = IoFilterType("dataProvider")
	// Lightweight Data Capture.
	IoFilterTypeDataCapture = IoFilterType("dataCapture")
)

func (e IoFilterType) Values() []IoFilterType {
	return []IoFilterType{
		IoFilterTypeCache,
		IoFilterTypeReplication,
		IoFilterTypeEncryption,
		IoFilterTypeCompression,
		IoFilterTypeInspection,
		IoFilterTypeDatastoreIoControl,
		IoFilterTypeDataProvider,
		IoFilterTypeDataCapture,
	}
}

func (e IoFilterType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["IoFilterType"] = reflect.TypeOf((*IoFilterType)(nil)).Elem()
	minAPIVersionForEnumValue["IoFilterType"] = map[string]string{
		"dataCapture": "7.0.2.1",
	}
}

type IscsiPortInfoPathStatus string

const (
	// There are no paths on this Virtual NIC
	IscsiPortInfoPathStatusNotUsed = IscsiPortInfoPathStatus("notUsed")
	// All paths on this Virtual NIC are standby paths from SCSI stack
	// perspective.
	IscsiPortInfoPathStatusActive = IscsiPortInfoPathStatus("active")
	// One or more paths on the Virtual NIC are active paths to
	// storage.
	//
	// Unbinding this Virtual NIC will cause storage path
	// transitions.
	IscsiPortInfoPathStatusStandBy = IscsiPortInfoPathStatus("standBy")
	// One or more paths on the Virtual NIC is the last active
	// path to a particular storage device.
	IscsiPortInfoPathStatusLastActive = IscsiPortInfoPathStatus("lastActive")
)

func (e IscsiPortInfoPathStatus) Values() []IscsiPortInfoPathStatus {
	return []IscsiPortInfoPathStatus{
		IscsiPortInfoPathStatusNotUsed,
		IscsiPortInfoPathStatusActive,
		IscsiPortInfoPathStatusStandBy,
		IscsiPortInfoPathStatusLastActive,
	}
}

func (e IscsiPortInfoPathStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["IscsiPortInfoPathStatus"] = reflect.TypeOf((*IscsiPortInfoPathStatus)(nil)).Elem()
}

// Key provider management type.
type KmipClusterInfoKmsManagementType string

const (
	KmipClusterInfoKmsManagementTypeUnknown        = KmipClusterInfoKmsManagementType("unknown")
	KmipClusterInfoKmsManagementTypeVCenter        = KmipClusterInfoKmsManagementType("vCenter")
	KmipClusterInfoKmsManagementTypeTrustAuthority = KmipClusterInfoKmsManagementType("trustAuthority")
	// `**Since:**` vSphere API Release 7.0.2.0
	KmipClusterInfoKmsManagementTypeNativeProvider = KmipClusterInfoKmsManagementType("nativeProvider")
)

func (e KmipClusterInfoKmsManagementType) Values() []KmipClusterInfoKmsManagementType {
	return []KmipClusterInfoKmsManagementType{
		KmipClusterInfoKmsManagementTypeUnknown,
		KmipClusterInfoKmsManagementTypeVCenter,
		KmipClusterInfoKmsManagementTypeTrustAuthority,
		KmipClusterInfoKmsManagementTypeNativeProvider,
	}
}

func (e KmipClusterInfoKmsManagementType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["KmipClusterInfoKmsManagementType"] = reflect.TypeOf((*KmipClusterInfoKmsManagementType)(nil)).Elem()
	minAPIVersionForEnumValue["KmipClusterInfoKmsManagementType"] = map[string]string{
		"nativeProvider": "7.0.2.0",
	}
}

// Enumeration of the nominal latency-sensitive values which can be
// used to specify the latency-sensitivity level of the application.
//
// In terms of latency-sensitivity the values relate:
// high&gt;medium&gt;normal&gt;low.
type LatencySensitivitySensitivityLevel string

const (
	// The relative latency-sensitivity low value.
	LatencySensitivitySensitivityLevelLow = LatencySensitivitySensitivityLevel("low")
	// The relative latency-sensitivity normal value.
	//
	// This is the default latency-sensitivity value.
	LatencySensitivitySensitivityLevelNormal = LatencySensitivitySensitivityLevel("normal")
	// The relative latency-sensitivity medium value.
	LatencySensitivitySensitivityLevelMedium = LatencySensitivitySensitivityLevel("medium")
	// The relative latency-sensitivity high value.
	LatencySensitivitySensitivityLevelHigh = LatencySensitivitySensitivityLevel("high")
	// Deprecated as of vSphere API Ver 6.0. Value will be ignored and
	// treated as "normal" latency sensitivity.
	//
	// The custom absolute latency-sensitivity specified in
	// `LatencySensitivity.sensitivity` property is used to
	// define the latency-sensitivity.
	//
	// When this value is set to `LatencySensitivity.level` the
	// `LatencySensitivity.sensitivity` property should be
	// set also.
	LatencySensitivitySensitivityLevelCustom = LatencySensitivitySensitivityLevel("custom")
)

func (e LatencySensitivitySensitivityLevel) Values() []LatencySensitivitySensitivityLevel {
	return []LatencySensitivitySensitivityLevel{
		LatencySensitivitySensitivityLevelLow,
		LatencySensitivitySensitivityLevelNormal,
		LatencySensitivitySensitivityLevelMedium,
		LatencySensitivitySensitivityLevelHigh,
		LatencySensitivitySensitivityLevelCustom,
	}
}

func (e LatencySensitivitySensitivityLevel) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LatencySensitivitySensitivityLevel"] = reflect.TypeOf((*LatencySensitivitySensitivityLevel)(nil)).Elem()
}

type LicenseAssignmentFailedReason string

const (
	// The license and the entity to which it is to be assigned are not compatible.
	LicenseAssignmentFailedReasonKeyEntityMismatch = LicenseAssignmentFailedReason("keyEntityMismatch")
	// The license downgrade is disallowed because some features are in use.
	LicenseAssignmentFailedReasonDowngradeDisallowed = LicenseAssignmentFailedReason("downgradeDisallowed")
	// The inventory has hosts which are not manageable by vCenter unless in evaluation.
	LicenseAssignmentFailedReasonInventoryNotManageableByVirtualCenter = LicenseAssignmentFailedReason("inventoryNotManageableByVirtualCenter")
	// The inventory has hosts that need the license server to be configured unless vCenter is in evaluation
	LicenseAssignmentFailedReasonHostsUnmanageableByVirtualCenterWithoutLicenseServer = LicenseAssignmentFailedReason("hostsUnmanageableByVirtualCenterWithoutLicenseServer")
)

func (e LicenseAssignmentFailedReason) Values() []LicenseAssignmentFailedReason {
	return []LicenseAssignmentFailedReason{
		LicenseAssignmentFailedReasonKeyEntityMismatch,
		LicenseAssignmentFailedReasonDowngradeDisallowed,
		LicenseAssignmentFailedReasonInventoryNotManageableByVirtualCenter,
		LicenseAssignmentFailedReasonHostsUnmanageableByVirtualCenterWithoutLicenseServer,
	}
}

func (e LicenseAssignmentFailedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseAssignmentFailedReason"] = reflect.TypeOf((*LicenseAssignmentFailedReason)(nil)).Elem()
}

// Some licenses may only be allowed to load from a specified source.
//
// This enum indicates what restrictions exist for this license if any.
type LicenseFeatureInfoSourceRestriction string

const (
	// The feature does not have a source restriction.
	LicenseFeatureInfoSourceRestrictionUnrestricted = LicenseFeatureInfoSourceRestriction("unrestricted")
	// The feature's license can only be served.
	LicenseFeatureInfoSourceRestrictionServed = LicenseFeatureInfoSourceRestriction("served")
	// The feature's license can only come from a file.
	LicenseFeatureInfoSourceRestrictionFile = LicenseFeatureInfoSourceRestriction("file")
)

func (e LicenseFeatureInfoSourceRestriction) Values() []LicenseFeatureInfoSourceRestriction {
	return []LicenseFeatureInfoSourceRestriction{
		LicenseFeatureInfoSourceRestrictionUnrestricted,
		LicenseFeatureInfoSourceRestrictionServed,
		LicenseFeatureInfoSourceRestrictionFile,
	}
}

func (e LicenseFeatureInfoSourceRestriction) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseFeatureInfoSourceRestriction"] = reflect.TypeOf((*LicenseFeatureInfoSourceRestriction)(nil)).Elem()
}

// Describes the state of the feature.
type LicenseFeatureInfoState string

const (
	// The current edition license has implicitly enabled this additional feature.
	LicenseFeatureInfoStateEnabled = LicenseFeatureInfoState("enabled")
	// The current edition license does not allow this additional feature.
	LicenseFeatureInfoStateDisabled = LicenseFeatureInfoState("disabled")
	// The current edition license allows this additional feature.
	//
	// The
	// `LicenseManager.EnableFeature` and `LicenseManager.DisableFeature` methods can be used to enable or disable
	// this feature.
	LicenseFeatureInfoStateOptional = LicenseFeatureInfoState("optional")
)

func (e LicenseFeatureInfoState) Values() []LicenseFeatureInfoState {
	return []LicenseFeatureInfoState{
		LicenseFeatureInfoStateEnabled,
		LicenseFeatureInfoStateDisabled,
		LicenseFeatureInfoStateOptional,
	}
}

func (e LicenseFeatureInfoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseFeatureInfoState"] = reflect.TypeOf((*LicenseFeatureInfoState)(nil)).Elem()
}

// Cost units apply to licenses for the purpose of determining
// how many licenses are needed.
type LicenseFeatureInfoUnit string

const (
	// One license is acquired per host.
	LicenseFeatureInfoUnitHost = LicenseFeatureInfoUnit("host")
	// One license is acquired per CPU core.
	LicenseFeatureInfoUnitCpuCore = LicenseFeatureInfoUnit("cpuCore")
	// One license is acquired per CPU package.
	LicenseFeatureInfoUnitCpuPackage = LicenseFeatureInfoUnit("cpuPackage")
	// One license is acquired per server.
	LicenseFeatureInfoUnitServer = LicenseFeatureInfoUnit("server")
	// One license is acquired per virtual machine.
	LicenseFeatureInfoUnitVm = LicenseFeatureInfoUnit("vm")
)

func (e LicenseFeatureInfoUnit) Values() []LicenseFeatureInfoUnit {
	return []LicenseFeatureInfoUnit{
		LicenseFeatureInfoUnitHost,
		LicenseFeatureInfoUnitCpuCore,
		LicenseFeatureInfoUnitCpuPackage,
		LicenseFeatureInfoUnitServer,
		LicenseFeatureInfoUnitVm,
	}
}

func (e LicenseFeatureInfoUnit) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseFeatureInfoUnit"] = reflect.TypeOf((*LicenseFeatureInfoUnit)(nil)).Elem()
}

// Deprecated as of VI API 2.5, use `LicenseManager.QueryLicenseSourceAvailability`
// to obtain an array of `LicenseAvailabilityInfo` data
// objects.
//
// Licensed features have unique keys to identify them.
type LicenseManagerLicenseKey string

const (
	// The edition license for the ESX Server, Standard edition.
	//
	// This is a per
	// CPU package license.
	LicenseManagerLicenseKeyEsxFull = LicenseManagerLicenseKey("esxFull")
	// The edition license for the ESX server, VMTN edition.
	//
	// This is a per CPU package
	// license.
	LicenseManagerLicenseKeyEsxVmtn = LicenseManagerLicenseKey("esxVmtn")
	// The edition license for the ESX server, Starter edition.
	//
	// This is a per CPU
	// package license.
	LicenseManagerLicenseKeyEsxExpress = LicenseManagerLicenseKey("esxExpress")
	// Enable use of SAN.
	//
	// This is a per CPU package license.
	LicenseManagerLicenseKeySan = LicenseManagerLicenseKey("san")
	// Enable use of iSCSI.
	//
	// This is a per CPU package license.
	LicenseManagerLicenseKeyIscsi = LicenseManagerLicenseKey("iscsi")
	// Enable use of NAS.
	//
	// This is a per CPU package license.
	LicenseManagerLicenseKeyNas = LicenseManagerLicenseKey("nas")
	// Enable up to 4-way VSMP feature.
	//
	// This is a per CPU package license.
	LicenseManagerLicenseKeyVsmp = LicenseManagerLicenseKey("vsmp")
	// Enable ESX Server consolidated backup feature.
	//
	// This is a per CPU package
	// license.
	LicenseManagerLicenseKeyBackup = LicenseManagerLicenseKey("backup")
	// The edition license for a VirtualCenter server, full edition.
	//
	// This license
	// is independent of the number of CPU packages for the VirtualCenter host.
	LicenseManagerLicenseKeyVc = LicenseManagerLicenseKey("vc")
	// The edition license for a VirtualCenter server, starter edition.
	//
	// This license
	// limits the number of hosts (esxHost or serverHost) that can be managed by the
	// VirtualCenter product.
	LicenseManagerLicenseKeyVcExpress = LicenseManagerLicenseKey("vcExpress")
	// Enable VirtualCenter ESX Server host management functionality.
	//
	// This is a per
	// ESX server CPU package license.
	LicenseManagerLicenseKeyEsxHost = LicenseManagerLicenseKey("esxHost")
	// Enable VirtualCenter GSX Server host management functionality.
	//
	// This is a per
	// GSX server CPU package license.
	LicenseManagerLicenseKeyGsxHost = LicenseManagerLicenseKey("gsxHost")
	// Enable VirtualCenter VMware server host management functionality.
	//
	// This is a per
	// VMware server CPU package license.
	LicenseManagerLicenseKeyServerHost = LicenseManagerLicenseKey("serverHost")
	// Enable VirtualCenter DRS Power Management Functionality.
	//
	// This is a per CPU package
	LicenseManagerLicenseKeyDrsPower = LicenseManagerLicenseKey("drsPower")
	// Enable VMotion.
	//
	// This is a per ESX server CPU package license.
	LicenseManagerLicenseKeyVmotion = LicenseManagerLicenseKey("vmotion")
	// Enable VirtualCenter Distributed Resource Scheduler.
	//
	// This is a per ESX server
	// CPU package license.
	LicenseManagerLicenseKeyDrs = LicenseManagerLicenseKey("drs")
	// Enable VirtualCenter HA.
	//
	// This is a per ESX server CPU package license.
	LicenseManagerLicenseKeyDas = LicenseManagerLicenseKey("das")
)

func (e LicenseManagerLicenseKey) Values() []LicenseManagerLicenseKey {
	return []LicenseManagerLicenseKey{
		LicenseManagerLicenseKeyEsxFull,
		LicenseManagerLicenseKeyEsxVmtn,
		LicenseManagerLicenseKeyEsxExpress,
		LicenseManagerLicenseKeySan,
		LicenseManagerLicenseKeyIscsi,
		LicenseManagerLicenseKeyNas,
		LicenseManagerLicenseKeyVsmp,
		LicenseManagerLicenseKeyBackup,
		LicenseManagerLicenseKeyVc,
		LicenseManagerLicenseKeyVcExpress,
		LicenseManagerLicenseKeyEsxHost,
		LicenseManagerLicenseKeyGsxHost,
		LicenseManagerLicenseKeyServerHost,
		LicenseManagerLicenseKeyDrsPower,
		LicenseManagerLicenseKeyVmotion,
		LicenseManagerLicenseKeyDrs,
		LicenseManagerLicenseKeyDas,
	}
}

func (e LicenseManagerLicenseKey) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseManagerLicenseKey"] = reflect.TypeOf((*LicenseManagerLicenseKey)(nil)).Elem()
}

// Deprecated as of vSphere API 4.0, this is not used by the system.
//
// State of licensing subsystem.
type LicenseManagerState string

const (
	// Setting or resetting configuration in progress.
	LicenseManagerStateInitializing = LicenseManagerState("initializing")
	// Running within operating parameters.
	LicenseManagerStateNormal = LicenseManagerState("normal")
	// License source unavailable, using license cache.
	LicenseManagerStateMarginal = LicenseManagerState("marginal")
	// Initialization has failed or grace period expired.
	LicenseManagerStateFault = LicenseManagerState("fault")
)

func (e LicenseManagerState) Values() []LicenseManagerState {
	return []LicenseManagerState{
		LicenseManagerStateInitializing,
		LicenseManagerStateNormal,
		LicenseManagerStateMarginal,
		LicenseManagerStateFault,
	}
}

func (e LicenseManagerState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseManagerState"] = reflect.TypeOf((*LicenseManagerState)(nil)).Elem()
}

// Describes the reservation state of a license.
type LicenseReservationInfoState string

const (
	// This license is currently unused by the system, or the feature does not
	// apply.
	//
	// For example, a DRS license appears as NotUsed if the host is not
	// part of a DRS-enabled cluster.
	LicenseReservationInfoStateNotUsed = LicenseReservationInfoState("notUsed")
	// This indicates that the license has expired or the system attempted to acquire
	// the license but was not successful in reserving it.
	LicenseReservationInfoStateNoLicense = LicenseReservationInfoState("noLicense")
	// The LicenseManager failed to acquire a license but the implementation
	// policy allows us to use the licensed feature anyway.
	//
	// This is possible, for
	// example, when a license server becomes unavailable after a license had been
	// successfully reserved from it.
	LicenseReservationInfoStateUnlicensedUse = LicenseReservationInfoState("unlicensedUse")
	// The required number of licenses have been acquired from the license source.
	LicenseReservationInfoStateLicensed = LicenseReservationInfoState("licensed")
)

func (e LicenseReservationInfoState) Values() []LicenseReservationInfoState {
	return []LicenseReservationInfoState{
		LicenseReservationInfoStateNotUsed,
		LicenseReservationInfoStateNoLicense,
		LicenseReservationInfoStateUnlicensedUse,
		LicenseReservationInfoStateLicensed,
	}
}

func (e LicenseReservationInfoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LicenseReservationInfoState"] = reflect.TypeOf((*LicenseReservationInfoState)(nil)).Elem()
}

// The Discovery Protocol operation.
type LinkDiscoveryProtocolConfigOperationType string

const (
	// Don't listen for incoming discovery packets and don't sent discover
	// packets for the switch either.
	LinkDiscoveryProtocolConfigOperationTypeNone = LinkDiscoveryProtocolConfigOperationType("none")
	// Listen for incoming discovery packets but don't sent discovery packet
	// for the switch.
	LinkDiscoveryProtocolConfigOperationTypeListen = LinkDiscoveryProtocolConfigOperationType("listen")
	// Sent discovery packets for the switch, but don't listen for incoming
	// discovery packets.
	LinkDiscoveryProtocolConfigOperationTypeAdvertise = LinkDiscoveryProtocolConfigOperationType("advertise")
	// Sent discovery packets for the switch and listen for incoming
	// discovery packets.
	LinkDiscoveryProtocolConfigOperationTypeBoth = LinkDiscoveryProtocolConfigOperationType("both")
)

func (e LinkDiscoveryProtocolConfigOperationType) Values() []LinkDiscoveryProtocolConfigOperationType {
	return []LinkDiscoveryProtocolConfigOperationType{
		LinkDiscoveryProtocolConfigOperationTypeNone,
		LinkDiscoveryProtocolConfigOperationTypeListen,
		LinkDiscoveryProtocolConfigOperationTypeAdvertise,
		LinkDiscoveryProtocolConfigOperationTypeBoth,
	}
}

func (e LinkDiscoveryProtocolConfigOperationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LinkDiscoveryProtocolConfigOperationType"] = reflect.TypeOf((*LinkDiscoveryProtocolConfigOperationType)(nil)).Elem()
}

// The Discovery Protocol types.
type LinkDiscoveryProtocolConfigProtocolType string

const (
	// Cisco Discovery Protocol
	LinkDiscoveryProtocolConfigProtocolTypeCdp = LinkDiscoveryProtocolConfigProtocolType("cdp")
	// Link Layer Discovery Protocol
	LinkDiscoveryProtocolConfigProtocolTypeLldp = LinkDiscoveryProtocolConfigProtocolType("lldp")
)

func (e LinkDiscoveryProtocolConfigProtocolType) Values() []LinkDiscoveryProtocolConfigProtocolType {
	return []LinkDiscoveryProtocolConfigProtocolType{
		LinkDiscoveryProtocolConfigProtocolTypeCdp,
		LinkDiscoveryProtocolConfigProtocolTypeLldp,
	}
}

func (e LinkDiscoveryProtocolConfigProtocolType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["LinkDiscoveryProtocolConfigProtocolType"] = reflect.TypeOf((*LinkDiscoveryProtocolConfigProtocolType)(nil)).Elem()
}

// The Status enumeration defines a general "health" value for a managed entity.
type ManagedEntityStatus string

const (
	// The status is unknown.
	ManagedEntityStatusGray = ManagedEntityStatus("gray")
	// The entity is OK.
	ManagedEntityStatusGreen = ManagedEntityStatus("green")
	// The entity might have a problem.
	ManagedEntityStatusYellow = ManagedEntityStatus("yellow")
	// The entity definitely has a problem.
	ManagedEntityStatusRed = ManagedEntityStatus("red")
)

func (e ManagedEntityStatus) Values() []ManagedEntityStatus {
	return []ManagedEntityStatus{
		ManagedEntityStatusGray,
		ManagedEntityStatusGreen,
		ManagedEntityStatusYellow,
		ManagedEntityStatusRed,
	}
}

func (e ManagedEntityStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ManagedEntityStatus"] = reflect.TypeOf((*ManagedEntityStatus)(nil)).Elem()
}

// The operation on the target metric item.
type MetricAlarmOperator string

const (
	// Test if the target metric item is above the given red or yellow values.
	MetricAlarmOperatorIsAbove = MetricAlarmOperator("isAbove")
	// Test if the target metric item is below the given red or yellow values.
	MetricAlarmOperatorIsBelow = MetricAlarmOperator("isBelow")
)

func (e MetricAlarmOperator) Values() []MetricAlarmOperator {
	return []MetricAlarmOperator{
		MetricAlarmOperatorIsAbove,
		MetricAlarmOperatorIsBelow,
	}
}

func (e MetricAlarmOperator) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["MetricAlarmOperator"] = reflect.TypeOf((*MetricAlarmOperator)(nil)).Elem()
}

// Set of constants defining the possible states of a multipath path.
type MultipathState string

const (
	MultipathStateStandby  = MultipathState("standby")
	MultipathStateActive   = MultipathState("active")
	MultipathStateDisabled = MultipathState("disabled")
	MultipathStateDead     = MultipathState("dead")
	MultipathStateUnknown  = MultipathState("unknown")
)

func (e MultipathState) Values() []MultipathState {
	return []MultipathState{
		MultipathStateStandby,
		MultipathStateActive,
		MultipathStateDisabled,
		MultipathStateDead,
		MultipathStateUnknown,
	}
}

func (e MultipathState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["MultipathState"] = reflect.TypeOf((*MultipathState)(nil)).Elem()
}

// NetBIOS configuration mode.
type NetBIOSConfigInfoMode string

const (
	// Mode of NetBIOS is unknown.
	NetBIOSConfigInfoModeUnknown = NetBIOSConfigInfoMode("unknown")
	// NetBIOS is enabled.
	NetBIOSConfigInfoModeEnabled = NetBIOSConfigInfoMode("enabled")
	// NetBIOS is disabled.
	NetBIOSConfigInfoModeDisabled = NetBIOSConfigInfoMode("disabled")
	// DHCP server decides whether or not to use NetBIOS.
	NetBIOSConfigInfoModeEnabledViaDHCP = NetBIOSConfigInfoMode("enabledViaDHCP")
)

func (e NetBIOSConfigInfoMode) Values() []NetBIOSConfigInfoMode {
	return []NetBIOSConfigInfoMode{
		NetBIOSConfigInfoModeUnknown,
		NetBIOSConfigInfoModeEnabled,
		NetBIOSConfigInfoModeDisabled,
		NetBIOSConfigInfoModeEnabledViaDHCP,
	}
}

func (e NetBIOSConfigInfoMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NetBIOSConfigInfoMode"] = reflect.TypeOf((*NetBIOSConfigInfoMode)(nil)).Elem()
}

// This specifies how an IP address was obtained for a given interface.
//
// See RFC 4293 IpAddressOriginTC.
type NetIpConfigInfoIpAddressOrigin string

const (
	// Any other type of address configuration other than the below
	// mentioned ones will fall under this category.
	//
	// For e.g., automatic
	// address configuration for the link local address falls under
	// this type.
	NetIpConfigInfoIpAddressOriginOther = NetIpConfigInfoIpAddressOrigin("other")
	// The address is configured manually.
	//
	// The term 'static' is a synonym.
	NetIpConfigInfoIpAddressOriginManual = NetIpConfigInfoIpAddressOrigin("manual")
	// The address is configured through dhcp.
	NetIpConfigInfoIpAddressOriginDhcp = NetIpConfigInfoIpAddressOrigin("dhcp")
	// The address is obtained through stateless autoconfiguration (autoconf).
	//
	// See RFC 4862, IPv6 Stateless Address Autoconfiguration.
	NetIpConfigInfoIpAddressOriginLinklayer = NetIpConfigInfoIpAddressOrigin("linklayer")
	// The address is chosen by the system at random
	// e.g., an IPv4 address within 169.254/16, or an RFC 3041 privacy address.
	NetIpConfigInfoIpAddressOriginRandom = NetIpConfigInfoIpAddressOrigin("random")
)

func (e NetIpConfigInfoIpAddressOrigin) Values() []NetIpConfigInfoIpAddressOrigin {
	return []NetIpConfigInfoIpAddressOrigin{
		NetIpConfigInfoIpAddressOriginOther,
		NetIpConfigInfoIpAddressOriginManual,
		NetIpConfigInfoIpAddressOriginDhcp,
		NetIpConfigInfoIpAddressOriginLinklayer,
		NetIpConfigInfoIpAddressOriginRandom,
	}
}

func (e NetIpConfigInfoIpAddressOrigin) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NetIpConfigInfoIpAddressOrigin"] = reflect.TypeOf((*NetIpConfigInfoIpAddressOrigin)(nil)).Elem()
}

type NetIpConfigInfoIpAddressStatus string

const (
	// Indicates that this is a valid address.
	NetIpConfigInfoIpAddressStatusPreferred = NetIpConfigInfoIpAddressStatus("preferred")
	// Indicates that this is a valid but deprecated address
	// that should no longer be used as a source address.
	NetIpConfigInfoIpAddressStatusDeprecated = NetIpConfigInfoIpAddressStatus("deprecated")
	// Indicates that this isn't a valid.
	NetIpConfigInfoIpAddressStatusInvalid = NetIpConfigInfoIpAddressStatus("invalid")
	// Indicates that the address is not accessible because
	// interface is not operational.
	NetIpConfigInfoIpAddressStatusInaccessible = NetIpConfigInfoIpAddressStatus("inaccessible")
	// Indicates that the status cannot be determined.
	NetIpConfigInfoIpAddressStatusUnknown = NetIpConfigInfoIpAddressStatus("unknown")
	// Indicates that the uniqueness of the
	// address on the link is presently being verified.
	NetIpConfigInfoIpAddressStatusTentative = NetIpConfigInfoIpAddressStatus("tentative")
	// Indicates the address has been determined to be non-unique
	// on the link, this address will not be reachable.
	NetIpConfigInfoIpAddressStatusDuplicate = NetIpConfigInfoIpAddressStatus("duplicate")
)

func (e NetIpConfigInfoIpAddressStatus) Values() []NetIpConfigInfoIpAddressStatus {
	return []NetIpConfigInfoIpAddressStatus{
		NetIpConfigInfoIpAddressStatusPreferred,
		NetIpConfigInfoIpAddressStatusDeprecated,
		NetIpConfigInfoIpAddressStatusInvalid,
		NetIpConfigInfoIpAddressStatusInaccessible,
		NetIpConfigInfoIpAddressStatusUnknown,
		NetIpConfigInfoIpAddressStatusTentative,
		NetIpConfigInfoIpAddressStatusDuplicate,
	}
}

func (e NetIpConfigInfoIpAddressStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NetIpConfigInfoIpAddressStatus"] = reflect.TypeOf((*NetIpConfigInfoIpAddressStatus)(nil)).Elem()
}

// IP Stack keeps state on entries in IpNetToMedia table to perform
// physical address lookups for IP addresses.
//
// Here are the standard
// states per @see RFC 4293 ipNetToMediaType.
type NetIpStackInfoEntryType string

const (
	// This implementation is reporting something other than
	// what states are listed below.
	NetIpStackInfoEntryTypeOther = NetIpStackInfoEntryType("other")
	// The IP Stack has marked this entry as not useable.
	NetIpStackInfoEntryTypeInvalid = NetIpStackInfoEntryType("invalid")
	// This entry has been learned using ARP or NDP.
	NetIpStackInfoEntryTypeDynamic = NetIpStackInfoEntryType("dynamic")
	// This entry was set manually.
	NetIpStackInfoEntryTypeManual = NetIpStackInfoEntryType("manual")
)

func (e NetIpStackInfoEntryType) Values() []NetIpStackInfoEntryType {
	return []NetIpStackInfoEntryType{
		NetIpStackInfoEntryTypeOther,
		NetIpStackInfoEntryTypeInvalid,
		NetIpStackInfoEntryTypeDynamic,
		NetIpStackInfoEntryTypeManual,
	}
}

func (e NetIpStackInfoEntryType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NetIpStackInfoEntryType"] = reflect.TypeOf((*NetIpStackInfoEntryType)(nil)).Elem()
}

// The set of values used to determine ordering of default routers.
//
// See RFC 4293 ipDefaultRouterPreference.
type NetIpStackInfoPreference string

const (
	NetIpStackInfoPreferenceReserved = NetIpStackInfoPreference("reserved")
	NetIpStackInfoPreferenceLow      = NetIpStackInfoPreference("low")
	NetIpStackInfoPreferenceMedium   = NetIpStackInfoPreference("medium")
	NetIpStackInfoPreferenceHigh     = NetIpStackInfoPreference("high")
)

func (e NetIpStackInfoPreference) Values() []NetIpStackInfoPreference {
	return []NetIpStackInfoPreference{
		NetIpStackInfoPreferenceReserved,
		NetIpStackInfoPreferenceLow,
		NetIpStackInfoPreferenceMedium,
		NetIpStackInfoPreferenceHigh,
	}
}

func (e NetIpStackInfoPreference) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NetIpStackInfoPreference"] = reflect.TypeOf((*NetIpStackInfoPreference)(nil)).Elem()
}

type NotSupportedDeviceForFTDeviceType string

const (
	// vmxnet3 virtual Ethernet adapter
	NotSupportedDeviceForFTDeviceTypeVirtualVmxnet3 = NotSupportedDeviceForFTDeviceType("virtualVmxnet3")
	// paravirtualized SCSI controller
	NotSupportedDeviceForFTDeviceTypeParaVirtualSCSIController = NotSupportedDeviceForFTDeviceType("paraVirtualSCSIController")
)

func (e NotSupportedDeviceForFTDeviceType) Values() []NotSupportedDeviceForFTDeviceType {
	return []NotSupportedDeviceForFTDeviceType{
		NotSupportedDeviceForFTDeviceTypeVirtualVmxnet3,
		NotSupportedDeviceForFTDeviceTypeParaVirtualSCSIController,
	}
}

func (e NotSupportedDeviceForFTDeviceType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NotSupportedDeviceForFTDeviceType"] = reflect.TypeOf((*NotSupportedDeviceForFTDeviceType)(nil)).Elem()
}

// Reasons why the number of virtual CPUs is incompatible.
type NumVirtualCpusIncompatibleReason string

const (
	// Deprecated as of vSphere API 6.0.
	//
	// The virtual machine needs to support record/replay functionality.
	NumVirtualCpusIncompatibleReasonRecordReplay = NumVirtualCpusIncompatibleReason("recordReplay")
	// The virtual machine is enabled for fault tolerance.
	NumVirtualCpusIncompatibleReasonFaultTolerance = NumVirtualCpusIncompatibleReason("faultTolerance")
)

func (e NumVirtualCpusIncompatibleReason) Values() []NumVirtualCpusIncompatibleReason {
	return []NumVirtualCpusIncompatibleReason{
		NumVirtualCpusIncompatibleReasonRecordReplay,
		NumVirtualCpusIncompatibleReasonFaultTolerance,
	}
}

func (e NumVirtualCpusIncompatibleReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NumVirtualCpusIncompatibleReason"] = reflect.TypeOf((*NumVirtualCpusIncompatibleReason)(nil)).Elem()
}

// State of interleave set
type NvdimmInterleaveSetState string

const (
	// Interleave set is invalid
	NvdimmInterleaveSetStateInvalid = NvdimmInterleaveSetState("invalid")
	// Interleave set is valid and active
	NvdimmInterleaveSetStateActive = NvdimmInterleaveSetState("active")
)

func (e NvdimmInterleaveSetState) Values() []NvdimmInterleaveSetState {
	return []NvdimmInterleaveSetState{
		NvdimmInterleaveSetStateInvalid,
		NvdimmInterleaveSetStateActive,
	}
}

func (e NvdimmInterleaveSetState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmInterleaveSetState"] = reflect.TypeOf((*NvdimmInterleaveSetState)(nil)).Elem()
}

// Overall health state for a namespace
type NvdimmNamespaceDetailsHealthStatus string

const (
	// Namespace health is normal
	NvdimmNamespaceDetailsHealthStatusNormal = NvdimmNamespaceDetailsHealthStatus("normal")
	// Namespace health is missing
	NvdimmNamespaceDetailsHealthStatusMissing = NvdimmNamespaceDetailsHealthStatus("missing")
	// Namespace health label is missing
	NvdimmNamespaceDetailsHealthStatusLabelMissing = NvdimmNamespaceDetailsHealthStatus("labelMissing")
	// Namespace health interleave broken
	NvdimmNamespaceDetailsHealthStatusInterleaveBroken = NvdimmNamespaceDetailsHealthStatus("interleaveBroken")
	// Namespace health label is inconsistent
	NvdimmNamespaceDetailsHealthStatusLabelInconsistent = NvdimmNamespaceDetailsHealthStatus("labelInconsistent")
)

func (e NvdimmNamespaceDetailsHealthStatus) Values() []NvdimmNamespaceDetailsHealthStatus {
	return []NvdimmNamespaceDetailsHealthStatus{
		NvdimmNamespaceDetailsHealthStatusNormal,
		NvdimmNamespaceDetailsHealthStatusMissing,
		NvdimmNamespaceDetailsHealthStatusLabelMissing,
		NvdimmNamespaceDetailsHealthStatusInterleaveBroken,
		NvdimmNamespaceDetailsHealthStatusLabelInconsistent,
	}
}

func (e NvdimmNamespaceDetailsHealthStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmNamespaceDetailsHealthStatus"] = reflect.TypeOf((*NvdimmNamespaceDetailsHealthStatus)(nil)).Elem()
}

// State of Namespace
type NvdimmNamespaceDetailsState string

const (
	// Namespace is invalid
	NvdimmNamespaceDetailsStateInvalid = NvdimmNamespaceDetailsState("invalid")
	// Namespace is valid but not in use
	NvdimmNamespaceDetailsStateNotInUse = NvdimmNamespaceDetailsState("notInUse")
	// Namespace is valid and is in use
	NvdimmNamespaceDetailsStateInUse = NvdimmNamespaceDetailsState("inUse")
)

func (e NvdimmNamespaceDetailsState) Values() []NvdimmNamespaceDetailsState {
	return []NvdimmNamespaceDetailsState{
		NvdimmNamespaceDetailsStateInvalid,
		NvdimmNamespaceDetailsStateNotInUse,
		NvdimmNamespaceDetailsStateInUse,
	}
}

func (e NvdimmNamespaceDetailsState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmNamespaceDetailsState"] = reflect.TypeOf((*NvdimmNamespaceDetailsState)(nil)).Elem()
}

// Overall health state for a namespace
type NvdimmNamespaceHealthStatus string

const (
	// Namespace health is normal
	NvdimmNamespaceHealthStatusNormal = NvdimmNamespaceHealthStatus("normal")
	// Namespace health is missing
	NvdimmNamespaceHealthStatusMissing = NvdimmNamespaceHealthStatus("missing")
	// Namespace health label is missing
	NvdimmNamespaceHealthStatusLabelMissing = NvdimmNamespaceHealthStatus("labelMissing")
	// Namespace health interleave broken
	NvdimmNamespaceHealthStatusInterleaveBroken = NvdimmNamespaceHealthStatus("interleaveBroken")
	// Namespace health label is inconsistent
	NvdimmNamespaceHealthStatusLabelInconsistent = NvdimmNamespaceHealthStatus("labelInconsistent")
	// Namespace health BTT is corrupt
	NvdimmNamespaceHealthStatusBttCorrupt = NvdimmNamespaceHealthStatus("bttCorrupt")
	// Namespace health encountered bad block
	NvdimmNamespaceHealthStatusBadBlockSize = NvdimmNamespaceHealthStatus("badBlockSize")
)

func (e NvdimmNamespaceHealthStatus) Values() []NvdimmNamespaceHealthStatus {
	return []NvdimmNamespaceHealthStatus{
		NvdimmNamespaceHealthStatusNormal,
		NvdimmNamespaceHealthStatusMissing,
		NvdimmNamespaceHealthStatusLabelMissing,
		NvdimmNamespaceHealthStatusInterleaveBroken,
		NvdimmNamespaceHealthStatusLabelInconsistent,
		NvdimmNamespaceHealthStatusBttCorrupt,
		NvdimmNamespaceHealthStatusBadBlockSize,
	}
}

func (e NvdimmNamespaceHealthStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmNamespaceHealthStatus"] = reflect.TypeOf((*NvdimmNamespaceHealthStatus)(nil)).Elem()
}

// State of Namespace
type NvdimmNamespaceState string

const (
	// Namespace is invalid
	NvdimmNamespaceStateInvalid = NvdimmNamespaceState("invalid")
	// Namespace is valid but not in use
	NvdimmNamespaceStateNotInUse = NvdimmNamespaceState("notInUse")
	// Namespace is valid and is in use
	NvdimmNamespaceStateInUse = NvdimmNamespaceState("inUse")
)

func (e NvdimmNamespaceState) Values() []NvdimmNamespaceState {
	return []NvdimmNamespaceState{
		NvdimmNamespaceStateInvalid,
		NvdimmNamespaceStateNotInUse,
		NvdimmNamespaceStateInUse,
	}
}

func (e NvdimmNamespaceState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmNamespaceState"] = reflect.TypeOf((*NvdimmNamespaceState)(nil)).Elem()
}

// Type of namespace.
type NvdimmNamespaceType string

const (
	// Block mode namespace
	NvdimmNamespaceTypeBlockNamespace = NvdimmNamespaceType("blockNamespace")
	// Persistent mode namespace
	NvdimmNamespaceTypePersistentNamespace = NvdimmNamespaceType("persistentNamespace")
)

func (e NvdimmNamespaceType) Values() []NvdimmNamespaceType {
	return []NvdimmNamespaceType{
		NvdimmNamespaceTypeBlockNamespace,
		NvdimmNamespaceTypePersistentNamespace,
	}
}

func (e NvdimmNamespaceType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmNamespaceType"] = reflect.TypeOf((*NvdimmNamespaceType)(nil)).Elem()
}

// Overall state of NVDIMM
type NvdimmNvdimmHealthInfoState string

const (
	// NVDIMM state is normal
	NvdimmNvdimmHealthInfoStateNormal = NvdimmNvdimmHealthInfoState("normal")
	// Error in NVDIMM state.
	//
	// Potential data loss.
	NvdimmNvdimmHealthInfoStateError = NvdimmNvdimmHealthInfoState("error")
)

func (e NvdimmNvdimmHealthInfoState) Values() []NvdimmNvdimmHealthInfoState {
	return []NvdimmNvdimmHealthInfoState{
		NvdimmNvdimmHealthInfoStateNormal,
		NvdimmNvdimmHealthInfoStateError,
	}
}

func (e NvdimmNvdimmHealthInfoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmNvdimmHealthInfoState"] = reflect.TypeOf((*NvdimmNvdimmHealthInfoState)(nil)).Elem()
}

// An indicator of how a memory range is being used
type NvdimmRangeType string

const (
	// Identifies the region to be volatile
	NvdimmRangeTypeVolatileRange = NvdimmRangeType("volatileRange")
	// Identifies the region to be persistent
	NvdimmRangeTypePersistentRange = NvdimmRangeType("persistentRange")
	// NVDIMM control region
	NvdimmRangeTypeControlRange = NvdimmRangeType("controlRange")
	// NVDIMM block data window region
	NvdimmRangeTypeBlockRange = NvdimmRangeType("blockRange")
	// NVDIMM volatile virtual disk region
	NvdimmRangeTypeVolatileVirtualDiskRange = NvdimmRangeType("volatileVirtualDiskRange")
	// NVDIMM volatile virtual CD region
	NvdimmRangeTypeVolatileVirtualCDRange = NvdimmRangeType("volatileVirtualCDRange")
	// NVDIMM persistent virtual disk region
	NvdimmRangeTypePersistentVirtualDiskRange = NvdimmRangeType("persistentVirtualDiskRange")
	// NVDIMM persistent virtual CD region
	NvdimmRangeTypePersistentVirtualCDRange = NvdimmRangeType("persistentVirtualCDRange")
)

func (e NvdimmRangeType) Values() []NvdimmRangeType {
	return []NvdimmRangeType{
		NvdimmRangeTypeVolatileRange,
		NvdimmRangeTypePersistentRange,
		NvdimmRangeTypeControlRange,
		NvdimmRangeTypeBlockRange,
		NvdimmRangeTypeVolatileVirtualDiskRange,
		NvdimmRangeTypeVolatileVirtualCDRange,
		NvdimmRangeTypePersistentVirtualDiskRange,
		NvdimmRangeTypePersistentVirtualCDRange,
	}
}

func (e NvdimmRangeType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["NvdimmRangeType"] = reflect.TypeOf((*NvdimmRangeType)(nil)).Elem()
}

// Enumeration of different kinds of updates.
type ObjectUpdateKind string

const (
	// A property of the managed object changed its value.
	ObjectUpdateKindModify = ObjectUpdateKind("modify")
	// A managed object became visible to a filter for the first time.
	//
	// For instance, this can happen if a virtual machine is added to a
	// folder.
	ObjectUpdateKindEnter = ObjectUpdateKind("enter")
	// A managed object left the set of objects visible to a filter.
	//
	// For
	// instance, this can happen when a virtual machine is destroyed.
	ObjectUpdateKindLeave = ObjectUpdateKind("leave")
)

func (e ObjectUpdateKind) Values() []ObjectUpdateKind {
	return []ObjectUpdateKind{
		ObjectUpdateKindModify,
		ObjectUpdateKindEnter,
		ObjectUpdateKindLeave,
	}
}

func (e ObjectUpdateKind) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ObjectUpdateKind"] = reflect.TypeOf((*ObjectUpdateKind)(nil)).Elem()
}

// The type of an OST node.
//
// Each OST node corresponds to an element in the OVF descriptor. See `OvfConsumerOstNode`
// for a description of the different node types.
type OvfConsumerOstNodeType string

const (
	OvfConsumerOstNodeTypeEnvelope                = OvfConsumerOstNodeType("envelope")
	OvfConsumerOstNodeTypeVirtualSystem           = OvfConsumerOstNodeType("virtualSystem")
	OvfConsumerOstNodeTypeVirtualSystemCollection = OvfConsumerOstNodeType("virtualSystemCollection")
)

func (e OvfConsumerOstNodeType) Values() []OvfConsumerOstNodeType {
	return []OvfConsumerOstNodeType{
		OvfConsumerOstNodeTypeEnvelope,
		OvfConsumerOstNodeTypeVirtualSystem,
		OvfConsumerOstNodeTypeVirtualSystemCollection,
	}
}

func (e OvfConsumerOstNodeType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["OvfConsumerOstNodeType"] = reflect.TypeOf((*OvfConsumerOstNodeType)(nil)).Elem()
}

// Types of disk provisioning that can be set for the disk in the deployed OVF
// package.
type OvfCreateImportSpecParamsDiskProvisioningType string

const (
	// A sparse (allocate on demand) monolithic disk.
	//
	// Disks in this format can
	// be used with other VMware products.
	OvfCreateImportSpecParamsDiskProvisioningTypeMonolithicSparse = OvfCreateImportSpecParamsDiskProvisioningType("monolithicSparse")
	// A preallocated monolithic disk.
	//
	// Disks in this format can be used with
	// other VMware products.
	OvfCreateImportSpecParamsDiskProvisioningTypeMonolithicFlat = OvfCreateImportSpecParamsDiskProvisioningType("monolithicFlat")
	// A sparse (allocate on demand) disk with 2GB maximum extent size.
	//
	// Disks in this format can be used with other VMware products. The 2GB
	// extent size makes these disks easier to burn to dvd or use on
	// filesystems that don't support large files.
	OvfCreateImportSpecParamsDiskProvisioningTypeTwoGbMaxExtentSparse = OvfCreateImportSpecParamsDiskProvisioningType("twoGbMaxExtentSparse")
	// A preallocated disk with 2GB maximum extent size.
	//
	// Disks in this format
	// can be used with other VMware products. The 2GB extent size
	// makes these disks easier to burn to dvd or use on filesystems that
	// don't support large files.
	OvfCreateImportSpecParamsDiskProvisioningTypeTwoGbMaxExtentFlat = OvfCreateImportSpecParamsDiskProvisioningType("twoGbMaxExtentFlat")
	// Space required for thin-provisioned virtual disk is allocated and
	// zeroed on demand as the space is used.
	OvfCreateImportSpecParamsDiskProvisioningTypeThin = OvfCreateImportSpecParamsDiskProvisioningType("thin")
	// A thick disk has all space allocated at creation time
	// and the space is zeroed on demand as the space is used.
	OvfCreateImportSpecParamsDiskProvisioningTypeThick = OvfCreateImportSpecParamsDiskProvisioningType("thick")
	// A sparse (allocate on demand) format with additional space
	// optimizations.
	OvfCreateImportSpecParamsDiskProvisioningTypeSeSparse = OvfCreateImportSpecParamsDiskProvisioningType("seSparse")
	// An eager zeroed thick disk has all space allocated and wiped clean
	// of any previous contents on the physical media at creation time.
	//
	// Such disks may take longer time during creation compared to other
	// disk formats.
	OvfCreateImportSpecParamsDiskProvisioningTypeEagerZeroedThick = OvfCreateImportSpecParamsDiskProvisioningType("eagerZeroedThick")
	// Depending on the host type, Sparse is mapped to either
	// MonolithicSparse or Thin.
	OvfCreateImportSpecParamsDiskProvisioningTypeSparse = OvfCreateImportSpecParamsDiskProvisioningType("sparse")
	// Depending on the host type, Flat is mapped to either
	// MonolithicFlat or Thick.
	OvfCreateImportSpecParamsDiskProvisioningTypeFlat = OvfCreateImportSpecParamsDiskProvisioningType("flat")
)

func (e OvfCreateImportSpecParamsDiskProvisioningType) Values() []OvfCreateImportSpecParamsDiskProvisioningType {
	return []OvfCreateImportSpecParamsDiskProvisioningType{
		OvfCreateImportSpecParamsDiskProvisioningTypeMonolithicSparse,
		OvfCreateImportSpecParamsDiskProvisioningTypeMonolithicFlat,
		OvfCreateImportSpecParamsDiskProvisioningTypeTwoGbMaxExtentSparse,
		OvfCreateImportSpecParamsDiskProvisioningTypeTwoGbMaxExtentFlat,
		OvfCreateImportSpecParamsDiskProvisioningTypeThin,
		OvfCreateImportSpecParamsDiskProvisioningTypeThick,
		OvfCreateImportSpecParamsDiskProvisioningTypeSeSparse,
		OvfCreateImportSpecParamsDiskProvisioningTypeEagerZeroedThick,
		OvfCreateImportSpecParamsDiskProvisioningTypeSparse,
		OvfCreateImportSpecParamsDiskProvisioningTypeFlat,
	}
}

func (e OvfCreateImportSpecParamsDiskProvisioningType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["OvfCreateImportSpecParamsDiskProvisioningType"] = reflect.TypeOf((*OvfCreateImportSpecParamsDiskProvisioningType)(nil)).Elem()
}

// The format in which performance counter data is returned.
type PerfFormat string

const (
	// Counters returned in an array of data objects.
	PerfFormatNormal = PerfFormat("normal")
	// Counters returned in comma-separate value (CSV) format.
	PerfFormatCsv = PerfFormat("csv")
)

func (e PerfFormat) Values() []PerfFormat {
	return []PerfFormat{
		PerfFormatNormal,
		PerfFormatCsv,
	}
}

func (e PerfFormat) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PerfFormat"] = reflect.TypeOf((*PerfFormat)(nil)).Elem()
}

// Indicates the type of statistical measurement that a counter&#146;s
// value represents.
//
// Valid types are &#147;absolute&#148;,
// &#147;delta&#148;, or &#147;rate&#148;.
type PerfStatsType string

const (
	// Represents an actual value, level, or state of the counter.
	//
	// For
	// example, the &#147;uptime&#148; counter (`*system*` group)
	// represents the actual number of seconds since startup. The
	// &#147;capacity&#148; counter represents the actual configured size
	// of the specified datastore. In other words, number of samples,
	// samplingPeriod, and intervals have no bearing on an
	// &#147;absolute&#148; counter&#147;s value.
	PerfStatsTypeAbsolute = PerfStatsType("absolute")
	// Represents an amount of change for the counter during the `PerfInterval.samplingPeriod` as compared to the previous
	// `interval`.
	//
	// The first sampling interval
	PerfStatsTypeDelta = PerfStatsType("delta")
	// Represents a value that has been normalized over the `PerfInterval.samplingPeriod`, enabling values for the same
	// counter type to be compared, regardless of interval.
	//
	// For example,
	// the number of reads per second.
	PerfStatsTypeRate = PerfStatsType("rate")
)

func (e PerfStatsType) Values() []PerfStatsType {
	return []PerfStatsType{
		PerfStatsTypeAbsolute,
		PerfStatsTypeDelta,
		PerfStatsTypeRate,
	}
}

func (e PerfStatsType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PerfStatsType"] = reflect.TypeOf((*PerfStatsType)(nil)).Elem()
}

// Indicates how multiple samples of a specific counter type are
// transformed into a single statistical value.
type PerfSummaryType string

const (
	// The actual value collected or the average of all values collected
	// during the summary period.
	PerfSummaryTypeAverage = PerfSummaryType("average")
	// The maximum value of the performance counter value over the
	// summarization period.
	PerfSummaryTypeMaximum = PerfSummaryType("maximum")
	// The minimum value of the performance counter value over the
	// summarization period.
	PerfSummaryTypeMinimum = PerfSummaryType("minimum")
	// The most recent value of the performance counter over the
	// summarization period.
	PerfSummaryTypeLatest = PerfSummaryType("latest")
	// The sum of all the values of the performance counter over the
	// summarization period.
	PerfSummaryTypeSummation = PerfSummaryType("summation")
	// The counter is never rolled up.
	PerfSummaryTypeNone = PerfSummaryType("none")
)

func (e PerfSummaryType) Values() []PerfSummaryType {
	return []PerfSummaryType{
		PerfSummaryTypeAverage,
		PerfSummaryTypeMaximum,
		PerfSummaryTypeMinimum,
		PerfSummaryTypeLatest,
		PerfSummaryTypeSummation,
		PerfSummaryTypeNone,
	}
}

func (e PerfSummaryType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PerfSummaryType"] = reflect.TypeOf((*PerfSummaryType)(nil)).Elem()
}

// Indicates the unit of measure represented by a counter or statistical
// value.
type PerformanceManagerUnit string

const (
	// Percentage values in units of 1/100th of a percent.
	//
	// For example 100
	// represents 1%.
	PerformanceManagerUnitPercent = PerformanceManagerUnit("percent")
	// Kilobytes.
	PerformanceManagerUnitKiloBytes = PerformanceManagerUnit("kiloBytes")
	// Megabytes.
	PerformanceManagerUnitMegaBytes = PerformanceManagerUnit("megaBytes")
	// Megahertz.
	PerformanceManagerUnitMegaHertz = PerformanceManagerUnit("megaHertz")
	// A quantity of items, for example, the number of CPUs.
	PerformanceManagerUnitNumber = PerformanceManagerUnit("number")
	// The time in microseconds.
	PerformanceManagerUnitMicrosecond = PerformanceManagerUnit("microsecond")
	// The time in milliseconds.
	PerformanceManagerUnitMillisecond = PerformanceManagerUnit("millisecond")
	// The time in seconds.
	PerformanceManagerUnitSecond = PerformanceManagerUnit("second")
	// Kilobytes per second.
	PerformanceManagerUnitKiloBytesPerSecond = PerformanceManagerUnit("kiloBytesPerSecond")
	// Megabytes per second.
	PerformanceManagerUnitMegaBytesPerSecond = PerformanceManagerUnit("megaBytesPerSecond")
	// Watts
	PerformanceManagerUnitWatt = PerformanceManagerUnit("watt")
	// Joules
	PerformanceManagerUnitJoule = PerformanceManagerUnit("joule")
	// Terabytes.
	PerformanceManagerUnitTeraBytes = PerformanceManagerUnit("teraBytes")
	// Temperature in celsius.
	PerformanceManagerUnitCelsius = PerformanceManagerUnit("celsius")
	// The time in nanoseconds.
	PerformanceManagerUnitNanosecond = PerformanceManagerUnit("nanosecond")
)

func (e PerformanceManagerUnit) Values() []PerformanceManagerUnit {
	return []PerformanceManagerUnit{
		PerformanceManagerUnitPercent,
		PerformanceManagerUnitKiloBytes,
		PerformanceManagerUnitMegaBytes,
		PerformanceManagerUnitMegaHertz,
		PerformanceManagerUnitNumber,
		PerformanceManagerUnitMicrosecond,
		PerformanceManagerUnitMillisecond,
		PerformanceManagerUnitSecond,
		PerformanceManagerUnitKiloBytesPerSecond,
		PerformanceManagerUnitMegaBytesPerSecond,
		PerformanceManagerUnitWatt,
		PerformanceManagerUnitJoule,
		PerformanceManagerUnitTeraBytes,
		PerformanceManagerUnitCelsius,
		PerformanceManagerUnitNanosecond,
	}
}

func (e PerformanceManagerUnit) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PerformanceManagerUnit"] = reflect.TypeOf((*PerformanceManagerUnit)(nil)).Elem()
	minAPIVersionForEnumValue["PerformanceManagerUnit"] = map[string]string{
		"nanosecond": "8.0.0.1",
	}
}

type PhysicalNicResourcePoolSchedulerDisallowedReason string

const (
	// Indicates that the user has opted out the Physical NIC from resource pool
	// based scheduling.
	PhysicalNicResourcePoolSchedulerDisallowedReasonUserOptOut = PhysicalNicResourcePoolSchedulerDisallowedReason("userOptOut")
	// Indicates that the NIC device does is not capable of resource pool
	// based scheduling.
	PhysicalNicResourcePoolSchedulerDisallowedReasonHardwareUnsupported = PhysicalNicResourcePoolSchedulerDisallowedReason("hardwareUnsupported")
)

func (e PhysicalNicResourcePoolSchedulerDisallowedReason) Values() []PhysicalNicResourcePoolSchedulerDisallowedReason {
	return []PhysicalNicResourcePoolSchedulerDisallowedReason{
		PhysicalNicResourcePoolSchedulerDisallowedReasonUserOptOut,
		PhysicalNicResourcePoolSchedulerDisallowedReasonHardwareUnsupported,
	}
}

func (e PhysicalNicResourcePoolSchedulerDisallowedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PhysicalNicResourcePoolSchedulerDisallowedReason"] = reflect.TypeOf((*PhysicalNicResourcePoolSchedulerDisallowedReason)(nil)).Elem()
}

// Set of possible values for `PhysicalNic.vmDirectPathGen2SupportedMode`.
type PhysicalNicVmDirectPathGen2SupportedMode string

const (
	PhysicalNicVmDirectPathGen2SupportedModeUpt = PhysicalNicVmDirectPathGen2SupportedMode("upt")
)

func (e PhysicalNicVmDirectPathGen2SupportedMode) Values() []PhysicalNicVmDirectPathGen2SupportedMode {
	return []PhysicalNicVmDirectPathGen2SupportedMode{
		PhysicalNicVmDirectPathGen2SupportedModeUpt,
	}
}

func (e PhysicalNicVmDirectPathGen2SupportedMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PhysicalNicVmDirectPathGen2SupportedMode"] = reflect.TypeOf((*PhysicalNicVmDirectPathGen2SupportedMode)(nil)).Elem()
}

// Rule scope determines conditions when an affinity rule is
// satisfied.
//
// The following uses affinity rule as example.
// cluster: All Vms in the rule list are placed in a single cluster.
// host: All Vms in the rule list are placed in a single host.
// storagePod: All Vms in the rule list are placed in a single storagePod.
// datastore: All Vms in the rule list are placed in a single datastore.
type PlacementAffinityRuleRuleScope string

const (
	// clusters are the scope
	PlacementAffinityRuleRuleScopeCluster = PlacementAffinityRuleRuleScope("cluster")
	// individual hosts are the scope
	PlacementAffinityRuleRuleScopeHost = PlacementAffinityRuleRuleScope("host")
	// datastore cluster is teh scope
	PlacementAffinityRuleRuleScopeStoragePod = PlacementAffinityRuleRuleScope("storagePod")
	// individual datastores are the scope
	PlacementAffinityRuleRuleScopeDatastore = PlacementAffinityRuleRuleScope("datastore")
)

func (e PlacementAffinityRuleRuleScope) Values() []PlacementAffinityRuleRuleScope {
	return []PlacementAffinityRuleRuleScope{
		PlacementAffinityRuleRuleScopeCluster,
		PlacementAffinityRuleRuleScopeHost,
		PlacementAffinityRuleRuleScopeStoragePod,
		PlacementAffinityRuleRuleScopeDatastore,
	}
}

func (e PlacementAffinityRuleRuleScope) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PlacementAffinityRuleRuleScope"] = reflect.TypeOf((*PlacementAffinityRuleRuleScope)(nil)).Elem()
}

// Rule type determines how the affinity rule is to be enforced:
// affinity: Vms in the list are kept together within the rule
// scope.
//
// anti-affinity: Vms in the rule list are kept separate
// across the objects in the rule scope.
// soft rule: The enforcement is best effort.
type PlacementAffinityRuleRuleType string

const (
	// Affinity
	PlacementAffinityRuleRuleTypeAffinity = PlacementAffinityRuleRuleType("affinity")
	// Anti-Affinity
	PlacementAffinityRuleRuleTypeAntiAffinity = PlacementAffinityRuleRuleType("antiAffinity")
	// Best-effort affinity
	PlacementAffinityRuleRuleTypeSoftAffinity = PlacementAffinityRuleRuleType("softAffinity")
	// Best-effort anti-affinity
	PlacementAffinityRuleRuleTypeSoftAntiAffinity = PlacementAffinityRuleRuleType("softAntiAffinity")
)

func (e PlacementAffinityRuleRuleType) Values() []PlacementAffinityRuleRuleType {
	return []PlacementAffinityRuleRuleType{
		PlacementAffinityRuleRuleTypeAffinity,
		PlacementAffinityRuleRuleTypeAntiAffinity,
		PlacementAffinityRuleRuleTypeSoftAffinity,
		PlacementAffinityRuleRuleTypeSoftAntiAffinity,
	}
}

func (e PlacementAffinityRuleRuleType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PlacementAffinityRuleRuleType"] = reflect.TypeOf((*PlacementAffinityRuleRuleType)(nil)).Elem()
}

// Defines the type of placement
type PlacementSpecPlacementType string

const (
	// Create a new VM
	PlacementSpecPlacementTypeCreate = PlacementSpecPlacementType("create")
	// Reconfigure a VM
	PlacementSpecPlacementTypeReconfigure = PlacementSpecPlacementType("reconfigure")
	// Relocate a VM
	PlacementSpecPlacementTypeRelocate = PlacementSpecPlacementType("relocate")
	// Clone a VM
	PlacementSpecPlacementTypeClone = PlacementSpecPlacementType("clone")
)

func (e PlacementSpecPlacementType) Values() []PlacementSpecPlacementType {
	return []PlacementSpecPlacementType{
		PlacementSpecPlacementTypeCreate,
		PlacementSpecPlacementTypeReconfigure,
		PlacementSpecPlacementTypeRelocate,
		PlacementSpecPlacementTypeClone,
	}
}

func (e PlacementSpecPlacementType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PlacementSpecPlacementType"] = reflect.TypeOf((*PlacementSpecPlacementType)(nil)).Elem()
}

// The type of component connected to a port group.
type PortGroupConnecteeType string

const (
	// A virtual machine is connected to this port group.
	PortGroupConnecteeTypeVirtualMachine = PortGroupConnecteeType("virtualMachine")
	// A system management entity (service console)
	// is connected to this port group.
	PortGroupConnecteeTypeSystemManagement = PortGroupConnecteeType("systemManagement")
	// The VMkernel is connected to this port group.
	PortGroupConnecteeTypeHost = PortGroupConnecteeType("host")
	// This port group serves an entity of unspecified kind.
	PortGroupConnecteeTypeUnknown = PortGroupConnecteeType("unknown")
)

func (e PortGroupConnecteeType) Values() []PortGroupConnecteeType {
	return []PortGroupConnecteeType{
		PortGroupConnecteeTypeVirtualMachine,
		PortGroupConnecteeTypeSystemManagement,
		PortGroupConnecteeTypeHost,
		PortGroupConnecteeTypeUnknown,
	}
}

func (e PortGroupConnecteeType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PortGroupConnecteeType"] = reflect.TypeOf((*PortGroupConnecteeType)(nil)).Elem()
}

// Defines the result status values for a
// `HostProfile*.*HostProfile.ExecuteHostProfile`
// operation.
//
// The result data is contained in the
// `ProfileExecuteResult` data object.
type ProfileExecuteResultStatus string

const (
	// Profile execution was successful.
	//
	// You can use the output configuration data
	// to apply the profile to a host.
	ProfileExecuteResultStatusSuccess = ProfileExecuteResultStatus("success")
	// Additional data is required to complete the operation.
	//
	// The data requirements are defined in the list of policy options for the profile
	// (`ApplyProfile*.*ApplyProfile.policy`\[\]).
	ProfileExecuteResultStatusNeedInput = ProfileExecuteResultStatus("needInput")
	// Profile execution generated an error.
	//
	// See `ProfileExecuteResult*.*ProfileExecuteResult.error`.
	ProfileExecuteResultStatusError = ProfileExecuteResultStatus("error")
)

func (e ProfileExecuteResultStatus) Values() []ProfileExecuteResultStatus {
	return []ProfileExecuteResultStatus{
		ProfileExecuteResultStatusSuccess,
		ProfileExecuteResultStatusNeedInput,
		ProfileExecuteResultStatusError,
	}
}

func (e ProfileExecuteResultStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ProfileExecuteResultStatus"] = reflect.TypeOf((*ProfileExecuteResultStatus)(nil)).Elem()
}

// Enumerates different operations supported for comparing
// numerical values.
type ProfileNumericComparator string

const (
	ProfileNumericComparatorLessThan         = ProfileNumericComparator("lessThan")
	ProfileNumericComparatorLessThanEqual    = ProfileNumericComparator("lessThanEqual")
	ProfileNumericComparatorEqual            = ProfileNumericComparator("equal")
	ProfileNumericComparatorNotEqual         = ProfileNumericComparator("notEqual")
	ProfileNumericComparatorGreaterThanEqual = ProfileNumericComparator("greaterThanEqual")
	ProfileNumericComparatorGreaterThan      = ProfileNumericComparator("greaterThan")
)

func (e ProfileNumericComparator) Values() []ProfileNumericComparator {
	return []ProfileNumericComparator{
		ProfileNumericComparatorLessThan,
		ProfileNumericComparatorLessThanEqual,
		ProfileNumericComparatorEqual,
		ProfileNumericComparatorNotEqual,
		ProfileNumericComparatorGreaterThanEqual,
		ProfileNumericComparatorGreaterThan,
	}
}

func (e ProfileNumericComparator) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ProfileNumericComparator"] = reflect.TypeOf((*ProfileNumericComparator)(nil)).Elem()
}

// The relation type to be supported.
type ProfileParameterMetadataRelationType string

const (
	// The relation to a subprofile or a parameter.
	ProfileParameterMetadataRelationTypeDynamic_relation = ProfileParameterMetadataRelationType("dynamic_relation")
	// The values from sources other than the parameter/profile or the static
	// value list are allowed.
	ProfileParameterMetadataRelationTypeExtensible_relation = ProfileParameterMetadataRelationType("extensible_relation")
	// The value list contains localization keys instead of values.
	ProfileParameterMetadataRelationTypeLocalizable_relation = ProfileParameterMetadataRelationType("localizable_relation")
	// The relation is defined by static valid value list.
	ProfileParameterMetadataRelationTypeStatic_relation = ProfileParameterMetadataRelationType("static_relation")
	// The relation is defined for validation purpose.
	ProfileParameterMetadataRelationTypeValidation_relation = ProfileParameterMetadataRelationType("validation_relation")
)

func (e ProfileParameterMetadataRelationType) Values() []ProfileParameterMetadataRelationType {
	return []ProfileParameterMetadataRelationType{
		ProfileParameterMetadataRelationTypeDynamic_relation,
		ProfileParameterMetadataRelationTypeExtensible_relation,
		ProfileParameterMetadataRelationTypeLocalizable_relation,
		ProfileParameterMetadataRelationTypeStatic_relation,
		ProfileParameterMetadataRelationTypeValidation_relation,
	}
}

func (e ProfileParameterMetadataRelationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ProfileParameterMetadataRelationType"] = reflect.TypeOf((*ProfileParameterMetadataRelationType)(nil)).Elem()
}

// Enumeration of possible changes to a property.
type PropertyChangeOp string

const (
	PropertyChangeOpAdd            = PropertyChangeOp("add")
	PropertyChangeOpRemove         = PropertyChangeOp("remove")
	PropertyChangeOpAssign         = PropertyChangeOp("assign")
	PropertyChangeOpIndirectRemove = PropertyChangeOp("indirectRemove")
)

func (e PropertyChangeOp) Values() []PropertyChangeOp {
	return []PropertyChangeOp{
		PropertyChangeOpAdd,
		PropertyChangeOpRemove,
		PropertyChangeOpAssign,
		PropertyChangeOpIndirectRemove,
	}
}

func (e PropertyChangeOp) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["PropertyChangeOp"] = reflect.TypeOf((*PropertyChangeOp)(nil)).Elem()
}

type QuarantineModeFaultFaultType string

const (
	// The cluster does not contain any non-quarantined host satisfying the
	// VM/host affinity rules for the VM.
	QuarantineModeFaultFaultTypeNoCompatibleNonQuarantinedHost = QuarantineModeFaultFaultType("NoCompatibleNonQuarantinedHost")
	// The current DRS migration priority setting disallows generating a
	// recommendation to prevent VMs on quarantined hosts.
	//
	// Thus, the
	// violation will not be corrected.
	QuarantineModeFaultFaultTypeCorrectionDisallowed = QuarantineModeFaultFaultType("CorrectionDisallowed")
	// DRS has determined that evacuation of VMs from quarantined hosts
	// impacts respecting cluster constraints or performance goals so they
	// are not evacuated.
	QuarantineModeFaultFaultTypeCorrectionImpact = QuarantineModeFaultFaultType("CorrectionImpact")
)

func (e QuarantineModeFaultFaultType) Values() []QuarantineModeFaultFaultType {
	return []QuarantineModeFaultFaultType{
		QuarantineModeFaultFaultTypeNoCompatibleNonQuarantinedHost,
		QuarantineModeFaultFaultTypeCorrectionDisallowed,
		QuarantineModeFaultFaultTypeCorrectionImpact,
	}
}

func (e QuarantineModeFaultFaultType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["QuarantineModeFaultFaultType"] = reflect.TypeOf((*QuarantineModeFaultFaultType)(nil)).Elem()
}

// Quiescing is a boolean flag in `ReplicationConfigSpec`
// and QuiesceModeType describes the supported quiesce mode
// for `VirtualMachine`.
//
// If application quiescing fails, HBR would attempt
// filesystem quiescing and if even filesystem quiescing
// fails, then we would just create a crash consistent
// instance.
type QuiesceMode string

const (
	// HBR supports application quescing for this
	// `VirtualMachine`.
	QuiesceModeApplication = QuiesceMode("application")
	// HBR supports filesystem quescing for this
	// `VirtualMachine`.
	QuiesceModeFilesystem = QuiesceMode("filesystem")
	// HBR does not support quescing for this
	// `VirtualMachine`.
	QuiesceModeNone = QuiesceMode("none")
)

func (e QuiesceMode) Values() []QuiesceMode {
	return []QuiesceMode{
		QuiesceModeApplication,
		QuiesceModeFilesystem,
		QuiesceModeNone,
	}
}

func (e QuiesceMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["QuiesceMode"] = reflect.TypeOf((*QuiesceMode)(nil)).Elem()
}

// List of defined migration reason codes:
type RecommendationReasonCode string

const (
	// Balance average CPU utilization.
	RecommendationReasonCodeFairnessCpuAvg = RecommendationReasonCode("fairnessCpuAvg")
	// Balance average memory utilization.
	RecommendationReasonCodeFairnessMemAvg = RecommendationReasonCode("fairnessMemAvg")
	// Fulfill affinity rule.
	RecommendationReasonCodeJointAffin = RecommendationReasonCode("jointAffin")
	// Fulfill anti-affinity rule.
	RecommendationReasonCodeAntiAffin = RecommendationReasonCode("antiAffin")
	// Host entering maintenance mode.
	RecommendationReasonCodeHostMaint = RecommendationReasonCode("hostMaint")
	// Host entering standby mode.
	RecommendationReasonCodeEnterStandby = RecommendationReasonCode("enterStandby")
	// balance CPU reservations
	RecommendationReasonCodeReservationCpu = RecommendationReasonCode("reservationCpu")
	// balance memory reservations
	RecommendationReasonCodeReservationMem = RecommendationReasonCode("reservationMem")
	// Power on virtual machine
	RecommendationReasonCodePowerOnVm = RecommendationReasonCode("powerOnVm")
	// Power off host for power savings
	RecommendationReasonCodePowerSaving = RecommendationReasonCode("powerSaving")
	// Power on host to increase cluster capacity
	RecommendationReasonCodeIncreaseCapacity = RecommendationReasonCode("increaseCapacity")
	// Sanity-check resource pool hierarchy
	RecommendationReasonCodeCheckResource = RecommendationReasonCode("checkResource")
	// Maintain unreserved capacity
	RecommendationReasonCodeUnreservedCapacity = RecommendationReasonCode("unreservedCapacity")
	// Fix hard VM/host affinity rule violation
	RecommendationReasonCodeVmHostHardAffinity = RecommendationReasonCode("vmHostHardAffinity")
	// Fix soft VM/host affinity rule violation
	RecommendationReasonCodeVmHostSoftAffinity = RecommendationReasonCode("vmHostSoftAffinity")
	// Balance datastore space usage.
	RecommendationReasonCodeBalanceDatastoreSpaceUsage = RecommendationReasonCode("balanceDatastoreSpaceUsage")
	// Deprecated as of vSphere8.0 U3, and there is no replacement for it.
	//
	// Balance datastore I/O workload.
	RecommendationReasonCodeBalanceDatastoreIOLoad = RecommendationReasonCode("balanceDatastoreIOLoad")
	// Deprecated as of vSphere8.0 U3, and there is no replacement for it.
	//
	// Balance datastore IOPS reservation
	RecommendationReasonCodeBalanceDatastoreIOPSReservation = RecommendationReasonCode("balanceDatastoreIOPSReservation")
	// Datastore entering maintenance mode.
	RecommendationReasonCodeDatastoreMaint = RecommendationReasonCode("datastoreMaint")
	// Fix virtual disk affinity rule violation.
	RecommendationReasonCodeVirtualDiskJointAffin = RecommendationReasonCode("virtualDiskJointAffin")
	// Fix virtual disk anti-affinity rule violation.
	RecommendationReasonCodeVirtualDiskAntiAffin = RecommendationReasonCode("virtualDiskAntiAffin")
	// Fix the issue that a datastore run out of space.
	RecommendationReasonCodeDatastoreSpaceOutage = RecommendationReasonCode("datastoreSpaceOutage")
	// Satisfy storage initial placement requests.
	RecommendationReasonCodeStoragePlacement = RecommendationReasonCode("storagePlacement")
	// Deprecated as of vSphere8.0 U3, and there is no replacement for it.
	//
	// IO load balancing was disabled internally.
	RecommendationReasonCodeIolbDisabledInternal = RecommendationReasonCode("iolbDisabledInternal")
	// Satisfy unified vmotion placement requests.
	RecommendationReasonCodeXvmotionPlacement = RecommendationReasonCode("xvmotionPlacement")
	// Fix network bandwidth reservation violation
	RecommendationReasonCodeNetworkBandwidthReservation = RecommendationReasonCode("networkBandwidthReservation")
	// Host is partially degraded.
	RecommendationReasonCodeHostInDegradation = RecommendationReasonCode("hostInDegradation")
	// Host is not degraded.
	RecommendationReasonCodeHostExitDegradation = RecommendationReasonCode("hostExitDegradation")
	// Fix maxVms constraint violation
	RecommendationReasonCodeMaxVmsConstraint = RecommendationReasonCode("maxVmsConstraint")
	// Fix ft maxVMs and maxVcpus constraint violations
	RecommendationReasonCodeFtConstraints = RecommendationReasonCode("ftConstraints")
	// Fix VM/host affinity policy violation
	RecommendationReasonCodeVmHostAffinityPolicy = RecommendationReasonCode("vmHostAffinityPolicy")
	// Fix VM/host anti-affinity policy violation
	RecommendationReasonCodeVmHostAntiAffinityPolicy = RecommendationReasonCode("vmHostAntiAffinityPolicy")
	// Fix VM-VM anti-affinity policy violations
	RecommendationReasonCodeVmAntiAffinityPolicy = RecommendationReasonCode("vmAntiAffinityPolicy")
	// `**Since:**` vSphere API Release 7.0.2.0
	RecommendationReasonCodeBalanceVsanUsage = RecommendationReasonCode("balanceVsanUsage")
	// Optimize assignable hardware resource orchestration
	RecommendationReasonCodeAhPlacementOptimization = RecommendationReasonCode("ahPlacementOptimization")
	// Upgrade virtual machine to new vmx binary
	RecommendationReasonCodeVmxUpgrade = RecommendationReasonCode("vmxUpgrade")
)

func (e RecommendationReasonCode) Values() []RecommendationReasonCode {
	return []RecommendationReasonCode{
		RecommendationReasonCodeFairnessCpuAvg,
		RecommendationReasonCodeFairnessMemAvg,
		RecommendationReasonCodeJointAffin,
		RecommendationReasonCodeAntiAffin,
		RecommendationReasonCodeHostMaint,
		RecommendationReasonCodeEnterStandby,
		RecommendationReasonCodeReservationCpu,
		RecommendationReasonCodeReservationMem,
		RecommendationReasonCodePowerOnVm,
		RecommendationReasonCodePowerSaving,
		RecommendationReasonCodeIncreaseCapacity,
		RecommendationReasonCodeCheckResource,
		RecommendationReasonCodeUnreservedCapacity,
		RecommendationReasonCodeVmHostHardAffinity,
		RecommendationReasonCodeVmHostSoftAffinity,
		RecommendationReasonCodeBalanceDatastoreSpaceUsage,
		RecommendationReasonCodeBalanceDatastoreIOLoad,
		RecommendationReasonCodeBalanceDatastoreIOPSReservation,
		RecommendationReasonCodeDatastoreMaint,
		RecommendationReasonCodeVirtualDiskJointAffin,
		RecommendationReasonCodeVirtualDiskAntiAffin,
		RecommendationReasonCodeDatastoreSpaceOutage,
		RecommendationReasonCodeStoragePlacement,
		RecommendationReasonCodeIolbDisabledInternal,
		RecommendationReasonCodeXvmotionPlacement,
		RecommendationReasonCodeNetworkBandwidthReservation,
		RecommendationReasonCodeHostInDegradation,
		RecommendationReasonCodeHostExitDegradation,
		RecommendationReasonCodeMaxVmsConstraint,
		RecommendationReasonCodeFtConstraints,
		RecommendationReasonCodeVmHostAffinityPolicy,
		RecommendationReasonCodeVmHostAntiAffinityPolicy,
		RecommendationReasonCodeVmAntiAffinityPolicy,
		RecommendationReasonCodeBalanceVsanUsage,
		RecommendationReasonCodeAhPlacementOptimization,
		RecommendationReasonCodeVmxUpgrade,
	}
}

func (e RecommendationReasonCode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["RecommendationReasonCode"] = reflect.TypeOf((*RecommendationReasonCode)(nil)).Elem()
	minAPIVersionForEnumValue["RecommendationReasonCode"] = map[string]string{
		"balanceVsanUsage":        "7.0.2.0",
		"ahPlacementOptimization": "8.0.2.0",
		"vmxUpgrade":              "8.0.3.0",
	}
}

// Pre-defined constants for possible recommendation types.
//
// Virtual Center
// uses this information to coordinate with the clients.
type RecommendationType string

const (
	RecommendationTypeV1 = RecommendationType("V1")
)

func (e RecommendationType) Values() []RecommendationType {
	return []RecommendationType{
		RecommendationTypeV1,
	}
}

func (e RecommendationType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["RecommendationType"] = reflect.TypeOf((*RecommendationType)(nil)).Elem()
}

type ReplicationDiskConfigFaultReasonForFault string

const (
	// Could not look up device by key
	ReplicationDiskConfigFaultReasonForFaultDiskNotFound = ReplicationDiskConfigFaultReasonForFault("diskNotFound")
	// Replication not supported for disk type or backend
	ReplicationDiskConfigFaultReasonForFaultDiskTypeNotSupported = ReplicationDiskConfigFaultReasonForFault("diskTypeNotSupported")
	// Invalid key value
	ReplicationDiskConfigFaultReasonForFaultInvalidDiskKey = ReplicationDiskConfigFaultReasonForFault("invalidDiskKey")
	// Invalid disk replication ID string
	ReplicationDiskConfigFaultReasonForFaultInvalidDiskReplicationId = ReplicationDiskConfigFaultReasonForFault("invalidDiskReplicationId")
	// Another disk in the VM has the same replication ID
	ReplicationDiskConfigFaultReasonForFaultDuplicateDiskReplicationId = ReplicationDiskConfigFaultReasonForFault("duplicateDiskReplicationId")
	// Invalid path (string) for the persistent file
	ReplicationDiskConfigFaultReasonForFaultInvalidPersistentFilePath = ReplicationDiskConfigFaultReasonForFault("invalidPersistentFilePath")
	// Attempting to re-configure the disk's replication ID
	ReplicationDiskConfigFaultReasonForFaultReconfigureDiskReplicationIdNotAllowed = ReplicationDiskConfigFaultReasonForFault("reconfigureDiskReplicationIdNotAllowed")
)

func (e ReplicationDiskConfigFaultReasonForFault) Values() []ReplicationDiskConfigFaultReasonForFault {
	return []ReplicationDiskConfigFaultReasonForFault{
		ReplicationDiskConfigFaultReasonForFaultDiskNotFound,
		ReplicationDiskConfigFaultReasonForFaultDiskTypeNotSupported,
		ReplicationDiskConfigFaultReasonForFaultInvalidDiskKey,
		ReplicationDiskConfigFaultReasonForFaultInvalidDiskReplicationId,
		ReplicationDiskConfigFaultReasonForFaultDuplicateDiskReplicationId,
		ReplicationDiskConfigFaultReasonForFaultInvalidPersistentFilePath,
		ReplicationDiskConfigFaultReasonForFaultReconfigureDiskReplicationIdNotAllowed,
	}
}

func (e ReplicationDiskConfigFaultReasonForFault) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ReplicationDiskConfigFaultReasonForFault"] = reflect.TypeOf((*ReplicationDiskConfigFaultReasonForFault)(nil)).Elem()
}

type ReplicationVmConfigFaultReasonForFault string

const (
	// Incompatible VM hardware version
	ReplicationVmConfigFaultReasonForFaultIncompatibleHwVersion = ReplicationVmConfigFaultReasonForFault("incompatibleHwVersion")
	// Invalid VM Replication ID string
	ReplicationVmConfigFaultReasonForFaultInvalidVmReplicationId = ReplicationVmConfigFaultReasonForFault("invalidVmReplicationId")
	// Invalid generation number in VM's configuration
	ReplicationVmConfigFaultReasonForFaultInvalidGenerationNumber = ReplicationVmConfigFaultReasonForFault("invalidGenerationNumber")
	// Invalid RPO value (out of bounds)
	ReplicationVmConfigFaultReasonForFaultOutOfBoundsRpoValue = ReplicationVmConfigFaultReasonForFault("outOfBoundsRpoValue")
	// Invalid destination IP address
	ReplicationVmConfigFaultReasonForFaultInvalidDestinationIpAddress = ReplicationVmConfigFaultReasonForFault("invalidDestinationIpAddress")
	// Invalid destination port
	ReplicationVmConfigFaultReasonForFaultInvalidDestinationPort = ReplicationVmConfigFaultReasonForFault("invalidDestinationPort")
	// Malformed extra options list
	ReplicationVmConfigFaultReasonForFaultInvalidExtraVmOptions = ReplicationVmConfigFaultReasonForFault("invalidExtraVmOptions")
	// Mis-matching generation number (stale)
	ReplicationVmConfigFaultReasonForFaultStaleGenerationNumber = ReplicationVmConfigFaultReasonForFault("staleGenerationNumber")
	// Attempting to re-configure the VM replication ID
	ReplicationVmConfigFaultReasonForFaultReconfigureVmReplicationIdNotAllowed = ReplicationVmConfigFaultReasonForFault("reconfigureVmReplicationIdNotAllowed")
	// Could not retrieve the VM configuration
	ReplicationVmConfigFaultReasonForFaultCannotRetrieveVmReplicationConfiguration = ReplicationVmConfigFaultReasonForFault("cannotRetrieveVmReplicationConfiguration")
	// Attempting to re-enable replication for the VM
	ReplicationVmConfigFaultReasonForFaultReplicationAlreadyEnabled = ReplicationVmConfigFaultReasonForFault("replicationAlreadyEnabled")
	// The existing replication configuration of the VM is broken
	// (applicable to re-configuration only).
	ReplicationVmConfigFaultReasonForFaultInvalidPriorConfiguration = ReplicationVmConfigFaultReasonForFault("invalidPriorConfiguration")
	// Attempting to re-configure or disable replication for a VM
	// for which replication has not been enabled.
	ReplicationVmConfigFaultReasonForFaultReplicationNotEnabled = ReplicationVmConfigFaultReasonForFault("replicationNotEnabled")
	// Failed to commit the new replication properties for the VM.
	ReplicationVmConfigFaultReasonForFaultReplicationConfigurationFailed = ReplicationVmConfigFaultReasonForFault("replicationConfigurationFailed")
	// VM is encrypted
	ReplicationVmConfigFaultReasonForFaultEncryptedVm = ReplicationVmConfigFaultReasonForFault("encryptedVm")
	// Remote certificate thumbprint is invalid
	ReplicationVmConfigFaultReasonForFaultInvalidThumbprint = ReplicationVmConfigFaultReasonForFault("invalidThumbprint")
	// VM hardware contains devices incompatible with replication
	ReplicationVmConfigFaultReasonForFaultIncompatibleDevice = ReplicationVmConfigFaultReasonForFault("incompatibleDevice")
)

func (e ReplicationVmConfigFaultReasonForFault) Values() []ReplicationVmConfigFaultReasonForFault {
	return []ReplicationVmConfigFaultReasonForFault{
		ReplicationVmConfigFaultReasonForFaultIncompatibleHwVersion,
		ReplicationVmConfigFaultReasonForFaultInvalidVmReplicationId,
		ReplicationVmConfigFaultReasonForFaultInvalidGenerationNumber,
		ReplicationVmConfigFaultReasonForFaultOutOfBoundsRpoValue,
		ReplicationVmConfigFaultReasonForFaultInvalidDestinationIpAddress,
		ReplicationVmConfigFaultReasonForFaultInvalidDestinationPort,
		ReplicationVmConfigFaultReasonForFaultInvalidExtraVmOptions,
		ReplicationVmConfigFaultReasonForFaultStaleGenerationNumber,
		ReplicationVmConfigFaultReasonForFaultReconfigureVmReplicationIdNotAllowed,
		ReplicationVmConfigFaultReasonForFaultCannotRetrieveVmReplicationConfiguration,
		ReplicationVmConfigFaultReasonForFaultReplicationAlreadyEnabled,
		ReplicationVmConfigFaultReasonForFaultInvalidPriorConfiguration,
		ReplicationVmConfigFaultReasonForFaultReplicationNotEnabled,
		ReplicationVmConfigFaultReasonForFaultReplicationConfigurationFailed,
		ReplicationVmConfigFaultReasonForFaultEncryptedVm,
		ReplicationVmConfigFaultReasonForFaultInvalidThumbprint,
		ReplicationVmConfigFaultReasonForFaultIncompatibleDevice,
	}
}

func (e ReplicationVmConfigFaultReasonForFault) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ReplicationVmConfigFaultReasonForFault"] = reflect.TypeOf((*ReplicationVmConfigFaultReasonForFault)(nil)).Elem()
}

type ReplicationVmFaultReasonForFault string

const (
	// `VirtualMachine` is not configured for replication
	ReplicationVmFaultReasonForFaultNotConfigured = ReplicationVmFaultReasonForFault("notConfigured")
	// `VirtualMachine` is powered off (and is not undergoing
	// offline replication)
	ReplicationVmFaultReasonForFaultPoweredOff = ReplicationVmFaultReasonForFault("poweredOff")
	// `VirtualMachine` is suspended (and is not undergoing
	// offline replication)
	ReplicationVmFaultReasonForFaultSuspended = ReplicationVmFaultReasonForFault("suspended")
	// `VirtualMachine` is powered on
	ReplicationVmFaultReasonForFaultPoweredOn = ReplicationVmFaultReasonForFault("poweredOn")
	// `VirtualMachine` is in the process of creating an
	// an offline instance.
	ReplicationVmFaultReasonForFaultOfflineReplicating = ReplicationVmFaultReasonForFault("offlineReplicating")
	// `VirtualMachine` is in an invalid state
	ReplicationVmFaultReasonForFaultInvalidState = ReplicationVmFaultReasonForFault("invalidState")
	// The specified instanceId does not match the `VirtualMachine`
	// instanceId
	ReplicationVmFaultReasonForFaultInvalidInstanceId = ReplicationVmFaultReasonForFault("invalidInstanceId")
	// `VirtualMachine` is in the process of creating an
	// offline instance and we are trying to disable it.
	//
	// The first step is to close the offline disk. If closing disks
	// is not successful, throw this fault.
	ReplicationVmFaultReasonForFaultCloseDiskError = ReplicationVmFaultReasonForFault("closeDiskError")
	// `VirtualMachine` is trying to create a group already
	// owned by another VM.
	ReplicationVmFaultReasonForFaultGroupExist = ReplicationVmFaultReasonForFault("groupExist")
)

func (e ReplicationVmFaultReasonForFault) Values() []ReplicationVmFaultReasonForFault {
	return []ReplicationVmFaultReasonForFault{
		ReplicationVmFaultReasonForFaultNotConfigured,
		ReplicationVmFaultReasonForFaultPoweredOff,
		ReplicationVmFaultReasonForFaultSuspended,
		ReplicationVmFaultReasonForFaultPoweredOn,
		ReplicationVmFaultReasonForFaultOfflineReplicating,
		ReplicationVmFaultReasonForFaultInvalidState,
		ReplicationVmFaultReasonForFaultInvalidInstanceId,
		ReplicationVmFaultReasonForFaultCloseDiskError,
		ReplicationVmFaultReasonForFaultGroupExist,
	}
}

func (e ReplicationVmFaultReasonForFault) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ReplicationVmFaultReasonForFault"] = reflect.TypeOf((*ReplicationVmFaultReasonForFault)(nil)).Elem()
}

type ReplicationVmInProgressFaultActivity string

const (
	// Initial synchronization with the remote site
	ReplicationVmInProgressFaultActivityFullSync = ReplicationVmInProgressFaultActivity("fullSync")
	// Delta updates to generate a consistent instance
	ReplicationVmInProgressFaultActivityDelta = ReplicationVmInProgressFaultActivity("delta")
)

func (e ReplicationVmInProgressFaultActivity) Values() []ReplicationVmInProgressFaultActivity {
	return []ReplicationVmInProgressFaultActivity{
		ReplicationVmInProgressFaultActivityFullSync,
		ReplicationVmInProgressFaultActivityDelta,
	}
}

func (e ReplicationVmInProgressFaultActivity) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ReplicationVmInProgressFaultActivity"] = reflect.TypeOf((*ReplicationVmInProgressFaultActivity)(nil)).Elem()
}

// Describes the current state of a replicated `VirtualMachine`
type ReplicationVmState string

const (
	// The `VirtualMachine` has no current replication state.
	//
	// This is a virtual machine that is configured for replication, but is
	// powered off and not undergoing offline replication.
	ReplicationVmStateNone = ReplicationVmState("none")
	// The `VirtualMachine` replication is paused.
	ReplicationVmStatePaused = ReplicationVmState("paused")
	// One or more of the `VirtualMachine` disks is in the
	// process of an initial synchronization with the remote site.
	ReplicationVmStateSyncing = ReplicationVmState("syncing")
	// The `VirtualMachine` is being replicated but is not
	// currently in the process of having a consistent instance created.
	ReplicationVmStateIdle = ReplicationVmState("idle")
	// The `VirtualMachine` is in the process of having
	// a consistent instance created.
	ReplicationVmStateActive = ReplicationVmState("active")
	// The `VirtualMachine` is unable to replicate due to
	// errors.
	//
	// XXX Currently unused.
	ReplicationVmStateError = ReplicationVmState("error")
)

func (e ReplicationVmState) Values() []ReplicationVmState {
	return []ReplicationVmState{
		ReplicationVmStateNone,
		ReplicationVmStatePaused,
		ReplicationVmStateSyncing,
		ReplicationVmStateIdle,
		ReplicationVmStateActive,
		ReplicationVmStateError,
	}
}

func (e ReplicationVmState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ReplicationVmState"] = reflect.TypeOf((*ReplicationVmState)(nil)).Elem()
}

type ResourceConfigSpecScaleSharesBehavior string

const (
	// Do not scale shares
	ResourceConfigSpecScaleSharesBehaviorDisabled = ResourceConfigSpecScaleSharesBehavior("disabled")
	// Scale both CPU and memory shares
	ResourceConfigSpecScaleSharesBehaviorScaleCpuAndMemoryShares = ResourceConfigSpecScaleSharesBehavior("scaleCpuAndMemoryShares")
)

func (e ResourceConfigSpecScaleSharesBehavior) Values() []ResourceConfigSpecScaleSharesBehavior {
	return []ResourceConfigSpecScaleSharesBehavior{
		ResourceConfigSpecScaleSharesBehaviorDisabled,
		ResourceConfigSpecScaleSharesBehaviorScaleCpuAndMemoryShares,
	}
}

func (e ResourceConfigSpecScaleSharesBehavior) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ResourceConfigSpecScaleSharesBehavior"] = reflect.TypeOf((*ResourceConfigSpecScaleSharesBehavior)(nil)).Elem()
}

// The policy setting used to determine when to perform scheduled
// upgrades for a virtual machine.
type ScheduledHardwareUpgradeInfoHardwareUpgradePolicy string

const (
	// No scheduled upgrades.
	ScheduledHardwareUpgradeInfoHardwareUpgradePolicyNever = ScheduledHardwareUpgradeInfoHardwareUpgradePolicy("never")
	// Run scheduled upgrades only on normal guest OS shutdown.
	ScheduledHardwareUpgradeInfoHardwareUpgradePolicyOnSoftPowerOff = ScheduledHardwareUpgradeInfoHardwareUpgradePolicy("onSoftPowerOff")
	// Always run scheduled upgrades.
	ScheduledHardwareUpgradeInfoHardwareUpgradePolicyAlways = ScheduledHardwareUpgradeInfoHardwareUpgradePolicy("always")
)

func (e ScheduledHardwareUpgradeInfoHardwareUpgradePolicy) Values() []ScheduledHardwareUpgradeInfoHardwareUpgradePolicy {
	return []ScheduledHardwareUpgradeInfoHardwareUpgradePolicy{
		ScheduledHardwareUpgradeInfoHardwareUpgradePolicyNever,
		ScheduledHardwareUpgradeInfoHardwareUpgradePolicyOnSoftPowerOff,
		ScheduledHardwareUpgradeInfoHardwareUpgradePolicyAlways,
	}
}

func (e ScheduledHardwareUpgradeInfoHardwareUpgradePolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScheduledHardwareUpgradeInfoHardwareUpgradePolicy"] = reflect.TypeOf((*ScheduledHardwareUpgradeInfoHardwareUpgradePolicy)(nil)).Elem()
}

// Status for last attempt to run scheduled hardware upgrade.
type ScheduledHardwareUpgradeInfoHardwareUpgradeStatus string

const (
	// No scheduled upgrade ever happened.
	ScheduledHardwareUpgradeInfoHardwareUpgradeStatusNone = ScheduledHardwareUpgradeInfoHardwareUpgradeStatus("none")
	// Upgrade is scheduled, but was not run yet.
	ScheduledHardwareUpgradeInfoHardwareUpgradeStatusPending = ScheduledHardwareUpgradeInfoHardwareUpgradeStatus("pending")
	// Upgrade succeeded.
	ScheduledHardwareUpgradeInfoHardwareUpgradeStatusSuccess = ScheduledHardwareUpgradeInfoHardwareUpgradeStatus("success")
	// Upgrade failed.
	//
	// # For more information about the failure
	//
	// See also `ScheduledHardwareUpgradeInfo.fault`.
	ScheduledHardwareUpgradeInfoHardwareUpgradeStatusFailed = ScheduledHardwareUpgradeInfoHardwareUpgradeStatus("failed")
)

func (e ScheduledHardwareUpgradeInfoHardwareUpgradeStatus) Values() []ScheduledHardwareUpgradeInfoHardwareUpgradeStatus {
	return []ScheduledHardwareUpgradeInfoHardwareUpgradeStatus{
		ScheduledHardwareUpgradeInfoHardwareUpgradeStatusNone,
		ScheduledHardwareUpgradeInfoHardwareUpgradeStatusPending,
		ScheduledHardwareUpgradeInfoHardwareUpgradeStatusSuccess,
		ScheduledHardwareUpgradeInfoHardwareUpgradeStatusFailed,
	}
}

func (e ScheduledHardwareUpgradeInfoHardwareUpgradeStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScheduledHardwareUpgradeInfoHardwareUpgradeStatus"] = reflect.TypeOf((*ScheduledHardwareUpgradeInfoHardwareUpgradeStatus)(nil)).Elem()
}

// The types of disk drives.
type ScsiDiskType string

const (
	// 512 native sector size drive.
	ScsiDiskTypeNative512 = ScsiDiskType("native512")
	// 4K sector size drive in 512 emulation mode.
	ScsiDiskTypeEmulated512 = ScsiDiskType("emulated512")
	// 4K native sector size drive.
	ScsiDiskTypeNative4k = ScsiDiskType("native4k")
	// Software emulated 4k.
	ScsiDiskTypeSoftwareEmulated4k = ScsiDiskType("SoftwareEmulated4k")
	// Unknown type.
	ScsiDiskTypeUnknown = ScsiDiskType("unknown")
)

func (e ScsiDiskType) Values() []ScsiDiskType {
	return []ScsiDiskType{
		ScsiDiskTypeNative512,
		ScsiDiskTypeEmulated512,
		ScsiDiskTypeNative4k,
		ScsiDiskTypeSoftwareEmulated4k,
		ScsiDiskTypeUnknown,
	}
}

func (e ScsiDiskType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScsiDiskType"] = reflect.TypeOf((*ScsiDiskType)(nil)).Elem()
}

// An indicator of the utility of Descriptor in being used as an
// identifier that is stable, unique, and correlatable.
type ScsiLunDescriptorQuality string

const (
	// The Descriptor has an identifier that is useful for identification
	// and correlation across hosts.
	ScsiLunDescriptorQualityHighQuality = ScsiLunDescriptorQuality("highQuality")
	// The Descriptor has an identifier that may be used for identification
	// and correlation across hosts.
	ScsiLunDescriptorQualityMediumQuality = ScsiLunDescriptorQuality("mediumQuality")
	// The Descriptor has an identifier that should not be used for
	// identification and correlation across hosts.
	ScsiLunDescriptorQualityLowQuality = ScsiLunDescriptorQuality("lowQuality")
	// The Descriptor has an identifier that may or may not be useful for
	// identification and correlation across hosts.
	ScsiLunDescriptorQualityUnknownQuality = ScsiLunDescriptorQuality("unknownQuality")
)

func (e ScsiLunDescriptorQuality) Values() []ScsiLunDescriptorQuality {
	return []ScsiLunDescriptorQuality{
		ScsiLunDescriptorQualityHighQuality,
		ScsiLunDescriptorQualityMediumQuality,
		ScsiLunDescriptorQualityLowQuality,
		ScsiLunDescriptorQualityUnknownQuality,
	}
}

func (e ScsiLunDescriptorQuality) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScsiLunDescriptorQuality"] = reflect.TypeOf((*ScsiLunDescriptorQuality)(nil)).Elem()
}

type ScsiLunLunReservationStatus string

const (
	ScsiLunLunReservationStatusLUN_RESERVED_UNKNOWN       = ScsiLunLunReservationStatus("LUN_RESERVED_UNKNOWN")
	ScsiLunLunReservationStatusLUN_RESERVED_YES           = ScsiLunLunReservationStatus("LUN_RESERVED_YES")
	ScsiLunLunReservationStatusLUN_RESERVED_NO            = ScsiLunLunReservationStatus("LUN_RESERVED_NO")
	ScsiLunLunReservationStatusLUN_RESERVED_NOT_SUPPORTED = ScsiLunLunReservationStatus("LUN_RESERVED_NOT_SUPPORTED")
)

func (e ScsiLunLunReservationStatus) Values() []ScsiLunLunReservationStatus {
	return []ScsiLunLunReservationStatus{
		ScsiLunLunReservationStatusLUN_RESERVED_UNKNOWN,
		ScsiLunLunReservationStatusLUN_RESERVED_YES,
		ScsiLunLunReservationStatusLUN_RESERVED_NO,
		ScsiLunLunReservationStatusLUN_RESERVED_NOT_SUPPORTED,
	}
}

func (e ScsiLunLunReservationStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScsiLunLunReservationStatus"] = reflect.TypeOf((*ScsiLunLunReservationStatus)(nil)).Elem()
	minAPIVersionForType["ScsiLunLunReservationStatus"] = "8.0.3.0"
}

// The Operational state of the LUN
type ScsiLunState string

const (
	// The LUN state is unknown.
	ScsiLunStateUnknownState = ScsiLunState("unknownState")
	// The LUN is on and available.
	ScsiLunStateOk = ScsiLunState("ok")
	// The LUN is dead and/or not reachable.
	ScsiLunStateError = ScsiLunState("error")
	// The LUN is off.
	ScsiLunStateOff = ScsiLunState("off")
	// The LUN is inactive.
	ScsiLunStateQuiesced = ScsiLunState("quiesced")
	// One or more paths to the LUN are down, but I/O
	// is still possible.
	//
	// Further path failures may
	// result in lost connectivity.
	ScsiLunStateDegraded = ScsiLunState("degraded")
	// No more paths are available to the LUN.
	ScsiLunStateLostCommunication = ScsiLunState("lostCommunication")
	// All Paths have been down for the timeout condition
	// determined by a user-configurable host advanced option.
	ScsiLunStateTimeout = ScsiLunState("timeout")
)

func (e ScsiLunState) Values() []ScsiLunState {
	return []ScsiLunState{
		ScsiLunStateUnknownState,
		ScsiLunStateOk,
		ScsiLunStateError,
		ScsiLunStateOff,
		ScsiLunStateQuiesced,
		ScsiLunStateDegraded,
		ScsiLunStateLostCommunication,
		ScsiLunStateTimeout,
	}
}

func (e ScsiLunState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScsiLunState"] = reflect.TypeOf((*ScsiLunState)(nil)).Elem()
}

// The list of SCSI device types.
//
// These values correspond to values
// published in the SCSI specification.
type ScsiLunType string

const (
	ScsiLunTypeDisk                   = ScsiLunType("disk")
	ScsiLunTypeTape                   = ScsiLunType("tape")
	ScsiLunTypePrinter                = ScsiLunType("printer")
	ScsiLunTypeProcessor              = ScsiLunType("processor")
	ScsiLunTypeWorm                   = ScsiLunType("worm")
	ScsiLunTypeCdrom                  = ScsiLunType("cdrom")
	ScsiLunTypeScanner                = ScsiLunType("scanner")
	ScsiLunTypeOpticalDevice          = ScsiLunType("opticalDevice")
	ScsiLunTypeMediaChanger           = ScsiLunType("mediaChanger")
	ScsiLunTypeCommunications         = ScsiLunType("communications")
	ScsiLunTypeStorageArrayController = ScsiLunType("storageArrayController")
	ScsiLunTypeEnclosure              = ScsiLunType("enclosure")
	ScsiLunTypeUnknown                = ScsiLunType("unknown")
)

func (e ScsiLunType) Values() []ScsiLunType {
	return []ScsiLunType{
		ScsiLunTypeDisk,
		ScsiLunTypeTape,
		ScsiLunTypePrinter,
		ScsiLunTypeProcessor,
		ScsiLunTypeWorm,
		ScsiLunTypeCdrom,
		ScsiLunTypeScanner,
		ScsiLunTypeOpticalDevice,
		ScsiLunTypeMediaChanger,
		ScsiLunTypeCommunications,
		ScsiLunTypeStorageArrayController,
		ScsiLunTypeEnclosure,
		ScsiLunTypeUnknown,
	}
}

func (e ScsiLunType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScsiLunType"] = reflect.TypeOf((*ScsiLunType)(nil)).Elem()
}

// Storage array hardware acceleration support status.
//
// When a host boots, the support status is unknown.
// As a host attempts hardware-accelerated operations,
// it determines whether the storage device supports hardware acceleration
// and sets the `ScsiLun.vStorageSupport` property accordingly.
type ScsiLunVStorageSupportStatus string

const (
	// Storage device supports hardware acceleration.
	//
	// The ESX host will use the feature to offload certain
	// storage-related operations to the device.
	ScsiLunVStorageSupportStatusVStorageSupported = ScsiLunVStorageSupportStatus("vStorageSupported")
	// Storage device does not support hardware acceleration.
	//
	// The ESX host will handle all storage-related operations.
	ScsiLunVStorageSupportStatusVStorageUnsupported = ScsiLunVStorageSupportStatus("vStorageUnsupported")
	// Initial support status value.
	ScsiLunVStorageSupportStatusVStorageUnknown = ScsiLunVStorageSupportStatus("vStorageUnknown")
)

func (e ScsiLunVStorageSupportStatus) Values() []ScsiLunVStorageSupportStatus {
	return []ScsiLunVStorageSupportStatus{
		ScsiLunVStorageSupportStatusVStorageSupported,
		ScsiLunVStorageSupportStatusVStorageUnsupported,
		ScsiLunVStorageSupportStatusVStorageUnknown,
	}
}

func (e ScsiLunVStorageSupportStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ScsiLunVStorageSupportStatus"] = reflect.TypeOf((*ScsiLunVStorageSupportStatus)(nil)).Elem()
}

type SessionManagerGenericServiceTicketTicketType string

const (
	// Ticket used for HttpNfc access to a file or disk on a datastore
	SessionManagerGenericServiceTicketTicketTypeHttpNfcServiceTicket = SessionManagerGenericServiceTicketTicketType("HttpNfcServiceTicket")
	// Ticket used for service request on a host
	SessionManagerGenericServiceTicketTicketTypeHostServiceTicket = SessionManagerGenericServiceTicketTicketType("HostServiceTicket")
	// Ticket used for service request on a VC
	SessionManagerGenericServiceTicketTicketTypeVcServiceTicket = SessionManagerGenericServiceTicketTicketType("VcServiceTicket")
)

func (e SessionManagerGenericServiceTicketTicketType) Values() []SessionManagerGenericServiceTicketTicketType {
	return []SessionManagerGenericServiceTicketTicketType{
		SessionManagerGenericServiceTicketTicketTypeHttpNfcServiceTicket,
		SessionManagerGenericServiceTicketTicketTypeHostServiceTicket,
		SessionManagerGenericServiceTicketTicketTypeVcServiceTicket,
	}
}

func (e SessionManagerGenericServiceTicketTicketType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SessionManagerGenericServiceTicketTicketType"] = reflect.TypeOf((*SessionManagerGenericServiceTicketTicketType)(nil)).Elem()
	minAPIVersionForType["SessionManagerGenericServiceTicketTicketType"] = "7.0.2.0"
}

// HTTP request methods.
type SessionManagerHttpServiceRequestSpecMethod string

const (
	SessionManagerHttpServiceRequestSpecMethodHttpOptions = SessionManagerHttpServiceRequestSpecMethod("httpOptions")
	SessionManagerHttpServiceRequestSpecMethodHttpGet     = SessionManagerHttpServiceRequestSpecMethod("httpGet")
	SessionManagerHttpServiceRequestSpecMethodHttpHead    = SessionManagerHttpServiceRequestSpecMethod("httpHead")
	SessionManagerHttpServiceRequestSpecMethodHttpPost    = SessionManagerHttpServiceRequestSpecMethod("httpPost")
	SessionManagerHttpServiceRequestSpecMethodHttpPut     = SessionManagerHttpServiceRequestSpecMethod("httpPut")
	SessionManagerHttpServiceRequestSpecMethodHttpDelete  = SessionManagerHttpServiceRequestSpecMethod("httpDelete")
	SessionManagerHttpServiceRequestSpecMethodHttpTrace   = SessionManagerHttpServiceRequestSpecMethod("httpTrace")
	SessionManagerHttpServiceRequestSpecMethodHttpConnect = SessionManagerHttpServiceRequestSpecMethod("httpConnect")
)

func (e SessionManagerHttpServiceRequestSpecMethod) Values() []SessionManagerHttpServiceRequestSpecMethod {
	return []SessionManagerHttpServiceRequestSpecMethod{
		SessionManagerHttpServiceRequestSpecMethodHttpOptions,
		SessionManagerHttpServiceRequestSpecMethodHttpGet,
		SessionManagerHttpServiceRequestSpecMethodHttpHead,
		SessionManagerHttpServiceRequestSpecMethodHttpPost,
		SessionManagerHttpServiceRequestSpecMethodHttpPut,
		SessionManagerHttpServiceRequestSpecMethodHttpDelete,
		SessionManagerHttpServiceRequestSpecMethodHttpTrace,
		SessionManagerHttpServiceRequestSpecMethodHttpConnect,
	}
}

func (e SessionManagerHttpServiceRequestSpecMethod) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SessionManagerHttpServiceRequestSpecMethod"] = reflect.TypeOf((*SessionManagerHttpServiceRequestSpecMethod)(nil)).Elem()
}

// Simplified shares notation.
//
// These designations have different meanings for different resources.
type SharesLevel string

const (
	// For CPU: Shares = 500 \* number of virtual CPUs
	// For Memory: Shares = 5 \* virtual machine memory size in megabytes
	// For Disk: Shares = 500
	// For Network: Shares = 0.25 \* `DVSFeatureCapability.networkResourcePoolHighShareValue`
	SharesLevelLow = SharesLevel("low")
	// For CPU: Shares = 1000 \* number of virtual CPUs
	// For Memory: Shares = 10 \* virtual machine memory size in megabytes
	// For Disk: Shares = 1000
	// For Network: Shares = 0.5 \* `DVSFeatureCapability.networkResourcePoolHighShareValue`
	SharesLevelNormal = SharesLevel("normal")
	// For CPU: Shares = 2000 \* number of virtual CPUs
	// For Memory: Shares = 20 \* virtual machine memory size in megabytes
	// For Disk: Shares = 2000
	// For Network: Shares = `DVSFeatureCapability.networkResourcePoolHighShareValue`
	SharesLevelHigh = SharesLevel("high")
	// If you specify <code>custom</code> for the `SharesInfo.level` property, when there is resource contention the Server uses the `SharesInfo.shares` value to determine resource allocation.
	SharesLevelCustom = SharesLevel("custom")
)

func (e SharesLevel) Values() []SharesLevel {
	return []SharesLevel{
		SharesLevelLow,
		SharesLevelNormal,
		SharesLevelHigh,
		SharesLevelCustom,
	}
}

func (e SharesLevel) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SharesLevel"] = reflect.TypeOf((*SharesLevel)(nil)).Elem()
}

// The encoding of the resultant return data.
//
// This is a hint to the client side
// to indicate the format of the information being returned.
type SimpleCommandEncoding string

const (
	// Comma separated values
	SimpleCommandEncodingCSV = SimpleCommandEncoding("CSV")
	// Hex encoded binary data
	SimpleCommandEncodingHEX    = SimpleCommandEncoding("HEX")
	SimpleCommandEncodingSTRING = SimpleCommandEncoding("STRING")
)

func (e SimpleCommandEncoding) Values() []SimpleCommandEncoding {
	return []SimpleCommandEncoding{
		SimpleCommandEncodingCSV,
		SimpleCommandEncodingHEX,
		SimpleCommandEncodingSTRING,
	}
}

func (e SimpleCommandEncoding) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SimpleCommandEncoding"] = reflect.TypeOf((*SimpleCommandEncoding)(nil)).Elem()
}

// The available SLP discovery methods.
type SlpDiscoveryMethod string

const (
	// Use DHCP to find the SLP DAs.
	SlpDiscoveryMethodSlpDhcp = SlpDiscoveryMethod("slpDhcp")
	// Use broadcasting to find SLP DAs.
	//
	// Only DAs on the current subnet will be found.
	SlpDiscoveryMethodSlpAutoUnicast = SlpDiscoveryMethod("slpAutoUnicast")
	// Use the well known multicast address to find DAs.
	SlpDiscoveryMethodSlpAutoMulticast = SlpDiscoveryMethod("slpAutoMulticast")
	// User specified address for a DA.
	SlpDiscoveryMethodSlpManual = SlpDiscoveryMethod("slpManual")
)

func (e SlpDiscoveryMethod) Values() []SlpDiscoveryMethod {
	return []SlpDiscoveryMethod{
		SlpDiscoveryMethodSlpDhcp,
		SlpDiscoveryMethodSlpAutoUnicast,
		SlpDiscoveryMethodSlpAutoMulticast,
		SlpDiscoveryMethodSlpManual,
	}
}

func (e SlpDiscoveryMethod) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SlpDiscoveryMethod"] = reflect.TypeOf((*SlpDiscoveryMethod)(nil)).Elem()
}

// These are the constraint relationships between software packages.
type SoftwarePackageConstraint string

const (
	SoftwarePackageConstraintEquals           = SoftwarePackageConstraint("equals")
	SoftwarePackageConstraintLessThan         = SoftwarePackageConstraint("lessThan")
	SoftwarePackageConstraintLessThanEqual    = SoftwarePackageConstraint("lessThanEqual")
	SoftwarePackageConstraintGreaterThanEqual = SoftwarePackageConstraint("greaterThanEqual")
	SoftwarePackageConstraintGreaterThan      = SoftwarePackageConstraint("greaterThan")
)

func (e SoftwarePackageConstraint) Values() []SoftwarePackageConstraint {
	return []SoftwarePackageConstraint{
		SoftwarePackageConstraintEquals,
		SoftwarePackageConstraintLessThan,
		SoftwarePackageConstraintLessThanEqual,
		SoftwarePackageConstraintGreaterThanEqual,
		SoftwarePackageConstraintGreaterThan,
	}
}

func (e SoftwarePackageConstraint) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SoftwarePackageConstraint"] = reflect.TypeOf((*SoftwarePackageConstraint)(nil)).Elem()
}

type SoftwarePackageVibType string

const (
	// This package is installed into bootbank in storage.
	SoftwarePackageVibTypeBootbank = SoftwarePackageVibType("bootbank")
	// This package is installed into tools partition in storage.
	SoftwarePackageVibTypeTools = SoftwarePackageVibType("tools")
	// This package contains install related data without
	// content to install.
	SoftwarePackageVibTypeMeta = SoftwarePackageVibType("meta")
)

func (e SoftwarePackageVibType) Values() []SoftwarePackageVibType {
	return []SoftwarePackageVibType{
		SoftwarePackageVibTypeBootbank,
		SoftwarePackageVibTypeTools,
		SoftwarePackageVibTypeMeta,
	}
}

func (e SoftwarePackageVibType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["SoftwarePackageVibType"] = reflect.TypeOf((*SoftwarePackageVibType)(nil)).Elem()
}

// The operation on the target state.
type StateAlarmOperator string

const (
	// Test if the target state matches the given red or yellow states.
	StateAlarmOperatorIsEqual = StateAlarmOperator("isEqual")
	// Test if the target state does not match the given red or yellow states.
	StateAlarmOperatorIsUnequal = StateAlarmOperator("isUnequal")
)

func (e StateAlarmOperator) Values() []StateAlarmOperator {
	return []StateAlarmOperator{
		StateAlarmOperatorIsEqual,
		StateAlarmOperatorIsUnequal,
	}
}

func (e StateAlarmOperator) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["StateAlarmOperator"] = reflect.TypeOf((*StateAlarmOperator)(nil)).Elem()
}

// Storage DRS behavior.
type StorageDrsPodConfigInfoBehavior string

const (
	// Specifies that VirtualCenter should generate recommendations for
	// virtual disk migration and for placement with a datastore,
	// but should not execute the recommendations automatically.
	StorageDrsPodConfigInfoBehaviorManual = StorageDrsPodConfigInfoBehavior("manual")
	// Specifies that VirtualCenter should generate recommendations
	// for virtual disk migration and for placement with a
	// datastore.
	//
	// The recommendations for virtual disk migrations
	// will be executed automatically, but the placement
	// recommendations will be done manually.
	StorageDrsPodConfigInfoBehaviorAutomated = StorageDrsPodConfigInfoBehavior("automated")
)

func (e StorageDrsPodConfigInfoBehavior) Values() []StorageDrsPodConfigInfoBehavior {
	return []StorageDrsPodConfigInfoBehavior{
		StorageDrsPodConfigInfoBehaviorManual,
		StorageDrsPodConfigInfoBehaviorAutomated,
	}
}

func (e StorageDrsPodConfigInfoBehavior) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["StorageDrsPodConfigInfoBehavior"] = reflect.TypeOf((*StorageDrsPodConfigInfoBehavior)(nil)).Elem()
}

// Defines the two ways a space utilization threshold can be specified.
type StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode string

const (
	// Default mode: threshold as a percentage of datastore capacity
	StorageDrsSpaceLoadBalanceConfigSpaceThresholdModeUtilization = StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode("utilization")
	// Threshold as an absolute value of free space in GBs
	StorageDrsSpaceLoadBalanceConfigSpaceThresholdModeFreeSpace = StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode("freeSpace")
)

func (e StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode) Values() []StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode {
	return []StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode{
		StorageDrsSpaceLoadBalanceConfigSpaceThresholdModeUtilization,
		StorageDrsSpaceLoadBalanceConfigSpaceThresholdModeFreeSpace,
	}
}

func (e StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode"] = reflect.TypeOf((*StorageDrsSpaceLoadBalanceConfigSpaceThresholdMode)(nil)).Elem()
}

// Deprecated as of vSphere8.0 U3, and there is no replacement for it.
//
// # User specification of congestion threshold mode on a given datastore
//
// For more information, see
// `StorageIORMInfo.congestionThreshold`
type StorageIORMThresholdMode string

const (
	// Storagage IO Control will choose appropriate congestion threshold value
	// for that datastore to operate at given percentage of peak throughput.
	//
	// This is the default setting
	StorageIORMThresholdModeAutomatic = StorageIORMThresholdMode("automatic")
	// Use user specified Storage IO Control congestion threshold value
	StorageIORMThresholdModeManual = StorageIORMThresholdMode("manual")
)

func (e StorageIORMThresholdMode) Values() []StorageIORMThresholdMode {
	return []StorageIORMThresholdMode{
		StorageIORMThresholdModeAutomatic,
		StorageIORMThresholdModeManual,
	}
}

func (e StorageIORMThresholdMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["StorageIORMThresholdMode"] = reflect.TypeOf((*StorageIORMThresholdMode)(nil)).Elem()
}

// Defines the storage placement operation type.
type StoragePlacementSpecPlacementType string

const (
	// Create a VM.
	StoragePlacementSpecPlacementTypeCreate = StoragePlacementSpecPlacementType("create")
	// Reconfigure a VM.
	StoragePlacementSpecPlacementTypeReconfigure = StoragePlacementSpecPlacementType("reconfigure")
	// Relocate a VM.
	StoragePlacementSpecPlacementTypeRelocate = StoragePlacementSpecPlacementType("relocate")
	// Clone a VM.
	StoragePlacementSpecPlacementTypeClone = StoragePlacementSpecPlacementType("clone")
)

func (e StoragePlacementSpecPlacementType) Values() []StoragePlacementSpecPlacementType {
	return []StoragePlacementSpecPlacementType{
		StoragePlacementSpecPlacementTypeCreate,
		StoragePlacementSpecPlacementTypeReconfigure,
		StoragePlacementSpecPlacementTypeRelocate,
		StoragePlacementSpecPlacementTypeClone,
	}
}

func (e StoragePlacementSpecPlacementType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["StoragePlacementSpecPlacementType"] = reflect.TypeOf((*StoragePlacementSpecPlacementType)(nil)).Elem()
}

// This option specifies how to select tasks based on child relationships
// in the inventory hierarchy.
//
// If a managed entity has children, their tasks
// can be retrieved with this filter option.
type TaskFilterSpecRecursionOption string

const (
	// Returns tasks that pertain only to the specified managed entity,
	// and not its children.
	TaskFilterSpecRecursionOptionSelf = TaskFilterSpecRecursionOption("self")
	// Returns tasks pertaining to child entities only.
	//
	// Excludes
	// tasks pertaining to the specified managed entity itself.
	TaskFilterSpecRecursionOptionChildren = TaskFilterSpecRecursionOption("children")
	// Returns tasks pertaining either to the specified managed entity
	// or to its child entities.
	TaskFilterSpecRecursionOptionAll = TaskFilterSpecRecursionOption("all")
)

func (e TaskFilterSpecRecursionOption) Values() []TaskFilterSpecRecursionOption {
	return []TaskFilterSpecRecursionOption{
		TaskFilterSpecRecursionOptionSelf,
		TaskFilterSpecRecursionOptionChildren,
		TaskFilterSpecRecursionOptionAll,
	}
}

func (e TaskFilterSpecRecursionOption) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["TaskFilterSpecRecursionOption"] = reflect.TypeOf((*TaskFilterSpecRecursionOption)(nil)).Elem()
}

// This option specifies a time stamp governing the selection of tasks.
type TaskFilterSpecTimeOption string

const (
	// The time stamp when the task was created and queued.
	TaskFilterSpecTimeOptionQueuedTime = TaskFilterSpecTimeOption("queuedTime")
	// The time stamp when the task started.
	TaskFilterSpecTimeOptionStartedTime = TaskFilterSpecTimeOption("startedTime")
	// The time stamp when the task finished.
	TaskFilterSpecTimeOptionCompletedTime = TaskFilterSpecTimeOption("completedTime")
)

func (e TaskFilterSpecTimeOption) Values() []TaskFilterSpecTimeOption {
	return []TaskFilterSpecTimeOption{
		TaskFilterSpecTimeOptionQueuedTime,
		TaskFilterSpecTimeOptionStartedTime,
		TaskFilterSpecTimeOptionCompletedTime,
	}
}

func (e TaskFilterSpecTimeOption) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["TaskFilterSpecTimeOption"] = reflect.TypeOf((*TaskFilterSpecTimeOption)(nil)).Elem()
}

// List of possible states of a task.
type TaskInfoState string

const (
	// When there are too many tasks for threads to handle.
	TaskInfoStateQueued = TaskInfoState("queued")
	// When the busy thread is freed from its current task by
	// finishing the task, it picks a queued task to run.
	//
	// Then the queued tasks are marked as running.
	TaskInfoStateRunning = TaskInfoState("running")
	// When a running task has completed.
	TaskInfoStateSuccess = TaskInfoState("success")
	// When a running task has encountered an error.
	TaskInfoStateError = TaskInfoState("error")
)

func (e TaskInfoState) Values() []TaskInfoState {
	return []TaskInfoState{
		TaskInfoStateQueued,
		TaskInfoStateRunning,
		TaskInfoStateSuccess,
		TaskInfoStateError,
	}
}

func (e TaskInfoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["TaskInfoState"] = reflect.TypeOf((*TaskInfoState)(nil)).Elem()
}

type ThirdPartyLicenseAssignmentFailedReason string

const (
	// A general failure has occurred during assigning license to the 3rd party module
	ThirdPartyLicenseAssignmentFailedReasonLicenseAssignmentFailed = ThirdPartyLicenseAssignmentFailedReason("licenseAssignmentFailed")
	// The 3rd party module we are trying to license is not installed.
	ThirdPartyLicenseAssignmentFailedReasonModuleNotInstalled = ThirdPartyLicenseAssignmentFailedReason("moduleNotInstalled")
)

func (e ThirdPartyLicenseAssignmentFailedReason) Values() []ThirdPartyLicenseAssignmentFailedReason {
	return []ThirdPartyLicenseAssignmentFailedReason{
		ThirdPartyLicenseAssignmentFailedReasonLicenseAssignmentFailed,
		ThirdPartyLicenseAssignmentFailedReasonModuleNotInstalled,
	}
}

func (e ThirdPartyLicenseAssignmentFailedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ThirdPartyLicenseAssignmentFailedReason"] = reflect.TypeOf((*ThirdPartyLicenseAssignmentFailedReason)(nil)).Elem()
}

// The policy setting used to determine when tools are auto-upgraded for
// a virtual machine
type UpgradePolicy string

const (
	// No auto-upgrades for tools will be performed for this
	// virtual machine.
	//
	// Users must manually invoke the UpgradeTools
	// operation to update the tools.
	UpgradePolicyManual = UpgradePolicy("manual")
	// When the virtual machine is power-cycled, the system checks
	// for a newer version of tools when the VM comes back up.
	//
	// If it
	// is available, a tools upgrade is automatically performed on the
	// virtual machine and it is rebooted if necessary.
	UpgradePolicyUpgradeAtPowerCycle = UpgradePolicy("upgradeAtPowerCycle")
)

func (e UpgradePolicy) Values() []UpgradePolicy {
	return []UpgradePolicy{
		UpgradePolicyManual,
		UpgradePolicyUpgradeAtPowerCycle,
	}
}

func (e UpgradePolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["UpgradePolicy"] = reflect.TypeOf((*UpgradePolicy)(nil)).Elem()
}

type VAppAutoStartAction string

const (
	// No action is taken for this virtual machine.
	//
	// This virtual machine is
	// not a part of the auto-start sequence. This can be used for both auto-start
	// and auto-start settings.
	VAppAutoStartActionNone = VAppAutoStartAction("none")
	// This virtual machine is powered on when it is next in the auto-start order.
	VAppAutoStartActionPowerOn = VAppAutoStartAction("powerOn")
	// This virtual machine is powered off when it is next in the auto-stop order.
	//
	// This is the default stopAction.
	VAppAutoStartActionPowerOff = VAppAutoStartAction("powerOff")
	// The guest operating system for a virtual machine is shut down when that
	// virtual machine in next in the auto-stop order.
	VAppAutoStartActionGuestShutdown = VAppAutoStartAction("guestShutdown")
	// This virtual machine is suspended when it is next in the auto-stop order.
	VAppAutoStartActionSuspend = VAppAutoStartAction("suspend")
)

func (e VAppAutoStartAction) Values() []VAppAutoStartAction {
	return []VAppAutoStartAction{
		VAppAutoStartActionNone,
		VAppAutoStartActionPowerOn,
		VAppAutoStartActionPowerOff,
		VAppAutoStartActionGuestShutdown,
		VAppAutoStartActionSuspend,
	}
}

func (e VAppAutoStartAction) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VAppAutoStartAction"] = reflect.TypeOf((*VAppAutoStartAction)(nil)).Elem()
}

// The cloned VMs can either be provisioned the same way as the VMs
// they are a clone of, thin provisioned or thick provisioned, or
// linked clones (i.e., using delta disks).
type VAppCloneSpecProvisioningType string

const (
	// Each disk in the cloned virtual machines will have the same
	// type of disk as the source vApp.
	VAppCloneSpecProvisioningTypeSameAsSource = VAppCloneSpecProvisioningType("sameAsSource")
	// Each disk in the cloned virtual machines is allocated in full
	// size now and committed on demand.
	//
	// This is only supported on
	// VMFS-3 and newer datastores. Other types of datastores may
	// create thick disks.
	VAppCloneSpecProvisioningTypeThin = VAppCloneSpecProvisioningType("thin")
	// Each disk in the cloned virtual machines are allocated and
	// committed in full size immediately.
	VAppCloneSpecProvisioningTypeThick = VAppCloneSpecProvisioningType("thick")
)

func (e VAppCloneSpecProvisioningType) Values() []VAppCloneSpecProvisioningType {
	return []VAppCloneSpecProvisioningType{
		VAppCloneSpecProvisioningTypeSameAsSource,
		VAppCloneSpecProvisioningTypeThin,
		VAppCloneSpecProvisioningTypeThick,
	}
}

func (e VAppCloneSpecProvisioningType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VAppCloneSpecProvisioningType"] = reflect.TypeOf((*VAppCloneSpecProvisioningType)(nil)).Elem()
}

// IP allocation schemes supported by the guest.
type VAppIPAssignmentInfoAllocationSchemes string

const (
	// The vApp supports DHCP to acquire IP configuration.
	VAppIPAssignmentInfoAllocationSchemesDhcp = VAppIPAssignmentInfoAllocationSchemes("dhcp")
	// The vApp supports setting the IP configuration through the
	// properties provided in the OVF environment.
	VAppIPAssignmentInfoAllocationSchemesOvfenv = VAppIPAssignmentInfoAllocationSchemes("ovfenv")
)

func (e VAppIPAssignmentInfoAllocationSchemes) Values() []VAppIPAssignmentInfoAllocationSchemes {
	return []VAppIPAssignmentInfoAllocationSchemes{
		VAppIPAssignmentInfoAllocationSchemesDhcp,
		VAppIPAssignmentInfoAllocationSchemesOvfenv,
	}
}

func (e VAppIPAssignmentInfoAllocationSchemes) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VAppIPAssignmentInfoAllocationSchemes"] = reflect.TypeOf((*VAppIPAssignmentInfoAllocationSchemes)(nil)).Elem()
}

// IP allocation policy for a deployment.
type VAppIPAssignmentInfoIpAllocationPolicy string

const (
	// Specifies that DHCP must be used to allocate IP addresses to the vApp
	VAppIPAssignmentInfoIpAllocationPolicyDhcpPolicy = VAppIPAssignmentInfoIpAllocationPolicy("dhcpPolicy")
	// Specifies that IP allocation is done through the range managed by the
	// vSphere platform.
	//
	// The IP addresses are allocated when needed, typically at
	// power-on, and deallocated during power-off. There is no guarantee that a
	// vApp will get the same IP address when restarted.
	VAppIPAssignmentInfoIpAllocationPolicyTransientPolicy = VAppIPAssignmentInfoIpAllocationPolicy("transientPolicy")
	// Specifies that IP addresses are configured manually when the vApp is deployed
	// and will be kept until reconfigured or the vApp destroyed.
	//
	// This will ensure
	// that a vApp gets a consistent IP for its life-time.
	VAppIPAssignmentInfoIpAllocationPolicyFixedPolicy = VAppIPAssignmentInfoIpAllocationPolicy("fixedPolicy")
	// Specifies that IP allocation is done through the range managed by the VI
	// platform.
	//
	// The IP addresses are allocated at first power-on, and remain
	// allocated at power-off. This will ensure that a vApp gets a consistent
	// IP for its life-time.
	VAppIPAssignmentInfoIpAllocationPolicyFixedAllocatedPolicy = VAppIPAssignmentInfoIpAllocationPolicy("fixedAllocatedPolicy")
)

func (e VAppIPAssignmentInfoIpAllocationPolicy) Values() []VAppIPAssignmentInfoIpAllocationPolicy {
	return []VAppIPAssignmentInfoIpAllocationPolicy{
		VAppIPAssignmentInfoIpAllocationPolicyDhcpPolicy,
		VAppIPAssignmentInfoIpAllocationPolicyTransientPolicy,
		VAppIPAssignmentInfoIpAllocationPolicyFixedPolicy,
		VAppIPAssignmentInfoIpAllocationPolicyFixedAllocatedPolicy,
	}
}

func (e VAppIPAssignmentInfoIpAllocationPolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VAppIPAssignmentInfoIpAllocationPolicy"] = reflect.TypeOf((*VAppIPAssignmentInfoIpAllocationPolicy)(nil)).Elem()
}

// IP protocols supported by the guest.
type VAppIPAssignmentInfoProtocols string

const (
	// The vApp supports IPv4 protocol.
	VAppIPAssignmentInfoProtocolsIPv4 = VAppIPAssignmentInfoProtocols("IPv4")
	// The vApp supports IPv6 protocol.
	VAppIPAssignmentInfoProtocolsIPv6 = VAppIPAssignmentInfoProtocols("IPv6")
)

func (e VAppIPAssignmentInfoProtocols) Values() []VAppIPAssignmentInfoProtocols {
	return []VAppIPAssignmentInfoProtocols{
		VAppIPAssignmentInfoProtocolsIPv4,
		VAppIPAssignmentInfoProtocolsIPv6,
	}
}

func (e VAppIPAssignmentInfoProtocols) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VAppIPAssignmentInfoProtocols"] = reflect.TypeOf((*VAppIPAssignmentInfoProtocols)(nil)).Elem()
}

type VFlashModuleNotSupportedReason string

const (
	VFlashModuleNotSupportedReasonCacheModeNotSupported            = VFlashModuleNotSupportedReason("CacheModeNotSupported")
	VFlashModuleNotSupportedReasonCacheConsistencyTypeNotSupported = VFlashModuleNotSupportedReason("CacheConsistencyTypeNotSupported")
	VFlashModuleNotSupportedReasonCacheBlockSizeNotSupported       = VFlashModuleNotSupportedReason("CacheBlockSizeNotSupported")
	VFlashModuleNotSupportedReasonCacheReservationNotSupported     = VFlashModuleNotSupportedReason("CacheReservationNotSupported")
	VFlashModuleNotSupportedReasonDiskSizeNotSupported             = VFlashModuleNotSupportedReason("DiskSizeNotSupported")
)

func (e VFlashModuleNotSupportedReason) Values() []VFlashModuleNotSupportedReason {
	return []VFlashModuleNotSupportedReason{
		VFlashModuleNotSupportedReasonCacheModeNotSupported,
		VFlashModuleNotSupportedReasonCacheConsistencyTypeNotSupported,
		VFlashModuleNotSupportedReasonCacheBlockSizeNotSupported,
		VFlashModuleNotSupportedReasonCacheReservationNotSupported,
		VFlashModuleNotSupportedReasonDiskSizeNotSupported,
	}
}

func (e VFlashModuleNotSupportedReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VFlashModuleNotSupportedReason"] = reflect.TypeOf((*VFlashModuleNotSupportedReason)(nil)).Elem()
}

// Types of a host's compatibility with a designated virtual machine
// that is a candidate for VMotion.
//
// Used with queryVMotionCompatibility
// both as inputs (to designate which compatibility types to test for)
// and as outputs (to specify which compatibility types apply for
// each host).
type VMotionCompatibilityType string

const (
	// The host's CPU features are compatible with the
	// the virtual machine's requirements.
	VMotionCompatibilityTypeCpu = VMotionCompatibilityType("cpu")
	// The software platform on the host supports VMotion
	// and is compatible with the virtual machine.
	VMotionCompatibilityTypeSoftware = VMotionCompatibilityType("software")
)

func (e VMotionCompatibilityType) Values() []VMotionCompatibilityType {
	return []VMotionCompatibilityType{
		VMotionCompatibilityTypeCpu,
		VMotionCompatibilityTypeSoftware,
	}
}

func (e VMotionCompatibilityType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMotionCompatibilityType"] = reflect.TypeOf((*VMotionCompatibilityType)(nil)).Elem()
}

// The teaming health check match status.
type VMwareDVSTeamingMatchStatus string

const (
	// The value of 'loadbalance\_ip' is used in a uplink teaming policy
	// `VmwareUplinkPortTeamingPolicy.policy`
	// in the vSphere Distributed Switch, and the external physical switch
	// has the matching EtherChannel configuration.
	VMwareDVSTeamingMatchStatusIphashMatch = VMwareDVSTeamingMatchStatus("iphashMatch")
	// The value of 'loadbalance\_ip' is not used in a uplink teaming policy
	// `VmwareUplinkPortTeamingPolicy.policy`
	// in the vSphere Distributed Switch, and the external physical switch
	// does not have EtherChannel configuration.
	VMwareDVSTeamingMatchStatusNonIphashMatch = VMwareDVSTeamingMatchStatus("nonIphashMatch")
	// The value of 'loadbalance\_ip' is used in a uplink teaming policy
	// `VmwareUplinkPortTeamingPolicy.policy`
	// in the vSphere Distributed Switch, but the external physical switch
	// does not have the matching EtherChannel configuration.
	VMwareDVSTeamingMatchStatusIphashMismatch = VMwareDVSTeamingMatchStatus("iphashMismatch")
	// The value of 'loadbalance\_ip' is not used in a uplink teaming policy
	// `VmwareUplinkPortTeamingPolicy.policy`
	// in the vSphere Distributed Switch, but the external physical switch
	// has EtherChannel configuration.
	VMwareDVSTeamingMatchStatusNonIphashMismatch = VMwareDVSTeamingMatchStatus("nonIphashMismatch")
)

func (e VMwareDVSTeamingMatchStatus) Values() []VMwareDVSTeamingMatchStatus {
	return []VMwareDVSTeamingMatchStatus{
		VMwareDVSTeamingMatchStatusIphashMatch,
		VMwareDVSTeamingMatchStatusNonIphashMatch,
		VMwareDVSTeamingMatchStatusIphashMismatch,
		VMwareDVSTeamingMatchStatusNonIphashMismatch,
	}
}

func (e VMwareDVSTeamingMatchStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareDVSTeamingMatchStatus"] = reflect.TypeOf((*VMwareDVSTeamingMatchStatus)(nil)).Elem()
}

// Distributed Port Mirroring session Encapsulation types.
type VMwareDVSVspanSessionEncapType string

const (
	// Encapsulate original packets with GRE protocol
	VMwareDVSVspanSessionEncapTypeGre = VMwareDVSVspanSessionEncapType("gre")
	// Encapsulate original packets with ERSPAN Type2 protocol
	VMwareDVSVspanSessionEncapTypeErspan2 = VMwareDVSVspanSessionEncapType("erspan2")
	// Encapsulate original packets with ERSPAN Type3 protocol
	VMwareDVSVspanSessionEncapTypeErspan3 = VMwareDVSVspanSessionEncapType("erspan3")
)

func (e VMwareDVSVspanSessionEncapType) Values() []VMwareDVSVspanSessionEncapType {
	return []VMwareDVSVspanSessionEncapType{
		VMwareDVSVspanSessionEncapTypeGre,
		VMwareDVSVspanSessionEncapTypeErspan2,
		VMwareDVSVspanSessionEncapTypeErspan3,
	}
}

func (e VMwareDVSVspanSessionEncapType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareDVSVspanSessionEncapType"] = reflect.TypeOf((*VMwareDVSVspanSessionEncapType)(nil)).Elem()
}

// Distributed Port Mirroring session types.
type VMwareDVSVspanSessionType string

const (
	// Deprecated as of vSphere API 5.1.
	//
	// In mixedDestMirror session, Distributed Ports can be used as source entities,
	// and both Distributed Ports and Uplink Ports Name can be used as destination entities.
	VMwareDVSVspanSessionTypeMixedDestMirror = VMwareDVSVspanSessionType("mixedDestMirror")
	// In dvPortMirror session, Distributed Ports can be used as both source
	// and destination entities.
	VMwareDVSVspanSessionTypeDvPortMirror = VMwareDVSVspanSessionType("dvPortMirror")
	// In remoteMirrorSource session, Distributed Ports can be used as source entities,
	// and uplink ports name can be used as destination entities.
	VMwareDVSVspanSessionTypeRemoteMirrorSource = VMwareDVSVspanSessionType("remoteMirrorSource")
	// In remoteMirrorDest session, vlan Ids can be used as source entities,
	// and Distributed Ports can be used as destination entities.
	VMwareDVSVspanSessionTypeRemoteMirrorDest = VMwareDVSVspanSessionType("remoteMirrorDest")
	// In encapsulatedRemoteMirrorSource session, Distributed Ports can be used as source entities,
	// and Ip address can be used as destination entities.
	VMwareDVSVspanSessionTypeEncapsulatedRemoteMirrorSource = VMwareDVSVspanSessionType("encapsulatedRemoteMirrorSource")
)

func (e VMwareDVSVspanSessionType) Values() []VMwareDVSVspanSessionType {
	return []VMwareDVSVspanSessionType{
		VMwareDVSVspanSessionTypeMixedDestMirror,
		VMwareDVSVspanSessionTypeDvPortMirror,
		VMwareDVSVspanSessionTypeRemoteMirrorSource,
		VMwareDVSVspanSessionTypeRemoteMirrorDest,
		VMwareDVSVspanSessionTypeEncapsulatedRemoteMirrorSource,
	}
}

func (e VMwareDVSVspanSessionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareDVSVspanSessionType"] = reflect.TypeOf((*VMwareDVSVspanSessionType)(nil)).Elem()
}

// Link Aggregation Control Protocol API versions.
type VMwareDvsLacpApiVersion string

const (
	// Deprecated as of vSphere API 7.0u1.
	//
	// One Link Aggregation Control Protocol group in the switch
	VMwareDvsLacpApiVersionSingleLag = VMwareDvsLacpApiVersion("singleLag")
	// Multiple Link Aggregation Control Protocol in the switch.
	VMwareDvsLacpApiVersionMultipleLag = VMwareDvsLacpApiVersion("multipleLag")
)

func (e VMwareDvsLacpApiVersion) Values() []VMwareDvsLacpApiVersion {
	return []VMwareDvsLacpApiVersion{
		VMwareDvsLacpApiVersionSingleLag,
		VMwareDvsLacpApiVersionMultipleLag,
	}
}

func (e VMwareDvsLacpApiVersion) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareDvsLacpApiVersion"] = reflect.TypeOf((*VMwareDvsLacpApiVersion)(nil)).Elem()
}

// Load balance algorithm in a Link Aggregation Control Protocol group.
type VMwareDvsLacpLoadBalanceAlgorithm string

const (
	// Source MAC address
	VMwareDvsLacpLoadBalanceAlgorithmSrcMac = VMwareDvsLacpLoadBalanceAlgorithm("srcMac")
	// Destination MAC address
	VMwareDvsLacpLoadBalanceAlgorithmDestMac = VMwareDvsLacpLoadBalanceAlgorithm("destMac")
	// Source and destination MAC address
	VMwareDvsLacpLoadBalanceAlgorithmSrcDestMac = VMwareDvsLacpLoadBalanceAlgorithm("srcDestMac")
	// Destination IP and VLAN
	VMwareDvsLacpLoadBalanceAlgorithmDestIpVlan = VMwareDvsLacpLoadBalanceAlgorithm("destIpVlan")
	// Source IP and VLAN
	VMwareDvsLacpLoadBalanceAlgorithmSrcIpVlan = VMwareDvsLacpLoadBalanceAlgorithm("srcIpVlan")
	// Source and destination IP and VLAN
	VMwareDvsLacpLoadBalanceAlgorithmSrcDestIpVlan = VMwareDvsLacpLoadBalanceAlgorithm("srcDestIpVlan")
	// Destination TCP/UDP port number
	VMwareDvsLacpLoadBalanceAlgorithmDestTcpUdpPort = VMwareDvsLacpLoadBalanceAlgorithm("destTcpUdpPort")
	// Source TCP/UDP port number
	VMwareDvsLacpLoadBalanceAlgorithmSrcTcpUdpPort = VMwareDvsLacpLoadBalanceAlgorithm("srcTcpUdpPort")
	// Source and destination TCP/UDP port number
	VMwareDvsLacpLoadBalanceAlgorithmSrcDestTcpUdpPort = VMwareDvsLacpLoadBalanceAlgorithm("srcDestTcpUdpPort")
	// Destination IP and TCP/UDP port number
	VMwareDvsLacpLoadBalanceAlgorithmDestIpTcpUdpPort = VMwareDvsLacpLoadBalanceAlgorithm("destIpTcpUdpPort")
	// Source IP and TCP/UDP port number
	VMwareDvsLacpLoadBalanceAlgorithmSrcIpTcpUdpPort = VMwareDvsLacpLoadBalanceAlgorithm("srcIpTcpUdpPort")
	// Source and destination IP and TCP/UDP port number
	VMwareDvsLacpLoadBalanceAlgorithmSrcDestIpTcpUdpPort = VMwareDvsLacpLoadBalanceAlgorithm("srcDestIpTcpUdpPort")
	// Destination IP, TCP/UDP port number and VLAN
	VMwareDvsLacpLoadBalanceAlgorithmDestIpTcpUdpPortVlan = VMwareDvsLacpLoadBalanceAlgorithm("destIpTcpUdpPortVlan")
	// Source IP, TCP/UDP port number and VLAN
	VMwareDvsLacpLoadBalanceAlgorithmSrcIpTcpUdpPortVlan = VMwareDvsLacpLoadBalanceAlgorithm("srcIpTcpUdpPortVlan")
	// Source and destination IP,
	// source and destination TCP/UDP port number and VLAN.
	VMwareDvsLacpLoadBalanceAlgorithmSrcDestIpTcpUdpPortVlan = VMwareDvsLacpLoadBalanceAlgorithm("srcDestIpTcpUdpPortVlan")
	// Destination IP
	VMwareDvsLacpLoadBalanceAlgorithmDestIp = VMwareDvsLacpLoadBalanceAlgorithm("destIp")
	// Source IP
	VMwareDvsLacpLoadBalanceAlgorithmSrcIp = VMwareDvsLacpLoadBalanceAlgorithm("srcIp")
	// Source and Destination IP
	VMwareDvsLacpLoadBalanceAlgorithmSrcDestIp = VMwareDvsLacpLoadBalanceAlgorithm("srcDestIp")
	// VLAN only
	VMwareDvsLacpLoadBalanceAlgorithmVlan = VMwareDvsLacpLoadBalanceAlgorithm("vlan")
	// Source Virtual Port Id
	VMwareDvsLacpLoadBalanceAlgorithmSrcPortId = VMwareDvsLacpLoadBalanceAlgorithm("srcPortId")
)

func (e VMwareDvsLacpLoadBalanceAlgorithm) Values() []VMwareDvsLacpLoadBalanceAlgorithm {
	return []VMwareDvsLacpLoadBalanceAlgorithm{
		VMwareDvsLacpLoadBalanceAlgorithmSrcMac,
		VMwareDvsLacpLoadBalanceAlgorithmDestMac,
		VMwareDvsLacpLoadBalanceAlgorithmSrcDestMac,
		VMwareDvsLacpLoadBalanceAlgorithmDestIpVlan,
		VMwareDvsLacpLoadBalanceAlgorithmSrcIpVlan,
		VMwareDvsLacpLoadBalanceAlgorithmSrcDestIpVlan,
		VMwareDvsLacpLoadBalanceAlgorithmDestTcpUdpPort,
		VMwareDvsLacpLoadBalanceAlgorithmSrcTcpUdpPort,
		VMwareDvsLacpLoadBalanceAlgorithmSrcDestTcpUdpPort,
		VMwareDvsLacpLoadBalanceAlgorithmDestIpTcpUdpPort,
		VMwareDvsLacpLoadBalanceAlgorithmSrcIpTcpUdpPort,
		VMwareDvsLacpLoadBalanceAlgorithmSrcDestIpTcpUdpPort,
		VMwareDvsLacpLoadBalanceAlgorithmDestIpTcpUdpPortVlan,
		VMwareDvsLacpLoadBalanceAlgorithmSrcIpTcpUdpPortVlan,
		VMwareDvsLacpLoadBalanceAlgorithmSrcDestIpTcpUdpPortVlan,
		VMwareDvsLacpLoadBalanceAlgorithmDestIp,
		VMwareDvsLacpLoadBalanceAlgorithmSrcIp,
		VMwareDvsLacpLoadBalanceAlgorithmSrcDestIp,
		VMwareDvsLacpLoadBalanceAlgorithmVlan,
		VMwareDvsLacpLoadBalanceAlgorithmSrcPortId,
	}
}

func (e VMwareDvsLacpLoadBalanceAlgorithm) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareDvsLacpLoadBalanceAlgorithm"] = reflect.TypeOf((*VMwareDvsLacpLoadBalanceAlgorithm)(nil)).Elem()
}

// Multicast Filtering mode.
type VMwareDvsMulticastFilteringMode string

const (
	// Legacy filtering mode
	VMwareDvsMulticastFilteringModeLegacyFiltering = VMwareDvsMulticastFilteringMode("legacyFiltering")
	// IGMP/MLD snooping mode
	VMwareDvsMulticastFilteringModeSnooping = VMwareDvsMulticastFilteringMode("snooping")
)

func (e VMwareDvsMulticastFilteringMode) Values() []VMwareDvsMulticastFilteringMode {
	return []VMwareDvsMulticastFilteringMode{
		VMwareDvsMulticastFilteringModeLegacyFiltering,
		VMwareDvsMulticastFilteringModeSnooping,
	}
}

func (e VMwareDvsMulticastFilteringMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareDvsMulticastFilteringMode"] = reflect.TypeOf((*VMwareDvsMulticastFilteringMode)(nil)).Elem()
}

// Link Aggregation Control Protocol policy modes.
type VMwareUplinkLacpMode string

const (
	// Link Aggregation Control Protocol always sends frames along the configured uplinks
	VMwareUplinkLacpModeActive = VMwareUplinkLacpMode("active")
	// Link Aggregation Control Protocol acts as "speak when spoken to".
	VMwareUplinkLacpModePassive = VMwareUplinkLacpMode("passive")
)

func (e VMwareUplinkLacpMode) Values() []VMwareUplinkLacpMode {
	return []VMwareUplinkLacpMode{
		VMwareUplinkLacpModeActive,
		VMwareUplinkLacpModePassive,
	}
}

func (e VMwareUplinkLacpMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareUplinkLacpMode"] = reflect.TypeOf((*VMwareUplinkLacpMode)(nil)).Elem()
}

type VMwareUplinkLacpTimeoutMode string

const (
	// Set long timeout for vmnics in one LACP LAG.
	//
	// Device send fast LACPDUs
	VMwareUplinkLacpTimeoutModeFast = VMwareUplinkLacpTimeoutMode("fast")
	// Set short timeout for vmnics in one LACP LAG.
	//
	// Device send slow LACPDUs
	VMwareUplinkLacpTimeoutModeSlow = VMwareUplinkLacpTimeoutMode("slow")
)

func (e VMwareUplinkLacpTimeoutMode) Values() []VMwareUplinkLacpTimeoutMode {
	return []VMwareUplinkLacpTimeoutMode{
		VMwareUplinkLacpTimeoutModeFast,
		VMwareUplinkLacpTimeoutModeSlow,
	}
}

func (e VMwareUplinkLacpTimeoutMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VMwareUplinkLacpTimeoutMode"] = reflect.TypeOf((*VMwareUplinkLacpTimeoutMode)(nil)).Elem()
	minAPIVersionForType["VMwareUplinkLacpTimeoutMode"] = "7.0.2.0"
}

// Consumption type constants.
//
// Consumption type describes how the virtual storage object is connected and
// consumed for data by the clients.
type VStorageObjectConsumptionType string

const (
	// Disk type.
	VStorageObjectConsumptionTypeDisk = VStorageObjectConsumptionType("disk")
)

func (e VStorageObjectConsumptionType) Values() []VStorageObjectConsumptionType {
	return []VStorageObjectConsumptionType{
		VStorageObjectConsumptionTypeDisk,
	}
}

func (e VStorageObjectConsumptionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VStorageObjectConsumptionType"] = reflect.TypeOf((*VStorageObjectConsumptionType)(nil)).Elem()
}

// Deprecated as of vSphere API 4.0, use `CheckTestType_enum` instead.
//
// Types of tests available for validateMigration.
type ValidateMigrationTestType string

const (
	// Tests that examine only the configuration
	// of the virtual machine and its current host; the destination
	// resource pool and host or cluster are irrelevant.
	ValidateMigrationTestTypeSourceTests = ValidateMigrationTestType("sourceTests")
	// Tests that examine both the virtual
	// machine and the destination host or cluster; the destination
	// resource pool is irrelevant.
	//
	// This set excludes tests that fall
	// into the diskAccessibilityTests group.
	ValidateMigrationTestTypeCompatibilityTests = ValidateMigrationTestType("compatibilityTests")
	// Tests that check that the
	// destination host or cluster can see the datastores where the virtual
	// machine's virtual disks are currently located.
	//
	// The destination
	// resource pool is irrelevant. If you are planning to relocate the
	// virtual disks, do not use these tests; instead examine the relevant
	// datastore objects for your planned disk locations to see if they
	// are accessible to the destination host.
	ValidateMigrationTestTypeDiskAccessibilityTests = ValidateMigrationTestType("diskAccessibilityTests")
	// Tests that check that the destination resource
	// pool can support the virtual machine if it is powered on.
	//
	// The
	// destination host or cluster is relevant because it will affect the
	// amount of overhead memory required to run the virtual machine.
	ValidateMigrationTestTypeResourceTests = ValidateMigrationTestType("resourceTests")
)

func (e ValidateMigrationTestType) Values() []ValidateMigrationTestType {
	return []ValidateMigrationTestType{
		ValidateMigrationTestTypeSourceTests,
		ValidateMigrationTestTypeCompatibilityTests,
		ValidateMigrationTestTypeDiskAccessibilityTests,
		ValidateMigrationTestTypeResourceTests,
	}
}

func (e ValidateMigrationTestType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["ValidateMigrationTestType"] = reflect.TypeOf((*ValidateMigrationTestType)(nil)).Elem()
}

// VchaClusterMode enum defines the possible modes for a VCHA Cluster.
type VchaClusterMode string

const (
	// VCHA Cluster is enabled.
	//
	// State replication between the Active and
	// Passive node is enabled and automatic failover is allowed.
	VchaClusterModeEnabled = VchaClusterMode("enabled")
	// VCHA Cluster is disabled.
	//
	// State replication between the Active and
	// Passive node is disabled and automatic failover is not allowed.
	VchaClusterModeDisabled = VchaClusterMode("disabled")
	// VCHA Cluster is in maintenance mode.
	//
	// State replication between the
	// Active and Passive node is enabled but automatic failover
	// is not allowed.
	VchaClusterModeMaintenance = VchaClusterMode("maintenance")
)

func (e VchaClusterMode) Values() []VchaClusterMode {
	return []VchaClusterMode{
		VchaClusterModeEnabled,
		VchaClusterModeDisabled,
		VchaClusterModeMaintenance,
	}
}

func (e VchaClusterMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VchaClusterMode"] = reflect.TypeOf((*VchaClusterMode)(nil)).Elem()
}

// VchaClusterState enum defines the possible states for a VCHA Cluster.
type VchaClusterState string

const (
	// All three nodes in a VCHA Cluster are healthy and connected.
	//
	// State
	// replication between Active and Passive node is working and both
	// nodes are in sync.
	VchaClusterStateHealthy = VchaClusterState("healthy")
	// A VCHA Cluster is said to be in a degraded state for
	// either or all of the following reasons:
	// \- There is a node loss.
	//
	// \- State replication between the Active and Passive node fails.
	VchaClusterStateDegraded = VchaClusterState("degraded")
	// All three nodes are isolated from each other.
	VchaClusterStateIsolated = VchaClusterState("isolated")
)

func (e VchaClusterState) Values() []VchaClusterState {
	return []VchaClusterState{
		VchaClusterStateHealthy,
		VchaClusterStateDegraded,
		VchaClusterStateIsolated,
	}
}

func (e VchaClusterState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VchaClusterState"] = reflect.TypeOf((*VchaClusterState)(nil)).Elem()
}

type VchaNodeRole string

const (
	// Node is having a role of Active.
	//
	// In this role, node runs a vCenter
	// Server that serves client requests.
	VchaNodeRoleActive = VchaNodeRole("active")
	// Node is having a role of Passive.
	//
	// In this role node, runs as a standby
	// for the Active vCenter Server and receives state updates. This node
	// takes over the role of Active vCenter Server upon failover.
	VchaNodeRolePassive = VchaNodeRole("passive")
	// Node is having a role of Witness.
	//
	// In this role, node acts as a quorom
	// node for avoiding the classic split-brain problem.
	VchaNodeRoleWitness = VchaNodeRole("witness")
)

func (e VchaNodeRole) Values() []VchaNodeRole {
	return []VchaNodeRole{
		VchaNodeRoleActive,
		VchaNodeRolePassive,
		VchaNodeRoleWitness,
	}
}

func (e VchaNodeRole) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VchaNodeRole"] = reflect.TypeOf((*VchaNodeRole)(nil)).Elem()
}

// VchaNodeState enum defines possible state a node can be in a
// VCHA Cluster.
type VchaNodeState string

const (
	// Node is up and has joined the VCHA Cluster.
	VchaNodeStateUp = VchaNodeState("up")
	// Node is down and has left the VCHA Cluster.
	VchaNodeStateDown = VchaNodeState("down")
)

func (e VchaNodeState) Values() []VchaNodeState {
	return []VchaNodeState{
		VchaNodeStateUp,
		VchaNodeStateDown,
	}
}

func (e VchaNodeState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VchaNodeState"] = reflect.TypeOf((*VchaNodeState)(nil)).Elem()
}

type VchaState string

const (
	// VCHA cluster is configured.
	VchaStateConfigured = VchaState("configured")
	// VCHA cluster is not configured.
	VchaStateNotConfigured = VchaState("notConfigured")
	// VCHA cluster is in an invalid/dirty state.
	VchaStateInvalid = VchaState("invalid")
	// VC appliance has been prepared for VCHA cluster configuration.
	VchaStatePrepared = VchaState("prepared")
)

func (e VchaState) Values() []VchaState {
	return []VchaState{
		VchaStateConfigured,
		VchaStateNotConfigured,
		VchaStateInvalid,
		VchaStatePrepared,
	}
}

func (e VchaState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VchaState"] = reflect.TypeOf((*VchaState)(nil)).Elem()
}

// The VAppState type defines the set of states a vApp can be
// in.
//
// The transitory states between started and stopped is modeled explicitly,
// since the starting or stopping of a vApp is typically a time-consuming
// process that might take minutes to complete.
type VirtualAppVAppState string

const (
	// The vApp is currently powered on .
	VirtualAppVAppStateStarted = VirtualAppVAppState("started")
	// The vApp is currently powered off or suspended.
	VirtualAppVAppStateStopped = VirtualAppVAppState("stopped")
	// The vApp is in the process of starting.
	VirtualAppVAppStateStarting = VirtualAppVAppState("starting")
	// The vApp is in the process of stopping.
	VirtualAppVAppStateStopping = VirtualAppVAppState("stopping")
)

func (e VirtualAppVAppState) Values() []VirtualAppVAppState {
	return []VirtualAppVAppState{
		VirtualAppVAppStateStarted,
		VirtualAppVAppStateStopped,
		VirtualAppVAppStateStarting,
		VirtualAppVAppStateStopping,
	}
}

func (e VirtualAppVAppState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualAppVAppState"] = reflect.TypeOf((*VirtualAppVAppState)(nil)).Elem()
}

// Describes the change mode of the device.
//
// Applies only to virtual disks during VirtualDeviceSpec.Operation "add"
type VirtualDeviceConfigSpecChangeMode string

const (
	VirtualDeviceConfigSpecChangeModeFail = VirtualDeviceConfigSpecChangeMode("fail")
	VirtualDeviceConfigSpecChangeModeSkip = VirtualDeviceConfigSpecChangeMode("skip")
)

func (e VirtualDeviceConfigSpecChangeMode) Values() []VirtualDeviceConfigSpecChangeMode {
	return []VirtualDeviceConfigSpecChangeMode{
		VirtualDeviceConfigSpecChangeModeFail,
		VirtualDeviceConfigSpecChangeModeSkip,
	}
}

func (e VirtualDeviceConfigSpecChangeMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceConfigSpecChangeMode"] = reflect.TypeOf((*VirtualDeviceConfigSpecChangeMode)(nil)).Elem()
	minAPIVersionForType["VirtualDeviceConfigSpecChangeMode"] = "8.0.0.1"
}

// The type of operation being performed on the backing of a virtual device.
//
// Valid values are:
type VirtualDeviceConfigSpecFileOperation string

const (
	// Specifies the creation of the device backing; for example,
	// the creation of a virtual disk or floppy image file.
	VirtualDeviceConfigSpecFileOperationCreate = VirtualDeviceConfigSpecFileOperation("create")
	// Specifies the destruction of a device backing.
	VirtualDeviceConfigSpecFileOperationDestroy = VirtualDeviceConfigSpecFileOperation("destroy")
	// Specifies the deletion of the existing backing for a virtual device
	// and the creation of a new backing.
	VirtualDeviceConfigSpecFileOperationReplace = VirtualDeviceConfigSpecFileOperation("replace")
)

func (e VirtualDeviceConfigSpecFileOperation) Values() []VirtualDeviceConfigSpecFileOperation {
	return []VirtualDeviceConfigSpecFileOperation{
		VirtualDeviceConfigSpecFileOperationCreate,
		VirtualDeviceConfigSpecFileOperationDestroy,
		VirtualDeviceConfigSpecFileOperationReplace,
	}
}

func (e VirtualDeviceConfigSpecFileOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceConfigSpecFileOperation"] = reflect.TypeOf((*VirtualDeviceConfigSpecFileOperation)(nil)).Elem()
}

// The type of operation being performed on the specified virtual device.
//
// Valid values are:
type VirtualDeviceConfigSpecOperation string

const (
	// Specifies the addition of a virtual device to the configuration.
	VirtualDeviceConfigSpecOperationAdd = VirtualDeviceConfigSpecOperation("add")
	// Specifies the removal of a virtual device.
	VirtualDeviceConfigSpecOperationRemove = VirtualDeviceConfigSpecOperation("remove")
	// Specifies changes to the virtual device specification.
	VirtualDeviceConfigSpecOperationEdit = VirtualDeviceConfigSpecOperation("edit")
)

func (e VirtualDeviceConfigSpecOperation) Values() []VirtualDeviceConfigSpecOperation {
	return []VirtualDeviceConfigSpecOperation{
		VirtualDeviceConfigSpecOperationAdd,
		VirtualDeviceConfigSpecOperationRemove,
		VirtualDeviceConfigSpecOperationEdit,
	}
}

func (e VirtualDeviceConfigSpecOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceConfigSpecOperation"] = reflect.TypeOf((*VirtualDeviceConfigSpecOperation)(nil)).Elem()
}

// Contains information about connectable virtual devices when
// the virtual machine restores from a migration.
type VirtualDeviceConnectInfoMigrateConnectOp string

const (
	// Attempt to connect the virtual device when the virtual machine
	// restores from a migration.
	//
	// This property has no effect if it
	// is set on a device that is already connected.
	VirtualDeviceConnectInfoMigrateConnectOpConnect = VirtualDeviceConnectInfoMigrateConnectOp("connect")
	// Attempt to disconnect the virtual device when the virtual machine
	// restores from a migration.
	//
	// This property has no effect if it
	// is set on a device that is already disconnected.
	VirtualDeviceConnectInfoMigrateConnectOpDisconnect = VirtualDeviceConnectInfoMigrateConnectOp("disconnect")
	// Unset the property, which resets the device to its default state.
	//
	// Under most circumstances, a device will return to the same
	// connection state before the migration was initiated.
	VirtualDeviceConnectInfoMigrateConnectOpUnset = VirtualDeviceConnectInfoMigrateConnectOp("unset")
)

func (e VirtualDeviceConnectInfoMigrateConnectOp) Values() []VirtualDeviceConnectInfoMigrateConnectOp {
	return []VirtualDeviceConnectInfoMigrateConnectOp{
		VirtualDeviceConnectInfoMigrateConnectOpConnect,
		VirtualDeviceConnectInfoMigrateConnectOpDisconnect,
		VirtualDeviceConnectInfoMigrateConnectOpUnset,
	}
}

func (e VirtualDeviceConnectInfoMigrateConnectOp) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceConnectInfoMigrateConnectOp"] = reflect.TypeOf((*VirtualDeviceConnectInfoMigrateConnectOp)(nil)).Elem()
}

// Specifies the connectable virtual device status.
type VirtualDeviceConnectInfoStatus string

const (
	// The device is working correctly.
	VirtualDeviceConnectInfoStatusOk = VirtualDeviceConnectInfoStatus("ok")
	// The device has reported a recoverable error.
	//
	// For example,
	// attempting to connect to floppy device that is being used by
	// another virtual machine or some other program would result in
	// this status.
	VirtualDeviceConnectInfoStatusRecoverableError = VirtualDeviceConnectInfoStatus("recoverableError")
	// The device cannot be used.
	//
	// For example, attempting to connect to
	// a floppy device that does not exist would result in this status.
	VirtualDeviceConnectInfoStatusUnrecoverableError = VirtualDeviceConnectInfoStatus("unrecoverableError")
	// The device status is unknown, or it has not been requested to
	// connect when the VM is powered on.
	VirtualDeviceConnectInfoStatusUntried = VirtualDeviceConnectInfoStatus("untried")
)

func (e VirtualDeviceConnectInfoStatus) Values() []VirtualDeviceConnectInfoStatus {
	return []VirtualDeviceConnectInfoStatus{
		VirtualDeviceConnectInfoStatusOk,
		VirtualDeviceConnectInfoStatusRecoverableError,
		VirtualDeviceConnectInfoStatusUnrecoverableError,
		VirtualDeviceConnectInfoStatusUntried,
	}
}

func (e VirtualDeviceConnectInfoStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceConnectInfoStatus"] = reflect.TypeOf((*VirtualDeviceConnectInfoStatus)(nil)).Elem()
}

// All known file extensions.
//
// Valid ones are:
type VirtualDeviceFileExtension string

const (
	// CD ISO Image backings
	VirtualDeviceFileExtensionIso = VirtualDeviceFileExtension("iso")
	// Floppy File Backings
	VirtualDeviceFileExtensionFlp = VirtualDeviceFileExtension("flp")
	// virtual disks
	VirtualDeviceFileExtensionVmdk = VirtualDeviceFileExtension("vmdk")
	// legacy virtual disks
	VirtualDeviceFileExtensionDsk = VirtualDeviceFileExtension("dsk")
	// pre 3.0 virtual disks using Raw Disk Maps
	VirtualDeviceFileExtensionRdm = VirtualDeviceFileExtension("rdm")
)

func (e VirtualDeviceFileExtension) Values() []VirtualDeviceFileExtension {
	return []VirtualDeviceFileExtension{
		VirtualDeviceFileExtensionIso,
		VirtualDeviceFileExtensionFlp,
		VirtualDeviceFileExtensionVmdk,
		VirtualDeviceFileExtensionDsk,
		VirtualDeviceFileExtensionRdm,
	}
}

func (e VirtualDeviceFileExtension) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceFileExtension"] = reflect.TypeOf((*VirtualDeviceFileExtension)(nil)).Elem()
}

// The <code>VirtualDeviceURIBackingOptionDirection</code> enum type
// provides values for the direction of a network connection.
type VirtualDeviceURIBackingOptionDirection string

const (
	// Indicates that the virtual machine can listen for a connection
	// on the specified `VirtualDeviceURIBackingInfo.serviceURI`.
	VirtualDeviceURIBackingOptionDirectionServer = VirtualDeviceURIBackingOptionDirection("server")
	// Indicates that the virtual machine can initiate a connection
	// with a system on the network using the specified
	// `VirtualDeviceURIBackingInfo.serviceURI`.
	VirtualDeviceURIBackingOptionDirectionClient = VirtualDeviceURIBackingOptionDirection("client")
)

func (e VirtualDeviceURIBackingOptionDirection) Values() []VirtualDeviceURIBackingOptionDirection {
	return []VirtualDeviceURIBackingOptionDirection{
		VirtualDeviceURIBackingOptionDirectionServer,
		VirtualDeviceURIBackingOptionDirectionClient,
	}
}

func (e VirtualDeviceURIBackingOptionDirection) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDeviceURIBackingOptionDirection"] = reflect.TypeOf((*VirtualDeviceURIBackingOptionDirection)(nil)).Elem()
}

// The types of virtual disk adapters used by virtual disks
type VirtualDiskAdapterType string

const (
	// Use IDE emulation for the virtual disk
	VirtualDiskAdapterTypeIde = VirtualDiskAdapterType("ide")
	// Use BusLogic emulation for the virtual disk
	VirtualDiskAdapterTypeBusLogic = VirtualDiskAdapterType("busLogic")
	// Use LSILogic emulation for the virtual disk
	VirtualDiskAdapterTypeLsiLogic = VirtualDiskAdapterType("lsiLogic")
)

func (e VirtualDiskAdapterType) Values() []VirtualDiskAdapterType {
	return []VirtualDiskAdapterType{
		VirtualDiskAdapterTypeIde,
		VirtualDiskAdapterTypeBusLogic,
		VirtualDiskAdapterTypeLsiLogic,
	}
}

func (e VirtualDiskAdapterType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskAdapterType"] = reflect.TypeOf((*VirtualDiskAdapterType)(nil)).Elem()
}

// All known compatibility modes for raw disk mappings.
//
// Valid compatibility
// modes are:
//   - virtualMode
//   - physicalMode
type VirtualDiskCompatibilityMode string

const (
	// A disk device backed by a virtual compatibility mode raw disk mapping can
	// use disk modes.
	//
	// See also `VirtualDiskMode_enum`.
	VirtualDiskCompatibilityModeVirtualMode = VirtualDiskCompatibilityMode("virtualMode")
	// A disk device backed by a physical compatibility mode raw disk mapping cannot
	// use disk modes, and commands are passed straight through to the LUN
	// indicated by the raw disk mapping.
	VirtualDiskCompatibilityModePhysicalMode = VirtualDiskCompatibilityMode("physicalMode")
)

func (e VirtualDiskCompatibilityMode) Values() []VirtualDiskCompatibilityMode {
	return []VirtualDiskCompatibilityMode{
		VirtualDiskCompatibilityModeVirtualMode,
		VirtualDiskCompatibilityModePhysicalMode,
	}
}

func (e VirtualDiskCompatibilityMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskCompatibilityMode"] = reflect.TypeOf((*VirtualDiskCompatibilityMode)(nil)).Elem()
}

// The delta disk format constants
type VirtualDiskDeltaDiskFormat string

const (
	// redo-log based format
	VirtualDiskDeltaDiskFormatRedoLogFormat = VirtualDiskDeltaDiskFormat("redoLogFormat")
	// native snapshot format
	VirtualDiskDeltaDiskFormatNativeFormat = VirtualDiskDeltaDiskFormat("nativeFormat")
	// Flex-SE redo-log based format
	VirtualDiskDeltaDiskFormatSeSparseFormat = VirtualDiskDeltaDiskFormat("seSparseFormat")
)

func (e VirtualDiskDeltaDiskFormat) Values() []VirtualDiskDeltaDiskFormat {
	return []VirtualDiskDeltaDiskFormat{
		VirtualDiskDeltaDiskFormatRedoLogFormat,
		VirtualDiskDeltaDiskFormatNativeFormat,
		VirtualDiskDeltaDiskFormatSeSparseFormat,
	}
}

func (e VirtualDiskDeltaDiskFormat) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskDeltaDiskFormat"] = reflect.TypeOf((*VirtualDiskDeltaDiskFormat)(nil)).Elem()
}

// The delta disk format variant constants
type VirtualDiskDeltaDiskFormatVariant string

const (
	// vmfsSparse based redo-log format
	VirtualDiskDeltaDiskFormatVariantVmfsSparseVariant = VirtualDiskDeltaDiskFormatVariant("vmfsSparseVariant")
	// vsanSparse based redo-log format
	VirtualDiskDeltaDiskFormatVariantVsanSparseVariant = VirtualDiskDeltaDiskFormatVariant("vsanSparseVariant")
)

func (e VirtualDiskDeltaDiskFormatVariant) Values() []VirtualDiskDeltaDiskFormatVariant {
	return []VirtualDiskDeltaDiskFormatVariant{
		VirtualDiskDeltaDiskFormatVariantVmfsSparseVariant,
		VirtualDiskDeltaDiskFormatVariantVsanSparseVariant,
	}
}

func (e VirtualDiskDeltaDiskFormatVariant) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskDeltaDiskFormatVariant"] = reflect.TypeOf((*VirtualDiskDeltaDiskFormatVariant)(nil)).Elem()
}

// The list of known disk modes.
//
// The list of supported disk modes varies by the backing type. The "persistent"
// mode is supported by every backing type.
type VirtualDiskMode string

const (
	// Changes are immediately and permanently written to the virtual disk.
	VirtualDiskModePersistent = VirtualDiskMode("persistent")
	// Changes to virtual disk are made to a redo log and discarded at power off.
	VirtualDiskModeNonpersistent = VirtualDiskMode("nonpersistent")
	// Changes are made to a redo log, but you are given the option to commit or undo.
	VirtualDiskModeUndoable = VirtualDiskMode("undoable")
	// Same as persistent, but not affected by snapshots.
	VirtualDiskModeIndependent_persistent = VirtualDiskMode("independent_persistent")
	// Same as nonpersistent, but not affected by snapshots.
	VirtualDiskModeIndependent_nonpersistent = VirtualDiskMode("independent_nonpersistent")
	// Changes are appended to the redo log; you revoke changes by removing the undo log.
	VirtualDiskModeAppend = VirtualDiskMode("append")
)

func (e VirtualDiskMode) Values() []VirtualDiskMode {
	return []VirtualDiskMode{
		VirtualDiskModePersistent,
		VirtualDiskModeNonpersistent,
		VirtualDiskModeUndoable,
		VirtualDiskModeIndependent_persistent,
		VirtualDiskModeIndependent_nonpersistent,
		VirtualDiskModeAppend,
	}
}

func (e VirtualDiskMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskMode"] = reflect.TypeOf((*VirtualDiskMode)(nil)).Elem()
}

// Rule type determines how the virtual disks in a vm can be grouped
// together.
type VirtualDiskRuleSpecRuleType string

const (
	// Virtual disks in the list are grouped together and placed on
	// the same data store.
	VirtualDiskRuleSpecRuleTypeAffinity = VirtualDiskRuleSpecRuleType("affinity")
	// Virtual disks in the list are placed on different data stores.
	VirtualDiskRuleSpecRuleTypeAntiAffinity = VirtualDiskRuleSpecRuleType("antiAffinity")
	// SDRS will be disabled for the disks in the list.
	VirtualDiskRuleSpecRuleTypeDisabled = VirtualDiskRuleSpecRuleType("disabled")
)

func (e VirtualDiskRuleSpecRuleType) Values() []VirtualDiskRuleSpecRuleType {
	return []VirtualDiskRuleSpecRuleType{
		VirtualDiskRuleSpecRuleTypeAffinity,
		VirtualDiskRuleSpecRuleTypeAntiAffinity,
		VirtualDiskRuleSpecRuleTypeDisabled,
	}
}

func (e VirtualDiskRuleSpecRuleType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskRuleSpecRuleType"] = reflect.TypeOf((*VirtualDiskRuleSpecRuleType)(nil)).Elem()
}

// The sharing mode of the virtual disk.
//
// Setting the value to sharingMultiWriter means that multiple virtual
// machines can write to the virtual disk. This sharing mode is allowed
// only for eagerly zeroed thick virtual disks.
type VirtualDiskSharing string

const (
	// The virtual disk is not shared.
	VirtualDiskSharingSharingNone = VirtualDiskSharing("sharingNone")
	// The virtual disk is shared between multiple virtual machines.
	VirtualDiskSharingSharingMultiWriter = VirtualDiskSharing("sharingMultiWriter")
)

func (e VirtualDiskSharing) Values() []VirtualDiskSharing {
	return []VirtualDiskSharing{
		VirtualDiskSharingSharingNone,
		VirtualDiskSharingSharingMultiWriter,
	}
}

func (e VirtualDiskSharing) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskSharing"] = reflect.TypeOf((*VirtualDiskSharing)(nil)).Elem()
}

// The types of virtual disks that can be created or cloned.
type VirtualDiskType string

const (
	// A preallocated disk has all space allocated at creation time
	// and the space is zeroed on demand as the space is used.
	VirtualDiskTypePreallocated = VirtualDiskType("preallocated")
	// Space required for thin-provisioned virtual disk is allocated and
	// zeroed on demand as the space is used.
	VirtualDiskTypeThin = VirtualDiskType("thin")
	// A sparse (allocate on demand) format with additional space
	// optimizations.
	VirtualDiskTypeSeSparse = VirtualDiskType("seSparse")
	// Virtual compatibility mode raw disk mapping.
	//
	// An rdm virtual disk
	// grants access to the entire raw disk and the virtual disk can
	// participate in snapshots.
	VirtualDiskTypeRdm = VirtualDiskType("rdm")
	// Physical compatibility mode (pass-through) raw disk mapping.
	//
	// An rdmp
	// virtual disk passes SCSI commands directly to the hardware, but the
	// virtual disk cannot participate in snapshots.
	VirtualDiskTypeRdmp = VirtualDiskType("rdmp")
	// Raw device.
	VirtualDiskTypeRaw = VirtualDiskType("raw")
	// A redo log disk.
	//
	// This format is only applicable as a destination format
	// in a clone operation, and not usable for disk creation.
	VirtualDiskTypeDelta = VirtualDiskType("delta")
	// A sparse disk with 2GB maximum extent size.
	//
	// Disks in this format
	// can be used with other VMware products. The 2GB extent size
	// makes these disks easier to burn to dvd or use on filesystems that
	// don't support large files. This format is only applicable as a
	// destination format in a clone operation, and not usable for disk
	// creation.
	VirtualDiskTypeSparse2Gb = VirtualDiskType("sparse2Gb")
	// A thick disk with 2GB maximum extent size.
	//
	// Disks in this format
	// can be used with other VMware products. The 2GB extent size
	// makes these disks easier to burn to dvd or use on filesystems that
	// don't support large files. This format is only applicable as a
	// destination format in a clone operation, and not usable for disk
	// creation.
	VirtualDiskTypeThick2Gb = VirtualDiskType("thick2Gb")
	// An eager zeroed thick disk has all space allocated and wiped clean
	// of any previous contents on the physical media at creation time.
	//
	// Such disks may take longer time during creation compared to other
	// disk formats.
	VirtualDiskTypeEagerZeroedThick = VirtualDiskType("eagerZeroedThick")
	// A sparse monolithic disk.
	//
	// Disks in this format can be used with other
	// VMware products. This format is only applicable as a destination
	// format in a clone operation, and not usable for disk creation.
	VirtualDiskTypeSparseMonolithic = VirtualDiskType("sparseMonolithic")
	// A preallocated monolithic disk.
	//
	// Disks in this format can be used with
	// other VMware products. This format is only applicable as a destination
	// format in a clone operation, and not usable for disk creation.
	VirtualDiskTypeFlatMonolithic = VirtualDiskType("flatMonolithic")
	// Deprecated as of vSphere API 4.x, use `eagerZeroedThick` instead
	// for clustering application, and `preallocated` for other applications.
	//
	// A thick disk has all space allocated at creation time.
	//
	// This
	// space may contain stale data on the physical media. Thick disks
	// are primarily used for virtual machine clustering, but they are
	// generally insecure and should not be used. Due to better performance
	// and security properties, the use of the 'preallocated' format is
	// preferred over this format.
	VirtualDiskTypeThick = VirtualDiskType("thick")
)

func (e VirtualDiskType) Values() []VirtualDiskType {
	return []VirtualDiskType{
		VirtualDiskTypePreallocated,
		VirtualDiskTypeThin,
		VirtualDiskTypeSeSparse,
		VirtualDiskTypeRdm,
		VirtualDiskTypeRdmp,
		VirtualDiskTypeRaw,
		VirtualDiskTypeDelta,
		VirtualDiskTypeSparse2Gb,
		VirtualDiskTypeThick2Gb,
		VirtualDiskTypeEagerZeroedThick,
		VirtualDiskTypeSparseMonolithic,
		VirtualDiskTypeFlatMonolithic,
		VirtualDiskTypeThick,
	}
}

func (e VirtualDiskType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskType"] = reflect.TypeOf((*VirtualDiskType)(nil)).Elem()
}

// Pre-defined constants for cache consistency types
type VirtualDiskVFlashCacheConfigInfoCacheConsistencyType string

const (
	// With strong consistency, it ensures that
	// a crash will leave the cache data consistent.
	VirtualDiskVFlashCacheConfigInfoCacheConsistencyTypeStrong = VirtualDiskVFlashCacheConfigInfoCacheConsistencyType("strong")
	// Cache data consistency is not guaranteed after a crash.
	VirtualDiskVFlashCacheConfigInfoCacheConsistencyTypeWeak = VirtualDiskVFlashCacheConfigInfoCacheConsistencyType("weak")
)

func (e VirtualDiskVFlashCacheConfigInfoCacheConsistencyType) Values() []VirtualDiskVFlashCacheConfigInfoCacheConsistencyType {
	return []VirtualDiskVFlashCacheConfigInfoCacheConsistencyType{
		VirtualDiskVFlashCacheConfigInfoCacheConsistencyTypeStrong,
		VirtualDiskVFlashCacheConfigInfoCacheConsistencyTypeWeak,
	}
}

func (e VirtualDiskVFlashCacheConfigInfoCacheConsistencyType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskVFlashCacheConfigInfoCacheConsistencyType"] = reflect.TypeOf((*VirtualDiskVFlashCacheConfigInfoCacheConsistencyType)(nil)).Elem()
}

// Pre-defined constants for cache modes.
type VirtualDiskVFlashCacheConfigInfoCacheMode string

const (
	// In write-through cache mode, writes to the cache cause writes
	// to the underlying storage.
	//
	// The cache acts as a facade to the underlying
	// storage.
	VirtualDiskVFlashCacheConfigInfoCacheModeWrite_thru = VirtualDiskVFlashCacheConfigInfoCacheMode("write_thru")
	// In write-back mode, writes to the cache do not go to the underlying storage
	// right away.
	//
	// Cache holds data temporarily till it can be permanently saved or
	// otherwise modified.
	VirtualDiskVFlashCacheConfigInfoCacheModeWrite_back = VirtualDiskVFlashCacheConfigInfoCacheMode("write_back")
)

func (e VirtualDiskVFlashCacheConfigInfoCacheMode) Values() []VirtualDiskVFlashCacheConfigInfoCacheMode {
	return []VirtualDiskVFlashCacheConfigInfoCacheMode{
		VirtualDiskVFlashCacheConfigInfoCacheModeWrite_thru,
		VirtualDiskVFlashCacheConfigInfoCacheModeWrite_back,
	}
}

func (e VirtualDiskVFlashCacheConfigInfoCacheMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualDiskVFlashCacheConfigInfoCacheMode"] = reflect.TypeOf((*VirtualDiskVFlashCacheConfigInfoCacheMode)(nil)).Elem()
}

// Possible device names for legacy network backing option are listed below.
//
// Note: This is not an exhaustive list. It is possible to specify
// a specific device as well.
// For example, on ESX hosts, the device name could be specified as "vmnic\[0-9\]"
// or vmnet\_\[0-9\].
// For VMware Server Windows hosts, the device name could be specified as "vmnet\[0-9\]"
// and for VMware Server Linux hosts, the device name could be specified as "/dev/vmnet\[0-9\]"
// depending on what devices are available on that particular host.
type VirtualEthernetCardLegacyNetworkDeviceName string

const (
	VirtualEthernetCardLegacyNetworkDeviceNameBridged  = VirtualEthernetCardLegacyNetworkDeviceName("bridged")
	VirtualEthernetCardLegacyNetworkDeviceNameNat      = VirtualEthernetCardLegacyNetworkDeviceName("nat")
	VirtualEthernetCardLegacyNetworkDeviceNameHostonly = VirtualEthernetCardLegacyNetworkDeviceName("hostonly")
)

func (e VirtualEthernetCardLegacyNetworkDeviceName) Values() []VirtualEthernetCardLegacyNetworkDeviceName {
	return []VirtualEthernetCardLegacyNetworkDeviceName{
		VirtualEthernetCardLegacyNetworkDeviceNameBridged,
		VirtualEthernetCardLegacyNetworkDeviceNameNat,
		VirtualEthernetCardLegacyNetworkDeviceNameHostonly,
	}
}

func (e VirtualEthernetCardLegacyNetworkDeviceName) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualEthernetCardLegacyNetworkDeviceName"] = reflect.TypeOf((*VirtualEthernetCardLegacyNetworkDeviceName)(nil)).Elem()
}

// The enumeration of all known valid MAC address types.
type VirtualEthernetCardMacType string

const (
	// A statistically assigned MAC address.
	VirtualEthernetCardMacTypeManual = VirtualEthernetCardMacType("manual")
	// An automatically generated MAC address.
	VirtualEthernetCardMacTypeGenerated = VirtualEthernetCardMacType("generated")
	// A MAC address assigned by VirtualCenter.
	VirtualEthernetCardMacTypeAssigned = VirtualEthernetCardMacType("assigned")
)

func (e VirtualEthernetCardMacType) Values() []VirtualEthernetCardMacType {
	return []VirtualEthernetCardMacType{
		VirtualEthernetCardMacTypeManual,
		VirtualEthernetCardMacTypeGenerated,
		VirtualEthernetCardMacTypeAssigned,
	}
}

func (e VirtualEthernetCardMacType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualEthernetCardMacType"] = reflect.TypeOf((*VirtualEthernetCardMacType)(nil)).Elem()
}

type VirtualHardwareMotherboardLayout string

const (
	// Single i440BX host bridge.
	VirtualHardwareMotherboardLayoutI440bxHostBridge = VirtualHardwareMotherboardLayout("i440bxHostBridge")
	// Multiple ACPI host bridges.
	VirtualHardwareMotherboardLayoutAcpiHostBridges = VirtualHardwareMotherboardLayout("acpiHostBridges")
)

func (e VirtualHardwareMotherboardLayout) Values() []VirtualHardwareMotherboardLayout {
	return []VirtualHardwareMotherboardLayout{
		VirtualHardwareMotherboardLayoutI440bxHostBridge,
		VirtualHardwareMotherboardLayoutAcpiHostBridges,
	}
}

func (e VirtualHardwareMotherboardLayout) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualHardwareMotherboardLayout"] = reflect.TypeOf((*VirtualHardwareMotherboardLayout)(nil)).Elem()
	minAPIVersionForType["VirtualHardwareMotherboardLayout"] = "8.0.0.1"
}

// Application heartbeat status type.
type VirtualMachineAppHeartbeatStatusType string

const (
	// Heartbeat status is disabled
	VirtualMachineAppHeartbeatStatusTypeAppStatusGray = VirtualMachineAppHeartbeatStatusType("appStatusGray")
	// Heartbeat status is OK
	VirtualMachineAppHeartbeatStatusTypeAppStatusGreen = VirtualMachineAppHeartbeatStatusType("appStatusGreen")
	// Heartbeating has stopped
	VirtualMachineAppHeartbeatStatusTypeAppStatusRed = VirtualMachineAppHeartbeatStatusType("appStatusRed")
)

func (e VirtualMachineAppHeartbeatStatusType) Values() []VirtualMachineAppHeartbeatStatusType {
	return []VirtualMachineAppHeartbeatStatusType{
		VirtualMachineAppHeartbeatStatusTypeAppStatusGray,
		VirtualMachineAppHeartbeatStatusTypeAppStatusGreen,
		VirtualMachineAppHeartbeatStatusTypeAppStatusRed,
	}
}

func (e VirtualMachineAppHeartbeatStatusType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineAppHeartbeatStatusType"] = reflect.TypeOf((*VirtualMachineAppHeartbeatStatusType)(nil)).Elem()
}

type VirtualMachineBootOptionsNetworkBootProtocolType string

const (
	// PXE (or Apple NetBoot) over IPv4.
	//
	// The default.
	VirtualMachineBootOptionsNetworkBootProtocolTypeIpv4 = VirtualMachineBootOptionsNetworkBootProtocolType("ipv4")
	// PXE over IPv6.
	//
	// Only meaningful for EFI virtual machines.
	VirtualMachineBootOptionsNetworkBootProtocolTypeIpv6 = VirtualMachineBootOptionsNetworkBootProtocolType("ipv6")
)

func (e VirtualMachineBootOptionsNetworkBootProtocolType) Values() []VirtualMachineBootOptionsNetworkBootProtocolType {
	return []VirtualMachineBootOptionsNetworkBootProtocolType{
		VirtualMachineBootOptionsNetworkBootProtocolTypeIpv4,
		VirtualMachineBootOptionsNetworkBootProtocolTypeIpv6,
	}
}

func (e VirtualMachineBootOptionsNetworkBootProtocolType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineBootOptionsNetworkBootProtocolType"] = reflect.TypeOf((*VirtualMachineBootOptionsNetworkBootProtocolType)(nil)).Elem()
}

type VirtualMachineCertThumbprintHashAlgorithm string

const (
	// SHA256
	VirtualMachineCertThumbprintHashAlgorithmSha256 = VirtualMachineCertThumbprintHashAlgorithm("sha256")
)

func (e VirtualMachineCertThumbprintHashAlgorithm) Values() []VirtualMachineCertThumbprintHashAlgorithm {
	return []VirtualMachineCertThumbprintHashAlgorithm{
		VirtualMachineCertThumbprintHashAlgorithmSha256,
	}
}

func (e VirtualMachineCertThumbprintHashAlgorithm) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineCertThumbprintHashAlgorithm"] = reflect.TypeOf((*VirtualMachineCertThumbprintHashAlgorithm)(nil)).Elem()
	minAPIVersionForType["VirtualMachineCertThumbprintHashAlgorithm"] = "7.0.3.1"
}

// TPM provisioning policies used when cloning a VM with a virtual TPM
type VirtualMachineCloneSpecTpmProvisionPolicy string

const (
	// The virtual TPM is copied.
	//
	// The virtual machine clone will have access
	// to the original virtual machine's TPM secrets.
	VirtualMachineCloneSpecTpmProvisionPolicyCopy = VirtualMachineCloneSpecTpmProvisionPolicy("copy")
	// The virtual TPM is replaced with a new one.
	//
	// The virtual machine clone
	// will not have access to the original virtual machine's TPM secrets.
	VirtualMachineCloneSpecTpmProvisionPolicyReplace = VirtualMachineCloneSpecTpmProvisionPolicy("replace")
)

func (e VirtualMachineCloneSpecTpmProvisionPolicy) Values() []VirtualMachineCloneSpecTpmProvisionPolicy {
	return []VirtualMachineCloneSpecTpmProvisionPolicy{
		VirtualMachineCloneSpecTpmProvisionPolicyCopy,
		VirtualMachineCloneSpecTpmProvisionPolicyReplace,
	}
}

func (e VirtualMachineCloneSpecTpmProvisionPolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineCloneSpecTpmProvisionPolicy"] = reflect.TypeOf((*VirtualMachineCloneSpecTpmProvisionPolicy)(nil)).Elem()
	minAPIVersionForType["VirtualMachineCloneSpecTpmProvisionPolicy"] = "8.0.0.1"
}

// The NPIV WWN source type.
type VirtualMachineConfigInfoNpivWwnType string

const (
	// This set of WWNs is generated by VC server.
	VirtualMachineConfigInfoNpivWwnTypeVc = VirtualMachineConfigInfoNpivWwnType("vc")
	// This set of WWNs is generated by Host Agent.
	VirtualMachineConfigInfoNpivWwnTypeHost = VirtualMachineConfigInfoNpivWwnType("host")
	// This set of WWNs is provided by the client.
	VirtualMachineConfigInfoNpivWwnTypeExternal = VirtualMachineConfigInfoNpivWwnType("external")
)

func (e VirtualMachineConfigInfoNpivWwnType) Values() []VirtualMachineConfigInfoNpivWwnType {
	return []VirtualMachineConfigInfoNpivWwnType{
		VirtualMachineConfigInfoNpivWwnTypeVc,
		VirtualMachineConfigInfoNpivWwnTypeHost,
		VirtualMachineConfigInfoNpivWwnTypeExternal,
	}
}

func (e VirtualMachineConfigInfoNpivWwnType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineConfigInfoNpivWwnType"] = reflect.TypeOf((*VirtualMachineConfigInfoNpivWwnType)(nil)).Elem()
}

// Available choices for virtual machine swapfile placement policy.
//
// This is
// the set of legal values for the virtual machine configuration's
// `swapPlacement` property. All
// values except for "inherit" and "vmConfigured" are also valid values for
// a compute resource configuration's
// `vmSwapPlacement`
// property.
type VirtualMachineConfigInfoSwapPlacementType string

const (
	// Honor the virtual machine swapfile placement policy of the compute
	// resource that contains this virtual machine.
	VirtualMachineConfigInfoSwapPlacementTypeInherit = VirtualMachineConfigInfoSwapPlacementType("inherit")
	// Store the swapfile in the same directory as the virtual machine.
	VirtualMachineConfigInfoSwapPlacementTypeVmDirectory = VirtualMachineConfigInfoSwapPlacementType("vmDirectory")
	// Store the swapfile in the datastore specified by the
	// `localSwapDatastore`
	// property of the virtual machine's host, if that property is set and
	// indicates a datastore with sufficient free space.
	//
	// Otherwise store the
	// swapfile in the same directory as the virtual machine.
	//
	// Note: This setting may degrade VMotion performance.
	VirtualMachineConfigInfoSwapPlacementTypeHostLocal = VirtualMachineConfigInfoSwapPlacementType("hostLocal")
)

func (e VirtualMachineConfigInfoSwapPlacementType) Values() []VirtualMachineConfigInfoSwapPlacementType {
	return []VirtualMachineConfigInfoSwapPlacementType{
		VirtualMachineConfigInfoSwapPlacementTypeInherit,
		VirtualMachineConfigInfoSwapPlacementTypeVmDirectory,
		VirtualMachineConfigInfoSwapPlacementTypeHostLocal,
	}
}

func (e VirtualMachineConfigInfoSwapPlacementType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineConfigInfoSwapPlacementType"] = reflect.TypeOf((*VirtualMachineConfigInfoSwapPlacementType)(nil)).Elem()
}

// The set of valid encrypted Fault Tolerance modes for a VM.
//
// If the VM is encrypted, its encrypted Fault Tolerance mode
type VirtualMachineConfigSpecEncryptedFtModes string

const (
	// Do not use encrypted Fault Tolerance, even if available.
	VirtualMachineConfigSpecEncryptedFtModesFtEncryptionDisabled = VirtualMachineConfigSpecEncryptedFtModes("ftEncryptionDisabled")
	// Use encrypted Fault Tolerance if source and destination hosts
	// support it, fall back to unencrypted Fault Tolerance otherwise.
	//
	// This is the default option.
	VirtualMachineConfigSpecEncryptedFtModesFtEncryptionOpportunistic = VirtualMachineConfigSpecEncryptedFtModes("ftEncryptionOpportunistic")
	// Allow only encrypted Fault Tolerance.
	//
	// If either the source or
	// destination host does not support encrypted Fault Tolerance,
	// do not allow the Fault Tolerance to occur.
	VirtualMachineConfigSpecEncryptedFtModesFtEncryptionRequired = VirtualMachineConfigSpecEncryptedFtModes("ftEncryptionRequired")
)

func (e VirtualMachineConfigSpecEncryptedFtModes) Values() []VirtualMachineConfigSpecEncryptedFtModes {
	return []VirtualMachineConfigSpecEncryptedFtModes{
		VirtualMachineConfigSpecEncryptedFtModesFtEncryptionDisabled,
		VirtualMachineConfigSpecEncryptedFtModesFtEncryptionOpportunistic,
		VirtualMachineConfigSpecEncryptedFtModesFtEncryptionRequired,
	}
}

func (e VirtualMachineConfigSpecEncryptedFtModes) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineConfigSpecEncryptedFtModes"] = reflect.TypeOf((*VirtualMachineConfigSpecEncryptedFtModes)(nil)).Elem()
	minAPIVersionForType["VirtualMachineConfigSpecEncryptedFtModes"] = "7.0.2.0"
}

// The set of valid encrypted vMotion modes for a VM.
//
// If the VM is encrypted, its encrypted vMotion mode will be required.
type VirtualMachineConfigSpecEncryptedVMotionModes string

const (
	// Do not use encrypted vMotion, even if available.
	VirtualMachineConfigSpecEncryptedVMotionModesDisabled = VirtualMachineConfigSpecEncryptedVMotionModes("disabled")
	// Use encrypted vMotion if source and destination hosts support it,
	// fall back to unencrypted vMotion otherwise.
	//
	// This is the default option.
	VirtualMachineConfigSpecEncryptedVMotionModesOpportunistic = VirtualMachineConfigSpecEncryptedVMotionModes("opportunistic")
	// Allow only encrypted vMotion.
	//
	// If the source or destination host does
	// not support vMotion encryption, do not allow the vMotion to occur.
	VirtualMachineConfigSpecEncryptedVMotionModesRequired = VirtualMachineConfigSpecEncryptedVMotionModes("required")
)

func (e VirtualMachineConfigSpecEncryptedVMotionModes) Values() []VirtualMachineConfigSpecEncryptedVMotionModes {
	return []VirtualMachineConfigSpecEncryptedVMotionModes{
		VirtualMachineConfigSpecEncryptedVMotionModesDisabled,
		VirtualMachineConfigSpecEncryptedVMotionModesOpportunistic,
		VirtualMachineConfigSpecEncryptedVMotionModesRequired,
	}
}

func (e VirtualMachineConfigSpecEncryptedVMotionModes) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineConfigSpecEncryptedVMotionModes"] = reflect.TypeOf((*VirtualMachineConfigSpecEncryptedVMotionModes)(nil)).Elem()
}

// The root WWN operation mode.
type VirtualMachineConfigSpecNpivWwnOp string

const (
	// Generate a new set of WWNs and assign it to the virtual machine.
	VirtualMachineConfigSpecNpivWwnOpGenerate = VirtualMachineConfigSpecNpivWwnOp("generate")
	// Take a client-specified set of WWNs (specified in "wwn" property) and
	// assign them to the virtual machine.
	//
	// If the new WWN quntity are more
	// than existing then we will append them to the existing list of WWNs.
	VirtualMachineConfigSpecNpivWwnOpSet = VirtualMachineConfigSpecNpivWwnOp("set")
	// Remove the currently assigned WWNs from the virtual machine.
	VirtualMachineConfigSpecNpivWwnOpRemove = VirtualMachineConfigSpecNpivWwnOp("remove")
	// Generate a new set of WWNs and append them to the existing list
	VirtualMachineConfigSpecNpivWwnOpExtend = VirtualMachineConfigSpecNpivWwnOp("extend")
)

func (e VirtualMachineConfigSpecNpivWwnOp) Values() []VirtualMachineConfigSpecNpivWwnOp {
	return []VirtualMachineConfigSpecNpivWwnOp{
		VirtualMachineConfigSpecNpivWwnOpGenerate,
		VirtualMachineConfigSpecNpivWwnOpSet,
		VirtualMachineConfigSpecNpivWwnOpRemove,
		VirtualMachineConfigSpecNpivWwnOpExtend,
	}
}

func (e VirtualMachineConfigSpecNpivWwnOp) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineConfigSpecNpivWwnOp"] = reflect.TypeOf((*VirtualMachineConfigSpecNpivWwnOp)(nil)).Elem()
}

// The connectivity state of a virtual machine.
//
// When the API is provided directly by
// a server product, such as ESX Server, then the disconnected state is not
// possible. However, when accessed through VirtualCenter, the state of a virtual
// machine is set to disconnected if the hosts that manage the virtual
// machine becomes unavailable.
type VirtualMachineConnectionState string

const (
	// The server has access to the virtual machine.
	VirtualMachineConnectionStateConnected = VirtualMachineConnectionState("connected")
	// The server is currently disconnected from the virtual machine, since its
	// host is disconnected.
	//
	// See general comment for this enumerated type for more
	// details.
	VirtualMachineConnectionStateDisconnected = VirtualMachineConnectionState("disconnected")
	// The virtual machine is no longer registered on the host it is associated
	// with.
	//
	// For example, a virtual machine that is unregistered or deleted
	// directly on a host managed by VirtualCenter shows up in this state.
	VirtualMachineConnectionStateOrphaned = VirtualMachineConnectionState("orphaned")
	// One or more of the virtual machine configuration files are inaccessible.
	//
	// For
	// example, this can be due to transient disk failures. In this case, no
	// configuration can be returned for a virtual machine.
	VirtualMachineConnectionStateInaccessible = VirtualMachineConnectionState("inaccessible")
	// The virtual machine configuration format is invalid.
	//
	// Thus, it is accessible
	// on disk, but corrupted in a way that does not allow the server to read the
	// content. In this case, no configuration can be returned for a virtual
	// machine.
	VirtualMachineConnectionStateInvalid = VirtualMachineConnectionState("invalid")
)

func (e VirtualMachineConnectionState) Values() []VirtualMachineConnectionState {
	return []VirtualMachineConnectionState{
		VirtualMachineConnectionStateConnected,
		VirtualMachineConnectionStateDisconnected,
		VirtualMachineConnectionStateOrphaned,
		VirtualMachineConnectionStateInaccessible,
		VirtualMachineConnectionStateInvalid,
	}
}

func (e VirtualMachineConnectionState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineConnectionState"] = reflect.TypeOf((*VirtualMachineConnectionState)(nil)).Elem()
}

// The crypto state of a encrypted virtual machine.
type VirtualMachineCryptoState string

const (
	// The virtual machine is in unlocked state.
	VirtualMachineCryptoStateUnlocked = VirtualMachineCryptoState("unlocked")
	// The virtual machine is in locked state for the configuration key missing
	// on the ESX host where the VM is registered.
	VirtualMachineCryptoStateLocked = VirtualMachineCryptoState("locked")
)

func (e VirtualMachineCryptoState) Values() []VirtualMachineCryptoState {
	return []VirtualMachineCryptoState{
		VirtualMachineCryptoStateUnlocked,
		VirtualMachineCryptoStateLocked,
	}
}

func (e VirtualMachineCryptoState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineCryptoState"] = reflect.TypeOf((*VirtualMachineCryptoState)(nil)).Elem()
}

type VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther string

const (
	// The virtual machine's host does not support VMDirectPath Gen 2.
	//
	// See also `HostCapability.vmDirectPathGen2Supported`.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOtherVmNptIncompatibleHost = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther("vmNptIncompatibleHost")
	// The configuration or state of the attached network prevents
	// VMDirectPath Gen 2.
	//
	// Refer to
	// `vmDirectPathGen2InactiveReasonNetwork`
	// and/or
	// `vmDirectPathGen2InactiveReasonExtended`
	// in the RuntimeInfo of the DistributedVirtualPort connected to this
	// device.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOtherVmNptIncompatibleNetwork = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther("vmNptIncompatibleNetwork")
)

func (e VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther) Values() []VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther {
	return []VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther{
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOtherVmNptIncompatibleHost,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOtherVmNptIncompatibleNetwork,
	}
}

func (e VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther"] = reflect.TypeOf((*VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonOther)(nil)).Elem()
}

type VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm string

const (
	// The virtual machine's guest OS does not support
	// VMDirectPath Gen 2.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleGuest = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptIncompatibleGuest")
	// The virtual machine's guest network driver does not support
	// VMDirectPath Gen 2.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleGuestDriver = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptIncompatibleGuestDriver")
	// The device type does not support VMDirectPath Gen 2.
	//
	// See also `VirtualEthernetCardOption.vmDirectPathGen2Supported`.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleAdapterType = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptIncompatibleAdapterType")
	// The virtual machine's network adapter is disabled or
	// disconnected, and thus is not participating in VMDirectPath Gen 2.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptDisabledOrDisconnectedAdapter = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptDisabledOrDisconnectedAdapter")
	// The virtual machine's network adapter has features enabled
	// which preclude it participating in VMDirectPath Gen 2 such
	// as INT-x or PXE booting.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleAdapterFeatures = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptIncompatibleAdapterFeatures")
	// The device backing is not a DistributedVirtualPortBacking.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleBackingType = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptIncompatibleBackingType")
	// The virtual machine does not have full memory reservation
	// required to activate VMDirectPath Gen 2.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptInsufficientMemoryReservation = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptInsufficientMemoryReservation")
	// Deprecated as of vSphere API 6.0.
	//
	// The virtual machine is configured for Fault Tolerance or
	// Record &amp; Replay, which prevents VMDirectPath Gen 2.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptFaultToleranceOrRecordReplayConfigured = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptFaultToleranceOrRecordReplayConfigured")
	// Some networking feature has placed a conflicting IOChain on
	// the network adapter, which prevents VMDirectPath Gen 2.
	//
	// Examples
	// include DVFilter.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptConflictingIOChainConfigured = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptConflictingIOChainConfigured")
	// The virtual machine monitor is exercising functionality which
	// which prevents VMDirectPath Gen 2.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptMonitorBlocks = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptMonitorBlocks")
	// VMDirectPath Gen 2 is temporarily suspended while the virtual
	// machine executes an operation such as suspend.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptConflictingOperationInProgress = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptConflictingOperationInProgress")
	// VMDirectPath Gen 2 is unavailable due to an unforeseen runtime error
	// in the virtualization platform (typically resource constraints.)
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptRuntimeError = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptRuntimeError")
	// VMDirectPath Gen 2 is unavailable due to host run out of intr
	// vector in host.
	//
	// Guest can configure the vNIC to use less rx/tx
	// queues or use MSI instead of MSIX.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptOutOfIntrVector = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptOutOfIntrVector")
	// VMDirectPath Gen 2 is unavailable due to Incompatibe feature
	// VMCI is active in the current VM.
	//
	// Kill the relevant VMCI
	// application(s) and restart the VM will allow the vNIC(s) to enter
	// passthrough mode.
	VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptVMCIActive = VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm("vmNptVMCIActive")
)

func (e VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm) Values() []VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm {
	return []VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm{
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleGuest,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleGuestDriver,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleAdapterType,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptDisabledOrDisconnectedAdapter,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleAdapterFeatures,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptIncompatibleBackingType,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptInsufficientMemoryReservation,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptFaultToleranceOrRecordReplayConfigured,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptConflictingIOChainConfigured,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptMonitorBlocks,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptConflictingOperationInProgress,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptRuntimeError,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptOutOfIntrVector,
		VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVmVmNptVMCIActive,
	}
}

func (e VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm"] = reflect.TypeOf((*VirtualMachineDeviceRuntimeInfoVirtualEthernetCardRuntimeStateVmDirectPathGen2InactiveReasonVm)(nil)).Elem()
}

// The FaultToleranceState type defines a simple set of states for a
// fault tolerant virtual machine:
// disabled, starting, and enabled.
type VirtualMachineFaultToleranceState string

const (
	// This state indicates that the virtual machine has not been
	// configured for fault tolerance.
	VirtualMachineFaultToleranceStateNotConfigured = VirtualMachineFaultToleranceState("notConfigured")
	// For a virtual machine that is the primary in a fault tolerant group,
	// this state indicates that the virtual machine has at least one
	// registered secondary, but no secondary is enabled.
	//
	// For a virtual machine that is the secondary in a fault tolerant
	// group, this state indicates that the secondary is disabled.
	VirtualMachineFaultToleranceStateDisabled = VirtualMachineFaultToleranceState("disabled")
	// For a virtual machine that is the primary in a fault tolerant group,
	// this state indicates that the virtual machine is not currently
	// powered on, but has at least one enabled secondary
	// For a virtual machine that is the secondary in a fault tolerant
	// group, this state indicates that the secondary is enabled, but is
	// not currently powered on.
	VirtualMachineFaultToleranceStateEnabled = VirtualMachineFaultToleranceState("enabled")
	// For a virtual machine that is the primary in a fault tolerant group,
	// this state indicates that the virtual machine is powered on and
	// has at least one enabled secondary, but no secondary is currently
	// active.
	//
	// This state is not valid for a virtual machine that is a secondary
	// in a fault tolerant group.
	VirtualMachineFaultToleranceStateNeedSecondary = VirtualMachineFaultToleranceState("needSecondary")
	// For a virtual machine that is the primary in a fault tolerant group,
	// this state indicates that the virtual machine is powered on and has
	// at least one secondary that is synchronizing its state with the
	// primary.
	//
	// For a virtual machine that is the secondary in a fault tolerant
	// group, this state indicates that the secondary is powered on and is
	// synchronizing its state with the primary virtual machine.
	VirtualMachineFaultToleranceStateStarting = VirtualMachineFaultToleranceState("starting")
	// This state indicates that the virtual machine is running with fault
	// tolerance protection.
	VirtualMachineFaultToleranceStateRunning = VirtualMachineFaultToleranceState("running")
)

func (e VirtualMachineFaultToleranceState) Values() []VirtualMachineFaultToleranceState {
	return []VirtualMachineFaultToleranceState{
		VirtualMachineFaultToleranceStateNotConfigured,
		VirtualMachineFaultToleranceStateDisabled,
		VirtualMachineFaultToleranceStateEnabled,
		VirtualMachineFaultToleranceStateNeedSecondary,
		VirtualMachineFaultToleranceStateStarting,
		VirtualMachineFaultToleranceStateRunning,
	}
}

func (e VirtualMachineFaultToleranceState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineFaultToleranceState"] = reflect.TypeOf((*VirtualMachineFaultToleranceState)(nil)).Elem()
}

// The FaultToleranceType defines the type of fault tolerance, if any,
// the virtual machine is configured for.
type VirtualMachineFaultToleranceType string

const (
	// FT not set
	VirtualMachineFaultToleranceTypeUnset = VirtualMachineFaultToleranceType("unset")
	// Record/replay
	VirtualMachineFaultToleranceTypeRecordReplay = VirtualMachineFaultToleranceType("recordReplay")
	// Checkpointing
	VirtualMachineFaultToleranceTypeCheckpointing = VirtualMachineFaultToleranceType("checkpointing")
)

func (e VirtualMachineFaultToleranceType) Values() []VirtualMachineFaultToleranceType {
	return []VirtualMachineFaultToleranceType{
		VirtualMachineFaultToleranceTypeUnset,
		VirtualMachineFaultToleranceTypeRecordReplay,
		VirtualMachineFaultToleranceTypeCheckpointing,
	}
}

func (e VirtualMachineFaultToleranceType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineFaultToleranceType"] = reflect.TypeOf((*VirtualMachineFaultToleranceType)(nil)).Elem()
}

// File-type constants.
type VirtualMachineFileLayoutExFileType string

const (
	// Config (vmx) file.
	VirtualMachineFileLayoutExFileTypeConfig = VirtualMachineFileLayoutExFileType("config")
	// Extended config (vmxf) file.
	VirtualMachineFileLayoutExFileTypeExtendedConfig = VirtualMachineFileLayoutExFileType("extendedConfig")
	// Disk descriptor (vmdk) file.
	VirtualMachineFileLayoutExFileTypeDiskDescriptor = VirtualMachineFileLayoutExFileType("diskDescriptor")
	// Disk extent (-flat/-delta/-s/-rdm/-rdmp.vmdk) file.
	VirtualMachineFileLayoutExFileTypeDiskExtent = VirtualMachineFileLayoutExFileType("diskExtent")
	// Disk digest descriptor file.
	VirtualMachineFileLayoutExFileTypeDigestDescriptor = VirtualMachineFileLayoutExFileType("digestDescriptor")
	// Disk digest extent file.
	VirtualMachineFileLayoutExFileTypeDigestExtent = VirtualMachineFileLayoutExFileType("digestExtent")
	// Host based replicated disk persistent state (psf) file.
	VirtualMachineFileLayoutExFileTypeDiskReplicationState = VirtualMachineFileLayoutExFileType("diskReplicationState")
	// Log (log) file.
	VirtualMachineFileLayoutExFileTypeLog = VirtualMachineFileLayoutExFileType("log")
	// Virtual machine statistics (stat) file.
	VirtualMachineFileLayoutExFileTypeStat = VirtualMachineFileLayoutExFileType("stat")
	// Namespace data file.
	VirtualMachineFileLayoutExFileTypeNamespaceData = VirtualMachineFileLayoutExFileType("namespaceData")
	// DataSets disk mode store (dsd) file.
	VirtualMachineFileLayoutExFileTypeDataSetsDiskModeStore = VirtualMachineFileLayoutExFileType("dataSetsDiskModeStore")
	// DataSets vm mode store (dsv) file.
	VirtualMachineFileLayoutExFileTypeDataSetsVmModeStore = VirtualMachineFileLayoutExFileType("dataSetsVmModeStore")
	// Non-volatile RAM (nvram) file.
	VirtualMachineFileLayoutExFileTypeNvram = VirtualMachineFileLayoutExFileType("nvram")
	// Snapshot data (vmsn) file.
	VirtualMachineFileLayoutExFileTypeSnapshotData = VirtualMachineFileLayoutExFileType("snapshotData")
	// Snapshot memory (vmem) file.
	VirtualMachineFileLayoutExFileTypeSnapshotMemory = VirtualMachineFileLayoutExFileType("snapshotMemory")
	// Snapshot metadata (vmsd) file.
	VirtualMachineFileLayoutExFileTypeSnapshotList = VirtualMachineFileLayoutExFileType("snapshotList")
	// Snapshot manifest metadata (-aux.xml) file.
	//
	// This file is still being created but is no longer necessary since
	// the manifest metadata is now available in the snapshot metadata
	// (vmsd) file in vSphere 5.0. This type will be deprecated when
	// vSphere 4.1 is no longer supported.
	VirtualMachineFileLayoutExFileTypeSnapshotManifestList = VirtualMachineFileLayoutExFileType("snapshotManifestList")
	// Suspend (vmss) file.
	VirtualMachineFileLayoutExFileTypeSuspend = VirtualMachineFileLayoutExFileType("suspend")
	// Suspend (vmem) file.
	VirtualMachineFileLayoutExFileTypeSuspendMemory = VirtualMachineFileLayoutExFileType("suspendMemory")
	// Swap (vswp) file.
	VirtualMachineFileLayoutExFileTypeSwap = VirtualMachineFileLayoutExFileType("swap")
	// File generated by VMware ESX kernel for a running virtual
	// machine.
	VirtualMachineFileLayoutExFileTypeUwswap = VirtualMachineFileLayoutExFileType("uwswap")
	// Core (core) file.
	VirtualMachineFileLayoutExFileTypeCore = VirtualMachineFileLayoutExFileType("core")
	// Screenshot file.
	VirtualMachineFileLayoutExFileTypeScreenshot = VirtualMachineFileLayoutExFileType("screenshot")
	// Fault Tolerance metadata file.
	VirtualMachineFileLayoutExFileTypeFtMetadata = VirtualMachineFileLayoutExFileType("ftMetadata")
	// Guest image customization file.
	VirtualMachineFileLayoutExFileTypeGuestCustomization = VirtualMachineFileLayoutExFileType("guestCustomization")
)

func (e VirtualMachineFileLayoutExFileType) Values() []VirtualMachineFileLayoutExFileType {
	return []VirtualMachineFileLayoutExFileType{
		VirtualMachineFileLayoutExFileTypeConfig,
		VirtualMachineFileLayoutExFileTypeExtendedConfig,
		VirtualMachineFileLayoutExFileTypeDiskDescriptor,
		VirtualMachineFileLayoutExFileTypeDiskExtent,
		VirtualMachineFileLayoutExFileTypeDigestDescriptor,
		VirtualMachineFileLayoutExFileTypeDigestExtent,
		VirtualMachineFileLayoutExFileTypeDiskReplicationState,
		VirtualMachineFileLayoutExFileTypeLog,
		VirtualMachineFileLayoutExFileTypeStat,
		VirtualMachineFileLayoutExFileTypeNamespaceData,
		VirtualMachineFileLayoutExFileTypeDataSetsDiskModeStore,
		VirtualMachineFileLayoutExFileTypeDataSetsVmModeStore,
		VirtualMachineFileLayoutExFileTypeNvram,
		VirtualMachineFileLayoutExFileTypeSnapshotData,
		VirtualMachineFileLayoutExFileTypeSnapshotMemory,
		VirtualMachineFileLayoutExFileTypeSnapshotList,
		VirtualMachineFileLayoutExFileTypeSnapshotManifestList,
		VirtualMachineFileLayoutExFileTypeSuspend,
		VirtualMachineFileLayoutExFileTypeSuspendMemory,
		VirtualMachineFileLayoutExFileTypeSwap,
		VirtualMachineFileLayoutExFileTypeUwswap,
		VirtualMachineFileLayoutExFileTypeCore,
		VirtualMachineFileLayoutExFileTypeScreenshot,
		VirtualMachineFileLayoutExFileTypeFtMetadata,
		VirtualMachineFileLayoutExFileTypeGuestCustomization,
	}
}

func (e VirtualMachineFileLayoutExFileType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineFileLayoutExFileType"] = reflect.TypeOf((*VirtualMachineFileLayoutExFileType)(nil)).Elem()
	minAPIVersionForEnumValue["VirtualMachineFileLayoutExFileType"] = map[string]string{
		"dataSetsDiskModeStore": "8.0.0.0",
		"dataSetsVmModeStore":   "8.0.0.0",
	}
}

// Set of possible values for `VirtualMachineFlagInfo.monitorType`.
type VirtualMachineFlagInfoMonitorType string

const (
	// Run vmx in default mode, matching the build type of vmkernel.
	VirtualMachineFlagInfoMonitorTypeRelease = VirtualMachineFlagInfoMonitorType("release")
	// Run vmx in debug mode.
	VirtualMachineFlagInfoMonitorTypeDebug = VirtualMachineFlagInfoMonitorType("debug")
	// Run vmx in stats mode.
	VirtualMachineFlagInfoMonitorTypeStats = VirtualMachineFlagInfoMonitorType("stats")
)

func (e VirtualMachineFlagInfoMonitorType) Values() []VirtualMachineFlagInfoMonitorType {
	return []VirtualMachineFlagInfoMonitorType{
		VirtualMachineFlagInfoMonitorTypeRelease,
		VirtualMachineFlagInfoMonitorTypeDebug,
		VirtualMachineFlagInfoMonitorTypeStats,
	}
}

func (e VirtualMachineFlagInfoMonitorType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineFlagInfoMonitorType"] = reflect.TypeOf((*VirtualMachineFlagInfoMonitorType)(nil)).Elem()
}

// Set of possible values for `VirtualMachineFlagInfo.virtualExecUsage`.
type VirtualMachineFlagInfoVirtualExecUsage string

const (
	// Determine automatically whether to use hardware virtualization (HV) support.
	VirtualMachineFlagInfoVirtualExecUsageHvAuto = VirtualMachineFlagInfoVirtualExecUsage("hvAuto")
	// Use hardware virtualization (HV) support if the physical hardware supports it.
	VirtualMachineFlagInfoVirtualExecUsageHvOn = VirtualMachineFlagInfoVirtualExecUsage("hvOn")
	// Do not use hardware virtualization (HV) support.
	VirtualMachineFlagInfoVirtualExecUsageHvOff = VirtualMachineFlagInfoVirtualExecUsage("hvOff")
)

func (e VirtualMachineFlagInfoVirtualExecUsage) Values() []VirtualMachineFlagInfoVirtualExecUsage {
	return []VirtualMachineFlagInfoVirtualExecUsage{
		VirtualMachineFlagInfoVirtualExecUsageHvAuto,
		VirtualMachineFlagInfoVirtualExecUsageHvOn,
		VirtualMachineFlagInfoVirtualExecUsageHvOff,
	}
}

func (e VirtualMachineFlagInfoVirtualExecUsage) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineFlagInfoVirtualExecUsage"] = reflect.TypeOf((*VirtualMachineFlagInfoVirtualExecUsage)(nil)).Elem()
}

// Set of possible values for `VirtualMachineFlagInfo.virtualMmuUsage`.
type VirtualMachineFlagInfoVirtualMmuUsage string

const (
	// Determine automatically whether to use nested page table hardware support.
	VirtualMachineFlagInfoVirtualMmuUsageAutomatic = VirtualMachineFlagInfoVirtualMmuUsage("automatic")
	// Use nested paging hardware support if the physical hardware supports it.
	VirtualMachineFlagInfoVirtualMmuUsageOn = VirtualMachineFlagInfoVirtualMmuUsage("on")
	// Do not use nested page table hardware support.
	VirtualMachineFlagInfoVirtualMmuUsageOff = VirtualMachineFlagInfoVirtualMmuUsage("off")
)

func (e VirtualMachineFlagInfoVirtualMmuUsage) Values() []VirtualMachineFlagInfoVirtualMmuUsage {
	return []VirtualMachineFlagInfoVirtualMmuUsage{
		VirtualMachineFlagInfoVirtualMmuUsageAutomatic,
		VirtualMachineFlagInfoVirtualMmuUsageOn,
		VirtualMachineFlagInfoVirtualMmuUsageOff,
	}
}

func (e VirtualMachineFlagInfoVirtualMmuUsage) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineFlagInfoVirtualMmuUsage"] = reflect.TypeOf((*VirtualMachineFlagInfoVirtualMmuUsage)(nil)).Elem()
}

// Fork child type.
//
// A child could be type of none, persistent, or
// nonpersistent.
type VirtualMachineForkConfigInfoChildType string

const (
	// The virtual machine is not a child.
	VirtualMachineForkConfigInfoChildTypeNone = VirtualMachineForkConfigInfoChildType("none")
	// The virtual machine is a persistent child.
	VirtualMachineForkConfigInfoChildTypePersistent = VirtualMachineForkConfigInfoChildType("persistent")
	// The virtual machine is a non-persistent child.
	VirtualMachineForkConfigInfoChildTypeNonpersistent = VirtualMachineForkConfigInfoChildType("nonpersistent")
)

func (e VirtualMachineForkConfigInfoChildType) Values() []VirtualMachineForkConfigInfoChildType {
	return []VirtualMachineForkConfigInfoChildType{
		VirtualMachineForkConfigInfoChildTypeNone,
		VirtualMachineForkConfigInfoChildTypePersistent,
		VirtualMachineForkConfigInfoChildTypeNonpersistent,
	}
}

func (e VirtualMachineForkConfigInfoChildType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineForkConfigInfoChildType"] = reflect.TypeOf((*VirtualMachineForkConfigInfoChildType)(nil)).Elem()
}

// Guest operating system family constants.
type VirtualMachineGuestOsFamily string

const (
	// Windows operating system
	VirtualMachineGuestOsFamilyWindowsGuest = VirtualMachineGuestOsFamily("windowsGuest")
	// Linux operating system
	VirtualMachineGuestOsFamilyLinuxGuest = VirtualMachineGuestOsFamily("linuxGuest")
	// Novell Netware
	VirtualMachineGuestOsFamilyNetwareGuest = VirtualMachineGuestOsFamily("netwareGuest")
	// Solaris operating system
	VirtualMachineGuestOsFamilySolarisGuest = VirtualMachineGuestOsFamily("solarisGuest")
	// Mac OS operating system
	VirtualMachineGuestOsFamilyDarwinGuestFamily = VirtualMachineGuestOsFamily("darwinGuestFamily")
	// Other operating systems
	VirtualMachineGuestOsFamilyOtherGuestFamily = VirtualMachineGuestOsFamily("otherGuestFamily")
)

func (e VirtualMachineGuestOsFamily) Values() []VirtualMachineGuestOsFamily {
	return []VirtualMachineGuestOsFamily{
		VirtualMachineGuestOsFamilyWindowsGuest,
		VirtualMachineGuestOsFamilyLinuxGuest,
		VirtualMachineGuestOsFamilyNetwareGuest,
		VirtualMachineGuestOsFamilySolarisGuest,
		VirtualMachineGuestOsFamilyDarwinGuestFamily,
		VirtualMachineGuestOsFamilyOtherGuestFamily,
	}
}

func (e VirtualMachineGuestOsFamily) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineGuestOsFamily"] = reflect.TypeOf((*VirtualMachineGuestOsFamily)(nil)).Elem()
}

// Guest operating system identifier.
type VirtualMachineGuestOsIdentifier string

const (
	// MS-DOS.
	VirtualMachineGuestOsIdentifierDosGuest = VirtualMachineGuestOsIdentifier("dosGuest")
	// Windows 3.1
	VirtualMachineGuestOsIdentifierWin31Guest = VirtualMachineGuestOsIdentifier("win31Guest")
	// Windows 95
	VirtualMachineGuestOsIdentifierWin95Guest = VirtualMachineGuestOsIdentifier("win95Guest")
	// Windows 98
	VirtualMachineGuestOsIdentifierWin98Guest = VirtualMachineGuestOsIdentifier("win98Guest")
	// Windows Millennium Edition
	VirtualMachineGuestOsIdentifierWinMeGuest = VirtualMachineGuestOsIdentifier("winMeGuest")
	// Windows NT 4
	VirtualMachineGuestOsIdentifierWinNTGuest = VirtualMachineGuestOsIdentifier("winNTGuest")
	// Windows 2000 Professional
	VirtualMachineGuestOsIdentifierWin2000ProGuest = VirtualMachineGuestOsIdentifier("win2000ProGuest")
	// Windows 2000 Server
	VirtualMachineGuestOsIdentifierWin2000ServGuest = VirtualMachineGuestOsIdentifier("win2000ServGuest")
	// Windows 2000 Advanced Server
	VirtualMachineGuestOsIdentifierWin2000AdvServGuest = VirtualMachineGuestOsIdentifier("win2000AdvServGuest")
	// Windows XP Home Edition
	VirtualMachineGuestOsIdentifierWinXPHomeGuest = VirtualMachineGuestOsIdentifier("winXPHomeGuest")
	// Windows XP Professional
	VirtualMachineGuestOsIdentifierWinXPProGuest = VirtualMachineGuestOsIdentifier("winXPProGuest")
	// Windows XP Professional Edition (64 bit)
	VirtualMachineGuestOsIdentifierWinXPPro64Guest = VirtualMachineGuestOsIdentifier("winXPPro64Guest")
	// Windows Server 2003, Web Edition
	VirtualMachineGuestOsIdentifierWinNetWebGuest = VirtualMachineGuestOsIdentifier("winNetWebGuest")
	// Windows Server 2003, Standard Edition
	VirtualMachineGuestOsIdentifierWinNetStandardGuest = VirtualMachineGuestOsIdentifier("winNetStandardGuest")
	// Windows Server 2003, Enterprise Edition
	VirtualMachineGuestOsIdentifierWinNetEnterpriseGuest = VirtualMachineGuestOsIdentifier("winNetEnterpriseGuest")
	// Windows Server 2003, Datacenter Edition
	VirtualMachineGuestOsIdentifierWinNetDatacenterGuest = VirtualMachineGuestOsIdentifier("winNetDatacenterGuest")
	// Windows Small Business Server 2003
	VirtualMachineGuestOsIdentifierWinNetBusinessGuest = VirtualMachineGuestOsIdentifier("winNetBusinessGuest")
	// Windows Server 2003, Standard Edition (64 bit)
	VirtualMachineGuestOsIdentifierWinNetStandard64Guest = VirtualMachineGuestOsIdentifier("winNetStandard64Guest")
	// Windows Server 2003, Enterprise Edition (64 bit)
	VirtualMachineGuestOsIdentifierWinNetEnterprise64Guest = VirtualMachineGuestOsIdentifier("winNetEnterprise64Guest")
	// Windows Longhorn
	VirtualMachineGuestOsIdentifierWinLonghornGuest = VirtualMachineGuestOsIdentifier("winLonghornGuest")
	// Windows Longhorn (64 bit)
	VirtualMachineGuestOsIdentifierWinLonghorn64Guest = VirtualMachineGuestOsIdentifier("winLonghorn64Guest")
	// Windows Server 2003, Datacenter Edition (64 bit)
	VirtualMachineGuestOsIdentifierWinNetDatacenter64Guest = VirtualMachineGuestOsIdentifier("winNetDatacenter64Guest")
	// Windows Vista
	VirtualMachineGuestOsIdentifierWinVistaGuest = VirtualMachineGuestOsIdentifier("winVistaGuest")
	// Windows Vista (64 bit)
	VirtualMachineGuestOsIdentifierWinVista64Guest = VirtualMachineGuestOsIdentifier("winVista64Guest")
	// Windows 7
	VirtualMachineGuestOsIdentifierWindows7Guest = VirtualMachineGuestOsIdentifier("windows7Guest")
	// Windows 7 (64 bit)
	VirtualMachineGuestOsIdentifierWindows7_64Guest = VirtualMachineGuestOsIdentifier("windows7_64Guest")
	// Windows Server 2008 R2 (64 bit)
	VirtualMachineGuestOsIdentifierWindows7Server64Guest = VirtualMachineGuestOsIdentifier("windows7Server64Guest")
	// Windows 8
	VirtualMachineGuestOsIdentifierWindows8Guest = VirtualMachineGuestOsIdentifier("windows8Guest")
	// Windows 8 (64 bit)
	VirtualMachineGuestOsIdentifierWindows8_64Guest = VirtualMachineGuestOsIdentifier("windows8_64Guest")
	// Windows 8 Server (64 bit)
	VirtualMachineGuestOsIdentifierWindows8Server64Guest = VirtualMachineGuestOsIdentifier("windows8Server64Guest")
	// Windows 10
	VirtualMachineGuestOsIdentifierWindows9Guest = VirtualMachineGuestOsIdentifier("windows9Guest")
	// Windows 10 (64 bit)
	VirtualMachineGuestOsIdentifierWindows9_64Guest = VirtualMachineGuestOsIdentifier("windows9_64Guest")
	// Windows 10 Server (64 bit)
	VirtualMachineGuestOsIdentifierWindows9Server64Guest = VirtualMachineGuestOsIdentifier("windows9Server64Guest")
	// Windows 11
	VirtualMachineGuestOsIdentifierWindows11_64Guest = VirtualMachineGuestOsIdentifier("windows11_64Guest")
	// Windows 12
	VirtualMachineGuestOsIdentifierWindows12_64Guest = VirtualMachineGuestOsIdentifier("windows12_64Guest")
	// Windows Hyper-V
	VirtualMachineGuestOsIdentifierWindowsHyperVGuest = VirtualMachineGuestOsIdentifier("windowsHyperVGuest")
	// Windows Server 2019
	VirtualMachineGuestOsIdentifierWindows2019srv_64Guest = VirtualMachineGuestOsIdentifier("windows2019srv_64Guest")
	// Windows Server 2022
	VirtualMachineGuestOsIdentifierWindows2019srvNext_64Guest = VirtualMachineGuestOsIdentifier("windows2019srvNext_64Guest")
	// Windows Server 2025
	VirtualMachineGuestOsIdentifierWindows2022srvNext_64Guest = VirtualMachineGuestOsIdentifier("windows2022srvNext_64Guest")
	// FreeBSD
	VirtualMachineGuestOsIdentifierFreebsdGuest = VirtualMachineGuestOsIdentifier("freebsdGuest")
	// FreeBSD x64
	VirtualMachineGuestOsIdentifierFreebsd64Guest = VirtualMachineGuestOsIdentifier("freebsd64Guest")
	// FreeBSD 11
	VirtualMachineGuestOsIdentifierFreebsd11Guest = VirtualMachineGuestOsIdentifier("freebsd11Guest")
	// FreeBSD 11 x64
	VirtualMachineGuestOsIdentifierFreebsd11_64Guest = VirtualMachineGuestOsIdentifier("freebsd11_64Guest")
	// FreeBSD 12
	VirtualMachineGuestOsIdentifierFreebsd12Guest = VirtualMachineGuestOsIdentifier("freebsd12Guest")
	// FreeBSD 12 x64
	VirtualMachineGuestOsIdentifierFreebsd12_64Guest = VirtualMachineGuestOsIdentifier("freebsd12_64Guest")
	// FreeBSD 13
	VirtualMachineGuestOsIdentifierFreebsd13Guest = VirtualMachineGuestOsIdentifier("freebsd13Guest")
	// FreeBSD 13 x64
	VirtualMachineGuestOsIdentifierFreebsd13_64Guest = VirtualMachineGuestOsIdentifier("freebsd13_64Guest")
	// FreeBSD 14
	VirtualMachineGuestOsIdentifierFreebsd14Guest = VirtualMachineGuestOsIdentifier("freebsd14Guest")
	// FreeBSD 14 x64
	VirtualMachineGuestOsIdentifierFreebsd14_64Guest = VirtualMachineGuestOsIdentifier("freebsd14_64Guest")
	// Red Hat Linux 2.1
	VirtualMachineGuestOsIdentifierRedhatGuest = VirtualMachineGuestOsIdentifier("redhatGuest")
	// Red Hat Enterprise Linux 2
	VirtualMachineGuestOsIdentifierRhel2Guest = VirtualMachineGuestOsIdentifier("rhel2Guest")
	// Red Hat Enterprise Linux 3
	VirtualMachineGuestOsIdentifierRhel3Guest = VirtualMachineGuestOsIdentifier("rhel3Guest")
	// Red Hat Enterprise Linux 3 (64 bit)
	VirtualMachineGuestOsIdentifierRhel3_64Guest = VirtualMachineGuestOsIdentifier("rhel3_64Guest")
	// Red Hat Enterprise Linux 4
	VirtualMachineGuestOsIdentifierRhel4Guest = VirtualMachineGuestOsIdentifier("rhel4Guest")
	// Red Hat Enterprise Linux 4 (64 bit)
	VirtualMachineGuestOsIdentifierRhel4_64Guest = VirtualMachineGuestOsIdentifier("rhel4_64Guest")
	// Red Hat Enterprise Linux 5
	VirtualMachineGuestOsIdentifierRhel5Guest = VirtualMachineGuestOsIdentifier("rhel5Guest")
	// Red Hat Enterprise Linux 5 (64 bit)
	VirtualMachineGuestOsIdentifierRhel5_64Guest = VirtualMachineGuestOsIdentifier("rhel5_64Guest")
	// Red Hat Enterprise Linux 6
	VirtualMachineGuestOsIdentifierRhel6Guest = VirtualMachineGuestOsIdentifier("rhel6Guest")
	// Red Hat Enterprise Linux 6 (64 bit)
	VirtualMachineGuestOsIdentifierRhel6_64Guest = VirtualMachineGuestOsIdentifier("rhel6_64Guest")
	// Red Hat Enterprise Linux 7
	VirtualMachineGuestOsIdentifierRhel7Guest = VirtualMachineGuestOsIdentifier("rhel7Guest")
	// Red Hat Enterprise Linux 7 (64 bit)
	VirtualMachineGuestOsIdentifierRhel7_64Guest = VirtualMachineGuestOsIdentifier("rhel7_64Guest")
	// Red Hat Enterprise Linux 8 (64 bit)
	VirtualMachineGuestOsIdentifierRhel8_64Guest = VirtualMachineGuestOsIdentifier("rhel8_64Guest")
	// Red Hat Enterprise Linux 9 (64 bit)
	VirtualMachineGuestOsIdentifierRhel9_64Guest = VirtualMachineGuestOsIdentifier("rhel9_64Guest")
	// CentOS 4/5
	VirtualMachineGuestOsIdentifierCentosGuest = VirtualMachineGuestOsIdentifier("centosGuest")
	// CentOS 4/5 (64-bit)
	VirtualMachineGuestOsIdentifierCentos64Guest = VirtualMachineGuestOsIdentifier("centos64Guest")
	// CentOS 6
	VirtualMachineGuestOsIdentifierCentos6Guest = VirtualMachineGuestOsIdentifier("centos6Guest")
	// CentOS 6 (64-bit)
	VirtualMachineGuestOsIdentifierCentos6_64Guest = VirtualMachineGuestOsIdentifier("centos6_64Guest")
	// CentOS 7
	VirtualMachineGuestOsIdentifierCentos7Guest = VirtualMachineGuestOsIdentifier("centos7Guest")
	// CentOS 7 (64-bit)
	VirtualMachineGuestOsIdentifierCentos7_64Guest = VirtualMachineGuestOsIdentifier("centos7_64Guest")
	// CentOS 8 (64-bit)
	VirtualMachineGuestOsIdentifierCentos8_64Guest = VirtualMachineGuestOsIdentifier("centos8_64Guest")
	// CentOS 9 (64-bit)
	VirtualMachineGuestOsIdentifierCentos9_64Guest = VirtualMachineGuestOsIdentifier("centos9_64Guest")
	// Oracle Linux 4/5
	VirtualMachineGuestOsIdentifierOracleLinuxGuest = VirtualMachineGuestOsIdentifier("oracleLinuxGuest")
	// Oracle Linux 4/5 (64-bit)
	VirtualMachineGuestOsIdentifierOracleLinux64Guest = VirtualMachineGuestOsIdentifier("oracleLinux64Guest")
	// Oracle 6
	VirtualMachineGuestOsIdentifierOracleLinux6Guest = VirtualMachineGuestOsIdentifier("oracleLinux6Guest")
	// Oracle 6 (64-bit)
	VirtualMachineGuestOsIdentifierOracleLinux6_64Guest = VirtualMachineGuestOsIdentifier("oracleLinux6_64Guest")
	// Oracle 7
	VirtualMachineGuestOsIdentifierOracleLinux7Guest = VirtualMachineGuestOsIdentifier("oracleLinux7Guest")
	// Oracle 7 (64-bit)
	VirtualMachineGuestOsIdentifierOracleLinux7_64Guest = VirtualMachineGuestOsIdentifier("oracleLinux7_64Guest")
	// Oracle 8 (64-bit)
	VirtualMachineGuestOsIdentifierOracleLinux8_64Guest = VirtualMachineGuestOsIdentifier("oracleLinux8_64Guest")
	// Oracle 9 (64-bit)
	VirtualMachineGuestOsIdentifierOracleLinux9_64Guest = VirtualMachineGuestOsIdentifier("oracleLinux9_64Guest")
	// Suse Linux
	VirtualMachineGuestOsIdentifierSuseGuest = VirtualMachineGuestOsIdentifier("suseGuest")
	// Suse Linux (64 bit)
	VirtualMachineGuestOsIdentifierSuse64Guest = VirtualMachineGuestOsIdentifier("suse64Guest")
	// Suse Linux Enterprise Server 9
	VirtualMachineGuestOsIdentifierSlesGuest = VirtualMachineGuestOsIdentifier("slesGuest")
	// Suse Linux Enterprise Server 9 (64 bit)
	VirtualMachineGuestOsIdentifierSles64Guest = VirtualMachineGuestOsIdentifier("sles64Guest")
	// Suse linux Enterprise Server 10
	VirtualMachineGuestOsIdentifierSles10Guest = VirtualMachineGuestOsIdentifier("sles10Guest")
	// Suse Linux Enterprise Server 10 (64 bit)
	VirtualMachineGuestOsIdentifierSles10_64Guest = VirtualMachineGuestOsIdentifier("sles10_64Guest")
	// Suse linux Enterprise Server 11
	VirtualMachineGuestOsIdentifierSles11Guest = VirtualMachineGuestOsIdentifier("sles11Guest")
	// Suse Linux Enterprise Server 11 (64 bit)
	VirtualMachineGuestOsIdentifierSles11_64Guest = VirtualMachineGuestOsIdentifier("sles11_64Guest")
	// Suse linux Enterprise Server 12
	VirtualMachineGuestOsIdentifierSles12Guest = VirtualMachineGuestOsIdentifier("sles12Guest")
	// Suse Linux Enterprise Server 12 (64 bit)
	VirtualMachineGuestOsIdentifierSles12_64Guest = VirtualMachineGuestOsIdentifier("sles12_64Guest")
	// Suse Linux Enterprise Server 15 (64 bit)
	VirtualMachineGuestOsIdentifierSles15_64Guest = VirtualMachineGuestOsIdentifier("sles15_64Guest")
	// Suse Linux Enterprise Server 16 (64 bit)
	VirtualMachineGuestOsIdentifierSles16_64Guest = VirtualMachineGuestOsIdentifier("sles16_64Guest")
	// Novell Linux Desktop 9
	VirtualMachineGuestOsIdentifierNld9Guest = VirtualMachineGuestOsIdentifier("nld9Guest")
	// Open Enterprise Server
	VirtualMachineGuestOsIdentifierOesGuest = VirtualMachineGuestOsIdentifier("oesGuest")
	// Sun Java Desktop System
	VirtualMachineGuestOsIdentifierSjdsGuest = VirtualMachineGuestOsIdentifier("sjdsGuest")
	// Mandrake Linux
	VirtualMachineGuestOsIdentifierMandrakeGuest = VirtualMachineGuestOsIdentifier("mandrakeGuest")
	// Mandriva Linux
	VirtualMachineGuestOsIdentifierMandrivaGuest = VirtualMachineGuestOsIdentifier("mandrivaGuest")
	// Mandriva Linux (64 bit)
	VirtualMachineGuestOsIdentifierMandriva64Guest = VirtualMachineGuestOsIdentifier("mandriva64Guest")
	// Turbolinux
	VirtualMachineGuestOsIdentifierTurboLinuxGuest = VirtualMachineGuestOsIdentifier("turboLinuxGuest")
	// Turbolinux (64 bit)
	VirtualMachineGuestOsIdentifierTurboLinux64Guest = VirtualMachineGuestOsIdentifier("turboLinux64Guest")
	// Ubuntu Linux
	VirtualMachineGuestOsIdentifierUbuntuGuest = VirtualMachineGuestOsIdentifier("ubuntuGuest")
	// Ubuntu Linux (64 bit)
	VirtualMachineGuestOsIdentifierUbuntu64Guest = VirtualMachineGuestOsIdentifier("ubuntu64Guest")
	// Debian GNU/Linux 4
	VirtualMachineGuestOsIdentifierDebian4Guest = VirtualMachineGuestOsIdentifier("debian4Guest")
	// Debian GNU/Linux 4 (64 bit)
	VirtualMachineGuestOsIdentifierDebian4_64Guest = VirtualMachineGuestOsIdentifier("debian4_64Guest")
	// Debian GNU/Linux 5
	VirtualMachineGuestOsIdentifierDebian5Guest = VirtualMachineGuestOsIdentifier("debian5Guest")
	// Debian GNU/Linux 5 (64 bit)
	VirtualMachineGuestOsIdentifierDebian5_64Guest = VirtualMachineGuestOsIdentifier("debian5_64Guest")
	// Debian GNU/Linux 6
	VirtualMachineGuestOsIdentifierDebian6Guest = VirtualMachineGuestOsIdentifier("debian6Guest")
	// Debian GNU/Linux 6 (64 bit)
	VirtualMachineGuestOsIdentifierDebian6_64Guest = VirtualMachineGuestOsIdentifier("debian6_64Guest")
	// Debian GNU/Linux 7
	VirtualMachineGuestOsIdentifierDebian7Guest = VirtualMachineGuestOsIdentifier("debian7Guest")
	// Debian GNU/Linux 7 (64 bit)
	VirtualMachineGuestOsIdentifierDebian7_64Guest = VirtualMachineGuestOsIdentifier("debian7_64Guest")
	// Debian GNU/Linux 8
	VirtualMachineGuestOsIdentifierDebian8Guest = VirtualMachineGuestOsIdentifier("debian8Guest")
	// Debian GNU/Linux 8 (64 bit)
	VirtualMachineGuestOsIdentifierDebian8_64Guest = VirtualMachineGuestOsIdentifier("debian8_64Guest")
	// Debian GNU/Linux 9
	VirtualMachineGuestOsIdentifierDebian9Guest = VirtualMachineGuestOsIdentifier("debian9Guest")
	// Debian GNU/Linux 9 (64 bit)
	VirtualMachineGuestOsIdentifierDebian9_64Guest = VirtualMachineGuestOsIdentifier("debian9_64Guest")
	// Debian GNU/Linux 10
	VirtualMachineGuestOsIdentifierDebian10Guest = VirtualMachineGuestOsIdentifier("debian10Guest")
	// Debian GNU/Linux 10 (64 bit)
	VirtualMachineGuestOsIdentifierDebian10_64Guest = VirtualMachineGuestOsIdentifier("debian10_64Guest")
	// Debian GNU/Linux 11
	VirtualMachineGuestOsIdentifierDebian11Guest = VirtualMachineGuestOsIdentifier("debian11Guest")
	// Debian GNU/Linux 11 (64 bit)
	VirtualMachineGuestOsIdentifierDebian11_64Guest = VirtualMachineGuestOsIdentifier("debian11_64Guest")
	// Debian GNU/Linux 12
	VirtualMachineGuestOsIdentifierDebian12Guest = VirtualMachineGuestOsIdentifier("debian12Guest")
	// Debian GNU/Linux 12 (64 bit)
	VirtualMachineGuestOsIdentifierDebian12_64Guest = VirtualMachineGuestOsIdentifier("debian12_64Guest")
	// Asianux Server 3
	VirtualMachineGuestOsIdentifierAsianux3Guest = VirtualMachineGuestOsIdentifier("asianux3Guest")
	// Asianux Server 3 (64 bit)
	VirtualMachineGuestOsIdentifierAsianux3_64Guest = VirtualMachineGuestOsIdentifier("asianux3_64Guest")
	// Asianux Server 4
	VirtualMachineGuestOsIdentifierAsianux4Guest = VirtualMachineGuestOsIdentifier("asianux4Guest")
	// Asianux Server 4 (64 bit)
	VirtualMachineGuestOsIdentifierAsianux4_64Guest = VirtualMachineGuestOsIdentifier("asianux4_64Guest")
	// Asianux Server 5 (64 bit)
	VirtualMachineGuestOsIdentifierAsianux5_64Guest = VirtualMachineGuestOsIdentifier("asianux5_64Guest")
	// Asianux Server 7 (64 bit)
	VirtualMachineGuestOsIdentifierAsianux7_64Guest = VirtualMachineGuestOsIdentifier("asianux7_64Guest")
	// Asianux Server 8 (64 bit)
	VirtualMachineGuestOsIdentifierAsianux8_64Guest = VirtualMachineGuestOsIdentifier("asianux8_64Guest")
	// Asianux Server 9 (64 bit)
	VirtualMachineGuestOsIdentifierAsianux9_64Guest = VirtualMachineGuestOsIdentifier("asianux9_64Guest")
	// OpenSUSE Linux
	VirtualMachineGuestOsIdentifierOpensuseGuest = VirtualMachineGuestOsIdentifier("opensuseGuest")
	// OpenSUSE Linux (64 bit)
	VirtualMachineGuestOsIdentifierOpensuse64Guest = VirtualMachineGuestOsIdentifier("opensuse64Guest")
	// Fedora Linux
	VirtualMachineGuestOsIdentifierFedoraGuest = VirtualMachineGuestOsIdentifier("fedoraGuest")
	// Fedora Linux (64 bit)
	VirtualMachineGuestOsIdentifierFedora64Guest = VirtualMachineGuestOsIdentifier("fedora64Guest")
	// CoreOS Linux (64 bit)
	VirtualMachineGuestOsIdentifierCoreos64Guest = VirtualMachineGuestOsIdentifier("coreos64Guest")
	// VMware Photon (64 bit)
	VirtualMachineGuestOsIdentifierVmwarePhoton64Guest = VirtualMachineGuestOsIdentifier("vmwarePhoton64Guest")
	// Linux 2.4x Kernel
	VirtualMachineGuestOsIdentifierOther24xLinuxGuest = VirtualMachineGuestOsIdentifier("other24xLinuxGuest")
	// Linux 2.6x Kernel
	VirtualMachineGuestOsIdentifierOther26xLinuxGuest = VirtualMachineGuestOsIdentifier("other26xLinuxGuest")
	// Linux 2.2x Kernel
	VirtualMachineGuestOsIdentifierOtherLinuxGuest = VirtualMachineGuestOsIdentifier("otherLinuxGuest")
	// Linux 3.x Kernel
	VirtualMachineGuestOsIdentifierOther3xLinuxGuest = VirtualMachineGuestOsIdentifier("other3xLinuxGuest")
	// Linux 4.x Kernel
	VirtualMachineGuestOsIdentifierOther4xLinuxGuest = VirtualMachineGuestOsIdentifier("other4xLinuxGuest")
	// Linux 5.x Kernel
	VirtualMachineGuestOsIdentifierOther5xLinuxGuest = VirtualMachineGuestOsIdentifier("other5xLinuxGuest")
	// Linux 6.x Kernel
	VirtualMachineGuestOsIdentifierOther6xLinuxGuest = VirtualMachineGuestOsIdentifier("other6xLinuxGuest")
	// Other Linux
	VirtualMachineGuestOsIdentifierGenericLinuxGuest = VirtualMachineGuestOsIdentifier("genericLinuxGuest")
	// Linux 2.4.x Kernel (64 bit)
	VirtualMachineGuestOsIdentifierOther24xLinux64Guest = VirtualMachineGuestOsIdentifier("other24xLinux64Guest")
	// Linux 2.6.x Kernel (64 bit)
	VirtualMachineGuestOsIdentifierOther26xLinux64Guest = VirtualMachineGuestOsIdentifier("other26xLinux64Guest")
	// Linux 3.x Kernel (64 bit)
	VirtualMachineGuestOsIdentifierOther3xLinux64Guest = VirtualMachineGuestOsIdentifier("other3xLinux64Guest")
	// Linux 4.x Kernel (64 bit)
	VirtualMachineGuestOsIdentifierOther4xLinux64Guest = VirtualMachineGuestOsIdentifier("other4xLinux64Guest")
	// Linux 5.x Kernel (64 bit)
	VirtualMachineGuestOsIdentifierOther5xLinux64Guest = VirtualMachineGuestOsIdentifier("other5xLinux64Guest")
	// Linux 6.x Kernel (64 bit)
	VirtualMachineGuestOsIdentifierOther6xLinux64Guest = VirtualMachineGuestOsIdentifier("other6xLinux64Guest")
	// Linux (64 bit)
	VirtualMachineGuestOsIdentifierOtherLinux64Guest = VirtualMachineGuestOsIdentifier("otherLinux64Guest")
	// Solaris 6
	VirtualMachineGuestOsIdentifierSolaris6Guest = VirtualMachineGuestOsIdentifier("solaris6Guest")
	// Solaris 7
	VirtualMachineGuestOsIdentifierSolaris7Guest = VirtualMachineGuestOsIdentifier("solaris7Guest")
	// Solaris 8
	VirtualMachineGuestOsIdentifierSolaris8Guest = VirtualMachineGuestOsIdentifier("solaris8Guest")
	// Solaris 9
	VirtualMachineGuestOsIdentifierSolaris9Guest = VirtualMachineGuestOsIdentifier("solaris9Guest")
	// Solaris 10 (32 bit)
	VirtualMachineGuestOsIdentifierSolaris10Guest = VirtualMachineGuestOsIdentifier("solaris10Guest")
	// Solaris 10 (64 bit)
	VirtualMachineGuestOsIdentifierSolaris10_64Guest = VirtualMachineGuestOsIdentifier("solaris10_64Guest")
	// Solaris 11 (64 bit)
	VirtualMachineGuestOsIdentifierSolaris11_64Guest = VirtualMachineGuestOsIdentifier("solaris11_64Guest")
	// OS/2
	VirtualMachineGuestOsIdentifierOs2Guest = VirtualMachineGuestOsIdentifier("os2Guest")
	// eComStation 1.x
	VirtualMachineGuestOsIdentifierEComStationGuest = VirtualMachineGuestOsIdentifier("eComStationGuest")
	// eComStation 2.0
	VirtualMachineGuestOsIdentifierEComStation2Guest = VirtualMachineGuestOsIdentifier("eComStation2Guest")
	// Novell NetWare 4
	VirtualMachineGuestOsIdentifierNetware4Guest = VirtualMachineGuestOsIdentifier("netware4Guest")
	// Novell NetWare 5.1
	VirtualMachineGuestOsIdentifierNetware5Guest = VirtualMachineGuestOsIdentifier("netware5Guest")
	// Novell NetWare 6.x
	VirtualMachineGuestOsIdentifierNetware6Guest = VirtualMachineGuestOsIdentifier("netware6Guest")
	// SCO OpenServer 5
	VirtualMachineGuestOsIdentifierOpenServer5Guest = VirtualMachineGuestOsIdentifier("openServer5Guest")
	// SCO OpenServer 6
	VirtualMachineGuestOsIdentifierOpenServer6Guest = VirtualMachineGuestOsIdentifier("openServer6Guest")
	// SCO UnixWare 7
	VirtualMachineGuestOsIdentifierUnixWare7Guest = VirtualMachineGuestOsIdentifier("unixWare7Guest")
	// Mac OS 10.5
	VirtualMachineGuestOsIdentifierDarwinGuest = VirtualMachineGuestOsIdentifier("darwinGuest")
	// Mac OS 10.5 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin64Guest = VirtualMachineGuestOsIdentifier("darwin64Guest")
	// Mac OS 10.6
	VirtualMachineGuestOsIdentifierDarwin10Guest = VirtualMachineGuestOsIdentifier("darwin10Guest")
	// Mac OS 10.6 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin10_64Guest = VirtualMachineGuestOsIdentifier("darwin10_64Guest")
	// Mac OS 10.7
	VirtualMachineGuestOsIdentifierDarwin11Guest = VirtualMachineGuestOsIdentifier("darwin11Guest")
	// Mac OS 10.7 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin11_64Guest = VirtualMachineGuestOsIdentifier("darwin11_64Guest")
	// Mac OS 10.8 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin12_64Guest = VirtualMachineGuestOsIdentifier("darwin12_64Guest")
	// Mac OS 10.9 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin13_64Guest = VirtualMachineGuestOsIdentifier("darwin13_64Guest")
	// Mac OS 10.10 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin14_64Guest = VirtualMachineGuestOsIdentifier("darwin14_64Guest")
	// Mac OS 10.11 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin15_64Guest = VirtualMachineGuestOsIdentifier("darwin15_64Guest")
	// Mac OS 10.12 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin16_64Guest = VirtualMachineGuestOsIdentifier("darwin16_64Guest")
	// macOS 10.13 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin17_64Guest = VirtualMachineGuestOsIdentifier("darwin17_64Guest")
	// macOS 10.14 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin18_64Guest = VirtualMachineGuestOsIdentifier("darwin18_64Guest")
	// macOS 10.15 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin19_64Guest = VirtualMachineGuestOsIdentifier("darwin19_64Guest")
	// macOS 11 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin20_64Guest = VirtualMachineGuestOsIdentifier("darwin20_64Guest")
	// macOS 12 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin21_64Guest = VirtualMachineGuestOsIdentifier("darwin21_64Guest")
	// macOS 13 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin22_64Guest = VirtualMachineGuestOsIdentifier("darwin22_64Guest")
	// macOS 14 (64 bit)
	VirtualMachineGuestOsIdentifierDarwin23_64Guest = VirtualMachineGuestOsIdentifier("darwin23_64Guest")
	// VMware ESX 4
	VirtualMachineGuestOsIdentifierVmkernelGuest = VirtualMachineGuestOsIdentifier("vmkernelGuest")
	// VMware ESX 5
	VirtualMachineGuestOsIdentifierVmkernel5Guest = VirtualMachineGuestOsIdentifier("vmkernel5Guest")
	// VMware ESX 6
	VirtualMachineGuestOsIdentifierVmkernel6Guest = VirtualMachineGuestOsIdentifier("vmkernel6Guest")
	// VMware ESXi 6.5 AND ESXi 6.7.
	VirtualMachineGuestOsIdentifierVmkernel65Guest = VirtualMachineGuestOsIdentifier("vmkernel65Guest")
	// VMware ESX 7
	VirtualMachineGuestOsIdentifierVmkernel7Guest = VirtualMachineGuestOsIdentifier("vmkernel7Guest")
	// VMware ESX 8
	VirtualMachineGuestOsIdentifierVmkernel8Guest = VirtualMachineGuestOsIdentifier("vmkernel8Guest")
	// Amazon Linux 2 (64 bit)
	VirtualMachineGuestOsIdentifierAmazonlinux2_64Guest = VirtualMachineGuestOsIdentifier("amazonlinux2_64Guest")
	// Amazon Linux 3 (64 bit)
	VirtualMachineGuestOsIdentifierAmazonlinux3_64Guest = VirtualMachineGuestOsIdentifier("amazonlinux3_64Guest")
	// CRX Pod 1
	VirtualMachineGuestOsIdentifierCrxPod1Guest = VirtualMachineGuestOsIdentifier("crxPod1Guest")
	// CRX Sys 1
	VirtualMachineGuestOsIdentifierCrxSys1Guest = VirtualMachineGuestOsIdentifier("crxSys1Guest")
	// Rocky Linux (64-bit)
	VirtualMachineGuestOsIdentifierRockylinux_64Guest = VirtualMachineGuestOsIdentifier("rockylinux_64Guest")
	// AlmaLinux (64-bit)
	VirtualMachineGuestOsIdentifierAlmalinux_64Guest = VirtualMachineGuestOsIdentifier("almalinux_64Guest")
	// Other Operating System
	VirtualMachineGuestOsIdentifierOtherGuest = VirtualMachineGuestOsIdentifier("otherGuest")
	// Other Operating System (64 bit)
	VirtualMachineGuestOsIdentifierOtherGuest64 = VirtualMachineGuestOsIdentifier("otherGuest64")
)

func (e VirtualMachineGuestOsIdentifier) Values() []VirtualMachineGuestOsIdentifier {
	return []VirtualMachineGuestOsIdentifier{
		VirtualMachineGuestOsIdentifierDosGuest,
		VirtualMachineGuestOsIdentifierWin31Guest,
		VirtualMachineGuestOsIdentifierWin95Guest,
		VirtualMachineGuestOsIdentifierWin98Guest,
		VirtualMachineGuestOsIdentifierWinMeGuest,
		VirtualMachineGuestOsIdentifierWinNTGuest,
		VirtualMachineGuestOsIdentifierWin2000ProGuest,
		VirtualMachineGuestOsIdentifierWin2000ServGuest,
		VirtualMachineGuestOsIdentifierWin2000AdvServGuest,
		VirtualMachineGuestOsIdentifierWinXPHomeGuest,
		VirtualMachineGuestOsIdentifierWinXPProGuest,
		VirtualMachineGuestOsIdentifierWinXPPro64Guest,
		VirtualMachineGuestOsIdentifierWinNetWebGuest,
		VirtualMachineGuestOsIdentifierWinNetStandardGuest,
		VirtualMachineGuestOsIdentifierWinNetEnterpriseGuest,
		VirtualMachineGuestOsIdentifierWinNetDatacenterGuest,
		VirtualMachineGuestOsIdentifierWinNetBusinessGuest,
		VirtualMachineGuestOsIdentifierWinNetStandard64Guest,
		VirtualMachineGuestOsIdentifierWinNetEnterprise64Guest,
		VirtualMachineGuestOsIdentifierWinLonghornGuest,
		VirtualMachineGuestOsIdentifierWinLonghorn64Guest,
		VirtualMachineGuestOsIdentifierWinNetDatacenter64Guest,
		VirtualMachineGuestOsIdentifierWinVistaGuest,
		VirtualMachineGuestOsIdentifierWinVista64Guest,
		VirtualMachineGuestOsIdentifierWindows7Guest,
		VirtualMachineGuestOsIdentifierWindows7_64Guest,
		VirtualMachineGuestOsIdentifierWindows7Server64Guest,
		VirtualMachineGuestOsIdentifierWindows8Guest,
		VirtualMachineGuestOsIdentifierWindows8_64Guest,
		VirtualMachineGuestOsIdentifierWindows8Server64Guest,
		VirtualMachineGuestOsIdentifierWindows9Guest,
		VirtualMachineGuestOsIdentifierWindows9_64Guest,
		VirtualMachineGuestOsIdentifierWindows9Server64Guest,
		VirtualMachineGuestOsIdentifierWindows11_64Guest,
		VirtualMachineGuestOsIdentifierWindows12_64Guest,
		VirtualMachineGuestOsIdentifierWindowsHyperVGuest,
		VirtualMachineGuestOsIdentifierWindows2019srv_64Guest,
		VirtualMachineGuestOsIdentifierWindows2019srvNext_64Guest,
		VirtualMachineGuestOsIdentifierWindows2022srvNext_64Guest,
		VirtualMachineGuestOsIdentifierFreebsdGuest,
		VirtualMachineGuestOsIdentifierFreebsd64Guest,
		VirtualMachineGuestOsIdentifierFreebsd11Guest,
		VirtualMachineGuestOsIdentifierFreebsd11_64Guest,
		VirtualMachineGuestOsIdentifierFreebsd12Guest,
		VirtualMachineGuestOsIdentifierFreebsd12_64Guest,
		VirtualMachineGuestOsIdentifierFreebsd13Guest,
		VirtualMachineGuestOsIdentifierFreebsd13_64Guest,
		VirtualMachineGuestOsIdentifierFreebsd14Guest,
		VirtualMachineGuestOsIdentifierFreebsd14_64Guest,
		VirtualMachineGuestOsIdentifierRedhatGuest,
		VirtualMachineGuestOsIdentifierRhel2Guest,
		VirtualMachineGuestOsIdentifierRhel3Guest,
		VirtualMachineGuestOsIdentifierRhel3_64Guest,
		VirtualMachineGuestOsIdentifierRhel4Guest,
		VirtualMachineGuestOsIdentifierRhel4_64Guest,
		VirtualMachineGuestOsIdentifierRhel5Guest,
		VirtualMachineGuestOsIdentifierRhel5_64Guest,
		VirtualMachineGuestOsIdentifierRhel6Guest,
		VirtualMachineGuestOsIdentifierRhel6_64Guest,
		VirtualMachineGuestOsIdentifierRhel7Guest,
		VirtualMachineGuestOsIdentifierRhel7_64Guest,
		VirtualMachineGuestOsIdentifierRhel8_64Guest,
		VirtualMachineGuestOsIdentifierRhel9_64Guest,
		VirtualMachineGuestOsIdentifierCentosGuest,
		VirtualMachineGuestOsIdentifierCentos64Guest,
		VirtualMachineGuestOsIdentifierCentos6Guest,
		VirtualMachineGuestOsIdentifierCentos6_64Guest,
		VirtualMachineGuestOsIdentifierCentos7Guest,
		VirtualMachineGuestOsIdentifierCentos7_64Guest,
		VirtualMachineGuestOsIdentifierCentos8_64Guest,
		VirtualMachineGuestOsIdentifierCentos9_64Guest,
		VirtualMachineGuestOsIdentifierOracleLinuxGuest,
		VirtualMachineGuestOsIdentifierOracleLinux64Guest,
		VirtualMachineGuestOsIdentifierOracleLinux6Guest,
		VirtualMachineGuestOsIdentifierOracleLinux6_64Guest,
		VirtualMachineGuestOsIdentifierOracleLinux7Guest,
		VirtualMachineGuestOsIdentifierOracleLinux7_64Guest,
		VirtualMachineGuestOsIdentifierOracleLinux8_64Guest,
		VirtualMachineGuestOsIdentifierOracleLinux9_64Guest,
		VirtualMachineGuestOsIdentifierSuseGuest,
		VirtualMachineGuestOsIdentifierSuse64Guest,
		VirtualMachineGuestOsIdentifierSlesGuest,
		VirtualMachineGuestOsIdentifierSles64Guest,
		VirtualMachineGuestOsIdentifierSles10Guest,
		VirtualMachineGuestOsIdentifierSles10_64Guest,
		VirtualMachineGuestOsIdentifierSles11Guest,
		VirtualMachineGuestOsIdentifierSles11_64Guest,
		VirtualMachineGuestOsIdentifierSles12Guest,
		VirtualMachineGuestOsIdentifierSles12_64Guest,
		VirtualMachineGuestOsIdentifierSles15_64Guest,
		VirtualMachineGuestOsIdentifierSles16_64Guest,
		VirtualMachineGuestOsIdentifierNld9Guest,
		VirtualMachineGuestOsIdentifierOesGuest,
		VirtualMachineGuestOsIdentifierSjdsGuest,
		VirtualMachineGuestOsIdentifierMandrakeGuest,
		VirtualMachineGuestOsIdentifierMandrivaGuest,
		VirtualMachineGuestOsIdentifierMandriva64Guest,
		VirtualMachineGuestOsIdentifierTurboLinuxGuest,
		VirtualMachineGuestOsIdentifierTurboLinux64Guest,
		VirtualMachineGuestOsIdentifierUbuntuGuest,
		VirtualMachineGuestOsIdentifierUbuntu64Guest,
		VirtualMachineGuestOsIdentifierDebian4Guest,
		VirtualMachineGuestOsIdentifierDebian4_64Guest,
		VirtualMachineGuestOsIdentifierDebian5Guest,
		VirtualMachineGuestOsIdentifierDebian5_64Guest,
		VirtualMachineGuestOsIdentifierDebian6Guest,
		VirtualMachineGuestOsIdentifierDebian6_64Guest,
		VirtualMachineGuestOsIdentifierDebian7Guest,
		VirtualMachineGuestOsIdentifierDebian7_64Guest,
		VirtualMachineGuestOsIdentifierDebian8Guest,
		VirtualMachineGuestOsIdentifierDebian8_64Guest,
		VirtualMachineGuestOsIdentifierDebian9Guest,
		VirtualMachineGuestOsIdentifierDebian9_64Guest,
		VirtualMachineGuestOsIdentifierDebian10Guest,
		VirtualMachineGuestOsIdentifierDebian10_64Guest,
		VirtualMachineGuestOsIdentifierDebian11Guest,
		VirtualMachineGuestOsIdentifierDebian11_64Guest,
		VirtualMachineGuestOsIdentifierDebian12Guest,
		VirtualMachineGuestOsIdentifierDebian12_64Guest,
		VirtualMachineGuestOsIdentifierAsianux3Guest,
		VirtualMachineGuestOsIdentifierAsianux3_64Guest,
		VirtualMachineGuestOsIdentifierAsianux4Guest,
		VirtualMachineGuestOsIdentifierAsianux4_64Guest,
		VirtualMachineGuestOsIdentifierAsianux5_64Guest,
		VirtualMachineGuestOsIdentifierAsianux7_64Guest,
		VirtualMachineGuestOsIdentifierAsianux8_64Guest,
		VirtualMachineGuestOsIdentifierAsianux9_64Guest,
		VirtualMachineGuestOsIdentifierOpensuseGuest,
		VirtualMachineGuestOsIdentifierOpensuse64Guest,
		VirtualMachineGuestOsIdentifierFedoraGuest,
		VirtualMachineGuestOsIdentifierFedora64Guest,
		VirtualMachineGuestOsIdentifierCoreos64Guest,
		VirtualMachineGuestOsIdentifierVmwarePhoton64Guest,
		VirtualMachineGuestOsIdentifierOther24xLinuxGuest,
		VirtualMachineGuestOsIdentifierOther26xLinuxGuest,
		VirtualMachineGuestOsIdentifierOtherLinuxGuest,
		VirtualMachineGuestOsIdentifierOther3xLinuxGuest,
		VirtualMachineGuestOsIdentifierOther4xLinuxGuest,
		VirtualMachineGuestOsIdentifierOther5xLinuxGuest,
		VirtualMachineGuestOsIdentifierOther6xLinuxGuest,
		VirtualMachineGuestOsIdentifierGenericLinuxGuest,
		VirtualMachineGuestOsIdentifierOther24xLinux64Guest,
		VirtualMachineGuestOsIdentifierOther26xLinux64Guest,
		VirtualMachineGuestOsIdentifierOther3xLinux64Guest,
		VirtualMachineGuestOsIdentifierOther4xLinux64Guest,
		VirtualMachineGuestOsIdentifierOther5xLinux64Guest,
		VirtualMachineGuestOsIdentifierOther6xLinux64Guest,
		VirtualMachineGuestOsIdentifierOtherLinux64Guest,
		VirtualMachineGuestOsIdentifierSolaris6Guest,
		VirtualMachineGuestOsIdentifierSolaris7Guest,
		VirtualMachineGuestOsIdentifierSolaris8Guest,
		VirtualMachineGuestOsIdentifierSolaris9Guest,
		VirtualMachineGuestOsIdentifierSolaris10Guest,
		VirtualMachineGuestOsIdentifierSolaris10_64Guest,
		VirtualMachineGuestOsIdentifierSolaris11_64Guest,
		VirtualMachineGuestOsIdentifierOs2Guest,
		VirtualMachineGuestOsIdentifierEComStationGuest,
		VirtualMachineGuestOsIdentifierEComStation2Guest,
		VirtualMachineGuestOsIdentifierNetware4Guest,
		VirtualMachineGuestOsIdentifierNetware5Guest,
		VirtualMachineGuestOsIdentifierNetware6Guest,
		VirtualMachineGuestOsIdentifierOpenServer5Guest,
		VirtualMachineGuestOsIdentifierOpenServer6Guest,
		VirtualMachineGuestOsIdentifierUnixWare7Guest,
		VirtualMachineGuestOsIdentifierDarwinGuest,
		VirtualMachineGuestOsIdentifierDarwin64Guest,
		VirtualMachineGuestOsIdentifierDarwin10Guest,
		VirtualMachineGuestOsIdentifierDarwin10_64Guest,
		VirtualMachineGuestOsIdentifierDarwin11Guest,
		VirtualMachineGuestOsIdentifierDarwin11_64Guest,
		VirtualMachineGuestOsIdentifierDarwin12_64Guest,
		VirtualMachineGuestOsIdentifierDarwin13_64Guest,
		VirtualMachineGuestOsIdentifierDarwin14_64Guest,
		VirtualMachineGuestOsIdentifierDarwin15_64Guest,
		VirtualMachineGuestOsIdentifierDarwin16_64Guest,
		VirtualMachineGuestOsIdentifierDarwin17_64Guest,
		VirtualMachineGuestOsIdentifierDarwin18_64Guest,
		VirtualMachineGuestOsIdentifierDarwin19_64Guest,
		VirtualMachineGuestOsIdentifierDarwin20_64Guest,
		VirtualMachineGuestOsIdentifierDarwin21_64Guest,
		VirtualMachineGuestOsIdentifierDarwin22_64Guest,
		VirtualMachineGuestOsIdentifierDarwin23_64Guest,
		VirtualMachineGuestOsIdentifierVmkernelGuest,
		VirtualMachineGuestOsIdentifierVmkernel5Guest,
		VirtualMachineGuestOsIdentifierVmkernel6Guest,
		VirtualMachineGuestOsIdentifierVmkernel65Guest,
		VirtualMachineGuestOsIdentifierVmkernel7Guest,
		VirtualMachineGuestOsIdentifierVmkernel8Guest,
		VirtualMachineGuestOsIdentifierAmazonlinux2_64Guest,
		VirtualMachineGuestOsIdentifierAmazonlinux3_64Guest,
		VirtualMachineGuestOsIdentifierCrxPod1Guest,
		VirtualMachineGuestOsIdentifierCrxSys1Guest,
		VirtualMachineGuestOsIdentifierRockylinux_64Guest,
		VirtualMachineGuestOsIdentifierAlmalinux_64Guest,
		VirtualMachineGuestOsIdentifierOtherGuest,
		VirtualMachineGuestOsIdentifierOtherGuest64,
	}
}

func (e VirtualMachineGuestOsIdentifier) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineGuestOsIdentifier"] = reflect.TypeOf((*VirtualMachineGuestOsIdentifier)(nil)).Elem()
	minAPIVersionForEnumValue["VirtualMachineGuestOsIdentifier"] = map[string]string{
		"windows11_64Guest":          "8.0.0.1",
		"windows12_64Guest":          "8.0.0.1",
		"windows2019srvNext_64Guest": "7.0.1.0",
		"windows2022srvNext_64Guest": "8.0.0.1",
		"freebsd13Guest":             "7.0.1.0",
		"freebsd13_64Guest":          "7.0.1.0",
		"freebsd14Guest":             "8.0.0.1",
		"freebsd14_64Guest":          "8.0.0.1",
		"rhel9_64Guest":              "7.0.1.0",
		"centos9_64Guest":            "7.0.1.0",
		"oracleLinux9_64Guest":       "7.0.1.0",
		"sles16_64Guest":             "7.0.1.0",
		"debian12Guest":              "8.0.0.1",
		"debian12_64Guest":           "8.0.0.1",
		"asianux9_64Guest":           "7.0.1.0",
		"other5xLinuxGuest":          "7.0.1.0",
		"other6xLinuxGuest":          "8.0.0.1",
		"other5xLinux64Guest":        "7.0.1.0",
		"other6xLinux64Guest":        "8.0.0.1",
		"darwin20_64Guest":           "7.0.1.0",
		"darwin21_64Guest":           "7.0.1.0",
		"darwin22_64Guest":           "8.0.0.1",
		"darwin23_64Guest":           "8.0.0.1",
		"vmkernel8Guest":             "8.0.0.1",
		"amazonlinux3_64Guest":       "7.0.1.0",
		"crxSys1Guest":               "8.0.3.0",
		"rockylinux_64Guest":         "8.0.0.1",
		"almalinux_64Guest":          "8.0.0.1",
	}
}

// The possible hints that the guest could display about current tasks
// inside the guest.
type VirtualMachineGuestState string

const (
	VirtualMachineGuestStateRunning      = VirtualMachineGuestState("running")
	VirtualMachineGuestStateShuttingDown = VirtualMachineGuestState("shuttingDown")
	VirtualMachineGuestStateResetting    = VirtualMachineGuestState("resetting")
	VirtualMachineGuestStateStandby      = VirtualMachineGuestState("standby")
	VirtualMachineGuestStateNotRunning   = VirtualMachineGuestState("notRunning")
	VirtualMachineGuestStateUnknown      = VirtualMachineGuestState("unknown")
)

func (e VirtualMachineGuestState) Values() []VirtualMachineGuestState {
	return []VirtualMachineGuestState{
		VirtualMachineGuestStateRunning,
		VirtualMachineGuestStateShuttingDown,
		VirtualMachineGuestStateResetting,
		VirtualMachineGuestStateStandby,
		VirtualMachineGuestStateNotRunning,
		VirtualMachineGuestStateUnknown,
	}
}

func (e VirtualMachineGuestState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineGuestState"] = reflect.TypeOf((*VirtualMachineGuestState)(nil)).Elem()
}

// Deprecated as of vSphere API 6.7.
//
// Set of possible values for `VirtualMachineFlagInfo.htSharing`.
type VirtualMachineHtSharing string

const (
	// VCPUs may freely share cores at any time with any other
	// VCPUs (default for all virtual machines on a hyperthreaded
	// system).
	VirtualMachineHtSharingAny = VirtualMachineHtSharing("any")
	// VCPUs should not share cores with each other or with VCPUs
	// from other virtual machines.
	//
	// That is, each VCPU from this
	// virtual machine should always get a whole core to itself,
	// with the other logical CPU on that core being placed into
	// the "halted" state.
	VirtualMachineHtSharingNone = VirtualMachineHtSharing("none")
	// Similar to "none", in that VCPUs from this virtual machine
	// will not be allowed to share cores with VCPUs from other
	// virtual machines.
	//
	// However, other VCPUs from the same virtual
	// machine will be allowed to share cores together. This
	// configuration option is only permitted for SMP virtual
	// machines. If applied to a uniprocessor virtual machine, it
	// will be converted to the "none" sharing option.
	VirtualMachineHtSharingInternal = VirtualMachineHtSharing("internal")
)

func (e VirtualMachineHtSharing) Values() []VirtualMachineHtSharing {
	return []VirtualMachineHtSharing{
		VirtualMachineHtSharingAny,
		VirtualMachineHtSharingNone,
		VirtualMachineHtSharingInternal,
	}
}

func (e VirtualMachineHtSharing) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineHtSharing"] = reflect.TypeOf((*VirtualMachineHtSharing)(nil)).Elem()
}

// Means for allocating additional memory for virtual machines.
type VirtualMachineMemoryAllocationPolicy string

const (
	// Fit all virtual machine memory into reserved host memory.
	VirtualMachineMemoryAllocationPolicySwapNone = VirtualMachineMemoryAllocationPolicy("swapNone")
	// Allow some virtual machine memory to be swapped.
	VirtualMachineMemoryAllocationPolicySwapSome = VirtualMachineMemoryAllocationPolicy("swapSome")
	// Allow most virtual machine memory to be swapped.
	VirtualMachineMemoryAllocationPolicySwapMost = VirtualMachineMemoryAllocationPolicy("swapMost")
)

func (e VirtualMachineMemoryAllocationPolicy) Values() []VirtualMachineMemoryAllocationPolicy {
	return []VirtualMachineMemoryAllocationPolicy{
		VirtualMachineMemoryAllocationPolicySwapNone,
		VirtualMachineMemoryAllocationPolicySwapSome,
		VirtualMachineMemoryAllocationPolicySwapMost,
	}
}

func (e VirtualMachineMemoryAllocationPolicy) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineMemoryAllocationPolicy"] = reflect.TypeOf((*VirtualMachineMemoryAllocationPolicy)(nil)).Elem()
}

// This enum represents the set of legal operations
type VirtualMachineMetadataManagerVmMetadataOp string

const (
	// Create or update the Metadata for the specified VM
	VirtualMachineMetadataManagerVmMetadataOpUpdate = VirtualMachineMetadataManagerVmMetadataOp("Update")
	// Remove the Metadata for the specified VM
	VirtualMachineMetadataManagerVmMetadataOpRemove = VirtualMachineMetadataManagerVmMetadataOp("Remove")
)

func (e VirtualMachineMetadataManagerVmMetadataOp) Values() []VirtualMachineMetadataManagerVmMetadataOp {
	return []VirtualMachineMetadataManagerVmMetadataOp{
		VirtualMachineMetadataManagerVmMetadataOpUpdate,
		VirtualMachineMetadataManagerVmMetadataOpRemove,
	}
}

func (e VirtualMachineMetadataManagerVmMetadataOp) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineMetadataManagerVmMetadataOp"] = reflect.TypeOf((*VirtualMachineMetadataManagerVmMetadataOp)(nil)).Elem()
}

// This enum contains a list of valid owner values for
// the name field
type VirtualMachineMetadataManagerVmMetadataOwnerOwner string

const (
	VirtualMachineMetadataManagerVmMetadataOwnerOwnerComVmwareVsphereHA = VirtualMachineMetadataManagerVmMetadataOwnerOwner("ComVmwareVsphereHA")
)

func (e VirtualMachineMetadataManagerVmMetadataOwnerOwner) Values() []VirtualMachineMetadataManagerVmMetadataOwnerOwner {
	return []VirtualMachineMetadataManagerVmMetadataOwnerOwner{
		VirtualMachineMetadataManagerVmMetadataOwnerOwnerComVmwareVsphereHA,
	}
}

func (e VirtualMachineMetadataManagerVmMetadataOwnerOwner) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineMetadataManagerVmMetadataOwnerOwner"] = reflect.TypeOf((*VirtualMachineMetadataManagerVmMetadataOwnerOwner)(nil)).Elem()
}

// MovePriority is an enumeration of values that indicate the priority of the task
// that moves a virtual machine from one host to another or one storage location
// to another.
//
// Note this priority can affect both the source and target hosts.
type VirtualMachineMovePriority string

const (
	// The task of moving this virtual machine is low priority.
	VirtualMachineMovePriorityLowPriority = VirtualMachineMovePriority("lowPriority")
	// The task of moving this virtual machine is high priority.
	VirtualMachineMovePriorityHighPriority = VirtualMachineMovePriority("highPriority")
	// The task of moving this virtual machine is the default priority.
	VirtualMachineMovePriorityDefaultPriority = VirtualMachineMovePriority("defaultPriority")
)

func (e VirtualMachineMovePriority) Values() []VirtualMachineMovePriority {
	return []VirtualMachineMovePriority{
		VirtualMachineMovePriorityLowPriority,
		VirtualMachineMovePriorityHighPriority,
		VirtualMachineMovePriorityDefaultPriority,
	}
}

func (e VirtualMachineMovePriority) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineMovePriority"] = reflect.TypeOf((*VirtualMachineMovePriority)(nil)).Elem()
}

// The NeedSecondaryReason type defines all reasons a virtual machine is
// in the needSecondary Fault Tolerance state following a failure.
type VirtualMachineNeedSecondaryReason string

const (
	// Initializing FT
	VirtualMachineNeedSecondaryReasonInitializing = VirtualMachineNeedSecondaryReason("initializing")
	// Divergence
	VirtualMachineNeedSecondaryReasonDivergence = VirtualMachineNeedSecondaryReason("divergence")
	// Lose connection to secondary
	VirtualMachineNeedSecondaryReasonLostConnection = VirtualMachineNeedSecondaryReason("lostConnection")
	// Partial hardware failure
	VirtualMachineNeedSecondaryReasonPartialHardwareFailure = VirtualMachineNeedSecondaryReason("partialHardwareFailure")
	// Terminated by user
	VirtualMachineNeedSecondaryReasonUserAction = VirtualMachineNeedSecondaryReason("userAction")
	// Checkpoint error
	VirtualMachineNeedSecondaryReasonCheckpointError = VirtualMachineNeedSecondaryReason("checkpointError")
	// All other reasons
	VirtualMachineNeedSecondaryReasonOther = VirtualMachineNeedSecondaryReason("other")
)

func (e VirtualMachineNeedSecondaryReason) Values() []VirtualMachineNeedSecondaryReason {
	return []VirtualMachineNeedSecondaryReason{
		VirtualMachineNeedSecondaryReasonInitializing,
		VirtualMachineNeedSecondaryReasonDivergence,
		VirtualMachineNeedSecondaryReasonLostConnection,
		VirtualMachineNeedSecondaryReasonPartialHardwareFailure,
		VirtualMachineNeedSecondaryReasonUserAction,
		VirtualMachineNeedSecondaryReasonCheckpointError,
		VirtualMachineNeedSecondaryReasonOther,
	}
}

func (e VirtualMachineNeedSecondaryReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineNeedSecondaryReason"] = reflect.TypeOf((*VirtualMachineNeedSecondaryReason)(nil)).Elem()
}

// Set of possible values for `VirtualMachineFlagInfo.snapshotPowerOffBehavior`.
type VirtualMachinePowerOffBehavior string

const (
	// Just power off the virtual machine.
	VirtualMachinePowerOffBehaviorPowerOff = VirtualMachinePowerOffBehavior("powerOff")
	// Revert to the snapshot.
	VirtualMachinePowerOffBehaviorRevert = VirtualMachinePowerOffBehavior("revert")
	// Prompt the user for instructions at power-off time.
	VirtualMachinePowerOffBehaviorPrompt = VirtualMachinePowerOffBehavior("prompt")
	// Take a new snapshot.
	VirtualMachinePowerOffBehaviorTake = VirtualMachinePowerOffBehavior("take")
)

func (e VirtualMachinePowerOffBehavior) Values() []VirtualMachinePowerOffBehavior {
	return []VirtualMachinePowerOffBehavior{
		VirtualMachinePowerOffBehaviorPowerOff,
		VirtualMachinePowerOffBehaviorRevert,
		VirtualMachinePowerOffBehaviorPrompt,
		VirtualMachinePowerOffBehaviorTake,
	}
}

func (e VirtualMachinePowerOffBehavior) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachinePowerOffBehavior"] = reflect.TypeOf((*VirtualMachinePowerOffBehavior)(nil)).Elem()
}

// The list of possible default power operations available for the virtual machine
type VirtualMachinePowerOpType string

const (
	VirtualMachinePowerOpTypeSoft   = VirtualMachinePowerOpType("soft")
	VirtualMachinePowerOpTypeHard   = VirtualMachinePowerOpType("hard")
	VirtualMachinePowerOpTypePreset = VirtualMachinePowerOpType("preset")
)

func (e VirtualMachinePowerOpType) Values() []VirtualMachinePowerOpType {
	return []VirtualMachinePowerOpType{
		VirtualMachinePowerOpTypeSoft,
		VirtualMachinePowerOpTypeHard,
		VirtualMachinePowerOpTypePreset,
	}
}

func (e VirtualMachinePowerOpType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachinePowerOpType"] = reflect.TypeOf((*VirtualMachinePowerOpType)(nil)).Elem()
}

// The PowerState type defines a simple set of states for a virtual machine:
// poweredOn, poweredOff, and suspended.
//
// This type does not model substates,
// such as when a task is running to change the virtual machine state.
// If the virtual machine is in a state with a task in progress, it
// transitions to a new state when the task completes. For example, a virtual
// machine continues to be in the poweredOn state while a suspend task
// is running, and changes to the suspended state once the task finishes.
//
// As a consequence of this approach, clients interested in monitoring
// the status of a virtual machine should typically track the
// `activeTask` data object in addition to the
// `powerState` object.
type VirtualMachinePowerState string

const (
	// The virtual machine is currently powered off.
	VirtualMachinePowerStatePoweredOff = VirtualMachinePowerState("poweredOff")
	// The virtual machine is currently powered on.
	VirtualMachinePowerStatePoweredOn = VirtualMachinePowerState("poweredOn")
	// The virtual machine is currently suspended.
	VirtualMachinePowerStateSuspended = VirtualMachinePowerState("suspended")
)

func (e VirtualMachinePowerState) Values() []VirtualMachinePowerState {
	return []VirtualMachinePowerState{
		VirtualMachinePowerStatePoweredOff,
		VirtualMachinePowerStatePoweredOn,
		VirtualMachinePowerStateSuspended,
	}
}

func (e VirtualMachinePowerState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachinePowerState"] = reflect.TypeOf((*VirtualMachinePowerState)(nil)).Elem()
}

// Deprecated as of vSphere API 6.0.
//
// The RecordReplayState type defines a simple set of record and replay
// states for a virtual machine.
type VirtualMachineRecordReplayState string

const (
	// The virtual machine is recording.
	VirtualMachineRecordReplayStateRecording = VirtualMachineRecordReplayState("recording")
	// The virtual machine is replaying.
	VirtualMachineRecordReplayStateReplaying = VirtualMachineRecordReplayState("replaying")
	// The virtual machine is currently not participating
	// in record or replay.
	VirtualMachineRecordReplayStateInactive = VirtualMachineRecordReplayState("inactive")
)

func (e VirtualMachineRecordReplayState) Values() []VirtualMachineRecordReplayState {
	return []VirtualMachineRecordReplayState{
		VirtualMachineRecordReplayStateRecording,
		VirtualMachineRecordReplayStateReplaying,
		VirtualMachineRecordReplayStateInactive,
	}
}

func (e VirtualMachineRecordReplayState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineRecordReplayState"] = reflect.TypeOf((*VirtualMachineRecordReplayState)(nil)).Elem()
}

// Specifies how a virtual disk is moved or copied to a
// datastore.
//
// In all cases after the move or copy the virtual machine's current running point
// will be placed on the target datastore. The current running point is defined
// as the disk backing which the virtual machine is currently
// writing to. This end state can be achieved in multiple
// ways, and the supported options are described in this
// enumeration.
//
// These options are only relevant when the backing of the
// specified disk is a *file backing*.
//
// Since disk backings may become shared as the result of
// either a *clone operation* or
// a *relocate operation*,
// `VirtualMachine.PromoteDisks_Task` has been provided as
// a way to unshare such disk backings.
//
// See also `VirtualDiskSparseVer1BackingInfo.parent`, `VirtualDiskSparseVer2BackingInfo.parent`, `VirtualDiskFlatVer1BackingInfo.parent`, `VirtualDiskFlatVer2BackingInfo.parent`, `VirtualDiskRawDiskMappingVer1BackingInfo.parent`, `VirtualMachineRelocateSpec.diskMoveType`, `VirtualMachineRelocateSpecDiskLocator.diskMoveType`.
type VirtualMachineRelocateDiskMoveOptions string

const (
	// All of the virtual disk's backings should be moved to the new datastore.
	//
	// If a disk backing is not the child-most backing of this virtual machine,
	// and there exists a read-only disk backing with the same content ID
	// on the target datastore, then this disk backing may not be copied. Instead
	// it is acceptable to attach to the read-only disk backing at the target
	// datastore. A read-only disk backing is defined as a virtual disk
	// backing which no virtual machine is currently writing to.
	//
	// See also `VirtualDiskSparseVer1BackingInfo.contentId`, `VirtualDiskSparseVer2BackingInfo.contentId`, `VirtualDiskFlatVer1BackingInfo.contentId`, `VirtualDiskFlatVer2BackingInfo.contentId`, `VirtualDiskRawDiskMappingVer1BackingInfo.contentId`.
	VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndAllowSharing = VirtualMachineRelocateDiskMoveOptions("moveAllDiskBackingsAndAllowSharing")
	// All of the virtual disk's backings should be moved to the new datastore.
	//
	// It is not acceptable to attach to a disk backing with the same content ID
	// on the destination datastore. During a *clone operation* any delta disk backings will be consolidated.
	VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndDisallowSharing = VirtualMachineRelocateDiskMoveOptions("moveAllDiskBackingsAndDisallowSharing")
	// Move only the child-most disk backing.
	//
	// Any parent disk backings should
	// be left in their current locations.
	//
	// This option only differs from `moveAllDiskBackingsAndAllowSharing` and
	// `moveAllDiskBackingsAndDisallowSharing` when the virtual
	// disk has a parent backing.
	//
	// Note that in the case of a *clone operation*,
	// this means that the parent disks will now be shared. This is safe as any
	// parent disks are always read-only.
	// Note that in the case of a `VirtualMachine.RelocateVM_Task` operation,
	// only the virtual disks in the current virtual machine configuration are moved.
	VirtualMachineRelocateDiskMoveOptionsMoveChildMostDiskBacking = VirtualMachineRelocateDiskMoveOptions("moveChildMostDiskBacking")
	// Create a new child disk backing on the destination datastore.
	//
	// None of the
	// virtual disk's existing files should be moved from their current locations.
	//
	// Note that in the case of a *clone operation*,
	// this means that the original virtual machine's disks are now all being shared.
	// This is only safe if the clone was taken from a snapshot point, because
	// snapshot points are always read-only. Thus for a clone this
	// option is only valid *when cloning from a snapshot*.
	// createNewChildDiskBacking is not a supported operation for
	// `VirtualMachine.RelocateVM_Task` operations unless all disks are moving.
	VirtualMachineRelocateDiskMoveOptionsCreateNewChildDiskBacking = VirtualMachineRelocateDiskMoveOptions("createNewChildDiskBacking")
	// All of the virtual disk's backings should be moved to the new datastore.
	//
	// During a *clone operation* or a
	// `VirtualMachine.MigrateVM_Task`, any delta disk backings will be
	// consolidated.
	VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndConsolidate = VirtualMachineRelocateDiskMoveOptions("moveAllDiskBackingsAndConsolidate")
)

func (e VirtualMachineRelocateDiskMoveOptions) Values() []VirtualMachineRelocateDiskMoveOptions {
	return []VirtualMachineRelocateDiskMoveOptions{
		VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndAllowSharing,
		VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndDisallowSharing,
		VirtualMachineRelocateDiskMoveOptionsMoveChildMostDiskBacking,
		VirtualMachineRelocateDiskMoveOptionsCreateNewChildDiskBacking,
		VirtualMachineRelocateDiskMoveOptionsMoveAllDiskBackingsAndConsolidate,
	}
}

func (e VirtualMachineRelocateDiskMoveOptions) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineRelocateDiskMoveOptions"] = reflect.TypeOf((*VirtualMachineRelocateDiskMoveOptions)(nil)).Elem()
}

// Deprecated as of vSphere API 5.0.
//
// The set of tranformations that can be performed on the virtual disks
// as part of the copy.
type VirtualMachineRelocateTransformation string

const (
	VirtualMachineRelocateTransformationFlat   = VirtualMachineRelocateTransformation("flat")
	VirtualMachineRelocateTransformationSparse = VirtualMachineRelocateTransformation("sparse")
)

func (e VirtualMachineRelocateTransformation) Values() []VirtualMachineRelocateTransformation {
	return []VirtualMachineRelocateTransformation{
		VirtualMachineRelocateTransformationFlat,
		VirtualMachineRelocateTransformationSparse,
	}
}

func (e VirtualMachineRelocateTransformation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineRelocateTransformation"] = reflect.TypeOf((*VirtualMachineRelocateTransformation)(nil)).Elem()
}

// Possible SCSI classes.
type VirtualMachineScsiPassthroughType string

const (
	VirtualMachineScsiPassthroughTypeDisk      = VirtualMachineScsiPassthroughType("disk")
	VirtualMachineScsiPassthroughTypeTape      = VirtualMachineScsiPassthroughType("tape")
	VirtualMachineScsiPassthroughTypePrinter   = VirtualMachineScsiPassthroughType("printer")
	VirtualMachineScsiPassthroughTypeProcessor = VirtualMachineScsiPassthroughType("processor")
	VirtualMachineScsiPassthroughTypeWorm      = VirtualMachineScsiPassthroughType("worm")
	VirtualMachineScsiPassthroughTypeCdrom     = VirtualMachineScsiPassthroughType("cdrom")
	VirtualMachineScsiPassthroughTypeScanner   = VirtualMachineScsiPassthroughType("scanner")
	VirtualMachineScsiPassthroughTypeOptical   = VirtualMachineScsiPassthroughType("optical")
	VirtualMachineScsiPassthroughTypeMedia     = VirtualMachineScsiPassthroughType("media")
	VirtualMachineScsiPassthroughTypeCom       = VirtualMachineScsiPassthroughType("com")
	VirtualMachineScsiPassthroughTypeRaid      = VirtualMachineScsiPassthroughType("raid")
	VirtualMachineScsiPassthroughTypeUnknown   = VirtualMachineScsiPassthroughType("unknown")
)

func (e VirtualMachineScsiPassthroughType) Values() []VirtualMachineScsiPassthroughType {
	return []VirtualMachineScsiPassthroughType{
		VirtualMachineScsiPassthroughTypeDisk,
		VirtualMachineScsiPassthroughTypeTape,
		VirtualMachineScsiPassthroughTypePrinter,
		VirtualMachineScsiPassthroughTypeProcessor,
		VirtualMachineScsiPassthroughTypeWorm,
		VirtualMachineScsiPassthroughTypeCdrom,
		VirtualMachineScsiPassthroughTypeScanner,
		VirtualMachineScsiPassthroughTypeOptical,
		VirtualMachineScsiPassthroughTypeMedia,
		VirtualMachineScsiPassthroughTypeCom,
		VirtualMachineScsiPassthroughTypeRaid,
		VirtualMachineScsiPassthroughTypeUnknown,
	}
}

func (e VirtualMachineScsiPassthroughType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineScsiPassthroughType"] = reflect.TypeOf((*VirtualMachineScsiPassthroughType)(nil)).Elem()
}

// Flexible Launch Enclave (FLC) modes.
type VirtualMachineSgxInfoFlcModes string

const (
	// FLC is available in the guest.
	//
	// The "launch Enclave MSRs" are locked and
	// initialized with the provided public key hash.
	VirtualMachineSgxInfoFlcModesLocked = VirtualMachineSgxInfoFlcModes("locked")
	// FLC is available in the guest.
	//
	// The "launch enclave MSRs" are writeable
	// and initialized with Intel's public key hash.
	VirtualMachineSgxInfoFlcModesUnlocked = VirtualMachineSgxInfoFlcModes("unlocked")
)

func (e VirtualMachineSgxInfoFlcModes) Values() []VirtualMachineSgxInfoFlcModes {
	return []VirtualMachineSgxInfoFlcModes{
		VirtualMachineSgxInfoFlcModesLocked,
		VirtualMachineSgxInfoFlcModesUnlocked,
	}
}

func (e VirtualMachineSgxInfoFlcModes) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineSgxInfoFlcModes"] = reflect.TypeOf((*VirtualMachineSgxInfoFlcModes)(nil)).Elem()
}

// The list of possible standby actions that the virtual machine can take
// for S1 ACPI.
type VirtualMachineStandbyActionType string

const (
	VirtualMachineStandbyActionTypeCheckpoint     = VirtualMachineStandbyActionType("checkpoint")
	VirtualMachineStandbyActionTypePowerOnSuspend = VirtualMachineStandbyActionType("powerOnSuspend")
)

func (e VirtualMachineStandbyActionType) Values() []VirtualMachineStandbyActionType {
	return []VirtualMachineStandbyActionType{
		VirtualMachineStandbyActionTypeCheckpoint,
		VirtualMachineStandbyActionTypePowerOnSuspend,
	}
}

func (e VirtualMachineStandbyActionType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineStandbyActionType"] = reflect.TypeOf((*VirtualMachineStandbyActionType)(nil)).Elem()
}

// Describes how widely the endpoint is available in a cluster.
//
// Note that these fields are not necessarily mutual-exclusive.
type VirtualMachineTargetInfoConfigurationTag string

const (
	// Indicates that this device is part of the cluster compliant
	// specification.
	VirtualMachineTargetInfoConfigurationTagCompliant = VirtualMachineTargetInfoConfigurationTag("compliant")
	// Indicates that this is available for all hosts in the cluster.
	VirtualMachineTargetInfoConfigurationTagClusterWide = VirtualMachineTargetInfoConfigurationTag("clusterWide")
)

func (e VirtualMachineTargetInfoConfigurationTag) Values() []VirtualMachineTargetInfoConfigurationTag {
	return []VirtualMachineTargetInfoConfigurationTag{
		VirtualMachineTargetInfoConfigurationTagCompliant,
		VirtualMachineTargetInfoConfigurationTagClusterWide,
	}
}

func (e VirtualMachineTargetInfoConfigurationTag) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineTargetInfoConfigurationTag"] = reflect.TypeOf((*VirtualMachineTargetInfoConfigurationTag)(nil)).Elem()
}

// The virtual machine ticket type.
type VirtualMachineTicketType string

const (
	// Deprecated as of vSphere API 8.0. Use `webmks` instead.
	//
	// Remote mouse-keyboard-screen ticket.
	VirtualMachineTicketTypeMks = VirtualMachineTicketType("mks")
	// Deprecated as of vSphere 8.0 API. Use `webRemoteDevice`
	// instead.
	//
	// Remote device ticket.
	VirtualMachineTicketTypeDevice = VirtualMachineTicketType("device")
	// Deprecated as of vSphere 6.6.3 API. Use
	// `GuestOperationsManager` instead.
	//
	// Guest operation ticket.
	VirtualMachineTicketTypeGuestControl = VirtualMachineTicketType("guestControl")
	// Mouse-keyboard-screen over WebSocket ticket.
	//
	// MKS protocol is VNC (a.k.a. RFB) protocol with
	// VMware extensions; the protocol gracefully degrades
	// to standard VNC if extensions are not available.
	// wss://{Ticket.host}/ticket/{Ticket.ticket}
	VirtualMachineTicketTypeWebmks = VirtualMachineTicketType("webmks")
	// Guest Integrity over WebSocket ticket.
	//
	// This ticket grants the client read-only access to guest integrity
	// messages and alerts.
	VirtualMachineTicketTypeGuestIntegrity = VirtualMachineTicketType("guestIntegrity")
	// Remote device over WebSocket ticket.
	VirtualMachineTicketTypeWebRemoteDevice = VirtualMachineTicketType("webRemoteDevice")
)

func (e VirtualMachineTicketType) Values() []VirtualMachineTicketType {
	return []VirtualMachineTicketType{
		VirtualMachineTicketTypeMks,
		VirtualMachineTicketTypeDevice,
		VirtualMachineTicketTypeGuestControl,
		VirtualMachineTicketTypeWebmks,
		VirtualMachineTicketTypeGuestIntegrity,
		VirtualMachineTicketTypeWebRemoteDevice,
	}
}

func (e VirtualMachineTicketType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineTicketType"] = reflect.TypeOf((*VirtualMachineTicketType)(nil)).Elem()
}

// The installation type of tools in the VM.
type VirtualMachineToolsInstallType string

const (
	// Installation type is not known.
	//
	// Most likely tools have been
	// installed by OSPs or open-vm-tools, but a version that does
	// not report its install type or an install type that we do
	// not recognize.
	VirtualMachineToolsInstallTypeGuestToolsTypeUnknown = VirtualMachineToolsInstallType("guestToolsTypeUnknown")
	// MSI is the installation type used for VMware Tools on Windows.
	VirtualMachineToolsInstallTypeGuestToolsTypeMSI = VirtualMachineToolsInstallType("guestToolsTypeMSI")
	// Tools have been installed by the tar installer.
	VirtualMachineToolsInstallTypeGuestToolsTypeTar = VirtualMachineToolsInstallType("guestToolsTypeTar")
	// OSPs are RPM or Debian packages tailored for the OS in the VM.
	//
	// See http://packages.vmware.com
	VirtualMachineToolsInstallTypeGuestToolsTypeOSP = VirtualMachineToolsInstallType("guestToolsTypeOSP")
	// open-vm-tools are the open-source version of VMware Tools, may have
	// been packaged by the OS vendor.
	VirtualMachineToolsInstallTypeGuestToolsTypeOpenVMTools = VirtualMachineToolsInstallType("guestToolsTypeOpenVMTools")
)

func (e VirtualMachineToolsInstallType) Values() []VirtualMachineToolsInstallType {
	return []VirtualMachineToolsInstallType{
		VirtualMachineToolsInstallTypeGuestToolsTypeUnknown,
		VirtualMachineToolsInstallTypeGuestToolsTypeMSI,
		VirtualMachineToolsInstallTypeGuestToolsTypeTar,
		VirtualMachineToolsInstallTypeGuestToolsTypeOSP,
		VirtualMachineToolsInstallTypeGuestToolsTypeOpenVMTools,
	}
}

func (e VirtualMachineToolsInstallType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineToolsInstallType"] = reflect.TypeOf((*VirtualMachineToolsInstallType)(nil)).Elem()
}

// Current running status of VMware Tools running in the guest
// operating system.
type VirtualMachineToolsRunningStatus string

const (
	// VMware Tools is not running.
	VirtualMachineToolsRunningStatusGuestToolsNotRunning = VirtualMachineToolsRunningStatus("guestToolsNotRunning")
	// VMware Tools is running.
	VirtualMachineToolsRunningStatusGuestToolsRunning = VirtualMachineToolsRunningStatus("guestToolsRunning")
	// VMware Tools is starting.
	VirtualMachineToolsRunningStatusGuestToolsExecutingScripts = VirtualMachineToolsRunningStatus("guestToolsExecutingScripts")
)

func (e VirtualMachineToolsRunningStatus) Values() []VirtualMachineToolsRunningStatus {
	return []VirtualMachineToolsRunningStatus{
		VirtualMachineToolsRunningStatusGuestToolsNotRunning,
		VirtualMachineToolsRunningStatusGuestToolsRunning,
		VirtualMachineToolsRunningStatusGuestToolsExecutingScripts,
	}
}

func (e VirtualMachineToolsRunningStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineToolsRunningStatus"] = reflect.TypeOf((*VirtualMachineToolsRunningStatus)(nil)).Elem()
}

// Deprecated as of vSphere API 4.0 use `VirtualMachineToolsVersionStatus_enum`
// and `VirtualMachineToolsRunningStatus_enum`.
//
// Current status of VMware Tools running in the guest operating system.
type VirtualMachineToolsStatus string

const (
	// VMware Tools has never been installed
	// or has not run in the virtual machine.
	VirtualMachineToolsStatusToolsNotInstalled = VirtualMachineToolsStatus("toolsNotInstalled")
	// VMware Tools is not running.
	VirtualMachineToolsStatusToolsNotRunning = VirtualMachineToolsStatus("toolsNotRunning")
	// VMware Tools is running, but the version is not current.
	VirtualMachineToolsStatusToolsOld = VirtualMachineToolsStatus("toolsOld")
	// VMware Tools is running and the version is current.
	VirtualMachineToolsStatusToolsOk = VirtualMachineToolsStatus("toolsOk")
)

func (e VirtualMachineToolsStatus) Values() []VirtualMachineToolsStatus {
	return []VirtualMachineToolsStatus{
		VirtualMachineToolsStatusToolsNotInstalled,
		VirtualMachineToolsStatusToolsNotRunning,
		VirtualMachineToolsStatusToolsOld,
		VirtualMachineToolsStatusToolsOk,
	}
}

func (e VirtualMachineToolsStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineToolsStatus"] = reflect.TypeOf((*VirtualMachineToolsStatus)(nil)).Elem()
}

// Current version status of VMware Tools installed in the guest operating
// system.
type VirtualMachineToolsVersionStatus string

const (
	// VMware Tools has never been installed.
	VirtualMachineToolsVersionStatusGuestToolsNotInstalled = VirtualMachineToolsVersionStatus("guestToolsNotInstalled")
	// Deprecated as of vSphere API 5.1 value is not reported by
	// toolsVersionStatus2, instead more detailed status is reported.
	//
	// VMware Tools is installed, but the version is not current.
	VirtualMachineToolsVersionStatusGuestToolsNeedUpgrade = VirtualMachineToolsVersionStatus("guestToolsNeedUpgrade")
	// VMware Tools is installed, and the version is current.
	VirtualMachineToolsVersionStatusGuestToolsCurrent = VirtualMachineToolsVersionStatus("guestToolsCurrent")
	// VMware Tools is installed, but it is not managed by VMWare.
	VirtualMachineToolsVersionStatusGuestToolsUnmanaged = VirtualMachineToolsVersionStatus("guestToolsUnmanaged")
	// VMware Tools is installed, but the version is too old.
	VirtualMachineToolsVersionStatusGuestToolsTooOld = VirtualMachineToolsVersionStatus("guestToolsTooOld")
	// VMware Tools is installed, supported, but a newer version is available.
	VirtualMachineToolsVersionStatusGuestToolsSupportedOld = VirtualMachineToolsVersionStatus("guestToolsSupportedOld")
	// VMware Tools is installed, supported, and newer
	// than the version available on the host.
	VirtualMachineToolsVersionStatusGuestToolsSupportedNew = VirtualMachineToolsVersionStatus("guestToolsSupportedNew")
	// VMware Tools is installed, and the version is known to be
	// too new to work correctly with this virtual machine.
	VirtualMachineToolsVersionStatusGuestToolsTooNew = VirtualMachineToolsVersionStatus("guestToolsTooNew")
	// VMware Tools is installed, but the installed version is
	// known to have a grave bug and should be immediately upgraded.
	VirtualMachineToolsVersionStatusGuestToolsBlacklisted = VirtualMachineToolsVersionStatus("guestToolsBlacklisted")
)

func (e VirtualMachineToolsVersionStatus) Values() []VirtualMachineToolsVersionStatus {
	return []VirtualMachineToolsVersionStatus{
		VirtualMachineToolsVersionStatusGuestToolsNotInstalled,
		VirtualMachineToolsVersionStatusGuestToolsNeedUpgrade,
		VirtualMachineToolsVersionStatusGuestToolsCurrent,
		VirtualMachineToolsVersionStatusGuestToolsUnmanaged,
		VirtualMachineToolsVersionStatusGuestToolsTooOld,
		VirtualMachineToolsVersionStatusGuestToolsSupportedOld,
		VirtualMachineToolsVersionStatusGuestToolsSupportedNew,
		VirtualMachineToolsVersionStatusGuestToolsTooNew,
		VirtualMachineToolsVersionStatusGuestToolsBlacklisted,
	}
}

func (e VirtualMachineToolsVersionStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineToolsVersionStatus"] = reflect.TypeOf((*VirtualMachineToolsVersionStatus)(nil)).Elem()
}

// Device class family.
type VirtualMachineUsbInfoFamily string

const (
	// Audio capable device.
	VirtualMachineUsbInfoFamilyAudio = VirtualMachineUsbInfoFamily("audio")
	// Human interface device.
	VirtualMachineUsbInfoFamilyHid = VirtualMachineUsbInfoFamily("hid")
	// Bootable human interface device, this is a subset of HID devices.
	VirtualMachineUsbInfoFamilyHid_bootable = VirtualMachineUsbInfoFamily("hid_bootable")
	// Physical interface device.
	VirtualMachineUsbInfoFamilyPhysical = VirtualMachineUsbInfoFamily("physical")
	// Communication device.
	VirtualMachineUsbInfoFamilyCommunication = VirtualMachineUsbInfoFamily("communication")
	// Still imaging device.
	VirtualMachineUsbInfoFamilyImaging = VirtualMachineUsbInfoFamily("imaging")
	// Printer device.
	VirtualMachineUsbInfoFamilyPrinter = VirtualMachineUsbInfoFamily("printer")
	// Mass storage device.
	VirtualMachineUsbInfoFamilyStorage = VirtualMachineUsbInfoFamily("storage")
	// USB hubs.
	VirtualMachineUsbInfoFamilyHub = VirtualMachineUsbInfoFamily("hub")
	// Smart card device.
	VirtualMachineUsbInfoFamilySmart_card = VirtualMachineUsbInfoFamily("smart_card")
	// Content security device.
	VirtualMachineUsbInfoFamilySecurity = VirtualMachineUsbInfoFamily("security")
	// Video device.
	VirtualMachineUsbInfoFamilyVideo = VirtualMachineUsbInfoFamily("video")
	// Wireless controller.
	VirtualMachineUsbInfoFamilyWireless = VirtualMachineUsbInfoFamily("wireless")
	// Standard bluetooth adapter that uses HCI protocol,
	// this is a subset of wireless controllers.
	VirtualMachineUsbInfoFamilyBluetooth = VirtualMachineUsbInfoFamily("bluetooth")
	// Wireless device related to the Wireless USB standard,
	// this is a subset of wireless controllers,
	VirtualMachineUsbInfoFamilyWusb = VirtualMachineUsbInfoFamily("wusb")
	// Palm PDA, and Micorsoft ActiveSync PDA.
	VirtualMachineUsbInfoFamilyPda = VirtualMachineUsbInfoFamily("pda")
	// Device that has an interface using a vendor-specific protocol.
	VirtualMachineUsbInfoFamilyVendor_specific = VirtualMachineUsbInfoFamily("vendor_specific")
	// Other miscellaneous device.
	VirtualMachineUsbInfoFamilyOther = VirtualMachineUsbInfoFamily("other")
	// There was an error in determining this device's classes
	// accurately.
	VirtualMachineUsbInfoFamilyUnknownFamily = VirtualMachineUsbInfoFamily("unknownFamily")
)

func (e VirtualMachineUsbInfoFamily) Values() []VirtualMachineUsbInfoFamily {
	return []VirtualMachineUsbInfoFamily{
		VirtualMachineUsbInfoFamilyAudio,
		VirtualMachineUsbInfoFamilyHid,
		VirtualMachineUsbInfoFamilyHid_bootable,
		VirtualMachineUsbInfoFamilyPhysical,
		VirtualMachineUsbInfoFamilyCommunication,
		VirtualMachineUsbInfoFamilyImaging,
		VirtualMachineUsbInfoFamilyPrinter,
		VirtualMachineUsbInfoFamilyStorage,
		VirtualMachineUsbInfoFamilyHub,
		VirtualMachineUsbInfoFamilySmart_card,
		VirtualMachineUsbInfoFamilySecurity,
		VirtualMachineUsbInfoFamilyVideo,
		VirtualMachineUsbInfoFamilyWireless,
		VirtualMachineUsbInfoFamilyBluetooth,
		VirtualMachineUsbInfoFamilyWusb,
		VirtualMachineUsbInfoFamilyPda,
		VirtualMachineUsbInfoFamilyVendor_specific,
		VirtualMachineUsbInfoFamilyOther,
		VirtualMachineUsbInfoFamilyUnknownFamily,
	}
}

func (e VirtualMachineUsbInfoFamily) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineUsbInfoFamily"] = reflect.TypeOf((*VirtualMachineUsbInfoFamily)(nil)).Elem()
}

// Device speed.
type VirtualMachineUsbInfoSpeed string

const (
	// This device operates at low speed (1.5Mb/s).
	VirtualMachineUsbInfoSpeedLow = VirtualMachineUsbInfoSpeed("low")
	// This device operates at full speed (12Mb/s).
	VirtualMachineUsbInfoSpeedFull = VirtualMachineUsbInfoSpeed("full")
	// This device can operate at high speed (480Mb/s)
	VirtualMachineUsbInfoSpeedHigh = VirtualMachineUsbInfoSpeed("high")
	// This device can operate at super speed (4.8Gb/s)
	VirtualMachineUsbInfoSpeedSuperSpeed = VirtualMachineUsbInfoSpeed("superSpeed")
	// This device can operate at super speed plus (10Gb/s)
	VirtualMachineUsbInfoSpeedSuperSpeedPlus = VirtualMachineUsbInfoSpeed("superSpeedPlus")
	// This device can operate at super speed gen 2x2 (20Gb/s)
	VirtualMachineUsbInfoSpeedSuperSpeed20Gbps = VirtualMachineUsbInfoSpeed("superSpeed20Gbps")
	// This device's speed is unknown.
	VirtualMachineUsbInfoSpeedUnknownSpeed = VirtualMachineUsbInfoSpeed("unknownSpeed")
)

func (e VirtualMachineUsbInfoSpeed) Values() []VirtualMachineUsbInfoSpeed {
	return []VirtualMachineUsbInfoSpeed{
		VirtualMachineUsbInfoSpeedLow,
		VirtualMachineUsbInfoSpeedFull,
		VirtualMachineUsbInfoSpeedHigh,
		VirtualMachineUsbInfoSpeedSuperSpeed,
		VirtualMachineUsbInfoSpeedSuperSpeedPlus,
		VirtualMachineUsbInfoSpeedSuperSpeed20Gbps,
		VirtualMachineUsbInfoSpeedUnknownSpeed,
	}
}

func (e VirtualMachineUsbInfoSpeed) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineUsbInfoSpeed"] = reflect.TypeOf((*VirtualMachineUsbInfoSpeed)(nil)).Elem()
	minAPIVersionForEnumValue["VirtualMachineUsbInfoSpeed"] = map[string]string{
		"superSpeed20Gbps": "7.0.3.2",
	}
}

// Set of possible values for action field in FilterSpec.
//
// Determines whether traffic is allowed or denied.
type VirtualMachineVMCIDeviceAction string

const (
	// Allow communication.
	VirtualMachineVMCIDeviceActionAllow = VirtualMachineVMCIDeviceAction("allow")
	// Deny communication.
	VirtualMachineVMCIDeviceActionDeny = VirtualMachineVMCIDeviceAction("deny")
)

func (e VirtualMachineVMCIDeviceAction) Values() []VirtualMachineVMCIDeviceAction {
	return []VirtualMachineVMCIDeviceAction{
		VirtualMachineVMCIDeviceActionAllow,
		VirtualMachineVMCIDeviceActionDeny,
	}
}

func (e VirtualMachineVMCIDeviceAction) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVMCIDeviceAction"] = reflect.TypeOf((*VirtualMachineVMCIDeviceAction)(nil)).Elem()
}

// Set of possible values for direction field in FilterSpec.
type VirtualMachineVMCIDeviceDirection string

const (
	// from host to guest
	VirtualMachineVMCIDeviceDirectionGuest = VirtualMachineVMCIDeviceDirection("guest")
	// from guest to host
	VirtualMachineVMCIDeviceDirectionHost = VirtualMachineVMCIDeviceDirection("host")
	// all of the above
	VirtualMachineVMCIDeviceDirectionAnyDirection = VirtualMachineVMCIDeviceDirection("anyDirection")
)

func (e VirtualMachineVMCIDeviceDirection) Values() []VirtualMachineVMCIDeviceDirection {
	return []VirtualMachineVMCIDeviceDirection{
		VirtualMachineVMCIDeviceDirectionGuest,
		VirtualMachineVMCIDeviceDirectionHost,
		VirtualMachineVMCIDeviceDirectionAnyDirection,
	}
}

func (e VirtualMachineVMCIDeviceDirection) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVMCIDeviceDirection"] = reflect.TypeOf((*VirtualMachineVMCIDeviceDirection)(nil)).Elem()
}

// Set of possible values for protocol field in FilterSpec.
type VirtualMachineVMCIDeviceProtocol string

const (
	// VMCI hypervisor datagram send op.
	//
	// Direction code is not applicable to this one.
	VirtualMachineVMCIDeviceProtocolHypervisor = VirtualMachineVMCIDeviceProtocol("hypervisor")
	// VMCI doorbell notification
	VirtualMachineVMCIDeviceProtocolDoorbell = VirtualMachineVMCIDeviceProtocol("doorbell")
	// VMCI queue pair alloc operation.
	//
	// Direction code not applicable to this one.
	VirtualMachineVMCIDeviceProtocolQueuepair = VirtualMachineVMCIDeviceProtocol("queuepair")
	// VMCI and VMCI Socket datagram send op.
	//
	// Since VMCI Socket datagrams map ports directly to resources,
	// there is no need to distinguish between the two.
	VirtualMachineVMCIDeviceProtocolDatagram = VirtualMachineVMCIDeviceProtocol("datagram")
	// VMCI Stream Socket connect op.
	VirtualMachineVMCIDeviceProtocolStream = VirtualMachineVMCIDeviceProtocol("stream")
	// All of the above.
	VirtualMachineVMCIDeviceProtocolAnyProtocol = VirtualMachineVMCIDeviceProtocol("anyProtocol")
)

func (e VirtualMachineVMCIDeviceProtocol) Values() []VirtualMachineVMCIDeviceProtocol {
	return []VirtualMachineVMCIDeviceProtocol{
		VirtualMachineVMCIDeviceProtocolHypervisor,
		VirtualMachineVMCIDeviceProtocolDoorbell,
		VirtualMachineVMCIDeviceProtocolQueuepair,
		VirtualMachineVMCIDeviceProtocolDatagram,
		VirtualMachineVMCIDeviceProtocolStream,
		VirtualMachineVMCIDeviceProtocolAnyProtocol,
	}
}

func (e VirtualMachineVMCIDeviceProtocol) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVMCIDeviceProtocol"] = reflect.TypeOf((*VirtualMachineVMCIDeviceProtocol)(nil)).Elem()
}

type VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType string

const (
	VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypePciPassthru = VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType("pciPassthru")
	VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypeNvidiaVgpu  = VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType("nvidiaVgpu")
	VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypeSriovNic    = VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType("sriovNic")
	VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypeDvx         = VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType("dvx")
)

func (e VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType) Values() []VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType {
	return []VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType{
		VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypePciPassthru,
		VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypeNvidiaVgpu,
		VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypeSriovNic,
		VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentTypeDvx,
	}
}

func (e VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType"] = reflect.TypeOf((*VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType)(nil)).Elem()
	minAPIVersionForType["VirtualMachineVendorDeviceGroupInfoComponentDeviceInfoComponentType"] = "8.0.0.1"
}

type VirtualMachineVgpuProfileInfoProfileClass string

const (
	VirtualMachineVgpuProfileInfoProfileClassCompute = VirtualMachineVgpuProfileInfoProfileClass("compute")
	VirtualMachineVgpuProfileInfoProfileClassQuadro  = VirtualMachineVgpuProfileInfoProfileClass("quadro")
)

func (e VirtualMachineVgpuProfileInfoProfileClass) Values() []VirtualMachineVgpuProfileInfoProfileClass {
	return []VirtualMachineVgpuProfileInfoProfileClass{
		VirtualMachineVgpuProfileInfoProfileClassCompute,
		VirtualMachineVgpuProfileInfoProfileClassQuadro,
	}
}

func (e VirtualMachineVgpuProfileInfoProfileClass) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVgpuProfileInfoProfileClass"] = reflect.TypeOf((*VirtualMachineVgpuProfileInfoProfileClass)(nil)).Elem()
	minAPIVersionForType["VirtualMachineVgpuProfileInfoProfileClass"] = "7.0.3.0"
}

type VirtualMachineVgpuProfileInfoProfileSharing string

const (
	// Time-sliced
	VirtualMachineVgpuProfileInfoProfileSharingTimeSliced = VirtualMachineVgpuProfileInfoProfileSharing("timeSliced")
	// Multi-instance GPU partitioning
	VirtualMachineVgpuProfileInfoProfileSharingMig = VirtualMachineVgpuProfileInfoProfileSharing("mig")
)

func (e VirtualMachineVgpuProfileInfoProfileSharing) Values() []VirtualMachineVgpuProfileInfoProfileSharing {
	return []VirtualMachineVgpuProfileInfoProfileSharing{
		VirtualMachineVgpuProfileInfoProfileSharingTimeSliced,
		VirtualMachineVgpuProfileInfoProfileSharingMig,
	}
}

func (e VirtualMachineVgpuProfileInfoProfileSharing) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVgpuProfileInfoProfileSharing"] = reflect.TypeOf((*VirtualMachineVgpuProfileInfoProfileSharing)(nil)).Elem()
	minAPIVersionForType["VirtualMachineVgpuProfileInfoProfileSharing"] = "7.0.3.0"
}

// Set of possible values for `VirtualMachineVideoCard.use3dRenderer`.
type VirtualMachineVideoCardUse3dRenderer string

const (
	// Determine automatically whether to render 3D with software or hardware.
	VirtualMachineVideoCardUse3dRendererAutomatic = VirtualMachineVideoCardUse3dRenderer("automatic")
	// Render 3D with software.
	VirtualMachineVideoCardUse3dRendererSoftware = VirtualMachineVideoCardUse3dRenderer("software")
	// Render 3D with graphics hardware.
	VirtualMachineVideoCardUse3dRendererHardware = VirtualMachineVideoCardUse3dRenderer("hardware")
)

func (e VirtualMachineVideoCardUse3dRenderer) Values() []VirtualMachineVideoCardUse3dRenderer {
	return []VirtualMachineVideoCardUse3dRenderer{
		VirtualMachineVideoCardUse3dRendererAutomatic,
		VirtualMachineVideoCardUse3dRendererSoftware,
		VirtualMachineVideoCardUse3dRendererHardware,
	}
}

func (e VirtualMachineVideoCardUse3dRenderer) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVideoCardUse3dRenderer"] = reflect.TypeOf((*VirtualMachineVideoCardUse3dRenderer)(nil)).Elem()
}

type VirtualMachineVirtualDeviceSwapDeviceSwapStatus string

const (
	// No operation active.
	VirtualMachineVirtualDeviceSwapDeviceSwapStatusNone = VirtualMachineVirtualDeviceSwapDeviceSwapStatus("none")
	// Device swap will be performed on next restart.
	VirtualMachineVirtualDeviceSwapDeviceSwapStatusScheduled = VirtualMachineVirtualDeviceSwapDeviceSwapStatus("scheduled")
	// Device swap is in progress.
	VirtualMachineVirtualDeviceSwapDeviceSwapStatusInprogress = VirtualMachineVirtualDeviceSwapDeviceSwapStatus("inprogress")
	// Device swap failed.
	VirtualMachineVirtualDeviceSwapDeviceSwapStatusFailed = VirtualMachineVirtualDeviceSwapDeviceSwapStatus("failed")
	// Device swap successfully completed.
	VirtualMachineVirtualDeviceSwapDeviceSwapStatusCompleted = VirtualMachineVirtualDeviceSwapDeviceSwapStatus("completed")
)

func (e VirtualMachineVirtualDeviceSwapDeviceSwapStatus) Values() []VirtualMachineVirtualDeviceSwapDeviceSwapStatus {
	return []VirtualMachineVirtualDeviceSwapDeviceSwapStatus{
		VirtualMachineVirtualDeviceSwapDeviceSwapStatusNone,
		VirtualMachineVirtualDeviceSwapDeviceSwapStatusScheduled,
		VirtualMachineVirtualDeviceSwapDeviceSwapStatusInprogress,
		VirtualMachineVirtualDeviceSwapDeviceSwapStatusFailed,
		VirtualMachineVirtualDeviceSwapDeviceSwapStatusCompleted,
	}
}

func (e VirtualMachineVirtualDeviceSwapDeviceSwapStatus) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVirtualDeviceSwapDeviceSwapStatus"] = reflect.TypeOf((*VirtualMachineVirtualDeviceSwapDeviceSwapStatus)(nil)).Elem()
	minAPIVersionForType["VirtualMachineVirtualDeviceSwapDeviceSwapStatus"] = "8.0.0.1"
}

type VirtualMachineVirtualPMemSnapshotMode string

const (
	// The data on virtual NVDIMMs are not affected by snapshot reverts.
	//
	// Writes to virtual NVDIMMs after a snapshot is taken cannot be
	// reverted to the snapshotted state.
	VirtualMachineVirtualPMemSnapshotModeIndependent_persistent = VirtualMachineVirtualPMemSnapshotMode("independent_persistent")
	// Virtual NVDIMMs are erased and recreated upon snapshot reverts.
	VirtualMachineVirtualPMemSnapshotModeIndependent_eraseonrevert = VirtualMachineVirtualPMemSnapshotMode("independent_eraseonrevert")
)

func (e VirtualMachineVirtualPMemSnapshotMode) Values() []VirtualMachineVirtualPMemSnapshotMode {
	return []VirtualMachineVirtualPMemSnapshotMode{
		VirtualMachineVirtualPMemSnapshotModeIndependent_persistent,
		VirtualMachineVirtualPMemSnapshotModeIndependent_eraseonrevert,
	}
}

func (e VirtualMachineVirtualPMemSnapshotMode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineVirtualPMemSnapshotMode"] = reflect.TypeOf((*VirtualMachineVirtualPMemSnapshotMode)(nil)).Elem()
	minAPIVersionForType["VirtualMachineVirtualPMemSnapshotMode"] = "7.0.3.0"
}

// The VSS Snapshot Context
// VSS\_SNAPSHOT\_CONTEXT values not listed below are not implemented.
type VirtualMachineWindowsQuiesceSpecVssBackupContext string

const (
	// The context value indicates auto selection of VSS snapshot context.
	//
	// The ctx\_backup may make Windows VSS-aware applications quiescing during
	// backup. The ctx\_auto makes VMTools select ctx\_file\_share\_backup context
	// if ctx\_backup is not available.
	VirtualMachineWindowsQuiesceSpecVssBackupContextCtx_auto = VirtualMachineWindowsQuiesceSpecVssBackupContext("ctx_auto")
	// Indicate VSS\_CTX\_BACKUP.
	VirtualMachineWindowsQuiesceSpecVssBackupContextCtx_backup = VirtualMachineWindowsQuiesceSpecVssBackupContext("ctx_backup")
	// Indicate VSS\_CTX\_FILE\_SHARE\_BACKUP.
	VirtualMachineWindowsQuiesceSpecVssBackupContextCtx_file_share_backup = VirtualMachineWindowsQuiesceSpecVssBackupContext("ctx_file_share_backup")
)

func (e VirtualMachineWindowsQuiesceSpecVssBackupContext) Values() []VirtualMachineWindowsQuiesceSpecVssBackupContext {
	return []VirtualMachineWindowsQuiesceSpecVssBackupContext{
		VirtualMachineWindowsQuiesceSpecVssBackupContextCtx_auto,
		VirtualMachineWindowsQuiesceSpecVssBackupContextCtx_backup,
		VirtualMachineWindowsQuiesceSpecVssBackupContextCtx_file_share_backup,
	}
}

func (e VirtualMachineWindowsQuiesceSpecVssBackupContext) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualMachineWindowsQuiesceSpecVssBackupContext"] = reflect.TypeOf((*VirtualMachineWindowsQuiesceSpecVssBackupContext)(nil)).Elem()
}

type VirtualNVMEControllerSharing string

const (
	VirtualNVMEControllerSharingNoSharing       = VirtualNVMEControllerSharing("noSharing")
	VirtualNVMEControllerSharingPhysicalSharing = VirtualNVMEControllerSharing("physicalSharing")
)

func (e VirtualNVMEControllerSharing) Values() []VirtualNVMEControllerSharing {
	return []VirtualNVMEControllerSharing{
		VirtualNVMEControllerSharingNoSharing,
		VirtualNVMEControllerSharingPhysicalSharing,
	}
}

func (e VirtualNVMEControllerSharing) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualNVMEControllerSharing"] = reflect.TypeOf((*VirtualNVMEControllerSharing)(nil)).Elem()
	minAPIVersionForType["VirtualNVMEControllerSharing"] = "8.0.2.0"
}

// The valid choices for host pointing devices are:
type VirtualPointingDeviceHostChoice string

const (
	// Automatically detects the host mouse type.
	VirtualPointingDeviceHostChoiceAutodetect = VirtualPointingDeviceHostChoice("autodetect")
	// The Microsoft IntelliMouse Explorer.
	VirtualPointingDeviceHostChoiceIntellimouseExplorer = VirtualPointingDeviceHostChoice("intellimouseExplorer")
	// The Microsoft Intellimouse with a PS2 connection.
	VirtualPointingDeviceHostChoiceIntellimousePs2 = VirtualPointingDeviceHostChoice("intellimousePs2")
	// The Logitech MouseMan.
	VirtualPointingDeviceHostChoiceLogitechMouseman = VirtualPointingDeviceHostChoice("logitechMouseman")
	// The Microsoft Serial Mouse.
	VirtualPointingDeviceHostChoiceMicrosoft_serial = VirtualPointingDeviceHostChoice("microsoft_serial")
	// The Mouse Systems Mouse.
	VirtualPointingDeviceHostChoiceMouseSystems = VirtualPointingDeviceHostChoice("mouseSystems")
	// The Logitech MouseMan Serial Bus Mouse.
	VirtualPointingDeviceHostChoiceMousemanSerial = VirtualPointingDeviceHostChoice("mousemanSerial")
	// A generic mouse with a PS2 connection.
	VirtualPointingDeviceHostChoicePs2 = VirtualPointingDeviceHostChoice("ps2")
)

func (e VirtualPointingDeviceHostChoice) Values() []VirtualPointingDeviceHostChoice {
	return []VirtualPointingDeviceHostChoice{
		VirtualPointingDeviceHostChoiceAutodetect,
		VirtualPointingDeviceHostChoiceIntellimouseExplorer,
		VirtualPointingDeviceHostChoiceIntellimousePs2,
		VirtualPointingDeviceHostChoiceLogitechMouseman,
		VirtualPointingDeviceHostChoiceMicrosoft_serial,
		VirtualPointingDeviceHostChoiceMouseSystems,
		VirtualPointingDeviceHostChoiceMousemanSerial,
		VirtualPointingDeviceHostChoicePs2,
	}
}

func (e VirtualPointingDeviceHostChoice) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualPointingDeviceHostChoice"] = reflect.TypeOf((*VirtualPointingDeviceHostChoice)(nil)).Elem()
}

// Sharing describes three possible ways of sharing the SCSI bus:
// One of these values is assigned to the sharedBus object to determine
// if or how the SCSI bus is shared.
type VirtualSCSISharing string

const (
	// The virtual SCSI bus is not shared.
	VirtualSCSISharingNoSharing = VirtualSCSISharing("noSharing")
	// The virtual SCSI bus is shared between two or more virtual machines.
	//
	// In this case, no physical machine is involved.
	VirtualSCSISharingVirtualSharing = VirtualSCSISharing("virtualSharing")
	// The virtual SCSI bus is shared between two or more virtual machines
	// residing on different physical hosts.
	VirtualSCSISharingPhysicalSharing = VirtualSCSISharing("physicalSharing")
)

func (e VirtualSCSISharing) Values() []VirtualSCSISharing {
	return []VirtualSCSISharing{
		VirtualSCSISharingNoSharing,
		VirtualSCSISharingVirtualSharing,
		VirtualSCSISharingPhysicalSharing,
	}
}

func (e VirtualSCSISharing) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualSCSISharing"] = reflect.TypeOf((*VirtualSCSISharing)(nil)).Elem()
}

// The <code>`VirtualSerialPortEndPoint_enum` enum defines
// endpoint values for virtual serial port pipe backing.
//
// When you use serial port pipe backing to connect a virtual machine
// to another process, you must define the endpoints.
// See the <code>`VirtualSerialPortPipeBackingInfo.endpoint`</code>
// property for the virtual serial port pipe backing information data object.
//
// The possible endpoint values are:
//   - client
//   - server
//
// For the supported choices, see the
// <code>`VirtualSerialPortPipeBackingOption.endpoint`</code>
// property for the virtual serial port pipe backing option data object.
type VirtualSerialPortEndPoint string

const (
	VirtualSerialPortEndPointClient = VirtualSerialPortEndPoint("client")
	VirtualSerialPortEndPointServer = VirtualSerialPortEndPoint("server")
)

func (e VirtualSerialPortEndPoint) Values() []VirtualSerialPortEndPoint {
	return []VirtualSerialPortEndPoint{
		VirtualSerialPortEndPointClient,
		VirtualSerialPortEndPointServer,
	}
}

func (e VirtualSerialPortEndPoint) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualSerialPortEndPoint"] = reflect.TypeOf((*VirtualSerialPortEndPoint)(nil)).Elem()
}

// The enumeration of all known valid VRDMA device protocols.
type VirtualVmxnet3VrdmaOptionDeviceProtocols string

const (
	// A RoCEv1 device.
	VirtualVmxnet3VrdmaOptionDeviceProtocolsRocev1 = VirtualVmxnet3VrdmaOptionDeviceProtocols("rocev1")
	// A RoCEv2 device.
	VirtualVmxnet3VrdmaOptionDeviceProtocolsRocev2 = VirtualVmxnet3VrdmaOptionDeviceProtocols("rocev2")
)

func (e VirtualVmxnet3VrdmaOptionDeviceProtocols) Values() []VirtualVmxnet3VrdmaOptionDeviceProtocols {
	return []VirtualVmxnet3VrdmaOptionDeviceProtocols{
		VirtualVmxnet3VrdmaOptionDeviceProtocolsRocev1,
		VirtualVmxnet3VrdmaOptionDeviceProtocolsRocev2,
	}
}

func (e VirtualVmxnet3VrdmaOptionDeviceProtocols) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VirtualVmxnet3VrdmaOptionDeviceProtocols"] = reflect.TypeOf((*VirtualVmxnet3VrdmaOptionDeviceProtocols)(nil)).Elem()
}

type VmDasBeingResetEventReasonCode string

const (
	// vmtools heartbeat failure
	VmDasBeingResetEventReasonCodeVmtoolsHeartbeatFailure = VmDasBeingResetEventReasonCode("vmtoolsHeartbeatFailure")
	// application heartbeat failure
	VmDasBeingResetEventReasonCodeAppHeartbeatFailure = VmDasBeingResetEventReasonCode("appHeartbeatFailure")
	// immediate reset request
	VmDasBeingResetEventReasonCodeAppImmediateResetRequest = VmDasBeingResetEventReasonCode("appImmediateResetRequest")
	// reset issued by VMCP when APD cleared
	VmDasBeingResetEventReasonCodeVmcpResetApdCleared = VmDasBeingResetEventReasonCode("vmcpResetApdCleared")
)

func (e VmDasBeingResetEventReasonCode) Values() []VmDasBeingResetEventReasonCode {
	return []VmDasBeingResetEventReasonCode{
		VmDasBeingResetEventReasonCodeVmtoolsHeartbeatFailure,
		VmDasBeingResetEventReasonCodeAppHeartbeatFailure,
		VmDasBeingResetEventReasonCodeAppImmediateResetRequest,
		VmDasBeingResetEventReasonCodeVmcpResetApdCleared,
	}
}

func (e VmDasBeingResetEventReasonCode) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VmDasBeingResetEventReasonCode"] = reflect.TypeOf((*VmDasBeingResetEventReasonCode)(nil)).Elem()
}

// The reason for the failure.
type VmFailedStartingSecondaryEventFailureReason string

const (
	// Remote host is incompatible for secondary virtual machine.
	//
	// For instance, the host doesn't have access to the virtual machine's
	// network or datastore.
	VmFailedStartingSecondaryEventFailureReasonIncompatibleHost = VmFailedStartingSecondaryEventFailureReason("incompatibleHost")
	// Login to remote host failed.
	VmFailedStartingSecondaryEventFailureReasonLoginFailed = VmFailedStartingSecondaryEventFailureReason("loginFailed")
	// Registration of the secondary virtual machine
	// on the remote host failed.
	VmFailedStartingSecondaryEventFailureReasonRegisterVmFailed = VmFailedStartingSecondaryEventFailureReason("registerVmFailed")
	// Migration failed.
	VmFailedStartingSecondaryEventFailureReasonMigrateFailed = VmFailedStartingSecondaryEventFailureReason("migrateFailed")
)

func (e VmFailedStartingSecondaryEventFailureReason) Values() []VmFailedStartingSecondaryEventFailureReason {
	return []VmFailedStartingSecondaryEventFailureReason{
		VmFailedStartingSecondaryEventFailureReasonIncompatibleHost,
		VmFailedStartingSecondaryEventFailureReasonLoginFailed,
		VmFailedStartingSecondaryEventFailureReasonRegisterVmFailed,
		VmFailedStartingSecondaryEventFailureReasonMigrateFailed,
	}
}

func (e VmFailedStartingSecondaryEventFailureReason) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VmFailedStartingSecondaryEventFailureReason"] = reflect.TypeOf((*VmFailedStartingSecondaryEventFailureReason)(nil)).Elem()
}

type VmFaultToleranceConfigIssueReasonForIssue string

const (
	// HA is not enabled on the cluster
	VmFaultToleranceConfigIssueReasonForIssueHaNotEnabled = VmFaultToleranceConfigIssueReasonForIssue("haNotEnabled")
	// There is already a secondary virtual machine for the primary
	// virtual machine
	VmFaultToleranceConfigIssueReasonForIssueMoreThanOneSecondary = VmFaultToleranceConfigIssueReasonForIssue("moreThanOneSecondary")
	// Deprecated as of vSphere API 6.0.
	//
	// The virtual machine does not support record/replay.
	//
	// Vm::Capability.RecordReplaySupported is false.
	VmFaultToleranceConfigIssueReasonForIssueRecordReplayNotSupported = VmFaultToleranceConfigIssueReasonForIssue("recordReplayNotSupported")
	// Deprecated as of vSphere API 6.0.
	//
	// It is not possible to turn on Fault Tolerance on this powered-on VM.
	//
	// The support for record/replay should be enabled or Fault Tolerance
	// turned on, when this VM is powered off.
	VmFaultToleranceConfigIssueReasonForIssueReplayNotSupported = VmFaultToleranceConfigIssueReasonForIssue("replayNotSupported")
	// The virtual machine is a template
	VmFaultToleranceConfigIssueReasonForIssueTemplateVm = VmFaultToleranceConfigIssueReasonForIssue("templateVm")
	// The virtual machine has more than one virtual CPU
	VmFaultToleranceConfigIssueReasonForIssueMultipleVCPU = VmFaultToleranceConfigIssueReasonForIssue("multipleVCPU")
	// The host is not active
	VmFaultToleranceConfigIssueReasonForIssueHostInactive = VmFaultToleranceConfigIssueReasonForIssue("hostInactive")
	// The host ftSupported flag is not set because of hardware issues
	VmFaultToleranceConfigIssueReasonForIssueFtUnsupportedHardware = VmFaultToleranceConfigIssueReasonForIssue("ftUnsupportedHardware")
	// The host ftSupported flag is not set because of it is a
	// VMware Server 2.0
	VmFaultToleranceConfigIssueReasonForIssueFtUnsupportedProduct = VmFaultToleranceConfigIssueReasonForIssue("ftUnsupportedProduct")
	// No VMotion license or VMotion nic is not configured on the host
	VmFaultToleranceConfigIssueReasonForIssueMissingVMotionNic = VmFaultToleranceConfigIssueReasonForIssue("missingVMotionNic")
	// FT logging nic is not configured on the host
	VmFaultToleranceConfigIssueReasonForIssueMissingFTLoggingNic = VmFaultToleranceConfigIssueReasonForIssue("missingFTLoggingNic")
	// The virtual machine has thin provisioned disks
	VmFaultToleranceConfigIssueReasonForIssueThinDisk = VmFaultToleranceConfigIssueReasonForIssue("thinDisk")
	// The "check host certificate" flag is not set
	VmFaultToleranceConfigIssueReasonForIssueVerifySSLCertificateFlagNotSet = VmFaultToleranceConfigIssueReasonForIssue("verifySSLCertificateFlagNotSet")
	// The virtual machine has one or more snapshots
	VmFaultToleranceConfigIssueReasonForIssueHasSnapshots = VmFaultToleranceConfigIssueReasonForIssue("hasSnapshots")
	// No configuration information is available for the virtual machine
	VmFaultToleranceConfigIssueReasonForIssueNoConfig = VmFaultToleranceConfigIssueReasonForIssue("noConfig")
	// The virtual machine is a fault tolerance secondary virtual machine
	VmFaultToleranceConfigIssueReasonForIssueFtSecondaryVm = VmFaultToleranceConfigIssueReasonForIssue("ftSecondaryVm")
	// The virtual machine has one or more disks on local datastore
	VmFaultToleranceConfigIssueReasonForIssueHasLocalDisk = VmFaultToleranceConfigIssueReasonForIssue("hasLocalDisk")
	// The virtual machine is an ESX agent VM
	VmFaultToleranceConfigIssueReasonForIssueEsxAgentVm = VmFaultToleranceConfigIssueReasonForIssue("esxAgentVm")
	// The virtual machine video device has 3D enabled
	VmFaultToleranceConfigIssueReasonForIssueVideo3dEnabled     = VmFaultToleranceConfigIssueReasonForIssue("video3dEnabled")
	VmFaultToleranceConfigIssueReasonForIssueHasUnsupportedDisk = VmFaultToleranceConfigIssueReasonForIssue("hasUnsupportedDisk")
	// FT logging nic does not have desired bandwidth
	VmFaultToleranceConfigIssueReasonForIssueInsufficientBandwidth = VmFaultToleranceConfigIssueReasonForIssue("insufficientBandwidth")
	// The host does not support fault tolerant VM with nested HV or VBS
	// enabled.
	VmFaultToleranceConfigIssueReasonForIssueHasNestedHVConfiguration = VmFaultToleranceConfigIssueReasonForIssue("hasNestedHVConfiguration")
	// The virtual machine has a vFlash memory device or/and disks with
	// vFlash cache configured.
	VmFaultToleranceConfigIssueReasonForIssueHasVFlashConfiguration = VmFaultToleranceConfigIssueReasonForIssue("hasVFlashConfiguration")
	// VMware product installed on the host does not support
	// fault tolerance
	VmFaultToleranceConfigIssueReasonForIssueUnsupportedProduct = VmFaultToleranceConfigIssueReasonForIssue("unsupportedProduct")
	// Host CPU does not support hardware virtualization
	VmFaultToleranceConfigIssueReasonForIssueCpuHvUnsupported = VmFaultToleranceConfigIssueReasonForIssue("cpuHvUnsupported")
	// Host CPU does not support hardware MMU virtualization
	VmFaultToleranceConfigIssueReasonForIssueCpuHwmmuUnsupported = VmFaultToleranceConfigIssueReasonForIssue("cpuHwmmuUnsupported")
	// Host CPU is compatible for replay-based FT, but hardware
	// virtualization has been disabled in the BIOS.
	VmFaultToleranceConfigIssueReasonForIssueCpuHvDisabled = VmFaultToleranceConfigIssueReasonForIssue("cpuHvDisabled")
	// The virtual machine firmware is of type EFI
	VmFaultToleranceConfigIssueReasonForIssueHasEFIFirmware = VmFaultToleranceConfigIssueReasonForIssue("hasEFIFirmware")
	// The host does not support fault tolerance virtual machines
	// with the specified number of virtual CPUs.
	VmFaultToleranceConfigIssueReasonForIssueTooManyVCPUs = VmFaultToleranceConfigIssueReasonForIssue("tooManyVCPUs")
	// The host does not support fault tolerance virtual machines
	// with the specified amount of memory.
	VmFaultToleranceConfigIssueReasonForIssueTooMuchMemory = VmFaultToleranceConfigIssueReasonForIssue("tooMuchMemory")
	// No VMotion license
	VmFaultToleranceConfigIssueReasonForIssueVMotionNotLicensed = VmFaultToleranceConfigIssueReasonForIssue("vMotionNotLicensed")
	// Host does not have proper FT license
	VmFaultToleranceConfigIssueReasonForIssueFtNotLicensed = VmFaultToleranceConfigIssueReasonForIssue("ftNotLicensed")
	// Host does not have HA agent running properly
	VmFaultToleranceConfigIssueReasonForIssueHaAgentIssue = VmFaultToleranceConfigIssueReasonForIssue("haAgentIssue")
	// The VM has unsupported storage policy
	VmFaultToleranceConfigIssueReasonForIssueUnsupportedSPBM = VmFaultToleranceConfigIssueReasonForIssue("unsupportedSPBM")
	// The virtual machine has virtual disk in linked-clone mode
	VmFaultToleranceConfigIssueReasonForIssueHasLinkedCloneDisk = VmFaultToleranceConfigIssueReasonForIssue("hasLinkedCloneDisk")
	// Virtual Machine with Pmem HA Failover is not supported
	VmFaultToleranceConfigIssueReasonForIssueUnsupportedPMemHAFailOver = VmFaultToleranceConfigIssueReasonForIssue("unsupportedPMemHAFailOver")
	// Virtual Machine with encrypted virtual disk is not supported.
	VmFaultToleranceConfigIssueReasonForIssueUnsupportedEncryptedDisk = VmFaultToleranceConfigIssueReasonForIssue("unsupportedEncryptedDisk")
	// The virtual machine does not allow to enable or disable FT Metro
	// Cluster while FT is turned on.
	VmFaultToleranceConfigIssueReasonForIssueFtMetroClusterNotEditable = VmFaultToleranceConfigIssueReasonForIssue("ftMetroClusterNotEditable")
	// Cannot turn on vSphere Fault Tolerance on a FT Metro Cluster enabled VM
	// with no Host Group configured.
	VmFaultToleranceConfigIssueReasonForIssueNoHostGroupConfigured = VmFaultToleranceConfigIssueReasonForIssue("noHostGroupConfigured")
)

func (e VmFaultToleranceConfigIssueReasonForIssue) Values() []VmFaultToleranceConfigIssueReasonForIssue {
	return []VmFaultToleranceConfigIssueReasonForIssue{
		VmFaultToleranceConfigIssueReasonForIssueHaNotEnabled,
		VmFaultToleranceConfigIssueReasonForIssueMoreThanOneSecondary,
		VmFaultToleranceConfigIssueReasonForIssueRecordReplayNotSupported,
		VmFaultToleranceConfigIssueReasonForIssueReplayNotSupported,
		VmFaultToleranceConfigIssueReasonForIssueTemplateVm,
		VmFaultToleranceConfigIssueReasonForIssueMultipleVCPU,
		VmFaultToleranceConfigIssueReasonForIssueHostInactive,
		VmFaultToleranceConfigIssueReasonForIssueFtUnsupportedHardware,
		VmFaultToleranceConfigIssueReasonForIssueFtUnsupportedProduct,
		VmFaultToleranceConfigIssueReasonForIssueMissingVMotionNic,
		VmFaultToleranceConfigIssueReasonForIssueMissingFTLoggingNic,
		VmFaultToleranceConfigIssueReasonForIssueThinDisk,
		VmFaultToleranceConfigIssueReasonForIssueVerifySSLCertificateFlagNotSet,
		VmFaultToleranceConfigIssueReasonForIssueHasSnapshots,
		VmFaultToleranceConfigIssueReasonForIssueNoConfig,
		VmFaultToleranceConfigIssueReasonForIssueFtSecondaryVm,
		VmFaultToleranceConfigIssueReasonForIssueHasLocalDisk,
		VmFaultToleranceConfigIssueReasonForIssueEsxAgentVm,
		VmFaultToleranceConfigIssueReasonForIssueVideo3dEnabled,
		VmFaultToleranceConfigIssueReasonForIssueHasUnsupportedDisk,
		VmFaultToleranceConfigIssueReasonForIssueInsufficientBandwidth,
		VmFaultToleranceConfigIssueReasonForIssueHasNestedHVConfiguration,
		VmFaultToleranceConfigIssueReasonForIssueHasVFlashConfiguration,
		VmFaultToleranceConfigIssueReasonForIssueUnsupportedProduct,
		VmFaultToleranceConfigIssueReasonForIssueCpuHvUnsupported,
		VmFaultToleranceConfigIssueReasonForIssueCpuHwmmuUnsupported,
		VmFaultToleranceConfigIssueReasonForIssueCpuHvDisabled,
		VmFaultToleranceConfigIssueReasonForIssueHasEFIFirmware,
		VmFaultToleranceConfigIssueReasonForIssueTooManyVCPUs,
		VmFaultToleranceConfigIssueReasonForIssueTooMuchMemory,
		VmFaultToleranceConfigIssueReasonForIssueVMotionNotLicensed,
		VmFaultToleranceConfigIssueReasonForIssueFtNotLicensed,
		VmFaultToleranceConfigIssueReasonForIssueHaAgentIssue,
		VmFaultToleranceConfigIssueReasonForIssueUnsupportedSPBM,
		VmFaultToleranceConfigIssueReasonForIssueHasLinkedCloneDisk,
		VmFaultToleranceConfigIssueReasonForIssueUnsupportedPMemHAFailOver,
		VmFaultToleranceConfigIssueReasonForIssueUnsupportedEncryptedDisk,
		VmFaultToleranceConfigIssueReasonForIssueFtMetroClusterNotEditable,
		VmFaultToleranceConfigIssueReasonForIssueNoHostGroupConfigured,
	}
}

func (e VmFaultToleranceConfigIssueReasonForIssue) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VmFaultToleranceConfigIssueReasonForIssue"] = reflect.TypeOf((*VmFaultToleranceConfigIssueReasonForIssue)(nil)).Elem()
	minAPIVersionForEnumValue["VmFaultToleranceConfigIssueReasonForIssue"] = map[string]string{
		"vMotionNotLicensed":        "8.0.3.0",
		"ftNotLicensed":             "8.0.3.0",
		"haAgentIssue":              "8.0.3.0",
		"unsupportedSPBM":           "8.0.3.0",
		"hasLinkedCloneDisk":        "8.0.3.0",
		"unsupportedPMemHAFailOver": "7.0.2.0",
		"unsupportedEncryptedDisk":  "8.0.3.0",
		"ftMetroClusterNotEditable": "8.0.3.0",
		"noHostGroupConfigured":     "8.0.3.0",
	}
}

type VmFaultToleranceInvalidFileBackingDeviceType string

const (
	// virtual floppy
	VmFaultToleranceInvalidFileBackingDeviceTypeVirtualFloppy = VmFaultToleranceInvalidFileBackingDeviceType("virtualFloppy")
	// virtual Cdrom
	VmFaultToleranceInvalidFileBackingDeviceTypeVirtualCdrom = VmFaultToleranceInvalidFileBackingDeviceType("virtualCdrom")
	// virtual serial port
	VmFaultToleranceInvalidFileBackingDeviceTypeVirtualSerialPort = VmFaultToleranceInvalidFileBackingDeviceType("virtualSerialPort")
	// virtual parallel port
	VmFaultToleranceInvalidFileBackingDeviceTypeVirtualParallelPort = VmFaultToleranceInvalidFileBackingDeviceType("virtualParallelPort")
	// virtual disk
	VmFaultToleranceInvalidFileBackingDeviceTypeVirtualDisk = VmFaultToleranceInvalidFileBackingDeviceType("virtualDisk")
)

func (e VmFaultToleranceInvalidFileBackingDeviceType) Values() []VmFaultToleranceInvalidFileBackingDeviceType {
	return []VmFaultToleranceInvalidFileBackingDeviceType{
		VmFaultToleranceInvalidFileBackingDeviceTypeVirtualFloppy,
		VmFaultToleranceInvalidFileBackingDeviceTypeVirtualCdrom,
		VmFaultToleranceInvalidFileBackingDeviceTypeVirtualSerialPort,
		VmFaultToleranceInvalidFileBackingDeviceTypeVirtualParallelPort,
		VmFaultToleranceInvalidFileBackingDeviceTypeVirtualDisk,
	}
}

func (e VmFaultToleranceInvalidFileBackingDeviceType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VmFaultToleranceInvalidFileBackingDeviceType"] = reflect.TypeOf((*VmFaultToleranceInvalidFileBackingDeviceType)(nil)).Elem()
}

type VmShutdownOnIsolationEventOperation string

const (
	// The virtual machine was shut down
	VmShutdownOnIsolationEventOperationShutdown = VmShutdownOnIsolationEventOperation("shutdown")
	// The virtual machine was powered off because shut down failed
	VmShutdownOnIsolationEventOperationPoweredOff = VmShutdownOnIsolationEventOperation("poweredOff")
)

func (e VmShutdownOnIsolationEventOperation) Values() []VmShutdownOnIsolationEventOperation {
	return []VmShutdownOnIsolationEventOperation{
		VmShutdownOnIsolationEventOperationShutdown,
		VmShutdownOnIsolationEventOperationPoweredOff,
	}
}

func (e VmShutdownOnIsolationEventOperation) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VmShutdownOnIsolationEventOperation"] = reflect.TypeOf((*VmShutdownOnIsolationEventOperation)(nil)).Elem()
}

// The PVLAN port types.
type VmwareDistributedVirtualSwitchPvlanPortType string

const (
	// The port can communicate with all other ports within the same PVLAN,
	// including the isolated and community ports .
	VmwareDistributedVirtualSwitchPvlanPortTypePromiscuous = VmwareDistributedVirtualSwitchPvlanPortType("promiscuous")
	// The port can only communicate with the promiscuous ports within the
	// same PVLAN, any other traffics are blocked.
	VmwareDistributedVirtualSwitchPvlanPortTypeIsolated = VmwareDistributedVirtualSwitchPvlanPortType("isolated")
	// The ports communicates with other community ports and with
	// promiscuous ports within the same PVLAN.
	//
	// any other traffics are
	// blocked.
	VmwareDistributedVirtualSwitchPvlanPortTypeCommunity = VmwareDistributedVirtualSwitchPvlanPortType("community")
)

func (e VmwareDistributedVirtualSwitchPvlanPortType) Values() []VmwareDistributedVirtualSwitchPvlanPortType {
	return []VmwareDistributedVirtualSwitchPvlanPortType{
		VmwareDistributedVirtualSwitchPvlanPortTypePromiscuous,
		VmwareDistributedVirtualSwitchPvlanPortTypeIsolated,
		VmwareDistributedVirtualSwitchPvlanPortTypeCommunity,
	}
}

func (e VmwareDistributedVirtualSwitchPvlanPortType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VmwareDistributedVirtualSwitchPvlanPortType"] = reflect.TypeOf((*VmwareDistributedVirtualSwitchPvlanPortType)(nil)).Elem()
}

// The list of disk issues.
type VsanDiskIssueType string

const (
	VsanDiskIssueTypeNonExist      = VsanDiskIssueType("nonExist")
	VsanDiskIssueTypeStampMismatch = VsanDiskIssueType("stampMismatch")
	VsanDiskIssueTypeUnknown       = VsanDiskIssueType("unknown")
)

func (e VsanDiskIssueType) Values() []VsanDiskIssueType {
	return []VsanDiskIssueType{
		VsanDiskIssueTypeNonExist,
		VsanDiskIssueTypeStampMismatch,
		VsanDiskIssueTypeUnknown,
	}
}

func (e VsanDiskIssueType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VsanDiskIssueType"] = reflect.TypeOf((*VsanDiskIssueType)(nil)).Elem()
}

// The action to take with regard to storage objects upon decommissioning
// a host from use with the VSAN service.
type VsanHostDecommissionModeObjectAction string

const (
	// No special action should take place regarding VSAN data.
	VsanHostDecommissionModeObjectActionNoAction = VsanHostDecommissionModeObjectAction("noAction")
	// VSAN data reconfiguration should be performed to ensure storage
	// object accessibility.
	VsanHostDecommissionModeObjectActionEnsureObjectAccessibility = VsanHostDecommissionModeObjectAction("ensureObjectAccessibility")
	// VSAN data evacuation should be performed such that all storage
	// object data is removed from the host.
	VsanHostDecommissionModeObjectActionEvacuateAllData = VsanHostDecommissionModeObjectAction("evacuateAllData")
)

func (e VsanHostDecommissionModeObjectAction) Values() []VsanHostDecommissionModeObjectAction {
	return []VsanHostDecommissionModeObjectAction{
		VsanHostDecommissionModeObjectActionNoAction,
		VsanHostDecommissionModeObjectActionEnsureObjectAccessibility,
		VsanHostDecommissionModeObjectActionEvacuateAllData,
	}
}

func (e VsanHostDecommissionModeObjectAction) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VsanHostDecommissionModeObjectAction"] = reflect.TypeOf((*VsanHostDecommissionModeObjectAction)(nil)).Elem()
}

// Values used for indicating a disk's status for use by the VSAN service.
//
// See also `VsanHostDiskResult.state`.
type VsanHostDiskResultState string

const (
	// Disk is currently in use by the VSAN service.
	//
	// A disk may be considered in use by the VSAN service regardless of
	// whether the VSAN service is enabled. As long as a disk is in use
	// by VSAN, it is reserved exclusively for VSAN and may not be used
	// for other purposes.
	//
	// See also `VsanHostDiskResult.error`.
	VsanHostDiskResultStateInUse = VsanHostDiskResultState("inUse")
	// Disk is considered eligible for use by the VSAN service,
	// but is not currently in use.
	VsanHostDiskResultStateEligible = VsanHostDiskResultState("eligible")
	// Disk is considered ineligible for use by the VSAN service,
	// and is not currently in use.
	//
	// See also `VsanHostDiskResult.error`.
	VsanHostDiskResultStateIneligible = VsanHostDiskResultState("ineligible")
)

func (e VsanHostDiskResultState) Values() []VsanHostDiskResultState {
	return []VsanHostDiskResultState{
		VsanHostDiskResultStateInUse,
		VsanHostDiskResultStateEligible,
		VsanHostDiskResultStateIneligible,
	}
}

func (e VsanHostDiskResultState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VsanHostDiskResultState"] = reflect.TypeOf((*VsanHostDiskResultState)(nil)).Elem()
}

// A `VsanHostHealthState_enum` represents the state of a participating
// host in the VSAN service.
//
// See also `VsanHostClusterStatus`.
type VsanHostHealthState string

const (
	// Node health is unknown.
	VsanHostHealthStateUnknown = VsanHostHealthState("unknown")
	// Node is considered healthy.
	VsanHostHealthStateHealthy = VsanHostHealthState("healthy")
	// Node is considered unhealthy.
	VsanHostHealthStateUnhealthy = VsanHostHealthState("unhealthy")
)

func (e VsanHostHealthState) Values() []VsanHostHealthState {
	return []VsanHostHealthState{
		VsanHostHealthStateUnknown,
		VsanHostHealthStateHealthy,
		VsanHostHealthStateUnhealthy,
	}
}

func (e VsanHostHealthState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VsanHostHealthState"] = reflect.TypeOf((*VsanHostHealthState)(nil)).Elem()
}

// A `VsanHostNodeState_enum` represents the state of participation of a host
// in the VSAN service.
//
// See also `VsanHostClusterStatus`, `VsanHostClusterStatusState`.
type VsanHostNodeState string

const (
	// The node is enabled for the VSAN service but has some configuration
	// error which prevents participation.
	VsanHostNodeStateError = VsanHostNodeState("error")
	// The node is disabled for the VSAN service.
	VsanHostNodeStateDisabled = VsanHostNodeState("disabled")
	// The node is enabled for the VSAN service and is serving as an agent.
	VsanHostNodeStateAgent = VsanHostNodeState("agent")
	// The node is enabled for the VSAN service and is serving as the master.
	VsanHostNodeStateMaster = VsanHostNodeState("master")
	// The node is enabled for the VSAN service and is serving as the backup.
	VsanHostNodeStateBackup = VsanHostNodeState("backup")
	// The node is starting the VSAN service; this state is considered
	// transitory.
	VsanHostNodeStateStarting = VsanHostNodeState("starting")
	// The node is stopping the VSAN service; this state is considered
	// transitory.
	VsanHostNodeStateStopping = VsanHostNodeState("stopping")
	// The node is entering maintenance mode; this state is considered
	// transitory.
	//
	// See also `HostSystem.EnterMaintenanceMode_Task`.
	VsanHostNodeStateEnteringMaintenanceMode = VsanHostNodeState("enteringMaintenanceMode")
	// The node is exiting maintenance mode; this state is considered
	// transitory.
	//
	// See also `HostSystem.ExitMaintenanceMode_Task`.
	VsanHostNodeStateExitingMaintenanceMode = VsanHostNodeState("exitingMaintenanceMode")
	// The node is being decommissioned from the VSAN service; this state is
	// considered transitory.
	VsanHostNodeStateDecommissioning = VsanHostNodeState("decommissioning")
)

func (e VsanHostNodeState) Values() []VsanHostNodeState {
	return []VsanHostNodeState{
		VsanHostNodeStateError,
		VsanHostNodeStateDisabled,
		VsanHostNodeStateAgent,
		VsanHostNodeStateMaster,
		VsanHostNodeStateBackup,
		VsanHostNodeStateStarting,
		VsanHostNodeStateStopping,
		VsanHostNodeStateEnteringMaintenanceMode,
		VsanHostNodeStateExitingMaintenanceMode,
		VsanHostNodeStateDecommissioning,
	}
}

func (e VsanHostNodeState) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VsanHostNodeState"] = reflect.TypeOf((*VsanHostNodeState)(nil)).Elem()
}

// Type of disk group operation performed.
type VsanUpgradeSystemUpgradeHistoryDiskGroupOpType string

const (
	// Disk group is being (re-)added.
	VsanUpgradeSystemUpgradeHistoryDiskGroupOpTypeAdd = VsanUpgradeSystemUpgradeHistoryDiskGroupOpType("add")
	// Disk group is being removed.
	VsanUpgradeSystemUpgradeHistoryDiskGroupOpTypeRemove = VsanUpgradeSystemUpgradeHistoryDiskGroupOpType("remove")
)

func (e VsanUpgradeSystemUpgradeHistoryDiskGroupOpType) Values() []VsanUpgradeSystemUpgradeHistoryDiskGroupOpType {
	return []VsanUpgradeSystemUpgradeHistoryDiskGroupOpType{
		VsanUpgradeSystemUpgradeHistoryDiskGroupOpTypeAdd,
		VsanUpgradeSystemUpgradeHistoryDiskGroupOpTypeRemove,
	}
}

func (e VsanUpgradeSystemUpgradeHistoryDiskGroupOpType) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["VsanUpgradeSystemUpgradeHistoryDiskGroupOpType"] = reflect.TypeOf((*VsanUpgradeSystemUpgradeHistoryDiskGroupOpType)(nil)).Elem()
}

type WeekOfMonth string

const (
	WeekOfMonthFirst  = WeekOfMonth("first")
	WeekOfMonthSecond = WeekOfMonth("second")
	WeekOfMonthThird  = WeekOfMonth("third")
	WeekOfMonthFourth = WeekOfMonth("fourth")
	WeekOfMonthLast   = WeekOfMonth("last")
)

func (e WeekOfMonth) Values() []WeekOfMonth {
	return []WeekOfMonth{
		WeekOfMonthFirst,
		WeekOfMonthSecond,
		WeekOfMonthThird,
		WeekOfMonthFourth,
		WeekOfMonthLast,
	}
}

func (e WeekOfMonth) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["WeekOfMonth"] = reflect.TypeOf((*WeekOfMonth)(nil)).Elem()
}

type WillLoseHAProtectionResolution string

const (
	// storage vmotion resolution
	WillLoseHAProtectionResolutionSvmotion = WillLoseHAProtectionResolution("svmotion")
	// relocate resolution
	WillLoseHAProtectionResolutionRelocate = WillLoseHAProtectionResolution("relocate")
)

func (e WillLoseHAProtectionResolution) Values() []WillLoseHAProtectionResolution {
	return []WillLoseHAProtectionResolution{
		WillLoseHAProtectionResolutionSvmotion,
		WillLoseHAProtectionResolutionRelocate,
	}
}

func (e WillLoseHAProtectionResolution) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["WillLoseHAProtectionResolution"] = reflect.TypeOf((*WillLoseHAProtectionResolution)(nil)).Elem()
}

type VslmDiskInfoFlag string

const (
	VslmDiskInfoFlagId                      = VslmDiskInfoFlag("id")
	VslmDiskInfoFlagDescriptorVersion       = VslmDiskInfoFlag("descriptorVersion")
	VslmDiskInfoFlagBackingObjectId         = VslmDiskInfoFlag("backingObjectId")
	VslmDiskInfoFlagPath                    = VslmDiskInfoFlag("path")
	VslmDiskInfoFlagParentPath              = VslmDiskInfoFlag("parentPath")
	VslmDiskInfoFlagName                    = VslmDiskInfoFlag("name")
	VslmDiskInfoFlagDeviceName              = VslmDiskInfoFlag("deviceName")
	VslmDiskInfoFlagCapacity                = VslmDiskInfoFlag("capacity")
	VslmDiskInfoFlagAllocated               = VslmDiskInfoFlag("allocated")
	VslmDiskInfoFlagType                    = VslmDiskInfoFlag("type")
	VslmDiskInfoFlagConsumers               = VslmDiskInfoFlag("consumers")
	VslmDiskInfoFlagTentativeState          = VslmDiskInfoFlag("tentativeState")
	VslmDiskInfoFlagCreateTime              = VslmDiskInfoFlag("createTime")
	VslmDiskInfoFlagIoFilter                = VslmDiskInfoFlag("ioFilter")
	VslmDiskInfoFlagControlFlags            = VslmDiskInfoFlag("controlFlags")
	VslmDiskInfoFlagKeepAfterVmDelete       = VslmDiskInfoFlag("keepAfterVmDelete")
	VslmDiskInfoFlagRelocationDisabled      = VslmDiskInfoFlag("relocationDisabled")
	VslmDiskInfoFlagKeyId                   = VslmDiskInfoFlag("keyId")
	VslmDiskInfoFlagKeyProviderId           = VslmDiskInfoFlag("keyProviderId")
	VslmDiskInfoFlagNativeSnapshotSupported = VslmDiskInfoFlag("nativeSnapshotSupported")
	VslmDiskInfoFlagCbtEnabled              = VslmDiskInfoFlag("cbtEnabled")
)

func (e VslmDiskInfoFlag) Values() []VslmDiskInfoFlag {
	return []VslmDiskInfoFlag{
		VslmDiskInfoFlagId,
		VslmDiskInfoFlagDescriptorVersion,
		VslmDiskInfoFlagBackingObjectId,
		VslmDiskInfoFlagPath,
		VslmDiskInfoFlagParentPath,
		VslmDiskInfoFlagName,
		VslmDiskInfoFlagDeviceName,
		VslmDiskInfoFlagCapacity,
		VslmDiskInfoFlagAllocated,
		VslmDiskInfoFlagType,
		VslmDiskInfoFlagConsumers,
		VslmDiskInfoFlagTentativeState,
		VslmDiskInfoFlagCreateTime,
		VslmDiskInfoFlagIoFilter,
		VslmDiskInfoFlagControlFlags,
		VslmDiskInfoFlagKeepAfterVmDelete,
		VslmDiskInfoFlagRelocationDisabled,
		VslmDiskInfoFlagKeyId,
		VslmDiskInfoFlagKeyProviderId,
		VslmDiskInfoFlagNativeSnapshotSupported,
		VslmDiskInfoFlagCbtEnabled,
	}
}

func (e VslmDiskInfoFlag) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["vslmDiskInfoFlag"] = reflect.TypeOf((*VslmDiskInfoFlag)(nil)).Elem()
}

type VslmVStorageObjectControlFlag string

const (
	VslmVStorageObjectControlFlagKeepAfterDeleteVm          = VslmVStorageObjectControlFlag("keepAfterDeleteVm")
	VslmVStorageObjectControlFlagDisableRelocation          = VslmVStorageObjectControlFlag("disableRelocation")
	VslmVStorageObjectControlFlagEnableChangedBlockTracking = VslmVStorageObjectControlFlag("enableChangedBlockTracking")
)

func (e VslmVStorageObjectControlFlag) Values() []VslmVStorageObjectControlFlag {
	return []VslmVStorageObjectControlFlag{
		VslmVStorageObjectControlFlagKeepAfterDeleteVm,
		VslmVStorageObjectControlFlagDisableRelocation,
		VslmVStorageObjectControlFlagEnableChangedBlockTracking,
	}
}

func (e VslmVStorageObjectControlFlag) Strings() []string {
	return EnumValuesAsStrings(e.Values())
}

func init() {
	t["vslmVStorageObjectControlFlag"] = reflect.TypeOf((*VslmVStorageObjectControlFlag)(nil)).Elem()
}
