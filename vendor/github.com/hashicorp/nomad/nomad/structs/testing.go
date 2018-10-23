package structs

// NodeResourcesToAllocatedResources converts a node resources to an allocated
// resources. The task name used is "web" and network is omitted. This is
// useful when trying to make an allocation fill an entire node.
func NodeResourcesToAllocatedResources(n *NodeResources) *AllocatedResources {
	if n == nil {
		return nil
	}

	return &AllocatedResources{
		Tasks: map[string]*AllocatedTaskResources{
			"web": {
				Cpu: AllocatedCpuResources{
					CpuShares: n.Cpu.CpuShares,
				},
				Memory: AllocatedMemoryResources{
					MemoryMB: n.Memory.MemoryMB,
				},
			},
		},
		Shared: AllocatedSharedResources{
			DiskMB: n.Disk.DiskMB,
		},
	}
}
