package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
)

const privateEndpointsPath = "groups/%s/privateEndpoint"

// PrivateEndpointsService is an interface for interfacing with the Private Endpoints
// of the MongoDB Atlas API.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint/
type PrivateEndpointsService interface {
	Create(context.Context, string, *PrivateEndpointConnection) (*PrivateEndpointConnection, *Response, error)
	Get(context.Context, string, string) (*PrivateEndpointConnection, *Response, error)
	List(context.Context, string, *ListOptions) ([]PrivateEndpointConnection, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
	AddOneInterfaceEndpoint(context.Context, string, string, string) (*InterfaceEndpointConnection, *Response, error)
	GetOneInterfaceEndpoint(context.Context, string, string, string) (*InterfaceEndpointConnection, *Response, error)
	DeleteOneInterfaceEndpoint(context.Context, string, string, string) (*Response, error)
}

// PrivateEndpointsServiceOp handles communication with the PrivateEndpoints related methods
// of the MongoDB Atlas API
type PrivateEndpointsServiceOp service

var _ PrivateEndpointsService = &PrivateEndpointsServiceOp{}

// PrivateEndpointConnection represents MongoDB Private Endpoint Connection.
type PrivateEndpointConnection struct {
	ID                  string   `json:"id,omitempty"`                  // Unique identifier of the AWS PrivateLink connection.
	ProviderName        string   `json:"providerName,omitempty"`        // Name of the cloud provider you want to create the private endpoint connection for. Must be AWS.
	Region              string   `json:"region,omitempty"`              // Cloud provider region in which you want to create the private endpoint connection.
	EndpointServiceName string   `json:"endpointServiceName,omitempty"` // Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
	ErrorMessage        string   `json:"errorMessage,omitempty"`        // Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
	InterfaceEndpoints  []string `json:"interfaceEndpoints,omitempty"`  // Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
	Status              string   `json:"status,omitempty"`              // Status of the AWS PrivateLink connection: INITIATING, WAITING_FOR_USER, FAILED, DELETING.
}

// InterfaceEndpointConnection represents MongoDB Interface Endpoint Connection.
type InterfaceEndpointConnection struct {
	ID               string `json:"interfaceEndpointId,omitempty"` // Unique identifier of the interface endpoint.
	DeleteRequested  *bool  `json:"deleteRequested,omitempty"`     // Indicates if Atlas received a request to remove the interface endpoint from the private endpoint connection.
	ErrorMessage     string `json:"errorMessage,omitempty"`        // Error message pertaining to the interface endpoint. Returns null if there are no errors.
	ConnectionStatus string `json:"connectionStatus,omitempty"`    // Status of the interface endpoint: NONE, PENDING_ACCEPTANCE, PENDING, AVAILABLE, REJECTED, DELETING.
}

// Create one private endpoint connection in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-create-one-private-endpoint-connection/
func (s *PrivateEndpointsServiceOp) Create(ctx context.Context, groupID string, createRequest *PrivateEndpointConnection) (*PrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(privateEndpointsPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Get retrieves details for one private endpoint connection by ID in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-get-one-private-endpoint-connection/
func (s *PrivateEndpointsServiceOp) Get(ctx context.Context, groupID, privateLinkID string) (*PrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if privateLinkID == "" {
		return nil, nil, NewArgError("privateLinkID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, privateLinkID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(PrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// List retrieves details for all private endpoint connections in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-get-all-private-endpoint-connections/
func (s *PrivateEndpointsServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]PrivateEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(privateEndpointsPath, groupID)

	// Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new([]PrivateEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return *root, resp, nil
}

// Delete removes one private endpoint connection in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-delete-one-private-endpoint-connection/
func (s *PrivateEndpointsServiceOp) Delete(ctx context.Context, groupID, privateLinkID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if privateLinkID == "" {
		return nil, NewArgError("privateLinkID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, privateLinkID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}

// AddOneInterfaceEndpoint adds one interface endpoint to a private endpoint connection in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-create-one-interface-endpoint/
func (s *PrivateEndpointsServiceOp) AddOneInterfaceEndpoint(ctx context.Context, groupID, privateLinkID, interfaceEndpointID string) (*InterfaceEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if privateLinkID == "" {
		return nil, nil, NewArgError("privateLinkID", "must be set")
	}
	if interfaceEndpointID == "" {
		return nil, nil, NewArgError("interfaceEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/interfaceEndpoints", basePath, privateLinkID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, &InterfaceEndpointConnection{ID: interfaceEndpointID})
	if err != nil {
		return nil, nil, err
	}

	root := new(InterfaceEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetOneInterfaceEndpoint retrieves one interface endpoint in a private endpoint connection in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-get-one-interface-endpoint/
func (s *PrivateEndpointsServiceOp) GetOneInterfaceEndpoint(ctx context.Context, groupID, privateLinkID, interfaceEndpointID string) (*InterfaceEndpointConnection, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if privateLinkID == "" {
		return nil, nil, NewArgError("privateLinkID", "must be set")
	}
	if interfaceEndpointID == "" {
		return nil, nil, NewArgError("interfaceEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/interfaceEndpoints/%s", basePath, privateLinkID, interfaceEndpointID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(InterfaceEndpointConnection)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// DeleteOneInterfaceEndpoint removes one interface endpoint from a private endpoint connection in an Atlas project.
// See more: https://docs.atlas.mongodb.com/reference/api/private-endpoint-delete-one-interface-endpoint/
func (s *PrivateEndpointsServiceOp) DeleteOneInterfaceEndpoint(ctx context.Context, groupID, privateLinkID, interfaceEndpointID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if privateLinkID == "" {
		return nil, NewArgError("privateLinkID", "must be set")
	}
	if interfaceEndpointID == "" {
		return nil, NewArgError("interfaceEndpointID", "must be set")
	}

	basePath := fmt.Sprintf(privateEndpointsPath, groupID)
	path := fmt.Sprintf("%s/%s/interfaceEndpoints/%s", basePath, privateLinkID, interfaceEndpointID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}
