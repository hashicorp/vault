package linodego

import (
	"context"
	"encoding/json"
	"time"

	"github.com/linode/linodego/internal/parseabletime"
)

// LinodeKernel represents a Linode Instance kernel object
type LinodeKernel struct {
	ID           string     `json:"id"`
	Label        string     `json:"label"`
	Version      string     `json:"version"`
	Architecture string     `json:"architecture"`
	Deprecated   bool       `json:"deprecated"`
	KVM          bool       `json:"kvm"`
	XEN          bool       `json:"xen"`
	PVOPS        bool       `json:"pvops"`
	Built        *time.Time `json:"-"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (i *LinodeKernel) UnmarshalJSON(b []byte) error {
	type Mask LinodeKernel

	p := struct {
		*Mask
		Built *parseabletime.ParseableTime `json:"built"`
	}{
		Mask: (*Mask)(i),
	}

	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}

	i.Built = (*time.Time)(p.Built)

	return nil
}

// ListKernels lists linode kernels. This endpoint is cached by default.
func (c *Client) ListKernels(ctx context.Context, opts *ListOptions) ([]LinodeKernel, error) {
	endpoint, err := generateListCacheURL("linode/kernels", opts)
	if err != nil {
		return nil, err
	}

	if result := c.getCachedResponse(endpoint); result != nil {
		return result.([]LinodeKernel), nil
	}

	response, err := getPaginatedResults[LinodeKernel](ctx, c, "linode/kernels", opts)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(endpoint, response, nil)

	return response, nil
}

// GetKernel gets the kernel with the provided ID. This endpoint is cached by default.
func (c *Client) GetKernel(ctx context.Context, kernelID string) (*LinodeKernel, error) {
	e := formatAPIPath("linode/kernels/%s", kernelID)

	if result := c.getCachedResponse(e); result != nil {
		result := result.(LinodeKernel)
		return &result, nil
	}

	response, err := doGETRequest[LinodeKernel](ctx, c, e)
	if err != nil {
		return nil, err
	}

	c.addCachedResponse(e, response, nil)

	return response, nil
}
