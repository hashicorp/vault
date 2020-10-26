package mongodbatlas

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const containersPath = "groups/%s/containers"

//ContainersService is an interface for interfacing with the Network Peering Containers
// endpoints of the MongoDB Atlas API.
//See more: https://docs.atlas.mongodb.com/reference/api/vpc/
type ContainersService interface {
	List(context.Context, string, *ContainersListOptions) ([]Container, *Response, error)
	Get(context.Context, string, string) (*Container, *Response, error)
	Create(context.Context, string, *Container) (*Container, *Response, error)
	Update(context.Context, string, string, *Container) (*Container, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

//ContainersServiceOp handles communication with the Network Peering Container related methods
// of the MongoDB Atlas API
type ContainersServiceOp struct {
	client *Client
}

var _ ContainersService = &ContainersServiceOp{}

type ContainersListOptions struct {
	ProviderName string `url:"providerName,omitempty"`
	ListOptions
}

// Container represents MongoDB network peering containter.
type Container struct {
	AtlasCIDRBlock      string `json:"atlasCidrBlock,omitempty"`
	AzureSubscriptionID string `json:"azureSubscriptionId,omitempty"`
	GCPProjectID        string `json:"gcpProjectId,omitempty"`
	ID                  string `json:"id,omitempty"`
	NetworkName         string `json:"networkName,omitempty"`
	ProviderName        string `json:"providerName,omitempty"`
	Provisioned         *bool  `json:"provisioned,omitempty"`
	Region              string `json:"region,omitempty"`
	RegionName          string `json:"regionName,omitempty"`
	VNetName            string `json:"vnetName,omitempty"`
	VPCID               string `json:"vpcId,omitempty"`
}

// containersResponse is the response from the ContainersService.List.
type containersResponse struct {
	Links      []*Link     `json:"links,omitempty"`
	Results    []Container `json:"results,omitempty"`
	TotalCount int         `json:"totalCount,omitempty"`
}

//List all containers in the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-containers-list/
func (s *ContainersServiceOp) List(ctx context.Context, groupID string, listOptions *ContainersListOptions) ([]Container, *Response, error) {
	path := fmt.Sprintf(containersPath, groupID)

	//Add query params from listOptions
	path, err := setListOptions(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(containersResponse)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

//Get gets the network peering container specified to {CONTAINER-ID} from the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/vpc-get-container/
func (s *ContainersServiceOp) Get(ctx context.Context, groupID string, containerID string) (*Container, *Response, error) {
	if containerID == "" {
		return nil, nil, NewArgError("perrID", "must be set")
	}

	basePath := fmt.Sprintf(containersPath, groupID)
	escapedEntry := url.PathEscape(containerID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Container)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Add a network peering container to the project associated to {GROUP-ID}.
//See more: https://docs.atlas.mongodb.com/reference/api/vpc-create-container/
func (s *ContainersServiceOp) Create(ctx context.Context, groupID string, createRequest *Container) (*Container, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(containersPath, groupID)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Container)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Update a network peering container in the project associated to {GROUP-ID}
//See more: https://docs.atlas.mongodb.com/reference/api/vpc-update-container/
func (s *ContainersServiceOp) Update(ctx context.Context, groupID string, containerID string, updateRequest *Container) (*Container, *Response, error) {
	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(containersPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, containerID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Container)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

//Delete the network peering container specified to {CONTAINER-ID} from the project associated to {GROUP-ID}.
// See more: https://docs.atlas.mongodb.com/reference/api/vpc-delete-one-container/
func (s *ContainersServiceOp) Delete(ctx context.Context, groupID string, containerID string) (*Response, error) {
	if containerID == "" {
		return nil, NewArgError("containerID", "must be set")
	}

	basePath := fmt.Sprintf(containersPath, groupID)
	escapedEntry := url.PathEscape(containerID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	//To avoid API Issues
	req.Header.Del("Content-Type")

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
