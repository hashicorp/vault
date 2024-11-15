/*
Copyright (c) 2024-2024 VMware, Inc. All Rights Reserved.

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

import (
	"fmt"
)

// EnsureDisksHaveControllers ensures that all disks in the provided
// ConfigSpec point to a controller. If no controller exists, LSILogic SCSI
// controllers are added to the ConfigSpec as necessary for the disks.
//
// Please note the following table for the number of controllers of each type
// that are supported as well as how many disks (per controller) each supports:
//
// SATA
//   - controllers                                    4
//   - disks                                         30
//
// SCSI
//   - controllers                                    4
//   - disks (non-paravirtual)                       16
//   - disks (paravirtual, hardware version <14)     16
//   - disks (paravirtual, hardware version >=14)   256
//
// NVME
//   - controllers                                    4
//   - disks (hardware version <20)                  15
//   - disks (hardware version >=21)                255
func (cs *VirtualMachineConfigSpec) EnsureDisksHaveControllers(
	existingDevices ...BaseVirtualDevice) error {

	if cs == nil {
		panic("configSpec is nil")
	}

	var (
		disks           []*VirtualDisk
		newDeviceKey    int32
		pciController   *VirtualPCIController
		diskControllers = ensureDiskControllerData{
			controllerKeys:                map[int32]BaseVirtualController{},
			controllerKeysToAttachedDisks: map[int32]int{},
		}
	)

	// Inspect the ConfigSpec
	for i := range cs.DeviceChange {
		var (
			bdc BaseVirtualDeviceConfigSpec
			bvd BaseVirtualDevice
			dc  *VirtualDeviceConfigSpec
			d   *VirtualDevice
		)

		if bdc = cs.DeviceChange[i]; bdc == nil {
			continue
		}

		if dc = bdc.GetVirtualDeviceConfigSpec(); dc == nil {
			continue
		}

		if dc.Operation == VirtualDeviceConfigSpecOperationRemove {
			// Do not consider devices being removed.
			continue
		}

		bvd = dc.Device
		if bvd == nil {
			continue
		}

		if d = bvd.GetVirtualDevice(); d == nil {
			continue
		}

		switch tvd := bvd.(type) {
		case *VirtualPCIController:
			pciController = tvd

		case
			// SCSI
			*ParaVirtualSCSIController,
			*VirtualBusLogicController,
			*VirtualLsiLogicController,
			*VirtualLsiLogicSASController,
			*VirtualSCSIController,

			// SATA
			*VirtualSATAController,
			*VirtualAHCIController,

			// NVME
			*VirtualNVMEController:

			diskControllers.add(bvd)

		case *VirtualDisk:

			disks = append(disks, tvd)

			if controllerKey := d.ControllerKey; controllerKey != 0 {
				// If the disk points to a controller key, then increment
				// the number of devices attached to that controller.
				//
				// Please note that at this point it is not yet known if the
				// controller key is a *valid* controller.
				diskControllers.attach(controllerKey)
			}
		}

		// Keep track of the smallest device key used. Please note, because
		// device keys in a ConfigSpec are negative numbers, -200 going to be
		// smaller than -1.
		if d.Key < newDeviceKey {
			newDeviceKey = d.Key
		}
	}

	if len(disks) == 0 {
		// If there are no disks, then go ahead and return early.
		return nil
	}

	// Categorize any controllers that already exist.
	for i := range existingDevices {
		var (
			d   *VirtualDevice
			bvd = existingDevices[i]
		)

		if bvd == nil {
			continue
		}

		if d = bvd.GetVirtualDevice(); d == nil {
			continue
		}

		switch tvd := bvd.(type) {
		case *VirtualPCIController:
			pciController = tvd
		case
			// SCSI
			*ParaVirtualSCSIController,
			*VirtualBusLogicController,
			*VirtualLsiLogicController,
			*VirtualLsiLogicSASController,
			*VirtualSCSIController,

			// SATA
			*VirtualSATAController,
			*VirtualAHCIController,

			// NVME
			*VirtualNVMEController:

			diskControllers.add(bvd)

		case *VirtualDisk:
			diskControllers.attach(tvd.ControllerKey)
		}
	}

	// Decrement the newDeviceKey so the next device has a unique key.
	newDeviceKey--

	if pciController == nil {
		// Add a PCI controller if one is not present.
		pciController = &VirtualPCIController{
			VirtualController: VirtualController{
				VirtualDevice: VirtualDevice{
					Key: newDeviceKey,
				},
			},
		}

		// Decrement the newDeviceKey so the next device has a unique key.
		newDeviceKey--

		// Add the new PCI controller to the ConfigSpec.
		cs.DeviceChange = append(
			cs.DeviceChange,
			&VirtualDeviceConfigSpec{
				Operation: VirtualDeviceConfigSpecOperationAdd,
				Device:    pciController,
			})
	}

	// Ensure all the recorded controller keys that point to disks are actually
	// valid controller keys.
	diskControllers.validateAttachments()

	for i := range disks {
		disk := disks[i]

		// If the disk already points to a controller then skip to the next
		// disk.
		if diskControllers.exists(disk.ControllerKey) {
			continue
		}

		// The disk does not point to a controller, so try to locate one.
		if ensureDiskControllerFind(disk, &diskControllers) {
			// A controller was located for the disk, so go ahead and skip to
			// the next disk.
			continue
		}

		// No controller was located for the disk, so a controller must be
		// created.
		if err := ensureDiskControllerCreate(
			cs,
			pciController,
			newDeviceKey,
			&diskControllers); err != nil {

			return err
		}

		// Point the disk to the new controller.
		disk.ControllerKey = newDeviceKey

		// Add the controller key to the map that tracks how many disks are
		// attached to a given controller.
		diskControllers.attach(newDeviceKey)

		// Decrement the newDeviceKey so the next device has a unique key.
		newDeviceKey--
	}

	return nil
}

const (
	maxSCSIControllers                     = 4
	maxSATAControllers                     = 4
	maxNVMEControllers                     = 4
	maxDisksPerSCSIController              = 16
	maxDisksPerPVSCSIControllerHWVersion14 = 256 // TODO(akutz)
	maxDisksPerSATAController              = 30
	maxDisksPerNVMEController              = 15
	maxDisksPerNVMEControllerHWVersion21   = 255 // TODO(akutz)
)

type ensureDiskControllerBusNumbers struct {
	zero bool
	one  bool
	two  bool
}

func (d ensureDiskControllerBusNumbers) free() int32 {
	switch {
	case !d.zero:
		return 0
	case !d.one:
		return 1
	case !d.two:
		return 2
	default:
		return 3
	}
}

func (d *ensureDiskControllerBusNumbers) set(busNumber int32) {
	switch busNumber {
	case 0:
		d.zero = true
	case 1:
		d.one = true
	case 2:
		d.two = true
	}
}

type ensureDiskControllerData struct {
	// TODO(akutz) Use the hardware version when calculating the max disks for
	//             a given controller type.
	// hardwareVersion int

	controllerKeys                map[int32]BaseVirtualController
	controllerKeysToAttachedDisks map[int32]int

	// SCSI
	scsiBusNumbers             ensureDiskControllerBusNumbers
	pvSCSIControllerKeys       []int32
	busLogicSCSIControllerKeys []int32
	lsiLogicControllerKeys     []int32
	lsiLogicSASControllerKeys  []int32
	scsiControllerKeys         []int32

	// SATA
	sataBusNumbers     ensureDiskControllerBusNumbers
	sataControllerKeys []int32
	ahciControllerKeys []int32

	// NVME
	nvmeBusNumbers     ensureDiskControllerBusNumbers
	nvmeControllerKeys []int32
}

func (d ensureDiskControllerData) numSCSIControllers() int {
	return len(d.pvSCSIControllerKeys) +
		len(d.busLogicSCSIControllerKeys) +
		len(d.lsiLogicControllerKeys) +
		len(d.lsiLogicSASControllerKeys) +
		len(d.scsiControllerKeys)
}

func (d ensureDiskControllerData) numSATAControllers() int {
	return len(d.sataControllerKeys) + len(d.ahciControllerKeys)
}

func (d ensureDiskControllerData) numNVMEControllers() int {
	return len(d.nvmeControllerKeys)
}

// validateAttachments ensures the attach numbers are correct by removing any
// keys from controllerKeysToAttachedDisks that do not also exist in
// controllerKeys.
func (d ensureDiskControllerData) validateAttachments() {
	// Remove any invalid controllers from controllerKeyToNumDiskMap.
	for key := range d.controllerKeysToAttachedDisks {
		if _, ok := d.controllerKeys[key]; !ok {
			delete(d.controllerKeysToAttachedDisks, key)
		}
	}
}

// exists returns true if a controller with the provided key exists.
func (d ensureDiskControllerData) exists(key int32) bool {
	return d.controllerKeys[key] != nil
}

// add records the provided controller in the map that relates keys to
// controllers as well as appends the key to the list of controllers of that
// given type.
func (d *ensureDiskControllerData) add(controller BaseVirtualDevice) {

	// Get the controller's device key.
	bvc := controller.(BaseVirtualController)
	key := bvc.GetVirtualController().Key
	busNumber := bvc.GetVirtualController().BusNumber

	// Record the controller's device key in the controller key map.
	d.controllerKeys[key] = bvc

	// Record the controller's device key in the list for that type of
	// controller.
	switch controller.(type) {

	// SCSI
	case *ParaVirtualSCSIController:
		d.pvSCSIControllerKeys = append(d.pvSCSIControllerKeys, key)
		d.scsiBusNumbers.set(busNumber)
	case *VirtualBusLogicController:
		d.busLogicSCSIControllerKeys = append(d.busLogicSCSIControllerKeys, key)
		d.scsiBusNumbers.set(busNumber)
	case *VirtualLsiLogicController:
		d.lsiLogicControllerKeys = append(d.lsiLogicControllerKeys, key)
		d.scsiBusNumbers.set(busNumber)
	case *VirtualLsiLogicSASController:
		d.lsiLogicSASControllerKeys = append(d.lsiLogicSASControllerKeys, key)
		d.scsiBusNumbers.set(busNumber)
	case *VirtualSCSIController:
		d.scsiControllerKeys = append(d.scsiControllerKeys, key)
		d.scsiBusNumbers.set(busNumber)

	// SATA
	case *VirtualSATAController:
		d.sataControllerKeys = append(d.sataControllerKeys, key)
		d.sataBusNumbers.set(busNumber)
	case *VirtualAHCIController:
		d.ahciControllerKeys = append(d.ahciControllerKeys, key)
		d.sataBusNumbers.set(busNumber)

	// NVME
	case *VirtualNVMEController:
		d.nvmeControllerKeys = append(d.nvmeControllerKeys, key)
		d.nvmeBusNumbers.set(busNumber)
	}
}

// attach increments the number of disks attached to the controller identified
// by the provided controller key.
func (d *ensureDiskControllerData) attach(controllerKey int32) {
	d.controllerKeysToAttachedDisks[controllerKey]++
}

// hasFreeSlot returns whether or not the controller identified by the provided
// controller key has a free slot to attach a disk.
//
// TODO(akutz) Consider the hardware version when calculating these values.
func (d *ensureDiskControllerData) hasFreeSlot(controllerKey int32) bool {

	var maxDisksForType int

	switch d.controllerKeys[controllerKey].(type) {
	case
		// SCSI (paravirtual)
		*ParaVirtualSCSIController:

		maxDisksForType = maxDisksPerSCSIController

	case
		// SCSI (non-paravirtual)
		*VirtualBusLogicController,
		*VirtualLsiLogicController,
		*VirtualLsiLogicSASController,
		*VirtualSCSIController:

		maxDisksForType = maxDisksPerSCSIController

	case
		// SATA
		*VirtualSATAController,
		*VirtualAHCIController:

		maxDisksForType = maxDisksPerSATAController

	case
		// NVME
		*VirtualNVMEController:

		maxDisksForType = maxDisksPerNVMEController
	}

	return d.controllerKeysToAttachedDisks[controllerKey] < maxDisksForType-1
}

// ensureDiskControllerFind attempts to locate a controller for the provided
// disk.
//
// Please note this function is written to preserve the order in which
// controllers are located by preferring controller types in the order in which
// they are listed in this function. This prevents the following situation:
//
//   - A ConfigSpec has three controllers in the following order: PVSCSI-1,
//     NVME-1, and PVSCSI-2.
//   - The controller PVSCSI-1 is full while NVME-1 and PVSCSI-2 have free
//     slots.
//   - The *desired* behavior is to look at all, possible PVSCSI controllers
//     before moving onto SATA and then finally NVME controllers.
//   - If the function iterated over the device list in list-order, then the
//     NVME-1 controller would be located first.
//   - Instead, this function iterates over each *type* of controller first
//     before moving onto the next type.
//   - This means that even though NVME-1 has free slots, PVSCSI-2 is checked
//     first.
//
// The order of preference is as follows:
//
// * SCSI
//   - ParaVirtualSCSIController
//   - VirtualBusLogicController
//   - VirtualLsiLogicController
//   - VirtualLsiLogicSASController
//   - VirtualSCSIController
//
// * SATA
//   - VirtualSATAController
//   - VirtualAHCIController
//
// * NVME
//   - VirtualNVMEController
func ensureDiskControllerFind(
	disk *VirtualDisk,
	diskControllers *ensureDiskControllerData) bool {

	return false ||
		// SCSI
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.pvSCSIControllerKeys) ||
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.busLogicSCSIControllerKeys) ||
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.lsiLogicControllerKeys) ||
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.lsiLogicSASControllerKeys) ||
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.scsiControllerKeys) ||

		// SATA
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.sataControllerKeys) ||
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.ahciControllerKeys) ||

		// NVME
		ensureDiskControllerFindWith(
			disk,
			diskControllers,
			diskControllers.nvmeControllerKeys)
}

func ensureDiskControllerFindWith(
	disk *VirtualDisk,
	diskControllers *ensureDiskControllerData,
	controllerKeys []int32) bool {

	for i := range controllerKeys {
		controllerKey := controllerKeys[i]
		if diskControllers.hasFreeSlot(controllerKey) {
			// If the controller has room for another disk, then use this
			// controller for the current disk.
			disk.ControllerKey = controllerKey
			diskControllers.attach(controllerKey)
			return true
		}
	}
	return false
}

func ensureDiskControllerCreate(
	configSpec *VirtualMachineConfigSpec,
	pciController *VirtualPCIController,
	newDeviceKey int32,
	diskControllers *ensureDiskControllerData) error {

	var controller BaseVirtualDevice
	switch {
	case diskControllers.numSCSIControllers() < maxSCSIControllers:
		// Prefer creating a new SCSI controller.
		controller = &ParaVirtualSCSIController{
			VirtualSCSIController: VirtualSCSIController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						ControllerKey: pciController.Key,
						Key:           newDeviceKey,
					},
					BusNumber: diskControllers.scsiBusNumbers.free(),
				},
				HotAddRemove: NewBool(true),
				SharedBus:    VirtualSCSISharingNoSharing,
			},
		}
	case diskControllers.numSATAControllers() < maxSATAControllers:
		// If there are no more SCSI controllers, create a SATA
		// controller.
		controller = &VirtualAHCIController{
			VirtualSATAController: VirtualSATAController{
				VirtualController: VirtualController{
					VirtualDevice: VirtualDevice{
						ControllerKey: pciController.Key,
						Key:           newDeviceKey,
					},
					BusNumber: diskControllers.sataBusNumbers.free(),
				},
			},
		}
	case diskControllers.numNVMEControllers() < maxNVMEControllers:
		// If there are no more SATA controllers, create an NVME
		// controller.
		controller = &VirtualNVMEController{
			VirtualController: VirtualController{
				VirtualDevice: VirtualDevice{
					ControllerKey: pciController.Key,
					Key:           newDeviceKey,
				},
				BusNumber: diskControllers.nvmeBusNumbers.free(),
			},
			SharedBus: string(VirtualNVMEControllerSharingNoSharing),
		}
	default:
		return fmt.Errorf("no controllers available")
	}

	// Add the new controller to the ConfigSpec.
	configSpec.DeviceChange = append(
		configSpec.DeviceChange,
		&VirtualDeviceConfigSpec{
			Operation: VirtualDeviceConfigSpecOperationAdd,
			Device:    controller,
		})

	// Record the new controller.
	diskControllers.add(controller)

	return nil
}
