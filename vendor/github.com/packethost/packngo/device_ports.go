package packngo

import (
	"fmt"
	"strings"
)

const portBasePath = "/ports"

// DevicePortService handles operations on a port which belongs to a particular device
//
// Deprecated: use PortService or Device methods
type DevicePortService interface {
	Assign(*PortAssignRequest) (*Port, *Response, error)
	Unassign(*PortAssignRequest) (*Port, *Response, error)
	AssignNative(*PortAssignRequest) (*Port, *Response, error)
	UnassignNative(string) (*Port, *Response, error)
	Bond(*Port, bool) (*Port, *Response, error)
	Disbond(*Port, bool) (*Port, *Response, error)
	DeviceToNetworkType(string, string) (*Device, error)
	DeviceNetworkType(string) (string, error)
	PortToLayerTwo(string, string) (*Port, *Response, error)
	PortToLayerThree(string, string) (*Port, *Response, error)
	GetPortByName(string, string) (*Port, error)
	GetOddEthPorts(*Device) (map[string]*Port, error)
	GetAllEthPorts(*Device) (map[string]*Port, error)
	ConvertDevice(*Device, string) error
}

// DevicePortServiceOp implements DevicePortService on the Equinix Metal API
//
// Deprecated: use PortServiceOp or Device methods
type DevicePortServiceOp struct {
	client *Client
}

// GetPortByName returns the matching Port on the specified device
//
// Deprecated: use Device.GetPortByName
func (i *DevicePortServiceOp) GetPortByName(deviceID, name string) (*Port, error) {
	device, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}
	return device.GetPortByName(name)
}

// Assign the specified VLAN to the specified Port
//
// Deprecated: use PortServiceOp.Assign
func (i *DevicePortServiceOp) Assign(par *PortAssignRequest) (*Port, *Response, error) {
	return i.client.Ports.Assign(par.PortID, par.VirtualNetworkID)
}

// AssignNative designates the specified VLAN as the native VLAN for the
// specified Port
//
// Deprecated: use PortServiceOp.AssignNative
func (i *DevicePortServiceOp) AssignNative(par *PortAssignRequest) (*Port, *Response, error) {
	return i.client.Ports.AssignNative(par.PortID, par.VirtualNetworkID)
}

// UnassignNative removes the native VLAN from the specified Port
//
// Deprecated: use PortServiceOp.UnassignNative
func (i *DevicePortServiceOp) UnassignNative(portID string) (*Port, *Response, error) {
	if validateErr := ValidateUUID(portID); validateErr != nil {
		return nil, nil, validateErr
	}
	return i.client.Ports.UnassignNative(portID)
}

// Unassign removes the specified VLAN from the specified Port
//
// Deprecated: use PortServiceOp.Unassign
func (i *DevicePortServiceOp) Unassign(par *PortAssignRequest) (*Port, *Response, error) {
	return i.client.Ports.Unassign(par.PortID, par.VirtualNetworkID)
}

// Bond enabled bonding on the specified port
//
// Deprecated: use PortServiceOp.Bond
func (i *DevicePortServiceOp) Bond(p *Port, bulk_enable bool) (*Port, *Response, error) {
	if p.Data.Bonded {
		return p, nil, nil
	}

	return i.client.Ports.Bond(p.ID, bulk_enable)
}

// Disbond disables bonding on the specified port
//
// Deprecated: use PortServiceOp.Disbond
func (i *DevicePortServiceOp) Disbond(p *Port, bulk_disable bool) (*Port, *Response, error) {
	if !p.Data.Bonded {
		return p, nil, nil
	}
	return i.client.Ports.Disbond(p.ID, bulk_disable)
}

// PortToLayerTwo fetches the specified device, finds the matching port by name,
// and converts it to layer2. A port may already be in a layer2 mode, in which
// case the port will be returned with a nil response and nil error with no
// additional action taking place.
//
// Deprecated: use PortServiceOp.ConvertToLayerTwo
func (i *DevicePortServiceOp) PortToLayerTwo(deviceID, portName string) (*Port, *Response, error) {
	p, err := i.GetPortByName(deviceID, portName)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasPrefix(p.NetworkType, "layer2") {
		return p, nil, nil
	}

	return i.client.Ports.ConvertToLayerTwo(p.ID)
}

// PortToLayerThree fetches the specified device, finds the matching port by
// name, and converts it to layer3. A port may already be in a layer3 mode, in
// which case the port will be returned with a nil response and nil error with
// no additional action taking place.
//
// When switching to Layer3, a new set of IP addresses will be requested
// including Public IPv4, Public IPv6, and Private IPv6 addresses.
//
// Deprecated: use PortServiceOp.ConvertToLayerTwo
func (i *DevicePortServiceOp) PortToLayerThree(deviceID, portName string) (*Port, *Response, error) {
	p, err := i.GetPortByName(deviceID, portName)
	if err != nil {
		return nil, nil, err
	}
	if (p.NetworkType == NetworkTypeL3) || (p.NetworkType == NetworkTypeHybrid) {
		return p, nil, nil
	}

	ips := []AddressRequest{
		{AddressFamily: 4, Public: true},
		{AddressFamily: 4, Public: false},
		{AddressFamily: 6, Public: true},
	}

	return i.client.Ports.ConvertToLayerThree(p.ID, ips)
}

