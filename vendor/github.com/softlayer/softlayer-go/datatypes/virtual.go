/**
 * Copyright 2016-2024 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed
 * on an "AS IS" BASIS,WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// AUTOMATICALLY GENERATED CODE - DO NOT MODIFY

package datatypes

// The SoftLayer_Virtual_BrowserConsoleAccessLog data type presents the data for events associated with accessing a VSIs console via the browser interface.
type Virtual_BrowserConsoleAccessLog struct {
	Entity

	// Date created
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The event type that was triggered
	EventType *string `json:"eventType,omitempty" xmlrpc:"eventType,omitempty"`

	// Unique ID for a console log event instance.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The message for the event
	Message *string `json:"message,omitempty" xmlrpc:"message,omitempty"`

	// Date modified
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The IP address of the client that triggered the event
	SourceIp *string `json:"sourceIp,omitempty" xmlrpc:"sourceIp,omitempty"`

	// The IP port of the client that triggered the event
	SourcePort *int `json:"sourcePort,omitempty" xmlrpc:"sourcePort,omitempty"`

	// The user id that triggered the event
	UserId *int `json:"userId,omitempty" xmlrpc:"userId,omitempty"`

	// The user type that triggered the event
	UserType *string `json:"userType,omitempty" xmlrpc:"userType,omitempty"`

	// The username of the user that triggered the event
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// This data type presents the structure for a dedicated host. The data type contains relational properties to distinguish a dedicated host and associate an account to it.
type Virtual_DedicatedHost struct {
	Entity

	// The account that the dedicated host belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The dedicated host's associated account id.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The container that represents allocations on the dedicated host.
	AllocationStatus *Container_Virtual_DedicatedHost_AllocationStatus `json:"allocationStatus,omitempty" xmlrpc:"allocationStatus,omitempty"`

	// The backend router behind dedicated host's pool of resources.
	BackendRouter *Hardware_Router_Backend `json:"backendRouter,omitempty" xmlrpc:"backendRouter,omitempty"`

	// The billing item for the dedicated host.
	BillingItem *Billing_Item_Virtual_DedicatedHost `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The capacity that the dedicated host's CPU allocation is restricted to.
	CpuCount *int `json:"cpuCount,omitempty" xmlrpc:"cpuCount,omitempty"`

	// The date that the dedicated host was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The datacenter that the dedicated host resides in.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// The capacity that the dedicated host's disk allocation is restricted to.
	DiskCapacity *int `json:"diskCapacity,omitempty" xmlrpc:"diskCapacity,omitempty"`

	// A count of the guests associated with the dedicated host.
	GuestCount *uint `json:"guestCount,omitempty" xmlrpc:"guestCount,omitempty"`

	// The guests associated with the dedicated host.
	Guests []Virtual_Guest `json:"guests,omitempty" xmlrpc:"guests,omitempty"`

	// The dedicated host's associated unique id.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of
	InternalTagReferenceCount *uint `json:"internalTagReferenceCount,omitempty" xmlrpc:"internalTagReferenceCount,omitempty"`

	// no documentation yet
	InternalTagReferences []Tag_Reference `json:"internalTagReferences,omitempty" xmlrpc:"internalTagReferences,omitempty"`

	// The capacity that the dedicated host's memory allocation is restricted to.
	MemoryCapacity *int `json:"memoryCapacity,omitempty" xmlrpc:"memoryCapacity,omitempty"`

	// The date that the dedicated host was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The dedicated host's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A note of up to 1,000 characters about a dedicated host.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// The container that represents PCI device allocations on the dedicated host.
	PciDeviceAllocationStatus *Container_Virtual_DedicatedHost_Pci_Device_AllocationStatus `json:"pciDeviceAllocationStatus,omitempty" xmlrpc:"pciDeviceAllocationStatus,omitempty"`

	// A count of a collection of SoftLayer_Virtual_Host_PciDevice objects on the host.
	PciDeviceCount *uint `json:"pciDeviceCount,omitempty" xmlrpc:"pciDeviceCount,omitempty"`

	// A collection of SoftLayer_Virtual_Host_PciDevice objects on the host.
	PciDevices []Virtual_Host_PciDevice `json:"pciDevices,omitempty" xmlrpc:"pciDevices,omitempty"`

	// A count of
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// no documentation yet
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`
}

// The virtual disk image data type presents the structure in which a virtual disk image will be presented.
//
// Virtual block devices are assigned to disk images.
type Virtual_Disk_Image struct {
	Entity

	// The billing item for a virtual disk image.
	BillingItem *Billing_Item_Virtual_Disk_Image `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// A count of the block devices that a disk image is attached to. Block devices connect computing instances to disk images.
	BlockDeviceCount *uint `json:"blockDeviceCount,omitempty" xmlrpc:"blockDeviceCount,omitempty"`

	// The block devices that a disk image is attached to. Block devices connect computing instances to disk images.
	BlockDevices []Virtual_Guest_Block_Device `json:"blockDevices,omitempty" xmlrpc:"blockDevices,omitempty"`

	// no documentation yet
	BootableVolumeFlag *bool `json:"bootableVolumeFlag,omitempty" xmlrpc:"bootableVolumeFlag,omitempty"`

	// A disk image's size measured in gigabytes.
	Capacity *int `json:"capacity,omitempty" xmlrpc:"capacity,omitempty"`

	// A disk image's unique md5 checksum.
	Checksum *string `json:"checksum,omitempty" xmlrpc:"checksum,omitempty"`

	// A column to save the Algorithm being used to store the checksum.
	ChecksumAlgorithm *string `json:"checksumAlgorithm,omitempty" xmlrpc:"checksumAlgorithm,omitempty"`

	// Check if cloud-init is enabled.
	CloudInitFlag *bool `json:"cloudInitFlag,omitempty" xmlrpc:"cloudInitFlag,omitempty"`

	// A count of
	CoalescedDiskImageCount *uint `json:"coalescedDiskImageCount,omitempty" xmlrpc:"coalescedDiskImageCount,omitempty"`

	// no documentation yet
	CoalescedDiskImages []Virtual_Disk_Image `json:"coalescedDiskImages,omitempty" xmlrpc:"coalescedDiskImages,omitempty"`

	// no documentation yet
	CopyOnWriteFlag *bool `json:"copyOnWriteFlag,omitempty" xmlrpc:"copyOnWriteFlag,omitempty"`

	// The date a disk image was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A brief description of a virtual disk image.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Return disk file extension
	DiskFileExtension *string `json:"diskFileExtension,omitempty" xmlrpc:"diskFileExtension,omitempty"`

	// no documentation yet
	DiskImageStorageGroup *Configuration_Storage_Group `json:"diskImageStorageGroup,omitempty" xmlrpc:"diskImageStorageGroup,omitempty"`

	// A disk image's unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// Return imported disk type
	ImportedDiskType *string `json:"importedDiskType,omitempty" xmlrpc:"importedDiskType,omitempty"`

	// Return if image is encrypted
	IsEncrypted *bool `json:"isEncrypted,omitempty" xmlrpc:"isEncrypted,omitempty"`

	// no documentation yet
	LocalDiskFlag *bool `json:"localDiskFlag,omitempty" xmlrpc:"localDiskFlag,omitempty"`

	// Whether this disk image is meant for storage of custom user data supplied with a Cloud Computing Instance order.
	MetadataFlag *bool `json:"metadataFlag,omitempty" xmlrpc:"metadataFlag,omitempty"`

	// The date a disk image was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// A descriptive name used to identify a disk image to a user.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The ID of the the disk image that this disk image is based on, if applicable.
	ParentId *int `json:"parentId,omitempty" xmlrpc:"parentId,omitempty"`

	// A count of references to the software that resides on a disk image.
	SoftwareReferenceCount *uint `json:"softwareReferenceCount,omitempty" xmlrpc:"softwareReferenceCount,omitempty"`

	// References to the software that resides on a disk image.
	SoftwareReferences []Virtual_Disk_Image_Software `json:"softwareReferences,omitempty" xmlrpc:"softwareReferences,omitempty"`

	// The original disk image that the current disk image was cloned from.
	SourceDiskImage *Virtual_Disk_Image `json:"sourceDiskImage,omitempty" xmlrpc:"sourceDiskImage,omitempty"`

	// A count of the storage group for a virtual disk image.
	StorageGroupCount *uint `json:"storageGroupCount,omitempty" xmlrpc:"storageGroupCount,omitempty"`

	// Return storage group details for symantec disk
	StorageGroupDetails *Container_Image_StorageGroupDetails `json:"storageGroupDetails,omitempty" xmlrpc:"storageGroupDetails,omitempty"`

	// The storage group for a virtual disk image.
	StorageGroups []Configuration_Storage_Group `json:"storageGroups,omitempty" xmlrpc:"storageGroups,omitempty"`

	// The storage repository that a disk image resides in.
	StorageRepository *Virtual_Storage_Repository `json:"storageRepository,omitempty" xmlrpc:"storageRepository,omitempty"`

	// The [[SoftLayer_Virtual_Storage_Repository|storage repository]] that a disk image is in.
	StorageRepositoryId *int `json:"storageRepositoryId,omitempty" xmlrpc:"storageRepositoryId,omitempty"`

	// The type of storage repository that a disk image resides in.
	StorageRepositoryType *Virtual_Storage_Repository_Type `json:"storageRepositoryType,omitempty" xmlrpc:"storageRepositoryType,omitempty"`

	// Return supported hardware component IDs for symantec disk
	SupportedHardware *string `json:"supportedHardware,omitempty" xmlrpc:"supportedHardware,omitempty"`

	// The template that attaches a disk image to a [[SoftLayer_Virtual_Guest_Block_Device_Template_Group|archive]].
	TemplateBlockDevice *Virtual_Guest_Block_Device_Template `json:"templateBlockDevice,omitempty" xmlrpc:"templateBlockDevice,omitempty"`

	// A virtual disk image's type.
	Type *Virtual_Disk_Image_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// A disk image's [[SoftLayer_Virtual_Disk_Image_Type|type]] ID
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// The unit of storage in which the size of the image is measured. Defaults to "GB" for gigabytes.
	Units *string `json:"units,omitempty" xmlrpc:"units,omitempty"`

	// A disk image's unique ID on a virtualization platform.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// A SoftLayer_Virtual_Disk_Image_Software record connects a computing instance's virtual disk images with software records. This can be useful if a disk image is directly associated with software such as operating systems.
type Virtual_Disk_Image_Software struct {
	Entity

	// The virtual disk image that is associated with software.
	DiskImage *Virtual_Disk_Image `json:"diskImage,omitempty" xmlrpc:"diskImage,omitempty"`

	// The unique identifier of a virtual disk image to software relationship.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of username/Password pairs used for access to a Software Installation.
	PasswordCount *uint `json:"passwordCount,omitempty" xmlrpc:"passwordCount,omitempty"`

	// Username/Password pairs used for access to a Software Installation.
	Passwords []Virtual_Disk_Image_Software_Password `json:"passwords,omitempty" xmlrpc:"passwords,omitempty"`

	// The software associated with a virtual disk image.
	SoftwareDescription *Software_Description `json:"softwareDescription,omitempty" xmlrpc:"softwareDescription,omitempty"`

	// The unique identifier of the software that a virtual disk image is associated with.
	SoftwareDescriptionId *int `json:"softwareDescriptionId,omitempty" xmlrpc:"softwareDescriptionId,omitempty"`
}

// This SoftLayer_Virtual_Disk_Image_Software_Password data type contains a password for a specific virtual disk image software instance.
type Virtual_Disk_Image_Software_Password struct {
	Entity

	// A virtual disk images' password.
	Password *string `json:"password,omitempty" xmlrpc:"password,omitempty"`

	// The instance that this username/password pair is valid for.
	Software *Virtual_Disk_Image_Software `json:"software,omitempty" xmlrpc:"software,omitempty"`

	// A virtual disk images' username.
	Username *string `json:"username,omitempty" xmlrpc:"username,omitempty"`
}

// SoftLayer_Virtual_Disk_Image_Type models the types of virtual disk images available to CloudLayer Computing Instances. Virtual disk image types describe if an image's data is preservable when upgraded, whether a disk contains a suspended virtual image, or if a disk contains crash dump information.
type Virtual_Disk_Image_Type struct {
	Entity

	// A brief description of a virtual disk image type's function.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A virtual disk image type's key name.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// A virtual disk image type's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The virtual guest data type presents the structure in which all virtual guests will be presented. Internally, the structure supports various virtualization platforms with no change to external interaction.
//
// A guest, also known as a virtual server, represents an allocation of resources on a virtual host.
type Virtual_Guest struct {
	Entity

	// The account that a virtual guest belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A computing instance's associated [[SoftLayer_Account|account]] id
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// no documentation yet
	AccountOwnedPoolFlag *bool `json:"accountOwnedPoolFlag,omitempty" xmlrpc:"accountOwnedPoolFlag,omitempty"`

	// A virtual guest's currently active network monitoring incidents.
	ActiveNetworkMonitorIncident []Network_Monitor_Version1_Incident `json:"activeNetworkMonitorIncident,omitempty" xmlrpc:"activeNetworkMonitorIncident,omitempty"`

	// A count of a virtual guest's currently active network monitoring incidents.
	ActiveNetworkMonitorIncidentCount *uint `json:"activeNetworkMonitorIncidentCount,omitempty" xmlrpc:"activeNetworkMonitorIncidentCount,omitempty"`

	// A count of
	ActiveTicketCount *uint `json:"activeTicketCount,omitempty" xmlrpc:"activeTicketCount,omitempty"`

	// no documentation yet
	ActiveTickets []Ticket `json:"activeTickets,omitempty" xmlrpc:"activeTickets,omitempty"`

	// A transaction that is still be performed on a cloud server.
	ActiveTransaction *Provisioning_Version1_Transaction `json:"activeTransaction,omitempty" xmlrpc:"activeTransaction,omitempty"`

	// A count of any active transaction(s) that are currently running for the server (example: os reload).
	ActiveTransactionCount *uint `json:"activeTransactionCount,omitempty" xmlrpc:"activeTransactionCount,omitempty"`

	// Any active transaction(s) that are currently running for the server (example: os reload).
	ActiveTransactions []Provisioning_Version1_Transaction `json:"activeTransactions,omitempty" xmlrpc:"activeTransactions,omitempty"`

	// The SoftLayer_Network_Storage_Allowed_Host information to connect this Virtual Guest to Network Storage volumes that require access control lists.
	AllowedHost *Network_Storage_Allowed_Host `json:"allowedHost,omitempty" xmlrpc:"allowedHost,omitempty"`

	// The SoftLayer_Network_Storage objects that this SoftLayer_Virtual_Guest has access to.
	AllowedNetworkStorage []Network_Storage `json:"allowedNetworkStorage,omitempty" xmlrpc:"allowedNetworkStorage,omitempty"`

	// A count of the SoftLayer_Network_Storage objects that this SoftLayer_Virtual_Guest has access to.
	AllowedNetworkStorageCount *uint `json:"allowedNetworkStorageCount,omitempty" xmlrpc:"allowedNetworkStorageCount,omitempty"`

	// A count of the SoftLayer_Network_Storage objects whose Replica that this SoftLayer_Virtual_Guest has access to.
	AllowedNetworkStorageReplicaCount *uint `json:"allowedNetworkStorageReplicaCount,omitempty" xmlrpc:"allowedNetworkStorageReplicaCount,omitempty"`

	// The SoftLayer_Network_Storage objects whose Replica that this SoftLayer_Virtual_Guest has access to.
	AllowedNetworkStorageReplicas []Network_Storage `json:"allowedNetworkStorageReplicas,omitempty" xmlrpc:"allowedNetworkStorageReplicas,omitempty"`

	// A antivirus / spyware software component object.
	AntivirusSpywareSoftwareComponent *Software_Component `json:"antivirusSpywareSoftwareComponent,omitempty" xmlrpc:"antivirusSpywareSoftwareComponent,omitempty"`

	// no documentation yet
	ApplicationDeliveryController *Network_Application_Delivery_Controller `json:"applicationDeliveryController,omitempty" xmlrpc:"applicationDeliveryController,omitempty"`

	// A count of
	AttributeCount *uint `json:"attributeCount,omitempty" xmlrpc:"attributeCount,omitempty"`

	// no documentation yet
	Attributes []Virtual_Guest_Attribute `json:"attributes,omitempty" xmlrpc:"attributes,omitempty"`

	// An object that stores the maximum level for the monitoring query types and response types.
	AvailableMonitoring []Network_Monitor_Version1_Query_Host_Stratum `json:"availableMonitoring,omitempty" xmlrpc:"availableMonitoring,omitempty"`

	// A count of an object that stores the maximum level for the monitoring query types and response types.
	AvailableMonitoringCount *uint `json:"availableMonitoringCount,omitempty" xmlrpc:"availableMonitoringCount,omitempty"`

	// The average daily private bandwidth usage for the current billing cycle.
	AverageDailyPrivateBandwidthUsage *Float64 `json:"averageDailyPrivateBandwidthUsage,omitempty" xmlrpc:"averageDailyPrivateBandwidthUsage,omitempty"`

	// The average daily public bandwidth usage for the current billing cycle.
	AverageDailyPublicBandwidthUsage *Float64 `json:"averageDailyPublicBandwidthUsage,omitempty" xmlrpc:"averageDailyPublicBandwidthUsage,omitempty"`

	// A count of a guests's backend network components.
	BackendNetworkComponentCount *uint `json:"backendNetworkComponentCount,omitempty" xmlrpc:"backendNetworkComponentCount,omitempty"`

	// A guests's backend network components.
	BackendNetworkComponents []Virtual_Guest_Network_Component `json:"backendNetworkComponents,omitempty" xmlrpc:"backendNetworkComponents,omitempty"`

	// A count of a guest's backend or private router.
	BackendRouterCount *uint `json:"backendRouterCount,omitempty" xmlrpc:"backendRouterCount,omitempty"`

	// A guest's backend or private router.
	BackendRouters []Hardware `json:"backendRouters,omitempty" xmlrpc:"backendRouters,omitempty"`

	// A computing instance's allotted bandwidth (measured in GB).
	BandwidthAllocation *Float64 `json:"bandwidthAllocation,omitempty" xmlrpc:"bandwidthAllocation,omitempty"`

	// A computing instance's allotted detail record. Allotment details link bandwidth allocation with allotments.
	BandwidthAllotmentDetail *Network_Bandwidth_Version1_Allotment_Detail `json:"bandwidthAllotmentDetail,omitempty" xmlrpc:"bandwidthAllotmentDetail,omitempty"`

	// The raw bandwidth usage data for the current billing cycle. One object will be returned for each network this server is attached to.
	BillingCycleBandwidthUsage []Network_Bandwidth_Usage `json:"billingCycleBandwidthUsage,omitempty" xmlrpc:"billingCycleBandwidthUsage,omitempty"`

	// A count of the raw bandwidth usage data for the current billing cycle. One object will be returned for each network this server is attached to.
	BillingCycleBandwidthUsageCount *uint `json:"billingCycleBandwidthUsageCount,omitempty" xmlrpc:"billingCycleBandwidthUsageCount,omitempty"`

	// The raw private bandwidth usage data for the current billing cycle.
	BillingCyclePrivateBandwidthUsage *Network_Bandwidth_Usage `json:"billingCyclePrivateBandwidthUsage,omitempty" xmlrpc:"billingCyclePrivateBandwidthUsage,omitempty"`

	// The raw public bandwidth usage data for the current billing cycle.
	BillingCyclePublicBandwidthUsage *Network_Bandwidth_Usage `json:"billingCyclePublicBandwidthUsage,omitempty" xmlrpc:"billingCyclePublicBandwidthUsage,omitempty"`

	// The billing item for a CloudLayer Compute Instance.
	BillingItem *Billing_Item_Virtual_Guest `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// Determines whether the instance is ineligible for cancellation because it is disconnected.
	BlockCancelBecauseDisconnectedFlag *bool `json:"blockCancelBecauseDisconnectedFlag,omitempty" xmlrpc:"blockCancelBecauseDisconnectedFlag,omitempty"`

	// A count of a computing instance's block devices. Block devices link [[SoftLayer_Virtual_Disk_Image|disk images]] to computing instances.
	BlockDeviceCount *uint `json:"blockDeviceCount,omitempty" xmlrpc:"blockDeviceCount,omitempty"`

	// The global identifier for the image template that was used to provision or reload a guest.
	BlockDeviceTemplateGroup *Virtual_Guest_Block_Device_Template_Group `json:"blockDeviceTemplateGroup,omitempty" xmlrpc:"blockDeviceTemplateGroup,omitempty"`

	// A computing instance's block devices. Block devices link [[SoftLayer_Virtual_Disk_Image|disk images]] to computing instances.
	BlockDevices []Virtual_Guest_Block_Device `json:"blockDevices,omitempty" xmlrpc:"blockDevices,omitempty"`

	// A count of a virtual guest's browser access logs.
	BrowserConsoleAccessLogCount *uint `json:"browserConsoleAccessLogCount,omitempty" xmlrpc:"browserConsoleAccessLogCount,omitempty"`

	// A virtual guest's browser access logs.
	BrowserConsoleAccessLogs []Virtual_BrowserConsoleAccessLog `json:"browserConsoleAccessLogs,omitempty" xmlrpc:"browserConsoleAccessLogs,omitempty"`

	// A container for a guest's console data
	ConsoleData *Container_Virtual_ConsoleData `json:"consoleData,omitempty" xmlrpc:"consoleData,omitempty"`

	// [DEPRECATED] A flag indicating a computing instance's console IP address is assigned.
	// Deprecated: This function has been marked as deprecated.
	ConsoleIpAddressFlag *bool `json:"consoleIpAddressFlag,omitempty" xmlrpc:"consoleIpAddressFlag,omitempty"`

	// [DEPRECATED] A record containing information about a computing instance's console IP and port number.
	// Deprecated: This function has been marked as deprecated.
	ConsoleIpAddressRecord *Virtual_Guest_Network_Component_IpAddress `json:"consoleIpAddressRecord,omitempty" xmlrpc:"consoleIpAddressRecord,omitempty"`

	// A continuous data protection software component object.
	ContinuousDataProtectionSoftwareComponent *Software_Component `json:"continuousDataProtectionSoftwareComponent,omitempty" xmlrpc:"continuousDataProtectionSoftwareComponent,omitempty"`

	// A guest's control panel.
	ControlPanel *Software_Component `json:"controlPanel,omitempty" xmlrpc:"controlPanel,omitempty"`

	// The date a virtual computing instance was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// An object that provides commonly used bandwidth summary components for the current billing cycle.
	CurrentBandwidthSummary *Metric_Tracking_Object_Bandwidth_Summary `json:"currentBandwidthSummary,omitempty" xmlrpc:"currentBandwidthSummary,omitempty"`

	// The datacenter that a virtual guest resides in.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// When true this flag specifies that a compute instance is to run on hosts that only have guests from the same account.
	DedicatedAccountHostOnlyFlag *bool `json:"dedicatedAccountHostOnlyFlag,omitempty" xmlrpc:"dedicatedAccountHostOnlyFlag,omitempty"`

	// The dedicated host associated with this guest.
	DedicatedHost *Virtual_DedicatedHost `json:"dedicatedHost,omitempty" xmlrpc:"dedicatedHost,omitempty"`

	// The device status of this virtual guest.
	DeviceStatus *Device_Status `json:"deviceStatus,omitempty" xmlrpc:"deviceStatus,omitempty"`

	// The device status ID of the virtual guest.
	DeviceStatusId *int `json:"deviceStatusId,omitempty" xmlrpc:"deviceStatusId,omitempty"`

	// A computing instance's domain name
	Domain *string `json:"domain,omitempty" xmlrpc:"domain,omitempty"`

	// A guest's associated EVault network storage service account.
	EvaultNetworkStorage []Network_Storage `json:"evaultNetworkStorage,omitempty" xmlrpc:"evaultNetworkStorage,omitempty"`

	// A count of a guest's associated EVault network storage service account.
	EvaultNetworkStorageCount *uint `json:"evaultNetworkStorageCount,omitempty" xmlrpc:"evaultNetworkStorageCount,omitempty"`

	// A computing instance's hardware firewall services.
	FirewallServiceComponent *Network_Component_Firewall `json:"firewallServiceComponent,omitempty" xmlrpc:"firewallServiceComponent,omitempty"`

	// A count of a guest's frontend network components.
	FrontendNetworkComponentCount *uint `json:"frontendNetworkComponentCount,omitempty" xmlrpc:"frontendNetworkComponentCount,omitempty"`

	// A guest's frontend network components.
	FrontendNetworkComponents []Virtual_Guest_Network_Component `json:"frontendNetworkComponents,omitempty" xmlrpc:"frontendNetworkComponents,omitempty"`

	// A guest's frontend or public router.
	FrontendRouters *Hardware `json:"frontendRouters,omitempty" xmlrpc:"frontendRouters,omitempty"`

	// A name reflecting the hostname and domain of the computing instance.
	FullyQualifiedDomainName *string `json:"fullyQualifiedDomainName,omitempty" xmlrpc:"fullyQualifiedDomainName,omitempty"`

	// A guest's universally unique identifier.
	GlobalIdentifier *string `json:"globalIdentifier,omitempty" xmlrpc:"globalIdentifier,omitempty"`

	// The number of GPUs attached to the guest.
	GpuCount *int `json:"gpuCount,omitempty" xmlrpc:"gpuCount,omitempty"`

	// The name of the GPU type attached to the guest.
	GpuType *string `json:"gpuType,omitempty" xmlrpc:"gpuType,omitempty"`

	// no documentation yet
	GuestBootParameter *Virtual_Guest_Boot_Parameter `json:"guestBootParameter,omitempty" xmlrpc:"guestBootParameter,omitempty"`

	// The object's function.
	HardwareFunctionDescription *string `json:"hardwareFunctionDescription,omitempty" xmlrpc:"hardwareFunctionDescription,omitempty"`

	// The virtual host on which a virtual guest resides (available only on private clouds).
	Host *Virtual_Host `json:"host,omitempty" xmlrpc:"host,omitempty"`

	// A host IPS software component object.
	HostIpsSoftwareComponent *Software_Component `json:"hostIpsSoftwareComponent,omitempty" xmlrpc:"hostIpsSoftwareComponent,omitempty"`

	// A virtual computing instance's hostname
	Hostname *string `json:"hostname,omitempty" xmlrpc:"hostname,omitempty"`

	// A guest's hourly billing status.
	HourlyBillingFlag *bool `json:"hourlyBillingFlag,omitempty" xmlrpc:"hourlyBillingFlag,omitempty"`

	// Unique ID for a computing instance.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The total private inbound bandwidth for this computing instance for the current billing cycle.
	InboundPrivateBandwidthUsage *Float64 `json:"inboundPrivateBandwidthUsage,omitempty" xmlrpc:"inboundPrivateBandwidthUsage,omitempty"`

	// The total public inbound bandwidth for this computing instance for the current billing cycle.
	InboundPublicBandwidthUsage *Float64 `json:"inboundPublicBandwidthUsage,omitempty" xmlrpc:"inboundPublicBandwidthUsage,omitempty"`

	// A count of
	InternalTagReferenceCount *uint `json:"internalTagReferenceCount,omitempty" xmlrpc:"internalTagReferenceCount,omitempty"`

	// no documentation yet
	InternalTagReferences []Tag_Reference `json:"internalTagReferences,omitempty" xmlrpc:"internalTagReferences,omitempty"`

	// The last known power state of a virtual guest in the event the guest is turned off outside of IMS or has gone offline.
	LastKnownPowerState *Virtual_Guest_Power_State `json:"lastKnownPowerState,omitempty" xmlrpc:"lastKnownPowerState,omitempty"`

	// The last transaction that a cloud server's operating system was loaded.
	LastOperatingSystemReload *Provisioning_Version1_Transaction `json:"lastOperatingSystemReload,omitempty" xmlrpc:"lastOperatingSystemReload,omitempty"`

	// no documentation yet
	LastPowerStateId *int `json:"lastPowerStateId,omitempty" xmlrpc:"lastPowerStateId,omitempty"`

	// The last transaction a cloud server had performed.
	LastTransaction *Provisioning_Version1_Transaction `json:"lastTransaction,omitempty" xmlrpc:"lastTransaction,omitempty"`

	// The last timestamp of when the guest was verified as a resident virtual machine on the host's hypervisor platform.
	LastVerifiedDate *Time `json:"lastVerifiedDate,omitempty" xmlrpc:"lastVerifiedDate,omitempty"`

	// A virtual guest's latest network monitoring incident.
	LatestNetworkMonitorIncident *Network_Monitor_Version1_Incident `json:"latestNetworkMonitorIncident,omitempty" xmlrpc:"latestNetworkMonitorIncident,omitempty"`

	// A flag indicating that the virtual guest has at least one disk which is local to the host it runs on. This does not include a SWAP device.
	LocalDiskFlag *bool `json:"localDiskFlag,omitempty" xmlrpc:"localDiskFlag,omitempty"`

	// Where guest is located within SoftLayer's location hierarchy.
	Location *Location `json:"location,omitempty" xmlrpc:"location,omitempty"`

	// A flag indicating that the virtual guest is a managed resource.
	ManagedResourceFlag *bool `json:"managedResourceFlag,omitempty" xmlrpc:"managedResourceFlag,omitempty"`

	// The maximum amount of CPU resources a computing instance may utilize.
	MaxCpu *int `json:"maxCpu,omitempty" xmlrpc:"maxCpu,omitempty"`

	// The unit of the maximum amount of CPU resources a computing instance may utilize.
	MaxCpuUnits *string `json:"maxCpuUnits,omitempty" xmlrpc:"maxCpuUnits,omitempty"`

	// The maximum amount of memory a computing instance may utilize.
	MaxMemory *int `json:"maxMemory,omitempty" xmlrpc:"maxMemory,omitempty"`

	// The date of the most recent metric tracking poll performed.
	MetricPollDate *Time `json:"metricPollDate,omitempty" xmlrpc:"metricPollDate,omitempty"`

	// A guest's metric tracking object.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// The metric tracking object id for this guest.
	MetricTrackingObjectId *int `json:"metricTrackingObjectId,omitempty" xmlrpc:"metricTrackingObjectId,omitempty"`

	// The date a virtual computing instance was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// no documentation yet
	MonitoringRobot *Monitoring_Robot `json:"monitoringRobot,omitempty" xmlrpc:"monitoringRobot,omitempty"`

	// A virtual guest's network monitoring services.
	MonitoringServiceComponent *Network_Monitor_Version1_Query_Host_Stratum `json:"monitoringServiceComponent,omitempty" xmlrpc:"monitoringServiceComponent,omitempty"`

	// no documentation yet
	MonitoringServiceEligibilityFlag *bool `json:"monitoringServiceEligibilityFlag,omitempty" xmlrpc:"monitoringServiceEligibilityFlag,omitempty"`

	// The monitoring notification objects for this guest. Each object links this guest instance to a user account that will be notified if monitoring on this guest object fails
	MonitoringUserNotification []User_Customer_Notification_Virtual_Guest `json:"monitoringUserNotification,omitempty" xmlrpc:"monitoringUserNotification,omitempty"`

	// A count of the monitoring notification objects for this guest. Each object links this guest instance to a user account that will be notified if monitoring on this guest object fails
	MonitoringUserNotificationCount *uint `json:"monitoringUserNotificationCount,omitempty" xmlrpc:"monitoringUserNotificationCount,omitempty"`

	// A count of a guests's network components.
	NetworkComponentCount *uint `json:"networkComponentCount,omitempty" xmlrpc:"networkComponentCount,omitempty"`

	// A guests's network components.
	NetworkComponents []Virtual_Guest_Network_Component `json:"networkComponents,omitempty" xmlrpc:"networkComponents,omitempty"`

	// A count of a guests's network monitors.
	NetworkMonitorCount *uint `json:"networkMonitorCount,omitempty" xmlrpc:"networkMonitorCount,omitempty"`

	// A count of all of a virtual guest's network monitoring incidents.
	NetworkMonitorIncidentCount *uint `json:"networkMonitorIncidentCount,omitempty" xmlrpc:"networkMonitorIncidentCount,omitempty"`

	// All of a virtual guest's network monitoring incidents.
	NetworkMonitorIncidents []Network_Monitor_Version1_Incident `json:"networkMonitorIncidents,omitempty" xmlrpc:"networkMonitorIncidents,omitempty"`

	// A guests's network monitors.
	NetworkMonitors []Network_Monitor_Version1_Query_Host `json:"networkMonitors,omitempty" xmlrpc:"networkMonitors,omitempty"`

	// A guest's associated network storage accounts.
	NetworkStorage []Network_Storage `json:"networkStorage,omitempty" xmlrpc:"networkStorage,omitempty"`

	// A count of a guest's associated network storage accounts.
	NetworkStorageCount *uint `json:"networkStorageCount,omitempty" xmlrpc:"networkStorageCount,omitempty"`

	// A count of the network Vlans that a guest's network components are associated with.
	NetworkVlanCount *uint `json:"networkVlanCount,omitempty" xmlrpc:"networkVlanCount,omitempty"`

	// The network Vlans that a guest's network components are associated with.
	NetworkVlans []Network_Vlan `json:"networkVlans,omitempty" xmlrpc:"networkVlans,omitempty"`

	// A note of up to 1,000 characters about a virtual server.
	Notes *string `json:"notes,omitempty" xmlrpc:"notes,omitempty"`

	// An open ticket requesting cancellation of this server, if one exists.
	OpenCancellationTicket *Ticket `json:"openCancellationTicket,omitempty" xmlrpc:"openCancellationTicket,omitempty"`

	// A guest's operating system.
	OperatingSystem *Software_Component_OperatingSystem `json:"operatingSystem,omitempty" xmlrpc:"operatingSystem,omitempty"`

	// A guest's operating system software description.
	OperatingSystemReferenceCode *string `json:"operatingSystemReferenceCode,omitempty" xmlrpc:"operatingSystemReferenceCode,omitempty"`

	// The original package id provided with the order for a Cloud Computing Instance.
	OrderedPackageId *string `json:"orderedPackageId,omitempty" xmlrpc:"orderedPackageId,omitempty"`

	// The total private outbound bandwidth for this computing instance for the current billing cycle.
	OutboundPrivateBandwidthUsage *Float64 `json:"outboundPrivateBandwidthUsage,omitempty" xmlrpc:"outboundPrivateBandwidthUsage,omitempty"`

	// The total public outbound bandwidth for this computing instance for the current billing cycle.
	OutboundPublicBandwidthUsage *Float64 `json:"outboundPublicBandwidthUsage,omitempty" xmlrpc:"outboundPublicBandwidthUsage,omitempty"`

	// Whether the bandwidth usage for this computing instance for the current billing cycle exceeds the allocation.
	OverBandwidthAllocationFlag *int `json:"overBandwidthAllocationFlag,omitempty" xmlrpc:"overBandwidthAllocationFlag,omitempty"`

	// When true this virtual guest must be migrated using SoftLayer_Virtual_Guest::migrate.
	PendingMigrationFlag *bool `json:"pendingMigrationFlag,omitempty" xmlrpc:"pendingMigrationFlag,omitempty"`

	// The placement group that a virtual guest belongs to.
	PlacementGroup *Virtual_PlacementGroup `json:"placementGroup,omitempty" xmlrpc:"placementGroup,omitempty"`

	// The placement group ID that the virtual guest belongs to.
	PlacementGroupId *int `json:"placementGroupId,omitempty" xmlrpc:"placementGroupId,omitempty"`

	// URI of the script to be downloaded and executed after installation is complete. This is deprecated in favor of supplementalCreateObjectOptions' postInstallScriptUri.
	// Deprecated: This function has been marked as deprecated.
	PostInstallScriptUri *string `json:"postInstallScriptUri,omitempty" xmlrpc:"postInstallScriptUri,omitempty"`

	// The current power state of a virtual guest.
	PowerState *Virtual_Guest_Power_State `json:"powerState,omitempty" xmlrpc:"powerState,omitempty"`

	// A guest's primary private IP address.
	PrimaryBackendIpAddress *string `json:"primaryBackendIpAddress,omitempty" xmlrpc:"primaryBackendIpAddress,omitempty"`

	// A guest's primary backend network component.
	PrimaryBackendNetworkComponent *Virtual_Guest_Network_Component `json:"primaryBackendNetworkComponent,omitempty" xmlrpc:"primaryBackendNetworkComponent,omitempty"`

	// The guest's primary public IP address.
	PrimaryIpAddress *string `json:"primaryIpAddress,omitempty" xmlrpc:"primaryIpAddress,omitempty"`

	// A guest's primary public network component.
	PrimaryNetworkComponent *Virtual_Guest_Network_Component `json:"primaryNetworkComponent,omitempty" xmlrpc:"primaryNetworkComponent,omitempty"`

	// Whether the computing instance only has access to the private network.
	PrivateNetworkOnlyFlag *bool `json:"privateNetworkOnlyFlag,omitempty" xmlrpc:"privateNetworkOnlyFlag,omitempty"`

	// Whether the bandwidth usage for this computing instance for the current billing cycle is projected to exceed the allocation.
	ProjectedOverBandwidthAllocationFlag *int `json:"projectedOverBandwidthAllocationFlag,omitempty" xmlrpc:"projectedOverBandwidthAllocationFlag,omitempty"`

	// The projected public outbound bandwidth for this computing instance for the current billing cycle.
	ProjectedPublicBandwidthUsage *Float64 `json:"projectedPublicBandwidthUsage,omitempty" xmlrpc:"projectedPublicBandwidthUsage,omitempty"`

	// no documentation yet
	ProvisionDate *Time `json:"provisionDate,omitempty" xmlrpc:"provisionDate,omitempty"`

	// A count of recent events that impact this computing instance.
	RecentEventCount *uint `json:"recentEventCount,omitempty" xmlrpc:"recentEventCount,omitempty"`

	// Recent events that impact this computing instance.
	RecentEvents []Notification_Occurrence_Event `json:"recentEvents,omitempty" xmlrpc:"recentEvents,omitempty"`

	// no documentation yet
	ReclaimDate *Time `json:"reclaimDate,omitempty" xmlrpc:"reclaimDate,omitempty"`

	// The regional group this guest is in.
	RegionalGroup *Location_Group_Regional `json:"regionalGroup,omitempty" xmlrpc:"regionalGroup,omitempty"`

	// no documentation yet
	RegionalInternetRegistry *Network_Regional_Internet_Registry `json:"regionalInternetRegistry,omitempty" xmlrpc:"regionalInternetRegistry,omitempty"`

	// The reserved capacity group the guest is associated with.
	ReservedCapacityGroup *Virtual_ReservedCapacityGroup `json:"reservedCapacityGroup,omitempty" xmlrpc:"reservedCapacityGroup,omitempty"`

	// Flag to indicate whether or not a guest is part of a reserved capacity group.
	ReservedCapacityGroupFlag *bool `json:"reservedCapacityGroupFlag,omitempty" xmlrpc:"reservedCapacityGroupFlag,omitempty"`

	// The reserved capacity group instance the guest is associated with.
	ReservedCapacityGroupInstance *Virtual_ReservedCapacityGroup_Instance `json:"reservedCapacityGroupInstance,omitempty" xmlrpc:"reservedCapacityGroupInstance,omitempty"`

	// A count of a guest's vulnerability scan requests.
	SecurityScanRequestCount *uint `json:"securityScanRequestCount,omitempty" xmlrpc:"securityScanRequestCount,omitempty"`

	// A guest's vulnerability scan requests.
	SecurityScanRequests []Network_Security_Scanner_Request `json:"securityScanRequests,omitempty" xmlrpc:"securityScanRequests,omitempty"`

	// The server room that a guest is located at. There may be more than one server room for every data center.
	ServerRoom *Location `json:"serverRoom,omitempty" xmlrpc:"serverRoom,omitempty"`

	// A count of a guest's installed software.
	SoftwareComponentCount *uint `json:"softwareComponentCount,omitempty" xmlrpc:"softwareComponentCount,omitempty"`

	// A guest's installed software.
	SoftwareComponents []Software_Component `json:"softwareComponents,omitempty" xmlrpc:"softwareComponents,omitempty"`

	// A count of sSH keys to be installed on the server during provisioning or an OS reload.
	SshKeyCount *uint `json:"sshKeyCount,omitempty" xmlrpc:"sshKeyCount,omitempty"`

	// SSH keys to be installed on the server during provisioning or an OS reload.
	SshKeys []Security_Ssh_Key `json:"sshKeys,omitempty" xmlrpc:"sshKeys,omitempty"`

	// The number of CPUs available to a computing instance upon startup.
	StartCpus *int `json:"startCpus,omitempty" xmlrpc:"startCpus,omitempty"`

	// A computing instance's status.
	Status *Virtual_Guest_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// A computing instances [[SoftLayer_Virtual_Guest_Status|status]] ID
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`

	// Extra options needed for [[SoftLayer_Virtual_Guest/createObject|createObject]] and [[SoftLayer_Virtual_Guest/createObjects|createObjects]].
	SupplementalCreateObjectOptions *Virtual_Guest_SupplementalCreateObjectOptions `json:"supplementalCreateObjectOptions,omitempty" xmlrpc:"supplementalCreateObjectOptions,omitempty"`

	// A count of
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// no documentation yet
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`

	// Whether or not a computing instance is a Transient Instance.
	TransientGuestFlag *bool `json:"transientGuestFlag,omitempty" xmlrpc:"transientGuestFlag,omitempty"`

	// The endpoint used to notify customers their transient guest is terminating.
	TransientWebhookURI *Virtual_Guest_Attribute `json:"transientWebhookURI,omitempty" xmlrpc:"transientWebhookURI,omitempty"`

	// The type of this virtual guest.
	Type *Virtual_Guest_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// Gives the type of guest categorized as PUBLIC, DEDICATED or PRIVATE.
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`

	// A computing instance's associated upgrade request object if any.
	UpgradeRequest *Product_Upgrade_Request `json:"upgradeRequest,omitempty" xmlrpc:"upgradeRequest,omitempty"`

	// A count of a list of users that have access to this computing instance.
	UserCount *uint `json:"userCount,omitempty" xmlrpc:"userCount,omitempty"`

	// A base64 encoded string containing custom user data for a Cloud Computing Instance order.
	UserData []Virtual_Guest_Attribute `json:"userData,omitempty" xmlrpc:"userData,omitempty"`

	// A count of a base64 encoded string containing custom user data for a Cloud Computing Instance order.
	UserDataCount *uint `json:"userDataCount,omitempty" xmlrpc:"userDataCount,omitempty"`

	// A list of users that have access to this computing instance.
	Users []User_Customer `json:"users,omitempty" xmlrpc:"users,omitempty"`

	// Unique ID for a computing instance's record on a virtualization platform.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`

	// The name of the bandwidth allotment that a hardware belongs too.
	VirtualRack *Network_Bandwidth_Version1_Allotment `json:"virtualRack,omitempty" xmlrpc:"virtualRack,omitempty"`

	// The id of the bandwidth allotment that a computing instance belongs too.
	VirtualRackId *int `json:"virtualRackId,omitempty" xmlrpc:"virtualRackId,omitempty"`

	// The name of the bandwidth allotment that a computing instance belongs too.
	VirtualRackName *string `json:"virtualRackName,omitempty" xmlrpc:"virtualRackName,omitempty"`
}

// no documentation yet
type Virtual_Guest_Attribute struct {
	Entity

	// no documentation yet
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// no documentation yet
	Type *Virtual_Guest_Attribute_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// A guest attribute's value.
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// no documentation yet
type Virtual_Guest_Attribute_Type struct {
	Entity

	// no documentation yet
	Keyname *string `json:"keyname,omitempty" xmlrpc:"keyname,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Virtual_Guest_Attribute_UserData struct {
	Virtual_Guest_Attribute
}

// The block device data type presents the structure in which all block devices will be presented. A block device attaches a disk image to a guest. Internally, the structure supports various virtualization platforms with no change to external interaction.
//
// A guest, also known as a virtual server, represents an allocation of resources on a virtual host.
type Virtual_Guest_Block_Device struct {
	Entity

	// A flag indicating if a block device can be booted from.
	BootableFlag *int `json:"bootableFlag,omitempty" xmlrpc:"bootableFlag,omitempty"`

	// The date a block device was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A name used to identify a block device.
	Device *string `json:"device,omitempty" xmlrpc:"device,omitempty"`

	// The disk image that a block device connects to in a computing instance.
	DiskImage *Virtual_Disk_Image `json:"diskImage,omitempty" xmlrpc:"diskImage,omitempty"`

	// A block device [[SoftLayer_Virtual_Disk_Image|disk image]]'s unique ID.
	DiskImageId *int `json:"diskImageId,omitempty" xmlrpc:"diskImageId,omitempty"`

	// The computing instance that this block device is attached to.
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// The [[SoftLayer_Virtual_Guest|computing instance]] that a block device is associated with.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// A flag indicating if a block device can be plugged into a computing instance without having to shut down the instance.
	HotPlugFlag *int `json:"hotPlugFlag,omitempty" xmlrpc:"hotPlugFlag,omitempty"`

	// A computing instance block device's unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The data a block device was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The writing mode that a virtual block device is mounted as, either "RO" for read-only mode or "RW" for read and write mode.
	MountMode *string `json:"mountMode,omitempty" xmlrpc:"mountMode,omitempty"`

	// The type of device that a virtual block device is mounted as, either "Disk" for a directly connected storage disk or "CD" for devices that are mounted as optical drives..
	MountType *string `json:"mountType,omitempty" xmlrpc:"mountType,omitempty"`

	// no documentation yet
	Status *Virtual_Guest_Block_Device_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// The status of the device, either disconnected or connected
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`

	// A block device's unique ID on a virtualization platform.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// no documentation yet
type Virtual_Guest_Block_Device_Status struct {
	Entity

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// The virtual block device template data type presents the structure in which all archived image templates are presented.
//
// A virtual block device template, also known as a image template, represents the image of a virtual guest instance.
type Virtual_Guest_Block_Device_Template struct {
	Entity

	// A name that identifies a block device template.
	Device *string `json:"device,omitempty" xmlrpc:"device,omitempty"`

	// A block device template's disk image.
	DiskImage *Virtual_Disk_Image `json:"diskImage,omitempty" xmlrpc:"diskImage,omitempty"`

	// A block device template's [[SoftLayer_Virtual_Disk_Image|disk image]] ID.
	DiskImageId *int `json:"diskImageId,omitempty" xmlrpc:"diskImageId,omitempty"`

	// The amount of disk space that a block device template is using.  Use this number along with the units property to obtain the correct space used.
	DiskSpace *Float64 `json:"diskSpace,omitempty" xmlrpc:"diskSpace,omitempty"`

	// A block device template's group. Several block device templates can be combined together into a group for archiving purposes.
	Group *Virtual_Guest_Block_Device_Template_Group `json:"group,omitempty" xmlrpc:"group,omitempty"`

	// A block device template's [[SoftLayer_Virtual_Guest_Block_Device_Template_Group|group]] ID.
	GroupId *int `json:"groupId,omitempty" xmlrpc:"groupId,omitempty"`

	// A block device template's unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The units that will be used with the disk space property to identify the amount of disk space used.
	Units *string `json:"units,omitempty" xmlrpc:"units,omitempty"`
}

// The virtual block device template group data type presents the structure in which a group of archived image templates will be presented. The structure consists of a parent template group which contain multiple child template group objects.  Each child template group object represents the image template in a particular location. Unless editing/deleting a specific child template group object, it is best to use the parent object.
//
// A virtual block device template group, also known as an image template group, represents an image of a virtual guest instance.
type Virtual_Guest_Block_Device_Template_Group struct {
	Entity

	// A block device template group's [[SoftLayer_Account|account]].
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A count of
	AccountContactCount *uint `json:"accountContactCount,omitempty" xmlrpc:"accountContactCount,omitempty"`

	// no documentation yet
	AccountContacts []Account_Contact `json:"accountContacts,omitempty" xmlrpc:"accountContacts,omitempty"`

	// A block device template group's [[SoftLayer_Account|account]] ID
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A count of the accounts which may have read-only access to an image template group. Will only be populated for parent template group objects.
	AccountReferenceCount *uint `json:"accountReferenceCount,omitempty" xmlrpc:"accountReferenceCount,omitempty"`

	// The accounts which may have read-only access to an image template group. Will only be populated for parent template group objects.
	AccountReferences []Virtual_Guest_Block_Device_Template_Group_Accounts `json:"accountReferences,omitempty" xmlrpc:"accountReferences,omitempty"`

	// A count of the block devices that are part of an image template group
	BlockDeviceCount *uint `json:"blockDeviceCount,omitempty" xmlrpc:"blockDeviceCount,omitempty"`

	// The block devices that are part of an image template group
	BlockDevices []Virtual_Guest_Block_Device_Template `json:"blockDevices,omitempty" xmlrpc:"blockDevices,omitempty"`

	// The total disk space of all images in a image template group.
	BlockDevicesDiskSpaceTotal *Float64 `json:"blockDevicesDiskSpaceTotal,omitempty" xmlrpc:"blockDevicesDiskSpaceTotal,omitempty"`

	// A flag indicating that customer is providing the software licenses.
	ByolFlag *bool `json:"byolFlag,omitempty" xmlrpc:"byolFlag,omitempty"`

	// The image template groups that are clones of an image template group.
	Children []Virtual_Guest_Block_Device_Template_Group `json:"children,omitempty" xmlrpc:"children,omitempty"`

	// A count of the image template groups that are clones of an image template group.
	ChildrenCount *uint `json:"childrenCount,omitempty" xmlrpc:"childrenCount,omitempty"`

	// The date a block device template group was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The location containing this image template group. Will only be populated for child template group objects.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// A count of a collection of locations containing a copy of this image template group. Will only be populated for parent template group objects.
	DatacenterCount *uint `json:"datacenterCount,omitempty" xmlrpc:"datacenterCount,omitempty"`

	// A collection of locations containing a copy of this image template group. Will only be populated for parent template group objects.
	Datacenters []Location `json:"datacenters,omitempty" xmlrpc:"datacenters,omitempty"`

	// The first clone of the image template group
	FirstChild *Virtual_Guest_Block_Device_Template_Group `json:"firstChild,omitempty" xmlrpc:"firstChild,omitempty"`

	// A flag indicating if this is a flex image.
	FlexImageFlag *bool `json:"flexImageFlag,omitempty" xmlrpc:"flexImageFlag,omitempty"`

	// An image template's universally unique identifier.
	GlobalIdentifier *string `json:"globalIdentifier,omitempty" xmlrpc:"globalIdentifier,omitempty"`

	// A block device template group's unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The virtual disk image type of this template. Value will be populated on parent and child, but only supports object filtering on the parent.
	ImageType *Virtual_Disk_Image_Type `json:"imageType,omitempty" xmlrpc:"imageType,omitempty"`

	// The virtual disk image type keyname (e.g. SYSTEM, DISK_CAPTURE, ISO, etc) of this template. Value will be populated on parent and child, but only supports object filtering on the parent.
	ImageTypeKeyName *string `json:"imageTypeKeyName,omitempty" xmlrpc:"imageTypeKeyName,omitempty"`

	// A user definable and optional name of a block device template group.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A flag indicating if this is a next generation image.
	NextGenFlag *bool `json:"nextGenFlag,omitempty" xmlrpc:"nextGenFlag,omitempty"`

	// A block device template group's user defined note.
	Note *string `json:"note,omitempty" xmlrpc:"note,omitempty"`

	// The image template group that another image template group was cloned from.
	Parent *Virtual_Guest_Block_Device_Template_Group `json:"parent,omitempty" xmlrpc:"parent,omitempty"`

	// A block device template group's [[SoftLayer_Virtual_Guest_Block_Device_Template_Group|parent]] ID.  This will only be set when a template group is created from a previously existing template group
	ParentId *int `json:"parentId,omitempty" xmlrpc:"parentId,omitempty"`

	// no documentation yet
	PublicFlag *int `json:"publicFlag,omitempty" xmlrpc:"publicFlag,omitempty"`

	// no documentation yet
	Region *Network_Service_Resource `json:"region,omitempty" xmlrpc:"region,omitempty"`

	// A count of
	RegionCount *uint `json:"regionCount,omitempty" xmlrpc:"regionCount,omitempty"`

	// no documentation yet
	Regions []Network_Service_Resource `json:"regions,omitempty" xmlrpc:"regions,omitempty"`

	// A count of the ssh keys to be implemented on the server when provisioned or reloaded from an image template group.
	SshKeyCount *uint `json:"sshKeyCount,omitempty" xmlrpc:"sshKeyCount,omitempty"`

	// The ssh keys to be implemented on the server when provisioned or reloaded from an image template group.
	SshKeys []Security_Ssh_Key `json:"sshKeys,omitempty" xmlrpc:"sshKeys,omitempty"`

	// A template group's status.
	Status *Virtual_Guest_Block_Device_Template_Group_Status `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// A block device template group's [[SoftLayer_Virtual_Guest_Block_Device_Template_Group_Status|status]] ID
	StatusId *int `json:"statusId,omitempty" xmlrpc:"statusId,omitempty"`

	// The storage repository that an image template group resides on.
	StorageRepository *Virtual_Storage_Repository `json:"storageRepository,omitempty" xmlrpc:"storageRepository,omitempty"`

	// A block device template group's user defined summary.
	Summary *string `json:"summary,omitempty" xmlrpc:"summary,omitempty"`

	// A count of the tags associated with this image template group.
	TagReferenceCount *uint `json:"tagReferenceCount,omitempty" xmlrpc:"tagReferenceCount,omitempty"`

	// The tags associated with this image template group.
	TagReferences []Tag_Reference `json:"tagReferences,omitempty" xmlrpc:"tagReferences,omitempty"`

	// A transaction that is being performed on a image template group.
	Transaction *Provisioning_Version1_Transaction `json:"transaction,omitempty" xmlrpc:"transaction,omitempty"`

	// A block device template group's [[SoftLayer_Provisioning_Version1_Transaction|transaction]] ID.  This will only be set when there is a transaction being performed on the block device template group.
	TransactionId *int `json:"transactionId,omitempty" xmlrpc:"transactionId,omitempty"`

	// A block device template group's [[SoftLayer_User_Customer]] ID
	UserRecordId *int `json:"userRecordId,omitempty" xmlrpc:"userRecordId,omitempty"`
}

// The SoftLayer_Virtual_Guest_Block_Device_Template_Group_Accounts data type represents the SoftLayer customer accounts which have access to provision CloudLayer Computing Instances from an image template group.
//
// All accounts other than the image template group owner have read-only access to that image template group.
//
// It is important to note that this data type should only exist to give accounts access to the parent template group object, not the child.  All image template sharing between accounts should occur on the parent object.
type Virtual_Guest_Block_Device_Template_Group_Accounts struct {
	Entity

	// The [[SoftLayer_Account|account]] that an image template group is shared with.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The [[SoftLayer_Account|account]] ID which will have access to an image.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The date access was granted to an account.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The [[SoftLayer_Virtual_Guest_Block_Device_Template_Group|image template group]] that is shared with an account.
	Group *Virtual_Guest_Block_Device_Template_Group `json:"group,omitempty" xmlrpc:"group,omitempty"`

	// The [[SoftLayer_Virtual_Guest_Block_Device_Template_Group|group]] ID which access will be granted to.
	GroupId *int `json:"groupId,omitempty" xmlrpc:"groupId,omitempty"`
}

// The virtual block device template group status data type represents the current status of the image template. Depending upon the status, the image template can be used for provisioning or reloading.
//
// For an operating system reload, the image template will need to have a status of 'Active' or 'Deprecated'. For a provision, the image template will need to have a status of 'Active'
type Virtual_Guest_Block_Device_Template_Group_Status struct {
	Entity

	// no documentation yet
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Virtual_Guest_Boot_Parameter struct {
	Entity

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// no documentation yet
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// no documentation yet
	GuestBootParameterType *Virtual_Guest_Boot_Parameter_Type `json:"guestBootParameterType,omitempty" xmlrpc:"guestBootParameterType,omitempty"`

	// no documentation yet
	GuestBootParameterTypeId *int `json:"guestBootParameterTypeId,omitempty" xmlrpc:"guestBootParameterTypeId,omitempty"`

	// no documentation yet
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`
}

// Describes a virtual guest boot parameter. In this the word class is used in the context of arguments sent to cloud computing instances such as single user mode and boot into bash.
type Virtual_Guest_Boot_Parameter_Type struct {
	Entity

	// Available boot options.
	BootOption *string `json:"bootOption,omitempty" xmlrpc:"bootOption,omitempty"`

	// no documentation yet
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A description of the boot parameter
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// Indentifier for record.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The key name of the boot parameter.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The common name of the boot parameter.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The virtual machine arguments
	Value *string `json:"value,omitempty" xmlrpc:"value,omitempty"`
}

// The virtual guest network component data type presents the structure in which all computing instance network components are presented. Internally, the structure supports various virtualization platforms with no change to external interaction.
//
// A guest, also known as a virtual server, represents an allocation of resources on a virtual host.
type Virtual_Guest_Network_Component struct {
	Entity

	// The date a computing instance's network component was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The computing instance that this network component exists on.
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// The unique ID of the [[SoftLayer_Virtual_Guest|computing instance]] that this network component belongs to.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// no documentation yet
	HighAvailabilityFirewallFlag *bool `json:"highAvailabilityFirewallFlag,omitempty" xmlrpc:"highAvailabilityFirewallFlag,omitempty"`

	// no documentation yet
	IcpBinding *Virtual_Guest_Network_Component_IcpBinding `json:"icpBinding,omitempty" xmlrpc:"icpBinding,omitempty"`

	// A computing instance's network component's unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the records of all IP addresses bound to a computing instance's network component.
	IpAddressBindingCount *uint `json:"ipAddressBindingCount,omitempty" xmlrpc:"ipAddressBindingCount,omitempty"`

	// The records of all IP addresses bound to a computing instance's network component.
	IpAddressBindings []Virtual_Guest_Network_Component_IpAddress `json:"ipAddressBindings,omitempty" xmlrpc:"ipAddressBindings,omitempty"`

	// A computing instance network component's unique MAC address.
	MacAddress *string `json:"macAddress,omitempty" xmlrpc:"macAddress,omitempty"`

	// A computing instance network component's maximum allowed speed, measured in Mbit per second. ''maxSpeed'' is determined by the capabilities of the network interface and the port speed purchased on your SoftLayer computing instance.
	MaxSpeed *int `json:"maxSpeed,omitempty" xmlrpc:"maxSpeed,omitempty"`

	// The date a computing instance's network component was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// A computing instance network component's short name. This is usually ''eth''. Use this in conjunction with the ''port'' property to identify a network component. For instance, the "eth0" interface on a server has the network component name "eth" and port 0.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The upstream network component firewall.
	NetworkComponentFirewall *Network_Component_Firewall `json:"networkComponentFirewall,omitempty" xmlrpc:"networkComponentFirewall,omitempty"`

	// A computing instance's network component's [[SoftLayer_Virtual_Network|network]] ID
	NetworkId *int `json:"networkId,omitempty" xmlrpc:"networkId,omitempty"`

	// The VLAN that a computing instance network component's subnet is associated with.
	NetworkVlan *Network_Vlan `json:"networkVlan,omitempty" xmlrpc:"networkVlan,omitempty"`

	// A computing instance network component's port number. Most computing instances have more than one network interface. The port property separates these interfaces. Use this in conjunction with the ''name'' property to identify a network component. For instance, the "eth0" interface on a server has the network component name "eth" and port 0.
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// A computing instance network component's primary IP address.
	PrimaryIpAddress *string `json:"primaryIpAddress,omitempty" xmlrpc:"primaryIpAddress,omitempty"`

	// no documentation yet
	PrimaryIpAddressRecord *Network_Subnet_IpAddress `json:"primaryIpAddressRecord,omitempty" xmlrpc:"primaryIpAddressRecord,omitempty"`

	// A network component's subnet for its primary IP address
	PrimarySubnet *Network_Subnet `json:"primarySubnet,omitempty" xmlrpc:"primarySubnet,omitempty"`

	// A network component's primary IPv6 IP address record.
	PrimaryVersion6IpAddressRecord *Network_Subnet_IpAddress `json:"primaryVersion6IpAddressRecord,omitempty" xmlrpc:"primaryVersion6IpAddressRecord,omitempty"`

	// A network component's routers.
	Router *Hardware_Router `json:"router,omitempty" xmlrpc:"router,omitempty"`

	// A count of the bindings associating security groups to this network component
	SecurityGroupBindingCount *uint `json:"securityGroupBindingCount,omitempty" xmlrpc:"securityGroupBindingCount,omitempty"`

	// The bindings associating security groups to this network component
	SecurityGroupBindings []Virtual_Network_SecurityGroup_NetworkComponentBinding `json:"securityGroupBindings,omitempty" xmlrpc:"securityGroupBindings,omitempty"`

	// A computing instance network component's speed, measured in Mbit per second.
	Speed *int `json:"speed,omitempty" xmlrpc:"speed,omitempty"`

	// A computing instance network component's status. This can be one of four possible values: "ACTIVE", "DISABLED", "INACTIVE", or "ABUSE_DISCONNECT". "ACTIVE" network components are enabled and in use on a cloud instance. "ABUSE_DISCONNECT" status components have been administratively disabled by SoftLayer accounting or abuse. "DISABLED" components have been administratively disabled by you, the user. You should never see a network interface in MACWAIT state. If you happen to see one please contact SoftLayer support.
	Status *string `json:"status,omitempty" xmlrpc:"status,omitempty"`

	// A count of a network component's subnets. A subnet is a group of IP addresses
	SubnetCount *uint `json:"subnetCount,omitempty" xmlrpc:"subnetCount,omitempty"`

	// A network component's subnets. A subnet is a group of IP addresses
	Subnets []Network_Subnet `json:"subnets,omitempty" xmlrpc:"subnets,omitempty"`

	// A computing instance's network component's unique ID on a virtualization platform.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// no documentation yet
type Virtual_Guest_Network_Component_IcpBinding struct {
	Entity

	// no documentation yet
	InterfaceId *string `json:"interfaceId,omitempty" xmlrpc:"interfaceId,omitempty"`

	// no documentation yet
	IpAllocationId *string `json:"ipAllocationId,omitempty" xmlrpc:"ipAllocationId,omitempty"`

	// The network component associated with this object.
	NetworkComponent *Virtual_Guest_Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`

	// The type of IP that this IP allocation id record references.  Set to PRIMARY for the first servic port.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// no documentation yet
	VpcId *string `json:"vpcId,omitempty" xmlrpc:"vpcId,omitempty"`
}

// The SoftLayer_Virtual_Guest_Network_Component_IpAddress data type contains general information relating to the binding of a single network component to a single SoftLayer IP address.
type Virtual_Guest_Network_Component_IpAddress struct {
	Entity

	// The IP address associated with this object's network component.
	IpAddress *Network_Subnet_IpAddress `json:"ipAddress,omitempty" xmlrpc:"ipAddress,omitempty"`

	// The unique ID of the [[SoftLayer_Network_Subnet_ipAddress|ip address]] this virtual IP address is associated with.
	IpAddressId *int `json:"ipAddressId,omitempty" xmlrpc:"ipAddressId,omitempty"`

	// The network component associated with this object's IP address.
	NetworkComponent *Virtual_Guest_Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`

	// The port that a network component has reserved.  This field is only required for some IP address types.
	Port *int `json:"port,omitempty" xmlrpc:"port,omitempty"`

	// The type of IP that this IP address record references.  Some examples are PRIMARY for the network component's primary IP address and CONSOLE_PROXY which represents the IP information for logging into a computing instance's console.
	Type *string `json:"type,omitempty" xmlrpc:"type,omitempty"`
}

// The power state class provides a common set of values for which a guest's power state will be presented in the SoftLayer API.
type Virtual_Guest_Power_State struct {
	Entity

	// The description of a power state
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The key name of a power state
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The name of a power state
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Virtual_Guest_Status struct {
	Entity

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Virtual_Guest_SupplementalCreateObjectOptions struct {
	Entity

	// The mode used to boot the [[SoftLayer_Virtual_Guest]].  Supported values are 'PV' and 'HVM'.
	BootMode *string `json:"bootMode,omitempty" xmlrpc:"bootMode,omitempty"`

	// When set the startCpus and maxMemory are defined by the flavor. If the flavor includes local storage blockDevice 0 is also defined by the flavor. When startCpus, maxMemory, or blockDevice 0 are also provided on the template object they are validated against the flavor provided.
	FlavorKeyName *string `json:"flavorKeyName,omitempty" xmlrpc:"flavorKeyName,omitempty"`

	// When explicitly set to true, createObject(s) will fail unless the order is started automatically. This can be used by automated systems to fail an order that might otherwise require manual approval. For multi-guest orders via [[SoftLayer_Virtual_Guest/createObjects|createObjects]], this value must be the exact same for every item.
	ImmediateApprovalOnlyFlag *bool `json:"immediateApprovalOnlyFlag,omitempty" xmlrpc:"immediateApprovalOnlyFlag,omitempty"`

	// URI of the script to be downloaded and executed after installation is complete. This can be different for each virtual guest when multiple are sent to [[SoftLayer_Virtual_Guest/createObjects|createObjects]].
	PostInstallScriptUri *string `json:"postInstallScriptUri,omitempty" xmlrpc:"postInstallScriptUri,omitempty"`
}

// SoftLayer_Virtual_Guest_Type models the type of a [[SoftLayer_Virtual_Guest]] (PUBLIC | DEDICATED | PRIVATE)
type Virtual_Guest_Type struct {
	Entity

	// no documentation yet
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// no documentation yet
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// no documentation yet
type Virtual_Guest_Vpc_IpAllocation struct {
	Entity
}

// no documentation yet
type Virtual_Guest_Vpc_Subnet struct {
	Entity
}

// The virtual host represents the platform on which virtual guests reside. At times a virtual host has no allocations on the physical server, however with many modern platforms it is a virtual machine with small CPU and Memory allocations that runs in the Control Domain.
type Virtual_Host struct {
	Entity

	// The account which a virtual host belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// A virtual host's associated account id
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The date a virtual host was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A virtual host's description.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// The enabled flag specifies whether a virtual host can run guests.
	EnabledFlag *int `json:"enabledFlag,omitempty" xmlrpc:"enabledFlag,omitempty"`

	// The hardware record which a virtual host resides on.
	Hardware *Hardware_Server `json:"hardware,omitempty" xmlrpc:"hardware,omitempty"`

	// A hardware device which a virtual host resides.
	HardwareId *int `json:"hardwareId,omitempty" xmlrpc:"hardwareId,omitempty"`

	// Unique ID for a virtual host.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The metric tracking object for this virtual host.
	MetricTrackingObject *Metric_Tracking_Object `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// The date a virtual host was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// A virtual host's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of
	PciDeviceCount *uint `json:"pciDeviceCount,omitempty" xmlrpc:"pciDeviceCount,omitempty"`

	// no documentation yet
	PciDevices []Virtual_Host_PciDevice `json:"pciDevices,omitempty" xmlrpc:"pciDevices,omitempty"`

	// The amount of memory physically available for a virtual host.
	PhysicalMemoryCapacity *int `json:"physicalMemoryCapacity,omitempty" xmlrpc:"physicalMemoryCapacity,omitempty"`

	// Unique ID for a virtual host's record on a virtualization platform.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`
}

// This type represents a PCI device on a host.
type Virtual_Host_PciDevice struct {
	Entity

	// no documentation yet
	HardwareComponentModel *Hardware_Component_Model `json:"hardwareComponentModel,omitempty" xmlrpc:"hardwareComponentModel,omitempty"`

	// no documentation yet
	Host *Virtual_Host `json:"host,omitempty" xmlrpc:"host,omitempty"`

	// ID of the PCI device.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The unique id of the PCI device's record on a virtualization platform.
	Uuid *string `json:"uuid,omitempty" xmlrpc:"uuid,omitempty"`

	// The BDF (Domain:Bus:Device.Function) id of the PCI device in XenServer.
	XenPciId *string `json:"xenPciId,omitempty" xmlrpc:"xenPciId,omitempty"`
}

// The SoftLayer_Virtual_Network_SecurityGroup_NetworkComponentBinding data type contains general information for a single binding. A binding associates a [[SoftLayer_Virtual_Guest_Network_Component]] with a [[SoftLayer_Network_SecurityGroup]].
type Virtual_Network_SecurityGroup_NetworkComponentBinding struct {
	Entity

	// The unique ID for a binding.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	NetworkComponent *Virtual_Guest_Network_Component `json:"networkComponent,omitempty" xmlrpc:"networkComponent,omitempty"`

	// The ID of the network component.
	NetworkComponentId *int `json:"networkComponentId,omitempty" xmlrpc:"networkComponentId,omitempty"`

	// no documentation yet
	SecurityGroup *Network_SecurityGroup `json:"securityGroup,omitempty" xmlrpc:"securityGroup,omitempty"`

	// The ID of the security group.
	SecurityGroupId *int `json:"securityGroupId,omitempty" xmlrpc:"securityGroupId,omitempty"`
}

// This data type presents the structure for a virtual guest placement group. The data type contains relational properties to the virtual guest placement group rule class.
type Virtual_PlacementGroup struct {
	Entity

	// The account that the placement group is implemented on.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The unique ID of the account that created the placement group.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// The router the placement group is implemented on.
	BackendRouter *Hardware_Router_Backend `json:"backendRouter,omitempty" xmlrpc:"backendRouter,omitempty"`

	// The placement group's backend router's associated unique ID.
	BackendRouterId *int `json:"backendRouterId,omitempty" xmlrpc:"backendRouterId,omitempty"`

	// The placement group's date of creation.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// A count of the virtual guests that are members of the placement group.
	GuestCount *uint `json:"guestCount,omitempty" xmlrpc:"guestCount,omitempty"`

	// The virtual guests that are members of the placement group.
	Guests []Virtual_Guest `json:"guests,omitempty" xmlrpc:"guests,omitempty"`

	// The placement group's associated unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The placement group's date of most recent modification.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The placement group's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The placement rule that the placement group is implementing.
	Rule *Virtual_PlacementGroup_Rule `json:"rule,omitempty" xmlrpc:"rule,omitempty"`

	// The associated unique ID of the placement group's rule.
	RuleId *int `json:"ruleId,omitempty" xmlrpc:"ruleId,omitempty"`
}

// This data type presents the structure of a virtual guest placement group rule.
type Virtual_PlacementGroup_Rule struct {
	Entity

	// The placement group rule's associated unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The placement group rule's key name.
	KeyName *string `json:"keyName,omitempty" xmlrpc:"keyName,omitempty"`

	// The placement group rule's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`
}

// This data type presents the structure for a virtual reserved capacity group.
type Virtual_ReservedCapacityGroup struct {
	Entity

	// The account that the reserved capacity group is implemented on.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The unique ID of the account that created the reserved capacity group.
	AccountId *int `json:"accountId,omitempty" xmlrpc:"accountId,omitempty"`

	// A count of the instances available for guest provisions on this reserved capacity group.
	AvailableInstanceCount *uint `json:"availableInstanceCount,omitempty" xmlrpc:"availableInstanceCount,omitempty"`

	// The instances available for guest provisions on this reserved capacity group.
	AvailableInstances []Virtual_ReservedCapacityGroup_Instance `json:"availableInstances,omitempty" xmlrpc:"availableInstances,omitempty"`

	// The router the reserved capacity group is implemented on.
	BackendRouter *Hardware_Router_Backend `json:"backendRouter,omitempty" xmlrpc:"backendRouter,omitempty"`

	// The reserved capacity group's backend router's associated unique ID.
	BackendRouterId *int `json:"backendRouterId,omitempty" xmlrpc:"backendRouterId,omitempty"`

	// The date that the reserved capacity group was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The reserved capacity group's associated unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// A count of the guest instances that are members of this reserved capacity group.
	InstanceCount *uint `json:"instanceCount,omitempty" xmlrpc:"instanceCount,omitempty"`

	// The guest instances that are members of this reserved capacity group.
	Instances []Virtual_ReservedCapacityGroup_Instance `json:"instances,omitempty" xmlrpc:"instances,omitempty"`

	// The number of instances that are members of this reserved capacity group.
	InstancesCount *uint `json:"instancesCount,omitempty" xmlrpc:"instancesCount,omitempty"`

	// The date that the reserved capacity group was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The reserved capacity group's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// A count of the instances already occupied by a guest on this reserved capacity group.
	OccupiedInstanceCount *uint `json:"occupiedInstanceCount,omitempty" xmlrpc:"occupiedInstanceCount,omitempty"`

	// The instances already occupied by a guest on this reserved capacity group.
	OccupiedInstances []Virtual_ReservedCapacityGroup_Instance `json:"occupiedInstances,omitempty" xmlrpc:"occupiedInstances,omitempty"`
}

// This data type presents the structure for a virtual reserved capacity group instance.
type Virtual_ReservedCapacityGroup_Instance struct {
	Entity

	// Flag to indecate whether or not the reserved instance is available or not.
	AvailableFlag *bool `json:"availableFlag,omitempty" xmlrpc:"availableFlag,omitempty"`

	// The billing item for the reserved capacity group instance.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// The date that the reserved capacity group instance was created.
	CreateDate *Time `json:"createDate,omitempty" xmlrpc:"createDate,omitempty"`

	// The virtual guest associated with this reserved capacity group instance.
	Guest *Virtual_Guest `json:"guest,omitempty" xmlrpc:"guest,omitempty"`

	// The virtual guest ID associated with this reserved capacity group instance.
	GuestId *int `json:"guestId,omitempty" xmlrpc:"guestId,omitempty"`

	// The reserved capacity group instance's associated unique ID.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// The date that the reserved capacity group instance was last modified.
	ModifyDate *Time `json:"modifyDate,omitempty" xmlrpc:"modifyDate,omitempty"`

	// The reserved instances that are members of this reserved capacity group.
	ReservedCapacityGroup *Virtual_ReservedCapacityGroup `json:"reservedCapacityGroup,omitempty" xmlrpc:"reservedCapacityGroup,omitempty"`

	// The ID of the reserved capacity group this instance is associated with.
	ReservedCapacityGroupId *int `json:"reservedCapacityGroupId,omitempty" xmlrpc:"reservedCapacityGroupId,omitempty"`
}

// The SoftLayer_Virtual_Storage_Repository represents a web based storage system that can be accessed through many types of devices, interfaces, and other resources.
type Virtual_Storage_Repository struct {
	Entity

	// The [[SoftLayer_Account|account]] that a storage repository belongs to.
	Account *Account `json:"account,omitempty" xmlrpc:"account,omitempty"`

	// The current billing item for a storage repository.
	BillingItem *Billing_Item `json:"billingItem,omitempty" xmlrpc:"billingItem,omitempty"`

	// A storage repositories capacity measured in Giga-Bytes (GB)
	Capacity *Float64 `json:"capacity,omitempty" xmlrpc:"capacity,omitempty"`

	// The datacenter that a virtual storage repository resides in.
	Datacenter *Location `json:"datacenter,omitempty" xmlrpc:"datacenter,omitempty"`

	// A storage repositories description that describes its purpose or contents
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A count of the [[SoftLayer_Virtual_Disk_Image|disk images]] that are in a storage repository. Disk images are the virtual hard drives for a virtual guest.
	DiskImageCount *uint `json:"diskImageCount,omitempty" xmlrpc:"diskImageCount,omitempty"`

	// The [[SoftLayer_Virtual_Disk_Image|disk images]] that are in a storage repository. Disk images are the virtual hard drives for a virtual guest.
	DiskImages []Virtual_Disk_Image `json:"diskImages,omitempty" xmlrpc:"diskImages,omitempty"`

	// A count of the computing instances that have disk images in a storage repository.
	GuestCount *uint `json:"guestCount,omitempty" xmlrpc:"guestCount,omitempty"`

	// The computing instances that have disk images in a storage repository.
	Guests []Virtual_Guest `json:"guests,omitempty" xmlrpc:"guests,omitempty"`

	// Unique ID for a storage repository.
	Id *int `json:"id,omitempty" xmlrpc:"id,omitempty"`

	// no documentation yet
	MetricTrackingObject *Metric_Tracking_Object_Virtual_Storage_Repository `json:"metricTrackingObject,omitempty" xmlrpc:"metricTrackingObject,omitempty"`

	// A storage repositories name that describes its purpose or contents
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// no documentation yet
	PublicFlag *int `json:"publicFlag,omitempty" xmlrpc:"publicFlag,omitempty"`

	// The current billing item for a public storage repository.
	PublicImageBillingItem *Billing_Item `json:"publicImageBillingItem,omitempty" xmlrpc:"publicImageBillingItem,omitempty"`

	// A storage repository's [[SoftLayer_Virtual_Storage_Repository_Type|type]].
	Type *Virtual_Storage_Repository_Type `json:"type,omitempty" xmlrpc:"type,omitempty"`

	// A storage repositories [[SoftLayer_Virtual_Storage_Repository_Type|type]] ID
	TypeId *int `json:"typeId,omitempty" xmlrpc:"typeId,omitempty"`
}

// SoftLayer employs many different types of repositories that computing instances use as their storage volume. SoftLayer_Virtual_Storage_Repository_Type models a single storage type. Common types of storage repositories include networked file systems, logical volume management, and local disk volumes for swap and page file management.
type Virtual_Storage_Repository_Type struct {
	Entity

	// A brief description os a storage repository type.
	Description *string `json:"description,omitempty" xmlrpc:"description,omitempty"`

	// A storage repository type's name.
	Name *string `json:"name,omitempty" xmlrpc:"name,omitempty"`

	// The storage repositories on a SoftLayer customer account that belong to this type.
	StorageRepositories []Virtual_Storage_Repository `json:"storageRepositories,omitempty" xmlrpc:"storageRepositories,omitempty"`

	// A count of the storage repositories on a SoftLayer customer account that belong to this type.
	StorageRepositoryCount *uint `json:"storageRepositoryCount,omitempty" xmlrpc:"storageRepositoryCount,omitempty"`
}
