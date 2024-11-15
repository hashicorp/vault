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

package services

import (
	"fmt"
	"strings"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
)

// This data type presents the structure for a dedicated host. The data type contains relational properties to distinguish a dedicated host and associate an account to it.
type Virtual_DedicatedHost struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualDedicatedHostService returns an instance of the Virtual_DedicatedHost SoftLayer service
func GetVirtualDedicatedHostService(sess session.SLSession) Virtual_DedicatedHost {
	return Virtual_DedicatedHost{Session: sess}
}

func (r Virtual_DedicatedHost) Id(id int) Virtual_DedicatedHost {
	r.Options.Id = &id
	return r
}

func (r Virtual_DedicatedHost) Mask(mask string) Virtual_DedicatedHost {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_DedicatedHost) Filter(filter string) Virtual_DedicatedHost {
	r.Options.Filter = filter
	return r
}

func (r Virtual_DedicatedHost) Limit(limit int) Virtual_DedicatedHost {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_DedicatedHost) Offset(offset int) Virtual_DedicatedHost {
	r.Options.Offset = &offset
	return r
}

// This method will cancel a dedicated host immediately.
func (r Virtual_DedicatedHost) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "deleteObject", nil, &r.Options, &resp)
	return
}

// Edit a dedicated host's properties.
func (r Virtual_DedicatedHost) EditObject(templateObject *datatypes.Virtual_DedicatedHost) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve The account that the dedicated host belongs to.
func (r Virtual_DedicatedHost) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve The container that represents allocations on the dedicated host.
func (r Virtual_DedicatedHost) GetAllocationStatus() (resp datatypes.Container_Virtual_DedicatedHost_AllocationStatus, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getAllocationStatus", nil, &r.Options, &resp)
	return
}

// This method will get the available backend routers to order a dedicated host.
func (r Virtual_DedicatedHost) GetAvailableRouters(dedicatedHost *datatypes.Virtual_DedicatedHost) (resp []datatypes.Hardware, err error) {
	params := []interface{}{
		dedicatedHost,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getAvailableRouters", params, &r.Options, &resp)
	return
}

// Retrieve The backend router behind dedicated host's pool of resources.
func (r Virtual_DedicatedHost) GetBackendRouter() (resp datatypes.Hardware_Router_Backend, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getBackendRouter", nil, &r.Options, &resp)
	return
}

// Retrieve The billing item for the dedicated host.
func (r Virtual_DedicatedHost) GetBillingItem() (resp datatypes.Billing_Item_Virtual_DedicatedHost, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve The datacenter that the dedicated host resides in.
func (r Virtual_DedicatedHost) GetDatacenter() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getDatacenter", nil, &r.Options, &resp)
	return
}

// Retrieve The guests associated with the dedicated host.
func (r Virtual_DedicatedHost) GetGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getGuests", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_DedicatedHost) GetInternalTagReferences() (resp []datatypes.Tag_Reference, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getInternalTagReferences", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_DedicatedHost) GetObject() (resp datatypes.Virtual_DedicatedHost, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The container that represents PCI device allocations on the dedicated host.
func (r Virtual_DedicatedHost) GetPciDeviceAllocationStatus() (resp datatypes.Container_Virtual_DedicatedHost_Pci_Device_AllocationStatus, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getPciDeviceAllocationStatus", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of SoftLayer_Virtual_Host_PciDevice objects on the host.
func (r Virtual_DedicatedHost) GetPciDevices() (resp []datatypes.Virtual_Host_PciDevice, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getPciDevices", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_DedicatedHost) GetTagReferences() (resp []datatypes.Tag_Reference, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "getTagReferences", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_DedicatedHost) SetTags(tags *string) (resp bool, err error) {
	params := []interface{}{
		tags,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_DedicatedHost", "setTags", params, &r.Options, &resp)
	return
}

// The virtual disk image data type presents the structure in which a virtual disk image will be presented.
//
// Virtual block devices are assigned to disk images.
type Virtual_Disk_Image struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualDiskImageService returns an instance of the Virtual_Disk_Image SoftLayer service
func GetVirtualDiskImageService(sess session.SLSession) Virtual_Disk_Image {
	return Virtual_Disk_Image{Session: sess}
}

func (r Virtual_Disk_Image) Id(id int) Virtual_Disk_Image {
	r.Options.Id = &id
	return r
}

func (r Virtual_Disk_Image) Mask(mask string) Virtual_Disk_Image {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Disk_Image) Filter(filter string) Virtual_Disk_Image {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Disk_Image) Limit(limit int) Virtual_Disk_Image {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Disk_Image) Offset(offset int) Virtual_Disk_Image {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Virtual_Disk_Image) EditObject(templateObject *datatypes.Virtual_Disk_Image) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "editObject", params, &r.Options, &resp)
	return
}

// Returns a collection of boot modes that are supported for primary disks.
func (r Virtual_Disk_Image) GetAvailableBootModes() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getAvailableBootModes", nil, &r.Options, &resp)
	return
}

// Retrieve The billing item for a virtual disk image.
func (r Virtual_Disk_Image) GetBillingItem() (resp datatypes.Billing_Item_Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve The block devices that a disk image is attached to. Block devices connect computing instances to disk images.
func (r Virtual_Disk_Image) GetBlockDevices() (resp []datatypes.Virtual_Guest_Block_Device, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getBlockDevices", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Disk_Image) GetBootableVolumeFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getBootableVolumeFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Check if cloud-init is enabled.
func (r Virtual_Disk_Image) GetCloudInitFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getCloudInitFlag", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Disk_Image) GetCoalescedDiskImages() (resp []datatypes.Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getCoalescedDiskImages", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Disk_Image) GetCopyOnWriteFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getCopyOnWriteFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Return disk file extension
func (r Virtual_Disk_Image) GetDiskFileExtension() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getDiskFileExtension", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Disk_Image) GetDiskImageStorageGroup() (resp datatypes.Configuration_Storage_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getDiskImageStorageGroup", nil, &r.Options, &resp)
	return
}

// Retrieve Return imported disk type
func (r Virtual_Disk_Image) GetImportedDiskType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getImportedDiskType", nil, &r.Options, &resp)
	return
}

// Retrieve Return if image is encrypted
func (r Virtual_Disk_Image) GetIsEncrypted() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getIsEncrypted", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Disk_Image) GetLocalDiskFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getLocalDiskFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether this disk image is meant for storage of custom user data supplied with a Cloud Computing Instance order.
func (r Virtual_Disk_Image) GetMetadataFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getMetadataFlag", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Disk_Image) GetObject() (resp datatypes.Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieves images from the public ISO repository
func (r Virtual_Disk_Image) GetPublicIsoImages() (resp []datatypes.Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getPublicIsoImages", nil, &r.Options, &resp)
	return
}

// Retrieve References to the software that resides on a disk image.
func (r Virtual_Disk_Image) GetSoftwareReferences() (resp []datatypes.Virtual_Disk_Image_Software, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getSoftwareReferences", nil, &r.Options, &resp)
	return
}

// Retrieve The original disk image that the current disk image was cloned from.
func (r Virtual_Disk_Image) GetSourceDiskImage() (resp datatypes.Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getSourceDiskImage", nil, &r.Options, &resp)
	return
}

// Retrieve Return storage group details for symantec disk
func (r Virtual_Disk_Image) GetStorageGroupDetails() (resp datatypes.Container_Image_StorageGroupDetails, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getStorageGroupDetails", nil, &r.Options, &resp)
	return
}

// Retrieve The storage group for a virtual disk image.
func (r Virtual_Disk_Image) GetStorageGroups() (resp []datatypes.Configuration_Storage_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getStorageGroups", nil, &r.Options, &resp)
	return
}

// Retrieve The storage repository that a disk image resides in.
func (r Virtual_Disk_Image) GetStorageRepository() (resp datatypes.Virtual_Storage_Repository, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getStorageRepository", nil, &r.Options, &resp)
	return
}

// Retrieve The type of storage repository that a disk image resides in.
func (r Virtual_Disk_Image) GetStorageRepositoryType() (resp datatypes.Virtual_Storage_Repository_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getStorageRepositoryType", nil, &r.Options, &resp)
	return
}

// Retrieve Return supported hardware component IDs for symantec disk
func (r Virtual_Disk_Image) GetSupportedHardware() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getSupportedHardware", nil, &r.Options, &resp)
	return
}

// Retrieve The template that attaches a disk image to a [[SoftLayer_Virtual_Guest_Block_Device_Template_Group|archive]].
func (r Virtual_Disk_Image) GetTemplateBlockDevice() (resp datatypes.Virtual_Guest_Block_Device_Template, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getTemplateBlockDevice", nil, &r.Options, &resp)
	return
}

// Retrieve A virtual disk image's type.
func (r Virtual_Disk_Image) GetType() (resp datatypes.Virtual_Disk_Image_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Disk_Image", "getType", nil, &r.Options, &resp)
	return
}

// The virtual guest data type presents the structure in which all virtual guests will be presented. Internally, the structure supports various virtualization platforms with no change to external interaction.
//
// A guest, also known as a virtual server, represents an allocation of resources on a virtual host.
type Virtual_Guest struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualGuestService returns an instance of the Virtual_Guest SoftLayer service
func GetVirtualGuestService(sess session.SLSession) Virtual_Guest {
	return Virtual_Guest{Session: sess}
}

func (r Virtual_Guest) Id(id int) Virtual_Guest {
	r.Options.Id = &id
	return r
}

func (r Virtual_Guest) Mask(mask string) Virtual_Guest {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Guest) Filter(filter string) Virtual_Guest {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Guest) Limit(limit int) Virtual_Guest {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Guest) Offset(offset int) Virtual_Guest {
	r.Options.Offset = &offset
	return r
}

// Activate the private network port
func (r Virtual_Guest) ActivatePrivatePort() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "activatePrivatePort", nil, &r.Options, &resp)
	return
}

// Activate the public network port
func (r Virtual_Guest) ActivatePublicPort() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "activatePublicPort", nil, &r.Options, &resp)
	return
}

// This method is used to allow access to a SoftLayer_Network_Storage volume that supports host- or network-level access control.
func (r Virtual_Guest) AllowAccessToNetworkStorage(networkStorageTemplateObject *datatypes.Network_Storage) (resp bool, err error) {
	params := []interface{}{
		networkStorageTemplateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "allowAccessToNetworkStorage", params, &r.Options, &resp)
	return
}

// This method is used to allow access to multiple SoftLayer_Network_Storage volumes that support host- or network-level access control.
func (r Virtual_Guest) AllowAccessToNetworkStorageList(networkStorageTemplateObjects []datatypes.Network_Storage) (resp bool, err error) {
	params := []interface{}{
		networkStorageTemplateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "allowAccessToNetworkStorageList", params, &r.Options, &resp)
	return
}

// Creates a transaction to attach a guest's disk image. If the disk image is already attached it will be ignored.
//
// WARNING: SoftLayer_Virtual_Guest::checkHostDiskAvailability should be called before this method. If the SoftLayer_Virtual_Guest::checkHostDiskAvailability method is not called before this method, the guest migration will happen automatically.
func (r Virtual_Guest) AttachDiskImage(imageId *int) (resp datatypes.Provisioning_Version1_Transaction, err error) {
	params := []interface{}{
		imageId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "attachDiskImage", params, &r.Options, &resp)
	return
}

// Reopens the public and/or private ports to reverse the changes made when the server was isolated for a destructive action.
func (r Virtual_Guest) CancelIsolationForDestructiveAction() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "cancelIsolationForDestructiveAction", nil, &r.Options, &resp)
	return
}

// Captures a Flex Image of the hard disk on the virtual machine, based on the capture template parameter. Returns the image template group containing the disk image.
func (r Virtual_Guest) CaptureImage(captureTemplate *datatypes.Container_Disk_Image_Capture_Template) (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	params := []interface{}{
		captureTemplate,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "captureImage", params, &r.Options, &resp)
	return
}

// Checks the associated host for available disk space to determine if guest migration is necessary. This method is only used with local disks. If this method returns false, calling attachDiskImage($imageId) will automatically migrate the destination guest to a new host before attaching the portable volume.
func (r Virtual_Guest) CheckHostDiskAvailability(diskCapacity *int) (resp bool, err error) {
	params := []interface{}{
		diskCapacity,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "checkHostDiskAvailability", params, &r.Options, &resp)
	return
}

// Creates a transaction to configure the guest's metadata disk. If the guest has user data associated with it, the transaction will create a small virtual drive and write the metadata to a file on the drive; if the drive already exists, the metadata will be rewritten. If the guest has no user data associated with it, the transaction will remove the virtual drive if it exists.
//
// WARNING: The transaction created by this service will shut down the guest while the metadata disk is configured. The guest will be turned back on once this process is complete.
func (r Virtual_Guest) ConfigureMetadataDisk() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "configureMetadataDisk", nil, &r.Options, &resp)
	return
}

// Create a transaction to archive a computing instance's block devices
func (r Virtual_Guest) CreateArchiveTemplate(groupName *string, blockDevices []datatypes.Virtual_Guest_Block_Device, note *string) (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	params := []interface{}{
		groupName,
		blockDevices,
		note,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "createArchiveTemplate", params, &r.Options, &resp)
	return
}

// Create a transaction to archive a computing instance's block devices
// Deprecated: This function has been marked as deprecated.
func (r Virtual_Guest) CreateArchiveTransaction(groupName *string, blockDevices []datatypes.Virtual_Guest_Block_Device, note *string) (resp datatypes.Provisioning_Version1_Transaction, err error) {
	params := []interface{}{
		groupName,
		blockDevices,
		note,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "createArchiveTransaction", params, &r.Options, &resp)
	return
}