// DeviceNetworkType fetches the specified Device and returns a heuristic single
// word network type consistent with the Equinix Metal console experience.
//
// Deprecated: use Device.GetNetworkType
func (i *DevicePortServiceOp) DeviceNetworkType(deviceID string) (string, error) {
	if validateErr := ValidateUUID(deviceID); validateErr != nil {
		return "", validateErr
	}
	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return "", err
	}
	return d.GetNetworkType(), nil
}

// GetAllEthPorts fetches the specified Device and returns a heuristic single
// word network type consistent with the Equinix Metal console experience.
//
// Deprecated: use Device.GetPhysicalPorts
func (i *DevicePortServiceOp) GetAllEthPorts(d *Device) (map[string]*Port, error) {
	d, _, err := i.client.Devices.Get(d.ID, nil)
	if err != nil {
		return nil, err
	}
	return d.GetPhysicalPorts(), nil
}

// GetOddEthPorts fetches the specified Device and returns physical
// ports eth1 and eth3.
//
// Deprecated: use Device.GetPhysicalPorts and filter the map to only the keys
// ending with odd digits
func (i *DevicePortServiceOp) GetOddEthPorts(d *Device) (map[string]*Port, error) {
	d, _, err := i.client.Devices.Get(d.ID, nil)
	if err != nil {
		return nil, err
	}
	ret := map[string]*Port{}
	eth1, err := d.GetPortByName("eth1")
	if err != nil {
		return nil, err
	}
	ret["eth1"] = eth1

	eth3, err := d.GetPortByName("eth3")
	if err != nil {
		return ret, nil
	}
	ret["eth3"] = eth3
	return ret, nil

}

// ConvertDevice converts the specified device's network ports (including
// addresses and vlans) to the named network type, consistent with the Equinix
// Metal console experience.
//
// Deprecated: Equinix Metal devices may support more than two ports and the
// whole-device single word network type can no longer capture the capabilities
// and permutations of device port configurations.
func (i *DevicePortServiceOp) ConvertDevice(d *Device, targetType string) error {
	bondPorts := d.GetBondPorts()

	if targetType == NetworkTypeL3 {
		// TODO: remove vlans from all the ports
		for _, p := range bondPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
		_, _, err := i.PortToLayerThree(d.ID, "bond0")
		if err != nil {
			return err
		}
		allEthPorts, err := i.GetAllEthPorts(d)
		if err != nil {
			return err
		}
		for _, p := range allEthPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeHybrid {
		for _, p := range bondPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}

		_, _, err := i.PortToLayerThree(d.ID, "bond0")
		if err != nil {
			return err
		}

		// ports need to be refreshed before bonding/disbonding
		oddEthPorts, err := i.GetOddEthPorts(d)
		if err != nil {
			return err
		}

		for _, p := range oddEthPorts {
			_, _, err := i.Disbond(p, false)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeL2Individual {
		_, _, err := i.PortToLayerTwo(d.ID, "bond0")
		if err != nil {
			return err
		}
		for _, p := range bondPorts {
			_, _, err = i.Disbond(p, true)
			if err != nil {
				return err
			}
		}
	}
	if targetType == NetworkTypeL2Bonded {

		for _, p := range bondPorts {
			_, _, err := i.PortToLayerTwo(d.ID, p.Name)
			if err != nil {
				return err
			}
		}
		allEthPorts, err := i.GetAllEthPorts(d)
		if err != nil {
			return err
		}
		for _, p := range allEthPorts {
			_, _, err := i.Bond(p, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeviceToNetworkType fetches the specified device and converts its network
// ports (including addresses and vlans) to the named network type, consistent
// with the Equinix Metal console experience.
//
// Deprecated: use DevicePortServiceOp.ConvertDevice which this function thinly
// wraps.
func (i *DevicePortServiceOp) DeviceToNetworkType(deviceID string, targetType string) (*Device, error) {
	if validateErr := ValidateUUID(deviceID); validateErr != nil {
		return nil, validateErr
	}
	d, _, err := i.client.Devices.Get(deviceID, nil)
	if err != nil {
		return nil, err
	}

	curType := d.GetNetworkType()

	if curType == targetType {
		return nil, fmt.Errorf("Device already is in state %s", targetType)
	}
	err = i.ConvertDevice(d, targetType)
	if err != nil {
		return nil, err
	}

	d, _, err = i.client.Devices.Get(deviceID, nil)

	if err != nil {
		return nil, err
	}

	finalType := d.GetNetworkType()

	if finalType != targetType {
		return nil, fmt.Errorf(
			"Failed to convert device %s from %s to %s. New type was %s",
			deviceID, curType, targetType, finalType)

	}
	return d, err
}
