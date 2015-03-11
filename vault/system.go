package vault

// SystemBackend implements the LogicalBackend interface but is used
// to interact with the core of the system. It acts like a "procfs"
// to provide a uniform interface to vault.
type SystemBackend struct {
	core *Core
}

func (s *SystemBackend) HandleRequest(req *Request) (*Response, error) {
	// Switch on the path to route to the appropriate handler
	switch {
	case req.Path == "mounts":
		return s.handleMountTable(req)
	default:
		return nil, ErrUnsupportedPath
	}
}

func (s *SystemBackend) RootPaths() []string {
	return []string{}
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (s *SystemBackend) handleMountTable(req *Request) (*Response, error) {
	switch req.Operation {
	case ReadOperation:
	case HelpOperation:
		return HelpResponse("logical backend mount table", nil), nil
	default:
		return nil, ErrUnsupportedOperation
	}

	s.core.mountsLock.RLock()
	defer s.core.mountsLock.RUnlock()

	resp := &Response{
		IsSecret: false,
		Data:     make(map[string]interface{}),
	}
	for _, entry := range s.core.mounts.Entries {
		info := map[string]string{
			"type":        entry.Type,
			"description": entry.Description,
		}
		resp.Data[entry.Path] = info
	}
	return resp, nil
}