// createObject() enables the creation of computing instances on an account. This method is a simplified alternative to interacting with the ordering system directly.
//
// In order to create a computing instance, a template object must be sent in with a few required values.
//
// When this method returns an order will have been placed for a computing instance of the specified configuration.
//
// To determine when the instance is available you can poll the instance via [[SoftLayer_Virtual_Guest/getObject]], with an object mask requesting the `provisionDate` relational property. When `provisionDate` is not `null`, the instance will be ready.
//
// > **Warning:** Computing instances created via this method will incur charges on your account. For testing input parameters see [[SoftLayer_Virtual_Guest/generateOrderTemplate]].
//
// ### Required Input [[SoftLayer_Virtual_Guest]]
//
// - `Hostname`  String **Required**
//   - Hostname for the computing instance.
//
// - `Domain` String **Required**
//   - Domain for the computing instance.
//
// - `startCpus` Integer **Required**
//   - The number of CPU cores to allocate.
//   - See [[SoftLayer_Virtual_Guest/getCreateObjectOptions]] for available options.
//
// - `maxMemory` Integer **Required**
//   - The amount of memory to allocate in megabytes.
//   - See [[SoftLayer_Virtual_Guest/getCreateObjectOptions]] for available options.
//
// - `datacenter.name` *String* **Required**
//   - Specifies which datacenter the instance is to be provisioned in. Needs to be a nested object.
//   - Example: `"datacenter": {"name": "dal05"}`
//   - See [[SoftLayer_Virtual_Guest/getCreateObjectOptions]] for available options.
//
// - `hourlyBillingFlag` Boolean **Required**
//   - Specifies the billing type for the instance.
//   - True for hourly billing, False for monthly billing.
//
// - `localDiskFlag` Boolean **Required**
//   - Specifies the disk type for the instance.
//   - True for local to the instance disks, False for SAN disks.
//
// - `dedicatedAccountHostOnlyFlag` Boolean
//   - When true this flag specifies that a compute instance is to run on hosts that only have guests from the same account.
//   - Default: False
//
// - `operatingSystemReferenceCode` String **Conditionally required**
//   - An identifier for the operating system to provision the computing instance with.
//   - Not required when using a `blockDeviceTemplateGroup.globalIdentifier`, as the template will have its own operating system.
//   - See [[SoftLayer_Virtual_Guest/getCreateObjectOptions]] for available options.
//   - **Notice**: Some operating systems are billed based on the number of CPUs the guest has. The price which is used can be determined by calling
//     [[SoftLayer_Virtual_Guest/generateOrderTemplate]] with your desired device specifications.
//
// - `blockDeviceTemplateGroup.globalIdentifier` String
//   - The GUID for the template to be used to provision the computing instance.
//   - Conflicts with `operatingSystemReferenceCode`
//   - **Notice**: Some operating systems are billed based on the number of CPUs the guest has. The price which is used can be determined by calling
//     [[SoftLayer_Virtual_Guest/generateOrderTemplate]] with your desired device specifications.
//   - A list of public images may be obtained via a request to [[SoftLayer_Virtual_Guest_Block_Device_Template_Group/getPublicImages]]
//   - A list of private images may be obtained via a request to [[SoftLayer_Account/getPrivateBlockDeviceTemplateGroups]]
//   - Example: `"blockDeviceTemplateGroup": { globalIdentifier": "07beadaa-1e11-476e-a188-3f7795feb9fb"`
//
// - `networkComponents.maxSpeed` Integer
//   - Specifies the connection speed for the instance's network components.
//   - The `networkComponents` property is an array with a single [[SoftLayer_Virtual_Guest_Network_Component]] structure.
//     The `maxSpeed` property must be set to specify the network uplink speed, in megabits per second, of the computing instance.
//   - See [[SoftLayer_Virtual_Guest/getCreateObjectOptions]] for available options.
//   - Default: 10
//   - Example: `"networkComponents": [{"maxSpeed": 1000}]`
//
// - `privateNetworkOnlyFlag` Boolean
//   - When true this flag specifies that a compute instance is to only have access to the private network.
//   - Default: False
//
// - `primaryNetworkComponent.networkVlan.id` Integer
//   - Specifies the network vlan which is to be used for the frontend interface of the computing instance.
//   - The `primaryNetworkComponent` property is a [[SoftLayer_Virtual_Guest_Network_Component]] structure with the `networkVlan` property populated with a i
//     [[SoftLayer_Network_Vlan]] structure. The `id` property must be set to specify the frontend network vlan of the computing instance.
//   - *NOTE* This is the VLAN `id`, NOT the vlan number.
//   - Example: `"primaryNetworkComponent":{"networkVlan": {"id": 1234567}}`
//
// - `backendNetworkComponent.networkVlan.id` Integer
//   - Specifies the network vlan which is to be used for the backend interface of the computing instance.
//   - The `backendNetworkComponent` property is a [[SoftLayer_Virtual_Guest_Network_Component]] structure with the `networkVlan` property populated with a
//     [[SoftLayer_Network_Vlan]] structure. The `id` property must be set to specify the backend network vlan of the computing instance.
//   - *NOTE* This is the VLAN `id`, NOT the vlan number.
//   - Example: `"backendNetworkComponent":{"networkVlan": {"id": 1234567}}`
//
// - `primaryNetworkComponent.securityGroupBindings` [[SoftLayer_Virtual_Network_SecurityGroup_NetworkComponentBinding]][]
//   - Specifies the security groups to be attached to this VSI's frontend network adapter
//   - The `primaryNetworkComponent` property is a [[SoftLayer_Virtual_Guest_Network_Component]] structure with the `securityGroupBindings` property populated
//     with an array of [[SoftLayer_Virtual_Network_SecurityGroup_NetworkComponentBinding]] structures. The `securityGroup` property in each must be set to
//     specify the security group to be attached to the primary frontend network component.
//   - Example:
//     ```
//     "primaryNetworkComponent": {
//     "securityGroupBindings": [
//     {"securityGroup":{"id": 5555555}},
//     {"securityGroup":{"id": 1112223}},
//     ]
//     }
//     ```
//
// - `primaryBackendNetworkComponent.securityGroupBindings` [[SoftLayer_Virtual_Network_SecurityGroup_NetworkComponentBinding]][]
//   - Specifies the security groups to be attached to this VSI's backend network adapter
//   - The `primaryNetworkComponent` property is a [[SoftLayer_Virtual_Guest_Network_Component]] structure with the `securityGroupBindings` property populated
//     with an array of [[SoftLayer_Virtual_Network_SecurityGroup_NetworkComponentBinding]] structures. The `securityGroup` property in each must be set to
//     specify the security group to be attached to the primary frontend network component.
//   - Example:
//     ```
//     "primaryBackendNetworkComponent": {
//     "securityGroupBindings": [
//     {"securityGroup":{"id": 33322211}},
//     {"securityGroup":{"id": 77777222}},
//     ]
//     }
//     ```
//
// - `blockDevices` [[SoftLayer_Virtual_Guest_Block_Device]][]
//   - Block device and disk image settings for the computing instance
//   - The `blockDevices` property is an array of [[SoftLayer_Virtual_Guest_Block_Device]] structures. Each block device must specify the `device` property
//     along with the `diskImage`  property, which is a [[SoftLayer_Virtual_Disk_Image]] structure with the `capacity` property set. The `device` number `'1'`
//     is reserved for the SWAP disk attached to the computing instance.
//   - Default: The smallest available capacity for the primary disk will be used. If an image template is specified the disk capacity will be be provided by the template.
//   - Example:
//     ```
//     "blockDevices":[{"device": "0", "diskImage": {"capacity": 100}}],
//     "localDiskFlag": true
//     ```
//   - See [[SoftLayer_Virtual_Guest/getCreateObjectOptions]] for available options.
//
// - `userData.value`  String
//   - Arbitrary data to be made available to the computing instance.
//   - The `userData` property is an array with a single [[SoftLayer_Virtual_Guest_Attribute]] structure with the `value` property set to an arbitrary value.
//     This value can be retrieved via the [[SoftLayer_Resource_Metadata/getUserMetadata]] method from a request originating from the computing instance.
//     This is primarily useful for providing data to software that may be on the instance and configured to execute upon first boot.
//   - Example: `"userData":[{"value": "testData"}]`
//
// - `sshKeys` [[SoftLayer_Security_Ssh_Key]][]
//   - The `sshKeys` property is an array of [[SoftLayer_Security_Ssh_Key]] structures with the `id` property set to the value of an existing SSH key.
//   - To create a new SSH key, call [[SoftLayer_Security_Ssh_Key/createObject|createObject]].
//   - To obtain a list of existing SSH keys, call [[SoftLayer_Account/getSshKeys]]
//   - Example: `"sshKeys":[{"id": 1234567}]`
//
// - `postInstallScriptUri` String
//   - Specifies the uri location of the script to be downloaded and run after installation is complete. Only scripts from HTTPS servers are executed on startup.
//
// REST Example:
// ```
//
//	curl -X POST -d '{
//	    "parameters":[
//	        {
//	            "hostname": "host1",
//	            "domain": "example.com",
//	            "startCpus": 1,
//	            "maxMemory": 1024,
//	            "hourlyBillingFlag": true,
//	            "localDiskFlag": true,
//	            "operatingSystemReferenceCode": "UBUNTU_LATEST"
//	        }
//	}' https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/createObject.json
//
// HTTP/1.1 201 Created
// Location: https://api.softlayer.com/rest/v3.1/SoftLayer_Virtual_Guest/1301396/getObject
//
//	{
//	  "accountId": 232298,
//	  "createDate": "2012-11-30T16:28:17-06:00",
//	  "dedicatedAccountHostOnlyFlag": false,
//	  "domain": "example.com",
//	  "hostname": "host1",
//	  "id": 1301396,
//	  "lastPowerStateId": null,
//	  "lastVerifiedDate": null,
//	  "maxCpu": 1,
//	  "maxCpuUnits": "CORE",
//	  "maxMemory": 1024,
//	  "metricPollDate": null,
//	  "modifyDate": null,
//	  "privateNetworkOnlyFlag": false,
//	  "startCpus": 1,
//	  "statusId": 1001,
//	  "globalIdentifier": "2d203774-0ee1-49f5-9599-6ef67358dd31"
//	}
//
// ```
func (r Virtual_Guest) CreateObject(templateObject *datatypes.Virtual_Guest) (resp datatypes.Virtual_Guest, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "createObject", params, &r.Options, &resp)
	return
}

// createObjects() enables the creation of multiple computing instances on an account in a single call. This
// method is a simplified alternative to interacting with the ordering system directly.
//
// In order to create a computing instance a set of template objects must be sent in with a few required
// values.
//
// <b>Warning:</b> Computing instances created via this method will incur charges on your account.
//
// See [[SoftLayer_Virtual_Guest/createObject|createObject]] for specifics on the requirements of each template object.
//
// <h1>Example</h1>
//
//	<http title="Request">curl -X POST -d '{
//	 "parameters":[
//	     [
//	         {
//	             "hostname": "host1",
//	             "domain": "example.com",
//	             "startCpus": 1,
//	             "maxMemory": 1024,
//	             "hourlyBillingFlag": true,
//	             "localDiskFlag": true,
//	             "operatingSystemReferenceCode": "UBUNTU_LATEST"
//	         },
//	         {
//	             "hostname": "host2",
//	             "domain": "example.com",
//	             "startCpus": 1,
//	             "maxMemory": 1024,
//	             "hourlyBillingFlag": true,
//	             "localDiskFlag": true,
//	             "operatingSystemReferenceCode": "UBUNTU_LATEST"
//	         }
//	     ]
//	 ]
//	}' https://api.softlayer.com/rest/v3/SoftLayer_Virtual_Guest/createObjects.json
//
// </http>
// <http title="Response">HTTP/1.1 200 OK
//
// [
//
//	{
//	    "accountId": 232298,
//	    "createDate": "2012-11-30T23:56:48-06:00",
//	    "dedicatedAccountHostOnlyFlag": false,
//	    "domain": "softlayer.com",
//	    "hostname": "ubuntu1",
//	    "id": 1301456,
//	    "lastPowerStateId": null,
//	    "lastVerifiedDate": null,
//	    "maxCpu": 1,
//	    "maxCpuUnits": "CORE",
//	    "maxMemory": 1024,
//	    "metricPollDate": null,
//	    "modifyDate": null,
//	    "privateNetworkOnlyFlag": false,
//	    "startCpus": 1,
//	    "statusId": 1001,
//	    "globalIdentifier": "fed4c822-48c0-45d0-85e2-90476aa0c542"
//	},
//	{
//	    "accountId": 232298,
//	    "createDate": "2012-11-30T23:56:49-06:00",
//	    "dedicatedAccountHostOnlyFlag": false,
//	    "domain": "softlayer.com",
//	    "hostname": "ubuntu2",
//	    "id": 1301457,
//	    "lastPowerStateId": null,
//	    "lastVerifiedDate": null,
//	    "maxCpu": 1,
//	    "maxCpuUnits": "CORE",
//	    "maxMemory": 1024,
//	    "metricPollDate": null,
//	    "modifyDate": null,
//	    "privateNetworkOnlyFlag": false,
//	    "startCpus": 1,
//	    "statusId": 1001,
//	    "globalIdentifier": "bed4c686-9562-4ade-9049-dc4d5b6b200c"
//	}
//
// ]
// </http>
func (r Virtual_Guest) CreateObjects(templateObjects []datatypes.Virtual_Guest) (resp []datatypes.Virtual_Guest, err error) {
	params := []interface{}{
		templateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "createObjects", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) CreatePostSoftwareInstallTransaction(data *string, returnBoolean *bool) (resp bool, err error) {
	params := []interface{}{
		data,
		returnBoolean,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "createPostSoftwareInstallTransaction", params, &r.Options, &resp)
	return
}

// This method will cancel a computing instance effective immediately. For instances billed hourly, the charges will stop immediately after the method returns.
func (r Virtual_Guest) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "deleteObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) DeleteTag(tagName *string) (resp bool, err error) {
	params := []interface{}{
		tagName,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "deleteTag", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) DeleteTransientWebhook() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "deleteTransientWebhook", nil, &r.Options, &resp)
	return
}

