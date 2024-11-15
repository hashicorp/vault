package godo

import (
	"context"
	"fmt"
	"net/http"
)

// ReservedIPActionsService is an interface for interfacing with the
// reserved IPs actions endpoints of the Digital Ocean API.
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/Reserved-IP-Actions
type ReservedIPActionsService interface {
	Assign(ctx context.Context, ip string, dropletID int) (*Action, *Response, error)
	Unassign(ctx context.Context, ip string) (*Action, *Response, error)
	Get(ctx context.Context, ip string, actionID int) (*Action, *Response, error)
	List(ctx context.Context, ip string, opt *ListOptions) ([]Action, *Response, error)
}

// ReservedIPActionsServiceOp handles communication with the reserved IPs
// action related methods of the DigitalOcean API.
type ReservedIPActionsServiceOp struct {
	client *Client
}

// Assign a reserved IP to a droplet.
func (s *ReservedIPActionsServiceOp) Assign(ctx context.Context, ip string, dropletID int) (*Action, *Response, error) {
	request := &ActionRequest{
		"type":       "assign",
		"droplet_id": dropletID,
	}
	return s.doAction(ctx, ip, request)
}

// Unassign a rerserved IP from the droplet it is currently assigned to.
func (s *ReservedIPActionsServiceOp) Unassign(ctx context.Context, ip string) (*Action, *Response, error) {
	request := &ActionRequest{"type": "unassign"}
	return s.doAction(ctx, ip, request)
}

// Get an action for a particular reserved IP by id.
func (s *ReservedIPActionsServiceOp) Get(ctx context.Context, ip string, actionID int) (*Action, *Response, error) {
	path := fmt.Sprintf("%s/%d", reservedIPActionPath(ip), actionID)
	return s.get(ctx, path)
}

// List the actions for a particular reserved IP.
func (s *ReservedIPActionsServiceOp) List(ctx context.Context, ip string, opt *ListOptions) ([]Action, *Response, error) {
	path := reservedIPActionPath(ip)
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	return s.list(ctx, path)
}

func (s *ReservedIPActionsServiceOp) doAction(ctx context.Context, ip string, request *ActionRequest) (*Action, *Response, error) {
	path := reservedIPActionPath(ip)

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

func (s *ReservedIPActionsServiceOp) get(ctx context.Context, path string) (*Action, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Event, resp, err
}

func (s *ReservedIPActionsServiceOp) list(ctx context.Context, path string) ([]Action, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(actionsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Actions, resp, err
}

func reservedIPActionPath(ip string) string {
	return fmt.Sprintf("%s/%s/actions", reservedIPsBasePath, ip)
}
