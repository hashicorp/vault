package godo

import (
	"context"
	"fmt"
	"net/http"
)

const keysBasePath = "v2/account/keys"

// KeysService is an interface for interfacing with the SSH keys
// endpoints of the DigitalOcean API
// See: https://docs.digitalocean.com/reference/api/api-reference/#tag/SSH-Keys
type KeysService interface {
	List(context.Context, *ListOptions) ([]Key, *Response, error)
	GetByID(context.Context, int) (*Key, *Response, error)
	GetByFingerprint(context.Context, string) (*Key, *Response, error)
	Create(context.Context, *KeyCreateRequest) (*Key, *Response, error)
	UpdateByID(context.Context, int, *KeyUpdateRequest) (*Key, *Response, error)
	UpdateByFingerprint(context.Context, string, *KeyUpdateRequest) (*Key, *Response, error)
	DeleteByID(context.Context, int) (*Response, error)
	DeleteByFingerprint(context.Context, string) (*Response, error)
}

// KeysServiceOp handles communication with SSH key related method of the
// DigitalOcean API.
type KeysServiceOp struct {
	client *Client
}

var _ KeysService = &KeysServiceOp{}

// Key represents a DigitalOcean Key.
type Key struct {
	ID          int    `json:"id,float64,omitempty"`
	Name        string `json:"name,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`
}

// KeyUpdateRequest represents a request to update an SSH key stored in a DigitalOcean account.
type KeyUpdateRequest struct {
	Name string `json:"name"`
}

type keysRoot struct {
	SSHKeys []Key  `json:"ssh_keys"`
	Links   *Links `json:"links"`
	Meta    *Meta  `json:"meta"`
}

type keyRoot struct {
	SSHKey *Key `json:"ssh_key"`
}

func (s Key) String() string {
	return Stringify(s)
}

// KeyCreateRequest represents a request to create a new SSH key.
type KeyCreateRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// List all SSH keys
func (s *KeysServiceOp) List(ctx context.Context, opt *ListOptions) ([]Key, *Response, error) {
	path := keysBasePath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(keysRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}
	if m := root.Meta; m != nil {
		resp.Meta = m
	}

	return root.SSHKeys, resp, err
}

// Performs a get given a path
func (s *KeysServiceOp) get(ctx context.Context, path string) (*Key, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKey, resp, err
}

// GetByID gets an SSH key by its ID
func (s *KeysServiceOp) GetByID(ctx context.Context, keyID int) (*Key, *Response, error) {
	if keyID < 1 {
		return nil, nil, NewArgError("keyID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", keysBasePath, keyID)
	return s.get(ctx, path)
}

// GetByFingerprint gets an SSH key by its fingerprint
func (s *KeysServiceOp) GetByFingerprint(ctx context.Context, fingerprint string) (*Key, *Response, error) {
	if len(fingerprint) < 1 {
		return nil, nil, NewArgError("fingerprint", "cannot not be empty")
	}

	path := fmt.Sprintf("%s/%s", keysBasePath, fingerprint)
	return s.get(ctx, path)
}

// Create an SSH key using a KeyCreateRequest
func (s *KeysServiceOp) Create(ctx context.Context, createRequest *KeyCreateRequest) (*Key, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, keysBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKey, resp, err
}

// UpdateByID updates an SSH key name by ID.
func (s *KeysServiceOp) UpdateByID(ctx context.Context, keyID int, updateRequest *KeyUpdateRequest) (*Key, *Response, error) {
	if keyID < 1 {
		return nil, nil, NewArgError("keyID", "cannot be less than 1")
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d", keysBasePath, keyID)
	req, err := s.client.NewRequest(ctx, "PUT", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKey, resp, err
}

// UpdateByFingerprint updates an SSH key name by fingerprint.
func (s *KeysServiceOp) UpdateByFingerprint(ctx context.Context, fingerprint string, updateRequest *KeyUpdateRequest) (*Key, *Response, error) {
	if len(fingerprint) < 1 {
		return nil, nil, NewArgError("fingerprint", "cannot be empty")
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s", keysBasePath, fingerprint)
	req, err := s.client.NewRequest(ctx, "PUT", path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(keyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKey, resp, err
}

// Delete an SSH key using a path
func (s *KeysServiceOp) delete(ctx context.Context, path string) (*Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// DeleteByID deletes an SSH key by its id
func (s *KeysServiceOp) DeleteByID(ctx context.Context, keyID int) (*Response, error) {
	if keyID < 1 {
		return nil, NewArgError("keyID", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d", keysBasePath, keyID)
	return s.delete(ctx, path)
}

// DeleteByFingerprint deletes an SSH key by its fingerprint
func (s *KeysServiceOp) DeleteByFingerprint(ctx context.Context, fingerprint string) (*Response, error) {
	if len(fingerprint) < 1 {
		return nil, NewArgError("fingerprint", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", keysBasePath, fingerprint)
	return s.delete(ctx, path)
}