// Creates a transaction to detach a guest's disk image. If the disk image is already detached it will be ignored.
//
// WARNING: The transaction created by this service will shut down the guest while the disk image is attached. The guest will be turned back on once this process is complete.
func (r Virtual_Guest) DetachDiskImage(imageId *int) (resp datatypes.Provisioning_Version1_Transaction, err error) {
	params := []interface{}{
		imageId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "detachDiskImage", params, &r.Options, &resp)
	return
}

// Edit a computing instance's properties
func (r Virtual_Guest) EditObject(templateObject *datatypes.Virtual_Guest) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "editObject", params, &r.Options, &resp)
	return
}

// Reboot a guest into the Idera Bare Metal Restore image.
func (r Virtual_Guest) ExecuteIderaBareMetalRestore() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "executeIderaBareMetalRestore", nil, &r.Options, &resp)
	return
}

// Reboot a guest into the R1Soft Bare Metal Restore image.
func (r Virtual_Guest) ExecuteR1SoftBareMetalRestore() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "executeR1SoftBareMetalRestore", nil, &r.Options, &resp)
	return
}

// Download and run remote script from uri on virtual guests.
func (r Virtual_Guest) ExecuteRemoteScript(uri *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		uri,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "executeRemoteScript", params, &r.Options, &resp)
	return
}

// Reboot a Linux guest into the Xen rescue image.
func (r Virtual_Guest) ExecuteRescueLayer() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "executeRescueLayer", nil, &r.Options, &resp)
	return
}

// Find VSIs by hostname.
func (r Virtual_Guest) FindByHostname(hostname *string) (resp []datatypes.Virtual_Guest, err error) {
	params := []interface{}{
		hostname,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "findByHostname", params, &r.Options, &resp)
	return
}

// Find CCI by only its primary public or private IP address. IP addresses within secondary subnets tied to the CCI will not return the CCI. If no CCI is found, no errors are generated and no data is returned.
func (r Virtual_Guest) FindByIpAddress(ipAddress *string) (resp datatypes.Virtual_Guest, err error) {
	params := []interface{}{
		ipAddress,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "findByIpAddress", params, &r.Options, &resp)
	return
}

// Obtain an [[SoftLayer_Container_Product_Order_Virtual_Guest (type)|order container]] that can be sent to [[SoftLayer_Product_Order/verifyOrder|verifyOrder]] or [[SoftLayer_Product_Order/placeOrder|placeOrder]].
//
// This is primarily useful when there is a necessity to confirm the price which will be charged for an order.
//
// See [[SoftLayer_Virtual_Guest/createObject|createObject]] for specifics on the requirements of the template object parameter.
func (r Virtual_Guest) GenerateOrderTemplate(templateObject *datatypes.Virtual_Guest) (resp datatypes.Container_Product_Order, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "generateOrderTemplate", params, &r.Options, &resp)
	return
}

// Retrieve The account that a virtual guest belongs to.
func (r Virtual_Guest) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetAccountOwnedPoolFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAccountOwnedPoolFlag", nil, &r.Options, &resp)
	return
}

// Retrieve A virtual guest's currently active network monitoring incidents.
func (r Virtual_Guest) GetActiveNetworkMonitorIncident() (resp []datatypes.Network_Monitor_Version1_Incident, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getActiveNetworkMonitorIncident", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetActiveTickets() (resp []datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getActiveTickets", nil, &r.Options, &resp)
	return
}

// Retrieve A transaction that is still be performed on a cloud server.
func (r Virtual_Guest) GetActiveTransaction() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getActiveTransaction", nil, &r.Options, &resp)
	return
}

// Retrieve Any active transaction(s) that are currently running for the server (example: os reload).
func (r Virtual_Guest) GetActiveTransactions() (resp []datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getActiveTransactions", nil, &r.Options, &resp)
	return
}

// Return a collection of SoftLayer_Item_Price objects for an OS reload
func (r Virtual_Guest) GetAdditionalRequiredPricesForOsReload(config *datatypes.Container_Hardware_Server_Configuration) (resp []datatypes.Product_Item_Price, err error) {
	params := []interface{}{
		config,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAdditionalRequiredPricesForOsReload", params, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer_Network_Storage_Allowed_Host information to connect this Virtual Guest to Network Storage volumes that require access control lists.
func (r Virtual_Guest) GetAllowedHost() (resp datatypes.Network_Storage_Allowed_Host, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAllowedHost", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer_Network_Storage objects that this SoftLayer_Virtual_Guest has access to.
func (r Virtual_Guest) GetAllowedNetworkStorage() (resp []datatypes.Network_Storage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAllowedNetworkStorage", nil, &r.Options, &resp)
	return
}

// Retrieve The SoftLayer_Network_Storage objects whose Replica that this SoftLayer_Virtual_Guest has access to.
func (r Virtual_Guest) GetAllowedNetworkStorageReplicas() (resp []datatypes.Network_Storage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAllowedNetworkStorageReplicas", nil, &r.Options, &resp)
	return
}

// Retrieve A antivirus / spyware software component object.
func (r Virtual_Guest) GetAntivirusSpywareSoftwareComponent() (resp datatypes.Software_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAntivirusSpywareSoftwareComponent", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetApplicationDeliveryController() (resp datatypes.Network_Application_Delivery_Controller, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getApplicationDeliveryController", nil, &r.Options, &resp)
	return
}

// This method is retrieve a list of SoftLayer_Network_Storage volumes that are authorized access to this SoftLayer_Virtual_Guest.
func (r Virtual_Guest) GetAttachedNetworkStorages(nasType *string) (resp []datatypes.Network_Storage, err error) {
	params := []interface{}{
		nasType,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAttachedNetworkStorages", params, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetAttributes() (resp []datatypes.Virtual_Guest_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAttributes", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) GetAvailableBlockDevicePositions() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAvailableBlockDevicePositions", nil, &r.Options, &resp)
	return
}

// Retrieve An object that stores the maximum level for the monitoring query types and response types.
func (r Virtual_Guest) GetAvailableMonitoring() (resp []datatypes.Network_Monitor_Version1_Query_Host_Stratum, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAvailableMonitoring", nil, &r.Options, &resp)
	return
}

// This method retrieves a list of SoftLayer_Network_Storage volumes that can be authorized to this SoftLayer_Virtual_Guest.
func (r Virtual_Guest) GetAvailableNetworkStorages(nasType *string) (resp []datatypes.Network_Storage, err error) {
	params := []interface{}{
		nasType,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAvailableNetworkStorages", params, &r.Options, &resp)
	return
}

// Retrieve The average daily private bandwidth usage for the current billing cycle.
func (r Virtual_Guest) GetAverageDailyPrivateBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAverageDailyPrivateBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve The average daily public bandwidth usage for the current billing cycle.
func (r Virtual_Guest) GetAverageDailyPublicBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getAverageDailyPublicBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve A guests's backend network components.
func (r Virtual_Guest) GetBackendNetworkComponents() (resp []datatypes.Virtual_Guest_Network_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBackendNetworkComponents", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's backend or private router.
func (r Virtual_Guest) GetBackendRouters() (resp []datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBackendRouters", nil, &r.Options, &resp)
	return
}

// Retrieve A computing instance's allotted bandwidth (measured in GB).
func (r Virtual_Guest) GetBandwidthAllocation() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthAllocation", nil, &r.Options, &resp)
	return
}

// Retrieve A computing instance's allotted detail record. Allotment details link bandwidth allocation with allotments.
func (r Virtual_Guest) GetBandwidthAllotmentDetail() (resp datatypes.Network_Bandwidth_Version1_Allotment_Detail, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthAllotmentDetail", nil, &r.Options, &resp)
	return
}

// Use this method when needing the metric data for bandwidth for a single guest.  It will gather the correct input parameters based on the date ranges
func (r Virtual_Guest) GetBandwidthDataByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time, networkType *string) (resp []datatypes.Metric_Tracking_Object_Data, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
		networkType,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthDataByDate", params, &r.Options, &resp)
	return
}

// Retrieve a collection of bandwidth data from an individual public or private network tracking object. Data is ideal if you with to employ your own traffic storage and graphing systems.
func (r Virtual_Guest) GetBandwidthForDateRange(startDate *datatypes.Time, endDate *datatypes.Time) (resp []datatypes.Metric_Tracking_Object_Data, err error) {
	params := []interface{}{
		startDate,
		endDate,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthForDateRange", params, &r.Options, &resp)
	return
}

// Use this method when needing a bandwidth image for a single guest.  It will gather the correct input parameters for the generic graphing utility automatically based on the snapshot specified.
func (r Virtual_Guest) GetBandwidthImage(networkType *string, snapshotRange *string, dateSpecified *datatypes.Time, dateSpecifiedEnd *datatypes.Time) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		networkType,
		snapshotRange,
		dateSpecified,
		dateSpecifiedEnd,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthImage", params, &r.Options, &resp)
	return
}

// Use this method when needing a bandwidth image for a single guest.  It will gather the correct input parameters for the generic graphing utility based on the date ranges
func (r Virtual_Guest) GetBandwidthImageByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time, networkType *string) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
		networkType,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthImageByDate", params, &r.Options, &resp)
	return
}

// Returns the total amount of bandwidth used during the time specified for a computing instance.
func (r Virtual_Guest) GetBandwidthTotal(startDateTime *datatypes.Time, endDateTime *datatypes.Time, direction *string, side *string) (resp uint, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
		direction,
		side,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBandwidthTotal", params, &r.Options, &resp)
	return
}

// Retrieve The raw bandwidth usage data for the current billing cycle. One object will be returned for each network this server is attached to.
func (r Virtual_Guest) GetBillingCycleBandwidthUsage() (resp []datatypes.Network_Bandwidth_Usage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBillingCycleBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve The raw private bandwidth usage data for the current billing cycle.
func (r Virtual_Guest) GetBillingCyclePrivateBandwidthUsage() (resp datatypes.Network_Bandwidth_Usage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBillingCyclePrivateBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve The raw public bandwidth usage data for the current billing cycle.
func (r Virtual_Guest) GetBillingCyclePublicBandwidthUsage() (resp datatypes.Network_Bandwidth_Usage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBillingCyclePublicBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve The billing item for a CloudLayer Compute Instance.
func (r Virtual_Guest) GetBillingItem() (resp datatypes.Billing_Item_Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve Determines whether the instance is ineligible for cancellation because it is disconnected.
func (r Virtual_Guest) GetBlockCancelBecauseDisconnectedFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBlockCancelBecauseDisconnectedFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The global identifier for the image template that was used to provision or reload a guest.
func (r Virtual_Guest) GetBlockDeviceTemplateGroup() (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBlockDeviceTemplateGroup", nil, &r.Options, &resp)
	return
}

// Retrieve A computing instance's block devices. Block devices link [[SoftLayer_Virtual_Disk_Image|disk images]] to computing instances.
func (r Virtual_Guest) GetBlockDevices() (resp []datatypes.Virtual_Guest_Block_Device, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBlockDevices", nil, &r.Options, &resp)
	return
}

// Retrieves the boot mode of the VSI.
func (r Virtual_Guest) GetBootMode() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBootMode", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) GetBootOrder() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBootOrder", nil, &r.Options, &resp)
	return
}

// Retrieve A virtual guest's browser access logs.
func (r Virtual_Guest) GetBrowserConsoleAccessLogs() (resp []datatypes.Virtual_BrowserConsoleAccessLog, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getBrowserConsoleAccessLogs", nil, &r.Options, &resp)
	return
}

// Gets the console access logs for a computing instance
func (r Virtual_Guest) GetConsoleAccessLog() (resp []datatypes.Network_Logging_Syslog, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getConsoleAccessLog", nil, &r.Options, &resp)
	return
}

// Retrieve A container for a guest's console data
func (r Virtual_Guest) GetConsoleData() (resp datatypes.Container_Virtual_ConsoleData, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getConsoleData", nil, &r.Options, &resp)
	return
}

// Retrieve [DEPRECATED] A flag indicating a computing instance's console IP address is assigned.
func (r Virtual_Guest) GetConsoleIpAddressFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getConsoleIpAddressFlag", nil, &r.Options, &resp)
	return
}

// Retrieve [DEPRECATED] A record containing information about a computing instance's console IP and port number.
func (r Virtual_Guest) GetConsoleIpAddressRecord() (resp datatypes.Virtual_Guest_Network_Component_IpAddress, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getConsoleIpAddressRecord", nil, &r.Options, &resp)
	return
}

// Retrieve A continuous data protection software component object.
func (r Virtual_Guest) GetContinuousDataProtectionSoftwareComponent() (resp datatypes.Software_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getContinuousDataProtectionSoftwareComponent", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's control panel.
func (r Virtual_Guest) GetControlPanel() (resp datatypes.Software_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getControlPanel", nil, &r.Options, &resp)
	return
}

// If the virtual server currently has an operating system that has a core capacity restriction, return the associated core-restricted operating system item price. Some operating systems (e.g., Red Hat Enterprise Linux) may be billed by the number of processor cores, so therefore require that a certain number of cores be present on the server.
func (r Virtual_Guest) GetCoreRestrictedOperatingSystemPrice() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCoreRestrictedOperatingSystemPrice", nil, &r.Options, &resp)
	return
}

// Use this method when needing the metric data for a single guest's CPUs.  It will gather the correct input parameters based on the date ranges
func (r Virtual_Guest) GetCpuMetricDataByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time, cpuIndexes []int) (resp []datatypes.Metric_Tracking_Object_Data, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
		cpuIndexes,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCpuMetricDataByDate", params, &r.Options, &resp)
	return
}

// Use this method when needing a cpu usage image for a single guest.  It will gather the correct input parameters for the generic graphing utility automatically based on the snapshot specified.
func (r Virtual_Guest) GetCpuMetricImage(snapshotRange *string, dateSpecified *datatypes.Time) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		snapshotRange,
		dateSpecified,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCpuMetricImage", params, &r.Options, &resp)
	return
}

