// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"slices"
	"strconv"
)

// Resources encapsulates the required resources of
// a given task or task group.
type Resources struct {
	CPU         *int               `hcl:"cpu,optional"`
	Cores       *int               `hcl:"cores,optional"`
	MemoryMB    *int               `mapstructure:"memory" hcl:"memory,optional"`
	MemoryMaxMB *int               `mapstructure:"memory_max" hcl:"memory_max,optional"`
	DiskMB      *int               `mapstructure:"disk" hcl:"disk,optional"`
	Networks    []*NetworkResource `hcl:"network,block"`
	Devices     []*RequestedDevice `hcl:"device,block"`
	NUMA        *NUMAResource      `hcl:"numa,block"`
	SecretsMB   *int               `mapstructure:"secrets" hcl:"secrets,optional"`

	// COMPAT(0.10)
	// XXX Deprecated. Please do not use. The field will be removed in Nomad
	// 0.10 and is only being kept to allow any references to be removed before
	// then.
	IOPS *int `hcl:"iops,optional"`
}

// Canonicalize will supply missing values in the cases
// where they are not provided.
func (r *Resources) Canonicalize() {
	defaultResources := DefaultResources()
	if r.Cores == nil {
		r.Cores = defaultResources.Cores

		// only set cpu to the default value if it and cores is not defined
		if r.CPU == nil {
			r.CPU = defaultResources.CPU
		}
	}

	// CPU will be set to the default if cores is nil above.
	// If cpu is nil here then cores has been set and cpu should be 0
	if r.CPU == nil {
		r.CPU = pointerOf(0)
	}

	if r.MemoryMB == nil {
		r.MemoryMB = defaultResources.MemoryMB
	}
	for _, d := range r.Devices {
		d.Canonicalize()
	}

	r.NUMA.Canonicalize()
}

// DefaultResources is a small resources object that contains the
// default resources requests that we will provide to an object.
// ---  THIS FUNCTION IS REPLICATED IN nomad/structs/structs.go
// and should be kept in sync.
func DefaultResources() *Resources {
	return &Resources{
		CPU:      pointerOf(100),
		Cores:    pointerOf(0),
		MemoryMB: pointerOf(300),
	}
}

// MinResources is a small resources object that contains the
// absolute minimum resources that we will provide to an object.
// This should not be confused with the defaults which are
// provided in DefaultResources() ---  THIS LOGIC IS REPLICATED
// IN nomad/structs/structs.go and should be kept in sync.
func MinResources() *Resources {
	return &Resources{
		CPU:      pointerOf(1),
		Cores:    pointerOf(0),
		MemoryMB: pointerOf(10),
	}
}

// Merge merges this resource with another resource.
func (r *Resources) Merge(other *Resources) {
	if other == nil {
		return
	}
	if other.CPU != nil {
		r.CPU = other.CPU
	}
	if other.MemoryMB != nil {
		r.MemoryMB = other.MemoryMB
	}
	if other.DiskMB != nil {
		r.DiskMB = other.DiskMB
	}
	if len(other.Networks) != 0 {
		r.Networks = other.Networks
	}
	if len(other.Devices) != 0 {
		r.Devices = other.Devices
	}
	if other.NUMA != nil {
		r.NUMA = other.NUMA.Copy()
	}
	if other.SecretsMB != nil {
		r.SecretsMB = other.SecretsMB
	}
}

// NUMAResource contains the NUMA affinity request for scheduling purposes.
//
// Applies only to Nomad Enterprise.
type NUMAResource struct {
	// Affinity must be one of "none", "prefer", "require".
	Affinity string `hcl:"affinity,optional"`

	// Devices is the subset of devices requested by the task that must share
	// the same numa node, along with the tasks reserved cpu cores.
	Devices []string `hcl:"devices,optional"`
}

func (n *NUMAResource) Copy() *NUMAResource {
	if n == nil {
		return nil
	}
	return &NUMAResource{
		Affinity: n.Affinity,
		Devices:  slices.Clone(n.Devices),
	}
}

func (n *NUMAResource) Canonicalize() {
	if n == nil {
		return
	}
	if n.Affinity == "" {
		n.Affinity = "none"
	}
	if len(n.Devices) == 0 {
		n.Devices = nil
	}
}

type Port struct {
	Label           string `hcl:",label"`
	Value           int    `hcl:"static,optional"`
	To              int    `hcl:"to,optional"`
	HostNetwork     string `hcl:"host_network,optional"`
	IgnoreCollision bool   `hcl:"ignore_collision,optional"`
}

type DNSConfig struct {
	Servers  []string `mapstructure:"servers" hcl:"servers,optional"`
	Searches []string `mapstructure:"searches" hcl:"searches,optional"`
	Options  []string `mapstructure:"options" hcl:"options,optional"`
}
type CNIConfig struct {
	Args map[string]string `hcl:"args,optional"`
}

