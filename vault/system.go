package vault

// SystemBackend implements the LogicalBackend interface but is used
// to interact with the core of the system. It acts like a "procfs"
// to provide a uniform interface to vault.
type SystemBackend struct {
	core *Core
}

func (s *SystemBackend) HandleRequest(*Request) (*Response, error) {
	return nil, nil
}

func (s *SystemBackend) RootPaths() []string {
	return []string{
		"acls*",    // Restrict all access to ACLs
		"auth/*",   // Restrict modifications to ACLs
		"mounts/*", // Restrict modifications to mounts
		"remount",  // Restrict modifications to mounts
		"seal",     // Restrict re-sealing to root
	}
}
