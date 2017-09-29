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

func (r *Resources) Canonicalize() {
	if r.CPU == nil {
		r.CPU = helper.IntToPtr(100)
	}
	if r.MemoryMB == nil {
		r.MemoryMB = helper.IntToPtr(10)
	}
	if r.IOPS == nil {
		r.IOPS = helper.IntToPtr(0)
	}
	for _, n := range r.Networks {
		n.Canonicalize()
	}
}

func MinResources() *Resources {
	return &Resources{
		CPU:      helper.IntToPtr(100),
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
