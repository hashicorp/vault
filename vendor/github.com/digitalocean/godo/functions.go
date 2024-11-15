package godo

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	functionsBasePath        = "/v2/functions/namespaces"
	functionsNamespacePath   = functionsBasePath + "/%s"
	functionsTriggerBasePath = functionsNamespacePath + "/triggers"
)

type FunctionsService interface {
	ListNamespaces(context.Context) ([]FunctionsNamespace, *Response, error)
	GetNamespace(context.Context, string) (*FunctionsNamespace, *Response, error)
	CreateNamespace(context.Context, *FunctionsNamespaceCreateRequest) (*FunctionsNamespace, *Response, error)
	DeleteNamespace(context.Context, string) (*Response, error)

	ListTriggers(context.Context, string) ([]FunctionsTrigger, *Response, error)
	GetTrigger(context.Context, string, string) (*FunctionsTrigger, *Response, error)
	CreateTrigger(context.Context, string, *FunctionsTriggerCreateRequest) (*FunctionsTrigger, *Response, error)
	UpdateTrigger(context.Context, string, string, *FunctionsTriggerUpdateRequest) (*FunctionsTrigger, *Response, error)
	DeleteTrigger(context.Context, string, string) (*Response, error)
}

type FunctionsServiceOp struct {
	client *Client
}

var _ FunctionsService = &FunctionsServiceOp{}

type namespacesRoot struct {
	Namespaces []FunctionsNamespace `json:"namespaces,omitempty"`
}

type namespaceRoot struct {
	Namespace *FunctionsNamespace `json:"namespace,omitempty"`
}

type FunctionsNamespace struct {
	ApiHost   string    `json:"api_host,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Label     string    `json:"label,omitempty"`
	Region    string    `json:"region,omitempty"`
	UUID      string    `json:"uuid,omitempty"`
	Key       string    `json:"key,omitempty"`
}

type FunctionsNamespaceCreateRequest struct {
	Label  string `json:"label"`
	Region string `json:"region"`
}

type triggersRoot struct {
	Triggers []FunctionsTrigger `json:"triggers,omitempty"`
}

type triggerRoot struct {
	Trigger *FunctionsTrigger `json:"trigger,omitempty"`
}

type FunctionsTrigger struct {
	Namespace        string                   `json:"namespace,omitempty"`
	Function         string                   `json:"function,omitempty"`
	Type             string                   `json:"type,omitempty"`
	Name             string                   `json:"name,omitempty"`
	IsEnabled        bool                     `json:"is_enabled"`
	CreatedAt        time.Time                `json:"created_at,omitempty"`
	UpdatedAt        time.Time                `json:"updated_at,omitempty"`
	ScheduledDetails *TriggerScheduledDetails `json:"scheduled_details,omitempty"`
	ScheduledRuns    *TriggerScheduledRuns    `json:"scheduled_runs,omitempty"`
}

type TriggerScheduledDetails struct {
	Cron string                 `json:"cron,omitempty"`
	Body map[string]interface{} `json:"body,omitempty"`
}

type TriggerScheduledRuns struct {
	LastRunAt time.Time `json:"last_run_at,omitempty"`
	NextRunAt time.Time `json:"next_run_at,omitempty"`
}

type FunctionsTriggerCreateRequest struct {
	Name             string                   `json:"name"`
	Type             string                   `json:"type"`
	Function         string                   `json:"function"`
	IsEnabled        bool                     `json:"is_enabled"`
	ScheduledDetails *TriggerScheduledDetails `json:"scheduled_details,omitempty"`
}

type FunctionsTriggerUpdateRequest struct {
	IsEnabled        *bool                    `json:"is_enabled,omitempty"`
	ScheduledDetails *TriggerScheduledDetails `json:"scheduled_details,omitempty"`
}

// Gets a list of namespaces
func (s *FunctionsServiceOp) ListNamespaces(ctx context.Context) ([]FunctionsNamespace, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, functionsBasePath, nil)
	if err != nil {
		return nil, nil, err
	}
	nsRoot := new(namespacesRoot)
	resp, err := s.client.Do(ctx, req, nsRoot)
	if err != nil {
		return nil, resp, err
	}
	return nsRoot.Namespaces, resp, nil
}

// Gets a single namespace
func (s *FunctionsServiceOp) GetNamespace(ctx context.Context, namespace string) (*FunctionsNamespace, *Response, error) {
	path := fmt.Sprintf(functionsNamespacePath, namespace)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	nsRoot := new(namespaceRoot)
	resp, err := s.client.Do(ctx, req, nsRoot)
	if err != nil {
		return nil, resp, err
	}
	return nsRoot.Namespace, resp, nil
}

// Creates a namespace
func (s *FunctionsServiceOp) CreateNamespace(ctx context.Context, opts *FunctionsNamespaceCreateRequest) (*FunctionsNamespace, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, functionsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}
	nsRoot := new(namespaceRoot)
	resp, err := s.client.Do(ctx, req, nsRoot)
	if err != nil {
		return nil, resp, err
	}
	return nsRoot.Namespace, resp, nil
}

// Delete a namespace
func (s *FunctionsServiceOp) DeleteNamespace(ctx context.Context, namespace string) (*Response, error) {
	path := fmt.Sprintf(functionsNamespacePath, namespace)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListTriggers gets a list of triggers
func (s *FunctionsServiceOp) ListTriggers(ctx context.Context, namespace string) ([]FunctionsTrigger, *Response, error) {
	path := fmt.Sprintf(functionsTriggerBasePath, namespace)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(triggersRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Triggers, resp, nil
}

// GetTrigger gets a single trigger
func (s *FunctionsServiceOp) GetTrigger(ctx context.Context, namespace string, trigger string) (*FunctionsTrigger, *Response, error) {
	path := fmt.Sprintf(functionsTriggerBasePath+"/%s", namespace, trigger)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(triggerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Trigger, resp, nil
}

// CreateTrigger creates a trigger
func (s *FunctionsServiceOp) CreateTrigger(ctx context.Context, namespace string, opts *FunctionsTriggerCreateRequest) (*FunctionsTrigger, *Response, error) {
	path := fmt.Sprintf(functionsTriggerBasePath, namespace)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, opts)
	if err != nil {
		return nil, nil, err
	}
	root := new(triggerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Trigger, resp, nil
}

// UpdateTrigger updates a trigger
func (s *FunctionsServiceOp) UpdateTrigger(ctx context.Context, namespace string, trigger string, opts *FunctionsTriggerUpdateRequest) (*FunctionsTrigger, *Response, error) {
	path := fmt.Sprintf(functionsTriggerBasePath+"/%s", namespace, trigger)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, opts)

	if err != nil {
		return nil, nil, err
	}
	root := new(triggerRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Trigger, resp, nil
}

// DeleteTrigger deletes a trigger
func (s *FunctionsServiceOp) DeleteTrigger(ctx context.Context, namespace string, trigger string) (*Response, error) {
	path := fmt.Sprintf(functionsTriggerBasePath+"/%s", namespace, trigger)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)

	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