// Use this method when needing a CPU usage image for a single guest.  It will gather the correct input parameters for the generic graphing utility based on the date ranges
func (r Virtual_Guest) GetCpuMetricImageByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time, cpuIndexes []int) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
		cpuIndexes,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCpuMetricImageByDate", params, &r.Options, &resp)
	return
}

// There are many options that may be provided while ordering a computing instance, this method can be used to determine what these options are.
//
// Detailed information on the return value can be found on the data type page for [[SoftLayer_Container_Virtual_Guest_Configuration (type)]].
func (r Virtual_Guest) GetCreateObjectOptions() (resp datatypes.Container_Virtual_Guest_Configuration, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCreateObjectOptions", nil, &r.Options, &resp)
	return
}

// Retrieve An object that provides commonly used bandwidth summary components for the current billing cycle.
func (r Virtual_Guest) GetCurrentBandwidthSummary() (resp datatypes.Metric_Tracking_Object_Bandwidth_Summary, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCurrentBandwidthSummary", nil, &r.Options, &resp)
	return
}

// Get the billing detail for this instance for the current billing period. This does not include bandwidth usage.
func (r Virtual_Guest) GetCurrentBillingDetail() (resp []datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCurrentBillingDetail", nil, &r.Options, &resp)
	return
}

// Get the total bill amount in US Dollars ($) for this instance in the current billing period. This includes all bandwidth used up to the point this method is called on the instance.
func (r Virtual_Guest) GetCurrentBillingTotal() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getCurrentBillingTotal", nil, &r.Options, &resp)
	return
}

// Retrieve The datacenter that a virtual guest resides in.
func (r Virtual_Guest) GetDatacenter() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getDatacenter", nil, &r.Options, &resp)
	return
}

// Retrieve The dedicated host associated with this guest.
func (r Virtual_Guest) GetDedicatedHost() (resp datatypes.Virtual_DedicatedHost, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getDedicatedHost", nil, &r.Options, &resp)
	return
}

// Retrieve The device status of this virtual guest.
func (r Virtual_Guest) GetDeviceStatus() (resp datatypes.Device_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getDeviceStatus", nil, &r.Options, &resp)
	return
}

// Return a drive retention SoftLayer_Item_Price object for a guest.
func (r Virtual_Guest) GetDriveRetentionItemPrice() (resp datatypes.Product_Item_Price, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getDriveRetentionItemPrice", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's associated EVault network storage service account.
func (r Virtual_Guest) GetEvaultNetworkStorage() (resp []datatypes.Network_Storage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getEvaultNetworkStorage", nil, &r.Options, &resp)
	return
}

// Get the subnets associated with this CloudLayer computing instance that are protectable by a network component firewall.
func (r Virtual_Guest) GetFirewallProtectableSubnets() (resp []datatypes.Network_Subnet, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getFirewallProtectableSubnets", nil, &r.Options, &resp)
	return
}

// Retrieve A computing instance's hardware firewall services.
func (r Virtual_Guest) GetFirewallServiceComponent() (resp datatypes.Network_Component_Firewall, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getFirewallServiceComponent", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) GetFirstAvailableBlockDevicePosition() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getFirstAvailableBlockDevicePosition", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's frontend network components.
func (r Virtual_Guest) GetFrontendNetworkComponents() (resp []datatypes.Virtual_Guest_Network_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getFrontendNetworkComponents", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's frontend or public router.
func (r Virtual_Guest) GetFrontendRouters() (resp datatypes.Hardware, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getFrontendRouters", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's universally unique identifier.
func (r Virtual_Guest) GetGlobalIdentifier() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getGlobalIdentifier", nil, &r.Options, &resp)
	return
}

// Retrieve The number of GPUs attached to the guest.
func (r Virtual_Guest) GetGpuCount() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getGpuCount", nil, &r.Options, &resp)
	return
}

// Retrieve The name of the GPU type attached to the guest.
func (r Virtual_Guest) GetGpuType() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getGpuType", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetGuestBootParameter() (resp datatypes.Virtual_Guest_Boot_Parameter, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getGuestBootParameter", nil, &r.Options, &resp)
	return
}

// Retrieve The object's function.
func (r Virtual_Guest) GetHardwareFunctionDescription() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getHardwareFunctionDescription", nil, &r.Options, &resp)
	return
}

// Retrieve The virtual host on which a virtual guest resides (available only on private clouds).
func (r Virtual_Guest) GetHost() (resp datatypes.Virtual_Host, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getHost", nil, &r.Options, &resp)
	return
}

// Retrieve A host IPS software component object.
func (r Virtual_Guest) GetHostIpsSoftwareComponent() (resp datatypes.Software_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getHostIpsSoftwareComponent", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's hourly billing status.
func (r Virtual_Guest) GetHourlyBillingFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getHourlyBillingFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The total private inbound bandwidth for this computing instance for the current billing cycle.
func (r Virtual_Guest) GetInboundPrivateBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getInboundPrivateBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve The total public inbound bandwidth for this computing instance for the current billing cycle.
func (r Virtual_Guest) GetInboundPublicBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getInboundPublicBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetInternalTagReferences() (resp []datatypes.Tag_Reference, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getInternalTagReferences", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) GetIsoBootImage() (resp datatypes.Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getIsoBootImage", nil, &r.Options, &resp)
	return
}

// Return a collection of SoftLayer_Item_Price objects from a collection of SoftLayer_Software_Description
func (r Virtual_Guest) GetItemPricesFromSoftwareDescriptions(softwareDescriptions []datatypes.Software_Description, includeTranslationsFlag *bool, returnAllPricesFlag *bool) (resp []datatypes.Product_Item, err error) {
	params := []interface{}{
		softwareDescriptions,
		includeTranslationsFlag,
		returnAllPricesFlag,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getItemPricesFromSoftwareDescriptions", params, &r.Options, &resp)
	return
}

// Retrieve The last known power state of a virtual guest in the event the guest is turned off outside of IMS or has gone offline.
func (r Virtual_Guest) GetLastKnownPowerState() (resp datatypes.Virtual_Guest_Power_State, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getLastKnownPowerState", nil, &r.Options, &resp)
	return
}

// Retrieve The last transaction that a cloud server's operating system was loaded.
func (r Virtual_Guest) GetLastOperatingSystemReload() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getLastOperatingSystemReload", nil, &r.Options, &resp)
	return
}

// Retrieve The last transaction a cloud server had performed.
func (r Virtual_Guest) GetLastTransaction() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getLastTransaction", nil, &r.Options, &resp)
	return
}

// Retrieve A virtual guest's latest network monitoring incident.
func (r Virtual_Guest) GetLatestNetworkMonitorIncident() (resp datatypes.Network_Monitor_Version1_Incident, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getLatestNetworkMonitorIncident", nil, &r.Options, &resp)
	return
}

// Retrieve A flag indicating that the virtual guest has at least one disk which is local to the host it runs on. This does not include a SWAP device.
func (r Virtual_Guest) GetLocalDiskFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getLocalDiskFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Where guest is located within SoftLayer's location hierarchy.
func (r Virtual_Guest) GetLocation() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getLocation", nil, &r.Options, &resp)
	return
}

// Retrieve A flag indicating that the virtual guest is a managed resource.
func (r Virtual_Guest) GetManagedResourceFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getManagedResourceFlag", nil, &r.Options, &resp)
	return
}

// Use this method when needing the metric data for memory for a single computing instance.
func (r Virtual_Guest) GetMemoryMetricDataByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time) (resp []datatypes.Metric_Tracking_Object_Data, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMemoryMetricDataByDate", params, &r.Options, &resp)
	return
}

// Use this method when needing a memory usage image for a single guest.  It will gather the correct input parameters for the generic graphing utility automatically based on the snapshot specified.
func (r Virtual_Guest) GetMemoryMetricImage(snapshotRange *string, dateSpecified *datatypes.Time) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		snapshotRange,
		dateSpecified,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMemoryMetricImage", params, &r.Options, &resp)
	return
}

// Use this method when needing a image displaying the amount of memory used over time for a single computing instance. It will gather the correct input parameters for the generic graphing utility based on the date ranges
func (r Virtual_Guest) GetMemoryMetricImageByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMemoryMetricImageByDate", params, &r.Options, &resp)
	return
}

// Retrieve A guest's metric tracking object.
func (r Virtual_Guest) GetMetricTrackingObject() (resp datatypes.Metric_Tracking_Object, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMetricTrackingObject", nil, &r.Options, &resp)
	return
}

// Retrieve The metric tracking object id for this guest.
func (r Virtual_Guest) GetMetricTrackingObjectId() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMetricTrackingObjectId", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetMonitoringRobot() (resp datatypes.Monitoring_Robot, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMonitoringRobot", nil, &r.Options, &resp)
	return
}

