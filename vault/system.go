package vault

import "strings"

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
	case strings.HasPrefix(req.Path, "mount/"):
		return s.handleMountOperation(req)
	case req.Path == "remount":
		return s.handleRemount(req)
	default:
		return nil, ErrUnsupportedPath
	}
}

func (s *SystemBackend) RootPaths() []string {
	return []string{
		"mount/*",
		"remount",
	}
}

// handleMountTable handles the "mounts" endpoint to provide the mount table
func (s *SystemBackend) handleMountTable(req *Request) (*Response, error) {
	switch req.Operation {
	case ReadOperation:
	case HelpOperation:
		return HelpResponse("logical backend mount table", []string{"sys/mount/"}), nil
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

// handleMountOperation is used to mount or unmount a path
func (s *SystemBackend) handleMountOperation(req *Request) (*Response, error) {
	switch req.Operation {
	case WriteOperation:
		return s.handleMount(req)
	case DeleteOperation:
		return s.handleUnmount(req)
	case HelpOperation:
		return HelpResponse("used to mount or unmount a path", []string{"sys/mounts"}), nil
	default:
		return nil, ErrUnsupportedOperation
	}
}

// handleMount is used to mount a new path
func (s *SystemBackend) handleMount(req *Request) (*Response, error) {
	suffix := strings.TrimPrefix(req.Path, "mount/")
	if len(suffix) == 0 {
		return ErrorResponse("path cannot be blank"), ErrInvalidRequest
	}

	// Get the type and description (optionally)
	logicalType := req.GetString("type")
	if logicalType == "" {
		return ErrorResponse("backend type must be specified as a string"), ErrInvalidRequest
	}
	description := req.GetString("description")

	// Create the mount entry
	me := &MountEntry{
		Path:        suffix,
		Type:        logicalType,
		Description: description,
	}

	// Attempt mount
	if err := s.core.mountEntry(me); err != nil {
		return ErrorResponse(err.Error()), ErrInvalidRequest
	}
	return nil, nil
}

// handleUnmount is used to unmount a path
func (s *SystemBackend) handleUnmount(req *Request) (*Response, error) {
	suffix := strings.TrimPrefix(req.Path, "mount/")
	if len(suffix) == 0 {
		return ErrorResponse("path cannot be blank"), ErrInvalidRequest
	}

	// Attempt unmount
	if err := s.core.unmountPath(suffix); err != nil {
		return ErrorResponse(err.Error()), ErrInvalidRequest
	}
	return nil, nil
}

// handleRemount is used to remount a path
func (s *SystemBackend) handleRemount(req *Request) (*Response, error) {
	// Only accept write operations
	switch req.Operation {
	case WriteOperation:
	case HelpOperation:
		return HelpResponse("remount a backend path", []string{"sys/mount/", "sys/mounts"}), nil
	default:
		return nil, ErrUnsupportedOperation
	}

	// Get the paths
	fromPath := req.GetString("from")
	toPath := req.GetString("to")
	if fromPath == "" || toPath == "" {
		return ErrorResponse("both 'from' and 'to' path must be specified as a string"), ErrInvalidRequest
	}

	// Attempt remount
	if err := s.core.remountPath(fromPath, toPath); err != nil {
		return ErrorResponse(err.Error()), ErrInvalidRequest
	}
	return nil, nil
}
