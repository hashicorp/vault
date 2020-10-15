package packngo

import "fmt"

const (
	sshKeyBasePath = "/ssh-keys"
)

// SSHKeyService interface defines available device methods
type SSHKeyService interface {
	List() ([]SSHKey, *Response, error)
	ProjectList(string) ([]SSHKey, *Response, error)
	Get(string) (*SSHKey, *Response, error)
	Create(*SSHKeyCreateRequest) (*SSHKey, *Response, error)
	Update(string, *SSHKeyUpdateRequest) (*SSHKey, *Response, error)
	Delete(string) (*Response, error)
}

type sshKeyRoot struct {
	SSHKeys []SSHKey `json:"ssh_keys"`
}

// SSHKey represents a user's ssh key
type SSHKey struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Key         string `json:"key"`
	FingerPrint string `json:"fingerprint"`
	Created     string `json:"created_at"`
	Updated     string `json:"updated_at"`
	User        User   `json:"user,omitempty"`
	URL         string `json:"href,omitempty"`
}

func (s SSHKey) String() string {
	return Stringify(s)
}

// SSHKeyCreateRequest type used to create an ssh key
type SSHKeyCreateRequest struct {
	Label     string `json:"label"`
	Key       string `json:"key"`
	ProjectID string `json:"-"`
}

func (s SSHKeyCreateRequest) String() string {
	return Stringify(s)
}

// SSHKeyUpdateRequest type used to update an ssh key
type SSHKeyUpdateRequest struct {
	Label *string `json:"label,omitempty"`
	Key   *string `json:"key,omitempty"`
}

func (s SSHKeyUpdateRequest) String() string {
	return Stringify(s)
}

// SSHKeyServiceOp implements SSHKeyService
type SSHKeyServiceOp struct {
	client *Client
}

func (s *SSHKeyServiceOp) list(url string) ([]SSHKey, *Response, error) {
	root := new(sshKeyRoot)

	resp, err := s.client.DoRequest("GET", url, nil, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKeys, resp, err
}

// ProjectList lists ssh keys of a project
func (s *SSHKeyServiceOp) ProjectList(projectID string) ([]SSHKey, *Response, error) {
	return s.list(fmt.Sprintf("%s/%s%s", projectBasePath, projectID, sshKeyBasePath))

}

// List returns a user's ssh keys
func (s *SSHKeyServiceOp) List() ([]SSHKey, *Response, error) {
	return s.list(sshKeyBasePath)
}

// Get returns an ssh key by id
func (s *SSHKeyServiceOp) Get(sshKeyID string) (*SSHKey, *Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, sshKeyID)
	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("GET", path, nil, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Create creates a new ssh key
func (s *SSHKeyServiceOp) Create(createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	path := sshKeyBasePath
	if createRequest.ProjectID != "" {
		path = fmt.Sprintf("%s/%s%s", projectBasePath, createRequest.ProjectID, sshKeyBasePath)
	}
	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("POST", path, createRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Update updates an ssh key
func (s *SSHKeyServiceOp) Update(id string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Response, error) {
	if updateRequest.Label == nil && updateRequest.Key == nil {
		return nil, nil, fmt.Errorf("You must set either Label or Key string for SSH Key update")
	}
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, id)

	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("PATCH", path, updateRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Delete deletes an ssh key
func (s *SSHKeyServiceOp) Delete(sshKeyID string) (*Response, error) {
	path := fmt.Sprintf("%s/%s", sshKeyBasePath, sshKeyID)

	return s.client.DoRequest("DELETE", path, nil, nil)
}