// Retrieve A virtual guest's network monitoring services.
func (r Virtual_Guest) GetMonitoringServiceComponent() (resp datatypes.Network_Monitor_Version1_Query_Host_Stratum, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMonitoringServiceComponent", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetMonitoringServiceEligibilityFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMonitoringServiceEligibilityFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The monitoring notification objects for this guest. Each object links this guest instance to a user account that will be notified if monitoring on this guest object fails
func (r Virtual_Guest) GetMonitoringUserNotification() (resp []datatypes.User_Customer_Notification_Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getMonitoringUserNotification", nil, &r.Options, &resp)
	return
}

// Get the IP addresses associated with this CloudLayer computing instance that are protectable by a network component firewall. Note, this may not return all values for IPv6 subnets for this CloudLayer computing instance. Please use getFirewallProtectableSubnets to get all protectable subnets.
func (r Virtual_Guest) GetNetworkComponentFirewallProtectableIpAddresses() (resp []datatypes.Network_Subnet_IpAddress, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getNetworkComponentFirewallProtectableIpAddresses", nil, &r.Options, &resp)
	return
}

// Retrieve A guests's network components.
func (r Virtual_Guest) GetNetworkComponents() (resp []datatypes.Virtual_Guest_Network_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getNetworkComponents", nil, &r.Options, &resp)
	return
}

// Retrieve All of a virtual guest's network monitoring incidents.
func (r Virtual_Guest) GetNetworkMonitorIncidents() (resp []datatypes.Network_Monitor_Version1_Incident, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getNetworkMonitorIncidents", nil, &r.Options, &resp)
	return
}

// Retrieve A guests's network monitors.
func (r Virtual_Guest) GetNetworkMonitors() (resp []datatypes.Network_Monitor_Version1_Query_Host, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getNetworkMonitors", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's associated network storage accounts.
func (r Virtual_Guest) GetNetworkStorage() (resp []datatypes.Network_Storage, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getNetworkStorage", nil, &r.Options, &resp)
	return
}

// Retrieve The network Vlans that a guest's network components are associated with.
func (r Virtual_Guest) GetNetworkVlans() (resp []datatypes.Network_Vlan, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getNetworkVlans", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) GetObject() (resp datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve An open ticket requesting cancellation of this server, if one exists.
func (r Virtual_Guest) GetOpenCancellationTicket() (resp datatypes.Ticket, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOpenCancellationTicket", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's operating system.
func (r Virtual_Guest) GetOperatingSystem() (resp datatypes.Software_Component_OperatingSystem, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOperatingSystem", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's operating system software description.
func (r Virtual_Guest) GetOperatingSystemReferenceCode() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOperatingSystemReferenceCode", nil, &r.Options, &resp)
	return
}

// Obtain an order container that is ready to be sent to the [[SoftLayer_Product_Order#placeOrder|SoftLayer_Product_Order::placeOrder]] method. This container will include all services that the selected computing instance has. If desired you may remove prices which were returned.
func (r Virtual_Guest) GetOrderTemplate(billingType *string, orderPrices []datatypes.Product_Item_Price) (resp datatypes.Container_Product_Order, err error) {
	params := []interface{}{
		billingType,
		orderPrices,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOrderTemplate", params, &r.Options, &resp)
	return
}

// Retrieve The original package id provided with the order for a Cloud Computing Instance.
func (r Virtual_Guest) GetOrderedPackageId() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOrderedPackageId", nil, &r.Options, &resp)
	return
}

// Retrieve The total private outbound bandwidth for this computing instance for the current billing cycle.
func (r Virtual_Guest) GetOutboundPrivateBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOutboundPrivateBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve The total public outbound bandwidth for this computing instance for the current billing cycle.
func (r Virtual_Guest) GetOutboundPublicBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOutboundPublicBandwidthUsage", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the bandwidth usage for this computing instance for the current billing cycle exceeds the allocation.
func (r Virtual_Guest) GetOverBandwidthAllocationFlag() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getOverBandwidthAllocationFlag", nil, &r.Options, &resp)
	return
}

// Returns a list of all the pending maintenance actions affecting this guest.
func (r Virtual_Guest) GetPendingMaintenanceActions() (resp []datatypes.Container_Virtual_Guest_PendingMaintenanceAction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPendingMaintenanceActions", nil, &r.Options, &resp)
	return
}

// Retrieve When true this virtual guest must be migrated using SoftLayer_Virtual_Guest::migrate.
func (r Virtual_Guest) GetPendingMigrationFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPendingMigrationFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The placement group that a virtual guest belongs to.
func (r Virtual_Guest) GetPlacementGroup() (resp datatypes.Virtual_PlacementGroup, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPlacementGroup", nil, &r.Options, &resp)
	return
}

// Retrieve The current power state of a virtual guest.
func (r Virtual_Guest) GetPowerState() (resp datatypes.Virtual_Guest_Power_State, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPowerState", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's primary private IP address.
func (r Virtual_Guest) GetPrimaryBackendIpAddress() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPrimaryBackendIpAddress", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's primary backend network component.
func (r Virtual_Guest) GetPrimaryBackendNetworkComponent() (resp datatypes.Virtual_Guest_Network_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPrimaryBackendNetworkComponent", nil, &r.Options, &resp)
	return
}

// Retrieve The guest's primary public IP address.
func (r Virtual_Guest) GetPrimaryIpAddress() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPrimaryIpAddress", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's primary public network component.
func (r Virtual_Guest) GetPrimaryNetworkComponent() (resp datatypes.Virtual_Guest_Network_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPrimaryNetworkComponent", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the computing instance only has access to the private network.
func (r Virtual_Guest) GetPrivateNetworkOnlyFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getPrivateNetworkOnlyFlag", nil, &r.Options, &resp)
	return
}

// Retrieve Whether the bandwidth usage for this computing instance for the current billing cycle is projected to exceed the allocation.
func (r Virtual_Guest) GetProjectedOverBandwidthAllocationFlag() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getProjectedOverBandwidthAllocationFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The projected public outbound bandwidth for this computing instance for the current billing cycle.
func (r Virtual_Guest) GetProjectedPublicBandwidthUsage() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getProjectedPublicBandwidthUsage", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) GetProvisionDate() (resp datatypes.Time, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getProvisionDate", nil, &r.Options, &resp)
	return
}

// Retrieve Recent events that impact this computing instance.
func (r Virtual_Guest) GetRecentEvents() (resp []datatypes.Notification_Occurrence_Event, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getRecentEvents", nil, &r.Options, &resp)
	return
}

// Recent metric data for a guest
func (r Virtual_Guest) GetRecentMetricData(time *uint) (resp []datatypes.Metric_Tracking_Object, err error) {
	params := []interface{}{
		time,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getRecentMetricData", params, &r.Options, &resp)
	return
}

// Retrieve The regional group this guest is in.
func (r Virtual_Guest) GetRegionalGroup() (resp datatypes.Location_Group_Regional, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getRegionalGroup", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetRegionalInternetRegistry() (resp datatypes.Network_Regional_Internet_Registry, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getRegionalInternetRegistry", nil, &r.Options, &resp)
	return
}

// Retrieve The reserved capacity group the guest is associated with.
func (r Virtual_Guest) GetReservedCapacityGroup() (resp datatypes.Virtual_ReservedCapacityGroup, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getReservedCapacityGroup", nil, &r.Options, &resp)
	return
}

// Retrieve Flag to indicate whether or not a guest is part of a reserved capacity group.
func (r Virtual_Guest) GetReservedCapacityGroupFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getReservedCapacityGroupFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The reserved capacity group instance the guest is associated with.
func (r Virtual_Guest) GetReservedCapacityGroupInstance() (resp datatypes.Virtual_ReservedCapacityGroup_Instance, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getReservedCapacityGroupInstance", nil, &r.Options, &resp)
	return
}

// Retrieve the reverse domain records associated with this server.
func (r Virtual_Guest) GetReverseDomainRecords() (resp []datatypes.Dns_Domain, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getReverseDomainRecords", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's vulnerability scan requests.
func (r Virtual_Guest) GetSecurityScanRequests() (resp []datatypes.Network_Security_Scanner_Request, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getSecurityScanRequests", nil, &r.Options, &resp)
	return
}

// Retrieve The server room that a guest is located at. There may be more than one server room for every data center.
func (r Virtual_Guest) GetServerRoom() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getServerRoom", nil, &r.Options, &resp)
	return
}

// Retrieve A guest's installed software.
func (r Virtual_Guest) GetSoftwareComponents() (resp []datatypes.Software_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getSoftwareComponents", nil, &r.Options, &resp)
	return
}

// Retrieve SSH keys to be installed on the server during provisioning or an OS reload.
func (r Virtual_Guest) GetSshKeys() (resp []datatypes.Security_Ssh_Key, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getSshKeys", nil, &r.Options, &resp)
	return
}

// Retrieve A computing instance's status.
func (r Virtual_Guest) GetStatus() (resp datatypes.Virtual_Guest_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getStatus", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest) GetTagReferences() (resp []datatypes.Tag_Reference, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getTagReferences", nil, &r.Options, &resp)
	return
}

// Retrieve Whether or not a computing instance is a Transient Instance.
func (r Virtual_Guest) GetTransientGuestFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getTransientGuestFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The endpoint used to notify customers their transient guest is terminating.
func (r Virtual_Guest) GetTransientWebhookURI() (resp datatypes.Virtual_Guest_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getTransientWebhookURI", nil, &r.Options, &resp)
	return
}

// Retrieve The type of this virtual guest.
func (r Virtual_Guest) GetType() (resp datatypes.Virtual_Guest_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getType", nil, &r.Options, &resp)
	return
}

// Retrieves a list of all upgrades available to a virtual server. Upgradeable items include, but are not limited to, number of cores, amount of RAM, storage configuration, and network port speed.
//
// This method exclude downgrade item prices by default. You can set the "includeDowngradeItemPrices" parameter to true so that it can include downgrade item prices.
func (r Virtual_Guest) GetUpgradeItemPrices(includeDowngradeItemPrices *bool) (resp []datatypes.Product_Item_Price, err error) {
	params := []interface{}{
		includeDowngradeItemPrices,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getUpgradeItemPrices", params, &r.Options, &resp)
	return
}

// Retrieve A computing instance's associated upgrade request object if any.
func (r Virtual_Guest) GetUpgradeRequest() (resp datatypes.Product_Upgrade_Request, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getUpgradeRequest", nil, &r.Options, &resp)
	return
}

// Retrieve A base64 encoded string containing custom user data for a Cloud Computing Instance order.
func (r Virtual_Guest) GetUserData() (resp []datatypes.Virtual_Guest_Attribute, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getUserData", nil, &r.Options, &resp)
	return
}

// Retrieve A list of users that have access to this computing instance.
func (r Virtual_Guest) GetUsers() (resp []datatypes.User_Customer, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getUsers", nil, &r.Options, &resp)
	return
}

// This method will return the list of block device template groups that are valid to the host. For instance, it will validate that the template groups returned are compatible with the size and number of disks on the host.
func (r Virtual_Guest) GetValidBlockDeviceTemplateGroups(visibility *string) (resp []datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	params := []interface{}{
		visibility,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getValidBlockDeviceTemplateGroups", params, &r.Options, &resp)
	return
}

// Retrieve The name of the bandwidth allotment that a hardware belongs too.
func (r Virtual_Guest) GetVirtualRack() (resp datatypes.Network_Bandwidth_Version1_Allotment, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getVirtualRack", nil, &r.Options, &resp)
	return
}

// Retrieve The id of the bandwidth allotment that a computing instance belongs too.
func (r Virtual_Guest) GetVirtualRackId() (resp int, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getVirtualRackId", nil, &r.Options, &resp)
	return
}

// Retrieve The name of the bandwidth allotment that a computing instance belongs too.
func (r Virtual_Guest) GetVirtualRackName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "getVirtualRackName", nil, &r.Options, &resp)
	return
}

// Issues a ping command and returns the success (true) or failure (false) of the ping command.
func (r Virtual_Guest) IsBackendPingable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "isBackendPingable", nil, &r.Options, &resp)
	return
}

// Determines if the virtual guest was provisioned from a cloud-init enabled image.
func (r Virtual_Guest) IsCloudInit() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "isCloudInit", nil, &r.Options, &resp)
	return
}

// Issues a ping command and returns the success (true) or failure (false) of the ping command.
func (r Virtual_Guest) IsPingable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "isPingable", nil, &r.Options, &resp)
	return
}

// Closes the public or private ports to isolate the instance before a destructive action.
func (r Virtual_Guest) IsolateInstanceForDestructiveAction() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "isolateInstanceForDestructiveAction", nil, &r.Options, &resp)
	return
}

// Creates a transaction to migrate a virtual guest to a new host. NOTE: Will only migrate if SoftLayer_Virtual_Guest property pendingMigrationFlag = true
func (r Virtual_Guest) Migrate() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "migrate", nil, &r.Options, &resp)
	return
}

// Create a transaction to migrate an instance from one dedicated host to another dedicated host
func (r Virtual_Guest) MigrateDedicatedHost(destinationHostId *int) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		destinationHostId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "migrateDedicatedHost", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) MountIsoImage(diskImageId *int) (resp datatypes.Provisioning_Version1_Transaction, err error) {
	params := []interface{}{
		diskImageId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "mountIsoImage", params, &r.Options, &resp)
	return
}

// Pause a virtual guest. This can only be called when the specified VM is in the Running state.
func (r Virtual_Guest) Pause() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "pause", nil, &r.Options, &resp)
	return
}

// Power cycle a virtual guest
func (r Virtual_Guest) PowerCycle() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "powerCycle", nil, &r.Options, &resp)
	return
}

// Power off a virtual guest
func (r Virtual_Guest) PowerOff() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "powerOff", nil, &r.Options, &resp)
	return
}

// Power off a virtual guest
func (r Virtual_Guest) PowerOffSoft() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "powerOffSoft", nil, &r.Options, &resp)
	return
}

// Power on a virtual guest
func (r Virtual_Guest) PowerOn() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "powerOn", nil, &r.Options, &resp)
	return
}

// Power cycle a virtual guest
func (r Virtual_Guest) RebootDefault() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "rebootDefault", nil, &r.Options, &resp)
	return
}

// Power cycle a guest.
func (r Virtual_Guest) RebootHard() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "rebootHard", nil, &r.Options, &resp)
	return
}

// Attempt to complete a soft reboot of a guest by shutting down the operating system.
func (r Virtual_Guest) RebootSoft() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "rebootSoft", nil, &r.Options, &resp)
	return
}

// no documentation yet
// Deprecated: This function has been marked as deprecated.
func (r Virtual_Guest) ReconfigureConsole() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "reconfigureConsole", nil, &r.Options, &resp)
	return
}

// Create a transaction to perform an OS reload
func (r Virtual_Guest) ReloadCurrentOperatingSystemConfiguration() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "reloadCurrentOperatingSystemConfiguration", nil, &r.Options, &resp)
	return
}

// Reloads current operating system configuration.
//
// This service has a confirmation protocol for proceeding with the reload. To proceed with the reload without confirmation, simply pass in 'FORCE' as the token parameter. To proceed with the reload with confirmation, simply call the service with no parameter. A token string will be returned by this service. The token will remain active for 10 minutes. Use this token as the parameter to confirm that a reload is to be performed for the server.
//
// As a precaution, we strongly  recommend backing up all data before reloading the operating system. The reload will format the primary disk and will reconfigure the computing instance to the current specifications on record.
//
// If reloading from an image template, we recommend first getting the list of valid private block device template groups, by calling the getOperatingSystemReloadImages method.
func (r Virtual_Guest) ReloadOperatingSystem(token *string, config *datatypes.Container_Hardware_Server_Configuration) (resp string, err error) {
	params := []interface{}{
		token,
		config,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "reloadOperatingSystem", params, &r.Options, &resp)
	return
}

// This method is used to remove access to a SoftLayer_Network_Storage volume that supports host- or network-level access control.
func (r Virtual_Guest) RemoveAccessToNetworkStorage(networkStorageTemplateObject *datatypes.Network_Storage) (resp bool, err error) {
	params := []interface{}{
		networkStorageTemplateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "removeAccessToNetworkStorage", params, &r.Options, &resp)
	return
}

// This method is used to allow access to multiple SoftLayer_Network_Storage volumes that support host- or network-level access control.
func (r Virtual_Guest) RemoveAccessToNetworkStorageList(networkStorageTemplateObjects []datatypes.Network_Storage) (resp bool, err error) {
	params := []interface{}{
		networkStorageTemplateObjects,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "removeAccessToNetworkStorageList", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) RemoveTags(tags *string) (resp bool, err error) {
	params := []interface{}{
		tags,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "removeTags", params, &r.Options, &resp)
	return
}

// Resume a virtual guest, this can only be called when a VSI is in Suspended state.
func (r Virtual_Guest) Resume() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "resume", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) SendTestReclaimScheduledAlert() (err error) {
	var resp datatypes.Void
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "sendTestReclaimScheduledAlert", nil, &r.Options, &resp)
	return
}

// Sets the private network interface speed to the new speed. Speed values can only be 0 (Disconnect), 10, 100, or 1000. The new speed must be equal to or less than the max speed of the interface.
//
// It will take less than a minute to update the port speed.
func (r Virtual_Guest) SetPrivateNetworkInterfaceSpeed(newSpeed *int) (resp bool, err error) {
	params := []interface{}{
		newSpeed,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "setPrivateNetworkInterfaceSpeed", params, &r.Options, &resp)
	return
}

// Sets the public network interface speed to the new speed. Speed values can only be 0 (Disconnect), 10, 100, or 1000. The new speed must be equal to or less than the max speed of the interface.
//
// It will take less than a minute to update the port speed.
func (r Virtual_Guest) SetPublicNetworkInterfaceSpeed(newSpeed *int) (resp bool, err error) {
	params := []interface{}{
		newSpeed,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "setPublicNetworkInterfaceSpeed", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) SetTags(tags *string) (resp bool, err error) {
	params := []interface{}{
		tags,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "setTags", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) SetTransientWebhook(uri *string, secret *string) (err error) {
	var resp datatypes.Void
	params := []interface{}{
		uri,
		secret,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "setTransientWebhook", params, &r.Options, &resp)
	return
}

// Sets the data that will be written to the configuration drive.
func (r Virtual_Guest) SetUserMetadata(metadata []string) (resp bool, err error) {
	params := []interface{}{
		metadata,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "setUserMetadata", params, &r.Options, &resp)
	return
}

// Shuts down the private network port
func (r Virtual_Guest) ShutdownPrivatePort() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "shutdownPrivatePort", nil, &r.Options, &resp)
	return
}

// Shuts down the public network port
func (r Virtual_Guest) ShutdownPublicPort() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "shutdownPublicPort", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest) UnmountIsoImage() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "unmountIsoImage", nil, &r.Options, &resp)
	return
}

// Validate an image template for OS Reload
func (r Virtual_Guest) ValidateImageTemplate(imageTemplateId *int) (resp bool, err error) {
	params := []interface{}{
		imageTemplateId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "validateImageTemplate", params, &r.Options, &resp)
	return
}

