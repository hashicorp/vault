package gocb

// ScopeEventingFunctionManager provides methods for performing scoped eventing function management operations.
// This manager is designed to work only against Couchbase Server 7.1+.
//
// # UNCOMMITTED
//
// This API is UNCOMMITTED and may change in the future.
type ScopeEventingFunctionManager struct {
	controller *providerController[eventingManagementProvider]

	scope *Scope
}

// UpsertFunction inserts or updates an eventing function.
func (efm *ScopeEventingFunctionManager) UpsertFunction(function EventingFunction, opts *UpsertEventingFunctionOptions) error {
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &UpsertEventingFunctionOptions{}
		}
		
		return provider.UpsertFunction(efm.scope, function, opts)
	})
}

// DropFunction drops an eventing function.
func (efm *ScopeEventingFunctionManager) DropFunction(name string, opts *DropEventingFunctionOptions) error {
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &DropEventingFunctionOptions{}
		}

		return provider.DropFunction(efm.scope, name, opts)
	})
}

// DeployFunction deploys an eventing function.
func (efm *ScopeEventingFunctionManager) DeployFunction(name string, opts *DeployEventingFunctionOptions) error {
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &DeployEventingFunctionOptions{}
		}

		return provider.DeployFunction(efm.scope, name, opts)
	})
}

// UndeployFunction undeploys an eventing function.
func (efm *ScopeEventingFunctionManager) UndeployFunction(name string, opts *UndeployEventingFunctionOptions) error {
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &UndeployEventingFunctionOptions{}
		}

		return provider.UndeployFunction(efm.scope, name, opts)
	})
}

// GetAllFunctions fetches all the eventing functions that are in this scope.
func (efm *ScopeEventingFunctionManager) GetAllFunctions(opts *GetAllEventingFunctionsOptions) ([]EventingFunction, error) {
	return autoOpControl(efm.controller, func(provider eventingManagementProvider) ([]EventingFunction, error) {
		if opts == nil {
			opts = &GetAllEventingFunctionsOptions{}
		}

		return provider.GetAllFunctions(efm.scope, opts)
	})
}

// GetFunction fetches an eventing function.
func (efm *ScopeEventingFunctionManager) GetFunction(name string, opts *GetEventingFunctionOptions) (*EventingFunction, error) {
	return autoOpControl(efm.controller, func(provider eventingManagementProvider) (*EventingFunction, error) {
		if opts == nil {
			opts = &GetEventingFunctionOptions{}
		}

		return provider.GetFunction(efm.scope, name, opts)
	})
}

// PauseFunction pauses an eventing function.
func (efm *ScopeEventingFunctionManager) PauseFunction(name string, opts *PauseEventingFunctionOptions) error {
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &PauseEventingFunctionOptions{}
		}

		return provider.PauseFunction(efm.scope, name, opts)
	})
}

// ResumeFunction resumes an eventing function.
func (efm *ScopeEventingFunctionManager) ResumeFunction(name string, opts *ResumeEventingFunctionOptions) error {
	return autoOpControlErrorOnly(efm.controller, func(provider eventingManagementProvider) error {
		if opts == nil {
			opts = &ResumeEventingFunctionOptions{}
		}

		return provider.ResumeFunction(efm.scope, name, opts)
	})
}

// FunctionsStatus fetches the current status of all eventing functions.
func (efm *ScopeEventingFunctionManager) FunctionsStatus(opts *EventingFunctionsStatusOptions) (*EventingStatus, error) {
	return autoOpControl(efm.controller, func(provider eventingManagementProvider) (*EventingStatus, error) {
		if opts == nil {
			opts = &EventingFunctionsStatusOptions{}
		}

		return provider.FunctionsStatus(efm.scope, opts)
	})
}
