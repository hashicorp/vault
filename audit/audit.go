package audit

// Backend interface must be implemented for an audit
// mechanism to be made available. Audit backends can be enabled to
// sink information to different backends such as logs, file, databases,
// or other external services.
type Backend interface {
}

// Factory is the factory function to create an audit backend.
type Factory func(map[string]string) (Backend, error)