// Verify that a virtual server can go through the operating system reload process. It may be useful to call this method before attempting to actually reload the operating system just to verify that the reload will go smoothly. If the server configuration is not setup correctly or there is some other issue, an exception will be thrown indicating the error. If there were no issues, this will just return true.
func (r Virtual_Guest) VerifyReloadOperatingSystem(config *datatypes.Container_Hardware_Server_Configuration) (resp bool, err error) {
	params := []interface{}{
		config,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest", "verifyReloadOperatingSystem", params, &r.Options, &resp)
	return
}

// The virtual block device template group data type presents the structure in which a group of archived image templates will be presented. The structure consists of a parent template group which contain multiple child template group objects.  Each child template group object represents the image template in a particular location. Unless editing/deleting a specific child template group object, it is best to use the parent object.
//
// A virtual block device template group, also known as an image template group, represents an image of a virtual guest instance.
type Virtual_Guest_Block_Device_Template_Group struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualGuestBlockDeviceTemplateGroupService returns an instance of the Virtual_Guest_Block_Device_Template_Group SoftLayer service
func GetVirtualGuestBlockDeviceTemplateGroupService(sess session.SLSession) Virtual_Guest_Block_Device_Template_Group {
	return Virtual_Guest_Block_Device_Template_Group{Session: sess}
}

func (r Virtual_Guest_Block_Device_Template_Group) Id(id int) Virtual_Guest_Block_Device_Template_Group {
	r.Options.Id = &id
	return r
}

func (r Virtual_Guest_Block_Device_Template_Group) Mask(mask string) Virtual_Guest_Block_Device_Template_Group {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Guest_Block_Device_Template_Group) Filter(filter string) Virtual_Guest_Block_Device_Template_Group {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Guest_Block_Device_Template_Group) Limit(limit int) Virtual_Guest_Block_Device_Template_Group {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Guest_Block_Device_Template_Group) Offset(offset int) Virtual_Guest_Block_Device_Template_Group {
	r.Options.Offset = &offset
	return r
}

// This method allows you to mark this image template as customer managed software license (BYOL)
func (r Virtual_Guest_Block_Device_Template_Group) AddByolAttribute() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "addByolAttribute", nil, &r.Options, &resp)
	return
}

// This method allows you to mark this image template as cloud init
func (r Virtual_Guest_Block_Device_Template_Group) AddCloudInitAttribute() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "addCloudInitAttribute", nil, &r.Options, &resp)
	return
}

// This method will create transaction(s) to add available locations to an archive image template.
func (r Virtual_Guest_Block_Device_Template_Group) AddLocations(locations []datatypes.Location) (resp bool, err error) {
	params := []interface{}{
		locations,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "addLocations", params, &r.Options, &resp)
	return
}

// This method allows you to mark this image's supported boot modes as 'HVM' or 'PV'.
func (r Virtual_Guest_Block_Device_Template_Group) AddSupportedBootMode(bootMode *string) (resp bool, err error) {
	params := []interface{}{
		bootMode,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "addSupportedBootMode", params, &r.Options, &resp)
	return
}

// Create a transaction to export/copy a template to an external source.
func (r Virtual_Guest_Block_Device_Template_Group) CopyToExternalSource(configuration *datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (resp bool, err error) {
	params := []interface{}{
		configuration,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "copyToExternalSource", params, &r.Options, &resp)
	return
}

// Create a transaction to export/copy a template to an ICOS.
func (r Virtual_Guest_Block_Device_Template_Group) CopyToIcos(configuration *datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (resp bool, err error) {
	params := []interface{}{
		configuration,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "copyToIcos", params, &r.Options, &resp)
	return
}

// Create a transaction to import a disk image from an external source and create a standard image template.
func (r Virtual_Guest_Block_Device_Template_Group) CreateFromExternalSource(configuration *datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	params := []interface{}{
		configuration,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "createFromExternalSource", params, &r.Options, &resp)
	return
}

// Create a process to import a disk image from ICOS and create a standard
func (r Virtual_Guest_Block_Device_Template_Group) CreateFromIcos(configuration *datatypes.Container_Virtual_Guest_Block_Device_Template_Configuration) (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	params := []interface{}{
		configuration,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "createFromIcos", params, &r.Options, &resp)
	return
}

// Create a transaction to copy archived block devices into public repository
func (r Virtual_Guest_Block_Device_Template_Group) CreatePublicArchiveTransaction(groupName *string, summary *string, note *string, locations []datatypes.Location) (resp int, err error) {
	params := []interface{}{
		groupName,
		summary,
		note,
		locations,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "createPublicArchiveTransaction", params, &r.Options, &resp)
	return
}

// This method allows you to remove BYOL attribute for a given image template.
func (r Virtual_Guest_Block_Device_Template_Group) DeleteByolAttribute() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "deleteByolAttribute", nil, &r.Options, &resp)
	return
}

// This method allows you to remove cloud init attribute for a given image template.
func (r Virtual_Guest_Block_Device_Template_Group) DeleteCloudInitAttribute() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "deleteCloudInitAttribute", nil, &r.Options, &resp)
	return
}

// Deleting a block device template group is different from the deletion of other objects.  A block device template group can contain several gigabytes of data in its disk images.  This may take some time to delete and requires a transaction to be created.  This method creates a transaction that will delete all resources associated with the block device template group.
func (r Virtual_Guest_Block_Device_Template_Group) DeleteObject() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "deleteObject", nil, &r.Options, &resp)
	return
}

// This method will deny another SoftLayer customer account's previously given access to provision CloudLayer Computing Instances from an image template group. Template access should only be removed from the parent template group object, not the child.
func (r Virtual_Guest_Block_Device_Template_Group) DenySharingAccess(accountId *int) (resp bool, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "denySharingAccess", params, &r.Options, &resp)
	return
}

// Edit an image template group's associated name and note. All other properties in the SoftLayer_Virtual_Guest_Block_Device_Template_Group data type are read-only.
func (r Virtual_Guest_Block_Device_Template_Group) EditObject(templateObject *datatypes.Virtual_Guest_Block_Device_Template_Group) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "editObject", params, &r.Options, &resp)
	return
}

// Find block device template groups containing a GC enabled cloudinit image for the current active user. A sorted collection of groups is returned. The Caller can optionally specify data center or region names to retrieve GC images from only those locations.
func (r Virtual_Guest_Block_Device_Template_Group) FindGcImagesByCurrentUser(dataCenters []string, regions []string) (resp []datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	params := []interface{}{
		dataCenters,
		regions,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "findGcImagesByCurrentUser", params, &r.Options, &resp)
	return
}

// Retrieve A block device template group's [[SoftLayer_Account|account]].
func (r Virtual_Guest_Block_Device_Template_Group) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Block_Device_Template_Group) GetAccountContacts() (resp []datatypes.Account_Contact, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getAccountContacts", nil, &r.Options, &resp)
	return
}

// Retrieve The accounts which may have read-only access to an image template group. Will only be populated for parent template group objects.
func (r Virtual_Guest_Block_Device_Template_Group) GetAccountReferences() (resp []datatypes.Virtual_Guest_Block_Device_Template_Group_Accounts, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getAccountReferences", nil, &r.Options, &resp)
	return
}

// Get all available compatible platform names that can be added to a template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetAllAvailableCompatiblePlatformNames() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getAllAvailableCompatiblePlatformNames", nil, &r.Options, &resp)
	return
}

// Retrieve The block devices that are part of an image template group
func (r Virtual_Guest_Block_Device_Template_Group) GetBlockDevices() (resp []datatypes.Virtual_Guest_Block_Device_Template, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getBlockDevices", nil, &r.Options, &resp)
	return
}

// Retrieve The total disk space of all images in a image template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetBlockDevicesDiskSpaceTotal() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getBlockDevicesDiskSpaceTotal", nil, &r.Options, &resp)
	return
}

// This method returns the boot mode, if any, set on a given image template.
func (r Virtual_Guest_Block_Device_Template_Group) GetBootMode() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getBootMode", nil, &r.Options, &resp)
	return
}

// Retrieve A flag indicating that customer is providing the software licenses.
func (r Virtual_Guest_Block_Device_Template_Group) GetByolFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getByolFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The image template groups that are clones of an image template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetChildren() (resp []datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getChildren", nil, &r.Options, &resp)
	return
}

// Get compatible platform names currently set on the template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetCurrentCompatiblePlatformNames() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getCurrentCompatiblePlatformNames", nil, &r.Options, &resp)
	return
}

// Retrieve The location containing this image template group. Will only be populated for child template group objects.
func (r Virtual_Guest_Block_Device_Template_Group) GetDatacenter() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getDatacenter", nil, &r.Options, &resp)
	return
}

// Retrieve A collection of locations containing a copy of this image template group. Will only be populated for parent template group objects.
func (r Virtual_Guest_Block_Device_Template_Group) GetDatacenters() (resp []datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getDatacenters", nil, &r.Options, &resp)
	return
}

// This method returns the default boot mode set by the software description
func (r Virtual_Guest_Block_Device_Template_Group) GetDefaultBootMode() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getDefaultBootMode", nil, &r.Options, &resp)
	return
}

// This method returns an array of encryption values, or empty array if none are found
func (r Virtual_Guest_Block_Device_Template_Group) GetEncryptionAttributes() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getEncryptionAttributes", nil, &r.Options, &resp)
	return
}

// Retrieve The first clone of the image template group
func (r Virtual_Guest_Block_Device_Template_Group) GetFirstChild() (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getFirstChild", nil, &r.Options, &resp)
	return
}

// Retrieve A flag indicating if this is a flex image.
func (r Virtual_Guest_Block_Device_Template_Group) GetFlexImageFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getFlexImageFlag", nil, &r.Options, &resp)
	return
}

// Retrieve An image template's universally unique identifier.
func (r Virtual_Guest_Block_Device_Template_Group) GetGlobalIdentifier() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getGlobalIdentifier", nil, &r.Options, &resp)
	return
}

// Retrieve The virtual disk image type of this template. Value will be populated on parent and child, but only supports object filtering on the parent.
func (r Virtual_Guest_Block_Device_Template_Group) GetImageType() (resp datatypes.Virtual_Disk_Image_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getImageType", nil, &r.Options, &resp)
	return
}

// Retrieve The virtual disk image type keyname (e.g. SYSTEM, DISK_CAPTURE, ISO, etc) of this template. Value will be populated on parent and child, but only supports object filtering on the parent.
func (r Virtual_Guest_Block_Device_Template_Group) GetImageTypeKeyName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getImageTypeKeyName", nil, &r.Options, &resp)
	return
}

// Retrieve A flag indicating if this is a next generation image.
func (r Virtual_Guest_Block_Device_Template_Group) GetNextGenFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getNextGenFlag", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Block_Device_Template_Group) GetObject() (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The image template group that another image template group was cloned from.
func (r Virtual_Guest_Block_Device_Template_Group) GetParent() (resp datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getParent", nil, &r.Options, &resp)
	return
}

// This method gets all public customer owned image templates that the user is allowed to see.
func (r Virtual_Guest_Block_Device_Template_Group) GetPublicCustomerOwnedImages() (resp []datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getPublicCustomerOwnedImages", nil, &r.Options, &resp)
	return
}

// This method gets all public image templates that the user is allowed to see.
func (r Virtual_Guest_Block_Device_Template_Group) GetPublicImages() (resp []datatypes.Virtual_Guest_Block_Device_Template_Group, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getPublicImages", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Block_Device_Template_Group) GetRegion() (resp datatypes.Network_Service_Resource, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getRegion", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Block_Device_Template_Group) GetRegions() (resp []datatypes.Network_Service_Resource, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getRegions", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Block_Device_Template_Group) GetRiasAccount(secret *string) (resp datatypes.Container_Virtual_Guest_Block_Device_Template_Group_RiasAccount, err error) {
	params := []interface{}{
		secret,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getRiasAccount", params, &r.Options, &resp)
	return
}

// Retrieve The ssh keys to be implemented on the server when provisioned or reloaded from an image template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetSshKeys() (resp []datatypes.Security_Ssh_Key, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getSshKeys", nil, &r.Options, &resp)
	return
}

// Retrieve A template group's status.
func (r Virtual_Guest_Block_Device_Template_Group) GetStatus() (resp datatypes.Virtual_Guest_Block_Device_Template_Group_Status, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getStatus", nil, &r.Options, &resp)
	return
}

// Returns the image storage locations.
func (r Virtual_Guest_Block_Device_Template_Group) GetStorageLocations() (resp []datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getStorageLocations", nil, &r.Options, &resp)
	return
}

// Retrieve The storage repository that an image template group resides on.
func (r Virtual_Guest_Block_Device_Template_Group) GetStorageRepository() (resp datatypes.Virtual_Storage_Repository, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getStorageRepository", nil, &r.Options, &resp)
	return
}

// This method indicates which boot modes are supported by the image.
func (r Virtual_Guest_Block_Device_Template_Group) GetSupportedBootModes() (resp []string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getSupportedBootModes", nil, &r.Options, &resp)
	return
}

// Retrieve The tags associated with this image template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetTagReferences() (resp []datatypes.Tag_Reference, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getTagReferences", nil, &r.Options, &resp)
	return
}

// This method allows you to grab the first data center that the image(s) reside on so we can pull it from there.
func (r Virtual_Guest_Block_Device_Template_Group) GetTemplateDataCenterName() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getTemplateDataCenterName", nil, &r.Options, &resp)
	return
}

// Retrieve A transaction that is being performed on a image template group.
func (r Virtual_Guest_Block_Device_Template_Group) GetTransaction() (resp datatypes.Provisioning_Version1_Transaction, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getTransaction", nil, &r.Options, &resp)
	return
}

// Returns an array of SoftLayer_Software_Description that are supported for VHD imports.
func (r Virtual_Guest_Block_Device_Template_Group) GetVhdImportSoftwareDescriptions() (resp []datatypes.Software_Description, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "getVhdImportSoftwareDescriptions", nil, &r.Options, &resp)
	return
}

// This method indicates whether or not this image is a customer supplied license image.
func (r Virtual_Guest_Block_Device_Template_Group) IsByol() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "isByol", nil, &r.Options, &resp)
	return
}

