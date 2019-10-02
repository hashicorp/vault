package template

import (
	"sync"
)

// TemplateServer
type TemplateServer struct {
	// config holds the template managers configuration
	// config *TaskTemplateManagerConfig

	// // lookup allows looking up the set of Nomad templates by their consul-template ID
	// lookup map[string][]*structs.Template

	// runner is the consul-template runner
	// runner *manager.Runner

	// 	// signals is a lookup map from the string representation of a signal to its
	// 	// actual signal
	// 	signals map[string]os.Signal

	// shutdownCh is used to signal and started goroutine to shutdown
	shutdownCh chan struct{}

	// shutdown marks whether the manager has been shutdown
	shutdown     bool
	shutdownLock sync.Mutex
}
