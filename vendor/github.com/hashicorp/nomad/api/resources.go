package api

import "github.com/hashicorp/nomad/helper"

// Resources encapsulates the required resources of
// a given task or task group.
type Resources struct {
	CPU      *int
	MemoryMB *int `mapstructure:"memory"`
	DiskMB   *int `mapstructure:"disk"`
	IOPS     *int
	Networks []*NetworkResource
}

// Canonicalize will supply missing values in the cases
// where they are not provided.
func (r *Resources) Canonicalize() {
	defaultResources := DefaultResources()
	if r.CPU == nil {
		r.CPU = defaultResources.CPU
	}
	if r.MemoryMB == nil {
		r.MemoryMB = defaultResources.MemoryMB
	}
	if r.IOPS == nil {
		r.IOPS = defaultResources.IOPS
	}
	for _, n := range r.Networks {
		n.Canonicalize()
	}
}

// DefaultResources is a small resources object that contains the
// default resources requests that we will provide to an object.
// ---  THIS FUNCTION IS REPLICATED IN nomad/structs/structs.go
// and should be kept in sync.
func DefaultResources() *Resources {
	return &Resources{
		CPU:      helper.IntToPtr(100),
		MemoryMB: helper.IntToPtr(300),
		IOPS:     helper.IntToPtr(0),
	}
}

// MinResources is a small resources object that contains the
// absolute minimum resources that we will provide to an object.
// This should not be confused with the defaults which are
// provided in DefaultResources() ---  THIS LOGIC IS REPLICATED
// IN nomad/structs/structs.go and should be kept in sync.
func MinResources() *Resources {
	return &Resources{
		CPU:      helper.IntToPtr(20),
		MemoryMB: helper.IntToPtr(10),
		IOPS:     helper.IntToPtr(0),
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
	if other.IOPS != nil {
		r.IOPS = other.IOPS
	}
	if len(other.Networks) != 0 {
		r.Networks = other.Networks
	}
}

type Port struct {
	Label string
	Value int `mapstructure:"static"`
}

// NetworkResource is used to describe required network
// resources of a given task.
type NetworkResource struct {
	Device        string
	CIDR          string
	IP            string
	MBits         *int
	ReservedPorts []Port
	DynamicPorts  []Port
}

func (n *NetworkResource) Canonicalize() {
	if n.MBits == nil {
		n.MBits = helper.IntToPtr(10)
	}
}