// This method indicates whether or not this image uses an operating system capable of using a customer supplied license image.
func (r Virtual_Guest_Block_Device_Template_Group) IsByolCapableOperatingSystem() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "isByolCapableOperatingSystem", nil, &r.Options, &resp)
	return
}

// This method indicates whether or not this image uses an operating system that requires using a customer supplied license image
func (r Virtual_Guest_Block_Device_Template_Group) IsByolOnlyOperatingSystem() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "isByolOnlyOperatingSystem", nil, &r.Options, &resp)
	return
}

// This method indicates whether or not this image is a cloud-init image.
func (r Virtual_Guest_Block_Device_Template_Group) IsCloudInit() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "isCloudInit", nil, &r.Options, &resp)
	return
}

// This method indicates whether or not this image uses an operating system that requires cloud init
func (r Virtual_Guest_Block_Device_Template_Group) IsCloudInitOnlyOperatingSystem() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "isCloudInitOnlyOperatingSystem", nil, &r.Options, &resp)
	return
}

// This method indicates whether this image template contains an encrypted disk image.
func (r Virtual_Guest_Block_Device_Template_Group) IsEncrypted() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "isEncrypted", nil, &r.Options, &resp)
	return
}

// This method will permit another SoftLayer customer account access to provision CloudLayer Computing Instances from an image template group. Template access should only be given to the parent template group object, not the child.
func (r Virtual_Guest_Block_Device_Template_Group) PermitSharingAccess(accountId *int) (resp bool, err error) {
	params := []interface{}{
		accountId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "permitSharingAccess", params, &r.Options, &resp)
	return
}

// Removes compatible platforms on the template group.
func (r Virtual_Guest_Block_Device_Template_Group) RemoveCompatiblePlatforms(platformNames []string) (resp bool, err error) {
	params := []interface{}{
		platformNames,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "removeCompatiblePlatforms", params, &r.Options, &resp)
	return
}

// This method will create transaction(s) to remove available locations from an archive image template.
func (r Virtual_Guest_Block_Device_Template_Group) RemoveLocations(locations []datatypes.Location) (resp bool, err error) {
	params := []interface{}{
		locations,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "removeLocations", params, &r.Options, &resp)
	return
}

// This method allows you to remove a supported boot mode attribute for a given image template.
func (r Virtual_Guest_Block_Device_Template_Group) RemoveSupportedBootMode(bootMode *string) (resp bool, err error) {
	params := []interface{}{
		bootMode,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "removeSupportedBootMode", params, &r.Options, &resp)
	return
}

// Create transaction(s) to set the archived block device available locations
func (r Virtual_Guest_Block_Device_Template_Group) SetAvailableLocations(locations []datatypes.Location) (resp bool, err error) {
	params := []interface{}{
		locations,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "setAvailableLocations", params, &r.Options, &resp)
	return
}

// This method allows you to specify the boot mode for a given image template.
func (r Virtual_Guest_Block_Device_Template_Group) SetBootMode(newBootMode *string) (resp bool, err error) {
	params := []interface{}{
		newBootMode,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "setBootMode", params, &r.Options, &resp)
	return
}

// Sets compatible platforms on the template group.
func (r Virtual_Guest_Block_Device_Template_Group) SetCompatiblePlatforms(platformNames []string) (resp bool, err error) {
	params := []interface{}{
		platformNames,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "setCompatiblePlatforms", params, &r.Options, &resp)
	return
}

// Set the tags for this template group.
func (r Virtual_Guest_Block_Device_Template_Group) SetTags(tags *string) (resp bool, err error) {
	params := []interface{}{
		tags,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Block_Device_Template_Group", "setTags", params, &r.Options, &resp)
	return
}

// no documentation yet
type Virtual_Guest_Boot_Parameter struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualGuestBootParameterService returns an instance of the Virtual_Guest_Boot_Parameter SoftLayer service
func GetVirtualGuestBootParameterService(sess session.SLSession) Virtual_Guest_Boot_Parameter {
	return Virtual_Guest_Boot_Parameter{Session: sess}
}

func (r Virtual_Guest_Boot_Parameter) Id(id int) Virtual_Guest_Boot_Parameter {
	r.Options.Id = &id
	return r
}

func (r Virtual_Guest_Boot_Parameter) Mask(mask string) Virtual_Guest_Boot_Parameter {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Guest_Boot_Parameter) Filter(filter string) Virtual_Guest_Boot_Parameter {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Guest_Boot_Parameter) Limit(limit int) Virtual_Guest_Boot_Parameter {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Guest_Boot_Parameter) Offset(offset int) Virtual_Guest_Boot_Parameter {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Virtual_Guest_Boot_Parameter) CreateObject(templateObject *datatypes.Virtual_Guest_Boot_Parameter) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter", "createObject", params, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Boot_Parameter) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter", "deleteObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Boot_Parameter) EditObject(templateObject *datatypes.Virtual_Guest_Boot_Parameter) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Boot_Parameter) GetGuest() (resp datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter", "getGuest", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Boot_Parameter) GetGuestBootParameterType() (resp datatypes.Virtual_Guest_Boot_Parameter_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter", "getGuestBootParameterType", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Boot_Parameter) GetObject() (resp datatypes.Virtual_Guest_Boot_Parameter, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter", "getObject", nil, &r.Options, &resp)
	return
}

// Describes a virtual guest boot parameter. In this the word class is used in the context of arguments sent to cloud computing instances such as single user mode and boot into bash.
type Virtual_Guest_Boot_Parameter_Type struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualGuestBootParameterTypeService returns an instance of the Virtual_Guest_Boot_Parameter_Type SoftLayer service
func GetVirtualGuestBootParameterTypeService(sess session.SLSession) Virtual_Guest_Boot_Parameter_Type {
	return Virtual_Guest_Boot_Parameter_Type{Session: sess}
}

func (r Virtual_Guest_Boot_Parameter_Type) Id(id int) Virtual_Guest_Boot_Parameter_Type {
	r.Options.Id = &id
	return r
}

func (r Virtual_Guest_Boot_Parameter_Type) Mask(mask string) Virtual_Guest_Boot_Parameter_Type {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Guest_Boot_Parameter_Type) Filter(filter string) Virtual_Guest_Boot_Parameter_Type {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Guest_Boot_Parameter_Type) Limit(limit int) Virtual_Guest_Boot_Parameter_Type {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Guest_Boot_Parameter_Type) Offset(offset int) Virtual_Guest_Boot_Parameter_Type {
	r.Options.Offset = &offset
	return r
}

// no documentation yet
func (r Virtual_Guest_Boot_Parameter_Type) GetAllObjects() (resp []datatypes.Virtual_Guest_Boot_Parameter_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter_Type", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Boot_Parameter_Type) GetObject() (resp datatypes.Virtual_Guest_Boot_Parameter_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Boot_Parameter_Type", "getObject", nil, &r.Options, &resp)
	return
}

// The virtual guest network component data type presents the structure in which all computing instance network components are presented. Internally, the structure supports various virtualization platforms with no change to external interaction.
//
// A guest, also known as a virtual server, represents an allocation of resources on a virtual host.
type Virtual_Guest_Network_Component struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualGuestNetworkComponentService returns an instance of the Virtual_Guest_Network_Component SoftLayer service
func GetVirtualGuestNetworkComponentService(sess session.SLSession) Virtual_Guest_Network_Component {
	return Virtual_Guest_Network_Component{Session: sess}
}

func (r Virtual_Guest_Network_Component) Id(id int) Virtual_Guest_Network_Component {
	r.Options.Id = &id
	return r
}

func (r Virtual_Guest_Network_Component) Mask(mask string) Virtual_Guest_Network_Component {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Guest_Network_Component) Filter(filter string) Virtual_Guest_Network_Component {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Guest_Network_Component) Limit(limit int) Virtual_Guest_Network_Component {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Guest_Network_Component) Offset(offset int) Virtual_Guest_Network_Component {
	r.Options.Offset = &offset
	return r
}

// Completely restrict all incoming and outgoing bandwidth traffic to a network component
func (r Virtual_Guest_Network_Component) Disable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "disable", nil, &r.Options, &resp)
	return
}

// Allow incoming and outgoing bandwidth traffic to a network component
func (r Virtual_Guest_Network_Component) Enable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "enable", nil, &r.Options, &resp)
	return
}

// Retrieve The computing instance that this network component exists on.
func (r Virtual_Guest_Network_Component) GetGuest() (resp datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getGuest", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Network_Component) GetHighAvailabilityFirewallFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getHighAvailabilityFirewallFlag", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Network_Component) GetIcpBinding() (resp datatypes.Virtual_Guest_Network_Component_IcpBinding, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getIcpBinding", nil, &r.Options, &resp)
	return
}

// Retrieve The records of all IP addresses bound to a computing instance's network component.
func (r Virtual_Guest_Network_Component) GetIpAddressBindings() (resp []datatypes.Virtual_Guest_Network_Component_IpAddress, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getIpAddressBindings", nil, &r.Options, &resp)
	return
}

// Retrieve The upstream network component firewall.
func (r Virtual_Guest_Network_Component) GetNetworkComponentFirewall() (resp datatypes.Network_Component_Firewall, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getNetworkComponentFirewall", nil, &r.Options, &resp)
	return
}

// Retrieve The VLAN that a computing instance network component's subnet is associated with.
func (r Virtual_Guest_Network_Component) GetNetworkVlan() (resp datatypes.Network_Vlan, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getNetworkVlan", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Network_Component) GetObject() (resp datatypes.Virtual_Guest_Network_Component, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve A computing instance network component's primary IP address.
func (r Virtual_Guest_Network_Component) GetPrimaryIpAddress() (resp string, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getPrimaryIpAddress", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Guest_Network_Component) GetPrimaryIpAddressRecord() (resp datatypes.Network_Subnet_IpAddress, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getPrimaryIpAddressRecord", nil, &r.Options, &resp)
	return
}

// Retrieve A network component's subnet for its primary IP address
func (r Virtual_Guest_Network_Component) GetPrimarySubnet() (resp datatypes.Network_Subnet, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getPrimarySubnet", nil, &r.Options, &resp)
	return
}

// Retrieve A network component's primary IPv6 IP address record.
func (r Virtual_Guest_Network_Component) GetPrimaryVersion6IpAddressRecord() (resp datatypes.Network_Subnet_IpAddress, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getPrimaryVersion6IpAddressRecord", nil, &r.Options, &resp)
	return
}

// Retrieve A network component's routers.
func (r Virtual_Guest_Network_Component) GetRouter() (resp datatypes.Hardware_Router, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getRouter", nil, &r.Options, &resp)
	return
}

// Retrieve The bindings associating security groups to this network component
func (r Virtual_Guest_Network_Component) GetSecurityGroupBindings() (resp []datatypes.Virtual_Network_SecurityGroup_NetworkComponentBinding, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getSecurityGroupBindings", nil, &r.Options, &resp)
	return
}

// Retrieve A network component's subnets. A subnet is a group of IP addresses
func (r Virtual_Guest_Network_Component) GetSubnets() (resp []datatypes.Network_Subnet, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "getSubnets", nil, &r.Options, &resp)
	return
}

// Issues a ping command and returns the success (true) or failure (false) of the ping command.
func (r Virtual_Guest_Network_Component) IsPingable() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "isPingable", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Guest_Network_Component) SecurityGroupsReady() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Guest_Network_Component", "securityGroupsReady", nil, &r.Options, &resp)
	return
}

// The virtual host represents the platform on which virtual guests reside. At times a virtual host has no allocations on the physical server, however with many modern platforms it is a virtual machine with small CPU and Memory allocations that runs in the Control Domain.
type Virtual_Host struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualHostService returns an instance of the Virtual_Host SoftLayer service
func GetVirtualHostService(sess session.SLSession) Virtual_Host {
	return Virtual_Host{Session: sess}
}

func (r Virtual_Host) Id(id int) Virtual_Host {
	r.Options.Id = &id
	return r
}

func (r Virtual_Host) Mask(mask string) Virtual_Host {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Host) Filter(filter string) Virtual_Host {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Host) Limit(limit int) Virtual_Host {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Host) Offset(offset int) Virtual_Host {
	r.Options.Offset = &offset
	return r
}

// Retrieve The account which a virtual host belongs to.
func (r Virtual_Host) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Host", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve The hardware record which a virtual host resides on.
func (r Virtual_Host) GetHardware() (resp datatypes.Hardware_Server, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Host", "getHardware", nil, &r.Options, &resp)
	return
}

// Retrieve The metric tracking object for this virtual host.
func (r Virtual_Host) GetMetricTrackingObject() (resp datatypes.Metric_Tracking_Object, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Host", "getMetricTrackingObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Host) GetObject() (resp datatypes.Virtual_Host, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Host", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Host) GetPciDevices() (resp []datatypes.Virtual_Host_PciDevice, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Host", "getPciDevices", nil, &r.Options, &resp)
	return
}

// This data type presents the structure for a virtual guest placement group. The data type contains relational properties to the virtual guest placement group rule class.
type Virtual_PlacementGroup struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualPlacementGroupService returns an instance of the Virtual_PlacementGroup SoftLayer service
func GetVirtualPlacementGroupService(sess session.SLSession) Virtual_PlacementGroup {
	return Virtual_PlacementGroup{Session: sess}
}

func (r Virtual_PlacementGroup) Id(id int) Virtual_PlacementGroup {
	r.Options.Id = &id
	return r
}

func (r Virtual_PlacementGroup) Mask(mask string) Virtual_PlacementGroup {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_PlacementGroup) Filter(filter string) Virtual_PlacementGroup {
	r.Options.Filter = filter
	return r
}

func (r Virtual_PlacementGroup) Limit(limit int) Virtual_PlacementGroup {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_PlacementGroup) Offset(offset int) Virtual_PlacementGroup {
	r.Options.Offset = &offset
	return r
}

