package gocb

type eventingManagementProvider interface {
	UpsertFunction(scope *Scope, function EventingFunction, opts *UpsertEventingFunctionOptions) error
	DropFunction(scope *Scope, name string, opts *DropEventingFunctionOptions) error
	DeployFunction(scope *Scope, name string, opts *DeployEventingFunctionOptions) error
	UndeployFunction(scope *Scope, name string, opts *UndeployEventingFunctionOptions) error
	GetAllFunctions(scope *Scope, opts *GetAllEventingFunctionsOptions) ([]EventingFunction, error)
	GetFunction(scope *Scope, name string, opts *GetEventingFunctionOptions) (*EventingFunction, error)
	PauseFunction(scope *Scope, name string, opts *PauseEventingFunctionOptions) error
	ResumeFunction(scope *Scope, name string, opts *ResumeEventingFunctionOptions) error
	FunctionsStatus(scope *Scope, opts *EventingFunctionsStatusOptions) (*EventingStatus, error)
}
