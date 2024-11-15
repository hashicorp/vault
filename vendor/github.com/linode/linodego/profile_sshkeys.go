package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// SSHKey represents a SSHKey object
type SSHKey struct {
	ID      int        `json:"id"`
	Label   string     `json:"label"`
	SSHKey  string     `json:"ssh_key"`
	Created *time.Time `json:"-"`
}

// SSHKeyCreateOptions fields are those accepted by CreateSSHKey
type SSHKeyCreateOptions struct {
	Label  string `json:"label"`
	SSHKey string `json:"ssh_key"`
}

// SSHKeyUpdateOptions fields are those accepted by UpdateSSHKey
type SSHKeyUpdateOptions struct {
	Label string `json:"label"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *SSHKey) UnmarshalJSON(b []byte) error {
	type Mask SSHKey

	p := struct {
		*Mask
		Created *parseabletime.ParseableTime `json:"created"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Created = (*time.Time)(p.Created)

	return nil
}

// GetCreateOptions converts a SSHKey to SSHKeyCreateOptions for use in CreateSSHKey
func (i SSHKey) GetCreateOptions() (o SSHKeyCreateOptions) {
	o.Label = i.Label
	o.SSHKey = i.SSHKey
	return
}

// GetUpdateOptions converts a SSHKey to SSHKeyCreateOptions for use in UpdateSSHKey
func (i SSHKey) GetUpdateOptions() (o SSHKeyUpdateOptions) {
	o.Label = i.Label
	return
}

// ListSSHKeys lists SSHKeys
func (c *Client) ListSSHKeys(ctx context.Context, opts *ListOptions) ([]SSHKey, error) {
	response, err := getPaginatedResults[SSHKey](ctx, c, "profile/sshkeys", opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetSSHKey gets the sshkey with the provided ID
func (c *Client) GetSSHKey(ctx context.Context, keyID int) (*SSHKey, error) {
	e := formatAPIPath("profile/sshkeys/%d", keyID)
	response, err := doGETRequest[SSHKey](ctx, c, e)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// CreateSSHKey creates a SSHKey
func (c *Client) CreateSSHKey(ctx context.Context, opts SSHKeyCreateOptions) (*SSHKey, error) {
	e := "profile/sshkeys"
	response, err := doPOSTRequest[SSHKey](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateSSHKey updates the SSHKey with the specified id
func (c *Client) UpdateSSHKey(ctx context.Context, keyID int, opts SSHKeyUpdateOptions) (*SSHKey, error) {
	e := formatAPIPath("profile/sshkeys/%d", keyID)
	response, err := doPUTRequest[SSHKey](ctx, c, e, opts)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteSSHKey deletes the SSHKey with the specified id
func (c *Client) DeleteSSHKey(ctx context.Context, keyID int) error {
	e := formatAPIPath("profile/sshkeys/%d", keyID)
	err := doDELETERequest(ctx, c, e)
	return err
}