// Add a placement group to your account for use during VSI provisioning.
func (r Virtual_PlacementGroup) CreateObject(templateObject *datatypes.Virtual_PlacementGroup) (resp datatypes.Virtual_PlacementGroup, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "createObject", params, &r.Options, &resp)
	return
}

// Delete a placement group from your account.
func (r Virtual_PlacementGroup) DeleteObject() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "deleteObject", nil, &r.Options, &resp)
	return
}

// Update a placement group.
func (r Virtual_PlacementGroup) EditObject(templateObject *datatypes.Virtual_PlacementGroup) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve The account that the placement group is implemented on.
func (r Virtual_PlacementGroup) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "getAccount", nil, &r.Options, &resp)
	return
}

// Returns all routers available for use with placement groups. If a datacenter location ID is provided, this method will further restrict the list of routers to ones contained within that datacenter.
func (r Virtual_PlacementGroup) GetAvailableRouters(datacenterId *int) (resp []datatypes.Hardware, err error) {
	params := []interface{}{
		datacenterId,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "getAvailableRouters", params, &r.Options, &resp)
	return
}

// Retrieve The router the placement group is implemented on.
func (r Virtual_PlacementGroup) GetBackendRouter() (resp datatypes.Hardware_Router_Backend, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "getBackendRouter", nil, &r.Options, &resp)
	return
}

// Retrieve The virtual guests that are members of the placement group.
func (r Virtual_PlacementGroup) GetGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "getGuests", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_PlacementGroup) GetObject() (resp datatypes.Virtual_PlacementGroup, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The placement rule that the placement group is implementing.
func (r Virtual_PlacementGroup) GetRule() (resp datatypes.Virtual_PlacementGroup_Rule, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup", "getRule", nil, &r.Options, &resp)
	return
}

// This data type presents the structure of a virtual guest placement group rule.
type Virtual_PlacementGroup_Rule struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualPlacementGroupRuleService returns an instance of the Virtual_PlacementGroup_Rule SoftLayer service
func GetVirtualPlacementGroupRuleService(sess session.SLSession) Virtual_PlacementGroup_Rule {
	return Virtual_PlacementGroup_Rule{Session: sess}
}

func (r Virtual_PlacementGroup_Rule) Id(id int) Virtual_PlacementGroup_Rule {
	r.Options.Id = &id
	return r
}

func (r Virtual_PlacementGroup_Rule) Mask(mask string) Virtual_PlacementGroup_Rule {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_PlacementGroup_Rule) Filter(filter string) Virtual_PlacementGroup_Rule {
	r.Options.Filter = filter
	return r
}

func (r Virtual_PlacementGroup_Rule) Limit(limit int) Virtual_PlacementGroup_Rule {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_PlacementGroup_Rule) Offset(offset int) Virtual_PlacementGroup_Rule {
	r.Options.Offset = &offset
	return r
}

// Get all placement group rules.
func (r Virtual_PlacementGroup_Rule) GetAllObjects() (resp []datatypes.Virtual_PlacementGroup_Rule, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup_Rule", "getAllObjects", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_PlacementGroup_Rule) GetObject() (resp datatypes.Virtual_PlacementGroup_Rule, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_PlacementGroup_Rule", "getObject", nil, &r.Options, &resp)
	return
}

// This data type presents the structure for a virtual reserved capacity group.
type Virtual_ReservedCapacityGroup struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualReservedCapacityGroupService returns an instance of the Virtual_ReservedCapacityGroup SoftLayer service
func GetVirtualReservedCapacityGroupService(sess session.SLSession) Virtual_ReservedCapacityGroup {
	return Virtual_ReservedCapacityGroup{Session: sess}
}

func (r Virtual_ReservedCapacityGroup) Id(id int) Virtual_ReservedCapacityGroup {
	r.Options.Id = &id
	return r
}

func (r Virtual_ReservedCapacityGroup) Mask(mask string) Virtual_ReservedCapacityGroup {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_ReservedCapacityGroup) Filter(filter string) Virtual_ReservedCapacityGroup {
	r.Options.Filter = filter
	return r
}

func (r Virtual_ReservedCapacityGroup) Limit(limit int) Virtual_ReservedCapacityGroup {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_ReservedCapacityGroup) Offset(offset int) Virtual_ReservedCapacityGroup {
	r.Options.Offset = &offset
	return r
}

// Update a reserved capacity group.
func (r Virtual_ReservedCapacityGroup) EditObject(templateObject *datatypes.Virtual_ReservedCapacityGroup) (resp bool, err error) {
	params := []interface{}{
		templateObject,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "editObject", params, &r.Options, &resp)
	return
}

// Retrieve The account that the reserved capacity group is implemented on.
func (r Virtual_ReservedCapacityGroup) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getAccount", nil, &r.Options, &resp)
	return
}

// Retrieve The instances available for guest provisions on this reserved capacity group.
func (r Virtual_ReservedCapacityGroup) GetAvailableInstances() (resp []datatypes.Virtual_ReservedCapacityGroup_Instance, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getAvailableInstances", nil, &r.Options, &resp)
	return
}

// Retrieve The router the reserved capacity group is implemented on.
func (r Virtual_ReservedCapacityGroup) GetBackendRouter() (resp datatypes.Hardware_Router_Backend, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getBackendRouter", nil, &r.Options, &resp)
	return
}

// Retrieve The guest instances that are members of this reserved capacity group.
func (r Virtual_ReservedCapacityGroup) GetInstances() (resp []datatypes.Virtual_ReservedCapacityGroup_Instance, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getInstances", nil, &r.Options, &resp)
	return
}

// Retrieve The number of instances that are members of this reserved capacity group.
func (r Virtual_ReservedCapacityGroup) GetInstancesCount() (resp uint, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getInstancesCount", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_ReservedCapacityGroup) GetObject() (resp datatypes.Virtual_ReservedCapacityGroup, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The instances already occupied by a guest on this reserved capacity group.
func (r Virtual_ReservedCapacityGroup) GetOccupiedInstances() (resp []datatypes.Virtual_ReservedCapacityGroup_Instance, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup", "getOccupiedInstances", nil, &r.Options, &resp)
	return
}

// This data type presents the structure for a virtual reserved capacity group instance.
type Virtual_ReservedCapacityGroup_Instance struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualReservedCapacityGroupInstanceService returns an instance of the Virtual_ReservedCapacityGroup_Instance SoftLayer service
func GetVirtualReservedCapacityGroupInstanceService(sess session.SLSession) Virtual_ReservedCapacityGroup_Instance {
	return Virtual_ReservedCapacityGroup_Instance{Session: sess}
}

func (r Virtual_ReservedCapacityGroup_Instance) Id(id int) Virtual_ReservedCapacityGroup_Instance {
	r.Options.Id = &id
	return r
}

func (r Virtual_ReservedCapacityGroup_Instance) Mask(mask string) Virtual_ReservedCapacityGroup_Instance {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_ReservedCapacityGroup_Instance) Filter(filter string) Virtual_ReservedCapacityGroup_Instance {
	r.Options.Filter = filter
	return r
}

func (r Virtual_ReservedCapacityGroup_Instance) Limit(limit int) Virtual_ReservedCapacityGroup_Instance {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_ReservedCapacityGroup_Instance) Offset(offset int) Virtual_ReservedCapacityGroup_Instance {
	r.Options.Offset = &offset
	return r
}

// Retrieve Flag to indecate whether or not the reserved instance is available or not.
func (r Virtual_ReservedCapacityGroup_Instance) GetAvailableFlag() (resp bool, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup_Instance", "getAvailableFlag", nil, &r.Options, &resp)
	return
}

// Retrieve The billing item for the reserved capacity group instance.
func (r Virtual_ReservedCapacityGroup_Instance) GetBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup_Instance", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve The virtual guest associated with this reserved capacity group instance.
func (r Virtual_ReservedCapacityGroup_Instance) GetGuest() (resp datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup_Instance", "getGuest", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_ReservedCapacityGroup_Instance) GetObject() (resp datatypes.Virtual_ReservedCapacityGroup_Instance, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup_Instance", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The reserved instances that are members of this reserved capacity group.
func (r Virtual_ReservedCapacityGroup_Instance) GetReservedCapacityGroup() (resp datatypes.Virtual_ReservedCapacityGroup, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_ReservedCapacityGroup_Instance", "getReservedCapacityGroup", nil, &r.Options, &resp)
	return
}

// The SoftLayer_Virtual_Storage_Repository represents a web based storage system that can be accessed through many types of devices, interfaces, and other resources.
type Virtual_Storage_Repository struct {
	Session session.SLSession
	Options sl.Options
}

// GetVirtualStorageRepositoryService returns an instance of the Virtual_Storage_Repository SoftLayer service
func GetVirtualStorageRepositoryService(sess session.SLSession) Virtual_Storage_Repository {
	return Virtual_Storage_Repository{Session: sess}
}

func (r Virtual_Storage_Repository) Id(id int) Virtual_Storage_Repository {
	r.Options.Id = &id
	return r
}

func (r Virtual_Storage_Repository) Mask(mask string) Virtual_Storage_Repository {
	if !strings.HasPrefix(mask, "mask[") && (strings.Contains(mask, "[") || strings.Contains(mask, ",")) {
		mask = fmt.Sprintf("mask[%s]", mask)
	}

	r.Options.Mask = mask
	return r
}

func (r Virtual_Storage_Repository) Filter(filter string) Virtual_Storage_Repository {
	r.Options.Filter = filter
	return r
}

func (r Virtual_Storage_Repository) Limit(limit int) Virtual_Storage_Repository {
	r.Options.Limit = &limit
	return r
}

func (r Virtual_Storage_Repository) Offset(offset int) Virtual_Storage_Repository {
	r.Options.Offset = &offset
	return r
}

// Retrieve The [[SoftLayer_Account|account]] that a storage repository belongs to.
func (r Virtual_Storage_Repository) GetAccount() (resp datatypes.Account, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getAccount", nil, &r.Options, &resp)
	return
}

// Returns the archive storage disk usage fee rate per gigabyte.
func (r Virtual_Storage_Repository) GetArchiveDiskUsageRatePerGb() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getArchiveDiskUsageRatePerGb", nil, &r.Options, &resp)
	return
}

// Returns the average disk space usage for a storage repository.
func (r Virtual_Storage_Repository) GetAverageDiskUsageMetricDataFromInfluxByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time) (resp datatypes.Float64, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getAverageDiskUsageMetricDataFromInfluxByDate", params, &r.Options, &resp)
	return
}

// Returns the average disk space usage for a storage repository.
func (r Virtual_Storage_Repository) GetAverageUsageMetricDataByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time) (resp datatypes.Float64, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getAverageUsageMetricDataByDate", params, &r.Options, &resp)
	return
}

// Retrieve The current billing item for a storage repository.
func (r Virtual_Storage_Repository) GetBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getBillingItem", nil, &r.Options, &resp)
	return
}

// Retrieve The datacenter that a virtual storage repository resides in.
func (r Virtual_Storage_Repository) GetDatacenter() (resp datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getDatacenter", nil, &r.Options, &resp)
	return
}

// Retrieve The [[SoftLayer_Virtual_Disk_Image|disk images]] that are in a storage repository. Disk images are the virtual hard drives for a virtual guest.
func (r Virtual_Storage_Repository) GetDiskImages() (resp []datatypes.Virtual_Disk_Image, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getDiskImages", nil, &r.Options, &resp)
	return
}

// Retrieve The computing instances that have disk images in a storage repository.
func (r Virtual_Storage_Repository) GetGuests() (resp []datatypes.Virtual_Guest, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getGuests", nil, &r.Options, &resp)
	return
}

// Retrieve
func (r Virtual_Storage_Repository) GetMetricTrackingObject() (resp datatypes.Metric_Tracking_Object_Virtual_Storage_Repository, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getMetricTrackingObject", nil, &r.Options, &resp)
	return
}

// no documentation yet
func (r Virtual_Storage_Repository) GetObject() (resp datatypes.Virtual_Storage_Repository, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getObject", nil, &r.Options, &resp)
	return
}

// Retrieve The current billing item for a public storage repository.
func (r Virtual_Storage_Repository) GetPublicImageBillingItem() (resp datatypes.Billing_Item, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getPublicImageBillingItem", nil, &r.Options, &resp)
	return
}

// Returns the public image storage disk usage fee rate per gigabyte.
func (r Virtual_Storage_Repository) GetPublicImageDiskUsageRatePerGb() (resp datatypes.Float64, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getPublicImageDiskUsageRatePerGb", nil, &r.Options, &resp)
	return
}

// Returns the public image storage locations.
func (r Virtual_Storage_Repository) GetStorageLocations() (resp []datatypes.Location, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getStorageLocations", nil, &r.Options, &resp)
	return
}

// Retrieve A storage repository's [[SoftLayer_Virtual_Storage_Repository_Type|type]].
func (r Virtual_Storage_Repository) GetType() (resp datatypes.Virtual_Storage_Repository_Type, err error) {
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getType", nil, &r.Options, &resp)
	return
}

// Retrieve disk usage data on a [[SoftLayer_Virtual_Guest|Cloud Computing Instance]] image for the time range you provide.  Each data entry objects contain ”dateTime” and ”counter” properties. ”dateTime” property indicates the time that the disk usage data was measured and ”counter” property holds the disk usage in bytes.
func (r Virtual_Storage_Repository) GetUsageMetricDataByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time) (resp []datatypes.Metric_Tracking_Object_Data, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getUsageMetricDataByDate", params, &r.Options, &resp)
	return
}

// Returns a disk usage image based on disk usage specified by the input parameters.
func (r Virtual_Storage_Repository) GetUsageMetricImageByDate(startDateTime *datatypes.Time, endDateTime *datatypes.Time) (resp datatypes.Container_Bandwidth_GraphOutputs, err error) {
	params := []interface{}{
		startDateTime,
		endDateTime,
	}
	err = r.Session.DoRequest("SoftLayer_Virtual_Storage_Repository", "getUsageMetricImageByDate", params, &r.Options, &resp)
	return
}