// NetworkResource is used to describe required network
// resources of a given task.
type NetworkResource struct {
	Mode          string     `hcl:"mode,optional"`
	Device        string     `hcl:"device,optional"`
	CIDR          string     `hcl:"cidr,optional"`
	IP            string     `hcl:"ip,optional"`
	DNS           *DNSConfig `hcl:"dns,block"`
	ReservedPorts []Port     `hcl:"reserved_ports,block"`
	DynamicPorts  []Port     `hcl:"port,block"`
	Hostname      string     `hcl:"hostname,optional"`

	// COMPAT(0.13)
	// XXX Deprecated. Please do not use. The field will be removed in Nomad
	// 0.13 and is only being kept to allow any references to be removed before
	// then.
	MBits *int       `hcl:"mbits,optional"`
	CNI   *CNIConfig `hcl:"cni,block"`
}

// COMPAT(0.13)
// XXX Deprecated. Please do not use. The method will be removed in Nomad
// 0.13 and is only being kept to allow any references to be removed before
// then.
func (n *NetworkResource) Megabits() int {
	if n == nil || n.MBits == nil {
		return 0
	}
	return *n.MBits
}

func (n *NetworkResource) Canonicalize() {
	// COMPAT(0.13)
	// Noop to maintain backwards compatibility
}

func (n *NetworkResource) HasPorts() bool {
	if n == nil {
		return false
	}

	return len(n.ReservedPorts)+len(n.DynamicPorts) > 0
}

// NodeDeviceResource captures a set of devices sharing a common
// vendor/type/device_name tuple.
type NodeDeviceResource struct {

	// Vendor specifies the vendor of device
	Vendor string

	// Type specifies the type of the device
	Type string

	// Name specifies the specific model of the device
	Name string

	// Instances are list of the devices matching the vendor/type/name
	Instances []*NodeDevice

	Attributes map[string]*Attribute
}

func (r NodeDeviceResource) ID() string {
	return r.Vendor + "/" + r.Type + "/" + r.Name
}

// NodeDevice is an instance of a particular device.
type NodeDevice struct {
	// ID is the ID of the device.
	ID string

	// Healthy captures whether the device is healthy.
	Healthy bool

	// HealthDescription is used to provide a human readable description of why
	// the device may be unhealthy.
	HealthDescription string

	// Locality stores HW locality information for the node to optionally be
	// used when making placement decisions.
	Locality *NodeDeviceLocality
}

// Attribute is used to describe the value of an attribute, optionally
// specifying units
type Attribute struct {
	// Float is the float value for the attribute
	FloatVal *float64 `json:"Float,omitempty"`

	// Int is the int value for the attribute
	IntVal *int64 `json:"Int,omitempty"`

	// String is the string value for the attribute
	StringVal *string `json:"String,omitempty"`

	// Bool is the bool value for the attribute
	BoolVal *bool `json:"Bool,omitempty"`

	// Unit is the optional unit for the set int or float value
	Unit string
}

func (a Attribute) String() string {
	switch {
	case a.FloatVal != nil:
		str := formatFloat(*a.FloatVal, 3)
		if a.Unit != "" {
			str += " " + a.Unit
		}
		return str
	case a.IntVal != nil:
		str := strconv.FormatInt(*a.IntVal, 10)
		if a.Unit != "" {
			str += " " + a.Unit
		}
		return str
	case a.StringVal != nil:
		return *a.StringVal
	case a.BoolVal != nil:
		return strconv.FormatBool(*a.BoolVal)
	default:
		return "<unknown>"
	}
}

// NodeDeviceLocality stores information about the devices hardware locality on
// the node.
type NodeDeviceLocality struct {
	// PciBusID is the PCI Bus ID for the device.
	PciBusID string
}

// RequestedDevice is used to request a device for a task.
type RequestedDevice struct {
	// Name is the request name. The possible values are as follows:
	// * <type>: A single value only specifies the type of request.
	// * <vendor>/<type>: A single slash delimiter assumes the vendor and type of device is specified.
	// * <vendor>/<type>/<name>: Two slash delimiters assume vendor, type and specific model are specified.
	//
	// Examples are as follows:
	// * "gpu"
	// * "nvidia/gpu"
	// * "nvidia/gpu/GTX2080Ti"
	Name string `hcl:",label"`

	// Count is the number of requested devices
	Count *uint64 `hcl:"count,optional"`

	// Constraints are a set of constraints to apply when selecting the device
	// to use.
	Constraints []*Constraint `hcl:"constraint,block"`

	// Affinities are a set of affinites to apply when selecting the device
	// to use.
	Affinities []*Affinity `hcl:"affinity,block"`
}

func (d *RequestedDevice) Canonicalize() {
	if d.Count == nil {
		d.Count = pointerOf(uint64(1))
	}

	for _, a := range d.Affinities {
		a.Canonicalize()
	}
}
